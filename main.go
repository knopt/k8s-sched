package main

import (
	"bytes"
	"encoding/json"
	"github.com/knopt/k8s-sched-extender/cmn"
	"io"
	"log"
	"net/http"

	"github.com/golang/glog"
	"github.com/knopt/k8s-sched-extender/sched"
	"github.com/knopt/k8s-sched-extender/stats"

	v1 "k8s.io/api/core/v1"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
)

const (
	versionPath      = "/version"
	apiPrefix        = "/scheduler"
	bindPath         = apiPrefix + "/bind"
	filterPrefix     = apiPrefix + "/filter"
	prioritizePrefix = apiPrefix + "/prioritize"
	preemptionPrefix = apiPrefix + "/preemption"
)

var (
	filter = sched.Filter{
		Name: "always_true",
		Func: func(node v1.Node, pod *v1.Pod) bool {
			return true
		},
	}
	
	prioritize = sched.Prioritize{
		Name: "always_1",
		Func: func(pod v1.Pod, nodes []v1.Node) (*schedulerapi.HostPriorityList, error) {
			var hostPriorityList schedulerapi.HostPriorityList
			hostPriorityList = make([]schedulerapi.HostPriority, len(nodes))
			for i, node := range nodes {
				glog.Warningf("node %s cpu capacity %s, allocatable %s", node.Name, node.Status.Capacity.Cpu(), node.Status.Allocatable.Cpu())
				hostPriorityList[i] = schedulerapi.HostPriority{
					Host:  node.Name,
					Score: i,
				}
			}

			return &hostPriorityList, nil
		},
	}
)

func main() {
	glog.Error("starting main")
	stats.Registry = stats.NewStatsRunner()
	go stats.Registry.Run()

	http.HandleFunc(filterPrefix, filterHandler)
	http.HandleFunc(prioritizePrefix, prioritizeHandler)
	http.HandleFunc(preemptionPrefix, preemptionHandler)
	http.HandleFunc(versionPath, versionHandler)

	glog.Error("serving at localhost:8080")

	glog.Fatal(http.ListenAndServe(":8080", nil))
}

func checkBody(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
}

func filterHandler(w http.ResponseWriter, r *http.Request) {
	glog.Errorf("filter: %v", r)

	checkBody(w, r)

	var buf bytes.Buffer
	body := io.TeeReader(r.Body, &buf)
	glog.Infof("ExtenderArgs = ", buf.String())

	var extenderArgs schedulerapi.ExtenderArgs
	var extenderFilterResult *schedulerapi.ExtenderFilterResult

	if err := json.NewDecoder(body).Decode(&extenderArgs); err != nil {
		extenderFilterResult = &schedulerapi.ExtenderFilterResult{
			Nodes:       nil,
			FailedNodes: nil,
			Error:       err.Error(),
		}
	} else {
		extenderFilterResult = filter.Handler(extenderArgs)
	}

	if resultBody, err := json.Marshal(extenderFilterResult); err != nil {
		panic(err)
	} else {
		log.Print("info: ", filter.Name, " extenderFilterResult = ", string(resultBody))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(resultBody)
		cmn.AssertNoErr(err)
	}
}

func prioritizeHandler(w http.ResponseWriter, r *http.Request) {
	glog.Errorf("prioritize %v", r)

	checkBody(w, r)

	var buf bytes.Buffer
	body := io.TeeReader(r.Body, &buf)
	glog.Info("info: ", prioritize.Name, " ExtenderArgs = ", buf.String())

	var extenderArgs schedulerapi.ExtenderArgs
	var hostPriorityList *schedulerapi.HostPriorityList

	if err := json.NewDecoder(body).Decode(&extenderArgs); err != nil {
		panic(err)
	}

	if list, err := prioritize.Handler(extenderArgs); err != nil {
		glog.Errorf("prioritize handler failed")
		panic(err)
	} else {
		hostPriorityList = list
	}

	glog.Errorf("priority list: %v", hostPriorityList)

	if resultBody, err := json.Marshal(hostPriorityList); err != nil {
		panic(err)
	} else {
		glog.Info("info: ", prioritize.Name, " hostPriorityList = ", string(resultBody))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resultBody)
	}
}

func preemptionHandler(w http.ResponseWriter, r *http.Request) {
	glog.Errorf("%v", r)

	panic("preemption handler not implemented")
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("0.0.0"))
}


package main

import (
	"bytes"
	"encoding/json"
	"github.com/golang/glog"
	"io"
	v1 "k8s.io/api/core/v1"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"

	"fmt"
	"log"
	"net/http"
)

const (
	versionPath      = "/version"
	apiPrefix        = "/scheduler"
	bindPath         = apiPrefix + "/bind"
	filterPrefix     = apiPrefix + "/filter"
	prioritizePrefix = apiPrefix + "/prioritize"
	preemptionPrexif = apiPrefix + "/preemption"
)

var (
	filter = Filter{
		Name: "always_true",
		Func: func(node v1.Node, pod *v1.Pod) bool {
			fmt.Printf("returning true in filtering nodes")
			return true
		},
	}
	
	prioritize = Prioritize{
		Name: "always_1",
		Func: func(pod v1.Pod, nodes []v1.Node) (*schedulerapi.HostPriorityList, error) {
			var hostpriorityList schedulerapi.HostPriorityList
			hostpriorityList = make([]schedulerapi.HostPriority, len(nodes))
			for i, node := range nodes {
				hostpriorityList[i] = schedulerapi.HostPriority{
					Host:  node.Name,
					Score: i,
				}
			}

			glog.Infof("assigned priorities: %v\n", hostpriorityList)

			return &hostpriorityList, nil
		},
	}
)

func main() {

	//flag.Set("v", "2")
	//flag.Parse()

	glog.Error("starting main")
	http.HandleFunc(filterPrefix, filterHandler)
	http.HandleFunc(prioritizePrefix, prioritizeHandler)
	http.HandleFunc(preemptionPrexif, preemptionHandler)
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
	glog.Error("filter: %v", r)
	glog.Error(w, "OK filter")

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
		w.Write(resultBody)
	}
}

func prioritizeHandler(w http.ResponseWriter, r *http.Request) {
	glog.Infof("prioritize %v", r)

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
		panic(err)
	} else {
		hostPriorityList = list
	}

	if resultBody, err := json.Marshal(hostPriorityList); err != nil {
		panic(err)
	} else {
		glog.Info("info: ", prioritize.Name, " hostPriorityList = ", string(resultBody))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resultBody)
	}
}

func preemptionHandler(W http.ResponseWriter, r *http.Request) {
	glog.Errorf("%v", r)

	panic("preemption handler not implemented")
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	glog.Errorf("pinged for version")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("0.0.0"))
}


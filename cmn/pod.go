package cmn

import (
	"github.com/golang/glog"
	"github.com/knopt/k8s-sched-extender/coreclient"
	v1 "k8s.io/api/core/v1"

	"math"
)

type PodStats struct {
	Mean, Std float64
	Pod       *v1.Pod
}

func PodMeanStd(pod *v1.Pod) (*PodStats, error) {
	dist, err := FloatsFromString(pod.Annotations["dist"])
	if err != nil {
		return nil, err
	}
	return &PodStats{Mean: Mean(dist), Std: Std(dist), Pod: pod}, nil
}

// Returns pods stats and total requests for pods which don't have dist annotation set
func NodePodsStats(node string) ([]*PodStats, float64, error) {
	pods, err := coreclient.GetNodePods(node)
	if err != nil {
		return nil, 0, err
	}

	res := make([]*PodStats, 0, len(pods))

	// 1 CPU core == 1000 m
	totalRequestCpu := int64(0)

	for _, pod := range pods {
		if pod.Status.Phase != "Running" {
			continue
		}
		if _, ok := pod.Annotations["dist"]; !ok {
			for _, c := range pod.Spec.Containers {
				cpuMCores, ok := c.Resources.Requests.Cpu().AsInt64()
				if !ok {
					cpuMCores = c.Resources.Requests.Cpu().AsDec().UnscaledBig().Int64()
				}
				totalRequestCpu += cpuMCores
				//glog.Errorf("[%q] pod %q, cont %q, status %v occupies %fM cpu", node, pod.Name, c.Name, pod.Status.Phase, cpuMCores)
			}
			continue
		}
		stats, err := PodMeanStd(&pod)
		if err != nil {
			return nil, 0, err
		}
		res = append(res, stats)
	}

	return res, float64(totalRequestCpu) / 1000, nil
}

func PodNodeP(pod *v1.Pod, node v1.Node, p float64) (float64, error) {
	podStats, err := PodMeanStd(pod)
	if err != nil {
		return 0, err
	}
	stats, occCapCpu, err := NodePodsStats(node.Name)
	if err != nil {
		return 0, err
	}

	var meanSum, stdSum float64
	for _, stat := range stats {
		m := stat.Mean / 100
		std := stat.Std / 100
		meanSum += m
		stdSum += math.Pow(std, 2)
	}
	meanSum += podStats.Mean / 100
	stdSum += math.Pow(podStats.Std/100, 2)

	cap := ToCpuCores(node.Status.Capacity.Cpu())
	nodeCap := float64(cap) - occCapCpu
	nCDF := NormalCDF(meanSum, stdSum, nodeCap)
	glog.Errorf("CDF %q %f (m %f, std %f) total (m %f, std %f) av capacity: %f, cap %f, occ cap %f", pod.Name, nCDF, podStats.Mean/100, podStats.Std/100, meanSum, stdSum, nodeCap, float64(cap), occCapCpu)
	actualP := 1 - nCDF
	fitsP := p - actualP

	return fitsP, nil
}

func PodFitsNode(pod *v1.Pod, node v1.Node, p float64) (bool, error) {
	fitsP, err := PodNodeP(pod, node, p)
	if err != nil {
		return false, err
	}

	return fitsP >= 0, nil

}

package cmn

import (
	"github.com/golang/glog"
	"github.com/knopt/k8s-sched-extender/coreclient"
	v1 "k8s.io/api/core/v1"
	"math"
)

type PodStats struct {
	Mean, Std float64
}

func PodMeanStd(pod *v1.Pod) (*PodStats, error) {
	dist, err := FloatsFromString(pod.Annotations["dist"])
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	return &PodStats{Mean: Mean(dist), Std: Std(dist)}, nil
}

func NodePodsStats(node string) ([]*PodStats, error) {
	pods, err := coreclient.GetNodePods(node)
	if err != nil {
		return nil, err
	}

	res := make([]*PodStats, 0, len(pods))

	for _, pod := range pods {
		stats, err := PodMeanStd(&pod)
		if err != nil {
			return nil, err
		}
		res = append(res, stats)
	}

	return res, nil
}

func PodNodeP(pod *v1.Pod, node v1.Node, p float64) (float64, error) {
	podStats, err := PodMeanStd(pod)
	if err != nil {
		return 0, err
	}
	stats, err := NodePodsStats(node.Name)
	if err != nil {
		return 0, err
	}

	var meanSum, stdSum float64
	for _, stat := range stats {
		meanSum += stat.Mean
		stdSum += math.Pow(stat.Std, 2)
	}
	meanSum += podStats.Mean
	stdSum += math.Pow(podStats.Std, 2)

	//mean := meanSum / float64(len(stats) + 1)
	//std := math.Sqrt(stdSum)

	nodeCap := ToCpuCores(node.Status.Capacity.Cpu())
	actualP := 1 - NormalCDF(meanSum, stdSum, nodeCap)
	fitsP := p - actualP

	glog.Errorf("pod %s fits node %s: %d\n", pod.Name, node.Name, fitsP)
	return fitsP, nil
}

func PodFitsNode(pod *v1.Pod, node v1.Node, p float64) (bool, error) {
	fitsP, err := PodNodeP(pod, node, p)
	if err != nil {
		return false, err
	}

	return fitsP >= 0, nil

}

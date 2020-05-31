package stats

import (
	"k8s.io/kubernetes/pkg/kubelet/kubeletconfig/util/log"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"

	"github.com/golang/glog"

	"time"
)

const fetchInterval = 10 * time.Second

var Registry *StatsRunner

type (
	StatsRunner struct {
		abortCh chan struct{}

		nodesMetrics map[string]*NodeMetric
		//nodesMetrics *v1beta1.NodeMetricsList
		podsMetrics *v1beta1.PodMetricsList
	}
)

func NewStatsRunner() *StatsRunner {
	return &StatsRunner{
		abortCh: make(chan struct{}, 1),
	}
}

func (s *StatsRunner) Abort() {
	s.abortCh <- struct{}{}
}

func (s *StatsRunner) Run() {
	var (
		ticker        = time.NewTicker(fetchInterval)
		metricsGetter = NewMetrics()
	)

	defer ticker.Stop()

	log.Infof("StatsRunner starting with interval %s\n", fetchInterval)

	for {
		s.fetch(metricsGetter)
		select {
		case <-s.abortCh:
			ticker.Stop()
			close(s.abortCh)
			log.Errorf("StatsRunner aborted")
			return
		case <-ticker.C:
			// fetch stats on next loop iteration
			break
		}
	}
}

func (s *StatsRunner) fetch(metricsGetter *Metrics) {
	if nodeMetrics, err := metricsGetter.GetNodesMetrics(); err != nil {
		glog.Errorf("Failed to get nodes metrics: %v", err)
	} else {
		s.nodesMetrics = NodeMetricsFromInternal(nodeMetrics)
		glog.Warning("node metrics:\n%s", NodeMetrics2S(s.nodesMetrics))
	}
	if podsMetrics, err := metricsGetter.GetPodsMetrics(); err != nil {
		glog.Errorf("Failed to get pods metrics: %v", err)
	} else {
		s.podsMetrics = podsMetrics
		//glog.Warning("pod metrics: %v", cmn.PodMetrics2S(s.podsMetrics))
	}
}

func (s *StatsRunner) PodsMetrics() *v1beta1.PodMetricsList {
	return s.podsMetrics
}

func (s *StatsRunner) NodesMetrics() map[string]*NodeMetric {
	return s.nodesMetrics
}

func (s *StatsRunner) NodeMetrics(name string) *NodeMetric {
	if val, ok := s.nodesMetrics[name]; ok {
		return val
	}
	return nil
}

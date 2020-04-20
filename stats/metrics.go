package stats

import (
	"github.com/knopt/k8s-sched-extender/cmn"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"

)


type Metrics struct {
	clientSet *metrics.Clientset
}

func NewMetrics() *Metrics {
	config, err := rest.InClusterConfig()
	cmn.AssertNoErr(err)

	return &Metrics{clientSet: metrics.NewForConfigOrDie(config)}
}

func (p *Metrics) GetPodsMetrics() (*v1beta1.PodMetricsList, error) {
	return p.clientSet.MetricsV1beta1().PodMetricses("").List(metav1.ListOptions{})
}

func (p *Metrics) GetNodesMetrics() (*v1beta1.NodeMetricsList, error) {
	return p.clientSet.MetricsV1beta1().NodeMetricses().List(metav1.ListOptions{})
}

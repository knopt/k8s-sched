package stats

import (
	"fmt"
	"github.com/knopt/k8s-sched-extender/cmn"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

type (
	MemMetric struct {
		Bytes int64
	}

	CpuMetric struct {
		FracUsed float64
	}

	NodeMetric struct {
		Mem MemMetric
		Cpu CpuMetric
		NodeName string
	}
)

func (m MemMetric) String() string {
	return fmt.Sprintf("MEM %s (%d)", cmn.B2S(m.Bytes), m.Bytes)
}

func (m CpuMetric) String() string {
	return fmt.Sprintf("CPU %.3f", m.FracUsed)
}

func (m NodeMetric) String() string {
	return fmt.Sprintf("Node %s: %s; %s", m.NodeName, m.Mem.String(), m.Cpu.String())
}

func NodeMetricsFromInternal(metrics *v1beta1.NodeMetricsList) map[string]*NodeMetric {
	res := make(map[string]*NodeMetric, len(metrics.Items))

	for _, metric := range metrics.Items {
		cpuNanoCores, ok := metric.Usage.Cpu().AsInt64()
		if !ok {
			cpuNanoCores = metric.Usage.Cpu().AsDec().UnscaledBig().Int64()
		}
		cpuFloat := float64(cpuNanoCores) / (1000 * 1000 * 1000)

		memInt, ok := metric.Usage.Memory().AsInt64()
		if !ok {
			memInt = metric.Usage.Memory().AsDec().UnscaledBig().Int64()
		}

		res[metric.Name] = &NodeMetric{
			Mem:      MemMetric{Bytes: memInt},
			Cpu:      CpuMetric{FracUsed: cpuFloat},
			NodeName: metric.Name,
		}
	}

	return res
}

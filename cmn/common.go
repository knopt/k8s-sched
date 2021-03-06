package cmn

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

func AssertNoErr(err error) {
	if err != nil {
		panic(err)
	}
}

func BetaNodeMetrics2S(metrics *v1beta1.NodeMetricsList) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("nodes items length %d\n\n", len(metrics.Items)))

	for _, metric := range metrics.Items {
		sb.WriteString(BetaNodeMetric2S(&metric))
	}
	sb.WriteString("\n")
	return sb.String()
}

func BetaNodeMetric2S(metric *v1beta1.NodeMetrics) string {
	sb := strings.Builder{}

	cpuNanoCores, ok := metric.Usage.Cpu().AsInt64()
	if !ok {
		cpuNanoCores = metric.Usage.Cpu().AsDec().UnscaledBig().Int64()
	}
	cpuFloat := float64(cpuNanoCores) / (1000 * 1000 * 1000)

	memInt, ok := metric.Usage.Memory().AsInt64()
	if !ok {
		memInt = metric.Usage.Memory().AsDec().UnscaledBig().Int64()
	}

	sb.WriteString(metric.String() + "\n\n")
	sb.WriteString(fmt.Sprintf("Node %s; Kind %s; Object %s; TypeMeta Kind: %s\n", metric.Name, metric.Kind, metric.ObjectMeta.Name, metric.TypeMeta.Kind))
	sb.WriteString(fmt.Sprintf("CPU: %f (%s)\n", cpuFloat, metric.Usage.Cpu().String()))
	sb.WriteString(fmt.Sprintf("MEM: %s (%d)\n\n", metric.Usage.Memory().String(), memInt))

	return sb.String()
}

func PodMetrics2S(metrics *v1beta1.PodMetricsList) string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("pods items length %d", len(metrics.Items)))
	for _, metric := range metrics.Items {
		sb.WriteString(PodMetric2S(&metric))
	}

	sb.WriteString("\n")
	return sb.String()
}

func PodMetric2S(metric *v1beta1.PodMetrics) string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("Pod %s; Kind %s; Object %s; TypeMeta Kind: %s\n", metric.Name, metric.Kind, metric.ObjectMeta.Name, metric.TypeMeta.Kind))
	for _, val := range metric.Containers {
		sb.WriteString(fmt.Sprintf("CPU: %s: %s\n", val.Name, val.Usage.Cpu()))
		sb.WriteString(fmt.Sprintf("MEM: %s: %s\n\n", val.Name, val.Usage.Memory()))
	}

	return sb.String()
}

func B2S(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func FloatsFromString(s string) ([]float64, error) {
	var floats []float64
	err := json.Unmarshal([]byte(s), &floats)
	return floats, err
}

func P(pct float64, floats []float64) float64 {
	sort.Float64s(floats)

	idx := int(math.Ceil(pct * float64(len(floats))))
	idx = Min(idx, len(floats)-1)
	idx = Max(idx, 0)

	return floats[idx]
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func Mean(fs []float64) float64 {
	total := float64(0)
	for _, f := range fs {
		total += f
	}
	return total / float64(len(fs))
}

func Std(fs []float64) float64 {
	var sd float64
	mean := Mean(fs)

	for _, f := range fs {
		sd += math.Pow(f-mean, 2)
	}

	return math.Sqrt(sd / float64(len(fs)))
}

func ToCpuCores(cpuResource *resource.Quantity) int64 {
	cpuCores, ok := cpuResource.AsInt64()
	if !ok {
		cpuCores = cpuResource.AsDec().UnscaledBig().Int64()
	}
	return cpuCores
}

func NormalCDF(mean, std, x float64) float64 {
	return 0.5 * math.Erfc(-(x-mean)/(std*math.Sqrt2))
}

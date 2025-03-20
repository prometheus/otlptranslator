package prometheus

import (
	"fmt"
	"testing"

	"go.opentelemetry.io/collector/pdata/pmetric"
)

func BenchmarkBuildMetricNames(b *testing.B) {
	metrics := createTestMetrics()

	for _, workerType := range []string{"BuildCompliantMetricName", "BuildMetricName"} {
		for _, withSuffixes := range []bool{true, false} {
			b.Run(fmt.Sprintf("%s/withSuffixes=%t", workerType, withSuffixes), func(b *testing.B) {
				var worker Worker
				if workerType == "BuildCompliantMetricName" {
					worker = &BuildCompliantMetricNameWorker{
						metrics:      metrics,
						namespace:    "test_namespace",
						withSuffixes: withSuffixes,
					}
				} else {
					worker = &BuildMetricNameWorker{
						metrics:      metrics,
						namespace:    "test_namespace",
						withSuffixes: withSuffixes,
					}
				}

				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					worker.Work()
				}
			})
		}
	}
}

type Worker interface {
	Work()
}

type BuildCompliantMetricNameWorker struct {
	metrics      []pmetric.Metric
	namespace    string
	withSuffixes bool
}

func (w *BuildCompliantMetricNameWorker) Work() {
	for _, metric := range w.metrics {
		BuildCompliantMetricName(metric, w.namespace, w.withSuffixes)
	}
}

type BuildMetricNameWorker struct {
	metrics      []pmetric.Metric
	namespace    string
	withSuffixes bool
}

func (w *BuildMetricNameWorker) Work() {
	for _, metric := range w.metrics {
		BuildMetricName(metric, w.namespace, w.withSuffixes)
	}
}

func createTestMetrics() []pmetric.Metric {
	metrics := make([]pmetric.Metric, 0)

	// Basic metric with no special characters
	metric := pmetric.NewMetric()
	metric.SetName("simple_metric")
	metric.SetUnit("")
	metric.SetEmptyGauge()
	metrics = append(metrics, metric)

	// Counter metric (should get _total suffix)
	metric = pmetric.NewMetric()
	metric.SetName("counter_metric")
	metric.SetUnit("")
	sum := metric.SetEmptySum()
	sum.SetIsMonotonic(true)
	metrics = append(metrics, metric)

	// Metric with unit "1" (should get _ratio suffix for gauge)
	metric = pmetric.NewMetric()
	metric.SetName("ratio_metric")
	metric.SetUnit("1")
	metric.SetEmptyGauge()
	metrics = append(metrics, metric)

	// Metric with per-unit notation
	metric = pmetric.NewMetric()
	metric.SetName("requests")
	metric.SetUnit("1/s")
	metric.SetEmptyGauge()
	metrics = append(metrics, metric)

	// Metric with special characters
	metric = pmetric.NewMetric()
	metric.SetName("metric@with#special$chars")
	metric.SetUnit("")
	metric.SetEmptyGauge()
	metrics = append(metrics, metric)

	// Metric starting with digit
	metric = pmetric.NewMetric()
	metric.SetName("123metric")
	metric.SetUnit("")
	metric.SetEmptyGauge()
	metrics = append(metrics, metric)

	// Metric with multiple consecutive underscores
	metric = pmetric.NewMetric()
	metric.SetName("metric__with__multiple__underscores")
	metric.SetUnit("")
	metric.SetEmptyGauge()
	metrics = append(metrics, metric)

	// Metric with complex unit
	metric = pmetric.NewMetric()
	metric.SetName("memory_usage")
	metric.SetUnit("MiBy/s")
	metric.SetEmptyGauge()
	metrics = append(metrics, metric)

	// Metric with percentage unit
	metric = pmetric.NewMetric()
	metric.SetName("cpu_usage")
	metric.SetUnit("%")
	metric.SetEmptyGauge()
	metrics = append(metrics, metric)

	return metrics
}

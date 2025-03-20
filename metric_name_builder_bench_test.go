package otlptranslator

import (
	"fmt"
	"testing"

	"go.opentelemetry.io/collector/pdata/pmetric"
)

func BenchmarkBuildCompliantMetricName(b *testing.B) {
	metrics := createTestScenarios()

	for _, scenario := range metrics {
		for _, withSuffixes := range []bool{true, false} {
			b.Run(fmt.Sprintf("%s/withSuffixes=%t", scenario.name, withSuffixes), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					BuildCompliantMetricName(scenario.metric, "test_namespace", withSuffixes)
				}
			})
		}
	}
}

func BenchmarkBuildMetricName(b *testing.B) {
	metrics := createTestScenarios()

	for _, scenario := range metrics {
		for _, withSuffixes := range []bool{true, false} {
			b.Run(fmt.Sprintf("%s/withSuffixes=%t", scenario.name, withSuffixes), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					BuildMetricName(scenario.metric, "test_namespace", withSuffixes)
				}
			})
		}
	}
}

type Scenario struct {
	name   string
	metric pmetric.Metric
}

func createTestScenarios() []Scenario {
	scenarios := make([]Scenario, 0)

	metric := pmetric.NewMetric()
	metric.SetName("simple_metric")
	metric.SetUnit("")
	metric.SetEmptyGauge()
	scenarios = append(scenarios, Scenario{
		name:   "Basic metric with no special characters",
		metric: metric,
	})

	metric = pmetric.NewMetric()
	metric.SetName("counter_metric")
	metric.SetUnit("")
	sum := metric.SetEmptySum()
	sum.SetIsMonotonic(true)
	scenarios = append(scenarios, Scenario{
		name:   "Counter metric",
		metric: metric,
	})

	metric = pmetric.NewMetric()
	metric.SetName("ratio_metric")
	metric.SetUnit("1")
	metric.SetEmptyGauge()
	scenarios = append(scenarios, Scenario{
		name:   "Gauge ratio metric",
		metric: metric,
	})

	metric = pmetric.NewMetric()
	metric.SetName("requests")
	metric.SetUnit("1/s")
	metric.SetEmptyGauge()
	scenarios = append(scenarios, Scenario{
		name:   "Metric with per-unit suffix notation",
		metric: metric,
	})

	metric = pmetric.NewMetric()
	metric.SetName("metric@with#special$chars")
	metric.SetUnit("")
	metric.SetEmptyGauge()
	scenarios = append(scenarios, Scenario{
		name:   "Metric with special characters",
		metric: metric,
	})

	metric = pmetric.NewMetric()
	metric.SetName("123metric")
	metric.SetUnit("")
	metric.SetEmptyGauge()
	scenarios = append(scenarios, Scenario{
		name:   "Metric starting with digit",
		metric: metric,
	})

	metric = pmetric.NewMetric()
	metric.SetName("metric__with__multiple__underscores")
	metric.SetUnit("")
	metric.SetEmptyGauge()
	scenarios = append(scenarios, Scenario{
		name:   "Metric with multiple underscores",
		metric: metric,
	})

	metric = pmetric.NewMetric()
	metric.SetName("memory_usage")
	metric.SetUnit("MiBy/s")
	metric.SetEmptyGauge()
	scenarios = append(scenarios, Scenario{
		name:   "Metric with complex unit",
		metric: metric,
	})

	return scenarios
}

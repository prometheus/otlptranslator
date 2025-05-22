package otlptranslator

import (
	"fmt"
	"testing"
)

func BenchmarkBuild(b *testing.B) {
	scenarios := createTestScenarios()
	builder := MetricNameBuilder{Namespace: "test_namespace"}

	for _, withSuffixes := range []bool{true, false} {
		builder.WithMetricSuffixes = withSuffixes
		b.Run(fmt.Sprintf("withSuffixes=%t", withSuffixes), func(b *testing.B) {
			for _, utf8Allowed := range []bool{true, false} {
				builder.UTF8Allowed = utf8Allowed
				b.Run(fmt.Sprintf("utf8Allowed=%t", utf8Allowed), func(b *testing.B) {
					for _, scenario := range scenarios {
						b.Run(scenario.name, func(b *testing.B) {
							for i := 0; i < b.N; i++ {
								builder.buildCompliantMetricName(scenario.metricName, scenario.metricUnit, scenario.metricType)
							}
						})
					}
				})
			}
		})
	}
}

type scenario struct {
	name       string
	metricName string
	metricUnit string
	metricType MetricType
}

func createTestScenarios() []scenario {
	scenarios := make([]scenario, 0)

	scenarios = append(scenarios, scenario{
		name:       "Basic metric with no special characters",
		metricName: "simple_metric",
		metricUnit: "",
		metricType: MetricTypeGauge,
	})

	scenarios = append(scenarios, scenario{
		name:       "Counter metric",
		metricName: "counter_metric",
		metricUnit: "",
		metricType: MetricTypeMonotonicCounter,
	})

	scenarios = append(scenarios, scenario{
		name:       "Gauge ratio metric",
		metricName: "ratio_metric",
		metricUnit: "1",
		metricType: MetricTypeGauge,
	})

	scenarios = append(scenarios, scenario{
		name:       "Metric with per-unit suffix notation",
		metricName: "requests",
		metricUnit: "1/s",
		metricType: MetricTypeGauge,
	})

	scenarios = append(scenarios, scenario{
		name:       "Metric with special characters",
		metricName: "metric@with#special$chars",
		metricUnit: "",
		metricType: MetricTypeGauge,
	})

	scenarios = append(scenarios, scenario{
		name:       "Metric starting with digit",
		metricName: "123metric",
		metricUnit: "",
		metricType: MetricTypeGauge,
	})

	scenarios = append(scenarios, scenario{
		name:       "Metric with multiple underscores",
		metricName: "metric__with__multiple__underscores",
		metricUnit: "",
		metricType: MetricTypeGauge,
	})

	scenarios = append(scenarios, scenario{
		name:       "Metric with complex unit",
		metricName: "memory_usage",
		metricUnit: "MiBy/s",
		metricType: MetricTypeGauge,
	})

	return scenarios
}

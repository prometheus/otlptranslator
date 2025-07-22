package otlptranslator

import (
	"fmt"
	"testing"
)

func BenchmarkBuild(b *testing.B) {
	scenarios := createTestScenarios()
	builder := MetricNamer{Namespace: "test_namespace"}

	for _, withSuffixes := range []bool{true, false} {
		builder.WithMetricSuffixes = withSuffixes
		b.Run(fmt.Sprintf("withSuffixes=%t", withSuffixes), func(b *testing.B) {
			for _, utf8Allowed := range []bool{true, false} {
				builder.UTF8Allowed = utf8Allowed
				b.Run(fmt.Sprintf("utf8Allowed=%t", utf8Allowed), func(b *testing.B) {
					for _, scenario := range scenarios {
						b.Run(scenario.name, func(b *testing.B) {
							for i := 0; i < b.N; i++ {
								//nolint:errcheck
								builder.Build(scenario.metric)
							}
						})
					}
				})
			}
		})
	}
}

type scenario struct {
	name   string
	metric Metric
}

func createTestScenarios() []scenario {
	scenarios := make([]scenario, 0)

	scenarios = append(scenarios, scenario{
		name: "Basic metric with no special characters",
		metric: Metric{
			Name: "simple_metric",
			Unit: "",
			Type: MetricTypeGauge,
		},
	})

	scenarios = append(scenarios, scenario{
		name: "Counter metric",
		metric: Metric{
			Name: "counter_metric",
			Unit: "",
			Type: MetricTypeMonotonicCounter,
		},
	})

	scenarios = append(scenarios, scenario{
		name: "Gauge ratio metric",
		metric: Metric{
			Name: "ratio_metric",
			Unit: "1",
			Type: MetricTypeGauge,
		},
	})

	scenarios = append(scenarios, scenario{
		name: "Metric with per-unit suffix notation",
		metric: Metric{
			Name: "requests",
			Unit: "1/s",
			Type: MetricTypeGauge,
		},
	})

	scenarios = append(scenarios, scenario{
		name: "Metric with special characters",
		metric: Metric{
			Name: "metric@with#special$chars",
			Unit: "",
			Type: MetricTypeGauge,
		},
	})

	scenarios = append(scenarios, scenario{
		name: "Metric starting with digit",
		metric: Metric{
			Name: "123metric",
			Unit: "",
			Type: MetricTypeGauge,
		},
	})

	scenarios = append(scenarios, scenario{
		name: "Metric with multiple underscores",
		metric: Metric{
			Name: "metric__with__multiple__underscores",
			Unit: "",
			Type: MetricTypeGauge,
		},
	})

	scenarios = append(scenarios, scenario{
		name: "Metric with complex unit",
		metric: Metric{
			Name: "memory_usage",
			Unit: "MiBy/s",
			Type: MetricTypeGauge,
		},
	})

	return scenarios
}

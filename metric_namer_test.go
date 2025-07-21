// Copyright 2025 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Provenance-includes-location: https://github.com/prometheus/prometheus/blob/93e991ef7ed19cc997a9360c8016cac3767b8057/storage/remote/otlptranslator/prometheus/metric_name_builder_test.go
// Provenance-includes-license: Apache-2.0
// Provenance-includes-copyright: Copyright The Prometheus Authors
// Provenance-includes-location: https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/95e8f8fdc2a9dc87230406c9a3cf02be4fd68bea/pkg/translator/prometheus/normalize_name_test.go
// Provenance-includes-license: Apache-2.0
// Provenance-includes-copyright: Copyright The OpenTelemetry Authors.

package otlptranslator

import (
	"strings"
	"testing"
)

func TestMetricNamer_Build(t *testing.T) {
	tests := []struct {
		name           string
		namer          MetricNamer
		metric         Metric
		wantMetricName string
		wantUnitName   string
	}{
		// UTF8Allowed = false, WithMetricSuffixes = false tests
		{
			name: "simple metric name without suffixes",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: false,
			},
			metric: Metric{
				Name: "simple_metric",
				Unit: "",
				Type: MetricTypeGauge,
			},
			wantMetricName: "simple_metric",
		},
		{
			name: "metric with special characters replaced",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: false,
			},
			metric: Metric{
				Name: "metric@with#special$chars",
				Unit: "",
				Type: MetricTypeGauge,
			},
			wantMetricName: "metric_with_special_chars",
		},
		{
			name: "metric starting with digit gets underscore prefix",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: false,
			},
			metric: Metric{
				Name: "123metric",
				Unit: "",
				Type: MetricTypeGauge,
			},
			wantMetricName: "_123metric",
		},
		{
			name: "metric with namespace without suffixes",
			namer: MetricNamer{
				Namespace:          "test_namespace",
				UTF8Allowed:        false,
				WithMetricSuffixes: false,
			},
			metric: Metric{
				Name: "simple_metric",
				Unit: "",
				Type: MetricTypeGauge,
			},
			wantMetricName: "test_namespace_simple_metric",
		},
		{
			name: "empty metric name without suffixes",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: false,
			},
			metric: Metric{
				Name: "",
				Unit: "",
				Type: MetricTypeGauge,
			},
			wantMetricName: "",
		},
		{
			name: "metric with multiple consecutive special chars",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: false,
			},
			metric: Metric{
				Name: "metric@@##$$name",
				Unit: "",
				Type: MetricTypeGauge,
			},
			wantMetricName: "metric_name",
		},
		{
			name: "metric with multiple consecutive special chars/keep multiple underscores",
			namer: MetricNamer{
				UTF8Allowed:             false,
				WithMetricSuffixes:      false,
				KeepMultipleUnderscores: true,
			},
			metric: Metric{
				Name: "metric@@##$$name",
				Unit: "",
				Type: MetricTypeGauge,
			},
			wantMetricName: "metric______name",
		},
		{
			name: "metric name with only special characters",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: false,
			},
			metric: Metric{
				Name: "@#$%",
				Unit: "",
				Type: MetricTypeGauge,
			},
			wantMetricName: "",
		},
		{
			name: "namespace with special characters",
			namer: MetricNamer{
				Namespace:          "test@namespace!!??",
				UTF8Allowed:        false,
				WithMetricSuffixes: false,
			},
			metric: Metric{
				Name: "metric",
				Unit: "",
				Type: MetricTypeGauge,
			},
			wantMetricName: "test_namespace_metric",
		},

		// UTF8Allowed = false, WithMetricSuffixes = true tests
		{
			name: "counter metric with total suffix",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "requests",
				Unit: "",
				Type: MetricTypeMonotonicCounter,
			},
			wantMetricName: "requests_total",
		},
		{
			name: "gauge with unit 1 gets ratio suffix",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "cpu_usage",
				Unit: "1",
				Type: MetricTypeGauge,
			},
			wantMetricName: "cpu_usage_ratio",
		},
		{
			name: "counter with unit 1 does not get ratio suffix",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "items",
				Unit: "1",
				Type: MetricTypeMonotonicCounter,
			},
			wantMetricName: "items_total",
		},
		{
			name: "metric with time unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "response_time",
				Unit: "ms",
				Type: MetricTypeGauge,
			},
			wantMetricName: "response_time_milliseconds",
			wantUnitName:   "milliseconds",
		},
		{
			name: "metric with bytes unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "memory_usage",
				Unit: "By",
				Type: MetricTypeGauge,
			},
			wantMetricName: "memory_usage_bytes",
			wantUnitName:   "bytes",
		},
		{
			name: "metric with per unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "requests",
				Unit: "1/s",
				Type: MetricTypeGauge,
			},
			wantMetricName: "requests_per_second",
			wantUnitName:   "per_second",
		},
		{
			name: "metric with complex per unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "throughput",
				Unit: "By/s",
				Type: MetricTypeGauge,
			},
			wantMetricName: "throughput_bytes_per_second",
			wantUnitName:   "bytes_per_second",
		},
		{
			name: "metric with unknown unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "custom_metric",
				Unit: "custom_unit",
				Type: MetricTypeGauge,
			},
			wantMetricName: "custom_metric_custom_unit",
			wantUnitName:   "custom_unit",
		},
		{
			name: "metric with unit containing braces is ignored",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "custom_metric",
				Unit: "{custom}",
				Type: MetricTypeGauge,
			},
			wantMetricName: "custom_metric",
		},
		{
			name: "metric with per unit containing braces is ignored",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "custom_metric",
				Unit: "By/{custom}",
				Type: MetricTypeGauge,
			},
			wantMetricName: "custom_metric_bytes",
			wantUnitName:   "bytes",
		},
		{
			name: "metric name already contains total suffix",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "requests_total",
				Unit: "",
				Type: MetricTypeMonotonicCounter,
			},
			wantMetricName: "requests_total",
		},
		{
			name: "metric name already contains ratio suffix",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "cpu_usage_ratio",
				Unit: "1",
				Type: MetricTypeGauge,
			},
			wantMetricName: "cpu_usage_ratio",
		},
		{
			name: "metric name already contains unit suffix",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "response_time_seconds",
				Unit: "s",
				Type: MetricTypeGauge,
			},
			wantMetricName: "response_time_seconds",
			wantUnitName:   "seconds",
		},
		{
			name: "metric with namespace and suffixes",
			namer: MetricNamer{
				Namespace:          "app",
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "requests",
				Unit: "1/s",
				Type: MetricTypeMonotonicCounter,
			},
			wantMetricName: "app_requests_per_second_total",
			wantUnitName:   "per_second",
		},
		{
			name: "metric starting with digit with namespace and suffixes",
			namer: MetricNamer{
				Namespace:          "app",
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "123_requests",
				Unit: "",
				Type: MetricTypeMonotonicCounter,
			},
			wantMetricName: "app_123_requests_total",
		},
		{
			name: "metric with multiple underscores normalized",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "metric__with__multiple__underscores",
				Unit: "unit__multiple__underscores",
				Type: MetricTypeGauge,
			},
			wantMetricName: "metric_with_multiple_underscores_unit_multiple_underscores",
			wantUnitName:   "unit_multiple_underscores",
		},
		{
			name: "metric with multiple underscores normalized/keep multiple underscores",
			namer: MetricNamer{
				UTF8Allowed:             false,
				WithMetricSuffixes:      true,
				KeepMultipleUnderscores: true,
			},
			metric: Metric{
				Name: "metric__with__multiple__underscores",
				Unit: "unit__multiple__underscores",
				Type: MetricTypeGauge,
			},
			wantMetricName: "metric__with__multiple__underscores_unit__multiple__underscores",
			wantUnitName:   "unit__multiple__underscores",
		},
		{
			name: "metric with special chars in unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "custom_metric",
				Unit: "unit@with#special/chars",
				Type: MetricTypeGauge,
			},
			wantMetricName: "custom_metric_unit_with_special_per_chars",
			wantUnitName:   "unit_with_special_per_chars",
		},
		{
			name: "metric name with only special characters",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "@#$%",
				Unit: "",
				Type: MetricTypeGauge,
			},
			wantMetricName: "",
		},

		// UTF8Allowed = true, WithMetricSuffixes = false tests
		{
			name: "utf8 metric without suffixes",
			namer: MetricNamer{
				UTF8Allowed:        true,
				WithMetricSuffixes: false,
			},
			metric: Metric{
				Name: "métric_with_ñ_chars",
				Unit: "",
				Type: MetricTypeGauge,
			},
			wantMetricName: "métric_with_ñ_chars",
		},
		{
			name: "utf8 metric with namespace without suffixes",
			namer: MetricNamer{
				Namespace:          "test_namespace",
				UTF8Allowed:        true,
				WithMetricSuffixes: false,
			},
			metric: Metric{
				Name: "métric_with_ñ_chars",
				Unit: "",
				Type: MetricTypeGauge,
			},
			wantMetricName: "test_namespace_métric_with_ñ_chars",
		},

		// UTF8Allowed = true, WithMetricSuffixes = true tests
		{
			name: "utf8 counter metric with total suffix",
			namer: MetricNamer{
				UTF8Allowed:        true,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "requêsts",
				Unit: "",
				Type: MetricTypeMonotonicCounter,
			},
			wantMetricName: "requêsts_total",
		},
		{
			name: "utf8 gauge with unit 1 gets ratio suffix",
			namer: MetricNamer{
				UTF8Allowed:        true,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "cpu_usagé",
				Unit: "1",
				Type: MetricTypeGauge,
			},
			wantMetricName: "cpu_usagé_ratio",
		},
		{
			name: "utf8 metric with time unit",
			namer: MetricNamer{
				UTF8Allowed:        true,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "respønse_time",
				Unit: "ms",
				Type: MetricTypeGauge,
			},
			wantMetricName: "respønse_time_milliseconds",
			wantUnitName:   "milliseconds",
		},
		{
			name: "utf8 metric with per unit",
			namer: MetricNamer{
				UTF8Allowed:        true,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "requêsts",
				Unit: "1/s",
				Type: MetricTypeGauge,
			},
			wantMetricName: "requêsts_per_second",
			wantUnitName:   "per_second",
		},
		{
			name: "utf8 metric with namespace and suffixes",
			namer: MetricNamer{
				Namespace:          "ñamespace",
				UTF8Allowed:        true,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "requêsts",
				Unit: "1/s",
				Type: MetricTypeMonotonicCounter,
			},
			wantMetricName: "ñamespace_requêsts_per_second_total",
			wantUnitName:   "per_second",
		},
		{
			name: "metric name with only special characters",
			namer: MetricNamer{
				UTF8Allowed:        true,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "@#$%",
				Unit: "",
				Type: MetricTypeMonotonicCounter,
			},
			wantMetricName: "@#$%_total",
		},
		{
			name: "namespace with special characters",
			namer: MetricNamer{
				Namespace:          "test@namespace",
				UTF8Allowed:        true,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "metric",
				Unit: "",
				Type: MetricTypeMonotonicCounter,
			},
			wantMetricName: "test@namespace_metric_total",
		},

		// Edge cases and different metric types
		{
			name: "histogram metric type",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "request_duration",
				Unit: "s",
				Type: MetricTypeHistogram,
			},
			wantMetricName: "request_duration_seconds",
			wantUnitName:   "seconds",
		},
		{
			name: "exponential histogram metric type",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "request_size",
				Unit: "By",
				Type: MetricTypeExponentialHistogram,
			},
			wantMetricName: "request_size_bytes",
			wantUnitName:   "bytes",
		},
		{
			name: "summary metric type",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "response_time",
				Unit: "ms",
				Type: MetricTypeSummary,
			},
			wantMetricName: "response_time_milliseconds",
			wantUnitName:   "milliseconds",
		},
		{
			name: "non-monotonic counter metric type",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "active_connections",
				Unit: "",
				Type: MetricTypeNonMonotonicCounter,
			},
			wantMetricName: "active_connections",
		},
		{
			name: "unknown metric type",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "unknown_metric",
				Unit: "",
				Type: MetricTypeUnknown,
			},
			wantMetricName: "unknown_metric",
		},

		// Additional unit mapping tests
		{
			name: "metric with days unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "uptime",
				Unit: "d",
				Type: MetricTypeGauge,
			},
			wantMetricName: "uptime_days",
			wantUnitName:   "days",
		},
		{
			name: "metric with hours unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "duration",
				Unit: "h",
				Type: MetricTypeGauge,
			},
			wantMetricName: "duration_hours",
			wantUnitName:   "hours",
		},
		{
			name: "metric with minutes unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "timeout",
				Unit: "min",
				Type: MetricTypeGauge,
			},
			wantMetricName: "timeout_minutes",
			wantUnitName:   "minutes",
		},
		{
			name: "metric with microseconds unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "latency",
				Unit: "us",
				Type: MetricTypeGauge,
			},
			wantMetricName: "latency_microseconds",
			wantUnitName:   "microseconds",
		},
		{
			name: "metric with nanoseconds unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "precision_time",
				Unit: "ns",
				Type: MetricTypeGauge,
			},
			wantMetricName: "precision_time_nanoseconds",
			wantUnitName:   "nanoseconds",
		},
		{
			name: "metric with kibibytes unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "cache_size",
				Unit: "KiBy",
				Type: MetricTypeGauge,
			},
			wantMetricName: "cache_size_kibibytes",
			wantUnitName:   "kibibytes",
		},
		{
			name: "metric with mebibytes unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "memory",
				Unit: "MiBy",
				Type: MetricTypeGauge,
			},
			wantMetricName: "memory_mebibytes",
			wantUnitName:   "mebibytes",
		},
		{
			name: "metric with gibibytes unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "storage",
				Unit: "GiBy",
				Type: MetricTypeGauge,
			},
			wantMetricName: "storage_gibibytes",
			wantUnitName:   "gibibytes",
		},
		{
			name: "metric with tibibytes unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "capacity",
				Unit: "TiBy",
				Type: MetricTypeGauge,
			},
			wantMetricName: "capacity_tibibytes",
			wantUnitName:   "tibibytes",
		},
		{
			name: "metric with kilobytes unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "transfer",
				Unit: "KBy",
				Type: MetricTypeGauge,
			},
			wantMetricName: "transfer_kilobytes",
			wantUnitName:   "kilobytes",
		},
		{
			name: "metric with megabytes unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "download",
				Unit: "MBy",
				Type: MetricTypeGauge,
			},
			wantMetricName: "download_megabytes",
			wantUnitName:   "megabytes",
		},
		{
			name: "metric with gigabytes unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "backup",
				Unit: "GBy",
				Type: MetricTypeGauge,
			},
			wantMetricName: "backup_gigabytes",
			wantUnitName:   "gigabytes",
		},
		{
			name: "metric with terabytes unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "archive",
				Unit: "TBy",
				Type: MetricTypeGauge,
			},
			wantMetricName: "archive_terabytes",
			wantUnitName:   "terabytes",
		},
		{
			name: "metric with meters unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "distance",
				Unit: "m",
				Type: MetricTypeGauge,
			},
			wantMetricName: "distance_meters",
			wantUnitName:   "meters",
		},
		{
			name: "metric with volts unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "voltage",
				Unit: "V",
				Type: MetricTypeGauge,
			},
			wantMetricName: "voltage_volts",
			wantUnitName:   "volts",
		},
		{
			name: "metric with amperes unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "current",
				Unit: "A",
				Type: MetricTypeGauge,
			},
			wantMetricName: "current_amperes",
			wantUnitName:   "amperes",
		},
		{
			name: "metric with joules unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "energy",
				Unit: "J",
				Type: MetricTypeGauge,
			},
			wantMetricName: "energy_joules",
			wantUnitName:   "joules",
		},
		{
			name: "metric with watts unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "power",
				Unit: "W",
				Type: MetricTypeGauge,
			},
			wantMetricName: "power_watts",
			wantUnitName:   "watts",
		},
		{
			name: "metric with grams unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "weight",
				Unit: "g",
				Type: MetricTypeGauge,
			},
			wantMetricName: "weight_grams",
			wantUnitName:   "grams",
		},
		{
			name: "metric with celsius unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "temperature",
				Unit: "Cel",
				Type: MetricTypeGauge,
			},
			wantMetricName: "temperature_celsius",
			wantUnitName:   "celsius",
		},
		{
			name: "metric with hertz unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "frequency",
				Unit: "Hz",
				Type: MetricTypeGauge,
			},
			wantMetricName: "frequency_hertz",
			wantUnitName:   "hertz",
		},
		{
			name: "metric with percent unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "cpu_usage",
				Unit: "%",
				Type: MetricTypeGauge,
			},
			wantMetricName: "cpu_usage_percent",
			wantUnitName:   "percent",
		},

		// Per unit mapping tests
		{
			name: "metric with per minute unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "requests",
				Unit: "1/m",
				Type: MetricTypeGauge,
			},
			wantMetricName: "requests_per_minute",
			wantUnitName:   "per_minute",
		},
		{
			name: "metric with per hour unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "events",
				Unit: "1/h",
				Type: MetricTypeGauge,
			},
			wantMetricName: "events_per_hour",
			wantUnitName:   "per_hour",
		},
		{
			name: "metric with per day unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "transactions",
				Unit: "1/d",
				Type: MetricTypeGauge,
			},
			wantMetricName: "transactions_per_day",
			wantUnitName:   "per_day",
		},
		{
			name: "metric with per week unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "reports",
				Unit: "1/w",
				Type: MetricTypeGauge,
			},
			wantMetricName: "reports_per_week",
			wantUnitName:   "per_week",
		},
		{
			name: "metric with per month unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "invoices",
				Unit: "1/mo",
				Type: MetricTypeGauge,
			},
			wantMetricName: "invoices_per_month",
			wantUnitName:   "per_month",
		},
		{
			name: "metric with per year unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "renewals",
				Unit: "1/y",
				Type: MetricTypeGauge,
			},
			wantMetricName: "renewals_per_year",
			wantUnitName:   "per_year",
		},
		{
			name: "metric with unknown per unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "custom",
				Unit: "1/custom_unit",
				Type: MetricTypeGauge,
			},
			wantMetricName: "custom_per_custom_unit",
			wantUnitName:   "per_custom_unit",
		},

		// Edge cases with empty and whitespace units
		{
			name: "metric with empty per unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "metric",
				Unit: "By/",
				Type: MetricTypeGauge,
			},
			wantMetricName: "metric_bytes",
			wantUnitName:   "bytes",
		},
		{
			name: "metric with whitespace in unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "metric",
				Unit: " By / s ",
				Type: MetricTypeGauge,
			},
			wantMetricName: "metric_bytes_per_second",
			wantUnitName:   "bytes_per_second",
		},
		{
			name: "metric with only slash in unit",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "metric",
				Unit: "/",
				Type: MetricTypeGauge,
			},
			wantMetricName: "metric",
		},

		// Common OTel metrics to showcase how the namer works
		{
			name: "http.request.duration/Prometheus-style",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "http.request.duration",
				Unit: "ms",
				Type: MetricTypeHistogram,
			},
			wantMetricName: "http_request_duration_milliseconds",
			wantUnitName:   "milliseconds",
		},
		{
			name: "http.request.duration/OTel-style",
			namer: MetricNamer{
				UTF8Allowed:        true,
				WithMetricSuffixes: false,
			},
			metric: Metric{
				Name: "http.request.duration",
				Unit: "ms",
				Type: MetricTypeHistogram,
			},
			wantMetricName: "http.request.duration",
			wantUnitName:   "milliseconds",
		},
		{
			name: "http.requests/Prometheus-style",
			namer: MetricNamer{
				UTF8Allowed:        false,
				WithMetricSuffixes: true,
			},
			metric: Metric{
				Name: "http.requests",
				Unit: "1",
				Type: MetricTypeMonotonicCounter,
			},
			wantMetricName: "http_requests_total",
		},
		{
			name: "http.requests/OTel-style",
			namer: MetricNamer{
				UTF8Allowed:        true,
				WithMetricSuffixes: false,
			},
			metric: Metric{
				Name: "http.requests",
				Unit: "1",
				Type: MetricTypeMonotonicCounter,
			},
			wantMetricName: "http.requests",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build metric name using MetricNamer
			gotMetricName := tt.namer.Build(tt.metric)
			if gotMetricName != tt.wantMetricName {
				t.Errorf("MetricNamer.Build(%v) = %q, want %q", tt.metric, gotMetricName, tt.wantMetricName)
			}

			// Build unit name using UnitNamer to verify correlation when suffixes are enabled
			if tt.namer.WithMetricSuffixes {
				unitNamer := UnitNamer{
					UTF8Allowed:             tt.namer.UTF8Allowed,
					KeepMultipleUnderscores: tt.namer.KeepMultipleUnderscores,
				}
				gotUnitName := unitNamer.Build(tt.metric.Unit)
				if gotUnitName != tt.wantUnitName {
					t.Errorf("UnitNamer.Build(%q) = %q, want %q", tt.metric.Unit, gotUnitName, tt.wantUnitName)
				}

				// Verify correlation: if UnitNamer produces a non-empty unit name,
				// it should be contained in the metric name when WithMetricSuffixes=true
				if tt.namer.WithMetricSuffixes && !strings.Contains(gotMetricName, gotUnitName) {
					t.Errorf("Metric name %q should contain unit name %q when WithMetricSuffixes=true", gotMetricName, gotUnitName)
				}
			}
		})
	}
}

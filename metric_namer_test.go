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
		name               string
		namer              MetricNamer
		metric             Metric
		expectedMetricName string
		expectedUnitName   string
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
			expectedMetricName: "simple_metric",
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
			expectedMetricName: "metric_with_special_chars",
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
			expectedMetricName: "_123metric",
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
			expectedMetricName: "test_namespace_simple_metric",
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
			expectedMetricName: "",
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
			expectedMetricName: "metric_name",
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
			expectedMetricName: "",
		},
		{
			name: "namespace with special characters",
			namer: MetricNamer{
				Namespace:          "test@namespace",
				UTF8Allowed:        false,
				WithMetricSuffixes: false,
			},
			metric: Metric{
				Name: "metric",
				Unit: "",
				Type: MetricTypeGauge,
			},
			expectedMetricName: "test@namespace_metric", // TODO: should be "test_namespace_metric"

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
			expectedMetricName: "requests_total",
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
			expectedMetricName: "cpu_usage_ratio",
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
			expectedMetricName: "items_total",
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
			expectedMetricName: "response_time_milliseconds",
			expectedUnitName:   "milliseconds",
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
			expectedMetricName: "memory_usage_bytes",
			expectedUnitName:   "bytes",
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
			expectedMetricName: "requests_per_second",
			expectedUnitName:   "per_second",
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
			expectedMetricName: "throughput_bytes_per_second",
			expectedUnitName:   "bytes_per_second",
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
			expectedMetricName: "custom_metric_custom_unit",
			expectedUnitName:   "custom_unit",
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
			expectedMetricName: "custom_metric",
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
			expectedMetricName: "custom_metric_bytes",
			expectedUnitName:   "bytes",
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
			expectedMetricName: "requests_total",
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
			expectedMetricName: "cpu_usage_ratio",
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
			expectedMetricName: "response_time_seconds",
			expectedUnitName:   "seconds",
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
			expectedMetricName: "app_requests_per_second_total",
			expectedUnitName:   "per_second",
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
			expectedMetricName: "app_123_requests_total",
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
			expectedMetricName: "metric_with_multiple_underscores_unit_multiple_underscores",
			expectedUnitName:   "unit_multiple_underscores",
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
			expectedMetricName: "custom_metric_unit_with_special_per_chars",
			expectedUnitName:   "unit_with_special_per_chars",
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
			expectedMetricName: "",
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
			expectedMetricName: "métric_with_ñ_chars",
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
			expectedMetricName: "test_namespace_métric_with_ñ_chars",
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
			expectedMetricName: "requêsts_total",
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
			expectedMetricName: "cpu_usagé_ratio",
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
			expectedMetricName: "respønse_time_milliseconds",
			expectedUnitName:   "milliseconds",
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
			expectedMetricName: "requêsts_per_second",
			expectedUnitName:   "per_second",
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
			expectedMetricName: "ñamespace_requêsts_per_second_total",
			expectedUnitName:   "per_second",
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
			expectedMetricName: "@#$%_total",
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
			expectedMetricName: "test@namespace_metric_total",
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
			expectedMetricName: "request_duration_seconds",
			expectedUnitName:   "seconds",
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
			expectedMetricName: "request_size_bytes",
			expectedUnitName:   "bytes",
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
			expectedMetricName: "response_time_milliseconds",
			expectedUnitName:   "milliseconds",
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
			expectedMetricName: "active_connections",
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
			expectedMetricName: "unknown_metric",
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
			expectedMetricName: "uptime_days",
			expectedUnitName:   "days",
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
			expectedMetricName: "duration_hours",
			expectedUnitName:   "hours",
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
			expectedMetricName: "timeout_minutes",
			expectedUnitName:   "minutes",
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
			expectedMetricName: "latency_microseconds",
			expectedUnitName:   "microseconds",
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
			expectedMetricName: "precision_time_nanoseconds",
			expectedUnitName:   "nanoseconds",
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
			expectedMetricName: "cache_size_kibibytes",
			expectedUnitName:   "kibibytes",
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
			expectedMetricName: "memory_mebibytes",
			expectedUnitName:   "mebibytes",
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
			expectedMetricName: "storage_gibibytes",
			expectedUnitName:   "gibibytes",
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
			expectedMetricName: "capacity_tibibytes",
			expectedUnitName:   "tibibytes",
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
			expectedMetricName: "transfer_kilobytes",
			expectedUnitName:   "kilobytes",
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
			expectedMetricName: "download_megabytes",
			expectedUnitName:   "megabytes",
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
			expectedMetricName: "backup_gigabytes",
			expectedUnitName:   "gigabytes",
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
			expectedMetricName: "archive_terabytes",
			expectedUnitName:   "terabytes",
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
			expectedMetricName: "distance_meters",
			expectedUnitName:   "meters",
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
			expectedMetricName: "voltage_volts",
			expectedUnitName:   "volts",
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
			expectedMetricName: "current_amperes",
			expectedUnitName:   "amperes",
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
			expectedMetricName: "energy_joules",
			expectedUnitName:   "joules",
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
			expectedMetricName: "power_watts",
			expectedUnitName:   "watts",
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
			expectedMetricName: "weight_grams",
			expectedUnitName:   "grams",
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
			expectedMetricName: "temperature_celsius",
			expectedUnitName:   "celsius",
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
			expectedMetricName: "frequency_hertz",
			expectedUnitName:   "hertz",
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
			expectedMetricName: "cpu_usage_percent",
			expectedUnitName:   "percent",
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
			expectedMetricName: "requests_per_minute",
			expectedUnitName:   "per_minute",
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
			expectedMetricName: "events_per_hour",
			expectedUnitName:   "per_hour",
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
			expectedMetricName: "transactions_per_day",
			expectedUnitName:   "per_day",
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
			expectedMetricName: "reports_per_week",
			expectedUnitName:   "per_week",
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
			expectedMetricName: "invoices_per_month",
			expectedUnitName:   "per_month",
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
			expectedMetricName: "renewals_per_year",
			expectedUnitName:   "per_year",
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
			expectedMetricName: "custom_per_custom_unit",
			expectedUnitName:   "per_custom_unit",
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
			expectedMetricName: "metric_bytes",
			expectedUnitName:   "bytes",
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
			expectedMetricName: "metric_bytes_per_second",
			expectedUnitName:   "bytes_per_second",
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
			expectedMetricName: "metric",
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
			expectedMetricName: "http_request_duration_milliseconds",
			expectedUnitName:   "milliseconds",
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
			expectedMetricName: "http.request.duration",
			expectedUnitName:   "milliseconds",
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
			expectedMetricName: "http_requests_total",
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
			expectedMetricName: "http.requests",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build metric name using MetricNamer
			gotMetricName := tt.namer.Build(tt.metric)
			if tt.expectedMetricName != gotMetricName {
				t.Errorf("namer.Build(%v), got %q, want %q", tt.metric, gotMetricName, tt.expectedMetricName)
			}

			// Build unit name using UnitNamer to verify correlation when suffixes are enabled
			if tt.namer.WithMetricSuffixes {
				unitNamer := UnitNamer{
					UTF8Allowed: tt.namer.UTF8Allowed,
				}
				gotUnitName := unitNamer.Build(tt.metric.Unit)
				if tt.expectedUnitName != gotUnitName {
					t.Errorf("unitNamer.Build(%q), got %q, want %q", tt.metric.Unit, gotUnitName, tt.expectedUnitName)
				}

				// Verify correlation: if UnitNamer produces a non-empty unit name,
				// it should be contained in the metric name when WithMetricSuffixes=true
				if gotUnitName != "" && tt.namer.WithMetricSuffixes {
					if !strings.Contains(gotMetricName, gotUnitName) {
						t.Errorf("Metric name '%q' should contain unit name '%s' when WithMetricSuffixes=true", gotMetricName, gotUnitName)
					}
				}
			}
		})
	}
}

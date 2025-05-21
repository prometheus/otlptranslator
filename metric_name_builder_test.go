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
	"testing"

	"github.com/stretchr/testify/require"
)

func TestByte(t *testing.T) {
	require.Equal(t, "system_filesystem_usage_bytes", normalizeName("system.filesystem.usage", "By", MetricTypeGauge, ""))
}

func TestByteCounter(t *testing.T) {
	require.Equal(t, "system_io_bytes_total", normalizeName("system.io", "By", MetricTypeMonotonicCounter, ""))
	require.Equal(t, "network_transmitted_bytes_total", normalizeName("network_transmitted_bytes_total", "By", MetricTypeMonotonicCounter, ""))
}

func TestWhiteSpaces(t *testing.T) {
	require.Equal(t, "system_filesystem_usage_bytes", normalizeName("\t system.filesystem.usage       ", "  By\t", MetricTypeGauge, ""))
}

func TestNonStandardUnit(t *testing.T) {
	require.Equal(t, "system_network_dropped", normalizeName("system.network.dropped", "{packets}", MetricTypeGauge, ""))
	// The normal metric name character set is allowed in non-standard units.
	require.Equal(t, "system_network_dropped_nonstandard:_1", normalizeName("system.network.dropped", "nonstandard:_1", MetricTypeGauge, ""))
}

func TestNonStandardUnitCounter(t *testing.T) {
	require.Equal(t, "system_network_dropped_total", normalizeName("system.network.dropped", "{packets}", MetricTypeMonotonicCounter, ""))
}

func TestBrokenUnit(t *testing.T) {
	require.Equal(t, "system_network_dropped_packets", normalizeName("system.network.dropped", "packets", MetricTypeGauge, ""))
	require.Equal(t, "system_network_packets_dropped", normalizeName("system.network.packets.dropped", "packets", MetricTypeGauge, ""))
	require.Equal(t, "system_network_packets", normalizeName("system.network.packets", "packets", MetricTypeGauge, ""))
}

func TestBrokenUnitCounter(t *testing.T) {
	require.Equal(t, "system_network_dropped_packets_total", normalizeName("system.network.dropped", "packets", MetricTypeMonotonicCounter, ""))
	require.Equal(t, "system_network_packets_dropped_total", normalizeName("system.network.packets.dropped", "packets", MetricTypeMonotonicCounter, ""))
	require.Equal(t, "system_network_packets_total", normalizeName("system.network.packets", "packets", MetricTypeMonotonicCounter, ""))
}

func TestRatio(t *testing.T) {
	require.Equal(t, "hw_gpu_memory_utilization_ratio", normalizeName("hw.gpu.memory.utilization", "1", MetricTypeGauge, ""))
	require.Equal(t, "hw_fan_speed_ratio", normalizeName("hw.fan.speed_ratio", "1", MetricTypeGauge, ""))
	require.Equal(t, "objects_total", normalizeName("objects", "1", MetricTypeMonotonicCounter, ""))
}

func TestHertz(t *testing.T) {
	require.Equal(t, "hw_cpu_speed_limit_hertz", normalizeName("hw.cpu.speed_limit", "Hz", MetricTypeGauge, ""))
}

func TestPer(t *testing.T) {
	require.Equal(t, "broken_metric_speed_km_per_hour", normalizeName("broken.metric.speed", "km/h", MetricTypeGauge, ""))
	require.Equal(t, "astro_light_speed_limit_meters_per_second", normalizeName("astro.light.speed_limit", "m/s", MetricTypeGauge, ""))
	// The normal metric name character set is allowed in non-standard units.
	require.Equal(t, "system_network_dropped_non_per_standard:_1", normalizeName("system.network.dropped", "non/standard:_1", MetricTypeGauge, ""))

	t.Run("invalid per unit", func(t *testing.T) {
		require.Equal(t, "broken_metric_speed_km", normalizeName("broken.metric.speed", "km/°", MetricTypeGauge, ""))
	})
}

func TestPercent(t *testing.T) {
	require.Equal(t, "broken_metric_success_ratio_percent", normalizeName("broken.metric.success_ratio", "%", MetricTypeGauge, ""))
	require.Equal(t, "broken_metric_success_percent", normalizeName("broken.metric.success_percent", "%", MetricTypeGauge, ""))
}

func TestEmpty(t *testing.T) {
	require.Equal(t, "test_metric_no_unit", normalizeName("test.metric.no_unit", "", MetricTypeGauge, ""))
	require.Equal(t, "test_metric_spaces", normalizeName("test.metric.spaces", "   \t  ", MetricTypeGauge, ""))
}

func TestOTelReceivers(t *testing.T) {
	require.Equal(t, "active_directory_ds_replication_network_io_bytes_total", normalizeName("active_directory.ds.replication.network.io", "By", MetricTypeMonotonicCounter, ""))
	require.Equal(t, "active_directory_ds_replication_sync_object_pending_total", normalizeName("active_directory.ds.replication.sync.object.pending", "{objects}", MetricTypeMonotonicCounter, ""))
	require.Equal(t, "active_directory_ds_replication_object_rate_per_second", normalizeName("active_directory.ds.replication.object.rate", "{objects}/s", MetricTypeGauge, ""))
	require.Equal(t, "active_directory_ds_name_cache_hit_rate_percent", normalizeName("active_directory.ds.name_cache.hit_rate", "%", MetricTypeGauge, ""))
	require.Equal(t, "active_directory_ds_ldap_bind_last_successful_time_milliseconds", normalizeName("active_directory.ds.ldap.bind.last_successful.time", "ms", MetricTypeGauge, ""))
	require.Equal(t, "apache_current_connections", normalizeName("apache.current_connections", "connections", MetricTypeGauge, ""))
	require.Equal(t, "apache_workers_connections", normalizeName("apache.workers", "connections", MetricTypeGauge, ""))
	require.Equal(t, "apache_requests_total", normalizeName("apache.requests", "1", MetricTypeMonotonicCounter, ""))
	require.Equal(t, "bigip_virtual_server_request_count_total", normalizeName("bigip.virtual_server.request.count", "{requests}", MetricTypeMonotonicCounter, ""))
	require.Equal(t, "system_cpu_utilization_ratio", normalizeName("system.cpu.utilization", "1", MetricTypeGauge, ""))
	require.Equal(t, "system_disk_operation_time_seconds_total", normalizeName("system.disk.operation_time", "s", MetricTypeMonotonicCounter, ""))
	require.Equal(t, "system_cpu_load_average_15m_ratio", normalizeName("system.cpu.load_average.15m", "1", MetricTypeGauge, ""))
	require.Equal(t, "memcached_operation_hit_ratio_percent", normalizeName("memcached.operation_hit_ratio", "%", MetricTypeGauge, ""))
	require.Equal(t, "mongodbatlas_process_asserts_per_second", normalizeName("mongodbatlas.process.asserts", "{assertions}/s", MetricTypeGauge, ""))
	require.Equal(t, "mongodbatlas_process_journaling_data_files_mebibytes", normalizeName("mongodbatlas.process.journaling.data_files", "MiBy", MetricTypeGauge, ""))
	require.Equal(t, "mongodbatlas_process_network_io_bytes_per_second", normalizeName("mongodbatlas.process.network.io", "By/s", MetricTypeGauge, ""))
	require.Equal(t, "mongodbatlas_process_oplog_rate_gibibytes_per_hour", normalizeName("mongodbatlas.process.oplog.rate", "GiBy/h", MetricTypeGauge, ""))
	require.Equal(t, "mongodbatlas_process_db_query_targeting_scanned_per_returned", normalizeName("mongodbatlas.process.db.query_targeting.scanned_per_returned", "{scanned}/{returned}", MetricTypeGauge, ""))
	require.Equal(t, "nginx_requests", normalizeName("nginx.requests", "requests", MetricTypeGauge, ""))
	require.Equal(t, "nginx_connections_accepted", normalizeName("nginx.connections_accepted", "connections", MetricTypeGauge, ""))
	require.Equal(t, "nsxt_node_memory_usage_kilobytes", normalizeName("nsxt.node.memory.usage", "KBy", MetricTypeGauge, ""))
	require.Equal(t, "redis_latest_fork_microseconds", normalizeName("redis.latest_fork", "us", MetricTypeGauge, ""))
}

func TestNamespace(t *testing.T) {
	require.Equal(t, "space_test", normalizeName("test", "", MetricTypeGauge, "space"))
	require.Equal(t, "space_test", normalizeName("#test", "", MetricTypeGauge, "space"))
}

func TestCleanUpUnit(t *testing.T) {
	require.Empty(t, cleanUpUnit(""))
	require.Equal(t, "a_b", cleanUpUnit("a b"))
	require.Equal(t, "hello_world", cleanUpUnit("hello, world"))
	require.Equal(t, "hello_you_2", cleanUpUnit("hello you 2"))
	require.Equal(t, "1000", cleanUpUnit("$1000"))
	require.Empty(t, cleanUpUnit("*+$^=)"))
}

func TestUnitMapGetOrDefault(t *testing.T) {
	require.Empty(t, unitMapGetOrDefault(""))
	require.Equal(t, "seconds", unitMapGetOrDefault("s"))
	require.Equal(t, "invalid", unitMapGetOrDefault("invalid"))
}

func TestPerUnitMapGetOrDefault(t *testing.T) {
	require.Empty(t, perUnitMapGetOrDefault(""))
	require.Equal(t, "second", perUnitMapGetOrDefault("s"))
	require.Equal(t, "invalid", perUnitMapGetOrDefault("invalid"))
}

func TestBuildUnitSuffixes(t *testing.T) {
	tests := []struct {
		unit         string
		expectedMain string
		expectedPer  string
	}{
		{"", "", ""},
		{"s", "seconds", ""},
		{"By/s", "bytes", "per_second"},
		{"requests/m", "requests", "per_minute"},
		{"{invalid}/second", "", "per_second"},
		{"bytes/{invalid}", "bytes", ""},
	}

	for _, test := range tests {
		mainUnitSuffix, perUnitSuffix := buildUnitSuffixes(test.unit)
		require.Equal(t, test.expectedMain, mainUnitSuffix)
		require.Equal(t, test.expectedPer, perUnitSuffix)
	}
}

func TestAddUnitTokens(t *testing.T) {
	tests := []struct {
		nameTokens     []string
		mainUnitSuffix string
		perUnitSuffix  string
		expected       []string
	}{
		{[]string{}, "", "", []string{}},
		{[]string{"token1"}, "main", "", []string{"token1", "main"}},
		{[]string{"token1"}, "", "per", []string{"token1", "per"}},
		{[]string{"token1"}, "main", "per", []string{"token1", "main", "per"}},
		{[]string{"token1", "per"}, "main", "per", []string{"token1", "per", "main"}},
		{[]string{"token1", "main"}, "main", "per", []string{"token1", "main", "per"}},
		{[]string{"token1"}, "main_", "per", []string{"token1", "main", "per"}},
		{[]string{"token1"}, "main_unit", "per_seconds_", []string{"token1", "main_unit", "per_seconds"}}, // trailing underscores are removed
		{[]string{"token1"}, "main_unit", "per_", []string{"token1", "main_unit"}},                        // 'per_' is removed entirely
	}

	for _, test := range tests {
		result := addUnitTokens(test.nameTokens, test.mainUnitSuffix, test.perUnitSuffix)
		require.Equal(t, test.expected, result)
	}
}

func TestRemoveItem(t *testing.T) {
	require.Equal(t, []string{}, removeItem([]string{}, "test"))
	require.Equal(t, []string{}, removeItem([]string{}, ""))
	require.Equal(t, []string{"a", "b", "c"}, removeItem([]string{"a", "b", "c"}, "d"))
	require.Equal(t, []string{"a", "b", "c"}, removeItem([]string{"a", "b", "c"}, ""))
	require.Equal(t, []string{"a", "b"}, removeItem([]string{"a", "b", "c"}, "c"))
	require.Equal(t, []string{"a", "c"}, removeItem([]string{"a", "b", "c"}, "b"))
	require.Equal(t, []string{"b", "c"}, removeItem([]string{"a", "b", "c"}, "a"))
}

func TestBuildCompliantMetricNameWithSuffixes(t *testing.T) {
	builder := NewMetricNameBuilder("", true)
	require.Equal(t, "system_io_bytes_total", builder.BuildCompliantMetricName("system.io", "By", MetricTypeMonotonicCounter))
	require.Equal(t, "_3_14_digits", builder.BuildCompliantMetricName("3.14 digits", "", MetricTypeGauge))
	require.Equal(t, "envoy_rule_engine_zlib_buf_error", builder.BuildCompliantMetricName("envoy__rule_engine_zlib_buf_error", "", MetricTypeGauge))
	require.Equal(t, ":foo::bar", builder.BuildCompliantMetricName(":foo::bar", "", MetricTypeGauge))
	require.Equal(t, ":foo::bar_total", builder.BuildCompliantMetricName(":foo::bar", "", MetricTypeMonotonicCounter))
	// Gauges with unit 1 are considered ratios.
	require.Equal(t, "foo_bar_ratio", builder.BuildCompliantMetricName("foo.bar", "1", MetricTypeGauge))
	// Slashes in units are converted.
	require.Equal(t, "system_io_foo_per_bar_total", builder.BuildCompliantMetricName("system.io", "foo/bar", MetricTypeMonotonicCounter))
	require.Equal(t, "metric_with_foreign_characters_total", builder.BuildCompliantMetricName("metric_with_字符_foreign_characters", "", MetricTypeMonotonicCounter))
	// Removes non aplhanumerical characters from units, but leaves colons.
	require.Equal(t, "temperature_:C", builder.BuildCompliantMetricName("temperature", "%*()°:C", MetricTypeGauge))
}

func TestBuildCompliantMetricNameWithoutSuffixes(t *testing.T) {
	builder := NewMetricNameBuilder("", false)
	require.Equal(t, "system_io", builder.BuildCompliantMetricName("system.io", "By", MetricTypeMonotonicCounter))
	require.Equal(t, "network_I_O", builder.BuildCompliantMetricName("network (I/O)", "By", MetricTypeMonotonicCounter))
	require.Equal(t, "_3_14_digits", builder.BuildCompliantMetricName("3.14 digits", "By", MetricTypeGauge))
	require.Equal(t, "envoy__rule_engine_zlib_buf_error", builder.BuildCompliantMetricName("envoy__rule_engine_zlib_buf_error", "", MetricTypeGauge))
	require.Equal(t, ":foo::bar", builder.BuildCompliantMetricName(":foo::bar", "", MetricTypeGauge))
	require.Equal(t, ":foo::bar", builder.BuildCompliantMetricName(":foo::bar", "", MetricTypeMonotonicCounter))
	require.Equal(t, "foo_bar", builder.BuildCompliantMetricName("foo.bar", "1", MetricTypeGauge))
	require.Equal(t, "system_io", builder.BuildCompliantMetricName("system.io", "foo/bar", MetricTypeMonotonicCounter))
	require.Equal(t, "metric_with___foreign_characters", builder.BuildCompliantMetricName("metric_with_字符_foreign_characters", "", MetricTypeMonotonicCounter))
}

func TestBuildCompliantMetricNameWithNamespace(t *testing.T) {
	builder := NewMetricNameBuilder("namespace", false)
	require.Equal(t, "namespace_system_io", builder.BuildCompliantMetricName("system.io", "", MetricTypeMonotonicCounter))
}

func TestBuildMetricNameWithSuffixes(t *testing.T) {
	builder := NewMetricNameBuilder("", true)
	require.Equal(t, "system.io_bytes_total", builder.BuildMetricName("system.io", "By", MetricTypeMonotonicCounter))
	require.Equal(t, "3.14 digits", builder.BuildMetricName("3.14 digits", "", MetricTypeGauge))
	require.Equal(t, "envoy__rule_engine_zlib_buf_error", builder.BuildMetricName("envoy__rule_engine_zlib_buf_error", "", MetricTypeGauge))
	require.Equal(t, ":foo::bar", builder.BuildMetricName(":foo::bar", "", MetricTypeGauge))
	require.Equal(t, ":foo::bar_total", builder.BuildMetricName(":foo::bar", "", MetricTypeMonotonicCounter))
	// Gauges with unit 1 are considered ratios.
	require.Equal(t, "foo.bar_ratio", builder.BuildMetricName("foo.bar", "1", MetricTypeGauge))
	// Slashes in units are converted.
	require.Equal(t, "system.io_foo_per_bar_total", builder.BuildMetricName("system.io", "foo/bar", MetricTypeMonotonicCounter))
	require.Equal(t, "metric_with_字符_foreign_characters_total", builder.BuildMetricName("metric_with_字符_foreign_characters", "", MetricTypeMonotonicCounter))
	require.Equal(t, "temperature_%*()°C", builder.BuildMetricName("temperature", "%*()°C", MetricTypeGauge)) // Keeps the all characters in unit
	// Tests below show weird interactions that users can have with the metric names.
	// With BuildMetricName we don't check if units/type suffixes are already present in the metric name, we always add them.
	require.Equal(t, "system_io_seconds_seconds", builder.BuildMetricName("system_io_seconds", "s", MetricTypeGauge))
	require.Equal(t, "system_io_total_total", builder.BuildMetricName("system_io_total", "", MetricTypeMonotonicCounter))
}

func TestBuildMetricNameWithoutSuffixes(t *testing.T) {
	builder := NewMetricNameBuilder("", false)
	require.Equal(t, "system.io", builder.BuildMetricName("system.io", "By", MetricTypeMonotonicCounter))
	require.Equal(t, "3.14 digits", builder.BuildMetricName("3.14 digits", "", MetricTypeGauge))
	require.Equal(t, "envoy__rule_engine_zlib_buf_error", builder.BuildMetricName("envoy__rule_engine_zlib_buf_error", "", MetricTypeGauge))
	require.Equal(t, ":foo::bar", builder.BuildMetricName(":foo::bar", "", MetricTypeGauge))
	require.Equal(t, ":foo::bar", builder.BuildMetricName(":foo::bar", "", MetricTypeMonotonicCounter))
	// Gauges with unit 1 are considered ratios.
	require.Equal(t, "foo.bar", builder.BuildMetricName("foo.bar", "1", MetricTypeGauge))
	require.Equal(t, "metric_with_字符_foreign_characters", builder.BuildMetricName("metric_with_字符_foreign_characters", "", MetricTypeMonotonicCounter))
	require.Equal(t, "system_io_seconds", builder.BuildMetricName("system_io_seconds", "s", MetricTypeGauge))
	require.Equal(t, "system_io_total", builder.BuildMetricName("system_io_total", "", MetricTypeMonotonicCounter))
}

func TestBuildMetricNameWithNamespace(t *testing.T) {
	builder := NewMetricNameBuilder("namespace", false)
	require.Equal(t, "namespace_system.io", builder.BuildMetricName("system.io", "", MetricTypeMonotonicCounter))
}

# OTLP Prometheus Translator

A Go library for converting [OpenTelemetry Protocol (OTLP)](https://opentelemetry.io/docs/specs/otlp/) metric and attribute names to [Prometheus](https://prometheus.io/)-compliant formats. This is an internal library for both Prometheus and Open Telemetry, without any stability guarantees for external usage.

Part of the [Prometheus](https://prometheus.io/) ecosystem, following the [OpenTelemetry to Prometheus compatibility specification](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/compatibility/prometheus_and_openmetrics.md).

## Features

- **Metric Name and Label Translation**: Convert OTLP metric names and attributes to Prometheus-compliant format
- **Unit Handling**: Translate OTLP units to Prometheus unit conventions
- **Type-Aware Suffixes**: Optionally append `_total`, `_ratio` based on metric type
- **Namespace Support**: Add configurable namespace prefixes
- **UTF-8 Support**: Choose between Prometheus legacy scheme compliant metric/label names (`[a-zA-Z0-9:_]`) or untranslated metric/label names
- **Translation Strategy Configuration**: Select a translation strategy with a standard set of strings.
- **Updated UCUM Unit Mappings**: Opt in to corrected UCUM unit suffixes (`TiBy`: `tebibytes`, `kBy`: `kilobytes`) by selecting the `UnderscoreEscapingWithUpdatedSuffixes` or `NoUTF8EscapingWithUpdatedSuffixes` translation strategy.

## Installation

```bash
go get github.com/prometheus/otlptranslator
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/prometheus/otlptranslator"
)

func main() {
    // Create a metric namer using traditional Prometheus name translation, with suffixes added and UTF-8 disallowed.
    // Use UnderscoreEscapingWithUpdatedSuffixes to opt into corrected UCUM unit mappings:
    // TiBy -> "tebibytes" instead of "tibibytes" and kBy -> "kilobytes".
    strategy := otlptranslator.UnderscoreEscapingWithUpdatedSuffixes
    namer := otlptranslator.NewMetricNamer("myapp", strategy)

    // Translate OTLP metric to Prometheus format
    metric := otlptranslator.Metric{
        Name: "http.server.request.duration",
        Unit: "s",
        Type: otlptranslator.MetricTypeHistogram,
    }
    metricName, err := namer.Build(metric)
    if err != nil {
        panic(err)
    }
    fmt.Println(metricName) // Output: myapp_http_server_request_duration_seconds

    // Translate label names
    labelNamer := otlptranslator.LabelNamer{UTF8Allowed: false}
    labelName, err := labelNamer.Build("http.method")
    if err != nil {
        panic(err)
    }
    fmt.Println(labelName) // Output: http_method
}
```

## Usage Examples

### Metric Name Translation

```go
namer := otlptranslator.MetricNamer{WithMetricSuffixes: true, UTF8Allowed: false}

// Counter gets _total suffix
counter := otlptranslator.Metric{
    Name: "requests.count", Unit: "1", Type: otlptranslator.MetricTypeMonotonicCounter,
}
counterName, _ := namer.Build(counter)
fmt.Println(counterName) // requests_count_total

// Gauge with unit conversion
gauge := otlptranslator.Metric{
    Name: "memory.usage", Unit: "By", Type: otlptranslator.MetricTypeGauge,
}
gaugeName, _ := namer.Build(gauge)
fmt.Println(gaugeName) // memory_usage_bytes

// Dimensionless gauge gets _ratio suffix
ratio := otlptranslator.Metric{
    Name: "cpu.utilization", Unit: "1", Type: otlptranslator.MetricTypeGauge,
}
ratioName, _ := namer.Build(ratio)
fmt.Println(ratioName) // cpu_utilization_ratio
```

### Label Translation

```go
labelNamer := otlptranslator.LabelNamer{UTF8Allowed: false}

labelNamer.Build("http.method")           // http_method, nil
labelNamer.Build("123invalid")            // key_123invalid, nil
labelNamer.Build("_private")              // _private, nil
labelNamer.Build("__reserved__")          // __reserved__, nil
labelNamer.Build("label@with$symbols")    // label_with_symbols, nil
```

### Unit Translation

```go
unitNamer := otlptranslator.UnitNamer{UTF8Allowed: false}

unitNamer.Build("s")           // seconds
unitNamer.Build("By")          // bytes
unitNamer.Build("requests/s")  // requests_per_second
unitNamer.Build("1")           // "" (dimensionless)
```

### Configuration Options

```go
// Prometheus-compliant mode - supports [a-zA-Z0-9:_]
compliantNamer := otlptranslator.MetricNamer{UTF8Allowed: false, WithMetricSuffixes: true}

// Transparent pass-through mode, aka "NoTranslation"
utf8Namer := otlptranslator.MetricNamer{UTF8Allowed: true, WithMetricSuffixes: false}
utf8Namer = otlptranslator.NewMetricNamer("", otlptranslator.NoTranslation)

// With namespace and suffixes
productionNamer := otlptranslator.MetricNamer{
    Namespace:          "myservice",
    WithMetricSuffixes: true,
    UTF8Allowed:        false,
}
```

## License

Licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

# otlptranslator

Library providing API to convert OTLP metric data to popular Prometheus metric formats:
* OTLP metric and attribute names to Prometheus metric and label names, optionally following the [Prometheus naming conventions](https://prometheus.io/docs/practices/naming/).
* OTLP metric to Prometheus Remote Write format.

# Problem statement

Throughout the years, several different libraries that translate OTLP metrics into Prometheus metric formats (e.g. OpenMetrics and Prometheus Remote Write) have been created, maintained and forked by different parties.

As new features are being introduced to Prometheus, e.g. UTF-8 support, Native Histograms and Native Histogram Custom Buckets, many of those old libraries start to lack behind and don't fulfill the expectations of users. To make things even worse, with different downstream projects adopting different translation libraries we're starting to see a fragmented ecosystem where data emitted from different projects start to look different.

# Vision

* Thrive to be the technical authority and go-to library for all projects that need to translate OpenTelemetry metric and attribute names to [Prometheus naming conventions](https://prometheus.io/docs/practices/naming/).
* Thrive to be the technical authority and go-to library for all projects that need to translate full OpenTelemetry metrics into the different versions of the Prometheus Remote-Write protocol.
* Focus on providing the most efficient library for said translation.
* Drive and evolve the [OpenTelemetry<->Prometheus compatibility specification](https://opentelemetry.io/docs/specs/otel/compatibility/prometheus_and_openmetrics/) with the experience we get from maintaining this library.

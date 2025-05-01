# otlptranslator

Library providing API to convert OTLP metric and attribute names to respectively Prometheus metric and label names. 

# Problem statement

Throughout the years, several different libraries that translate OTLP metrics into Prometheus metric formats (e.g. OpenMetrics) have been created and maintained by different parties.

As new feature are being introduced to Prometheus, e.g. UTF-8 support, many of those old libraries start to lack behind and don't fulfill the expectations of users. To make things even worse, with different downstream projects adopting different translation libraries we're starting to see a fragmented ecosystem where data emitted from different projects start to look different.

# Vision

* Thrive to be the technical authority and go-to library for all projects that need to translate OpenTelemetry metric and attribute names to [Prometheus naming conventions](https://prometheus.io/docs/practices/naming/).
* Focus on providing the most efficient library for said translation.

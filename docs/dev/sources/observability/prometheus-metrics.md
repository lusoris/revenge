# Prometheus Metric Types

> Auto-fetched from [https://prometheus.io/docs/concepts/metric_types/](https://prometheus.io/docs/concepts/metric_types/)
> Last Updated: 2026-01-28T21:47:22.123670+00:00

---

Show nav
Introduction
Overview
First steps
Comparison to alternatives
FAQ
Roadmap
Design documents
Media
Glossary
Long-term support
Concepts
Data model
Metric types
Jobs and instances
Prometheus Server
Getting started
Installation
Configuration
Configuration
Recording rules
Alerting rules
Template examples
Template reference
HTTP configuration for promtool
Unit testing for rules
HTTPS and authentication
Agent Mode
Querying
Basics
Operators
Functions
Examples
HTTP API
Remote Read API
Storage
Federation
HTTP SD
Management API
Command Line
prometheus
promtool
Migration
API stability
Feature flags
Visualization
Expression browser
Grafana
Perses
Console templates
Instrumenting
Client libraries
Writing client libraries
Pushing metrics
Exporters and integrations
Writing exporters
Exposition formats
UTF-8 escaping schemes
Content negotiation
Operating
Security
Integrations
Alertmanager
Alerting overview
Alertmanager
Configuration
High Availability
Alerts API
Notification template reference
Notification template examples
Management API
HTTPS and authentication
Best practices
Metric and label naming
Consoles and dashboards
Instrumentation
Histograms and summaries
Alerting
Recording rules
When to use the Pushgateway
Remote write tuning
Guides
Basic auth
Monitoring Docker container metrics using cAdvisor
Use file-based service discovery to discover scrape targets
Instrumenting a Go application
Understanding and using the multi-target exporter pattern
Monitoring Linux host metrics with the Node Exporter
OpenTelemetry
UTF-8 in Prometheus
Docker Swarm
Query log
TLS encryption
Tutorials
Getting started with Prometheus
Understanding metric types
Instrumenting HTTP server written in Go
Visualizing metrics using Grafana
Alerting based on metrics
Specifications
Native Histograms
OpenMetrics
1.0
Remote Write
1.0
2.0
Metric types
The Prometheus client libraries offer four core metric types. These are
currently only differentiated in the client libraries (to enable APIs tailored
to the usage of the specific types) and in the wire protocol. The Prometheus
server does not yet make use of the type information and flattens all data into
untyped time series. This may change in the future.
Counter
A
counter
is a cumulative metric that represents a single
monotonically
increasing counter
whose
value can only increase or be reset to zero on restart. For example, you can
use a counter to represent the number of requests served, tasks completed, or
errors.
Do not use a counter to expose a value that can decrease. For example, do not
use a counter for the number of currently running processes; instead use a gauge.
Client library usage documentation for counters:
Go
Java
Python
Ruby
.Net
Rust
Gauge
A
gauge
is a metric that represents a single numerical value that can
arbitrarily go up and down.
Gauges are typically used for measured values like temperatures or current
memory usage, but also "counts" that can go up and down, like the number of
concurrent requests.
Client library usage documentation for gauges:
Go
Java
Python
Ruby
.Net
Rust
Histogram
A
histogram
samples observations (usually things like request durations or
response sizes) and counts them in configurable buckets. It also provides a sum
of all observed values.
A histogram with a base metric name of
<basename>
exposes multiple time series
during a scrape:
cumulative counters for the observation buckets, exposed as
<basename>_bucket{le="<upper inclusive bound>"}
the
total sum
of all observed values, exposed as
<basename>_sum
the
count
of events that have been observed, exposed as
<basename>_count
(identical to
<basename>_bucket{le="+Inf"}
above)
Use the
histogram_quantile()
function
to calculate quantiles from histograms or even aggregations of histograms. A
histogram is also suitable to calculate an
Apdex score
. When operating on buckets,
remember that the histogram is
cumulative
. See
histograms and summaries
for details of histogram
usage and differences to
summaries
.
NOTE
:
Beginning with Prometheus v2.40, there is experimental support for native
histograms. A native histogram requires only one time series, which includes a
dynamic number of buckets in addition to the sum and count of
observations. Native histograms allow much higher resolution at a fraction of
the cost. Detailed documentation will follow once native histograms are closer
to becoming a stable feature.
NOTE
:
Beginning with Prometheus v3.0, the values of the
le
label of classic
histograms are normalized during ingestion to follow the format of
OpenMetrics Canonical Numbers
.
Client library usage documentation for histograms:
Go
Java
Python
Ruby
.Net
Rust
Summary
Similar to a
histogram
, a
summary
samples observations (usually things like
request durations and response sizes). While it also provides a total count of
observations and a sum of all observed values, it calculates configurable
quantiles over a sliding time window.
A summary with a base metric name of
<basename>
exposes multiple time series
during a scrape:
streaming
φ-quantiles
(0 ≤ φ ≤ 1) of observed events, exposed as
<basename>{quantile="<φ>"}
the
total sum
of all observed values, exposed as
<basename>_sum
the
count
of events that have been observed, exposed as
<basename>_count
See
histograms and summaries
for
detailed explanations of φ-quantiles, summary usage, and differences
to
histograms
.
NOTE
:
Beginning with Prometheus v3.0, the values of the
quantile
label are normalized during
ingestion to follow the format of
OpenMetrics Canonical Numbers
.
Client library usage documentation for summaries:
Go
Java
Python
Ruby
.Net
Previous
Data model
Edit
Edit this page
Next
Jobs and instances
On this page
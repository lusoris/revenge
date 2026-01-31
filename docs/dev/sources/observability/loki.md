# Grafana Loki

> Source: https://grafana.com/docs/loki/latest/
> Fetched: 2026-01-31T16:02:34.831612+00:00
> Content-Hash: fa9cdd4b84777109
> Type: html

---

[Documentation](/docs/) Grafana Loki

Open source 

# Grafana Loki

Grafana Loki is a set of open source components that can be composed into a fully featured logging stack. A small index and highly compressed chunks simplifies the operation and significantly lowers the cost of Loki.

* * *

## Overview

Unlike other logging systems, Loki is built around the idea of only indexing metadata about your logs’ labels (just like Prometheus labels). Log data itself is then compressed and stored in chunks in object stores such as Amazon Simple Storage Service (S3) or Google Cloud Storage (GCS), or even locally on the filesystem.

## Explore

[Learn about LokiLearn about the Loki architecture and components, the various deployment modes, and best practices for labels.](/docs/loki/latest/get-started/)[Set up LokiView instructions for how to configure and install Loki, migrate from previous deployments, and upgrade your Loki environment.](/docs/loki/latest/setup/)[Configure LokiView the Loki configuration reference and configuration examples.](/docs/loki/latest/configure/)[Send logs to LokiSelect one or more clients to use to send your logs to Loki.](/docs/loki/latest/send-data/)[Manage LokiLearn how to manage tenants, log ingestion, storage, queries, and more.](/docs/loki/latest/operations/)[Query with LogQLInspired by PromQL, LogQL is Grafana Lokiâs query language. LogQL uses labels and operators for filtering.](/docs/loki/latest/query/)

## Was this page helpful?

Yes No

[Suggest an edit in GitHub ](https://github.com/grafana/loki/edit/main/docs/sources/_index.md)[Create a GitHub issue ](https://github.com/grafana/loki/issues/new?title=Documentation%20feedback:%20/docs/sources/_index.md)[Email docs@grafana.com ](mailto:docs@grafana.com)[Help and support ](/help/)[Community](/community/)

## Related resources from Grafana Labs

Additional helpful documentation, links, and articles:

[VideoGetting started with logging and Grafana LokiSee a demo of the updated features in Loki, and how to create metrics from logs and alert on your logs with powerful Prometheus-style alerting rules.](/go/webinar/getting-started-with-logging-and-grafana-loki/)[VideoEssential Grafana Loki configuration settingsThis webinar focuses on Grafana Loki configuration including agents Promtail and Docker; the Loki server; and Loki storage for popular backends.](/go/webinar/logging-with-loki-essential-configuration-settings/)[VideoScaling and securing your logs with Grafana LokiThis webinar covers the challenges of scaling and securing logs, and how Grafana Cloud Logs powered by Grafana Loki can help, cost-effectively.](/go/webinar/scaling-and-securing-your-logs-with-grafana-loki/)

Select page language

Is this page helpful?

Yes No

On this page

Scroll for more
  *[↑]: Back to Top
  *[v]: View this template
  *[t]: Discuss this template
  *[e]: Edit this template

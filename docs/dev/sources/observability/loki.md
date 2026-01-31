# Grafana Loki

> Source: https://grafana.com/docs/loki/latest/
> Fetched: 2026-01-30T23:55:19.506529+00:00
> Content-Hash: af03fd14ca673705
> Type: html

---

Documentation

Grafana Loki

Open source

Grafana Loki

Grafana Loki is a set of open source components that can be composed into a fully featured logging stack. A small index and highly compressed chunks simplifies the operation and significantly lowers the cost of Loki.

Overview

Unlike other logging systems, Loki is built around the idea of only indexing metadata about your logs’ labels (just like Prometheus labels).
Log data itself is then compressed and stored in chunks in object stores such as Amazon Simple Storage Service (S3) or Google Cloud Storage (GCS), or even locally on the filesystem.

Explore

Learn about Loki

Learn about the Loki architecture and components, the various deployment modes, and best practices for labels.

Set up Loki

View instructions for how to configure and install Loki, migrate from previous deployments, and upgrade your Loki environment.

Configure Loki

View the Loki configuration reference and configuration examples.

Send logs to Loki

Select one or more clients to use to send your logs to Loki.

Manage Loki

Learn how to manage tenants, log ingestion, storage, queries, and more.

Query with LogQL

Inspired by PromQL, LogQL is Grafana Lokiâs query language. LogQL uses labels and operators for filtering.

Was this page helpful?

Yes

No

Suggest an edit in GitHub

Create a GitHub issue

Email docs@grafana.com

Help and support

Community

Related resources from Grafana Labs

Additional helpful documentation, links, and articles:

Video

Getting started with logging and Grafana Loki

See a demo of the updated features in Loki, and how to create metrics from logs and alert on your logs with powerful Prometheus-style alerting rules.

Video

Essential Grafana Loki configuration settings

This webinar focuses on Grafana Loki configuration including agents Promtail and Docker; the Loki server; and Loki storage for popular backends.

Video

Scaling and securing your logs with Grafana Loki

This webinar covers the challenges of scaling and securing logs, and how Grafana Cloud Logs powered by Grafana Loki can help, cost-effectively.

Select page language

Is this page helpful?

Yes

No

On this page

Scroll for more
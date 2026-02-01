## Table of Contents

- [Observability - Metrics, Tracing, and Logging](#observability-metrics-tracing-and-logging)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [See Also](#see-also)

---
sources:
- name: Prometheus Go Client
    url: ../sources/observability/prometheus.md
    note: Metrics instrumentation
- name: Prometheus Metric Types
    url: ../sources/observability/prometheus-metrics.md
    note: Counter, Gauge, Histogram, Summary
- name: Jaeger Go Client
    url: ../sources/observability/jaeger-go.md
    note: Distributed tracing client
- name: OpenTelemetry Go
    url: https://pkg.go.dev/go.opentelemetry.io/otel
    note: Tracing and metrics SDK
- name: Loki
    url: ../sources/observability/loki.md
    note: Log aggregation system
- name: Grafana
    url: ../sources/observability/grafana.md
    note: Visualization and dashboards
- name: slog-multi
    url: ../sources/observability/slog-multi.md
    note: Multi-handler slog setup
- name: Go slog
    url: ../sources/go/stdlib/slog.md
    note: Structured logging
design_refs:
- title: technical
    path: technical.md
- title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
- title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
- title: OFFLOADING
    path: technical/OFFLOADING.md
---

# Observability - Metrics, Tracing, and Logging

> Monitor, trace, and debug your Revenge instance with comprehensive observability

The Observability Stack provides complete visibility into your Revenge server. Prometheus collects metrics like request rates and error counts. Jaeger with OpenTelemetry traces requests across services for debugging. Loki aggregates logs from all components with structured search. Grafana dashboards visualize metrics and logs in real-time. Pre-built dashboards for common scenarios included. Alert rules notify you of issues before users notice.

---

## Contents

<!-- TOC will be auto-generated here by markdown-toc -->

---

## How It Works

<!-- User-friendly explanation -->

## Features
<!-- Feature list placeholder -->

## Configuration
<!-- User-friendly configuration guide -->

## Related Documentation
### See Also
<!-- Related wiki pages -->

---

**Need Help?** [Open an issue](https://github.com/revenge-project/revenge/issues) or [Join the discussion](https://github.com/revenge-project/revenge/discussions)

## Table of Contents

- [Observability - Metrics, Tracing, and Logging](#observability-metrics-tracing-and-logging)
  - [Status](#status)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Observability - Metrics, Tracing, and Logging


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: technical


> > Complete observability stack with metrics, distributed tracing, and structured logging

Three pillars of observability:
- **Metrics**: Prometheus for metrics collection and alerting
- **Tracing**: Jaeger + OpenTelemetry for distributed request tracing
- **Logging**: Loki for log aggregation, slog for structured logging
- **Dashboards**: Grafana for visualization and alerting
- **Instrumentation**: Automatic + manual instrumentation patterns

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete observability patterns |
| Sources | ✅ | All observability tools documented |
| Instructions | ✅ | Generated from design |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete




## Related Documentation
### Design Documents
- [technical](INDEX.md)
- [01_ARCHITECTURE](../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../architecture/02_DESIGN_PRINCIPLES.md)
- [OFFLOADING](../technical/OFFLOADING.md)

### External Sources
- [Prometheus Go Client](../../sources/observability/prometheus.md) - Metrics instrumentation
- [Prometheus Metric Types](../../sources/observability/prometheus-metrics.md) - Counter, Gauge, Histogram, Summary
- [Jaeger Go Client](../../sources/observability/jaeger-go.md) - Distributed tracing client
- [OpenTelemetry Go](https://pkg.go.dev/go.opentelemetry.io/otel) - Tracing and metrics SDK
- [Loki](../../sources/observability/loki.md) - Log aggregation system
- [Grafana](../../sources/observability/grafana.md) - Visualization and dashboards
- [slog-multi](../../sources/observability/slog-multi.md) - Multi-handler slog setup
- [Go slog](../../sources/go/stdlib/slog.md) - Structured logging


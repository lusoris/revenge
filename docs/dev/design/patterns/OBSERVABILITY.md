## Table of Contents

- [Observability Pattern](#observability-pattern)
  - [Status](#status)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Observability Pattern

<!-- DESIGN: patterns, README, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES -->


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: pattern


> > Metrics, tracing, and logging patterns with Prometheus, OpenTelemetry, and structured logging

Three pillars of observability:
- **Metrics**: Prometheus for RED metrics (Rate, Errors, Duration)
- **Tracing**: OpenTelemetry for distributed traces (Jaeger export)
- **Logging**: Structured logging with slog (tint for dev, zap for prod)

All integrated via OpenTelemetry SDK for unified observability.


---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | ✅ | - |
| Instructions | ✅ | Generated from design |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete


## Related Documentation
### Design Documents
- [patterns](INDEX.md)
- [TECH_STACK](../technical/TECH_STACK.md)
- [BEST_PRACTICES](../operations/BEST_PRACTICES.md)

### External Sources
- [Prometheus](https://prometheus.io/docs/introduction/overview/) - Metrics collection
- [Prometheus Metrics Types](../../sources/observability/prometheus-metrics.md) - Counter, Gauge, Histogram, Summary
- [OpenTelemetry Go](../../sources/observability/opentelemetry.md) - Tracing SDK
- [slog](../../sources/go/stdlib/slog.md) - Structured logging
- [tint](https://github.com/lmittmann/tint) - Colorized slog handler (dev)
- [zap](https://github.com/uber-go/zap) - High-performance JSON logs (prod)
- [Jaeger](../../sources/observability/jaeger.md) - Distributed tracing backend


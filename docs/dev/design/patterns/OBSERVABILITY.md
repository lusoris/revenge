## Table of Contents

- [Observability Pattern](#observability-pattern)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Components](#components)
  - [Implementation](#implementation)
    - [File Structure](#file-structure)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
    - [Config Keys](#config-keys)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)



---
sources:
  - name: Prometheus
    url: https://prometheus.io/docs/introduction/overview/
    note: Metrics collection
  - name: Prometheus Metrics Types
    url: https://prometheus.io/docs/concepts/metric_types/
    note: Counter, Gauge, Histogram, Summary
  - name: OpenTelemetry Go
    url: https://opentelemetry.io/docs/languages/go/
    note: Tracing SDK
  - name: slog
    url: https://pkg.go.dev/log/slog
    note: Structured logging
  - name: tint
    url: https://github.com/lmittmann/tint
    note: Colorized slog handler (dev)
  - name: zap
    url: https://github.com/uber-go/zap
    note: High-performance JSON logs (prod)
  - name: Jaeger
    url: https://www.jaegertracing.io/docs/
    note: Distributed tracing backend
design_refs:
  - title: patterns
    path: patterns/INDEX.md
  - title: TECH_STACK
    path: technical/TECH_STACK.md
  - title: BEST_PRACTICES
    path: operations/BEST_PRACTICES.md
---

# Observability Pattern


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
| Instructions | 🔴 | - |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---


## Architecture

<!-- Architecture diagram placeholder -->

### Components

<!-- Component description -->


## Implementation

### File Structure

<!-- File structure -->

### Key Interfaces

<!-- Interface definitions -->

### Dependencies

<!-- Dependency list -->





## Configuration
### Environment Variables

<!-- Environment variables -->

### Config Keys

<!-- Configuration keys -->




## Testing Strategy

### Unit Tests

<!-- Unit test strategy -->

### Integration Tests

<!-- Integration test strategy -->

### Test Coverage

Target: **80% minimum**







## Related Documentation
### Design Documents
- [patterns](patterns/INDEX.md)
- [TECH_STACK](technical/TECH_STACK.md)
- [BEST_PRACTICES](operations/BEST_PRACTICES.md)

### External Sources
- [Prometheus](https://prometheus.io/docs/introduction/overview/) - Metrics collection
- [Prometheus Metrics Types](https://prometheus.io/docs/concepts/metric_types/) - Counter, Gauge, Histogram, Summary
- [OpenTelemetry Go](https://opentelemetry.io/docs/languages/go/) - Tracing SDK
- [slog](https://pkg.go.dev/log/slog) - Structured logging
- [tint](https://github.com/lmittmann/tint) - Colorized slog handler (dev)
- [zap](https://github.com/uber-go/zap) - High-performance JSON logs (prod)
- [Jaeger](https://www.jaegertracing.io/docs/) - Distributed tracing backend


## Table of Contents

- [Advanced Offloading Architecture](#advanced-offloading-architecture)
  - [Status](#status)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Advanced Offloading Architecture


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: technical


> > Offload heavy operations to background workers and external services

Complete offloading strategy:
- **Background Jobs**: River queue for async tasks (transcoding, metadata enrichment)
- **Caching**: Dragonfly/Rueidis for session storage, rate limiting, API caching
- **Search**: Typesense for full-text search offloading
- **Metrics**: Prometheus for monitoring and alerting
- **Pattern**: Fast HTTP response, queue heavy work, notify on completion

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete offloading architecture |
| Sources | ✅ | All offloading tools documented |
| Instructions | ✅ | Generated from design |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---



















## Related Documentation
### Design Documents
- [technical](INDEX.md)
- [01_ARCHITECTURE](../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../architecture/02_DESIGN_PRINCIPLES.md)
- [METADATA_ENRICHMENT](../patterns/METADATA_ENRICHMENT.md)

### External Sources
- [Dragonfly Documentation](../../sources/infrastructure/dragonfly.md) - Auto-resolved from dragonfly
- [Uber fx](../../sources/tooling/fx.md) - Auto-resolved from fx
- [koanf](../../sources/tooling/koanf.md) - Auto-resolved from koanf
- [Prometheus Go Client](../../sources/observability/prometheus.md) - Auto-resolved from prometheus
- [Prometheus Metric Types](../../sources/observability/prometheus-metrics.md) - Auto-resolved from prometheus-metrics
- [River Job Queue](../../sources/tooling/river.md) - Auto-resolved from river
- [rueidis](../../sources/tooling/rueidis.md) - Auto-resolved from rueidis
- [rueidis GitHub README](../../sources/tooling/rueidis-guide.md) - Auto-resolved from rueidis-docs
- [Typesense API](../../sources/infrastructure/typesense.md) - Auto-resolved from typesense
- [Typesense Go Client](../../sources/infrastructure/typesense-go.md) - Auto-resolved from typesense-go


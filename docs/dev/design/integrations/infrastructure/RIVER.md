## Table of Contents

- [River](#river)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Integration Structure](#integration-structure)
    - [Data Flow](#data-flow)
    - [Provides](#provides)
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
  - name: Dragonfly Documentation
    url: https://www.dragonflydb.io/docs
    note: Auto-resolved from dragonfly
  - name: Uber fx
    url: https://pkg.go.dev/go.uber.org/fx
    note: Auto-resolved from fx
  - name: google/uuid
    url: https://pkg.go.dev/github.com/google/uuid
    note: Auto-resolved from google-uuid
  - name: pgx PostgreSQL Driver
    url: https://pkg.go.dev/github.com/jackc/pgx/v5
    note: Auto-resolved from pgx
  - name: pgxpool Connection Pool
    url: https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool
    note: Auto-resolved from pgxpool
  - name: PostgreSQL Arrays
    url: https://www.postgresql.org/docs/current/arrays.html
    note: Auto-resolved from postgresql-arrays
  - name: PostgreSQL JSON Functions
    url: https://www.postgresql.org/docs/current/functions-json.html
    note: Auto-resolved from postgresql-json
  - name: Prometheus Go Client
    url: https://pkg.go.dev/github.com/prometheus/client_golang/prometheus
    note: Auto-resolved from prometheus
  - name: Prometheus Metric Types
    url: https://prometheus.io/docs/concepts/metric_types/
    note: Auto-resolved from prometheus-metrics
  - name: River Job Queue
    url: https://pkg.go.dev/github.com/riverqueue/river
    note: Auto-resolved from river
  - name: River Documentation
    url: https://riverqueue.com/docs
    note: Auto-resolved from river-docs
  - name: rueidis
    url: https://pkg.go.dev/github.com/redis/rueidis
    note: Auto-resolved from rueidis
  - name: rueidis GitHub README
    url: https://github.com/redis/rueidis
    note: Auto-resolved from rueidis-docs
  - name: Typesense API
    url: https://typesense.org/docs/latest/api/
    note: Auto-resolved from typesense
  - name: Typesense Go Client
    url: https://github.com/typesense/typesense-go
    note: Auto-resolved from typesense-go
design_refs:
  - title: integrations/infrastructure
    path: integrations/infrastructure.md
  - title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# River


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration


> Integration with River

> PostgreSQL-native job queue

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | ✅ | - |
| Instructions | ✅ | - |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---


## Architecture

### Integration Structure

```
internal/integration/river/
├── client.go              # API client
├── types.go               # Response types
├── mapper.go              # Map external → internal types
├── cache.go               # Response caching
└── client_test.go         # Tests
```

### Data Flow

<!-- Data flow diagram -->

### Provides

This integration provides:
<!-- Data provided by integration -->


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
- [integrations/infrastructure](integrations/infrastructure.md)
- [01_ARCHITECTURE](architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Dragonfly Documentation](https://www.dragonflydb.io/docs) - Auto-resolved from dragonfly
- [Uber fx](https://pkg.go.dev/go.uber.org/fx) - Auto-resolved from fx
- [google/uuid](https://pkg.go.dev/github.com/google/uuid) - Auto-resolved from google-uuid
- [pgx PostgreSQL Driver](https://pkg.go.dev/github.com/jackc/pgx/v5) - Auto-resolved from pgx
- [pgxpool Connection Pool](https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool) - Auto-resolved from pgxpool
- [PostgreSQL Arrays](https://www.postgresql.org/docs/current/arrays.html) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](https://www.postgresql.org/docs/current/functions-json.html) - Auto-resolved from postgresql-json
- [Prometheus Go Client](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus) - Auto-resolved from prometheus
- [Prometheus Metric Types](https://prometheus.io/docs/concepts/metric_types/) - Auto-resolved from prometheus-metrics
- [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) - Auto-resolved from river
- [River Documentation](https://riverqueue.com/docs) - Auto-resolved from river-docs
- [rueidis](https://pkg.go.dev/github.com/redis/rueidis) - Auto-resolved from rueidis
- [rueidis GitHub README](https://github.com/redis/rueidis) - Auto-resolved from rueidis-docs
- [Typesense API](https://typesense.org/docs/latest/api/) - Auto-resolved from typesense
- [Typesense Go Client](https://github.com/typesense/typesense-go) - Auto-resolved from typesense-go


## Table of Contents

- [PostgreSQL](#postgresql)
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
    url: ../../sources/infrastructure/dragonfly.md
    note: Auto-resolved from dragonfly
- name: pgx PostgreSQL Driver
    url: ../../sources/database/pgx.md
    note: Auto-resolved from pgx
- name: PostgreSQL Arrays
    url: ../../sources/database/postgresql-arrays.md
    note: Auto-resolved from postgresql-arrays
- name: PostgreSQL JSON Functions
    url: ../../sources/database/postgresql-json.md
    note: Auto-resolved from postgresql-json
- name: Prometheus Go Client
    url: ../../sources/observability/prometheus.md
    note: Auto-resolved from prometheus
- name: Prometheus Metric Types
    url: ../../sources/observability/prometheus-metrics.md
    note: Auto-resolved from prometheus-metrics
- name: River Job Queue
    url: ../../sources/tooling/river.md
    note: Auto-resolved from river
- name: River Documentation
    url: ../../sources/tooling/river-guide.md
    note: Auto-resolved from river-docs
- name: sqlc
    url: ../../sources/database/sqlc.md
    note: Auto-resolved from sqlc
- name: sqlc Configuration
    url: ../../sources/database/sqlc-config.md
    note: Auto-resolved from sqlc-config
- name: Typesense API
    url: ../../sources/infrastructure/typesense.md
    note: Auto-resolved from typesense
- name: Typesense Go Client
    url: ../../sources/infrastructure/typesense-go.md
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

# PostgreSQL

<!-- DESIGN: integrations/infrastructure, README, SCAFFOLD_TEMPLATE, test_output_claude -->

**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration

> Integration with PostgreSQL

> Primary database for all Revenge data

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
internal/integration/postgresql/
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
- [01_ARCHITECTURE](../../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Dragonfly Documentation](../sources/infrastructure/dragonfly.md) - Auto-resolved from dragonfly
- [pgx PostgreSQL Driver](../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [Prometheus Go Client](../sources/observability/prometheus.md) - Auto-resolved from prometheus
- [Prometheus Metric Types](../sources/observability/prometheus-metrics.md) - Auto-resolved from prometheus-metrics
- [River Job Queue](../sources/tooling/river.md) - Auto-resolved from river
- [River Documentation](../sources/tooling/river-guide.md) - Auto-resolved from river-docs
- [sqlc](../sources/database/sqlc.md) - Auto-resolved from sqlc
- [sqlc Configuration](../sources/database/sqlc-config.md) - Auto-resolved from sqlc-config
- [Typesense API](../sources/infrastructure/typesense.md) - Auto-resolved from typesense
- [Typesense Go Client](../sources/infrastructure/typesense-go.md) - Auto-resolved from typesense-go

## Table of Contents

- [Search Service](#search-service)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Service Structure](#service-structure)
    - [Dependencies](#dependencies)
    - [Provides](#provides)
    - [Component Diagram](#component-diagram)
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
  - name: Uber fx
    url: https://pkg.go.dev/go.uber.org/fx
    note: Auto-resolved from fx
  - name: River Job Queue
    url: https://pkg.go.dev/github.com/riverqueue/river
    note: Auto-resolved from river
  - name: Typesense API
    url: https://typesense.org/docs/latest/api/
    note: Auto-resolved from typesense
  - name: Typesense Go Client
    url: https://github.com/typesense/typesense-go
    note: Auto-resolved from typesense-go
design_refs:
  - title: services
    path: services/INDEX.md
  - title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# Search Service


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: service


> > Full-text search via Typesense with per-module collections

**Package**: `internal/service/search`
**fx Module**: `search.Module`

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

### Service Structure

```
internal/service/search/
├── module.go              # fx module definition
├── service.go             # Service implementation
├── repository.go          # Data access (if needed)
├── handler.go             # HTTP handlers (if exposed)
├── middleware.go          # Middleware (if needed)
├── types.go               # Domain types
└── service_test.go        # Tests
```

### Dependencies
No external service dependencies.

### Provides
<!-- Service provides -->

### Component Diagram

<!-- Component diagram -->


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
- [services](services/INDEX.md)
- [01_ARCHITECTURE](architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Uber fx](https://pkg.go.dev/go.uber.org/fx) - Auto-resolved from fx
- [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) - Auto-resolved from river
- [Typesense API](https://typesense.org/docs/latest/api/) - Auto-resolved from typesense
- [Typesense Go Client](https://github.com/typesense/typesense-go) - Auto-resolved from typesense-go


## Table of Contents

- [Lidarr](#lidarr)
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
  - name: Uber fx
    url: https://pkg.go.dev/go.uber.org/fx
    note: Auto-resolved from fx
  - name: Last.fm API
    url: https://www.last.fm/api/intro
    note: Auto-resolved from lastfm-api
  - name: Lidarr API Docs
    url: https://lidarr.audio/docs/api/
    note: Auto-resolved from lidarr-docs
  - name: pgx PostgreSQL Driver
    url: https://pkg.go.dev/github.com/jackc/pgx/v5
    note: Auto-resolved from pgx
  - name: PostgreSQL Arrays
    url: https://www.postgresql.org/docs/current/arrays.html
    note: Auto-resolved from postgresql-arrays
  - name: PostgreSQL JSON Functions
    url: https://www.postgresql.org/docs/current/functions-json.html
    note: Auto-resolved from postgresql-json
  - name: River Job Queue
    url: https://pkg.go.dev/github.com/riverqueue/river
    note: Auto-resolved from river
  - name: Servarr Wiki
    url: https://wiki.servarr.com/
    note: Auto-resolved from servarr-wiki
  - name: Typesense API
    url: https://typesense.org/docs/latest/api/
    note: Auto-resolved from typesense
  - name: Typesense Go Client
    url: https://github.com/typesense/typesense-go
    note: Auto-resolved from typesense-go
design_refs:
  - title: integrations/servarr
    path: integrations/servarr.md
  - title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# Lidarr


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration


> Integration with Lidarr

> Music management automation
**Authentication**: api_key

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | ✅ | - |
| Instructions | 🟡 | - |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---


## Architecture

### Integration Structure

```
internal/integration/lidarr/
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
- [integrations/servarr](integrations/servarr.md)
- [01_ARCHITECTURE](architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Uber fx](https://pkg.go.dev/go.uber.org/fx) - Auto-resolved from fx
- [Last.fm API](https://www.last.fm/api/intro) - Auto-resolved from lastfm-api
- [Lidarr API Docs](https://lidarr.audio/docs/api/) - Auto-resolved from lidarr-docs
- [pgx PostgreSQL Driver](https://pkg.go.dev/github.com/jackc/pgx/v5) - Auto-resolved from pgx
- [PostgreSQL Arrays](https://www.postgresql.org/docs/current/arrays.html) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](https://www.postgresql.org/docs/current/functions-json.html) - Auto-resolved from postgresql-json
- [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) - Auto-resolved from river
- [Servarr Wiki](https://wiki.servarr.com/) - Auto-resolved from servarr-wiki
- [Typesense API](https://typesense.org/docs/latest/api/) - Auto-resolved from typesense
- [Typesense Go Client](https://github.com/typesense/typesense-go) - Auto-resolved from typesense-go


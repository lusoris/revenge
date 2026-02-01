---
sources:
  - name: Khan/genqlient
    url: ../../../../sources/tooling/genqlient.md
    note: Auto-resolved from genqlient
  - name: genqlient GitHub README
    url: ../../../../sources/tooling/genqlient-guide.md
    note: Auto-resolved from genqlient-docs
  - name: gohlslib (HLS)
    url: ../../../../sources/media/gohlslib.md
    note: Auto-resolved from gohlslib
  - name: pgx PostgreSQL Driver
    url: ../../../../sources/database/pgx.md
    note: Auto-resolved from pgx
  - name: PostgreSQL Arrays
    url: ../../../../sources/database/postgresql-arrays.md
    note: Auto-resolved from postgresql-arrays
  - name: PostgreSQL JSON Functions
    url: ../../../../sources/database/postgresql-json.md
    note: Auto-resolved from postgresql-json
  - name: River Job Queue
    url: ../../../../sources/tooling/river.md
    note: Auto-resolved from river
  - name: Typesense API
    url: ../../../../sources/infrastructure/typesense.md
    note: Auto-resolved from typesense
  - name: Typesense Go Client
    url: ../../../../sources/infrastructure/typesense-go.md
    note: Auto-resolved from typesense-go
design_refs:
  - title: ADULT_CONTENT_SYSTEM (QAR module)
    path: ../../../features/adult/ADULT_CONTENT_SYSTEM.md
  - title: STASHDB (community metadata)
    path: ./STASHDB.md
  - title: WHISPARR (PRIMARY for QAR)
    path: ../../servarr/WHISPARR.md
---

## Table of Contents

- [Stash](#stash)
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


# Stash


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration


> Integration with Stash

> Migration/sync tool for self-hosted Stash libraries
**Authentication**: api_key

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
internal/integration/stash/
├── client.go              # API client
├── types.go               # Response types
├── mapper.go              # Map external → internal types
├── cache.go               # Response caching
└── client_test.go         # Tests
```

### Data Flow

<!-- Data flow diagram -->

### Provides
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
- [ADULT_CONTENT_SYSTEM (QAR module)](../../../features/adult/ADULT_CONTENT_SYSTEM.md)
- [STASHDB (community metadata)](./STASHDB.md)
- [WHISPARR (PRIMARY for QAR)](../../servarr/WHISPARR.md)

### External Sources
- [Khan/genqlient](../../../../sources/tooling/genqlient.md) - Auto-resolved from genqlient
- [genqlient GitHub README](../../../../sources/tooling/genqlient-guide.md) - Auto-resolved from genqlient-docs
- [gohlslib (HLS)](../../../../sources/media/gohlslib.md) - Auto-resolved from gohlslib
- [pgx PostgreSQL Driver](../../../../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../../../../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../../../../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [River Job Queue](../../../../sources/tooling/river.md) - Auto-resolved from river
- [Typesense API](../../../../sources/infrastructure/typesense.md) - Auto-resolved from typesense
- [Typesense Go Client](../../../../sources/infrastructure/typesense-go.md) - Auto-resolved from typesense-go


## Table of Contents

- [TV Show Module](#tv-show-module)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Database Schema](#database-schema)
    - [Module Structure](#module-structure)
    - [Component Interaction](#component-interaction)
  - [Implementation](#implementation)
    - [File Structure](#file-structure)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
    - [Config Keys](#config-keys)
  - [API Endpoints](#api-endpoints)
    - [Content Management](#content-management)
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
    url: ../sources/tooling/fx.md
    note: Auto-resolved from fx
  - name: go-blurhash
    url: ../sources/media/go-blurhash.md
    note: Auto-resolved from go-blurhash
  - name: ogen OpenAPI Generator
    url: ../sources/tooling/ogen.md
    note: Auto-resolved from ogen
  - name: pgx PostgreSQL Driver
    url: ../sources/database/pgx.md
    note: Auto-resolved from pgx
  - name: PostgreSQL Arrays
    url: ../sources/database/postgresql-arrays.md
    note: Auto-resolved from postgresql-arrays
  - name: PostgreSQL JSON Functions
    url: ../sources/database/postgresql-json.md
    note: Auto-resolved from postgresql-json
  - name: River Job Queue
    url: ../sources/tooling/river.md
    note: Auto-resolved from river
  - name: Sonarr API Docs
    url: ../sources/apis/sonarr-docs.md
    note: Auto-resolved from sonarr-docs
  - name: sqlc
    url: ../sources/database/sqlc.md
    note: Auto-resolved from sqlc
  - name: sqlc Configuration
    url: ../sources/database/sqlc-config.md
    note: Auto-resolved from sqlc-config
  - name: TheTVDB API
    url: ../sources/apis/thetvdb.md
    note: Auto-resolved from thetvdb
design_refs:
  - title: features/video
    path: features/video.md
  - title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# TV Show Module


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: feature


> Content module for TV Shows, Seasons, Episodes

> TV series, seasons, and episodes management

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

### Database Schema

**Schema**: `public`

<!-- Schema diagram -->

### Module Structure

```
internal/content/tv_show/
├── module.go              # fx module definition
├── repository.go          # Database operations
├── service.go             # Business logic
├── handler.go             # HTTP handlers (ogen)
├── types.go               # Domain types
└── tv_show_test.go
```

### Component Interaction

<!-- Component interaction diagram -->


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


## API Endpoints

### Content Management
<!-- API endpoints placeholder -->


## Testing Strategy

### Unit Tests

<!-- Unit test strategy -->

### Integration Tests

<!-- Integration test strategy -->

### Test Coverage

Target: **80% minimum**







## Related Documentation
### Design Documents
- [features/video](features/video.md)
- [01_ARCHITECTURE](architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Uber fx](../sources/tooling/fx.md) - Auto-resolved from fx
- [go-blurhash](../sources/media/go-blurhash.md) - Auto-resolved from go-blurhash
- [ogen OpenAPI Generator](../sources/tooling/ogen.md) - Auto-resolved from ogen
- [pgx PostgreSQL Driver](../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [River Job Queue](../sources/tooling/river.md) - Auto-resolved from river
- [Sonarr API Docs](../sources/apis/sonarr-docs.md) - Auto-resolved from sonarr-docs
- [sqlc](../sources/database/sqlc.md) - Auto-resolved from sqlc
- [sqlc Configuration](../sources/database/sqlc-config.md) - Auto-resolved from sqlc-config
- [TheTVDB API](../sources/apis/thetvdb.md) - Auto-resolved from thetvdb


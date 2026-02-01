## Table of Contents

- [TMDb (The Movie Database)](#tmdb-the-movie-database)
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
  - name: go-blurhash
    url: ../sources/media/go-blurhash.md
    note: Auto-resolved from go-blurhash
  - name: golang.org/x/time
    url: ../sources/go/x/time.md
    note: Auto-resolved from golang-x-time
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
  - name: Typesense API
    url: ../sources/infrastructure/typesense.md
    note: Auto-resolved from typesense
  - name: Typesense Go Client
    url: ../sources/infrastructure/typesense-go.md
    note: Auto-resolved from typesense-go
design_refs:
  - title: integrations/metadata/video
    path: integrations/metadata/video.md
  - title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# TMDb (The Movie Database)

<!-- DESIGN: integrations/metadata/video, README, SCAFFOLD_TEMPLATE, test_output_claude -->


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration


> Integration with TMDb (The Movie Database)

> Primary metadata provider for movies and TV shows
**API Base URL**: `https://api.themoviedb.org/3`
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
internal/integration/tmdb_the_movie_database/
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
- [integrations/metadata/video](integrations/metadata/video.md)
- [01_ARCHITECTURE](../../../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../../../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../../../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [go-blurhash](../sources/media/go-blurhash.md) - Auto-resolved from go-blurhash
- [golang.org/x/time](../sources/go/x/time.md) - Auto-resolved from golang-x-time
- [pgx PostgreSQL Driver](../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [River Job Queue](../sources/tooling/river.md) - Auto-resolved from river
- [Typesense API](../sources/infrastructure/typesense.md) - Auto-resolved from typesense
- [Typesense Go Client](../sources/infrastructure/typesense-go.md) - Auto-resolved from typesense-go


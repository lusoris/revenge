## Table of Contents

- [StashDB](#stashdb)
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
- name: FFmpeg Documentation
    url: ../sources/media/ffmpeg.md
    note: Auto-resolved from ffmpeg
- name: FFmpeg Codecs
    url: ../sources/media/ffmpeg-codecs.md
    note: Auto-resolved from ffmpeg-codecs
- name: FFmpeg Formats
    url: ../sources/media/ffmpeg-formats.md
    note: Auto-resolved from ffmpeg-formats
- name: Khan/genqlient
    url: ../sources/tooling/genqlient.md
    note: Auto-resolved from genqlient
- name: genqlient GitHub README
    url: ../sources/tooling/genqlient-guide.md
    note: Auto-resolved from genqlient-docs
- name: go-astiav (FFmpeg bindings)
    url: ../sources/media/go-astiav.md
    note: Auto-resolved from go-astiav
- name: go-astiav GitHub README
    url: ../sources/media/go-astiav-guide.md
    note: Auto-resolved from go-astiav-docs
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
- name: StashDB GraphQL API
    url: ../sources/apis/stashdb-schema.graphql
    note: Auto-resolved from stashdb
- name: Typesense API
    url: ../sources/infrastructure/typesense.md
    note: Auto-resolved from typesense
- name: Typesense Go Client
    url: ../sources/infrastructure/typesense-go.md
    note: Auto-resolved from typesense-go
design_refs:
- title: integrations/metadata/adult
    path: integrations/metadata/adult.md
- title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
- title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
- title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# StashDB

<!-- DESIGN: integrations/metadata/adult, README, SCAFFOLD_TEMPLATE, test_output_claude -->

**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration

> Integration with StashDB

> Adult metadata database for performers, studios, and scenes
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
internal/integration/stashdb/
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
- [integrations/metadata/adult](integrations/metadata/adult.md)
- [01_ARCHITECTURE](../../../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../../../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../../../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [FFmpeg Documentation](../sources/media/ffmpeg.md) - Auto-resolved from ffmpeg
- [FFmpeg Codecs](../sources/media/ffmpeg-codecs.md) - Auto-resolved from ffmpeg-codecs
- [FFmpeg Formats](../sources/media/ffmpeg-formats.md) - Auto-resolved from ffmpeg-formats
- [Khan/genqlient](../sources/tooling/genqlient.md) - Auto-resolved from genqlient
- [genqlient GitHub README](../sources/tooling/genqlient-guide.md) - Auto-resolved from genqlient-docs
- [go-astiav (FFmpeg bindings)](../sources/media/go-astiav.md) - Auto-resolved from go-astiav
- [go-astiav GitHub README](../sources/media/go-astiav-guide.md) - Auto-resolved from go-astiav-docs
- [pgx PostgreSQL Driver](../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [River Job Queue](../sources/tooling/river.md) - Auto-resolved from river
- [StashDB GraphQL API](../sources/apis/stashdb-schema.graphql) - Auto-resolved from stashdb
- [Typesense API](../sources/infrastructure/typesense.md) - Auto-resolved from typesense
- [Typesense Go Client](../sources/infrastructure/typesense-go.md) - Auto-resolved from typesense-go

## Table of Contents

- [AniList](#anilist)
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
- name: AniList GraphQL API
    url: ../../sources/apis/anilist.md
    note: Auto-resolved from anilist
- name: AniList GraphQL Schema
    url: ../../sources/apis/anilist-schema.graphql
    note: Auto-resolved from anilist-graphql
- name: go-blurhash
    url: ../../sources/media/go-blurhash.md
    note: Auto-resolved from go-blurhash
- name: River Job Queue
    url: ../../sources/tooling/river.md
    note: Auto-resolved from river
design_refs:
- title: integrations/anime
    path: integrations/anime.md
- title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
- title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
- title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# AniList

<!-- DESIGN: integrations/anime, README, SCAFFOLD_TEMPLATE, test_output_claude -->

**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration

> Integration with AniList

> Primary metadata and tracking provider for anime and manga
**Authentication**: oauth

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
internal/integration/anilist/
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
- [integrations/anime](integrations/anime.md)
- [01_ARCHITECTURE](../../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [AniList GraphQL API](../sources/apis/anilist.md) - Auto-resolved from anilist
- [AniList GraphQL Schema](../sources/apis/anilist-schema.graphql) - Auto-resolved from anilist-graphql
- [go-blurhash](../sources/media/go-blurhash.md) - Auto-resolved from go-blurhash
- [River Job Queue](../sources/tooling/river.md) - Auto-resolved from river

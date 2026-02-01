## Table of Contents

- [Discogs](#discogs)
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
  - name: Discogs API
    url: ../sources/apis/discogs.md
    note: Auto-resolved from discogs
  - name: Last.fm API
    url: ../sources/apis/lastfm.md
    note: Auto-resolved from lastfm-api
design_refs:
  - title: integrations/metadata/music
    path: integrations/metadata/music.md
  - title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# Discogs

<!-- DESIGN: integrations/metadata/music, README, SCAFFOLD_TEMPLATE, test_output_claude -->


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration


> Integration with Discogs

> Music marketplace and database - vinyl releases, marketplace data, detailed credits
**API Base URL**: `https://api.discogs.com`
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
internal/integration/discogs/
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
- [integrations/metadata/music](integrations/metadata/music.md)
- [01_ARCHITECTURE](../../../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../../../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../../../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Discogs API](../sources/apis/discogs.md) - Auto-resolved from discogs
- [Last.fm API](../sources/apis/lastfm.md) - Auto-resolved from lastfm-api


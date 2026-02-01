## Table of Contents

- [Simkl](#simkl)
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
  - name: Last.fm API
    url: ../sources/apis/lastfm.md
    note: Auto-resolved from lastfm-api
  - name: River Job Queue
    url: ../sources/tooling/river.md
    note: Auto-resolved from river
  - name: Simkl API
    url: ../sources/apis/simkl.md
    note: Auto-resolved from simkl
design_refs:
  - title: integrations/scrobbling
    path: integrations/scrobbling.md
  - title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# Simkl


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration


> Integration with Simkl

> TV tracker and movie scrobbler (alternative to Trakt)
**API Base URL**: `https://api.simkl.com/oauth/token`
**Authentication**: oauth

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
internal/integration/simkl/
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
- [integrations/scrobbling](integrations/scrobbling.md)
- [01_ARCHITECTURE](architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Last.fm API](../sources/apis/lastfm.md) - Auto-resolved from lastfm-api
- [River Job Queue](../sources/tooling/river.md) - Auto-resolved from river
- [Simkl API](../sources/apis/simkl.md) - Auto-resolved from simkl


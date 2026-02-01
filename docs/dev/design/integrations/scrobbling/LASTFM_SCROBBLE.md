## Table of Contents

- [Last.fm Scrobbling](#lastfm-scrobbling)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Integration Structure](#integration-structure)
    - [Data Flow](#data-flow)
    - [Provides](#provides)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Last.fm Scrobbling


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration


> Integration with Last.fm Scrobbling

> Music scrobbling and listening history tracking
**API Base URL**: `https://www.last.fm/api/account/create`
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




## Architecture

### Integration Structure

```
internal/integration/lastfm_scrobbling/
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
## Related Documentation
### Design Documents
- [01_ARCHITECTURE](../../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Last.fm API](../../../sources/apis/lastfm.md) - Auto-resolved from lastfm-api
- [River Job Queue](../../../sources/tooling/river.md) - Auto-resolved from river


## Table of Contents

- [ErsatzTV](#ersatztv)
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
  - name: ErsatzTV Documentation
    url: https://ersatztv.org/docs/
    note: Auto-resolved from ersatztv-docs
  - name: gohlslib (HLS)
    url: https://pkg.go.dev/github.com/bluenviron/gohlslib/v2
    note: Auto-resolved from gohlslib
  - name: M3U8 Extended Format
    url: https://datatracker.ietf.org/doc/html/rfc8216
    note: Auto-resolved from m3u8
  - name: River Job Queue
    url: https://pkg.go.dev/github.com/riverqueue/river
    note: Auto-resolved from river
  - name: XMLTV Format
    url: https://github.com/XMLTV/xmltv/blob/master/xmltv.dtd
    note: Auto-resolved from xmltv
design_refs:
  - title: integrations/livetv
    path: integrations/livetv.md
  - title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# ErsatzTV


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration


> Integration with ErsatzTV

> Custom IPTV channel creation from your media library

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
internal/integration/ersatztv/
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
- [integrations/livetv](integrations/livetv.md)
- [01_ARCHITECTURE](architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](architecture/03_METADATA_SYSTEM.md)

### External Sources
- [ErsatzTV Documentation](https://ersatztv.org/docs/) - Auto-resolved from ersatztv-docs
- [gohlslib (HLS)](https://pkg.go.dev/github.com/bluenviron/gohlslib/v2) - Auto-resolved from gohlslib
- [M3U8 Extended Format](https://datatracker.ietf.org/doc/html/rfc8216) - Auto-resolved from m3u8
- [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) - Auto-resolved from river
- [XMLTV Format](https://github.com/XMLTV/xmltv/blob/master/xmltv.dtd) - Auto-resolved from xmltv


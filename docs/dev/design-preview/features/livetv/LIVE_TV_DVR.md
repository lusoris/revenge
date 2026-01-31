## Table of Contents

- [Live TV & DVR](#live-tv-dvr)
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
  - name: FFmpeg Documentation
    url: https://ffmpeg.org/ffmpeg.html
    note: Auto-resolved from ffmpeg
  - name: FFmpeg Codecs
    url: https://ffmpeg.org/ffmpeg-codecs.html
    note: Auto-resolved from ffmpeg-codecs
  - name: FFmpeg Formats
    url: https://ffmpeg.org/ffmpeg-formats.html
    note: Auto-resolved from ffmpeg-formats
  - name: Uber fx
    url: https://pkg.go.dev/go.uber.org/fx
    note: Auto-resolved from fx
  - name: go-astiav (FFmpeg bindings)
    url: https://pkg.go.dev/github.com/asticode/go-astiav
    note: Auto-resolved from go-astiav
  - name: go-astiav GitHub README
    url: https://github.com/asticode/go-astiav
    note: Auto-resolved from go-astiav-docs
  - name: gohlslib (HLS)
    url: https://pkg.go.dev/github.com/bluenviron/gohlslib/v2
    note: Auto-resolved from gohlslib
  - name: M3U8 Extended Format
    url: https://datatracker.ietf.org/doc/html/rfc8216
    note: Auto-resolved from m3u8
  - name: ogen OpenAPI Generator
    url: https://pkg.go.dev/github.com/ogen-go/ogen
    note: Auto-resolved from ogen
  - name: River Job Queue
    url: https://pkg.go.dev/github.com/riverqueue/river
    note: Auto-resolved from river
  - name: sqlc
    url: https://docs.sqlc.dev/en/stable/
    note: Auto-resolved from sqlc
  - name: sqlc Configuration
    url: https://docs.sqlc.dev/en/stable/reference/config.html
    note: Auto-resolved from sqlc-config
  - name: XMLTV Format
    url: https://github.com/XMLTV/xmltv/blob/master/xmltv.dtd
    note: Auto-resolved from xmltv
  - name: XMLTV Wiki
    url: https://wiki.xmltv.org/index.php/XMLTVFormat
    note: Auto-resolved from xmltv-wiki
design_refs:
  - title: features/livetv
    path: features/livetv.md
  - title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# Live TV & DVR


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: feature


> Content module for TV Shows, Seasons, Episodes

> Live television streaming and digital video recording

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
internal/content/live_tv_&_dvr/
├── module.go              # fx module definition
├── repository.go          # Database operations
├── service.go             # Business logic
├── handler.go             # HTTP handlers (ogen)
├── types.go               # Domain types
└── live_tv_&_dvr_test.go
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
- [features/livetv](features/livetv.md)
- [01_ARCHITECTURE](architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](architecture/03_METADATA_SYSTEM.md)

### External Sources
- [FFmpeg Documentation](https://ffmpeg.org/ffmpeg.html) - Auto-resolved from ffmpeg
- [FFmpeg Codecs](https://ffmpeg.org/ffmpeg-codecs.html) - Auto-resolved from ffmpeg-codecs
- [FFmpeg Formats](https://ffmpeg.org/ffmpeg-formats.html) - Auto-resolved from ffmpeg-formats
- [Uber fx](https://pkg.go.dev/go.uber.org/fx) - Auto-resolved from fx
- [go-astiav (FFmpeg bindings)](https://pkg.go.dev/github.com/asticode/go-astiav) - Auto-resolved from go-astiav
- [go-astiav GitHub README](https://github.com/asticode/go-astiav) - Auto-resolved from go-astiav-docs
- [gohlslib (HLS)](https://pkg.go.dev/github.com/bluenviron/gohlslib/v2) - Auto-resolved from gohlslib
- [M3U8 Extended Format](https://datatracker.ietf.org/doc/html/rfc8216) - Auto-resolved from m3u8
- [ogen OpenAPI Generator](https://pkg.go.dev/github.com/ogen-go/ogen) - Auto-resolved from ogen
- [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) - Auto-resolved from river
- [sqlc](https://docs.sqlc.dev/en/stable/) - Auto-resolved from sqlc
- [sqlc Configuration](https://docs.sqlc.dev/en/stable/reference/config.html) - Auto-resolved from sqlc-config
- [XMLTV Format](https://github.com/XMLTV/xmltv/blob/master/xmltv.dtd) - Auto-resolved from xmltv
- [XMLTV Wiki](https://wiki.xmltv.org/index.php/XMLTVFormat) - Auto-resolved from xmltv-wiki


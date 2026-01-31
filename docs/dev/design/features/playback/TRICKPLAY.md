## Table of Contents

- [Trickplay (Timeline Thumbnails)](#trickplay-timeline-thumbnails)
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
  - name: Roku BIF Format
    url: https://developer.roku.com/docs/developer-program/media-playback/trick-mode/bif-file-creation.md
    note: Auto-resolved from bif-spec
  - name: FFmpeg Documentation
    url: https://ffmpeg.org/ffmpeg.html
    note: Auto-resolved from ffmpeg
  - name: FFmpeg Codecs
    url: https://ffmpeg.org/ffmpeg-codecs.html
    note: Auto-resolved from ffmpeg-codecs
  - name: FFmpeg Formats
    url: https://ffmpeg.org/ffmpeg-formats.html
    note: Auto-resolved from ffmpeg-formats
  - name: go-astiav (FFmpeg bindings)
    url: https://pkg.go.dev/github.com/asticode/go-astiav
    note: Auto-resolved from go-astiav
  - name: go-astiav GitHub README
    url: https://github.com/asticode/go-astiav
    note: Auto-resolved from go-astiav-docs
  - name: Jellyfin Trickplay
    url: https://jellyfin.org/docs/general/server/media/trickplay/
    note: Auto-resolved from jellyfin-trickplay
  - name: River Job Queue
    url: https://pkg.go.dev/github.com/riverqueue/river
    note: Auto-resolved from river
  - name: WebVTT Specification
    url: https://www.w3.org/TR/webvtt1/
    note: Auto-resolved from webvtt
design_refs:
  - title: features/playback
    path: features/playback.md
  - title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# Trickplay (Timeline Thumbnails)


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: feature


> Content module for 

> Thumbnail previews on video seek bar

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
internal/content/trickplay_(timeline_thumbnails)/
├── module.go              # fx module definition
├── repository.go          # Database operations
├── service.go             # Business logic
├── handler.go             # HTTP handlers (ogen)
├── types.go               # Domain types
└── trickplay_(timeline_thumbnails)_test.go
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
- [features/playback](features/playback.md)
- [01_ARCHITECTURE](architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Roku BIF Format](https://developer.roku.com/docs/developer-program/media-playback/trick-mode/bif-file-creation.md) - Auto-resolved from bif-spec
- [FFmpeg Documentation](https://ffmpeg.org/ffmpeg.html) - Auto-resolved from ffmpeg
- [FFmpeg Codecs](https://ffmpeg.org/ffmpeg-codecs.html) - Auto-resolved from ffmpeg-codecs
- [FFmpeg Formats](https://ffmpeg.org/ffmpeg-formats.html) - Auto-resolved from ffmpeg-formats
- [go-astiav (FFmpeg bindings)](https://pkg.go.dev/github.com/asticode/go-astiav) - Auto-resolved from go-astiav
- [go-astiav GitHub README](https://github.com/asticode/go-astiav) - Auto-resolved from go-astiav-docs
- [Jellyfin Trickplay](https://jellyfin.org/docs/general/server/media/trickplay/) - Auto-resolved from jellyfin-trickplay
- [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) - Auto-resolved from river
- [WebVTT Specification](https://www.w3.org/TR/webvtt1/) - Auto-resolved from webvtt


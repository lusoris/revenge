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
    url: ../../sources/media/ffmpeg.md
    note: Auto-resolved from ffmpeg
- name: FFmpeg Codecs
    url: ../../sources/media/ffmpeg-codecs.md
    note: Auto-resolved from ffmpeg-codecs
- name: FFmpeg Formats
    url: ../../sources/media/ffmpeg-formats.md
    note: Auto-resolved from ffmpeg-formats
- name: Uber fx
    url: ../../sources/tooling/fx.md
    note: Auto-resolved from fx
- name: go-astiav (FFmpeg bindings)
    url: ../../sources/media/go-astiav.md
    note: Auto-resolved from go-astiav
- name: go-astiav GitHub README
    url: ../../sources/media/go-astiav-guide.md
    note: Auto-resolved from go-astiav-docs
- name: gohlslib (HLS)
    url: ../../sources/media/gohlslib.md
    note: Auto-resolved from gohlslib
- name: M3U8 Extended Format
    url: ../../sources/protocols/m3u8.md
    note: Auto-resolved from m3u8
- name: ogen OpenAPI Generator
    url: ../../sources/tooling/ogen.md
    note: Auto-resolved from ogen
- name: River Job Queue
    url: ../../sources/tooling/river.md
    note: Auto-resolved from river
- name: sqlc
    url: ../../sources/database/sqlc.md
    note: Auto-resolved from sqlc
- name: sqlc Configuration
    url: ../../sources/database/sqlc-config.md
    note: Auto-resolved from sqlc-config
- name: XMLTV Format
    url: ../../sources/protocols/xmltv.md
    note: Auto-resolved from xmltv
- name: XMLTV Wiki
    url: ../../sources/protocols/xmltv-format.md
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

<!-- DESIGN: features/livetv, README, SCAFFOLD_TEMPLATE, test_output_claude -->

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
- [01_ARCHITECTURE](../../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [FFmpeg Documentation](../../sources/media/ffmpeg.md) - Auto-resolved from ffmpeg
- [FFmpeg Codecs](../../sources/media/ffmpeg-codecs.md) - Auto-resolved from ffmpeg-codecs
- [FFmpeg Formats](../../sources/media/ffmpeg-formats.md) - Auto-resolved from ffmpeg-formats
- [Uber fx](../../sources/tooling/fx.md) - Auto-resolved from fx
- [go-astiav (FFmpeg bindings)](../../sources/media/go-astiav.md) - Auto-resolved from go-astiav
- [go-astiav GitHub README](../../sources/media/go-astiav-guide.md) - Auto-resolved from go-astiav-docs
- [gohlslib (HLS)](../../sources/media/gohlslib.md) - Auto-resolved from gohlslib
- [M3U8 Extended Format](../../sources/protocols/m3u8.md) - Auto-resolved from m3u8
- [ogen OpenAPI Generator](../../sources/tooling/ogen.md) - Auto-resolved from ogen
- [River Job Queue](../../sources/tooling/river.md) - Auto-resolved from river
- [sqlc](../../sources/database/sqlc.md) - Auto-resolved from sqlc
- [sqlc Configuration](../../sources/database/sqlc-config.md) - Auto-resolved from sqlc-config
- [XMLTV Format](../../sources/protocols/xmltv.md) - Auto-resolved from xmltv
- [XMLTV Wiki](../../sources/protocols/xmltv-format.md) - Auto-resolved from xmltv-wiki

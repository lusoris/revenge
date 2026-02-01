## Table of Contents

- [Skip Intro / Credits Detection](#skip-intro-credits-detection)
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
- name: Chromaprint/AcoustID
    url: ../../sources/standards/chromaprint.md
    note: Auto-resolved from chromaprint-acoustid
- name: FFmpeg Documentation
    url: ../../sources/media/ffmpeg.md
    note: Auto-resolved from ffmpeg
- name: FFmpeg Codecs
    url: ../../sources/media/ffmpeg-codecs.md
    note: Auto-resolved from ffmpeg-codecs
- name: FFmpeg Formats
    url: ../../sources/media/ffmpeg-formats.md
    note: Auto-resolved from ffmpeg-formats
- name: go-astiav (FFmpeg bindings)
    url: ../../sources/media/go-astiav.md
    note: Auto-resolved from go-astiav
- name: go-astiav GitHub README
    url: ../../sources/media/go-astiav-guide.md
    note: Auto-resolved from go-astiav-docs
- name: River Job Queue
    url: ../../sources/tooling/river.md
    note: Auto-resolved from river
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

# Skip Intro / Credits Detection

<!-- DESIGN: features/playback, README, SCAFFOLD_TEMPLATE, test_output_claude -->

**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: feature

> Content module for

> Automatic intro and credits detection with one-click skip

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
internal/content/skip_intro_/_credits_detection/
├── module.go              # fx module definition
├── repository.go          # Database operations
├── service.go             # Business logic
├── handler.go             # HTTP handlers (ogen)
├── types.go               # Domain types
└── skip_intro_/_credits_detection_test.go
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
- [01_ARCHITECTURE](../../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Chromaprint/AcoustID](../sources/standards/chromaprint.md) - Auto-resolved from chromaprint-acoustid
- [FFmpeg Documentation](../sources/media/ffmpeg.md) - Auto-resolved from ffmpeg
- [FFmpeg Codecs](../sources/media/ffmpeg-codecs.md) - Auto-resolved from ffmpeg-codecs
- [FFmpeg Formats](../sources/media/ffmpeg-formats.md) - Auto-resolved from ffmpeg-formats
- [go-astiav (FFmpeg bindings)](../sources/media/go-astiav.md) - Auto-resolved from go-astiav
- [go-astiav GitHub README](../sources/media/go-astiav-guide.md) - Auto-resolved from go-astiav-docs
- [River Job Queue](../sources/tooling/river.md) - Auto-resolved from river

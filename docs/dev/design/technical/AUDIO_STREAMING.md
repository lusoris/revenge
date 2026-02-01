---
sources:
  - name: gohlslib (HLS)
    url: ../../sources/media/gohlslib.md
    note: HLS server implementation
  - name: go-astiav (FFmpeg)
    url: ../../sources/media/go-astiav.md
    note: Audio transcoding
  - name: Dragonfly
    url: ../../sources/infrastructure/dragonfly.md
    note: Progress tracking cache
design_refs:
  - title: technical
    path: INDEX.md
  - title: 01_ARCHITECTURE
    path: ../architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: ../architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: ../architecture/03_METADATA_SYSTEM.md
---

## Table of Contents

- [Audio Streaming & Progress Tracking](#audio-streaming-progress-tracking)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Components](#components)
  - [Implementation](#implementation)
    - [File Structure](#file-structure)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
    - [Config Keys](#config-keys)
  - [API Endpoints](#api-endpoints)
    - [GET /api/v1/stream/:track_id/playlist.m3u8](#get-apiv1streamtrack_idplaylistm3u8)
    - [GET /api/v1/stream/:track_id/segment:N.ts](#get-apiv1streamtrack_idsegmentnts)
    - [GET /api/v1/playback/progress/:track_id](#get-apiv1playbackprogresstrack_id)
    - [PUT /api/v1/playback/progress/:track_id](#put-apiv1playbackprogresstrack_id)
    - [POST /api/v1/playback/scrobble](#post-apiv1playbackscrobble)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)


# Audio Streaming & Progress Tracking


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: technical


> > HLS audio streaming with adaptive bitrate and progress tracking

Audio streaming architecture:
- **Protocol**: HLS (HTTP Live Streaming) via gohlslib
- **Codecs**: AAC, MP3, FLAC, Opus (transcode on-demand)
- **Adaptive Bitrate**: Multiple quality streams (64k, 128k, 256k, 320k)
- **Progress Tracking**: Per-second accuracy with real-time sync
- **Session Management**: Resume playback across devices

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete audio streaming design |
| Sources | ✅ | All streaming tools documented |
| Instructions | ✅ | Generated from design |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---


## Architecture

<!-- Architecture diagram placeholder -->

### Components

<!-- Component description -->


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
### GET /api/v1/stream/:track_id/playlist.m3u8

HLS master playlist

**Request**:
```json
{}
```

**Response**:
```json
{}
```
### GET /api/v1/stream/:track_id/segment:N.ts

HLS segment

**Request**:
```json
{}
```

**Response**:
```json
{}
```
### GET /api/v1/playback/progress/:track_id

Get current progress

**Request**:
```json
{}
```

**Response**:
```json
{}
```
### PUT /api/v1/playback/progress/:track_id

Update progress

**Request**:
```json
{"position_seconds": 123.45, "duration_seconds": 245.0}
```

**Response**:
```json
{}
```
### POST /api/v1/playback/scrobble

Submit scrobble

**Request**:
```json
{}
```

**Response**:
```json
{}
```


## Testing Strategy

### Unit Tests

<!-- Unit test strategy -->

### Integration Tests

<!-- Integration test strategy -->

### Test Coverage

Target: **80% minimum**







## Related Documentation
### Design Documents
- [technical](INDEX.md)
- [01_ARCHITECTURE](../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [gohlslib (HLS)](../../sources/media/gohlslib.md) - HLS server implementation
- [go-astiav (FFmpeg)](../../sources/media/go-astiav.md) - Audio transcoding
- [Dragonfly](../../sources/infrastructure/dragonfly.md) - Progress tracking cache


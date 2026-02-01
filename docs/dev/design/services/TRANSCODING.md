## Table of Contents

- [Transcoding Service](#transcoding-service)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Service Structure](#service-structure)
    - [Dependencies](#dependencies)
    - [Provides](#provides)
    - [Component Diagram](#component-diagram)
  - [Implementation](#implementation)
    - [File Structure](#file-structure)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
    - [Config Keys](#config-keys)
  - [API Endpoints](#api-endpoints)
    - [POST /api/v1/transcode/video/:id](#post-apiv1transcodevideoid)
    - [GET /api/v1/transcode/job/:job_id](#get-apiv1transcodejobjob_id)
    - [DELETE /api/v1/transcode/job/:job_id](#delete-apiv1transcodejobjob_id)
    - [GET /api/v1/transcode/cache/stats](#get-apiv1transcodecachestats)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

---
sources:
- name: go-astiav (FFmpeg)
    url: ../../sources/media/go-astiav.md
    note: FFmpeg Go bindings
- name: gohlslib
    url: ../../sources/media/gohlslib.md
    note: HLS streaming library
- name: River Job Queue
    url: ../../sources/tooling/river.md
    note: Background job processing
- name: Uber fx
    url: ../../sources/tooling/fx.md
    note: Dependency injection
design_refs:
- title: services
    path: services.md
- title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
- title: OFFLOADING
    path: technical/OFFLOADING.md
- title: AUDIO_STREAMING
    path: technical/AUDIO_STREAMING.md
---

# Transcoding Service

<!-- DESIGN: services, README, SCAFFOLD_TEMPLATE, test_output_claude -->

**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: service

> > On-demand video/audio transcoding service with hardware acceleration

Transcoding capabilities:
- **Primary**: Blackbeard service for offloaded transcoding
- **Fallback**: Local FFmpeg with go-astiav bindings
- **Hardware Acceleration**: NVENC (NVIDIA), QSV (Intel), VAAPI (AMD)
- **Adaptive Streaming**: HLS with multiple quality tiers
- **Queue Management**: River background jobs for async processing

**Package**: `internal/service/transcoding`
**fx Module**: `transcoding.Module`

---

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete transcoding service design |
| Sources | ✅ | All transcoding tools documented |
| Instructions | ✅ | Generated from design |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |
**Overall**: ✅ Complete

---

## Architecture

### Service Structure

```
internal/service/transcoding/
├── module.go              # fx module definition
├── service.go             # Service implementation
├── repository.go          # Data access (if needed)
├── handler.go             # HTTP handlers (if exposed)
├── middleware.go          # Middleware (if needed)
├── types.go               # Domain types
└── service_test.go        # Tests
```

### Dependencies
No external service dependencies.

### Provides
<!-- Service provides -->

### Component Diagram

<!-- Component diagram -->

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
### POST /api/v1/transcode/video/:id

Start video transcoding

**Request**:
```json
{}
```

**Response**:
```json
{}
```
### GET /api/v1/transcode/job/:job_id

Get transcode job status

**Request**:
```json
{}
```

**Response**:
```json
{}
```
### DELETE /api/v1/transcode/job/:job_id

Cancel transcode job

**Request**:
```json
{}
```

**Response**:
```json
{}
```
### GET /api/v1/transcode/cache/stats

Get cache statistics

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
- [services](services.md)
- [01_ARCHITECTURE](../architecture/01_ARCHITECTURE.md)
- [OFFLOADING](../technical/OFFLOADING.md)
- [AUDIO_STREAMING](../technical/AUDIO_STREAMING.md)

### External Sources
- [go-astiav (FFmpeg)](../../sources/media/go-astiav.md) - FFmpeg Go bindings
- [gohlslib](../../sources/media/gohlslib.md) - HLS streaming library
- [River Job Queue](../../sources/tooling/river.md) - Background job processing
- [Uber fx](../../sources/tooling/fx.md) - Dependency injection

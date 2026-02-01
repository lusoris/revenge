---
sources:
  - name: Dragonfly Documentation
    url: ../../sources/infrastructure/dragonfly.md
    note: Auto-resolved from dragonfly
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
  - name: gohlslib (HLS)
    url: ../../sources/media/gohlslib.md
    note: Auto-resolved from gohlslib
  - name: M3U8 Extended Format
    url: ../../sources/protocols/m3u8.md
    note: Auto-resolved from m3u8
  - name: River Job Queue
    url: ../../sources/tooling/river.md
    note: Auto-resolved from river
  - name: shadcn-svelte
    url: ../../sources/frontend/shadcn-svelte.md
    note: Auto-resolved from shadcn-svelte
  - name: Svelte 5 Runes
    url: ../../sources/frontend/svelte-runes.md
    note: Auto-resolved from svelte-runes
  - name: Svelte 5 Documentation
    url: ../../sources/frontend/svelte5.md
    note: Auto-resolved from svelte5
  - name: SvelteKit Documentation
    url: ../../sources/frontend/sveltekit.md
    note: Auto-resolved from sveltekit
  - name: TanStack Query
    url: ../../sources/frontend/tanstack-query.md
    note: Auto-resolved from tanstack-query
  - name: Typesense API
    url: ../../sources/infrastructure/typesense.md
    note: Auto-resolved from typesense
  - name: Typesense Go Client
    url: ../../sources/infrastructure/typesense-go.md
    note: Auto-resolved from typesense-go
design_refs:
  - title: architecture
    path: INDEX.md
  - title: 01_ARCHITECTURE
    path: 01_ARCHITECTURE.md
  - title: 03_METADATA_SYSTEM
    path: 03_METADATA_SYSTEM.md
  - title: METADATA (service)
    path: ../services/METADATA.md
  - title: HTTP_CLIENT (service)
    path: ../services/HTTP_CLIENT.md
  - title: ADULT_CONTENT_SYSTEM
    path: ../features/adult/ADULT_CONTENT_SYSTEM.md
  - title: ADULT_METADATA
    path: ../features/adult/ADULT_METADATA.md
  - title: DATA_RECONCILIATION
    path: ../features/adult/DATA_RECONCILIATION.md
---

## Table of Contents

- [Revenge - Design Principles](#revenge-design-principles)
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
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)


# Revenge - Design Principles


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: architecture


> > Core design principles and patterns guiding Revenge development

Key principles:
- **PostgreSQL Only**: No SQLite - simplifies codebase, pgx is best driver
- **80% Test Coverage**: Required minimum for all packages
- **Local-First Metadata**: Always prefer cached/local data over external APIs
- **Table-Driven Tests**: testify + mockery for consistent testing
- **Sentinel Errors**: Type-safe error handling with wrapped contexts


---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | ⚪ | - |
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




## Testing Strategy

### Unit Tests

<!-- Unit test strategy -->

### Integration Tests

<!-- Integration test strategy -->

### Test Coverage

Target: **80% minimum**







## Related Documentation
### Design Documents
- [architecture](INDEX.md)
- [01_ARCHITECTURE](01_ARCHITECTURE.md)
- [03_METADATA_SYSTEM](03_METADATA_SYSTEM.md)
- [METADATA (service)](../services/METADATA.md)
- [HTTP_CLIENT (service)](../services/HTTP_CLIENT.md)
- [ADULT_CONTENT_SYSTEM](../features/adult/ADULT_CONTENT_SYSTEM.md)
- [ADULT_METADATA](../features/adult/ADULT_METADATA.md)
- [DATA_RECONCILIATION](../features/adult/DATA_RECONCILIATION.md)

### External Sources
- [Dragonfly Documentation](../../sources/infrastructure/dragonfly.md) - Auto-resolved from dragonfly
- [FFmpeg Documentation](../../sources/media/ffmpeg.md) - Auto-resolved from ffmpeg
- [FFmpeg Codecs](../../sources/media/ffmpeg-codecs.md) - Auto-resolved from ffmpeg-codecs
- [FFmpeg Formats](../../sources/media/ffmpeg-formats.md) - Auto-resolved from ffmpeg-formats
- [go-astiav (FFmpeg bindings)](../../sources/media/go-astiav.md) - Auto-resolved from go-astiav
- [go-astiav GitHub README](../../sources/media/go-astiav-guide.md) - Auto-resolved from go-astiav-docs
- [gohlslib (HLS)](../../sources/media/gohlslib.md) - Auto-resolved from gohlslib
- [M3U8 Extended Format](../../sources/protocols/m3u8.md) - Auto-resolved from m3u8
- [River Job Queue](../../sources/tooling/river.md) - Auto-resolved from river
- [shadcn-svelte](../../sources/frontend/shadcn-svelte.md) - Auto-resolved from shadcn-svelte
- [Svelte 5 Runes](../../sources/frontend/svelte-runes.md) - Auto-resolved from svelte-runes
- [Svelte 5 Documentation](../../sources/frontend/svelte5.md) - Auto-resolved from svelte5
- [SvelteKit Documentation](../../sources/frontend/sveltekit.md) - Auto-resolved from sveltekit
- [TanStack Query](../../sources/frontend/tanstack-query.md) - Auto-resolved from tanstack-query
- [Typesense API](../../sources/infrastructure/typesense.md) - Auto-resolved from typesense
- [Typesense Go Client](../../sources/infrastructure/typesense-go.md) - Auto-resolved from typesense-go


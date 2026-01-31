## Table of Contents

- [Revenge - Architecture v2](#revenge-architecture-v2)
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



---
sources:
  - name: Dragonfly Documentation
    url: https://www.dragonflydb.io/docs
    note: Auto-resolved from dragonfly
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
  - name: go-blurhash
    url: https://pkg.go.dev/github.com/bbrks/go-blurhash
    note: Auto-resolved from go-blurhash
  - name: gohlslib (HLS)
    url: https://pkg.go.dev/github.com/bluenviron/gohlslib/v2
    note: Auto-resolved from gohlslib
  - name: koanf
    url: https://pkg.go.dev/github.com/knadh/koanf/v2
    note: Auto-resolved from koanf
  - name: Last.fm API
    url: https://www.last.fm/api/intro
    note: Auto-resolved from lastfm-api
  - name: M3U8 Extended Format
    url: https://datatracker.ietf.org/doc/html/rfc8216
    note: Auto-resolved from m3u8
  - name: ogen OpenAPI Generator
    url: https://pkg.go.dev/github.com/ogen-go/ogen
    note: Auto-resolved from ogen
  - name: pgx PostgreSQL Driver
    url: https://pkg.go.dev/github.com/jackc/pgx/v5
    note: Auto-resolved from pgx
  - name: PostgreSQL Arrays
    url: https://www.postgresql.org/docs/current/arrays.html
    note: Auto-resolved from postgresql-arrays
  - name: PostgreSQL JSON Functions
    url: https://www.postgresql.org/docs/current/functions-json.html
    note: Auto-resolved from postgresql-json
  - name: River Job Queue
    url: https://pkg.go.dev/github.com/riverqueue/river
    note: Auto-resolved from river
  - name: rueidis
    url: https://pkg.go.dev/github.com/redis/rueidis
    note: Auto-resolved from rueidis
  - name: rueidis GitHub README
    url: https://github.com/redis/rueidis
    note: Auto-resolved from rueidis-docs
  - name: shadcn-svelte
    url: https://www.shadcn-svelte.com/docs
    note: Auto-resolved from shadcn-svelte
  - name: sqlc
    url: https://docs.sqlc.dev/en/stable/
    note: Auto-resolved from sqlc
  - name: sqlc Configuration
    url: https://docs.sqlc.dev/en/stable/reference/config.html
    note: Auto-resolved from sqlc-config
  - name: Svelte 5 Runes
    url: https://svelte.dev/docs/svelte/$state
    note: Auto-resolved from svelte-runes
  - name: Svelte 5 Documentation
    url: https://svelte.dev/docs/svelte/overview
    note: Auto-resolved from svelte5
  - name: SvelteKit Documentation
    url: https://svelte.dev/docs/kit/introduction
    note: Auto-resolved from sveltekit
  - name: TanStack Query
    url: https://tanstack.com/query/latest/docs/framework/svelte/overview
    note: Auto-resolved from tanstack-query
  - name: Typesense API
    url: https://typesense.org/docs/latest/api/
    note: Auto-resolved from typesense
  - name: Typesense Go Client
    url: https://github.com/typesense/typesense-go
    note: Auto-resolved from typesense-go
design_refs:
  - title: architecture
    path: architecture/INDEX.md
  - title: ADULT_CONTENT_SYSTEM
    path: ADULT_CONTENT_SYSTEM.md
  - title: ADULT_METADATA
    path: ADULT_METADATA.md
  - title: DATA_RECONCILIATION
    path: DATA_RECONCILIATION.md
---

# Revenge - Architecture v2


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: architecture


> PLACEHOLDER: Brief technical summary

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | ⚪ | - |
| Instructions | 🔴 | - |
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
- [architecture](architecture/INDEX.md)
- [ADULT_CONTENT_SYSTEM](ADULT_CONTENT_SYSTEM.md)
- [ADULT_METADATA](ADULT_METADATA.md)
- [DATA_RECONCILIATION](DATA_RECONCILIATION.md)

### External Sources
- [Dragonfly Documentation](https://www.dragonflydb.io/docs) - Auto-resolved from dragonfly
- [FFmpeg Documentation](https://ffmpeg.org/ffmpeg.html) - Auto-resolved from ffmpeg
- [FFmpeg Codecs](https://ffmpeg.org/ffmpeg-codecs.html) - Auto-resolved from ffmpeg-codecs
- [FFmpeg Formats](https://ffmpeg.org/ffmpeg-formats.html) - Auto-resolved from ffmpeg-formats
- [Uber fx](https://pkg.go.dev/go.uber.org/fx) - Auto-resolved from fx
- [go-astiav (FFmpeg bindings)](https://pkg.go.dev/github.com/asticode/go-astiav) - Auto-resolved from go-astiav
- [go-astiav GitHub README](https://github.com/asticode/go-astiav) - Auto-resolved from go-astiav-docs
- [go-blurhash](https://pkg.go.dev/github.com/bbrks/go-blurhash) - Auto-resolved from go-blurhash
- [gohlslib (HLS)](https://pkg.go.dev/github.com/bluenviron/gohlslib/v2) - Auto-resolved from gohlslib
- [koanf](https://pkg.go.dev/github.com/knadh/koanf/v2) - Auto-resolved from koanf
- [Last.fm API](https://www.last.fm/api/intro) - Auto-resolved from lastfm-api
- [M3U8 Extended Format](https://datatracker.ietf.org/doc/html/rfc8216) - Auto-resolved from m3u8
- [ogen OpenAPI Generator](https://pkg.go.dev/github.com/ogen-go/ogen) - Auto-resolved from ogen
- [pgx PostgreSQL Driver](https://pkg.go.dev/github.com/jackc/pgx/v5) - Auto-resolved from pgx
- [PostgreSQL Arrays](https://www.postgresql.org/docs/current/arrays.html) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](https://www.postgresql.org/docs/current/functions-json.html) - Auto-resolved from postgresql-json
- [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) - Auto-resolved from river
- [rueidis](https://pkg.go.dev/github.com/redis/rueidis) - Auto-resolved from rueidis
- [rueidis GitHub README](https://github.com/redis/rueidis) - Auto-resolved from rueidis-docs
- [shadcn-svelte](https://www.shadcn-svelte.com/docs) - Auto-resolved from shadcn-svelte
- [sqlc](https://docs.sqlc.dev/en/stable/) - Auto-resolved from sqlc
- [sqlc Configuration](https://docs.sqlc.dev/en/stable/reference/config.html) - Auto-resolved from sqlc-config
- [Svelte 5 Runes](https://svelte.dev/docs/svelte/$state) - Auto-resolved from svelte-runes
- [Svelte 5 Documentation](https://svelte.dev/docs/svelte/overview) - Auto-resolved from svelte5
- [SvelteKit Documentation](https://svelte.dev/docs/kit/introduction) - Auto-resolved from sveltekit
- [TanStack Query](https://tanstack.com/query/latest/docs/framework/svelte/overview) - Auto-resolved from tanstack-query
- [Typesense API](https://typesense.org/docs/latest/api/) - Auto-resolved from typesense
- [Typesense Go Client](https://github.com/typesense/typesense-go) - Auto-resolved from typesense-go


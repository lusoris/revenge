## Table of Contents

- [Music Module](#music-module)
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
      - [GET /api/v1/music/artists](#get-apiv1musicartists)
      - [GET /api/v1/music/artists/:id](#get-apiv1musicartistsid)
      - [GET /api/v1/music/artists/:id/albums](#get-apiv1musicartistsidalbums)
      - [GET /api/v1/music/albums](#get-apiv1musicalbums)
      - [GET /api/v1/music/albums/:id](#get-apiv1musicalbumsid)
      - [GET /api/v1/music/albums/:id/tracks](#get-apiv1musicalbumsidtracks)
      - [GET /api/v1/music/tracks](#get-apiv1musictracks)
      - [GET /api/v1/music/tracks/:id](#get-apiv1musictracksid)
      - [GET /api/v1/music/tracks/:id/stream](#get-apiv1musictracksidstream)
      - [POST /api/v1/music/playlists](#post-apiv1musicplaylists)
      - [GET /api/v1/music/genres](#get-apiv1musicgenres)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)



---
sources:
  - name: Uber fx
    url: ../sources/tooling/fx.md
    note: Auto-resolved from fx
  - name: Last.fm API
    url: ../sources/apis/lastfm.md
    note: Auto-resolved from lastfm-api
  - name: ogen OpenAPI Generator
    url: ../sources/tooling/ogen.md
    note: Auto-resolved from ogen
  - name: pgx PostgreSQL Driver
    url: ../sources/database/pgx.md
    note: Auto-resolved from pgx
  - name: PostgreSQL Arrays
    url: ../sources/database/postgresql-arrays.md
    note: Auto-resolved from postgresql-arrays
  - name: PostgreSQL JSON Functions
    url: ../sources/database/postgresql-json.md
    note: Auto-resolved from postgresql-json
  - name: River Job Queue
    url: ../sources/tooling/river.md
    note: Auto-resolved from river
  - name: sqlc
    url: ../sources/database/sqlc.md
    note: Auto-resolved from sqlc
  - name: sqlc Configuration
    url: ../sources/database/sqlc-config.md
    note: Auto-resolved from sqlc-config
design_refs:
  - title: features/music
    path: features/music.md
  - title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# Music Module


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: feature


> Content module for Artists, Albums, Tracks

> Music content management with metadata enrichment from MusicBrainz and Last.fm

Complete music library management:
- **Lidarr Integration**: Two-way sync for music library automation
- **Metadata Sources**: MusicBrainz (primary), Last.fm (scrobbling, tags)
- **Supported Formats**: MP3, FLAC, AAC, OGG, ALAC, Opus
- **Playback**: HLS adaptive streaming with gapless playback
- **Features**: Playlists, smart collections, album art, lyrics

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete music module design |
| Sources | ✅ | All music APIs documented |
| Instructions | ✅ | Generated from design |
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
internal/content/music/
├── module.go              # fx module definition
├── repository.go          # Database operations
├── service.go             # Business logic
├── handler.go             # HTTP handlers (ogen)
├── types.go               # Domain types
└── music_test.go
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
#### GET /api/v1/music/artists

---
#### GET /api/v1/music/artists/:id

---
#### GET /api/v1/music/artists/:id/albums

---
#### GET /api/v1/music/albums

---
#### GET /api/v1/music/albums/:id

---
#### GET /api/v1/music/albums/:id/tracks

---
#### GET /api/v1/music/tracks

---
#### GET /api/v1/music/tracks/:id

---
#### GET /api/v1/music/tracks/:id/stream

---
#### POST /api/v1/music/playlists

---
#### GET /api/v1/music/genres

---


## Testing Strategy

### Unit Tests

<!-- Unit test strategy -->

### Integration Tests

<!-- Integration test strategy -->

### Test Coverage

Target: **80% minimum**







## Related Documentation
### Design Documents
- [features/music](features/music.md)
- [01_ARCHITECTURE](architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Uber fx](../sources/tooling/fx.md) - Auto-resolved from fx
- [Last.fm API](../sources/apis/lastfm.md) - Auto-resolved from lastfm-api
- [ogen OpenAPI Generator](../sources/tooling/ogen.md) - Auto-resolved from ogen
- [pgx PostgreSQL Driver](../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [River Job Queue](../sources/tooling/river.md) - Auto-resolved from river
- [sqlc](../sources/database/sqlc.md) - Auto-resolved from sqlc
- [sqlc Configuration](../sources/database/sqlc-config.md) - Auto-resolved from sqlc-config


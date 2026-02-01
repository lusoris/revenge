## Table of Contents

- [Audiobook Module](#audiobook-module)
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
      - [GET /api/v1/audiobooks](#get-apiv1audiobooks)
      - [GET /api/v1/audiobooks/:id](#get-apiv1audiobooksid)
      - [GET /api/v1/audiobooks/:id/chapters](#get-apiv1audiobooksidchapters)
      - [GET /api/v1/audiobooks/:id/stream](#get-apiv1audiobooksidstream)
      - [GET /api/v1/audiobooks/:id/progress](#get-apiv1audiobooksidprogress)
      - [PUT /api/v1/audiobooks/:id/progress](#put-apiv1audiobooksidprogress)
      - [POST /api/v1/audiobooks/:id/bookmarks](#post-apiv1audiobooksidbookmarks)
      - [GET /api/v1/audiobooks/authors](#get-apiv1audiobooksauthors)
      - [GET /api/v1/audiobooks/series](#get-apiv1audiobooksseries)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)



---
sources:
  - name: Audnexus API
    url: ../../../sources/apis/audnexus.md
    note: Auto-resolved from audnexus
  - name: Uber fx
    url: ../../../sources/tooling/fx.md
    note: Auto-resolved from fx
  - name: ogen OpenAPI Generator
    url: ../../../sources/tooling/ogen.md
    note: Auto-resolved from ogen
  - name: Open Library API
    url: ../../../sources/apis/openlibrary.md
    note: Auto-resolved from openlibrary
  - name: pgx PostgreSQL Driver
    url: ../../../sources/database/pgx.md
    note: Auto-resolved from pgx
  - name: PostgreSQL Arrays
    url: ../../../sources/database/postgresql-arrays.md
    note: Auto-resolved from postgresql-arrays
  - name: PostgreSQL JSON Functions
    url: ../../../sources/database/postgresql-json.md
    note: Auto-resolved from postgresql-json
  - name: River Job Queue
    url: ../../../sources/tooling/river.md
    note: Auto-resolved from river
  - name: sqlc
    url: ../../../sources/database/sqlc.md
    note: Auto-resolved from sqlc
  - name: sqlc Configuration
    url: ../../../sources/database/sqlc-config.md
    note: Auto-resolved from sqlc-config
design_refs:
  - title: 01_ARCHITECTURE
    path: ../../architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: ../../architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: ../../architecture/03_METADATA_SYSTEM.md
---

# Audiobook Module


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: feature


> Content module for Books, Authors, Series

> Audiobook content management with metadata enrichment from Audnexus and OpenLibrary

Complete audiobook library:
- **Chaptarr Integration**: Automated audiobook management and metadata sync
- **Metadata Sources**: Audnexus (primary), OpenLibrary, Goodreads
- **Supported Formats**: M4B (with chapters), MP3 (multi-file), AAC
- **Chapter Navigation**: Jump to chapters, bookmarks, progress tracking
- **Playback**: Variable speed (0.5x-3x), sleep timer, per-user resume

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete audiobook module design |
| Sources | ✅ | All audiobook APIs documented |
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
internal/content/audiobook/
├── module.go              # fx module definition
├── repository.go          # Database operations
├── service.go             # Business logic
├── handler.go             # HTTP handlers (ogen)
├── types.go               # Domain types
└── audiobook_test.go
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
#### GET /api/v1/audiobooks

List all audiobooks with pagination and filters

---
#### GET /api/v1/audiobooks/:id

Get audiobook details by ID

---
#### GET /api/v1/audiobooks/:id/chapters

Get chapter list for an audiobook

---
#### GET /api/v1/audiobooks/:id/stream

Get HLS streaming URL for audiobook playback

---
#### GET /api/v1/audiobooks/:id/progress

Get user playback progress for an audiobook

---
#### PUT /api/v1/audiobooks/:id/progress

Update user playback progress

---
#### POST /api/v1/audiobooks/:id/bookmarks

Create a bookmark at current position

---
#### GET /api/v1/audiobooks/authors

List all audiobook authors

---
#### GET /api/v1/audiobooks/series

List all audiobook series

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
- [01_ARCHITECTURE](../../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Audnexus API](../../../sources/apis/audnexus.md) - Auto-resolved from audnexus
- [Uber fx](../../../sources/tooling/fx.md) - Auto-resolved from fx
- [ogen OpenAPI Generator](../../../sources/tooling/ogen.md) - Auto-resolved from ogen
- [Open Library API](../../../sources/apis/openlibrary.md) - Auto-resolved from openlibrary
- [pgx PostgreSQL Driver](../../../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../../../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../../../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [River Job Queue](../../../sources/tooling/river.md) - Auto-resolved from river
- [sqlc](../../../sources/database/sqlc.md) - Auto-resolved from sqlc
- [sqlc Configuration](../../../sources/database/sqlc-config.md) - Auto-resolved from sqlc-config


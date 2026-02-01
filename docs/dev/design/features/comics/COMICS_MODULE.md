## Table of Contents

- [Comics Module](#comics-module)
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
      - [GET /api/v1/comics/series](#get-apiv1comicsseries)
      - [GET /api/v1/comics/series/:id](#get-apiv1comicsseriesid)
      - [GET /api/v1/comics/series/:id/issues](#get-apiv1comicsseriesidissues)
      - [GET /api/v1/comics/issues](#get-apiv1comicsissues)
      - [GET /api/v1/comics/issues/:id](#get-apiv1comicsissuesid)
      - [GET /api/v1/comics/issues/:id/read](#get-apiv1comicsissuesidread)
      - [GET /api/v1/comics/issues/:id/pages/:page](#get-apiv1comicsissuesidpagespage)
      - [GET /api/v1/comics/issues/:id/progress](#get-apiv1comicsissuesidprogress)
      - [PUT /api/v1/comics/issues/:id/progress](#put-apiv1comicsissuesidprogress)
      - [GET /api/v1/comics/publishers](#get-apiv1comicspublishers)
      - [POST /api/v1/comics/pull-list](#post-apiv1comicspull-list)
      - [GET /api/v1/comics/pull-list](#get-apiv1comicspull-list)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)



---
sources:
  - name: AniList GraphQL API
    url: ../sources/apis/anilist.md
    note: Auto-resolved from anilist
  - name: ComicVine API
    url: ../sources/apis/comicvine.md
    note: Auto-resolved from comicvine
  - name: Uber fx
    url: ../sources/tooling/fx.md
    note: Auto-resolved from fx
  - name: MyAnimeList API
    url: ../sources/apis/myanimelist.md
    note: Auto-resolved from myanimelist
  - name: ogen OpenAPI Generator
    url: ../sources/tooling/ogen.md
    note: Auto-resolved from ogen
  - name: River Job Queue
    url: ../sources/tooling/river.md
    note: Auto-resolved from river
  - name: sqlc
    url: ../sources/database/sqlc.md
    note: Auto-resolved from sqlc
  - name: sqlc Configuration
    url: ../sources/database/sqlc-config.md
    note: Auto-resolved from sqlc-config
  - name: Svelte 5 Runes
    url: ../sources/frontend/svelte-runes.md
    note: Auto-resolved from svelte-runes
  - name: Svelte 5 Documentation
    url: ../sources/frontend/svelte5.md
    note: Auto-resolved from svelte5
  - name: SvelteKit Documentation
    url: ../sources/frontend/sveltekit.md
    note: Auto-resolved from sveltekit
design_refs:
  - title: features/comics
    path: features/comics.md
  - title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# Comics Module


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: feature


> Content module for Comics, Issues, Series

> Digital comics/manga/graphic novel support with metadata from ComicVine, Marvel API, GCD

Complete comics library:
- **Metadata Sources**: ComicVine (primary), Marvel API, Grand Comics Database, AniList/MAL (manga)
- **Supported Formats**: CBZ, CBR, CB7, CBT, PDF
- **Reader Features**: Page-by-page viewing, two-page spread, webtoon scroll mode
- **Progress Tracking**: Per-user reading progress with sync across devices
- **Collection Management**: Pull lists, reading lists, series tracking

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
internal/content/comics/
├── module.go              # fx module definition
├── repository.go          # Database operations
├── service.go             # Business logic
├── handler.go             # HTTP handlers (ogen)
├── types.go               # Domain types
└── comics_test.go
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
#### GET /api/v1/comics/series

List all comic series with pagination and filters

---
#### GET /api/v1/comics/series/:id

Get comic series details by ID

---
#### GET /api/v1/comics/series/:id/issues

List all issues in a series

---
#### GET /api/v1/comics/issues

List all comic issues with pagination and filters

---
#### GET /api/v1/comics/issues/:id

Get comic issue details by ID

---
#### GET /api/v1/comics/issues/:id/read

Get issue content for comic reader

---
#### GET /api/v1/comics/issues/:id/pages/:page

Get a specific page image from an issue

---
#### GET /api/v1/comics/issues/:id/progress

Get user reading progress for an issue

---
#### PUT /api/v1/comics/issues/:id/progress

Update user reading progress for an issue

---
#### GET /api/v1/comics/publishers

List all comic publishers

---
#### POST /api/v1/comics/pull-list

Add a series to user pull list

---
#### GET /api/v1/comics/pull-list

Get user pull list of series

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
- [features/comics](features/comics.md)
- [01_ARCHITECTURE](../../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [AniList GraphQL API](../sources/apis/anilist.md) - Auto-resolved from anilist
- [ComicVine API](../sources/apis/comicvine.md) - Auto-resolved from comicvine
- [Uber fx](../sources/tooling/fx.md) - Auto-resolved from fx
- [MyAnimeList API](../sources/apis/myanimelist.md) - Auto-resolved from myanimelist
- [ogen OpenAPI Generator](../sources/tooling/ogen.md) - Auto-resolved from ogen
- [River Job Queue](../sources/tooling/river.md) - Auto-resolved from river
- [sqlc](../sources/database/sqlc.md) - Auto-resolved from sqlc
- [sqlc Configuration](../sources/database/sqlc-config.md) - Auto-resolved from sqlc-config
- [Svelte 5 Runes](../sources/frontend/svelte-runes.md) - Auto-resolved from svelte-runes
- [Svelte 5 Documentation](../sources/frontend/svelte5.md) - Auto-resolved from svelte5
- [SvelteKit Documentation](../sources/frontend/sveltekit.md) - Auto-resolved from sveltekit


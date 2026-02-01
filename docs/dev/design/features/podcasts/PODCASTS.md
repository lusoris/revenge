## Table of Contents

- [Podcasts](#podcasts)
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
      - [GET /api/v1/podcasts](#get-apiv1podcasts)
      - [POST /api/v1/podcasts](#post-apiv1podcasts)
      - [GET /api/v1/podcasts/:id](#get-apiv1podcastsid)
      - [DELETE /api/v1/podcasts/:id](#delete-apiv1podcastsid)
      - [GET /api/v1/podcasts/:id/episodes](#get-apiv1podcastsidepisodes)
      - [GET /api/v1/podcasts/episodes/:id](#get-apiv1podcastsepisodesid)
      - [GET /api/v1/podcasts/episodes/:id/stream](#get-apiv1podcastsepisodesidstream)
      - [POST /api/v1/podcasts/episodes/:id/download](#post-apiv1podcastsepisodesiddownload)
      - [GET /api/v1/podcasts/episodes/:id/progress](#get-apiv1podcastsepisodesidprogress)
      - [PUT /api/v1/podcasts/episodes/:id/progress](#put-apiv1podcastsepisodesidprogress)
      - [GET /api/v1/podcasts/search](#get-apiv1podcastssearch)
      - [POST /api/v1/podcasts/import-opml](#post-apiv1podcastsimport-opml)
      - [GET /api/v1/podcasts/export-opml](#get-apiv1podcastsexport-opml)
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
design_refs:
  - title: features/podcasts
    path: features/podcasts.md
  - title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# Podcasts

<!-- DESIGN: features/podcasts, README, SCAFFOLD_TEMPLATE, test_output_claude -->


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: feature


> Content module for Podcasts, Episodes

> RSS podcast subscription and playback

Complete podcast experience:
- **RSS Feed Support**: Subscribe to any podcast via RSS/Atom feeds
- **Automatic Updates**: Background jobs refresh feeds and download new episodes
- **Playback Features**: Variable speed, chapter navigation, sleep timer
- **Offline Support**: Download episodes for offline listening
- **Discovery**: Search and browse podcasts via Podcast Index API

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
internal/content/podcasts/
├── module.go              # fx module definition
├── repository.go          # Database operations
├── service.go             # Business logic
├── handler.go             # HTTP handlers (ogen)
├── types.go               # Domain types
└── podcasts_test.go
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
#### GET /api/v1/podcasts

List all subscribed podcasts

---
#### POST /api/v1/podcasts

Subscribe to a podcast by RSS URL

---
#### GET /api/v1/podcasts/:id

Get podcast details by ID

---
#### DELETE /api/v1/podcasts/:id

Unsubscribe from a podcast

---
#### GET /api/v1/podcasts/:id/episodes

List all episodes for a podcast

---
#### GET /api/v1/podcasts/episodes/:id

Get episode details by ID

---
#### GET /api/v1/podcasts/episodes/:id/stream

Get streaming URL for an episode

---
#### POST /api/v1/podcasts/episodes/:id/download

Download an episode for offline listening

---
#### GET /api/v1/podcasts/episodes/:id/progress

Get user playback progress for an episode

---
#### PUT /api/v1/podcasts/episodes/:id/progress

Update user playback progress for an episode

---
#### GET /api/v1/podcasts/search

Search podcasts via Podcast Index API

---
#### POST /api/v1/podcasts/import-opml

Import podcast subscriptions from OPML file

---
#### GET /api/v1/podcasts/export-opml

Export podcast subscriptions as OPML file

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
- [features/podcasts](features/podcasts.md)
- [01_ARCHITECTURE](../../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Uber fx](../sources/tooling/fx.md) - Auto-resolved from fx
- [ogen OpenAPI Generator](../sources/tooling/ogen.md) - Auto-resolved from ogen
- [River Job Queue](../sources/tooling/river.md) - Auto-resolved from river
- [sqlc](../sources/database/sqlc.md) - Auto-resolved from sqlc
- [sqlc Configuration](../sources/database/sqlc-config.md) - Auto-resolved from sqlc-config


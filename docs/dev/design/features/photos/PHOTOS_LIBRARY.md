## Table of Contents

- [Photos Library](#photos-library)
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
      - [GET /api/v1/photos](#get-apiv1photos)
      - [GET /api/v1/photos/:id](#get-apiv1photosid)
      - [POST /api/v1/photos](#post-apiv1photos)
      - [DELETE /api/v1/photos/:id](#delete-apiv1photosid)
      - [GET /api/v1/photos/:id/download](#get-apiv1photosiddownload)
      - [GET /api/v1/photos/:id/thumbnail/:size](#get-apiv1photosidthumbnailsize)
      - [PUT /api/v1/photos/:id/tags](#put-apiv1photosidtags)
      - [PUT /api/v1/photos/:id/favorite](#put-apiv1photosidfavorite)
      - [GET /api/v1/photos/albums](#get-apiv1photosalbums)
      - [POST /api/v1/photos/albums](#post-apiv1photosalbums)
      - [GET /api/v1/photos/albums/:id](#get-apiv1photosalbumsid)
      - [POST /api/v1/photos/albums/:id/photos](#post-apiv1photosalbumsidphotos)
      - [GET /api/v1/photos/people](#get-apiv1photospeople)
      - [POST /api/v1/photos/people](#post-apiv1photospeople)
      - [GET /api/v1/photos/timeline](#get-apiv1photostimeline)
      - [GET /api/v1/photos/map](#get-apiv1photosmap)
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
    url: ../../../sources/tooling/fx.md
    note: Auto-resolved from fx
  - name: go-blurhash
    url: ../../../sources/media/go-blurhash.md
    note: Auto-resolved from go-blurhash
  - name: ogen OpenAPI Generator
    url: ../../../sources/tooling/ogen.md
    note: Auto-resolved from ogen
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
  - title: features/photos
    path: features/photos.md
  - title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# Photos Library


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: feature


> Content module for Albums, Photos

> Photo organization, viewing, and management

Complete photo library:
- **Supported Formats**: JPEG, PNG, WebP, HEIC, RAW (CR2, NEF, ARW, DNG)
- **EXIF Extraction**: GPS, camera, lens, settings metadata
- **Organization**: Albums, people/faces, places, events, tags
- **Image Processing**: Thumbnails, blurhash placeholders, format conversion
- **Viewing**: Lightbox gallery, slideshow, map view

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
internal/content/photos_library/
├── module.go              # fx module definition
├── repository.go          # Database operations
├── service.go             # Business logic
├── handler.go             # HTTP handlers (ogen)
├── types.go               # Domain types
└── photos_library_test.go
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
#### GET /api/v1/photos

List all photos with pagination and filters

---
#### GET /api/v1/photos/:id

Get photo details by ID

---
#### POST /api/v1/photos

Upload a new photo

---
#### DELETE /api/v1/photos/:id

Delete a photo

---
#### GET /api/v1/photos/:id/download

Download original photo file

---
#### GET /api/v1/photos/:id/thumbnail/:size

Get photo thumbnail at specified size

---
#### PUT /api/v1/photos/:id/tags

Update tags for a photo

---
#### PUT /api/v1/photos/:id/favorite

Toggle favorite status for a photo

---
#### GET /api/v1/photos/albums

List all photo albums

---
#### POST /api/v1/photos/albums

Create a new photo album

---
#### GET /api/v1/photos/albums/:id

Get album details by ID

---
#### POST /api/v1/photos/albums/:id/photos

Add photos to an album

---
#### GET /api/v1/photos/people

List all tagged people

---
#### POST /api/v1/photos/people

Create a new person tag

---
#### GET /api/v1/photos/timeline

Get photos organized by timeline

---
#### GET /api/v1/photos/map

Get photos with location data for map view

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
- [features/photos](features/photos.md)
- [01_ARCHITECTURE](architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Uber fx](../../../sources/tooling/fx.md) - Auto-resolved from fx
- [go-blurhash](../../../sources/media/go-blurhash.md) - Auto-resolved from go-blurhash
- [ogen OpenAPI Generator](../../../sources/tooling/ogen.md) - Auto-resolved from ogen
- [pgx PostgreSQL Driver](../../../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../../../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../../../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [River Job Queue](../../../sources/tooling/river.md) - Auto-resolved from river
- [sqlc](../../../sources/database/sqlc.md) - Auto-resolved from sqlc
- [sqlc Configuration](../../../sources/database/sqlc-config.md) - Auto-resolved from sqlc-config


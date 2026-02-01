## Table of Contents

- [Collections & Playlists](#collections-playlists)
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
      - [POST /api/v1/collections](#post-apiv1collections)
      - [GET /api/v1/collections](#get-apiv1collections)
      - [GET /api/v1/collections/{id}](#get-apiv1collectionsid)
      - [PUT /api/v1/collections/{id}](#put-apiv1collectionsid)
      - [DELETE /api/v1/collections/{id}](#delete-apiv1collectionsid)
      - [POST /api/v1/collections/{id}/items](#post-apiv1collectionsiditems)
      - [DELETE /api/v1/collections/{id}/items/{item_id}](#delete-apiv1collectionsiditemsitem_id)
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
  - name: ogen OpenAPI Generator
    url: ../../../sources/tooling/ogen.md
    note: Auto-resolved from ogen
design_refs:
  - title: features/shared
    path: features/shared/INDEX.md
  - title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
  - title: LIBRARY_TYPES
    path: features/shared/LIBRARY_TYPES.md
---

# Collections & Playlists


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: feature


> Content module for Movies, TV Shows, Music, All media types

> User-created collections and playlists for organizing media across content types

Collections allow users to group media items together:
- **Manual Collections**: User-curated lists (e.g., "Best Sci-Fi Movies")
- **Smart Collections**: Filter-based dynamic lists (e.g., "All 4K Movies")
- **Cross-Type Collections**: Mix movies, TV shows, music in one collection
- **User-Specific**: Each user can create their own collections
- **Shared Collections**: Optional sharing between users


---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | ✅ | - |
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
internal/content/collections/
├── module.go              # fx module definition
├── repository.go          # Database operations
├── service.go             # Business logic
├── handler.go             # HTTP handlers (ogen)
├── types.go               # Domain types
└── collections_test.go
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
#### POST /api/v1/collections

Create a new collection
**Request**:
```json
{
  "name": "Best Sci-Fi",
  "description": "My favorite science fiction movies",
  "type": "manual",
  "is_public": false
}

```
**Response**:
```json
{
  "id": "uuid-123",
  "name": "Best Sci-Fi",
  "type": "manual",
  "is_public": false,
  "created_at": "2026-01-31T12:00:00Z"
}

```

---
#### GET /api/v1/collections

List user's collections
**Request**:
```json
{}
```
**Response**:
```json
{
  "collections": [
    {"id": "uuid-123", "name": "Best Sci-Fi", "type": "manual", "item_count": 12}
  ]
}

```

---
#### GET /api/v1/collections/{id}

Get collection details with items
**Request**:
```json
{}
```
**Response**:
```json
{
  "id": "uuid-123",
  "name": "Best Sci-Fi",
  "items": [
    {"type": "movie", "id": "uuid-456", "title": "Inception"}
  ]
}

```

---
#### PUT /api/v1/collections/{id}

Update collection metadata
**Request**:
```json
{
  "name": "Updated Name",
  "description": "New description"
}

```
**Response**:
```json
{
  "id": "uuid-123",
  "name": "Updated Name",
  "updated_at": "2026-01-31T13:00:00Z"
}

```

---
#### DELETE /api/v1/collections/{id}

Delete collection
**Request**:
```json
{}
```
**Response**:
```json
204 No Content
```

---
#### POST /api/v1/collections/{id}/items

Add items to collection
**Request**:
```json
{
  "items": [
    {"type": "movie", "id": "uuid-123"},
    {"type": "tvshow", "id": "uuid-456"}
  ]
}

```
**Response**:
```json
{
  "added": 2,
  "item_count": 14
}

```

---
#### DELETE /api/v1/collections/{id}/items/{item_id}

Remove item from collection
**Request**:
```json
{}
```
**Response**:
```json
204 No Content
```

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
- [features/shared](features/shared/INDEX.md)
- [01_ARCHITECTURE](architecture/01_ARCHITECTURE.md)
- [LIBRARY_TYPES](features/shared/LIBRARY_TYPES.md)

### External Sources
- [Uber fx](../../../sources/tooling/fx.md) - Auto-resolved from fx
- [ogen OpenAPI Generator](../../../sources/tooling/ogen.md) - Auto-resolved from ogen


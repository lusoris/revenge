## Table of Contents

- [Book Module](#book-module)
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
      - [GET /api/v1/books](#get-apiv1books)
      - [GET /api/v1/books/:id](#get-apiv1booksid)
      - [GET /api/v1/books/:id/read](#get-apiv1booksidread)
      - [GET /api/v1/books/:id/download](#get-apiv1booksiddownload)
      - [GET /api/v1/books/:id/progress](#get-apiv1booksidprogress)
      - [PUT /api/v1/books/:id/progress](#put-apiv1booksidprogress)
      - [POST /api/v1/books/:id/bookmarks](#post-apiv1booksidbookmarks)
      - [GET /api/v1/books/:id/bookmarks](#get-apiv1booksidbookmarks)
      - [POST /api/v1/books/:id/highlights](#post-apiv1booksidhighlights)
      - [GET /api/v1/books/:id/highlights](#get-apiv1booksidhighlights)
      - [GET /api/v1/books/authors](#get-apiv1booksauthors)
      - [GET /api/v1/books/series](#get-apiv1booksseries)
      - [GET /api/v1/books/collections](#get-apiv1bookscollections)
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
  - name: Google Books API
    url: ../../../sources/apis/google-books.md
    note: Auto-resolved from google-books
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
  - title: features/book
    path: features/book.md
  - title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# Book Module


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: feature


> Content module for Books, Authors, Series

> Book/eBook content management with metadata enrichment from OpenLibrary and Goodreads

Complete eBook library:
- **Metadata Sources**: OpenLibrary (primary), Google Books, Goodreads, Hardcover
- **Supported Formats**: EPUB, PDF, MOBI, AZW3, CBZ (comics)
- **Reading Features**: Web reader, progress tracking, bookmarks, highlights
- **Collections**: Reading lists, series tracking, genre collections
- **Sync**: Multi-device reading progress synchronization

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete book module design |
| Sources | ✅ | All book APIs documented |
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
internal/content/book/
├── module.go              # fx module definition
├── repository.go          # Database operations
├── service.go             # Business logic
├── handler.go             # HTTP handlers (ogen)
├── types.go               # Domain types
└── book_test.go
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
#### GET /api/v1/books

List all books with pagination and filters

---
#### GET /api/v1/books/:id

Get book details by ID

---
#### GET /api/v1/books/:id/read

Get book content for web reader

---
#### GET /api/v1/books/:id/download

Download book file in original format

---
#### GET /api/v1/books/:id/progress

Get user reading progress for a book

---
#### PUT /api/v1/books/:id/progress

Update user reading progress

---
#### POST /api/v1/books/:id/bookmarks

Create a bookmark at current position

---
#### GET /api/v1/books/:id/bookmarks

Get all user bookmarks for a book

---
#### POST /api/v1/books/:id/highlights

Create a text highlight with optional note

---
#### GET /api/v1/books/:id/highlights

Get all user highlights for a book

---
#### GET /api/v1/books/authors

List all book authors

---
#### GET /api/v1/books/series

List all book series

---
#### GET /api/v1/books/collections

List all book collections

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
- [features/book](features/book.md)
- [01_ARCHITECTURE](../../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Uber fx](../../../sources/tooling/fx.md) - Auto-resolved from fx
- [Google Books API](../../../sources/apis/google-books.md) - Auto-resolved from google-books
- [ogen OpenAPI Generator](../../../sources/tooling/ogen.md) - Auto-resolved from ogen
- [Open Library API](../../../sources/apis/openlibrary.md) - Auto-resolved from openlibrary
- [pgx PostgreSQL Driver](../../../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../../../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../../../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [River Job Queue](../../../sources/tooling/river.md) - Auto-resolved from river
- [sqlc](../../../sources/database/sqlc.md) - Auto-resolved from sqlc
- [sqlc Configuration](../../../sources/database/sqlc-config.md) - Auto-resolved from sqlc-config


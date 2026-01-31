# Book Module

<!-- SOURCES: fx, google-books, ogen, openlibrary, pgx, postgresql-arrays, postgresql-json, river, sqlc, sqlc-config -->

<!-- DESIGN: features/book, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> Book/eBook content management with metadata enrichment from OpenLibrary and Goodreads


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Developer Resources](#developer-resources)
- [Overview](#overview)
- [Architecture](#architecture)
- [Files (Planned)](#files-planned)
- [Entities (Planned)](#entities-planned)
  - [Book](#book)
  - [Author](#author)
  - [Series](#series)
- [Supported Formats](#supported-formats)
- [Metadata Priority Chain](#metadata-priority-chain)
- [Arr Integration](#arr-integration)
- [Reading Progress](#reading-progress)
- [Database Schema (Planned)](#database-schema-planned)
- [API Endpoints (Planned)](#api-endpoints-planned)
- [Implementation Checklist](#implementation-checklist)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Documents](#related-documents)

<!-- TOC-END -->

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | 🟡 | Scaffold - needs detailed spec |
| Sources | 🔴 | OpenLibrary, Goodreads, Google Books API docs needed |
| Instructions | 🔴 |  |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |**Location**: `internal/content/book/`

---

## Developer Resources

| Source | URL | Purpose |
|--------|-----|---------|
| OpenLibrary API | [openlibrary.org/developers](https://openlibrary.org/developers/api) | Primary book metadata |
| Goodreads | Via Hardcover | Community ratings, reviews |
| Google Books API | [developers.google.com/books](https://developers.google.com/books) | Alternative metadata |
| Chaptarr (Readarr) | See [integrations/servarr/CHAPTARR.md](../../integrations/servarr/CHAPTARR.md) | Servarr integration |

---

## Overview

The Book module provides complete eBook library management:

- Entity definitions (Book, Author, Publisher, Series, etc.)
- Repository pattern with PostgreSQL implementation
- Service layer with otter caching
- Background jobs for metadata enrichment via River
- User data (reading progress, ratings, favorites)
- Multiple format support (EPUB, PDF, MOBI, etc.)

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                       API Layer                              │
│                    (ogen handlers)                           │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                     Book Service                             │
│   - Local cache (otter)                                      │
│   - Business logic                                           │
│   - Reading progress                                         │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                    Repository Layer                          │
│   - PostgreSQL queries (sqlc)                                │
│   - User data (progress, annotations)                        │
│   - Relations (authors, series, genres)                      │
└─────────────────────────────────────────────────────────────┘
```

---

## Files (Planned)

| File | Description |
|------|-------------|
| `entity.go` | Domain entities (Book, Author, Series, etc.) |
| `repository.go` | Repository interface definition |
| `repository_pg.go` | PostgreSQL implementation |
| `service.go` | Business logic with caching |
| `jobs.go` | River background jobs |
| `metadata_provider.go` | OpenLibrary/Goodreads interface |
| `module.go` | fx dependency injection |

---

## Entities (Planned)

### Book

```go
type Book struct {
    shared.ContentEntity

    Title          string
    Subtitle       string
    ISBN10         *string
    ISBN13         *string
    OpenLibraryID  *string
    GoodreadsID    *string

    // Content
    PageCount      int
    FilePath       string
    Format         string // epub, pdf, mobi, azw3
    FileSizeBytes  int64

    // Metadata
    Publisher      string
    PublishDate    *time.Time
    Language       string
    Description    string

    // Relations
    SeriesID       *uuid.UUID
    SeriesPosition *float32
}
```

### Author

```go
type Author struct {
    shared.ContentEntity

    Name          string
    SortName      string
    OpenLibraryID *string
    GoodreadsID   *string
    Biography     string
    BirthDate     *time.Time
    DeathDate     *time.Time
    Website       string
}
```

### Series

```go
type Series struct {
    shared.ContentEntity

    Name          string
    OpenLibraryID *string
    GoodreadsID   *string
    Description   string
    BookCount     int
}
```

---

## Supported Formats

| Format | Extension | Reader Support |
|--------|-----------|----------------|
| EPUB | `.epub` | Full (reflowable) |
| PDF | `.pdf` | Full (fixed layout) |
| MOBI | `.mobi` | Convert to EPUB |
| AZW3 | `.azw3` | Convert to EPUB |
| CBZ/CBR | `.cbz`, `.cbr` | See Comics module |

---

## Metadata Priority Chain

See [00_SOURCE_OF_TRUTH.md](../../00_SOURCE_OF_TRUTH.md) for the core metadata priority principle.

```
1. LOCAL CACHE     → First, instant UI display
2. CHAPTARR        → Arr-first metadata (Readarr API)
3. OPENLIBRARY     → Primary book metadata
4. GOODREADS       → Community ratings/reviews (via Hardcover)
5. GOOGLE_BOOKS    → Alternative metadata source
```

---

## Arr Integration

**Primary**: Chaptarr (Readarr-compatible)

See [integrations/servarr/CHAPTARR.md](../../integrations/servarr/CHAPTARR.md) for:
- Webhook handling
- Import notifications
- Library sync patterns

---

## Reading Progress

Books require page-based or percentage-based progress tracking:

- Current page / total pages
- Percentage complete
- Last read timestamp
- Reading speed statistics
- Annotations and highlights (optional)

---

## Database Schema (Planned)

Tables in `public` schema:

- `books` - Book entities
- `book_authors` - Author entities (shared with audiobooks)
- `book_publishers` - Publisher entities
- `book_series` - Series entities
- `book_author` - Author relationships
- `book_genres` - Genre mappings
- `user_book_progress` - Reading progress
- `user_book_annotations` - Highlights, notes

---

## API Endpoints (Planned)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/books` | List books |
| GET | `/api/v1/books/{id}` | Get book details |
| GET | `/api/v1/books/{id}/download` | Download book file |
| GET | `/api/v1/books/authors` | List authors |
| GET | `/api/v1/books/series` | List series |
| PUT | `/api/v1/books/{id}/progress` | Update reading progress |
| GET | `/api/v1/books/{id}/annotations` | Get annotations |
| POST | `/api/v1/books/{id}/annotations` | Add annotation |

---

## Implementation Checklist

- [ ] Define entity structs in `entity.go`
- [ ] Create repository interface
- [ ] Implement PostgreSQL repository
- [ ] Create database migrations
- [ ] Implement service layer with caching
- [ ] Add River jobs for metadata enrichment
- [ ] Integrate OpenLibrary provider
- [ ] Integrate Goodreads provider (via Hardcover)
- [ ] Add Chaptarr webhook handlers
- [ ] Implement format detection
- [ ] Implement reading progress tracking
- [ ] Write unit tests
- [ ] Write integration tests

---


## Related Documents

- [OpenLibrary Integration](../../integrations/metadata/books/OPENLIBRARY.md)
- [Goodreads Integration](../../integrations/metadata/books/GOODREADS.md)
- [Hardcover Integration](../../integrations/metadata/books/HARDCOVER.md)
- [Chaptarr Integration](../../integrations/servarr/CHAPTARR.md)
- [Audiobook Module](../audiobook/AUDIOBOOK_MODULE.md)

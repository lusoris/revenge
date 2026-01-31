# Audiobook Module

> Audiobook content management with metadata enrichment from Audnexus and OpenLibrary


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Developer Resources](#developer-resources)
- [Overview](#overview)
- [Architecture](#architecture)
- [Files (Planned)](#files-planned)
- [Entities (Planned)](#entities-planned)
  - [Audiobook](#audiobook)
  - [Author](#author)
  - [Narrator](#narrator)
  - [Chapter](#chapter)
- [Metadata Priority Chain](#metadata-priority-chain)
- [Arr Integration](#arr-integration)
- [Progress Tracking](#progress-tracking)
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
| Sources | 🔴 | Audnexus, OpenLibrary, Audible API docs needed |
| Instructions | 🔴 |  |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |
**Location**: `internal/content/audiobook/`

---

## Developer Resources

| Source | URL | Purpose |
|--------|-----|---------|
| Audnexus API | [audnex.us](https://audnex.us/) | Primary audiobook metadata |
| OpenLibrary API | [openlibrary.org/developers](https://openlibrary.org/developers/api) | Book metadata, covers |
| Chaptarr (Readarr) | See [integrations/servarr/CHAPTARR.md](../../integrations/servarr/CHAPTARR.md) | Servarr integration |

---

## Overview

The Audiobook module provides complete audiobook library management:

- Entity definitions (Audiobook, Author, Narrator, Series, etc.)
- Repository pattern with PostgreSQL implementation
- Service layer with otter caching
- Background jobs for metadata enrichment via River
- User data (progress tracking, ratings, favorites)
- Chapter navigation and bookmarks

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                       API Layer                              │
│                    (ogen handlers)                           │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                  Audiobook Service                           │
│   - Local cache (otter)                                      │
│   - Business logic                                           │
│   - Progress tracking                                        │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                    Repository Layer                          │
│   - PostgreSQL queries (sqlc)                                │
│   - User data (progress, bookmarks)                          │
│   - Relations (authors, narrators, series)                   │
└─────────────────────────────────────────────────────────────┘
```

---

## Files (Planned)

| File | Description |
|------|-------------|
| `entity.go` | Domain entities (Audiobook, Author, Narrator, etc.) |
| `repository.go` | Repository interface definition |
| `repository_pg.go` | PostgreSQL implementation |
| `service.go` | Business logic with caching |
| `jobs.go` | River background jobs |
| `metadata_provider.go` | Audnexus/OpenLibrary interface |
| `module.go` | fx dependency injection |

---

## Entities (Planned)

### Audiobook

```go
type Audiobook struct {
    shared.ContentEntity

    Title          string
    Subtitle       string
    ASIN           *string
    ISBN           *string
    OpenLibraryID  *string

    // Content
    DurationMs     int64
    ChapterCount   int
    FilePath       string
    Container      string

    // Metadata
    Publisher      string
    PublishDate    *time.Time
    Language       string
    Abridged       bool
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
    ASIN          *string
    OpenLibraryID *string
    Biography     string
    BirthDate     *time.Time
    DeathDate     *time.Time
}
```

### Narrator

```go
type Narrator struct {
    shared.ContentEntity

    Name     string
    SortName string
    ASIN     *string
}
```

### Chapter

```go
type Chapter struct {
    ID          uuid.UUID
    AudiobookID uuid.UUID
    Title       string
    StartMs     int64
    EndMs       int64
    ChapterNum  int
}
```

---

## Metadata Priority Chain

See [00_SOURCE_OF_TRUTH.md](../../00_SOURCE_OF_TRUTH.md) for the core metadata priority principle.

```
1. LOCAL CACHE     → First, instant UI display
2. CHAPTARR        → Arr-first metadata (Readarr API)
3. AUDNEXUS        → Primary audiobook metadata
4. OPENLIBRARY     → Fallback book metadata
5. GOODREADS       → Additional ratings/reviews
```

---

## Arr Integration

**Primary**: Chaptarr (Readarr-compatible)

See [integrations/servarr/CHAPTARR.md](../../integrations/servarr/CHAPTARR.md) for:
- Webhook handling
- Import notifications
- Library sync patterns

---

## Progress Tracking

Audiobooks require precise progress tracking:

- Current chapter
- Position within chapter (ms)
- Playback speed preference
- Bookmarks with notes

---

## Database Schema (Planned)

Tables in `public` schema:

- `audiobooks` - Audiobook entities
- `authors` - Author entities
- `narrators` - Narrator entities
- `series` - Series entities
- `chapters` - Chapter markers
- `audiobook_author` - Author relationships
- `audiobook_narrator` - Narrator relationships
- `user_audiobook_progress` - Progress tracking
- `user_audiobook_bookmarks` - User bookmarks

---

## API Endpoints (Planned)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/audiobooks` | List audiobooks |
| GET | `/api/v1/audiobooks/{id}` | Get audiobook details |
| GET | `/api/v1/audiobooks/{id}/chapters` | Get chapters |
| GET | `/api/v1/audiobooks/authors` | List authors |
| GET | `/api/v1/audiobooks/narrators` | List narrators |
| PUT | `/api/v1/audiobooks/{id}/progress` | Update progress |
| POST | `/api/v1/audiobooks/{id}/bookmarks` | Add bookmark |

---

## Implementation Checklist

- [ ] Define entity structs in `entity.go`
- [ ] Create repository interface
- [ ] Implement PostgreSQL repository
- [ ] Create database migrations
- [ ] Implement service layer with caching
- [ ] Add River jobs for metadata enrichment
- [ ] Integrate Audnexus provider
- [ ] Integrate OpenLibrary provider
- [ ] Add Chaptarr webhook handlers
- [ ] Implement chapter extraction
- [ ] Implement progress tracking
- [ ] Write unit tests
- [ ] Write integration tests

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [Audnexus API](https://api.audnex.us/) | [Local](../../../sources/apis/audnexus.md) |
| [Open Library API](https://openlibrary.org/developers/api) | [Local](../../../sources/apis/openlibrary.md) |
| [PostgreSQL Arrays](https://www.postgresql.org/docs/current/arrays.html) | [Local](../../../sources/database/postgresql-arrays.md) |
| [PostgreSQL JSON Functions](https://www.postgresql.org/docs/current/functions-json.html) | [Local](../../../sources/database/postgresql-json.md) |
| [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) | [Local](../../../sources/tooling/river.md) |
| [Uber fx](https://pkg.go.dev/go.uber.org/fx) | [Local](../../../sources/tooling/fx.md) |
| [ogen OpenAPI Generator](https://pkg.go.dev/github.com/ogen-go/ogen) | [Local](../../../sources/tooling/ogen.md) |
| [pgx PostgreSQL Driver](https://pkg.go.dev/github.com/jackc/pgx/v5) | [Local](../../../sources/database/pgx.md) |
| [sqlc](https://docs.sqlc.dev/en/stable/) | [Local](../../../sources/database/sqlc.md) |
| [sqlc Configuration](https://docs.sqlc.dev/en/stable/reference/config.html) | [Local](../../../sources/database/sqlc-config.md) |

<!-- SOURCE-BREADCRUMBS-END -->

## Related Documents

- [Audible Integration](../../integrations/metadata/books/AUDIBLE.md)
- [OpenLibrary Integration](../../integrations/metadata/books/OPENLIBRARY.md)
- [Chaptarr Integration](../../integrations/servarr/CHAPTARR.md)

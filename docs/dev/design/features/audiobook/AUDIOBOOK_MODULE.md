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
  - title: CHAPTARR (PRIMARY metadata + downloads)
    path: ../../integrations/servarr/CHAPTARR.md
  - title: AUDNEXUS (supplementary audiobook metadata)
    path: ../../integrations/metadata/AUDNEXUS.md
  - title: OPENLIBRARY (supplemental book metadata)
    path: ../../integrations/metadata/OPENLIBRARY.md
  - title: GOODREADS (ratings + reviews)
    path: ../../integrations/metadata/GOODREADS.md
---

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


# Audiobook Module


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: feature


> Content module for Books, Authors, Series

> Audiobook content management with metadata enrichment from Chaptarr and external providers

Complete audiobook library:
- **Metadata Sources**: Chaptarr (PRIMARY - aggregates Audnexus/OpenLibrary), with direct APIs as supplementary
- **Download Automation**: Chaptarr integration for automated audiobook management
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

```
internal/content/audiobook/
├── module.go              # fx.Module with all providers
├── repository.go          # Database layer
├── repository_test.go     # Repository tests (testcontainers)
├── service.go             # Business logic
├── service_test.go        # Service tests (mocks)
├── handler.go             # HTTP handlers
├── handler_test.go        # Handler tests (httptest)
├── types.go               # Domain types
├── cache.go               # Caching logic
├── cache_test.go          # Cache tests
├── chapters/
│   ├── extractor.go       # Chapter extraction from M4B/MP3
│   ├── extractor_test.go  # Chapter extraction tests
│   └── generator.go       # Chapter generation for multi-file
├── streaming/
│   ├── transcoder.go      # FFmpeg transcoding logic
│   ├── hls.go             # HLS manifest generation
│   └── chapters.go        # Chapter-aware HLS segments
├── metadata/
│   ├── provider.go        # Interface: MetadataProvider
│   ├── chaptarr.go        # Chaptarr API integration
│   ├── audnexus.go        # Audnexus API integration
│   ├── openlibrary.go     # OpenLibrary API integration
│   └── enricher.go        # Enrichment orchestration
└── progress/
    ├── tracker.go         # Playback progress tracking
    ├── sync.go            # Multi-device sync logic
    └── statistics.go      # Listening statistics

migrations/
└── audiobooks/
    ├── 001_audiobooks.sql # Audiobooks schema
    ├── 002_chapters.sql   # Chapters schema
    └── 003_progress.sql   # Progress tracking schema

api/
└── openapi.yaml           # OpenAPI spec (audiobooks/* endpoints)
```


### Key Interfaces

```go
// Repository defines database operations for audiobooks
type Repository interface {
    // Audiobook CRUD
    GetAudiobook(ctx context.Context, id uuid.UUID) (*Audiobook, error)
    ListAudiobooks(ctx context.Context, filters ListFilters) ([]Audiobook, error)
    CreateAudiobook(ctx context.Context, audiobook *Audiobook) error
    UpdateAudiobook(ctx context.Context, audiobook *Audiobook) error
    DeleteAudiobook(ctx context.Context, id uuid.UUID) error

    // Chapter operations
    GetChapters(ctx context.Context, audiobookID uuid.UUID) ([]Chapter, error)
    CreateChapter(ctx context.Context, chapter *Chapter) error
    UpdateChapter(ctx context.Context, chapter *Chapter) error

    // Author and narrator operations
    GetAuthor(ctx context.Context, id uuid.UUID) (*Author, error)
    GetNarrator(ctx context.Context, id uuid.UUID) (*Narrator, error)
    ListNarrators(ctx context.Context) ([]Narrator, error)

    // Progress tracking
    GetProgress(ctx context.Context, userID, audiobookID uuid.UUID) (*PlaybackProgress, error)
    UpdateProgress(ctx context.Context, progress *PlaybackProgress) error
    GetListeningStatistics(ctx context.Context, userID uuid.UUID) (*Statistics, error)

    // Bookmarks
    CreateBookmark(ctx context.Context, bookmark *Bookmark) error
    ListBookmarks(ctx context.Context, userID, audiobookID uuid.UUID) ([]Bookmark, error)
}

// Service defines business logic for audiobooks
type Service interface {
    // Audiobook operations
    GetAudiobook(ctx context.Context, id uuid.UUID) (*Audiobook, error)
    SearchAudiobooks(ctx context.Context, query string, filters SearchFilters) ([]Audiobook, error)
    EnrichAudiobook(ctx context.Context, id uuid.UUID) error

    // Streaming operations
    GetStreamURL(ctx context.Context, audiobookID uuid.UUID, format string) (string, error)
    GetChapters(ctx context.Context, audiobookID uuid.UUID) ([]Chapter, error)

    // Progress operations
    UpdateProgress(ctx context.Context, userID, audiobookID uuid.UUID, progress ProgressUpdate) error
    GetResumePoint(ctx context.Context, userID, audiobookID uuid.UUID) (*ResumePoint, error)
}

// MetadataProvider fetches audiobook metadata from external sources
type MetadataProvider interface {
    GetAudiobookByASIN(ctx context.Context, asin string) (*AudiobookMetadata, error)
    GetAudiobookByISBN(ctx context.Context, isbn string) (*AudiobookMetadata, error)
    GetAuthorByID(ctx context.Context, providerID string) (*AuthorMetadata, error)
    GetNarratorByID(ctx context.Context, providerID string) (*NarratorMetadata, error)
    SearchAudiobooks(ctx context.Context, query string) ([]AudiobookMetadata, error)
}

// ChapterExtractor extracts chapter information from audiobook files
type ChapterExtractor interface {
    // ExtractFromM4B extracts chapters from M4B file metadata
    ExtractFromM4B(ctx context.Context, filePath string) ([]Chapter, error)

    // GenerateFromMultiFile creates chapters from multi-file structure
    GenerateFromMultiFile(ctx context.Context, files []string) ([]Chapter, error)

    // SupportedFormats returns formats this extractor can handle
    SupportedFormats() []string
}

// ProgressTracker manages playback progress synchronization
type ProgressTracker interface {
    // UpdateProgress updates user's playback position
    UpdateProgress(ctx context.Context, update ProgressUpdate) error

    // GetProgress retrieves current playback position
    GetProgress(ctx context.Context, userID, audiobookID uuid.UUID) (*PlaybackProgress, error)

    // SyncProgress synchronizes progress across devices
    SyncProgress(ctx context.Context, userID uuid.UUID) error

    // RecordListeningTime updates total listening time statistics
    RecordListeningTime(ctx context.Context, userID, audiobookID uuid.UUID, seconds int) error
}
```


### Dependencies
**Go Dependencies**:
- `github.com/jackc/pgx/v5/pgxpool` - PostgreSQL connection pool
- `github.com/google/uuid` - UUID generation
- `github.com/maypok86/otter` - In-memory cache
- `github.com/asticode/go-astiav` - FFmpeg bindings for audio processing
- `github.com/dhowden/tag` - Audio file metadata reading (for M4B chapters)
- `github.com/go-resty/resty/v2` - HTTP client for external APIs
- `go.uber.org/fx` - Dependency injection
- `github.com/riverqueue/river` - Background job queue
- `golang.org/x/net/proxy` - SOCKS5 proxy support for external metadata calls

**External APIs** (priority order):
- **Chaptarr API** - PRIMARY metadata source (local Audnexus/OpenLibrary cache) + download automation
- **Audnexus API** - Supplementary metadata (via proxy/VPN when Chaptarr lacks data)
- **OpenLibrary API** - Supplemental book metadata (via proxy/VPN)
- **Goodreads API** - Ratings and reviews (via proxy/VPN)

**External Tools**:
- FFmpeg 7.0+ - Audio transcoding and chapter extraction

**Database**:
- PostgreSQL 18+ with trigram extension for fuzzy search






## Configuration
### Environment Variables

**Environment Variables**:
- `REVENGE_AUDIOBOOK_CACHE_TTL` - Cache TTL duration (default: 20m)
- `REVENGE_AUDIOBOOK_CACHE_SIZE` - Cache size in MB (default: 150)
- `REVENGE_CHAPTARR_URL` - Chaptarr instance URL (optional)
- `REVENGE_CHAPTARR_API_KEY` - Chaptarr API key (optional)
- `REVENGE_METADATA_AUDNEXUS_RATE_LIMIT` - Rate limit per second (default: 10)
- `REVENGE_AUDIOBOOK_STREAMING_QUALITY` - Streaming quality preset (low, medium, high)
- `REVENGE_AUDIOBOOK_PROGRESS_SYNC_INTERVAL` - Progress sync interval in seconds (default: 15)
- `REVENGE_AUDIOBOOK_CHAPTER_EXTRACTION_ENABLED` - Enable automatic chapter extraction (default: true)


### Config Keys
**config.yaml keys**:
```yaml
audiobook:
  cache:
    ttl: 20m
    size_mb: 150

  metadata:
    priority:
      - chaptarr     # PRIMARY: Local Audnexus/OpenLibrary cache
      - audnexus     # Supplementary: Direct API (via proxy/VPN)
      - openlibrary  # Supplemental (via proxy/VPN)

    chaptarr:
      enabled: true       # Should be enabled for PRIMARY metadata
      url: ${REVENGE_CHAPTARR_URL}
      api_key: ${REVENGE_CHAPTARR_API_KEY}
      sync_interval: 30m

    audnexus:
      api_url: https://api.audnex.us
      rate_limit: 10  # Requests per second
      timeout: 15s
      proxy: tor  # Route through proxy/VPN (see HTTP_CLIENT service)

    openlibrary:
      rate_limit: 5
      enabled: true
      proxy: tor  # Route through proxy/VPN

  streaming:
    quality: high  # low, medium, high
    formats:
      - aac: {bitrate: 128, codec: aac}
      - opus: {bitrate: 96, codec: libopus}
    cache_segments: true
    segment_duration: 10  # seconds

  chapters:
    auto_extract: true
    multi_file_detection: true
    cache_chapter_data: true

  progress:
    sync_interval: 15s
    autosave_interval: 5s
    track_listening_time: true
    sync_to_chaptarr: false

  playback:
    speed_options: [0.5, 0.75, 1.0, 1.25, 1.5, 1.75, 2.0, 2.5, 3.0]
    sleep_timer_options: [5, 10, 15, 30, 45, 60, "end_of_chapter"]
    offline_downloads_enabled: true
```



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
- [CHAPTARR (PRIMARY metadata + downloads)](../../integrations/servarr/CHAPTARR.md)
- [AUDNEXUS (supplementary audiobook metadata)](../../integrations/metadata/AUDNEXUS.md)
- [OPENLIBRARY (supplemental book metadata)](../../integrations/metadata/OPENLIBRARY.md)
- [GOODREADS (ratings + reviews)](../../integrations/metadata/GOODREADS.md)

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


# TODO A12: Shared Metadata Service

**Phase**: A12
**Priority**: P1 (High - Enables clean metadata architecture)
**Effort**: 24-32 hours
**Status**: ✅ Complete (7/7 tasks complete)
**Dependencies**: A11 (TV Module - for testing)
**Created**: 2026-02-05
**Updated**: 2026-02-05

---

## Overview

**Goal**: Create a shared metadata service that abstracts all external metadata providers (TMDb, TVDb, etc.) from content modules.

**Current Problem**:
- TMDb adapter duplicated in movie and tvshow modules
- Content modules know about external APIs
- Hard to add new sources (Fanart.tv, OMDb, MusicBrainz)
- No centralized rate limiting, caching, fallback

**Solution**:
- Shared metadata service in `internal/service/metadata/`
- Content modules call `metadataService.RefreshMovie()` - don't know about providers
- Centralized rate limiting, caching, provider fallback
- Easy to add new sources later

**Benefits**:
- Content modules stay clean
- Single place for all metadata logic
- Provider-agnostic enrichment
- Cross-content enrichment (TVDb person data for movies)
- Centralized error handling, retries, fallback

---

## Decision Log

| Topic | Decision | Reason |
|-------|----------|--------|
| Architecture | Option B (Shared Service) | Content modules don't know about providers |
| Location | `internal/service/metadata/` | Follows service pattern |
| Providers | Interface-based | Easy to add new sources |
| Jobs | River integration | Async metadata updates |
| Caching | Provider-level caching | Reduce API calls |

---

## Tasks

### A12.1: Provider Interface & Types 🔴 CRITICAL

**Priority**: P0
**Effort**: 4-6h
**Status**: Pending

#### A12.1.1: Define Provider Interface

**Location**: `internal/service/metadata/provider.go`

```go
// Provider represents an external metadata source
type Provider interface {
    Name() string
    Priority() int  // Higher = preferred

    // Movies
    SearchMovie(ctx context.Context, query string, year *int) ([]MovieSearchResult, error)
    GetMovie(ctx context.Context, id string) (*MovieMetadata, error)
    GetMovieCredits(ctx context.Context, id string) (*Credits, error)
    GetMovieImages(ctx context.Context, id string) (*Images, error)

    // TV Shows
    SearchSeries(ctx context.Context, query string, year *int) ([]SeriesSearchResult, error)
    GetSeries(ctx context.Context, id string) (*SeriesMetadata, error)
    GetSeason(ctx context.Context, seriesID string, seasonNum int) (*SeasonMetadata, error)
    GetEpisode(ctx context.Context, seriesID string, seasonNum, episodeNum int) (*EpisodeMetadata, error)
    GetSeriesCredits(ctx context.Context, id string) (*Credits, error)
    GetSeriesImages(ctx context.Context, id string) (*Images, error)

    // People
    GetPerson(ctx context.Context, id string) (*PersonMetadata, error)
    GetPersonCredits(ctx context.Context, id string) (*PersonCredits, error)

    // Configuration
    GetImageBaseURL() string
    SupportsLanguage(lang string) bool
}
```

#### A12.1.2: Define Metadata Types

**Location**: `internal/service/metadata/types.go`

Provider-agnostic metadata types:
- `MovieMetadata` - Movie with all fields
- `SeriesMetadata` - Series with all fields
- `SeasonMetadata` - Season with episodes
- `EpisodeMetadata` - Episode details
- `PersonMetadata` - Person info
- `Credits` - Cast + Crew
- `Images` - Posters, backdrops, stills
- `SearchResult` - Generic search result

---

### A12.2: TMDb Provider 🔴 CRITICAL

**Priority**: P0
**Effort**: 8-10h
**Status**: Pending

#### A12.2.1: TMDb Client

**Location**: `internal/service/metadata/providers/tmdb/client.go`

- HTTP client with resty
- Rate limiting (40 req/10s for TMDb)
- Response caching (1h for details, 15min for search)
- Multi-language support
- Error handling with retries

#### A12.2.2: TMDb Provider Implementation

**Location**: `internal/service/metadata/providers/tmdb/provider.go`

Implement Provider interface:
- All movie endpoints
- All TV endpoints
- Person endpoints
- Image configuration
- Language support

---

### A12.3: TVDb Provider

**Priority**: P1
**Effort**: 6-8h
**Status**: Pending

#### A12.3.1: TVDb Client

**Location**: `internal/service/metadata/providers/tvdb/client.go`

- Auth token management (JWT)
- Rate limiting
- Response caching

#### A12.3.2: TVDb Provider Implementation

**Location**: `internal/service/metadata/providers/tvdb/provider.go`

- TV-focused endpoints
- Artwork endpoints (often better than TMDb)
- Episode ordering variants

---

### A12.4: Metadata Service 🔴 CRITICAL

**Priority**: P0
**Effort**: 6-8h
**Status**: Pending

#### A12.4.1: Service Interface

**Location**: `internal/service/metadata/service.go`

```go
type Service interface {
    // Movies
    SearchMovie(ctx context.Context, query string, year *int, lang string) ([]MovieSearchResult, error)
    GetMovieMetadata(ctx context.Context, tmdbID int32, languages []string) (*MovieMetadata, error)
    RefreshMovie(ctx context.Context, movieID uuid.UUID) error  // Async job

    // TV Shows
    SearchSeries(ctx context.Context, query string, year *int, lang string) ([]SeriesSearchResult, error)
    GetSeriesMetadata(ctx context.Context, tmdbID int32, languages []string) (*SeriesMetadata, error)
    GetSeasonMetadata(ctx context.Context, seriesTmdbID int32, seasonNum int, languages []string) (*SeasonMetadata, error)
    GetEpisodeMetadata(ctx context.Context, seriesTmdbID int32, seasonNum, episodeNum int, languages []string) (*EpisodeMetadata, error)
    RefreshSeries(ctx context.Context, seriesID uuid.UUID) error  // Async job

    // People
    GetPersonMetadata(ctx context.Context, tmdbPersonID int32) (*PersonMetadata, error)

    // Images
    GetImageURL(ctx context.Context, path string, size string) string

    // Enrichment (combines multiple providers)
    EnrichMovie(ctx context.Context, movieID uuid.UUID) error
    EnrichSeries(ctx context.Context, seriesID uuid.UUID) error
}
```

#### A12.4.2: Provider Aggregation

- Try primary provider first (TMDb)
- Fallback to secondary providers
- Merge results for enrichment
- Handle partial failures gracefully

---

### A12.5: Jobs Integration

**Priority**: P1
**Effort**: 4-6h
**Status**: Pending

#### A12.5.1: Metadata Jobs

**Location**: `internal/service/metadata/jobs/`

```go
// Job types for River
type RefreshMovieArgs struct {
    MovieID   uuid.UUID
    Force     bool
    Languages []string
}

type RefreshSeriesArgs struct {
    SeriesID  uuid.UUID
    Force     bool
    Languages []string
    IncludeSeasons  bool
    IncludeEpisodes bool
}

type EnrichContentArgs struct {
    ContentType string  // "movie" or "series"
    ContentID   uuid.UUID
    Providers   []string  // Which providers to use
}
```

---

### A12.6: Refactor Content Modules

**Priority**: P1
**Effort**: 4-6h
**Status**: ✅ Complete

#### A12.6.1: Create Movie Adapter

**Location**: `internal/service/metadata/adapters/movie/adapter.go`

- Created adapter that implements `movie.MetadataProvider` interface
- Uses shared metadata service (no direct TMDb dependency)
- Maps shared metadata types to movie domain types
- Handles multi-language translations and age ratings

#### A12.6.2: TV Module

- TV module doesn't have MetadataProvider interface yet
- Will add tvshow adapter when needed
- Shared service ready for TV operations

---

### A12.7: fx Module & Wiring

**Priority**: P1
**Effort**: 2-3h
**Status**: Pending

**Location**: `internal/service/metadata/module.go`

```go
var Module = fx.Module("metadata",
    fx.Provide(
        NewTMDbProvider,
        NewTVDbProvider,  // Optional
        NewService,
    ),
)
```

---

## Files Created

```
internal/service/metadata/
├── doc.go              # Package documentation ✅
├── provider.go         # Provider interface ✅
├── types.go            # Metadata types ✅
├── service.go          # Service implementation ✅
├── errors.go           # Error types ✅
├── providers/
│   ├── tmdb/
│   │   ├── client.go   # HTTP client with rate limiting ✅
│   │   ├── provider.go # Provider impl ✅
│   │   ├── types.go    # TMDb API types ✅
│   │   └── mapping.go  # TMDb → Metadata mapping ✅
│   └── tvdb/
│       ├── client.go   # HTTP client with JWT auth ✅
│       ├── provider.go # Provider impl ✅
│       ├── types.go    # TVDb API types ✅
│       └── mapping.go  # TVDb → Metadata mapping ✅
├── jobs/
│   ├── refresh.go      # Job argument types ✅
│   └── queue.go        # Queue helper ✅
├── adapters/
│   └── movie/
│       └── adapter.go  # Movie MetadataProvider adapter ✅
└── metadatafx/
    └── module.go       # fx module ✅
```

---

## Testing Strategy

1. **Unit Tests**: Mock providers, test service logic
2. **Integration Tests**: Real API calls (with VCR/cassettes)
3. **Provider Tests**: Each provider independently

---

## Questions

| # | Question | Status | Answer |
|---|----------|--------|--------|
| 1 | TVDb API key required - do we have one? | Open | |
| 2 | Should we cache at service level too? | Open | |
| 3 | Rate limit strategy per provider? | Open | TMDb: 40/10s, TVDb: ? |

---

## Bugs Found

| # | Bug | Severity | Status | Fix |
|---|-----|----------|--------|-----|
| 1 | Job kind mismatch: metadata service enqueued `metadata_refresh_movie` but worker handled `movie_metadata_refresh` | Critical | ✅ Fixed | Updated worker to use `metadatajobs.RefreshMovieArgs` |
| 2 | MetadataProvider interface missing `ClearCache()` method | Medium | ✅ Fixed | Added method to interface |

---

## Progress Tracking

| Task | Status | Notes |
|------|--------|-------|
| A12.1: Provider Interface | ✅ Complete | Provider interfaces, types, errors |
| A12.2: TMDb Provider | ✅ Complete | Full implementation with all endpoints |
| A12.3: TVDb Provider | ✅ Complete | JWT auth, TV focus |
| A12.4: Metadata Service | ✅ Complete | Aggregates providers with fallback |
| A12.5: Jobs Integration | ✅ Complete | River jobs for async refresh |
| A12.6: Refactor Content Modules | ✅ Complete | Movie adapter in service layer; tvshow has no MetadataProvider yet |
| A12.7: fx Module | ✅ Complete | metadatafx package |

---

## Notes

- TMDb has 40 requests per 10 seconds limit
- TVDb requires JWT authentication
- Consider Fanart.tv for high-quality artwork later
- OMDb has limited free tier but good movie data

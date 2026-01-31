# Movie Module

> Movie content management with metadata enrichment


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Developer Resources](#developer-resources)
- [Overview](#overview)
- [Architecture](#architecture)
- [Files](#files)
- [Entity: Movie](#entity-movie)
- [Service Configuration](#service-configuration)
- [Service Operations](#service-operations)
  - [Get Movie](#get-movie)
  - [List Movies](#list-movies)
  - [Create/Update/Delete](#createupdatedelete)
- [User Data](#user-data)
  - [Favorites](#favorites)
  - [Watch Progress](#watch-progress)
  - [Ratings](#ratings)
- [Metadata Flow](#metadata-flow)
- [Background Jobs](#background-jobs)
- [Database Schema](#database-schema)
- [Implementation Checklist](#implementation-checklist)
  - [Phase 1: Core Infrastructure](#phase-1-core-infrastructure)
  - [Phase 2: Database](#phase-2-database)
  - [Phase 3: Service Layer](#phase-3-service-layer)
  - [Phase 4: User Data](#phase-4-user-data)
  - [Phase 5: Background Jobs](#phase-5-background-jobs)
  - [Phase 6: API Integration](#phase-6-api-integration)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related](#related)

<!-- TOC-END -->

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Comprehensive spec with architecture, entities, operations |
| Sources | ✅ | TMDb, Radarr API docs linked |
| Instructions | ✅ | Implementation checklist added |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |
**Location**: `internal/content/movie/`

---

## Developer Resources

| Source | URL | Purpose |
|--------|-----|---------|
| TMDb API | [developers.themoviedb.org](https://developers.themoviedb.org/3) | Primary movie metadata |
| Radarr API | [radarr.video/docs/api](https://radarr.video/docs/api/) | Servarr integration (Radarr-first) |
| TMDb Design Doc | [integrations/metadata/video/TMDB.md](../../integrations/metadata/video/TMDB.md) | TMDb integration spec |
| Radarr Design Doc | [integrations/servarr/RADARR.md](../../integrations/servarr/RADARR.md) | Radarr integration spec |

---

## Overview

The Movie module provides complete movie library management:

- Entity definitions (Movie, Collection, Cast, Crew, etc.)
- Repository pattern with PostgreSQL implementation
- Service layer with otter caching
- Background jobs for metadata enrichment via River
- User data (ratings, watch history, favorites)

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                       API Layer                              │
│                    (ogen handlers)                           │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                    Movie Service                             │
│   - Local cache (otter)                                      │
│   - Business logic                                           │
│   - Resilience patterns                                      │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                    Repository Layer                          │
│   - PostgreSQL queries (sqlc)                                │
│   - User data (ratings, watch history)                       │
│   - Relations (cast, crew, genres)                           │
└─────────────────────────────────────────────────────────────┘
```

---

## Files

| File | Description |
|------|-------------|
| `entity.go` | Domain entities (Movie, Collection, Cast, etc.) |
| `repository.go` | Repository interface definition |
| `repository_pg.go` | PostgreSQL implementation |
| `repository_pg_user_data.go` | User ratings, favorites, watch history |
| `repository_pg_relations.go` | Cast, crew, genres, studios |
| `service.go` | Business logic with caching |
| `jobs.go` | River background jobs |
| `metadata_provider.go` | TMDb metadata interface |
| `module.go` | fx dependency injection |

---

## Entity: Movie

```go
type Movie struct {
    shared.ContentEntity

    // File info
    Container    string
    SizeBytes    int64
    RuntimeTicks int64

    // Metadata
    OriginalTitle string
    Tagline       string
    Overview      string
    ReleaseDate   *time.Time
    Year          int
    ContentRating string
    RatingLevel   int

    // Financials
    Budget  int64
    Revenue int64

    // Ratings
    CommunityRating float64
    VoteCount       int
    CriticRating    float64

    // Images
    PosterPath       string
    PosterBlurhash   string
    BackdropPath     string
    BackdropBlurhash string
    LogoPath         string

    // External IDs
    TmdbID int
    ImdbID string
    TvdbID int

    // Collection
    CollectionID    *uuid.UUID
    CollectionOrder int

    // Loaded on demand
    Collection *Collection
    Cast       []CastMember
    Crew       []CrewMember
    Directors  []CrewMember
    Writers    []CrewMember
    Studios    []Studio
    Genres     []Genre
    Images     []Image
    Videos     []Video
}
```

---

## Service Configuration

```go
type ServiceConfig struct {
    CacheMaxEntries int           // default: 10,000
    CacheTTL        time.Duration // default: 5 minutes
}
```

---

## Service Operations

### Get Movie

```go
func (s *Service) GetMovie(ctx context.Context, id uuid.UUID) (*Movie, error)
```

- Checks local otter cache first
- Falls back to PostgreSQL repository
- Caches result for subsequent requests

### List Movies

```go
func (s *Service) ListMovies(ctx context.Context, libraryID uuid.UUID, opts ListOptions) ([]*Movie, error)
```

### Create/Update/Delete

```go
func (s *Service) CreateMovie(ctx context.Context, movie *Movie) error
func (s *Service) UpdateMovie(ctx context.Context, movie *Movie) error
func (s *Service) DeleteMovie(ctx context.Context, id uuid.UUID) error
```

---

## User Data

### Favorites

```go
func (r *Repository) AddToFavorites(ctx context.Context, userID, movieID uuid.UUID) error
func (r *Repository) RemoveFromFavorites(ctx context.Context, userID, movieID uuid.UUID) error
func (r *Repository) IsFavorite(ctx context.Context, userID, movieID uuid.UUID) (bool, error)
```

### Watch Progress

```go
type WatchHistory struct {
    ID               uuid.UUID
    UserID           uuid.UUID
    ProfileID        *uuid.UUID
    MovieID          uuid.UUID
    PositionTicks    int64
    DurationTicks    int64
    PlayedPercentage float64
    Completed        bool
    CompletedAt      *time.Time
    DeviceName       string
    DeviceType       string
    ClientName       string
    PlayMethod       string
    StartedAt        time.Time
    LastUpdatedAt    time.Time
}
```

### Ratings

```go
type UserRating struct {
    UserID    uuid.UUID
    MovieID   uuid.UUID
    Rating    float64
    Review    string
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

---

## Metadata Flow

> See [00_SOURCE_OF_TRUTH.md](../../00_SOURCE_OF_TRUTH.md#metadata-priority-chain) for full priority chain

**Priority Order:**
1. **LOCAL CACHE** → otter cache, instant display
2. **ARR SERVICE** → Radarr (cached TMDb metadata)
3. **EXTERNAL** → Direct TMDb API (if Radarr unavailable)
4. **ENRICHMENT** → Background jobs for additional data

**Primary Metadata Source:** TMDb
**Arr Integration:** Radarr

```go
type MetadataProvider interface {
    SearchMovies(ctx context.Context, query string, year int) ([]MovieSearchResult, error)
    GetMovieMetadata(ctx context.Context, id int) (*MovieMetadata, error)
    MatchMovie(ctx context.Context, title string, year int, imdbID string) (*MovieMetadata, error)
}
```

---

## Background Jobs

Metadata enrichment runs via River:

- `MovieMetadataRefreshJob` - Refresh metadata from TMDb
- `MovieImageDownloadJob` - Download and cache images
- `MovieCollectionSyncJob` - Sync collection membership

---

## Database Schema

```sql
-- movie schema
CREATE SCHEMA IF NOT EXISTS movie;

CREATE TABLE movie.movies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    library_id UUID NOT NULL REFERENCES libraries(id),
    path TEXT NOT NULL,
    title TEXT NOT NULL,
    sort_title TEXT,
    original_title TEXT,
    tagline TEXT,
    overview TEXT,
    release_date DATE,
    year SMALLINT,
    runtime_ticks BIGINT,
    content_rating TEXT,
    rating_level SMALLINT,
    community_rating NUMERIC(3,1),
    vote_count INTEGER,
    poster_path TEXT,
    backdrop_path TEXT,
    tmdb_id INTEGER,
    imdb_id TEXT,
    collection_id UUID REFERENCES movie.collections(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## Implementation Checklist

### Phase 1: Core Infrastructure
- [ ] Create `internal/content/movie/` package structure
- [ ] Define `entity.go` with Movie, Collection, Cast, Crew structs
- [ ] Create `repository.go` interface definition
- [ ] Implement `repository_pg.go` with sqlc queries
- [ ] Add fx module wiring in `module.go`

### Phase 2: Database
- [ ] Create migration `000XXX_create_movie_schema.up.sql`
- [ ] Create `movie.movies` table with all columns
- [ ] Create `movie.collections` table
- [ ] Create `movie.cast` and `movie.crew` junction tables
- [ ] Create `movie.genres` and `movie.studios` tables
- [ ] Add indexes (library_id, tmdb_id, imdb_id, title search)
- [ ] Write sqlc queries in `queries/movie/`

### Phase 3: Service Layer
- [ ] Implement `service.go` with otter caching
- [ ] Add GetMovie, ListMovies, CreateMovie, UpdateMovie, DeleteMovie
- [ ] Implement cache invalidation on mutations
- [ ] Add resilience patterns (circuit breaker, retries)

### Phase 4: User Data
- [ ] Implement `repository_pg_user_data.go`
- [ ] Add favorites (add, remove, list, check)
- [ ] Add watch history tracking
- [ ] Add user ratings and reviews
- [ ] Implement watch progress persistence

### Phase 5: Background Jobs
- [ ] Create River job definitions in `jobs.go`
- [ ] Implement `MovieMetadataRefreshJob`
- [ ] Implement `MovieImageDownloadJob`
- [ ] Implement `MovieCollectionSyncJob`
- [ ] Add job scheduling and retry logic

### Phase 6: API Integration
- [ ] Define OpenAPI endpoints for movies
- [ ] Generate ogen handlers
- [ ] Wire handlers to service layer
- [ ] Add authentication/authorization checks

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
| [PostgreSQL Arrays](https://www.postgresql.org/docs/current/arrays.html) | [Local](../../../sources/database/postgresql-arrays.md) |
| [PostgreSQL JSON Functions](https://www.postgresql.org/docs/current/functions-json.html) | [Local](../../../sources/database/postgresql-json.md) |
| [Radarr API Docs](https://radarr.video/docs/api/) | [Local](../../../sources/apis/radarr-docs.md) |
| [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) | [Local](../../../sources/tooling/river.md) |
| [Uber fx](https://pkg.go.dev/go.uber.org/fx) | [Local](../../../sources/tooling/fx.md) |
| [go-blurhash](https://pkg.go.dev/github.com/bbrks/go-blurhash) | [Local](../../../sources/media/go-blurhash.md) |
| [ogen OpenAPI Generator](https://pkg.go.dev/github.com/ogen-go/ogen) | [Local](../../../sources/tooling/ogen.md) |
| [pgx PostgreSQL Driver](https://pkg.go.dev/github.com/jackc/pgx/v5) | [Local](../../../sources/database/pgx.md) |
| [sqlc](https://docs.sqlc.dev/en/stable/) | [Local](../../../sources/database/sqlc.md) |
| [sqlc Configuration](https://docs.sqlc.dev/en/stable/reference/config.html) | [Local](../../../sources/database/sqlc-config.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Video](INDEX.md)

### In This Section

- [TV Show Module](TVSHOW_MODULE.md)

### Related Topics

- [Revenge - Architecture v2](../../architecture/01_ARCHITECTURE.md) _Architecture_
- [Revenge - Design Principles](../../architecture/02_DESIGN_PRINCIPLES.md) _Architecture_
- [Revenge - Metadata System](../../architecture/03_METADATA_SYSTEM.md) _Architecture_
- [Revenge - Player Architecture](../../architecture/04_PLAYER_ARCHITECTURE.md) _Architecture_
- [Plugin Architecture Decision](../../architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) _Architecture_

### Indexes

- [Design Index](../../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

## Related

- [TV Show Module](TVSHOW_MODULE.md) - TV series management
- [Library Service](../../services/LIBRARY.md) - Library management
- [Metadata Service](../../services/METADATA.md) - TMDb/Radarr providers
- [Integrations: Radarr](../../integrations/servarr/RADARR.md) - Radarr integration

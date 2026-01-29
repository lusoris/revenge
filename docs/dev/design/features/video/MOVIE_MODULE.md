# Movie Module

> Movie content management with metadata enrichment

**Location**: `internal/content/movie/`

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

## Metadata Provider

Primary: **Radarr** (Servarr-first principle)
Fallback: **TMDb** (via background jobs)

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

## Related

- [TV Show Module](TVSHOW_MODULE.md) - TV series management
- [Library Service](../../services/LIBRARY.md) - Library management
- [Metadata Service](../../services/METADATA.md) - TMDb/Radarr providers
- [Integrations: Radarr](../../integrations/servarr/RADARR.md) - Radarr integration

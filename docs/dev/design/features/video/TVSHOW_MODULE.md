# TV Show Module

> TV series, seasons, and episodes management

**Location**: `internal/content/tvshow/`

---

## Overview

The TV Show module provides complete television library management:

- Hierarchical structure (Series → Seasons → Episodes)
- Entity definitions with full metadata support
- Repository pattern with PostgreSQL implementation
- Service layer with otter caching
- Background jobs for metadata enrichment
- User data (ratings, watch history, progress tracking)

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                       API Layer                              │
│                    (ogen handlers)                           │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                   TVShow Service                             │
│   - Local cache (otter)                                      │
│   - Business logic                                           │
│   - Series/Season/Episode operations                         │
└─────────────────────────┬───────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────┐
│                    Repository Layer                          │
│   - PostgreSQL queries (sqlc)                                │
│   - Hierarchical queries                                     │
│   - User data, relations                                     │
└─────────────────────────────────────────────────────────────┘
```

---

## Files

| File | Description |
|------|-------------|
| `entity.go` | Domain entities (Series, Season, Episode, etc.) |
| `repository.go` | Repository interface definition |
| `repository_pg.go` | PostgreSQL implementation |
| `repository_pg_user_data.go` | User ratings, watch history |
| `repository_pg_relations.go` | Cast, crew, genres |
| `service.go` | Business logic with caching |
| `jobs.go` | River background jobs |
| `metadata_provider.go` | TMDb metadata interface |
| `module.go` | fx dependency injection |

---

## Entity Hierarchy

### Series

```go
type Series struct {
    shared.ContentEntity

    // Metadata
    OriginalTitle  string
    Tagline        string
    Overview       string
    FirstAirDate   *time.Time
    LastAirDate    *time.Time
    Status         string  // Continuing, Ended, Canceled
    ContentRating  string
    RatingLevel    int

    // Counts
    SeasonCount   int
    EpisodeCount  int

    // Ratings
    CommunityRating float64
    VoteCount       int

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

    // Loaded on demand
    Seasons  []Season
    Cast     []CastMember
    Crew     []CrewMember
    Genres   []Genre
    Studios  []Studio
}
```

### Season

```go
type Season struct {
    ID            uuid.UUID
    SeriesID      uuid.UUID
    SeasonNumber  int
    Name          string
    Overview      string
    AirDate       *time.Time
    EpisodeCount  int
    PosterPath    string
    PosterBlurhash string

    // External IDs
    TmdbID int
    TvdbID int

    // Loaded on demand
    Episodes []Episode
}
```

### Episode

```go
type Episode struct {
    shared.ContentEntity

    // Hierarchy
    SeriesID      uuid.UUID
    SeasonID      uuid.UUID
    SeasonNumber  int
    EpisodeNumber int

    // Metadata
    Overview      string
    AirDate       *time.Time
    RuntimeTicks  int64

    // Ratings
    CommunityRating float64
    VoteCount       int

    // Images
    StillPath      string
    StillBlurhash  string

    // External IDs
    TmdbID int
    ImdbID string
    TvdbID int

    // File info
    Path      string
    Container string
    SizeBytes int64
}
```

---

## Service Operations

### Series Operations

```go
func (s *Service) GetSeries(ctx context.Context, id uuid.UUID) (*Series, error)
func (s *Service) GetSeriesWithSeasons(ctx context.Context, id uuid.UUID) (*Series, error)
func (s *Service) ListSeries(ctx context.Context, libraryID uuid.UUID, opts ListOptions) ([]*Series, error)
func (s *Service) CreateSeries(ctx context.Context, series *Series) error
func (s *Service) UpdateSeries(ctx context.Context, series *Series) error
func (s *Service) DeleteSeries(ctx context.Context, id uuid.UUID) error
```

### Season Operations

```go
func (s *Service) GetSeason(ctx context.Context, id uuid.UUID) (*Season, error)
func (s *Service) GetSeasonWithEpisodes(ctx context.Context, id uuid.UUID) (*Season, error)
func (s *Service) ListSeasons(ctx context.Context, seriesID uuid.UUID) ([]*Season, error)
```

### Episode Operations

```go
func (s *Service) GetEpisode(ctx context.Context, id uuid.UUID) (*Episode, error)
func (s *Service) ListEpisodes(ctx context.Context, seasonID uuid.UUID) ([]*Episode, error)
func (s *Service) GetNextEpisode(ctx context.Context, userID, seriesID uuid.UUID) (*Episode, error)
```

---

## User Data

### Watch Progress

Tracks progress at episode level with series-wide aggregation:

```go
type EpisodeWatchHistory struct {
    ID               uuid.UUID
    UserID           uuid.UUID
    EpisodeID        uuid.UUID
    PositionTicks    int64
    DurationTicks    int64
    PlayedPercentage float64
    Completed        bool
    CompletedAt      *time.Time
    StartedAt        time.Time
    LastUpdatedAt    time.Time
}
```

### Series Progress

Calculated from episode watch history:

- Episodes watched
- Episodes remaining
- Next episode to watch
- Continue watching position

---

## Metadata Provider

Primary: **Sonarr** (Servarr-first principle)
Fallback: **TMDb** / **TheTVDB** (via background jobs)

```go
type MetadataProvider interface {
    SearchSeries(ctx context.Context, query string) ([]SeriesSearchResult, error)
    GetSeriesMetadata(ctx context.Context, id int) (*SeriesMetadata, error)
    GetSeasonMetadata(ctx context.Context, seriesID, seasonNumber int) (*SeasonMetadata, error)
    GetEpisodeMetadata(ctx context.Context, seriesID, seasonNumber, episodeNumber int) (*EpisodeMetadata, error)
}
```

---

## Background Jobs

Metadata enrichment runs via River:

- `SeriesMetadataRefreshJob` - Refresh series metadata
- `SeasonMetadataRefreshJob` - Refresh season metadata
- `EpisodeMetadataRefreshJob` - Refresh episode metadata
- `SeriesImageDownloadJob` - Download series images
- `NewEpisodeCheckJob` - Check for new episodes

---

## Database Schema

```sql
-- tvshow schema
CREATE SCHEMA IF NOT EXISTS tvshow;

CREATE TABLE tvshow.series (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    library_id UUID NOT NULL REFERENCES libraries(id),
    path TEXT NOT NULL,
    title TEXT NOT NULL,
    sort_title TEXT,
    original_title TEXT,
    overview TEXT,
    first_air_date DATE,
    last_air_date DATE,
    status TEXT,
    season_count SMALLINT,
    episode_count SMALLINT,
    tmdb_id INTEGER,
    tvdb_id INTEGER,
    imdb_id TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE tvshow.seasons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    series_id UUID NOT NULL REFERENCES tvshow.series(id) ON DELETE CASCADE,
    season_number SMALLINT NOT NULL,
    name TEXT,
    overview TEXT,
    air_date DATE,
    episode_count SMALLINT,
    tmdb_id INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (series_id, season_number)
);

CREATE TABLE tvshow.episodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    series_id UUID NOT NULL REFERENCES tvshow.series(id) ON DELETE CASCADE,
    season_id UUID NOT NULL REFERENCES tvshow.seasons(id) ON DELETE CASCADE,
    season_number SMALLINT NOT NULL,
    episode_number SMALLINT NOT NULL,
    path TEXT,
    title TEXT NOT NULL,
    overview TEXT,
    air_date DATE,
    runtime_ticks BIGINT,
    tmdb_id INTEGER,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (series_id, season_number, episode_number)
);
```

---

## Related

- [Movie Module](MOVIE_MODULE.md) - Movie management
- [Watch Next](../playback/WATCH_NEXT_CONTINUE_WATCHING.md) - Continue watching
- [Library Service](../../services/LIBRARY.md) - Library management
- [Integrations: Sonarr](../../integrations/servarr/SONARR.md) - Sonarr integration

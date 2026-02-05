# TODO A11: TV Shows Module Implementation

**Phase**: A11
**Priority**: P1 (High - First expansion module)
**Effort**: 32-48 hours (with shared abstractions: <1 week!)
**Status**: Pending
**Dependencies**: A9 (Multi-Language), A10 (Shared Abstractions)
**Created**: 2026-02-05

---

## Overview

**Goal**: Implement TV shows module using shared abstractions from A10.

**Benefits of doing it now**:
- Validates shared abstractions work correctly
- Multi-language support from day 1 (A9)
- Cluster-ready from day 1 (A8)
- <1 week implementation (vs 3-4 weeks without A10)

**Approach**:
- Reuse 60-70% of code via shared abstractions
- TV-specific adapters for parsing, matching, metadata
- Separate database tables (tvshow schema)
- TMDb as primary metadata source (not TheTVDB initially)
- Sonarr integration similar to Radarr

---

## Decision Log

| Topic | Decision | Reason |
|-------|----------|--------|
| Metadata Provider | TMDb primary | Already have client, multi-language support |
| Database Schema | Separate tvshow schema | Clear isolation, independent migrations |
| Sonarr Integration | Parallel to TV implementation | Similar to Radarr pattern |
| Filename Parsing | S##E## + 1x05 + 105 formats | Cover common naming conventions |

---

## Tasks

### A11.1: Database Schema 🔴 CRITICAL

**Priority**: P0
**Effort**: 8-12h

#### A11.1.1: Create TV Schema

**Location**: `migrations/000031_create_tvshow_schema.up.sql`

```sql
-- Create tvshow schema
CREATE SCHEMA IF NOT EXISTS tvshow;

-- Series table
CREATE TABLE tvshow.series (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- External IDs
    tmdb_id INTEGER UNIQUE,
    tvdb_id INTEGER UNIQUE,
    imdb_id TEXT,
    sonarr_id INTEGER,

    -- Default language fields
    title TEXT NOT NULL,
    tagline TEXT,
    overview TEXT,

    -- Multi-language (from A9)
    titles_i18n JSONB DEFAULT '{}',
    taglines_i18n JSONB DEFAULT '{}',
    overviews_i18n JSONB DEFAULT '{}',
    age_ratings JSONB DEFAULT '{}',

    -- Original
    original_language TEXT NOT NULL DEFAULT 'en',
    original_title TEXT,

    -- Series metadata
    status TEXT, -- 'Returning Series', 'Ended', 'Canceled'
    type TEXT,   -- 'Scripted', 'Documentary', 'Reality'
    first_air_date DATE,
    last_air_date DATE,

    -- Media info
    poster_path TEXT,
    backdrop_path TEXT,

    -- Stats
    total_seasons INTEGER DEFAULT 0,
    total_episodes INTEGER DEFAULT 0,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Seasons table
CREATE TABLE tvshow.seasons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    series_id UUID NOT NULL REFERENCES tvshow.series(id) ON DELETE CASCADE,

    tmdb_id INTEGER,
    season_number INTEGER NOT NULL,

    -- Multi-language
    name TEXT NOT NULL,
    overview TEXT,
    names_i18n JSONB DEFAULT '{}',
    overviews_i18n JSONB DEFAULT '{}',

    -- Media
    poster_path TEXT,

    -- Stats
    episode_count INTEGER DEFAULT 0,
    air_date DATE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(series_id, season_number)
);

-- Episodes table
CREATE TABLE tvshow.episodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    series_id UUID NOT NULL REFERENCES tvshow.series(id) ON DELETE CASCADE,
    season_id UUID NOT NULL REFERENCES tvshow.seasons(id) ON DELETE CASCADE,

    tmdb_id INTEGER,
    tvdb_id INTEGER,
    imdb_id TEXT,

    season_number INTEGER NOT NULL,
    episode_number INTEGER NOT NULL,

    -- Multi-language
    title TEXT NOT NULL,
    overview TEXT,
    titles_i18n JSONB DEFAULT '{}',
    overviews_i18n JSONB DEFAULT '{}',

    -- Episode metadata
    air_date DATE,
    runtime INTEGER, -- minutes

    -- Media
    still_path TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(series_id, season_number, episode_number)
);

-- Episode files table
CREATE TABLE tvshow.episode_files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    episode_id UUID NOT NULL REFERENCES tvshow.episodes(id) ON DELETE CASCADE,
    library_id UUID NOT NULL,

    path TEXT NOT NULL UNIQUE,
    size BIGINT NOT NULL,

    -- Media info
    duration INTEGER,
    video_codec TEXT,
    audio_codec TEXT,
    resolution TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Watch progress
CREATE TABLE tvshow.episode_watched (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    episode_id UUID NOT NULL REFERENCES tvshow.episodes(id) ON DELETE CASCADE,

    watched BOOLEAN NOT NULL DEFAULT FALSE,
    progress_seconds INTEGER DEFAULT 0,
    watched_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(user_id, episode_id)
);

-- Indexes
CREATE INDEX idx_series_tmdb_id ON tvshow.series(tmdb_id);
CREATE INDEX idx_series_titles_i18n ON tvshow.series USING GIN (titles_i18n);
CREATE INDEX idx_seasons_series ON tvshow.seasons(series_id);
CREATE INDEX idx_episodes_series ON tvshow.episodes(series_id);
CREATE INDEX idx_episodes_season ON tvshow.episodes(season_id);
CREATE INDEX idx_episode_files_episode ON tvshow.episode_files(episode_id);
CREATE INDEX idx_episode_watched_user_episode ON tvshow.episode_watched(user_id, episode_id);
```

**Subtasks**:
- [ ] Create migration (up/down)
- [ ] Test on dev database
- [ ] Verify indexes created
- [ ] Add to sqlc config

---

### A11.2: TV Domain Models 🔴 CRITICAL

**Priority**: P0
**Effort**: 4-6h
**Location**: `internal/content/tvshow/types.go`

```go
package tvshow

type Series struct {
    ID uuid.UUID `json:"id"`

    TMDbID   *int    `json:"tmdb_id,omitempty"`
    TVDbID   *int    `json:"tvdb_id,omitempty"`
    IMDbID   *string `json:"imdb_id,omitempty"`
    SonarrID *int    `json:"sonarr_id,omitempty"`

    // Default language
    Title    string `json:"title"`
    Tagline  string `json:"tagline,omitempty"`
    Overview string `json:"overview,omitempty"`

    // Multi-language (A9)
    TitlesI18n    map[string]string `json:"titles_i18n,omitempty"`
    TaglinesI18n  map[string]string `json:"taglines_i18n,omitempty"`
    OverviewsI18n map[string]string `json:"overviews_i18n,omitempty"`
    AgeRatings    map[string]map[string]string `json:"age_ratings,omitempty"`

    OriginalLanguage string `json:"original_language"`
    OriginalTitle    string `json:"original_title"`

    Status         string    `json:"status"`
    Type           string    `json:"type"`
    FirstAirDate   time.Time `json:"first_air_date"`
    LastAirDate    time.Time `json:"last_air_date,omitempty"`

    PosterPath   string `json:"poster_path,omitempty"`
    BackdropPath string `json:"backdrop_path,omitempty"`

    TotalSeasons  int `json:"total_seasons"`
    TotalEpisodes int `json:"total_episodes"`

    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type Season struct {
    ID       uuid.UUID `json:"id"`
    SeriesID uuid.UUID `json:"series_id"`

    TMDbID       *int `json:"tmdb_id,omitempty"`
    SeasonNumber int  `json:"season_number"`

    Name     string `json:"name"`
    Overview string `json:"overview,omitempty"`

    NamesI18n     map[string]string `json:"names_i18n,omitempty"`
    OverviewsI18n map[string]string `json:"overviews_i18n,omitempty"`

    PosterPath   string    `json:"poster_path,omitempty"`
    EpisodeCount int       `json:"episode_count"`
    AirDate      time.Time `json:"air_date,omitempty"`
}

type Episode struct {
    ID       uuid.UUID `json:"id"`
    SeriesID uuid.UUID `json:"series_id"`
    SeasonID uuid.UUID `json:"season_id"`

    TMDbID  *int    `json:"tmdb_id,omitempty"`
    TVDbID  *int    `json:"tvdb_id,omitempty"`
    IMDbID  *string `json:"imdb_id,omitempty"`

    SeasonNumber  int `json:"season_number"`
    EpisodeNumber int `json:"episode_number"`

    Title    string `json:"title"`
    Overview string `json:"overview,omitempty"`

    TitlesI18n    map[string]string `json:"titles_i18n,omitempty"`
    OverviewsI18n map[string]string `json:"overviews_i18n,omitempty"`

    AirDate   time.Time `json:"air_date,omitempty"`
    Runtime   int       `json:"runtime"` // minutes
    StillPath string    `json:"still_path,omitempty"`
}

type EpisodeFile struct {
    ID        uuid.UUID `json:"id"`
    EpisodeID uuid.UUID `json:"episode_id"`
    LibraryID uuid.UUID `json:"library_id"`

    Path string `json:"path"`
    Size int64  `json:"size"`

    Duration    int    `json:"duration,omitempty"`
    VideoCodec  string `json:"video_codec,omitempty"`
    AudioCodec  string `json:"audio_codec,omitempty"`
    Resolution  string `json:"resolution,omitempty"`
}
```

**Subtasks**:
- [ ] Define all TV types
- [ ] Add helper methods (GetTitle, GetOverview)
- [ ] Implement ContentItem interface
- [ ] Write unit tests

---

### A11.3: TV Adapters (Using A10) 🟠 HIGH

**Priority**: P1
**Effort**: 8-12h

#### A11.3.1: Scanner Adapter

**Location**: `internal/content/tvshow/adapters/scanner_adapter.go`

```go
type TVFileParser struct{}

func (p *TVFileParser) Parse(filename string) (*scanner.ParseResult, error) {
    // Parse patterns:
    // - Series.Name.S01E05.Episode.Title.1080p.mkv
    // - Series Name - 1x05 - Episode Title.mkv
    // - Series.Name.105.Episode.Title.mkv

    series, season, episode := extractTVInfo(filename)

    return &scanner.ParseResult{
        Title: series,
        Metadata: map[string]any{
            "season":  season,
            "episode": episode,
        },
    }, nil
}
```

**Subtasks**:
- [ ] Implement S##E## parsing
- [ ] Implement 1x05 parsing
- [ ] Implement 105 parsing
- [ ] Handle multi-episode files (S01E01-E02)
- [ ] Write comprehensive tests

---

#### A11.3.2: Matcher Adapter

**Location**: `internal/content/tvshow/adapters/matcher_adapter.go`

```go
type TVMatchStrategy struct {
    repo     *tvshow.Repository
    provider metadata.Provider[*tvshow.Episode]
}

func (s *TVMatchStrategy) CalculateConfidence(parse scanner.ParseResult, candidate *tvshow.Episode) float64 {
    // Must match season AND episode
    season := parse.Metadata["season"].(int)
    episode := parse.Metadata["episode"].(int)

    if candidate.SeasonNumber != season || candidate.EpisodeNumber != episode {
        return 0.0
    }

    // Series name fuzzy match
    seriesSimilarity := matcher.LevenshteinNormalized(parse.Title, candidate.Series.Title)
    return seriesSimilarity
}
```

**Subtasks**:
- [ ] Implement TV matching strategy
- [ ] Handle multi-episode files
- [ ] Write tests

---

#### A11.3.3: TMDb TV Adapter

**Location**: `internal/content/tvshow/adapters/tmdb_tv_adapter.go`

```go
type TMDbTVProvider struct {
    metadata.HTTPProvider[*tvshow.Series]
}

func NewTMDbTVProvider(apiKey string, cache *cache.Cache) *TMDbTVProvider {
    return &TMDbTVProvider{
        HTTPProvider: metadata.HTTPProvider[*tvshow.Series]{
            BaseURL: "https://api.themoviedb.org/3",
            APIKey:  apiKey,
            Limiter: rate.NewLimiter(rate.Every(250*time.Millisecond), 40),
            Cache:   cache,
            Mapper:  &TMDbTVMapper{},
        },
    }
}

// TMDb TV endpoints:
// GET /tv/{id}
// GET /tv/{id}/season/{season_number}
// GET /tv/{id}/season/{season_number}/episode/{episode_number}
```

**Subtasks**:
- [ ] Implement TMDb TV client
- [ ] Fetch series, seasons, episodes
- [ ] Multi-language support (A9)
- [ ] Write integration tests

---

### A11.4: TV Repository 🟠 HIGH

**Priority**: P1
**Effort**: 6-8h
**Location**: `internal/content/tvshow/repository_postgres.go`

**Implement**: Full CRUD for Series, Seasons, Episodes, Files

**Subtasks**:
- [ ] Write sqlc queries
- [ ] Generate Go code
- [ ] Implement repository interface
- [ ] Write repository tests
- [ ] Add transactions support

---

### A11.5: TV Service 🟠 HIGH

**Priority**: P1
**Effort**: 6-8h
**Location**: `internal/content/tvshow/service.go`

**Extend**: BaseLibraryService[*Series] from A10

**Subtasks**:
- [ ] Implement TV-specific business logic
- [ ] Add season/episode management
- [ ] Integrate with shared library service
- [ ] Write service tests

---

### A11.6: TV Jobs 🟡 MEDIUM

**Priority**: P2
**Effort**: 4-6h
**Location**: `internal/content/tvshow/tvjobs/`

**Jobs**:
- Library scan
- File match
- Metadata refresh
- Search indexing

**Subtasks**:
- [ ] Extend base job workers from A10
- [ ] Implement TV-specific job logic
- [ ] Register with River
- [ ] Write job tests

---

### A11.7: Sonarr Integration 🟡 MEDIUM

**Priority**: P2
**Effort**: 8-12h
**Location**: `internal/integration/sonarr/`

**Clone from Radarr**:
- Client
- Sync service
- Webhooks
- Jobs
- Mapper

**Subtasks**:
- [ ] Implement Sonarr client
- [ ] Sync service
- [ ] Webhook handler
- [ ] Write integration tests

---

### A11.8: TV API Handlers 🟡 MEDIUM

**Priority**: P2
**Effort**: 4-6h
**Location**: `internal/api/tv_handlers.go`

**Endpoints**:
- GET /api/v1/tv/series
- GET /api/v1/tv/series/{id}
- GET /api/v1/tv/series/{id}/seasons
- GET /api/v1/tv/series/{id}/seasons/{season}/episodes
- POST /api/v1/tv/watch-progress

**Subtasks**:
- [ ] Implement handlers
- [ ] Add localization support (A9)
- [ ] Write handler tests
- [ ] Update OpenAPI spec

---

## Testing

- [ ] Unit tests: 80%+ coverage
- [ ] Integration tests with TMDb
- [ ] Repository tests with testcontainers
- [ ] Scanner adapter tests (all patterns)
- [ ] End-to-end library scan test

---

## Documentation

- [ ] TV module architecture
- [ ] TMDb TV API usage
- [ ] Sonarr integration guide
- [ ] API documentation
- [ ] User guide for TV libraries

---

## Verification Checklist

- [ ] Database migrations successful
- [ ] TV types and models complete
- [ ] All adapters implemented
- [ ] Repository fully functional
- [ ] Service layer working
- [ ] Jobs registered and working
- [ ] API endpoints functional
- [ ] Sonarr integration working
- [ ] Tests passing (80%+)
- [ ] Can scan TV library end-to-end
- [ ] Multi-language support working
- [ ] Cluster-ready (A8)

---

**Completion Criteria**:
✅ TV module fully functional
✅ Can scan TV library
✅ TMDb metadata fetching works
✅ Sonarr integration works
✅ Multi-language from day 1
✅ Tests passing
✅ Implemented in <1 week (thanks to A10!)

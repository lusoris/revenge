# Module Implementation TODO

> Concrete coding tasks for implementing Revenge modules

**Last Updated**: 2026-01-29
**Current Phase**: Phase 2 - Movie Module Implementation

---

## Phase 1: Core Infrastructure Completion

### 1.1 Migrate Cache to rueidis ✅

**Files modified:**
- [x] `internal/infra/cache/cache.go` - Using rueidis client
- [x] `go.mod` - Has `github.com/redis/rueidis`

### 1.2 Register Missing Modules in main.go ✅

**File:** `cmd/revenge/main.go` - cache.Module, search.Module, jobs.Module registered

### 1.3 Add Local Cache (otter) ✅

**File:** `internal/infra/cache/local.go`
- [x] LocalCache struct with otter
- [x] Config integration (LocalCapacity, LocalTTL)
- [x] GetOrSet, GetJSON, SetJSON methods
- [x] Stats and HitRate tracking

### 1.4 Add API Cache (sturdyc) ✅

**File:** `internal/infra/cache/api.go`
- [x] APICache with request coalescing
- [x] MetadataCache with provider-specific configs
- [x] Config integration (APICapacity, APINumShards, APITTL)

### 1.5 Implement Radarr Provider ✅

**Files:**
- [x] `internal/service/metadata/radarr/types.go` - Full API types
- [x] `internal/service/metadata/radarr/client.go` - HTTP client with circuit breaker
- [x] `internal/service/metadata/radarr/provider.go` - Metadata provider
- [x] `internal/service/metadata/radarr/module.go` - fx module

### 1.6 Implement Central MetadataService ✅

**Files:**
- [x] `internal/service/metadata/service.go` - Orchestration with fallback
- [x] `internal/service/metadata/module.go` - fx module with all providers

---

## Phase 2: Movie Module Implementation

### 2.1 Database Migration

**New file:** `internal/infra/database/migrations/movie/000001_movies.up.sql`

```sql
-- Movies table
CREATE TABLE movies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    original_title TEXT,
    sort_title TEXT NOT NULL,
    overview TEXT,
    tagline TEXT,
    release_date DATE,
    runtime_minutes INT,
    status TEXT DEFAULT 'released',
    -- External IDs
    tmdb_id INT UNIQUE,
    imdb_id TEXT UNIQUE,
    radarr_id INT,
    -- Metadata
    certification TEXT,
    original_language TEXT,
    budget BIGINT,
    revenue BIGINT,
    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_movies_tmdb_id ON movies(tmdb_id);
CREATE INDEX idx_movies_release_date ON movies(release_date);
CREATE INDEX idx_movies_sort_title ON movies(sort_title);
```

### 2.2 sqlc Queries

**New file:** `internal/infra/database/queries/movie/movies.sql`

```sql
-- name: GetMovieByID :one
SELECT * FROM movies WHERE id = $1;

-- name: GetMovieByTmdbID :one
SELECT * FROM movies WHERE tmdb_id = $1;

-- name: ListMovies :many
SELECT * FROM movies
ORDER BY sort_title
LIMIT $1 OFFSET $2;

-- name: CreateMovie :one
INSERT INTO movies (title, original_title, sort_title, overview, tmdb_id, imdb_id, release_date, runtime_minutes)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateMovie :one
UPDATE movies SET
    title = COALESCE($2, title),
    overview = COALESCE($3, overview),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteMovie :exec
DELETE FROM movies WHERE id = $1;
```

### 2.3 Domain Entity

**New file:** `internal/content/movie/entity.go`

```go
package movie

import (
    "time"
    "github.com/google/uuid"
)

type Movie struct {
    ID              uuid.UUID
    Title           string
    OriginalTitle   string
    SortTitle       string
    Overview        string
    Tagline         string
    ReleaseDate     *time.Time
    RuntimeMinutes  int
    Status          string
    TmdbID          *int
    ImdbID          *string
    RadarrID        *int
    Certification   string
    OriginalLanguage string
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

type CreateMovieParams struct {
    Title          string
    OriginalTitle  string
    Overview       string
    TmdbID         *int
    ImdbID         *string
    ReleaseDate    *time.Time
    RuntimeMinutes int
}
```

### 2.4 Repository

**New file:** `internal/content/movie/repository.go`

```go
package movie

import (
    "context"
    "github.com/google/uuid"
)

type Repository interface {
    GetByID(ctx context.Context, id uuid.UUID) (*Movie, error)
    GetByTmdbID(ctx context.Context, tmdbID int) (*Movie, error)
    List(ctx context.Context, limit, offset int) ([]*Movie, error)
    Create(ctx context.Context, params CreateMovieParams) (*Movie, error)
    Update(ctx context.Context, id uuid.UUID, params UpdateMovieParams) (*Movie, error)
    Delete(ctx context.Context, id uuid.UUID) error
}

type repositoryImpl struct {
    queries *db.Queries
}

func NewRepository(queries *db.Queries) Repository {
    return &repositoryImpl{queries: queries}
}
```

### 2.5 Service

**New file:** `internal/content/movie/service.go`

```go
package movie

import (
    "context"
    "log/slog"
    "github.com/google/uuid"
)

type Service struct {
    repo   Repository
    logger *slog.Logger
}

func NewService(repo Repository, logger *slog.Logger) *Service {
    return &Service{
        repo:   repo,
        logger: logger.With(slog.String("service", "movie")),
    }
}

func (s *Service) GetMovie(ctx context.Context, id uuid.UUID) (*Movie, error) {
    return s.repo.GetByID(ctx, id)
}

func (s *Service) CreateMovie(ctx context.Context, params CreateMovieParams) (*Movie, error) {
    // Generate sort title
    params.SortTitle = generateSortTitle(params.Title)
    return s.repo.Create(ctx, params)
}
```

### 2.6 HTTP Handler (ogen)

**New file:** `internal/content/movie/handler.go`

```go
package movie

import (
    "context"
    "github.com/google/uuid"
    api "revenge/api/generated"
)

type Handler struct {
    svc *Service
}

func NewHandler(svc *Service) *Handler {
    return &Handler{svc: svc}
}

// Implements ogen interface
func (h *Handler) GetMovie(ctx context.Context, params api.GetMovieParams) (api.GetMovieRes, error) {
    id, err := uuid.Parse(params.ID)
    if err != nil {
        return &api.GetMovieBadRequest{}, nil
    }

    movie, err := h.svc.GetMovie(ctx, id)
    if err != nil {
        return &api.GetMovieNotFound{}, nil
    }

    return &api.Movie{
        ID:    movie.ID.String(),
        Title: movie.Title,
        // ... map fields
    }, nil
}
```

### 2.7 fx Module

**New file:** `internal/content/movie/module.go`

```go
package movie

import "go.uber.org/fx"

var Module = fx.Module("movie",
    fx.Provide(
        NewRepository,
        NewService,
        NewHandler,
    ),
)
```

### 2.8 River Jobs

**New file:** `internal/content/movie/jobs.go`

```go
package movie

import (
    "context"
    "github.com/google/uuid"
    "github.com/riverqueue/river"
)

// ScanLibraryArgs for library scanning job
type ScanLibraryArgs struct {
    LibraryID uuid.UUID `json:"library_id"`
    FullScan  bool      `json:"full_scan"`
}

func (ScanLibraryArgs) Kind() string { return "movie.scan_library" }

type ScanLibraryWorker struct {
    river.WorkerDefaults[ScanLibraryArgs]
    svc *Service
}

func (w *ScanLibraryWorker) Work(ctx context.Context, job *river.Job[ScanLibraryArgs]) error {
    // Scan library for new movies
    return nil
}

// FetchMetadataArgs for metadata fetching job
type FetchMetadataArgs struct {
    MovieID uuid.UUID `json:"movie_id"`
}

func (FetchMetadataArgs) Kind() string { return "movie.fetch_metadata" }
```

---

## Phase 3: TV Show Module

Same pattern as Movie module with:
- `series`, `seasons`, `episodes` tables
- Sonarr integration
- TheTVDB/TMDb metadata providers

---

## Phase 4: Music Module

Same pattern with:
- `artists`, `albums`, `tracks`, `music_videos` tables
- Lidarr integration
- MusicBrainz/Last.fm metadata

---

## Phase 5: QAR Module (Adult Content)

> Queen Anne's Revenge - isolated adult content in `qar` schema
> See [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#qar-obfuscation-terminology) for terminology mapping

### 5.1 Schema Already Complete ✅

The QAR module uses isolated `qar` PostgreSQL schema:
- `qar.expeditions` - Movies
- `qar.voyages` - Scenes
- `qar.crew` - Performers
- `qar.ports` - Studios
- `qar.flags` - Tags
- `qar.fleets` - Libraries

See [M5 in TODO.md](../../../../TODO.md#m5-qar-module-adult-content) for implementation status.

---

## Checklist Summary

### Phase 1: Core Completion ✅
- [x] Migrate cache.go to rueidis
- [x] Add otter local cache
- [x] Register all modules in main.go
- [x] Add sturdyc for API response caching
- [x] Implement Radarr provider (full client)
- [x] Implement TMDb provider (already existed)
- [x] Implement central MetadataService

### Phase 1.5: Library Refactor ✅
- [x] Create per-module library tables (movie_libraries, tv_libraries)
- [x] Update FK: `movies.movie_library_id` → `movie_libraries(id)`
- [x] Update FK: `series.tv_library_id` → `tv_libraries(id)`
- [x] Update sqlc queries for new library columns
- [x] Create deprecation migration for shared libraries
- [x] Implement `LibraryProvider` interface in shared
- [x] Implement `LibraryService` for movie module

### Phase 2: Movie Module
- [x] Migration: 000001_movie_core.up.sql
- [x] Migration: 000002_movie_genres.up.sql
- [x] Migration: 000003_movie_people.up.sql
- [x] Migration: 000004_movie_images.up.sql
- [x] Migration: 000005_movie_libraries.up.sql
- [x] Migration: 000006_migrate_library_fk.up.sql
- [x] sqlc queries
- [x] entity.go
- [x] repository.go
- [x] service.go
- [x] library_service.go
- [ ] handler.go (API handlers)
- [x] jobs.go
- [x] module.go
- [x] Radarr client
- [x] TMDb client
- [ ] Unit tests

### Phase 3-9: Remaining Modules
See [TODO.md](../../../TODO.md) for high-level phases.

---

## References

- [Architecture](../architecture/01_ARCHITECTURE.md)
- [Content Modules Instructions](../../../.github/instructions/content-modules.instructions.md)
- [sqlc Instructions](../../../.github/instructions/sqlc-database.instructions.md)
- [fx Instructions](../../../.github/instructions/fx-dependency-injection.instructions.md)

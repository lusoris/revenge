# Content Module Development Instructions

> Instructions for developing content modules in the Revenge modular media server architecture.

## Module Structure

Every content module MUST follow this directory structure:

```
internal/
  content/
    {module}/
      entity.go           # Domain entities (Movie, Track, etc.)
      repository.go       # Repository interface
      repository_pg.go    # PostgreSQL implementation
      service.go          # Business logic
      handler.go          # HTTP handlers (implement ogen interfaces)
      scanner.go          # File scanner (if applicable)
      provider_{name}.go  # Metadata providers (tmdb, musicbrainz, etc.)
      jobs.go             # River job definitions
      module.go           # fx.Module registration
```

## Database Conventions

### Migration Files

Location: `internal/infra/database/migrations/{module}/`

Naming:

- `000001_{module}_core.up.sql` - Main content tables
- `000001_{module}_core.down.sql` - Rollback
- `000002_{module}_credits.up.sql` - Credits/cast (references shared or module-specific people)
- `000003_{module}_streams.up.sql` - Media streams
- `000004_{module}_user_data.up.sql` - Ratings, favorites, history

**People Tables Strategy:**
- **Video modules (movie, tvshow):** Use shared `video_people` from `shared/000017_video_people.up.sql`
  - Data overlaps 100% after background worker enrichment (TMDB, TVDB, IMDB)
  - `movie_credits` and `series_credits` reference `video_people`
- **Other modules (music, books, comics):** Module-specific people tables
  - Different metadata schemas (discography vs bibliography vs comic credits)
- **Adult module:** Completely isolated `c.performers` in schema `c`
  - NSFW images, different metadata sources (StashDB, ThePornDB)

### Table Naming

- Main table: plural (`movies`, `tracks`, `episodes`)
- Junction tables: `{parent}_{child}` (`movie_cast`, `track_artists`)
- User data: `{item}_user_ratings`, `{item}_favorites`, `{item}_history`
- External data: `{item}_external_ratings`

### Column Conventions

```sql
-- Always include
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()

-- For content items
library_id UUID NOT NULL REFERENCES libraries(id) ON DELETE CASCADE,
path TEXT NOT NULL,
title VARCHAR(500) NOT NULL,
sort_title VARCHAR(500),

-- For durations (use ticks: 1 tick = 100 nanoseconds)
runtime_ticks BIGINT,
duration_ticks BIGINT,
position_ticks BIGINT,

-- For external IDs (per provider, not JSONB)
tmdb_id INT,
imdb_id VARCHAR(20),
tvdb_id INT,
musicbrainz_id UUID,
```

### SQLC Queries

Location: `internal/infra/database/queries/{module}/`

Files:

- `{module}.sql` - Core CRUD operations
- `{module}_people.sql` - People/credits queries
- `{module}_ratings.sql` - Rating queries
- `{module}_user_data.sql` - User-specific queries

Query naming:

```sql
-- name: Get{Entity}ByID :one
-- name: List{Entities} :many
-- name: List{Entities}ByLibrary :many
-- name: Create{Entity} :one
-- name: Update{Entity} :one
-- name: Delete{Entity} :exec
-- name: Search{Entities} :many
```

## Entity Design

### Base Fields (embed in all entities)

```go
type BaseEntity struct {
    ID        uuid.UUID
    CreatedAt time.Time
    UpdatedAt time.Time
}

type ContentEntity struct {
    BaseEntity
    LibraryID uuid.UUID
    Path      string
    Title     string
    SortTitle string
}
```

### Module-Specific Entity

```go
// content/movie/entity.go
package movie

type Movie struct {
    ContentEntity

    // Movie-specific fields
    OriginalTitle string
    Tagline       string
    Overview      string
    RuntimeTicks  int64
    ReleaseDate   *time.Time
    Year          int
    Budget        int64
    Revenue       int64

    // External IDs
    TMDbID  *int
    IMDbID  *string

    // Relationships (loaded on demand)
    Cast    []MovieCastMember
    Crew    []MovieCrewMember
    Studios []MovieStudio
    Genres  []Genre
    Images  []MovieImage
}
```

## Repository Pattern

### Interface Definition

```go
// content/movie/repository.go
package movie

type Repository interface {
    // Core CRUD
    GetByID(ctx context.Context, id uuid.UUID) (*Movie, error)
    List(ctx context.Context, params ListParams) ([]*Movie, error)
    ListByLibrary(ctx context.Context, libraryID uuid.UUID, params ListParams) ([]*Movie, error)
    Create(ctx context.Context, movie *Movie) error
    Update(ctx context.Context, movie *Movie) error
    Delete(ctx context.Context, id uuid.UUID) error

    // Search
    Search(ctx context.Context, query string, params ListParams) ([]*Movie, error)

    // Relationships
    GetCast(ctx context.Context, movieID uuid.UUID) ([]MovieCastMember, error)
    GetCrew(ctx context.Context, movieID uuid.UUID) ([]MovieCrewMember, error)
    GetGenres(ctx context.Context, movieID uuid.UUID) ([]Genre, error)

    // User data
    GetUserRating(ctx context.Context, userID, movieID uuid.UUID) (*UserRating, error)
    SetUserRating(ctx context.Context, userID, movieID uuid.UUID, score float64) error
    IsFavorite(ctx context.Context, userID, movieID uuid.UUID) (bool, error)
    AddFavorite(ctx context.Context, userID, movieID uuid.UUID) error
    RemoveFavorite(ctx context.Context, userID, movieID uuid.UUID) error
    GetHistory(ctx context.Context, userID, movieID uuid.UUID) (*WatchHistory, error)
    UpdateHistory(ctx context.Context, userID, movieID uuid.UUID, positionTicks int64, completed bool) error
}
```

### PostgreSQL Implementation

```go
// content/movie/repository_pg.go
package movie

type pgRepository struct {
    pool    *pgxpool.Pool
    queries *db.Queries
}

func NewRepository(pool *pgxpool.Pool) Repository {
    return &pgRepository{
        pool:    pool,
        queries: db.New(pool),
    }
}
```

## Service Layer

```go
// content/movie/service.go
package movie

type Service struct {
    repo     Repository
    scanner  *Scanner
    tmdb     *TMDbProvider
    logger   *slog.Logger
}

func NewService(repo Repository, scanner *Scanner, tmdb *TMDbProvider, logger *slog.Logger) *Service {
    return &Service{
        repo:    repo,
        scanner: scanner,
        tmdb:    tmdb,
        logger:  logger.With(slog.String("module", "movie")),
    }
}

// Business logic methods
func (s *Service) GetMovie(ctx context.Context, id uuid.UUID) (*Movie, error)
func (s *Service) ListMovies(ctx context.Context, params ListParams) ([]*Movie, error)
func (s *Service) RefreshMetadata(ctx context.Context, id uuid.UUID) error
func (s *Service) RateMovie(ctx context.Context, userID, movieID uuid.UUID, score float64) error
```

## HTTP Handlers (ogen-generated)

Handlers implement interfaces generated by ogen from OpenAPI specs:

```go
// content/movie/handler.go
package movie

// Handler implements the ogen-generated MoviesHandler interface
type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

// Implement ogen-generated interface methods
func (h *Handler) ListMovies(ctx context.Context, params api.ListMoviesParams) (api.ListMoviesRes, error) {
    movies, err := h.service.ListMovies(ctx, params)
    if err != nil {
        return nil, err
    }
    return &api.MovieListResponse{Movies: movies}, nil
}

func (h *Handler) GetMovie(ctx context.Context, params api.GetMovieParams) (api.GetMovieRes, error) {
    movie, err := h.service.GetMovie(ctx, params.ID)
    if err != nil {
        if errors.Is(err, ErrMovieNotFound) {
            return &api.GetMovieNotFound{}, nil
        }
        return nil, err
    }
    return movie, nil
}
```

## River Jobs

Each module defines its background jobs:

```go
// content/movie/jobs.go
package movie

import (
    "github.com/riverqueue/river"
)

// ScanLibraryArgs - Scan a movie library for new files
type ScanLibraryArgs struct {
    LibraryID uuid.UUID `json:"library_id"`
    FullScan  bool      `json:"full_scan"`
}

func (ScanLibraryArgs) Kind() string { return "movie.scan_library" }

type ScanLibraryWorker struct {
    river.WorkerDefaults[ScanLibraryArgs]
    scanner *Scanner
}

func (w *ScanLibraryWorker) Work(ctx context.Context, job *river.Job[ScanLibraryArgs]) error {
    return w.scanner.Scan(ctx, job.Args.LibraryID, job.Args.FullScan)
}

// FetchMetadataArgs - Fetch metadata for a movie
type FetchMetadataArgs struct {
    MovieID uuid.UUID `json:"movie_id"`
}

func (FetchMetadataArgs) Kind() string { return "movie.fetch_metadata" }

type FetchMetadataWorker struct {
    river.WorkerDefaults[FetchMetadataArgs]
    service *Service
}

func (w *FetchMetadataWorker) Work(ctx context.Context, job *river.Job[FetchMetadataArgs]) error {
    return w.service.RefreshMetadata(ctx, job.Args.MovieID)
}

// IndexMovieArgs - Index movie in Typesense
type IndexMovieArgs struct {
    MovieID uuid.UUID `json:"movie_id"`
}

func (IndexMovieArgs) Kind() string { return "movie.index" }

type IndexMovieWorker struct {
    river.WorkerDefaults[IndexMovieArgs]
    search *search.Service
    repo   Repository
}

func (w *IndexMovieWorker) Work(ctx context.Context, job *river.Job[IndexMovieArgs]) error {
    movie, err := w.repo.GetByID(ctx, job.Args.MovieID)
    if err != nil {
        return err
    }
    return w.search.IndexMovie(ctx, movie)
}
```

## fx Module Registration

```go
// content/movie/module.go
package movie

import (
    "github.com/riverqueue/river"
    "go.uber.org/fx"
)

var Module = fx.Module("content/movie",
    // Repository
    fx.Provide(NewRepository),

    // Providers
    fx.Provide(NewTMDbProvider),

    // Scanner
    fx.Provide(NewScanner),

    // Service
    fx.Provide(NewService),

    // Handler (implements ogen interface)
    fx.Provide(NewHandler),

    // River workers
    fx.Provide(func(scanner *Scanner) *ScanLibraryWorker {
        return &ScanLibraryWorker{scanner: scanner}
    }),
    fx.Provide(func(service *Service) *FetchMetadataWorker {
        return &FetchMetadataWorker{service: service}
    }),

    // Register workers with River
    fx.Invoke(func(
        workers *river.Workers,
        scanWorker *ScanLibraryWorker,
        metaWorker *FetchMetadataWorker,
    ) {
        river.AddWorker(workers, scanWorker)
        river.AddWorker(workers, metaWorker)
    }),
)
```

## Adult Module Special Considerations

Adult modules use the `c` PostgreSQL schema (obscured name):

```go
// content/c/movie/repository_pg.go
package movie

// Queries reference c.* tables
const getMovieByID = `
SELECT id, title, path, ...
FROM c.movies
WHERE id = $1
`
```

Adult modules MUST:

- Use `c.` schema prefix for all tables
- Have completely isolated data (no FK to public schema)
- Have own performers, studios, tags tables
- Require `adult:read` scope for all endpoints
- Use `/api/v1/c/` namespace

## Testing

Each module should have:

```
internal/
  content/
    movie/
      service_test.go      # Service unit tests
      handler_test.go      # Handler unit tests
      repository_test.go   # Repository integration tests
```

Use table-driven tests:

```go
func TestMovieService_GetMovie(t *testing.T) {
    tests := []struct {
        name    string
        movieID uuid.UUID
        want    *Movie
        wantErr error
    }{
        // test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```
---

## Related

- [INDEX.instructions.md](INDEX.instructions.md) - Main instruction index with all cross-references
- [ARCHITECTURE_V2.md](../../docs/dev/design/architecture/ARCHITECTURE_V2.md) - System architecture
- [BEST_PRACTICES.md](../../docs/dev/design/operations/BEST_PRACTICES.md) - Best practices

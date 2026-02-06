# Movie Content Module

<!-- DESIGN: features/video -->

**Package**: `internal/content/movie`
**fx Module**: `movie.Module` + `moviejobs.Module`

> Movie content management with TMDb metadata, library scanning, file matching, search indexing, and watch progress tracking

---

## Module Structure

```
internal/content/movie/
├── service.go             # Service interface (23 methods) + movieService implementation
├── repository.go          # Repository interface (40 methods)
├── repository_postgres.go # PostgreSQL implementation (pgxpool + sqlc moviedb.Queries)
├── handler.go             # Internal handler (17 methods, called by API layer)
├── types.go               # Domain types: Movie, MovieFile, MovieCredit, MovieCollection, etc.
├── cached_service.go      # CachedService decorator (13 cached methods, L1+L2 cache)
├── library_service.go     # LibraryService: scan, match, refresh, probe
├── library_scanner.go     # Scanner wrapping shared FilesystemScanner
├── library_matcher.go     # Matcher with confidence scoring (title 60% + year 40%)
├── mediainfo.go           # FFmpeg media probing via go-astiav (Unix/Linux/macOS)
├── mediainfo_types.go     # MediaInfo, AudioStreamInfo, SubtitleStreamInfo types
├── mediainfo_windows.go   # Windows stub (no FFmpeg)
├── metadata_provider.go   # MetadataProvider interface (4 methods)
├── errors.go              # 5 sentinel errors
├── module.go              # fx.Module("movie") wiring
├── adapters/
│   ├── metadata_adapter.go  # TMDb client setup (4 req/sec, burst 10)
│   └── scanner_adapter.go   # MovieFileParser (regex: "Title (YEAR)", "Title.YEAR")
├── moviejobs/
│   ├── module.go            # fx.Module("moviejobs"), RegisterWorkers
│   ├── library_scan.go      # MovieLibraryScanWorker (30m timeout)
│   ├── metadata_refresh.go  # MovieMetadataRefreshWorker (5m timeout)
│   ├── file_match.go        # MovieFileMatchWorker (5m timeout)
│   └── search_index.go      # MovieSearchIndexWorker (15m timeout)
└── db/                      # sqlc-generated (moviedb package)
```

## Domain Types

### Movie

19 fields + i18n maps:

| Field Group | Fields |
|-------------|--------|
| Core | ID, Title, OriginalTitle, Year, ReleaseDate, Runtime |
| External IDs | TMDbID, IMDbID, RadarrID |
| Metadata | Status, OriginalLanguage, Overview, Tagline, TrailerURL |
| Media | PosterPath, BackdropPath, VoteAverage, VoteCount, Popularity, Budget, Revenue |
| Timestamps | CreatedAt, UpdatedAt, LibraryAddedAt, MetadataUpdatedAt |
| i18n | TitlesI18n, TaglinesI18n, OverviewsI18n (`map[string]string`), AgeRatings (`map[string]map[string]string`) |

Methods: `GetTitle(lang)`, `GetTagline(lang)`, `GetOverview(lang)`, `GetAgeRating(country, system)`, `GetAvailableLanguages()`, `GetAvailableAgeRatingCountries()`

### Supporting Types

| Type | Key Fields |
|------|-----------|
| MovieFile | FilePath, Resolution, VideoCodec, AudioCodec, BitrateKbps, DurationSeconds, AudioLanguages, SubtitleLanguages, RadarrFileID |
| MovieCredit | TMDbPersonID, Name, CreditType (cast/crew), Character, Job, Department, CastOrder |
| MovieCollection | TMDbCollectionID, Name, Overview, PosterPath, BackdropPath |
| MovieGenre | TMDbGenreID, Name |
| MovieWatched | UserID, MovieID, ProgressSeconds, DurationSeconds, ProgressPercent, IsCompleted, WatchCount |
| ContinueWatchingItem | Movie + progress fields + LastWatchedAt |
| WatchedMovieItem | Movie + WatchCount + LastWatchedAt |
| UserMovieStats | WatchedCount, InProgressCount, TotalWatches |

### Errors

`ErrMovieNotFound`, `ErrMovieFileNotFound`, `ErrProgressNotFound`, `ErrNotInCollection`, `ErrCollectionNotFound`

## Service Interface

23 exported methods on the `Service` interface:

**Movie CRUD**: GetMovie, GetMovieByTMDbID, GetMovieByIMDbID, ListMovies, SearchMovies, ListRecentlyAdded, ListTopRated, CreateMovie, UpdateMovie, DeleteMovie

**Files**: GetMovieFiles, CreateMovieFile, DeleteMovieFile

**Credits & Genres**: GetMovieCast, GetMovieCrew, GetMovieGenres, GetMoviesByGenre

**Collections**: GetMovieCollection, GetMoviesByCollection, GetCollectionForMovie

**Watch Progress**: UpdateWatchProgress, GetWatchProgress, MarkAsWatched, RemoveWatchProgress, GetContinueWatching, GetWatchHistory, GetUserStats

**Metadata**: RefreshMovieMetadata(ctx, id, opts)

Implementation: `movieService` struct with `repo Repository` + `metadataProvider MetadataProvider` fields.

## Repository Interface

40 exported methods on the `Repository` interface:

| Category | Count | Key Methods |
|----------|-------|-------------|
| Movie CRUD | 14 | GetMovie, GetMovieByTMDbID/IMDbID/RadarrID, ListMovies, CountMovies, SearchByTitle, SearchByTitleAnyLanguage, ListByYear, RecentlyAdded, TopRated, Create, Update, Delete |
| Files | 7 | CreateFile, GetFile, GetByPath, GetByRadarrID, ListByMovieID, UpdateFile, DeleteFile |
| Credits | 4 | CreateCredit, ListCast, ListCrew, DeleteCredits |
| Collections | 8 | Create, Get, GetByTMDbID, Update, AddMovie, RemoveMovie, ListMoviesByCollection, GetCollectionForMovie |
| Genres | 4 | AddGenre, ListGenres, DeleteGenres, ListMoviesByGenre |
| Watch Progress | 6 | CreateOrUpdateProgress, GetProgress, DeleteProgress, ListContinueWatching, ListWatchedMovies, GetUserStats |

Implementation: `postgresRepository` wrapping `pgxpool.Pool` + sqlc `moviedb.Queries`. Conversion helpers: `dbMovieToMovie()`, `dbMovieFileToMovieFile()`, etc. Custom JSON marshaling for i18n maps.

## CachedService

Decorator wrapping `Service` with `*cache.Cache` (L1 otter + L2 Dragonfly):

| Cached Method | TTL | Strategy |
|---------------|-----|----------|
| GetMovie | 5 min | Async cache |
| ListMovies | 1 min | Key: SHA256 of filters |
| ListRecentlyAdded | 2 min | Async cache |
| ListTopRated | 5 min | Async cache |
| GetMovieCast | 10 min | Async cache |
| GetMovieCrew | 10 min | Async cache |
| GetMovieGenres | 10 min | Async cache |
| GetMovieCollection | 10 min | Async cache |
| GetContinueWatching | 1 min | Per-user key |

Write operations (`UpdateMovie`, `DeleteMovie`, `UpdateWatchProgress`, `MarkAsWatched`) invalidate related cache patterns.

## LibraryService

Orchestrates scanning, matching, probing, and metadata refresh:

```go
type LibraryService struct {
    repo            Repository
    metadataService MetadataProvider
    scanner         *Scanner
    matcher         *Matcher
    prober          Prober  // FFmpeg via go-astiav
}
```

| Method | Purpose |
|--------|---------|
| ScanLibrary | Full library scan: discover files, match to movies, probe media info |
| RefreshMovie | Re-fetch metadata from TMDb for a specific movie |
| GetLibraryStats | Return counts (total movies, files, unmatched) |
| MatchFile | Match a single file path to a movie entity |

### Scanner

Wraps `shared/scanner.FilesystemScanner` with `MovieFileParser` (from adapters). Regex patterns:
- `Title (YEAR)` - parenthesized year
- `Title.YEAR` / `Title YEAR` / `Title_YEAR` - dot/space/underscore separated

### Matcher

Confidence scoring: **60% title similarity** (Levenshtein) + **40% year match** + popularity bonus. Thresholds: exact (1.0), title (>=0.8), fuzzy (>=0.5), unmatched (<0.5).

### MediaInfo Prober

FFmpeg integration via `go-astiav` (Unix/Linux/macOS only, Windows stub):
- Extracts: video codec, profile, resolution, framerate, HDR detection (SDR/HDR10/Dolby Vision/HLG)
- Audio: codec, channels, sample rate, language per stream
- Subtitles: language, codec, forced/default flags
- Output: `MediaInfo` struct with `ToMovieFileInfo()` converter

## MetadataProvider Interface

```go
type MetadataProvider interface {
    SearchMovies(ctx, query, year) ([]*Movie, error)
    EnrichMovie(ctx, mov, opts) error
    GetMovieCredits(ctx, movieID, tmdbID) ([]MovieCredit, error)
    GetMovieGenres(ctx, movieID, tmdbID) ([]MovieGenre, error)
    ClearCache() error
}
```

Implementation injected via `metadatafx` module as `MovieMetadataAdapter`. Uses shared `metadata.BaseClient` with TMDb API (4 req/sec rate limit, burst 10, sync.Map cache).

## Background Workers

4 River workers registered via `moviejobs.RegisterWorkers()`:

| Worker | Kind | Queue | Timeout | Trigger |
|--------|------|-------|---------|---------|
| MovieLibraryScanWorker | `movie_library_scan` | bulk | 30m | Library scan request |
| MovieMetadataRefreshWorker | `metadata_refresh_movie` | default | 5m | Metadata refresh job (from metadata service) |
| MovieFileMatchWorker | `movie_file_match` | default | 5m | Individual file match request |
| MovieSearchIndexWorker | `movie_search_index` | default | 15m | Index/remove/reindex operations |

SearchIndexWorker supports 3 operations: `index` (single movie), `remove` (from index), `reindex` (full). Fetches movie + genres + credits + files before indexing to Typesense.

## API Endpoints

21 endpoints in `internal/api/movie_handlers.go` (ogen-generated types):

| Endpoint | Handler Method |
|----------|---------------|
| `GET /movies` | ListMovies |
| `GET /movies/search` | SearchMovies |
| `GET /movies/recently-added` | GetRecentlyAdded |
| `GET /movies/top-rated` | GetTopRated |
| `GET /movies/continue-watching` | GetContinueWatching |
| `GET /movies/watch-history` | GetWatchHistory |
| `GET /movies/stats` | GetUserMovieStats |
| `GET /movies/{id}` | GetMovie |
| `GET /movies/{id}/files` | GetMovieFiles |
| `GET /movies/{id}/cast` | GetMovieCast |
| `GET /movies/{id}/crew` | GetMovieCrew |
| `GET /movies/{id}/genres` | GetMovieGenres |
| `GET /movies/{id}/collection` | GetMovieCollection |
| `GET /movies/{id}/progress` | GetWatchProgress |
| `PUT /movies/{id}/progress` | UpdateWatchProgress |
| `DELETE /movies/{id}/progress` | DeleteWatchProgress |
| `POST /movies/{id}/watched` | MarkAsWatched |
| `POST /movies/{id}/refresh` | RefreshMovieMetadata |
| `GET /movies/{id}/similar` | GetSimilarMovies |
| `GET /collections/{id}` | GetCollection |
| `GET /collections/{id}/movies` | GetCollectionMovies |

Converter functions in `movie_converters.go` bridge domain types to ogen API types.

## Dependencies

- `github.com/jackc/pgx/v5/pgxpool` - PostgreSQL (via repository)
- `github.com/imroc/req/v3` - HTTP client for TMDb (via shared metadata.BaseClient)
- `github.com/asticode/go-astiav` - FFmpeg bindings for media probing
- `github.com/riverqueue/river` - Background job processing
- `github.com/google/uuid` - UUID generation
- `github.com/shopspring/decimal` - Decimal types for ratings
- `go.uber.org/zap` - Structured logging
- `go.uber.org/fx` - Dependency injection
- Shared packages: `content/shared/scanner`, `content/shared/matcher`, `content/shared/metadata`, `content/shared/library`, `content/shared/jobs`

## fx Wiring

```go
// movie.Module provides:
fx.Provide(NewPostgresRepository)    // → Repository
fx.Provide(provideService)           // → Service (repo + metadataProvider)
fx.Provide(NewHandler)               // → *Handler
fx.Provide(provideLibraryService)    // → *LibraryService (repo + metadata + scanner + matcher + prober)

// moviejobs.Module provides:
fx.Provide(NewMovieLibraryScanWorker, NewMovieMetadataRefreshWorker,
           NewMovieFileMatchWorker, NewMovieSearchIndexWorker)
```

## Related Documentation

- [TVSHOW_MODULE.md](TVSHOW_MODULE.md) - TV show content module (similar architecture)
- [../../architecture/METADATA_SYSTEM.md](../../architecture/METADATA_SYSTEM.md) - Provider chain and caching
- [../../infrastructure/JOBS.md](../../infrastructure/JOBS.md) - River job queue setup
- [../../infrastructure/CACHE.md](../../infrastructure/CACHE.md) - L1/L2 caching infrastructure
- [../../infrastructure/SEARCH_INFRA.md](../../infrastructure/SEARCH_INFRA.md) - Typesense search
- [../../services/LIBRARY.md](../../services/LIBRARY.md) - Library management service

# Radarr Integration

<!-- DESIGN: integrations/servarr -->

**Package**: `internal/integration/radarr`
**fx Module**: `radarr.Module`

> Movie library synchronization, webhook processing, and Radarr API v3 client

---

## Module Structure

```
internal/integration/radarr/
├── client.go           # HTTP client (req, rate limiter, sync.Map cache)
├── types.go            # Radarr API v3 response types
├── mapper.go           # Radarr types → movie domain types
├── service.go          # SyncService (full + single movie sync)
├── jobs.go             # RadarrSyncWorker + RadarrWebhookWorker
├── webhook_handler.go  # Webhook event routing (10 event types)
├── errors.go           # Sentinel errors
├── module.go           # fx.Module("radarr") wiring
└── doc.go              # Package documentation
```

## Client

HTTP client wrapping `imroc/req` with rate limiting and caching:

```go
type Client struct {
    client      *req.Client
    baseURL     string
    apiKey      string
    rateLimiter *rate.Limiter
    cache       sync.Map
    cacheTTL    time.Duration
}
```

### API Methods

| Method | Radarr Endpoint | Description |
|--------|----------------|-------------|
| GetSystemStatus | `GET /api/v3/system/status` | Health check |
| GetAllMovies | `GET /api/v3/movie` | List all movies |
| GetMovie | `GET /api/v3/movie/{id}` | Get specific movie |
| GetMovieByTMDbID | `GET /api/v3/movie` (filter) | Lookup by TMDb ID |
| GetMovieFiles | `GET /api/v3/moviefile` | List movie files |
| GetQualityProfiles | `GET /api/v3/qualityprofile` | Quality profiles |
| GetRootFolders | `GET /api/v3/rootfolder` | Root folders |
| GetTags | `GET /api/v3/tag` | Tags |
| GetCalendar | `GET /api/v3/calendar` | Upcoming releases |
| GetHistory | `GET /api/v3/history` | Download history |
| AddMovie | `POST /api/v3/movie` | Add movie to Radarr |
| DeleteMovie | `DELETE /api/v3/movie/{id}` | Delete movie |
| RefreshMovie | `POST /api/v3/command` | Trigger metadata refresh |
| RescanMovie | `POST /api/v3/command` | Trigger file rescan |
| SearchMovie | `POST /api/v3/command` | Trigger download search |
| ClearCache | - | Flush client cache |
| IsHealthy | `GET /api/v3/system/status` | Quick health check |

## Mapper

Converts Radarr API types to movie domain types:

| Method | Converts |
|--------|----------|
| ToMovie | Radarr `Movie` → `movie.Movie` |
| ToMovieFile | Radarr `MovieFile` → `movie.MovieFile` |
| ToMovieCollection | Radarr `Collection` → `movie.MovieCollection` |
| ToGenres | Radarr `Movie.Genres` → `[]movie.MovieGenre` |

## SyncService

Orchestrates library synchronization between Radarr and Revenge:

```go
type SyncService struct {
    client     *Client
    mapper     *Mapper
    movieRepo  movie.Repository
    logger     *slog.Logger
    syncMu     sync.Mutex
    syncStatus SyncStatus
}
```

| Method | Purpose |
|--------|---------|
| SyncLibrary | Full library sync (add/update/remove movies) |
| SyncMovie | Sync single movie by Radarr ID |
| RefreshMovie | Trigger Radarr metadata refresh by TMDb ID |
| GetStatus | Return current sync status |
| IsHealthy | Check Radarr connectivity |

`SyncResult`: MoviesAdded, MoviesUpdated, MoviesRemoved, MoviesSkipped, Errors, Duration.

## Background Workers

2 River workers:

| Worker | Kind | Queue | Timeout | Purpose |
|--------|------|-------|---------|---------|
| RadarrSyncWorker | `radarr_sync` | high | 10m | Full or single-movie sync |
| RadarrWebhookWorker | `radarr_webhook` | high | 1m | Process webhook events |

## Webhook Handler

Routes Radarr webhook events to handlers:

| Event | Handler | Action |
|-------|---------|--------|
| `Grab` | handleGrab | Movie grabbed for download |
| `Download` | handleDownload | Movie downloaded, trigger sync |
| `Rename` | handleRename | File renamed |
| `MovieDelete` | handleMovieDelete | Movie deleted from Radarr |
| `MovieFileDelete` | handleMovieFileDelete | Movie file deleted |
| `Health` | handleHealthIssue | Radarr health warning |
| `HealthRestored` | - | Health restored |
| `ApplicationUpdate` | handleApplicationUpdate | Radarr updated |
| `ManualInteractionRequired` | handleManualInteraction | Manual action needed |
| `Test` | handleTest | Webhook test event |

## Errors

`ErrMovieNotFound`, `ErrNotConfigured`, `ErrUnauthorized`

## Dependencies

- `github.com/imroc/req/v3` - HTTP client
- `golang.org/x/time/rate` - Rate limiting
- `github.com/riverqueue/river` - Background jobs
- `github.com/google/uuid` - UUID generation
- `go.uber.org/zap` + `log/slog` - Logging
- `go.uber.org/fx` - Dependency injection
- `internal/content/movie` - Movie domain types and repository

## Configuration

From `config.go` (koanf namespace `radarr.*`):

```go
type Config struct {
    BaseURL   string        // Radarr base URL (e.g., http://localhost:7878)
    APIKey    string        // Radarr API key
    RateLimit rate.Limit    // Requests per second
    CacheTTL  time.Duration // Client cache TTL
    Timeout   time.Duration // HTTP timeout
}
```

## Related Documentation

- [SONARR.md](SONARR.md) - Sonarr integration (identical architecture for TV shows)
- [../../features/video/MOVIE_MODULE.md](../../features/video/MOVIE_MODULE.md) - Movie domain types
- [../../infrastructure/JOBS.md](../../infrastructure/JOBS.md) - River job queue

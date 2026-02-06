# Sonarr Integration

<!-- DESIGN: integrations/servarr -->

**Package**: `internal/integration/sonarr`
**fx Module**: `sonarr.Module`

> TV show library synchronization, webhook processing, and Sonarr API v3 client

---

## Module Structure

```
internal/integration/sonarr/
├── client.go           # HTTP client (req, rate limiter, sync.Map cache)
├── types.go            # Sonarr API v3 response types
├── mapper.go           # Sonarr types → tvshow domain types
├── service.go          # SyncService (full + single series sync)
├── jobs.go             # SonarrSyncWorker + SonarrWebhookWorker
├── webhook_handler.go  # Webhook event routing (11 event types)
├── errors.go           # Sentinel errors
├── module.go           # fx.Module("sonarr") wiring
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

| Method | Sonarr Endpoint | Description |
|--------|----------------|-------------|
| GetSystemStatus | `GET /api/v3/system/status` | Health check |
| GetAllSeries | `GET /api/v3/series` | List all series |
| GetSeries | `GET /api/v3/series/{id}` | Get specific series |
| GetSeriesByTVDbID | `GET /api/v3/series` (filter) | Lookup by TVDb ID |
| GetEpisodes | `GET /api/v3/episode` | List episodes by series |
| GetEpisode | `GET /api/v3/episode/{id}` | Get specific episode |
| GetEpisodesBySeason | `GET /api/v3/episode` (filter) | Episodes by season |
| GetEpisodeFiles | `GET /api/v3/episodefile` | List episode files |
| GetEpisodeFile | `GET /api/v3/episodefile/{id}` | Get specific file |
| GetQualityProfiles | `GET /api/v3/qualityprofile` | Quality profiles |
| GetRootFolders | `GET /api/v3/rootfolder` | Root folders |
| GetTags | `GET /api/v3/tag` | Tags |
| GetCalendar | `GET /api/v3/calendar` | Upcoming episodes |
| GetHistory | `GET /api/v3/history` | Download history |
| AddSeries | `POST /api/v3/series` | Add series to Sonarr |
| DeleteSeries | `DELETE /api/v3/series/{id}` | Delete series |
| RefreshSeries | `POST /api/v3/command` | Trigger metadata refresh |
| RescanSeries | `POST /api/v3/command` | Trigger file rescan |
| SearchSeries | `POST /api/v3/command` | Trigger series search |
| SearchSeason | `POST /api/v3/command` | Trigger season search |
| SearchEpisodes | `POST /api/v3/command` | Trigger episode search |
| ClearCache | - | Flush client cache |
| IsHealthy | `GET /api/v3/system/status` | Quick health check |

## Mapper

Converts Sonarr API types to tvshow domain types:

| Method | Converts |
|--------|----------|
| ToSeries | Sonarr `Series` → `tvshow.Series` |
| ToSeason | Sonarr `SeasonInfo` → `tvshow.Season` |
| ToEpisode | Sonarr `Episode` → `tvshow.Episode` |
| ToEpisodeFile | Sonarr `EpisodeFile` → `tvshow.EpisodeFile` |
| ToGenres | Sonarr `Series.Genres` → `[]tvshow.SeriesGenre` |

## SyncService

Orchestrates library synchronization between Sonarr and Revenge:

```go
type SyncService struct {
    client      *Client
    mapper      *Mapper
    tvshowRepo  tvshow.Repository
    logger      *slog.Logger
    syncMu      sync.Mutex
    syncStatus  SyncStatus
}
```

| Method | Purpose |
|--------|---------|
| SyncLibrary | Full library sync (series + seasons + episodes + files) |
| SyncSeries | Sync single series by Sonarr ID |
| GetStatus | Return current sync status |
| IsHealthy | Check Sonarr connectivity |

`SyncResult`: SeriesAdded/Updated/Removed/Skipped, EpisodesAdded/Updated, Errors, Duration.

## Background Workers

2 River workers:

| Worker | Kind | Queue | Timeout | Purpose |
|--------|------|-------|---------|---------|
| SonarrSyncWorker | `sonarr_sync` | high | 10m | Full or single-series sync |
| SonarrWebhookWorker | `sonarr_webhook` | high | 1m | Process webhook events |

## Webhook Handler

Routes Sonarr webhook events to handlers:

| Event | Handler | Action |
|-------|---------|--------|
| `Grab` | handleGrab | Episode grabbed for download |
| `Download` | handleDownload | Episode downloaded, trigger sync |
| `Rename` | handleRename | File renamed |
| `SeriesAdd` | handleSeriesAdd | New series added to Sonarr |
| `SeriesDelete` | handleSeriesDelete | Series deleted from Sonarr |
| `EpisodeFileDelete` | handleEpisodeFileDelete | Episode file deleted |
| `Health` | handleHealth | Sonarr health warning |
| `HealthRestored` | handleHealthRestored | Health restored |
| `ApplicationUpdate` | handleApplicationUpdate | Sonarr updated |
| `ManualInteractionRequired` | handleManualInteractionRequired | Manual action needed |
| `Test` | handleTest | Webhook test event |

## Sonarr-Specific Constants

Series types: `standard`, `daily`, `anime`
Series status: `continuing`, `ended`, `upcoming`, `deleted`
Monitor options: `all`, `future`, `missing`, `existing`, `pilot`, `firstSeason`, `lastSeason`, `monitorSpecials`, `unmonitorSpecials`, `none`

## Errors

`ErrSeriesNotFound`, `ErrEpisodeNotFound`, `ErrEpisodeFileNotFound`, `ErrUnauthorized`, `ErrConnectionFailed`, `ErrRateLimited`

## Dependencies

- `github.com/imroc/req/v3` - HTTP client
- `golang.org/x/time/rate` - Rate limiting
- `github.com/riverqueue/river` - Background jobs
- `github.com/google/uuid` - UUID generation
- `go.uber.org/zap` + `log/slog` - Logging
- `go.uber.org/fx` - Dependency injection
- `internal/content/tvshow` - TV show domain types and repository

## Configuration

From `config.go` (koanf namespace `sonarr.*`):

```go
type Config struct {
    BaseURL   string        // Sonarr base URL (e.g., http://localhost:8989)
    APIKey    string        // Sonarr API key
    RateLimit rate.Limit    // Requests per second
    CacheTTL  time.Duration // Client cache TTL
    Timeout   time.Duration // HTTP timeout
}
```

## Related Documentation

- [RADARR.md](RADARR.md) - Radarr integration (identical architecture for movies)
- [../../features/video/TVSHOW_MODULE.md](../../features/video/TVSHOW_MODULE.md) - TV show domain types
- [../../infrastructure/JOBS.md](../../infrastructure/JOBS.md) - River job queue

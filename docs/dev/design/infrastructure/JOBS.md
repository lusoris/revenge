# Job Queue Infrastructure

<!-- DESIGN: infrastructure -->

**Package**: `internal/infra/jobs`
**fx Module**: `jobs.Module`

> River-based background job processing with 5-level priority queues and 17 workers

---

## Service Structure

```
internal/infra/jobs/
├── river.go               # Client wrapper (Insert, Start/Stop, Subscribe)
├── module.go              # fx module, lifecycle hooks
├── queues.go              # 5-level priority queues, backoff strategies
├── progress.go            # JobProgress tracking for frontend polling
├── notification_job.go    # NotificationWorker (dispatches to agents)
└── cleanup_job.go         # Generic CleanupWorker (leader-aware)
```

## Queue Configuration

| Queue | Workers | Purpose |
|-------|---------|---------|
| `critical` | 20 | Security events, auth failures, urgent tasks |
| `high` | 15 | Notifications, webhooks, user actions |
| `default` | 10 | Metadata fetching, sync operations |
| `low` | 5 | Cleanup, maintenance, session pruning |
| `bulk` | 3 | Library scans, batch operations, reindexing |

## Client Interface

```go
type Client struct { /* wraps river.Client[pgx.Tx] */ }

func NewClient(pool *pgxpool.Pool, workers *river.Workers, cfg Config) (*Client, error)
func (c *Client) Start(ctx context.Context) error
func (c *Client) Stop(ctx context.Context) error
func (c *Client) Insert(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) (*river.JobInsertResult, error)
func (c *Client) InsertMany(ctx context.Context, params []river.InsertManyParams) ([]*river.JobInsertResult, error)
func (c *Client) JobGet(ctx context.Context, id int64) (*river.Job[any], error)
func (c *Client) JobCancel(ctx context.Context, id int64) (*river.Job[any], error)
func (c *Client) Subscribe(handler river.EventHandler) error
```

## Progress Tracking

```go
type JobProgress struct {
    Phase      string  // Current phase name
    Current    int64   // Items processed
    Total      int64   // Total items
    Percentage float64 // Auto-calculated
    Message    string  // Status message
}

func ReportProgress(ctx context.Context, job river.JobRow, progress JobProgress) error
func GetJobProgress(ctx context.Context, client *Client, jobID int64) (*JobProgress, error)
```

## Backoff Strategies

```go
func ExponentialBackoff(attempt int) time.Duration  // min(1s * 2^attempt, 1h)
func LinearBackoff(attempt int) time.Duration       // min(30s * attempt, 30m)
```

Default: 25 max attempts with exponential backoff.

## All 17 Workers

### Infrastructure Workers (in `internal/infra/jobs/`)

| Worker | Kind | Queue | Timeout | Purpose |
|--------|------|-------|---------|---------|
| NotificationWorker | `notification` | high | 2m | Dispatch events to notification agents |
| CleanupWorker | `cleanup` | low | 2m | Generic cleanup (leader-aware) |

### Service Workers

| Worker | Kind | Location | Queue | Timeout |
|--------|------|----------|-------|---------|
| LibraryScanCleanup | `library_scan_cleanup` | `service/library` | low | 2m |
| ActivityCleanup | `activity_cleanup` | `service/activity` | low | 2m |

### Movie Workers (in `content/movie/moviejobs/`)

| Worker | Kind | Queue | Timeout | Purpose |
|--------|------|-------|---------|---------|
| MovieLibraryScan | `movie_library_scan` | bulk | 30m | Scan directories for movie files |
| MovieMetadataRefresh | `metadata_refresh_movie` | default | 5m | Refresh movie metadata from TMDb |
| MovieFileMatch | `movie_file_match` | default | 5m | Match files to movie entities |
| MovieSearchIndex | `movie_search_index` | default | 15m | Index/remove/reindex movies in Typesense |

### TV Show Workers (in `content/tvshow/jobs/`)

| Worker | Kind | Queue | Timeout | Purpose |
|--------|------|-------|---------|---------|
| TVShowLibraryScan | `tvshow_library_scan` | bulk | 30m | Scan directories for TV show files |
| TVShowMetadataRefresh | `tvshow_metadata_refresh` | default | 15m | Refresh series/season/episode metadata |
| TVShowFileMatch | `tvshow_file_match` | default | 5m | Match files to episodes |
| TVShowSearchIndex | `tvshow_search_index` | bulk | 10m | Index series for search (stub) |
| TVShowSeriesRefresh | `tvshow_series_refresh` | default | 10m | Refresh single series cascade |

### Integration Workers

| Worker | Kind | Location | Queue | Timeout | Purpose |
|--------|------|----------|-------|---------|---------|
| RadarrSync | `radarr_sync` | `integration/radarr` | high | 10m | Sync Radarr library |
| RadarrWebhook | `radarr_webhook` | `integration/radarr` | high | 1m | Process Radarr webhooks |
| SonarrSync | `sonarr_sync` | `integration/sonarr` | high | 10m | Sync Sonarr library |
| SonarrWebhook | `sonarr_webhook` | `integration/sonarr` | high | 1m | Process Sonarr webhooks |

### Metadata Job Args (in `service/metadata/jobs/`)

Defines 7 job arg types + batch enqueue helpers for metadata refresh operations. Workers are implemented in content modules.

## Configuration

From `config.go` `JobsConfig` (koanf namespace `jobs.*`):
```yaml
jobs:
  max_workers: 10            # Per-queue default
  fetch_cooldown: 100ms      # Poll interval
```

## Dependencies

- `github.com/riverqueue/river` - Job queue framework
- `github.com/jackc/pgx/v5/pgxpool` - PostgreSQL (River uses pgx for job storage)

## Related Documentation

- [DATABASE.md](DATABASE.md) - River stores jobs in PostgreSQL
- [../services/NOTIFICATION.md](../services/NOTIFICATION.md) - NotificationWorker dispatches to agents

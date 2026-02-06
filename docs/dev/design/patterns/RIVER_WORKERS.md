# River Worker Guide

> How to add background job workers. Written from code as of 2026-02-06.

---

## Architecture

```
Handler/Service                    River (PostgreSQL-backed)              Worker
     |                                    |                                 |
     +-- client.Insert(args) ----------->  job stored in river_job table     |
                                          |                                 |
                                          +-- poll + dequeue ------------> Work(ctx, job)
                                          |                                 |
                                          +-- on error: retry with backoff  |
                                          +-- on success: mark completed    |
```

**Queue system** (5-tier, defined in `internal/infra/jobs/queues.go`):

| Queue | Max Workers | Purpose | Use for |
|-------|-------------|---------|---------|
| `critical` | 20 | Security events, auth failures | Urgent system tasks |
| `high` | 15 | Notifications, webhooks | User-initiated, external integrations |
| `default` | 10 | Metadata refresh, file matching | General async work |
| `low` | 5 | Cleanup, maintenance | Leader-aware maintenance |
| `bulk` | 3 | Library scans, search reindexing | Resource-intensive batch ops |

---

## Step-by-Step: Add a New Worker

### 1. Define Job Args

```go
package mymodule

import (
    "github.com/riverqueue/river"
    infrajobs "github.com/lusoris/revenge/internal/infra/jobs"
)

const MyJobKind = "my_job"

type MyJobArgs struct {
    EntityID  uuid.UUID `json:"entity_id"`
    Force     bool      `json:"force,omitempty"`
}

func (MyJobArgs) Kind() string {
    return MyJobKind
}

// InsertOpts sets the default queue for this job type.
func (MyJobArgs) InsertOpts() river.InsertOpts {
    return river.InsertOpts{
        Queue: infrajobs.QueueDefault,
    }
}
```

**Rules:**
- All fields must be JSON-serializable (struct tags required)
- `Kind()` returns a unique string identifier — use a package-level constant
- `InsertOpts()` sets the default queue — choose based on the table above
- Optional: set `MaxAttempts` in InsertOpts to override the global default (25)

### 2. Create the Worker

```go
type MyWorker struct {
    river.WorkerDefaults[MyJobArgs]
    myService  *MyService
    logger     *zap.Logger
}

func NewMyWorker(svc *MyService, logger *zap.Logger) *MyWorker {
    return &MyWorker{
        myService: svc,
        logger:    logger.Named("my-worker"),
    }
}

func (w *MyWorker) Timeout(job *river.Job[MyJobArgs]) time.Duration {
    return 5 * time.Minute
}

func (w *MyWorker) Work(ctx context.Context, job *river.Job[MyJobArgs]) error {
    args := job.Args

    w.logger.Info("starting job",
        zap.String("entity_id", args.EntityID.String()),
        zap.Bool("force", args.Force),
        zap.Int("attempt", job.Attempt),
    )

    if err := w.myService.DoWork(ctx, args.EntityID); err != nil {
        w.logger.Error("job failed", zap.Error(err))
        return fmt.Errorf("my job failed: %w", err)
    }

    w.logger.Info("job completed")
    return nil
}
```

**Timeout guidelines:**

| Job type | Timeout |
|----------|---------|
| Webhooks | 1-2 min |
| Cleanup / maintenance | 2 min |
| Metadata refresh | 5-15 min |
| External API sync | 10 min |
| Library scan / batch | 30 min |

**Error handling:** Return an error to trigger River's retry (exponential backoff, max 25 attempts). Return `nil` for success or when skipping (e.g., not leader).

### 3. Register via fx Module

```go
// module.go
package mymodule

import "go.uber.org/fx"

var Module = fx.Module("mymodule",
    fx.Provide(NewMyWorker),
)
```

Workers must be registered with River's worker registry. This happens in the content/integration module's registration function:

```go
func RegisterWorkers(workers *river.Workers, myWorker *MyWorker) error {
    river.AddWorker(workers, myWorker)
    return nil
}
```

Then wire into `internal/app/module.go`:

```go
var Module = fx.Module("app",
    // ...
    mymodule.Module,
    // ...
)
```

### 4. Enqueue Jobs

**Simple insert** (from a handler or service):

```go
_, err := h.riverClient.Insert(ctx, &mymodule.MyJobArgs{
    EntityID: entityID,
    Force:    true,
}, nil)  // nil = use default InsertOpts from args
```

**With uniqueness constraint** (prevent duplicate jobs):

```go
_, err := client.Insert(ctx, MyJobArgs{
    EntityID: id,
}, &river.InsertOpts{
    UniqueOpts: river.UniqueOpts{
        ByArgs:   true,
        ByPeriod: 24 * time.Hour,
    },
})
```

**Batch insert** (for bulk operations):

```go
params := make([]river.InsertManyParams, len(ids))
for i, id := range ids {
    params[i] = river.InsertManyParams{
        Args: MyJobArgs{EntityID: id},
    }
}
_, err := client.InsertMany(ctx, params)
```

---

## Progress Tracking (optional)

For long-running jobs, report progress so the frontend can poll:

```go
func (w *MyWorker) Work(ctx context.Context, job *river.Job[MyJobArgs]) error {
    items := loadItems()

    for i, item := range items {
        if err := processItem(ctx, item); err != nil {
            return err
        }

        _ = w.jobClient.ReportProgress(ctx, job.ID, &infrajobs.JobProgress{
            Phase:   "processing",
            Current: i + 1,
            Total:   len(items),
            Message: item.Name,
        })
    }
    return nil
}
```

`JobProgress` auto-calculates `Percent` when `Total > 0`. Poll via `client.GetJobProgress(ctx, jobID)`.

---

## Leader-Aware Workers

For maintenance jobs that should only run on one node in a cluster:

```go
func (w *CleanupWorker) Work(ctx context.Context, job *river.Job[CleanupArgs]) error {
    if w.leaderElection != nil && !w.leaderElection.IsLeader() {
        w.logger.Info("skipping: not the leader node")
        return nil  // Complete without error — don't retry
    }
    // ... actual work ...
}
```

Used by: `CleanupWorker`, `ActivityCleanupWorker`, `LibraryScanCleanupWorker`.

---

## All 17 Workers

### Infrastructure (`internal/infra/jobs/`)

| Worker | Kind | Queue | Timeout |
|--------|------|-------|---------|
| NotificationWorker | `notification` | high | 2m |
| CleanupWorker | `cleanup` | low | 2m |

### Service (`internal/service/`)

| Worker | Kind | Queue | Timeout |
|--------|------|-------|---------|
| ActivityCleanupWorker | `activity_cleanup` | low | 2m |
| LibraryScanCleanupWorker | `library_scan_cleanup` | low | 2m |

### Movie (`internal/content/movie/moviejobs/`)

| Worker | Kind | Queue | Timeout |
|--------|------|-------|---------|
| MovieLibraryScanWorker | `movie_library_scan` | bulk | 30m |
| MovieMetadataRefreshWorker | `metadata_refresh_movie` | default | 5m |
| MovieFileMatchWorker | `movie_file_match` | default | 5m |
| MovieSearchIndexWorker | `movie_search_index` | default | 15m |

### TV Show (`internal/content/tvshow/jobs/`)

| Worker | Kind | Queue | Timeout |
|--------|------|-------|---------|
| LibraryScanWorker | `tvshow_library_scan` | bulk | 30m |
| MetadataRefreshWorker | `tvshow_metadata_refresh` | default | 15m |
| FileMatchWorker | `tvshow_file_match` | default | 5m |
| SearchIndexWorker | `tvshow_search_index` | bulk | 10m |
| SeriesRefreshWorker | `tvshow_series_refresh` | default | 10m |

### Integration (`internal/integration/`)

| Worker | Kind | Queue | Timeout |
|--------|------|-------|---------|
| RadarrSyncWorker | `radarr_sync` | high | 10m |
| RadarrWebhookWorker | `radarr_webhook` | high | 1m |
| SonarrSyncWorker | `sonarr_sync` | high | 10m |
| SonarrWebhookWorker | `sonarr_webhook` | high | 1m |

---

## Configuration

From `config.go` (`jobs.*` koanf namespace):

```yaml
jobs:
  max_workers: 10       # Per-queue default (overridden by DefaultQueueConfig)
  fetch_cooldown: 100ms # Poll interval
```

Global max attempts: 25 (set in `module.go`). Override per job via `InsertOpts{MaxAttempts: N}`.

Retry: exponential backoff built into River. Custom backoff available via `queues.go`:
- `ExponentialBackoff(attempt)` — `min(1s * 2^attempt, 1h)`
- `LinearBackoff(attempt)` — `min(30s * attempt, 30m)`

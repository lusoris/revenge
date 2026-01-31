# River Integration

<!-- SOURCES: dragonfly, fx, google-uuid, pgx, pgxpool, postgresql-arrays, postgresql-json, prometheus, prometheus-metrics, river, river-docs, rueidis, rueidis-docs, typesense, typesense-go -->

<!-- DESIGN: integrations/infrastructure, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> PostgreSQL-native job queue


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
- [Developer Resources](#developer-resources)
- [Configuration](#configuration)
- [Job Definitions](#job-definitions)
  - [Library Scan Job](#library-scan-job)
  - [Metadata Fetch Job](#metadata-fetch-job)
  - [Transcode Job](#transcode-job)
  - [Scheduled Jobs](#scheduled-jobs)
- [River Client Setup](#river-client-setup)
- [Worker Registration](#worker-registration)
- [Enqueueing Jobs](#enqueueing-jobs)
  - [Simple Enqueue](#simple-enqueue)
  - [Transactional Enqueue](#transactional-enqueue)
  - [Delayed Enqueue](#delayed-enqueue)
  - [Bulk Enqueue](#bulk-enqueue)
- [Implementation Checklist](#implementation-checklist)
- [Database Migration](#database-migration)
- [River UI](#river-ui)
- [Monitoring](#monitoring)
  - [Key Metrics](#key-metrics)
  - [Query Job Stats](#query-job-stats)
- [Error Handling](#error-handling)
  - [Retry Configuration](#retry-configuration)
  - [Error States](#error-states)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documentation](#related-documentation)

<!-- TOC-END -->

## Status

| Dimension | Status |
|-----------|--------|
| Design | ✅ |
| Sources | ✅ |
| Instructions | ✅ |
| Code | 🔴 |
| Linting | 🔴 |
| Unit Testing | 🔴 |
| Integration Testing | 🔴 |**Priority**: 🟡 MEDIUM (Phase 1 - Core Infrastructure)
**Type**: Background job queue

---

## Overview

River is Revenge's background job processing system. It uses PostgreSQL as the job store, eliminating the need for additional infrastructure like RabbitMQ or Redis-based queues.

**Key Features**:
- PostgreSQL-native (ACID guarantees)
- Type-safe job arguments (Go generics)
- Automatic retries with exponential backoff
- Scheduled jobs (cron-like)
- Job priorities and queues
- Transactional job enqueueing
- Built-in job observability

**Why River?**:
- No additional infrastructure (uses existing PostgreSQL)
- Transactional enqueueing (job inserted with data in same tx)
- Type-safe Go API
- Built for Go from the ground up
- Active development

---

## Developer Resources

- 📚 **Docs**: https://riverqueue.com/docs
- 🔗 **GitHub**: https://github.com/riverqueue/river
- 📖 **Go Docs**: https://pkg.go.dev/github.com/riverqueue/river
- 🔗 **River UI**: https://github.com/riverqueue/riverui

---

## Configuration

```yaml
# configs/config.yaml
jobs:
  enabled: true
  driver: "river"

  # Worker settings
  workers:
    count: 10  # Concurrent workers per queue

  # Queue definitions
  queues:
    default:
      workers: 5
    high:
      workers: 10
    low:
      workers: 2
    scheduled:
      workers: 2

  # Retry settings
  retry:
    max_attempts: 5
    backoff_multiplier: 2
    initial_interval: "1s"
    max_interval: "1h"

  # Job cleanup
  cleanup:
    completed_jobs_retention: "7d"
    cancelled_jobs_retention: "24h"
    discarded_jobs_retention: "30d"
```

---

## Job Definitions

### Library Scan Job

```go
package jobs

import (
    "context"
    "github.com/google/uuid"
    "github.com/riverqueue/river"
)

// ScanLibraryArgs defines arguments for library scanning
type ScanLibraryArgs struct {
    LibraryID uuid.UUID `json:"library_id"`
    FullScan  bool      `json:"full_scan"`
}

func (ScanLibraryArgs) Kind() string { return "library.scan" }

// InsertOpts returns job-specific options
func (a ScanLibraryArgs) InsertOpts() river.InsertOpts {
    return river.InsertOpts{
        Queue:    "default",
        Priority: 2,
        UniqueOpts: river.UniqueOpts{
            ByArgs:  true,
            ByState: []rivertype.JobState{rivertype.JobStateAvailable, rivertype.JobStateRunning},
        },
    }
}

// ScanLibraryWorker processes library scan jobs
type ScanLibraryWorker struct {
    river.WorkerDefaults[ScanLibraryArgs]
    libraryService *library.Service
    logger         *slog.Logger
}

func (w *ScanLibraryWorker) Work(ctx context.Context, job *river.Job[ScanLibraryArgs]) error {
    w.logger.Info("starting library scan",
        "library_id", job.Args.LibraryID,
        "full_scan", job.Args.FullScan,
    )

    return w.libraryService.Scan(ctx, job.Args.LibraryID, job.Args.FullScan)
}
```

### Metadata Fetch Job

```go
type FetchMetadataArgs struct {
    ItemID     uuid.UUID `json:"item_id"`
    ItemType   string    `json:"item_type"`  // movie, tvshow, album, etc.
    ExternalID string    `json:"external_id"`
    Provider   string    `json:"provider"`   // tmdb, musicbrainz, etc.
}

func (FetchMetadataArgs) Kind() string { return "metadata.fetch" }

func (a FetchMetadataArgs) InsertOpts() river.InsertOpts {
    return river.InsertOpts{
        Queue:    "default",
        Priority: 3,
    }
}

type FetchMetadataWorker struct {
    river.WorkerDefaults[FetchMetadataArgs]
    metadataService *metadata.Service
}

func (w *FetchMetadataWorker) Work(ctx context.Context, job *river.Job[FetchMetadataArgs]) error {
    return w.metadataService.Fetch(ctx, job.Args.ItemID, job.Args.Provider, job.Args.ExternalID)
}
```

### Transcode Job

```go
type TranscodeArgs struct {
    SessionID string    `json:"session_id"`
    MediaID   uuid.UUID `json:"media_id"`
    Profile   string    `json:"profile"`  // 1080p, 720p, etc.
}

func (TranscodeArgs) Kind() string { return "media.transcode" }

func (a TranscodeArgs) InsertOpts() river.InsertOpts {
    return river.InsertOpts{
        Queue:    "high",  // High priority queue
        Priority: 1,
    }
}
```

### Scheduled Jobs

```go
// Clean up expired sessions every hour
type CleanupSessionsArgs struct{}

func (CleanupSessionsArgs) Kind() string { return "cleanup.sessions" }

func (a CleanupSessionsArgs) InsertOpts() river.InsertOpts {
    return river.InsertOpts{
        Queue: "scheduled",
    }
}

// Register periodic job
periodicJobs := []*river.PeriodicJob{
    river.NewPeriodicJob(
        river.PeriodicInterval(1*time.Hour),
        func() (river.JobArgs, *river.InsertOpts) {
            return CleanupSessionsArgs{}, nil
        },
        &river.PeriodicJobOpts{RunOnStart: true},
    ),
    river.NewPeriodicJob(
        river.PeriodicInterval(24*time.Hour),
        func() (river.JobArgs, *river.InsertOpts) {
            return RefreshMetadataArgs{}, nil
        },
        nil,
    ),
}
```

---

## River Client Setup

```go
package jobs

import (
    "context"
    "log/slog"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/riverqueue/river"
    "github.com/riverqueue/river/riverdriver/riverpgxv5"
    "go.uber.org/fx"
)

type Client struct {
    *river.Client[pgx.Tx]
    pool   *pgxpool.Pool
    logger *slog.Logger
}

func NewClient(
    lc fx.Lifecycle,
    pool *pgxpool.Pool,
    logger *slog.Logger,
    workers *river.Workers,
    periodicJobs []*river.PeriodicJob,
) (*Client, error) {
    riverClient, err := river.NewClient(riverpgxv5.New(pool), &river.Config{
        Queues: map[string]river.QueueConfig{
            "default":   {MaxWorkers: 5},
            "high":      {MaxWorkers: 10},
            "low":       {MaxWorkers: 2},
            "scheduled": {MaxWorkers: 2},
        },
        Workers:      workers,
        PeriodicJobs: periodicJobs,
        Logger:       slogadapter.New(logger),
    })
    if err != nil {
        return nil, err
    }

    client := &Client{
        Client: riverClient,
        pool:   pool,
        logger: logger,
    }

    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            return riverClient.Start(ctx)
        },
        OnStop: func(ctx context.Context) error {
            return riverClient.Stop(ctx)
        },
    })

    return client, nil
}
```

---

## Worker Registration

```go
func NewWorkers(
    scanWorker *ScanLibraryWorker,
    metadataWorker *FetchMetadataWorker,
    transcodeWorker *TranscodeWorker,
    cleanupWorker *CleanupSessionsWorker,
) *river.Workers {
    workers := river.NewWorkers()

    river.AddWorker(workers, scanWorker)
    river.AddWorker(workers, metadataWorker)
    river.AddWorker(workers, transcodeWorker)
    river.AddWorker(workers, cleanupWorker)

    return workers
}
```

---

## Enqueueing Jobs

### Simple Enqueue

```go
func (s *LibraryService) RequestScan(ctx context.Context, libraryID uuid.UUID) error {
    _, err := s.riverClient.Insert(ctx, ScanLibraryArgs{
        LibraryID: libraryID,
        FullScan:  false,
    }, nil)
    return err
}
```

### Transactional Enqueue

```go
func (s *MovieService) Create(ctx context.Context, movie *Movie) error {
    tx, err := s.pool.Begin(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)

    // Insert movie
    if err := s.queries.WithTx(tx).InsertMovie(ctx, movie); err != nil {
        return err
    }

    // Enqueue metadata fetch in same transaction
    _, err = s.riverClient.InsertTx(ctx, tx, FetchMetadataArgs{
        ItemID:   movie.ID,
        ItemType: "movie",
        Provider: "tmdb",
    }, nil)
    if err != nil {
        return err
    }

    return tx.Commit(ctx)
}
```

### Delayed Enqueue

```go
// Schedule job for later
_, err := s.riverClient.Insert(ctx, args, &river.InsertOpts{
    ScheduledAt: time.Now().Add(5 * time.Minute),
})
```

### Bulk Enqueue

```go
params := make([]river.InsertManyParams, len(items))
for i, item := range items {
    params[i] = river.InsertManyParams{
        Args: FetchMetadataArgs{
            ItemID:   item.ID,
            ItemType: item.Type,
        },
    }
}

_, err := s.riverClient.InsertMany(ctx, params)
```

---

## Implementation Checklist

- [ ] **River Client** (`internal/infra/jobs/client.go`)
  - [ ] Client initialization
  - [ ] Queue configuration
  - [ ] Lifecycle management

- [ ] **Worker Registration** (`internal/infra/jobs/workers.go`)
  - [ ] fx module
  - [ ] Worker providers

- [ ] **Library Jobs** (`internal/infra/jobs/library/`)
  - [ ] ScanLibraryJob
  - [ ] RefreshMetadataJob
  - [ ] AnalyzeMediaJob

- [ ] **Metadata Jobs** (`internal/infra/jobs/metadata/`)
  - [ ] FetchMetadataJob
  - [ ] DownloadImagesJob
  - [ ] RefreshExternalIDsJob

- [ ] **Playback Jobs** (`internal/infra/jobs/playback/`)
  - [ ] TranscodeJob
  - [ ] GenerateThumbnailsJob
  - [ ] ExtractChaptersJob

- [ ] **Scheduled Jobs** (`internal/infra/jobs/scheduled/`)
  - [ ] CleanupSessionsJob
  - [ ] RefreshTokensJob
  - [ ] DatabaseMaintenanceJob

---

## Database Migration

```sql
-- River requires its own schema
-- Run: go run github.com/riverqueue/river/cmd/river migrate-up --database-url=...

-- Or include in migrations
CREATE SCHEMA IF NOT EXISTS river;

-- River will create its tables automatically
-- Main tables: river_job, river_leader, river_queue
```

---

## River UI

```yaml
# docker-compose.yml
services:
  river-ui:
    image: ghcr.io/riverqueue/riverui:latest
    container_name: revenge-river-ui
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgresql://revenge:password@postgres:5432/revenge
    depends_on:
      - postgres
```

---

## Monitoring

### Key Metrics

```go
var (
    jobsEnqueued = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "revenge_jobs_enqueued_total",
            Help: "Total jobs enqueued",
        },
        []string{"kind", "queue"},
    )

    jobsCompleted = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "revenge_jobs_completed_total",
            Help: "Total jobs completed",
        },
        []string{"kind", "queue", "status"},
    )

    jobDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "revenge_job_duration_seconds",
            Help:    "Job execution duration",
            Buckets: []float64{.1, .5, 1, 5, 10, 30, 60, 300},
        },
        []string{"kind"},
    )

    queueDepth = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "revenge_job_queue_depth",
            Help: "Current queue depth",
        },
        []string{"queue", "state"},
    )
)
```

### Query Job Stats

```sql
-- Pending jobs by queue
SELECT queue, count(*)
FROM river_job
WHERE state = 'available'
GROUP BY queue;

-- Failed jobs in last hour
SELECT kind, count(*), max(errors) as last_error
FROM river_job
WHERE state = 'discarded'
  AND finalized_at > now() - interval '1 hour'
GROUP BY kind;

-- Job throughput
SELECT
    date_trunc('minute', finalized_at) as minute,
    count(*) as completed
FROM river_job
WHERE state = 'completed'
  AND finalized_at > now() - interval '1 hour'
GROUP BY 1
ORDER BY 1;
```

---

## Error Handling

### Retry Configuration

```go
func (a FetchMetadataArgs) InsertOpts() river.InsertOpts {
    return river.InsertOpts{
        MaxAttempts: 5,
    }
}

// Worker can return specific errors
func (w *FetchMetadataWorker) Work(ctx context.Context, job *river.Job[FetchMetadataArgs]) error {
    err := w.fetch(ctx, job.Args)

    if errors.Is(err, ErrRateLimited) {
        // Snooze job for later
        return river.JobSnooze(5 * time.Minute)
    }

    if errors.Is(err, ErrNotFound) {
        // Cancel job, don't retry
        return river.JobCancel(err)
    }

    return err  // Will retry with backoff
}
```

### Error States

| State | Description |
|-------|-------------|
| `available` | Ready to run |
| `running` | Currently executing |
| `completed` | Successfully finished |
| `retryable` | Failed, will retry |
| `scheduled` | Waiting for scheduled time |
| `discarded` | Failed all retries |
| `cancelled` | Manually cancelled |

---


## Related Documentation

- [PostgreSQL](POSTGRESQL.md)
- [Dragonfly](DRAGONFLY.md)
- [Typesense](TYPESENSE.md)

---
applyTo: "**/internal/infra/jobs/**/*.go,**/internal/content/**/jobs.go"
---

# River Job Queue - Patterns & Best Practices

> PostgreSQL-native job queue using River for background processing.
> Source: https://riverqueue.com/docs

## Overview

River uses PostgreSQL as the job queue backend, providing ACID guarantees and eliminating the need for external job queue infrastructure. Jobs are retried up to 25 times by default with exponential backoff (`attempts^4 + rand(±10%)` seconds).

## Installation

```bash
go get github.com/riverqueue/river
go get github.com/riverqueue/river/riverdriver/riverpgxv5
```

Run migrations:

```bash
go install github.com/riverqueue/river/cmd/river@latest
river migrate-up --database-url "$DATABASE_URL"
```

## Job Args and Workers

Each job requires a `JobArgs` struct with `Kind()` method and a `Worker[T]`:

```go
// Job arguments - serialized to JSON
type SortArgs struct {
    Strings []string `json:"strings"`
}

func (SortArgs) Kind() string { return "sort" }

// Worker implementation
type SortWorker struct {
    river.WorkerDefaults[SortArgs]
}

func (w *SortWorker) Work(ctx context.Context, job *river.Job[SortArgs]) error {
    sort.Strings(job.Args.Strings)
    fmt.Printf("Sorted strings: %+v\n", job.Args.Strings)
    return nil
}
```

## Registering Workers

```go
workers := river.NewWorkers()

// AddWorker panics if invalid (acceptable for startup)
river.AddWorker(workers, &SortWorker{})

// Or use AddWorkerSafely for error handling
if err := river.AddWorkerSafely(workers, &SortWorker{}); err != nil {
    return err
}
```

## Starting a Client

```go
dbPool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
if err != nil {
    return err
}

riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
    Queues: map[string]river.QueueConfig{
        river.QueueDefault: {MaxWorkers: 100},
    },
    Workers: workers,
})
if err != nil {
    return err
}

// Start the client
if err := riverClient.Start(ctx); err != nil {
    return err
}

// Stop on shutdown
defer riverClient.Stop(ctx)
```

## Inserting Jobs

Prefer transactional insertion to avoid bugs:

```go
// With transaction (recommended)
_, err = riverClient.InsertTx(ctx, tx, SortArgs{
    Strings: []string{"whale", "tiger", "bear"},
}, nil)

// Without transaction
_, err = riverClient.Insert(ctx, SortArgs{
    Strings: []string{"whale", "tiger", "bear"},
}, nil)
```

## Multiple Queues

Configure different queues for workload isolation:

```go
riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
    Queues: map[string]river.QueueConfig{
        river.QueueDefault:  {MaxWorkers: 100},
        "high_priority":     {MaxWorkers: 100},
        "metadata":          {MaxWorkers: 5},   // Rate-limited APIs
        "scan":              {MaxWorkers: 2},   // I/O intensive
    },
    Workers: workers,
})
```

Override queue per job kind:

```go
func (args AlwaysHighPriorityArgs) InsertOpts() river.InsertOpts {
    return river.InsertOpts{
        Queue: "high_priority",
    }
}
```

Or per insertion:

```go
_, err = riverClient.Insert(ctx, MyArgs{}, &river.InsertOpts{
    Queue: "high_priority",
})
```

## Job Retries

Default: 25 attempts with exponential backoff. Customize per job kind:

```go
type RetryOnceJobArgs struct{}

func (RetryOnceJobArgs) Kind() string { return "retry_once" }

func (RetryOnceJobArgs) InsertOpts() river.InsertOpts {
    return river.InsertOpts{MaxAttempts: 1}
}
```

Custom retry delay per worker:

```go
type ConstantRetryWorker struct {
    river.WorkerDefaults[MyArgs]
}

func (w *ConstantRetryWorker) NextRetry(job *river.Job[MyArgs]) time.Time {
    return time.Now().Add(10 * time.Second)
}
```

## Unique Jobs

Prevent duplicate jobs using `UniqueOpts`:

```go
type ReconcileAccountArgs struct {
    AccountID int `json:"account_id"`
}

func (ReconcileAccountArgs) Kind() string { return "reconcile_account" }

func (ReconcileAccountArgs) InsertOpts() river.InsertOpts {
    return river.InsertOpts{
        UniqueOpts: river.UniqueOpts{
            ByArgs:   true,              // Unique per args
            ByPeriod: 24 * time.Hour,    // Once per 24h
        },
    }
}
```

Check if insert was skipped:

```go
insertRes, err := riverClient.Insert(ctx, MyArgs{}, nil)
if insertRes.UniqueSkippedAsDuplicate {
    // Job already exists
}
```

## Periodic Jobs

```go
periodicJobs := []*river.PeriodicJob{
    river.NewPeriodicJob(
        river.PeriodicInterval(15*time.Minute),
        func() (river.JobArgs, *river.InsertOpts) {
            return MyPeriodicJobArgs{}, nil
        },
        &river.PeriodicJobOpts{RunOnStart: true},
    ),
}

riverClient, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{
    PeriodicJobs: periodicJobs,
    // ...
})
```

For complex cron schedules, use `robfig/cron`:

```go
import "github.com/robfig/cron/v3"

schedule, _ := cron.ParseStandard("*/15 * * * *") // every 15 minutes

river.NewPeriodicJob(schedule, func() (river.JobArgs, *river.InsertOpts) {
    return MyArgs{}, nil
}, nil)
```

## Revenge-Specific Patterns

### Job Naming Convention

Use `{module}.{action}` format:

- `library.scan` - Scan a library
- `movie.refresh_metadata` - Refresh movie metadata
- `music.import_playlist` - Import music playlist
- `search.reindex` - Reindex search collection

### Recommended Queue Setup

```go
Queues: map[string]river.QueueConfig{
    river.QueueDefault: {MaxWorkers: 10},
    "metadata":         {MaxWorkers: 5},  // External API rate limits
    "scan":             {MaxWorkers: 2},  // I/O intensive
    "search":           {MaxWorkers: 3},  // Typesense indexing
}
```

## Do's and Don'ts

### DO

- ✅ Use `InsertTx` for transactional job insertion
- ✅ Use unique constraints to prevent duplicate jobs
- ✅ Keep job arguments JSON-serializable and small
- ✅ Use separate queues for different workload types
- ✅ Combine periodic jobs with unique jobs for reliability

### DON'T

- ❌ Store large data in job arguments (use references/IDs)
- ❌ Ignore context cancellation in workers
- ❌ Rename/remove queues without handling existing jobs
- ❌ Use River for real-time/low-latency tasks
---

## Related

- [INDEX.instructions.md](INDEX.instructions.md) - Main instruction index with all cross-references
- [ARCHITECTURE_V2.md](../../docs/dev/design/architecture/ARCHITECTURE_V2.md) - System architecture
- [BEST_PRACTICES.md](../../docs/dev/design/operations/BEST_PRACTICES.md) - Best practices

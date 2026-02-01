

---
sources:
  - name: Dragonfly Documentation
    url: ../../../sources/infrastructure/dragonfly.md
    note: Auto-resolved from dragonfly
  - name: Uber fx
    url: ../../../sources/tooling/fx.md
    note: Auto-resolved from fx
  - name: google/uuid
    url: ../../../sources/tooling/uuid.md
    note: Auto-resolved from google-uuid
  - name: pgx PostgreSQL Driver
    url: ../../../sources/database/pgx.md
    note: Auto-resolved from pgx
  - name: pgxpool Connection Pool
    url: ../../../sources/database/pgxpool.md
    note: Auto-resolved from pgxpool
  - name: PostgreSQL Arrays
    url: ../../../sources/database/postgresql-arrays.md
    note: Auto-resolved from postgresql-arrays
  - name: PostgreSQL JSON Functions
    url: ../../../sources/database/postgresql-json.md
    note: Auto-resolved from postgresql-json
  - name: Prometheus Go Client
    url: ../../../sources/observability/prometheus.md
    note: Auto-resolved from prometheus
  - name: Prometheus Metric Types
    url: ../../../sources/observability/prometheus-metrics.md
    note: Auto-resolved from prometheus-metrics
  - name: River Job Queue
    url: ../../../sources/tooling/river.md
    note: Auto-resolved from river
  - name: River Documentation
    url: ../../../sources/tooling/river-guide.md
    note: Auto-resolved from river-docs
  - name: rueidis
    url: ../../../sources/tooling/rueidis.md
    note: Auto-resolved from rueidis
  - name: rueidis GitHub README
    url: ../../../sources/tooling/rueidis-guide.md
    note: Auto-resolved from rueidis-docs
  - name: Typesense API
    url: ../../../sources/infrastructure/typesense.md
    note: Auto-resolved from typesense
  - name: Typesense Go Client
    url: ../../../sources/infrastructure/typesense-go.md
    note: Auto-resolved from typesense-go
design_refs:
  - title: 01_ARCHITECTURE
    path: ../../architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: ../../architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: ../../architecture/03_METADATA_SYSTEM.md
---

## Table of Contents

- [River](#river)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Integration Structure](#integration-structure)
    - [Data Flow](#data-flow)
    - [Provides](#provides)
  - [Implementation](#implementation)
    - [File Structure](#file-structure)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
- [River configuration](#river-configuration)
- [Queue priorities](#queue-priorities)
    - [Config Keys](#config-keys)
  - [API Endpoints](#api-endpoints)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)


# River


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration


> Integration with River

> PostgreSQL-native job queue

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | ✅ | - |
| Instructions | ✅ | - |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---


## Architecture

### Integration Structure

```
internal/integration/river/
├── client.go              # API client
├── types.go               # Response types
├── mapper.go              # Map external → internal types
├── cache.go               # Response caching
└── client_test.go         # Tests
```

### Data Flow

<!-- Data flow diagram -->

### Provides
<!-- Data provided by integration -->


## Implementation

### File Structure

<!-- File structure -->

### Key Interfaces

```go
// River client wrapper
type JobQueue interface {
  Enqueue(ctx context.Context, job river.JobArgs) (*river.JobInsertResult, error)
  EnqueueWithPriority(ctx context.Context, job river.JobArgs, priority int) (*river.JobInsertResult, error)
  EnqueueScheduled(ctx context.Context, job river.JobArgs, scheduledAt time.Time) (*river.JobInsertResult, error)
  Cancel(ctx context.Context, jobID int64) error
  Start(ctx context.Context) error
  Stop(ctx context.Context) error
}

// Example worker
type LibraryScanWorker struct {
  river.WorkerDefaults[LibraryScanArgs]
  libraryService LibraryService
}

type LibraryScanArgs struct {
  LibraryID uuid.UUID `json:"library_id"`
  FullScan  bool      `json:"full_scan"`
}

func (w *LibraryScanWorker) Work(ctx context.Context, job *river.Job[LibraryScanArgs]) error {
  return w.libraryService.ScanLibrary(ctx, job.Args.LibraryID, job.Args.FullScan)
}
```


### Dependencies

**Go Packages**:
- `github.com/riverqueue/river` - Job queue
- `github.com/riverqueue/river/riverdriver/riverpgxv5` - PostgreSQL driver
- `github.com/jackc/pgx/v5/pgxpool` - Connection pool
- `github.com/google/uuid` - UUID support
- `go.uber.org/fx` - Dependency injection






## Configuration
### Environment Variables

```bash
# River configuration
RIVER_WORKERS=10
RIVER_MAX_ATTEMPTS=25
RIVER_POLL_INTERVAL=1s
RIVER_SHUTDOWN_TIMEOUT=30s

# Queue priorities
RIVER_QUEUE_DEFAULT_PRIORITY=1
RIVER_QUEUE_HIGH_PRIORITY=10
```


### Config Keys

```yaml
jobs:
  river:
    workers: 10
    max_attempts: 25
    poll_interval: 1s
    shutdown_timeout: 30s

    queues:
      default:
        max_workers: 10
      high_priority:
        max_workers: 20
      low_priority:
        max_workers: 5
```



## API Endpoints
**List Jobs**:
```
GET /api/v1/admin/jobs?state=running&limit=50
```

**Response**:
```json
{
  "jobs": [
    {
      "id": 12345,
      "state": "running",
      "queue": "default",
      "kind": "library_scan",
      "args": {
        "library_id": "uuid-123",
        "full_scan": false
      },
      "attempt": 1,
      "max_attempts": 25,
      "created_at": "2026-02-01T10:00:00Z",
      "scheduled_at": "2026-02-01T10:00:00Z"
    }
  ],
  "total": 1
}
```

**Cancel Job**:
```
DELETE /api/v1/admin/jobs/:id
```



## Testing Strategy

### Unit Tests

<!-- Unit test strategy -->

### Integration Tests

<!-- Integration test strategy -->

### Test Coverage

Target: **80% minimum**







## Related Documentation
### Design Documents
- [01_ARCHITECTURE](../../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Dragonfly Documentation](../../../sources/infrastructure/dragonfly.md) - Auto-resolved from dragonfly
- [Uber fx](../../../sources/tooling/fx.md) - Auto-resolved from fx
- [google/uuid](../../../sources/tooling/uuid.md) - Auto-resolved from google-uuid
- [pgx PostgreSQL Driver](../../../sources/database/pgx.md) - Auto-resolved from pgx
- [pgxpool Connection Pool](../../../sources/database/pgxpool.md) - Auto-resolved from pgxpool
- [PostgreSQL Arrays](../../../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../../../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [Prometheus Go Client](../../../sources/observability/prometheus.md) - Auto-resolved from prometheus
- [Prometheus Metric Types](../../../sources/observability/prometheus-metrics.md) - Auto-resolved from prometheus-metrics
- [River Job Queue](../../../sources/tooling/river.md) - Auto-resolved from river
- [River Documentation](../../../sources/tooling/river-guide.md) - Auto-resolved from river-docs
- [rueidis](../../../sources/tooling/rueidis.md) - Auto-resolved from rueidis
- [rueidis GitHub README](../../../sources/tooling/rueidis-guide.md) - Auto-resolved from rueidis-docs
- [Typesense API](../../../sources/infrastructure/typesense.md) - Auto-resolved from typesense
- [Typesense Go Client](../../../sources/infrastructure/typesense-go.md) - Auto-resolved from typesense-go


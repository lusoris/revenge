## Table of Contents

- [Analytics Service](#analytics-service)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Service Structure](#service-structure)
    - [Dependencies](#dependencies)
    - [Provides](#provides)
    - [Component Diagram](#component-diagram)
  - [Implementation](#implementation)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
    - [Config Keys](#config-keys)
  - [API Endpoints](#api-endpoints)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Analytics Service


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: service


> > Usage analytics, playback statistics, and library insights

**Package**: `internal/service/analytics`
**fx Module**: `analytics.Module`

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | ✅ | - |
| Instructions | ✅ | Generated from design |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---


## Architecture

```mermaid
flowchart TD
    node1([Client<br/>(Web/App)])
    node2[[API Handler<br/>(ogen)]]
    node3[[Service<br/>(Logic)]]
    node4["Repository<br/>(sqlc)"]
    node5["River<br/>Queue"]
    node6[(Cache<br/>(otter))]
    node7[(PostgreSQL<br/>(pgx))]
    node1 --> node2
    node2 --> node3
    node4 --> node5
    node5 --> node6
    node3 --> node4
    node6 --> node7
```

### Service Structure

```
internal/service/analytics/
├── module.go              # fx module definition
├── service.go             # Service implementation
├── repository.go          # Data access (if needed)
├── handler.go             # HTTP handlers (if exposed)
├── middleware.go          # Middleware (if needed)
├── types.go               # Domain types
└── service_test.go        # Tests
```

### Dependencies
**Go Packages**:
- `github.com/google/uuid`
- `github.com/jackc/pgx/v5`
- `github.com/riverqueue/river` - Aggregation jobs
- `github.com/maypok86/otter` - Stats cache
- `go.uber.org/fx`


### Provides
<!-- Service provides -->

### Component Diagram

<!-- Component diagram -->
## Implementation

### Key Interfaces

```go
type AnalyticsService interface {
  // Library stats
  GetLibraryStats(ctx context.Context, libraryID uuid.UUID, dateRange DateRange) (*LibraryStats, error)
  GetServerStats(ctx context.Context, dateRange DateRange) (*ServerStats, error)

  // Popular content
  GetMostWatched(ctx context.Context, contentType string, period TimePeriod, limit int) ([]PopularItem, error)
  GetTopGenres(ctx context.Context, period TimePeriod) ([]GenreStats, error)

  // User insights
  GetUserActivity(ctx context.Context, userID uuid.UUID, dateRange DateRange) (*UserActivity, error)
  GetActiveUsers(ctx context.Context, period TimePeriod) (int, error)

  // Aggregation
  AggregateDaily(ctx context.Context, date time.Time) error
}

type ServerStats struct {
  TotalLibraries    int     `json:"total_libraries"`
  TotalItems        int     `json:"total_items"`
  TotalSizeGB       float64 `json:"total_size_gb"`
  TotalUsers        int     `json:"total_users"`
  ActiveUsers24h    int     `json:"active_users_24h"`
  TotalPlays        int     `json:"total_plays"`
  TotalWatchHours   float64 `json:"total_watch_hours"`
}
```


### Dependencies
**Go Packages**:
- `github.com/google/uuid`
- `github.com/jackc/pgx/v5`
- `github.com/riverqueue/river` - Aggregation jobs
- `github.com/maypok86/otter` - Stats cache
- `go.uber.org/fx`

## Configuration

### Environment Variables

```bash
ANALYTICS_AGGREGATION_INTERVAL=1h
ANALYTICS_RETENTION_DAYS=365
```


### Config Keys
```yaml
analytics:
  aggregation_interval: 1h
  retention_days: 365
  cache_ttl: 5m
```

## API Endpoints
```
GET    /api/v1/analytics/server              # Server overview
GET    /api/v1/analytics/libraries/:id       # Library stats
GET    /api/v1/analytics/popular/:type       # Most watched content
GET    /api/v1/analytics/users/:id           # User activity
```

**Example Server Stats Response**:
```json
{
  "total_libraries": 5,
  "total_items": 1543,
  "total_size_gb": 8542.3,
  "total_users": 12,
  "active_users_24h": 8,
  "total_plays": 23456,
  "total_watch_hours": 4532.5
}
```

## Related Documentation
### Design Documents
- [services](INDEX.md)
- [01_ARCHITECTURE](../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Uber fx](../../sources/tooling/fx.md) - Auto-resolved from fx
- [pgx PostgreSQL Driver](../../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [River Job Queue](../../sources/tooling/river.md) - Auto-resolved from river


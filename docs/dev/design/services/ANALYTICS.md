# Analytics Service

> Usage analytics, playback statistics, and library insights


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Developer Resources](#developer-resources)
- [Overview](#overview)
- [Goals](#goals)
- [Non-Goals](#non-goals)
- [Technical Design](#technical-design)
  - [Analytics Types](#analytics-types)
  - [Data Model](#data-model)
  - [Service Interface](#service-interface)
- [Database Schema](#database-schema)
- [River Jobs](#river-jobs)
- [Configuration](#configuration)
- [Checklist](#checklist)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documents](#related-documents)

<!-- TOC-END -->

**Module**: `internal/service/analytics`
**Dependencies**: [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#go-dependencies-observability)

## Status

| Dimension | Status |
|-----------|--------|
| Design | ✅ |
| Sources | ✅ |
| Instructions | 🔴 |
| Code | 🔴 |
| Linting | 🔴 |
| Unit Testing | 🔴 |
| Integration Testing | 🔴 |
## Developer Resources

> Package versions: [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#go-dependencies-core)

| Package | Purpose |
|---------|---------|
| River | Background job processing |
| pgx | PostgreSQL driver |
| otter | In-memory caching |
| fx | Dependency injection |

---

## Overview

The Analytics service collects and aggregates usage data for insights:
- Playback statistics (most watched, watch time)
- Library statistics (size, growth, health)
- User activity patterns
- System performance metrics

All analytics are privacy-respecting and stored locally.

## Goals

- Provide library owners with usage insights
- Enable "most popular" and "trending" features
- Track library health (missing metadata, orphaned files)
- Support admin dashboards

## Non-Goals

- Send data to external services
- Track individual user behavior for profiling
- Real-time analytics (batch processing is fine)

---

## Technical Design

### Analytics Types

| Type | Granularity | Retention |
|------|-------------|-----------|
| Playback events | Per-play | 90 days raw, aggregated forever |
| Library stats | Daily snapshot | Forever |
| User activity | Daily aggregate | 30 days |
| System metrics | Hourly | 7 days |

### Data Model

```go
type PlaybackEvent struct {
    ID          uuid.UUID
    UserID      uuid.UUID
    MediaID     uuid.UUID
    MediaType   string
    StartedAt   time.Time
    EndedAt     *time.Time
    DurationSec int
    Completed   bool
    ClientType  string
}

type LibraryStats struct {
    LibraryID     uuid.UUID
    Date          time.Time
    TotalItems    int
    TotalSizeGB   float64
    MissingMeta   int
    RecentlyAdded int
}

type UserActivityAggregate struct {
    UserID      uuid.UUID
    Date        time.Time
    PlayCount   int
    WatchTimeMin int
    UniqueItems int
}
```

### Service Interface

```go
type AnalyticsService interface {
    // Recording
    RecordPlayback(ctx context.Context, event PlaybackEvent) error
    RecordLibraryScan(ctx context.Context, libraryID uuid.UUID, stats ScanStats) error

    // Queries
    GetMostWatched(ctx context.Context, libraryID uuid.UUID, period time.Duration, limit int) ([]MediaStats, error)
    GetRecentlyWatched(ctx context.Context, userID uuid.UUID, limit int) ([]MediaStats, error)
    GetLibraryStats(ctx context.Context, libraryID uuid.UUID) (*LibraryStats, error)
    GetUserStats(ctx context.Context, userID uuid.UUID, period time.Duration) (*UserStats, error)

    // Admin
    GetSystemStats(ctx context.Context) (*SystemStats, error)
    GetActiveUsers(ctx context.Context, since time.Duration) ([]UserActivity, error)
}
```

---

## Database Schema

```sql
CREATE TABLE playback_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    media_id UUID NOT NULL,
    media_type VARCHAR(20) NOT NULL,
    started_at TIMESTAMPTZ NOT NULL,
    ended_at TIMESTAMPTZ,
    duration_sec INT,
    completed BOOLEAN DEFAULT FALSE,
    client_type VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (started_at);

CREATE TABLE library_stats_daily (
    library_id UUID NOT NULL,
    date DATE NOT NULL,
    total_items INT,
    total_size_bytes BIGINT,
    missing_metadata INT,
    recently_added INT,
    PRIMARY KEY (library_id, date)
);

CREATE TABLE user_activity_daily (
    user_id UUID NOT NULL,
    date DATE NOT NULL,
    play_count INT DEFAULT 0,
    watch_time_min INT DEFAULT 0,
    unique_items INT DEFAULT 0,
    PRIMARY KEY (user_id, date)
);
```

---

## River Jobs

```go
type AggregateAnalyticsArgs struct {
    Date time.Time `json:"date"`
}

func (AggregateAnalyticsArgs) Kind() string { return "analytics.aggregate_daily" }

type CleanupAnalyticsArgs struct {
    RetentionDays int `json:"retention_days"`
}

func (CleanupAnalyticsArgs) Kind() string { return "analytics.cleanup" }
```

---

## Configuration

```yaml
analytics:
  enabled: true
  retention:
    raw_events: 90d
    daily_aggregates: 365d
  aggregation_schedule: "0 3 * * *"  # 3 AM daily
```

---

## Checklist

- [ ] Database migrations created
- [ ] Playback event recording
- [ ] Daily aggregation job
- [ ] Query methods implemented
- [ ] Admin dashboard endpoints
- [ ] Cleanup job for retention
- [ ] Tests written

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [PostgreSQL Arrays](https://www.postgresql.org/docs/current/arrays.html) | [Local](../../sources/database/postgresql-arrays.md) |
| [PostgreSQL JSON Functions](https://www.postgresql.org/docs/current/functions-json.html) | [Local](../../sources/database/postgresql-json.md) |
| [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) | [Local](../../sources/tooling/river.md) |
| [Uber fx](https://pkg.go.dev/go.uber.org/fx) | [Local](../../sources/tooling/fx.md) |
| [pgx PostgreSQL Driver](https://pkg.go.dev/github.com/jackc/pgx/v5) | [Local](../../sources/database/pgx.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Services](INDEX.md)

### In This Section

- [Activity Service](ACTIVITY.md)
- [API Keys Service](APIKEYS.md)
- [Auth Service](AUTH.md)
- [Fingerprint Service](FINGERPRINT.md)
- [Grants Service](GRANTS.md)
- [Library Service](LIBRARY.md)
- [Metadata Service](METADATA.md)
- [Notification Service](NOTIFICATION.md)

### Related Topics

- [Revenge - Architecture v2](../architecture/01_ARCHITECTURE.md) _Architecture_
- [Revenge - Design Principles](../architecture/02_DESIGN_PRINCIPLES.md) _Architecture_
- [Revenge - Metadata System](../architecture/03_METADATA_SYSTEM.md) _Architecture_
- [Revenge - Player Architecture](../architecture/04_PLAYER_ARCHITECTURE.md) _Architecture_
- [Plugin Architecture Decision](../architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) _Architecture_

### Indexes

- [Design Index](../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

## Related Documents

- [Activity Service](ACTIVITY.md) - Event logging
- [Library Service](LIBRARY.md) - Library statistics
- [Session Service](SESSION.md) - User session tracking
- [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) - Service inventory

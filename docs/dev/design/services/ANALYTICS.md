# Analytics Service

> Usage analytics, playback statistics, and library insights

**Status**: 🔴 PLANNED
**Priority**: 🟢 LOW
**Module**: `internal/service/analytics`
**Dependencies**: [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#go-dependencies-observability)

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

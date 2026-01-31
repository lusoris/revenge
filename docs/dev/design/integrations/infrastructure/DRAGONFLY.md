# Dragonfly Integration

<!-- SOURCES: dragonfly, pgx, postgresql-arrays, postgresql-json, prometheus, prometheus-metrics, river, rueidis, rueidis-docs, typesense, typesense-go -->

<!-- DESIGN: integrations/infrastructure, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> High-performance Redis-compatible cache


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
- [Developer Resources](#developer-resources)
- [Connection Details](#connection-details)
- [Configuration](#configuration)
- [Cache Usage Patterns](#cache-usage-patterns)
  - [Session Storage](#session-storage)
  - [API Response Caching](#api-response-caching)
  - [Rate Limiting](#rate-limiting)
  - [Pub/Sub for Real-time](#pubsub-for-real-time)
  - [Distributed Locking](#distributed-locking)
- [Key Naming Convention](#key-naming-convention)
- [Implementation Checklist](#implementation-checklist)
- [Docker Compose](#docker-compose)
  - [Dragonfly vs Redis Configuration](#dragonfly-vs-redis-configuration)
- [Health Checks](#health-checks)
- [Monitoring](#monitoring)
  - [Key Metrics](#key-metrics)
  - [Dragonfly INFO](#dragonfly-info)
- [Error Handling](#error-handling)
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
**Type**: In-memory cache/session store

---

## Overview

Dragonfly is a modern, Redis-compatible in-memory data store that serves as Revenge's caching layer. It provides:
- Session storage
- API response caching
- Rate limiting state
- Real-time features (pub/sub)
- Distributed locking

**Why Dragonfly over Redis?**:
- 25x faster than Redis on same hardware
- Lower memory footprint
- Native multi-threading
- Full Redis API compatibility
- Drop-in replacement

---

## Developer Resources

- 📚 **Docs**: https://www.dragonflydb.io/docs
- 🔗 **GitHub**: https://github.com/dragonflydb/dragonfly
- 🔗 **Go Client**: `github.com/redis/rueidis` (14x faster, auto-pipelining, server-assisted client-side caching)
- 🔗 **Commands**: https://www.dragonflydb.io/docs/command-reference

---

## Connection Details

**Default Settings**:
| Setting | Value |
|---------|-------|
| Host | `localhost` |
| Port | `6379` |
| Database | `0` |
| Password | (none by default) |

**Connection String**:
```
redis://localhost:6379/0
```

---

## Configuration

```yaml
# configs/config.yaml
cache:
  enabled: true
  driver: "dragonfly"  # or "redis"

  connection:
    host: "${REVENGE_CACHE_HOST:localhost}"
    port: ${REVENGE_CACHE_PORT:6379}
    password: "${REVENGE_CACHE_PASSWORD:}"
    database: 0

  pool:
    size: 100
    min_idle: 10
    max_idle_time: "5m"

  # TLS (optional)
  tls:
    enabled: false
    cert_file: ""
    key_file: ""
    ca_file: ""

  # Timeouts
  timeouts:
    dial: "5s"
    read: "3s"
    write: "3s"
```

---

## Cache Usage Patterns

### Session Storage

```go
// Store session
func (s *SessionStore) Set(ctx context.Context, sessionID string, data *Session) error {
    jsonData, _ := json.Marshal(data)
    return s.cache.Set(ctx, "session:"+sessionID, jsonData, 24*time.Hour).Err()
}

// Get session
func (s *SessionStore) Get(ctx context.Context, sessionID string) (*Session, error) {
    data, err := s.cache.Get(ctx, "session:"+sessionID).Bytes()
    if errors.Is(err, redis.Nil) {
        return nil, ErrSessionNotFound
    }
    var session Session
    json.Unmarshal(data, &session)
    return &session, nil
}
```

### API Response Caching

```go
// Cache movie metadata
key := fmt.Sprintf("movie:%s:metadata", movieID)
cache.Set(ctx, key, metadata, 1*time.Hour)

// Cache TMDb response
key := fmt.Sprintf("tmdb:movie:%d", tmdbID)
cache.Set(ctx, key, response, 24*time.Hour)
```

### Rate Limiting

```go
// Sliding window rate limiter
func (r *RateLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) bool {
    now := time.Now().UnixMilli()
    windowStart := now - window.Milliseconds()

    pipe := r.cache.Pipeline()
    pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprint(windowStart))
    pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: now})
    pipe.ZCard(ctx, key)
    pipe.Expire(ctx, key, window)

    results, _ := pipe.Exec(ctx)
    count := results[2].(*redis.IntCmd).Val()

    return count <= int64(limit)
}
```

### Pub/Sub for Real-time

```go
// Publish playback progress
func (p *PlaybackPublisher) PublishProgress(ctx context.Context, sessionID string, progress Progress) {
    data, _ := json.Marshal(progress)
    p.cache.Publish(ctx, "playback:"+sessionID, data)
}

// Subscribe to progress updates
func (s *PlaybackSubscriber) Subscribe(ctx context.Context, sessionID string) <-chan Progress {
    sub := s.cache.Subscribe(ctx, "playback:"+sessionID)
    ch := make(chan Progress)

    go func() {
        defer close(ch)
        for msg := range sub.Channel() {
            var progress Progress
            json.Unmarshal([]byte(msg.Payload), &progress)
            ch <- progress
        }
    }()

    return ch
}
```

### Distributed Locking

```go
// Acquire lock for library scan
func (l *Locker) AcquireScanLock(ctx context.Context, libraryID string) (func(), error) {
    key := "lock:scan:" + libraryID

    // Try to set lock with NX (only if not exists)
    ok, err := l.cache.SetNX(ctx, key, "locked", 5*time.Minute).Result()
    if err != nil {
        return nil, err
    }
    if !ok {
        return nil, ErrLockNotAcquired
    }

    // Return release function
    return func() {
        l.cache.Del(ctx, key)
    }, nil
}
```

---

## Key Naming Convention

```
# Sessions
session:{session_id}

# User data
user:{user_id}:profile
user:{user_id}:preferences

# Media metadata cache
movie:{movie_id}:metadata
movie:{movie_id}:streams
tvshow:{show_id}:metadata

# External API cache
tmdb:movie:{tmdb_id}
tmdb:search:{query_hash}
musicbrainz:artist:{mbid}

# Rate limiting
ratelimit:{ip}:{endpoint}
ratelimit:user:{user_id}:{endpoint}

# Locks
lock:scan:{library_id}
lock:metadata:{item_id}

# Real-time channels
playback:{session_id}
notifications:{user_id}
```

---

## Implementation Checklist

- [ ] **Cache Client** (`internal/infra/cache/client.go`)
  - [ ] Connection management
  - [ ] Connection pooling
  - [ ] Health checks
  - [ ] Graceful shutdown

- [ ] **Session Store** (`internal/infra/cache/sessions.go`)
  - [ ] Set/Get/Delete sessions
  - [ ] Session refresh
  - [ ] Bulk invalidation

- [ ] **Response Cache** (`internal/infra/cache/response.go`)
  - [ ] HTTP response caching
  - [ ] Cache-Control header support
  - [ ] Conditional caching (ETag)

- [ ] **Rate Limiter** (`internal/infra/cache/ratelimit.go`)
  - [ ] Fixed window
  - [ ] Sliding window
  - [ ] Token bucket
  - [ ] Per-IP and per-user limits

- [ ] **Pub/Sub** (`internal/infra/cache/pubsub.go`)
  - [ ] Publish/Subscribe patterns
  - [ ] Channel management
  - [ ] Reconnection handling

---

## Docker Compose

```yaml
services:
  dragonfly:
    image: docker.dragonflydb.io/dragonflydb/dragonfly
    container_name: revenge-dragonfly
    ulimits:
      memlock: -1
    ports:
      - "6379:6379"
    volumes:
      - dragonfly_data:/data
    command: >
      --maxmemory 512mb
      --proactor_threads 2
      --cache_mode true
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  dragonfly_data:
```

### Dragonfly vs Redis Configuration

```yaml
# Dragonfly-specific options
command: >
  --maxmemory 512mb          # Memory limit
  --proactor_threads 2       # Thread count
  --cache_mode true          # Eviction when full
  --snapshot_cron "0 */6 * * *"  # Periodic snapshots
```

---

## Health Checks

```go
func (c *Cache) HealthCheck(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()

    pong, err := c.client.Ping(ctx).Result()
    if err != nil {
        return fmt.Errorf("cache health check failed: %w", err)
    }
    if pong != "PONG" {
        return fmt.Errorf("unexpected ping response: %s", pong)
    }
    return nil
}
```

---

## Monitoring

### Key Metrics

```go
// Prometheus metrics
var (
    cacheHits = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "revenge_cache_hits_total",
            Help: "Total cache hits",
        },
        []string{"key_prefix"},
    )

    cacheMisses = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "revenge_cache_misses_total",
            Help: "Total cache misses",
        },
        []string{"key_prefix"},
    )

    cacheLatency = prometheus.NewHistogram(
        prometheus.HistogramOpts{
            Name:    "revenge_cache_operation_duration_seconds",
            Help:    "Cache operation latency",
            Buckets: []float64{.001, .005, .01, .025, .05, .1},
        },
    )
)
```

### Dragonfly INFO

```bash
# Get server info
redis-cli INFO

# Memory stats
redis-cli INFO memory

# Client connections
redis-cli CLIENT LIST
```

---

## Error Handling

| Error | Cause | Solution |
|-------|-------|----------|
| Connection refused | Dragonfly not running | Check Docker container |
| OOM | Memory limit reached | Increase maxmemory or enable eviction |
| Timeout | Network or load issues | Check network, increase timeout |
| NOAUTH | Password required | Configure password |

---


## Related Documentation

- [PostgreSQL](POSTGRESQL.md)
- [Typesense](TYPESENSE.md)
- [River](RIVER.md)

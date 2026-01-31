# Advanced Patterns & Best Practices

> Comprehensive guide for professional-grade implementations in Revenge.

## Table of Contents

1. [Disk-Based Transcode Cache](#disk-based-transcode-cache)
2. [Resilience Patterns](#resilience-patterns)
3. [Self-Healing & Supervision](#self-healing--supervision)
4. [Graceful Shutdown](#graceful-shutdown)
5. [Hot Reload](#hot-reload)
6. [Observability](#observability)
7. [Memory Management](#memory-management)
8. [Database Patterns](#database-patterns)
9. [API Best Practices](#api-best-practices)

---

## Disk-Based Transcode Cache

### Problem
Transcoding is expensive. If the same content is requested with the same profile, we shouldn't transcode again.

### Solution
Persist transcoded segments to disk with quota management:

```
/var/cache/revenge/transcodes/
  index.json                    # Cache index (survives restart)
  ab/                           # First 2 chars of cache key
    ab1234.../
      master.m3u8               # HLS manifest
      segment_0.ts              # Video segments
      segment_1.ts
      ...
```

### Cache Key
Deterministic key based on:
- Media ID
- Transcode profile ID
- Source file hash (for invalidation when source changes)

```go
cacheKey := sha256(mediaID + profileID + sourceHash)[:32]
```

### Quotas

| Quota Type | Default | Purpose |
|------------|---------|---------|
| Global | 50 GB | Total cache size |
| Per-user | 10 GB | Prevent one user monopolizing |
| Per-media | 5 GB | Limit per title (multiple profiles) |
| Min free space | 10 GB | Reserve disk space |
| Max age | 72h | Auto-expire old transcodes |

### Cache Hit Flow

```
1. Client requests stream
2. Generate cache key from (mediaID, profile, sourceHash)
3. Check disk cache index
4. If hit → Serve from disk (no Blackbeard needed)
5. If miss → Transcode via Blackbeard, write to cache
```

### Configuration

```yaml
playback:
  disk_cache:
    enabled: true
    base_path: /var/cache/revenge/transcodes
    max_size_bytes: 53687091200  # 50 GB
    max_age_hours: 72
    min_free_space_bytes: 10737418240  # 10 GB
    per_user_quota_bytes: 10737418240  # 10 GB
    per_media_quota_bytes: 5368709120  # 5 GB
    eviction_check_interval: 5m
```

---

## Resilience Patterns

> **Implementation**: Uses [failsafe-go](https://github.com/failsafe-go/failsafe-go) for circuit breakers, retries, bulkheads, and rate limiting.

### Circuit Breaker

Prevent cascade failures when external services fail:

```go
cb := resilience.NewCircuitBreaker(resilience.CircuitBreakerConfig{
    Name:        "blackbeard",
    MaxFailures: 5,
    Timeout:     30 * time.Second,
})

err := cb.Execute(func() error {
    return blackbeard.StartTranscode(ctx, request)
})

if errors.Is(err, resilience.ErrCircuitOpen) {
    // Serve direct stream or cached version
}
```

**States:**
- **Closed** → Normal operation
- **Open** → Rejecting requests (service is down)
- **Half-Open** → Testing if recovered

### Bulkhead

Isolate failures by limiting concurrent operations:

```go
bulkhead := resilience.NewBulkhead(resilience.BulkheadConfig{
    Name:          "transcoding",
    MaxConcurrent: 10,
    MaxWait:       5 * time.Second,
    QueueSize:     100,
})

err := bulkhead.Execute(func() error {
    return startTranscode(ctx, req)
})

if errors.Is(err, resilience.ErrBulkheadFull) {
    // Return 503 Service Unavailable
}
```

### Rate Limiting

Protect APIs from abuse:

```go
// Per-user rate limiter
limiter := resilience.NewPerKeyLimiter(
    resilience.RateLimiterConfig{
        Rate:  100,  // requests per second
        Burst: 200,
    },
    time.Minute, // cleanup interval
)

func rateLimitMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        userID := getUserID(r)
        if !limiter.Allow(userID) {
            http.Error(w, "Rate limited", http.StatusTooManyRequests)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

### Retry with Backoff

Retry transient failures:

```go
retry := resilience.Retry{
    MaxAttempts: 3,
    InitialWait: 100 * time.Millisecond,
    MaxWait:     10 * time.Second,
    Multiplier:  2.0,
    Jitter:      0.1,
}

err := retry.DoWithContext(ctx, func(ctx context.Context) error {
    return fetchMetadata(ctx, mediaID)
})
```

---

## Self-Healing & Supervision

### Service Supervisor

Automatically restart failed services:

```go
supervisor := supervisor.NewSupervisor(
    supervisor.SupervisorConfig{
        Name:             "revenge",
        Strategy:         supervisor.StrategyOneForOne,
        MaxRestarts:      5,
        MaxRestartWindow: time.Minute,
        RestartDelay:     100 * time.Millisecond,
        MaxRestartDelay:  30 * time.Second,
    },
    logger,
)

// Add services
supervisor.Add(metadataFetcher)
supervisor.Add(libraryScanner)
supervisor.Add(searchIndexer)

// Start supervision
supervisor.Start()
defer supervisor.Stop()

// Manual restart if needed
supervisor.RestartService("metadata-fetcher")
```

### Supervision Strategies

| Strategy | Behavior |
|----------|----------|
| OneForOne | Only restart failed service |
| OneForAll | Restart all if one fails |
| RestForOne | Restart failed + services started after it |

### Health Integration

```go
func (s *Supervisor) HealthCheck() error {
    for _, svc := range s.services {
        if svc.state == StateFailed {
            return fmt.Errorf("service %s failed", svc.Name())
        }
    }
    return nil
}
```

---

## Graceful Shutdown

### Shutdown Hooks

Register cleanup in priority order:

```go
shutdowner := graceful.NewShutdowner(
    graceful.DefaultShutdownConfig(),
    logger,
)

// Priority 0: Stop accepting requests
shutdowner.RegisterFunc("http-server", 0, func(ctx context.Context) error {
    return httpServer.Shutdown(ctx)
})

// Priority 10: Drain active streams
shutdowner.RegisterFunc("playback-sessions", 10, func(ctx context.Context) error {
    return playbackService.DrainSessions(ctx)
})

// Priority 20: Flush caches
shutdowner.RegisterFunc("cache-flush", 20, func(ctx context.Context) error {
    transcodeCache.SaveIndex()
    return diskCache.SaveIndex()
})

// Priority 30: Close database
shutdowner.RegisterFunc("database", 30, func(ctx context.Context) error {
    return db.Close()
})

// Start and wait
done := shutdowner.Start()
<-done
```

### Connection Draining

```go
type DrainableServer struct {
    server *http.Server
    active sync.WaitGroup
}

func (s *DrainableServer) Shutdown(ctx context.Context) error {
    // Stop accepting new connections
    if err := s.server.Shutdown(ctx); err != nil {
        return err
    }

    // Wait for in-flight requests
    done := make(chan struct{})
    go func() {
        s.active.Wait()
        close(done)
    }()

    select {
    case <-done:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

---

## Hot Reload

### Configuration Hot Reload

```go
watcher := hotreload.NewConfigWatcher(
    hotreload.WatcherConfig{
        Files:        []string{"config.yaml", "config.local.yaml"},
        PollInterval: 5 * time.Second,
        Debounce:     time.Second,
        OnReload: func(err error) {
            if err == nil {
                logger.Info("config reloaded")
            }
        },
    },
    configLoader,
    logger,
)

watcher.Start(ctx)
```

### Feature Flags

Runtime feature toggles without deployment:

```go
flags := hotreload.NewFeatureFlags()

flags.Set(hotreload.FeatureFlagConfig{
    Name:       "new-player",
    Enabled:    true,
    Percentage: 10,  // 10% rollout
})

if flags.IsEnabledForUser("new-player", userID) {
    // Use new player
}
```

### Atomic Configuration Swap

```go
type Config struct { /* ... */ }

var currentConfig = hotreload.NewAtomicValue(loadConfig())

// Read config (lock-free)
cfg := currentConfig.Load()

// Update config
currentConfig.Store(newConfig)
```

---

## Observability

### Structured Logging

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
    AddSource: true,
}))

// Request-scoped logger
requestLogger := logger.With(
    slog.String("request_id", requestID),
    slog.String("user_id", userID),
)

requestLogger.Info("playback started",
    slog.String("media_id", mediaID),
    slog.String("profile", profile),
)
```

### Metrics

```go
registry := metrics.NewRegistry()

// Counters
registry.Counter("playback_started").Inc()
registry.Counter("transcode_cache_hits").Inc()

// Gauges
registry.Gauge("active_streams").Inc()
defer registry.Gauge("active_streams").Dec()

// Timers
done := registry.Timer("transcode_latency").Time()
defer done()

// HTTP metrics middleware
httpMetrics := metrics.NewHTTPMetrics()
mux.Use(httpMetrics.Middleware)
```

### Health Endpoint

```go
mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
    status := healthChecker.Check(r.Context())

    if status.Status == health.StatusUnhealthy {
        w.WriteHeader(http.StatusServiceUnavailable)
    }

    json.NewEncoder(w).Encode(status)
})
```

---

## Memory Management

### Overview

Media server UI requires caching substantial metadata for responsive browsing:

| Data Type | Typical Size | Cache Location | TTL |
|-----------|--------------|----------------|-----|
| Library index | 50-200 MB | Memory + Dragonfly | 5m |
| Movie metadata | 5-10 KB/item | Dragonfly | 1h |
| Cover images (URLs) | 500 B/item | Memory | 24h |
| User sessions | 1-2 KB/session | Memory | 24h |
| Search results | 1-5 KB/query | Dragonfly | 30s |
| Transcode segments | 1-5 MB/segment | Memory → Disk | Dynamic |

**Estimated Total Memory:** 1-2 GB for large libraries (50k+ items)

### Memory Budgeting

```go
type MemoryBudget struct {
    TotalLimit     int64 // Total memory budget
    MetadataLimit  int64 // Metadata cache (40%)
    TranscodeLimit int64 // Transcode segments (25%)
    SearchLimit    int64 // Search cache (15%)
    SessionLimit   int64 // Sessions (10%)
    BufferLimit    int64 // Buffers/pools (10%)
}

func NewMemoryBudget(totalMB int64) *MemoryBudget {
    total := totalMB * 1024 * 1024
    return &MemoryBudget{
        TotalLimit:     total,
        MetadataLimit:  int64(float64(total) * 0.40),
        TranscodeLimit: int64(float64(total) * 0.25),
        SearchLimit:    int64(float64(total) * 0.15),
        SessionLimit:   int64(float64(total) * 0.10),
        BufferLimit:    int64(float64(total) * 0.10),
    }
}
```

### Tiered Caching Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     Memory Cache (L1)                           │
│  otter (W-TinyLFU) - Hot data, sub-ms access                   │
│  Size: ~500MB, TTL: 5-30 minutes                               │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Dragonfly Cache (L2)                          │
│  Redis-compatible - Warm data, <10ms access                    │
│  Size: 2-8GB, TTL: 1-24 hours                                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                   PostgreSQL (L3)                               │
│  Persistent storage - All data, 10-50ms access                 │
└─────────────────────────────────────────────────────────────────┘
```

### UI Metadata Cache

```go
type MetadataCache struct {
    mu          sync.RWMutex
    items       map[uuid.UUID]*CachedMetadata
    index       []uuid.UUID          // LRU order
    currentSize int64
    maxSize     int64
    evictionCh  chan struct{}
}

type CachedMetadata struct {
    ID          uuid.UUID
    Title       string
    Year        int
    PosterURL   string
    Size        int64
    AccessedAt  time.Time
    ExpiresAt   time.Time
}

func (c *MetadataCache) Get(id uuid.UUID) (*CachedMetadata, bool) {
    c.mu.RLock()
    item, ok := c.items[id]
    c.mu.RUnlock()

    if !ok || time.Now().After(item.ExpiresAt) {
        return nil, false
    }

    // Update LRU (async to avoid write lock)
    go c.touch(id)
    return item, true
}

func (c *MetadataCache) Set(item *CachedMetadata) {
    c.mu.Lock()
    defer c.mu.Unlock()

    // Check if we need to evict
    for c.currentSize+item.Size > c.maxSize {
        c.evictOldest()
    }

    c.items[item.ID] = item
    c.index = append(c.index, item.ID)
    c.currentSize += item.Size
}

func (c *MetadataCache) evictOldest() {
    if len(c.index) == 0 {
        return
    }

    // Remove oldest
    oldest := c.index[0]
    c.index = c.index[1:]

    if item, ok := c.items[oldest]; ok {
        c.currentSize -= item.Size
        delete(c.items, oldest)
    }
}
```

### Memory Pressure Monitoring

```go
type MemoryMonitor struct {
    budget    *MemoryBudget
    caches    []MemoryAwareCache
    threshold float64  // Trigger cleanup at 85%
    interval  time.Duration
}

func (m *MemoryMonitor) Start(ctx context.Context) {
    ticker := time.NewTicker(m.interval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            m.checkPressure()
        }
    }
}

func (m *MemoryMonitor) checkPressure() {
    var stats runtime.MemStats
    runtime.ReadMemStats(&stats)

    used := int64(stats.Alloc)
    pressure := float64(used) / float64(m.budget.TotalLimit)

    if pressure > m.threshold {
        // Trigger cache evictions
        for _, cache := range m.caches {
            cache.Evict(0.20) // Evict 20% of each cache
        }

        // Force GC after eviction
        runtime.GC()
    }
}
```

### Memory-Aware Transcode Cache

```go
cache := playback.NewTranscodeCache(playback.TranscodeCacheConfig{
    MaxMemoryBytes:     0,     // Auto-detect (25% RAM)
    MaxSegmentsPerSession: 50, // ~5 minutes
    MinRetentionTime:  30 * time.Second,
    HighPressureThreshold: 0.80,  // Start evicting
    CriticalPressureThreshold: 0.95,  // Aggressive eviction
})
```

### Pool Patterns

Reuse allocations:

```go
var bufferPool = sync.Pool{
    New: func() any {
        buf := make([]byte, 64*1024)
        return &buf
    },
}

func readSegment(r io.Reader) ([]byte, error) {
    buf := bufferPool.Get().(*[]byte)
    defer bufferPool.Put(buf)

    n, err := r.Read(*buf)
    if err != nil {
        return nil, err
    }

    result := make([]byte, n)
    copy(result, (*buf)[:n])
    return result, nil
}
```

---

## Database Patterns

### Connection Pool Sizing

```go
pool, _ := pgxpool.NewWithConfig(ctx, &pgxpool.Config{
    MaxConns:          int32(runtime.NumCPU() * 4),
    MinConns:          int32(runtime.NumCPU()),
    MaxConnLifetime:   time.Hour,
    MaxConnIdleTime:   30 * time.Minute,
    HealthCheckPeriod: time.Minute,
})
```

### Query Timeouts

```go
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()

rows, err := queries.GetMediaItems(ctx, params)
```

### Transaction Patterns

```go
func (s *Service) UpdateWithHistory(ctx context.Context, id uuid.UUID, update Update) error {
    tx, err := s.db.Begin(ctx)
    if err != nil {
        return err
    }
    defer tx.Rollback(ctx)

    qtx := s.queries.WithTx(tx)

    if err := qtx.UpdateMedia(ctx, id, update); err != nil {
        return err
    }

    if err := qtx.InsertHistory(ctx, id, update); err != nil {
        return err
    }

    return tx.Commit(ctx)
}
```

---

## API Best Practices

### Request Validation

```go
type CreateLibraryRequest struct {
    Name     string   `json:"name" validate:"required,min=1,max=255"`
    Type     string   `json:"type" validate:"required,oneof=movie tvshow music"`
    Paths    []string `json:"paths" validate:"required,min=1,dive,dirpath"`
}

func (r *CreateLibraryRequest) Validate() error {
    return validator.Struct(r)
}
```

### Error Handling

```go
type APIError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details any    `json:"details,omitempty"`
}

func writeError(w http.ResponseWriter, status int, code, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(APIError{
        Code:    code,
        Message: message,
    })
}

// Usage
switch {
case errors.Is(err, domain.ErrNotFound):
    writeError(w, 404, "NOT_FOUND", "Resource not found")
case errors.Is(err, domain.ErrUnauthorized):
    writeError(w, 401, "UNAUTHORIZED", "Invalid credentials")
default:
    writeError(w, 500, "INTERNAL_ERROR", "Internal server error")
}
```

### Pagination

```go
type PaginationParams struct {
    Limit  int    `query:"limit" validate:"min=1,max=100"`
    Offset int    `query:"offset" validate:"min=0"`
    Cursor string `query:"cursor"`
}

type PaginatedResponse[T any] struct {
    Data       []T    `json:"data"`
    Total      int64  `json:"total"`
    Limit      int    `json:"limit"`
    Offset     int    `json:"offset"`
    NextCursor string `json:"next_cursor,omitempty"`
}
```

### Versioning

```
/api/v1/movies          # Current stable
/api/v2/movies          # New version (breaking changes)
/api/unstable/movies    # Experimental
```

---

## Package Summary

| Package | Purpose |
|---------|---------|
| `github.com/failsafe-go/failsafe-go` | Circuit breaker, bulkhead, rate limiting, retry |
| `internal/infra/supervisor` | Service supervision, self-healing |
| `internal/infra/graceful` | Graceful shutdown with hooks |
| `github.com/knadh/koanf/v2` + `fsnotify` | Config hot reload |
| `github.com/maypok86/otter/v2` | L1 cache (W-TinyLFU, sub-ms) |
| `github.com/viccon/sturdyc` | API response caching |
| `github.com/redis/rueidis` | L2 cache (Dragonfly, <10ms) |
| `go.opentelemetry.io/otel` | Metrics, tracing |
| `internal/infra/health` | Health checks |


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [koanf](https://pkg.go.dev/github.com/knadh/koanf/v2) | [Local](../../sources/tooling/koanf.md) |
| [otter Cache](https://pkg.go.dev/github.com/maypok86/otter/v2) | [Local](../../sources/tooling/otter.md) |
| [rueidis](https://pkg.go.dev/github.com/redis/rueidis) | [Local](../../sources/tooling/rueidis.md) |
| [rueidis GitHub README](https://github.com/redis/rueidis) | [Local](../../sources/tooling/rueidis-guide.md) |
| [sturdyc](https://pkg.go.dev/github.com/viccon/sturdyc) | [Local](../../sources/tooling/sturdyc.md) |
| [sturdyc GitHub README](https://github.com/viccon/sturdyc) | [Local](../../sources/tooling/sturdyc-guide.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Operations](INDEX.md)

### In This Section

- [Branch Protection Rules](BRANCH_PROTECTION.md)
- [Database Auto-Healing & Consistency Restoration](DATABASE_AUTO_HEALING.md)
- [Clone repository](DEVELOPMENT.md)
- [GitFlow Workflow Guide](GITFLOW.md)
- [Revenge - Reverse Proxy & Deployment Best Practices](REVERSE_PROXY.md)
- [revenge - Setup Guide](SETUP.md)

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

---

## Quick Reference

### Startup Order
1. Load configuration
2. Initialize database pool
3. Start supervisor with services
4. Register shutdown hooks
5. Start HTTP server
6. Start config watcher

### Shutdown Order
1. Stop accepting HTTP requests
2. Drain active sessions
3. Stop supervisor (services)
4. Flush caches to disk
5. Close database connections

### Key Timeouts

| Operation | Timeout |
|-----------|---------|
| HTTP request | 30s |
| Database query | 5s |
| Shutdown | 30s |
| Circuit breaker | 30s |
| Bulkhead wait | 5s |
| Cache eviction | 30s min retention |

# Revenge - Advanced Offloading Architecture

> Keep only essential services hot, offload everything else with fast spin-up.



<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Philosophy](#philosophy)
- [Service Categories](#service-categories)
  - [🔴 Always Hot (Never Offload)](#always-hot-never-offload)
  - [🟡 Warm Standby (Connection Pooling)](#warm-standby-connection-pooling)
  - [🟢 Cold Start (Lazy Initialization)](#cold-start-lazy-initialization)
  - [⚫ Background (Async Processing)](#background-async-processing)
- [Implementation Patterns](#implementation-patterns)
  - [1. Lazy Service Factory](#1-lazy-service-factory)
  - [2. Connection Pool with Auto-Scaling](#2-connection-pool-with-auto-scaling)
  - [3. Service Health & Readiness](#3-service-health-readiness)
  - [4. Graceful Degradation](#4-graceful-degradation)
- [Offloading by Module](#offloading-by-module)
  - [Playback Service](#playback-service)
  - [Content Modules (Movie, TV, etc.)](#content-modules-movie-tv-etc)
  - [Infrastructure](#infrastructure)
- [Memory Budget](#memory-budget)
- [Configuration](#configuration)
- [Monitoring](#monitoring)
  - [Metrics to Track](#metrics-to-track)
  - [Alerts](#alerts)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Best Practices](#best-practices)
  - [DO](#do)
  - [DON'T](#dont)

<!-- TOC-END -->

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | 🔴 |  |
| Sources | 🔴 |  |
| Instructions | 🔴 |  |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |

---

## Philosophy

**Core Principle:** Minimize resource usage during idle periods while ensuring instant availability when needed.

| Category | Strategy | Target Latency |
|----------|----------|----------------|
| **Always Hot** | Never offload | 0ms |
| **Warm Standby** | Keep connections, pause processing | <10ms |
| **Cold Start** | Full offload, lazy initialization | <100ms |
| **Background** | Async, no latency requirements | N/A |

---

## Service Categories

### 🔴 Always Hot (Never Offload)

These services must always be ready for instant response:

| Service | Reason |
|---------|--------|
| HTTP Server | User-facing API responses |
| Auth Middleware | Every request needs auth |
| Session Cache | Session validation on every request |
| Config | Needed by all services |

```go
// These are always initialized at startup
fx.Provide(
    NewHTTPServer,
    NewAuthMiddleware,
    NewSessionCache,
    config.Load,
)
```

### 🟡 Warm Standby (Connection Pooling)

Keep connections alive, but minimize active processing:

| Service | Strategy |
|---------|----------|
| Database Pool | Maintain min connections, scale on demand |
| Redis/Dragonfly | Keep connection, lazy reconnect |
| Typesense | Connection pool with health checks |

```go
// Database with dynamic pool sizing
type DatabaseConfig struct {
    MinConns        int32         `koanf:"min_conns"`        // 2
    MaxConns        int32         `koanf:"max_conns"`        // 50
    MaxConnIdleTime time.Duration `koanf:"max_conn_idle"`    // 5m
    HealthCheckPeriod time.Duration `koanf:"health_check"`   // 30s
}
```

### 🟢 Cold Start (Lazy Initialization)

Only initialize when first needed:

| Service | Trigger |
|---------|---------|
| Transcoder Client | First playback request |
| Metadata Providers | First library scan |
| Search Indexer | First index operation |
| Email Service | First notification |
| OIDC Providers | First SSO login |

```go
// Lazy initialization pattern
type LazyService[T any] struct {
    once     sync.Once
    instance T
    init     func() (T, error)
    err      error
}

func (l *LazyService[T]) Get() (T, error) {
    l.once.Do(func() {
        l.instance, l.err = l.init()
    })
    return l.instance, l.err
}
```

### ⚫ Background (Async Processing)

No latency requirements, can be fully offloaded:

| Service | Pattern |
|---------|---------|
| Library Scanner | River job queue |
| Metadata Fetcher | River job queue |
| Image Processor | River job queue |
| Search Reindexer | River job queue |
| Cleanup Tasks | Scheduled River jobs |
| Analytics | Batch processing |

---

## Implementation Patterns

### 1. Lazy Service Factory

```go
// pkg/lazy/lazy.go
package lazy

import (
    "sync"
    "sync/atomic"
)

// Service wraps a lazily-initialized service.
type Service[T any] struct {
    init     func() (T, error)
    instance atomic.Pointer[T]
    once     sync.Once
    err      error
}

// New creates a lazy service wrapper.
func New[T any](init func() (T, error)) *Service[T] {
    return &Service[T]{init: init}
}

// Get returns the service instance, initializing if needed.
func (s *Service[T]) Get() (T, error) {
    s.once.Do(func() {
        instance, err := s.init()
        if err != nil {
            s.err = err
            return
        }
        s.instance.Store(&instance)
    })

    if s.err != nil {
        var zero T
        return zero, s.err
    }

    return *s.instance.Load(), nil
}

// IsInitialized returns whether the service has been initialized.
func (s *Service[T]) IsInitialized() bool {
    return s.instance.Load() != nil
}
```

### 2. Connection Pool with Auto-Scaling

```go
// internal/infra/database/pool.go
type AdaptivePool struct {
    pool    *pgxpool.Pool
    config  AdaptivePoolConfig
    metrics PoolMetrics
    mu      sync.RWMutex
}

type AdaptivePoolConfig struct {
    MinConns           int32
    MaxConns           int32
    IdleTimeout        time.Duration
    ScaleUpThreshold   float64 // 0.8 = scale when 80% utilized
    ScaleDownThreshold float64 // 0.2 = scale down when 20% utilized
    ScaleCheckInterval time.Duration
}

func (p *AdaptivePool) adjustPoolSize() {
    stats := p.pool.Stat()
    utilization := float64(stats.AcquiredConns()) / float64(stats.MaxConns())

    if utilization > p.config.ScaleUpThreshold {
        // Increase pool size (handled by pgx automatically up to MaxConns)
        p.metrics.ScaleUpEvents++
    } else if utilization < p.config.ScaleDownThreshold {
        // Let idle connections timeout naturally
        p.metrics.ScaleDownEvents++
    }
}
```

### 3. Service Health & Readiness

```go
// internal/health/health.go
type HealthChecker struct {
    checks map[string]Check
    mu     sync.RWMutex
}

type Check struct {
    Name     string
    Category string // "critical", "warm", "cold"
    Check    func(ctx context.Context) error
    Timeout  time.Duration
}

type HealthStatus struct {
    Status   string            `json:"status"` // "healthy", "degraded", "unhealthy"
    Services map[string]Status `json:"services"`
}

type Status struct {
    Healthy     bool          `json:"healthy"`
    Initialized bool          `json:"initialized"`
    Latency     time.Duration `json:"latency_ms"`
    Error       string        `json:"error,omitempty"`
}
```

### 4. Graceful Degradation

```go
// When non-critical services fail, degrade gracefully
type DegradableService struct {
    primary   Service
    fallback  Service
    degraded  atomic.Bool
}

func (d *DegradableService) Execute(ctx context.Context, op Operation) (Result, error) {
    if d.degraded.Load() {
        return d.fallback.Execute(ctx, op)
    }

    result, err := d.primary.Execute(ctx, op)
    if err != nil && isCriticalError(err) {
        d.degraded.Store(true)
        go d.scheduleRecovery()
        return d.fallback.Execute(ctx, op)
    }

    return result, err
}
```

---

## Offloading by Module

### Playback Service

| Component | Strategy | Notes |
|-----------|----------|-------|
| Session Manager | Always Hot | Active sessions in memory |
| Transcode Cache | Warm | Memory-pressure aware eviction |
| Stream Buffer | Warm | Per-session, evict on idle |
| Transcoder Client | Cold | Init on first playback |
| Bandwidth Monitor | Warm | Only for active sessions |

### Content Modules (Movie, TV, etc.)

| Component | Strategy | Notes |
|-----------|----------|-------|
| Handler | Always Hot | HTTP handlers |
| Service | Warm | Business logic |
| Repository | Warm | Database queries |
| Scanner | Cold | Only during library scan |
| Provider | Cold | Only during metadata fetch |

### Infrastructure

| Component | Strategy | Notes |
|-----------|----------|-------|
| Database | Warm | Min 2 connections |
| Cache | Warm | Persistent connection |
| Search | Cold | Init on first search |
| Jobs (River) | Warm | Polling for jobs |

---

## Memory Budget

Target memory allocation by category:

| Category | Budget | Notes |
|----------|--------|-------|
| HTTP/API | 50MB | Request handling |
| Session Cache | 100MB | Active sessions |
| Transcode Cache | 25% RAM | Configurable |
| Database Pool | 50MB | Connection overhead |
| River Workers | 100MB | Job processing |
| Misc | 200MB | Buffers, etc. |

**Total Baseline:** ~500MB + 25% RAM for transcode cache

---

## Configuration

```yaml
# configs/offloading.yaml
offloading:
  # Connection pools
  database:
    min_conns: 2
    max_conns: 50
    idle_timeout: 5m

  cache:
    pool_size: 10
    min_idle: 2

  # Lazy services
  lazy_init:
    transcoder: true
    metadata_providers: true
    email: true
    oidc: true

  # Memory limits
  memory:
    transcode_cache_percent: 25
    session_cache_max_mb: 100

  # Health checks
  health:
    interval: 30s
    timeout: 5s
```

---

## Monitoring

### Metrics to Track

```go
type OffloadingMetrics struct {
    // Service initialization
    LazyInitCount     prometheus.Counter
    LazyInitLatency   prometheus.Histogram

    // Connection pools
    PoolActiveConns   prometheus.Gauge
    PoolIdleConns     prometheus.Gauge
    PoolWaitCount     prometheus.Counter

    // Memory
    HeapInUse         prometheus.Gauge
    TranscodeCacheSize prometheus.Gauge

    // Degradation
    DegradedServices  prometheus.Gauge
    FallbackCount     prometheus.Counter
}
```

### Alerts

| Metric | Threshold | Action |
|--------|-----------|--------|
| Pool utilization | >90% | Scale up |
| Lazy init latency | >500ms | Investigate |
| Degraded services | >0 | Alert on-call |
| Memory pressure | >80% | Evict caches |


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [Dragonfly Documentation](https://www.dragonflydb.io/docs) | [Local](../../sources/infrastructure/dragonfly.md) |
| [Prometheus Go Client](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus) | [Local](../../sources/observability/prometheus.md) |
| [Prometheus Metric Types](https://prometheus.io/docs/concepts/metric_types/) | [Local](../../sources/observability/prometheus-metrics.md) |
| [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) | [Local](../../sources/tooling/river.md) |
| [Typesense API](https://typesense.org/docs/latest/api/) | [Local](../../sources/infrastructure/typesense.md) |
| [Typesense Go Client](https://github.com/typesense/typesense-go) | [Local](../../sources/infrastructure/typesense-go.md) |
| [Uber fx](https://pkg.go.dev/go.uber.org/fx) | [Local](../../sources/tooling/fx.md) |
| [koanf](https://pkg.go.dev/github.com/knadh/koanf/v2) | [Local](../../sources/tooling/koanf.md) |
| [rueidis](https://pkg.go.dev/github.com/redis/rueidis) | [Local](../../sources/tooling/rueidis.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Technical](INDEX.md)

### In This Section

- [API Reference](API.md)
- [Revenge - Audio Streaming & Progress Tracking](AUDIO_STREAMING.md)
- [Configuration Reference](CONFIGURATION.md)
- [Revenge - Frontend Architecture](FRONTEND.md)
- [Revenge - Technology Stack](TECH_STACK.md)

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

## Best Practices

### DO

- ✅ Use lazy initialization for non-critical services
- ✅ Implement graceful degradation
- ✅ Monitor initialization latency
- ✅ Set appropriate connection pool limits
- ✅ Use memory-pressure-aware caching
- ✅ Background process heavy operations
- ✅ Health check all dependencies

### DON'T

- ❌ Eagerly initialize everything at startup
- ❌ Keep unused connections open indefinitely
- ❌ Fail hard when optional services are down
- ❌ Ignore memory pressure signals
- ❌ Block request handling for background work
- ❌ Over-allocate connection pools

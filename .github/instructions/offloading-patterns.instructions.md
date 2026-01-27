---
applyTo: "**/internal/**/*.go,**/pkg/**/*.go"
---

# Advanced Offloading Patterns

> Keep essential services hot, offload everything else with fast spin-up.

## Service Categories

### 🔴 Always Hot (Never Offload)

```go
// These must be initialized at startup
// - HTTP Server
// - Auth Middleware
// - Session Cache
// - Config

// Good: Always available
fx.Provide(
    NewHTTPServer,
    NewAuthMiddleware,
    NewSessionCache,
)
```

### 🟡 Warm Standby (Connection Pooling)

```go
// Good: Maintain minimum connections, scale on demand
type DatabaseConfig struct {
    MinConns        int32         `koanf:"min_conns"`        // 2
    MaxConns        int32         `koanf:"max_conns"`        // 50
    MaxConnIdleTime time.Duration `koanf:"max_conn_idle"`    // 5m
}
```

### 🟢 Cold Start (Lazy Initialization)

```go
// Good: Only initialize when first needed
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

// Usage
var transcoderClient = &LazyService[*TranscoderClient]{
    init: func() (*TranscoderClient, error) {
        return NewTranscoderClient(config)
    },
}

// First playback request triggers init
client, err := transcoderClient.Get()
```

### ⚫ Background (Async Processing)

```go
// Good: Use River for all background work
type ScanLibraryArgs struct {
    LibraryID uuid.UUID `json:"library_id"`
}

func (ScanLibraryArgs) Kind() string { return "library.scan" }

// Queue job, don't block
client.Insert(ctx, ScanLibraryArgs{LibraryID: id}, nil)
```

## Lazy Initialization Pattern

### Basic Implementation

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

// IsInitialized checks if service has been initialized.
func (s *Service[T]) IsInitialized() bool {
    return s.instance.Load() != nil
}
```

### fx Integration

```go
// Good: Lazy provider for fx
func ProvideLazyTranscoder(config TranscoderConfig, logger *slog.Logger) *lazy.Service[*TranscoderClient] {
    return lazy.New(func() (*TranscoderClient, error) {
        logger.Info("initializing transcoder client (lazy)")
        return NewTranscoderClient(config, logger), nil
    })
}

// In module
fx.Provide(ProvideLazyTranscoder)

// Usage in handler
type PlaybackHandler struct {
    transcoder *lazy.Service[*TranscoderClient]
}

func (h *PlaybackHandler) StartPlayback(w http.ResponseWriter, r *http.Request) {
    client, err := h.transcoder.Get() // Initializes on first call
    if err != nil {
        http.Error(w, "transcoder unavailable", http.StatusServiceUnavailable)
        return
    }
    // Use client
}
```

## Connection Pool Management

### Adaptive Pool Sizing

```go
// Good: Scale based on utilization
type AdaptivePoolConfig struct {
    MinConns           int32
    MaxConns           int32
    IdleTimeout        time.Duration
    ScaleUpThreshold   float64 // 0.8
    ScaleDownThreshold float64 // 0.2
}

func (p *AdaptivePool) Monitor(ctx context.Context) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            stats := p.pool.Stat()
            utilization := float64(stats.AcquiredConns()) / float64(stats.MaxConns())

            if utilization > p.config.ScaleUpThreshold {
                p.logger.Debug("high pool utilization", "util", utilization)
            }
        }
    }
}
```

## Memory-Aware Caching

### Pressure-Based Eviction

```go
// Good: Evict based on memory pressure, not just size
type MemoryAwareCache struct {
    maxBytes     int64
    currentBytes int64

    highThreshold     float64 // 0.8 - start evicting
    criticalThreshold float64 // 0.95 - aggressive eviction
}

func (c *MemoryAwareCache) checkPressure() {
    pressure := float64(c.currentBytes) / float64(c.maxBytes)

    if pressure >= c.criticalThreshold {
        c.evictToTarget(int64(float64(c.maxBytes) * 0.8))
    } else if pressure >= c.highThreshold {
        c.evictLRU(10) // Evict 10 oldest items
    }
}
```

### Priority-Based Retention

```go
// Good: Keep high-priority items longer
type CacheItem struct {
    Key        string
    Data       []byte
    Priority   int       // 1=low, 2=normal, 3=high
    LastAccess time.Time
    IsActive   bool      // Currently in use
}

func (c *Cache) selectForEviction() *CacheItem {
    var candidate *CacheItem

    for _, item := range c.items {
        // Never evict active items unless critical
        if item.IsActive && c.pressure() < c.criticalThreshold {
            continue
        }

        // Prefer low priority, then oldest
        if candidate == nil ||
           item.Priority < candidate.Priority ||
           (item.Priority == candidate.Priority && item.LastAccess.Before(candidate.LastAccess)) {
            candidate = item
        }
    }

    return candidate
}
```

## Graceful Degradation

### Fallback Pattern

```go
// Good: Degrade gracefully when services fail
type DegradableService struct {
    primary  Service
    fallback Service
    degraded atomic.Bool
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

func (d *DegradableService) scheduleRecovery() {
    time.Sleep(30 * time.Second)
    if d.primary.HealthCheck() == nil {
        d.degraded.Store(false)
    } else {
        go d.scheduleRecovery() // Try again
    }
}
```

## Health Checks

### Tiered Health Checking

```go
// Good: Different health check strategies by category
type HealthChecker struct {
    checks map[string]Check
}

type Check struct {
    Name     string
    Category string // "critical", "warm", "cold"
    Check    func(ctx context.Context) error
    Timeout  time.Duration
}

func (h *HealthChecker) CheckAll(ctx context.Context) HealthStatus {
    status := HealthStatus{Services: make(map[string]ServiceStatus)}

    for name, check := range h.checks {
        checkCtx, cancel := context.WithTimeout(ctx, check.Timeout)

        start := time.Now()
        err := check.Check(checkCtx)
        latency := time.Since(start)
        cancel()

        status.Services[name] = ServiceStatus{
            Healthy:  err == nil,
            Category: check.Category,
            Latency:  latency,
        }
    }

    // Overall status based on critical services
    status.Status = "healthy"
    for _, svc := range status.Services {
        if !svc.Healthy && svc.Category == "critical" {
            status.Status = "unhealthy"
            break
        } else if !svc.Healthy {
            status.Status = "degraded"
        }
    }

    return status
}
```

## DO's and DON'Ts

### DO

```go
// ✅ Use lazy initialization for non-critical services
var searchClient = lazy.New(NewSearchClient)

// ✅ Set appropriate connection pool limits
pool.Config.MinConns = 2
pool.Config.MaxConns = 50

// ✅ Implement graceful degradation
if err := primary.Do(); err != nil {
    return fallback.Do()
}

// ✅ Monitor and adapt based on utilization
if poolUtilization > 0.8 {
    scaleUp()
}

// ✅ Use background jobs for heavy work
client.Insert(ctx, HeavyWorkArgs{}, nil)
```

### DON'T

```go
// ❌ Initialize everything at startup
func init() {
    searchClient = NewSearchClient() // May not be needed
    emailClient = NewEmailClient()   // May not be needed
}

// ❌ Keep connections forever
pool.Config.MaxConnIdleTime = 0 // Never close!

// ❌ Fail hard on optional service errors
if err := optionalService.Do(); err != nil {
    panic(err) // Don't do this
}

// ❌ Block requests with background work
func HandleRequest(w http.ResponseWriter, r *http.Request) {
    scanLibrary() // Blocks for minutes!
    w.Write([]byte("done"))
}
```

## Configuration

```yaml
offloading:
  # Database pool
  database:
    min_conns: 2
    max_conns: 50
    idle_timeout: 5m

  # Lazy services
  lazy_init:
    transcoder: true
    metadata_providers: true
    email: true
    search: true

  # Memory limits
  memory:
    transcode_cache_percent: 25
    high_pressure_threshold: 0.8
    critical_pressure_threshold: 0.95
```

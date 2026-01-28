---
applyTo: "**/internal/infra/cache/**/*.go"
---

# Local Cache - otter

> W-TinyLFU local cache for hot data (50% less memory than ristretto)

## Overview

Use `otter` for local in-memory caching. Combined with rueidis (remote cache), this provides a two-tier caching strategy.

**Package**: `github.com/maypok86/otter`

## When to Use

| Use Case                  | Cache                             |
| ------------------------- | --------------------------------- |
| Session tokens            | rueidis (shared across instances) |
| Metadata (TMDb responses) | otter (local) + rueidis (shared)  |
| User preferences          | otter (local, frequent access)    |
| Search suggestions        | otter (local, fast access)        |
| Rate limit counters       | rueidis (shared)                  |

## Installation

```bash
go get github.com/maypok86/otter
```

## Basic Usage

### Create Cache

```go
import "github.com/maypok86/otter"

// Create cache with max 10,000 entries and 1GB memory limit
cache, err := otter.MustBuilder[string, []byte](10_000).
    CollectStats().
    Cost(func(key string, value []byte) uint32 {
        return uint32(len(value))
    }).
    WithTTL(5 * time.Minute).
    Build()
if err != nil {
    return err
}
defer cache.Close()
```

### Get/Set Operations

```go
// Set with default TTL
cache.Set("user:123", userData)

// Set with custom TTL
cache.SetWithTTL("session:abc", sessionData, 30*time.Minute)

// Get
value, found := cache.Get("user:123")
if !found {
    // Cache miss - fetch from source
}

// Delete
cache.Delete("user:123")

// Check existence
if cache.Has("user:123") {
    // Key exists
}
```

### Cache Stats

```go
stats := cache.Stats()
log.Info("cache stats",
    "hits", stats.Hits(),
    "misses", stats.Misses(),
    "ratio", stats.Ratio(),
    "evictions", stats.EvictedCount(),
)
```

## fx Module Integration

```go
package cache

import (
    "github.com/maypok86/otter"
    "go.uber.org/fx"
)

type LocalCache struct {
    cache otter.Cache[string, []byte]
}

func NewLocalCache() (*LocalCache, error) {
    cache, err := otter.MustBuilder[string, []byte](100_000).
        CollectStats().
        Cost(func(key string, value []byte) uint32 {
            return uint32(len(value))
        }).
        WithTTL(5 * time.Minute).
        Build()
    if err != nil {
        return nil, err
    }
    return &LocalCache{cache: cache}, nil
}

func (c *LocalCache) Close() {
    c.cache.Close()
}

var Module = fx.Module("local-cache",
    fx.Provide(NewLocalCache),
    fx.Invoke(func(lc fx.Lifecycle, cache *LocalCache) {
        lc.Append(fx.Hook{
            OnStop: func(ctx context.Context) error {
                cache.Close()
                return nil
            },
        })
    }),
)
```

## Two-Tier Caching Pattern

```go
type TieredCache struct {
    local  *LocalCache    // otter
    remote *RemoteCache   // rueidis
}

func (c *TieredCache) Get(ctx context.Context, key string) ([]byte, error) {
    // Check local cache first (fast)
    if value, found := c.local.Get(key); found {
        return value, nil
    }

    // Check remote cache (shared)
    value, err := c.remote.Get(ctx, key)
    if err != nil {
        return nil, err
    }

    // Populate local cache for next access
    c.local.Set(key, value)
    return value, nil
}

func (c *TieredCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
    // Set in both caches
    c.local.SetWithTTL(key, value, ttl)
    return c.remote.Set(ctx, key, value, ttl)
}

func (c *TieredCache) Invalidate(ctx context.Context, key string) error {
    c.local.Delete(key)
    return c.remote.Delete(ctx, key)
}
```

## Configuration

```go
type LocalCacheConfig struct {
    MaxEntries int           `koanf:"max_entries"`
    MaxMemory  int64         `koanf:"max_memory_mb"`
    DefaultTTL time.Duration `koanf:"default_ttl"`
}

// Default values
var DefaultLocalCacheConfig = LocalCacheConfig{
    MaxEntries: 100_000,
    MaxMemory:  512, // 512 MB
    DefaultTTL: 5 * time.Minute,
}
```

## DO's and DON'Ts

### DO

- ✅ Use otter for frequently accessed, read-heavy data
- ✅ Set appropriate max entries based on memory
- ✅ Use cost function for variable-size values
- ✅ Collect stats for monitoring
- ✅ Close cache on shutdown

### DON'T

- ❌ Use otter for data that must be shared across instances
- ❌ Store very large values (>1MB)
- ❌ Forget to invalidate when source changes
- ❌ Use otter alone for session data (not shared)

---

## Related

- [INDEX.instructions.md](INDEX.instructions.md) - Main instruction index
- [dragonfly-cache.instructions.md](dragonfly-cache.instructions.md) - Remote cache (Tier 2)
- [sturdyc-api-cache.instructions.md](sturdyc-api-cache.instructions.md) - API cache (Tier 3)
- [otter.md](../../docs/dev/sources/tooling/otter.md) - Live otter docs

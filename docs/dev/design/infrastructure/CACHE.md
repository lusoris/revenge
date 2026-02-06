# Cache Infrastructure

<!-- DESIGN: infrastructure -->

**Package**: `internal/infra/cache`
**fx Module**: `cache.Module`

> Two-layer cache: L1 in-memory (otter/W-TinyLFU) + L2 distributed (rueidis/Dragonfly)

---

## Service Structure

```
internal/infra/cache/
├── cache.go               # Unified Cache struct (L1+L2 operations)
├── module.go              # rueidis client init, fx lifecycle hooks
├── otter.go               # L1Cache[K,V] generic wrapper (otter v2)
└── keys.go                # 20+ key prefixes, 13 TTL constants, invalidation helpers
```

## Architecture

```
Request → Cache.Get(key)
           │
           ├─ L1 hit → return (fastest, ~ns)
           │
           ├─ L1 miss → L2 check
           │              │
           │              ├─ L2 hit → populate L1, return (~ms)
           │              │
           │              └─ L2 miss → call loader, populate both layers
           │
           └─ Smart TTL: if TTL < L1 TTL, skip L1 (prevents stale reads)
```

## Cache Interface

```go
type Cache struct {
    l1     *L1Cache[string, []byte]
    l1TTL  time.Duration
    client *Client              // rueidis wrapper
    name   string               // metrics label
}

func (c *Cache) Get(ctx context.Context, key string) ([]byte, error)
func (c *Cache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
func (c *Cache) Delete(ctx context.Context, key string) error
func (c *Cache) Exists(ctx context.Context, key string) (bool, error)
func (c *Cache) Invalidate(ctx context.Context, pattern string) error

// JSON convenience
func (c *Cache) GetJSON(ctx context.Context, key string, dest interface{}) error
func (c *Cache) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error

// Cache-aside helper
func (c *Cache) CacheAside(ctx context.Context, key string, ttl time.Duration, loader func() (interface{}, error), dest interface{}) error
```

## L1 Cache (In-Memory)

```go
type L1Cache[K comparable, V any] struct { /* otter v2 */ }

func NewL1Cache[K comparable, V any](maxSize int, ttl time.Duration) (*L1Cache[K, V], error)
func (c *L1Cache[K, V]) Get(key K) (V, bool)
func (c *L1Cache[K, V]) Set(key K, value V)
func (c *L1Cache[K, V]) Delete(key K)
func (c *L1Cache[K, V]) Clear()
func (c *L1Cache[K, V]) Has(key K) bool
func (c *L1Cache[K, V]) Size() int
func (c *L1Cache[K, V]) Close()
```

**Defaults**: 10,000 entries max, 5-minute TTL, W-TinyLFU eviction

## L2 Cache (Distributed)

rueidis client connecting to Dragonfly/Redis:
- Client-side cache: 16 MiB per connection
- Dial timeout: 5s, read/write timeout: 3s
- Ring scale: 10 (1024 slots per connection)
- Blocking pool size: 128

## Cache Keys & TTLs

| Category | Prefix | TTL | Notes |
|----------|--------|-----|-------|
| Session | `session:` | 30s | Frequent validation |
| RBAC Policy | `rbac:policy:` | 5m | Infrequent changes |
| RBAC Enforce | `rbac:enforce:` | 30s | Balance perf vs freshness |
| Server Settings | `settings:server:` | 5m | Rarely changed |
| User Settings | `settings:user:` | 2m | |
| User | `user:` | 1m | |
| Movie | `movie:` | 5m | |
| Movie Metadata | `movie:meta:` | 10m | |
| Library Stats | `library:stats:` | 10m | Expensive computation |
| Search Results | `search:` | 30s | Index-driven freshness |
| Image Metadata | `image:` | 24h | Rarely changes |
| Recently Added | `movie:recent` | 2m | Homepage hot path |
| Top Rated | `movie:toprated` | 5m | |

**Invalidation helpers**: `InvalidateSession()`, `InvalidateUser()`, `InvalidateRBACForUser()`, `InvalidateAllRBAC()`, `InvalidateMovie()`, `InvalidateMovieLists()`, `InvalidateSearch()`, `InvalidateLibrary()`

## Cached Service Pattern

6 services use cache wrappers (session, user, rbac, settings, library, search):

```go
type CachedService struct {
    *Service                  // Embedded underlying service
    cache  *cache.Cache       // Unified cache
    logger *zap.Logger
}
```

- **Read**: Check cache → miss → call service → async populate cache
- **Write**: Call service → async invalidate cache
- **Nil-safe**: All wrappers check `cache == nil` before operating

## Configuration

From `config.go` `CacheConfig` (koanf namespace `cache.*`):
```yaml
cache:
  url: redis://localhost:6379    # Dragonfly/Redis URL
  enabled: false                 # Disabled by default
```

**Graceful degradation**: If L2 unavailable, L1-only mode still works.

## Dependencies

- `github.com/maypok86/otter/v2` - W-TinyLFU in-memory cache
- `github.com/redis/rueidis` - Redis/Dragonfly client
- `internal/infra/observability` - Cache hit/miss metrics

## Related Documentation

- [DATABASE.md](DATABASE.md) - Primary data store
- [JOBS.md](JOBS.md) - Job queue (uses same Dragonfly instance optionally)

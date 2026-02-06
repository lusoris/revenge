# Cache Strategy

> L1/L2 caching architecture, key conventions, and CachedService pattern. Written from code as of 2026-02-06.

---

## Architecture

```
Request
    |
    v
CachedService
    |
    +-- cache.GetJSON(key) ─────────┐
    |                               v
    |                          L1 (otter)          in-process, W-TinyLFU eviction
    |                          max 10K entries      5 min default TTL
    |                               |
    |                          miss ↓
    |                          L2 (rueidis)        Dragonfly/Redis
    |                          16 MiB client cache  per-key TTL
    |                               |
    |                          miss ↓
    +-- Service.Get() ─────────── database
    |
    +-- async: cache.SetJSON(key, result, ttl)
    |
    v
Response
```

**Three cache tiers:**

| Tier | Technology | Location | Eviction | TTL |
|------|-----------|----------|----------|-----|
| L0 | `sync.Map` | Per-client (TMDb, TVDb, Radarr, Sonarr, image) | Manual clear | Application-managed |
| L1 | otter (W-TinyLFU) | In-process | Size-based + TTL | 5 min default |
| L2 | rueidis → Dragonfly | Network (shared) | Per-key TTL | Varies (30s–24h) |

L0 is used in integration clients for HTTP response caching. L1+L2 is the formal cache infrastructure used by CachedService wrappers.

---

## Key Conventions

Format: `{domain}:{qualifier}:{id}`

```go
// internal/infra/cache/keys.go

// Prefix constants
KeyPrefixSession       = "session:"
KeyPrefixRBACEnforce   = "rbac:enforce:"
KeyPrefixServerSetting = "settings:server:"
KeyPrefixUser          = "user:"
KeyPrefixMovie         = "movie:"
KeyPrefixLibrary       = "library:"
KeyPrefixSearch        = "search:"
KeyPrefixImage         = "image:"

// Key constructors (type-safe, prevent typos)
cache.SessionKey(tokenHash)                     // "session:{hash}"
cache.RBACEnforceKey(sub, obj, act)             // "rbac:enforce:{sub}:{obj}:{act}"
cache.MovieKey(id)                              // "movie:{id}"
cache.MovieListKey(filterHash)                  // "movie:list:{hash}"
cache.MovieRecentKey(limit, offset)             // "movie:recent:{limit}:{offset}"
cache.SearchMoviesKey(queryHash)                // "search:movies:{hash}"
cache.ImageKey(imageType, size, path)           // "image:{type}:{size}:{path}"
```

Always use the key constructors from `cache/keys.go` — never build keys manually.

---

## TTL Reference

```go
// internal/infra/cache/keys.go

SessionTTL          = 30 * time.Second   // Validated frequently
RBACPolicyTTL       = 5 * time.Minute    // Policies change rarely
RBACEnforceTTL      = 30 * time.Second   // Balance perf vs freshness
ServerSettingsTTL   = 5 * time.Minute    // Settings change rarely
UserSettingsTTL     = 2 * time.Minute
UserTTL             = 1 * time.Minute
MovieTTL            = 5 * time.Minute    // Read-heavy
MovieMetaTTL        = 10 * time.Minute
LibraryStatsTTL     = 10 * time.Minute   // Expensive to compute
SearchResultsTTL    = 30 * time.Second   // Changes with indexing
ImageMetaTTL        = 24 * time.Hour     // Metadata only, not bytes
ContinueWatchingTTL = 1 * time.Minute    // Per-user, changes often
RecentlyAddedTTL    = 2 * time.Minute    // Homepage hot path
TopRatedTTL         = 5 * time.Minute
```

**Guidelines:** Short TTL (30s) for user-specific or frequently changing data. Long TTL (5–10 min) for read-heavy, rarely-changing data. 24h for immutable metadata.

---

## CachedService Wrapper

Every service that needs caching follows this pattern:

```go
// internal/service/{name}/cached_service.go

type CachedService struct {
    Service                    // Embed interface or *Service
    cache  *cache.Cache
    logger *zap.Logger
}

func NewCachedService(svc Service, c *cache.Cache, logger *zap.Logger) *CachedService {
    return &CachedService{
        Service: svc,
        cache:   c,
        logger:  logger.Named("myservice-cache"),
    }
}
```

**Read method pattern:**

```go
func (s *CachedService) Get(ctx context.Context, id uuid.UUID) (*Entity, error) {
    // 1. Graceful degradation
    if s.cache == nil {
        return s.Service.Get(ctx, id)
    }

    // 2. Try cache
    cacheKey := cache.EntityKey(id.String())
    var entity Entity
    if err := s.cache.GetJSON(ctx, cacheKey, &entity); err == nil {
        return &entity, nil  // Hit
    }

    // 3. Load from DB
    result, err := s.Service.Get(ctx, id)
    if err != nil {
        return nil, err
    }

    // 4. Populate cache async (never block response)
    go func() {
        cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
        defer cancel()
        _ = s.cache.SetJSON(cacheCtx, cacheKey, result, cache.EntityTTL)
    }()

    return result, nil
}
```

**Write method pattern** (invalidate on mutation):

```go
func (s *CachedService) Update(ctx context.Context, id uuid.UUID, params UpdateParams) (*Entity, error) {
    result, err := s.Service.Update(ctx, id, params)
    if err != nil {
        return nil, err
    }

    // Invalidate after successful write
    _ = s.cache.InvalidateEntity(ctx, id.String())

    return result, nil
}
```

**Services with CachedService wrappers:** movie, user, settings, session, rbac, library.

---

## CacheAside Helper

For one-off caching without a full CachedService wrapper:

```go
var result MyType
err := cache.CacheAside(ctx, cacheKey, ttl, func() (interface{}, error) {
    return service.LoadExpensiveData(ctx)
}, &result)
```

---

## Invalidation

Domain-specific invalidation methods on `*cache.Cache`:

```go
cache.InvalidateMovie(ctx, movieID)          // Movie + cast + crew + genres + files + lists
cache.InvalidateMovieLists(ctx)              // All list/recent/toprated caches
cache.InvalidateSession(ctx, tokenHash)      // Single session
cache.InvalidateUserSessions(ctx, userID)    // All user sessions
cache.InvalidateRBACForUser(ctx, userID)     // Roles + perms + enforce results
cache.InvalidateAllRBAC(ctx)                 // All RBAC (after policy reload)
cache.InvalidateServerSettings(ctx)          // All server settings
cache.InvalidateUserSettings(ctx, userID)    // User's settings
cache.InvalidateUser(ctx, userID)            // User + continue watching
cache.InvalidateLibrary(ctx, libraryID)      // Library + stats
cache.InvalidateSearch(ctx)                  // All search results
```

Pattern: invalidation cascades to related keys (e.g., movie invalidation also clears movie lists).

---

## L0: sync.Map Caching (Integration Clients)

Used for HTTP response caching in external API clients. Not part of the formal L1/L2 infrastructure.

```go
// internal/service/metadata/providers/tmdb/client.go
type Client struct {
    httpClient  *http.Client
    cache       sync.Map  // key: URL, value: cached response
    // ...
}
```

Present in: TMDb client, TVDb client, Radarr client, Sonarr client, shared metadata client, image service, rate limiter.

These caches are cleared via `c.cache = sync.Map{}` (full reset) or per-key deletion.

---

## Configuration

```yaml
cache:
  enabled: true
  url: "redis://localhost:6379"  # Dragonfly URL
```

L1 defaults: 10,000 entries max, 5 min TTL.
L2 defaults: 16 MiB client-side cache per connection, 128 blocking pool size.

If `cache.enabled = false`, the `Client` is created without a rueidis connection. All CachedService wrappers degrade gracefully (the `if s.cache == nil` check).

---

## Adding Caching to a New Service

1. Add key prefix + TTL + key constructor to `internal/infra/cache/keys.go`
2. Add invalidation method to `internal/infra/cache/cache.go`
3. Create `cached_service.go` in your service package (see pattern above)
4. Add `NewCachedService` to your fx module
5. Ensure write operations call the appropriate invalidation method

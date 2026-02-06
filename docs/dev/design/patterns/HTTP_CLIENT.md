# HTTP Client Pattern

> How all external HTTP clients are built: req/v3, rate limiting, caching, proxy support, and retry logic. Written from code as of 2026-02-06.

---

## Overview

All HTTP clients in the codebase use `imroc/req/v3`. There are 5 clients, each with the same core pattern but different auth, rate limits, and caching strategies.

| Client | Package | Auth | Rate Limit | Cache | Proxy |
|--------|---------|------|------------|-------|-------|
| **TMDb** | `service/metadata/providers/tmdb` | Bearer token or API key | 4 req/s, burst 10 | sync.Map (15m search, 24h detail) | Yes |
| **TVDb** | `service/metadata/providers/tvdb` | JWT token (auto-refresh) | 5 req/s, burst 10 | sync.Map (configurable TTL) | Yes |
| **Radarr** | `integration/radarr` | `X-Api-Key` header | 10 req/s, burst 20 | sync.Map (5m TTL) | No |
| **Sonarr** | `integration/sonarr` | `X-Api-Key` header | 10 req/s, burst 20 | sync.Map (5m TTL) | No |
| **Image** | `infra/image` | None | None | Disk (file-based) | Yes |

---

## Client Construction

Every client follows the same builder pattern:

```go
client := req.C().
    SetBaseURL(baseURL).
    SetTimeout(30 * time.Second).
    SetCommonRetryCount(3).
    SetCommonRetryBackoffInterval(1*time.Second, 10*time.Second)

// Optional proxy
if config.ProxyURL != "" {
    client.SetProxyURL(config.ProxyURL)
}
```

### Rate Limiting

All rate-limited clients use `golang.org/x/time/rate`:

```go
rateLimiter := rate.NewLimiter(rate.Limit(config.RateLimit), config.Burst)

// Before every request:
if err := c.rateLimiter.Wait(ctx); err != nil {
    return fmt.Errorf("rate limiter: %w", err)
}
```

- **External APIs** (TMDb, TVDb): Conservative (4-5 req/s) to respect API limits
- **Local services** (Radarr, Sonarr): Higher (10 req/s) since they're on the same network
- **Image proxy**: No rate limiting (CDN can handle it)

### Retry Logic

All clients: 3 retries with exponential backoff (1s to 10s). Exception: Image client uses fixed 1s intervals.

---

## Authentication Patterns

### API Key (TMDb v3, Radarr, Sonarr)

```go
// TMDb: query parameter
client.SetCommonQueryParam("api_key", config.APIKey)

// Radarr/Sonarr: header
client.SetCommonHeader("X-Api-Key", config.APIKey)
```

### Bearer Token (TMDb v4)

```go
client.SetBearerAuthToken(config.AccessToken)
```

### JWT with Auto-Refresh (TVDb)

TVDb uses JWT tokens obtained from a `/login` endpoint:

```go
type Client struct {
    token     string
    tokenExp  time.Time
    tokenMu   sync.RWMutex
}

func (c *Client) authenticate(ctx context.Context) error {
    c.tokenMu.RLock()
    if c.token != "" && time.Now().Before(c.tokenExp) {
        c.tokenMu.RUnlock()
        return nil
    }
    c.tokenMu.RUnlock()

    c.tokenMu.Lock()
    defer c.tokenMu.Unlock()
    // POST /login -> get JWT -> store token + expiry
}
```

Token refresh buffer: 1 hour before expiry. Mutex protects concurrent access.

---

## Caching Patterns

### In-Memory (sync.Map)

Used by TMDb, TVDb, Radarr, Sonarr clients:

```go
type CacheEntry struct {
    Data      interface{}
    ExpiresAt time.Time
}

func (c *Client) getFromCache(key string) (interface{}, bool) {
    entry, ok := c.cache.Load(key)
    if !ok { return nil, false }
    if time.Now().After(entry.ExpiresAt) {
        c.cache.Delete(key)
        return nil, false
    }
    return entry.Data, true
}
```

Cache key format: `"{type}:{query}:{params}"` with colons.

**Known issue**: `sync.Map` has unbounded memory growth. Should migrate to otter (same W-TinyLFU cache used in L1 service caching). Tracked in `.workingdir3/CODEBASE_TODOS.md` item #17.

### Disk Cache (Image Proxy)

Image service uses file-based disk cache:

```
{cacheDir}/{imageType}/{size}/{path}
```

Features:
- Content validation: MIME type whitelist (jpeg, png, gif, webp, svg+xml)
- Size limit: 10MB default
- ETag + cache-control headers for browser caching
- No TTL — files persist until manually cleared

---

## Proxy Support

External-facing clients (TMDb, TVDb, Image) support configurable proxy:

```go
// Config
type TMDbConfig struct {
    ProxyURL string  // HTTP, HTTPS, or SOCKS5 proxy URL
}

// Applied at client creation
if config.ProxyURL != "" {
    client.SetProxyURL(config.ProxyURL)
}
```

Proxy is needed for accessing metadata providers in restricted regions. Local service clients (Radarr, Sonarr) don't support proxy since they're always on the same network.

---

## Error Handling

All clients wrap errors with `fmt.Errorf("context: %w", err)` for traceability.

Common patterns:
- **404**: Return domain-specific sentinel error (`ErrMovieNotFound`, `ErrSeriesNotFound`)
- **401**: Return auth error, log for debugging
- **429**: Rate limited — client-side rate limiting should prevent this
- **5xx**: Retry via req/v3 built-in retry logic

```go
if resp.IsErrorState() {
    if resp.StatusCode == http.StatusNotFound {
        return nil, ErrMovieNotFound
    }
    return nil, fmt.Errorf("radarr API error: %d %s", resp.StatusCode, resp.Status)
}
```

---

## Adding a New HTTP Client

1. Create client struct with `req.Client`, rate limiter, and cache
2. Use `req.C()` builder pattern for construction
3. Add rate limiting with `golang.org/x/time/rate`
4. Add caching with `sync.Map` (or otter when migrated)
5. Add proxy support if the service is external
6. Wrap errors with context using `%w`
7. Add config section to `internal/config/config.go`

---

## Related Documentation

- [Arr Integration Pattern](SERVARR.md) — Radarr/Sonarr client details
- [Metadata Enrichment](METADATA.md) — TMDb/TVDb client usage
- [Metadata System (Architecture)](../architecture/METADATA_SYSTEM.md) — Provider chain
- [Cache Strategy](CACHE_STRATEGY.md) — L1/L2 service-level caching (separate from client caching)

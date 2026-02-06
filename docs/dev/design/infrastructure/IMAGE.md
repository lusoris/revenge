# Image Service

<!-- DESIGN: infrastructure -->

**Package**: `internal/infra/image`
**fx Module**: `image.Module`

> TMDb image download, dual-layer caching, and HTTP proxy with ETag support

---

## Service Structure

```
internal/infra/image/
├── service.go             # Image download, cache, proxy (http.Handler)
└── module.go              # fx module
```

## Service Interface

```go
type Service struct {
    httpClient *req.Client      // imroc/req with retries
    cacheDir   string           // Filesystem cache directory
    cacheTTL   time.Duration    // Default 7 days
    maxSize    int64            // Default 10MB
    memCache   sync.Map         // In-memory cache layer
    logger     *zap.Logger
}

func (s *Service) FetchImage(ctx context.Context, imageType, size, path string) ([]byte, string, error)
func (s *Service) StreamImage(ctx context.Context, imageType, size, path string, w http.ResponseWriter) error
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request)  // http.Handler
func (s *Service) ClearCache() error
```

## Image Types & Sizes

TMDb standard sizes:

| Type | Sizes |
|------|-------|
| Poster | w92, w154, w185, w342, w500, w780, original |
| Backdrop | w300, w780, w1280, original |
| Profile | w45, w185, h632, original |
| Logo | w45, w92, w154, w185, w300, w500, original |

## Caching

**Dual-layer**:
1. **Memory** (`sync.Map`): Fast in-process cache, no TTL eviction
2. **Filesystem**: Persistent cache with configurable TTL (default 7 days) and size limits

**Cache key**: `{type}/{size}/{path}` mapped to filesystem path under `cacheDir`.

## HTTP Proxy

Implements `http.Handler` for endpoint pattern: `/images/{type}/{size}/{path}`

**Response headers**:
- `ETag` - Generated from content hash
- `Cache-Control: public, max-age=604800, immutable`
- `Content-Type` - Validated MIME type
- Supports `If-None-Match` for conditional requests (304 Not Modified)
- CORS headers included

**Security**:
- Validates MIME types: JPEG, PNG, GIF, WebP, SVG
- Size limit: 10MB default
- Directory traversal prevention

## Configuration

No dedicated config struct. Values passed at construction:
- `BaseURL` - TMDb image base URL
- `CacheDir` - Filesystem cache directory
- `CacheTTL` - 7 days default
- `MaxSize` - 10MB default

## Dependencies

- `github.com/imroc/req/v3` - HTTP client (30s timeout, 3 retries)
- `go.uber.org/zap` - Logging

## Related Documentation

- [../services/METADATA.md](../services/METADATA.md) - Metadata service triggers image downloads

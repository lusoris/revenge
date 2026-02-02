## Table of Contents

- [Pornhub](#pornhub)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Integration Structure](#integration-structure)
    - [Data Flow](#data-flow)
    - [Provides](#provides)
  - [Implementation](#implementation)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
- [Pornhub integration](#pornhub-integration)
- [Rate limiting (very conservative)](#rate-limiting-very-conservative)
- [Caching](#caching)
- [Proxy/VPN (REQUIRED for Cloudflare)](#proxyvpn-required-for-cloudflare)
- [Headless browser (for Cloudflare bypass)](#headless-browser-for-cloudflare-bypass)
    - [Config Keys](#config-keys)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Pornhub


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration


> Integration with Pornhub

> LINK-ONLY performer channel references for QAR content
**API Base URL**: `https://www.pornhub.com`

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | ✅ | - |
| Instructions | ✅ | - |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---


## Architecture

```mermaid
flowchart TD
    node1["Performer<br/>Profile Page<br/>(Revenge UI)"]
    node2["Pornhub Link<br/>(verified URL)"]
    node3["Pornhub.com<br/>(Cloudflare)"]
    node4["View Count<br/>Subscribers<br/>Video Count"]
    node5([HTTP_CLIENT<br/>(REQUIRED<br/>proxy/VPN)])
    node6["Headless<br/>Browser<br/>(Cloudflare)"]
    node3 --> node4
    node1 --> node2
    node2 --> node3
    node4 --> node5
    node5 --> node6
```

### Integration Structure

```
internal/integration/pornhub/
├── client.go              # API client
├── types.go               # Response types
├── mapper.go              # Map external → internal types
├── cache.go               # Response caching
└── client_test.go         # Tests
```

### Data Flow

<!-- Data flow diagram -->

### Provides
<!-- Data provided by integration -->
## Implementation

### Key Interfaces

```go
// Pornhub link provider (LINK-ONLY, no content)
type PornhubProvider struct {
  httpFactory  httpclient.ClientFactory
  rateLimiter  *rate.Limiter
  cache        Cache
  browserPool  *chromedp.Pool  // For Cloudflare bypass
}

// External link provider interface
type ExternalLinkProvider interface {
  // Search for performer channel
  SearchChannel(ctx context.Context, performerName string) (*ChannelResult, error)

  // Get channel metrics
  GetChannelMetrics(ctx context.Context, channelURL string) (*ChannelMetrics, error)

  // Verify link is still valid
  VerifyLink(ctx context.Context, url string) (*LinkStatus, error)

  // Provider metadata
  ProviderName() string
  ProvidesContent() bool  // Always returns false
}

// Channel search result
type ChannelResult struct {
  URL           string `json:"url"`
  Slug          string `json:"slug"`
  Name          string `json:"name"`
  ChannelType   string `json:"type"`      // pornstar, model, channel
  ProfileImage  string `json:"image"`
  Verified      bool   `json:"verified"`  // Pornhub verified badge
  SubscriberCount int64 `json:"subscribers,omitempty"`
}

// Channel public metrics
type ChannelMetrics struct {
  URL             string    `json:"url"`
  SubscriberCount int64     `json:"subscribers"`
  VideoCount      int       `json:"videos"`
  TotalViews      int64     `json:"total_views"`
  ProfileImage    string    `json:"image"`
  FetchedAt       time.Time `json:"fetched_at"`
}

// Link verification result
type LinkStatus struct {
  URL         string `json:"url"`
  IsValid     bool   `json:"valid"`
  HTTPStatus  int    `json:"status"`
  RedirectURL string `json:"redirect,omitempty"`
  CheckedAt   time.Time `json:"checked_at"`
}
```


### Dependencies
**Go Packages**:
- `net/http` - HTTP client
- `github.com/PuerkitoBio/goquery` - HTML parsing
- `github.com/chromedp/chromedp` - Headless browser (Cloudflare)
- `golang.org/x/time/rate` - Rate limiting (0.5 req/sec)
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `github.com/riverqueue/river` - Background verification jobs
- `go.uber.org/fx` - Dependency injection

**External**:
- Pornhub website (web scraping, Cloudflare protected)

**Internal Services**:
- HTTP_CLIENT - Proxy/VPN routing (REQUIRED)

## Configuration

### Environment Variables

```bash
# Pornhub integration
PORNHUB_ENABLED=true

# Rate limiting (very conservative)
PORNHUB_RATE_LIMIT=0.5            # 1 request per 2 seconds
PORNHUB_RATE_WINDOW=1s

# Caching
PORNHUB_CACHE_TTL=168h            # 7 days

# Proxy/VPN (REQUIRED for Cloudflare)
PORNHUB_PROXY_ENABLED=true
PORNHUB_PROXY_URL=socks5://127.0.0.1:9050

# Headless browser (for Cloudflare bypass)
PORNHUB_USE_BROWSER=true
PORNHUB_BROWSER_TIMEOUT=30s
```


### Config Keys
```yaml
metadata:
  providers:
    pornhub:
      enabled: true
      rate_limit: 0.5
      rate_window: 1s
      cache_ttl: 168h

      # LINK role only (no content)
      role: link
      provides_content: false     # Explicit: NO content streaming

      # Proxy/VPN support (REQUIRED)
      proxy:
        enabled: true             # REQUIRED for Pornhub
        type: tor
        url: socks5://127.0.0.1:9050

      # Cloudflare bypass
      cloudflare:
        use_browser: true         # Use headless browser
        browser_timeout: 30s
        retry_on_challenge: true
        max_retries: 3

      # Link verification
      verification:
        enabled: true
        check_interval: 168h      # Re-verify weekly
        remove_broken_links: false  # Keep but mark as broken
```

## Related Documentation
### Design Documents
- [03_METADATA_SYSTEM](../../../architecture/03_METADATA_SYSTEM.md)
- [WHISPARR (PRIMARY for QAR)](../../servarr/WHISPARR.md)
- [STASHDB](./STASHDB.md)
- [FREEONES](./FREEONES.md)
- [HTTP_CLIENT (proxy/VPN support)](../../../services/HTTP_CLIENT.md)
- [ADULT_CONTENT_SYSTEM (QAR module)](../../../features/adult/ADULT_CONTENT_SYSTEM.md)

### External Sources
- [Go io](../../../../sources/go/stdlib/io.md) - Auto-resolved from go-io
- [River Job Queue](../../../../sources/tooling/river.md) - Auto-resolved from river
- [PuerkitoBio/goquery](https://pkg.go.dev/github.com/PuerkitoBio/goquery) - HTML parsing
- [golang.org/x/time](../../../../sources/go/x/time.md) - Rate limiting
- [chromedp](https://pkg.go.dev/github.com/chromedp/chromedp) - Headless browser for Cloudflare bypass


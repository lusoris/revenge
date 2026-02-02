## Table of Contents

- [Twitter/X](#twitterx)
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
- [Twitter/X integration (NO API)](#twitterx-integration-no-api)
- [Proxy (recommended)](#proxy-recommended)
    - [Config Keys](#config-keys)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Twitter/X


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration


> Integration with Twitter/X

> LINK-ONLY performer social media profiles for QAR content
**API Base URL**: `https://x.com`
**Authentication**: none

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
    node1["Performer<br/>Profile Page<br/>[Revenge UI]"]
    node2["Twitter/X Link<br/>[verified URL]"]
    node3["x.com<br/>[login wall]"]
    node4["Followers<br/>[limited]"]
    node5([HTTP_CLIENT<br/>[recommended<br/>proxy/VPN]])
    node3 --> node4
    node1 --> node2
    node2 --> node3
    node4 --> node5
```

### Integration Structure

```
internal/integration/twitterx/
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
// Twitter/X provider (no API, scraping only)
type TwitterProvider struct {
  httpFactory httpclient.ClientFactory
  rateLimiter *rate.Limiter
  cache       Cache
}

func (p *TwitterProvider) Platform() string { return "twitter" }
func (p *TwitterProvider) BaseURL() string { return "https://x.com" }

// Verify username exists
func (p *TwitterProvider) VerifyUsername(
  ctx context.Context,
  username string,
) (*ProfileInfo, error) {
  // Wait for rate limiter
  if err := p.rateLimiter.Wait(ctx); err != nil {
    return nil, err
  }

  url := fmt.Sprintf("https://x.com/%s", username)

  client, err := p.httpFactory.GetClientForService(ctx, "twitter")
  if err != nil {
    return nil, err
  }

  req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
  req.Header.Set("User-Agent", userAgent)

  resp, err := client.Do(req)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()

  if resp.StatusCode == 404 {
    return nil, ErrProfileNotFound
  }

  info := &ProfileInfo{
    Platform:   "twitter",
    Username:   username,
    ProfileURL: url,
    FetchedAt:  time.Now(),
  }

  // Try to extract basic info from page
  if resp.StatusCode == 200 {
    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err == nil {
      p.extractProfileInfo(doc, info)
    }
  }

  return info, nil
}

// Best-effort extraction from page
func (p *TwitterProvider) extractProfileInfo(
  doc *goquery.Document,
  info *ProfileInfo,
) {
  // Twitter uses JS rendering, so this is unreliable
  // May need to use meta tags instead

  // Try Open Graph tags
  if name, exists := doc.Find(`meta[property="og:title"]`).Attr("content"); exists {
    info.DisplayName = strings.TrimSuffix(name, " (@"+info.Username+")")
  }

  if desc, exists := doc.Find(`meta[property="og:description"]`).Attr("content"); exists {
    info.Bio = desc
  }

  if img, exists := doc.Find(`meta[property="og:image"]`).Attr("content"); exists {
    info.ProfileImage = img
  }
}
```


### Dependencies
**Go Packages**:
- `net/http` - HTTP client
- `github.com/PuerkitoBio/goquery` - HTML/meta tag parsing
- `golang.org/x/time/rate` - Rate limiting
- `github.com/riverqueue/river` - Background jobs
- `go.uber.org/fx` - Dependency injection

**External**:
- Twitter/X website (web scraping, limited)

**Internal Services**:
- HTTP_CLIENT - Proxy/VPN routing

**NOT Used**:
- Twitter API v2 (too expensive, $100/mo minimum)

## Configuration

### Environment Variables

```bash
# Twitter/X integration (NO API)
TWITTER_ENABLED=true
TWITTER_RATE_LIMIT=0.5
TWITTER_CACHE_TTL=168h

# Proxy (recommended)
TWITTER_PROXY_ENABLED=true
TWITTER_PROXY_URL=socks5://127.0.0.1:9050
```


### Config Keys
```yaml
metadata:
  providers:
    twitter:
      enabled: true
      rate_limit: 0.5
      rate_window: 1s
      cache_ttl: 168h

      role: link
      provides_content: false
      use_api: false              # API too expensive

      proxy:
        enabled: true
        type: tor
        url: socks5://127.0.0.1:9050

      verification:
        enabled: true
        check_interval: 168h
```

## Related Documentation
### Design Documents
- [03_METADATA_SYSTEM](../../../architecture/03_METADATA_SYSTEM.md)
- [FREEONES](./FREEONES.md)
- [HTTP_CLIENT (proxy/VPN support)](../../../services/HTTP_CLIENT.md)
- [ADULT_CONTENT_SYSTEM (QAR module)](../../../features/adult/ADULT_CONTENT_SYSTEM.md)

### External Sources
- [River Job Queue](../../../../sources/tooling/river.md) - Auto-resolved from river
- [PuerkitoBio/goquery](https://pkg.go.dev/github.com/PuerkitoBio/goquery) - HTML parsing
- [golang.org/x/time](../../../../sources/go/x/time.md) - Rate limiting


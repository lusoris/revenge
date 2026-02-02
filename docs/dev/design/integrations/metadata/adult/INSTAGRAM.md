## Table of Contents

- [Instagram](#instagram)
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
- [Instagram integration](#instagram-integration)
- [Proxy (recommended due to aggressive blocking)](#proxy-recommended-due-to-aggressive-blocking)
    - [Config Keys](#config-keys)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Instagram


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration


> Integration with Instagram

> LINK-ONLY performer social media profiles for QAR content
**API Base URL**: `https://www.instagram.com`

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
    node2["Instagram Link<br/>(verified URL)"]
    node3["Instagram.com<br/>(login wall)"]
    node4["Followers<br/>Post Count<br/>Bio"]
    node5([HTTP_CLIENT<br/>(recommended<br/>proxy/VPN)])
    node3 --> node4
    node1 --> node2
    node2 --> node3
    node4 --> node5
```

### Integration Structure

```
internal/integration/instagram/
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
// Social profile link provider interface
type SocialProfileProvider interface {
  // Verify username exists on platform
  VerifyUsername(ctx context.Context, username string) (*ProfileInfo, error)

  // Get basic profile info (if publicly accessible)
  GetProfileInfo(ctx context.Context, username string) (*ProfileInfo, error)

  // Verify link is still valid
  VerifyLink(ctx context.Context, url string) (bool, error)

  // Provider metadata
  Platform() string
  BaseURL() string
}

// Profile information (what's publicly accessible)
type ProfileInfo struct {
  Platform     string `json:"platform"`
  Username     string `json:"username"`
  ProfileURL   string `json:"url"`
  DisplayName  string `json:"display_name,omitempty"`
  Bio          string `json:"bio,omitempty"`
  ProfileImage string `json:"image,omitempty"`

  // Metrics (may not always be available)
  FollowerCount  *int64 `json:"followers,omitempty"`
  FollowingCount *int64 `json:"following,omitempty"`
  PostCount      *int   `json:"posts,omitempty"`

  // Verification
  IsVerified    bool      `json:"verified"`        // Platform badge
  FetchedAt     time.Time `json:"fetched_at"`
}

// Instagram provider
type InstagramProvider struct {
  httpFactory httpclient.ClientFactory
  rateLimiter *rate.Limiter
  cache       Cache
}
```


### Dependencies
**Go Packages**:
- `net/http` - HTTP client
- `github.com/PuerkitoBio/goquery` - HTML parsing
- `golang.org/x/time/rate` - Rate limiting
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `github.com/riverqueue/river` - Background jobs
- `go.uber.org/fx` - Dependency injection

**External**:
- Instagram website (very limited scraping)

**Internal Services**:
- HTTP_CLIENT - Proxy/VPN routing

## Configuration

### Environment Variables

```bash
# Instagram integration
INSTAGRAM_ENABLED=true
INSTAGRAM_RATE_LIMIT=0.5
INSTAGRAM_CACHE_TTL=168h

# Proxy (recommended due to aggressive blocking)
INSTAGRAM_PROXY_ENABLED=true
INSTAGRAM_PROXY_URL=socks5://127.0.0.1:9050
```


### Config Keys
```yaml
metadata:
  providers:
    instagram:
      enabled: true
      rate_limit: 0.5
      rate_window: 1s
      cache_ttl: 168h

      role: link                   # LINK only
      provides_content: false

      proxy:
        enabled: true              # Recommended
        type: tor
        url: socks5://127.0.0.1:9050

      verification:
        enabled: true
        check_interval: 168h       # Weekly
```

## Related Documentation
### Design Documents
- [03_METADATA_SYSTEM](../../../architecture/03_METADATA_SYSTEM.md)
- [FREEONES](./FREEONES.md)
- [HTTP_CLIENT (proxy/VPN support)](../../../services/HTTP_CLIENT.md)
- [ADULT_CONTENT_SYSTEM (QAR module)](../../../features/adult/ADULT_CONTENT_SYSTEM.md)

### External Sources
- [PuerkitoBio/goquery](https://pkg.go.dev/github.com/PuerkitoBio/goquery) - HTML parsing
- [golang.org/x/time](../../../../sources/go/x/time.md) - Rate limiting
- [River Job Queue](../../../../sources/tooling/river.md) - Background verification jobs


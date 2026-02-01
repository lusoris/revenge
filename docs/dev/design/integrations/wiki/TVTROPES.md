---
sources:
  - name: TVTropes
    url: https://tvtropes.org
    note: Main site
  - name: PuerkitoBio/goquery
    url: https://pkg.go.dev/github.com/PuerkitoBio/goquery
    note: HTML parsing for verification
  - name: golang.org/x/time
    url: ../../../sources/go/x/time.md
    note: Rate limiting
  - name: River Job Queue
    url: ../../../sources/tooling/river.md
    note: Auto-resolved from river
design_refs:
  - title: 03_METADATA_SYSTEM
    path: ../../architecture/03_METADATA_SYSTEM.md
  - title: WIKI_SYSTEM
    path: ../../features/shared/WIKI_SYSTEM.md
  - title: MOVIE_MODULE
    path: ../../features/video/MOVIE_MODULE.md
  - title: TVSHOW_MODULE
    path: ../../features/video/TVSHOW_MODULE.md
  - title: HTTP_CLIENT
    path: ../../services/HTTP_CLIENT.md
---

## Table of Contents

- [TVTropes](#tvtropes)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Integration Structure](#integration-structure)
    - [Data Flow](#data-flow)
    - [Provides](#provides)
  - [Implementation](#implementation)
    - [File Structure](#file-structure)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
    - [Config Keys](#config-keys)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)


# TVTropes


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration


> Integration with TVTropes

> ENRICHMENT link provider for trope analysis and storytelling patterns
**API Base URL**: `https://tvtropes.org`
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

### Integration Structure

```
internal/integration/tvtropes/
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

### File Structure

<!-- File structure -->

### Key Interfaces

```go
// TVTropes link provider
type TVTropesProvider struct {
  httpClient  *http.Client
  rateLimiter *rate.Limiter
  cache       Cache
}

// Link provider interface
type TropesLinkProvider interface {
  GenerateLink(ctx context.Context, content *Content) (*TropesLink, error)
  VerifyPage(ctx context.Context, url string) (bool, error)
  GetNamespace(contentType string) string
}

// Tropes link
type TropesLink struct {
  PageTitle string `json:"title"`
  URL       string `json:"url"`
  Namespace string `json:"namespace"`
  Verified  bool   `json:"verified"`
}

// Title formatter
func FormatTropesTitle(title string) string {
  // Remove articles
  title = strings.TrimPrefix(title, "The ")
  title = strings.TrimPrefix(title, "A ")
  title = strings.TrimPrefix(title, "An ")

  // Remove punctuation and special characters
  reg := regexp.MustCompile(`[^a-zA-Z0-9\s]`)
  title = reg.ReplaceAllString(title, "")

  // CamelCase
  words := strings.Fields(title)
  for i, word := range words {
    words[i] = strings.Title(strings.ToLower(word))
  }
  return strings.Join(words, "")
}
```


### Dependencies
**Go Packages**:
- `net/http` - HTTP client
- `github.com/PuerkitoBio/goquery` - HTML parsing for verification
- `golang.org/x/time/rate` - Rate limiting
- `github.com/jackc/pgx/v5` - PostgreSQL
- `github.com/riverqueue/river` - Background jobs
- `go.uber.org/fx` - DI

**External**:
- TVTropes website (no API)






## Configuration
### Environment Variables

```bash
TVTROPES_ENABLED=true
TVTROPES_RATE_LIMIT=0.5
TVTROPES_CACHE_TTL=168h
```


### Config Keys
```yaml
metadata:
  providers:
    tvtropes:
      enabled: true
      rate_limit: 0.5
      rate_window: 1s
      cache_ttl: 168h
      role: enrichment

      # Namespace mapping
      namespaces:
        movie: Film
        tvshow: Series
        anime: Anime
        animation: WesternAnimation
        book: Literature
        game: VideoGame
```





## Testing Strategy

### Unit Tests

<!-- Unit test strategy -->

### Integration Tests

<!-- Integration test strategy -->

### Test Coverage

Target: **80% minimum**







## Related Documentation
### Design Documents
- [03_METADATA_SYSTEM](../../architecture/03_METADATA_SYSTEM.md)
- [WIKI_SYSTEM](../../features/shared/WIKI_SYSTEM.md)
- [MOVIE_MODULE](../../features/video/MOVIE_MODULE.md)
- [TVSHOW_MODULE](../../features/video/TVSHOW_MODULE.md)
- [HTTP_CLIENT](../../services/HTTP_CLIENT.md)

### External Sources
- [TVTropes](https://tvtropes.org) - Main site
- [PuerkitoBio/goquery](https://pkg.go.dev/github.com/PuerkitoBio/goquery) - HTML parsing for verification
- [golang.org/x/time](../../../sources/go/x/time.md) - Rate limiting
- [River Job Queue](../../../sources/tooling/river.md) - Auto-resolved from river


## Table of Contents

- [Wikipedia](#wikipedia)
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
    - [Config Keys](#config-keys)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Wikipedia


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration


> Integration with Wikipedia

> ENRICHMENT source for encyclopedic context - MediaWiki API
**API Base URL**: `https://en.wikipedia.org/w/api.php`
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

```
┌──────────────┐
│  Revenge     │
│  Detail Page │
│  (Movie/TV)  │
└──────┬───────┘
       │ Enrichment request
       ▼
┌──────────────┐     ┌───────────────────┐
│  Wikipedia   │────▶│  Article Content  │
│  MediaWiki   │     │  - Summary        │
│  API         │     │  - Plot           │
└──────┬───────┘     │  - Production     │
       │             │  - Reception      │
┌──────┴───────┐     │  - Cast info      │
│ Rate Limiter │     └───────────────────┘
│ (polite)     │
└──────────────┘

Languages: en, de, fr, es, ja, etc.
Content: CC BY-SA license
```


### Integration Structure

```
internal/integration/wikipedia/
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
// Wikipedia enrichment provider
type WikipediaProvider struct {
  language    string
  client      *http.Client
  rateLimiter *rate.Limiter
  cache       Cache
}

// Wiki enrichment provider interface
type WikiEnrichmentProvider interface {
  Search(ctx context.Context, query string) ([]SearchResult, error)
  GetArticle(ctx context.Context, title string) (*Article, error)
  GetArticleByID(ctx context.Context, pageID int) (*Article, error)
  GetExtract(ctx context.Context, title string, sentences int) (string, error)
  GetPersonBio(ctx context.Context, name string) (*PersonBio, error)
}

// Wikipedia article
type Article struct {
  PageID       int       `json:"pageid"`
  Title        string    `json:"title"`
  Extract      string    `json:"extract"`       // Intro summary
  FullContent  string    `json:"content"`       // Full article (if fetched)
  ImageURL     string    `json:"thumbnail.source"`
  LastModified time.Time `json:"touched"`
  URL          string    `json:"fullurl"`
  Categories   []string  `json:"categories"`
}

// Extracted person info
type PersonBio struct {
  Name       string `json:"name"`
  BirthDate  string `json:"birthdate"`
  BirthPlace string `json:"birthplace"`
  Bio        string `json:"bio"`
  ImageURL   string `json:"image"`
}
```


### Dependencies
**Go Packages**:
- `net/http` - HTTP client
- `golang.org/x/time/rate` - Polite rate limiting
- `github.com/jackc/pgx/v5` - PostgreSQL
- `github.com/riverqueue/river` - Background jobs
- `go.uber.org/fx` - DI

**External**:
- Wikipedia MediaWiki API (free, no key)







## Configuration

### Environment Variables

```bash
WIKIPEDIA_ENABLED=true
WIKIPEDIA_DEFAULT_LANGUAGE=en
WIKIPEDIA_CACHE_TTL=168h    # 7 days
WIKIPEDIA_RATE_LIMIT=1      # req/sec
```


### Config Keys
```yaml
metadata:
  providers:
    wikipedia:
      enabled: true
      default_language: en
      supported_languages:
        - en
        - de
        - fr
        - es
        - ja
      rate_limit: 1
      cache_ttl: 168h
      role: enrichment
      extract_sentences: 5
```










## Related Documentation
### Design Documents
- [03_METADATA_SYSTEM](../../architecture/03_METADATA_SYSTEM.md)
- [MOVIE_MODULE](../../features/video/MOVIE_MODULE.md)
- [TVSHOW_MODULE](../../features/video/TVSHOW_MODULE.md)
- [WIKI_SYSTEM](../../features/shared/WIKI_SYSTEM.md)

### External Sources
- [Dragonfly Documentation](../../../sources/infrastructure/dragonfly.md) - Auto-resolved from dragonfly
- [MediaWiki API](../../../sources/wiki/mediawiki.md) - Auto-resolved from mediawiki-api
- [River Job Queue](../../../sources/tooling/river.md) - Auto-resolved from river
- [golang.org/x/time](../../../sources/go/x/time.md) - Rate limiting


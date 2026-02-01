## Table of Contents

- [Grand Comics Database (GCD)](#grand-comics-database-gcd)
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

# Grand Comics Database (GCD)


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration


> Integration with Grand Comics Database

> SUPPLEMENTARY historical comics database - Golden/Silver Age specialist
**API Base URL**: `https://www.comics.org/api`
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
│  Comics      │
│  Library     │
└──────┬───────┘
       │
       ├─────────────────────────────────────────┐
       │ PRIMARY                                  │ SUPPLEMENTARY
       ▼                                          ▼
┌──────────────┐                           ┌──────────────┐
│  ComicVine   │                           │     GCD      │
│  (modern)    │                           │ (historical) │
└──────────────┘                           └──────┬───────┘
                                                  │
                                           ┌──────┴───────┐
                                           │ Rate Limiter │
                                           │ (polite)     │
                                           └──────────────┘

Focus Areas:
- Golden Age (1938-1956)
- Silver Age (1956-1970)
- Defunct publishers
- Bibliographic details
```


### Integration Structure

```
internal/integration/gcd/
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
// GCD provider (supplementary to ComicVine)
type GCDProvider struct {
  client      *http.Client
  rateLimiter *rate.Limiter
  cache       Cache
}

// Supplementary comics metadata provider
type ComicsSupplementaryProvider interface {
  SearchSeries(ctx context.Context, query string, year int) ([]SeriesResult, error)
  GetSeries(ctx context.Context, gcdID int) (*Series, error)
  GetIssue(ctx context.Context, gcdID int) (*Issue, error)
  MatchHistoricalComic(ctx context.Context, title string, year int, issueNum string) (*Issue, error)
  Priority() int  // Returns 20 (after ComicVine=10)
}

// GCD Series
type Series struct {
  ID           int    `json:"id"`
  Name         string `json:"name"`
  YearBegan    int    `json:"year_began"`
  YearEnded    int    `json:"year_ended,omitempty"`
  Publisher    string `json:"publisher_name"`
  Country      string `json:"country_code"`
  Language     string `json:"language_code"`
  IssueCount   int    `json:"issue_count"`
  Notes        string `json:"notes,omitempty"`
}

// GCD Issue
type Issue struct {
  ID              int    `json:"id"`
  SeriesID        int    `json:"series"`
  Number          string `json:"number"`
  PublicationDate string `json:"publication_date"`
  Price           string `json:"price"`
  PageCount       int    `json:"page_count"`
  Notes           string `json:"notes,omitempty"`
  Indicia         string `json:"indicia_publisher"`
  EditorialCredit string `json:"editing"`
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
- GCD REST API (free, no key required)







## Configuration

### Environment Variables

```bash
GCD_ENABLED=true
GCD_RATE_LIMIT=1
GCD_CACHE_TTL=168h    # 7 days
```


### Config Keys
```yaml
metadata:
  providers:
    gcd:
      enabled: true
      rate_limit: 1
      rate_window: 1s
      cache_ttl: 168h
      role: supplementary
      priority: 20          # After ComicVine (10)
      focus:
        - golden_age        # 1938-1956
        - silver_age        # 1956-1970
        - defunct_publishers
```










## Related Documentation
### Design Documents
- [03_METADATA_SYSTEM](../../../architecture/03_METADATA_SYSTEM.md)
- [COMICS_MODULE](../../../features/comics/COMICS_MODULE.md)
- [COMICVINE (PRIMARY for comics)](./COMICVINE.md)
- [HTTP_CLIENT](../../../services/HTTP_CLIENT.md)

### External Sources
- [GCD REST API](https://www.comics.org/api/) - REST API documentation
- [pgx PostgreSQL Driver](../../../../sources/database/pgx.md) - Auto-resolved from pgx
- [golang.org/x/time](../../../../sources/go/x/time.md) - Rate limiting
- [River Job Queue](../../../../sources/tooling/river.md) - Auto-resolved from river


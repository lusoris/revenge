## Table of Contents

- [OMDb (Open Movie Database)](#omdb-open-movie-database)
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
- [OMDb API](#omdb-api)
- [Rate limiting](#rate-limiting)
- [Caching](#caching)
    - [Config Keys](#config-keys)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# OMDb (Open Movie Database)


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration


> Integration with OMDb (Open Movie Database)

> SUPPLEMENTARY ratings enrichment provider (IMDb/RT/Metacritic)
**Authentication**: api_key

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
    node1[[Metadata<br/>Service]]
    node2[(OMDb<br/>Provider)]
    node3[(OMDb API<br/>(External))]
    node4["Rate Limiter<br/>(1000/day)"]
    node1 --> node2
    node2 --> node3
    node3 --> node4
```

### Integration Structure

```
internal/integration/omdb_open_movie_database/
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
// OMDb provider implementation
type OMDbProvider struct {
  client      *OMDbClient
  rateLimiter *DailyRateLimiter
  cache       Cache
}

// Daily rate limiter (1000 requests per day)
type DailyRateLimiter struct {
  limit     int
  used      int
  resetTime time.Time
  mu        sync.Mutex
}

// Metadata provider interface
type MetadataProvider interface {
  // Fetch by IMDb ID
  GetByIMDbID(ctx context.Context, imdbID string) (*OMDbMetadata, error)

  // Search by title
  SearchByTitle(ctx context.Context, title string, year *int) (*OMDbMetadata, error)

  // Get ratings only (lightweight)
  GetRatings(ctx context.Context, imdbID string) (*Ratings, error)
}

// OMDb metadata structure
type OMDbMetadata struct {
  Title          string  `json:"Title"`
  Year           string  `json:"Year"`
  Rated          string  `json:"Rated"`
  Released       string  `json:"Released"`
  Runtime        string  `json:"Runtime"`
  Genre          string  `json:"Genre"`
  Director       string  `json:"Director"`
  Actors         string  `json:"Actors"`
  Plot           string  `json:"Plot"`
  Awards         string  `json:"Awards"`
  IMDbRating     string  `json:"imdbRating"`
  IMDbVotes      string  `json:"imdbVotes"`
  IMDbID         string  `json:"imdbID"`
  BoxOffice      string  `json:"BoxOffice"`
  Ratings        []Rating `json:"Ratings"`
}

type Rating struct {
  Source string `json:"Source"`   // "Internet Movie Database", "Rotten Tomatoes", "Metacritic"
  Value  string `json:"Value"`    // "8.8/10", "87%", "82/100"
}
```


### Dependencies
**Go Packages**:
- `net/http` - HTTP client
- `github.com/google/uuid` - UUID support
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `go.uber.org/fx` - Dependency injection

**External APIs**:
- OMDb API (free tier: 1,000 requests/day)

## Configuration

### Environment Variables

```bash
# OMDb API
OMDB_API_KEY=your_api_key_here

# Rate limiting
OMDB_DAILY_LIMIT=1000

# Caching
OMDB_CACHE_TTL=168h  # 7 days
```


### Config Keys
```yaml
metadata:
  providers:
    omdb:
      enabled: true
      api_key: ${OMDB_API_KEY}
      daily_limit: 1000
      cache_ttl: 168h  # 7 days
```

## Related Documentation
### Design Documents
- [03_METADATA_SYSTEM](../../../architecture/03_METADATA_SYSTEM.md)
- [HTTP_CLIENT (proxy/VPN support)](../../../services/HTTP_CLIENT.md)
- [MOVIE_MODULE](../../../features/video/MOVIE_MODULE.md)
- [TVSHOW_MODULE](../../../features/video/TVSHOW_MODULE.md)
- [TRAKT (alternative ratings source)](../../scrobbling/TRAKT.md)

### External Sources
- [OMDb API](../../../../sources/apis/omdb.md) - Auto-resolved from omdb
- [pgx PostgreSQL Driver](../../../../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../../../../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../../../../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [River Job Queue](../../../../sources/tooling/river.md) - Auto-resolved from river


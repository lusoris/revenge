## Table of Contents

- [TheTVDB](#thetvdb)
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
- [TheTVDB API](#thetvdb-api)
- [Episode ordering preference](#episode-ordering-preference)
- [Caching](#caching)
    - [Config Keys](#config-keys)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# TheTVDB


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration


> Integration with TheTVDB

> SUPPLEMENTARY metadata provider (fallback + enrichment) for TV shows
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

```
┌──────────────┐
│  Revenge     │
│  Metadata    │
│  Service     │
└──────┬───────┘
       │
       ├─────────────────────────────────────┐
       │ PRIMARY                             │ SUPPLEMENTARY
       ▼                                     ▼
┌──────────────┐                      ┌──────────────┐
│   Sonarr     │                      │ TheTVDB API  │
│ (LOCAL cache)│                      │  (fallback + │
│              │                      │  enrichment) │
└──────┬───────┘                      └──────┬───────┘
       │                                     │
       ▼                              ┌──────┴────────┐
┌──────────────┐                     │  HTTP_CLIENT  │
│ TheTVDB API  │                     │  (optional    │
│  (external)  │                     │   proxy/VPN)  │
└──────────────┘                     └───────┬───────┘
                                            │
                                     ┌──────┴───────┐
                                     │  JWT Token   │
                                     │   Manager    │
                                     └──────────────┘
```


### Integration Structure

```
internal/integration/thetvdb/
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
// TheTVDB provider implementation
type TVDBProvider struct {
  client      *TVDBClient
  tokenMgr    *TokenManager
  cache       Cache
}

// Token manager (auto-refresh JWT)
type TokenManager struct {
  apiKey      string
  token       string
  expiresAt   time.Time
  mu          sync.RWMutex
}

// Metadata provider interface
type MetadataProvider interface {
  // Search
  SearchSeries(ctx context.Context, query string, year *int) ([]SeriesSearchResult, error)

  // Fetch details
  GetSeriesDetails(ctx context.Context, tvdbID int) (*SeriesMetadata, error)
  GetSeriesExtended(ctx context.Context, tvdbID int) (*SeriesExtendedMetadata, error)
  GetSeasonDetails(ctx context.Context, seasonID int) (*SeasonMetadata, error)
  GetEpisodeDetails(ctx context.Context, episodeID int) (*EpisodeMetadata, error)

  // Episodes (paginated)
  GetAllEpisodes(ctx context.Context, tvdbID int, ordering string) ([]EpisodeMetadata, error)

  // Images
  GetSeriesArtwork(ctx context.Context, tvdbID int) (*ArtworkSet, error)
}

// Series metadata structure
type SeriesMetadata struct {
  TVDBID        int       `json:"id"`
  Name          string    `json:"name"`
  Overview      string    `json:"overview"`
  FirstAired    string    `json:"firstAired"`
  Status        string    `json:"status"`
  Genres        []Genre   `json:"genres"`
  Networks      []Network `json:"networks"`
  Image         string    `json:"image"`
  Banner        string    `json:"banner"`
  Rating        float64   `json:"rating"`
}

// Episode metadata
type EpisodeMetadata struct {
  TVDBID         int     `json:"id"`
  SeriesID       int     `json:"seriesId"`
  Name           string  `json:"name"`
  Overview       string  `json:"overview"`
  Aired          string  `json:"aired"`
  Runtime        int     `json:"runtime"`
  AiredSeason    int     `json:"airedSeason"`
  AiredEpisode   int     `json:"airedEpisodeNumber"`
  DVDSeason      int     `json:"dvdSeason"`
  DVDEpisode     int     `json:"dvdEpisodeNumber"`
  AbsoluteNumber int     `json:"absoluteNumber"`
  Image          string  `json:"image"`
}
```


### Dependencies
**Go Packages**:
- `net/http` - HTTP client
- `github.com/google/uuid` - UUID support
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `github.com/riverqueue/river` - Background jobs
- `go.uber.org/fx` - Dependency injection

**External APIs**:
- TheTVDB API v4 (free tier with API key)







## Configuration

### Environment Variables

```bash
# TheTVDB API
TVDB_API_KEY=your_api_key_here
TVDB_PIN=optional_pin_for_premium

# Episode ordering preference
TVDB_DEFAULT_ORDERING=default  # 'default', 'dvd', 'absolute'

# Caching
TVDB_CACHE_TTL=24h
```


### Config Keys
```yaml
metadata:
  providers:
    tvdb:
      enabled: true
      api_key: ${TVDB_API_KEY}
      pin: ${TVDB_PIN}
      default_ordering: default
      cache_ttl: 24h

      # SUPPLEMENTARY role configuration
      role: supplementary  # fallback + enrichment

      # Proxy/VPN support (OPTIONAL - must be setup and enabled)
      proxy:
        enabled: false           # Must explicitly enable
        type: tor                # 'http', 'socks5', 'tor', 'vpn'
        url: socks5://127.0.0.1:9050  # Tor SOCKS5 proxy (if type=tor/socks5)
        interface: tun0          # VPN interface (if type=vpn)
```










## Related Documentation
### Design Documents
- [03_METADATA_SYSTEM](../../../architecture/03_METADATA_SYSTEM.md)
- [SONARR (PRIMARY for TV shows)](../../servarr/SONARR.md)
- [HTTP_CLIENT (proxy/VPN support)](../../../services/HTTP_CLIENT.md)
- [TVSHOW_MODULE](../../../features/video/TVSHOW_MODULE.md)

### External Sources
- [go-blurhash](../../../../sources/media/go-blurhash.md) - Auto-resolved from go-blurhash
- [pgx PostgreSQL Driver](../../../../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../../../../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../../../../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [River Job Queue](../../../../sources/tooling/river.md) - Auto-resolved from river
- [TheTVDB API](../../../../sources/apis/thetvdb.md) - Auto-resolved from thetvdb
- [Typesense API](../../../../sources/infrastructure/typesense.md) - Auto-resolved from typesense
- [Typesense Go Client](../../../../sources/infrastructure/typesense-go.md) - Auto-resolved from typesense-go


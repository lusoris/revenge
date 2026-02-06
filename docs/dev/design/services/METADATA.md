## Table of Contents

- [Metadata Service](#metadata-service)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Service Structure](#service-structure)
    - [Dependencies](#dependencies)
    - [Provides](#provides)
    - [Component Diagram](#component-diagram)
  - [Implementation](#implementation)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
    - [Config Keys](#config-keys)
  - [API Endpoints](#api-endpoints)
- [Search](#search)
- [Match](#match)
- [Refresh](#refresh)
- [Providers](#providers)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Metadata Service

<!-- DESIGN: services, README, test_output_claude, test_output_wiki -->


**Created**: 2026-01-31
**Status**: 🟡 In Progress
**Category**: service


> > Metadata aggregation with PRIMARY (Arr) and SUPPLEMENTARY (external) providers

**Package**: `internal/service/metadata`
**fx Module**: `metadatafx.Module`

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | ✅ | - |
| Instructions | ✅ | - |
| Code | 🟡 Partial | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: 🟡 In Progress


---


## Architecture

```mermaid
flowchart LR
    subgraph Layer1["Layer 1"]
        node1(["Client<br/>(Web/App)"])
        node2[["API Handler<br/>(ogen)"]]
        node3[["Metadata<br/>Service"]]
    end

    subgraph Layer2["Layer 2"]
        node4["Priority<br/>Chain"]
    end

    subgraph Layer3["Layer 3"]
        node5[("L1 Cache L2 Cache<br/>Dragonfly<br/>(distrib)")]
    end

    subgraph Layer4["Layer 4"]
        node6["Otter<br/>(in-mem)"]
        node7["Dragonfly<br/>(distrib)"]
    end

    subgraph Layer5["Layer 5"]
        node8[["PRIMARY (local Arr services)<br/>Lidarr Chaptarr Whisparr<br/>(music) (books) (QAR)"]]
    end

    subgraph Layer6["Layer 6"]
        node9[("Radarr/Sonarr<br/>(LOCAL cache)")]
        node10["Lidarr<br/>(music)"]
        node11["Chaptarr<br/>(books)"]
        node12["Whisparr<br/>(QAR)"]
    end

    subgraph Layer7["Layer 7"]
        node13[("TMDb/TVDB<br/>(via proxy)")]
        node14["MusicBrainz<br/>(via proxy)"]
        node15["OpenLibrary<br/>(via proxy)"]
        node16[("StashDB<br/>(via proxy)")]
    end

    subgraph Layer8["Layer 8"]
        node17(["HTTP_CLIENT<br/>(optional<br/>proxy/VPN)"])
    end

    %% Connections
    node3 --> node4
    node4 --> node5
    node5 --> node6
    node7 --> node8
    node8 --> node9
    node12 --> node13
    node16 --> node17

    %% Styling
    style Layer1 fill:#1976D2,stroke:#1976D2,color:#fff
    style Layer2 fill:#388E3C,stroke:#388E3C,color:#fff
    style Layer3 fill:#7B1FA2,stroke:#7B1FA2,color:#fff
    style Layer4 fill:#F57C00,stroke:#F57C00,color:#fff
    style Layer5 fill:#C2185B,stroke:#C2185B,color:#fff
    style Layer6 fill:#00796B,stroke:#00796B,color:#fff
    style Layer7 fill:#1976D2,stroke:#1976D2,color:#fff
    style Layer8 fill:#388E3C,stroke:#388E3C,color:#fff
```

### Service Structure

```
internal/service/metadata/
├── doc.go                 # Package documentation
├── errors.go              # Sentinel errors (ErrNotFound, ErrRateLimited, etc.)
├── provider.go            # Provider + specialized interfaces
├── service.go             # Service interface (31 methods)
├── types.go               # Domain types (MovieMetadata, Credits, etc.)
├── adapters/              # Content module adapters
│   ├── movie/adapter.go   # movie.MetadataProvider implementation
│   └── tvshow/adapter.go  # tvshow.MetadataProvider implementation
├── jobs/                  # River job workers for async refresh
├── metadatafx/module.go   # fx module + provider registration
└── providers/
    ├── tmdb/              # TMDb provider (priority 100)
    └── tvdb/              # TVDb provider (priority 80)
```

### Dependencies
**Go Packages**:
- `github.com/google/uuid`
- `github.com/jackc/pgx/v5`
- `github.com/imroc/req/v3` - HTTP client for provider APIs (✅ implemented)
- `github.com/riverqueue/river` - Background metadata refresh jobs (✅ implemented)
- `sync.Map` (Go stdlib) - Per-provider in-memory cache with TTL (✅ current IST)
- `github.com/maypok86/otter` - L1 cache (🔴 planned, replaces sync.Map)
- `github.com/redis/rueidis` - L2 cache via Dragonfly (🔴 planned)
- `go.uber.org/fx`

**PRIMARY APIs** (Arr services - local, no proxy):
- Radarr API v3 (http://localhost:7878/api/v3) - Movies ✅
- Sonarr API v3 (http://localhost:8989/api/v3) - TV Shows ✅
- Lidarr API v1 (http://localhost:8686/api/v1) - Music (🔴 planned)
- Chaptarr/Readarr API (http://localhost:8787/api/v1) - Books/Audiobooks (🔴 planned)
- Whisparr API v3 (http://localhost:6969/api/v3) - QAR/Adult (🔴 planned, separate doc)

**SUPPLEMENTARY APIs** (external, via optional proxy/VPN):
- TMDb API v3 (https://api.themoviedb.org/3/) ✅ (priority 100, req/v3 + rate limiter)
- TheTVDB API v4 (https://api4.thetvdb.com/v4/) ✅ (priority 80, req/v3 + rate limiter)
- MusicBrainz API v2 (https://musicbrainz.org/ws/2/) (🔴 planned)
- OpenLibrary API (https://openlibrary.org/api/) (🔴 planned)
- StashDB GraphQL API (https://stashdb.org/graphql) (🔴 planned, QAR-specific)
- Last.fm API (https://ws.audioscrobbler.com/2.0/) (🔴 planned)

### Provides
<!-- Service provides -->

### Component Diagram

<!-- Component diagram -->
## Implementation

### Key Interfaces (from code)

```go
// Service interface (27 methods across 7 categories)
// Source: internal/service/metadata/service.go
type Service interface {
  // Movie (8 methods)
  SearchMovie(ctx context.Context, query string, opts SearchOptions) ([]MovieSearchResult, error)
  GetMovieMetadata(ctx context.Context, tmdbID int32, languages []string) (*MovieMetadata, error)
  GetMovieCredits(ctx context.Context, tmdbID int32) (*Credits, error)
  GetMovieImages(ctx context.Context, tmdbID int32) (*Images, error)
  GetMovieReleaseDates(ctx context.Context, tmdbID int32) ([]ReleaseDate, error)
  GetMovieExternalIDs(ctx context.Context, tmdbID int32) (*ExternalIDs, error)
  GetSimilarMovies(ctx context.Context, tmdbID int32, opts SearchOptions) ([]MovieSearchResult, int, error)
  GetMovieRecommendations(ctx context.Context, tmdbID int32, opts SearchOptions) ([]MovieSearchResult, int, error)

  // TV Show (8 methods)
  SearchTVShow(ctx context.Context, query string, opts SearchOptions) ([]TVShowSearchResult, error)
  GetTVShowMetadata(ctx context.Context, tmdbID int32, languages []string) (*TVShowMetadata, error)
  GetTVShowCredits(ctx context.Context, tmdbID int32) (*Credits, error)
  GetTVShowImages(ctx context.Context, tmdbID int32) (*Images, error)
  GetTVShowContentRatings(ctx context.Context, tmdbID int32) ([]ContentRating, error)
  GetTVShowExternalIDs(ctx context.Context, tmdbID int32) (*ExternalIDs, error)
  GetSeasonMetadata(ctx context.Context, tmdbID int32, seasonNum int, languages []string) (*SeasonMetadata, error)
  GetEpisodeMetadata(ctx context.Context, tmdbID int32, seasonNum, episodeNum int, languages []string) (*EpisodeMetadata, error)

  // Person (4 methods)
  SearchPerson(ctx context.Context, query string, opts SearchOptions) ([]PersonSearchResult, error)
  GetPersonMetadata(ctx context.Context, tmdbID int32, languages []string) (*PersonMetadata, error)
  GetPersonCredits(ctx context.Context, tmdbID int32) (*PersonCredits, error)
  GetPersonImages(ctx context.Context, tmdbID int32) (*Images, error)

  // Collection (1 method)
  GetCollectionMetadata(ctx context.Context, tmdbID int32, languages []string) (*CollectionMetadata, error)

  // Image (1 method)
  GetImageURL(path string, size ImageSize) string

  // Refresh (2 methods) - triggers async River jobs
  RefreshMovie(ctx context.Context, movieID uuid.UUID) error
  RefreshTVShow(ctx context.Context, seriesID uuid.UUID) error

  // Management (3 methods)
  ClearCache()
  RegisterProvider(provider Provider)
  GetProviders() []Provider
}

// Base Provider interface (all providers implement this)
// Source: internal/service/metadata/provider.go
type Provider interface {
  ID() ProviderID              // "tmdb", "tvdb"
  Name() string
  Priority() int               // Higher = preferred (TMDb=100, TVDb=80)
  SupportsMovies() bool
  SupportsTVShows() bool
  SupportsPeople() bool
  SupportsLanguage(lang string) bool
  ClearCache()
}

// Specialized interfaces: MovieProvider, TVShowProvider, PersonProvider,
// ImageProvider, CollectionProvider (see METADATA_SYSTEM.md)
```


## Configuration

### Current Config (from code) ✅

**Application-level** (`config.go` via koanf):
```yaml
movie:
  tmdb:
    api_key: ""                # TMDb API key (v3)
    rate_limit: 40             # Requests per 10 seconds
    cache_ttl: 5m              # TMDb response cache TTL
    proxy_url: ""              # Optional SOCKS5/HTTP proxy

integrations:
  radarr:
    enabled: false
    base_url: http://localhost:7878
    api_key: ""
    auto_sync: false
    sync_interval: 300         # seconds
  sonarr:
    enabled: false
    base_url: http://localhost:8989
    api_key: ""
    auto_sync: false
    sync_interval: 300
```

**Provider-level** (in-code Config structs, passed via fx):

`tmdb.Config` (`providers/tmdb/client.go`):
- `APIKey`, `AccessToken` (v4, optional)
- `RateLimit` (default: 4.0 req/s = 40/10s), `Burst` (default: 10)
- `CacheTTL` (default: 24h), `Timeout` (default: 30s)
- `ProxyURL`, `RetryCount` (default: 3)

`tvdb.Config` (`providers/tvdb/client.go`):
- `APIKey`, `PIN` (optional subscriber PIN)
- `RateLimit` (default: 5.0 req/s), `Burst` (default: 10)
- `CacheTTL` (default: 24h), `Timeout` (default: 30s)
- `ProxyURL`, `RetryCount` (default: 3)

`metadatafx.Config` (`metadatafx/module.go`):
- `DefaultLanguages` (default: `["en", "de", "fr", "es", "ja"]`)
- `EnableProviderFallback` (default: `true`)
- `EnableEnrichment` (default: `false`)
- `TMDbAPIKey`, `TMDbProxyURL`
- `TVDbAPIKey`, `TVDbPIN`

### Planned Config (🔴 not yet in config.go)

**Environment Variables** (🔴 planned - currently env vars are mapped via koanf to `movie.tmdb.*`):
```bash
METADATA_TMDB_API_KEY=your_api_key
METADATA_TVDB_API_KEY=your_api_key
METADATA_LASTFM_API_KEY=your_api_key
METADATA_REFRESH_INTERVAL=168h  # 7 days
METADATA_CACHE_TTL=24h
```

**Planned unified `metadata:` namespace** (🔴 not yet implemented - TMDb currently under `movie.tmdb`, Arr under `integrations.*`):
```yaml
metadata:
  # Priority chain configuration
  priority:
    - cache_l1      # Otter (in-memory)
    - cache_l2      # Dragonfly (distributed)
    - arr           # Arr services (PRIMARY)
    - external      # External APIs (SUPPLEMENTARY)

  # PRIMARY providers (Arr services)
  arr:
    radarr:
      enabled: true
      url: http://localhost:7878
      api_key: ${RADARR_API_KEY}
    sonarr:
      enabled: true
      url: http://localhost:8989
      api_key: ${SONARR_API_KEY}
    lidarr:
      enabled: true
      url: http://localhost:8686
      api_key: ${LIDARR_API_KEY}
    chaptarr:
      enabled: true
      url: http://localhost:8787
      api_key: ${CHAPTARR_API_KEY}
    whisparr:
      enabled: true
      url: http://localhost:6969
      api_key: ${WHISPARR_API_KEY}

  # SUPPLEMENTARY providers (external APIs)
  external:
    tmdb:
      api_key: ${TMDB_API_KEY}
      enabled: true
      rate_limit: 40
      proxy: tor  # Optional proxy/VPN routing
    tvdb:
      api_key: ${TVDB_API_KEY}
      enabled: true
      proxy: tor
    musicbrainz:
      enabled: true
      rate_limit: 1
      proxy: false  # No proxy needed
    openlib:
      enabled: true
      proxy: false
    stashdb:
      enabled: true
      proxy: tor  # Privacy for adult content

  refresh:
    interval: 168h
    auto_refresh: true

  cache:
    l1_ttl: 5m   # Otter cache
    l2_ttl: 24h  # Dragonfly cache
```

## API Endpoints
```
# Search
GET    /api/v1/metadata/search?provider=tmdb&query=Inception # Search provider

# Match
POST   /api/v1/metadata/match            # Match content to external ID
DELETE /api/v1/metadata/match/:id        # Remove match

# Refresh
POST   /api/v1/metadata/refresh/:id      # Trigger metadata refresh

# Providers
GET    /api/v1/metadata/providers        # List providers
PUT    /api/v1/metadata/providers/:name  # Configure provider
```

**Example Search Response**:
```json
{
  "results": [
    {
      "external_id": "27205",
      "title": "Inception",
      "release_date": "2010-07-16",
      "overview": "A thief who steals corporate secrets...",
      "poster_url": "https://image.tmdb.org/t/p/w500/9gk7adHYeDvHkCSEqAvQNLV5Uge.jpg",
      "rating": 8.4
    }
  ]
}
```

## Related Documentation
### Design Documents
- [03_METADATA_SYSTEM (PRIMARY architecture doc)](../architecture/METADATA_SYSTEM.md)
- [DRAGONFLY (cache architecture)](../integrations/infrastructure/DRAGONFLY.md)
- [HTTP_CLIENT (proxy/VPN support)](./HTTP_CLIENT.md)
- [RADARR (PRIMARY for movies)](../integrations/servarr/RADARR.md)
- [SONARR (PRIMARY for TV shows)](../integrations/servarr/SONARR.md)
- [LIDARR (PRIMARY for music)](../integrations/servarr/LIDARR.md)
- [CHAPTARR (PRIMARY for books/audiobooks)](../integrations/servarr/CHAPTARR.md)
- [WHISPARR (PRIMARY for QAR/adult)](../integrations/servarr/WHISPARR.md)

### External Sources
- [Uber fx](../../sources/tooling/fx.md) - Auto-resolved from fx
- [Last.fm API](../../sources/apis/lastfm.md) - Auto-resolved from lastfm-api
- [River Job Queue](../../sources/tooling/river.md) - Auto-resolved from river


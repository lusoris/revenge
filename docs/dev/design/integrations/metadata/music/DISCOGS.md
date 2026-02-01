## Table of Contents

- [Discogs](#discogs)
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
- [Discogs API](#discogs-api)
- [Rate limiting](#rate-limiting)
- [Caching](#caching)
    - [Config Keys](#config-keys)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Discogs


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration


> Integration with Discogs

> SUPPLEMENTARY enrichment provider (vinyl/CD releases, marketplace, credits)
**API Base URL**: `https://api.discogs.com`
**Authentication**: oauth

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
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│  Metadata   │────▶│   Discogs    │────▶│  Discogs    │
│  Service    │◀────│   Provider   │◀────│     API     │
└─────────────┘     └──────────────┘     └─────────────┘
```


### Integration Structure

```
internal/integration/discogs/
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
// Discogs provider implementation
type DiscogsProvider struct {
  client      *DiscogsClient
  token       string  // Personal Access Token
  cache       Cache
}

// Metadata provider interface
type MetadataProvider interface {
  // Search
  SearchRelease(ctx context.Context, artist, title string, year *int) ([]ReleaseSearchResult, error)
  SearchArtist(ctx context.Context, name string) ([]ArtistSearchResult, error)

  // Fetch details
  GetRelease(ctx context.Context, releaseID int) (*ReleaseDetails, error)
  GetMaster(ctx context.Context, masterID int) (*MasterRelease, error)
  GetArtist(ctx context.Context, artistID int) (*ArtistDetails, error)
}

// Release details from Discogs
type ReleaseDetails struct {
  ID          int      `json:"id"`
  Title       string   `json:"title"`
  Artists     []Artist `json:"artists"`
  Year        int      `json:"year"`
  Country     string   `json:"country"`
  Genres      []string `json:"genres"`
  Styles      []string `json:"styles"`
  Formats     []Format `json:"formats"`
  Labels      []Label  `json:"labels"`
  Tracklist   []Track  `json:"tracklist"`
  Credits     []Credit `json:"extraartists"`  // Full credits
  Images      []Image  `json:"images"`
}

type Credit struct {
  Name string `json:"name"`
  Role string `json:"role"`  // "Producer", "Engineer", "Mastering"
}
```


### Dependencies
**Go Packages**:
- `net/http` - HTTP client
- `github.com/google/uuid` - UUID support
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `go.uber.org/fx` - Dependency injection

**External APIs**:
- Discogs API v2 (free with registration)







## Configuration

### Environment Variables

```bash
# Discogs API
DISCOGS_TOKEN=your_personal_access_token_here

# Rate limiting
DISCOGS_RATE_LIMIT=60  # requests per minute

# Caching
DISCOGS_CACHE_TTL=168h  # 7 days
```


### Config Keys
```yaml
metadata:
  providers:
    discogs:
      enabled: true
      token: ${DISCOGS_TOKEN}
      rate_limit: 60  # requests/minute
      cache_ttl: 168h
```










## Related Documentation
### Design Documents
- [03_METADATA_SYSTEM](../../../architecture/03_METADATA_SYSTEM.md)
- [LIDARR (PRIMARY for music)](../../servarr/LIDARR.md)
- [HTTP_CLIENT (proxy/VPN support)](../../../services/HTTP_CLIENT.md)
- [MUSIC_MODULE](../../../features/music/MUSIC_MODULE.md)

### External Sources
- [Discogs API](../../../../sources/apis/discogs.md) - Auto-resolved from discogs
- [Last.fm API](../../../../sources/apis/lastfm.md) - Auto-resolved from lastfm-api


---
sources:
  - name: Discogs API
    url: ../../../../sources/apis/discogs.md
    note: Auto-resolved from discogs
  - name: Last.fm API
    url: ../../../../sources/apis/lastfm.md
    note: Auto-resolved from lastfm-api
design_refs:
  - title: 03_METADATA_SYSTEM
    path: ../../../architecture/03_METADATA_SYSTEM.md
  - title: LIDARR (PRIMARY for music)
    path: ../../servarr/LIDARR.md
  - title: HTTP_CLIENT (proxy/VPN support)
    path: ../../../services/HTTP_CLIENT.md
  - title: MUSIC_MODULE
    path: ../../../features/music/MUSIC_MODULE.md
---

## Table of Contents

- [Discogs](#discogs)
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
- [Discogs API](#discogs-api)
- [Rate limiting](#rate-limiting)
- [Caching](#caching)
    - [Config Keys](#config-keys)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
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

### File Structure

<!-- File structure -->

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





## Testing Strategy

### Unit Tests

<!-- Unit test strategy -->

### Integration Tests

<!-- Integration test strategy -->

### Test Coverage

Target: **80% minimum**







## Related Documentation
### Design Documents
- [03_METADATA_SYSTEM](../../../architecture/03_METADATA_SYSTEM.md)
- [LIDARR (PRIMARY for music)](../../servarr/LIDARR.md)
- [HTTP_CLIENT (proxy/VPN support)](../../../services/HTTP_CLIENT.md)
- [MUSIC_MODULE](../../../features/music/MUSIC_MODULE.md)

### External Sources
- [Discogs API](../../../../sources/apis/discogs.md) - Auto-resolved from discogs
- [Last.fm API](../../../../sources/apis/lastfm.md) - Auto-resolved from lastfm-api


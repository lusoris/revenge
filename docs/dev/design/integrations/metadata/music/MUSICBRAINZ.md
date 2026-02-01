

---
sources:
  - name: go-blurhash
    url: ../../../../sources/media/go-blurhash.md
    note: Auto-resolved from go-blurhash
  - name: Last.fm API
    url: ../../../../sources/apis/lastfm.md
    note: Auto-resolved from lastfm-api
  - name: MusicBrainz API
    url: ../../../../sources/apis/musicbrainz.md
    note: Auto-resolved from musicbrainz
design_refs:
  - title: 03_METADATA_SYSTEM
    path: ../../../architecture/03_METADATA_SYSTEM.md
  - title: LIDARR (PRIMARY for music)
    path: ../../servarr/LIDARR.md
  - title: HTTP_CLIENT (proxy/VPN support)
    path: ../../../services/HTTP_CLIENT.md
  - title: MUSIC_MODULE
    path: ../../../features/music/MUSIC_MODULE.md
  - title: LASTFM (enrichment metadata)
    path: ./LASTFM.md
  - title: LISTENBRAINZ (scrobbling)
    path: ../../scrobbling/LISTENBRAINZ.md
---

## Table of Contents

- [MusicBrainz](#musicbrainz)
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
- [MusicBrainz](#musicbrainz)
- [Cover Art Archive](#cover-art-archive)
- [AcoustID (fingerprinting)](#acoustid-fingerprinting)
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


# MusicBrainz


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: integration


> Integration with MusicBrainz

> SUPPLEMENTARY metadata provider (fallback + enrichment) for music
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
internal/integration/musicbrainz/
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
// MusicBrainz provider implementation
type MusicBrainzProvider struct {
  client       *MBClient
  rateLimiter  *rate.Limiter
  coverArtClient *CoverArtClient
  acoustIDClient *AcoustIDClient
  cache        Cache
}

// Metadata provider interface
type MetadataProvider interface {
  // Search
  SearchArtist(ctx context.Context, query string) ([]ArtistSearchResult, error)
  SearchAlbum(ctx context.Context, query string, artist string) ([]AlbumSearchResult, error)
  SearchTrack(ctx context.Context, query string, artist string) ([]TrackSearchResult, error)

  // Lookup by MBID
  GetArtist(ctx context.Context, mbid uuid.UUID) (*ArtistMetadata, error)
  GetAlbum(ctx context.Context, mbid uuid.UUID) (*AlbumMetadata, error)
  GetTrack(ctx context.Context, mbid uuid.UUID) (*TrackMetadata, error)

  // Cover art
  GetCoverArt(ctx context.Context, releaseMBID uuid.UUID) ([]CoverImage, error)

  // Fingerprinting
  LookupByFingerprint(ctx context.Context, fingerprint string, duration int) (*TrackMatch, error)
}

// Artist metadata structure
type ArtistMetadata struct {
  MBID       uuid.UUID `json:"id"`
  Name       string    `json:"name"`
  SortName   string    `json:"sort-name"`
  Country    string    `json:"country"`
  Type       string    `json:"type"`        // "Person", "Group"
  Gender     string    `json:"gender"`
  LifeSpan   *LifeSpan `json:"life-span"`
  Genres     []Genre   `json:"genres"`
  Tags       []Tag     `json:"tags"`
}

// Album metadata
type AlbumMetadata struct {
  MBID          uuid.UUID   `json:"id"`
  Title         string      `json:"title"`
  ArtistCredit  []Artist    `json:"artist-credit"`
  FirstReleased string      `json:"first-release-date"`
  PrimaryType   string      `json:"primary-type"`    // "Album", "EP", "Single"
  SecondaryTypes []string   `json:"secondary-types"` // "Compilation", "Live"
  TrackCount    int         `json:"track-count"`
}
```


### Dependencies

**Go Packages**:
- `net/http` - HTTP client
- `golang.org/x/time/rate` - Rate limiting (1 request/second)
- `github.com/google/uuid` - UUID support (MBIDs are UUIDs)
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `github.com/bbrks/go-blurhash` - Blurhash for cover art
- `go.uber.org/fx` - Dependency injection

**External APIs**:
- MusicBrainz API v2 (free, no key required)
- Cover Art Archive (free)
- AcoustID API (free with API key)






## Configuration
### Environment Variables

```bash
# MusicBrainz
MUSICBRAINZ_USER_AGENT="Revenge/1.0.0 (https://example.com)"

# Cover Art Archive
COVERART_ENABLED=true

# AcoustID (fingerprinting)
ACOUSTID_API_KEY=your_api_key_here
ACOUSTID_ENABLED=true

# Rate limiting
MUSICBRAINZ_RATE_LIMIT=1  # requests per second

# Caching
MUSICBRAINZ_CACHE_TTL=168h  # 7 days
```


### Config Keys

```yaml
metadata:
  providers:
    musicbrainz:
      enabled: true
      user_agent: "Revenge/1.0.0 (https://example.com)"
      rate_limit: 1  # requests/second
      cache_ttl: 168h

      coverart:
        enabled: true

      acoustid:
        enabled: true
        api_key: ${ACOUSTID_API_KEY}
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
- [LASTFM (enrichment metadata)](./LASTFM.md)
- [LISTENBRAINZ (scrobbling)](../../scrobbling/LISTENBRAINZ.md)

### External Sources
- [go-blurhash](../../../../sources/media/go-blurhash.md) - Auto-resolved from go-blurhash
- [Last.fm API](../../../../sources/apis/lastfm.md) - Auto-resolved from lastfm-api
- [MusicBrainz API](../../../../sources/apis/musicbrainz.md) - Auto-resolved from musicbrainz


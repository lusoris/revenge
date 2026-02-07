# TheTVDB Integration

<!-- DESIGN: integrations/metadata/video -->

**Package**: `internal/service/metadata/providers/tvdb`
**API**: TVDb API v4 (`https://api4.thetvdb.com/v4`)

> Secondary metadata provider for TV shows and people, with JWT authentication

---

## Module Structure

```
internal/service/metadata/providers/tvdb/
├── client.go    # HTTP client (req, rate limiter, L1Cache/otter, JWT auth)
├── types.go     # TVDb API v4 response types
├── provider.go  # Provider interface implementation
└── mapping.go   # TVDb responses → metadata domain types
```

## Provider

Implements 3 metadata interfaces:

```go
var (
    _ metadata.Provider       = (*Provider)(nil)
    _ metadata.TVShowProvider = (*Provider)(nil)
    _ metadata.PersonProvider = (*Provider)(nil)
)

type Provider struct {
    client   *Client
    priority int  // 80 (secondary to TMDb)
}
```

- `ID()` → `metadata.ProviderTVDb`
- `Name()` → `"TheTVDB"`
- `Priority()` → `80`
- `SupportsMovies()` → `false` (TMDb used for movies)
- `SupportsTVShows()` → `true`
- `SupportsPeople()` → `true`
- `SupportsLanguage(lang)` → `true` (all non-empty languages)

### TV Show Methods

| Method | Returns | Notes |
|--------|---------|-------|
| SearchTVShow | `[]TVShowSearchResult` | Filters by `type=="series"` |
| GetTVShow | `*TVShowMetadata` | Uses extended endpoint, includes translations |
| GetTVShowCredits | `*Credits` | From characters array |
| GetTVShowImages | `*Images` | From artworks endpoint |
| GetTVShowContentRatings | `[]ContentRating` | From extended response |
| GetTVShowTranslations | `[]Translation` | From overviews + name translations |
| GetTVShowExternalIDs | `*ExternalIDs` | Maps remote IDs (IMDb, TMDb) |
| GetSeason | `*SeasonMetadata` | Finds season ID from series, then extended |
| GetSeasonCredits | `*Credits` | Falls back to show credits |
| GetSeasonImages | `*Images` | Returns empty (no dedicated endpoint) |
| GetEpisode | `*EpisodeMetadata` | Fetches via series episodes endpoint |
| GetEpisodeCredits | `*Credits` | From episode characters array |
| GetEpisodeImages | `*Images` | Single image from episode |

### Person Methods

| Method | Returns |
|--------|---------|
| SearchPerson | `[]PersonSearchResult` |
| GetPerson | `*PersonMetadata` |
| GetPersonCredits | `*PersonCredits` |
| GetPersonImages | `*Images` |
| GetPersonExternalIDs | `*ExternalIDs` |

## Client

HTTP client with JWT authentication, rate limiting, and L1Cache (otter) caching:

```go
type Client struct {
    httpClient  *req.Client
    apiKey      string
    pin         string
    rateLimiter *rate.Limiter
    cache       *cache.L1Cache[string, any]
    cacheTTL    time.Duration

    // JWT token management
    token       string
    tokenExpiry time.Time
    tokenMutex  sync.RWMutex
}
```

### JWT Authentication

TVDb API v4 requires JWT tokens obtained via `/login`. The client handles this transparently:

1. `authenticate()` called before each request
2. Token cached with 24h expiry (actual validity: 30 days)
3. Refreshed when within `TokenRefreshBuffer` (1h) of expiry
4. Thread-safe via `sync.RWMutex`
5. `Logout()` invalidates current token

### Constants

```go
BaseURL            = "https://api4.thetvdb.com/v4"
DefaultRateLimit   = rate.Limit(5.0)
DefaultBurst       = 10
DefaultCacheTTL    = 24 * time.Hour
SearchCacheTTL     = 15 * time.Minute
TokenRefreshBuffer = 1 * time.Hour
```

### Client Methods

| Method | TVDb Endpoint | Description |
|--------|--------------|-------------|
| Search | `GET /search` | Search by query + type filter |
| GetSeries | `GET /series/{id}` | Basic series info |
| GetSeriesExtended | `GET /series/{id}/extended` | Full series with seasons, artworks, characters |
| GetSeriesArtworks | `GET /series/{id}/artworks` | Series artwork by type/language |
| GetSeriesEpisodes | `GET /series/{id}/episodes/{type}` | Episodes by season type + season filter |
| GetSeason | `GET /seasons/{id}` | Basic season info |
| GetSeasonExtended | `GET /seasons/{id}/extended` | Full season with episodes |
| GetEpisode | `GET /episodes/{id}` | Basic episode info |
| GetEpisodeExtended | `GET /episodes/{id}/extended` | Full episode with characters |
| GetMovie | `GET /movies/{id}` | Basic movie info |
| GetMovieExtended | `GET /movies/{id}/extended` | Full movie info |
| GetPerson | `GET /people/{id}` | Basic person info |
| GetPersonExtended | `GET /people/{id}/extended` | Full person with characters, remote IDs |
| ClearCache | - | Flush L1Cache |
| Logout | - | Invalidate JWT token |

Generic response wrappers: `BaseResponse[T]` for single items, `ListResponse[T]` for lists.

## Mapper

15 mapping functions converting TVDb response types to `metadata.*` domain types:

| Function | Converts |
|----------|----------|
| mapTVSearchResult | `SearchResult` → `TVShowSearchResult` |
| mapTVShowMetadata | `SeriesResponse` → `*TVShowMetadata` |
| mapSeasonMetadata | `SeasonResponse` → `*SeasonMetadata` |
| mapEpisodeMetadata | `EpisodeResponse` → `*EpisodeMetadata` |
| mapPersonSearchResult | `SearchResult` → `PersonSearchResult` |
| mapPersonMetadata | `PersonResponse` → `*PersonMetadata` |
| mapPersonCredits | `PersonResponse` → `*PersonCredits` |
| mapCharactersToCredits | `[]CharacterResponse` → `*Credits` |
| mapArtworksToImages | `[]ArtworkResponse` → `*Images` |
| mapContentRatings | `[]ContentRatingResponse` → `[]ContentRating` |
| mapOverviewsToTranslations | `map[string]string` → `[]Translation` |
| mapRemoteIDsToExternalIDs | `[]RemoteIDResponse` → `*ExternalIDs` |

### Artwork Type Mapping

| TVDb Artwork Type | Maps To |
|-------------------|---------|
| Poster | `Images.Posters` |
| Background | `Images.Backdrops` |
| Banner | `Images.Logos` |
| ClearLogo, ClearArt | `Images.Logos` |

### Character Type Mapping

| TVDb Character Type | Maps To |
|--------------------|---------|
| Actor | `CastMember` |
| Director | `CrewMember` (Directing) |
| Writer | `CrewMember` (Writing) |
| Producer | `CrewMember` (Production) |

### Remote ID Mapping

TVDb remote IDs are cross-referenced to populate `ExternalIDs` with IMDb and TMDb identifiers.

## Configuration

```go
type Config struct {
    APIKey     string        // TVDb API key
    PIN        string        // Optional subscriber PIN
    RateLimit  rate.Limit    // Requests per second (default: 5.0)
    Burst      int           // Burst capacity (default: 10)
    CacheTTL   time.Duration // Cache duration (default: 24h)
    Timeout    time.Duration // HTTP timeout (default: 30s)
    ProxyURL   string        // Optional HTTP proxy
    RetryCount int           // Retry count (default: 3)
}
```

## Dependencies

- `github.com/imroc/req/v3` - HTTP client
- `golang.org/x/time/rate` - Rate limiting
- `internal/service/metadata` - Domain types and interfaces

## Related Documentation

- [TMDB.md](TMDB.md) - TMDb provider (primary, movies + TV shows + collections)
- [../../../architecture/METADATA_SYSTEM.md](../../../architecture/METADATA_SYSTEM.md) - Metadata system architecture
- [../../../features/video/TVSHOW_MODULE.md](../../../features/video/TVSHOW_MODULE.md) - TV show content module

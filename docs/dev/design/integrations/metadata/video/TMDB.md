# TMDb Integration

<!-- DESIGN: integrations/metadata/video -->

**Package**: `internal/service/metadata/providers/tmdb`
**API**: TMDb API v3 (`https://api.themoviedb.org/3`)

> Primary external metadata provider for movies and TV shows, with person, image, and collection support

---

## Module Structure

```
internal/service/metadata/providers/tmdb/
├── client.go    # HTTP client (req, rate limiter, sync.Map cache)
├── types.go     # TMDb API v3 response types (500+ lines)
├── provider.go  # Provider interface implementation
└── mapping.go   # TMDb responses → metadata domain types
```

## Provider

Implements 6 metadata interfaces:

```go
var (
    _ metadata.Provider           = (*Provider)(nil)
    _ metadata.MovieProvider      = (*Provider)(nil)
    _ metadata.TVShowProvider     = (*Provider)(nil)
    _ metadata.PersonProvider     = (*Provider)(nil)
    _ metadata.ImageProvider      = (*Provider)(nil)
    _ metadata.CollectionProvider = (*Provider)(nil)
)

type Provider struct {
    client   *Client
    priority int  // 100 (primary)
}
```

- `ID()` → `metadata.ProviderTMDb`
- `Name()` → `"The Movie Database"`
- `Priority()` → `100`
- `SupportsMovies()` → `true`
- `SupportsTVShows()` → `true`
- `SupportsPeople()` → `true`
- `SupportsLanguage(lang)` → `true` (all languages)

### Movie Methods

| Method | Returns |
|--------|---------|
| SearchMovie | `[]MovieSearchResult` |
| GetMovie | `*MovieMetadata` |
| GetMovieCredits | `*Credits` |
| GetMovieImages | `*Images` |
| GetMovieReleaseDates | `[]ReleaseDate` |
| GetMovieTranslations | `[]Translation` |
| GetMovieExternalIDs | `*ExternalIDs` |
| GetSimilarMovies | `[]MovieSearchResult, int` |
| GetMovieRecommendations | `[]MovieSearchResult, int` |

### TV Show Methods

| Method | Returns |
|--------|---------|
| SearchTVShow | `[]TVShowSearchResult` |
| GetTVShow | `*TVShowMetadata` |
| GetTVShowCredits | `*Credits` |
| GetTVShowImages | `*Images` |
| GetTVShowContentRatings | `[]ContentRating` |
| GetTVShowTranslations | `[]Translation` |
| GetTVShowExternalIDs | `*ExternalIDs` |
| GetSeason | `*SeasonMetadata` |
| GetSeasonCredits | `*Credits` |
| GetSeasonImages | `*Images` |
| GetEpisode | `*EpisodeMetadata` |
| GetEpisodeCredits | `*Credits` |
| GetEpisodeImages | `*Images` |

### Person Methods

| Method | Returns |
|--------|---------|
| SearchPerson | `[]PersonSearchResult` |
| GetPerson | `*PersonMetadata` |
| GetPersonCredits | `*PersonCredits` |
| GetPersonImages | `*Images` |
| GetPersonExternalIDs | `*ExternalIDs` |

### Image Methods

| Method | Purpose |
|--------|---------|
| GetImageURL | Build full URL from path + size |
| GetImageBaseURL | Returns `https://image.tmdb.org/t/p` |
| DownloadImage | Download image bytes |

### Collection Methods

| Method | Returns |
|--------|---------|
| GetCollection | `*CollectionMetadata` |
| GetCollectionImages | `*Images` |

## Client

HTTP client wrapping `imroc/req` with rate limiting and sync.Map caching:

```go
type Client struct {
    httpClient  *req.Client
    apiKey      string
    accessToken string
    rateLimiter *rate.Limiter
    cache       sync.Map
    cacheTTL    time.Duration
}
```

### Constants

```go
BaseURL         = "https://api.themoviedb.org/3"
ImageBaseURL    = "https://image.tmdb.org/t/p"
DefaultRateLimit = rate.Limit(4.0)  // 40 requests per 10 seconds
DefaultBurst     = 10
DefaultCacheTTL  = 24 * time.Hour
SearchCacheTTL   = 15 * time.Minute
```

### Client Methods (~28)

Maps 1:1 to TMDb API v3 endpoints. Each method checks cache first, waits on rate limiter, then makes the HTTP request. Key methods:

- `SearchMovie`, `SearchTV`, `SearchPerson`
- `GetMovie`, `GetTV`, `GetSeason`, `GetEpisode`, `GetPerson`
- `GetMovieCredits`, `GetMovieImages`, `GetMovieReleaseDates`, `GetMovieTranslations`, `GetMovieExternalIDs`
- `GetSimilarMovies`, `GetMovieRecommendations`
- `GetTVCredits`, `GetTVImages`, `GetTVContentRatings`, `GetTVTranslations`, `GetTVExternalIDs`
- `GetSeasonCredits`, `GetSeasonImages`
- `GetEpisodeCredits`, `GetEpisodeImages`
- `GetPersonCredits`, `GetPersonImages`, `GetPersonExternalIDs`
- `GetCollection`, `GetImageURL`, `DownloadImage`, `ClearCache`

Uses TMDb's `append_to_response` parameter to batch sub-requests.

## Mapper

20 mapping functions converting TMDb response types to `metadata.*` domain types:

| Function | Converts |
|----------|----------|
| mapMovieSearchResults | `SearchResultsResponse` → `[]MovieSearchResult` |
| mapMovieMetadata | `MovieResponse` → `*MovieMetadata` |
| mapTVSearchResult | `TVSearchResponse` → `TVShowSearchResult` |
| mapTVShowMetadata | `TVResponse` → `*TVShowMetadata` |
| mapSeasonMetadata | `SeasonResponse` → `*SeasonMetadata` |
| mapEpisodeMetadata | `EpisodeResponse` → `*EpisodeMetadata` |
| mapPersonSearchResult | `PersonSearchResponse` → `PersonSearchResult` |
| mapPersonMetadata | `PersonResponse` → `*PersonMetadata` |
| mapPersonCredits | `PersonCreditsResponse` → `*PersonCredits` |
| mapCredits | `CreditsResponse` → `*Credits` |
| mapImages | `ImagesResponse` → `*Images` |
| mapReleaseDates | `ReleaseDatesWrapper` → `[]ReleaseDate` |
| mapContentRatings | `ContentRatingsWrapper` → `[]ContentRating` |
| mapTranslations | `TranslationsWrapper` → `[]Translation` |
| mapExternalIDs | `ExternalIDsResponse` → `*ExternalIDs` |
| mapCollectionMetadata | `CollectionResponse` → `*CollectionMetadata` |
| normalizeLang | ISO 639-1 → TMDb format (en → en-US) |

## Configuration

```go
type Config struct {
    APIKey      string        // TMDb API key (v3)
    AccessToken string        // TMDb access token (v4, alternative auth)
    RateLimit   rate.Limit    // Requests per second (default: 4.0)
    Burst       int           // Burst capacity (default: 10)
    CacheTTL    time.Duration // Cache duration (default: 24h)
    Timeout     time.Duration // HTTP timeout (default: 30s)
    ProxyURL    string        // Optional HTTP proxy
    RetryCount  int           // Retry count (default: 3)
}
```

## Dependencies

- `github.com/imroc/req/v3` - HTTP client
- `golang.org/x/time/rate` - Rate limiting
- `internal/service/metadata` - Domain types and interfaces

## Related Documentation

- [THETVDB.md](THETVDB.md) - TVDb provider (secondary, TV shows + people)
- [../../../architecture/METADATA_SYSTEM.md](../../../architecture/METADATA_SYSTEM.md) - Metadata system architecture
- [../../../features/video/MOVIE_MODULE.md](../../../features/video/MOVIE_MODULE.md) - Movie content module
- [../../../features/video/TVSHOW_MODULE.md](../../../features/video/TVSHOW_MODULE.md) - TV show content module

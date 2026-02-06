# Metadata System

**Last Updated**: 2026-02-06

Multi-provider metadata aggregation with priority-based fallback, per-provider caching, and adapter pattern for content modules.

---

## Overview

```
Content Modules                    Metadata Service                    External APIs
┌──────────┐                      ┌──────────────┐                   ┌──────────┐
│  movie   │──MetadataProvider──→│              │──MovieProvider───→│  TMDb    │
│  Service  │                      │   metadata   │                   │ (pri 100)│
└──────────┘                      │   .Service   │                   └──────────┘
                                   │              │                   ┌──────────┐
┌──────────┐                      │  (interface, │──TVShowProvider──→│  TVDb    │
│  tvshow  │──MetadataProvider──→│   31 methods)│                   │ (pri 80) │
│  Service  │                      │              │                   └──────────┘
└──────────┘                      └──────────────┘
       ↑                                 ↑
   Adapters                         Provider
   (type mapping)                   registration
```

Content modules never call external APIs directly. They receive a `MetadataProvider` interface injected by fx. Behind that interface sits an adapter that delegates to the shared `metadata.Service`, which dispatches to registered providers by priority.

---

## Service Interface

The `metadata.Service` interface (`internal/service/metadata/service.go`) exposes 31 methods in 7 categories:

| Category | Methods | Notes |
|----------|---------|-------|
| Movie | `SearchMovie`, `GetMovieMetadata`, `GetMovieCredits`, `GetMovieImages`, `GetMovieReleaseDates`, `GetMovieExternalIDs`, `GetSimilarMovies`, `GetMovieRecommendations` | 8 methods |
| TV Show | `SearchTVShow`, `GetTVShowMetadata`, `GetTVShowCredits`, `GetTVShowImages`, `GetTVShowContentRatings`, `GetTVShowExternalIDs`, `GetSeasonMetadata`, `GetEpisodeMetadata` | 8 methods |
| Person | `SearchPerson`, `GetPersonMetadata`, `GetPersonCredits`, `GetPersonImages` | 4 methods |
| Collection | `GetCollectionMetadata` | 1 method |
| Image | `GetImageURL` | URL construction |
| Refresh | `RefreshMovie`, `RefreshTVShow` | Enqueue River jobs |
| Management | `ClearCache`, `RegisterProvider`, `GetProviders` | 3 methods |

---

## Provider Architecture

### Base Provider Interface

Every provider implements the base `Provider` interface:

```go
type Provider interface {
    ID() ProviderID              // "tmdb", "tvdb", etc.
    Name() string                // Human-readable name
    Priority() int               // Higher = preferred
    SupportsMovies() bool
    SupportsTVShows() bool
    SupportsPeople() bool
    SupportsLanguage(lang string) bool
    ClearCache()
}
```

### Specialized Interfaces

Providers additionally implement capability-specific interfaces:

| Interface | Methods | Implemented by |
|-----------|---------|---------------|
| `MovieProvider` | `SearchMovie`, `GetMovie`, `GetMovieCredits`, `GetMovieImages`, `GetMovieReleaseDates`, `GetMovieTranslations`, `GetMovieExternalIDs`, `GetSimilarMovies`, `GetMovieRecommendations` | TMDb |
| `TVShowProvider` | `SearchTVShow`, `GetTVShow`, `GetTVShowCredits`, `GetTVShowImages`, `GetTVShowContentRatings`, `GetTVShowTranslations`, `GetTVShowExternalIDs`, `GetSeason`, `GetSeasonCredits`, `GetSeasonImages`, `GetEpisode`, `GetEpisodeCredits`, `GetEpisodeImages` | TMDb, TVDb |
| `PersonProvider` | `SearchPerson`, `GetPerson`, `GetPersonCredits`, `GetPersonImages`, `GetPersonExternalIDs` | TMDb |
| `ImageProvider` | `GetImageURL`, `GetImageBaseURL`, `DownloadImage` | TMDb |
| `CollectionProvider` | `GetCollection`, `GetCollectionImages` | TMDb |

### Current Providers

| Provider | ID | Priority | Capabilities | API Client |
|----------|-----|---------|-------------|------------|
| TMDb | `tmdb` | 100 | Movie, TV, Person, Image, Collection | resty + rate limiter |
| TVDb | `tvdb` | 80 | TV (series, seasons, episodes) | resty + rate limiter |

Reserved provider IDs (defined but not yet implemented): `fanarttv`, `omdb`.

### Provider Registration

Providers are registered during fx module initialization. The service sorts them by priority (highest first) and categorizes by capability using type assertions:

```go
svc.RegisterProvider(tmdbProvider)  // auto-sorted by Priority()
svc.RegisterProvider(tvdbProvider)  // categorized by interface
```

---

## Priority and Fallback

### Request Flow

1. Service receives a request (e.g., `SearchMovie`)
2. Selects providers for that capability (e.g., `movieProviders`), sorted by priority
3. Calls the highest-priority provider first
4. If it fails and `EnableProviderFallback` is true, tries the next provider
5. Returns the first successful result or the last error

### Multi-Language Fetching

For metadata methods that accept `languages []string`:

1. Iterates over each requested language
2. First language response becomes the base result
3. Subsequent languages are merged as translations into a `Translations` map
4. Each translation contains localized fields (title, overview, tagline)

```go
// Example: GetMovieMetadata with languages ["en", "de", "fr"]
// → First call: base metadata in English
// → Second call: German title/overview/tagline added to Translations["de"]
// → Third call: French title/overview/tagline added to Translations["fr"]
```

Default languages when none specified: `["en"]` (from `ServiceConfig.DefaultLanguages`).

---

## Caching

Each provider maintains its own in-memory cache using `sync.Map`:

```go
type Client struct {
    cache    sync.Map      // key → *CacheEntry
    cacheTTL time.Duration // default: 24h for metadata, 15m for search
}

type CacheEntry struct {
    Data      any
    ExpiresAt time.Time
}
```

### Cache Behavior

- **Per-provider**: TMDb and TVDb each have independent caches
- **TTL-based**: Entries expire after a configurable duration
- **Lazy expiration**: Expired entries are checked on read, not proactively evicted
- **Key format**: Provider-specific (e.g., `movie:{id}:{lang}`, `search:{query}`)

### Cache Clearing

`ClearCache()` cascades through the system:

```
metadata.Service.ClearCache()
  → for each provider:
      provider.ClearCache()
        → client.ClearCache()
          → sync.Map range + delete
```

The movie adapter also delegates `ClearCache()` to the shared service, and triggers it on force-refresh operations.

---

## Adapters

Adapters bridge the shared metadata types to content module domain types. They are created during fx module initialization.

### Movie Adapter

`internal/service/metadata/adapters/movie/adapter.go`

Implements `movie.MetadataProvider`:
- `SearchMovies` - searches via service, maps `MovieSearchResult` → `movie.Movie`
- `EnrichMovie` - fetches metadata + release dates, maps to movie domain type
- `GetMovieCredits` - maps `Credits` → `[]movie.MovieCredit` (cast + crew)
- `GetMovieGenres` - extracts genres from metadata
- `GetMovieByTMDbID` - creates a new movie from TMDb metadata
- `GetMovieImages` - delegates directly
- `GetImageURL` - delegates directly
- `ClearCache` - delegates to service

Type mapping includes:
- `float64` → `decimal.Decimal` (vote average, popularity)
- Release dates → age ratings map (country → system → rating, e.g., US → MPAA → PG-13)
- Translations → i18n maps (`TitlesI18n`, `TaglinesI18n`, `OverviewsI18n`)

### TV Show Adapter

`internal/service/metadata/adapters/tvshow/adapter.go`

Implements `tvshow.MetadataProvider` with the same delegation pattern for series, seasons, and episodes.

---

## Errors

### Sentinel Errors

```go
var (
    ErrNotFound             = errors.New("metadata: not found")
    ErrProviderUnavailable  = errors.New("metadata: provider unavailable")
    ErrRateLimited          = errors.New("metadata: rate limited")
    ErrUnauthorized         = errors.New("metadata: unauthorized")
    ErrInvalidID            = errors.New("metadata: invalid id")
    ErrNoProviders          = errors.New("metadata: no providers configured")
    ErrUnsupported          = errors.New("metadata: operation not supported")
)
```

### Structured Error Types

| Type | Purpose |
|------|---------|
| `ProviderError` | Wraps errors from a specific provider with ID, status code, message |
| `AggregateError` | Collects errors from multiple providers/languages, exposes `First()` and `HasNotFound()` |

---

## Async Refresh

The service delegates refresh operations to River jobs via a `JobQueue` interface:

```go
type JobQueue interface {
    EnqueueRefreshMovie(ctx context.Context, movieID uuid.UUID, force bool, languages []string) error
    EnqueueRefreshTVShow(ctx context.Context, seriesID uuid.UUID, force bool, languages []string) error
}
```

Related River workers (defined in content module `jobs/` packages):

| Worker | Args Type | Purpose |
|--------|-----------|---------|
| MetadataRefreshMovie | `RefreshMovieArgs` | Refresh single movie metadata |
| MetadataRefreshTVShow | `RefreshTVShowArgs` | Refresh TV show metadata |
| MetadataRefreshSeason | `RefreshSeasonArgs` | Refresh season metadata |
| MetadataRefreshEpisode | `RefreshEpisodeArgs` | Refresh episode metadata |
| SeriesRefresh | - | Full series refresh (all seasons + episodes) |

---

## fx Wiring

`internal/service/metadata/metadatafx/module.go` provides:

| Output | Type | Description |
|--------|------|-------------|
| `Service` | `metadata.Service` | Shared metadata service |
| `MovieMetadataAdapter` | `movie.MetadataProvider` | Injected into movie module |
| `TVShowMetadataAdapter` | `tvshow.MetadataProvider` | Injected into tvshow module |
| `TMDbProvider` | `*tmdb.Provider` | Optional, only if API key configured |
| `TVDbProvider` | `*tvdb.Provider` | Optional, only if API key configured |

Configuration is read from the app config. If no API key is provided for a provider, it is not registered and the service operates without it.

---

## Configuration

```go
type ServiceConfig struct {
    DefaultLanguages       []string  // default: ["en"]
    EnableProviderFallback bool      // default: true
    EnableEnrichment       bool      // default: false
}
```

Provider configs:

```yaml
metadata:
  tmdb:
    api_key: "..."       # Required for TMDb provider
    proxy_url: "..."     # Optional proxy
  tvdb:
    api_key: "..."       # Required for TVDb provider
    pin: "..."           # TVDb subscriber PIN
```

---

## Related Documentation

- [Architecture](ARCHITECTURE.md) - System structure and layers
- [Design Principles](DESIGN_PRINCIPLES.md) - Adapter pattern, error handling
- [Dragonfly](../integrations/infrastructure/DRAGONFLY.md) - Infrastructure cache (separate from metadata provider cache)

# Metadata Providers Instructions

> Source: https://developer.themoviedb.org/docs, https://musicbrainz.org/doc/MusicBrainz_API

Apply to: `**/internal/service/metadata/**/*.go`, `**/internal/content/**/provider*.go`

## Overview

Revenge uses external metadata providers for content enrichment. Each content module has its own provider implementations.

**Provider Priority Chain:**

1. Primary provider (module-specific)
2. Fallback providers
3. Cache (Dragonfly)

## Rate Limiting

All providers enforce strict rate limits:

| Provider    | Rate Limit           | User-Agent Required |
| ----------- | -------------------- | ------------------- |
| TMDb        | ~50 req/s (varies)   | Yes                 |
| MusicBrainz | **1 req/s** (STRICT) | **Mandatory**       |
| TheTVDB     | 100 req/10s          | Yes (token)         |
| OpenLibrary | 1 req/s              | Yes                 |
| ComicVine   | 200 req/hr           | Yes (API key)       |
| StashDB     | 1 req/s              | Yes (API key)       |

### User-Agent Format

```
Revenge/1.0.0 (https://github.com/lusoris/revenge; contact@example.com)
```

**CRITICAL for MusicBrainz**: They will block requests without a meaningful User-Agent!

## TMDb (The Movie Database)

### API Endpoint

```
https://api.themoviedb.org/3/
```

### Authentication

```go
// Bearer token in Authorization header
req.Header.Set("Authorization", "Bearer "+apiKey)
req.Header.Set("Accept", "application/json")
```

### Key Endpoints

| Endpoint                          | Method | Description                      |
| --------------------------------- | ------ | -------------------------------- |
| `/search/movie`                   | GET    | Search movies                    |
| `/search/tv`                      | GET    | Search TV shows                  |
| `/movie/{id}`                     | GET    | Movie details                    |
| `/movie/{id}/credits`             | GET    | Movie cast/crew                  |
| `/movie/{id}/images`              | GET    | Movie artwork                    |
| `/movie/{id}/videos`              | GET    | Movie trailers                   |
| `/tv/{id}`                        | GET    | TV series details                |
| `/tv/{id}/season/{n}`             | GET    | Season details                   |
| `/tv/{id}/season/{n}/episode/{e}` | GET    | Episode details                  |
| `/find/{id}`                      | GET    | Find by external ID (IMDb, TVDb) |
| `/configuration`                  | GET    | Image base URLs                  |

### Append To Response

Reduce API calls with `append_to_response`:

```go
// Get movie with credits, images, and videos in one call
url := fmt.Sprintf(
    "https://api.themoviedb.org/3/movie/%d?append_to_response=credits,images,videos",
    movieID,
)
```

### Image Configuration

```go
// Get base URL from /configuration
// baseURL typically: https://image.tmdb.org/t/p/
// Sizes: w92, w154, w185, w342, w500, w780, original

posterURL := fmt.Sprintf("%s%s%s", baseURL, "w500", posterPath)
```

### Example: Search Movie

```go
type TMDbClient struct {
    client  *http.Client
    apiKey  string
    baseURL string
}

func (c *TMDbClient) SearchMovie(ctx context.Context, query string, year int) ([]TMDbMovie, error) {
    u := fmt.Sprintf("%s/search/movie?query=%s&year=%d",
        c.baseURL, url.QueryEscape(query), year)

    req, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
    req.Header.Set("Authorization", "Bearer "+c.apiKey)
    req.Header.Set("Accept", "application/json")

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode == 429 {
        return nil, ErrRateLimited
    }

    var result TMDbSearchResult
    json.NewDecoder(resp.Body).Decode(&result)
    return result.Results, nil
}
```

## MusicBrainz

### API Endpoint

```
https://musicbrainz.org/ws/2/
```

### Request Format

```go
// JSON format
req.Header.Set("Accept", "application/json")
// OR: add ?fmt=json to URL

// MANDATORY User-Agent
req.Header.Set("User-Agent", "Revenge/1.0.0 (https://github.com/lusoris/revenge)")
```

### Core Entities

| Entity          | Description                  |
| --------------- | ---------------------------- |
| `artist`        | Musicians, bands, orchestras |
| `release-group` | Album concept (all editions) |
| `release`       | Specific edition of album    |
| `recording`     | Unique audio recording       |
| `work`          | Musical composition          |

### Request Types

```
# Lookup (by MBID)
/ws/2/{entity}/{mbid}?inc={includes}&fmt=json

# Browse (linked entities)
/ws/2/{entity}?{linked_entity}={mbid}&limit={n}&offset={n}&fmt=json

# Search
/ws/2/{entity}?query={lucene_query}&limit={n}&offset={n}&fmt=json
```

### Include Parameters

```
# Artist lookup with releases and recordings
/ws/2/artist/{mbid}?inc=releases+recordings&fmt=json

# Release with media, recordings, and artist credits
/ws/2/release/{mbid}?inc=media+recordings+artist-credits&fmt=json

# Recording with artist credits and ISRC
/ws/2/recording/{mbid}?inc=artist-credits+isrcs&fmt=json
```

### Common Includes

| Entity        | Available Includes                                                 |
| ------------- | ------------------------------------------------------------------ |
| artist        | recordings, releases, release-groups, works, aliases, tags, genres |
| release       | media, recordings, artist-credits, labels, release-groups          |
| recording     | releases, artist-credits, isrcs, tags                              |
| release-group | releases, artist-credits, tags, genres                             |

### Example: Artist Lookup

```go
type MusicBrainzClient struct {
    client    *http.Client
    baseURL   string
    userAgent string
    limiter   *rate.Limiter  // 1 req/sec
}

func (c *MusicBrainzClient) GetArtist(ctx context.Context, mbid string) (*MBArtist, error) {
    // Wait for rate limit
    if err := c.limiter.Wait(ctx); err != nil {
        return nil, err
    }

    u := fmt.Sprintf("%s/artist/%s?inc=aliases+genres+tags&fmt=json",
        c.baseURL, mbid)

    req, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
    req.Header.Set("User-Agent", c.userAgent)
    req.Header.Set("Accept", "application/json")

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode == 503 {
        return nil, ErrRateLimited  // Back off!
    }

    var artist MBArtist
    json.NewDecoder(resp.Body).Decode(&artist)
    return &artist, nil
}
```

### Browse Pagination

```go
// Browse releases by artist with pagination
func (c *MusicBrainzClient) GetArtistReleases(ctx context.Context, artistMBID string) ([]MBRelease, error) {
    var allReleases []MBRelease
    offset := 0
    limit := 100  // max is 100

    for {
        if err := c.limiter.Wait(ctx); err != nil {
            return nil, err
        }

        u := fmt.Sprintf(
            "%s/release?artist=%s&inc=media+artist-credits&limit=%d&offset=%d&fmt=json",
            c.baseURL, artistMBID, limit, offset)

        // ... make request ...

        allReleases = append(allReleases, page.Releases...)

        if offset + len(page.Releases) >= page.TotalCount {
            break
        }
        offset += limit
    }

    return allReleases, nil
}
```

## Provider Interface

```go
type MetadataProvider interface {
    // Search for content by name
    Search(ctx context.Context, query string, opts SearchOpts) ([]SearchResult, error)

    // Get full details by provider ID
    GetDetails(ctx context.Context, providerID string) (*ContentMetadata, error)

    // Get images for content
    GetImages(ctx context.Context, providerID string) ([]Image, error)

    // Provider name for attribution
    Name() string

    // Provider priority (lower = higher priority)
    Priority() int
}

type SearchOpts struct {
    Year     int
    Language string
    Limit    int
}

type ContentMetadata struct {
    Title       string
    OrigTitle   string
    Overview    string
    Year        int
    Genres      []string
    Cast        []Person
    Crew        []Person
    ExternalIDs map[string]string
}
```

## Provider by Module

| Module      | Primary     | Fallback                |
| ----------- | ----------- | ----------------------- |
| movie       | TMDb        | OMDb, Fanart.tv         |
| tvshow      | TMDb        | TheTVDB, Fanart.tv      |
| music       | MusicBrainz | Last.fm, Fanart.tv      |
| audiobook   | Audible     | Google Books            |
| book        | OpenLibrary | Google Books            |
| podcast     | iTunes      | Podcast Index           |
| comics      | ComicVine   | Marvel, Grand Comics DB |
| adult_movie | StashDB     | TPDB                    |
| adult_show  | StashDB     | TPDB                    |

## Caching Strategy

```go
// Cache metadata responses
func (s *MetadataService) GetMovie(ctx context.Context, tmdbID int) (*Movie, error) {
    // Check cache first
    key := fmt.Sprintf("meta:tmdb:movie:%d", tmdbID)
    if cached, err := s.cache.Get(ctx, key); err == nil {
        var movie Movie
        json.Unmarshal(cached, &movie)
        return &movie, nil
    }

    // Fetch from provider
    movie, err := s.tmdb.GetMovie(ctx, tmdbID)
    if err != nil {
        return nil, err
    }

    // Cache for 1 hour
    data, _ := json.Marshal(movie)
    s.cache.Set(ctx, key, data, 1*time.Hour)

    return movie, nil
}
```

## Error Handling

```go
var (
    ErrProviderNotFound   = errors.New("content not found in provider")
    ErrRateLimited        = errors.New("provider rate limit exceeded")
    ErrProviderUnavailable = errors.New("provider temporarily unavailable")
    ErrInvalidResponse    = errors.New("invalid response from provider")
)

// Retry with exponential backoff for transient errors
func (s *MetadataService) fetchWithRetry(ctx context.Context, fn func() error) error {
    backoff := 1 * time.Second
    for attempt := 0; attempt < 3; attempt++ {
        err := fn()
        if err == nil {
            return nil
        }
        if errors.Is(err, ErrProviderNotFound) {
            return err  // Don't retry not found
        }
        if errors.Is(err, ErrRateLimited) {
            select {
            case <-time.After(backoff):
                backoff *= 2
            case <-ctx.Done():
                return ctx.Err()
            }
            continue
        }
        return err
    }
    return ErrProviderUnavailable
}
```

## River Jobs for Metadata

```go
type FetchMetadataArgs struct {
    ContentType string    `json:"content_type"`  // "movie", "tvshow", etc.
    ContentID   uuid.UUID `json:"content_id"`
    ProviderID  string    `json:"provider_id"`
    Provider    string    `json:"provider"`      // "tmdb", "musicbrainz", etc.
}

func (FetchMetadataArgs) Kind() string { return "metadata.fetch" }

type FetchMetadataWorker struct {
    river.WorkerDefaults[FetchMetadataArgs]
    providers map[string]MetadataProvider
    repo      ContentRepository
}

func (w *FetchMetadataWorker) Work(ctx context.Context, job *river.Job[FetchMetadataArgs]) error {
    provider, ok := w.providers[job.Args.Provider]
    if !ok {
        return fmt.Errorf("unknown provider: %s", job.Args.Provider)
    }

    metadata, err := provider.GetDetails(ctx, job.Args.ProviderID)
    if err != nil {
        return err
    }

    return w.repo.UpdateMetadata(ctx, job.Args.ContentID, metadata)
}
```

## DO's and DON'Ts

### DO

- ✅ Implement rate limiting per provider
- ✅ Set meaningful User-Agent (especially for MusicBrainz)
- ✅ Cache responses appropriately
- ✅ Use `append_to_response` for TMDb
- ✅ Handle 429/503 responses gracefully
- ✅ Fetch metadata via River jobs (async)
- ✅ Store provider IDs for re-fetching

### DON'T

- ❌ Make requests without rate limiting
- ❌ Ignore User-Agent requirements
- ❌ Fetch metadata synchronously in HTTP handlers
- ❌ Forget to handle provider downtime
- ❌ Make multiple requests when one will do
- ❌ Store provider responses without transformation
---

## Related

- [INDEX.instructions.md](INDEX.instructions.md) - Main instruction index with all cross-references
- [ARCHITECTURE_V2.md](../../docs/dev/design/architecture/ARCHITECTURE_V2.md) - System architecture
- [BEST_PRACTICES.md](../../docs/dev/design/operations/BEST_PRACTICES.md) - Best practices

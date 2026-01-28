---
applyTo: "**/internal/service/**/*.go,**/internal/content/**/*.go"
---

# API Response Caching - sturdyc

> Request coalescing and caching for external API calls (90% call reduction)

## Overview

Use `sturdyc` for caching external API responses (TMDb, MusicBrainz, etc.). It provides:

- **Request coalescing**: Multiple concurrent requests for same key → single backend call
- **Stale-while-revalidate**: Serve stale data while refreshing in background
- **Batch fetching**: Aggregate multiple keys into single backend call
- **Passthrough recording**: Track non-cacheable requests

**Package**: `github.com/viccon/sturdyc`

## When to Use

| Use Case                 | Use sturdyc?              |
| ------------------------ | ------------------------- |
| TMDb movie metadata      | ✅ Yes - rate limited API |
| MusicBrainz lookups      | ✅ Yes - 1 req/sec limit! |
| User session data        | ❌ No - use rueidis       |
| Real-time playback state | ❌ No - too dynamic       |
| Servarr API calls        | ✅ Yes - reduce load      |

## Installation

```bash
go get github.com/viccon/sturdyc
```

## Basic Usage

### Create Client

```go
import "github.com/viccon/sturdyc"

client := sturdyc.New[any](
    10_000,                    // Max entries
    10,                        // Shards (power of 2)
    5*time.Minute,             // TTL
    10,                        // Eviction percentage
    sturdyc.WithStampedeProtection(
        1000,                  // Max refresh goroutines
        2*time.Minute,         // Min refresh delay
        4*time.Minute,         // Max refresh delay
        3*time.Minute,         // Retry base delay
        true,                  // Storefront mode
    ),
)
```

### Single Key Fetch

```go
// Fetch with caching and request coalescing
movie, err := sturdyc.GetOrFetch(ctx, client, fmt.Sprintf("movie:%d", tmdbID),
    func(ctx context.Context) (*Movie, error) {
        // This only runs once even with 100 concurrent requests
        return tmdbClient.GetMovie(ctx, tmdbID)
    },
)
```

### Batch Fetch

```go
// Fetch multiple items efficiently
ids := []int{123, 456, 789}
movies, err := sturdyc.GetOrFetchBatch(ctx, client, ids,
    func(id int) string {
        return fmt.Sprintf("movie:%d", id)
    },
    func(ctx context.Context, missingIDs []int) (map[int]*Movie, error) {
        // Fetch only missing items from backend
        return tmdbClient.GetMovies(ctx, missingIDs)
    },
)
```

## Metadata Provider Pattern

```go
package tmdb

import (
    "context"
    "time"
    "github.com/viccon/sturdyc"
)

type Client struct {
    httpClient *resty.Client
    cache      *sturdyc.Client[any]
}

func NewClient(httpClient *resty.Client) *Client {
    cache := sturdyc.New[any](
        50_000,               // Cache up to 50k items
        16,                   // 16 shards
        1*time.Hour,          // 1 hour TTL for metadata
        5,                    // Evict 5% when full
        sturdyc.WithStampedeProtection(500, 30*time.Minute, 45*time.Minute, 5*time.Minute, true),
    )

    return &Client{
        httpClient: httpClient,
        cache:      cache,
    }
}

func (c *Client) GetMovie(ctx context.Context, tmdbID int) (*Movie, error) {
    return sturdyc.GetOrFetch(ctx, c.cache, fmt.Sprintf("tmdb:movie:%d", tmdbID),
        func(ctx context.Context) (*Movie, error) {
            var movie Movie
            resp, err := c.httpClient.R().
                SetContext(ctx).
                SetResult(&movie).
                Get(fmt.Sprintf("/movie/%d", tmdbID))
            if err != nil {
                return nil, err
            }
            if resp.IsError() {
                return nil, fmt.Errorf("TMDb error: %d", resp.StatusCode())
            }
            return &movie, nil
        },
    )
}

func (c *Client) GetMovies(ctx context.Context, tmdbIDs []int) (map[int]*Movie, error) {
    return sturdyc.GetOrFetchBatch(ctx, c.cache, tmdbIDs,
        func(id int) string {
            return fmt.Sprintf("tmdb:movie:%d", id)
        },
        func(ctx context.Context, missingIDs []int) (map[int]*Movie, error) {
            // Fetch missing movies from TMDb
            result := make(map[int]*Movie)
            for _, id := range missingIDs {
                movie, err := c.fetchMovie(ctx, id)
                if err != nil {
                    continue // Skip failures
                }
                result[id] = movie
            }
            return result, nil
        },
    )
}
```

## MusicBrainz Rate Limiting

MusicBrainz has strict 1 req/sec limit. sturdyc helps:

```go
type MusicBrainzClient struct {
    httpClient *resty.Client
    cache      *sturdyc.Client[any]
    limiter    *rate.Limiter // 1 req/sec
}

func NewMusicBrainzClient() *MusicBrainzClient {
    return &MusicBrainzClient{
        httpClient: resty.New().
            SetBaseURL("https://musicbrainz.org/ws/2").
            SetHeader("User-Agent", "Revenge/1.0 (contact@example.com)"),
        cache: sturdyc.New[any](
            100_000,          // Large cache - reduce API calls
            16,
            24*time.Hour,     // Long TTL - music metadata rarely changes
            5,
            sturdyc.WithStampedeProtection(100, 12*time.Hour, 18*time.Hour, 1*time.Hour, true),
        ),
        limiter: rate.NewLimiter(rate.Every(time.Second), 1), // 1 req/sec
    }
}

func (c *MusicBrainzClient) GetArtist(ctx context.Context, mbid string) (*Artist, error) {
    return sturdyc.GetOrFetch(ctx, c.cache, fmt.Sprintf("mb:artist:%s", mbid),
        func(ctx context.Context) (*Artist, error) {
            // Wait for rate limiter
            if err := c.limiter.Wait(ctx); err != nil {
                return nil, err
            }
            // Fetch from API
            return c.fetchArtist(ctx, mbid)
        },
    )
}
```

## fx Module Integration

```go
package metadata

import (
    "github.com/viccon/sturdyc"
    "go.uber.org/fx"
)

type CacheConfig struct {
    MaxEntries       int           `koanf:"max_entries"`
    TTL              time.Duration `koanf:"ttl"`
    RefreshMinDelay  time.Duration `koanf:"refresh_min_delay"`
    RefreshMaxDelay  time.Duration `koanf:"refresh_max_delay"`
}

var DefaultCacheConfig = CacheConfig{
    MaxEntries:      50_000,
    TTL:             1 * time.Hour,
    RefreshMinDelay: 30 * time.Minute,
    RefreshMaxDelay: 45 * time.Minute,
}

func NewMetadataCache(cfg CacheConfig) *sturdyc.Client[any] {
    return sturdyc.New[any](
        cfg.MaxEntries,
        16,
        cfg.TTL,
        5,
        sturdyc.WithStampedeProtection(
            500,
            cfg.RefreshMinDelay,
            cfg.RefreshMaxDelay,
            5*time.Minute,
            true,
        ),
    )
}

var Module = fx.Module("metadata-cache",
    fx.Provide(NewMetadataCache),
)
```

## DO's and DON'Ts

### DO

- ✅ Use for external API calls (rate-limited services)
- ✅ Use batch fetching when fetching multiple items
- ✅ Enable stampede protection for high-traffic endpoints
- ✅ Use long TTLs for stable data (metadata)
- ✅ Combine with rate limiting for strict APIs (MusicBrainz)

### DON'T

- ❌ Use for user-specific data (sessions)
- ❌ Use for rapidly changing data (playback state)
- ❌ Forget to handle partial batch failures
- ❌ Set TTL shorter than refresh delay
---

## Related

- [INDEX.instructions.md](INDEX.instructions.md) - Main instruction index with all cross-references
- [ARCHITECTURE_V2.md](../../docs/dev/design/architecture/ARCHITECTURE_V2.md) - System architecture
- [BEST_PRACTICES.md](../../docs/dev/design/operations/BEST_PRACTICES.md) - Best practices

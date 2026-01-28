---
applyTo: "**/internal/service/**/*.go,**/internal/content/**/*.go"
---

# HTTP Client - resty v3

> HTTP client for external API calls (metadata providers, Servarr, etc.)

## Overview

Use `resty` for all external HTTP calls. It provides:

- Fluent API
- Automatic retry with backoff
- Circuit breaker integration
- Request/response middleware
- Automatic JSON marshaling

**Package**: `resty.dev/v3`

## Installation

```bash
go get resty.dev/v3
```

## Basic Usage

### Create Client

```go
import "resty.dev/v3"

client := resty.New().
    SetBaseURL("https://api.themoviedb.org/3").
    SetHeader("Accept", "application/json").
    SetQueryParam("api_key", cfg.APIKey).
    SetTimeout(30 * time.Second).
    SetRetryCount(3).
    SetRetryWaitTime(1 * time.Second).
    SetRetryMaxWaitTime(5 * time.Second)
```

### GET Request

```go
type Movie struct {
    ID       int    `json:"id"`
    Title    string `json:"title"`
    Overview string `json:"overview"`
}

var movie Movie
resp, err := client.R().
    SetContext(ctx).
    SetResult(&movie).
    SetPathParam("id", "550").
    Get("/movie/{id}")

if err != nil {
    return nil, fmt.Errorf("request failed: %w", err)
}
if resp.IsError() {
    return nil, fmt.Errorf("API error: %d - %s", resp.StatusCode(), resp.String())
}

return &movie, nil
```

### POST Request

```go
type CreatePlaylistRequest struct {
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Items       []string `json:"items"`
}

type Playlist struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

var playlist Playlist
resp, err := client.R().
    SetContext(ctx).
    SetBody(&CreatePlaylistRequest{
        Name:        "My Playlist",
        Description: "A test playlist",
        Items:       []string{"item1", "item2"},
    }).
    SetResult(&playlist).
    Post("/playlists")
```

## Metadata Provider Client

```go
package tmdb

import (
    "context"
    "fmt"
    "log/slog"
    "time"
    "resty.dev/v3"
)

type Client struct {
    http   *resty.Client
    logger *slog.Logger
}

type Config struct {
    APIKey  string        `koanf:"api_key"`
    BaseURL string        `koanf:"base_url"`
    Timeout time.Duration `koanf:"timeout"`
}

var DefaultConfig = Config{
    BaseURL: "https://api.themoviedb.org/3",
    Timeout: 30 * time.Second,
}

func NewClient(cfg Config, logger *slog.Logger) *Client {
    http := resty.New().
        SetBaseURL(cfg.BaseURL).
        SetQueryParam("api_key", cfg.APIKey).
        SetHeader("Accept", "application/json").
        SetTimeout(cfg.Timeout).
        SetRetryCount(3).
        SetRetryWaitTime(500 * time.Millisecond).
        SetRetryMaxWaitTime(5 * time.Second).
        AddRetryCondition(func(r *resty.Response, err error) bool {
            // Retry on 429 (rate limit) and 5xx errors
            return err != nil || r.StatusCode() == 429 || r.StatusCode() >= 500
        }).
        OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
            logger.Debug("TMDb request", "method", r.Method, "url", r.URL)
            return nil
        }).
        OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
            logger.Debug("TMDb response",
                "status", r.StatusCode(),
                "duration", r.Time(),
            )
            return nil
        })

    return &Client{http: http, logger: logger}
}

func (c *Client) GetMovie(ctx context.Context, tmdbID int) (*Movie, error) {
    var movie Movie
    var apiError APIError

    resp, err := c.http.R().
        SetContext(ctx).
        SetResult(&movie).
        SetError(&apiError).
        SetPathParam("id", fmt.Sprint(tmdbID)).
        SetQueryParam("append_to_response", "credits,images,videos").
        Get("/movie/{id}")

    if err != nil {
        return nil, fmt.Errorf("TMDb request failed: %w", err)
    }
    if resp.IsError() {
        return nil, fmt.Errorf("TMDb error %d: %s", resp.StatusCode(), apiError.Message)
    }

    return &movie, nil
}

func (c *Client) SearchMovies(ctx context.Context, query string, page int) (*SearchResult, error) {
    var result SearchResult

    resp, err := c.http.R().
        SetContext(ctx).
        SetResult(&result).
        SetQueryParams(map[string]string{
            "query": query,
            "page":  fmt.Sprint(page),
        }).
        Get("/search/movie")

    if err != nil {
        return nil, err
    }
    if resp.IsError() {
        return nil, fmt.Errorf("search failed: %d", resp.StatusCode())
    }

    return &result, nil
}
```

## Rate Limiting

For APIs with strict rate limits (MusicBrainz: 1 req/sec):

```go
import "golang.org/x/time/rate"

type RateLimitedClient struct {
    http    *resty.Client
    limiter *rate.Limiter
}

func NewRateLimitedClient(rps float64) *RateLimitedClient {
    return &RateLimitedClient{
        http:    resty.New(),
        limiter: rate.NewLimiter(rate.Limit(rps), 1),
    }
}

func (c *RateLimitedClient) Do(ctx context.Context, req *resty.Request) (*resty.Response, error) {
    // Wait for rate limiter
    if err := c.limiter.Wait(ctx); err != nil {
        return nil, fmt.Errorf("rate limit wait: %w", err)
    }
    return req.SetContext(ctx).Send()
}
```

## Servarr Client (Radarr/Sonarr/Lidarr)

```go
package servarr

type RadarrClient struct {
    http *resty.Client
}

func NewRadarrClient(baseURL, apiKey string) *RadarrClient {
    http := resty.New().
        SetBaseURL(baseURL).
        SetHeader("X-Api-Key", apiKey).
        SetHeader("Accept", "application/json").
        SetTimeout(30 * time.Second)

    return &RadarrClient{http: http}
}

func (c *RadarrClient) GetMovies(ctx context.Context) ([]RadarrMovie, error) {
    var movies []RadarrMovie

    resp, err := c.http.R().
        SetContext(ctx).
        SetResult(&movies).
        Get("/api/v3/movie")

    if err != nil {
        return nil, err
    }
    if resp.IsError() {
        return nil, fmt.Errorf("Radarr error: %d", resp.StatusCode())
    }

    return movies, nil
}

func (c *RadarrClient) RefreshMovie(ctx context.Context, movieID int) error {
    resp, err := c.http.R().
        SetContext(ctx).
        SetBody(map[string]any{
            "name":    "RefreshMovie",
            "movieId": movieID,
        }).
        Post("/api/v3/command")

    if err != nil {
        return err
    }
    if resp.IsError() {
        return fmt.Errorf("refresh failed: %d", resp.StatusCode())
    }

    return nil
}
```

## Error Handling

```go
type APIError struct {
    StatusCode int    `json:"status_code"`
    Message    string `json:"status_message"`
}

func (e APIError) Error() string {
    return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

// Use with SetError
resp, err := client.R().
    SetContext(ctx).
    SetResult(&result).
    SetError(&APIError{}).
    Get("/endpoint")

if resp.IsError() {
    if apiErr, ok := resp.Error().(*APIError); ok {
        return nil, apiErr
    }
    return nil, fmt.Errorf("unknown error: %d", resp.StatusCode())
}
```

## fx Module Integration

```go
package http

import (
    "resty.dev/v3"
    "go.uber.org/fx"
)

func NewRestyClient(cfg Config) *resty.Client {
    return resty.New().
        SetTimeout(cfg.Timeout).
        SetRetryCount(cfg.RetryCount).
        SetRetryWaitTime(cfg.RetryWaitTime)
}

var Module = fx.Module("http",
    fx.Provide(NewRestyClient),
)
```

## DO's and DON'Ts

### DO

- ✅ Always use `SetContext(ctx)` for cancellation
- ✅ Use `SetResult()` and `SetError()` for automatic unmarshaling
- ✅ Configure retry with appropriate conditions
- ✅ Use rate limiting for strict APIs
- ✅ Log requests/responses for debugging
- ✅ Set reasonable timeouts

### DON'T

- ❌ Use standard `net/http` for external APIs - use resty
- ❌ Forget error checking on response status
- ❌ Use unbounded retries
- ❌ Ignore rate limits (especially MusicBrainz!)
- ❌ Store API keys in code - use config
---

## Related

- [INDEX.instructions.md](INDEX.instructions.md) - Main instruction index with all cross-references
- [ARCHITECTURE_V2.md](../../docs/dev/design/architecture/ARCHITECTURE_V2.md) - System architecture
- [BEST_PRACTICES.md](../../docs/dev/design/operations/BEST_PRACTICES.md) - Best practices

package tmdb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/viccon/sturdyc"

	"github.com/lusoris/revenge/pkg/resilience"
)

// Client errors.
var (
	ErrNotFound    = errors.New("movie not found")
	ErrRateLimited = errors.New("rate limit exceeded")
	ErrUnavailable = errors.New("service unavailable")
	ErrInvalidKey  = errors.New("invalid API key")
)

// Config holds TMDb client configuration.
type Config struct {
	APIKey     string        `koanf:"api_key"`
	BaseURL    string        `koanf:"base_url"`
	ImageURL   string        `koanf:"image_url"`
	Timeout    time.Duration `koanf:"timeout"`
	CacheTTL   time.Duration `koanf:"cache_ttl"`
	CacheSize  int           `koanf:"cache_size"`
	RetryCount int           `koanf:"retry_count"`
}

// DefaultConfig returns sensible defaults.
var DefaultConfig = Config{
	BaseURL:    "https://api.themoviedb.org/3",
	ImageURL:   "https://image.tmdb.org/t/p",
	Timeout:    30 * time.Second,
	CacheTTL:   1 * time.Hour,
	CacheSize:  50_000,
	RetryCount: 3,
}

// Client is a TMDb API client with caching and resilience.
type Client struct {
	http    *resty.Client
	cache   *sturdyc.Client[any]
	breaker *resilience.CircuitBreaker
	config  Config
	logger  *slog.Logger
}

// NewClient creates a new TMDb client.
func NewClient(cfg Config, logger *slog.Logger) *Client {
	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultConfig.BaseURL
	}
	if cfg.ImageURL == "" {
		cfg.ImageURL = DefaultConfig.ImageURL
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = DefaultConfig.Timeout
	}
	if cfg.CacheTTL == 0 {
		cfg.CacheTTL = DefaultConfig.CacheTTL
	}
	if cfg.CacheSize == 0 {
		cfg.CacheSize = DefaultConfig.CacheSize
	}
	if cfg.RetryCount == 0 {
		cfg.RetryCount = DefaultConfig.RetryCount
	}

	// Create resty client with retry and logging
	http := resty.New().
		SetBaseURL(cfg.BaseURL).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", "Bearer "+cfg.APIKey).
		SetTimeout(cfg.Timeout).
		SetRetryCount(cfg.RetryCount).
		SetRetryWaitTime(500 * time.Millisecond).
		SetRetryMaxWaitTime(5 * time.Second).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			// Retry on 429 (rate limit) and 5xx errors
			if err != nil {
				return true
			}
			return r.StatusCode() == 429 || r.StatusCode() >= 500
		}).
		OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
			logger.Debug("TMDb request", "method", r.Method, "url", r.URL)
			return nil
		}).
		OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
			logger.Debug("TMDb response", "status", r.StatusCode())
			return nil
		})

	// Create sturdyc cache with early refreshes for stale-while-revalidate
	cache := sturdyc.New[any](
		cfg.CacheSize,
		16, // shards
		cfg.CacheTTL,
		5, // eviction percentage
		sturdyc.WithEarlyRefreshes(
			cfg.CacheTTL/2,    // min async refresh time
			cfg.CacheTTL*3/4,  // max async refresh time
			cfg.CacheTTL*9/10, // sync refresh time
			5*time.Minute,     // retry base delay
		),
	)

	// Create circuit breaker
	breaker := resilience.NewCircuitBreaker(resilience.CircuitBreakerConfig{
		Name:                "tmdb",
		MaxFailures:         5,
		Timeout:             30 * time.Second,
		MaxHalfOpenRequests: 3,
		OnStateChange: func(name string, from, to int) {
			states := []string{"closed", "open", "half-open"}
			logger.Warn("circuit breaker state change",
				"name", name,
				"from", states[from],
				"to", states[to],
			)
		},
		IsSuccessful: func(err error) bool {
			// Don't count 404s as failures
			return err == nil || errors.Is(err, ErrNotFound)
		},
	})

	return &Client{
		http:    http,
		cache:   cache,
		breaker: breaker,
		config:  cfg,
		logger:  logger,
	}
}

// GetMovie fetches a movie by TMDb ID with full details.
func (c *Client) GetMovie(ctx context.Context, tmdbID int) (*Movie, error) {
	key := fmt.Sprintf("tmdb:movie:%d", tmdbID)

	result, err := sturdyc.GetOrFetch(ctx, c.cache, key, func(ctx context.Context) (*Movie, error) {
		return c.fetchMovie(ctx, tmdbID)
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) fetchMovie(ctx context.Context, tmdbID int) (*Movie, error) {
	var movie Movie
	var apiErr APIError

	err := c.breaker.ExecuteWithContext(ctx, func(ctx context.Context) error {
		resp, err := c.http.R().
			SetContext(ctx).
			SetResult(&movie).
			SetError(&apiErr).
			SetPathParam("id", fmt.Sprint(tmdbID)).
			SetQueryParam("append_to_response", "credits,images,videos,external_ids").
			Get("/movie/{id}")

		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}

		return c.handleResponse(resp, &apiErr)
	})

	if err != nil {
		return nil, err
	}

	return &movie, nil
}

// SearchMovies searches for movies by query.
func (c *Client) SearchMovies(ctx context.Context, query string, year int, page int) (*SearchResult, error) {
	key := fmt.Sprintf("tmdb:search:%s:%d:%d", query, year, page)

	result, err := sturdyc.GetOrFetch(ctx, c.cache, key, func(ctx context.Context) (*SearchResult, error) {
		return c.searchMovies(ctx, query, year, page)
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) searchMovies(ctx context.Context, query string, year int, page int) (*SearchResult, error) {
	var result SearchResult
	var apiErr APIError

	err := c.breaker.ExecuteWithContext(ctx, func(ctx context.Context) error {
		req := c.http.R().
			SetContext(ctx).
			SetResult(&result).
			SetError(&apiErr).
			SetQueryParam("query", query).
			SetQueryParam("page", fmt.Sprint(page))

		if year > 0 {
			req.SetQueryParam("year", fmt.Sprint(year))
		}

		resp, err := req.Get("/search/movie")
		if err != nil {
			return fmt.Errorf("search request failed: %w", err)
		}

		return c.handleResponse(resp, &apiErr)
	})

	if err != nil {
		return nil, err
	}

	return &result, nil
}

// FindByIMDbID finds a movie by IMDb ID.
func (c *Client) FindByIMDbID(ctx context.Context, imdbID string) (*MovieResult, error) {
	key := fmt.Sprintf("tmdb:find:imdb:%s", imdbID)

	result, err := sturdyc.GetOrFetch(ctx, c.cache, key, func(ctx context.Context) (*MovieResult, error) {
		return c.findByIMDbID(ctx, imdbID)
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) findByIMDbID(ctx context.Context, imdbID string) (*MovieResult, error) {
	var result FindResult
	var apiErr APIError

	err := c.breaker.ExecuteWithContext(ctx, func(ctx context.Context) error {
		resp, err := c.http.R().
			SetContext(ctx).
			SetResult(&result).
			SetError(&apiErr).
			SetPathParam("id", imdbID).
			SetQueryParam("external_source", "imdb_id").
			Get("/find/{id}")

		if err != nil {
			return fmt.Errorf("find request failed: %w", err)
		}

		return c.handleResponse(resp, &apiErr)
	})

	if err != nil {
		return nil, err
	}

	if len(result.MovieResults) == 0 {
		return nil, ErrNotFound
	}

	return &result.MovieResults[0], nil
}

// GetCollection fetches a collection by ID.
func (c *Client) GetCollection(ctx context.Context, collectionID int) (*Collection, error) {
	key := fmt.Sprintf("tmdb:collection:%d", collectionID)

	result, err := sturdyc.GetOrFetch(ctx, c.cache, key, func(ctx context.Context) (*Collection, error) {
		return c.fetchCollection(ctx, collectionID)
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) fetchCollection(ctx context.Context, collectionID int) (*Collection, error) {
	var collection Collection
	var apiErr APIError

	err := c.breaker.ExecuteWithContext(ctx, func(ctx context.Context) error {
		resp, err := c.http.R().
			SetContext(ctx).
			SetResult(&collection).
			SetError(&apiErr).
			SetPathParam("id", fmt.Sprint(collectionID)).
			Get("/collection/{id}")

		if err != nil {
			return fmt.Errorf("collection request failed: %w", err)
		}

		return c.handleResponse(resp, &apiErr)
	})

	if err != nil {
		return nil, err
	}

	return &collection, nil
}

// ImageURL builds a full image URL.
func (c *Client) ImageURL(path string, size string) string {
	if path == "" {
		return ""
	}
	if size == "" {
		size = "original"
	}
	return fmt.Sprintf("%s/%s%s", c.config.ImageURL, size, path)
}

// PosterURL builds a poster URL with default size.
func (c *Client) PosterURL(path string) string {
	return c.ImageURL(path, "w500")
}

// BackdropURL builds a backdrop URL with default size.
func (c *Client) BackdropURL(path string) string {
	return c.ImageURL(path, "w1280")
}

// ProfileURL builds a profile URL with default size.
func (c *Client) ProfileURL(path string) string {
	return c.ImageURL(path, "w185")
}

// handleResponse handles HTTP response and converts to errors.
func (c *Client) handleResponse(resp *resty.Response, apiErr *APIError) error {
	switch resp.StatusCode() {
	case 200, 201:
		return nil
	case 401:
		return ErrInvalidKey
	case 404:
		return ErrNotFound
	case 429:
		return ErrRateLimited
	case 500, 502, 503, 504:
		return ErrUnavailable
	default:
		if apiErr != nil && apiErr.StatusMessage != "" {
			return apiErr
		}
		return fmt.Errorf("unexpected status: %d", resp.StatusCode())
	}
}

// Stats returns cache and circuit breaker statistics.
func (c *Client) Stats() ClientStats {
	cbStats := c.breaker.Stats()
	states := []string{"closed", "open", "half-open"}

	return ClientStats{
		CircuitBreakerState:    states[cbStats.State],
		CircuitBreakerFailures: cbStats.Failures,
		CircuitBreakerRequests: cbStats.Requests,
	}
}

// ClientStats contains client statistics.
type ClientStats struct {
	CircuitBreakerState    string
	CircuitBreakerFailures int
	CircuitBreakerRequests int
}

// Close releases client resources.
func (c *Client) Close() error {
	// resty.Client doesn't require explicit cleanup
	return nil
}

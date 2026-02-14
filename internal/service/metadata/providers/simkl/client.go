package simkl

import (
	"context"
	"fmt"
	"time"

	"github.com/imroc/req/v3"
	"golang.org/x/time/rate"

	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/circuitbreaker"
)

const (
	baseURL      = "https://api.simkl.com"
	imageBaseURL = "https://wsrv.nl/?url=https://simkl.in"
)

// Config contains configuration for the Simkl client.
type Config struct {
	// Enabled activates Simkl as a metadata provider.
	Enabled bool

	// ClientID is the Simkl API client ID (required).
	// Obtain from https://simkl.com/settings/developer/
	ClientID string

	// RateLimit is requests per second (default: 2.0).
	RateLimit rate.Limit

	// Burst is the burst capacity (default: 5).
	Burst int

	// CacheTTL is the cache duration (default: 1h).
	CacheTTL time.Duration

	// Timeout is the request timeout (default: 15s).
	Timeout time.Duration
}

// Client is a client for the Simkl API.
type Client struct {
	client      *req.Client
	rateLimiter *rate.Limiter
	cache       *cache.L1Cache[string, any]
}

// NewClient creates a new Simkl API client.
func NewClient(config Config) (*Client, error) {
	if config.ClientID == "" {
		return nil, fmt.Errorf("simkl client ID is required")
	}
	if config.Timeout == 0 {
		config.Timeout = 15 * time.Second
	}
	if config.RateLimit == 0 {
		config.RateLimit = rate.Limit(2.0)
	}
	if config.Burst == 0 {
		config.Burst = 5
	}
	if config.CacheTTL == 0 {
		config.CacheTTL = 1 * time.Hour
	}

	client := req.C().
		SetBaseURL(baseURL).
		SetTimeout(config.Timeout).
		SetCommonHeader("Content-Type", "application/json").
		SetCommonHeader("simkl-api-key", config.ClientID).
		SetCommonRetryCount(2).
		SetCommonRetryBackoffInterval(1*time.Second, 5*time.Second).
		SetCommonRetryCondition(func(resp *req.Response, err error) bool {
			if err != nil {
				return true
			}
			return resp.StatusCode >= 500
		})

	// Circuit breaker
	circuitbreaker.WrapReqClient(client, "simkl", circuitbreaker.TierExternal)

	l1, err := cache.NewL1Cache[string, any](2000, config.CacheTTL, cache.WithExpiryAccessing[string, any]())
	if err != nil {
		l1, _ = cache.NewL1Cache[string, any](0, 0)
	}

	return &Client{
		client:      client,
		rateLimiter: rate.NewLimiter(config.RateLimit, config.Burst),
		cache:       l1,
	}, nil
}

func (c *Client) getFromCache(key string) any {
	if v, ok := c.cache.Get(key); ok {
		return v
	}
	return nil
}

func (c *Client) setCache(key string, data any) {
	c.cache.Set(key, data)
}

func (c *Client) clearCache() {
	c.cache.Clear()
}

// SearchMovies searches for movies by text query.
func (c *Client) SearchMovies(ctx context.Context, query string) ([]SearchResult, error) {
	cacheKey := fmt.Sprintf("search:movie:%s", query)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.([]SearchResult); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result []SearchResult
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParam("q", query).
		SetQueryParam("type", "movie").
		SetQueryParam("extended", "full").
		SetSuccessResult(&result).
		Get("/search/text")

	if err != nil {
		return nil, fmt.Errorf("simkl search: %w", err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("simkl search error: %s", resp.Status)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// SearchShows searches for TV shows by text query.
func (c *Client) SearchShows(ctx context.Context, query string) ([]SearchResult, error) {
	cacheKey := fmt.Sprintf("search:tv:%s", query)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.([]SearchResult); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result []SearchResult
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParam("q", query).
		SetQueryParam("type", "tv").
		SetQueryParam("extended", "full").
		SetSuccessResult(&result).
		Get("/search/text")

	if err != nil {
		return nil, fmt.Errorf("simkl search: %w", err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("simkl search error: %s", resp.Status)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// SearchAnime searches for anime by text query.
func (c *Client) SearchAnime(ctx context.Context, query string) ([]SearchResult, error) {
	cacheKey := fmt.Sprintf("search:anime:%s", query)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.([]SearchResult); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result []SearchResult
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParam("q", query).
		SetQueryParam("type", "anime").
		SetQueryParam("extended", "full").
		SetSuccessResult(&result).
		Get("/search/text")

	if err != nil {
		return nil, fmt.Errorf("simkl search: %w", err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("simkl search error: %s", resp.Status)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// GetMovie returns movie details by Simkl ID.
func (c *Client) GetMovie(ctx context.Context, id string) (*Movie, error) {
	cacheKey := fmt.Sprintf("movie:%s", id)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.(*Movie); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result Movie
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParam("extended", "full").
		SetSuccessResult(&result).
		Get(fmt.Sprintf("/movies/%s", id))

	if err != nil {
		return nil, fmt.Errorf("simkl get movie: %w", err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("simkl get movie error: %s", resp.Status)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

// GetShow returns TV show details by Simkl ID.
func (c *Client) GetShow(ctx context.Context, id string) (*Show, error) {
	cacheKey := fmt.Sprintf("show:%s", id)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.(*Show); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result Show
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParam("extended", "full").
		SetSuccessResult(&result).
		Get(fmt.Sprintf("/tv/%s", id))

	if err != nil {
		return nil, fmt.Errorf("simkl get show: %w", err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("simkl get show error: %s", resp.Status)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

// GetAnime returns anime details by Simkl ID.
func (c *Client) GetAnime(ctx context.Context, id string) (*Show, error) {
	cacheKey := fmt.Sprintf("anime:%s", id)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.(*Show); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result Show
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParam("extended", "full").
		SetSuccessResult(&result).
		Get(fmt.Sprintf("/anime/%s", id))

	if err != nil {
		return nil, fmt.Errorf("simkl get anime: %w", err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("simkl get anime error: %s", resp.Status)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

// GetShowEpisodes returns all episodes for a TV show.
func (c *Client) GetShowEpisodes(ctx context.Context, id string) ([]Episode, error) {
	cacheKey := fmt.Sprintf("show:%s:episodes", id)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.([]Episode); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result []Episode
	resp, err := c.client.R().
		SetContext(ctx).
		SetSuccessResult(&result).
		Get(fmt.Sprintf("/tv/%s/episodes", id))

	if err != nil {
		return nil, fmt.Errorf("simkl get episodes: %w", err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("simkl get episodes error: %s", resp.Status)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// GetAnimeEpisodes returns all episodes for an anime.
func (c *Client) GetAnimeEpisodes(ctx context.Context, id string) ([]Episode, error) {
	cacheKey := fmt.Sprintf("anime:%s:episodes", id)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.([]Episode); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result []Episode
	resp, err := c.client.R().
		SetContext(ctx).
		SetSuccessResult(&result).
		Get(fmt.Sprintf("/anime/%s/episodes", id))

	if err != nil {
		return nil, fmt.Errorf("simkl get anime episodes: %w", err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("simkl get anime episodes error: %s", resp.Status)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// LookupByIMDb looks up items by IMDb ID.
func (c *Client) LookupByIMDb(ctx context.Context, imdbID string) ([]IDLookupResult, error) {
	cacheKey := fmt.Sprintf("lookup:imdb:%s", imdbID)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.([]IDLookupResult); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result []IDLookupResult
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParam("imdb", imdbID).
		SetSuccessResult(&result).
		Get("/search/id")

	if err != nil {
		return nil, fmt.Errorf("simkl lookup: %w", err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("simkl lookup error: %s", resp.Status)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// LookupByTMDb looks up items by TMDb ID.
func (c *Client) LookupByTMDb(ctx context.Context, tmdbID int) ([]IDLookupResult, error) {
	cacheKey := fmt.Sprintf("lookup:tmdb:%d", tmdbID)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.([]IDLookupResult); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result []IDLookupResult
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParam("tmdb", fmt.Sprint(tmdbID)).
		SetSuccessResult(&result).
		Get("/search/id")

	if err != nil {
		return nil, fmt.Errorf("simkl lookup: %w", err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("simkl lookup error: %s", resp.Status)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// GetRating returns rating information for a Simkl item.
func (c *Client) GetRating(ctx context.Context, simklID int) (*Ratings, error) {
	cacheKey := fmt.Sprintf("rating:%d", simklID)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.(*Ratings); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result Ratings
	resp, err := c.client.R().
		SetContext(ctx).
		SetSuccessResult(&result).
		Get(fmt.Sprintf("/ratings/%d", simklID))

	if err != nil {
		return nil, fmt.Errorf("simkl get rating: %w", err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("simkl get rating error: %s", resp.Status)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

// Close stops the cache's background goroutines.
func (c *Client) Close() {
	if c.cache != nil {
		c.cache.Close()
	}
}

// ImageURL constructs a full Simkl image URL from a relative path with a size suffix.
func ImageURL(path, suffix string) string {
	if path == "" {
		return ""
	}
	return fmt.Sprintf("%s/posters/%s%s.jpg", imageBaseURL, path, suffix)
}

// FanartURL constructs a full Simkl fanart URL.
func FanartURL(path, suffix string) string {
	if path == "" {
		return ""
	}
	return fmt.Sprintf("%s/fanart/%s%s.jpg", imageBaseURL, path, suffix)
}

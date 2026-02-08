package trakt

import (
	"context"
	"fmt"
	"time"

	"github.com/imroc/req/v3"
	"golang.org/x/time/rate"

	"github.com/lusoris/revenge/internal/infra/cache"
)

const (
	baseURL    = "https://api.trakt.tv"
	apiVersion = "2"
)

// Config contains configuration for the Trakt client.
type Config struct {
	// Enabled activates Trakt as a metadata provider.
	Enabled bool

	// ClientID is the Trakt API client ID (required).
	// Obtain from https://trakt.tv/oauth/applications
	ClientID string

	// RateLimit is requests per second (default: 3.0, Trakt allows 1000/5min).
	RateLimit rate.Limit

	// Burst is the burst capacity (default: 10).
	Burst int

	// CacheTTL is the cache duration (default: 1h).
	CacheTTL time.Duration

	// Timeout is the request timeout (default: 15s).
	Timeout time.Duration
}

// Client is a client for the Trakt API v2.
type Client struct {
	client      *req.Client
	rateLimiter *rate.Limiter
	cache       *cache.L1Cache[string, any]
}

// NewClient creates a new Trakt API client.
func NewClient(config Config) (*Client, error) {
	if config.ClientID == "" {
		return nil, fmt.Errorf("trakt client ID is required")
	}
	if config.Timeout == 0 {
		config.Timeout = 15 * time.Second
	}
	if config.RateLimit == 0 {
		config.RateLimit = rate.Limit(3.0) // 1000 per 5 min = ~3.3/s
	}
	if config.Burst == 0 {
		config.Burst = 10
	}
	if config.CacheTTL == 0 {
		config.CacheTTL = 1 * time.Hour
	}

	client := req.C().
		SetBaseURL(baseURL).
		SetTimeout(config.Timeout).
		SetCommonHeader("Content-Type", "application/json").
		SetCommonHeader("trakt-api-version", apiVersion).
		SetCommonHeader("trakt-api-key", config.ClientID).
		SetCommonRetryCount(2).
		SetCommonRetryBackoffInterval(1*time.Second, 5*time.Second)

	l1, err := cache.NewL1Cache[string, any](2000, config.CacheTTL)
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
		SetQueryParam("query", query).
		SetSuccessResult(&result).
		Get("/search/movie")

	if err != nil {
		return nil, fmt.Errorf("trakt search: %w", err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("trakt search error: %s", resp.Status)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// SearchShows searches for TV shows by text query.
func (c *Client) SearchShows(ctx context.Context, query string) ([]SearchResult, error) {
	cacheKey := fmt.Sprintf("search:show:%s", query)
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
		SetQueryParam("query", query).
		SetSuccessResult(&result).
		Get("/search/show")

	if err != nil {
		return nil, fmt.Errorf("trakt search: %w", err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("trakt search error: %s", resp.Status)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// GetMovie returns extended movie details by Trakt slug or ID.
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
		return nil, fmt.Errorf("trakt get movie: %w", err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("trakt get movie error: %s", resp.Status)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

// GetShow returns extended show details by Trakt slug or ID.
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
		Get(fmt.Sprintf("/shows/%s", id))

	if err != nil {
		return nil, fmt.Errorf("trakt get show: %w", err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("trakt get show error: %s", resp.Status)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

// GetShowSeasons returns all seasons for a show.
func (c *Client) GetShowSeasons(ctx context.Context, id string) ([]Season, error) {
	cacheKey := fmt.Sprintf("show:%s:seasons", id)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.([]Season); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result []Season
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParam("extended", "full").
		SetSuccessResult(&result).
		Get(fmt.Sprintf("/shows/%s/seasons", id))

	if err != nil {
		return nil, fmt.Errorf("trakt get seasons: %w", err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("trakt get seasons error: %s", resp.Status)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// GetSeasonEpisodes returns all episodes for a season.
func (c *Client) GetSeasonEpisodes(ctx context.Context, showID string, season int) ([]Episode, error) {
	cacheKey := fmt.Sprintf("show:%s:season:%d:episodes", showID, season)
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
		SetQueryParam("extended", "full").
		SetSuccessResult(&result).
		Get(fmt.Sprintf("/shows/%s/seasons/%d", showID, season))

	if err != nil {
		return nil, fmt.Errorf("trakt get episodes: %w", err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("trakt get episodes error: %s", resp.Status)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// GetMovieCredits returns cast and crew for a movie.
func (c *Client) GetMovieCredits(ctx context.Context, id string) (*Credits, error) {
	cacheKey := fmt.Sprintf("movie:%s:credits", id)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.(*Credits); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result Credits
	resp, err := c.client.R().
		SetContext(ctx).
		SetSuccessResult(&result).
		Get(fmt.Sprintf("/movies/%s/people", id))

	if err != nil {
		return nil, fmt.Errorf("trakt get movie credits: %w", err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("trakt get movie credits error: %s", resp.Status)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

// GetShowCredits returns cast and crew for a show.
func (c *Client) GetShowCredits(ctx context.Context, id string) (*Credits, error) {
	cacheKey := fmt.Sprintf("show:%s:credits", id)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.(*Credits); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result Credits
	resp, err := c.client.R().
		SetContext(ctx).
		SetSuccessResult(&result).
		Get(fmt.Sprintf("/shows/%s/people", id))

	if err != nil {
		return nil, fmt.Errorf("trakt get show credits: %w", err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("trakt get show credits error: %s", resp.Status)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

// GetMovieTranslations returns translations for a movie.
func (c *Client) GetMovieTranslations(ctx context.Context, id string) ([]Translation, error) {
	cacheKey := fmt.Sprintf("movie:%s:translations", id)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.([]Translation); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result []Translation
	resp, err := c.client.R().
		SetContext(ctx).
		SetSuccessResult(&result).
		Get(fmt.Sprintf("/movies/%s/translations", id))

	if err != nil {
		return nil, fmt.Errorf("trakt get translations: %w", err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("trakt get translations error: %s", resp.Status)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// GetShowTranslations returns translations for a show.
func (c *Client) GetShowTranslations(ctx context.Context, id string) ([]Translation, error) {
	cacheKey := fmt.Sprintf("show:%s:translations", id)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if result, ok := cached.([]Translation); ok {
			return result, nil
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit: %w", err)
	}

	var result []Translation
	resp, err := c.client.R().
		SetContext(ctx).
		SetSuccessResult(&result).
		Get(fmt.Sprintf("/shows/%s/translations", id))

	if err != nil {
		return nil, fmt.Errorf("trakt get translations: %w", err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("trakt get translations error: %s", resp.Status)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

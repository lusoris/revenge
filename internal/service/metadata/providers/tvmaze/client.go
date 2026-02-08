package tvmaze

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/imroc/req/v3"
	"golang.org/x/time/rate"

	"github.com/lusoris/revenge/internal/infra/cache"
)

const (
	// BaseURL is the TVmaze API base URL.
	BaseURL = "https://api.tvmaze.com"

	// DefaultRateLimit is 2 requests per second (TVmaze allows 20/10s).
	DefaultRateLimit = rate.Limit(2.0)

	// DefaultBurst is the burst capacity.
	DefaultBurst = 5

	// DefaultCacheTTL is 1 hour (matches TVmaze's server cache).
	DefaultCacheTTL = 1 * time.Hour
)

// Config configures the TVmaze client.
type Config struct {
	// RateLimit is requests per second (default: 2).
	RateLimit rate.Limit

	// Burst is the burst capacity (default: 5).
	Burst int

	// CacheTTL is the cache duration (default: 1h).
	CacheTTL time.Duration

	// Timeout is the request timeout (default: 10s).
	Timeout time.Duration
}

// Client is the TVmaze API client with rate limiting and caching.
type Client struct {
	httpClient  *req.Client
	rateLimiter *rate.Limiter
	cache       *cache.L1Cache[string, any]
	cacheTTL    time.Duration
}

// NewClient creates a new TVmaze client.
func NewClient(config Config) (*Client, error) {
	if config.RateLimit == 0 {
		config.RateLimit = DefaultRateLimit
	}
	if config.Burst == 0 {
		config.Burst = DefaultBurst
	}
	if config.CacheTTL == 0 {
		config.CacheTTL = DefaultCacheTTL
	}
	if config.Timeout == 0 {
		config.Timeout = 10 * time.Second
	}

	l1, err := cache.NewL1Cache[string, any](10000, config.CacheTTL)
	if err != nil {
		return nil, fmt.Errorf("create tvmaze cache: %w", err)
	}

	client := req.C().
		SetBaseURL(BaseURL).
		SetTimeout(config.Timeout).
		SetCommonRetryCount(2).
		SetCommonRetryBackoffInterval(1*time.Second, 5*time.Second)

	return &Client{
		httpClient:  client,
		rateLimiter: rate.NewLimiter(config.RateLimit, config.Burst),
		cache:       l1,
		cacheTTL:    config.CacheTTL,
	}, nil
}

func (c *Client) waitRateLimit(ctx context.Context) error {
	return c.rateLimiter.Wait(ctx)
}

func (c *Client) getFromCache(key string) any {
	val, ok := c.cache.Get(key)
	if !ok {
		return nil
	}
	return val
}

func (c *Client) setCache(key string, data any) {
	c.cache.Set(key, data)
}

func (c *Client) clearCache() {
	c.cache.Clear()
}

// SearchShows searches for TV shows by query.
func (c *Client) SearchShows(ctx context.Context, query string) ([]ShowSearchResult, error) {
	cacheKey := "search:" + query
	if cached := c.getFromCache(cacheKey); cached != nil {
		if v, ok := cached.([]ShowSearchResult); ok {
			return v, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, err
	}

	var result []ShowSearchResult
	resp, err := c.httpClient.R().SetContext(ctx).
		SetQueryParam("q", query).
		SetSuccessResult(&result).
		Get("/search/shows")
	if err != nil {
		return nil, fmt.Errorf("tvmaze search: %w", err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("tvmaze search: status %d", resp.StatusCode)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// GetShow fetches show details by TVmaze ID.
func (c *Client) GetShow(ctx context.Context, id int) (*Show, error) {
	cacheKey := "show:" + strconv.Itoa(id)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if v, ok := cached.(*Show); ok {
			return v, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, err
	}

	var result Show
	resp, err := c.httpClient.R().SetContext(ctx).
		SetSuccessResult(&result).
		Get("/shows/" + strconv.Itoa(id))
	if err != nil {
		return nil, fmt.Errorf("tvmaze show: %w", err)
	}
	if resp.IsErrorState() {
		if resp.StatusCode == 404 {
			return nil, nil
		}
		return nil, fmt.Errorf("tvmaze show: status %d", resp.StatusCode)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

// LookupByTVDbID looks up a show by TVDb ID.
func (c *Client) LookupByTVDbID(ctx context.Context, tvdbID int) (*Show, error) {
	cacheKey := "lookup:tvdb:" + strconv.Itoa(tvdbID)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if v, ok := cached.(*Show); ok {
			return v, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, err
	}

	var result Show
	resp, err := c.httpClient.R().SetContext(ctx).
		SetQueryParam("thetvdb", strconv.Itoa(tvdbID)).
		SetSuccessResult(&result).
		Get("/lookup/shows")
	if err != nil {
		return nil, fmt.Errorf("tvmaze lookup: %w", err)
	}
	if resp.IsErrorState() {
		if resp.StatusCode == 404 {
			return nil, nil
		}
		return nil, fmt.Errorf("tvmaze lookup: status %d", resp.StatusCode)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

// LookupByIMDbID looks up a show by IMDb ID.
func (c *Client) LookupByIMDbID(ctx context.Context, imdbID string) (*Show, error) {
	cacheKey := "lookup:imdb:" + imdbID
	if cached := c.getFromCache(cacheKey); cached != nil {
		if v, ok := cached.(*Show); ok {
			return v, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, err
	}

	var result Show
	resp, err := c.httpClient.R().SetContext(ctx).
		SetQueryParam("imdb", imdbID).
		SetSuccessResult(&result).
		Get("/lookup/shows")
	if err != nil {
		return nil, fmt.Errorf("tvmaze lookup: %w", err)
	}
	if resp.IsErrorState() {
		if resp.StatusCode == 404 {
			return nil, nil
		}
		return nil, fmt.Errorf("tvmaze lookup: status %d", resp.StatusCode)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

// GetSeasons fetches all seasons for a show.
func (c *Client) GetSeasons(ctx context.Context, showID int) ([]Season, error) {
	cacheKey := "seasons:" + strconv.Itoa(showID)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if v, ok := cached.([]Season); ok {
			return v, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, err
	}

	var result []Season
	resp, err := c.httpClient.R().SetContext(ctx).
		SetSuccessResult(&result).
		Get("/shows/" + strconv.Itoa(showID) + "/seasons")
	if err != nil {
		return nil, fmt.Errorf("tvmaze seasons: %w", err)
	}
	if resp.IsErrorState() {
		if resp.StatusCode == 404 {
			return nil, nil
		}
		return nil, fmt.Errorf("tvmaze seasons: status %d", resp.StatusCode)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// GetEpisodes fetches all episodes for a show.
func (c *Client) GetEpisodes(ctx context.Context, showID int) ([]Episode, error) {
	cacheKey := "episodes:" + strconv.Itoa(showID)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if v, ok := cached.([]Episode); ok {
			return v, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, err
	}

	var result []Episode
	resp, err := c.httpClient.R().SetContext(ctx).
		SetSuccessResult(&result).
		Get("/shows/" + strconv.Itoa(showID) + "/episodes")
	if err != nil {
		return nil, fmt.Errorf("tvmaze episodes: %w", err)
	}
	if resp.IsErrorState() {
		if resp.StatusCode == 404 {
			return nil, nil
		}
		return nil, fmt.Errorf("tvmaze episodes: status %d", resp.StatusCode)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// GetCast fetches the cast for a show.
func (c *Client) GetCast(ctx context.Context, showID int) ([]CastMember, error) {
	cacheKey := "cast:" + strconv.Itoa(showID)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if v, ok := cached.([]CastMember); ok {
			return v, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, err
	}

	var result []CastMember
	resp, err := c.httpClient.R().SetContext(ctx).
		SetSuccessResult(&result).
		Get("/shows/" + strconv.Itoa(showID) + "/cast")
	if err != nil {
		return nil, fmt.Errorf("tvmaze cast: %w", err)
	}
	if resp.IsErrorState() {
		if resp.StatusCode == 404 {
			return nil, nil
		}
		return nil, fmt.Errorf("tvmaze cast: status %d", resp.StatusCode)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// GetCrew fetches the crew for a show.
func (c *Client) GetCrew(ctx context.Context, showID int) ([]CrewMember, error) {
	cacheKey := "crew:" + strconv.Itoa(showID)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if v, ok := cached.([]CrewMember); ok {
			return v, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, err
	}

	var result []CrewMember
	resp, err := c.httpClient.R().SetContext(ctx).
		SetSuccessResult(&result).
		Get("/shows/" + strconv.Itoa(showID) + "/crew")
	if err != nil {
		return nil, fmt.Errorf("tvmaze crew: %w", err)
	}
	if resp.IsErrorState() {
		if resp.StatusCode == 404 {
			return nil, nil
		}
		return nil, fmt.Errorf("tvmaze crew: status %d", resp.StatusCode)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

// GetImages fetches images for a show.
func (c *Client) GetImages(ctx context.Context, showID int) ([]ShowImage, error) {
	cacheKey := "images:" + strconv.Itoa(showID)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if v, ok := cached.([]ShowImage); ok {
			return v, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, err
	}

	var result []ShowImage
	resp, err := c.httpClient.R().SetContext(ctx).
		SetSuccessResult(&result).
		Get("/shows/" + strconv.Itoa(showID) + "/images")
	if err != nil {
		return nil, fmt.Errorf("tvmaze images: %w", err)
	}
	if resp.IsErrorState() {
		if resp.StatusCode == 404 {
			return nil, nil
		}
		return nil, fmt.Errorf("tvmaze images: status %d", resp.StatusCode)
	}

	c.setCache(cacheKey, result)
	return result, nil
}

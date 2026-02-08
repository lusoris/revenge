package fanarttv

import (
	"context"
	"fmt"
	"time"

	"github.com/imroc/req/v3"
	"golang.org/x/time/rate"

	"github.com/lusoris/revenge/internal/infra/cache"
)

const (
	// BaseURL is the Fanart.tv API base URL.
	BaseURL = "https://webservice.fanart.tv/v3"

	// DefaultRateLimit is 1 request per second (conservative).
	DefaultRateLimit = rate.Limit(1.0)

	// DefaultBurst is the burst capacity.
	DefaultBurst = 5

	// DefaultCacheTTL is 7 days (images rarely change).
	DefaultCacheTTL = 7 * 24 * time.Hour
)

// Config configures the Fanart.tv client.
type Config struct {
	// APIKey is the Fanart.tv project API key (required).
	APIKey string

	// ClientKey is the personal API key for faster image access (optional).
	ClientKey string

	// RateLimit is requests per second (default: 1).
	RateLimit rate.Limit

	// Burst is the burst capacity (default: 5).
	Burst int

	// CacheTTL is the cache duration (default: 7 days).
	CacheTTL time.Duration

	// Timeout is the request timeout (default: 15s).
	Timeout time.Duration
}

// Client is the Fanart.tv API client with rate limiting and caching.
type Client struct {
	httpClient  *req.Client
	apiKey      string
	clientKey   string
	rateLimiter *rate.Limiter
	cache       *cache.L1Cache[string, any]
	cacheTTL    time.Duration
}

// NewClient creates a new Fanart.tv client.
func NewClient(config Config) (*Client, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("fanart.tv API key is required")
	}

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
		config.Timeout = 15 * time.Second
	}

	l1, err := cache.NewL1Cache[string, any](5000, config.CacheTTL)
	if err != nil {
		return nil, fmt.Errorf("create fanarttv cache: %w", err)
	}

	client := req.C().
		SetBaseURL(BaseURL).
		SetTimeout(config.Timeout).
		SetCommonRetryCount(2).
		SetCommonRetryBackoffInterval(1*time.Second, 5*time.Second)

	return &Client{
		httpClient:  client,
		apiKey:      config.APIKey,
		clientKey:   config.ClientKey,
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

func (c *Client) request(ctx context.Context) *req.Request {
	r := c.httpClient.R().SetContext(ctx).
		SetQueryParam("api_key", c.apiKey)

	if c.clientKey != "" {
		r.SetQueryParam("client_key", c.clientKey)
	}

	return r
}

// GetMovieImages fetches images for a movie by TMDb ID.
func (c *Client) GetMovieImages(ctx context.Context, tmdbID string) (*MovieResponse, error) {
	cacheKey := "movie:" + tmdbID
	if cached := c.getFromCache(cacheKey); cached != nil {
		if v, ok := cached.(*MovieResponse); ok {
			return v, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, err
	}

	var result MovieResponse
	resp, err := c.request(ctx).
		SetSuccessResult(&result).
		Get("/movies/" + tmdbID)
	if err != nil {
		return nil, fmt.Errorf("fanart.tv movie request: %w", err)
	}
	if resp.IsErrorState() {
		if resp.StatusCode == 404 {
			return nil, nil
		}
		return nil, fmt.Errorf("fanart.tv movie request: status %d", resp.StatusCode)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

// GetTVShowImages fetches images for a TV show by TVDb ID.
func (c *Client) GetTVShowImages(ctx context.Context, tvdbID string) (*TVShowResponse, error) {
	cacheKey := "tv:" + tvdbID
	if cached := c.getFromCache(cacheKey); cached != nil {
		if v, ok := cached.(*TVShowResponse); ok {
			return v, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, err
	}

	var result TVShowResponse
	resp, err := c.request(ctx).
		SetSuccessResult(&result).
		Get("/tv/" + tvdbID)
	if err != nil {
		return nil, fmt.Errorf("fanart.tv tv request: %w", err)
	}
	if resp.IsErrorState() {
		if resp.StatusCode == 404 {
			return nil, nil
		}
		return nil, fmt.Errorf("fanart.tv tv request: status %d", resp.StatusCode)
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

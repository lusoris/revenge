package omdb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/imroc/req/v3"
	"golang.org/x/time/rate"

	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/circuitbreaker"
	"github.com/lusoris/revenge/internal/infra/observability"
)

const (
	// BaseURL is the OMDb API base URL.
	BaseURL = "https://www.omdbapi.com"

	// DefaultRateLimit is 1 request per second (conservative for free tier).
	DefaultRateLimit = rate.Limit(1.0)

	// DefaultBurst is the burst capacity.
	DefaultBurst = 3

	// DefaultCacheTTL is 24 hours (ratings change slowly).
	DefaultCacheTTL = 24 * time.Hour
)

// Config configures the OMDb client.
type Config struct {
	// APIKey is the OMDb API key (required).
	APIKey string

	// RateLimit is requests per second (default: 1).
	RateLimit rate.Limit

	// Burst is the burst capacity (default: 3).
	Burst int

	// CacheTTL is the cache duration (default: 24h).
	CacheTTL time.Duration

	// Timeout is the request timeout (default: 10s).
	Timeout time.Duration
}

// Client is the OMDb API client with rate limiting and caching.
type Client struct {
	httpClient  *req.Client
	apiKey      string
	rateLimiter *rate.Limiter
	cache       *cache.L1Cache[string, any]
	cacheTTL    time.Duration
}

// NewClient creates a new OMDb client.
func NewClient(config Config) (*Client, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("omdb API key is required")
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
		config.Timeout = 10 * time.Second
	}

	l1, err := cache.NewL1Cache[string, any](10000, config.CacheTTL, cache.WithExpiryAccessing[string, any]())
	if err != nil {
		return nil, fmt.Errorf("create omdb cache: %w", err)
	}

	client := req.C().
		SetBaseURL(BaseURL).
		SetTimeout(config.Timeout).
		SetCommonRetryCount(2).
		SetCommonRetryBackoffInterval(1*time.Second, 5*time.Second).
		SetCommonRetryCondition(func(resp *req.Response, err error) bool {
			if err != nil {
				return true
			}
			return resp.StatusCode >= 500
		}).
		OnAfterResponse(func(_ *req.Client, resp *req.Response) error {
			// OMDB returns type in query param or response
			mediaType := "movie" // OMDB is primarily movie-focused
			if strings.Contains(resp.Request.RawURL, "type=series") {
				mediaType = "tvshow"
			}
			status := "success"
			if resp.IsErrorState() {
				status = "error"
				if resp.StatusCode == 429 {
					status = "rate_limited"
					observability.RecordMetadataRateLimited("omdb")
				}
			}
			observability.RecordMetadataFetch("omdb", mediaType, status, resp.TotalTime().Seconds())
			return nil
		})

	// Circuit breaker
	circuitbreaker.WrapReqClient(client, "omdb", circuitbreaker.TierExternal)

	return &Client{
		httpClient:  client,
		apiKey:      config.APIKey,
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

// GetByIMDbID fetches title details by IMDb ID (e.g., "tt0111161").
func (c *Client) GetByIMDbID(ctx context.Context, imdbID string) (*Response, error) {
	cacheKey := "imdb:" + imdbID
	if cached := c.getFromCache(cacheKey); cached != nil {
		if v, ok := cached.(*Response); ok {
			return v, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, err
	}

	var result Response
	resp, err := c.httpClient.R().SetContext(ctx).
		SetQueryParam("apikey", c.apiKey).
		SetQueryParam("i", imdbID).
		SetQueryParam("plot", "full").
		SetSuccessResult(&result).
		Get("/")
	if err != nil {
		return nil, fmt.Errorf("omdb request: %w", err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("omdb request: status %d", resp.StatusCode)
	}
	if result.Response == "False" {
		return nil, nil
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

// GetByTitle fetches title details by exact title match.
func (c *Client) GetByTitle(ctx context.Context, title string, year string, mediaType string) (*Response, error) {
	cacheKey := "title:" + title + ":" + year + ":" + mediaType
	if cached := c.getFromCache(cacheKey); cached != nil {
		if v, ok := cached.(*Response); ok {
			return v, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, err
	}

	r := c.httpClient.R().SetContext(ctx).
		SetQueryParam("apikey", c.apiKey).
		SetQueryParam("t", title).
		SetQueryParam("plot", "full")

	if year != "" {
		r.SetQueryParam("y", year)
	}
	if mediaType != "" {
		r.SetQueryParam("type", mediaType)
	}

	var result Response
	resp, err := r.SetSuccessResult(&result).Get("/")
	if err != nil {
		return nil, fmt.Errorf("omdb request: %w", err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("omdb request: status %d", resp.StatusCode)
	}
	if result.Response == "False" {
		return nil, nil
	}

	c.setCache(cacheKey, &result)
	return &result, nil
}

// Search searches for titles by query string.
func (c *Client) Search(ctx context.Context, query string, year string, mediaType string, page int) (*SearchResponse, error) {
	cacheKey := fmt.Sprintf("search:%s:%s:%s:%d", query, year, mediaType, page)
	if cached := c.getFromCache(cacheKey); cached != nil {
		if v, ok := cached.(*SearchResponse); ok {
			return v, nil
		}
	}

	if err := c.waitRateLimit(ctx); err != nil {
		return nil, err
	}

	r := c.httpClient.R().SetContext(ctx).
		SetQueryParam("apikey", c.apiKey).
		SetQueryParam("s", query)

	if year != "" {
		r.SetQueryParam("y", year)
	}
	if mediaType != "" {
		r.SetQueryParam("type", mediaType)
	}
	if page > 0 {
		r.SetQueryParam("page", fmt.Sprintf("%d", page))
	}

	var result SearchResponse
	resp, err := r.SetSuccessResult(&result).Get("/")
	if err != nil {
		return nil, fmt.Errorf("omdb search: %w", err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("omdb search: status %d", resp.StatusCode)
	}
	if result.Response == "False" {
		return nil, nil
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

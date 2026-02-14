package metadata

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/imroc/req/v3"
	"golang.org/x/time/rate"

	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/circuitbreaker"
)

// ClientConfig configures a metadata API client.
type ClientConfig struct {
	// BaseURL is the API base URL.
	BaseURL string

	// APIKey is the authentication key/token.
	APIKey string

	// RateLimit is requests per second (default: 4).
	RateLimit rate.Limit

	// RateBurst is the maximum burst size (default: 10).
	RateBurst int

	// CacheTTL is the cache duration (default: 24h).
	CacheTTL time.Duration

	// CacheMaxSize is the maximum number of cache entries (default: 10000).
	CacheMaxSize int

	// Timeout is the request timeout (default: 30s).
	Timeout time.Duration

	// RetryCount is the number of retries (default: 3).
	RetryCount int

	// ProxyURL for HTTP proxy (optional).
	ProxyURL string
}

// DefaultClientConfig returns a ClientConfig with sensible defaults.
func DefaultClientConfig() ClientConfig {
	return ClientConfig{
		RateLimit:    rate.Limit(4.0),
		RateBurst:    10,
		CacheTTL:     24 * time.Hour,
		CacheMaxSize: 10000,
		Timeout:      30 * time.Second,
		RetryCount:   3,
	}
}

// BaseClient provides shared HTTP client functionality for metadata providers.
// It includes rate limiting, caching, and retry logic.
type BaseClient struct {
	client      *req.Client
	apiKey      string
	rateLimiter *rate.Limiter
	cache       *cache.L1Cache[string, any]
	cacheTTL    time.Duration
	baseURL     string
}

// NewBaseClient creates a new BaseClient with the given configuration.
func NewBaseClient(config ClientConfig) *BaseClient {
	// Apply defaults
	if config.RateLimit == 0 {
		config.RateLimit = rate.Limit(4.0)
	}
	if config.RateBurst == 0 {
		config.RateBurst = 10
	}
	if config.CacheTTL == 0 {
		config.CacheTTL = 24 * time.Hour
	}
	if config.CacheMaxSize == 0 {
		config.CacheMaxSize = 10000
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.RetryCount == 0 {
		config.RetryCount = 3
	}

	client := req.C().
		SetBaseURL(config.BaseURL).
		SetTimeout(config.Timeout).
		SetCommonRetryCount(config.RetryCount).
		SetCommonRetryBackoffInterval(1*time.Second, 10*time.Second).
		SetCommonRetryCondition(func(resp *req.Response, err error) bool {
			if err != nil {
				return true
			}
			return resp.StatusCode >= 500
		})

	if config.ProxyURL != "" {
		client.SetProxyURL(config.ProxyURL)
	}

	// Circuit breaker (uses base URL as breaker name)
	breakerName := "metadata"
	if config.BaseURL != "" {
		// Extract a short name from the URL for metric labels
		if idx := strings.Index(config.BaseURL, "://"); idx >= 0 {
			host := config.BaseURL[idx+3:]
			if slashIdx := strings.Index(host, "/"); slashIdx >= 0 {
				host = host[:slashIdx]
			}
			breakerName = host
		}
	}
	circuitbreaker.WrapReqClient(client, breakerName, circuitbreaker.TierExternal)

	l1, err := cache.NewL1Cache[string, any](config.CacheMaxSize, config.CacheTTL, cache.WithExpiryAccessing[string, any]())
	if err != nil {
		// Fallback: create with defaults if custom config fails
		l1, _ = cache.NewL1Cache[string, any](cache.DefaultL1MaxSize, cache.DefaultL1TTL)
	}

	return &BaseClient{
		client:      client,
		apiKey:      config.APIKey,
		rateLimiter: rate.NewLimiter(config.RateLimit, config.RateBurst),
		cache:       l1,
		cacheTTL:    config.CacheTTL,
		baseURL:     config.BaseURL,
	}
}

// GetAPIKey returns the configured API key.
func (c *BaseClient) GetAPIKey() string {
	return c.apiKey
}

// GetBaseURL returns the configured base URL.
func (c *BaseClient) GetBaseURL() string {
	return c.baseURL
}

// WaitForRateLimit blocks until the rate limiter allows a request.
func (c *BaseClient) WaitForRateLimit(ctx context.Context) error {
	if err := c.rateLimiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limit wait: %w", err)
	}
	return nil
}

// GetFromCache retrieves a value from the cache if it exists.
func (c *BaseClient) GetFromCache(key string) any {
	if val, ok := c.cache.Get(key); ok {
		return val
	}
	return nil
}

// SetCache stores a value in the cache with the configured TTL.
func (c *BaseClient) SetCache(key string, data any) {
	c.cache.Set(key, data)
}

// SetCacheWithTTL stores a value in the cache.
// Note: L1Cache uses a fixed TTL set at creation. For different TTLs,
// use separate cache instances. This method exists for API compatibility.
func (c *BaseClient) SetCacheWithTTL(key string, data any, _ time.Duration) {
	c.cache.Set(key, data)
}

// ClearCache removes all entries from the cache.
func (c *BaseClient) ClearCache() {
	c.cache.Clear()
}

// Close closes the cache and stops background goroutines.
func (c *BaseClient) Close() {
	if c.cache != nil {
		c.cache.Close()
	}
}

// Request returns a new request builder configured with the base settings.
// The caller should set result/error types and make the actual request.
func (c *BaseClient) Request() *req.Request {
	return c.client.R()
}

// GetClient returns the underlying resty client for advanced use cases.
func (c *BaseClient) GetClient() *req.Client {
	return c.client
}

// CacheKey generates a cache key from components.
func CacheKey(parts ...any) string {
	var key strings.Builder
	for i, part := range parts {
		if i > 0 {
			key.WriteString(":")
		}
		key.WriteString(fmt.Sprintf("%v", part))
	}
	return key.String()
}

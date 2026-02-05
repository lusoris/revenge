package metadata

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/time/rate"
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
		RateLimit:  rate.Limit(4.0),
		RateBurst:  10,
		CacheTTL:   24 * time.Hour,
		Timeout:    30 * time.Second,
		RetryCount: 3,
	}
}

// BaseClient provides shared HTTP client functionality for metadata providers.
// It includes rate limiting, caching, and retry logic.
type BaseClient struct {
	client      *resty.Client
	apiKey      string
	rateLimiter *rate.Limiter
	cache       sync.Map
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
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.RetryCount == 0 {
		config.RetryCount = 3
	}

	client := resty.New().
		SetBaseURL(config.BaseURL).
		SetTimeout(config.Timeout).
		SetRetryCount(config.RetryCount).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second)

	if config.ProxyURL != "" {
		client.SetProxy(config.ProxyURL)
	}

	return &BaseClient{
		client:      client,
		apiKey:      config.APIKey,
		rateLimiter: rate.NewLimiter(config.RateLimit, config.RateBurst),
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

// GetFromCache retrieves a value from the cache if it exists and hasn't expired.
func (c *BaseClient) GetFromCache(key string) any {
	if val, ok := c.cache.Load(key); ok {
		entry, ok := val.(*CacheEntry)
		if !ok {
			return nil
		}
		if !entry.IsExpired() {
			return entry.Data
		}
		c.cache.Delete(key)
	}
	return nil
}

// SetCache stores a value in the cache with the configured TTL.
func (c *BaseClient) SetCache(key string, data any) {
	entry := &CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(c.cacheTTL),
	}
	c.cache.Store(key, entry)
}

// SetCacheWithTTL stores a value in the cache with a custom TTL.
func (c *BaseClient) SetCacheWithTTL(key string, data any, ttl time.Duration) {
	entry := &CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(ttl),
	}
	c.cache.Store(key, entry)
}

// ClearCache removes all entries from the cache.
func (c *BaseClient) ClearCache() {
	c.cache.Range(func(key, value any) bool {
		c.cache.Delete(key)
		return true
	})
}

// Request returns a new request builder configured with the base settings.
// The caller should set result/error types and make the actual request.
func (c *BaseClient) Request() *resty.Request {
	return c.client.R()
}

// GetClient returns the underlying resty client for advanced use cases.
func (c *BaseClient) GetClient() *resty.Client {
	return c.client
}

// CacheKey generates a cache key from components.
func CacheKey(parts ...any) string {
	key := ""
	for i, part := range parts {
		if i > 0 {
			key += ":"
		}
		key += fmt.Sprintf("%v", part)
	}
	return key
}

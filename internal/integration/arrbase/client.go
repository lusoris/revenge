package arrbase

import (
	"context"
	"fmt"
	"time"

	"github.com/imroc/req/v3"
	"golang.org/x/time/rate"

	"github.com/lusoris/revenge/internal/infra/cache"
)

// ClientConfig contains configuration for an arr API client.
type ClientConfig struct {
	BaseURL   string
	APIKey    string
	RateLimit rate.Limit    // requests per second (0 = default 10)
	CacheTTL  time.Duration // L1 cache TTL (0 = default 5 min)
	Timeout   time.Duration // HTTP timeout (0 = default 30s)
}

// BaseClient provides shared HTTP infrastructure for arr API clients.
// It encapsulates the HTTP client, rate limiter, and L1 cache that are
// identical between Radarr, Sonarr, and future arr integrations.
//
// Embed this in concrete clients:
//
//	type Client struct {
//	    *arrbase.BaseClient
//	}
type BaseClient struct {
	HTTP        *req.Client
	RateLimiter *rate.Limiter
	Cache       *cache.L1Cache[string, any]
	serviceName string // "radarr" or "sonarr" — for error messages
}

// NewBaseClient creates a new base arr API client with standard configuration.
func NewBaseClient(config ClientConfig, serviceName string) *BaseClient {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.RateLimit == 0 {
		config.RateLimit = rate.Limit(10.0) // 10 req/s default for local service
	}
	if config.CacheTTL == 0 {
		config.CacheTTL = 5 * time.Minute
	}

	client := req.C().
		SetBaseURL(config.BaseURL+"/api/v3").
		SetTimeout(config.Timeout).
		SetCommonHeader("X-Api-Key", config.APIKey).
		SetCommonHeader("Content-Type", "application/json").
		SetCommonRetryCount(3).
		SetCommonRetryBackoffInterval(1*time.Second, 10*time.Second).
		SetCommonRetryCondition(func(resp *req.Response, err error) bool {
			if err != nil {
				return true
			}
			return resp.StatusCode >= 500
		})

	l1, err := cache.NewL1Cache[string, any](1000, config.CacheTTL, cache.WithExpiryAccessing[string, any]())
	if err != nil {
		l1, _ = cache.NewL1Cache[string, any](0, 0)
	}

	return &BaseClient{
		HTTP:        client,
		RateLimiter: rate.NewLimiter(config.RateLimit, 20),
		Cache:       l1,
		serviceName: serviceName,
	}
}

// CachedGet performs a GET request with L1 cache and rate limiting.
// This encapsulates the cache → rate-limit → HTTP → cache pattern
// that is repeated for every arr API method.
func CachedGet[T any](ctx context.Context, c *BaseClient, cacheKey string, path string) (*T, error) {
	// Check cache
	if val, ok := c.Cache.Get(cacheKey); ok {
		if result, ok := val.(*T); ok {
			return result, nil
		}
	}

	// Rate limit
	if err := c.RateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	// HTTP request
	var result T
	resp, err := c.HTTP.R().
		SetContext(ctx).
		SetSuccessResult(&result).
		Get(path)

	if err != nil {
		return nil, fmt.Errorf("%s api request: %w", c.serviceName, err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("%s api error: %s", c.serviceName, resp.Status)
	}

	// Store in cache
	c.Cache.Set(cacheKey, &result)
	return &result, nil
}

// CachedGetList performs a GET request returning a slice, with L1 cache and rate limiting.
func CachedGetList[T any](ctx context.Context, c *BaseClient, cacheKey string, path string) ([]T, error) {
	// Check cache
	if val, ok := c.Cache.Get(cacheKey); ok {
		if result, ok := val.([]T); ok {
			return result, nil
		}
	}

	// Rate limit
	if err := c.RateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	// HTTP request
	var result []T
	resp, err := c.HTTP.R().
		SetContext(ctx).
		SetSuccessResult(&result).
		Get(path)

	if err != nil {
		return nil, fmt.Errorf("%s api request: %w", c.serviceName, err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("%s api error: %s", c.serviceName, resp.Status)
	}

	// Store in cache
	c.Cache.Set(cacheKey, result)
	return result, nil
}

// PostCommand sends a command POST request with rate limiting (no cache).
func PostCommand[T any](ctx context.Context, c *BaseClient, path string, body any) (*T, error) {
	if err := c.RateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait: %w", err)
	}

	var result T
	resp, err := c.HTTP.R().
		SetContext(ctx).
		SetBody(body).
		SetSuccessResult(&result).
		Post(path)

	if err != nil {
		return nil, fmt.Errorf("%s api request: %w", c.serviceName, err)
	}
	if resp.IsErrorState() {
		return nil, fmt.Errorf("%s api error: %s", c.serviceName, resp.Status)
	}

	return &result, nil
}

// Delete performs a DELETE request with rate limiting (no cache).
func Delete(ctx context.Context, c *BaseClient, path string) error {
	if err := c.RateLimiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limit wait: %w", err)
	}

	resp, err := c.HTTP.R().
		SetContext(ctx).
		Delete(path)

	if err != nil {
		return fmt.Errorf("%s api request: %w", c.serviceName, err)
	}
	if resp.IsErrorState() {
		return fmt.Errorf("%s api error: %s", c.serviceName, resp.Status)
	}

	return nil
}

// ClearCache clears all entries from the L1 cache.
func (c *BaseClient) ClearCache() {
	c.Cache.Clear()
}

// Close closes the HTTP client and clears the cache.
func (c *BaseClient) Close() {
	c.Cache.Clear()
	c.HTTP.CloseIdleConnections()
}

// IsHealthy checks if the arr service is reachable.
func (c *BaseClient) IsHealthy(ctx context.Context) bool {
	_, err := CachedGet[SystemStatus](ctx, c, "system:status", "/system/status")
	return err == nil
}

// GetSystemStatus returns the system status of the arr instance.
func (c *BaseClient) GetSystemStatus(ctx context.Context) (*SystemStatus, error) {
	return CachedGet[SystemStatus](ctx, c, "system:status", "/system/status")
}

// GetQualityProfiles returns quality profiles from the arr instance.
func (c *BaseClient) GetQualityProfiles(ctx context.Context) ([]QualityProfile, error) {
	return CachedGetList[QualityProfile](ctx, c, "qualityprofiles:all", "/qualityprofile")
}

// GetRootFolders returns root folders from the arr instance.
func (c *BaseClient) GetRootFolders(ctx context.Context) ([]RootFolder, error) {
	return CachedGetList[RootFolder](ctx, c, "rootfolders:all", "/rootfolder")
}

// GetTags returns tags from the arr instance.
func (c *BaseClient) GetTags(ctx context.Context) ([]Tag, error) {
	return CachedGetList[Tag](ctx, c, "tags:all", "/tag")
}


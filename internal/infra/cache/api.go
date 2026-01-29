package cache

import (
	"context"
	"log/slog"
	"time"

	"github.com/viccon/sturdyc"

	"github.com/lusoris/revenge/internal/config"
)

// APICache provides request coalescing and stale-while-revalidate for API responses.
// This is Tier 3 of the caching hierarchy - optimized for external API calls.
//
// sturdyc provides:
// - Request coalescing: multiple concurrent requests for the same key result in single fetch
// - Stale-while-revalidate: serves stale data while refreshing in background
// - Batch fetching: groups multiple keys into single batch request
// - Perfect for rate-limited APIs (TMDb, MusicBrainz, etc.)
type APICache struct {
	cache  *sturdyc.Client[any]
	logger *slog.Logger
}

// APICacheConfig holds configuration for the API cache.
type APICacheConfig struct {
	// Capacity is the maximum number of entries.
	Capacity int
	// NumShards is the number of shards for concurrent access.
	NumShards int
	// TTL is the time-to-live for entries.
	TTL time.Duration
	// EvictionPercentage is the percentage of entries to evict when full.
	EvictionPercentage int
}

// DefaultAPICacheConfig returns sensible defaults for API cache.
func DefaultAPICacheConfig() APICacheConfig {
	return APICacheConfig{
		Capacity:           10_000,
		NumShards:          32,
		TTL:                1 * time.Hour,
		EvictionPercentage: 10,
	}
}

// NewAPICache creates a new API response cache using sturdyc.
func NewAPICache(cfg *config.Config, logger *slog.Logger) *APICache {
	apiCfg := DefaultAPICacheConfig()

	cache := sturdyc.New[any](
		apiCfg.Capacity,
		apiCfg.NumShards,
		apiCfg.TTL,
		apiCfg.EvictionPercentage,
	)

	return &APICache{
		cache:  cache,
		logger: logger.With(slog.String("component", "cache.api")),
	}
}

// GetOrFetch retrieves a value from cache or fetches it using the provided function.
// Multiple concurrent requests for the same key will be coalesced into a single fetch.
func (c *APICache) GetOrFetch(ctx context.Context, key string, fetch func(ctx context.Context) (any, error)) (any, error) {
	return c.cache.GetOrFetch(ctx, key, fetch)
}

// Set manually sets a value in the cache.
func (c *APICache) Set(key string, value any) {
	c.cache.Set(key, value)
}

// Get retrieves a value from the cache without fetching.
func (c *APICache) Get(key string) (any, bool) {
	return c.cache.Get(key)
}

// Delete removes a key from the cache.
func (c *APICache) Delete(key string) {
	c.cache.Delete(key)
}

// Size returns the number of entries in the cache.
func (c *APICache) Size() int {
	return c.cache.Size()
}

// MetadataCache is a specialized API cache for metadata providers.
// It includes provider-specific configurations for rate limiting.
type MetadataCache struct {
	*APICache
	providerConfigs map[string]ProviderConfig
}

// ProviderConfig holds provider-specific cache configuration.
type ProviderConfig struct {
	// Name is the provider name (e.g., "tmdb", "musicbrainz").
	Name string
	// TTL is the cache TTL for this provider.
	TTL time.Duration
	// RateLimit is the maximum requests per second (0 = unlimited).
	RateLimit float64
}

// DefaultProviderConfigs returns default configurations for known metadata providers.
func DefaultProviderConfigs() map[string]ProviderConfig {
	return map[string]ProviderConfig{
		"tmdb": {
			Name:      "tmdb",
			TTL:       24 * time.Hour,
			RateLimit: 4, // 40 req/10 sec
		},
		"musicbrainz": {
			Name:      "musicbrainz",
			TTL:       24 * time.Hour,
			RateLimit: 1, // 1 req/sec STRICT
		},
		"thetvdb": {
			Name:      "thetvdb",
			TTL:       24 * time.Hour,
			RateLimit: 0, // No published limit
		},
		"fanart": {
			Name:      "fanart",
			TTL:       7 * 24 * time.Hour, // Images don't change often
			RateLimit: 0,
		},
		"lastfm": {
			Name:      "lastfm",
			TTL:       1 * time.Hour,
			RateLimit: 5, // Reasonable usage
		},
		"stashdb": {
			Name:      "stashdb",
			TTL:       24 * time.Hour,
			RateLimit: 2, // Be respectful
		},
		"tpdb": {
			Name:      "tpdb",
			TTL:       24 * time.Hour,
			RateLimit: 2,
		},
	}
}

// NewMetadataCache creates a cache specifically for metadata providers.
func NewMetadataCache(cfg *config.Config, logger *slog.Logger) *MetadataCache {
	return &MetadataCache{
		APICache:        NewAPICache(cfg, logger),
		providerConfigs: DefaultProviderConfigs(),
	}
}

// GetProviderConfig returns the configuration for a provider.
func (c *MetadataCache) GetProviderConfig(provider string) (ProviderConfig, bool) {
	cfg, ok := c.providerConfigs[provider]
	return cfg, ok
}

// CacheKey generates a cache key for a provider and resource.
func CacheKey(provider, resourceType, id string) string {
	return "meta:" + provider + ":" + resourceType + ":" + id
}

// CacheKeyWithLang generates a cache key including language.
func CacheKeyWithLang(provider, resourceType, id, lang string) string {
	return "meta:" + provider + ":" + resourceType + ":" + id + ":" + lang
}

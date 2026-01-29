package cache

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/maypok86/otter"

	"github.com/lusoris/revenge/internal/config"
)

// LocalCache provides a high-performance in-memory cache using otter.
// This is Tier 1 of the caching hierarchy - microsecond latency.
//
// otter uses W-TinyLFU eviction policy which provides:
// - Better hit rates than LRU
// - 50% less memory than ristretto
// - Lock-free reads
type LocalCache struct {
	cache  otter.Cache[string, cacheEntry]
	logger *slog.Logger
	ttl    time.Duration
}

// cacheEntry wraps a value with its expiration time.
type cacheEntry struct {
	Data      []byte
	ExpiresAt time.Time
}

// isExpired checks if the entry has expired.
func (e cacheEntry) isExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// LocalCacheConfig holds configuration for the local cache.
type LocalCacheConfig struct {
	// MaxSize is the maximum number of entries in the cache.
	MaxSize int
	// DefaultTTL is the default time-to-live for entries.
	DefaultTTL time.Duration
}

// DefaultLocalCacheConfig returns sensible defaults for local cache.
func DefaultLocalCacheConfig() LocalCacheConfig {
	return LocalCacheConfig{
		MaxSize:    100_000, // 100k entries
		DefaultTTL: 5 * time.Minute,
	}
}

// NewLocalCache creates a new local in-memory cache using otter.
func NewLocalCache(cfg *config.Config, logger *slog.Logger) (*LocalCache, error) {
	localCfg := DefaultLocalCacheConfig()

	// Override defaults with config values if set
	if cfg.Cache.LocalCapacity > 0 {
		localCfg.MaxSize = cfg.Cache.LocalCapacity
	}
	if cfg.Cache.LocalTTL > 0 {
		localCfg.DefaultTTL = time.Duration(cfg.Cache.LocalTTL) * time.Second
	}

	cache, err := otter.MustBuilder[string, cacheEntry](localCfg.MaxSize).
		Build()
	if err != nil {
		return nil, err
	}

	return &LocalCache{
		cache:  cache,
		logger: logger.With(slog.String("component", "cache.local")),
		ttl:    localCfg.DefaultTTL,
	}, nil
}

// defaultTTL returns the default TTL.
func (c *LocalCache) defaultTTL() time.Duration {
	if c.ttl > 0 {
		return c.ttl
	}
	return 5 * time.Minute
}

// Get retrieves a value from the local cache.
func (c *LocalCache) Get(key string) ([]byte, bool) {
	entry, ok := c.cache.Get(key)
	if !ok || entry.isExpired() {
		if ok {
			c.cache.Delete(key) // Clean up expired entry
		}
		return nil, false
	}
	return entry.Data, true
}

// GetString retrieves a string value from the local cache.
func (c *LocalCache) GetString(key string) (string, bool) {
	data, ok := c.Get(key)
	if !ok {
		return "", false
	}
	return string(data), true
}

// GetJSON retrieves and unmarshals a JSON value from the local cache.
func (c *LocalCache) GetJSON(key string, dest any) bool {
	data, ok := c.Get(key)
	if !ok {
		return false
	}
	if err := json.Unmarshal(data, dest); err != nil {
		c.logger.Warn("failed to unmarshal cached json",
			slog.String("key", key),
			slog.Any("error", err))
		return false
	}
	return true
}

// Set stores a value in the local cache with the default TTL.
func (c *LocalCache) Set(key string, value []byte) {
	c.SetWithTTL(key, value, c.defaultTTL())
}

// SetWithTTL stores a value in the local cache with a custom TTL.
func (c *LocalCache) SetWithTTL(key string, value []byte, ttl time.Duration) {
	entry := cacheEntry{
		Data:      value,
		ExpiresAt: time.Now().Add(ttl),
	}
	c.cache.Set(key, entry)
}

// SetString stores a string value in the local cache.
func (c *LocalCache) SetString(key string, value string) {
	c.Set(key, []byte(value))
}

// SetStringWithTTL stores a string value with a custom TTL.
func (c *LocalCache) SetStringWithTTL(key string, value string, ttl time.Duration) {
	c.SetWithTTL(key, []byte(value), ttl)
}

// SetJSON stores a JSON-serializable value in the local cache.
func (c *LocalCache) SetJSON(key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	c.Set(key, data)
	return nil
}

// SetJSONWithTTL stores a JSON value with a custom TTL.
func (c *LocalCache) SetJSONWithTTL(key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	c.SetWithTTL(key, data, ttl)
	return nil
}

// Delete removes a key from the local cache.
func (c *LocalCache) Delete(key string) {
	c.cache.Delete(key)
}

// Has checks if a key exists in the local cache (and is not expired).
func (c *LocalCache) Has(key string) bool {
	entry, ok := c.cache.Get(key)
	if !ok || entry.isExpired() {
		return false
	}
	return true
}

// Clear removes all entries from the local cache.
func (c *LocalCache) Clear() {
	c.cache.Clear()
}

// Size returns the number of entries in the cache.
func (c *LocalCache) Size() int {
	return c.cache.Size()
}

// Close releases resources used by the cache.
func (c *LocalCache) Close() {
	c.cache.Close()
}

// Stats returns cache statistics.
func (c *LocalCache) Stats() LocalCacheStats {
	return LocalCacheStats{
		Size:   c.cache.Size(),
		Hits:   c.cache.Stats().Hits(),
		Misses: c.cache.Stats().Misses(),
	}
}

// LocalCacheStats holds cache statistics.
type LocalCacheStats struct {
	Size   int
	Hits   int64
	Misses int64
}

// HitRate returns the cache hit rate as a percentage.
func (s LocalCacheStats) HitRate() float64 {
	total := s.Hits + s.Misses
	if total == 0 {
		return 0
	}
	return float64(s.Hits) / float64(total) * 100
}

// GetOrSet retrieves a value from cache, or computes and stores it if not present.
// This is useful for cache-aside pattern.
func (c *LocalCache) GetOrSet(ctx context.Context, key string, ttl time.Duration, fetch func(context.Context) ([]byte, error)) ([]byte, error) {
	// Check cache first
	if data, ok := c.Get(key); ok {
		return data, nil
	}

	// Fetch from source
	data, err := fetch(ctx)
	if err != nil {
		return nil, err
	}

	// Store in cache
	c.SetWithTTL(key, data, ttl)

	return data, nil
}

// GetOrSetJSON is like GetOrSet but for JSON values.
func (c *LocalCache) GetOrSetJSON(ctx context.Context, key string, ttl time.Duration, dest any, fetch func(context.Context) (any, error)) error {
	// Check cache first
	if c.GetJSON(key, dest) {
		return nil
	}

	// Fetch from source
	value, err := fetch(ctx)
	if err != nil {
		return err
	}

	// Store in cache
	if err := c.SetJSONWithTTL(key, value, ttl); err != nil {
		return err
	}

	// Copy to destination
	data, _ := json.Marshal(value)
	return json.Unmarshal(data, dest)
}

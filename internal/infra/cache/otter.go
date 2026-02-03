package cache

import (
	"time"

	"github.com/maypok86/otter/v2"
)

const (
	// DefaultL1MaxSize is the default maximum number of entries in L1 cache.
	DefaultL1MaxSize = 10000

	// DefaultL1TTL is the default TTL for L1 cache entries.
	DefaultL1TTL = 5 * time.Minute
)

// L1Cache wraps otter cache for in-memory L1 caching using W-TinyLFU eviction.
//
// IMPORTANT: Evictions are performed asynchronously in background goroutines
// for performance reasons. The cache size may temporarily exceed MaximumSize
// during eviction processing, typically settling within milliseconds.
//
// When configuring memory limits, plan for ~20% headroom above the max size
// to account for this async behavior. Monitor cache size in production to
// verify memory usage stays within acceptable bounds.
type L1Cache[K comparable, V any] struct {
	cache *otter.Cache[K, V]
}

// NewL1Cache creates a new L1 in-memory cache with W-TinyLFU eviction policy.
func NewL1Cache[K comparable, V any](maxSize int, ttl time.Duration) (*L1Cache[K, V], error) {
	if maxSize <= 0 {
		maxSize = DefaultL1MaxSize
	}
	if ttl <= 0 {
		ttl = DefaultL1TTL
	}

	cache, err := otter.New(&otter.Options[K, V]{
		MaximumSize:      maxSize,
		ExpiryCalculator: otter.ExpiryWriting[K, V](ttl),
	})
	if err != nil {
		return nil, err
	}

	return &L1Cache[K, V]{
		cache: cache,
	}, nil
}

// Get retrieves a value from the cache.
func (l *L1Cache[K, V]) Get(key K) (V, bool) {
	return l.cache.GetIfPresent(key)
}

// Set stores a value in the cache.
func (l *L1Cache[K, V]) Set(key K, value V) {
	l.cache.Set(key, value)
}

// Delete removes a value from the cache.
func (l *L1Cache[K, V]) Delete(key K) {
	l.cache.Invalidate(key)
}

// Clear removes all entries from the cache.
func (l *L1Cache[K, V]) Clear() {
	l.cache.InvalidateAll()
}

// Size returns the estimated number of entries in the cache.
func (l *L1Cache[K, V]) Size() int {
	return l.cache.EstimatedSize()
}

// Close closes the cache and stops all background goroutines.
func (l *L1Cache[K, V]) Close() {
	l.cache.StopAllGoroutines()
}

// Has checks if a key exists in the cache without updating access time.
func (l *L1Cache[K, V]) Has(key K) bool {
	_, ok := l.cache.GetIfPresent(key)
	return ok
}

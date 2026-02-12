package cache

import (
	"context"
	"time"
)

// Get implements the type-safe cache-aside pattern using generics.
// It first checks the cache (via GetJSON), and on miss calls the loader
// function, stores the result asynchronously, and returns it.
//
// If c is nil, the loader is called directly (no caching).
//
// Usage:
//
//	movie, err := cache.Get(ctx, s.cache, cacheKey, cache.MovieTTL, func(ctx context.Context) (*Movie, error) {
//	    return s.Service.GetMovie(ctx, id)
//	})
func Get[T any](ctx context.Context, c *Cache, key string, ttl time.Duration, loader func(context.Context) (T, error)) (T, error) {
	var zero T
	if c == nil {
		return loader(ctx)
	}

	// Try cache first
	var result T
	if err := c.GetJSON(ctx, key, &result); err == nil {
		return result, nil
	}

	// Cache miss — load from source
	result, err := loader(ctx)
	if err != nil {
		return zero, err
	}

	// Store in cache asynchronously (fire-and-forget)
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		_ = c.SetJSON(cacheCtx, key, result, ttl)
	}()

	return result, nil
}

// Pair holds items together with a total count for paginated cache results.
// Use with [Get] to cache paginated queries that return (items, total, error).
//
// Usage:
//
//	result, err := cache.Get(ctx, s.cache, key, ttl, func(ctx context.Context) (cache.Pair[[]Movie], error) {
//	    items, total, err := s.Service.ListRecentlyAdded(ctx, limit, offset)
//	    return cache.Pair[[]Movie]{Items: items, Total: total}, err
//	})
//	return result.Items, result.Total, err
type Pair[T any] struct {
	Items T     `json:"items"`
	Total int64 `json:"total"`
}

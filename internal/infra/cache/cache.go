package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/redis/rueidis"

	"github.com/lusoris/revenge/internal/infra/observability"
)

// Cache provides unified L1 (otter) + L2 (rueidis) caching operations.
type Cache struct {
	l1     *L1Cache[string, []byte]
	l1TTL  time.Duration
	client *Client
	name   string // cache name for metrics
}

// NewCache creates a new unified cache with L1 and L2 layers.
func NewCache(client *Client, l1MaxSize int, l1TTL time.Duration) (*Cache, error) {
	return NewNamedCache(client, l1MaxSize, l1TTL, "default")
}

// NewNamedCache creates a new unified cache with a specific name for metrics.
func NewNamedCache(client *Client, l1MaxSize int, l1TTL time.Duration, name string) (*Cache, error) {
	l1, err := NewL1Cache[string, []byte](l1MaxSize, l1TTL)
	if err != nil {
		return nil, fmt.Errorf("failed to create L1 cache: %w", err)
	}

	return &Cache{
		l1:     l1,
		l1TTL:  l1TTL,
		client: client,
		name:   name,
	}, nil
}

// Get retrieves a value from cache (L1 first, then L2).
func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	start := time.Now()
	defer func() {
		observability.CacheOperationDuration.WithLabelValues(c.name, "get").Observe(time.Since(start).Seconds())
	}()

	// Try L1 first
	if val, ok := c.l1.Get(key); ok {
		observability.RecordCacheHit(c.name, "l1")
		return val, nil
	}
	observability.RecordCacheMiss(c.name, "l1")

	// L1 miss - try L2 (rueidis)
	if c.client == nil || c.client.rueidisClient == nil {
		observability.RecordCacheMiss(c.name, "l2")
		return nil, fmt.Errorf("cache miss: key not found in L1 and L2 unavailable")
	}

	cmd := c.client.rueidisClient.B().Get().Key(key).Build()
	resp := c.client.rueidisClient.Do(ctx, cmd)

	if err := resp.Error(); err != nil {
		observability.RecordCacheMiss(c.name, "l2")
		return nil, fmt.Errorf("L2 cache get failed: %w", err)
	}

	val, err := resp.AsBytes()
	if err != nil {
		observability.RecordCacheMiss(c.name, "l2")
		return nil, fmt.Errorf("L2 cache value decode failed: %w", err)
	}

	observability.RecordCacheHit(c.name, "l2")

	// Populate L1 on L2 hit
	c.l1.Set(key, val)

	return val, nil
}

// Set stores a value in both L1 and L2 caches.
// For TTLs shorter than L1's TTL, skips L1 to ensure accurate expiration.
func (c *Cache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	start := time.Now()
	defer func() {
		observability.CacheOperationDuration.WithLabelValues(c.name, "set").Observe(time.Since(start).Seconds())
	}()

	// Only use L1 if TTL is longer than or equal to L1's TTL
	// This ensures short-lived items expire accurately in L2
	if ttl == 0 || ttl >= c.l1TTL {
		c.l1.Set(key, value)
		observability.CacheSize.WithLabelValues(c.name).Set(float64(c.l1.Size()))
	} else {
		// For short TTLs, remove from L1 to prevent stale reads
		c.l1.Delete(key)
	}

	// Store in L2 if available
	if c.client != nil && c.client.rueidisClient != nil {
		// Use PX (milliseconds) for sub-second precision, EX (seconds) otherwise
		if ttl < time.Second && ttl > 0 {
			// Use PX for millisecond precision
			cmd := c.client.rueidisClient.B().Set().Key(key).Value(string(value)).Px(ttl).Build()
			if err := c.client.rueidisClient.Do(ctx, cmd).Error(); err != nil {
				return fmt.Errorf("L2 cache set failed: %w", err)
			}
		} else if ttl > 0 {
			// Use EX for second precision
			cmd := c.client.rueidisClient.B().Set().Key(key).Value(string(value)).Ex(ttl).Build()
			if err := c.client.rueidisClient.Do(ctx, cmd).Error(); err != nil {
				return fmt.Errorf("L2 cache set failed: %w", err)
			}
		} else {
			// No expiration
			cmd := c.client.rueidisClient.B().Set().Key(key).Value(string(value)).Build()
			if err := c.client.rueidisClient.Do(ctx, cmd).Error(); err != nil {
				return fmt.Errorf("L2 cache set failed: %w", err)
			}
		}
	}

	return nil
}

// Delete removes a value from both L1 and L2 caches.
func (c *Cache) Delete(ctx context.Context, key string) error {
	start := time.Now()
	defer func() {
		observability.CacheOperationDuration.WithLabelValues(c.name, "delete").Observe(time.Since(start).Seconds())
	}()

	// Delete from L1
	c.l1.Delete(key)
	observability.CacheSize.WithLabelValues(c.name).Set(float64(c.l1.Size()))

	// Delete from L2 if available
	if c.client != nil && c.client.rueidisClient != nil {
		cmd := c.client.rueidisClient.B().Del().Key(key).Build()
		if err := c.client.rueidisClient.Do(ctx, cmd).Error(); err != nil {
			return fmt.Errorf("L2 cache delete failed: %w", err)
		}
	}

	return nil
}

// Exists checks if a key exists in L1 or L2 cache.
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	// Check L1 first
	if c.l1.Has(key) {
		return true, nil
	}

	// Check L2 if available
	if c.client != nil && c.client.rueidisClient != nil {
		cmd := c.client.rueidisClient.B().Exists().Key(key).Build()
		resp := c.client.rueidisClient.Do(ctx, cmd)

		if err := resp.Error(); err != nil {
			return false, fmt.Errorf("L2 cache exists check failed: %w", err)
		}

		count, err := resp.AsInt64()
		if err != nil {
			return false, fmt.Errorf("L2 cache exists response decode failed: %w", err)
		}

		return count > 0, nil
	}

	return false, nil
}

// Invalidate removes all keys matching a pattern from both cache layers.
//
// L1: If the pattern is a simple prefix glob (e.g. "movie:*"), only matching
// keys are evicted. Otherwise L1 is fully cleared as a safe fallback.
//
// L2: Uses SCAN (non-blocking, cursor-based) instead of KEYS to avoid
// blocking Redis on large keyspaces. Matching keys are deleted in batches.
func (c *Cache) Invalidate(ctx context.Context, pattern string) error {
	// L1: targeted prefix delete when pattern is "prefix*"
	if prefix, ok := simpleGlobPrefix(pattern); ok {
		deleted := c.l1.DeleteByPrefix(prefix)
		observability.CacheSize.WithLabelValues(c.name).Set(float64(c.l1.Size()))
		_ = deleted
	} else {
		c.l1.Clear()
		observability.CacheSize.WithLabelValues(c.name).Set(0)
	}

	// L2: use SCAN to find matching keys (non-blocking)
	if c.client != nil && c.client.rueidisClient != nil {
		client := c.client.rueidisClient

		scanner := rueidis.NewScanner(func(cursor uint64) (rueidis.ScanEntry, error) {
			cmd := client.B().Scan().Cursor(cursor).Match(pattern).Count(100).Build()
			return client.Do(ctx, cmd).AsScanEntry()
		})

		var batch []string
		for key := range scanner.Iter() {
			batch = append(batch, key)
			if len(batch) >= 100 {
				if err := c.deleteBatch(ctx, batch); err != nil {
					return err
				}
				batch = batch[:0]
			}
		}
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("L2 cache SCAN failed: %w", err)
		}
		if len(batch) > 0 {
			if err := c.deleteBatch(ctx, batch); err != nil {
				return err
			}
		}
	}

	return nil
}

// deleteBatch deletes a slice of keys from L2.
func (c *Cache) deleteBatch(ctx context.Context, keys []string) error {
	cmd := c.client.rueidisClient.B().Del().Key(keys...).Build()
	if err := c.client.rueidisClient.Do(ctx, cmd).Error(); err != nil {
		return fmt.Errorf("L2 cache batch delete failed: %w", err)
	}
	return nil
}

// simpleGlobPrefix checks if pattern is a simple "prefix*" glob (no other
// wildcards or special chars). Returns the prefix and true if so.
func simpleGlobPrefix(pattern string) (string, bool) {
	if !strings.HasSuffix(pattern, "*") {
		return "", false
	}
	prefix := strings.TrimSuffix(pattern, "*")
	// Reject patterns with wildcards or special glob chars in the prefix
	if strings.ContainsAny(prefix, "*?[\\") {
		return "", false
	}
	return prefix, true
}

// Close closes both L1 and L2 caches.
func (c *Cache) Close() {
	if c.l1 != nil {
		c.l1.Close()
	}
}

// GetJSON retrieves and unmarshals a JSON value from cache.
func (c *Cache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	data, err := c.Get(ctx, key)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal cached JSON: %w", err)
	}

	return nil
}

// SetJSON marshals and stores a JSON value in cache.
func (c *Cache) SetJSON(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value to JSON: %w", err)
	}

	return c.Set(ctx, key, data, ttl)
}

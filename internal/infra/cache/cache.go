// Package cache provides a three-tier caching system:
// - Tier 1: Local in-memory cache (otter) - microsecond latency
// - Tier 2: Distributed cache (Dragonfly via rueidis) - millisecond latency
// - Tier 3: API response cache (sturdyc) - request coalescing
package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/rueidis"
	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/config"
)

var (
	// ErrCacheMiss is returned when a key is not found in cache.
	ErrCacheMiss = errors.New("cache miss")
)

// Client wraps the rueidis client for distributed caching operations.
// This is Tier 2 of the caching hierarchy - Dragonfly/Redis compatible.
type Client struct {
	client rueidis.Client
	logger *slog.Logger
}

// NewClient creates a new distributed cache client using rueidis.
// rueidis provides 14x better performance than go-redis through:
// - Auto-pipelining (batches commands automatically)
// - Server-assisted client caching (RESP3)
// - Zero-allocation design
func NewClient(cfg *config.Config, logger *slog.Logger) (*Client, error) {
	opts := rueidis.ClientOption{
		InitAddress: []string{cfg.Cache.Addr},
		Password:    cfg.Cache.Password,
		SelectDB:    cfg.Cache.DB,
		// Enable client-side caching for frequently accessed keys
		ClientTrackingOptions: []string{"OPTIN"},
	}

	client, err := rueidis.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("create rueidis client: %w", err)
	}

	return &Client{
		client: client,
		logger: logger.With(slog.String("component", "cache.distributed")),
	}, nil
}

// Ping checks if the cache is reachable.
func (c *Client) Ping(ctx context.Context) error {
	cmd := c.client.B().Ping().Build()
	return c.client.Do(ctx, cmd).Error()
}

// Close closes the cache connection.
func (c *Client) Close() {
	c.client.Close()
}

// Get retrieves a string value from cache.
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	cmd := c.client.B().Get().Key(key).Build()
	result, err := c.client.Do(ctx, cmd).ToString()
	if err != nil {
		if rueidis.IsRedisNil(err) {
			return "", ErrCacheMiss
		}
		return "", fmt.Errorf("cache get: %w", err)
	}
	return result, nil
}

// GetBytes retrieves bytes from cache.
func (c *Client) GetBytes(ctx context.Context, key string) ([]byte, error) {
	cmd := c.client.B().Get().Key(key).Build()
	result, err := c.client.Do(ctx, cmd).AsBytes()
	if err != nil {
		if rueidis.IsRedisNil(err) {
			return nil, ErrCacheMiss
		}
		return nil, fmt.Errorf("cache get bytes: %w", err)
	}
	return result, nil
}

// Set stores a string value in cache with expiration.
func (c *Client) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	cmd := c.client.B().Set().Key(key).Value(value).Ex(expiration).Build()
	return c.client.Do(ctx, cmd).Error()
}

// SetBytes stores bytes in cache with expiration.
func (c *Client) SetBytes(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	cmd := c.client.B().Set().Key(key).Value(string(value)).Ex(expiration).Build()
	return c.client.Do(ctx, cmd).Error()
}

// SetJSON stores a JSON-serializable value in cache.
func (c *Client) SetJSON(ctx context.Context, key string, value any, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}
	return c.SetBytes(ctx, key, data, expiration)
}

// GetJSON retrieves and unmarshals a JSON value from cache.
func (c *Client) GetJSON(ctx context.Context, key string, dest any) error {
	data, err := c.GetBytes(ctx, key)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("unmarshal json: %w", err)
	}
	return nil
}

// Delete removes keys from cache.
func (c *Client) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	cmd := c.client.B().Del().Key(keys...).Build()
	return c.client.Do(ctx, cmd).Error()
}

// Exists checks if a key exists in cache.
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	cmd := c.client.B().Exists().Key(key).Build()
	count, err := c.client.Do(ctx, cmd).AsInt64()
	if err != nil {
		return false, fmt.Errorf("cache exists: %w", err)
	}
	return count > 0, nil
}

// Expire sets a new expiration on a key.
func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	cmd := c.client.B().Expire().Key(key).Seconds(int64(expiration.Seconds())).Build()
	return c.client.Do(ctx, cmd).Error()
}

// TTL returns the remaining time-to-live for a key.
func (c *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	cmd := c.client.B().Ttl().Key(key).Build()
	seconds, err := c.client.Do(ctx, cmd).AsInt64()
	if err != nil {
		return 0, fmt.Errorf("cache ttl: %w", err)
	}
	if seconds < 0 {
		return 0, ErrCacheMiss
	}
	return time.Duration(seconds) * time.Second, nil
}

// Incr increments a counter and returns the new value.
func (c *Client) Incr(ctx context.Context, key string) (int64, error) {
	cmd := c.client.B().Incr().Key(key).Build()
	return c.client.Do(ctx, cmd).AsInt64()
}

// IncrBy increments a counter by a specific amount.
func (c *Client) IncrBy(ctx context.Context, key string, delta int64) (int64, error) {
	cmd := c.client.B().Incrby().Key(key).Increment(delta).Build()
	return c.client.Do(ctx, cmd).AsInt64()
}

// HSet sets a hash field.
func (c *Client) HSet(ctx context.Context, key, field string, value string) error {
	cmd := c.client.B().Hset().Key(key).FieldValue().FieldValue(field, value).Build()
	return c.client.Do(ctx, cmd).Error()
}

// HGet retrieves a hash field.
func (c *Client) HGet(ctx context.Context, key, field string) (string, error) {
	cmd := c.client.B().Hget().Key(key).Field(field).Build()
	result, err := c.client.Do(ctx, cmd).ToString()
	if err != nil {
		if rueidis.IsRedisNil(err) {
			return "", ErrCacheMiss
		}
		return "", fmt.Errorf("cache hget: %w", err)
	}
	return result, nil
}

// HGetAll retrieves all fields of a hash.
func (c *Client) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	cmd := c.client.B().Hgetall().Key(key).Build()
	return c.client.Do(ctx, cmd).AsStrMap()
}

// HDel removes hash fields.
func (c *Client) HDel(ctx context.Context, key string, fields ...string) error {
	if len(fields) == 0 {
		return nil
	}
	cmd := c.client.B().Hdel().Key(key).Field(fields...).Build()
	return c.client.Do(ctx, cmd).Error()
}

// Underlying returns the underlying rueidis client for advanced operations.
func (c *Client) Underlying() rueidis.Client {
	return c.client
}

// Module provides distributed cache dependencies for fx.
var Module = fx.Module("cache",
	fx.Provide(NewClient),
	fx.Provide(NewLocalCache),
	fx.Invoke(func(lc fx.Lifecycle, client *Client, local *LocalCache) {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				client.Close()
				local.Close()
				return nil
			},
		})
	}),
)

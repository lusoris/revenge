// Package cache provides a Redis/Dragonfly cache client using rueidis.
package cache

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/rueidis"
	"go.uber.org/fx"

	"github.com/lusoris/revenge/pkg/config"
)

// Config holds cache configuration.
type Config struct {
	Addr     string
	Password string
	DB       int
}

// Client wraps the rueidis client for caching operations.
type Client struct {
	rdb    rueidis.Client
	logger *slog.Logger
}

// NewClient creates a new cache client using rueidis.
func NewClient(cfg Config, logger *slog.Logger) (*Client, error) {
	opts := rueidis.ClientOption{
		InitAddress: []string{cfg.Addr},
		Password:    cfg.Password,
		SelectDB:    cfg.DB,
	}

	rdb, err := rueidis.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create rueidis client: %w", err)
	}

	return &Client{
		rdb:    rdb,
		logger: logger.With(slog.String("component", "cache")),
	}, nil
}

// Ping checks if the cache is reachable.
func (c *Client) Ping(ctx context.Context) error {
	cmd := c.rdb.B().Ping().Build()
	return c.rdb.Do(ctx, cmd).Error()
}

// Close closes the cache connection.
func (c *Client) Close() error {
	c.rdb.Close()
	return nil
}

// Get retrieves a value from cache.
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	cmd := c.rdb.B().Get().Key(key).Build()
	return c.rdb.Do(ctx, cmd).ToString()
}

// Set stores a value in cache with expiration.
func (c *Client) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	var cmd rueidis.Completed
	if expiration > 0 {
		cmd = c.rdb.B().Set().Key(key).Value(fmt.Sprintf("%v", value)).Ex(expiration).Build()
	} else {
		cmd = c.rdb.B().Set().Key(key).Value(fmt.Sprintf("%v", value)).Build()
	}
	return c.rdb.Do(ctx, cmd).Error()
}

// Delete removes keys from cache.
func (c *Client) Delete(ctx context.Context, keys ...string) error {
	cmd := c.rdb.B().Del().Key(keys...).Build()
	return c.rdb.Do(ctx, cmd).Error()
}

// GetBytes retrieves bytes from cache.
func (c *Client) GetBytes(ctx context.Context, key string) ([]byte, error) {
	cmd := c.rdb.B().Get().Key(key).Build()
	return c.rdb.Do(ctx, cmd).AsBytes()
}

// SetBytes stores bytes in cache with expiration.
func (c *Client) SetBytes(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	var cmd rueidis.Completed
	if expiration > 0 {
		cmd = c.rdb.B().Set().Key(key).Value(string(value)).Ex(expiration).Build()
	} else {
		cmd = c.rdb.B().Set().Key(key).Value(string(value)).Build()
	}
	return c.rdb.Do(ctx, cmd).Error()
}

// Exists checks if a key exists in cache.
func (c *Client) Exists(ctx context.Context, keys ...string) (int64, error) {
	cmd := c.rdb.B().Exists().Key(keys...).Build()
	return c.rdb.Do(ctx, cmd).ToInt64()
}

// Expire sets expiration on a key.
func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) error {
	cmd := c.rdb.B().Expire().Key(key).Seconds(int64(expiration.Seconds())).Build()
	return c.rdb.Do(ctx, cmd).Error()
}

// TTL gets the remaining TTL of a key.
func (c *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	cmd := c.rdb.B().Ttl().Key(key).Build()
	secs, err := c.rdb.Do(ctx, cmd).ToInt64()
	if err != nil {
		return 0, err
	}
	return time.Duration(secs) * time.Second, nil
}

// Underlying returns the underlying rueidis client for advanced operations.
func (c *Client) Underlying() rueidis.Client {
	return c.rdb
}

// Module provides cache dependencies for fx.
var Module = fx.Module("cache",
	fx.Provide(func(cfg *config.Config, logger *slog.Logger) (*Client, error) {
		addr := cfg.Cache.Addr
		if addr == "" {
			addr = "localhost:6379"
		}
		cacheConfig := Config{
			Addr:     addr,
			Password: cfg.Cache.Password,
			DB:       cfg.Cache.DB,
		}
		return NewClient(cacheConfig, logger)
	}),
)

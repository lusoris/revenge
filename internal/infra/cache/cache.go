// Package cache provides a Redis/Dragonfly cache client.
package cache

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/fx"
)

// Config holds cache configuration.
type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// Client wraps the Redis client for caching operations.
type Client struct {
	rdb    *redis.Client
	logger *slog.Logger
}

// NewClient creates a new cache client.
func NewClient(cfg Config, logger *slog.Logger) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return &Client{
		rdb:    rdb,
		logger: logger.With(slog.String("component", "cache")),
	}, nil
}

// Ping checks if the cache is reachable.
func (c *Client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

// Close closes the cache connection.
func (c *Client) Close() error {
	return c.rdb.Close()
}

// Get retrieves a value from cache.
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.rdb.Get(ctx, key).Result()
}

// Set stores a value in cache with expiration.
func (c *Client) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return c.rdb.Set(ctx, key, value, expiration).Err()
}

// Delete removes a key from cache.
func (c *Client) Delete(ctx context.Context, keys ...string) error {
	return c.rdb.Del(ctx, keys...).Err()
}

// GetBytes retrieves bytes from cache.
func (c *Client) GetBytes(ctx context.Context, key string) ([]byte, error) {
	return c.rdb.Get(ctx, key).Bytes()
}

// Module provides cache dependencies for fx.
var Module = fx.Module("cache",
	fx.Provide(func(logger *slog.Logger) (*Client, error) {
		// TODO: Get config from koanf
		cfg := Config{
			Host: "localhost",
			Port: 6379,
			DB:   0,
		}
		return NewClient(cfg, logger)
	}),
)

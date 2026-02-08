package cache

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/rueidis"
	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/config"
)

const (
	// DefaultDialTimeout is the default connection timeout.
	DefaultDialTimeout = 5 * time.Second

	// DefaultReadTimeout is the default read timeout.
	DefaultReadTimeout = 3 * time.Second

	// DefaultWriteTimeout is the default write timeout.
	DefaultWriteTimeout = 3 * time.Second

	// DefaultCacheSizeEachConn is the default client-side cache size per connection (16 MiB).
	DefaultCacheSizeEachConn = 16 * 1024 * 1024

	// DefaultRingScale is the default ring buffer scale (2^10 = 1024 slots per connection).
	DefaultRingScale = 10

	// DefaultBlockingPoolSize is the default blocking pool size.
	DefaultBlockingPoolSize = 128
)

// Module provides cache dependencies.
var Module = fx.Module("cache",
	fx.Provide(NewClient),
	fx.Provide(provideRueidisClient),
	fx.Invoke(registerHooks),
)

// provideRueidisClient extracts the underlying rueidis.Client from our Client
// wrapper so other packages (e.g. observability) can consume it without
// importing cache (which would create an import cycle).
func provideRueidisClient(c *Client) rueidis.Client {
	return c.RueidisClient()
}

// Client represents the cache client (rueidis + otter).
type Client struct {
	config        *config.Config
	logger        *slog.Logger
	rueidisClient rueidis.Client
}

// NewClient creates a new cache client with rueidis integration.
func NewClient(cfg *config.Config, logger *slog.Logger) (*Client, error) {
	client := &Client{
		config: cfg,
		logger: logger,
	}

	// Skip initialization if cache is disabled
	if !cfg.Cache.Enabled {
		logger.Info("cache disabled, client created without rueidis connection")
		return client, nil
	}

	// Validate cache URL
	if cfg.Cache.URL == "" {
		return nil, fmt.Errorf("cache URL is required when cache is enabled")
	}

	// Parse URL to extract connection options
	opts, err := rueidis.ParseURL(cfg.Cache.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cache URL: %w", err)
	}

	// Configure timeouts
	opts.Dialer.Timeout = DefaultDialTimeout
	opts.Dialer.KeepAlive = 30 * time.Second

	// Configure client-side caching
	opts.CacheSizeEachConn = DefaultCacheSizeEachConn
	opts.RingScaleEachConn = DefaultRingScale

	// Configure connection pool and timeouts
	opts.BlockingPoolSize = DefaultBlockingPoolSize
	opts.ConnWriteTimeout = DefaultWriteTimeout

	// Disable auto-pipelining for more predictable behavior
	opts.DisableAutoPipelining = false

	// Create the rueidis client
	rueidisClient, err := rueidis.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create rueidis client: %w", err)
	}

	client.rueidisClient = rueidisClient
	logger.Info("rueidis client initialized", "url", cfg.Cache.URL)

	return client, nil
}

// Close closes the Redis/Dragonfly client and releases resources.
func (c *Client) Close() {
	if c.rueidisClient != nil {
		c.rueidisClient.Close()
		c.logger.Info("rueidis client closed")
	}
}

// Ping checks the connection to Redis/Dragonfly.
func (c *Client) Ping(ctx context.Context) error {
	if c.rueidisClient == nil {
		return fmt.Errorf("rueidis client not initialized")
	}

	cmd := c.rueidisClient.B().Ping().Build()
	return c.rueidisClient.Do(ctx, cmd).Error()
}

// RueidisClient returns the underlying rueidis.Client for advanced operations.
// Returns nil if cache is disabled or client not initialized.
func (c *Client) RueidisClient() rueidis.Client {
	return c.rueidisClient
}

// registerHooks registers lifecycle hooks for the cache client.
func registerHooks(lc fx.Lifecycle, client *Client, logger *slog.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if !client.config.Cache.Enabled {
				logger.Info("cache disabled, skipping startup")
				return nil
			}

			// Test connection with timeout
			pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			if err := client.Ping(pingCtx); err != nil {
				logger.Warn("cache ping failed during startup", "error", err)
				// Don't fail startup if cache is unavailable
				return nil
			}

			logger.Info("cache client started and connected")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if !client.config.Cache.Enabled {
				return nil
			}

			client.Close()
			logger.Info("cache client stopped")
			return nil
		},
	})
}

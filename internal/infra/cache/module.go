package cache

import (
	"context"
	"log/slog"

	"github.com/lusoris/revenge/internal/config"
	"go.uber.org/fx"
)

// Module provides cache dependencies.
var Module = fx.Module("cache",
	fx.Provide(NewClient),
	fx.Invoke(registerHooks),
)

// Client represents the cache client (rueidis + otter).
// This is a placeholder stub for v0.1.0 skeleton.
type Client struct {
	config *config.Config
	logger *slog.Logger
}

// NewClient creates a new cache client.
func NewClient(cfg *config.Config, logger *slog.Logger) (*Client, error) {
	return &Client{
		config: cfg,
		logger: logger,
	}, nil
}

// registerHooks registers lifecycle hooks for the cache client.
func registerHooks(lc fx.Lifecycle, client *Client, logger *slog.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if !client.config.Cache.Enabled {
				logger.Info("cache disabled, skipping startup")
				return nil
			}
			logger.Info("cache client started (stub)")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if !client.config.Cache.Enabled {
				return nil
			}
			logger.Info("cache client stopped")
			return nil
		},
	})
}

package search

import (
	"context"
	"log/slog"

	"github.com/lusoris/revenge/internal/config"
	"go.uber.org/fx"
)

// Module provides search dependencies.
var Module = fx.Module("search",
	fx.Provide(NewClient),
	fx.Invoke(registerHooks),
)

// Client represents the Typesense search client.
// This is a placeholder stub for v0.1.0 skeleton.
type Client struct {
	config *config.Config
	logger *slog.Logger
}

// NewClient creates a new search client.
func NewClient(cfg *config.Config, logger *slog.Logger) (*Client, error) {
	return &Client{
		config: cfg,
		logger: logger,
	}, nil
}

// registerHooks registers lifecycle hooks for the search client.
func registerHooks(lc fx.Lifecycle, client *Client, logger *slog.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if !client.config.Search.Enabled {
				logger.Info("search disabled, skipping startup")
				return nil
			}
			logger.Info("search client started (stub)")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if !client.config.Search.Enabled {
				return nil
			}
			logger.Info("search client stopped")
			return nil
		},
	})
}

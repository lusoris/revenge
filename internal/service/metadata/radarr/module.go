package radarr

import (
	"log/slog"

	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/config"
)

// Module provides the Radarr metadata provider.
var Module = fx.Module("radarr",
	fx.Provide(func(cfg *config.Config, logger *slog.Logger) *Provider {
		return NewProvider(Config{
			BaseURL: cfg.Metadata.Radarr.BaseURL,
			APIKey:  cfg.Metadata.Radarr.APIKey,
		}, logger)
	}),
)

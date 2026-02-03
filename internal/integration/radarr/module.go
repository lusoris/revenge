package radarr

import (
	"log/slog"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/content/movie"
)

// Module provides the Radarr integration for fx.
var Module = fx.Module("radarr",
	fx.Provide(
		NewClientFromConfig,
		NewMapper,
		NewSyncServiceFromDeps,
	),
)

// NewClientFromConfig creates a new Radarr client from configuration.
func NewClientFromConfig(cfg *config.Config, logger *zap.Logger) *Client {
	radarrCfg := cfg.GetRadarrConfig()
	if !radarrCfg.Enabled {
		logger.Info("radarr integration disabled")
		return nil
	}

	if radarrCfg.BaseURL == "" || radarrCfg.APIKey == "" {
		logger.Warn("radarr integration enabled but not configured",
			zap.String("base_url", radarrCfg.BaseURL),
			zap.Bool("has_api_key", radarrCfg.APIKey != ""),
		)
		return nil
	}

	client := NewClient(Config{
		BaseURL:   radarrCfg.BaseURL,
		APIKey:    radarrCfg.APIKey,
		RateLimit: 10.0, // 10 req/s for local service
		CacheTTL:  5 * time.Minute,
		Timeout:   30 * time.Second,
	})

	logger.Info("radarr integration initialized",
		zap.String("base_url", radarrCfg.BaseURL),
	)

	return client
}

// SyncServiceDeps contains dependencies for the sync service.
type SyncServiceDeps struct {
	fx.In

	Client    *Client           `optional:"true"`
	Mapper    *Mapper
	MovieRepo movie.Repository
	Logger    *slog.Logger
}

// NewSyncServiceFromDeps creates a new sync service from dependencies.
func NewSyncServiceFromDeps(deps SyncServiceDeps) *SyncService {
	if deps.Client == nil {
		return nil
	}
	return NewSyncService(deps.Client, deps.Mapper, deps.MovieRepo, deps.Logger)
}

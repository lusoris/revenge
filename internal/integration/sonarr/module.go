package sonarr

import (
	"log/slog"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/content/tvshow"
)

// Module provides the Sonarr integration for fx.
var Module = fx.Module("sonarr",
	fx.Provide(
		NewClientFromConfig,
		NewMapper,
		NewSyncServiceFromDeps,
	),
)

// NewClientFromConfig creates a new Sonarr client from configuration.
func NewClientFromConfig(cfg *config.Config, logger *zap.Logger) *Client {
	sonarrCfg := cfg.GetSonarrConfig()
	if !sonarrCfg.Enabled {
		logger.Info("sonarr integration disabled")
		return nil
	}

	if sonarrCfg.BaseURL == "" || sonarrCfg.APIKey == "" {
		logger.Warn("sonarr integration enabled but not configured",
			zap.String("base_url", sonarrCfg.BaseURL),
			zap.Bool("has_api_key", sonarrCfg.APIKey != ""),
		)
		return nil
	}

	client := NewClient(Config{
		BaseURL:   sonarrCfg.BaseURL,
		APIKey:    sonarrCfg.APIKey,
		RateLimit: 10.0, // 10 req/s for local service
		CacheTTL:  5 * time.Minute,
		Timeout:   30 * time.Second,
	})

	logger.Info("sonarr integration initialized",
		zap.String("base_url", sonarrCfg.BaseURL),
	)

	return client
}

// SyncServiceDeps contains dependencies for the sync service.
type SyncServiceDeps struct {
	fx.In

	Client     *Client            `optional:"true"`
	Mapper     *Mapper
	TVShowRepo tvshow.Repository
	Logger     *slog.Logger
}

// NewSyncServiceFromDeps creates a new sync service from dependencies.
func NewSyncServiceFromDeps(deps SyncServiceDeps) *SyncService {
	if deps.Client == nil {
		return nil
	}
	return NewSyncService(deps.Client, deps.Mapper, deps.TVShowRepo, deps.Logger)
}

package sonarr

import (
	"log/slog"
	"time"

	"github.com/riverqueue/river"
	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/content/tvshow"
)

// Module provides the Sonarr integration for fx.
var Module = fx.Module("sonarr",
	fx.Provide(
		NewClientFromConfig,
		NewMapper,
		NewSyncServiceFromDeps,
		NewWebhookHandlerFromDeps,
		NewSonarrSyncWorkerFromDeps,
		NewSonarrWebhookWorkerFromDeps,
	),
	fx.Invoke(registerSonarrWorkers),
)

// WebhookHandlerDeps contains dependencies for the webhook handler.
type WebhookHandlerDeps struct {
	fx.In

	SyncService *SyncService `optional:"true"`
	Logger      *slog.Logger
}

// NewWebhookHandlerFromDeps creates a new webhook handler from fx dependencies.
func NewWebhookHandlerFromDeps(deps WebhookHandlerDeps) *WebhookHandler {
	return NewWebhookHandler(deps.SyncService, deps.Logger)
}

// WorkerDeps contains dependencies for Sonarr River workers.
type WorkerDeps struct {
	fx.In

	SyncService    *SyncService `optional:"true"`
	WebhookHandler *WebhookHandler
	Logger         *slog.Logger
}

// NewSonarrSyncWorkerFromDeps creates a new Sonarr sync worker from fx dependencies.
func NewSonarrSyncWorkerFromDeps(deps WorkerDeps) *SonarrSyncWorker {
	return NewSonarrSyncWorker(deps.SyncService, deps.Logger)
}

// NewSonarrWebhookWorkerFromDeps creates a new Sonarr webhook worker from fx dependencies.
func NewSonarrWebhookWorkerFromDeps(deps WorkerDeps) *SonarrWebhookWorker {
	return NewSonarrWebhookWorker(deps.WebhookHandler, deps.Logger)
}

// registerSonarrWorkers registers Sonarr workers with the River workers registry.
func registerSonarrWorkers(workers *river.Workers, syncWorker *SonarrSyncWorker, webhookWorker *SonarrWebhookWorker) {
	river.AddWorker(workers, syncWorker)
	river.AddWorker(workers, webhookWorker)
}

// NewClientFromConfig creates a new Sonarr client from configuration.
func NewClientFromConfig(cfg *config.Config, logger *slog.Logger) *Client {
	sonarrCfg := cfg.GetSonarrConfig()
	if !sonarrCfg.Enabled {
		logger.Info("sonarr integration disabled")
		return nil
	}

	if sonarrCfg.BaseURL == "" || sonarrCfg.APIKey == "" {
		logger.Warn("sonarr integration enabled but not configured",
			slog.String("base_url", sonarrCfg.BaseURL),
			slog.Bool("has_api_key", sonarrCfg.APIKey != ""),
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
		slog.String("base_url", sonarrCfg.BaseURL),
	)

	return client
}

// SyncServiceDeps contains dependencies for the sync service.
type SyncServiceDeps struct {
	fx.In

	Client     *Client `optional:"true"`
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

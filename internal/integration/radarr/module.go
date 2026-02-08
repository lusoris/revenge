package radarr

import (
	"log/slog"
	"time"

	"github.com/riverqueue/river"
	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/content/movie"
)

// Module provides the Radarr integration for fx.
var Module = fx.Module("radarr",
	fx.Provide(
		NewClientFromConfig,
		NewMapper,
		NewSyncServiceFromDeps,
		NewWebhookHandlerFromDeps,
		NewRadarrSyncWorkerFromDeps,
		NewRadarrWebhookWorkerFromDeps,
	),
	fx.Invoke(registerRadarrWorkers),
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

// WorkerDeps contains dependencies for Radarr River workers.
type WorkerDeps struct {
	fx.In

	SyncService    *SyncService    `optional:"true"`
	WebhookHandler *WebhookHandler
	Logger         *slog.Logger
}

// NewRadarrSyncWorkerFromDeps creates a new Radarr sync worker from fx dependencies.
func NewRadarrSyncWorkerFromDeps(deps WorkerDeps) *RadarrSyncWorker {
	return NewRadarrSyncWorker(deps.SyncService, deps.Logger)
}

// NewRadarrWebhookWorkerFromDeps creates a new Radarr webhook worker from fx dependencies.
func NewRadarrWebhookWorkerFromDeps(deps WorkerDeps) *RadarrWebhookWorker {
	return NewRadarrWebhookWorker(deps.WebhookHandler, deps.Logger)
}

// registerRadarrWorkers registers Radarr workers with the River workers registry.
func registerRadarrWorkers(workers *river.Workers, syncWorker *RadarrSyncWorker, webhookWorker *RadarrWebhookWorker) {
	river.AddWorker(workers, syncWorker)
	river.AddWorker(workers, webhookWorker)
}

// NewClientFromConfig creates a new Radarr client from configuration.
func NewClientFromConfig(cfg *config.Config, logger *slog.Logger) *Client {
	radarrCfg := cfg.GetRadarrConfig()
	if !radarrCfg.Enabled {
		logger.Info("radarr integration disabled")
		return nil
	}

	if radarrCfg.BaseURL == "" || radarrCfg.APIKey == "" {
		logger.Warn("radarr integration enabled but not configured",
			slog.String("base_url", radarrCfg.BaseURL),
			slog.Bool("has_api_key", radarrCfg.APIKey != ""),
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
		slog.String("base_url", radarrCfg.BaseURL),
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

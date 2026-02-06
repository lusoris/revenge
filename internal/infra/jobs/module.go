package jobs

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/config"
)

// Module provides job queue dependencies.
var Module = fx.Module("jobs",
	fx.Provide(
		NewRiverWorkers,
		NewRiverClient,
	),
	fx.Invoke(registerHooks),
)

// NewRiverWorkers creates a new River workers registry.
func NewRiverWorkers() *river.Workers {
	return river.NewWorkers()
}

// NewRiverClient creates a new River client.
func NewRiverClient(
	pool *pgxpool.Pool,
	workers *river.Workers,
	cfg *config.Config,
	logger *slog.Logger,
) (*Client, error) {
	queueCfg := DefaultQueueConfig()
	jobsConfig := &Config{
		Queues:        queueCfg.Queues,
		FetchCooldown: cfg.Jobs.FetchCooldown,
		MaxAttempts:   25, // TODO: Make configurable
	}

	return NewClient(pool, workers, jobsConfig, logger)
}

// registerHooks registers lifecycle hooks for the job client.
func registerHooks(lc fx.Lifecycle, client *Client, logger *slog.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("starting River job queue client")
			return client.Start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("stopping River job queue client")
			return client.Stop(ctx)
		},
	})
}

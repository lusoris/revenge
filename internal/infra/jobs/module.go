package jobs

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lusoris/revenge/internal/config"
	"go.uber.org/fx"
)

// Module provides job queue dependencies.
var Module = fx.Module("jobs",
	fx.Provide(NewWorkers),
	fx.Invoke(registerHooks),
)

// Workers represents the River job queue workers.
// This is a placeholder stub for v0.1.0 skeleton.
type Workers struct {
	config *config.Config
	pool   *pgxpool.Pool
	logger *slog.Logger
}

// NewWorkers creates new job queue workers.
func NewWorkers(cfg *config.Config, pool *pgxpool.Pool, logger *slog.Logger) (*Workers, error) {
	return &Workers{
		config: cfg,
		pool:   pool,
		logger: logger,
	}, nil
}

// registerHooks registers lifecycle hooks for the job workers.
func registerHooks(lc fx.Lifecycle, workers *Workers, logger *slog.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("job workers started (stub)")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("job workers stopped")
			return nil
		},
	})
}

package jobs

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivermigrate"
	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/config"
)

// Module provides job queue dependencies.
var Module = fx.Module("jobs",
	fx.Provide(
		NewRiverWorkers,
		NewRiverClient,
		NewNotificationWorker,
	),
	fx.Invoke(
		registerHooks,
		registerCleanupWorker,
		registerNotificationWorker,
	),
)

// registerCleanupWorker registers the cleanup worker with the River workers registry.
func registerCleanupWorker(workers *river.Workers, cleanupWorker *CleanupWorker) {
	river.AddWorker(workers, cleanupWorker)
}

// registerNotificationWorker registers the notification worker with the River workers registry.
func registerNotificationWorker(workers *river.Workers, notificationWorker *NotificationWorker) {
	river.AddWorker(workers, notificationWorker)
}

// NewRiverWorkers creates a new River workers registry.
func NewRiverWorkers() *river.Workers {
	return river.NewWorkers()
}

// NewRiverClient creates a new River client.
// It runs River's schema migrations before creating the client.
func NewRiverClient(
	pool *pgxpool.Pool,
	workers *river.Workers,
	cfg *config.Config,
	logger *slog.Logger,
	periodicJobs []*river.PeriodicJob,
) (*Client, error) {
	// Run River schema migrations (creates river_job, river_queue, river_leader tables)
	migrator, err := rivermigrate.New(riverpgxv5.New(pool), nil)
	if err != nil {
		return nil, err
	}
	if _, err := migrator.Migrate(context.Background(), rivermigrate.DirectionUp, nil); err != nil {
		return nil, err
	}
	logger.Info("River schema migrations applied")

	queueCfg := DefaultQueueConfig()

	// Apply configured MaxWorkers as a total cap across all queues.
	// Scale per-queue workers proportionally if MaxWorkers is set.
	if cfg.Jobs.MaxWorkers > 0 {
		scaleQueueWorkers(queueCfg.Queues, cfg.Jobs.MaxWorkers)
	}

	jobsConfig := &Config{
		Queues:               queueCfg.Queues,
		FetchCooldown:        cfg.Jobs.FetchCooldown,
		FetchPollInterval:    cfg.Jobs.FetchPollInterval,
		RescueStuckJobsAfter: cfg.Jobs.RescueStuckJobsAfter,
		MaxAttempts:          cfg.Jobs.MaxAttempts,
		PeriodicJobs:         periodicJobs,
	}

	return NewClient(pool, workers, jobsConfig, logger)
}

// scaleQueueWorkers distributes maxWorkers across queues proportionally
// based on their default worker allocation.
func scaleQueueWorkers(queues map[string]river.QueueConfig, maxWorkers int) {
	// Calculate total default workers
	var totalDefault int
	for _, qc := range queues {
		totalDefault += qc.MaxWorkers
	}
	if totalDefault == 0 {
		return
	}

	// Scale each queue proportionally, ensuring at least 1 worker per queue
	for name, qc := range queues {
		scaled := (qc.MaxWorkers * maxWorkers) / totalDefault
		if scaled < 1 {
			scaled = 1
		}
		qc.MaxWorkers = scaled
		queues[name] = qc
	}
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

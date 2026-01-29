// Package jobs provides River job queue setup.
package jobs

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivertype"
	"go.uber.org/fx"
)

// Config holds job queue configuration.
type Config struct {
	// MaxWorkers is the default max workers per queue.
	MaxWorkers int
}

// Service wraps the River client for job queue operations.
type Service struct {
	client  *river.Client[pgx.Tx]
	workers *river.Workers
	logger  *slog.Logger
}

// NewService creates a new job queue service.
func NewService(pool *pgxpool.Pool, workers *river.Workers, logger *slog.Logger) (*Service, error) {
	riverClient, err := river.NewClient(riverpgxv5.New(pool), &river.Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 100},
			QueueScanning:      {MaxWorkers: 2},  // Heavy IO
			QueueMetadata:      {MaxWorkers: 10}, // Network bound
			QueueIndexing:      {MaxWorkers: 5},  // CPU bound
			QueueImages:        {MaxWorkers: 5},  // Network bound
			QueueCleanup:       {MaxWorkers: 1},  // Background maintenance
		},
		Workers: workers,
		Logger:  logger.With(slog.String("component", "river")),
	})
	if err != nil {
		return nil, err
	}

	return &Service{
		client:  riverClient,
		workers: workers,
		logger:  logger.With(slog.String("component", "jobs")),
	}, nil
}

// Start starts the job queue worker.
func (s *Service) Start(ctx context.Context) error {
	return s.client.Start(ctx)
}

// Stop stops the job queue worker.
func (s *Service) Stop(ctx context.Context) error {
	return s.client.Stop(ctx)
}

// Insert inserts a new job into the queue.
func (s *Service) Insert(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error) {
	return s.client.Insert(ctx, args, opts)
}

// InsertTx inserts a new job into the queue within a transaction.
func (s *Service) InsertTx(ctx context.Context, tx pgx.Tx, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error) {
	return s.client.InsertTx(ctx, tx, args, opts)
}

// Healthy returns nil if the job service is healthy.
func (s *Service) Healthy(ctx context.Context) error {
	// River's health is tied to the database connection.
	// The client maintains its own connection health internally.
	// If the service started successfully, it's considered healthy.
	return nil
}

// NewWorkers creates a new River workers registry.
func NewWorkers() *river.Workers {
	return river.NewWorkers()
}

// WorkerDepsInput contains optional worker dependencies for fx injection.
type WorkerDepsInput struct {
	fx.In
	Scanner    LibraryScanner     `optional:"true"`
	Fetcher    MetadataFetcher    `optional:"true"`
	Downloader ImageDownloader    `optional:"true"`
	Indexer    SearchIndexer      `optional:"true"`
	Cleaner    CleanupService     `optional:"true"`
	Trickplay  TrickplayGenerator `optional:"true"`
	Logger     *slog.Logger
}

// Module provides job queue dependencies for fx.
var Module = fx.Module("jobs",
	fx.Provide(NewWorkers),
	fx.Invoke(registerAllWorkers),
	fx.Provide(NewService),
	fx.Provide(func(svc *Service) *river.Client[pgx.Tx] {
		return svc.client
	}),
	fx.Invoke(func(lc fx.Lifecycle, svc *Service) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				return svc.Start(ctx)
			},
			OnStop: func(ctx context.Context) error {
				return svc.Stop(ctx)
			},
		})
	}),
)

// registerAllWorkers registers all workers with their dependencies.
func registerAllWorkers(workers *river.Workers, deps WorkerDepsInput) error {
	return RegisterWorkers(workers, WorkerDeps{
		Scanner:    deps.Scanner,
		Fetcher:    deps.Fetcher,
		Downloader: deps.Downloader,
		Indexer:    deps.Indexer,
		Cleaner:    deps.Cleaner,
		Trickplay:  deps.Trickplay,
		Logger:     deps.Logger,
	})
}

package jobs

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lusoris/revenge/internal/infra/observability"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivertype"
)

// Config holds River client configuration.
type Config struct {
	Queues                      map[string]river.QueueConfig
	FetchCooldown               time.Duration
	FetchPollInterval           time.Duration
	RescueStuckJobsAfter        time.Duration
	MaxAttempts                 int
	PeriodicJobs                []*river.PeriodicJob
	CompletedJobRetentionPeriod time.Duration
	DiscardedJobRetentionPeriod time.Duration
}

// DefaultConfig returns default River client configuration.
func DefaultConfig() *Config {
	return &Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 100},
		},
		FetchCooldown:               100 * time.Millisecond,
		FetchPollInterval:           500 * time.Millisecond,
		RescueStuckJobsAfter:        1 * time.Hour,
		MaxAttempts:                 5,
		CompletedJobRetentionPeriod: 24 * time.Hour,
		DiscardedJobRetentionPeriod: 7 * 24 * time.Hour,
	}
}

// Client wraps River job queue client.
type Client struct {
	client *river.Client[pgx.Tx]
	config *Config
	logger *slog.Logger
}

// NewClient creates a new River client.
func NewClient(pool *pgxpool.Pool, workers *river.Workers, config *Config, logger *slog.Logger) (*Client, error) {
	if pool == nil {
		return nil, errors.New("database pool is required")
	}
	if workers == nil {
		return nil, errors.New("workers registry is required")
	}
	if config == nil {
		config = DefaultConfig()
	}
	if logger == nil {
		logger = slog.Default()
	}

	// Create a leveled logger for River that only logs WARN+ to reduce polling spam
	riverLogger := slog.New(&leveledHandler{
		handler:  logger.Handler(),
		minLevel: slog.LevelWarn,
	})

	riverConfig := &river.Config{
		Queues:                      config.Queues,
		FetchCooldown:               config.FetchCooldown,
		FetchPollInterval:           config.FetchPollInterval,
		RescueStuckJobsAfter:        config.RescueStuckJobsAfter,
		MaxAttempts:                 config.MaxAttempts,
		PeriodicJobs:                config.PeriodicJobs,
		Workers:                     workers,
		Logger:                      riverLogger,
		JobTimeout:                  -1, // Per-worker Timeout() methods handle this
		CompletedJobRetentionPeriod: config.CompletedJobRetentionPeriod,
		DiscardedJobRetentionPeriod: config.DiscardedJobRetentionPeriod,
	}

	riverClient, err := river.NewClient(riverpgxv5.New(pool), riverConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create River client: %w", err)
	}

	return &Client{
		client: riverClient,
		config: config,
		logger: logger,
	}, nil
}

// Start starts the River client.
func (c *Client) Start(ctx context.Context) error {
	if c.client == nil {
		return errors.New("river client not initialized")
	}
	c.logger.Info("starting River job queue client")
	if err := c.client.Start(ctx); err != nil {
		return fmt.Errorf("failed to start River client: %w", err)
	}

	// Subscribe to job lifecycle events for metrics
	eventCh, cancel := c.client.Subscribe(river.EventKindJobCompleted, river.EventKindJobFailed)
	go func() {
		defer cancel()
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-eventCh:
				if !ok {
					return
				}
				jobKind := event.Job.Kind

				switch event.Kind {
				case river.EventKindJobCompleted:
					observability.RecordJobCompleted(jobKind, "success")
				case river.EventKindJobFailed:
					observability.RecordJobCompleted(jobKind, "failure")
				}

				if event.Job.FinalizedAt != nil && event.Job.AttemptedAt != nil {
					duration := event.Job.FinalizedAt.Sub(*event.Job.AttemptedAt).Seconds()
					observability.JobDuration.WithLabelValues(jobKind).Observe(duration)
				}
			}
		}
	}()

	c.logger.Info("River job queue client started")
	return nil
}

// Stop gracefully stops the River client.
func (c *Client) Stop(ctx context.Context) error {
	if c.client == nil {
		return errors.New("river client not initialized")
	}
	c.logger.Info("stopping River job queue client")
	if err := c.client.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop River client: %w", err)
	}
	c.logger.Info("River job queue client stopped")
	return nil
}

// Insert enqueues a new job.
func (c *Client) Insert(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error) {
	if c.client == nil {
		return nil, errors.New("river client not initialized")
	}
	result, err := c.client.Insert(ctx, args, opts)
	if err == nil {
		observability.RecordJobEnqueued(args.Kind())
	}
	return result, err
}

// InsertMany enqueues multiple jobs in a single transaction.
func (c *Client) InsertMany(ctx context.Context, params []river.InsertManyParams) ([]*rivertype.JobInsertResult, error) {
	if c.client == nil {
		return nil, errors.New("river client not initialized")
	}
	results, err := c.client.InsertMany(ctx, params)
	if err == nil {
		for _, p := range params {
			observability.RecordJobEnqueued(p.Args.Kind())
		}
	}
	return results, err
}

// JobGet retrieves a job by ID.
func (c *Client) JobGet(ctx context.Context, id int64) (*rivertype.JobRow, error) {
	if c.client == nil {
		return nil, errors.New("river client not initialized")
	}
	return c.client.JobGet(ctx, id)
}

// JobCancel cancels a job by ID.
func (c *Client) JobCancel(ctx context.Context, id int64) (*rivertype.JobRow, error) {
	if c.client == nil {
		return nil, errors.New("river client not initialized")
	}
	return c.client.JobCancel(ctx, id)
}

// Subscribe subscribes to job lifecycle events.
func (c *Client) Subscribe(kinds ...river.EventKind) (<-chan *river.Event, func()) {
	if c.client == nil {
		panic("river client not initialized")
	}
	return c.client.Subscribe(kinds...)
}

// RiverClient returns the underlying River client for advanced usage.
func (c *Client) RiverClient() *river.Client[pgx.Tx] {
	return c.client
}

// leveledHandler wraps an slog.Handler to filter logs below a minimum level.
// Used to reduce River's verbose polling logs.
type leveledHandler struct {
	handler  slog.Handler
	minLevel slog.Level
}

func (h *leveledHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.minLevel
}

func (h *leveledHandler) Handle(ctx context.Context, r slog.Record) error {
	return h.handler.Handle(ctx, r)
}

func (h *leveledHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &leveledHandler{handler: h.handler.WithAttrs(attrs), minLevel: h.minLevel}
}

func (h *leveledHandler) WithGroup(name string) slog.Handler {
	return &leveledHandler{handler: h.handler.WithGroup(name), minLevel: h.minLevel}
}

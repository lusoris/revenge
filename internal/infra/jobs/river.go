package jobs

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/riverqueue/river/rivertype"
)

// Config holds River client configuration.
type Config struct {
	Queues               map[string]river.QueueConfig
	FetchCooldown        time.Duration
	FetchPollInterval    time.Duration
	RescueStuckJobsAfter time.Duration
	MaxAttempts          int
}

// DefaultConfig returns default River client configuration.
func DefaultConfig() *Config {
	return &Config{
		Queues: map[string]river.QueueConfig{
			river.QueueDefault: {MaxWorkers: 100},
		},
		FetchCooldown:        100 * time.Millisecond,
		FetchPollInterval:    500 * time.Millisecond,
		RescueStuckJobsAfter: 1 * time.Hour,
		MaxAttempts:          25,
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

	riverConfig := &river.Config{
		Queues:               config.Queues,
		FetchCooldown:        config.FetchCooldown,
		FetchPollInterval:    config.FetchPollInterval,
		RescueStuckJobsAfter: config.RescueStuckJobsAfter,
		MaxAttempts:          config.MaxAttempts,
		Workers:              workers,
		Logger:               logger,
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
	return c.client.Insert(ctx, args, opts)
}

// InsertMany enqueues multiple jobs in a single transaction.
func (c *Client) InsertMany(ctx context.Context, params []river.InsertManyParams) ([]*rivertype.JobInsertResult, error) {
	if c.client == nil {
		return nil, errors.New("river client not initialized")
	}
	return c.client.InsertMany(ctx, params)
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

// Package database provides PostgreSQL connection pooling and migration support.
package database

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/errors"
)

// PoolConfig converts application config to pgxpool config.
func PoolConfig(cfg *config.Config) (*pgxpool.Config, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.Database.URL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse database URL")
	}

	// Apply connection pool settings from config
	if cfg.Database.MaxConns > 0 {
		poolConfig.MaxConns = int32(cfg.Database.MaxConns)
	} else {
		// Default: (CPU * 2) + 1
		poolConfig.MaxConns = int32((runtime.NumCPU() * 2) + 1)
	}

	if cfg.Database.MinConns > 0 {
		poolConfig.MinConns = int32(cfg.Database.MinConns)
	}

	if cfg.Database.MaxConnLifetime > 0 {
		poolConfig.MaxConnLifetime = cfg.Database.MaxConnLifetime
	}

	if cfg.Database.MaxConnIdleTime > 0 {
		poolConfig.MaxConnIdleTime = cfg.Database.MaxConnIdleTime
	}

	if cfg.Database.HealthCheckPeriod > 0 {
		poolConfig.HealthCheckPeriod = cfg.Database.HealthCheckPeriod
	}

	return poolConfig, nil
}

// NewPool creates a new PostgreSQL connection pool.
func NewPool(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*pgxpool.Pool, error) {
	poolConfig, err := PoolConfig(cfg)
	if err != nil {
		return nil, err
	}

	logger.Info("connecting to database",
		slog.String("database", poolConfig.ConnConfig.Database),
		slog.String("host", poolConfig.ConnConfig.Host),
		slog.Int("max_conns", int(poolConfig.MaxConns)),
		slog.Int("min_conns", int(poolConfig.MinConns)),
	)

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create connection pool")
	}

	// Verify connection
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, errors.Wrap(err, "failed to ping database")
	}

	logger.Info("database connection established")

	return pool, nil
}

// Stats returns human-readable pool statistics.
func Stats(pool *pgxpool.Pool) map[string]interface{} {
	stat := pool.Stat()
	return map[string]interface{}{
		"acquire_count":              stat.AcquireCount(),
		"acquire_duration_ms":        stat.AcquireDuration().Milliseconds(),
		"acquired_conns":             stat.AcquiredConns(),
		"canceled_acquire_count":     stat.CanceledAcquireCount(),
		"constructing_conns":         stat.ConstructingConns(),
		"empty_acquire_count":        stat.EmptyAcquireCount(),
		"idle_conns":                 stat.IdleConns(),
		"max_conns":                  stat.MaxConns(),
		"total_conns":                stat.TotalConns(),
		"new_conns_count":            stat.NewConnsCount(),
		"max_lifetime_destroy_count": stat.MaxLifetimeDestroyCount(),
		"max_idle_destroy_count":     stat.MaxIdleDestroyCount(),
	}
}

// Health checks if the database is healthy.
func Health(ctx context.Context, pool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return errors.Wrap(err, "database health check failed")
	}

	// Check if we can acquire a connection
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to acquire connection")
	}
	defer conn.Release()

	// Simple query to verify database is responding
	var result int
	if err := conn.QueryRow(ctx, "SELECT 1").Scan(&result); err != nil {
		return errors.Wrap(err, "database query failed")
	}

	if result != 1 {
		return fmt.Errorf("unexpected query result: %d", result)
	}

	return nil
}

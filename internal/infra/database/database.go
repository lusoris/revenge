// Package database provides PostgreSQL database connectivity for Jellyfin Go.
// It uses pgxpool for connection pooling and integrates with fx for lifecycle management.
package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	"github.com/jellyfin/jellyfin-go/pkg/config"
)

// Module provides database dependencies for fx.
var Module = fx.Module("database",
	fx.Provide(
		NewPool,
	),
)

// PoolParams contains dependencies for creating a database pool.
type PoolParams struct {
	fx.In

	Config *config.Config
	Logger *slog.Logger
	LC     fx.Lifecycle
}

// PoolResult contains the outputs from creating a database pool.
type PoolResult struct {
	fx.Out

	Pool *pgxpool.Pool
}

// NewPool creates a new PostgreSQL connection pool with lifecycle management.
func NewPool(p PoolParams) (PoolResult, error) {
	cfg := p.Config.Database
	logger := p.Logger.With(slog.String("component", "database"))

	// Build connection string
	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s pool_max_conns=%d",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode,
		cfg.MaxConns,
	)

	// Parse config (allows additional pgxpool configuration)
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return PoolResult{}, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Configure pool settings
	poolConfig.MaxConns = int32(cfg.MaxConns)
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = 1 * time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	// Create pool (doesn't connect yet)
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return PoolResult{}, fmt.Errorf("failed to create database pool: %w", err)
	}

	logger.Info("Database pool created",
		slog.String("host", cfg.Host),
		slog.Int("port", cfg.Port),
		slog.String("database", cfg.Name),
		slog.Int("max_conns", cfg.MaxConns),
	)

	// Register lifecycle hooks
	p.LC.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("Connecting to PostgreSQL...")

			// Ping to verify connection
			if err := pool.Ping(ctx); err != nil {
				return fmt.Errorf("failed to connect to database: %w", err)
			}

			logger.Info("Database connection established")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Closing database connections...")
			pool.Close()
			logger.Info("Database connections closed")
			return nil
		},
	})

	return PoolResult{Pool: pool}, nil
}

// HealthCheck performs a database health check.
func HealthCheck(ctx context.Context, pool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var result int
	err := pool.QueryRow(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// Stats returns current pool statistics.
func Stats(pool *pgxpool.Pool) map[string]any {
	stat := pool.Stat()
	return map[string]any{
		"total_conns":       stat.TotalConns(),
		"acquired_conns":    stat.AcquiredConns(),
		"idle_conns":        stat.IdleConns(),
		"max_conns":         stat.MaxConns(),
		"constructing_conns": stat.ConstructingConns(),
	}
}

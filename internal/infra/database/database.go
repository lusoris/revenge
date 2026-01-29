// Package database provides PostgreSQL database access using pgx v5.
package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/config"
)

// NewPool creates a new PostgreSQL connection pool.
func NewPool(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.Database.DSN())
	if err != nil {
		return nil, fmt.Errorf("parse database config: %w", err)
	}

	poolConfig.MaxConns = cfg.Database.MaxConns
	poolConfig.MinConns = cfg.Database.MinConns

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	logger.Info("database connected",
		slog.String("host", cfg.Database.Host),
		slog.Int("port", cfg.Database.Port),
		slog.String("database", cfg.Database.Name),
	)

	return pool, nil
}

// Module provides database dependencies for fx.
var Module = fx.Module("database",
	fx.Provide(func(lc fx.Lifecycle, cfg *config.Config, logger *slog.Logger) (*pgxpool.Pool, error) {
		ctx := context.Background()
		pool, err := NewPool(ctx, cfg, logger)
		if err != nil {
			return nil, err
		}

		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				logger.Info("closing database connection pool")
				pool.Close()
				return nil
			},
		})

		return pool, nil
	}),
)

package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lusoris/revenge/internal/infra/database/db"
	"go.uber.org/fx"
)

// Module provides database dependencies.
var Module = fx.Module("database",
	fx.Provide(NewPool),
	fx.Provide(NewQueries),
	fx.Invoke(registerHooks),
)

// NewQueries creates a new Queries instance from the pool.
func NewQueries(pool *pgxpool.Pool) *db.Queries {
	return db.New(pool)
}

// registerHooks registers lifecycle hooks for the database pool.
func registerHooks(lc fx.Lifecycle, pool *pgxpool.Pool, logger *slog.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// pgxpool.New returns without waiting for connections to be established.
			// Ping ensures at least one connection is ready before other services start.
			// See: https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool#hdr-Creating_a_Pool
			if err := pool.Ping(ctx); err != nil {
				return fmt.Errorf("database ping failed: %w", err)
			}
			logger.Info("database pool started",
				slog.Int("total_conns", int(pool.Stat().TotalConns())),
			)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("closing database pool")
			pool.Close()
			logger.Info("database pool closed")
			return nil
		},
	})
}

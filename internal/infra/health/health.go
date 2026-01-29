// Package health provides infrastructure health check registration.
package health

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"

	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/lusoris/revenge/internal/infra/search"
	"github.com/lusoris/revenge/pkg/health"
)

// Deps contains dependencies for health checks.
type Deps struct {
	fx.In

	Logger *slog.Logger
	Pool   *pgxpool.Pool
	Cache  *cache.Client   `optional:"true"`
	Jobs   *jobs.Service   `optional:"true"`
	Search *search.Client  `optional:"true"`
}

// NewChecker creates a health checker with registered infrastructure checks.
func NewChecker(deps Deps) *health.Checker {
	checker := health.NewChecker(deps.Logger)

	// Register database health check (critical)
	checker.RegisterFunc("database", health.CategoryCritical, func(ctx context.Context) error {
		return deps.Pool.Ping(ctx)
	})

	// Register cache health check (warm - system works without it, but degraded)
	if deps.Cache != nil {
		checker.RegisterFunc("cache", health.CategoryWarm, func(ctx context.Context) error {
			return deps.Cache.Ping(ctx)
		})
	}

	// Register jobs health check (warm - background processing degraded without it)
	if deps.Jobs != nil {
		checker.RegisterFunc("jobs", health.CategoryWarm, func(ctx context.Context) error {
			return deps.Jobs.Healthy(ctx)
		})
	}

	// Register search health check (warm - search degraded without it)
	if deps.Search != nil {
		checker.RegisterFunc("search", health.CategoryWarm, func(ctx context.Context) error {
			return deps.Search.Health(ctx)
		})
	}

	deps.Logger.Info("health checks registered")

	return checker
}

// Module provides health checking dependencies for fx.
var Module = fx.Module("health",
	fx.Provide(NewChecker),
)

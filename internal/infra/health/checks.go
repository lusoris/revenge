package health

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/database"
	"github.com/lusoris/revenge/internal/infra/jobs"
)

// CheckDatabase checks if the PostgreSQL database is healthy.
func CheckDatabase(ctx context.Context, pool *pgxpool.Pool) CheckResult {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if pool == nil {
		return CheckResult{
			Name:    "database",
			Status:  StatusUnhealthy,
			Message: "database pool not initialized",
		}
	}

	if err := database.Health(ctx, pool); err != nil {
		return CheckResult{
			Name:    "database",
			Status:  StatusUnhealthy,
			Message: err.Error(),
		}
	}

	// Get pool stats
	stats := database.Stats(pool)

	return CheckResult{
		Name:    "database",
		Status:  StatusHealthy,
		Message: "database is healthy",
		Details: stats,
	}
}

// CheckCache checks if the Dragonfly/Redis cache is healthy.
func CheckCache(ctx context.Context, client *cache.Client) CheckResult {
	if client == nil {
		return CheckResult{
			Name:    "cache",
			Status:  StatusDegraded,
			Message: "cache client not initialized (disabled)",
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx); err != nil {
		return CheckResult{
			Name:    "cache",
			Status:  StatusUnhealthy,
			Message: err.Error(),
		}
	}

	return CheckResult{
		Name:    "cache",
		Status:  StatusHealthy,
		Message: "cache is healthy",
	}
}

// CheckJobs checks if the River job queue client is healthy.
func CheckJobs(ctx context.Context, client *jobs.Client) CheckResult {
	if client == nil {
		return CheckResult{
			Name:    "jobs",
			Status:  StatusDegraded,
			Message: "job queue client not initialized",
		}
	}

	// Check if the underlying River client is available
	riverClient := client.RiverClient()
	if riverClient == nil {
		return CheckResult{
			Name:    "jobs",
			Status:  StatusUnhealthy,
			Message: "river client not initialized",
		}
	}

	return CheckResult{
		Name:    "jobs",
		Status:  StatusHealthy,
		Message: "job queue is healthy",
	}
}

// CheckAll runs all dependency health checks.
func CheckAll(ctx context.Context, pool *pgxpool.Pool, cacheClient *cache.Client, jobClient *jobs.Client) map[string]CheckResult {
	return map[string]CheckResult{
		"database": CheckDatabase(ctx, pool),
		"cache":    CheckCache(ctx, cacheClient),
		"jobs":     CheckJobs(ctx, jobClient),
	}
}

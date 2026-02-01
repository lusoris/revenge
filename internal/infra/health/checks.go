package health

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/database"
	"github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/lusoris/revenge/internal/infra/search"
)

// CheckDatabase checks if the PostgreSQL database is healthy.
func CheckDatabase(ctx context.Context, pool *pgxpool.Pool) CheckResult {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

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
// This is a stub for v0.1.0 skeleton.
func CheckCache(ctx context.Context, client *cache.Client) CheckResult {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// TODO: Implement actual cache health check when cache client is implemented
	return CheckResult{
		Name:    "cache",
		Status:  StatusHealthy,
		Message: "cache check not implemented (stub)",
	}
}

// CheckSearch checks if the Typesense search service is healthy.
// This is a stub for v0.1.0 skeleton.
func CheckSearch(ctx context.Context, client *search.Client) CheckResult {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// TODO: Implement actual search health check when search client is implemented
	return CheckResult{
		Name:    "search",
		Status:  StatusHealthy,
		Message: "search check not implemented (stub)",
	}
}

// CheckJobs checks if the River job queue workers are healthy.
// This is a stub for v0.1.0 skeleton.
func CheckJobs(ctx context.Context, workers *jobs.Workers) CheckResult {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// TODO: Implement actual job workers health check when workers are implemented
	return CheckResult{
		Name:    "jobs",
		Status:  StatusHealthy,
		Message: "jobs check not implemented (stub)",
	}
}

// CheckAll runs all dependency health checks.
func CheckAll(ctx context.Context, pool *pgxpool.Pool, cacheClient *cache.Client, searchClient *search.Client, jobWorkers *jobs.Workers) map[string]CheckResult {
	return map[string]CheckResult{
		"database": CheckDatabase(ctx, pool),
		"cache":    CheckCache(ctx, cacheClient),
		"search":   CheckSearch(ctx, searchClient),
		"jobs":     CheckJobs(ctx, jobWorkers),
	}
}

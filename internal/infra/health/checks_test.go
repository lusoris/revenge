package health

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/jobs"
)

func TestCheckCache_NilClient(t *testing.T) {
	ctx := context.Background()

	result := CheckCache(ctx, nil)

	assert.Equal(t, "cache", result.Name)
	assert.Equal(t, StatusDegraded, result.Status)
	assert.Contains(t, result.Message, "not initialized")
}

func TestCheckCache_WithClient(t *testing.T) {
	ctx := context.Background()
	// Create an empty client (rueidis not connected)
	client := &cache.Client{}

	result := CheckCache(ctx, client)

	assert.Equal(t, "cache", result.Name)
	// Without actual rueidis connection, Ping will fail
	assert.Equal(t, StatusUnhealthy, result.Status)
}

func TestCheckJobs_NilClient(t *testing.T) {
	ctx := context.Background()

	result := CheckJobs(ctx, nil)

	assert.Equal(t, "jobs", result.Name)
	assert.Equal(t, StatusDegraded, result.Status)
	assert.Contains(t, result.Message, "not initialized")
}

func TestCheckJobs_NilRiverClient(t *testing.T) {
	ctx := context.Background()
	// Client exists but River client inside is nil
	client := &jobs.Client{}

	result := CheckJobs(ctx, client)

	assert.Equal(t, "jobs", result.Name)
	assert.Equal(t, StatusUnhealthy, result.Status)
	assert.Contains(t, result.Message, "river client not initialized")
}

func TestCheckDatabase_NilPool(t *testing.T) {
	ctx := context.Background()

	result := CheckDatabase(ctx, nil)

	assert.Equal(t, "database", result.Name)
	assert.Equal(t, StatusUnhealthy, result.Status)
	assert.Contains(t, result.Message, "not initialized")
}

func TestStatusConstants(t *testing.T) {
	// Test that status constants have expected values
	assert.Equal(t, Status("healthy"), StatusHealthy)
	assert.Equal(t, Status("unhealthy"), StatusUnhealthy)
	assert.Equal(t, Status("degraded"), StatusDegraded)
}

func TestCheckResult_Fields(t *testing.T) {
	result := CheckResult{
		Name:    "test",
		Status:  StatusHealthy,
		Message: "test message",
		Details: map[string]any{
			"key": "value",
		},
	}

	assert.Equal(t, "test", result.Name)
	assert.Equal(t, StatusHealthy, result.Status)
	assert.Equal(t, "test message", result.Message)
	assert.Equal(t, "value", result.Details["key"])
}

func TestCheckAll(t *testing.T) {
	ctx := context.Background()

	// Test with all nil dependencies
	results := CheckAll(ctx, nil, nil, nil)

	assert.Contains(t, results, "database")
	assert.Contains(t, results, "cache")
	assert.Contains(t, results, "jobs")

	// Nil pool means unhealthy database
	assert.Equal(t, StatusUnhealthy, results["database"].Status)
	// Nil cache client means degraded
	assert.Equal(t, StatusDegraded, results["cache"].Status)
	// Nil jobs client means degraded
	assert.Equal(t, StatusDegraded, results["jobs"].Status)
}

func TestCheckCache_Concurrent(t *testing.T) {
	ctx := context.Background()

	var wg sync.WaitGroup
	for range 100 {
		wg.Go(func() {
			result := CheckCache(ctx, nil)
			assert.Equal(t, StatusDegraded, result.Status)
		})
	}
	wg.Wait()
}

func TestCheckJobs_Concurrent(t *testing.T) {
	ctx := context.Background()

	var wg sync.WaitGroup
	for range 100 {
		wg.Go(func() {
			result := CheckJobs(ctx, nil)
			assert.Equal(t, StatusDegraded, result.Status)
		})
	}
	wg.Wait()
}

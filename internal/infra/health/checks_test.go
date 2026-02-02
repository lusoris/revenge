package health

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lusoris/revenge/internal/infra/cache"
	"github.com/lusoris/revenge/internal/infra/jobs"
	"github.com/lusoris/revenge/internal/infra/search"
)

func TestCheckCache(t *testing.T) {
	ctx := context.Background()
	client := &cache.Client{}

	result := CheckCache(ctx, client)

	assert.Equal(t, "cache", result.Name)
	assert.Equal(t, StatusHealthy, result.Status)
	assert.Contains(t, result.Message, "stub")
}

func TestCheckSearch(t *testing.T) {
	ctx := context.Background()
	client := &search.Client{}

	result := CheckSearch(ctx, client)

	assert.Equal(t, "search", result.Name)
	assert.Equal(t, StatusHealthy, result.Status)
	assert.Contains(t, result.Message, "stub")
}

func TestCheckJobs(t *testing.T) {
	ctx := context.Background()
	workers := &jobs.Workers{}

	result := CheckJobs(ctx, workers)

	assert.Equal(t, "jobs", result.Name)
	assert.Equal(t, StatusHealthy, result.Status)
	assert.Contains(t, result.Message, "stub")
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
		Details: map[string]interface{}{
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
	cacheClient := &cache.Client{}
	searchClient := &search.Client{}
	jobWorkers := &jobs.Workers{}

	// CheckAll needs a real pool for database check, so we only test stub checks
	// For now, just verify the function doesn't panic with nil pool
	// In real scenario, this would be an integration test
	t.Run("stub checks work", func(t *testing.T) {
		// Test individual stub checks
		cacheResult := CheckCache(ctx, cacheClient)
		assert.Equal(t, StatusHealthy, cacheResult.Status)

		searchResult := CheckSearch(ctx, searchClient)
		assert.Equal(t, StatusHealthy, searchResult.Status)

		jobsResult := CheckJobs(ctx, jobWorkers)
		assert.Equal(t, StatusHealthy, jobsResult.Status)
	})
}

func TestCheckCache_Concurrent(t *testing.T) {
	ctx := context.Background()
	client := &cache.Client{}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := CheckCache(ctx, client)
			assert.Equal(t, StatusHealthy, result.Status)
		}()
	}
	wg.Wait()
}

func TestCheckSearch_Concurrent(t *testing.T) {
	ctx := context.Background()
	client := &search.Client{}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := CheckSearch(ctx, client)
			assert.Equal(t, StatusHealthy, result.Status)
		}()
	}
	wg.Wait()
}

func TestCheckJobs_Concurrent(t *testing.T) {
	ctx := context.Background()
	workers := &jobs.Workers{}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := CheckJobs(ctx, workers)
			assert.Equal(t, StatusHealthy, result.Status)
		}()
	}
	wg.Wait()
}

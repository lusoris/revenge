//go:build integration

package integration

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/lusoris/revenge/internal/api/oas"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOgenClientGeneration(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	// Create ogen client
	client, err := oas.NewClient(ts.BaseURL, oas.WithClient(ts.HTTPClient))
	require.NoError(t, err, "should be able to create ogen client")
	require.NotNil(t, client, "client should not be nil")
}

func TestHealthEndpointsViaClient(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	// Create ogen client
	client, err := oas.NewClient(ts.BaseURL, oas.WithClient(ts.HTTPClient))
	require.NoError(t, err, "should be able to create ogen client")

	ctx := context.Background()

	// Test liveness endpoint
	t.Run("Liveness", func(t *testing.T) {
		resp, err := client.GetLiveness(ctx)
		require.NoError(t, err, "should be able to call liveness endpoint")

		// Type assertion to get the actual health check
		healthCheck, ok := resp.(*oas.HealthCheck)
		require.True(t, ok, "response should be *oas.HealthCheck")
		require.NotNil(t, healthCheck, "health check should not be nil")

		assert.Equal(t, "healthy", healthCheck.Status, "status should be healthy")
		assert.False(t, healthCheck.Timestamp.IsZero(), "timestamp should be set")
	})

	// Test readiness endpoint
	t.Run("Readiness", func(t *testing.T) {
		resp, err := client.GetReadiness(ctx)
		require.NoError(t, err, "should be able to call readiness endpoint")

		healthCheck, ok := resp.(*oas.HealthCheck)
		require.True(t, ok, "response should be *oas.HealthCheck")
		require.NotNil(t, healthCheck, "health check should not be nil")

		assert.Equal(t, "ready", healthCheck.Status, "status should be ready")
		assert.False(t, healthCheck.Timestamp.IsZero(), "timestamp should be set")

		// Verify checks are present
		require.NotNil(t, healthCheck.Checks.Value, "checks should be present")
		checks := healthCheck.Checks.Value
		assert.Contains(t, checks, "database", "should have database check")
	})

	// Test startup endpoint
	t.Run("Startup", func(t *testing.T) {
		resp, err := client.GetStartup(ctx)
		require.NoError(t, err, "should be able to call startup endpoint")

		healthCheck, ok := resp.(*oas.HealthCheck)
		require.True(t, ok, "response should be *oas.HealthCheck")
		require.NotNil(t, healthCheck, "health check should not be nil")

		assert.Equal(t, "started", healthCheck.Status, "status should be started")
		assert.False(t, healthCheck.Timestamp.IsZero(), "timestamp should be set")
	})
}

func TestClientTypeSafety(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	// Create ogen client
	client, err := oas.NewClient(ts.BaseURL, oas.WithClient(ts.HTTPClient))
	require.NoError(t, err, "should be able to create client")

	ctx := context.Background()

	// Call liveness endpoint
	resp, err := client.GetLiveness(ctx)
	require.NoError(t, err, "should be able to call endpoint")

	// Verify type safety - response should be HealthCheck
	healthCheck, ok := resp.(*oas.HealthCheck)
	require.True(t, ok, "response should be *oas.HealthCheck type")

	// Verify fields are accessible with type safety
	_ = healthCheck.Status    // string
	_ = healthCheck.Timestamp // time.Time
	_ = healthCheck.Checks    // OptHealthCheckChecks
	_ = healthCheck.Details   // OptHealthCheckDetails
}

func TestClientErrorHandling(t *testing.T) {
	// Create client pointing to non-existent server
	client, err := oas.NewClient("http://localhost:9999", oas.WithClient(&http.Client{
		Timeout: 1 * time.Second,
	}))
	require.NoError(t, err, "should be able to create client")

	ctx := context.Background()

	// Call endpoint (should fail with connection error)
	_, err = client.GetLiveness(ctx)
	assert.Error(t, err, "should get error when server is unreachable")
}

func TestClientConcurrentRequests(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	// Create ogen client
	client, err := oas.NewClient(ts.BaseURL, oas.WithClient(ts.HTTPClient))
	require.NoError(t, err, "should be able to create client")

	ctx := context.Background()

	// Make concurrent requests
	const numRequests = 10
	results := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			_, err := client.GetLiveness(ctx)
			results <- err
		}()
	}

	// Collect results
	for i := 0; i < numRequests; i++ {
		err := <-results
		assert.NoError(t, err, "concurrent request %d should succeed", i)
	}
}

func TestClientTimeout(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	// Create client with very short timeout
	client, err := oas.NewClient(ts.BaseURL, oas.WithClient(&http.Client{
		Timeout: 1 * time.Nanosecond, // Extremely short timeout
	}))
	require.NoError(t, err, "should be able to create client")

	ctx := context.Background()

	// Call endpoint (should timeout)
	_, err = client.GetLiveness(ctx)
	assert.Error(t, err, "should get timeout error")
}

func TestClientContextCancellation(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	// Create ogen client
	client, err := oas.NewClient(ts.BaseURL, oas.WithClient(ts.HTTPClient))
	require.NoError(t, err, "should be able to create client")

	// Create context with immediate cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Call endpoint (should fail with context cancelled)
	_, err = client.GetLiveness(ctx)
	assert.Error(t, err, "should get error when context is cancelled")
	assert.ErrorIs(t, err, context.Canceled, "error should be context.Canceled")
}

func TestClientWithContextTimeout(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	// Create ogen client
	client, err := oas.NewClient(ts.BaseURL, oas.WithClient(ts.HTTPClient))
	require.NoError(t, err, "should be able to create client")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Wait for timeout
	time.Sleep(10 * time.Millisecond)

	// Call endpoint (should fail with deadline exceeded)
	_, err = client.GetLiveness(ctx)
	assert.Error(t, err, "should get error when context times out")
}

func TestClientResponseValidation(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	// Create ogen client
	client, err := oas.NewClient(ts.BaseURL, oas.WithClient(ts.HTTPClient))
	require.NoError(t, err, "should be able to create client")

	ctx := context.Background()

	// Call endpoint and validate response structure
	resp, err := client.GetLiveness(ctx)
	require.NoError(t, err, "should be able to call endpoint")

	healthCheck, ok := resp.(*oas.HealthCheck)
	require.True(t, ok, "response should be *oas.HealthCheck")

	// Validate required fields are present
	assert.NotEmpty(t, healthCheck.Status, "status should not be empty")
	assert.False(t, healthCheck.Timestamp.IsZero(), "timestamp should not be zero")

	// Validate status is one of expected values
	validStatuses := []string{"healthy", "unhealthy", "degraded"}
	assert.Contains(t, validStatuses, healthCheck.Status, "status should be valid")

	// Validate timestamp is recent (within last minute)
	assert.WithinDuration(t, time.Now(), healthCheck.Timestamp, 1*time.Minute, "timestamp should be recent")
}

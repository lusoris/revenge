//go:build integration

package integration

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOgenClientGeneration(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	// Create ogen client
	client, err := ogen.NewClient(ts.BaseURL, ogen.WithClient(ts.HTTPClient))
	require.NoError(t, err, "should be able to create ogen client")
	require.NotNil(t, client, "client should not be nil")
}

func TestHealthEndpointsViaClient(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	// Create ogen client
	client, err := ogen.NewClient(ts.BaseURL, ogen.WithClient(ts.HTTPClient))
	require.NoError(t, err, "should be able to create ogen client")

	ctx := context.Background()

	// Test liveness endpoint
	t.Run("Liveness", func(t *testing.T) {
		resp, err := client.GetLiveness(ctx)
		require.NoError(t, err, "should be able to call liveness endpoint")
		require.NotNil(t, resp, "response should not be nil")

		// GetLiveness returns *HealthCheck directly
		assert.Equal(t, "liveness", resp.Name, "name should be liveness")
		assert.Equal(t, ogen.HealthCheckStatusHealthy, resp.Status, "status should be healthy")
		assert.True(t, resp.Message.IsSet(), "message should be set")
	})

	// Test readiness endpoint
	t.Run("Readiness", func(t *testing.T) {
		resp, err := client.GetReadiness(ctx)
		require.NoError(t, err, "should be able to call readiness endpoint")
		require.NotNil(t, resp, "response should not be nil")

		// GetReadiness returns GetReadinessRes interface (either OK or ServiceUnavailable)
		switch r := resp.(type) {
		case *ogen.GetReadinessOK:
			healthCheck := (*ogen.HealthCheck)(r)
			assert.Equal(t, "readiness", healthCheck.Name, "name should be readiness")
			assert.Equal(t, ogen.HealthCheckStatusHealthy, healthCheck.Status, "status should be healthy")
		case *ogen.GetReadinessServiceUnavailable:
			healthCheck := (*ogen.HealthCheck)(r)
			assert.Equal(t, "readiness", healthCheck.Name, "name should be readiness")
			assert.Equal(t, ogen.HealthCheckStatusUnhealthy, healthCheck.Status, "status should be unhealthy")
		default:
			t.Fatalf("unexpected response type: %T", resp)
		}
	})

	// Test startup endpoint
	t.Run("Startup", func(t *testing.T) {
		resp, err := client.GetStartup(ctx)
		require.NoError(t, err, "should be able to call startup endpoint")
		require.NotNil(t, resp, "response should not be nil")

		// GetStartup returns GetStartupRes interface (either OK or ServiceUnavailable)
		switch r := resp.(type) {
		case *ogen.GetStartupOK:
			healthCheck := (*ogen.HealthCheck)(r)
			assert.Equal(t, "startup", healthCheck.Name, "name should be startup")
			assert.Equal(t, ogen.HealthCheckStatusHealthy, healthCheck.Status, "status should be healthy")
		case *ogen.GetStartupServiceUnavailable:
			healthCheck := (*ogen.HealthCheck)(r)
			assert.Equal(t, "startup", healthCheck.Name, "name should be startup")
			assert.Equal(t, ogen.HealthCheckStatusUnhealthy, healthCheck.Status, "status should be unhealthy")
		default:
			t.Fatalf("unexpected response type: %T", resp)
		}
	})
}

func TestClientTypeSafety(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	// Create ogen client
	client, err := ogen.NewClient(ts.BaseURL, ogen.WithClient(ts.HTTPClient))
	require.NoError(t, err, "should be able to create client")

	ctx := context.Background()

	// Call liveness endpoint
	resp, err := client.GetLiveness(ctx)
	require.NoError(t, err, "should be able to call endpoint")
	require.NotNil(t, resp, "response should not be nil")

	// Verify fields are accessible with type safety
	assert.NotEmpty(t, resp.Name, "name should not be empty")
	assert.NotEmpty(t, resp.Status, "status should not be empty")
	// message and details are optional (OptString, OptHealthCheckDetails)
}

func TestClientErrorHandling(t *testing.T) {
	// Create client pointing to non-existent server
	client, err := ogen.NewClient("http://localhost:9999", ogen.WithClient(&http.Client{
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
	client, err := ogen.NewClient(ts.BaseURL, ogen.WithClient(ts.HTTPClient))
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
		assert.NoError(t, err, "concurrent request should succeed")
	}
}

func TestClientTimeout(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	// Create client with very short timeout
	client, err := ogen.NewClient(ts.BaseURL, ogen.WithClient(&http.Client{
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
	client, err := ogen.NewClient(ts.BaseURL, ogen.WithClient(ts.HTTPClient))
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
	client, err := ogen.NewClient(ts.BaseURL, ogen.WithClient(ts.HTTPClient))
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
	client, err := ogen.NewClient(ts.BaseURL, ogen.WithClient(ts.HTTPClient))
	require.NoError(t, err, "should be able to create client")

	ctx := context.Background()

	// Call endpoint and validate response structure
	resp, err := client.GetLiveness(ctx)
	require.NoError(t, err, "should be able to call endpoint")
	require.NotNil(t, resp, "response should not be nil")

	// Validate required fields are present
	assert.NotEmpty(t, resp.Name, "name should not be empty")
	assert.NotEmpty(t, resp.Status, "status should not be empty")

	// Validate status is one of expected values
	validStatuses := []ogen.HealthCheckStatus{
		ogen.HealthCheckStatusHealthy,
		ogen.HealthCheckStatusUnhealthy,
		ogen.HealthCheckStatusDegraded,
	}
	assert.Contains(t, validStatuses, resp.Status, "status should be valid")
}

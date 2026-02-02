//go:build integration

package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// HealthCheck matches the OpenAPI schema for health endpoints.
type HealthCheck struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Checks    map[string]Check  `json:"checks,omitempty"`
	Details   map[string]string `json:"details,omitempty"`
}

// Check represents an individual health check result.
type Check struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func TestHealthLivenessEndpoint(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	// Make request to liveness endpoint
	resp, err := ts.HTTPClient.Get(ts.BaseURL + "/health/live")
	require.NoError(t, err, "should be able to reach liveness endpoint")
	defer resp.Body.Close()

	// Verify status code
	assert.Equal(t, http.StatusOK, resp.StatusCode, "liveness should return 200 OK")

	// Verify Content-Type
	contentType := resp.Header.Get("Content-Type")
	assert.Contains(t, contentType, "application/json", "should return JSON")

	// Parse response body
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "should be able to read response body")

	var health HealthCheck
	err = json.Unmarshal(body, &health)
	require.NoError(t, err, "should be able to parse JSON response")

	// Verify response structure
	assert.Equal(t, "healthy", health.Status, "status should be healthy")
	assert.False(t, health.Timestamp.IsZero(), "timestamp should be set")
	assert.WithinDuration(t, time.Now(), health.Timestamp, 5*time.Second, "timestamp should be recent")
}

func TestHealthReadinessEndpointWithDB(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	// Make request to readiness endpoint
	resp, err := ts.HTTPClient.Get(ts.BaseURL + "/health/ready")
	require.NoError(t, err, "should be able to reach readiness endpoint")
	defer resp.Body.Close()

	// Verify status code
	assert.Equal(t, http.StatusOK, resp.StatusCode, "readiness should return 200 OK when DB is up")

	// Parse response body
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "should be able to read response body")

	var health HealthCheck
	err = json.Unmarshal(body, &health)
	require.NoError(t, err, "should be able to parse JSON response")

	// Verify response structure
	assert.Equal(t, "ready", health.Status, "status should be ready")
	assert.False(t, health.Timestamp.IsZero(), "timestamp should be set")

	// Verify database check is present
	require.NotNil(t, health.Checks, "checks should be present")
	dbCheck, exists := health.Checks["database"]
	require.True(t, exists, "database check should be present")
	assert.Equal(t, "healthy", dbCheck.Status, "database check should be healthy")
}

func TestHealthReadinessEndpointWithoutDB(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	// Close database connection to simulate DB failure
	ts.DB.Pool.Close()

	// Wait a moment for health check to detect the issue
	time.Sleep(100 * time.Millisecond)

	// Make request to readiness endpoint
	resp, err := ts.HTTPClient.Get(ts.BaseURL + "/health/ready")
	require.NoError(t, err, "should be able to reach readiness endpoint")
	defer resp.Body.Close()

	// Verify status code
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode, "readiness should return 503 when DB is down")

	// Parse response body
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "should be able to read response body")

	var health HealthCheck
	err = json.Unmarshal(body, &health)
	require.NoError(t, err, "should be able to parse JSON response")

	// Verify response structure
	assert.Equal(t, "not_ready", health.Status, "status should be not_ready")

	// Verify database check shows unhealthy
	require.NotNil(t, health.Checks, "checks should be present")
	dbCheck, exists := health.Checks["database"]
	require.True(t, exists, "database check should be present")
	assert.Equal(t, "unhealthy", dbCheck.Status, "database check should be unhealthy")
	assert.NotEmpty(t, dbCheck.Message, "error message should be present")
}

func TestHealthStartupEndpoint(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	// Make request to startup endpoint
	resp, err := ts.HTTPClient.Get(ts.BaseURL + "/health/startup")
	require.NoError(t, err, "should be able to reach startup endpoint")
	defer resp.Body.Close()

	// Verify status code
	assert.Equal(t, http.StatusOK, resp.StatusCode, "startup should return 200 OK after startup")

	// Parse response body
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "should be able to read response body")

	var health HealthCheck
	err = json.Unmarshal(body, &health)
	require.NoError(t, err, "should be able to parse JSON response")

	// Verify response structure
	assert.Equal(t, "started", health.Status, "status should be started")
	assert.False(t, health.Timestamp.IsZero(), "timestamp should be set")
}

func TestHealthResponseFormat(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	// Test all health endpoints for consistent format
	endpoints := []struct {
		path           string
		expectedStatus string
	}{
		{"/health/live", "healthy"},
		{"/health/ready", "ready"},
		{"/health/startup", "started"},
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint.path, func(t *testing.T) {
			resp, err := ts.HTTPClient.Get(ts.BaseURL + endpoint.path)
			require.NoError(t, err, "should be able to reach endpoint")
			defer resp.Body.Close()

			// Verify response format
			assert.Equal(t, http.StatusOK, resp.StatusCode, "should return 200 OK")

			contentType := resp.Header.Get("Content-Type")
			assert.Contains(t, contentType, "application/json", "should return JSON")

			// Parse and verify structure
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err, "should be able to read body")

			var health HealthCheck
			err = json.Unmarshal(body, &health)
			require.NoError(t, err, "should be able to parse JSON")

			// Verify required fields
			assert.NotEmpty(t, health.Status, "status should not be empty")
			assert.False(t, health.Timestamp.IsZero(), "timestamp should be set")

			// Verify status matches expected
			assert.Equal(t, endpoint.expectedStatus, health.Status, "status should match expected value")
		})
	}
}

func TestHealthConcurrentRequests(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	// Make concurrent requests to health endpoints
	const numRequests = 20
	results := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		go func(idx int) {
			endpoint := []string{"/health/live", "/health/ready", "/health/startup"}[idx%3]
			resp, err := ts.HTTPClient.Get(ts.BaseURL + endpoint)
			if err != nil {
				results <- err
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				results <- assert.AnError
				return
			}

			// Verify response is valid JSON
			var health HealthCheck
			err = json.NewDecoder(resp.Body).Decode(&health)
			results <- err
		}(i)
	}

	// Collect results
	for i := 0; i < numRequests; i++ {
		err := <-results
		assert.NoError(t, err, "concurrent request %d should succeed", i)
	}
}

func TestHealthEndpointsResponseTime(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	// Health endpoints should respond quickly (< 100ms)
	endpoints := []string{"/health/live", "/health/ready", "/health/startup"}

	for _, endpoint := range endpoints {
		t.Run(endpoint, func(t *testing.T) {
			start := time.Now()
			resp, err := ts.HTTPClient.Get(ts.BaseURL + endpoint)
			elapsed := time.Since(start)

			require.NoError(t, err, "should be able to reach endpoint")
			resp.Body.Close()

			assert.Less(t, elapsed, 100*time.Millisecond, "endpoint should respond within 100ms")
			assert.Equal(t, http.StatusOK, resp.StatusCode, "should return 200 OK")
		})
	}
}

func TestHealthEndpointsIdempotency(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	// Health endpoints should return consistent results
	endpoint := "/health/live"

	// Make multiple requests
	var responses []HealthCheck
	for i := 0; i < 3; i++ {
		resp, err := ts.HTTPClient.Get(ts.BaseURL + endpoint)
		require.NoError(t, err, "request %d should succeed", i)
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err, "should be able to read body")

		var health HealthCheck
		err = json.Unmarshal(body, &health)
		require.NoError(t, err, "should be able to parse JSON")

		responses = append(responses, health)

		// Small delay between requests
		time.Sleep(10 * time.Millisecond)
	}

	// All responses should have the same status
	for i := 1; i < len(responses); i++ {
		assert.Equal(t, responses[0].Status, responses[i].Status, "status should be consistent")
	}
}

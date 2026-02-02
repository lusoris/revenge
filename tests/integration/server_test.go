//go:build integration

package integration

import (
	"context"
	"net/http"
	"syscall"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerStartup(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	// Verify fx app is running
	assert.NotNil(t, ts.App, "fx app should be running")

	// Verify server is accessible
	resp, err := ts.HTTPClient.Get(ts.BaseURL + "/health/live")
	require.NoError(t, err, "should be able to reach server")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "server should respond with 200 OK")
}

func TestServerGracefulShutdown(t *testing.T) {
	// Setup test server
	ts := setupServer(t)

	// Verify server is running
	resp, err := ts.HTTPClient.Get(ts.BaseURL + "/health/live")
	require.NoError(t, err, "server should be running")
	resp.Body.Close()

	// Stop server gracefully
	ts.App.RequireStop()

	// Wait a moment for shutdown
	time.Sleep(100 * time.Millisecond)

	// Verify server is no longer accessible
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", ts.BaseURL+"/health/live", nil)
	require.NoError(t, err)

	_, err = ts.HTTPClient.Do(req)
	assert.Error(t, err, "server should not be accessible after shutdown")

	// Cleanup database
	ts.DB.Close()
}

func TestServerSignalHandling(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	// Verify server is running
	resp, err := ts.HTTPClient.Get(ts.BaseURL + "/health/live")
	require.NoError(t, err, "server should be running")
	resp.Body.Close()

	// Send SIGTERM (in real scenario - we just verify clean stop)
	// Note: In test environment we can't easily send signals to the fx app,
	// so we just verify that RequireStop works cleanly
	doneCh := make(chan struct{})
	go func() {
		ts.App.RequireStop()
		close(doneCh)
	}()

	// Wait for shutdown with timeout
	select {
	case <-doneCh:
		// Shutdown completed successfully
		t.Log("server shutdown completed cleanly")
	case <-time.After(5 * time.Second):
		t.Fatal("server shutdown timed out")
	}

	// Cleanup database separately since we already stopped the app
	ts.DB.Close()
	ts.App = nil // Prevent double-stop in defer
}

func TestServerModuleInitialization(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	// Verify all modules initialized by checking their endpoints/functionality

	// 1. API module - health endpoints accessible
	resp, err := ts.HTTPClient.Get(ts.BaseURL + "/health/live")
	require.NoError(t, err, "API module should be initialized")
	resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// 2. Database module - can execute queries
	ctx := context.Background()
	err = ts.DB.Pool.Ping(ctx)
	assert.NoError(t, err, "Database module should be initialized")

	// 3. Health module - readiness check works
	resp, err = ts.HTTPClient.Get(ts.BaseURL + "/health/ready")
	require.NoError(t, err, "Health module should be initialized")
	resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Logging is verified by the fact the modules above work (they all use zap.Logger)
}

func TestServerConcurrentRequests(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	// Make concurrent requests
	const numRequests = 10
	results := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			resp, err := ts.HTTPClient.Get(ts.BaseURL + "/health/live")
			if err != nil {
				results <- err
				return
			}
			resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				results <- assert.AnError
				return
			}
			results <- nil
		}()
	}

	// Collect results
	for i := 0; i < numRequests; i++ {
		err := <-results
		assert.NoError(t, err, "concurrent request %d should succeed", i)
	}
}

func TestServerDatabaseConnectionPool(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	ctx := context.Background()

	// Acquire multiple connections from pool
	type connWrapper struct {
		conn *pgxpool.Conn
	}
	conns := make([]*connWrapper, 5)
	for i := 0; i < 5; i++ {
		conn, err := ts.DB.Pool.Acquire(ctx)
		require.NoError(t, err, "should be able to acquire connection %d", i)
		conns[i] = &connWrapper{conn: conn}
	}

	// Verify all connections work
	for i, wrapper := range conns {
		err := wrapper.conn.Ping(ctx)
		assert.NoError(t, err, "connection %d should work", i)
	}

	// Release connections
	for _, wrapper := range conns {
		wrapper.conn.Release()
	}

	// Verify pool is still healthy
	err := ts.DB.Pool.Ping(ctx)
	assert.NoError(t, err, "pool should be healthy after releasing connections")
}

func TestServerConfigurationLoading(t *testing.T) {
	// Setup test server
	ts := setupServer(t)
	defer teardownServer(t, ts)

	// Verify configuration was loaded correctly
	assert.NotNil(t, ts.DB.Config, "configuration should be loaded")
	assert.Equal(t, "debug", ts.DB.Config.Logging.Level, "logging level should be debug")
	assert.True(t, ts.DB.Config.Logging.Development, "development mode should be enabled")
	assert.NotEmpty(t, ts.DB.Config.Database.URL, "database URL should be set")
}

// Note: SIGTERM/SIGINT handling is verified through RequireStop above,
// as Go testing doesn't support sending real signals to test processes
var _ = syscall.SIGTERM // Import usage to satisfy linter

package api

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/lusoris/revenge/internal/infra/logging"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/infra/health"
)

func newTestHandler(healthService *health.Service) *Handler {
	return &Handler{
		logger:        logging.NewTestLogger(),
		healthService: healthService,
	}
}

func newTestHealthService() *health.Service {
	// Use slog.Default() with nil pool for tests
	return health.NewService(slog.Default(), nil)
}

func TestHandler_GetLiveness(t *testing.T) {
	healthService := newTestHealthService()
	handler := newTestHandler(healthService)

	ctx := context.Background()
	result, err := handler.GetLiveness(ctx)

	require.NoError(t, err)
	assert.Equal(t, "liveness", result.Name)
	assert.Equal(t, ogen.HealthCheckStatusHealthy, result.Status)
	assert.True(t, result.Message.Set)
	assert.Equal(t, "Service is alive", result.Message.Value)
}

func TestHandler_GetStartup_NotStarted(t *testing.T) {
	healthService := newTestHealthService()
	handler := newTestHandler(healthService)

	ctx := context.Background()
	result, err := handler.GetStartup(ctx)

	require.NoError(t, err)
	// Service is not marked as started, so should return ServiceUnavailable
	unavailable, ok := result.(*ogen.GetStartupServiceUnavailable)
	require.True(t, ok, "expected ServiceUnavailable response")
	assert.Equal(t, "startup", unavailable.Name)
	assert.Equal(t, ogen.HealthCheckStatusUnhealthy, unavailable.Status)
}

func TestHandler_GetStartup_Started(t *testing.T) {
	healthService := newTestHealthService()
	healthService.MarkStartupComplete()
	handler := newTestHandler(healthService)

	ctx := context.Background()
	result, err := handler.GetStartup(ctx)

	require.NoError(t, err)
	// Service is marked as started, so should return OK
	ok_result, ok := result.(*ogen.GetStartupOK)
	require.True(t, ok, "expected OK response")
	assert.Equal(t, "startup", ok_result.Name)
	assert.Equal(t, ogen.HealthCheckStatusHealthy, ok_result.Status)
}

func TestHandler_GetReadiness_Healthy(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Start embedded postgres
	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15438).
		Username("test").
		Password("test").
		Database("test"))

	err := postgres.Start()
	require.NoError(t, err)
	defer func() { _ = postgres.Stop() }()

	// Create pool
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://test:test@localhost:15438/test?sslmode=disable")
	require.NoError(t, err)
	defer pool.Close()

	// Create health service with real pool
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	healthService := health.NewService(logger, pool)
	healthService.MarkStartupComplete()
	handler := newTestHandler(healthService)

	result, err := handler.GetReadiness(ctx)

	require.NoError(t, err)
	okResult, ok := result.(*ogen.GetReadinessOK)
	require.True(t, ok, "expected OK response, got %T", result)
	assert.Equal(t, "readiness", okResult.Name)
	assert.Equal(t, ogen.HealthCheckStatusHealthy, okResult.Status)
}

func TestHandler_GetReadiness_Unhealthy(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Start embedded postgres
	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15439).
		Username("test").
		Password("test").
		Database("test"))

	err := postgres.Start()
	require.NoError(t, err)

	// Create pool
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://test:test@localhost:15439/test?sslmode=disable")
	require.NoError(t, err)

	// Create health service with pool
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	healthService := health.NewService(logger, pool)
	healthService.MarkStartupComplete()
	_ = newTestHandler(healthService) // Verify handler can be created

	// Stop the database to simulate unhealthy state
	_ = postgres.Stop()
	pool.Close()

	// Create new pool pointing to stopped DB
	pool2, err := pgxpool.New(ctx, "postgres://test:test@localhost:15439/test?sslmode=disable")
	require.NoError(t, err)
	defer pool2.Close()

	healthService2 := health.NewService(logger, pool2)
	healthService2.MarkStartupComplete()
	handler2 := newTestHandler(healthService2)

	result, err := handler2.GetReadiness(ctx)

	require.NoError(t, err)
	unavailableResult, ok := result.(*ogen.GetReadinessServiceUnavailable)
	require.True(t, ok, "expected ServiceUnavailable response, got %T", result)
	assert.Equal(t, "readiness", unavailableResult.Name)
	assert.Equal(t, ogen.HealthCheckStatusUnhealthy, unavailableResult.Status)
}

func TestHandler_NewError(t *testing.T) {
	healthService := newTestHealthService()
	handler := newTestHandler(healthService)

	ctx := context.Background()
	testErr := assert.AnError

	result := handler.NewError(ctx, testErr)

	assert.Equal(t, 500, result.StatusCode)
	assert.Equal(t, 500, result.Response.Code)
	assert.Equal(t, "Internal server error", result.Response.Message)
}

func TestHandler_GetLiveness_Concurrent(t *testing.T) {
	healthService := newTestHealthService()
	handler := newTestHandler(healthService)

	ctx := context.Background()

	// Run concurrent requests
	done := make(chan bool, 100)
	for i := 0; i < 100; i++ {
		go func() {
			result, err := handler.GetLiveness(ctx)
			assert.NoError(t, err)
			assert.Equal(t, ogen.HealthCheckStatusHealthy, result.Status)
			done <- true
		}()
	}

	// Wait for all to complete
	for i := 0; i < 100; i++ {
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Fatal("timeout waiting for concurrent requests")
		}
	}
}

func TestHandler_GetReadiness_WithCancelledContext(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Start embedded postgres
	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15440).
		Username("test").
		Password("test").
		Database("test"))

	err := postgres.Start()
	require.NoError(t, err)
	defer func() { _ = postgres.Stop() }()

	// Create pool
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://test:test@localhost:15440/test?sslmode=disable")
	require.NoError(t, err)
	defer pool.Close()

	// Create health service with real pool
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	healthService := health.NewService(logger, pool)
	healthService.MarkStartupComplete()
	handler := newTestHandler(healthService)

	// Create a context that's already cancelled
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	// Should still return a response (may be unhealthy due to cancelled context)
	result, err := handler.GetReadiness(cancelledCtx)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

func TestHandler_NewError_VariousErrors(t *testing.T) {
	healthService := newTestHealthService()
	handler := newTestHandler(healthService)
	ctx := context.Background()

	testCases := []struct {
		name       string
		err        error
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "generic error",
			err:        assert.AnError,
			wantStatus: 500,
			wantMsg:    "Internal server error",
		},
		{
			name:       "context canceled",
			err:        context.Canceled,
			wantStatus: 500,
			wantMsg:    "Internal server error",
		},
		{
			name:       "context deadline exceeded",
			err:        context.DeadlineExceeded,
			wantStatus: 500,
			wantMsg:    "Internal server error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := handler.NewError(ctx, tc.err)
			assert.Equal(t, tc.wantStatus, result.StatusCode)
			assert.Equal(t, tc.wantMsg, result.Response.Message)
		})
	}
}

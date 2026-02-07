package health_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/infra/health"
	"github.com/lusoris/revenge/internal/testutil"
)

func TestNewService(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	pool := &pgxpool.Pool{} // Mock pool

	service := health.NewService(logger, pool)

	require.NotNil(t, service)
}

func TestService_Liveness(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	pool := &pgxpool.Pool{}
	service := health.NewService(logger, pool)

	ctx := context.Background()
	result := service.Liveness(ctx)

	assert.Equal(t, "liveness", result.Name)
	assert.Equal(t, health.StatusHealthy, result.Status)
	assert.NotEmpty(t, result.Message)
}

func TestService_Startup(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	pool := &pgxpool.Pool{}
	service := health.NewService(logger, pool)

	ctx := context.Background()

	t.Run("before startup complete", func(t *testing.T) {
		result := service.Startup(ctx)

		assert.Equal(t, "startup", result.Name)
		assert.Equal(t, health.StatusUnhealthy, result.Status)
		assert.Contains(t, result.Message, "initialization in progress")
	})

	t.Run("after startup complete", func(t *testing.T) {
		service.MarkStartupComplete()
		result := service.Startup(ctx)

		assert.Equal(t, "startup", result.Name)
		assert.Equal(t, health.StatusHealthy, result.Status)
		assert.Contains(t, result.Message, "startup complete")
	})
}

func TestService_Readiness(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	t.Run("unhealthy before startup complete", func(t *testing.T) {
		pool := &pgxpool.Pool{}
		service := health.NewService(logger, pool)

		ctx := context.Background()
		result := service.Readiness(ctx)

		assert.Equal(t, "readiness", result.Name)
		assert.Equal(t, health.StatusUnhealthy, result.Status)
		assert.Contains(t, result.Message, "startup not complete")
	})

	t.Run("with real database", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping database integration test in short mode")
		}

		// Use shared test DB infrastructure instead of starting own embedded postgres
		testDB := testutil.NewTestDB(t)

		service := health.NewService(logger, testDB.Pool())
		service.MarkStartupComplete()

		ctx := context.Background()
		result := service.Readiness(ctx)

		assert.Equal(t, "readiness", result.Name)
		assert.Equal(t, health.StatusHealthy, result.Status)
		assert.Contains(t, result.Message, "service is ready")
		assert.NotNil(t, result.Details)
		assert.Contains(t, result.Details, "database")
	})
}

func TestService_FullCheck(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping FullCheck integration test in short mode")
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Use shared test DB infrastructure instead of starting own embedded postgres
	testDB := testutil.NewTestDB(t)

	service := health.NewService(logger, testDB.Pool())
	service.MarkStartupComplete()

	ctx := context.Background()
	results := service.FullCheck(ctx)

	require.NotNil(t, results)
	assert.Contains(t, results, "liveness")
	assert.Contains(t, results, "readiness")
	assert.Contains(t, results, "startup")

	// Check that database check was performed
	readiness := results["readiness"]
	assert.Equal(t, health.StatusHealthy, readiness.Status)
}

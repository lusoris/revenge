package health_test

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/fergusstrange/embedded-postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/infra/health"
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

func TestService_MarkStartupComplete(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	pool := &pgxpool.Pool{}
	service := health.NewService(logger, pool)

	ctx := context.Background()

	// Should be unhealthy initially
	result := service.Startup(ctx)
	assert.Equal(t, health.StatusUnhealthy, result.Status)

	// Mark startup complete
	service.MarkStartupComplete()

	// Should be healthy now
	result = service.Startup(ctx)
	assert.Equal(t, health.StatusHealthy, result.Status)
}

func TestService_Readiness(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	t.Run("before startup complete", func(t *testing.T) {
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

		// Start embedded PostgreSQL
		pg := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
			Port(5433).
			Database("test_health").
			Username("test").
			Password("test"))

		err := pg.Start()
		require.NoError(t, err)
		defer func() {
			_ = pg.Stop()
		}()

		// Wait for database to be ready
		time.Sleep(2 * time.Second)

		// Create connection pool
		ctx := context.Background()
		connStr := "postgres://test:test@localhost:5433/test_health?sslmode=disable"
		pool, err := pgxpool.New(ctx, connStr)
		require.NoError(t, err)
		defer pool.Close()

		service := health.NewService(logger, pool)
		service.MarkStartupComplete()

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

	// Start embedded PostgreSQL
	pg := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(5434).
		Database("test_fullcheck").
		Username("test").
		Password("test"))

	err := pg.Start()
	require.NoError(t, err)
	defer func() {
		_ = pg.Stop()
	}()

	// Wait for database to be ready
	time.Sleep(2 * time.Second)

	// Create connection pool
	ctx := context.Background()
	connStr := "postgres://test:test@localhost:5434/test_fullcheck?sslmode=disable"
	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)
	defer pool.Close()

	service := health.NewService(logger, pool)
	service.MarkStartupComplete()

	results := service.FullCheck(ctx)

	require.NotNil(t, results)
	assert.Contains(t, results, "liveness")
	assert.Contains(t, results, "readiness")
	assert.Contains(t, results, "startup")

	// Liveness should always be healthy
	assert.Equal(t, health.StatusHealthy, results["liveness"].Status)

	// Startup should be healthy after marking complete
	assert.Equal(t, health.StatusHealthy, results["startup"].Status)

	// Readiness should be healthy with working database
	assert.Equal(t, health.StatusHealthy, results["readiness"].Status)
}

func TestStatus_Constants(t *testing.T) {
	assert.Equal(t, health.Status("healthy"), health.StatusHealthy)
	assert.Equal(t, health.Status("unhealthy"), health.StatusUnhealthy)
	assert.Equal(t, health.Status("degraded"), health.StatusDegraded)
}

func TestCheckResult_Structure(t *testing.T) {
	result := health.CheckResult{
		Name:    "test",
		Status:  health.StatusHealthy,
		Message: "test message",
		Details: map[string]interface{}{
			"key": "value",
		},
	}

	assert.Equal(t, "test", result.Name)
	assert.Equal(t, health.StatusHealthy, result.Status)
	assert.Equal(t, "test message", result.Message)
	assert.Equal(t, "value", result.Details["key"])
}

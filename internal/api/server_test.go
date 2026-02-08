package api

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"testing"
	"time"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"github.com/lusoris/revenge/internal/infra/logging"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/health"
)

// testConfig creates a minimal config for testing
func testConfig(port int) *config.Config {
	return &config.Config{
		Server: config.ServerConfig{
			Host:         "127.0.0.1",
			Port:         port,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}
}

func TestNewServer_WithFxLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Start embedded postgres for health service
	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15450).
		Username("test").
		Password("test").
		Database("test"))

	err := postgres.Start()
	require.NoError(t, err)
	defer func() { _ = postgres.Stop() }()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://test:test@localhost:15450/test?sslmode=disable")
	require.NoError(t, err)
	defer pool.Close()

	logger := logging.NewTestLogger()
	healthService := health.NewService(logger, pool)

	// Use fxtest for proper lifecycle testing
	var server *Server
	_ = server

	app := fxtest.New(t,
		fx.Provide(
			func() *config.Config { return testConfig(15451) },
			func() *slog.Logger { return logger },
			func() *health.Service { return healthService },
		),
		fx.Invoke(func(lc fx.Lifecycle, cfg *config.Config, log *slog.Logger, hs *health.Service) error {
			var err error
			server, err = NewServer(ServerParams{
				Config:        cfg,
				Logger:        log,
				HealthService: hs,
				Lifecycle:     lc,
			})
			return err
		}),
	)

	// Start the app (triggers OnStart hooks)
	app.RequireStart()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Verify server was created
	require.NotNil(t, server)
	assert.NotNil(t, server.httpServer)
	assert.NotNil(t, server.ogenServer)

	// Test that server is actually listening
	resp, err := http.Get("http://127.0.0.1:15451/healthz")
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Stop the app (triggers OnStop hooks)
	app.RequireStop()
}

func TestNewServer_StartsAndStops(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Start embedded postgres
	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15452).
		Username("test").
		Password("test").
		Database("test"))

	err := postgres.Start()
	require.NoError(t, err)
	defer func() { _ = postgres.Stop() }()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://test:test@localhost:15452/test?sslmode=disable")
	require.NoError(t, err)
	defer pool.Close()

	logger := logging.NewTestLogger()
	healthService := health.NewService(logger, pool)
	healthService.MarkStartupComplete()

	var server *Server
	_ = server

	app := fxtest.New(t,
		fx.Provide(
			func() *config.Config { return testConfig(15453) },
			func() *slog.Logger { return logger },
			func() *health.Service { return healthService },
		),
		fx.Invoke(func(lc fx.Lifecycle, cfg *config.Config, log *slog.Logger, hs *health.Service) error {
			var err error
			server, err = NewServer(ServerParams{
				Config:        cfg,
				Logger:        log,
				HealthService: hs,
				Lifecycle:     lc,
			})
			return err
		}),
	)

	app.RequireStart()
	time.Sleep(100 * time.Millisecond)

	// Test all health endpoints
	t.Run("liveness endpoint", func(t *testing.T) {
		resp, err := http.Get("http://127.0.0.1:15453/healthz")
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("readiness endpoint", func(t *testing.T) {
		resp, err := http.Get("http://127.0.0.1:15453/readyz")
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("startup endpoint", func(t *testing.T) {
		resp, err := http.Get("http://127.0.0.1:15453/startupz")
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	app.RequireStop()

	// After stop, server should not respond
	time.Sleep(100 * time.Millisecond)
	_, err = http.Get("http://127.0.0.1:15453/healthz")
	assert.Error(t, err) // Connection should be refused
}

func TestNewServer_ReadinessUnhealthyWhenDBDown(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Start embedded postgres
	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15454).
		Username("test").
		Password("test").
		Database("test"))

	err := postgres.Start()
	require.NoError(t, err)

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://test:test@localhost:15454/test?sslmode=disable")
	require.NoError(t, err)

	logger := logging.NewTestLogger()
	healthService := health.NewService(logger, pool)
	healthService.MarkStartupComplete()

	var server *Server
	_ = server

	app := fxtest.New(t,
		fx.Provide(
			func() *config.Config { return testConfig(15455) },
			func() *slog.Logger { return logger },
			func() *health.Service { return healthService },
		),
		fx.Invoke(func(lc fx.Lifecycle, cfg *config.Config, log *slog.Logger, hs *health.Service) error {
			var err error
			server, err = NewServer(ServerParams{
				Config:        cfg,
				Logger:        log,
				HealthService: hs,
				Lifecycle:     lc,
			})
			return err
		}),
	)

	app.RequireStart()
	time.Sleep(100 * time.Millisecond)

	// First verify readiness is healthy
	resp, err := http.Get("http://127.0.0.1:15455/readyz")
	require.NoError(t, err)
	_ = resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Now stop postgres and close pool
	pool.Close()
	_ = postgres.Stop()

	// Give it a moment
	time.Sleep(100 * time.Millisecond)

	// Readiness should now return 503 Service Unavailable
	resp, err = http.Get("http://127.0.0.1:15455/readyz")
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

	// But liveness should still work (process is alive)
	resp2, err := http.Get("http://127.0.0.1:15455/healthz")
	require.NoError(t, err)
	defer func() { _ = resp2.Body.Close() }()
	assert.Equal(t, http.StatusOK, resp2.StatusCode)

	app.RequireStop()
}

func TestNewServer_StartupUnhealthyBeforeMarkComplete(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Start embedded postgres
	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15456).
		Username("test").
		Password("test").
		Database("test"))

	err := postgres.Start()
	require.NoError(t, err)
	defer func() { _ = postgres.Stop() }()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://test:test@localhost:15456/test?sslmode=disable")
	require.NoError(t, err)
	defer pool.Close()

	logger := logging.NewTestLogger()
	healthService := health.NewService(logger, pool)
	// Note: NOT calling MarkStartupComplete()

	var server *Server
	_ = server

	app := fxtest.New(t,
		fx.Provide(
			func() *config.Config { return testConfig(15457) },
			func() *slog.Logger { return logger },
			func() *health.Service { return healthService },
		),
		fx.Invoke(func(lc fx.Lifecycle, cfg *config.Config, log *slog.Logger, hs *health.Service) error {
			var err error
			server, err = NewServer(ServerParams{
				Config:        cfg,
				Logger:        log,
				HealthService: hs,
				Lifecycle:     lc,
			})
			return err
		}),
	)

	app.RequireStart()
	time.Sleep(100 * time.Millisecond)

	// Startup should be unhealthy (503)
	resp, err := http.Get("http://127.0.0.1:15457/startupz")
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

	// Mark startup complete
	healthService.MarkStartupComplete()

	// Now startup should be healthy
	resp2, err := http.Get("http://127.0.0.1:15457/startupz")
	require.NoError(t, err)
	defer func() { _ = resp2.Body.Close() }()
	assert.Equal(t, http.StatusOK, resp2.StatusCode)

	app.RequireStop()
}

func TestNewServer_ResponseBodyContent(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Start embedded postgres
	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15458).
		Username("test").
		Password("test").
		Database("test"))

	err := postgres.Start()
	require.NoError(t, err)
	defer func() { _ = postgres.Stop() }()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://test:test@localhost:15458/test?sslmode=disable")
	require.NoError(t, err)
	defer pool.Close()

	logger := logging.NewTestLogger()
	healthService := health.NewService(logger, pool)
	healthService.MarkStartupComplete()

	var server *Server
	_ = server

	app := fxtest.New(t,
		fx.Provide(
			func() *config.Config { return testConfig(15459) },
			func() *slog.Logger { return logger },
			func() *health.Service { return healthService },
		),
		fx.Invoke(func(lc fx.Lifecycle, cfg *config.Config, log *slog.Logger, hs *health.Service) error {
			var err error
			server, err = NewServer(ServerParams{
				Config:        cfg,
				Logger:        log,
				HealthService: hs,
				Lifecycle:     lc,
			})
			return err
		}),
	)

	app.RequireStart()
	time.Sleep(100 * time.Millisecond)

	t.Run("liveness response contains correct JSON", func(t *testing.T) {
		resp, err := http.Get("http://127.0.0.1:15459/healthz")
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		bodyStr := string(body)
		assert.Contains(t, bodyStr, `"name":"liveness"`)
		assert.Contains(t, bodyStr, `"status":"healthy"`)
	})

	t.Run("readiness response contains correct JSON", func(t *testing.T) {
		resp, err := http.Get("http://127.0.0.1:15459/readyz")
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		bodyStr := string(body)
		assert.Contains(t, bodyStr, `"name":"readiness"`)
		assert.Contains(t, bodyStr, `"status":"healthy"`)
	})

	t.Run("startup response contains correct JSON", func(t *testing.T) {
		resp, err := http.Get("http://127.0.0.1:15459/startupz")
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		bodyStr := string(body)
		assert.Contains(t, bodyStr, `"name":"startup"`)
		assert.Contains(t, bodyStr, `"status":"healthy"`)
	})

	app.RequireStop()
}

func TestNewServer_ConcurrentRequests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Start embedded postgres
	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15460).
		Username("test").
		Password("test").
		Database("test"))

	err := postgres.Start()
	require.NoError(t, err)
	defer func() { _ = postgres.Stop() }()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://test:test@localhost:15460/test?sslmode=disable")
	require.NoError(t, err)
	defer pool.Close()

	logger := logging.NewTestLogger()
	healthService := health.NewService(logger, pool)
	healthService.MarkStartupComplete()

	var server *Server
	_ = server

	app := fxtest.New(t,
		fx.Provide(
			func() *config.Config { return testConfig(15461) },
			func() *slog.Logger { return logger },
			func() *health.Service { return healthService },
		),
		fx.Invoke(func(lc fx.Lifecycle, cfg *config.Config, log *slog.Logger, hs *health.Service) error {
			var err error
			server, err = NewServer(ServerParams{
				Config:        cfg,
				Logger:        log,
				HealthService: hs,
				Lifecycle:     lc,
			})
			return err
		}),
	)

	app.RequireStart()
	time.Sleep(100 * time.Millisecond)

	// Send 50 concurrent requests
	const numRequests = 50
	results := make(chan int, numRequests)
	errors := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		go func() {
			resp, err := http.Get("http://127.0.0.1:15461/healthz")
			if err != nil {
				errors <- err
				return
			}
			defer func() { _ = resp.Body.Close() }()
			results <- resp.StatusCode
		}()
	}

	// Collect results
	for i := 0; i < numRequests; i++ {
		select {
		case status := <-results:
			assert.Equal(t, http.StatusOK, status)
		case err := <-errors:
			t.Errorf("Request failed: %v", err)
		case <-time.After(10 * time.Second):
			t.Fatal("Timeout waiting for concurrent requests")
		}
	}

	app.RequireStop()
}

func TestNewServer_ServerConfigApplied(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Start embedded postgres
	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15462).
		Username("test").
		Password("test").
		Database("test"))

	err := postgres.Start()
	require.NoError(t, err)
	defer func() { _ = postgres.Stop() }()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://test:test@localhost:15462/test?sslmode=disable")
	require.NoError(t, err)
	defer pool.Close()

	logger := logging.NewTestLogger()
	healthService := health.NewService(logger, pool)

	// Create config with specific timeouts
	cfg := &config.Config{
		Server: config.ServerConfig{
			Host:         "127.0.0.1",
			Port:         15463,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
	}

	var server *Server
	_ = server

	app := fxtest.New(t,
		fx.Provide(
			func() *config.Config { return cfg },
			func() *slog.Logger { return logger },
			func() *health.Service { return healthService },
		),
		fx.Invoke(func(lc fx.Lifecycle, cfg *config.Config, log *slog.Logger, hs *health.Service) error {
			var err error
			server, err = NewServer(ServerParams{
				Config:        cfg,
				Logger:        log,
				HealthService: hs,
				Lifecycle:     lc,
			})
			return err
		}),
	)

	app.RequireStart()
	time.Sleep(100 * time.Millisecond)

	// Verify server configuration was applied
	require.NotNil(t, server)
	assert.Equal(t, "127.0.0.1:15463", server.httpServer.Addr)
	assert.Equal(t, 15*time.Second, server.httpServer.ReadTimeout)
	assert.Equal(t, 30*time.Second, server.httpServer.WriteTimeout)
	assert.Equal(t, 120*time.Second, server.httpServer.IdleTimeout)

	app.RequireStop()
}

func TestNewServer_GracefulShutdown(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Start embedded postgres
	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15464).
		Username("test").
		Password("test").
		Database("test"))

	err := postgres.Start()
	require.NoError(t, err)
	defer func() { _ = postgres.Stop() }()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://test:test@localhost:15464/test?sslmode=disable")
	require.NoError(t, err)
	defer pool.Close()

	logger := logging.NewTestLogger()
	healthService := health.NewService(logger, pool)

	var server *Server
	_ = server

	app := fxtest.New(t,
		fx.Provide(
			func() *config.Config { return testConfig(15465) },
			func() *slog.Logger { return logger },
			func() *health.Service { return healthService },
		),
		fx.Invoke(func(lc fx.Lifecycle, cfg *config.Config, log *slog.Logger, hs *health.Service) error {
			var err error
			server, err = NewServer(ServerParams{
				Config:        cfg,
				Logger:        log,
				HealthService: hs,
				Lifecycle:     lc,
			})
			return err
		}),
	)

	app.RequireStart()
	time.Sleep(100 * time.Millisecond)

	// Verify server is running
	resp, err := http.Get("http://127.0.0.1:15465/healthz")
	require.NoError(t, err)
	_ = resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Stop should complete gracefully (not hang)
	stopDone := make(chan struct{})
	go func() {
		app.RequireStop()
		close(stopDone)
	}()

	select {
	case <-stopDone:
		// Good - shutdown completed
	case <-time.After(5 * time.Second):
		t.Fatal("Server shutdown timed out - not graceful")
	}
}

func TestNewServer_MultiplePortsInSequence(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// This test ensures we can start/stop servers on different ports sequentially
	// without port conflicts

	// Start embedded postgres (shared)
	postgres := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
		Port(15466).
		Username("test").
		Password("test").
		Database("test"))

	err := postgres.Start()
	require.NoError(t, err)
	defer func() { _ = postgres.Stop() }()

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://test:test@localhost:15466/test?sslmode=disable")
	require.NoError(t, err)
	defer pool.Close()

	logger := logging.NewTestLogger()

	ports := []int{15467, 15468, 15469}

	for _, port := range ports {
		t.Run(fmt.Sprintf("port_%d", port), func(t *testing.T) {
			healthService := health.NewService(logger, pool)

			var server *Server
			_ = server

			app := fxtest.New(t,
				fx.Provide(
					func() *config.Config { return testConfig(port) },
					func() *slog.Logger { return logger },
					func() *health.Service { return healthService },
				),
				fx.Invoke(func(lc fx.Lifecycle, cfg *config.Config, log *slog.Logger, hs *health.Service) error {
					var err error
					server, err = NewServer(ServerParams{
						Config:        cfg,
						Logger:        log,
						HealthService: hs,
						Lifecycle:     lc,
					})
					return err
				}),
			)

			app.RequireStart()
			time.Sleep(100 * time.Millisecond)

			// Verify server is running
			resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/healthz", port))
			require.NoError(t, err)
			_ = resp.Body.Close()
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			app.RequireStop()
			time.Sleep(100 * time.Millisecond) // Ensure port is released
		})
	}
}

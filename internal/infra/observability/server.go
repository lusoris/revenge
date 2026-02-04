package observability

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/lusoris/revenge/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// ServerParams defines dependencies for the observability server.
type ServerParams struct {
	fx.In

	Config    *config.Config
	Logger    *zap.Logger
	Lifecycle fx.Lifecycle
}

// Server provides observability endpoints (metrics, pprof).
type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
}

// NewServer creates a new observability server.
// It exposes:
// - /metrics for Prometheus scraping
// - /debug/pprof/* for profiling (only in development mode)
// - /health/live and /health/ready for k8s probes
func NewServer(p ServerParams) *Server {
	mux := http.NewServeMux()

	// Always expose metrics endpoint
	mux.Handle("/metrics", promhttp.Handler())

	// Health check endpoints (simple, for the observability port)
	mux.HandleFunc("/health/live", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
	mux.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Only expose pprof in development mode
	if p.Config.Logging.Development {
		RegisterPprofHandlers(mux)
		p.Logger.Info("pprof endpoints enabled (development mode)")
	}

	// Use a different port for observability (main port + 1000 or configurable)
	port := p.Config.Server.Port + 1000
	addr := fmt.Sprintf("%s:%d", p.Config.Server.Host, port)

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second, // Prevent slowloris attacks
	}

	server := &Server{
		httpServer: httpServer,
		logger:     p.Logger.Named("observability"),
	}

	// Register lifecycle hooks
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				server.logger.Info("Starting observability server",
					zap.String("address", addr),
					zap.Bool("pprof_enabled", p.Config.Logging.Development),
				)
				if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					server.logger.Error("Observability server error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			server.logger.Info("Stopping observability server")
			return httpServer.Shutdown(ctx)
		},
	})

	return server
}

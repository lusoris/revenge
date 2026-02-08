package observability

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/fx"
)

// Infrastructure gauges collected periodically.
var (
	riverQueueSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "revenge",
		Subsystem: "river",
		Name:      "queue_size",
		Help:      "Number of River jobs by state.",
	}, []string{"state"})

	// pgxpool gauges (collected periodically, not per-query)
	pgxPoolAcquiredConns = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "revenge",
		Subsystem: "pgxpool",
		Name:      "acquired_conns",
		Help:      "Number of currently acquired connections in the pool.",
	})
	pgxPoolIdleConns = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "revenge",
		Subsystem: "pgxpool",
		Name:      "idle_conns",
		Help:      "Number of idle connections in the pool.",
	})
	pgxPoolTotalConns = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "revenge",
		Subsystem: "pgxpool",
		Name:      "total_conns",
		Help:      "Total number of connections in the pool.",
	})
	pgxPoolMaxConns = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "revenge",
		Subsystem: "pgxpool",
		Name:      "max_conns",
		Help:      "Maximum number of connections allowed in the pool.",
	})
	pgxPoolConstructingConns = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "revenge",
		Subsystem: "pgxpool",
		Name:      "constructing_conns",
		Help:      "Number of connections being constructed.",
	})
)

// CollectorParams contains dependencies for the periodic metrics collector.
type CollectorParams struct {
	fx.In

	Pool      *pgxpool.Pool
	Logger    *slog.Logger
	Lifecycle fx.Lifecycle
}

// StartCollector registers a periodic metrics collector that runs every 15 seconds.
// It collects pgxpool stats and River queue depth.
func StartCollector(p CollectorParams) {
	logger := p.Logger.With("component", "metrics-collector")

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go runCollector(p.Pool, logger)
			return nil
		},
	})
}

func runCollector(pool *pgxpool.Pool, logger *slog.Logger) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		collectPoolStats(pool)
		collectRiverStats(pool, logger)
	}
}

func collectPoolStats(pool *pgxpool.Pool) {
	stat := pool.Stat()
	pgxPoolAcquiredConns.Set(float64(stat.AcquiredConns()))
	pgxPoolIdleConns.Set(float64(stat.IdleConns()))
	pgxPoolTotalConns.Set(float64(stat.TotalConns()))
	pgxPoolMaxConns.Set(float64(stat.MaxConns()))
	pgxPoolConstructingConns.Set(float64(stat.ConstructingConns()))
}

func collectRiverStats(pool *pgxpool.Pool, logger *slog.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := pool.Query(ctx,
		`SELECT state, COUNT(*) FROM river_job GROUP BY state`)
	if err != nil {
		logger.Debug("failed to query river_job stats", slog.Any("error", err))
		return
	}
	defer rows.Close()

	// Reset all known states to 0 first so disappeared states show 0
	for _, state := range []string{"available", "running", "retryable", "scheduled", "completed", "discarded", "cancelled", "pending"} {
		riverQueueSize.WithLabelValues(state).Set(0)
	}

	for rows.Next() {
		var state string
		var count int64
		if err := rows.Scan(&state, &count); err != nil {
			continue
		}
		riverQueueSize.WithLabelValues(state).Set(float64(count))
	}
}

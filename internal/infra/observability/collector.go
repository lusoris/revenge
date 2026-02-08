package observability

import (
	"context"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/redis/rueidis"
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

	// Dragonfly/Redis gauges
	dragonflyUsedMemory = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "revenge",
		Subsystem: "dragonfly",
		Name:      "used_memory_bytes",
		Help:      "Memory used by Dragonfly in bytes.",
	})
	dragonflyMaxMemory = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "revenge",
		Subsystem: "dragonfly",
		Name:      "max_memory_bytes",
		Help:      "Maximum memory configured for Dragonfly in bytes.",
	})
	dragonflyConnectedClients = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "revenge",
		Subsystem: "dragonfly",
		Name:      "connected_clients",
		Help:      "Number of connected clients to Dragonfly.",
	})
	dragonflyTotalKeys = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "revenge",
		Subsystem: "dragonfly",
		Name:      "total_keys",
		Help:      "Total number of keys stored in Dragonfly.",
	})
	dragonflyCommandsProcessed = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "revenge",
		Subsystem: "dragonfly",
		Name:      "commands_processed_total",
		Help:      "Total number of commands processed by Dragonfly.",
	})
	dragonflyHitRate = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "revenge",
		Subsystem: "dragonfly",
		Name:      "hit_rate",
		Help:      "Dragonfly keyspace hit rate (0-1).",
	})
	dragonflyEvictedKeys = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "revenge",
		Subsystem: "dragonfly",
		Name:      "evicted_keys_total",
		Help:      "Total number of keys evicted by Dragonfly.",
	})
)

// CollectorParams contains dependencies for the periodic metrics collector.
type CollectorParams struct {
	fx.In

	Pool          *pgxpool.Pool
	RueidisClient rueidis.Client `optional:"true"`
	Logger        *slog.Logger
	Lifecycle     fx.Lifecycle
}

// StartCollector registers a periodic metrics collector that runs every 15 seconds.
// It collects pgxpool stats, River queue depth, and Dragonfly stats.
func StartCollector(p CollectorParams) {
	logger := p.Logger.With("component", "metrics-collector")

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go runCollector(p.Pool, p.RueidisClient, logger)
			return nil
		},
	})
}

func runCollector(pool *pgxpool.Pool, redis rueidis.Client, logger *slog.Logger) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		collectPoolStats(pool)
		collectRiverStats(pool, logger)
		if redis != nil {
			collectDragonflyStats(redis, logger)
		}
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

// collectDragonflyStats runs INFO commands against Dragonfly/Redis and exports
// memory, client, keyspace, and command stats as Prometheus gauges.
func collectDragonflyStats(client rueidis.Client, logger *slog.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Fetch all INFO sections in one call
	cmd := client.B().Info().Build()
	resp := client.Do(ctx, cmd)
	if err := resp.Error(); err != nil {
		logger.Debug("failed to collect dragonfly INFO", slog.Any("error", err))
		return
	}

	info, err := resp.ToString()
	if err != nil {
		logger.Debug("failed to read dragonfly INFO response", slog.Any("error", err))
		return
	}

	parsed := parseRedisInfo(info)

	// Memory
	if v, ok := parsed["used_memory"]; ok {
		dragonflyUsedMemory.Set(parseFloat(v))
	}
	if v, ok := parsed["maxmemory"]; ok {
		dragonflyMaxMemory.Set(parseFloat(v))
	}

	// Clients
	if v, ok := parsed["connected_clients"]; ok {
		dragonflyConnectedClients.Set(parseFloat(v))
	}

	// Stats
	if v, ok := parsed["total_commands_processed"]; ok {
		dragonflyCommandsProcessed.Set(parseFloat(v))
	}
	if v, ok := parsed["evicted_keys"]; ok {
		dragonflyEvictedKeys.Set(parseFloat(v))
	}

	// Hit rate: keyspace_hits / (keyspace_hits + keyspace_misses)
	hits := parseFloat(parsed["keyspace_hits"])
	misses := parseFloat(parsed["keyspace_misses"])
	if total := hits + misses; total > 0 {
		dragonflyHitRate.Set(hits / total)
	}

	// Keyspace: parse db0 line like "db0:keys=123,expires=45,avg_ttl=6789"
	if db0, ok := parsed["db0"]; ok {
		for _, part := range strings.Split(db0, ",") {
			kv := strings.SplitN(part, "=", 2)
			if len(kv) == 2 && kv[0] == "keys" {
				dragonflyTotalKeys.Set(parseFloat(kv[1]))
			}
		}
	}
}

// parseRedisInfo parses the output of Redis/Dragonfly INFO into a map.
func parseRedisInfo(info string) map[string]string {
	result := make(map[string]string)
	for _, line := range strings.Split(info, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if k, v, ok := strings.Cut(line, ":"); ok {
			result[k] = v
		}
	}
	return result
}

// parseFloat converts a string to float64, returning 0 on failure.
func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return v
}

package database

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Pool connection metrics
	poolAcquireCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "db_pool_acquire_total",
		Help: "Total number of successful connection acquisitions from the pool",
	})

	poolAcquireDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "db_pool_acquire_duration_seconds",
		Help:    "Duration of connection acquisitions from the pool",
		Buckets: prometheus.DefBuckets,
	})

	poolAcquiredConns = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "db_pool_acquired_conns",
		Help: "Number of currently acquired connections in the pool",
	})

	poolCanceledAcquireCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "db_pool_canceled_acquire_total",
		Help: "Total number of canceled connection acquisitions",
	})

	poolConstructingConns = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "db_pool_constructing_conns",
		Help: "Number of connections being constructed",
	})

	poolEmptyAcquireCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "db_pool_empty_acquire_total",
		Help: "Total number of successful acquires from an empty pool",
	})

	poolIdleConns = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "db_pool_idle_conns",
		Help: "Number of idle connections in the pool",
	})

	poolMaxConns = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "db_pool_max_conns",
		Help: "Maximum number of connections allowed in the pool",
	})

	poolTotalConns = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "db_pool_total_conns",
		Help: "Total number of connections in the pool",
	})

	poolNewConnsCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "db_pool_new_conns_total",
		Help: "Total number of new connections created",
	})

	poolMaxLifetimeDestroyCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "db_pool_max_lifetime_destroy_total",
		Help: "Total number of connections destroyed due to max lifetime",
	})

	poolMaxIdleDestroyCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "db_pool_max_idle_destroy_total",
		Help: "Total number of connections destroyed due to max idle time",
	})
)

// RecordPoolMetrics updates Prometheus metrics from pool statistics.
func RecordPoolMetrics(pool *pgxpool.Pool) {
	stat := pool.Stat()

	// Counter metrics (cumulative)
	poolAcquireCount.Add(float64(stat.AcquireCount()))
	poolCanceledAcquireCount.Add(float64(stat.CanceledAcquireCount()))
	poolEmptyAcquireCount.Add(float64(stat.EmptyAcquireCount()))
	poolNewConnsCount.Add(float64(stat.NewConnsCount()))
	poolMaxLifetimeDestroyCount.Add(float64(stat.MaxLifetimeDestroyCount()))
	poolMaxIdleDestroyCount.Add(float64(stat.MaxIdleDestroyCount()))

	// Gauge metrics (current values)
	poolAcquiredConns.Set(float64(stat.AcquiredConns()))
	poolConstructingConns.Set(float64(stat.ConstructingConns()))
	poolIdleConns.Set(float64(stat.IdleConns()))
	poolMaxConns.Set(float64(stat.MaxConns()))
	poolTotalConns.Set(float64(stat.TotalConns()))

	// Histogram for acquire duration
	poolAcquireDuration.Observe(stat.AcquireDuration().Seconds())
}

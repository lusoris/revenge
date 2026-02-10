// Package observability provides metrics and profiling for production monitoring.
package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// HTTP request metrics
var (
	// HTTPRequestsTotal counts total HTTP requests by method, path, and status.
	HTTPRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "revenge",
		Subsystem: "http",
		Name:      "requests_total",
		Help:      "Total number of HTTP requests by method, path, and status code.",
	}, []string{"method", "path", "status"})

	// HTTPRequestDuration measures HTTP request latency in seconds.
	HTTPRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "revenge",
		Subsystem: "http",
		Name:      "request_duration_seconds",
		Help:      "HTTP request latency in seconds.",
		Buckets:   []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
	}, []string{"method", "path"})

	// HTTPRequestsInFlight tracks current in-flight requests.
	HTTPRequestsInFlight = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "revenge",
		Subsystem: "http",
		Name:      "requests_in_flight",
		Help:      "Number of HTTP requests currently being processed.",
	})
)

// Session metrics
var (
	// ActiveSessions tracks the number of active user sessions.
	ActiveSessions = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "revenge",
		Subsystem: "sessions",
		Name:      "active_total",
		Help:      "Number of active user sessions.",
	})
)

// Cache metrics
var (
	// CacheHitsTotal counts total cache hits by cache name.
	CacheHitsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "revenge",
		Subsystem: "cache",
		Name:      "hits_total",
		Help:      "Total number of cache hits.",
	}, []string{"cache", "layer"})

	// CacheMissesTotal counts total cache misses by cache name.
	CacheMissesTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "revenge",
		Subsystem: "cache",
		Name:      "misses_total",
		Help:      "Total number of cache misses.",
	}, []string{"cache", "layer"})

	// CacheOperationDuration measures cache operation latency.
	CacheOperationDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "revenge",
		Subsystem: "cache",
		Name:      "operation_duration_seconds",
		Help:      "Cache operation latency in seconds.",
		Buckets:   []float64{.0001, .0005, .001, .005, .01, .025, .05, .1},
	}, []string{"cache", "operation"})

	// CacheSize tracks the number of items in cache (L1 only).
	CacheSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "revenge",
		Subsystem: "cache",
		Name:      "size",
		Help:      "Number of items in cache.",
	}, []string{"cache"})
)

// Database query metrics (additional to pool metrics in database package)
var (
	// DBQueryDuration measures database query latency.
	DBQueryDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "revenge",
		Subsystem: "db",
		Name:      "query_duration_seconds",
		Help:      "Database query latency in seconds.",
		Buckets:   []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5},
	}, []string{"operation"})

	// DBQueryErrorsTotal counts database query errors.
	DBQueryErrorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "revenge",
		Subsystem: "db",
		Name:      "query_errors_total",
		Help:      "Total number of database query errors.",
	}, []string{"operation"})
)

// Job queue metrics (River)
var (
	// JobsEnqueuedTotal counts total jobs enqueued by job type.
	JobsEnqueuedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "revenge",
		Subsystem: "jobs",
		Name:      "enqueued_total",
		Help:      "Total number of jobs enqueued.",
	}, []string{"job_type"})

	// JobsCompletedTotal counts total jobs completed by job type and status.
	JobsCompletedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "revenge",
		Subsystem: "jobs",
		Name:      "completed_total",
		Help:      "Total number of jobs completed.",
	}, []string{"job_type", "status"})

	// JobDuration measures job execution time.
	JobDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "revenge",
		Subsystem: "jobs",
		Name:      "duration_seconds",
		Help:      "Job execution duration in seconds.",
		Buckets:   []float64{.1, .5, 1, 5, 10, 30, 60, 120, 300, 600},
	}, []string{"job_type"})

	// JobsQueueSize is DEPRECATED — use the periodically-collected riverQueueSize
	// in collector.go instead, which queries actual River job states from the DB.
	// Kept for backward-compatible metric exposition but no longer populated.
	JobsQueueSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "revenge",
		Subsystem: "jobs",
		Name:      "queue_size",
		Help:      "DEPRECATED: Use revenge_river_queue_size instead. Number of jobs in queue by state.",
	}, []string{"state"})
)

// Library scanner metrics
var (
	// LibraryScanDuration measures library scan duration.
	LibraryScanDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "revenge",
		Subsystem: "library",
		Name:      "scan_duration_seconds",
		Help:      "Library scan duration in seconds.",
		Buckets:   []float64{1, 5, 10, 30, 60, 120, 300, 600, 1800, 3600},
	}, []string{"library_id"})

	// LibraryFilesScanned counts files scanned per library.
	LibraryFilesScanned = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "revenge",
		Subsystem: "library",
		Name:      "files_scanned_total",
		Help:      "Total number of files scanned.",
	}, []string{"library_id"})

	// LibraryScanErrorsTotal counts scan errors.
	LibraryScanErrorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "revenge",
		Subsystem: "library",
		Name:      "scan_errors_total",
		Help:      "Total number of scan errors.",
	}, []string{"library_id", "error_type"})
)

// Search metrics
var (
	// SearchQueriesTotal counts total search queries.
	SearchQueriesTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "revenge",
		Subsystem: "search",
		Name:      "queries_total",
		Help:      "Total number of search queries.",
	}, []string{"type"})

	// SearchQueryDuration measures search query latency.
	SearchQueryDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "revenge",
		Subsystem: "search",
		Name:      "query_duration_seconds",
		Help:      "Search query latency in seconds.",
		Buckets:   []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
	}, []string{"type"})
)

// Auth metrics
var (
	// AuthAttemptsTotal counts authentication attempts.
	AuthAttemptsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "revenge",
		Subsystem: "auth",
		Name:      "attempts_total",
		Help:      "Total number of authentication attempts.",
	}, []string{"method", "status"})

	// RateLimitHitsTotal counts rate limit hits.
	RateLimitHitsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "revenge",
		Subsystem: "ratelimit",
		Name:      "hits_total",
		Help:      "Total number of rate limit hits.",
	}, []string{"limiter", "action"})
)

// RecordCacheHit records a cache hit.
func RecordCacheHit(cacheName, layer string) {
	CacheHitsTotal.WithLabelValues(cacheName, layer).Inc()
}

// RecordCacheMiss records a cache miss.
func RecordCacheMiss(cacheName, layer string) {
	CacheMissesTotal.WithLabelValues(cacheName, layer).Inc()
}

// RecordJobEnqueued records a job being enqueued.
func RecordJobEnqueued(jobType string) {
	JobsEnqueuedTotal.WithLabelValues(jobType).Inc()
}

// RecordJobCompleted records a job completion.
func RecordJobCompleted(jobType, status string) {
	JobsCompletedTotal.WithLabelValues(jobType, status).Inc()
}

// RecordAuthAttempt records an authentication attempt.
func RecordAuthAttempt(method, status string) {
	AuthAttemptsTotal.WithLabelValues(method, status).Inc()
}

// RecordRateLimitHit records a rate limit hit.
func RecordRateLimitHit(limiter, action string) {
	RateLimitHitsTotal.WithLabelValues(limiter, action).Inc()
}

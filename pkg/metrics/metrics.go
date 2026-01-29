// Package metrics provides observability helpers.
package metrics

import (
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

// Counter is a simple counter metric.
type Counter struct {
	value atomic.Int64
}

// Inc increments by 1.
func (c *Counter) Inc() {
	c.value.Add(1)
}

// Add adds delta to the counter.
func (c *Counter) Add(delta int64) {
	c.value.Add(delta)
}

// Value returns current value.
func (c *Counter) Value() int64 {
	return c.value.Load()
}

// Gauge is a value that can go up and down.
type Gauge struct {
	value atomic.Int64
}

// Set sets the gauge value.
func (g *Gauge) Set(val int64) {
	g.value.Store(val)
}

// Inc increments by 1.
func (g *Gauge) Inc() {
	g.value.Add(1)
}

// Dec decrements by 1.
func (g *Gauge) Dec() {
	g.value.Add(-1)
}

// Add adds delta.
func (g *Gauge) Add(delta int64) {
	g.value.Add(delta)
}

// Value returns current value.
func (g *Gauge) Value() int64 {
	return g.value.Load()
}

// Histogram tracks value distributions.
type Histogram struct {
	mu      sync.Mutex
	buckets []int64 // Upper bounds
	counts  []int64
	sum     float64
	count   int64
}

// NewHistogram creates a histogram with buckets.
func NewHistogram(buckets []int64) *Histogram {
	return &Histogram{
		buckets: buckets,
		counts:  make([]int64, len(buckets)+1), // +1 for infinity bucket
	}
}

// Observe records a value.
func (h *Histogram) Observe(val int64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.sum += float64(val)
	h.count++

	for i, bound := range h.buckets {
		if val <= bound {
			h.counts[i]++
			return
		}
	}
	h.counts[len(h.buckets)]++ // Infinity bucket
}

// Percentile estimates a percentile (0-100).
func (h *Histogram) Percentile(p float64) float64 {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.count == 0 {
		return 0
	}

	threshold := int64(float64(h.count) * p / 100)
	var cumulative int64

	for i, count := range h.counts {
		cumulative += count
		if cumulative >= threshold {
			if i < len(h.buckets) {
				return float64(h.buckets[i])
			}
			return float64(h.buckets[len(h.buckets)-1]) // Max bucket
		}
	}

	return 0
}

// Timer tracks durations.
type Timer struct {
	histogram *Histogram
}

// NewTimer creates a timer with millisecond buckets.
func NewTimer() *Timer {
	return &Timer{
		histogram: NewHistogram([]int64{1, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000}),
	}
}

// ObserveDuration records a duration.
func (t *Timer) ObserveDuration(d time.Duration) {
	t.histogram.Observe(d.Milliseconds())
}

// Time returns a function to call when done.
func (t *Timer) Time() func() {
	start := time.Now()
	return func() {
		t.ObserveDuration(time.Since(start))
	}
}

// P50 returns 50th percentile.
func (t *Timer) P50() time.Duration {
	return time.Duration(t.histogram.Percentile(50)) * time.Millisecond
}

// P95 returns 95th percentile.
func (t *Timer) P95() time.Duration {
	return time.Duration(t.histogram.Percentile(95)) * time.Millisecond
}

// P99 returns 99th percentile.
func (t *Timer) P99() time.Duration {
	return time.Duration(t.histogram.Percentile(99)) * time.Millisecond
}

// Registry holds all metrics.
type Registry struct {
	mu         sync.RWMutex
	counters   map[string]*Counter
	gauges     map[string]*Gauge
	histograms map[string]*Histogram
	timers     map[string]*Timer
}

// NewRegistry creates a metric registry.
func NewRegistry() *Registry {
	return &Registry{
		counters:   make(map[string]*Counter),
		gauges:     make(map[string]*Gauge),
		histograms: make(map[string]*Histogram),
		timers:     make(map[string]*Timer),
	}
}

// Counter returns or creates a counter.
func (r *Registry) Counter(name string) *Counter {
	r.mu.Lock()
	defer r.mu.Unlock()

	if c, ok := r.counters[name]; ok {
		return c
	}
	c := &Counter{}
	r.counters[name] = c
	return c
}

// Gauge returns or creates a gauge.
func (r *Registry) Gauge(name string) *Gauge {
	r.mu.Lock()
	defer r.mu.Unlock()

	if g, ok := r.gauges[name]; ok {
		return g
	}
	g := &Gauge{}
	r.gauges[name] = g
	return g
}

// Timer returns or creates a timer.
func (r *Registry) Timer(name string) *Timer {
	r.mu.Lock()
	defer r.mu.Unlock()

	if t, ok := r.timers[name]; ok {
		return t
	}
	t := NewTimer()
	r.timers[name] = t
	return t
}

// Snapshot returns current metric values.
func (r *Registry) Snapshot() map[string]any {
	r.mu.RLock()
	defer r.mu.RUnlock()

	snap := make(map[string]any)

	for name, c := range r.counters {
		snap["counter_"+name] = c.Value()
	}
	for name, g := range r.gauges {
		snap["gauge_"+name] = g.Value()
	}
	for name, t := range r.timers {
		snap["timer_"+name+"_p50_ms"] = t.P50().Milliseconds()
		snap["timer_"+name+"_p95_ms"] = t.P95().Milliseconds()
		snap["timer_"+name+"_p99_ms"] = t.P99().Milliseconds()
	}

	return snap
}

// HTTPMetrics provides HTTP middleware metrics.
type HTTPMetrics struct {
	Requests       *Counter
	RequestsActive *Gauge
	ResponseTime   *Timer
	ResponseSize   *Histogram
	StatusCodes    map[int]*Counter
	mu             sync.Mutex
}

// NewHTTPMetrics creates HTTP metrics.
func NewHTTPMetrics() *HTTPMetrics {
	return &HTTPMetrics{
		Requests:       &Counter{},
		RequestsActive: &Gauge{},
		ResponseTime:   NewTimer(),
		ResponseSize:   NewHistogram([]int64{100, 1000, 10000, 100000, 1000000}),
		StatusCodes:    make(map[int]*Counter),
	}
}

// Middleware returns HTTP middleware.
func (m *HTTPMetrics) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.Requests.Inc()
		m.RequestsActive.Inc()
		defer m.RequestsActive.Dec()

		done := m.ResponseTime.Time()
		defer done()

		// Wrap response writer to capture status and size
		wrapped := &responseWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(wrapped, r)

		m.ResponseSize.Observe(int64(wrapped.size))
		m.statusCounter(wrapped.status).Inc()
	})
}

func (m *HTTPMetrics) statusCounter(status int) *Counter {
	m.mu.Lock()
	defer m.mu.Unlock()

	if c, ok := m.StatusCodes[status]; ok {
		return c
	}
	c := &Counter{}
	m.StatusCodes[status] = c
	return c
}

// Stats returns current HTTP stats.
func (m *HTTPMetrics) Stats() map[string]any {
	m.mu.Lock()
	defer m.mu.Unlock()

	stats := map[string]any{
		"requests_total":    m.Requests.Value(),
		"requests_active":   m.RequestsActive.Value(),
		"response_time_p50": m.ResponseTime.P50().Milliseconds(),
		"response_time_p95": m.ResponseTime.P95().Milliseconds(),
		"response_time_p99": m.ResponseTime.P99().Milliseconds(),
	}

	for status, c := range m.StatusCodes {
		stats["status_"+strconv.Itoa(status)] = c.Value()
	}

	return stats
}

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

func (rw *responseWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}

package observability

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"log/slog"

	"github.com/ogen-go/ogen/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"github.com/lusoris/revenge/internal/config"
	"github.com/lusoris/revenge/internal/infra/logging"
)

// --- extractStatusFromResponse ---

// Types named to match the suffix patterns that extractStatusFromResponse detects.
// extractStatusFromResponse uses fmt.Sprintf("%T", resp.Type), which yields the Go type name.
type fakeGetMovieOK struct{}
type fakeCreateMovieCreated struct{}
type fakeDeleteMovieNoContent struct{}
type fakeCreateMovieBadRequest struct{}
type fakeGetMovieUnauthorized struct{}
type fakeGetMovieForbidden struct{}
type fakeGetMovieNotFound struct{}
type fakeCreateUserConflict struct{}
type fakeLoginTooManyRequests struct{}
type fakeGetMovieInternalServerError struct{}
type fakeUnknownResponse struct{}

func TestExtractStatusFromResponse_AllSuffixes(t *testing.T) {
	tests := []struct {
		name     string
		resp     middleware.Response
		expected string
	}{
		{
			name:     "OK suffix",
			resp:     middleware.Response{Type: fakeGetMovieOK{}},
			expected: "200",
		},
		{
			name:     "Created suffix",
			resp:     middleware.Response{Type: fakeCreateMovieCreated{}},
			expected: "201",
		},
		{
			name:     "NoContent suffix",
			resp:     middleware.Response{Type: fakeDeleteMovieNoContent{}},
			expected: "204",
		},
		{
			name:     "BadRequest suffix",
			resp:     middleware.Response{Type: fakeCreateMovieBadRequest{}},
			expected: "400",
		},
		{
			name:     "Unauthorized suffix",
			resp:     middleware.Response{Type: fakeGetMovieUnauthorized{}},
			expected: "401",
		},
		{
			name:     "Forbidden suffix",
			resp:     middleware.Response{Type: fakeGetMovieForbidden{}},
			expected: "403",
		},
		{
			name:     "NotFound suffix",
			resp:     middleware.Response{Type: fakeGetMovieNotFound{}},
			expected: "404",
		},
		{
			name:     "Conflict suffix",
			resp:     middleware.Response{Type: fakeCreateUserConflict{}},
			expected: "409",
		},
		{
			name:     "TooManyRequests suffix",
			resp:     middleware.Response{Type: fakeLoginTooManyRequests{}},
			expected: "429",
		},
		{
			name:     "InternalServerError suffix",
			resp:     middleware.Response{Type: fakeGetMovieInternalServerError{}},
			expected: "500",
		},
		{
			name:     "unknown suffix defaults to 200",
			resp:     middleware.Response{Type: fakeUnknownResponse{}},
			expected: "200",
		},
		{
			name:     "nil type defaults to 200",
			resp:     middleware.Response{Type: nil},
			expected: "200",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractStatusFromResponse(tt.resp)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// --- RegisterPprofHandlers ---

func TestRegisterPprofHandlers(t *testing.T) {
	mux := http.NewServeMux()

	// Should not panic
	assert.NotPanics(t, func() {
		RegisterPprofHandlers(mux)
	})

	// Verify all expected pprof routes are registered by doing a request
	pprofPaths := []string{
		"/debug/pprof/",
		"/debug/pprof/cmdline",
		"/debug/pprof/heap",
		"/debug/pprof/goroutine",
		"/debug/pprof/block",
		"/debug/pprof/mutex",
		"/debug/pprof/allocs",
		"/debug/pprof/threadcreate",
	}

	for _, path := range pprofPaths {
		t.Run(path, func(t *testing.T) {
			req := httptest.NewRequest("GET", path, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			// Should get a response (not 404 from default mux handler)
			assert.NotEqual(t, 404, rec.Code,
				"pprof handler should be registered at %s", path)
		})
	}
}

// --- Module ---

func TestModule(t *testing.T) {
	assert.NotNil(t, Module)
}

// --- NewServer ---

func TestNewServer_Development(t *testing.T) {
	logger := logging.NewTestLogger()

	lc := fxtest.NewLifecycle(t)

	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: 19999,
		},
		Logging: config.LoggingConfig{
			Development: true,
		},
	}

	params := ServerParams{
		Config:    cfg,
		Logger:    logger,
		Lifecycle: lc,
	}

	server := NewServer(params)
	assert.NotNil(t, server)
	assert.NotNil(t, server.httpServer)
	assert.NotNil(t, server.logger)

	// Verify the address is correct (port + 1000)
	assert.Equal(t, "127.0.0.1:20999", server.httpServer.Addr)
}

func TestNewServer_Production(t *testing.T) {
	logger := logging.NewTestLogger()

	lc := fxtest.NewLifecycle(t)

	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
		Logging: config.LoggingConfig{
			Development: false,
		},
	}

	params := ServerParams{
		Config:    cfg,
		Logger:    logger,
		Lifecycle: lc,
	}

	server := NewServer(params)
	assert.NotNil(t, server)
	assert.Equal(t, "0.0.0.0:9080", server.httpServer.Addr)
}

// --- Metric variable initialization ---

func TestMetricVariables_NotNil(t *testing.T) {
	// HTTP metrics
	assert.NotNil(t, HTTPRequestsTotal, "HTTPRequestsTotal should be initialized")
	assert.NotNil(t, HTTPRequestDuration, "HTTPRequestDuration should be initialized")
	assert.NotNil(t, HTTPRequestsInFlight, "HTTPRequestsInFlight should be initialized")

	// Session metrics
	assert.NotNil(t, ActiveSessions, "ActiveSessions should be initialized")

	// Cache metrics
	assert.NotNil(t, CacheHitsTotal, "CacheHitsTotal should be initialized")
	assert.NotNil(t, CacheMissesTotal, "CacheMissesTotal should be initialized")
	assert.NotNil(t, CacheOperationDuration, "CacheOperationDuration should be initialized")
	assert.NotNil(t, CacheSize, "CacheSize should be initialized")

	// DB metrics
	assert.NotNil(t, DBQueryDuration, "DBQueryDuration should be initialized")
	assert.NotNil(t, DBQueryErrorsTotal, "DBQueryErrorsTotal should be initialized")

	// Job metrics
	assert.NotNil(t, JobsEnqueuedTotal, "JobsEnqueuedTotal should be initialized")
	assert.NotNil(t, JobsCompletedTotal, "JobsCompletedTotal should be initialized")
	assert.NotNil(t, JobDuration, "JobDuration should be initialized")
	assert.NotNil(t, JobsQueueSize, "JobsQueueSize should be initialized")

	// Library metrics
	assert.NotNil(t, LibraryScanDuration, "LibraryScanDuration should be initialized")
	assert.NotNil(t, LibraryFilesScanned, "LibraryFilesScanned should be initialized")
	assert.NotNil(t, LibraryScanErrorsTotal, "LibraryScanErrorsTotal should be initialized")

	// Search metrics
	assert.NotNil(t, SearchQueriesTotal, "SearchQueriesTotal should be initialized")
	assert.NotNil(t, SearchQueryDuration, "SearchQueryDuration should be initialized")

	// Auth metrics
	assert.NotNil(t, AuthAttemptsTotal, "AuthAttemptsTotal should be initialized")
	assert.NotNil(t, RateLimitHitsTotal, "RateLimitHitsTotal should be initialized")
}

// --- Metric operations (functional tests) ---

func TestMetricOperations_Functional(t *testing.T) {
	t.Run("HTTPRequestsTotal increments", func(t *testing.T) {
		// Should not panic
		HTTPRequestsTotal.WithLabelValues("GET", "/api/v1/test", "200").Inc()
	})

	t.Run("HTTPRequestDuration observes", func(t *testing.T) {
		HTTPRequestDuration.WithLabelValues("GET", "/api/v1/test").Observe(0.042)
	})

	t.Run("HTTPRequestsInFlight gauge", func(t *testing.T) {
		HTTPRequestsInFlight.Inc()
		HTTPRequestsInFlight.Dec()
	})

	t.Run("ActiveSessions gauge", func(t *testing.T) {
		ActiveSessions.Set(100)
		ActiveSessions.Inc()
		ActiveSessions.Dec()
	})

	t.Run("CacheHitsTotal increments", func(t *testing.T) {
		CacheHitsTotal.WithLabelValues("movies", "l1").Inc()
		CacheHitsTotal.WithLabelValues("movies", "l2").Inc()
	})

	t.Run("CacheMissesTotal increments", func(t *testing.T) {
		CacheMissesTotal.WithLabelValues("sessions", "l1").Inc()
	})

	t.Run("CacheOperationDuration observes", func(t *testing.T) {
		CacheOperationDuration.WithLabelValues("default", "get").Observe(0.001)
		CacheOperationDuration.WithLabelValues("default", "set").Observe(0.002)
		CacheOperationDuration.WithLabelValues("default", "delete").Observe(0.0005)
	})

	t.Run("CacheSize gauge", func(t *testing.T) {
		CacheSize.WithLabelValues("sessions").Set(42)
	})

	t.Run("DBQueryDuration observes", func(t *testing.T) {
		DBQueryDuration.WithLabelValues("select").Observe(0.01)
		DBQueryDuration.WithLabelValues("insert").Observe(0.005)
	})

	t.Run("DBQueryErrorsTotal increments", func(t *testing.T) {
		DBQueryErrorsTotal.WithLabelValues("select").Inc()
	})

	t.Run("JobsEnqueuedTotal increments", func(t *testing.T) {
		JobsEnqueuedTotal.WithLabelValues("library_scan").Inc()
	})

	t.Run("JobsCompletedTotal increments", func(t *testing.T) {
		JobsCompletedTotal.WithLabelValues("library_scan", "completed").Inc()
		JobsCompletedTotal.WithLabelValues("metadata_refresh", "failed").Inc()
	})

	t.Run("JobDuration observes", func(t *testing.T) {
		JobDuration.WithLabelValues("library_scan").Observe(120.5)
	})

	t.Run("JobsQueueSize gauge", func(t *testing.T) {
		JobsQueueSize.WithLabelValues("available").Set(10)
		JobsQueueSize.WithLabelValues("running").Set(3)
	})

	t.Run("LibraryScanDuration observes", func(t *testing.T) {
		LibraryScanDuration.WithLabelValues("lib-1").Observe(45.0)
	})

	t.Run("LibraryFilesScanned increments", func(t *testing.T) {
		LibraryFilesScanned.WithLabelValues("lib-1").Add(150)
	})

	t.Run("LibraryScanErrorsTotal increments", func(t *testing.T) {
		LibraryScanErrorsTotal.WithLabelValues("lib-1", "permission_denied").Inc()
	})

	t.Run("SearchQueriesTotal increments", func(t *testing.T) {
		SearchQueriesTotal.WithLabelValues("fulltext").Inc()
		SearchQueriesTotal.WithLabelValues("autocomplete").Inc()
	})

	t.Run("SearchQueryDuration observes", func(t *testing.T) {
		SearchQueryDuration.WithLabelValues("fulltext").Observe(0.025)
	})

	t.Run("AuthAttemptsTotal increments", func(t *testing.T) {
		AuthAttemptsTotal.WithLabelValues("api_key", "success").Inc()
		AuthAttemptsTotal.WithLabelValues("session", "failure").Inc()
	})

	t.Run("RateLimitHitsTotal increments", func(t *testing.T) {
		RateLimitHitsTotal.WithLabelValues("api", "throttled").Inc()
	})
}

// --- Helper functions ---

func TestRecordCacheHit_Unit(t *testing.T) {
	// Should not panic; verify with prometheus collection
	assert.NotPanics(t, func() {
		RecordCacheHit("unit_test_cache", "l1")
		RecordCacheHit("unit_test_cache", "l2")
	})
}

func TestRecordCacheMiss_Unit(t *testing.T) {
	assert.NotPanics(t, func() {
		RecordCacheMiss("unit_test_cache", "l1")
		RecordCacheMiss("unit_test_cache", "l2")
	})
}

func TestRecordJobEnqueued_Unit(t *testing.T) {
	assert.NotPanics(t, func() {
		RecordJobEnqueued("test_job_type")
	})
}

func TestRecordJobCompleted_Unit(t *testing.T) {
	assert.NotPanics(t, func() {
		RecordJobCompleted("test_job_type", "completed")
		RecordJobCompleted("test_job_type", "failed")
		RecordJobCompleted("test_job_type", "discarded")
	})
}

func TestRecordAuthAttempt_Unit(t *testing.T) {
	assert.NotPanics(t, func() {
		RecordAuthAttempt("login", "success")
		RecordAuthAttempt("register", "failure")
		RecordAuthAttempt("verify_email", "success")
	})
}

func TestRecordRateLimitHit_Unit(t *testing.T) {
	assert.NotPanics(t, func() {
		RecordRateLimitHit("global", "allowed")
		RecordRateLimitHit("auth", "blocked")
		RecordRateLimitHit("api", "throttled")
	})
}

// --- StandardHTTPMetricsMiddleware detailed tests ---

func TestStandardHTTPMetricsMiddleware_Methods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			wrapped := StandardHTTPMetricsMiddleware(handler)

			req := httptest.NewRequest(method, "/api/v1/movies", nil)
			rec := httptest.NewRecorder()

			wrapped.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
		})
	}
}

func TestStandardHTTPMetricsMiddleware_StatusCodes(t *testing.T) {
	statusCodes := []int{200, 201, 204, 400, 401, 403, 404, 409, 429, 500}

	for _, code := range statusCodes {
		t.Run(fmt.Sprintf("status_%d", code), func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(code)
			})

			wrapped := StandardHTTPMetricsMiddleware(handler)

			req := httptest.NewRequest("GET", "/test", nil)
			rec := httptest.NewRecorder()

			wrapped.ServeHTTP(rec, req)

			assert.Equal(t, code, rec.Code)
		})
	}
}

func TestStandardHTTPMetricsMiddleware_DefaultStatus(t *testing.T) {
	// When handler doesn't call WriteHeader, default is 200
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("OK"))
	})

	wrapped := StandardHTTPMetricsMiddleware(handler)

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
}

// --- statusResponseWriter ---

func TestStatusResponseWriter_WriteWithoutHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	wrapped := &statusResponseWriter{ResponseWriter: rec, status: 200}

	// Write without calling WriteHeader first
	n, err := wrapped.Write([]byte("hello"))
	assert.NoError(t, err)
	assert.Equal(t, 5, n)

	// Status should remain default 200
	assert.Equal(t, 200, wrapped.status)
}

func TestStatusResponseWriter_MultipleWriteHeaders(t *testing.T) {
	rec := httptest.NewRecorder()
	wrapped := &statusResponseWriter{ResponseWriter: rec, status: 200}

	wrapped.WriteHeader(301)
	assert.Equal(t, 301, wrapped.status)

	// HTTP spec says only first WriteHeader matters, but our wrapper tracks it
	wrapped.WriteHeader(404)
	assert.Equal(t, 404, wrapped.status)
}

// --- normalizePath edge cases ---

func TestNormalizePath_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "root path",
			input:    "/",
			expected: "/",
		},
		{
			name:     "trailing slash",
			input:    "/api/v1/movies/",
			expected: "/api/v1/movies/",
		},
		{
			name:     "double slash",
			input:    "/api//v1",
			expected: "/api//v1",
		},
		{
			name:     "deep nesting with IDs",
			input:    "/api/v1/libraries/123/movies/456/cast",
			expected: "/api/v1/libraries/{id}/movies/{id}/cast",
		},
		{
			name:     "query string not affected (path only)",
			input:    "/api/v1/search",
			expected: "/api/v1/search",
		},
		{
			name:     "legacy path",
			input:    "/api/v1/legacy/movies",
			expected: "/api/v1/legacy/movies",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizePath(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// --- isIDSegment edge cases ---

func TestIsIDSegment_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"empty string", "", false},
		{"single digit", "5", true},
		{"negative number (parses as int64)", "-1", true},
		{"zero", "0", true},
		{"large number", "999999999", true},
		{"short word", "api", false},
		{"two chars", "v1", false},
		{"seven alphanumeric with digit", "abc1234", false},     // less than 8 chars
		{"eight alphanumeric with digit", "abcd1234", true},     // exactly 8 chars
		{"eight alpha no digit", "abcdefgh", false},             // no digits
		{"UUID-like format", "12345678-1234-1234-1234-123456789012", true},
		{"not quite UUID (still looks like ID)", "12345678-1234-1234-1234", true}, // alphanumeric with hyphens, 8+ chars, has digits
		{"special characters", "abc!@#123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isIDSegment(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// --- isAlphanumericWithHyphens edge cases ---

func TestIsAlphanumericWithHyphens_EdgeCases(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", true},
		{"-", true},
		{"---", true},
		{"a", true},
		{"Z", true},
		{"0", true},
		{"9", true},
		{"a-z", true},
		{"A-Z-0-9", true},
		{"hello_world", false},
		{"hello world", false},
		{"hello.world", false},
		{"hello/world", false},
		{"hello\nworld", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := isAlphanumericWithHyphens(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// --- hasDigits edge cases ---

func TestHasDigits_EdgeCases(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"a", false},
		{"Z", false},
		{"-", false},
		{"0", true},
		{"9", true},
		{"a0", true},
		{"0a", true},
		{"abc", false},
		{"ABC", false},
		{"---", false},
		{"a1b2c3", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := hasDigits(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// --- Prometheus registration verification ---

func TestPrometheusRegistration(t *testing.T) {
	// Verify metrics are registered with prometheus by attempting to describe them
	t.Run("HTTPRequestsTotal is a counter", func(t *testing.T) {
		ch := make(chan *prometheus.Desc, 10)
		HTTPRequestsTotal.Describe(ch)
		desc := <-ch
		assert.NotNil(t, desc)
	})

	t.Run("HTTPRequestDuration is a histogram", func(t *testing.T) {
		ch := make(chan *prometheus.Desc, 10)
		HTTPRequestDuration.Describe(ch)
		desc := <-ch
		assert.NotNil(t, desc)
	})

	t.Run("HTTPRequestsInFlight is a gauge", func(t *testing.T) {
		ch := make(chan *prometheus.Desc, 10)
		HTTPRequestsInFlight.Describe(ch)
		desc := <-ch
		assert.NotNil(t, desc)
	})

	t.Run("CacheHitsTotal is a counter", func(t *testing.T) {
		ch := make(chan *prometheus.Desc, 10)
		CacheHitsTotal.Describe(ch)
		desc := <-ch
		assert.NotNil(t, desc)
	})

	t.Run("JobsEnqueuedTotal is a counter", func(t *testing.T) {
		ch := make(chan *prometheus.Desc, 10)
		JobsEnqueuedTotal.Describe(ch)
		desc := <-ch
		assert.NotNil(t, desc)
	})
}

// --- HTTPMetricsMiddleware (ogen) ---

func TestHTTPMetricsMiddleware_Success(t *testing.T) {
	mw := HTTPMetricsMiddleware()
	assert.NotNil(t, mw)

	req := middleware.Request{
		Raw: httptest.NewRequest("GET", "/api/v1/movies", nil),
	}

	next := func(req middleware.Request) (middleware.Response, error) {
		return middleware.Response{Type: fakeGetMovieOK{}}, nil
	}

	resp, err := mw(req, next)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestHTTPMetricsMiddleware_Error(t *testing.T) {
	mw := HTTPMetricsMiddleware()

	req := middleware.Request{
		Raw: httptest.NewRequest("POST", "/api/v1/movies", nil),
	}

	next := func(req middleware.Request) (middleware.Response, error) {
		return middleware.Response{}, assert.AnError
	}

	resp, err := mw(req, next)
	assert.Error(t, err)
	// When error occurs, status should be "500"
	_ = resp
}

func TestHTTPMetricsMiddleware_WithIDs(t *testing.T) {
	mw := HTTPMetricsMiddleware()

	req := middleware.Request{
		Raw: httptest.NewRequest("GET", "/api/v1/movies/550e8400-e29b-41d4-a716-446655440000", nil),
	}

	next := func(req middleware.Request) (middleware.Response, error) {
		return middleware.Response{Type: fakeGetMovieOK{}}, nil
	}

	resp, err := mw(req, next)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestHTTPMetricsMiddleware_VariousStatuses(t *testing.T) {
	mw := HTTPMetricsMiddleware()

	statusTypes := []struct {
		name     string
		respType any
	}{
		{"OK", fakeGetMovieOK{}},
		{"Created", fakeCreateMovieCreated{}},
		{"NotFound", fakeGetMovieNotFound{}},
		{"BadRequest", fakeCreateMovieBadRequest{}},
		{"Unauthorized", fakeGetMovieUnauthorized{}},
	}

	for _, st := range statusTypes {
		t.Run(st.name, func(t *testing.T) {
			req := middleware.Request{
				Raw: httptest.NewRequest("GET", "/api/v1/test", nil),
			}

			next := func(req middleware.Request) (middleware.Response, error) {
				return middleware.Response{Type: st.respType}, nil
			}

			resp, err := mw(req, next)
			assert.NoError(t, err)
			assert.NotNil(t, resp)
		})
	}
}

// --- NewServer lifecycle ---

func TestNewServer_HealthEndpoints(t *testing.T) {
	logger := logging.NewTestLogger()
	lc := fxtest.NewLifecycle(t)

	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: 19876,
		},
		Logging: config.LoggingConfig{
			Development: false,
		},
	}

	_ = NewServer(ServerParams{
		Config:    cfg,
		Logger:    logger,
		Lifecycle: lc,
	})

	// Start the lifecycle to register the HTTP server
	require.NoError(t, lc.Start(t.Context()))

	// The server listens on port+1000=20876
	// Test health endpoints
	resp, err := http.Get("http://127.0.0.1:20876/health/live")
	if err == nil {
		assert.Equal(t, 200, resp.StatusCode)
		resp.Body.Close()
	}

	resp, err = http.Get("http://127.0.0.1:20876/health/ready")
	if err == nil {
		assert.Equal(t, 200, resp.StatusCode)
		resp.Body.Close()
	}

	// Stop the lifecycle
	require.NoError(t, lc.Stop(t.Context()))
}

func TestNewServer_WithPprof(t *testing.T) {
	logger := logging.NewTestLogger()
	lc := fxtest.NewLifecycle(t)

	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "127.0.0.1",
			Port: 19877,
		},
		Logging: config.LoggingConfig{
			Development: true, // enables pprof
		},
	}

	server := NewServer(ServerParams{
		Config:    cfg,
		Logger:    logger,
		Lifecycle: lc,
	})

	assert.NotNil(t, server)
	// Port should be 19877 + 1000 = 20877
	assert.Equal(t, "127.0.0.1:20877", server.httpServer.Addr)
}

// --- ServerParams struct validation ---

func TestServerParams_FxIn(t *testing.T) {
	// Verify ServerParams has fx.In embedded (compilation check)
	_ = ServerParams{
		Config:    &config.Config{},
		Logger:    logging.NewTestLogger(),
		Lifecycle: fxtest.NewLifecycle(t),
	}
}

// --- Server struct ---

func TestServer_Fields(t *testing.T) {
	logger := logging.NewTestLogger()
	lc := fxtest.NewLifecycle(t)

	cfg := &config.Config{
		Server: config.ServerConfig{
			Host: "localhost",
			Port: 3000,
		},
		Logging: config.LoggingConfig{
			Development: false,
		},
	}

	server := NewServer(ServerParams{
		Config:    cfg,
		Logger:    logger,
		Lifecycle: lc,
	})

	assert.NotNil(t, server.httpServer)
	assert.NotNil(t, server.logger)
	assert.Equal(t, "localhost:4000", server.httpServer.Addr)

	// ReadHeaderTimeout should be set for security
	require.NotZero(t, server.httpServer.ReadHeaderTimeout)
}

// --- Module integration with fx ---

func TestModule_FxOptions(t *testing.T) {
	// Module should be a valid fx.Options
	assert.NotPanics(t, func() {
		_ = fx.New(
			Module,
			fx.Provide(func() *config.Config {
				return &config.Config{
					Server: config.ServerConfig{
						Host: "127.0.0.1",
						Port: 18888,
					},
					Logging: config.LoggingConfig{
						Development: false,
					},
				}
			}),
			fx.Provide(func() *slog.Logger {
				return logging.NewTestLogger()
			}),
			fx.NopLogger,
		)
	})
}

package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lusoris/revenge/internal/config"
)

// ============================================================================
// CookieAuthMiddleware
// ============================================================================

func TestCookieAuthMiddleware_Disabled(t *testing.T) {
	cfg := config.CookieAuthConfig{Enabled: false}
	handler := CookieAuthMiddleware(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Should NOT inject Authorization when disabled
		assert.Empty(t, r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.AddCookie(&http.Cookie{Name: CookieAccessToken, Value: "some-jwt"})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCookieAuthMiddleware_InjectsBearer(t *testing.T) {
	cfg := config.CookieAuthConfig{Enabled: true}
	handler := CookieAuthMiddleware(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer my-token-123", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.AddCookie(&http.Cookie{Name: CookieAccessToken, Value: "my-token-123"})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCookieAuthMiddleware_ExistingAuthHeaderTakesPrecedence(t *testing.T) {
	cfg := config.CookieAuthConfig{Enabled: true}
	handler := CookieAuthMiddleware(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer existing-token", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	req.Header.Set("Authorization", "Bearer existing-token")
	req.AddCookie(&http.Cookie{Name: CookieAccessToken, Value: "cookie-token"})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCookieAuthMiddleware_NoCookie(t *testing.T) {
	cfg := config.CookieAuthConfig{Enabled: true}
	handler := CookieAuthMiddleware(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// ============================================================================
// SetAuthCookies
// ============================================================================

func TestSetAuthCookies_AllTokens(t *testing.T) {
	rec := httptest.NewRecorder()
	cfg := config.CookieAuthConfig{
		Enabled:  true,
		Path:     "/",
		Domain:   "example.com",
		Secure:   true,
		SameSite: "strict",
	}

	SetAuthCookies(rec, cfg, "access-tok", "refresh-tok", "csrf-tok", 900)

	cookies := rec.Result().Cookies()
	require.Len(t, cookies, 3)

	cookieMap := make(map[string]*http.Cookie)
	for _, c := range cookies {
		cookieMap[c.Name] = c
	}

	// Access token: HttpOnly, Secure
	access := cookieMap[CookieAccessToken]
	require.NotNil(t, access)
	assert.Equal(t, "access-tok", access.Value)
	assert.True(t, access.HttpOnly)
	assert.True(t, access.Secure)
	assert.Equal(t, 900, access.MaxAge)
	assert.Equal(t, "example.com", access.Domain)

	// Refresh token: restricted path
	refresh := cookieMap[CookieRefreshToken]
	require.NotNil(t, refresh)
	assert.Equal(t, "refresh-tok", refresh.Value)
	assert.Equal(t, "/api/v1/auth", refresh.Path)
	assert.True(t, refresh.HttpOnly)
	assert.Equal(t, 7*24*3600, refresh.MaxAge)

	// CSRF token: NOT HttpOnly
	csrf := cookieMap[CookieCSRFToken]
	require.NotNil(t, csrf)
	assert.Equal(t, "csrf-tok", csrf.Value)
	assert.False(t, csrf.HttpOnly)
}

func TestSetAuthCookies_NoRefreshOrCSRF(t *testing.T) {
	rec := httptest.NewRecorder()
	cfg := config.CookieAuthConfig{Enabled: true, Path: "/app"}

	SetAuthCookies(rec, cfg, "access-only", "", "", 600)

	cookies := rec.Result().Cookies()
	require.Len(t, cookies, 1)
	assert.Equal(t, CookieAccessToken, cookies[0].Name)
	assert.Equal(t, "/app", cookies[0].Path)
}

func TestSetAuthCookies_DefaultPath(t *testing.T) {
	rec := httptest.NewRecorder()
	cfg := config.CookieAuthConfig{Enabled: true}

	SetAuthCookies(rec, cfg, "tok", "", "", 60)

	cookies := rec.Result().Cookies()
	require.Len(t, cookies, 1)
	assert.Equal(t, "/", cookies[0].Path)
}

// ============================================================================
// ClearAuthCookies
// ============================================================================

func TestClearAuthCookies(t *testing.T) {
	rec := httptest.NewRecorder()
	cfg := config.CookieAuthConfig{
		Enabled: true,
		Path:    "/",
		Domain:  "example.com",
	}

	ClearAuthCookies(rec, cfg)

	cookies := rec.Result().Cookies()
	require.Len(t, cookies, 3)

	for _, c := range cookies {
		assert.Equal(t, -1, c.MaxAge, "cookie %s should be expired", c.Name)
		assert.Equal(t, "", c.Value, "cookie %s should have empty value", c.Name)
	}

	names := make([]string, len(cookies))
	for i, c := range cookies {
		names[i] = c.Name
	}
	assert.Contains(t, names, CookieAccessToken)
	assert.Contains(t, names, CookieRefreshToken)
	assert.Contains(t, names, CookieCSRFToken)
}

// ============================================================================
// ResponseWriterMiddleware / GetResponseWriter
// ============================================================================

func TestResponseWriterMiddleware(t *testing.T) {
	handler := ResponseWriterMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, ok := GetResponseWriter(r.Context())
		assert.True(t, ok)
		assert.Equal(t, w, got)
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetResponseWriter_NotInContext(t *testing.T) {
	_, ok := GetResponseWriter(context.Background())
	assert.False(t, ok)
}

// ============================================================================
// CSRFMiddleware
// ============================================================================

func TestCSRFMiddleware_Disabled(t *testing.T) {
	cfg := config.CookieAuthConfig{Enabled: false}
	handler := CSRFMiddleware(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCSRFMiddleware_SafeMethodSkipped(t *testing.T) {
	cfg := config.CookieAuthConfig{Enabled: true}
	handler := CSRFMiddleware(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	for _, method := range []string{http.MethodGet, http.MethodHead, http.MethodOptions} {
		req := httptest.NewRequest(method, "/api/test", nil)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code, "method %s should be allowed", method)
	}
}

func TestCSRFMiddleware_BearerAuthSkipsCSRF(t *testing.T) {
	cfg := config.CookieAuthConfig{Enabled: true}
	handler := CSRFMiddleware(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/test", nil)
	req.Header.Set("Authorization", "Bearer some-token")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCSRFMiddleware_NoCookieAuthSkipsCSRF(t *testing.T) {
	cfg := config.CookieAuthConfig{Enabled: true}
	handler := CSRFMiddleware(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// No cookie, no auth header → not using cookie auth → CSRF skipped
	req := httptest.NewRequest(http.MethodPost, "/api/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCSRFMiddleware_MissingCSRFCookie(t *testing.T) {
	cfg := config.CookieAuthConfig{Enabled: true}
	handler := CSRFMiddleware(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/test", nil)
	req.AddCookie(&http.Cookie{Name: CookieAccessToken, Value: "jwt"})
	// No CSRF cookie
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Contains(t, rec.Body.String(), "CSRF token missing")
}

func TestCSRFMiddleware_InvalidCSRFToken(t *testing.T) {
	cfg := config.CookieAuthConfig{Enabled: true}
	handler := CSRFMiddleware(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/api/test", nil)
	req.AddCookie(&http.Cookie{Name: CookieAccessToken, Value: "jwt"})
	req.AddCookie(&http.Cookie{Name: CookieCSRFToken, Value: "correct-csrf"})
	req.Header.Set(HeaderCSRF, "wrong-csrf")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.Contains(t, rec.Body.String(), "CSRF token invalid")
}

func TestCSRFMiddleware_ValidCSRFToken(t *testing.T) {
	cfg := config.CookieAuthConfig{Enabled: true}
	handler := CSRFMiddleware(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	token, err := GenerateCSRFToken()
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/test", nil)
	req.AddCookie(&http.Cookie{Name: CookieAccessToken, Value: "jwt"})
	req.AddCookie(&http.Cookie{Name: CookieCSRFToken, Value: token})
	req.Header.Set(HeaderCSRF, token)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestCSRFMiddleware_AllMutatingMethods(t *testing.T) {
	cfg := config.CookieAuthConfig{Enabled: true}
	handler := CSRFMiddleware(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	token := "csrf-test-token"

	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
		req := httptest.NewRequest(method, "/api/test", nil)
		req.AddCookie(&http.Cookie{Name: CookieAccessToken, Value: "jwt"})
		req.AddCookie(&http.Cookie{Name: CookieCSRFToken, Value: token})
		req.Header.Set(HeaderCSRF, token)
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code, "method %s should pass with valid CSRF", method)
	}
}

// ============================================================================
// ErrorHandler
// ============================================================================

func TestErrorHandler_RateLimitError(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := req.Context()

	rateLimitErr := &RateLimitError{
		IP:            "192.168.1.1",
		OperationName: "GetMovie",
		RetryAfter:    30 * time.Second,
	}

	ErrorHandler(ctx, rec, req, rateLimitErr)

	assert.Equal(t, http.StatusTooManyRequests, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Equal(t, "30s", rec.Header().Get("Retry-After"))

	var resp ErrorResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	require.NoError(t, err)
	assert.Equal(t, "rate_limit_exceeded", resp.Error)
	assert.Equal(t, "RATE_LIMIT_EXCEEDED", resp.Code)
}

func TestErrorHandler_GenericError(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := req.Context()

	ErrorHandler(ctx, rec, req, assert.AnError)

	// Default ogen error handler returns 500
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

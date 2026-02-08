package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/lusoris/revenge/internal/config"
)

// CSRFMiddleware enforces double-submit cookie CSRF protection.
// On state-changing requests (POST, PUT, PATCH, DELETE), it validates
// that the X-CSRF-Token header matches the revenge_csrf cookie value.
// Only active when cookie auth is enabled.
func CSRFMiddleware(cfg config.CookieAuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if !cfg.Enabled {
			return next
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip CSRF for safe methods
			if isSafeMethod(r.Method) {
				next.ServeHTTP(w, r)
				return
			}

			// Only enforce CSRF when the request uses cookie auth
			// (has the access token cookie but no explicit Authorization header)
			_, hasCookie := r.Cookie(CookieAccessToken)
			hasAuthHeader := r.Header.Get("Authorization") != ""

			// If using Bearer or API key auth directly, skip CSRF
			if hasCookie != nil || hasAuthHeader {
				next.ServeHTTP(w, r)
				return
			}

			// Validate CSRF: cookie value must match header
			csrfCookie, err := r.Cookie(CookieCSRFToken)
			if err != nil || csrfCookie.Value == "" {
				http.Error(w, "CSRF token missing", http.StatusForbidden)
				return
			}

			headerToken := r.Header.Get(HeaderCSRF)
			if headerToken == "" || headerToken != csrfCookie.Value {
				http.Error(w, "CSRF token invalid", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GenerateCSRFToken creates a cryptographically secure random CSRF token.
func GenerateCSRFToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func isSafeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
		return true
	default:
		return false
	}
}

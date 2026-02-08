package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/lusoris/revenge/internal/config"
)

// CORSMiddleware returns an HTTP middleware that handles Cross-Origin Resource Sharing.
//
// Behavior:
//   - Sets Access-Control-Allow-Origin to the request origin if it matches the allowed list
//   - When AllowedOrigins contains "*", any origin is accepted (reflected, not literal "*")
//   - Handles OPTIONS preflight requests with 204 No Content
//   - Adds Vary: Origin header for correct cache behavior behind CDNs/proxies
func CORSMiddleware(cfg config.CORSConfig) func(http.Handler) http.Handler {
	originsSet := make(map[string]bool, len(cfg.AllowedOrigins))
	allowAll := false
	for _, o := range cfg.AllowedOrigins {
		if o == "*" {
			allowAll = true
		}
		originsSet[strings.TrimRight(o, "/")] = true
	}

	maxAge := fmt.Sprintf("%d", int(cfg.MaxAge.Seconds()))
	if cfg.MaxAge == 0 {
		maxAge = "43200" // 12h default
	}

	const (
		allowMethods = "GET, POST, PUT, PATCH, DELETE, OPTIONS"
		allowHeaders = "Authorization, Content-Type, X-API-Key, X-Request-ID, Accept, Accept-Language"
		exposeHeaders = "X-Request-ID, Retry-After, Content-Disposition"
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// No Origin header = same-origin or non-browser request, skip CORS
			if origin == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Check if origin is allowed
			if !allowAll && !originsSet[strings.TrimRight(origin, "/")] {
				next.ServeHTTP(w, r)
				return
			}

			// Always set Vary: Origin so caches key on it
			w.Header().Add("Vary", "Origin")

			// Reflect the actual origin (never literal "*" — browsers reject it with credentials)
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", allowMethods)
			w.Header().Set("Access-Control-Allow-Headers", allowHeaders)
			w.Header().Set("Access-Control-Expose-Headers", exposeHeaders)
			w.Header().Set("Access-Control-Max-Age", maxAge)

			if cfg.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			// Preflight: respond immediately, don't hit the router
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

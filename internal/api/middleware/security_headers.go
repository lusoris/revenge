package middleware

import "net/http"

// SecurityHeadersMiddleware adds standard security headers to all responses.
//
// These headers protect against common web vulnerabilities:
//   - Content-Security-Policy: restricts resource loading (protects /api/docs page)
//   - X-Content-Type-Options: prevents MIME-type sniffing
//   - X-Frame-Options: prevents clickjacking via iframes
//   - Referrer-Policy: controls referrer information leakage
//   - Permissions-Policy: disables unnecessary browser features
//
// For the API server, these are defense-in-depth. The SvelteKit frontend
// should set its own CSP via server hooks with frontend-specific directives.
func SecurityHeadersMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h := w.Header()

			// CSP: restrictive policy suitable for an API server.
			// - 'self' allows the Scalar docs page to load its own scripts/styles
			// - img-src includes image.tmdb.org for proxied poster/backdrop images
			// - media-src blob: for HLS playback segments
			// - connect-src 'self' for API calls and SSE from the docs page
			// - style-src 'unsafe-inline' required by Scalar's inline styles
			h.Set("Content-Security-Policy",
				"default-src 'self'; "+
					"script-src 'self'; "+
					"style-src 'self' 'unsafe-inline'; "+
					"img-src 'self' image.tmdb.org *.image.tmdb.org data:; "+
					"media-src 'self' blob:; "+
					"connect-src 'self'; "+
					"font-src 'self'; "+
					"object-src 'none'; "+
					"frame-ancestors 'none'; "+
					"base-uri 'self'; "+
					"form-action 'self'")

			// Prevent MIME-type sniffing (e.g., treating JSON as HTML)
			h.Set("X-Content-Type-Options", "nosniff")

			// Prevent embedding in iframes (clickjacking protection)
			h.Set("X-Frame-Options", "DENY")

			// Control referrer information sent with requests
			h.Set("Referrer-Policy", "strict-origin-when-cross-origin")

			// Disable unnecessary browser features
			h.Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

			next.ServeHTTP(w, r)
		})
	}
}

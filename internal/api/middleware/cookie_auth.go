package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/lusoris/revenge/internal/config"
)

const (
	// CookieAccessToken is the name of the HttpOnly cookie for the JWT access token.
	CookieAccessToken = "revenge_access_token"

	// CookieRefreshToken is the name of the HttpOnly cookie for the refresh token.
	CookieRefreshToken = "revenge_refresh_token"

	// CookieCSRFToken is the name of the non-HttpOnly cookie for CSRF double-submit.
	CookieCSRFToken = "revenge_csrf"

	// HeaderCSRF must match the CSRF cookie value on state-changing requests.
	HeaderCSRF = "X-CSRF-Token"
)

// CookieAuthMiddleware reads the access token from an HttpOnly cookie
// and injects it as a Bearer Authorization header so that ogen's
// security handler can validate it normally.
// If the request already has an Authorization header, it takes precedence.
func CookieAuthMiddleware(cfg config.CookieAuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if !cfg.Enabled {
			return next
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only inject if no Authorization header already present
			if r.Header.Get("Authorization") == "" {
				if c, err := r.Cookie(CookieAccessToken); err == nil && c.Value != "" {
					r.Header.Set("Authorization", "Bearer "+c.Value)
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// SetAuthCookies sets HttpOnly cookies for access and refresh tokens and a CSRF cookie.
func SetAuthCookies(w http.ResponseWriter, cfg config.CookieAuthConfig, accessToken, refreshToken, csrfToken string, maxAge int) {
	path := cfg.Path
	if path == "" {
		path = "/"
	}
	sameSite := parseSameSite(cfg.SameSite)
	secure := cfg.Secure

	// Access token cookie (HttpOnly, Secure)
	http.SetCookie(w, &http.Cookie{
		Name:     CookieAccessToken,
		Value:    accessToken,
		Path:     path,
		Domain:   cfg.Domain,
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
	})

	// Refresh token cookie (HttpOnly, Secure, restricted path)
	if refreshToken != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     CookieRefreshToken,
			Value:    refreshToken,
			Path:     "/api/v1/auth", // Only sent to auth endpoints
			Domain:   cfg.Domain,
			MaxAge:   7 * 24 * 3600, // 7 days
			HttpOnly: true,
			Secure:   secure,
			SameSite: sameSite,
		})
	}

	// CSRF token cookie (NOT HttpOnly - JS must read it)
	if csrfToken != "" {
		http.SetCookie(w, &http.Cookie{
			Name:     CookieCSRFToken,
			Value:    csrfToken,
			Path:     path,
			Domain:   cfg.Domain,
			MaxAge:   maxAge,
			HttpOnly: false, // JS needs to read this
			Secure:   secure,
			SameSite: sameSite,
		})
	}
}

// ClearAuthCookies expires all auth cookies.
func ClearAuthCookies(w http.ResponseWriter, cfg config.CookieAuthConfig) {
	path := cfg.Path
	if path == "" {
		path = "/"
	}
	for _, name := range []string{CookieAccessToken, CookieCSRFToken} {
		http.SetCookie(w, &http.Cookie{
			Name:   name,
			Value:  "",
			Path:   path,
			Domain: cfg.Domain,
			MaxAge: -1,
		})
	}
	http.SetCookie(w, &http.Cookie{
		Name:   CookieRefreshToken,
		Value:  "",
		Path:   "/api/v1/auth",
		Domain: cfg.Domain,
		MaxAge: -1,
	})
}

// responseWriterKey is the context key for the http.ResponseWriter.
type responseWriterKey struct{}

// ResponseWriterMiddleware injects the http.ResponseWriter into the request context
// so that ogen handlers can set response headers (e.g. Set-Cookie) before ogen encodes the body.
func ResponseWriterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), responseWriterKey{}, w)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetResponseWriter retrieves the http.ResponseWriter from the context.
func GetResponseWriter(ctx context.Context) (http.ResponseWriter, bool) {
	w, ok := ctx.Value(responseWriterKey{}).(http.ResponseWriter)
	return w, ok
}

func parseSameSite(s string) http.SameSite {
	switch strings.ToLower(s) {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}

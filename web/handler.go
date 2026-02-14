package web

import (
	"io/fs"
	"net/http"
	"path"
	"strings"
)

// Enabled reports whether frontend assets are embedded in the binary.
// This is true when built with -tags frontend, false otherwise.
func Enabled() bool {
	return assets != nil
}

// NewSPAHandler creates an http.Handler that serves the SvelteKit build
// output as a single-page application.
//
// Static files (JS, CSS, images, fonts) are served directly with
// appropriate cache headers. All other paths fall back to index.html
// for SvelteKit client-side routing.
//
// Returns nil if no frontend assets are embedded (built without -tags frontend).
func NewSPAHandler() http.Handler {
	if assets == nil {
		return nil
	}

	// The embed.FS roots at "build/", strip that prefix.
	var sub fs.FS
	if embedFS, ok := assets.(fs.SubFS); ok {
		var err error
		sub, err = embedFS.Sub("build")
		if err != nil {
			sub = assets
		}
	} else {
		sub = assets
	}

	fileServer := http.FileServer(http.FS(sub))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if p == "" {
			p = "index.html"
		}

		// Try to serve static file
		if f, err := sub.Open(p); err == nil {
			f.Close()
			if isImmutableAsset(p) {
				w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			}
			fileServer.ServeHTTP(w, r)
			return
		}

		// Missing static file (has extension) -> 404
		if path.Ext(p) != "" {
			http.NotFound(w, r)
			return
		}

		// SPA fallback: serve index.html for client-side routing
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	})
}

// isImmutableAsset returns true for SvelteKit hashed asset paths
// that can be cached indefinitely (_app/immutable/...).
func isImmutableAsset(p string) bool {
	return strings.HasPrefix(p, "_app/immutable/") ||
		strings.HasSuffix(p, ".js") ||
		strings.HasSuffix(p, ".css") ||
		strings.HasSuffix(p, ".woff2") ||
		strings.HasSuffix(p, ".woff")
}

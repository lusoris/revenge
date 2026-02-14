//go:build !frontend

// Package web provides embedded frontend assets for the SvelteKit SPA.
//
// This stub is compiled when the "frontend" build tag is NOT set (default).
// Backend-only development does not require frontend assets.
package web

import "io/fs"

// assets is nil when frontend is not embedded.
var assets fs.FS

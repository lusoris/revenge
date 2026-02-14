//go:build frontend

// Package web provides embedded frontend assets for the SvelteKit SPA.
//
// This file is only compiled when the "frontend" build tag is set.
// Build with: go build -tags frontend
//
// The SvelteKit project must be built first with @sveltejs/adapter-static
// which outputs to the build/ directory:
//
//cd web && npm run build
//go build -tags frontend ./cmd/revenge
package web

import "embed"

// assets contains the SvelteKit static build output.
//
//go:embed all:build
var assets embed.FS

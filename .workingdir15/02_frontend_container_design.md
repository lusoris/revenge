# Frontend Container Design

**Date**: 2026-02-14

## Architecture Decision

The frontend runs as a **separate container** from the backend:

1. **`revenge-frontend`** — nginx serving SvelteKit static adapter output
2. **`revenge`** (existing) — Go API server

### Why Separate Containers?
- Independent scaling (frontend is stateless, cacheable at CDN)
- Independent deployment cycles (UI changes don't rebuild Go binary)
- Smaller attack surface (frontend container has no DB access)
- Standard pattern: nginx serves static + reverse proxies `/api/` to backend
- The embedded SPA (`web/embed.go` with `-tags frontend`) remains available as an **alternative** single-binary deployment

### Container Strategy
- **Build stage**: `node:22-alpine` — pnpm + SvelteKit build
- **Runtime stage**: `nginx:alpine` — serves `build/` static output
- nginx handles: gzip, static caching, SPA fallback, `/api/` proxy to backend

## Files Created

| File | Purpose |
|------|---------|
| `web/Dockerfile` | Multi-stage: node build + nginx runtime |
| `web/nginx.conf` | nginx config: SPA routing, caching, API proxy |
| `web/.dockerignore` | Exclude node_modules etc from context |
| `web/package.json` | SvelteKit project root |
| `web/svelte.config.js` | SvelteKit adapter-static config |
| `web/vite.config.ts` | Vite config |
| `web/tsconfig.json` | TypeScript config |
| `web/tailwind.config.ts` | Tailwind CSS 4 config |
| `web/src/app.html` | HTML shell |
| `web/src/app.css` | Global styles (Tailwind directives) |
| `web/src/routes/+layout.svelte` | Root layout |
| `web/src/routes/+page.svelte` | Home page placeholder |
| `web/static/favicon.png` | Placeholder |

## Compose Integration

- `docker-compose.yml` + `docker-compose.dev.yml` get `frontend` service
- Frontend exposes port 3000
- API calls proxied to `revenge:8096` via nginx

## Helm Integration

- `charts/revenge/values.yaml` gets `frontend:` section
- New deployment + service + ingress path for frontend

## Makefile Integration

- `make frontend-build` — build SvelteKit
- `make frontend-dev` — dev server with HMR
- `make frontend-docker` — build frontend Docker image
- `make docker-build-all` — build both images

# Frontend Container Implementation

**Date**: 2026-02-14
**Status**: Complete

## Files Created

### SvelteKit Project (`web/`)
| File | Purpose |
|------|---------|
| `web/package.json` | SvelteKit 2 + Svelte 5 + Tailwind 4 + shadcn deps |
| `web/svelte.config.js` | adapter-static → `build/`, fallback `index.html`, precompress |
| `web/vite.config.ts` | Tailwind vite plugin, dev proxy `/api/` → `localhost:8096` |
| `web/tsconfig.json` | Strict TS, bundler module resolution |
| `web/src/app.html` | HTML shell, dark theme, sveltekit placeholders |
| `web/src/app.css` | `@import 'tailwindcss'` (v4 syntax) |
| `web/src/app.d.ts` | SvelteKit type declarations |
| `web/src/routes/+layout.svelte` | Root layout, imports app.css |
| `web/src/routes/+layout.ts` | `prerender=true`, `ssr=false` (SPA mode) |
| `web/src/routes/+page.svelte` | Placeholder home page |
| `web/src/lib/utils.ts` | `cn()` utility for shadcn-svelte (clsx + twMerge) |
| `web/.prettierrc` | Prettier config with svelte + tailwind plugins |
| `web/.prettierignore` | Ignore build output |
| `web/.gitignore` | node_modules, build, .svelte-kit |
| `web/.dockerignore` | Exclude node_modules, build from Docker context |

### Docker Configuration
| File | Purpose |
|------|---------|
| `web/Dockerfile` | Multi-stage: `node:22-alpine` build → `nginx:alpine` runtime |
| `web/nginx.conf` | SPA routing, gzip, immutable cache, `/api/` + `/hls/` proxy |
| `web/static/favicon.png` | Minimal valid 1x1 PNG placeholder |

### Updated Files
| File | Change |
|------|--------|
| `docker-compose.yml` | Added `frontend` service (port 3000, depends on revenge) |
| `docker-compose.dev.yml` | Added `frontend` service for dev stack |
| `docker-compose.prod.yml` | Added `frontend` service with GHCR image |
| `Makefile` | Added `frontend-*` targets + `docker-build-all` |
| `charts/revenge/values.yaml` | Added `frontend:` section |
| `charts/revenge/templates/frontend-deployment.yaml` | Frontend Deployment |
| `charts/revenge/templates/frontend-service.yaml` | Frontend Service |
| `charts/revenge/templates/ingress.yaml` | Split routing: `/api/`,`/hls/` → backend, `/` → frontend |

## Architecture

```
┌─────────────┐     ┌──────────────┐
│   Browser    │────▶│   Frontend   │  (nginx:3000)
│              │     │  SvelteKit   │
└─────────────┘     └──────┬───────┘
                           │ /api/, /hls/
                    ┌──────▼───────┐
                    │   Backend    │  (Go:8096)
                    │   revenge    │
                    └──────────────┘
```

## Makefile Targets
- `make frontend-install` — pnpm install
- `make frontend-dev` — dev server with HMR
- `make frontend-build` — production build
- `make frontend-check` — svelte-check
- `make frontend-docker` — build frontend Docker image
- `make docker-build-all` — build both backend + frontend images

## Deployment Modes
1. **Separate containers** (default) — frontend nginx + backend Go
2. **Embedded SPA** — `go build -tags frontend` bakes SPA into Go binary
3. **Dev mode** — `make frontend-dev` with vite proxy to `localhost:8096`

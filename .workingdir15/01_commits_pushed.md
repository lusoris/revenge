# Commits Pushed to develop

**Date**: 2026-02-14
**Total**: 10 commits, 102+ files changed

## Commit Log

| Hash | Type | Summary |
|------|------|---------|
| `445838cd` | fix(cache) | Correct L1 cache TTL and improve cache layer tests |
| `b5d6172b` | feat(session) | Improve session management with cleanup job |
| `414deb47` | feat(auth) | Link auth tokens to sessions with migration 000042 |
| `6f0b0b24` | feat(playback) | Migrate transcoding from FFmpeg CLI to astiav |
| `e215bae4` | fix(api) | Fix HTTP context bug and improve handler tests |
| `c8d27135` | feat(infra) | Add circuit breaker for all 16 external API clients |
| `151337d8` | feat(web) | Add SvelteKit SPA serving with build-tag gating |
| `a3a56182` | test | Improve load tests and live smoke tests |
| `84679941` | chore | Regenerate mocks and sqlc models |
| `2800164b` | chore | Update deps, Makefile, docker-compose, project config |

## Key Features

### Circuit Breaker (gobreaker)
- `internal/infra/circuitbreaker/` — new package
- Three tiers: External (5 fails/30s), Local (3 fails/15s), CDN (10 fails/60s)
- Wired into all 16 req.Client instances via WrapRoundTripFunc
- Prometheus metrics: state change counter + current state gauge

### SvelteKit SPA Serving
- `web/embed.go` — `//go:build frontend`, embeds `web/build/`
- `web/embed_stub.go` — `//go:build !frontend`, nil stub
- `web/handler.go` — SPA handler with immutable cache headers + index.html fallback
- `internal/api/server.go` — conditional routing: frontend → ogen on `/api/`, SPA on `/`

### astiav Migration
- Deleted FFmpeg CLI subprocess spawning
- Native libav bindings via go-astiav
- Eliminated FFmpeg process leak

### Session + Auth
- Session maintenance job for expired session cleanup
- Auth tokens linked to sessions (migration 000042)

## Full Test Suite
All packages pass (`go test ./... -count=1 -short`), 0 failures.

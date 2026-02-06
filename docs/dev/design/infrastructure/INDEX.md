# Infrastructure Documentation

> Internal infrastructure packages in `internal/infra/`

## Documents

| Document | Package | Purpose |
|----------|---------|---------|
| [DATABASE.md](DATABASE.md) | `internal/infra/database` | PostgreSQL pooling, migrations, sqlc, query logging |
| [CACHE.md](CACHE.md) | `internal/infra/cache` | L1 (otter) + L2 (rueidis/Dragonfly) unified cache |
| [JOBS.md](JOBS.md) | `internal/infra/jobs` | River job queue, 5-level priority, 17 workers |
| [HEALTH.md](HEALTH.md) | `internal/infra/health` | K8s probes, dependency checks |
| [IMAGE.md](IMAGE.md) | `internal/infra/image` | TMDb image download, proxy, caching |
| [LOGGING.md](LOGGING.md) | `internal/infra/logging` | slog + zap, tint dev mode |
| [OBSERVABILITY.md](OBSERVABILITY.md) | `internal/infra/observability` | Prometheus metrics, pprof, HTTP middleware |
| [SEARCH_INFRA.md](SEARCH_INFRA.md) | `internal/infra/search` | Typesense client wrapper |

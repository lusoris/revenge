# Architecture

**Last Updated**: 2026-02-06

Go media server with layered architecture, fx dependency injection, and ogen-generated API.

---

## System Overview

```
                    ┌──────────────┐
                    │   Clients    │  (REST API consumers)
                    └──────┬───────┘
                           │
                    ┌──────▼───────┐
                    │  API Layer   │  ogen-generated handlers + middleware
                    │  (net/http)  │  auth, metrics, CORS, request ID
                    └──────┬───────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
       ┌──────▼──────┐ ┌──▼───────┐ ┌──▼──────────┐
       │   Content    │ │ Services │ │ Integrations │
       │   Modules    │ │          │ │              │
       │ movie,tv,qar │ │ auth,    │ │ tmdb, tvdb,  │
       │              │ │ user,    │ │ radarr,      │
       │              │ │ metadata │ │ sonarr       │
       └──────┬───────┘ └────┬─────┘ └──────┬──────┘
              │              │               │
       ┌──────▼──────────────▼───────────────▼──────┐
       │              Infrastructure                 │
       │  PostgreSQL (pgx) │ Dragonfly (rueidis)    │
       │  Typesense        │ River (job queue)       │
       └────────────────────────────────────────────┘
```

---

## Project Structure

```
cmd/revenge/              Application entrypoint (main.go → fx.New)
internal/
  api/
    ogen/                 Generated API server, handlers, types (from OpenAPI spec)
    middleware/            Auth, metrics, CORS, request ID
  content/
    movie/                Movie module (handler, service, repository, jobs, types)
    tvshow/               TV show module (series/season/episode hierarchy)
    qar/                  Adult content module (pirate-themed obfuscation)
  service/
    auth/                 Authentication (JWT + refresh tokens)
    user/                 User management
    session/              Session tracking
    mfa/                  Multi-factor auth (TOTP + WebAuthn)
    oidc/                 OpenID Connect / SSO
    rbac/                 Role-based access control (Casbin)
    apikeys/              API key management
    settings/             Application settings
    activity/             Activity logging
    library/              Library management
    metadata/             Metadata aggregation service
      providers/tmdb/     TMDb provider (priority 100)
      providers/tvdb/     TVDb provider (priority 80)
      adapters/movie/     Movie adapter (metadata → movie domain types)
      adapters/tvshow/    TV show adapter (metadata → tvshow domain types)
      jobs/               River job definitions (refresh, enrich, download)
    search/               Typesense full-text search
    email/                SMTP email (go-mail)
    notification/         Notification dispatch
    storage/              File storage abstraction (local + S3)
    http/                 HTTP client factory (resty, rate limiting)
  infra/
    database/
      migrations/shared/  SQL migrations (embedded via go:embed)
      db/                 sqlc-generated queries and types
      migrate.go          Migration runner (iofs + pgx)
    cache/                Distributed cache client (rueidis → Dragonfly)
    jobs/                 River job client and worker registration
  config/                 Configuration (koanf: YAML + env vars)
api/openapi/              OpenAPI specification (source of truth for API)
charts/revenge/           Helm chart (full stack)
scripts/                  Docker entrypoint, build scripts
```

---

## Dependency Injection (fx)

All components are wired through [uber-go/fx](https://uber-go.github.io/fx/) modules. Each package exposes a `Module` variable.

**Startup sequence** (`cmd/revenge/main.go`):

```
fx.New(app.Module)
  ├─ config.Module         → *config.Config (koanf loader)
  ├─ logging.Module        → *slog.Logger (slog + zap handler)
  ├─ database.Module       → *pgxpool.Pool, *db.Queries, migrations
  ├─ cache.Module          → *cache.Client (rueidis)
  ├─ search.Module         → *search.Client (typesense)
  ├─ jobs.Module           → *river.Client, river.Workers
  ├─ auth.Module           → auth.Service, auth.TokenManager
  ├─ user.Module           → user.Service
  ├─ session.Module        → session.Service
  ├─ mfa.Module            → mfa.Service
  ├─ oidc.Module           → oidc.Service
  ├─ rbac.Module           → rbac.Service (Casbin)
  ├─ apikeys.Module        → apikeys.Service
  ├─ settings.Module       → settings.Service
  ├─ activity.Module       → activity.Service
  ├─ library.Module        → library.Service
  ├─ metadatafx.Module     → metadata.Service + providers + adapters
  ├─ radarr.Module         → radarr.SyncService
  ├─ sonarr.Module         → sonarr.SyncService
  ├─ movie.Module          → movie.Handler, movie.Service
  ├─ tvshow.Module         → tvshow.Handler, tvshow.Service
  ├─ notification.Module   → notification.Service
  ├─ email.Module          → email.Service
  ├─ storage.Module        → storage.Service
  ├─ health.Module         → health check handlers
  ├─ observability.Module  → metrics, tracing (Prometheus, OTLP)
  ├─ raft.Module           → leader election (scaffold, not operational)
  └─ api.Module            → *http.Server (ogen, middleware, lifecycle hooks)
```

**Lifecycle**: fx manages graceful startup and shutdown. OnStart hooks run in registration order (database first, HTTP server last). OnStop hooks run in reverse order.

---

## Layer Details

### API Layer

The API is **spec-first**: an OpenAPI YAML spec in `api/openapi/` generates the entire server skeleton via [ogen](https://ogen.dev/). This includes:

- Request/response types
- Route registration
- Parameter validation
- Authentication hooks (`HandleBearerAuth`)

Custom code lives in the `Handler` struct which implements ogen's generated interface. Each API method delegates to the appropriate content handler or service.

**Middleware stack** (applied in order):
1. Request ID generation
2. Structured logging (slog)
3. HTTP metrics (Prometheus)
4. CORS
5. Bearer auth (JWT validation via `TokenManager`)

### Content Modules

Each content module (`movie`, `tvshow`, `qar`) follows the same pattern:

| Component | File | Purpose |
|-----------|------|---------|
| Handler | `handler.go` | Accepts parsed API requests, delegates to service |
| Service | `service.go` | Business logic, orchestrates repository + metadata |
| Repository | `repository.go` | Interface for database operations |
| Repository (PG) | `repository_pg.go` | PostgreSQL implementation using sqlc queries |
| Types | `types.go` | Domain types (Movie, Series, Episode, etc.) |
| Jobs | `jobs/` | River worker definitions for async tasks |
| Module | `module.go` | fx module wiring |

**Movie module** is the most complete: full CRUD, metadata enrichment (TMDb + TVDb), library scanning, file matching, search indexing, credits, genres, i18n translations, age ratings, and watch progress tracking.

**TV Show module** adds hierarchical structure: Series → Season → Episode, each with independent metadata refresh.

### Service Layer

Backend services are shared across content modules:

| Service | Package | Key Responsibility |
|---------|---------|-------------------|
| auth | `service/auth` | JWT access/refresh tokens, login, password hashing (bcrypt) |
| user | `service/user` | User CRUD, profile management |
| session | `service/session` | Session tracking, device management |
| mfa | `service/mfa` | TOTP + WebAuthn second factor |
| oidc | `service/oidc` | External identity provider integration |
| rbac | `service/rbac` | Casbin policy enforcement |
| apikeys | `service/apikeys` | API key generation and validation |
| settings | `service/settings` | Application configuration storage |
| activity | `service/activity` | User activity audit log |
| library | `service/library` | Library scanning, file discovery |
| metadata | `service/metadata` | Multi-provider metadata aggregation (see [Metadata System](METADATA_SYSTEM.md)) |
| search | `service/search` | Typesense indexing and querying |
| email | `service/email` | SMTP delivery via go-mail |
| notification | `service/notification` | Multi-channel notification dispatch |
| storage | `service/storage` | File storage abstraction (local + S3) |
| http | `service/http` | HTTP client factory with rate limiting |

### Infrastructure Layer

| Component | Package | Technology |
|-----------|---------|-----------|
| Database | `infra/database` | PostgreSQL 18+ via pgxpool, sqlc codegen, embedded migrations |
| Cache | `infra/cache` | Dragonfly (Redis-compatible) via rueidis |
| Jobs | `infra/jobs` | River PostgreSQL-native job queue |
| Search | via `service/search` | Typesense full-text search |

### Integration Layer

External service integrations:

| Integration | Location | Purpose |
|-------------|----------|---------|
| TMDb | `service/metadata/providers/tmdb/` | Movie + TV metadata (primary, priority 100) |
| TVDb | `service/metadata/providers/tvdb/` | TV metadata (secondary, priority 80) |
| Radarr | `integration/servarr/radarr/` | Movie library sync |
| Sonarr | `integration/servarr/sonarr/` | TV show library sync |
| Generic OIDC | `integration/auth/oidc/` | External identity providers |

---

## Request Flow

A typical API request flows through these layers:

```
HTTP Request
  → ogen router (generated)
  → Middleware (auth, logging, metrics)
  → Handler.Method() (implements ogen interface)
    → content handler (movie/tvshow)
      → service (business logic)
        → repository (sqlc queries)
          → PostgreSQL
        → metadata provider (if enrichment needed)
          → TMDb/TVDb API
  → ogen response encoder (generated)
HTTP Response
```

---

## Background Jobs

17 River workers across 5 priority queues handle async processing:

| Queue | Workers | Purpose |
|-------|---------|---------|
| `critical` | 20 | Security events, auth failures, urgent tasks |
| `high` | 15 | Notifications, webhooks, Radarr/Sonarr sync |
| `default` | 10 | Metadata refresh, file matching, general tasks |
| `low` | 5 | Cleanup, maintenance (leader-aware) |
| `bulk` | 3 | Library scans, search reindexing |

Key workers: MovieLibraryScan, TVShowLibraryScan (bulk), MetadataRefresh (default), RadarrSync/SonarrSync (high), NotificationWorker (high), CleanupWorker (low).

See [JOBS.md](../infrastructure/JOBS.md) for the full worker list and queue assignments.

Workers are registered via fx and managed by River's lifecycle. Each worker has a configurable timeout.

---

## Configuration

Configuration uses [koanf](https://github.com/knadh/koanf) with YAML files + environment variable overrides.

**Key config areas:**

```yaml
server:
  address: ":8080"
  log_level: "info"

database:
  url: "postgres://user:pass@localhost:5432/revenge"
  max_conns: 25

cache:
  url: "redis://localhost:6379"

search:
  url: "http://localhost:8108"
  api_key: "..."

metadata:
  tmdb:
    api_key: "..."
  tvdb:
    api_key: "..."

radarr:
  base_url: "http://localhost:7878"
  api_key: "..."

sonarr:
  base_url: "http://localhost:8989"
  api_key: "..."
```

Environment variables override YAML with `REVENGE_` prefix (e.g., `REVENGE_DATABASE_URL`).

---

## Health Checks

Kubernetes-ready health endpoints:

| Endpoint | Purpose | Checks |
|----------|---------|--------|
| `/healthz` | Liveness | Process alive |
| `/readyz` | Readiness | Database + cache connectivity |
| `/startupz` | Startup | Migrations complete, services initialized |

---

## Related Documentation

- [Design Principles](DESIGN_PRINCIPLES.md) - Core patterns and conventions
- [Metadata System](METADATA_SYSTEM.md) - Multi-provider metadata aggregation
- [Tech Stack](../technical/TECH_STACK.md) - Technology choices and rationale
- [API](../technical/API.md) - OpenAPI specification details

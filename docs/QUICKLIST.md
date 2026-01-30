# Revenge - Quick Reference by Task

> Load the right docs into context based on what you're working on.

**Last Updated**: 2026-01-30

---

## 🎯 Task-Based Quick Links

### Adding a New API Endpoint

**Read these:**
1. [API.md](dev/design/technical/API.md) - OpenAPI spec-first approach
2. `api/openapi/*.yaml` - Existing OpenAPI specs (source of truth)
3. `internal/api/` - Handler implementations

**Code locations:**
- OpenAPI specs: `api/openapi/`
- Handlers: `internal/api/`
- Converters: `internal/api/converters.go`
- Module wiring: `internal/api/module.go`

---

### Working on RBAC / Permissions

**Read these:**
1. [RBAC_CASBIN.md](dev/design/features/shared/RBAC_CASBIN.md) - Full Casbin RBAC design
2. [ACCESS_CONTROLS.md](dev/design/features/shared/ACCESS_CONTROLS.md) - Permission system overview

**Code locations:**
- Service: `internal/service/rbac/casbin.go`
- Model: `internal/service/rbac/model.conf`
- Queries: `internal/infra/database/queries/shared/roles.sql`
- Migration: `internal/infra/database/migrations/shared/`

**Key concepts:**
- Casbin enforces role→permission mapping
- Roles stored in `roles` table, policies in `casbin_rules`
- Resource grants for per-item permissions (libraries, content)

---

### Working on Request System

**Read these:**
1. [REQUEST_SYSTEM.md](dev/design/features/shared/REQUEST_SYSTEM.md) - Native request system design

**Code locations:**
- Schema: `qar` for adult, `shared` for general
- Queries: `internal/infra/database/queries/shared/requests.sql`
- Service: `internal/service/request/` (to create)

**Key concepts:**
- Requests can be movies, TV shows, or adult content
- Auto-approval rules, quotas, voting, and polls
- Integration with Radarr/Sonarr/Whisparr

---

### Working on Adult Content (QAR)

**Read these:**
1. [ADULT_CONTENT_SYSTEM.md](dev/design/features/adult/ADULT_CONTENT_SYSTEM.md) - QAR obfuscation design
2. [ADULT_METADATA.md](dev/design/features/adult/ADULT_METADATA.md) - Metadata providers

**Code locations:**
- Domain: `internal/content/qar/` (expedition, voyage, crew, port, flag, fleet)
- Queries: `internal/infra/database/queries/qar/`
- Handlers: `internal/api/adult.go`
- Schema: `qar.*` tables (PostgreSQL schema isolation)

**Obfuscation mapping:**
| Real Term | QAR Term |
|-----------|----------|
| Movie | Expedition |
| Scene | Voyage |
| Performer | Crew |
| Studio | Port |
| Tag | Flag |
| Library | Fleet |

---

### Working on Playback / Continue Watching

**Read these:**
1. [WATCH_NEXT_CONTINUE_WATCHING.md](dev/design/features/playback/WATCH_NEXT_CONTINUE_WATCHING.md) - Playback continuation
2. [INDEX.md (playback)](dev/design/features/playback/INDEX.md) - All playback features

**Code locations:**
- Service: `internal/service/playback/`
- Types: `internal/service/playback/types.go`
- Movie queries: `internal/infra/database/queries/movies/user_data.sql`
- TV queries: `internal/infra/database/queries/tvshows/user_data.sql`

**Key concepts:**
- 30-day window for continue watching
- `mark_watched_percent` (default 90%) determines completion
- In-memory session tracking with database persistence

---

### Working on Search (Typesense)

**Read these:**
1. [TYPESENSE.md](dev/design/integrations/infrastructure/TYPESENSE.md) - Search engine setup

**Code locations:**
- Service: `internal/service/search/`
- Collections: defined per content type
- Indexer workers: `internal/jobs/`

**Key concepts:**
- Separate collections per content type (movies, tvshows, qar_expeditions, etc.)
- Real-time indexing via River jobs
- Multi-search for combined results

---

### Working on Metadata Providers

**Read these:**
1. [METADATA_SYSTEM.md](dev/design/architecture/METADATA_SYSTEM.md) - Overall metadata architecture
2. Provider-specific: [TMDB.md](dev/design/integrations/metadata/video/TMDB.md), [STASHDB.md](dev/design/integrations/metadata/adult/STASHDB.md)

**Code locations:**
- TMDb: `internal/service/metadata/tmdb/`
- StashDB: `internal/service/metadata/stashdb/`
- Whisparr: `internal/service/metadata/whisparr/`

**Key concepts:**
- Rate limiting with circuit breakers
- Metadata refresh jobs via River
- Provider priority for fallbacks

---

### Working on Jobs (River)

**Read these:**
1. [RIVER.md](dev/design/integrations/infrastructure/RIVER.md) - Job queue design

**Code locations:**
- Workers: `internal/jobs/`
- Module: `internal/infra/jobs/module.go`

**Key concepts:**
- PostgreSQL-native job queue
- Workers registered via fx
- Retry policies per job type

---

### Working on Caching

**Read these:**
1. [DRAGONFLY.md](dev/design/integrations/infrastructure/DRAGONFLY.md) - Redis-compatible cache

**Code locations:**
- Service: `internal/infra/cache/`
- Uses: `rueidis` (NOT go-redis)
- Local cache: `otter` for in-memory
- API cache: `sturdyc` for request coalescing

---

### Working on Database / Migrations

**Read these:**
1. [POSTGRESQL.md](dev/design/integrations/infrastructure/POSTGRESQL.md) - Database design

**Code locations:**
- Migrations: `internal/infra/database/migrations/`
  - `shared/` - Common tables
  - `movies/` - Movie module
  - `tvshows/` - TV module
  - `qar/` - Adult content
- Queries: `internal/infra/database/queries/`
- sqlc config: `sqlc.yaml`

**Commands:**
```bash
sqlc generate          # Generate Go code from SQL
go generate ./api/...  # Generate OpenAPI code
```

---

### Working on Authentication

**Read these:**
1. [Auth integrations](dev/design/integrations/auth/INDEX.md) - OIDC providers

**Code locations:**
- Session service: `internal/service/session/`
- Auth middleware: `internal/api/middleware/`
- User queries: `internal/infra/database/queries/shared/users.sql`

**Key concepts:**
- Session tokens hashed with SHA-256
- Support for external OIDC (Authelia, Authentik, Keycloak)
- Device tracking per session

---

### Working on Configuration

**Read these:**
1. [CONFIGURATION.md](dev/design/technical/CONFIGURATION.md) - Config structure

**Code locations:**
- Config: `pkg/config/`
- Uses: `koanf` (NOT viper)

---

## 📁 Directory Structure

```
internal/
├── api/                    # HTTP handlers (ogen-generated + custom)
├── content/
│   ├── movies/             # Movie domain
│   ├── tvshows/            # TV show domain
│   └── qar/                # Adult content (QAR obfuscation)
│       ├── expedition/     # Adult movies
│       ├── voyage/         # Adult scenes
│       ├── crew/           # Performers
│       ├── port/           # Studios
│       ├── flag/           # Tags
│       └── fleet/          # Libraries
├── infra/
│   ├── cache/              # Redis/Otter caching
│   ├── database/           # PostgreSQL + sqlc
│   │   ├── migrations/     # SQL migrations by module
│   │   ├── queries/        # sqlc query files
│   │   └── db/             # Generated sqlc code
│   ├── health/             # Health checks
│   ├── jobs/               # River job queue
│   └── search/             # Typesense
├── jobs/                   # River workers
└── service/                # Business logic services
    ├── metadata/           # Metadata providers
    ├── playback/           # Playback session management
    ├── rbac/               # Casbin RBAC
    ├── request/            # Content requests
    └── session/            # User sessions
```

---

## 🔧 Tech Stack Quick Reference

| Component | Package | Notes |
|-----------|---------|-------|
| Cache (distributed) | `github.com/redis/rueidis` | NOT go-redis |
| Cache (local) | `github.com/maypok86/otter` v1.2.4 | W-TinyLFU |
| Cache (API) | `github.com/viccon/sturdyc` v1.1.5 | Request coalescing |
| Search | `github.com/typesense/typesense-go/v4` | NOT v3 |
| Config | `github.com/knadh/koanf/v2` | NOT viper |
| Logging | `log/slog` | NOT zap |
| Jobs | `github.com/riverqueue/river` | PostgreSQL-native |
| RBAC | `github.com/casbin/casbin/v2` | Dynamic roles |
| DI | `go.uber.org/fx` | Dependency injection |
| HTTP client | `github.com/go-resty/resty/v2` | External APIs |

---

## 📋 Common Commands

```bash
# Build with experiments
GOEXPERIMENT=greenteagc,jsonv2 go build -o bin/revenge ./cmd/revenge

# Generate code
sqlc generate
go generate ./api/...

# Lint
golangci-lint run

# Test
go test ./...
```

---

## 🔗 Design Doc Index Files

Use these to find all docs in a category:

- [Main Index](INDEX.md)
- [Features Index](dev/design/features/INDEX.md)
- [Integrations Index](dev/design/integrations/INDEX.md)
- [Technical Index](dev/design/technical/INDEX.md)
- [Architecture Index](dev/design/architecture/INDEX.md)
- [Operations Index](dev/design/operations/INDEX.md)

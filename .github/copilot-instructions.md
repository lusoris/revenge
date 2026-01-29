# Revenge - Copilot Instructions

> Modular media server with complete content isolation

## Project Overview

**Revenge** is a ground-up media server built in Go with a fully modular architecture.

| Aspect       | Value                        |
| ------------ | ---------------------------- |
| Language     | Go 1.25                      |
| Database     | PostgreSQL 18+               |
| Cache        | Dragonfly (Redis-compatible) |
| Search       | Typesense 0.25+              |
| Job Queue    | River (PostgreSQL-native)    |
| API Docs     | ogen (OpenAPI spec-first)    |
| Architecture | Modular (12 content modules) |

## Architecture

See [docs/dev/design/architecture/ARCHITECTURE_V2.md](../docs/dev/design/architecture/ARCHITECTURE_V2.md) for complete design.

### Content Modules (12)

- `movie`, `tvshow` - Video content (public schema)
- `music` - Artists, albums, tracks
- `audiobook`, `book` - Reading content
- `podcast` - Podcasts & episodes
- `photo` - Photos & albums
- `livetv` - Channels, programs, DVR
- `comics` - Comics, manga, graphic novels
- `collection` - Cross-module collections
- `adult_movie`, `adult_show` - Adult content (isolated `c` schema, `/c/` API namespace)

### Key Design Principles

1. **No shared content tables** - Each module has own optimized tables
2. **Per-module everything** - Ratings, history, favorites, metadata
3. **Adult isolation** - Separate PostgreSQL schema `c` and API namespace `/c/`
4. **Module enable/disable** - Independent activation
5. **External transcoding** - Delegate to "Blackbeard" service

## Quick Commands

```bash
# Development
go run ./cmd/revenge
make dev                    # Docker Compose + hot reload

# Build & Test
go build -o bin/revenge ./cmd/revenge
go test ./...
go test -tags=integration ./tests/integration/...

# Generate
sqlc generate               # Database queries
go generate ./api/...       # ogen API handlers
go generate ./...           # Stringer, etc.

# Lint
golangci-lint run
```

## Project Structure

```
cmd/revenge/main.go           # Entry point with fx.New()
api/
  openapi/                    # OpenAPI specs (ogen)
  generated/                  # ogen-generated handlers
internal/
  content/                    # Content modules (isolated)
    movie/                    # Movie module
    tvshow/                   # TV show module
    music/                    # Music module
    c/                        # Adult modules (obscured)
      movie/
      show/
  service/                    # Shared services
    auth/                     # Authentication
    user/                     # User management
    oidc/                     # SSO/OIDC
    library/                  # Library management
    playback/                 # Playback session management
      client.go               # Client detection & capabilities
      bandwidth.go            # Bandwidth monitoring (external)
      transcoder.go           # Blackbeard integration
      session.go              # Playback session state
  infra/
    database/                 # PostgreSQL + sqlc
      migrations/             # Per-module migrations
      queries/                # Per-module SQL queries
    cache/                    # Dragonfly client
    search/                   # Typesense client
    jobs/                     # River job queue
pkg/config/                   # koanf configuration
```

## Core Stack

- **DI**: `go.uber.org/fx` v1.24+ (see `fx-dependency-injection.instructions.md`)
- **Config**: `github.com/knadh/koanf/v2` (see `koanf-configuration.instructions.md`)
- **Database**: `pgx/v5` + `sqlc` (see `sqlc-database.instructions.md`)
- **Cache (Dragonfly)**: `github.com/redis/rueidis` (14x faster than go-redis)
- **Cache (Local)**: `github.com/maypok86/otter` (W-TinyLFU local cache)
- **Cache (API)**: `github.com/viccon/sturdyc` (request coalescing)
- **Search**: `github.com/typesense/typesense-go/v3`
- **Job Queue**: `github.com/riverqueue/river`
- **API**: `github.com/ogen-go/ogen` (OpenAPI spec-first)
- **HTTP Client**: `resty.dev/v3` (metadata provider calls)
- **WebSocket**: `github.com/coder/websocket` (Watch Party, live updates)
- **File Watch**: `github.com/fsnotify/fsnotify` (library scanning)
- **Routing**: `net/http` stdlib (Go 1.22+ patterns)
- **Logging**: `log/slog` stdlib

## Do's and Don'ts

### DO

- ✅ Use `context.Context` as first parameter
- ✅ Use `slog` for logging, `errors.Is/As` for error checking
- ✅ Use Go 1.22+ HTTP routing: `mux.HandleFunc("GET /api/movies/{id}", h.GetMovie)`
- ✅ Use `sync.WaitGroup.Go` (Go 1.25) instead of `wg.Add(1); go func()`
- ✅ Use `testing.B.Loop()` for benchmarks (Go 1.24+)
- ✅ Keep modules isolated - no cross-module imports
- ✅ Per-module tables, services, handlers
- ✅ Use River for background jobs
- ✅ Use ogen for API handlers

### DON'T

- ❌ Use `init()` - use fx constructors
- ❌ Use global variables - inject dependencies
- ❌ Use `panic` for errors
- ❌ Use gorilla/mux, gorilla/websocket, viper, zap, logrus, lib/pq
- ❌ Use go-redis/v9 - use rueidis instead
- ❌ Use ristretto - use otter instead
- ❌ Use `automaxprocs` - Go 1.25 has built-in container support
- ❌ Share tables between content modules
- ❌ Use polymorphic references
- ❌ Transcode internally - use external Blackbeard service

## Commit Convention

```
type(scope): description

Types: feat, fix, docs, refactor, perf, test, ci, chore
Scope: api, db, auth, movie, tvshow, music, config, etc.

Example: feat(movie): add movie metadata endpoints
```

## Module Development

See `.github/instructions/content-modules.instructions.md` for patterns.
See `.github/instructions/adult-modules.instructions.md` for adult content isolation.

## Detailed Instructions

Path-specific instructions in `.github/instructions/`:

### Core Patterns

- `go-features.instructions.md` - Go 1.25 features
- `fx-dependency-injection.instructions.md` - DI patterns
- `testing-patterns.instructions.md` - Test patterns

### Data & Storage

- `sqlc-database.instructions.md` - Database queries
- `migrations.instructions.md` - Database migrations
- `dragonfly-cache.instructions.md` - Remote caching (rueidis)
- `otter-local-cache.instructions.md` - Local caching (W-TinyLFU)
- `sturdyc-api-cache.instructions.md` - API response caching
- `typesense-search.instructions.md` - Search indexing

### API & HTTP

- `ogen-api.instructions.md` - OpenAPI code generation
- `koanf-configuration.instructions.md` - Config management
- `revenge-api-compatibility.instructions.md` - API design
- `resty-http-client.instructions.md` - External API calls
- `websocket.instructions.md` - WebSocket (Watch Party, live updates)

### Content Modules

- `content-modules.instructions.md` - Module development
- `adult-modules.instructions.md` - Adult content isolation (`/c/` namespace)
- `metadata-providers.instructions.md` - TMDb, MusicBrainz, etc.
- `fsnotify-file-watching.instructions.md` - Library file watching

### Services & Jobs

- `river-job-queue.instructions.md` - Background jobs
- `oidc-authentication.instructions.md` - OIDC/SSO
- `external-services.instructions.md` - Scrobbling, sync

### Playback & Streaming

- `streaming-best-practices.instructions.md` - Streaming patterns
- `player-architecture.instructions.md` - Player components
- `client-detection.instructions.md` - Client capabilities
- `disk-cache.instructions.md` - Transcode caching
- `offloading-patterns.instructions.md` - Blackbeard integration

### Resilience & Operations

- `resilience-patterns.instructions.md` - Circuit breakers, retries
- `self-healing.instructions.md` - Supervision, graceful shutdown
- `health-checks.instructions.md` - Health check patterns
- `hotreload.instructions.md` - Runtime config reload
- `lazy-initialization.instructions.md` - Lazy services
- `observability.instructions.md` - Metrics, logging

### Frontend

- `frontend-architecture.instructions.md` - SvelteKit patterns

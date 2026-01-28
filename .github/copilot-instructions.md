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
| Architecture | Modular (11 content modules) |

## Architecture

See [docs/dev/design/architecture/ARCHITECTURE_V2.md](../docs/dev/design/architecture/ARCHITECTURE_V2.md) for complete design.

### Content Modules (11)

- `movie`, `tvshow` - Video content (public schema)
- `music` - Artists, albums, tracks
- `audiobook`, `book` - Reading content
- `podcast` - Podcasts & episodes
- `photo` - Photos & albums
- `livetv` - Channels, programs, DVR
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
- **Cache**: `github.com/redis/go-redis/v9` (Dragonfly)
- **Search**: `github.com/typesense/typesense-go/v4`
- **Job Queue**: `github.com/riverqueue/river`
- **API**: `github.com/ogen-go/ogen` (OpenAPI spec-first)
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
- ❌ Use gorilla/mux, viper, zap, logrus, lib/pq
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

- `go-features.instructions.md` - Go 1.25 features
- `fx-dependency-injection.instructions.md` - DI patterns
- `sqlc-database.instructions.md` - Database queries
- `koanf-configuration.instructions.md` - Config management
- `testing-patterns.instructions.md` - Test patterns
- `revenge-api-compatibility.instructions.md` - API design
- `content-modules.instructions.md` - Module development
- `adult-modules.instructions.md` - Adult content isolation (`/c/` namespace)
- `oidc-authentication.instructions.md` - OIDC/SSO

# Jellyfin Go - Copilot Instructions

> Complete rewrite of Jellyfin Media Server in Go

## Project Overview

**Jellyfin Go** rewrites [Jellyfin](https://github.com/jellyfin/jellyfin) from C# to Go with 100% API compatibility.

| Aspect | Value |
|--------|-------|
| Source | C# Jellyfin (upstream: `Jellyfin.*`, `MediaBrowser.*`, `Emby.*`) |
| Target | Go 1.25 |
| Database | PostgreSQL 18+ |
| Cache | Dragonfly (Redis-compatible) |
| Search | Typesense 0.25+ |

## Quick Commands

```bash
# Development
go run ./cmd/jellyfin
make dev                    # Docker Compose + hot reload

# Build & Test
go build -o bin/jellyfin ./cmd/jellyfin
go test ./...
go test -tags=integration ./tests/integration/...

# Generate
sqlc generate               # Database queries
go generate ./...           # Stringer, etc.

# Lint
golangci-lint run
```

## Project Structure

```
cmd/jellyfin/main.go       # Entry point with fx.New()
internal/
  api/handlers/            # HTTP handlers (Go 1.22+ routing)
  api/middleware/          # Auth, logging, CORS
  domain/                  # Core entities (User, Media, Library)
  service/                 # Business logic
  infra/database/          # PostgreSQL + sqlc
  infra/cache/             # Dragonfly
pkg/config/                # koanf configuration
pkg/logger/                # slog setup
migrations/                # golang-migrate SQL files
```

## Core Stack

- **DI**: `go.uber.org/fx` v1.24+ (see `fx-dependency-injection.instructions.md`)
- **Config**: `github.com/knadh/koanf/v2` (see `koanf-configuration.instructions.md`)
- **Database**: `pgx/v5` + `sqlc` (see `sqlc-database.instructions.md`)
- **Routing**: `net/http` stdlib (Go 1.22+ patterns)
- **Logging**: `log/slog` stdlib

## Do's and Don'ts

### DO
- ✅ Use `context.Context` as first parameter
- ✅ Use `slog` for logging, `errors.Is/As` for error checking
- ✅ Use Go 1.22+ HTTP routing: `mux.HandleFunc("GET /api/users/{id}", h.GetUser)`
- ✅ Use `sync.WaitGroup.Go` (Go 1.25) instead of `wg.Add(1); go func()`
- ✅ Use `testing.B.Loop()` for benchmarks (Go 1.24+)
- ✅ Use `net/http.CrossOriginProtection` for CSRF (Go 1.25)

### DON'T
- ❌ Use `init()` - use fx constructors
- ❌ Use global variables - inject dependencies
- ❌ Use `panic` for errors
- ❌ Use gorilla/mux, viper, zap, logrus, lib/pq
- ❌ Use `automaxprocs` - Go 1.25 has built-in container support

## API Compatibility

Match original C# Jellyfin exactly:
1. Same routes, methods, query params
2. Same JSON response structure
3. Same validation and error codes

```csharp
// C# Original
[HttpGet("Users/{userId}")]
public ActionResult<UserDto> GetUserById([FromRoute] Guid userId)
```

```go
// Go Equivalent
// Route: GET /Users/{userId}
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
    userID, _ := uuid.Parse(r.PathValue("userId"))
    // ... return same JSON structure as UserDto
}
```

## Commit Convention

```
type(scope): description

Types: feat, fix, docs, refactor, perf, test, ci, chore
Scope: api, db, auth, media, config, etc.

Example: feat(api): add user authentication endpoints
```

## Detailed Instructions

Path-specific instructions in `.github/instructions/`:
- `go-features.instructions.md` - Go 1.25 features
- `fx-dependency-injection.instructions.md` - DI patterns
- `sqlc-database.instructions.md` - Database queries
- `koanf-configuration.instructions.md` - Config management
- `testing-patterns.instructions.md` - Test patterns
- `jellyfin-api-compatibility.instructions.md` - C# to Go translation
- `oidc-authentication.instructions.md` - OIDC/SSO

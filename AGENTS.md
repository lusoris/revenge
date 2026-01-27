# Jellyfin Go - Agent Instructions

> Instructions for automated coding agents (Copilot coding agent, Claude, etc.)

## Project Context

This is a **ground-up rewrite** of Jellyfin Media Server from C# to Go. The original C# code is in this repo under `Jellyfin.*`, `MediaBrowser.*`, `Emby.*` directories - use these as reference for API compatibility.

## Build & Test Commands

```bash
# ALWAYS run these before committing
go build ./...              # Must pass
go test ./...               # Must pass
golangci-lint run           # Must pass

# Generate code after schema/query changes
sqlc generate               # Regenerates internal/infra/database/db/

# Integration tests (requires Docker)
go test -tags=integration ./tests/integration/...
```

## Key Files & Locations

| Purpose | Path |
|---------|------|
| Entry point | `cmd/jellyfin/main.go` |
| HTTP handlers | `internal/api/handlers/` |
| Middleware | `internal/api/middleware/` |
| Business logic | `internal/service/` |
| Domain entities | `internal/domain/` |
| Database (sqlc) | `internal/infra/database/` |
| SQL queries | `internal/infra/database/queries/` |
| SQL migrations | `migrations/` |
| Configuration | `pkg/config/` |
| Config files | `configs/*.yaml` |

## Code Patterns to Follow

### 1. Dependency Injection (fx)

```go
// All constructors must be registered with fx
func NewUserService(db *pgxpool.Pool, logger *slog.Logger) *UserService {
    return &UserService{db: db, logger: logger}
}

// In module registration
fx.Provide(NewUserService)
```

### 2. HTTP Handlers (Go 1.22+)

```go
// Use stdlib routing patterns
mux.HandleFunc("GET /api/users/{id}", h.GetUser)
mux.HandleFunc("POST /api/users", h.CreateUser)

// Get path values
id := r.PathValue("id")
```

### 3. Database Queries (sqlc)

```sql
-- Add queries to internal/infra/database/queries/*.sql
-- name: GetUser :one
SELECT * FROM users WHERE id = $1;
```

Then run `sqlc generate`.

### 4. Error Handling

```go
// Always wrap errors with context
if err != nil {
    return fmt.Errorf("failed to get user %s: %w", id, err)
}

// Use sentinel errors
var ErrUserNotFound = errors.New("user not found")
```

### 5. Logging (slog)

```go
slog.Info("user created",
    slog.String("user_id", id),
    slog.String("email", email),
)
```

## CI/CD Checks

Pull requests must pass:
1. `go build ./...` - Compilation
2. `go test ./...` - Unit tests  
3. `golangci-lint run` - Linting
4. `go vet ./...` - Static analysis

## What NOT to Do

- ❌ Don't use `init()` functions - use fx
- ❌ Don't use global variables - inject via fx
- ❌ Don't use `panic()` for errors
- ❌ Don't use gorilla/mux - use stdlib
- ❌ Don't use Viper - use koanf
- ❌ Don't use zap/logrus - use slog
- ❌ Don't use lib/pq - use pgx/v5

## API Compatibility Requirement

When implementing API endpoints, match the original C# Jellyfin exactly:
- Same route paths and HTTP methods
- Same query parameters
- Same JSON response structure
- Same error codes

Reference: `Jellyfin.Api/Controllers/` for original implementations.

## Testing Requirements

- Write table-driven tests with `t.Run()`
- Use `testing.B.Loop()` for benchmarks (Go 1.24+)
- Test coverage targets: services 80%+, handlers 70%+
- Integration tests require `//go:build integration` tag

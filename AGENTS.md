# Revenge - Agent Instructions

> Instructions for automated coding agents (Copilot coding agent, Claude, etc.)

## Project Context

**Revenge** is a modular media server with 11 isolated content modules. Each module has its own tables, services, and handlers - no shared content tables.

See [docs/ARCHITECTURE_V2.md](docs/ARCHITECTURE_V2.md) for architecture.
See [docs/MODULE_IMPLEMENTATION_TODO.md](docs/MODULE_IMPLEMENTATION_TODO.md) for implementation phases.

## Build & Test Commands

```bash
# ALWAYS run these before committing
go build ./...              # Must pass
go test ./...               # Must pass
golangci-lint run           # Must pass

# Generate code after schema/query changes
sqlc generate               # Regenerates internal/infra/database/db/
go generate ./api/...       # Regenerates ogen API handlers

# Integration tests (requires Docker)
go test -tags=integration ./tests/integration/...
```

## Key Files & Locations

| Purpose | Path |
|---------|------|
| Entry point | `cmd/revenge/main.go` |
| OpenAPI specs | `api/openapi/` |
| Generated handlers | `api/generated/` |
| Middleware | `internal/api/middleware/` |
| Content modules | `internal/content/<module>/` |
| Shared services | `internal/service/<module>/` |
| Domain entities | `internal/domain/` |
| Database (sqlc) | `internal/infra/database/` |
| SQL queries | `internal/infra/database/queries/` |
| SQL migrations | `internal/infra/database/migrations/` |
| Cache client | `internal/infra/cache/` |
| Search client | `internal/infra/search/` |
| Job queue | `internal/infra/jobs/` |
| Configuration | `pkg/config/` |
| Config files | `configs/*.yaml` |

## Module Structure

Each content module follows this pattern:

```
internal/content/<module>/
  entity.go      # Domain entities
  repository.go  # Repository interface
  service.go     # Business logic
  handler.go     # HTTP handlers (ogen interface)
  jobs.go        # River job definitions
  module.go      # fx module registration

internal/infra/database/
  migrations/<module>/
    000001_<module>.up.sql
    000001_<module>.down.sql
  queries/<module>/
    <module>.sql
```

## Code Patterns

### 1. Dependency Injection (fx)

```go
func NewMovieService(db *pgxpool.Pool, logger *slog.Logger) *MovieService {
    return &MovieService{db: db, logger: logger}
}

// In module.go
fx.Provide(NewMovieService)
```

### 2. HTTP Handlers (Go 1.22+)

```go
mux.HandleFunc("GET /api/movies/{id}", h.GetMovie)
mux.HandleFunc("POST /api/movies", h.CreateMovie)

id := r.PathValue("id")
```

### 3. Database Queries (sqlc)

```sql
-- queries/movie/movies.sql
-- name: GetMovie :one
SELECT * FROM movies WHERE id = $1;
```

### 4. River Jobs

```go
type ScanLibraryArgs struct {
    LibraryID uuid.UUID `json:"library_id"`
}

func (ScanLibraryArgs) Kind() string { return "movie.scan_library" }
```

### 5. Error Handling

```go
if err != nil {
    return fmt.Errorf("failed to get movie %s: %w", id, err)
}

var ErrMovieNotFound = errors.New("movie not found")
```

## CI/CD Checks

Pull requests must pass:
1. `go build ./...` - Compilation
2. `go test ./...` - Unit tests
3. `golangci-lint run` - Linting
4. `go vet ./...` - Static analysis

## Critical Rules

### DO
- ✅ Keep modules isolated - no cross-module imports for content
- ✅ Per-module tables (movies, movie_genres, movie_people, etc.)
- ✅ Per-module user data (movie_user_ratings, movie_watch_history, etc.)
- ✅ Use `context.Context` as first parameter
- ✅ Use `slog` for logging
- ✅ Use River for background jobs
- ✅ Use ogen for API handlers

### DON'T
- ❌ Share content tables between modules
- ❌ Use polymorphic references (item_type + item_id)
- ❌ Use `init()` functions
- ❌ Use global variables
- ❌ Use `panic()` for errors
- ❌ Transcode internally (use Blackbeard service)

## Adult Content

Adult modules (`adult_movie`, `adult_show`) use isolated PostgreSQL schema `c`:

```sql
CREATE SCHEMA IF NOT EXISTS c;
-- All adult tables in c.* schema
```

API namespace: `/api/v1/c/movies`, `/api/v1/c/shows`

See `.github/instructions/adult-modules.instructions.md` for details.

## Testing

- Write table-driven tests with `t.Run()`
- Use `testing.B.Loop()` for benchmarks (Go 1.24+)
- Coverage targets: services 80%+, handlers 70%+
- Integration tests: `//go:build integration` tag

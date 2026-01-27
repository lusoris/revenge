# Jellyfin Go - Copilot Instructions

> Complete rewrite of Jellyfin Media Server in Go.
> This document provides modular instructions for AI-assisted development.

---

## Project Overview

**Jellyfin Go** is a ground-up rewrite of the [Jellyfin](https://github.com/jellyfin/jellyfin) media server from C# to Go. The goal is 100% API compatibility with the original Jellyfin while leveraging Go's performance, simplicity, and deployment advantages.

### Key Facts

- **Source**: C# Jellyfin (upstream in this repo under Jellyfin.*/MediaBrowser.*)
- **Target**: Go 1.24 with modern idioms
- **Versioning**: SemVer 0.x until feature parity, then v1.0.0
- **Architecture**: Clean Architecture (Hexagonal) with DI
- **Database**: PostgreSQL 16+ (required)
- **Cache**: Dragonfly (Redis-compatible, required)
- **Search**: Typesense 0.25+ (required)

---

## Technology Stack (Go 1.24 - February 2025)

### Go 1.24 New Features (USE THESE!)

```go
// Generic type aliases - fully supported
type Set[T comparable] = map[T]struct{}

// Tool directives in go.mod - no more tools.go
//go:tool golang.org/x/tools/cmd/stringer

// Use testing.B.Loop for benchmarks (not b.N)
func BenchmarkFoo(b *testing.B) {
    for b.Loop() {
        // benchmark code
    }
}

// runtime.AddCleanup instead of SetFinalizer
runtime.AddCleanup(obj, func(ptr *MyType) {
    ptr.Close()
})

// encoding/json omitzero tag
type Config struct {
    Timeout time.Duration `json:"timeout,omitzero"`
}

// os.Root for directory-limited filesystem access
root, _ := os.OpenRoot("/data")
f, _ := root.Open("file.txt") // Can't escape /data
```

### Core Dependencies

| Package                     | Version | Purpose                          |
| --------------------------- | ------- | -------------------------------- |
| `go.uber.org/fx`            | v1.24+  | Dependency injection             |
| `github.com/knadh/koanf/v2` | v2.3+   | Configuration management         |
| `github.com/jackc/pgx/v5`   | v5.x    | PostgreSQL driver                |
| `log/slog`                  | stdlib  | Structured logging               |
| `net/http`                  | stdlib  | HTTP routing (Go 1.22+ patterns) |

### Configuration with koanf

```go
import (
    "github.com/knadh/koanf/v2"
    "github.com/knadh/koanf/providers/file"
    "github.com/knadh/koanf/providers/env/v2"
    "github.com/knadh/koanf/parsers/yaml"
)

var k = koanf.New(".")

// Load with merge: file first, then env overrides
k.Load(file.Provider("config.yaml"), yaml.Parser())
k.Load(env.Provider(".", env.Opt{
    Prefix: "JELLYFIN_",
    TransformFunc: func(key, val string) (string, any) {
        return strings.ToLower(strings.TrimPrefix(key, "JELLYFIN_")), val
    },
}), nil)

// Unmarshal to struct
var cfg Config
k.Unmarshal("", &cfg)
```

### Dependency Injection with fx

```go
import "go.uber.org/fx"

func main() {
    fx.New(
        fx.Provide(
            NewConfig,
            NewLogger,
            NewDatabase,
            NewHTTPServer,
        ),
        fx.Invoke(StartServer),
    ).Run()
}

// Use fx.In for parameter structs
type ServerParams struct {
    fx.In
    Config   *Config
    Logger   *slog.Logger
    DB       *pgxpool.Pool
    Handlers []http.Handler `group:"handlers"`
}

// Use fx.Out for result structs
type ServerResult struct {
    fx.Out
    Server   *http.Server
    Handlers http.Handler `group:"handlers"`
}

// Lifecycle hooks
func NewHTTPServer(lc fx.Lifecycle, p ServerParams) *http.Server {
    srv := &http.Server{...}
    lc.Append(fx.Hook{
        OnStart: func(ctx context.Context) error {
            go srv.ListenAndServe()
            return nil
        },
        OnStop: func(ctx context.Context) error {
            return srv.Shutdown(ctx)
        },
    })
    return srv
}
```

### SQL with sqlc

```yaml
# sqlc.yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "internal/infra/database/queries/"
    schema: "migrations/"
    gen:
      go:
        package: "db"
        out: "internal/infra/database/db"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_interface: true
        emit_empty_slices: true
        json_tags_case_style: "camel"
```

```sql
-- queries/users.sql
-- name: GetUser :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users ORDER BY name;

-- name: CreateUser :one
INSERT INTO users (name, email) VALUES ($1, $2) RETURNING *;

-- name: UpdateUser :exec
UPDATE users SET name = $2 WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;
```

---

## Project Structure

```
jellyfin-go/
├── cmd/
│   └── jellyfin/
│       └── main.go           # Entry point with fx.New()
├── internal/
│   ├── api/
│   │   ├── handlers/         # HTTP handlers by domain
│   │   └── middleware/       # Auth, logging, rate limiting
│   ├── domain/               # Core business entities
│   │   ├── user.go
│   │   ├── media.go
│   │   └── library.go
│   ├── service/              # Business logic layer
│   └── infra/
│       ├── database/         # PostgreSQL + sqlc generated
│       │   ├── queries/      # SQL files for sqlc
│       │   └── db/           # Generated code
│       └── cache/            # Redis/Dragonfly
├── pkg/
│   ├── config/               # Configuration types
│   └── logger/               # slog setup
├── migrations/               # Database migrations
├── configs/
│   ├── config.yaml
│   └── defaults.yaml
└── tests/
    └── integration/
```

---

## Code Style Guidelines

### Error Handling

```go
// Always wrap errors with context
if err != nil {
    return fmt.Errorf("failed to fetch user %d: %w", id, err)
}

// Use errors.Is and errors.As for checking
if errors.Is(err, sql.ErrNoRows) {
    return ErrUserNotFound
}

// Define sentinel errors
var (
    ErrUserNotFound = errors.New("user not found")
    ErrUnauthorized = errors.New("unauthorized")
)
```

### Context Usage

```go
// All functions accept context as first parameter
func (s *UserService) GetUser(ctx context.Context, id int64) (*User, error) {
    // Check context cancellation early
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
    }

    return s.repo.GetUser(ctx, id)
}
```

### Logging with slog

```go
import "log/slog"

// Use structured logging
slog.Info("user created",
    slog.Int64("user_id", user.ID),
    slog.String("email", user.Email),
)

// Group related attributes
slog.Info("request completed",
    slog.Group("request",
        slog.String("method", r.Method),
        slog.String("path", r.URL.Path),
    ),
    slog.Group("response",
        slog.Int("status", status),
        slog.Duration("duration", duration),
    ),
)

// Error logging with stack context
slog.Error("operation failed",
    slog.String("operation", "create_user"),
    slog.Any("error", err),
)
```

### HTTP Handlers (Go 1.22+ patterns)

```go
mux := http.NewServeMux()

// Method + path pattern matching
mux.HandleFunc("GET /api/users", s.ListUsers)
mux.HandleFunc("GET /api/users/{id}", s.GetUser)
mux.HandleFunc("POST /api/users", s.CreateUser)
mux.HandleFunc("PUT /api/users/{id}", s.UpdateUser)
mux.HandleFunc("DELETE /api/users/{id}", s.DeleteUser)

// Handler implementation
func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
    if err != nil {
        http.Error(w, "invalid user id", http.StatusBadRequest)
        return
    }

    user, err := s.users.GetUser(r.Context(), id)
    if errors.Is(err, ErrUserNotFound) {
        http.Error(w, "user not found", http.StatusNotFound)
        return
    }
    if err != nil {
        slog.Error("failed to get user", slog.Any("error", err))
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}
```

### Testing

```go
// Table-driven tests
func TestGetUser(t *testing.T) {
    tests := []struct {
        name    string
        id      int64
        want    *User
        wantErr error
    }{
        {
            name: "existing user",
            id:   1,
            want: &User{ID: 1, Name: "Alice"},
        },
        {
            name:    "non-existing user",
            id:      999,
            wantErr: ErrUserNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := svc.GetUser(context.Background(), tt.id)
            if !errors.Is(err, tt.wantErr) {
                t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("GetUser() = %v, want %v", got, tt.want)
            }
        })
    }
}

// Use testing.B.Loop for benchmarks (Go 1.24)
func BenchmarkGetUser(b *testing.B) {
    svc := setupTestService(b)
    ctx := context.Background()

    for b.Loop() {
        _, _ = svc.GetUser(ctx, 1)
    }
}
```

---

## Jellyfin API Compatibility

### Goal: 100% API Compatibility

When implementing API endpoints, always reference the original C# Jellyfin:

1. **Route must match exactly** - Same path, method, query params
2. **Response structure must match** - Same JSON field names, types
3. **Behavior must match** - Same validation, defaults, errors

### Key API Patterns from C# Jellyfin

```csharp
// Original C# endpoint
[HttpGet("Users/{userId}")]
[Authorize(Policy = Policies.DefaultAuthorization)]
public ActionResult<UserDto> GetUserById([FromRoute] Guid userId)
```

```go
// Equivalent Go implementation
// Route: GET /Users/{userId}
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
    userID, err := uuid.Parse(r.PathValue("userId"))
    if err != nil {
        writeError(w, http.StatusBadRequest, "Invalid user ID")
        return
    }

    // Authorization check via middleware
    user, err := h.service.GetUser(r.Context(), userID)
    // ... response matches UserDto JSON structure
}
```

### Common Jellyfin Types

```go
// Match Jellyfin's GUID usage
type ItemID = uuid.UUID

// Match Jellyfin's nullable patterns
type NullableString struct {
    Value string
    Valid bool
}

// Match Jellyfin's date format
const JellyfinDateFormat = "2006-01-02T15:04:05.0000000Z"
```

---

## Do's and Don'ts

### DO

- ✅ Use `context.Context` as first parameter
- ✅ Use `slog` for all logging
- ✅ Use `errors.Is`/`errors.As` for error checking
- ✅ Use `%w` for error wrapping
- ✅ Use Go 1.22+ HTTP routing patterns
- ✅ Use `sqlc` for database queries
- ✅ Use `fx` for dependency injection
- ✅ Use `koanf` for configuration
- ✅ Write table-driven tests
- ✅ Use `testing.B.Loop` for benchmarks (Go 1.24)

### DON'T

- ❌ Use `init()` functions - use fx instead
- ❌ Use global variables - inject dependencies
- ❌ Use `panic` for error handling
- ❌ Use `interface{}` - use `any` (Go 1.18+)
- ❌ Use gorilla/mux - use stdlib http.ServeMux
- ❌ Use Viper - use koanf v2
- ❌ Use zap/logrus - use slog
- ❌ Use lib/pq - use pgx/v5
- ❌ Use `b.N` in benchmarks - use `b.Loop()` (Go 1.24)

---

## Git Workflow

### Branch Strategy

- `main` - Production releases only
- `develop` - Integration branch
- `feature/*` - New features
- `bugfix/*` - Bug fixes
- `release/*` - Release preparation

### Commit Convention

```
type(scope): description

Types: feat, fix, docs, style, refactor, perf, test, build, ci, chore
Scope: api, db, auth, media, config, etc.

Examples:
feat(api): add user authentication endpoints
fix(db): handle null values in media queries
docs(readme): update installation instructions
perf(media): optimize thumbnail generation
```

---

## Build & Run

```bash
# Development
go run ./cmd/jellyfin

# Build
go build -o bin/jellyfin ./cmd/jellyfin

# With version info
go build -ldflags "-X main.version=0.1.0 -X main.commit=$(git rev-parse HEAD)" ./cmd/jellyfin

# Test
go test ./...

# Generate sqlc
sqlc generate

# Lint
golangci-lint run
```

---

## References

- [Go 1.24 Release Notes](https://go.dev/doc/go1.24)
- [uber-go/fx Documentation](https://uber-go.github.io/fx/)
- [koanf v2 Documentation](https://github.com/knadh/koanf)
- [sqlc Documentation](https://docs.sqlc.dev/)
- [Jellyfin API Documentation](https://api.jellyfin.org/)
- [Original Jellyfin Source](https://github.com/jellyfin/jellyfin)

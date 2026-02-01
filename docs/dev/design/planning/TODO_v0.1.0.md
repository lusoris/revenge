# TODO v0.1.0 - Skeleton

<!-- DESIGN: planning, README, SCAFFOLD_TEMPLATE, test_output_claude -->

> Project Structure

**Status**: 🔴 Not Started
**Tag**: `v0.1.0`
**Focus**: Project Structure & Foundation

**Depends On**: [v0.0.0](TODO_v0.0.0.md) (CI/CD must work)

---

## Overview

This milestone establishes the Go project structure, dependency injection framework, configuration system, database migrations, and basic infrastructure. No business logic - just the skeleton.

---

## Deliverables

### Go Module Structure

- [ ] **Root Module** (`go.mod`)
  - [ ] Initialize module: `github.com/lusoris/revenge`
  - [ ] Add core dependencies from SOURCE_OF_TRUTH
  - [ ] Configure Go 1.25.6
  - [ ] Set up GOEXPERIMENT flags

- [ ] **Directory Structure**
  ```
  revenge/
  ├── cmd/revenge/main.go      # Entry point
  ├── internal/
  │   ├── api/                 # ogen-generated handlers
  │   ├── config/              # koanf configuration
  │   ├── content/             # Content module stubs
  │   ├── service/             # Service stubs
  │   └── infra/               # Infrastructure
  │       ├── database/        # pgx, migrations
  │       ├── cache/           # rueidis, otter
  │       ├── search/          # typesense
  │       ├── jobs/            # river
  │       └── health/          # health checks
  ├── api/openapi/             # OpenAPI specs
  ├── pkg/                     # Public packages
  └── migrations/              # SQL migrations
  ```

### fx Dependency Injection

- [ ] **Root Module** (`internal/app/module.go`)
  - [ ] fx.App configuration
  - [ ] Lifecycle hooks setup
  - [ ] Graceful shutdown handling

- [ ] **Config Module** (`internal/config/module.go`)
  - [ ] fx.Provide for config loading
  - [ ] Environment variable binding

- [ ] **Infrastructure Modules**
  - [ ] `internal/infra/database/module.go` - pgxpool provider
  - [ ] `internal/infra/cache/module.go` - rueidis provider
  - [ ] `internal/infra/search/module.go` - typesense provider
  - [ ] `internal/infra/jobs/module.go` - river provider
  - [ ] `internal/infra/health/module.go` - health provider

### Configuration System (koanf)

- [ ] **Config Struct** (`internal/config/config.go`)
  - [ ] Server configuration (port, host)
  - [ ] Database configuration (url, pool settings)
  - [ ] Cache configuration (url)
  - [ ] Search configuration (url, api_key)
  - [ ] Logging configuration (level, format)

- [ ] **Config Loading** (`internal/config/loader.go`)
  - [ ] Default values
  - [ ] YAML file loading
  - [ ] Environment variable override (REVENGE_*)
  - [ ] Validation with go-playground/validator

- [ ] **Config Files**
  - [ ] `config/config.yaml` - Default configuration
  - [ ] `config/config.example.yaml` - Example with comments

### Database Infrastructure

- [ ] **pgxpool Setup** (`internal/infra/database/pool.go`)
  - [ ] Connection string parsing
  - [ ] Pool configuration from config
  - [ ] Health check integration
  - [ ] Graceful shutdown

- [ ] **Migration Framework** (`internal/infra/database/migrate.go`)
  - [ ] golang-migrate integration
  - [ ] Embedded migrations support
  - [ ] Up/Down/Version commands

- [ ] **Initial Migrations** (`migrations/`)
  - [ ] `000001_create_schemas.up.sql` - Create public, shared, qar schemas
  - [ ] `000001_create_schemas.down.sql`
  - [ ] `000002_create_users_table.up.sql` - Basic users table
  - [ ] `000002_create_users_table.down.sql`
  - [ ] `000003_create_sessions_table.up.sql`
  - [ ] `000003_create_sessions_table.down.sql`

- [ ] **sqlc Configuration** (`sqlc.yaml`)
  - [ ] Database connection settings
  - [ ] Query file locations
  - [ ] Code generation settings

### OpenAPI Skeleton (ogen)

- [ ] **Base Spec** (`api/openapi/openapi.yaml`)
  - [ ] OpenAPI 3.1 header
  - [ ] Server definitions
  - [ ] Security schemes (Bearer JWT)
  - [ ] Common components (Error, Pagination)

- [ ] **Health Endpoints**
  - [ ] `GET /health/live` - Liveness probe
  - [ ] `GET /health/ready` - Readiness probe
  - [ ] `GET /health/startup` - Startup probe

- [ ] **ogen Generation**
  - [ ] Configure ogen.yaml
  - [ ] Generate server interfaces
  - [ ] Generate client (for testing)
  - [ ] Makefile target: `make generate`

### Logging Infrastructure

- [ ] **Development Logging** (`internal/infra/logging/tint.go`)
  - [ ] tint handler for colorized output
  - [ ] slog integration
  - [ ] Log level from config

- [ ] **Production Logging** (`internal/infra/logging/zap.go`)
  - [ ] zap JSON handler
  - [ ] Structured fields
  - [ ] Performance optimization

- [ ] **Logging Module** (`internal/infra/logging/module.go`)
  - [ ] fx provider
  - [ ] Environment-based selection

### Error Handling Patterns

- [ ] **Sentinel Errors** (`internal/errors/errors.go`)
  - [ ] ErrNotFound
  - [ ] ErrUnauthorized
  - [ ] ErrForbidden
  - [ ] ErrConflict
  - [ ] ErrValidation

- [ ] **API Errors** (`internal/api/errors.go`)
  - [ ] APIError struct
  - [ ] Error response formatting
  - [ ] Error code mapping

- [ ] **Error Wrapping** (`internal/errors/wrap.go`)
  - [ ] go-faster/errors integration
  - [ ] Stack trace preservation
  - [ ] Error unwrapping helpers

### Basic Health Endpoints

- [ ] **Health Service** (`internal/infra/health/service.go`)
  - [ ] Liveness check (always healthy if running)
  - [ ] Readiness check (all dependencies ready)
  - [ ] Startup check (initialization complete)

- [ ] **Dependency Checks** (`internal/infra/health/checks.go`)
  - [ ] PostgreSQL ping check
  - [ ] Dragonfly ping check
  - [ ] Typesense health check
  - [ ] River worker check

- [ ] **Health Handler** (`internal/infra/health/handler.go`)
  - [ ] HTTP endpoint implementation
  - [ ] JSON response format
  - [ ] Prometheus metrics integration

### Main Entry Point

- [ ] **cmd/revenge/main.go**
  - [ ] fx.New() with all modules
  - [ ] Signal handling (SIGINT, SIGTERM)
  - [ ] Version flag
  - [ ] Config path flag

- [ ] **cmd/revenge/migrate.go**
  - [ ] Subcommand: `revenge migrate up`
  - [ ] Subcommand: `revenge migrate down`
  - [ ] Subcommand: `revenge migrate version`

### Testing Infrastructure

- [ ] **Test Helpers** (`internal/testutil/`)
  - [ ] `database.go` - embedded-postgres setup
  - [ ] `fixtures.go` - Common test fixtures
  - [ ] `assertions.go` - Custom assertions

- [ ] **Integration Test Setup** (`internal/testutil/containers.go`)
  - [ ] testcontainers-go PostgreSQL
  - [ ] testcontainers-go Dragonfly
  - [ ] testcontainers-go Typesense

- [ ] **Test Configuration**
  - [ ] `config/config.test.yaml`
  - [ ] CI-specific test settings

### Makefile

- [ ] **Build Targets**
  - [ ] `make build` - Build binary
  - [ ] `make test` - Run tests
  - [ ] `make test-integration` - Run integration tests
  - [ ] `make generate` - Generate code (ogen, sqlc, mockery)
  - [ ] `make lint` - Run linters
  - [ ] `make migrate-up` - Run migrations
  - [ ] `make migrate-down` - Rollback migrations

- [ ] **Development Targets**
  - [ ] `make dev` - Start with hot reload (air)
  - [ ] `make docker-up` - Start Docker Compose stack
  - [ ] `make docker-down` - Stop Docker Compose stack

### Air Configuration

- [ ] **.air.toml**
  - [ ] Watch directories
  - [ ] Exclude patterns
  - [ ] Build command with GOEXPERIMENT
  - [ ] Binary output path

---

## Verification Checklist

- [ ] `go build ./...` succeeds
- [ ] `make test` passes
- [ ] `make lint` passes (when golangci-lint supports Go 1.25)
- [ ] Health endpoints respond correctly
- [ ] Migrations run successfully
- [ ] Docker Compose stack starts
- [ ] CI pipeline passes

---

## Dependencies from SOURCE_OF_TRUTH

| Package | Version | Purpose |
|---------|---------|---------|
| go.uber.org/fx | v1.23.0 | Dependency injection |
| github.com/jackc/pgx/v5 | v5.7.5 | PostgreSQL driver |
| github.com/knadh/koanf/v2 | v2.3.0 | Configuration |
| github.com/ogen-go/ogen | v1.18.0 | OpenAPI codegen |
| github.com/golang-migrate/migrate/v4 | v4.19.1 | Migrations |
| github.com/lmittmann/tint | v1.1.2 | Dev logging |
| go.uber.org/zap | v1.28.0 | Prod logging |
| github.com/go-faster/errors | v0.7.1 | Error handling |
| github.com/stretchr/testify | v1.11.1 | Testing |
| github.com/fergusstrange/embedded-postgres | v1.30.0 | Test database |

---

## Related Documentation

- [ROADMAP.md](ROADMAP.md) - Full roadmap overview
- [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) - Authoritative versions
- [TECH_STACK.md](../technical/TECH_STACK.md) - Technology rationale
- [01_ARCHITECTURE.md](../architecture/01_ARCHITECTURE.md) - System architecture

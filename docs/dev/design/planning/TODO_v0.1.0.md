# TODO v0.1.0 - Skeleton

<!-- DESIGN: planning, README, test_output_claude, test_output_wiki -->


<!-- TOC-START -->

## Table of Contents

- [Overview](#overview)
- [Deliverables](#deliverables)
  - [Go Module Structure](#go-module-structure)
  - [fx Dependency Injection](#fx-dependency-injection)
  - [Configuration System (koanf)](#configuration-system-koanf)
  - [Database Infrastructure](#database-infrastructure)
  - [OpenAPI Skeleton (ogen)](#openapi-skeleton-ogen)
  - [Logging Infrastructure](#logging-infrastructure)
  - [Error Handling Patterns](#error-handling-patterns)
  - [Basic Health Endpoints](#basic-health-endpoints)
  - [Main Entry Point](#main-entry-point)
  - [Testing Infrastructure](#testing-infrastructure)
  - [Makefile](#makefile)
  - [Air Configuration](#air-configuration)
- [Verification Checklist](#verification-checklist)
- [Dependencies](#dependencies)
- [Related Documentation](#related-documentation)

<!-- TOC-END -->


> Project Structure

**Status**: ✅ Complete
**Tag**: `v0.1.0`
**Completed**: 2026-02-02
**Focus**: Project Structure & Foundation

**Depends On**: [v0.0.0](TODO_v0.0.0.md) (CI/CD must work)

### Patch Releases

| Version | Focus | Status |
|---------|-------|--------|
| v0.1.0 | Skeleton - Project Structure | ✅ Complete |
| v0.1.1 | Test Coverage Sprint (Database 78%) | ✅ Complete |
| v0.1.2 | Errors Package 100% Coverage | ✅ Complete |
| v0.1.3 | CI Fixes (Port Conflicts, Lint, macOS/Windows) | ✅ Complete |

---

## Overview

This milestone establishes the Go project structure, dependency injection framework, configuration system, database migrations, and basic infrastructure. No business logic - just the skeleton.

---

## Deliverables

### Go Module Structure

- [x] **Root Module** (`go.mod`)
  - [x] Initialize module: `github.com/lusoris/revenge`
  - [x] Add core dependencies from SOURCE_OF_TRUTH
  - [x] Configure Go 1.25.6
  - [x] Set up GOEXPERIMENT flags

- [x] **Directory Structure**
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

- [x] **Root Module** (`internal/app/module.go`)
  - [x] fx.App configuration
  - [x] Lifecycle hooks setup
  - [x] Graceful shutdown handling

- [x] **Config Module** (`internal/config/module.go`)
  - [x] fx.Provide for config loading
  - [x] Environment variable binding

- [x] **Infrastructure Modules**
  - [x] `internal/infra/database/module.go` - pgxpool provider
  - [ ] `internal/infra/cache/module.go` - rueidis provider (deferred to v0.2.0)
  - [ ] `internal/infra/search/module.go` - typesense provider (deferred to v0.2.0)
  - [ ] `internal/infra/jobs/module.go` - river provider (deferred to v0.2.0)
  - [x] `internal/infra/health/module.go` - health provider

### Configuration System (koanf)

- [x] **Config Struct** (`internal/config/config.go`)
  - [x] Server configuration (port, host)
  - [x] Database configuration (url, pool settings)
  - [x] Cache configuration (url)
  - [x] Search configuration (url, api_key)
  - [x] Logging configuration (level, format)

- [x] **Config Loading** (`internal/config/loader.go`)
  - [x] Default values
  - [x] YAML file loading
  - [x] Environment variable override (REVENGE_*)
  - [x] Validation with go-playground/validator

- [x] **Config Files**
  - [x] `config/config.yaml` - Default configuration
  - [x] `config/config.example.yaml` - Example with comments

### Database Infrastructure

- [x] **pgxpool Setup** (`internal/infra/database/pool.go`)
  - [x] Connection string parsing
  - [x] Pool configuration from config
  - [x] Health check integration
  - [x] Graceful shutdown

- [x] **Migration Framework** (`internal/infra/database/migrate.go`)
  - [x] golang-migrate integration
  - [x] Embedded migrations support
  - [x] Up/Down/Version commands

- [x] **Initial Migrations** (`migrations/`)
  - [x] `000001_create_schemas.up.sql` - Create public, shared, qar schemas
  - [x] `000001_create_schemas.down.sql`
  - [x] `000002_create_users_table.up.sql` - Basic users table
  - [x] `000002_create_users_table.down.sql`
  - [x] `000003_create_sessions_table.up.sql`
  - [x] `000003_create_sessions_table.down.sql`

- [x] **sqlc Configuration** (`sqlc.yaml`)
  - [x] Database connection settings
  - [x] Query file locations
  - [x] Code generation settings

### OpenAPI Skeleton (ogen)

- [x] **Base Spec** (`api/openapi/openapi.yaml`)
  - [x] OpenAPI 3.1 header
  - [x] Server definitions
  - [x] Security schemes (Bearer JWT)
  - [x] Common components (Error, Pagination)

- [x] **Health Endpoints**
  - [x] `GET /health/live` - Liveness probe
  - [x] `GET /health/ready` - Readiness probe
  - [x] `GET /health/startup` - Startup probe

- [x] **ogen Generation**
  - [x] Configure ogen.yaml
  - [x] Generate server interfaces
  - [x] Generate client (for testing)
  - [x] Makefile target: `make generate`

### Logging Infrastructure

- [x] **Development Logging** (`internal/infra/logging/logger.go`)
  - [x] tint handler for colorized output
  - [x] slog integration
  - [x] Log level from config

- [x] **Production Logging** (`internal/infra/logging/logger.go`)
  - [x] JSON handler for production
  - [x] Structured fields
  - [x] Performance optimization

- [x] **Logging Module** (`internal/infra/logging/module.go`)
  - [x] fx provider
  - [x] Environment-based selection

### Error Handling Patterns

- [x] **Sentinel Errors** (`internal/errors/errors.go`)
  - [x] ErrNotFound
  - [x] ErrUnauthorized
  - [x] ErrForbidden
  - [x] ErrConflict
  - [x] ErrValidation
  - [x] ErrInternal, ErrBadRequest, ErrUnavailable, ErrTimeout

- [x] **API Errors** (`internal/api/handler.go`)
  - [x] APIError struct
  - [x] Error response formatting
  - [x] Error code mapping

- [x] **Error Wrapping** (`internal/errors/wrap.go`)
  - [x] go-faster/errors integration
  - [x] Stack trace preservation
  - [x] Error unwrapping helpers

### Basic Health Endpoints

- [x] **Health Service** (`internal/infra/health/service.go`)
  - [x] Liveness check (always healthy if running)
  - [x] Readiness check (all dependencies ready)
  - [x] Startup check (initialization complete)
  - [x] Full check (detailed dependency status)

- [x] **Dependency Checks** (`internal/infra/health/service.go`)
  - [x] PostgreSQL ping check
  - [ ] Dragonfly ping check (deferred to v0.2.0)
  - [ ] Typesense health check (deferred to v0.2.0)
  - [ ] River worker check (deferred to v0.2.0)

- [x] **Health Handler** (`internal/infra/health/handler.go`)
  - [x] HTTP endpoint implementation
  - [x] JSON response format
  - [ ] Prometheus metrics integration (deferred to v0.2.0)

### Main Entry Point

- [x] **cmd/revenge/main.go**
  - [x] fx.New() with all modules
  - [x] Signal handling (SIGINT, SIGTERM)
  - [x] Version flag
  - [x] Config path flag

- [ ] **cmd/revenge/migrate.go** (deferred - using golang-migrate CLI)
  - [ ] Subcommand: `revenge migrate up`
  - [ ] Subcommand: `revenge migrate down`
  - [ ] Subcommand: `revenge migrate version`

### Testing Infrastructure

- [x] **Test Helpers** (`internal/testutil/`)
  - [x] `containers.go` - testcontainers PostgreSQL setup
  - [x] `fixtures.go` - Common test fixtures
  - [x] Custom assertions via testify

- [x] **Integration Test Setup** (`internal/testutil/containers.go`)
  - [x] testcontainers-go PostgreSQL
  - [ ] testcontainers-go Dragonfly (deferred to v0.2.0)
  - [ ] testcontainers-go Typesense (deferred to v0.2.0)

- [x] **Test Configuration**
  - [x] `config/config.test.yaml`
  - [x] CI-specific test settings

### Makefile

- [x] **Build Targets**
  - [x] `make build` - Build binary
  - [x] `make test` - Run tests
  - [x] `make test-integration` - Run integration tests
  - [x] `make generate` - Generate code (ogen, sqlc, mockery)
  - [x] `make lint` - Run linters
  - [x] `make migrate-up` - Run migrations
  - [x] `make migrate-down` - Rollback migrations

- [x] **Development Targets**
  - [x] `make dev` - Start with hot reload (air)
  - [x] `make docker-up` - Start Docker Compose stack
  - [x] `make docker-down` - Stop Docker Compose stack

### Air Configuration

- [x] **.air.toml**
  - [x] Watch directories
  - [x] Exclude patterns
  - [x] Build command with GOEXPERIMENT
  - [x] Binary output path

---

## Verification Checklist

- [x] `go build ./...` succeeds
- [x] `make test` passes
- [x] `make lint` passes (golangci-lint v2)
- [x] Health endpoints respond correctly
- [x] Migrations run successfully
- [x] Docker Compose stack starts
- [x] CI pipeline passes (all platforms: Ubuntu, macOS, Windows)

---

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| go.uber.org/fx | v1.24.0 | Dependency injection |
| github.com/jackc/pgx/v5 | v5.8.0 | PostgreSQL driver |
| github.com/knadh/koanf/v2 | v2.3.0 | Configuration |
| github.com/ogen-go/ogen | v1.18.0 | OpenAPI codegen |
| github.com/golang-migrate/migrate/v4 | v4.19.1 | Migrations |
| github.com/lmittmann/tint | v1.1.2 | Dev logging |
| go.uber.org/zap | v1.27.1 | Prod logging |
| github.com/go-faster/errors | v0.7.1 | Error handling |
| github.com/stretchr/testify | v1.11.1 | Testing |
| github.com/fergusstrange/embedded-postgres | v1.30.0 | Test database |

---

## Design Documentation

> **Note**: Design work for v0.1.0 scope is **COMPLETE**. The following design documents exist and should be referenced during implementation:

### Architecture & Foundation
- [ARCHITECTURE.md](../architecture/ARCHITECTURE.md) - Complete system architecture with all layers and components
- [DESIGN_PRINCIPLES.md](../architecture/DESIGN_PRINCIPLES.md) - Design principles guiding all implementation
- [TECH_STACK.md](../technical/TECH_STACK.md) - Full technology stack with rationale

### Integration Designs (Foundation)
- [POSTGRESQL.md](../integrations/infrastructure/POSTGRESQL.md) - Database patterns, migrations, pgx usage
- [DRAGONFLY.md](../integrations/infrastructure/DRAGONFLY.md) - Cache architecture (rueidis, otter)
- [RIVER.md](../integrations/infrastructure/RIVER.md) - Job queue setup with fx
- [TYPESENSE.md](../integrations/infrastructure/TYPESENSE.md) - Search infrastructure

### Technical Patterns
- [API.md](../technical/API.md) - OpenAPI design patterns with ogen
- [CONFIGURATION.md](../technical/CONFIGURATION.md) - Configuration system with koanf
- [OBSERVABILITY.md](../technical/OBSERVABILITY.md) - Logging, metrics, tracing setup
- [TESTING.md](../technical/TESTING.md) - Testing patterns and infrastructure

### Operations
- [DEVELOPMENT.md](../operations/DEVELOPMENT.md) - Development environment setup
- [CODING_STANDARDS.md](../operations/CODING_STANDARDS.md) - Go best practices

---

## Related Documentation

- [ROADMAP.md](ROADMAP.md) - Full roadmap overview
- [DESIGN_INDEX.md](../DESIGN_INDEX.md) - Full design documentation index

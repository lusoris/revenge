## Table of Contents

- [Revenge - Technology Stack](#revenge-technology-stack)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Components](#components)
  - [Implementation](#implementation)
    - [File Structure](#file-structure)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
    - [Config Keys](#config-keys)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)



---
sources:
  - name: Dragonfly Documentation
    url: ../../sources/infrastructure/dragonfly.md
    note: L2 cache backend
  - name: Uber fx
    url: ../../sources/tooling/fx.md
    note: Dependency injection framework
  - name: Go log/slog
    url: ../../sources/go/stdlib/slog.md
    note: Structured logging (dev)
  - name: gohlslib (HLS)
    url: ../../sources/media/gohlslib.md
    note: HLS streaming support
  - name: koanf
    url: ../../sources/tooling/koanf.md
    note: Configuration management
  - name: ogen OpenAPI Generator
    url: ../../sources/tooling/ogen.md
    note: Type-safe API code generation
  - name: ogen Documentation
    url: ../../sources/tooling/ogen-guide.md
    note: Official ogen docs
  - name: pgx PostgreSQL Driver
    url: ../../sources/database/pgx.md
    note: PostgreSQL native driver
  - name: PostgreSQL Arrays
    url: ../../sources/database/postgresql-arrays.md
    note: Array type support
  - name: PostgreSQL JSON Functions
    url: ../../sources/database/postgresql-json.md
    note: JSON/JSONB operations
  - name: River Job Queue
    url: ../../sources/tooling/river.md
    note: PostgreSQL-backed jobs
  - name: River Documentation
    url: ../../sources/tooling/river-guide.md
    note: Official River docs
  - name: rueidis
    url: ../../sources/tooling/rueidis.md
    note: Redis/Dragonfly client
  - name: rueidis GitHub README
    url: ../../sources/tooling/rueidis-guide.md
    note: Client documentation
  - name: shadcn-svelte
    url: ../../sources/frontend/shadcn-svelte.md
    note: UI component library
  - name: sqlc
    url: ../../sources/database/sqlc.md
    note: SQL code generator
  - name: sqlc Configuration
    url: ../../sources/database/sqlc-config.md
    note: sqlc.yaml reference
  - name: Svelte 5 Runes
    url: ../../sources/frontend/svelte-runes.md
    note: Runes-based reactivity
  - name: Svelte 5 Documentation
    url: ../../sources/frontend/svelte5.md
    note: Svelte 5 API reference
  - name: SvelteKit Documentation
    url: ../../sources/frontend/sveltekit.md
    note: SvelteKit framework
  - name: TanStack Query
    url: ../../sources/frontend/tanstack-query.md
    note: Server state management
  - name: Typesense API
    url: ../../sources/infrastructure/typesense.md
    note: Search engine API
  - name: Typesense Go Client
    url: ../../sources/infrastructure/typesense-go.md
    note: Go client library
  - name: otter Cache
    url: https://pkg.go.dev/github.com/maypok86/otter
    note: In-memory L1 cache
  - name: sturdyc
    url: ../../sources/tooling/sturdyc.md
    note: Request coalescing
  - name: zap Logger
    url: ../../sources/tooling/zap.md
    note: Production logging
  - name: tint Logger
    url: ../../sources/tooling/tint.md
    note: Development logging
  - name: golang-migrate
    url: https://pkg.go.dev/github.com/golang-migrate/migrate/v4
    note: Database migrations
  - name: testify
    url: ../../sources/testing/testify.md
    note: Testing framework
  - name: mockery
    url: ../../sources/testing/mockery-guide.md
    note: Mock generation
  - name: testcontainers-go
    url: ../../sources/testing/testcontainers.md
    note: Integration testing
  - name: golangci-lint
    url: ../../sources/go_dev_tools/golangci-lint/main.md
    note: Go linting
  - name: markdownlint-cli2
    url: https://github.com/DavidAnson/markdownlint-cli2
    note: Markdown linting
design_refs:
  - title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
  - title: 00_SOURCE_OF_TRUTH
    path: 00_SOURCE_OF_TRUTH.md
---

# Revenge - Technology Stack


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: technical


> > Complete technology stack powering Revenge

Stack overview:
- **Backend**: Go 1.25.6 with GOEXPERIMENT=greenteagc,jsonv2
- **DI Framework**: fx v1.24.0 for dependency injection and lifecycle management
- **API Layer**: ogen v1.18.0 for type-safe OpenAPI code generation
- **Database**: PostgreSQL 18+ with pgx v5.8.0 driver, sqlc v1.30.0 code generation
- **Caching**: L1 (otter v2.3.1 in-memory), L2 (Dragonfly v1.26.1 via rueidis v1.0.54)
- **Search**: Typesense 27.1 for full-text search with typo tolerance
- **Jobs**: River v0.20.1 for PostgreSQL-backed job queue
- **Logging**: slog/tint v1.1.2 (dev), zap v1.27.1 (prod)
- **Frontend**: SvelteKit 2, Svelte 5, Tailwind CSS 4, shadcn-svelte
- **Testing**: testify v1.11.1, mockery v3.3.0, testcontainers v0.40.0


---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | 🟢 | Complete stack documented |
| Sources | ✅ | All sources documented |
| Instructions | 🟢 | Implementation patterns included |
| Code | 🟢 | Stack implemented in codebase |
| Linting | ✅ | golangci-lint v2.8.0 |
| Unit Testing | 🟢 | 80%+ coverage target |
| Integration Testing | 🟢 | testcontainers in use |

**Overall**: ✅ Complete



---


## Architecture

┌─────────────────────────────────────────────────────────────────┐
│                    SvelteKit Frontend                           │
│          (Svelte 5 + Tailwind 4 + shadcn-svelte)               │
└─────────────────────┬───────────────────────────────────────────┘
                      │ HTTP/JSON
┌─────────────────────▼───────────────────────────────────────────┐
│                  ogen OpenAPI Handlers                          │
│             (Type-safe from openapi.yaml)                       │
└─────────────────────┬───────────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────────┐
│                   Service Layer (fx)                            │
│         Movies │ TV │ Music │ QAR │ User │ Auth                 │
└─────┬───────┬──┴──┬────────┬──────┬──────────────────────┬─────┘
      │       │     │        │      │                      │
┌─────▼───┐ ┌─▼─────▼──┐ ┌──▼──────▼───┐          ┌───────▼──────┐
│  Cache  │ │  Search  │ │   Jobs      │          │  Repository  │
│  (L1+L2)│ │(Typesense)│ │   (River)   │          │   (sqlc)     │
└─────────┘ └──────────┘ └─────────────┘          └───────┬──────┘
                                                            │
                                                  ┌─────────▼──────┐
                                                  │  PostgreSQL 18 │
                                                  └────────────────┘


### Components

<!-- Component description -->


## Implementation

### File Structure

<!-- File structure -->

### Key Interfaces

<!-- Interface definitions -->

### Dependencies

<!-- Dependency list -->





## Configuration
### Environment Variables

[{'name': 'REVENGE_SERVER_HOST', 'type': 'string', 'default': '0.0.0.0', 'description': 'Server bind address'}, {'name': 'REVENGE_SERVER_PORT', 'type': 'int', 'default': 8080, 'description': 'Server listen port'}, {'name': 'REVENGE_DATABASE_URL', 'type': 'string', 'required': True, 'description': 'PostgreSQL connection string (postgres://user:pass@host:5432/db)'}, {'name': 'REVENGE_CACHE_REDIS_ADDR', 'type': 'string', 'default': 'localhost:6379', 'description': 'Dragonfly/Redis address for L2 cache'}, {'name': 'REVENGE_SEARCH_HOST', 'type': 'string', 'default': 'localhost:8108', 'description': 'Typesense server host'}, {'name': 'REVENGE_JOBS_ENABLED', 'type': 'bool', 'default': True, 'description': 'Enable River job queue'}]

### Config Keys

<!-- Configuration keys -->




## Testing Strategy

### Unit Tests

<!-- Unit test strategy -->

### Integration Tests

<!-- Integration test strategy -->

### Test Coverage

Target: **80% minimum**







## Related Documentation
### Design Documents
- [01_ARCHITECTURE](architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](architecture/03_METADATA_SYSTEM.md)
- [00_SOURCE_OF_TRUTH](00_SOURCE_OF_TRUTH.md)

### External Sources
- [Dragonfly Documentation](../../sources/infrastructure/dragonfly.md) - L2 cache backend
- [Uber fx](../../sources/tooling/fx.md) - Dependency injection framework
- [Go log/slog](../../sources/go/stdlib/slog.md) - Structured logging (dev)
- [gohlslib (HLS)](../../sources/media/gohlslib.md) - HLS streaming support
- [koanf](../../sources/tooling/koanf.md) - Configuration management
- [ogen OpenAPI Generator](../../sources/tooling/ogen.md) - Type-safe API code generation
- [ogen Documentation](../../sources/tooling/ogen-guide.md) - Official ogen docs
- [pgx PostgreSQL Driver](../../sources/database/pgx.md) - PostgreSQL native driver
- [PostgreSQL Arrays](../../sources/database/postgresql-arrays.md) - Array type support
- [PostgreSQL JSON Functions](../../sources/database/postgresql-json.md) - JSON/JSONB operations
- [River Job Queue](../../sources/tooling/river.md) - PostgreSQL-backed jobs
- [River Documentation](../../sources/tooling/river-guide.md) - Official River docs
- [rueidis](../../sources/tooling/rueidis.md) - Redis/Dragonfly client
- [rueidis GitHub README](../../sources/tooling/rueidis-guide.md) - Client documentation
- [shadcn-svelte](../../sources/frontend/shadcn-svelte.md) - UI component library
- [sqlc](../../sources/database/sqlc.md) - SQL code generator
- [sqlc Configuration](../../sources/database/sqlc-config.md) - sqlc.yaml reference
- [Svelte 5 Runes](../../sources/frontend/svelte-runes.md) - Runes-based reactivity
- [Svelte 5 Documentation](../../sources/frontend/svelte5.md) - Svelte 5 API reference
- [SvelteKit Documentation](../../sources/frontend/sveltekit.md) - SvelteKit framework
- [TanStack Query](../../sources/frontend/tanstack-query.md) - Server state management
- [Typesense API](../../sources/infrastructure/typesense.md) - Search engine API
- [Typesense Go Client](../../sources/infrastructure/typesense-go.md) - Go client library
- [otter Cache](https://pkg.go.dev/github.com/maypok86/otter) - In-memory L1 cache
- [sturdyc](../../sources/tooling/sturdyc.md) - Request coalescing
- [zap Logger](../../sources/tooling/zap.md) - Production logging
- [tint Logger](../../sources/tooling/tint.md) - Development logging
- [golang-migrate](https://pkg.go.dev/github.com/golang-migrate/migrate/v4) - Database migrations
- [testify](../../sources/testing/testify.md) - Testing framework
- [mockery](../../sources/testing/mockery-guide.md) - Mock generation
- [testcontainers-go](../../sources/testing/testcontainers.md) - Integration testing
- [golangci-lint](../../sources/go_dev_tools/golangci-lint/main.md) - Go linting
- [markdownlint-cli2](https://github.com/DavidAnson/markdownlint-cli2) - Markdown linting


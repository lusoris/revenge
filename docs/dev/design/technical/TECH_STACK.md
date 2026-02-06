## Table of Contents

- [Revenge - Technology Stack](#revenge-technology-stack)
  - [Status](#status)
  - [Architecture](#architecture)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# Revenge - Technology Stack

<!-- DESIGN: technical, README, test_output_claude, test_output_wiki -->


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
| Design | ✅ | Complete stack documented |
| Sources | ✅ | All sources documented |
| Instructions | ✅ | Implementation patterns included |
| Code | ✅ | Stack implemented in codebase |
| Linting | ✅ | golangci-lint v2.8.0 |
| Unit Testing | ✅ | 80%+ coverage target |
| Integration Testing | ✅ | testcontainers in use |

**Overall**: ✅ Complete


---


## Architecture

```mermaid
flowchart LR
    subgraph Layer1["Layer 1"]
        node1["SvelteKit Frontend<br/>(Svelte 5 + Tailwind 4 + shadcn-svelte)"]
    end

    subgraph Layer2["Layer 2"]
        node2[["ogen OpenAPI Handlers<br/>(Type-safe from openapi.yaml)"]]
    end

    subgraph Service["Service Layer (fx)"]
        node3[" "]
    end

    subgraph Layer4["Layer 4"]
        node4[("Cache<br/>(L1+L2)")]
        node5["Search<br/>(Typesense)"]
        node6["Jobs<br/>(River)"]
        node7["Repository<br/>(sqlc)"]
    end

    subgraph Layer5["Layer 5"]
        node8[("PostgreSQL 18")]
    end

    %% Connections
    node1 --> node2
    node2 --> node3
    node3 --> node4
    node7 --> node8

    %% Styling
    style Layer1 fill:#1976D2,stroke:#1976D2,color:#fff
    style Layer2 fill:#388E3C,stroke:#388E3C,color:#fff
    style Service fill:#7B1FA2,stroke:#7B1FA2,color:#fff
    style Layer4 fill:#F57C00,stroke:#F57C00,color:#fff
    style Layer5 fill:#C2185B,stroke:#C2185B,color:#fff
```
## Configuration

### Environment Variables

[{'name': 'REVENGE_SERVER_HOST', 'type': 'string', 'default': '0.0.0.0', 'description': 'Server bind address'}, {'name': 'REVENGE_SERVER_PORT', 'type': 'int', 'default': 8080, 'description': 'Server listen port'}, {'name': 'REVENGE_DATABASE_URL', 'type': 'string', 'required': True, 'description': 'PostgreSQL connection string (postgres://user:pass@host:5432/db)'}, {'name': 'REVENGE_CACHE_REDIS_ADDR', 'type': 'string', 'default': 'localhost:6379', 'description': 'Dragonfly/Redis address for L2 cache'}, {'name': 'REVENGE_SEARCH_HOST', 'type': 'string', 'default': 'localhost:8108', 'description': 'Typesense server host'}, {'name': 'REVENGE_JOBS_ENABLED', 'type': 'bool', 'default': True, 'description': 'Enable River job queue'}]
## Related Documentation
### Design Documents
- [01_ARCHITECTURE](../architecture/ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../architecture/DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../architecture/METADATA_SYSTEM.md)
- [00_SOURCE_OF_TRUTH](../00_SOURCE_OF_TRUTH.md)

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


## Table of Contents

- [Revenge - Technology Stack](#revenge-technology-stack)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [See Also](#see-also)



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
    path: ../architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: ../architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: ../architecture/03_METADATA_SYSTEM.md
  - title: 00_SOURCE_OF_TRUTH
    path: 00_SOURCE_OF_TRUTH.md
---

# Revenge - Technology Stack




> Modern, performant stack for self-hosted media serving

Revenge uses a carefully selected stack optimized for performance and maintainability. The Go backend leverages fx for dependency injection, ogen for type-safe API generation from OpenAPI specs, and pgx for high-performance PostgreSQL access. Multi-tier caching with otter (L1) and Dragonfly (L2) ensures low latency. Typesense provides instant full-text search. River handles background jobs using PostgreSQL for ACID guarantees. The SvelteKit frontend with Svelte 5 runes delivers reactive UIs with minimal overhead.

---




## Contents

<!-- TOC will be auto-generated here by markdown-toc -->

---


## How It Works

<!-- User-friendly explanation -->




## Features
<!-- Feature list placeholder -->



## Configuration
<!-- User-friendly configuration guide -->









## Related Documentation
### See Also
<!-- Related wiki pages -->



---

**Need Help?** [Open an issue](https://github.com/revenge-project/revenge/issues) or [Join the discussion](https://github.com/revenge-project/revenge/discussions)
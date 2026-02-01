## Table of Contents

- [Revenge - Technology Stack](#revenge-technology-stack)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [Related Pages](#related-pages)
    - [Learn More](#learn-more)

# Revenge - Technology Stack




> Modern, performant stack for self-hosted media serving

Revenge uses a carefully selected stack optimized for performance and maintainability. The Go backend leverages fx for dependency injection, ogen for type-safe API generation from OpenAPI specs, and pgx for high-performance PostgreSQL access. Multi-tier caching with otter (L1) and Dragonfly (L2) ensures low latency. Typesense provides instant full-text search. River handles background jobs using PostgreSQL for ACID guarantees. The SvelteKit frontend with Svelte 5 runes delivers reactive UIs with minimal overhead.

---





---






## Features
<!-- Feature list placeholder -->




## Configuration









## Related Documentation
### Related Pages
<!-- Related wiki pages -->

### Learn More

Official documentation and guides:
- [Dragonfly Documentation](../../sources/infrastructure/dragonfly.md)
- [Uber fx](../../sources/tooling/fx.md)
- [Go log/slog](../../sources/go/stdlib/slog.md)
- [gohlslib (HLS)](../../sources/media/gohlslib.md)
- [koanf](../../sources/tooling/koanf.md)
- [ogen OpenAPI Generator](../../sources/tooling/ogen.md)
- [ogen Documentation](../../sources/tooling/ogen-guide.md)
- [pgx PostgreSQL Driver](../../sources/database/pgx.md)
- [PostgreSQL Arrays](../../sources/database/postgresql-arrays.md)
- [PostgreSQL JSON Functions](../../sources/database/postgresql-json.md)
- [River Job Queue](../../sources/tooling/river.md)
- [River Documentation](../../sources/tooling/river-guide.md)
- [rueidis](../../sources/tooling/rueidis.md)
- [rueidis GitHub README](../../sources/tooling/rueidis-guide.md)
- [shadcn-svelte](../../sources/frontend/shadcn-svelte.md)
- [sqlc](../../sources/database/sqlc.md)
- [sqlc Configuration](../../sources/database/sqlc-config.md)
- [Svelte 5 Runes](../../sources/frontend/svelte-runes.md)
- [Svelte 5 Documentation](../../sources/frontend/svelte5.md)
- [SvelteKit Documentation](../../sources/frontend/sveltekit.md)
- [TanStack Query](../../sources/frontend/tanstack-query.md)
- [Typesense API](../../sources/infrastructure/typesense.md)
- [Typesense Go Client](../../sources/infrastructure/typesense-go.md)
- [otter Cache](https://pkg.go.dev/github.com/maypok86/otter)
- [sturdyc](../../sources/tooling/sturdyc.md)
- [zap Logger](../../sources/tooling/zap.md)
- [tint Logger](../../sources/tooling/tint.md)
- [golang-migrate](https://pkg.go.dev/github.com/golang-migrate/migrate/v4)
- [testify](../../sources/testing/testify.md)
- [mockery](../../sources/testing/mockery-guide.md)
- [testcontainers-go](../../sources/testing/testcontainers.md)
- [golangci-lint](../../sources/go_dev_tools/golangci-lint/main.md)
- [markdownlint-cli2](https://github.com/DavidAnson/markdownlint-cli2)



---

**Need Help?** [Open an issue](https://github.com/revenge-project/revenge/issues) or [Join the discussion](https://github.com/revenge-project/revenge/discussions)
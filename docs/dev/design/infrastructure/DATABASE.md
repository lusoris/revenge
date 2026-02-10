# Database Infrastructure

<!-- DESIGN: infrastructure -->

**Package**: `internal/infra/database`
**fx Module**: `database.Module`

> PostgreSQL connection pooling, migrations, query generation, logging, and metrics

---

## Service Structure

```
internal/infra/database/
├── pool.go                # pgxpool config, health check, stats
├── module.go              # fx module (pool, queries, lifecycle hooks)
├── migrate.go             # golang-migrate with embedded FS (iofs)
├── logger.go              # QueryLogger (slog) with slow query detection
├── testing.go             # Embedded PostgreSQL for tests
├── queries/               # Hand-written SQL source (one dir per schema)
│   ├── shared/            # Auth, sessions, settings, activity, library, MFA, OIDC, API keys
│   ├── movie/             # Movie metadata queries
│   ├── tvshow/            # TV show series/season/episode queries
│   ├── qar/               # Adult content queries (placeholder)
│   └── <module>/          # Each new content module gets its own directory
└── migrations/shared/     # Versioned migrations (all schemas in one sequence)

# sqlc-generated code lives in each content module's db/ package:
internal/infra/database/db/        # Shared queries (auth, sessions, settings, etc.)
internal/content/<module>/db/      # Module-specific queries (one per content schema)
```

## Connection Pool

```go
// pool.go
func PoolConfig(cfg *config.Config) (*pgxpool.Config, error)
func NewPool(ctx context.Context, cfg *pgxpool.Config) (*pgxpool.Pool, error)
func Health(ctx context.Context, pool *pgxpool.Pool) error
func Stats(pool *pgxpool.Pool) map[string]interface{}
```

**Pool defaults**:
- MaxConns: `CPU * 2 + 1` (configurable)
- MinConns: Optional minimum idle connections
- Health check: `SELECT 1` with 2-second timeout
- Stats: acquire count/duration, active/idle/total conns, error counts

## Migrations

36 versioned migrations embedded via `//go:embed migrations/shared/*.sql`:

| Range | Coverage |
|-------|----------|
| 001-002 | Schemas (shared, public, qar) + users table |
| 003-013 | Sessions, settings, avatars, auth tokens, password reset, email verification, casbin rules, API keys, OIDC |
| 014-015 | Activity logs, libraries + permissions + scans |
| 016-020 | MFA (TOTP, WebAuthn, backup codes, MFA settings, session MFA tracking) |
| 021-026 | Movies (table, files, credits, collections, genres, watched) |
| 027-032 | Moderator role, TOTP cleanup, fine-grained permissions, failed logins, multi-language movies, TV shows |
| 033-035 | Metadata language, movie soft delete, external ratings |
| 036 | Create `movie` schema (migrate movie tables from `public`) |

**Schema-per-module model**: each content module owns a dedicated PostgreSQL schema.
All schemas live in the same database to allow cross-schema JOINs (e.g. `user_id` FKs into `shared.users`) and share a single connection pool.

| Schema | Module | Purpose |
|--------|--------|---------|
| `shared` | infra | Auth, sessions, settings, RBAC, API keys, OIDC, MFA, libraries, activity |
| `movie` | content/movie | Movies, files, credits, collections, genres, watch progress |
| `tvshow` | content/tvshow | Series, seasons, episodes, files, credits, genres, networks, watch progress |
| `qar` | content/qar | Adult content (placeholder) |
| *`<name>`* | *content/`<name>`* | *New content modules follow the same convention* |

**Convention for adding a new content module**:
1. Create a migration: `CREATE SCHEMA IF NOT EXISTS <name>;` + tables
2. Add SQL queries in `queries/<name>/`
3. Add a sqlc target outputting to `internal/content/<name>/db/`
4. Add repository, service, and handler layers in `internal/content/<name>/`

```go
func MigrateUp(databaseURL string, logger *slog.Logger) error
func MigrateDown(databaseURL string, logger *slog.Logger) error
func MigrateVersion(databaseURL string) (version uint, dirty bool, err error)
func MigrateTo(databaseURL string, version uint, logger *slog.Logger) error
```

## Query Generation (sqlc)

**sqlc v1.30.0** generates typed Go methods from hand-written SQL queries.
Each content schema has its own sqlc target (separate `queries/` dir → separate `db/` package).

| Target | Queries | Output | Scope |
|--------|---------|--------|-------|
| shared | `queries/shared/` | `internal/infra/database/db/` | Users, sessions, auth tokens, API keys, MFA, OIDC, settings, activity, libraries |
| movie | `queries/movie/` | `internal/content/movie/db/` | Movies, files, credits, collections, genres, watch progress |
| tvshow | `queries/tvshow/` | `internal/content/tvshow/db/` | Series, seasons, episodes, files, credits, genres, networks, watch progress |
| qar | `queries/qar/` | `internal/content/qar/db/` | *(placeholder)* |

All targets share the same migration source (`migrations/shared/`) so sqlc can resolve cross-schema references.

## Query Logging

```go
// logger.go
type QueryLogger struct {
    logger             *slog.Logger
    slowQueryThreshold time.Duration
}
```

Integrates with pgx tracer. Flags queries exceeding `slowQueryThreshold` as warnings. Logs all query parameters as structured attributes.

## Prometheus Metrics

12 pool metrics (namespace `db_pool`):
- **Counters**: acquire_total, canceled_acquire_total, empty_acquire_total, new_conns_total, max_lifetime_destroy_total, max_idle_destroy_total
- **Gauges**: acquired_conns, constructing_conns, idle_conns, total_conns, max_conns
- **Histograms**: acquire_duration_seconds

## Configuration

From `config.go` `DatabaseConfig` (koanf namespace `database.*`):
```yaml
database:
  url: postgres://user:pass@localhost:5432/revenge?sslmode=disable
  max_conns: 0          # 0 = CPU*2+1
  min_conns: 0
  max_conn_lifetime: 1h
  max_conn_idle_time: 30m
  health_check_period: 1m
```

## Dependencies

- `github.com/jackc/pgx/v5` / `pgxpool` - PostgreSQL driver + pool
- `github.com/golang-migrate/migrate/v4` - Migration runner (iofs source, pgx driver)
- `github.com/prometheus/client_golang` - Metrics
- `log/slog` - Query logging

## Related Documentation

- [ARCHITECTURE.md](../architecture/ARCHITECTURE.md) - System layers
- [CACHE.md](CACHE.md) - Distributed cache (companion to database)
- [JOBS.md](JOBS.md) - River job queue (uses pgx pool)

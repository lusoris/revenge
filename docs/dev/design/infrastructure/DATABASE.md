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
├── metrics.go             # 12 Prometheus pool metrics
├── testing.go             # Embedded PostgreSQL for tests
├── db/                    # sqlc-generated code (425 methods)
│   ├── db.go              # DBTX interface, Queries struct
│   ├── models.go          # 20+ model structs
│   ├── querier.go         # Querier interface (425 methods)
│   └── *.sql.go           # Generated query implementations
├── queries/               # Hand-written SQL source
│   ├── shared/            # Auth, sessions, settings, activity, library, MFA, OIDC, API keys
│   ├── movie/             # Movie metadata queries
│   ├── tvshow/            # TV show series/season/episode queries
│   └── qar/               # Adult content queries (placeholder)
└── migrations/shared/     # 32 versioned migrations (64 .sql files)
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

32 versioned migrations embedded via `//go:embed migrations/shared/*.sql`:

| Range | Coverage |
|-------|----------|
| 001-002 | Schemas (shared, public, qar) + users table |
| 003-013 | Sessions, settings, avatars, auth tokens, password reset, email verification, casbin rules, API keys, OIDC |
| 014-015 | Activity logs, libraries + permissions + scans |
| 016-020 | MFA (TOTP, WebAuthn, backup codes, MFA settings, session MFA tracking) |
| 021-026 | Movies (table, files, credits, collections, genres, watched) |
| 027-032 | Moderator role, TOTP cleanup, fine-grained permissions, failed logins, multi-language movies, TV shows |

**Three-schema model**: `shared` (auth/config), `public` (content), `qar` (adult content)

```go
func MigrateUp(databaseURL string, logger *slog.Logger) error
func MigrateDown(databaseURL string, logger *slog.Logger) error
func MigrateVersion(databaseURL string) (version uint, dirty bool, err error)
func MigrateTo(databaseURL string, version uint, logger *slog.Logger) error
```

## Query Generation (sqlc)

**sqlc v1.30.0** generates 425 typed Go methods from hand-written SQL queries.

Query domains:
- **Users** (749 lines): CRUD, search, count, avatar management
- **Sessions** (522): Token management, MFA tracking, device fingerprints
- **Auth Tokens** (610): Lifecycle, device filtering, revocation
- **API Keys** (305): CRUD, scope management, usage tracking
- **MFA** (782): TOTP, WebAuthn, backup codes, MFA settings
- **OIDC** (876): Providers, states, user links, OAuth2 flow
- **Settings** (619): Server/user settings, upsert operations
- **Activity** (643): Audit logs, search/filter, stats
- **Library** (923): Libraries, permissions, scans, progress

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

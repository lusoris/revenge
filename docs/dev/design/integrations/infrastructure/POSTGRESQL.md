# PostgreSQL Integration

> Primary database for all Revenge data

**Status**: ✅ IMPLEMENTED
**Priority**: 🔴 CRITICAL (Phase 1 - Core Infrastructure)
**Type**: Primary data store

---

## Overview

PostgreSQL is Revenge's primary database, storing all application data including:
- User accounts and sessions
- Media metadata (movies, TV, music, etc.)
- User preferences, ratings, and watch history
- Library configuration
- Application settings

**Version Requirements**:
- **Minimum**: PostgreSQL 15+
- **Recommended**: PostgreSQL 18+ (for improved performance)
- **Extensions**: `pgcrypto`, `pg_trgm` (optional)

---

## Developer Resources

- 📚 **PostgreSQL Docs**: https://www.postgresql.org/docs/current/
- 🔗 **pgx Driver**: https://github.com/jackc/pgx
- 🔗 **sqlc**: https://sqlc.dev/
- 🔗 **River (Job Queue)**: https://riverqueue.com/

---

## Connection Details

**Default Settings**:
| Setting | Value |
|---------|-------|
| Host | `localhost` |
| Port | `5432` |
| Database | `revenge` |
| User | `revenge` |
| SSL Mode | `prefer` |
| Pool Size | `25` |

**Connection String**:
```
postgresql://revenge:password@localhost:5432/revenge?sslmode=prefer
```

---

## Configuration

```yaml
# configs/config.yaml
database:
  host: "${REVENGE_DB_HOST:localhost}"
  port: ${REVENGE_DB_PORT:5432}
  name: "${REVENGE_DB_NAME:revenge}"
  user: "${REVENGE_DB_USER:revenge}"
  password: "${REVENGE_DB_PASSWORD}"

  ssl:
    mode: "prefer"  # disable, prefer, require, verify-ca, verify-full
    ca_file: ""
    cert_file: ""
    key_file: ""

  pool:
    min_conns: 5
    max_conns: 25
    max_conn_lifetime: "1h"
    max_conn_idle_time: "30m"
    health_check_period: "1m"

  # Auto-healing (see DATABASE_AUTO_HEALING.md)
  auto_healing:
    enabled: true
    corruption_detection: true
    auto_reindex: true
```

---

## Schema Design Principles

### Module Isolation

Each content module has its own tables:

```sql
-- Movies module
CREATE TABLE movies (...);
CREATE TABLE movie_genres (...);
CREATE TABLE movie_people (...);
CREATE TABLE movie_user_ratings (...);

-- TV Shows module (separate tables)
CREATE TABLE series (...);
CREATE TABLE seasons (...);
CREATE TABLE episodes (...);
CREATE TABLE episode_user_ratings (...);

-- NO shared content tables!
```

### QAR Content Isolation

Adult content uses separate schema `qar` (Queen Anne's Revenge):

```sql
CREATE SCHEMA IF NOT EXISTS qar;

CREATE TABLE qar.expeditions (...);  -- Movies
CREATE TABLE qar.voyages (...);      -- Scenes
CREATE TABLE qar.crew (...);         -- Performers
```

See [00_SOURCE_OF_TRUTH.md](../../00_SOURCE_OF_TRUTH.md#qar-obfuscation-terminology) for terminology mapping.

### User Data Separation

User-specific data in separate tables:

```sql
-- Per-module user tables
CREATE TABLE movie_user_ratings (
    user_id UUID NOT NULL REFERENCES users(id),
    movie_id UUID NOT NULL REFERENCES movies(id),
    rating INTEGER CHECK (rating BETWEEN 1 AND 100),
    PRIMARY KEY (user_id, movie_id)
);

CREATE TABLE movie_watch_history (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    movie_id UUID NOT NULL REFERENCES movies(id),
    position_ms BIGINT NOT NULL,
    watched_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## sqlc Integration

All database queries use sqlc for type-safe Go code generation.

**sqlc.yaml**:
```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "internal/infra/database/queries"
    schema: "internal/infra/database/migrations"
    gen:
      go:
        package: "db"
        out: "internal/infra/database/db"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_empty_slices: true
        emit_pointers_for_null_types: true
```

**Query Example**:
```sql
-- name: GetMovie :one
SELECT * FROM movies WHERE id = $1;

-- name: ListMoviesByGenre :many
SELECT m.* FROM movies m
JOIN movie_genres mg ON m.id = mg.movie_id
WHERE mg.genre_id = $1
ORDER BY m.title
LIMIT $2 OFFSET $3;
```

---

## Migration Strategy

Migrations are per-module:

```
internal/infra/database/migrations/
├── 000001_users.up.sql
├── 000001_users.down.sql
├── 000002_libraries.up.sql
├── 000002_libraries.down.sql
├── 000010_movies.up.sql
├── 000010_movies.down.sql
├── 000020_tvshows.up.sql
├── 000020_tvshows.down.sql
└── ...
```

**Migration Commands**:
```bash
# Run migrations
go run ./cmd/revenge migrate up

# Rollback
go run ./cmd/revenge migrate down 1

# Create new migration
go run ./cmd/revenge migrate create movie_tags
```

---

## Performance Optimization

### Indexes

```sql
-- Common query patterns
CREATE INDEX idx_movies_title ON movies(title);
CREATE INDEX idx_movies_release_date ON movies(release_date);
CREATE INDEX idx_movies_added_at ON movies(added_at DESC);

-- Full-text search (if not using Typesense)
CREATE INDEX idx_movies_search ON movies
USING gin(to_tsvector('english', title || ' ' || overview));

-- User-specific queries
CREATE INDEX idx_watch_history_user_time ON movie_watch_history(user_id, watched_at DESC);
```

### Partitioning

For large tables (watch history, activity logs):

```sql
CREATE TABLE activity_log (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    action VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
) PARTITION BY RANGE (created_at);

-- Monthly partitions
CREATE TABLE activity_log_2024_01 PARTITION OF activity_log
    FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');
```

### Connection Pooling

pgx's built-in pool with pgxpool:

```go
config, err := pgxpool.ParseConfig(connString)
if err != nil {
    return nil, err
}

config.MinConns = 5
config.MaxConns = 25
config.MaxConnLifetime = time.Hour
config.MaxConnIdleTime = 30 * time.Minute
config.HealthCheckPeriod = time.Minute

pool, err := pgxpool.NewWithConfig(ctx, config)
```

---

## Health Checks

```go
func (db *Database) HealthCheck(ctx context.Context) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    var result int
    err := db.pool.QueryRow(ctx, "SELECT 1").Scan(&result)
    if err != nil {
        return fmt.Errorf("database health check failed: %w", err)
    }
    return nil
}
```

---

## Backup Strategy

### pg_dump

```bash
# Full backup
pg_dump -Fc revenge > revenge_$(date +%Y%m%d).dump

# Restore
pg_restore -d revenge revenge_20240128.dump
```

### Continuous Archiving (WAL)

```ini
# postgresql.conf
archive_mode = on
archive_command = 'cp %p /archive/%f'
```

---

## Docker Compose

```yaml
services:
  postgres:
    image: postgres:18-alpine
    container_name: revenge-postgres
    environment:
      POSTGRES_USER: revenge
      POSTGRES_PASSWORD: ${REVENGE_DB_PASSWORD}
      POSTGRES_DB: revenge
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U revenge"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  postgres_data:
```

---

## Monitoring

### Key Metrics

```sql
-- Active connections
SELECT count(*) FROM pg_stat_activity WHERE state = 'active';

-- Table sizes
SELECT relname, pg_size_pretty(pg_total_relation_size(relid))
FROM pg_stat_user_tables ORDER BY pg_total_relation_size(relid) DESC;

-- Slow queries (pg_stat_statements)
SELECT query, mean_exec_time, calls
FROM pg_stat_statements ORDER BY mean_exec_time DESC LIMIT 10;
```

### Prometheus Exporter

Use `postgres_exporter` for metrics:

```yaml
services:
  postgres-exporter:
    image: prometheuscommunity/postgres-exporter
    environment:
      DATA_SOURCE_NAME: "postgresql://revenge:pass@postgres:5432/revenge?sslmode=disable"
```

---

## Related Documentation

- [Dragonfly (Cache)](DRAGONFLY.md)
- [Typesense (Search)](TYPESENSE.md)
- [River (Job Queue)](RIVER.md)
- [Database Auto-Healing](../../operations/DATABASE_AUTO_HEALING.md)
- [sqlc Instructions](../../.github/instructions/sqlc-database.instructions.md)

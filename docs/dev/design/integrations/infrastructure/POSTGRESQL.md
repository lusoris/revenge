# PostgreSQL Integration

> Primary database for all Revenge data


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
- [Developer Resources](#developer-resources)
- [Connection Details](#connection-details)
- [Configuration](#configuration)
- [Schema Design Principles](#schema-design-principles)
  - [Module Isolation](#module-isolation)
  - [QAR Content Isolation](#qar-content-isolation)
  - [User Data Separation](#user-data-separation)
- [sqlc Integration](#sqlc-integration)
- [Migration Strategy](#migration-strategy)
- [Performance Optimization](#performance-optimization)
  - [Indexes](#indexes)
  - [Partitioning](#partitioning)
  - [Connection Pooling](#connection-pooling)
- [Health Checks](#health-checks)
- [Backup Strategy](#backup-strategy)
  - [pg_dump](#pg-dump)
  - [Continuous Archiving (WAL)](#continuous-archiving-wal)
- [Docker Compose](#docker-compose)
- [Monitoring](#monitoring)
  - [Key Metrics](#key-metrics)
  - [Prometheus Exporter](#prometheus-exporter)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documentation](#related-documentation)

<!-- TOC-END -->

## Status

| Dimension | Status |
|-----------|--------|
| Design | ✅ |
| Sources | ✅ |
| Instructions | ✅ |
| Code | 🔴 |
| Linting | 🔴 |
| Unit Testing | 🔴 |
| Integration Testing | 🔴 |
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


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [Dragonfly Documentation](https://www.dragonflydb.io/docs) | [Local](../../../sources/infrastructure/dragonfly.md) |
| [PostgreSQL Arrays](https://www.postgresql.org/docs/current/arrays.html) | [Local](../../../sources/database/postgresql-arrays.md) |
| [PostgreSQL JSON Functions](https://www.postgresql.org/docs/current/functions-json.html) | [Local](../../../sources/database/postgresql-json.md) |
| [Prometheus Go Client](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus) | [Local](../../../sources/observability/prometheus.md) |
| [Prometheus Metric Types](https://prometheus.io/docs/concepts/metric_types/) | [Local](../../../sources/observability/prometheus-metrics.md) |
| [River Documentation](https://riverqueue.com/docs) | [Local](../../../sources/tooling/river-guide.md) |
| [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) | [Local](../../../sources/tooling/river.md) |
| [Typesense API](https://typesense.org/docs/latest/api/) | [Local](../../../sources/infrastructure/typesense.md) |
| [Typesense Go Client](https://github.com/typesense/typesense-go) | [Local](../../../sources/infrastructure/typesense-go.md) |
| [pgx PostgreSQL Driver](https://pkg.go.dev/github.com/jackc/pgx/v5) | [Local](../../../sources/database/pgx.md) |
| [sqlc](https://docs.sqlc.dev/en/stable/) | [Local](../../../sources/database/sqlc.md) |
| [sqlc Configuration](https://docs.sqlc.dev/en/stable/reference/config.html) | [Local](../../../sources/database/sqlc-config.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Infrastructure](INDEX.md)

### In This Section

- [Dragonfly Integration](DRAGONFLY.md)
- [River Integration](RIVER.md)
- [Typesense Integration](TYPESENSE.md)

### Related Topics

- [Revenge - Architecture v2](../../architecture/01_ARCHITECTURE.md) _Architecture_
- [Revenge - Design Principles](../../architecture/02_DESIGN_PRINCIPLES.md) _Architecture_
- [Revenge - Metadata System](../../architecture/03_METADATA_SYSTEM.md) _Architecture_
- [Revenge - Player Architecture](../../architecture/04_PLAYER_ARCHITECTURE.md) _Architecture_
- [Plugin Architecture Decision](../../architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) _Architecture_

### Indexes

- [Design Index](../../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

## Related Documentation

- [Dragonfly (Cache)](DRAGONFLY.md)
- [Typesense (Search)](TYPESENSE.md)
- [River (Job Queue)](RIVER.md)
- [Database Auto-Healing](../../operations/DATABASE_AUTO_HEALING.md)
- [sqlc Instructions](../../.github/instructions/sqlc-database.instructions.md)

# Database Auto-Healing & Consistency Restoration

> Automatic recovery strategies for PostgreSQL database corruption and inconsistency



<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Context](#context)
- [Detection Tools](#detection-tools)
  - [1. `amcheck` Extension (PostgreSQL 10+)](#1-amcheck-extension-postgresql-10)
  - [2. `pg_visibility` Extension (9.6+)](#2-pg-visibility-extension-96)
  - [3. `pg_checksums` (PostgreSQL 12+)](#3-pg-checksums-postgresql-12)
  - [4. Custom Consistency Checks](#4-custom-consistency-checks)
    - [FK Integrity](#fk-integrity)
    - [Enum Constraints](#enum-constraints)
    - [JSON Schema Validation](#json-schema-validation)
- [Automatic Repair Strategies](#automatic-repair-strategies)
  - [1. Index Corruption → REINDEX](#1-index-corruption-reindex)
  - [2. Visibility Map Corruption → VACUUM FULL](#2-visibility-map-corruption-vacuum-full)
  - [3. Orphaned Data → Cascade Delete](#3-orphaned-data-cascade-delete)
  - [4. Invalid Data → Nullify + Log](#4-invalid-data-nullify-log)
  - [5. Constraint Violations → Rollback Transaction](#5-constraint-violations-rollback-transaction)
- [Background Jobs (River)](#background-jobs-river)
  - [Health Check Job (Every 10 minutes)](#health-check-job-every-10-minutes)
  - [Index Check Job (Weekly)](#index-check-job-weekly)
  - [Orphan Cleanup Job (Daily)](#orphan-cleanup-job-daily)
- [Recovery from Corruption](#recovery-from-corruption)
  - [CRITICAL: Snapshot Before Repair](#critical-snapshot-before-repair)
  - [Repair Workflow](#repair-workflow)
- [Monitoring & Alerts](#monitoring-alerts)
  - [Metrics to Track (via pkg/metrics)](#metrics-to-track-via-pkgmetrics)
  - [Alerts](#alerts)
- [Configuration (config.yaml)](#configuration-configyaml)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Tools](#related-tools)
- [References](#references)

<!-- TOC-END -->

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | 🔴 |  |
| Sources | 🔴 |  |
| Instructions | 🔴 |  |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |

---

## Context

Revenge's PostgreSQL database may encounter corruption or inconsistency due to:
- Hardware issues (RAID battery failures, SSD power-loss protection, defective RAM/CPU)
- Software bugs (PostgreSQL, OS, app-level)
- Admin errors (manual file modifications, improper failover)
- OS crashes with `fsync=off`
- glibc locale collation changes across OS upgrades

## Detection Tools

### 1. `amcheck` Extension (PostgreSQL 10+)
```sql
-- General smoke test (B-Tree indexes)
CREATE EXTENSION IF NOT EXISTS amcheck;
SELECT bt_index_check(c.oid, true) FROM pg_class c WHERE c.relkind = 'i';

-- Heap verification (with heapallindexed option)
SELECT bt_index_parent_check(c.oid, true, true) FROM pg_class c WHERE c.relkind = 'i';
```

**Usage**: Run weekly as background job (River task)

### 2. `pg_visibility` Extension (9.6+)
```sql
-- Visibility map integrity check
CREATE EXTENSION IF NOT EXISTS pg_visibility;
SELECT * FROM pg_check_visible('table_name');
```

### 3. `pg_checksums` (PostgreSQL 12+)
```bash
# Enable checksums (offline, rewrites all files)
pg_checksums --enable -D /var/lib/postgresql/data

# Verify checksums
pg_checksums --check -D /var/lib/postgresql/data
```

**Implementation**: Enable at installation, verify weekly via cron.

### 4. Custom Consistency Checks

#### FK Integrity
```sql
-- Check orphaned foreign keys (example: movies.tmdb_id)
SELECT m.id, m.title FROM movies m
LEFT JOIN external_metadata em ON m.id = em.movie_id
WHERE em.id IS NULL AND m.tmdb_id IS NOT NULL;
```

#### Enum Constraints
```sql
-- Check invalid enum values
SELECT * FROM movies WHERE content_rating NOT IN (
  SELECT enumlabel::text FROM pg_enum WHERE enumtypid = 'content_rating'::regtype
);
```

#### JSON Schema Validation
```sql
-- Example: Validate metadata_json structure
SELECT id, title FROM movies WHERE NOT (
  metadata_json ? 'release_date' AND
  jsonb_typeof(metadata_json->'release_date') = 'string'
);
```

## Automatic Repair Strategies

### 1. Index Corruption → REINDEX
```sql
-- Detect corrupted indexes via bt_index_check failures
-- Auto-repair via REINDEX CONCURRENTLY (no downtime)
REINDEX INDEX CONCURRENTLY movies_pkey;
REINDEX TABLE CONCURRENTLY movies;
```

**Implementation**: River job triggered by amcheck failures.

### 2. Visibility Map Corruption → VACUUM FULL
```sql
-- If pg_check_visible reports issues:
VACUUM FULL ANALYZE movies;
```

**Caveat**: Requires `ACCESS EXCLUSIVE` lock → schedule maintenance window.

### 3. Orphaned Data → Cascade Delete
```sql
-- Example: Remove orphaned movie_genres entries
DELETE FROM movie_genres WHERE movie_id NOT IN (SELECT id FROM movies);
```

**Implementation**: Weekly cleanup job.

### 4. Invalid Data → Nullify + Log
```sql
-- Mark corrupted records for manual review
UPDATE movies SET metadata_json = NULL, needs_review = TRUE
WHERE NOT jsonb_typeof(metadata_json) = 'object';

INSERT INTO audit_log (event_type, severity, message, metadata)
VALUES ('data_corruption', 'error', 'Invalid metadata_json detected', jsonb_build_object('movie_id', id));
```

### 5. Constraint Violations → Rollback Transaction
```go
// In service layer (Go):
tx, _ := db.Begin(ctx)
defer tx.Rollback(ctx)

// Attempt repair
if err := repairData(ctx, tx); err != nil {
    log.Error("Repair failed, rolling back", "error", err)
    return err
}

tx.Commit(ctx)
```

## Background Jobs (River)

### Health Check Job (Every 10 minutes)
```go
type DBHealthCheckArgs struct{}

func (DBHealthCheckArgs) Kind() string { return "db.health_check" }

type DBHealthCheckWorker struct {
    river.WorkerDefaults[DBHealthCheckArgs]
    db *pgxpool.Pool
}

func (w *DBHealthCheckWorker) Work(ctx context.Context, job *river.Job[DBHealthCheckArgs]) error {
    // 1. Check pg_stat_database for anomalies
    // 2. Run amcheck on critical indexes
    // 3. Log warnings if issues detected
    // 4. Enqueue repair jobs if needed
    return nil
}
```

### Index Check Job (Weekly)
```go
type IndexCheckArgs struct {
    TableName string `json:"table_name"`
}

func (IndexCheckArgs) Kind() string { return "db.index_check" }

// Run bt_index_check on all indexes for table
```

### Orphan Cleanup Job (Daily)
```go
type OrphanCleanupArgs struct{}

func (OrphanCleanupArgs) Kind() string { return "db.orphan_cleanup" }

// Run FK integrity checks + cleanup
```

## Recovery from Corruption

### CRITICAL: Snapshot Before Repair
```bash
# File-system level copy BEFORE any repair attempts
sudo systemctl stop postgresql
sudo cp -a /var/lib/postgresql/data /backup/pg-snapshot-$(date +%s)
sudo systemctl start postgresql
```

### Repair Workflow
1. **Detect** → amcheck, pg_visibility, custom checks
2. **Snapshot** → File-system copy (not pg_dump)
3. **Attempt Auto-Repair** → REINDEX, VACUUM, DELETE orphans
4. **Verify** → Re-run checks
5. **Log + Alert** → If auto-repair fails, escalate to admin
6. **Manual Intervention** → Use `pg_hexedit` on snapshot copy for forensics

## Monitoring & Alerts

### Metrics to Track (via pkg/metrics)
- `db.corruption.events` (counter)
- `db.repair.attempts` (counter)
- `db.repair.success_rate` (gauge)
- `db.health_check.duration` (histogram)

### Alerts
- **Critical**: `amcheck` failures > 0
- **Warning**: Orphaned records > 1000
- **Info**: Checksums verified successfully

## Configuration (config.yaml)

```yaml
database:
  auto_healing:
    enabled: true
    health_check_interval: 10m
    index_check_interval: 1w
    orphan_cleanup_interval: 1d
    snapshot_before_repair: true
    max_auto_repair_attempts: 3
    alert_on_failure: true
```


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [PostgreSQL Arrays](https://www.postgresql.org/docs/current/arrays.html) | [Local](../../sources/database/postgresql-arrays.md) |
| [PostgreSQL JSON Functions](https://www.postgresql.org/docs/current/functions-json.html) | [Local](../../sources/database/postgresql-json.md) |
| [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) | [Local](../../sources/tooling/river.md) |
| [pgx PostgreSQL Driver](https://pkg.go.dev/github.com/jackc/pgx/v5) | [Local](../../sources/database/pgx.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Operations](INDEX.md)

### In This Section

- [Advanced Patterns & Best Practices](BEST_PRACTICES.md)
- [Branch Protection Rules](BRANCH_PROTECTION.md)
- [Clone repository](DEVELOPMENT.md)
- [GitFlow Workflow Guide](GITFLOW.md)
- [Revenge - Reverse Proxy & Deployment Best Practices](REVERSE_PROXY.md)
- [revenge - Setup Guide](SETUP.md)

### Related Topics

- [Revenge - Architecture v2](../architecture/01_ARCHITECTURE.md) _Architecture_
- [Revenge - Design Principles](../architecture/02_DESIGN_PRINCIPLES.md) _Architecture_
- [Revenge - Metadata System](../architecture/03_METADATA_SYSTEM.md) _Architecture_
- [Revenge - Player Architecture](../architecture/04_PLAYER_ARCHITECTURE.md) _Architecture_
- [Plugin Architecture Decision](../architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) _Architecture_

### Indexes

- [Design Index](../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

## Related Tools

| Tool | Purpose |
|------|---------|
| `pg_hexedit` | Low-level hex editor for PostgreSQL files (forensics) |
| `pg_dirtyread` | Read data from corrupted tables |
| `pg_filedump` | Dump PostgreSQL relation files (debugging) |

## References

- [PostgreSQL Corruption Wiki](https://wiki.postgresql.org/wiki/Corruption)
- [amcheck Documentation](https://www.postgresql.org/docs/current/amcheck.html)
- [pg_visibility Documentation](https://www.postgresql.org/docs/current/pgvisibility.html)
- [Routine Reindexing](https://www.postgresql.org/docs/current/routine-reindex.html)


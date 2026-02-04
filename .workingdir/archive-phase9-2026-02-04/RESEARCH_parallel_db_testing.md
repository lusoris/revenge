# Parallel Database Testing - Research & Implementation

## Problem

Database integration tests are slow (~3-5s per test) because each test:
1. Starts a new embedded PostgreSQL instance
2. Runs all migrations
3. Executes test
4. Stops PostgreSQL

With 20+ database tests, this takes 60-100+ seconds.

## Solution: Template Database Pattern

PostgreSQL can clone databases from templates **instantly** (~10ms).

### How It Works

1. **Once per test run**: Start ONE shared PostgreSQL instance
2. **Once per test run**: Create template database with all migrations applied
3. **Per test**: `CREATE DATABASE test_xxx TEMPLATE template_db` (instant clone)
4. **Per test cleanup**: `DROP DATABASE test_xxx`

### Performance Comparison

| Approach | Time per Test | 20 Tests |
|----------|---------------|----------|
| New PG instance per test | 3-5s | 60-100s |
| Template DB cloning | 10-50ms | 200ms-1s |

**Speed improvement: 60-500x faster**

## Key PostgreSQL Concepts

### Template Databases
- Every database can be a template for `CREATE DATABASE`
- `CREATE DATABASE newdb TEMPLATE templatedb` copies instantly
- Template must have no active connections during clone
- All objects (tables, indexes, data) are copied

### Configuration
- `datistemplate = true` marks DB as template
- `datallowconn = false` prevents connections (for protection)

## Implementation in Go

### embedded-postgres Configuration

```go
// Use unique RuntimePath per process to avoid conflicts
runtimePath := fmt.Sprintf("/tmp/embedded-postgres-%d", os.Getpid())

postgres := embeddedpostgres.NewDatabase(
    embeddedpostgres.DefaultConfig().
        Port(15555).
        RuntimePath(runtimePath).
        StartTimeout(60*time.Second),
)
```

### Template DB Creation

```go
// Create template once
conn.Exec(ctx, "CREATE DATABASE revenge_template")

// Run migrations on template
MigrateUp(templateURL, logger)

// For each test - instant clone
conn.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s TEMPLATE revenge_template", testDBName))
```

### Test Cleanup

```go
// Force disconnect before drop
conn.Exec(ctx, `
    SELECT pg_terminate_backend(pid)
    FROM pg_stat_activity
    WHERE datname = $1 AND pid <> pg_backend_pid()
`, dbName)

// Drop test database
conn.Exec(ctx, fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
```

## Best Practices

1. **Use sync.Once** for shared resources (postgres instance, template DB)
2. **Unique DB names** per test: `test_{pid}_{counter}`
3. **t.Cleanup()** for automatic database cleanup
4. **t.Parallel()** safe - each test gets isolated database

## Alternative Approaches Considered

### 1. Testcontainers
- Docker-based, more isolation
- Slower startup than embedded-postgres
- Good for CI where Docker is available

### 2. Schema Isolation (same DB, different schemas)
- Even faster than template cloning
- Requires careful schema management
- Risk of test pollution

### 3. Transaction Rollback
- Wrap each test in transaction, rollback at end
- Fastest, but limits what tests can do (no DDL in tests)
- Can't test transaction behavior

## Files Created

- `internal/testutil/testdb.go` - Main TestDB helper
- `internal/testutil/testdb_migrate.go` - Migration runner
- `internal/testutil/migrations/` - Embedded migrations copy

## Usage Example

```go
func TestUserCreate(t *testing.T) {
    t.Parallel() // Safe!

    db := testutil.NewTestDB(t)
    pool := db.Pool()

    // Test code using pool
    _, err := pool.Exec(ctx, "INSERT INTO users ...")
    require.NoError(t, err)

    // Automatic cleanup via t.Cleanup()
}
```

## Sources

- PostgreSQL Template Databases: https://www.postgresql.org/docs/current/manage-ag-templatedbs.html
- embedded-postgres: https://github.com/fergusstrange/embedded-postgres
- testcontainers-go: https://github.com/testcontainers/testcontainers-go
- pgx testing patterns: https://github.com/jackc/pgx

## Date
2026-02-02

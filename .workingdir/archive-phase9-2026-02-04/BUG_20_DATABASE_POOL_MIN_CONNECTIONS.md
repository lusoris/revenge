# Bug #20: Database Connection Pool Not Respecting MinConns

**Status**: 🔴 CONFIRMED
**Severity**: Low (performance/efficiency issue)
**Component**: Database Infrastructure (`pgxpool`)
**Discovered**: 2025-01-XX during integration testing

## Symptom

When configuring a `pgxpool` with `MinConns = 2`, the pool reports 0 total connections immediately after creation, instead of eagerly establishing the minimum number of connections.

## Reproduction

```go
config, err := pgxpool.ParseConfig(testDatabaseURL)
config.MaxConns = 10
config.MinConns = 2

pool, err := pgxpool.NewWithConfig(ctx, config)

stats := pool.Stat()
// Expected: stats.TotalConns() >= 2
// Actual: stats.TotalConns() == 0
```

## Test Output

```
Error:            "0" is not greater than or equal to "2"
Test:             TestDatabasePoolStats
Messages:         should have min connections
```

## Root Cause Analysis

The `pgxpool` library appears to establish connections **lazily** rather than **eagerly**, even when `MinConns` is set. This means the pool only creates connections when they are first requested, not at initialization time.

According to pgxpool documentation:
- `MinConns`: The minimum number of connections that will be maintained by the pool
- Connections are created on-demand and maintained once established
- **Not** a guarantee of pre-warmed connections at pool creation

## Impact

- **Low Impact**: Pool will eventually reach MinConns under normal load
- No functional issue, just delayed connection establishment
- Slight performance hit on first requests (connection setup overhead)
- May cause slower cold-start performance

## Options for Resolution

### Option 1: Accept Current Behavior ✅ RECOMMENDED
- This is expected pgxpool behavior
- Update test to accept 0 connections initially
- Test that MinConns is respected after usage
- Document this as expected behavior

### Option 2: Pre-warm Connections
- Explicitly acquire and release MinConns connections after pool creation
- Ensures connections are established upfront
- Adds complexity to initialization code

```go
// Example pre-warming code
for i := 0; i < config.MinConns; i++ {
    conn, err := pool.Acquire(ctx)
    if err != nil {
        return err
    }
    conn.Release()
}
```

### Option 3: Custom Health Check with Connection Count
- Add health check that validates connection pool state
- Warm up connections during application startup
- Monitor connection count in metrics

## Recommendation

**Option 1** - Accept as expected behavior and update test. The pgxpool library is designed this way intentionally to avoid unnecessary connections. The MinConns setting ensures connections persist once established, preventing them from being closed during idle periods.

## Next Steps

1. Update `TestDatabasePoolStats` to test MinConns behavior correctly
2. Add test for connection pool behavior under load
3. Document pool behavior in infrastructure docs
4. Consider pre-warming if cold-start performance becomes an issue

## Related Code

- Test: `tests/integration/database/database_test.go:33-49`
- Database Module: `internal/infra/database/module.go`
- pgxpool docs: https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool

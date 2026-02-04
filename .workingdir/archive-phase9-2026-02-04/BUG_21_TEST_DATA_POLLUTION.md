# Bug #21: Integration Test Data Pollution

**Status**: 🔴 CONFIRMED
**Severity**: Medium (breaks test reliability)
**Component**: Integration Tests - Database Layer
**Discovered**: 2025-01-XX during comprehensive integration testing

## Symptom

Multiple integration tests failing with duplicate key constraint violations:

```
ERROR: duplicate key value violates unique constraint "users_username_key" (SQLSTATE 23505)
Test: TestDatabaseConcurrentOperations
Test: TestDatabaseNullableTypes
```

Tests are attempting to create users with the same usernames (`concurrent_test`, `nullable_test`) that already exist in the database from previous test runs.

## Root Cause

**Test cleanup is not working properly**:
1. Tests use `defer` to delete test users after completion
2. If test fails or panics, deferred cleanup may not execute
3. Database persists between test runs (Docker volume)
4. Subsequent test runs encounter existing data

**Lack of test isolation**:
- No test-specific prefixes or UUIDs in usernames
- No truncate/reset between tests
- No unique constraint handling in test setup

## Impact

- **High**: Tests become flaky and unreliable
- Cannot run tests multiple times without manual DB cleanup
- Test failures cascade - one failure breaks future runs
- Integration CI/CD would be broken

## Failing Tests

1. `TestDatabaseConcurrentOperations` - Line 113
   - Username: `concurrent_test`
   - Violates unique constraint on second run

2. `TestDatabaseNullableTypes` - Line 207
   - Username: `nullable_test`
   - Violates unique constraint on second run

## Solutions

### Option 1: Unique Test Data ✅ RECOMMENDED
Use timestamps or UUIDs to make test data unique on each run:

```go
username := fmt.Sprintf("test_user_%d", time.Now().UnixNano())
// OR
username := fmt.Sprintf("test_user_%s", uuid.New().String()[:8])
```

**Pros**: Simple, no DB state dependencies
**Cons**: Leaves test data in DB (may accumulate)

### Option 2: Robust Cleanup with Idempotent Setup
Delete before create (idempotent):

```go
// Try to delete if exists (ignore errors)
_ = queries.DeleteUserByUsername(ctx, "concurrent_test")

// Now safe to create
user, err := queries.CreateUser(ctx, ...)
```

**Pros**: Clean database state, predictable usernames
**Cons**: Requires DeleteByUsername query (may not exist)

### Option 3: Transaction-Based Tests
Wrap entire test in transaction and rollback:

```go
tx, _ := pool.Begin(ctx)
defer tx.Rollback(ctx)
txQueries := queries.WithTx(tx)
// All operations in transaction
```

**Pros**: Perfect isolation, no cleanup needed
**Cons**: Can't test cross-transaction behavior, some operations may not work in TX

### Option 4: Database Reset Before Each Test
Truncate tables in `BeforeEach` or `SetupTest`:

```go
func cleanupTestData(t *testing.T, pool *pgxpool.Pool) {
    _, _ = pool.Exec(ctx, "TRUNCATE shared.users CASCADE")
}
```

**Pros**: Fresh state for every test
**Cons**: Slow, CASCADE may delete too much

## Recommended Approach

**Combination of Option 1 + Option 3**:
1. Use unique usernames with timestamps for non-transactional tests
2. Use transaction-wrapped tests where possible for full isolation
3. Add cleanup helper that handles both delete and rollback

```go
func createTestUser(t *testing.T, ctx context.Context, queries *db.Queries, suffix string) db.SharedUser {
    username := fmt.Sprintf("test_%s_%d", suffix, time.Now().UnixNano())
    user, err := queries.CreateUser(ctx, db.CreateUserParams{
        Username:     username,
        Email:        fmt.Sprintf("%s@example.com", username),
        PasswordHash: "test_hash",
    })
    require.NoError(t, err)

    t.Cleanup(func() {
        _ = queries.DeleteUser(ctx, user.ID)
    })

    return user
}
```

## Next Steps

1. ✅ Implement unique username generation for affected tests
2. Create test helper functions for common operations
3. Add transaction-based variants where applicable
4. Document test data management strategy
5. Consider DB reset script for local development

## Related Code

- `tests/integration/database/database_test.go:90-115` (TestDatabaseConcurrentOperations)
- `tests/integration/database/database_test.go:185-207` (TestDatabaseNullableTypes)
- Need: `DeleteUserByUsername` query or similar cleanup mechanism

## Test Fix Required

Immediate fix: Use timestamp-based usernames to avoid conflicts.

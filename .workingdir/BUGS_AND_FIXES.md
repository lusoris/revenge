# Bugs and Fixes Log

**Purpose**: Document all bugs found and fixes implemented during Phase A development.

**Format**: Each entry includes bug description, root cause, fix implemented, files changed, and testing status.

---

## A7.1.1: Missing Transaction in User Registration

**Status**: ✅ FIXED
**Date**: 2026-02-05
**Priority**: P0 (Critical - Data Integrity)
**Category**: Security / Data Integrity

### Bug Description

User registration was not atomic. The Register method performed two database operations without a transaction:
1. CreateUser - insert user into `shared.users`
2. CreateEmailVerificationToken - insert token into `shared.email_verification_tokens`

If CreateEmailVerificationToken failed after CreateUser succeeded, the user account would exist in the database without a verification token, creating an orphaned user account that could never verify their email.

### Root Cause

The service method called repository methods sequentially without wrapping them in a transaction:

```go
// Create user in database
user, err := s.repo.CreateUser(ctx, ...)
// ...
// Store verification token in database
_, err = s.repo.CreateEmailVerificationToken(ctx, ...)
```

### Fix Implemented

**Approach**: Wrap both operations in a PostgreSQL transaction using pgxpool.Pool.

**Changes**:
1. Added `pool *pgxpool.Pool` field to auth Service struct
2. Modified NewService constructor to accept pool parameter
3. Refactored Register method to use transaction pattern:
   - Begin transaction with `tx, err := s.pool.Begin(ctx)`
   - Defer rollback with `defer tx.Rollback(ctx)`
   - Create transaction-scoped queries with `db.New(tx)`
   - Execute both CreateUser and CreateEmailVerificationToken within transaction
   - Commit transaction at end

**Transaction Pattern** (following RBAC adapter pattern):
```go
tx, err := s.pool.Begin(ctx)
if err != nil {
    return nil, fmt.Errorf("failed to begin transaction: %w", err)
}
defer func() {
    _ = tx.Rollback(ctx)
}()

txQueries := db.New(tx)

// Database operations...
user, err := txQueries.CreateUser(ctx, ...)
// ...
_, err = txQueries.CreateEmailVerificationToken(ctx, ...)

if err := tx.Commit(ctx); err != nil {
    return nil, fmt.Errorf("failed to commit transaction: %w", err)
}
```

### Files Changed

1. **internal/service/auth/service.go**
   - Added `pool *pgxpool.Pool` field to Service struct
   - Added `"github.com/jackc/pgx/v5/pgxpool"` import
   - Modified NewService signature to accept pool as first parameter
   - Refactored Register method to use transactions (lines 65-114)

2. **internal/service/auth/module.go**
   - Added `"github.com/jackc/pgx/v5/pgxpool"` import
   - Updated Service provider to inject pool parameter

3. **internal/service/auth/service_testing.go**
   - Added `"github.com/jackc/pgx/v5/pgxpool"` import
   - Updated NewServiceForTesting to accept pool parameter
   - Updated NewServiceForTestingWithEmail to accept pool parameter

4. **internal/service/auth/service_integration_test.go**
   - Updated all test functions to pass testDB.Pool() to NewServiceForTesting
   - Added TestService_Register_TransactionAtomicity integration test

5. **internal/service/auth/service_exhaustive_test.go**
   - Updated setupMockService to pass nil for pool (mock tests)
   - Added comment explaining transaction-based methods can't be mocked

### Testing

**Integration Test Added**: `TestService_Register_TransactionAtomicity`
- Verifies successful registration creates both user and token
- Verifies failed registration (duplicate username) prevents orphaned records
- Documents transaction atomicity behavior

**Test Location**: internal/service/auth/service_integration_test.go:402-466

**Note**: Integration tests require Docker/testcontainers and may not run in all environments.

### Impact

**Before**: Race condition could create orphaned user accounts
**After**: User creation and token generation are atomic - both succeed or both fail

**Related Security Issues**: A7.1.2 (Avatar Upload), A7.1.3 (Session Refresh)

### References

- Source: [TODO_A7_SECURITY_FIXES.md](TODO_A7_SECURITY_FIXES.md) lines 27-76
- Pattern Reference: internal/service/rbac/adapter.go lines 80-127
- Report: [REPORT_2_IMPLEMENTATION_VERIFICATION.md](REPORT_2_IMPLEMENTATION_VERIFICATION.md)

---

**Next Bug**: A7.1.2 - Avatar Upload Transaction

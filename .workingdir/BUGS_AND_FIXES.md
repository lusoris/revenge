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

## A7.1.2: Missing Transaction in Avatar Upload

**Status**: ✅ FIXED
**Date**: 2026-02-05
**Priority**: P0 (Critical - Data Integrity)
**Category**: Security / Data Integrity

### Bug Description

Avatar upload was not atomic. The UploadAvatar method performed four database operations without a transaction:
1. GetLatestAvatarVersion - get next version number
2. UnsetCurrentAvatars - mark existing avatars as not current
3. CreateAvatar - insert new avatar record
4. UpdateUser - update user's avatar_url field

If any operation failed after previous ones succeeded, the database would be in an inconsistent state. For example:
- Avatars could be unset but no new avatar created
- New avatar created but user's avatar_url not updated
- Multiple avatars marked as current due to race conditions

Additionally, the UpdateUser call was wrapped in error handling that logged but didn't fail the operation, meaning avatar upload could "succeed" without actually updating the user's avatar URL.

### Root Cause

The service method called repository methods sequentially without wrapping them in a transaction, and also stored the file to storage before any database operations, requiring manual cleanup on any failure.

### Fix Implemented

**Approach**: Wrap all database operations in a PostgreSQL transaction, store file first (outside transaction), cleanup file if transaction fails.

**Transaction Flow**:
1. Store file to storage (outside transaction)
2. Begin transaction
3. Get latest version (within transaction)
4. Unset current avatars (within transaction)
5. Create new avatar record (within transaction)
6. Update user avatar_url (within transaction, now fails entire operation on error)
7. Commit transaction
8. If any DB operation fails, cleanup stored file and rollback transaction

**Changes**:
1. Added `pool *pgxpool.Pool` field to user Service struct
2. Modified NewService constructor to accept pool parameter
3. Refactored UploadAvatar method to use transaction pattern:
   - Begin transaction with `tx, err := s.pool.Begin(ctx)`
   - Defer rollback with `defer tx.Rollback(ctx)`
   - Create transaction-scoped queries with `db.New(tx)`
   - Execute all DB operations within transaction
   - Properly handle IP address parsing from *string to netip.Addr
   - Use correct UpdateUser parameters (UserID not ID)
   - Commit transaction at end
   - Cleanup stored file on any DB error

### Files Changed

1. **internal/service/user/service.go**
   - Added `pool *pgxpool.Pool` field to Service struct
   - Added `"github.com/jackc/pgx/v5/pgxpool"` import
   - Added `"net/netip"` import for IP address handling
   - Modified NewService signature to accept pool as first parameter
   - Refactored UploadAvatar method to use transactions (lines 329-413)
   - Added IP address parsing logic (ParseAddr from *string to netip.Addr)
   - Fixed UpdateUser call to use correct parameters and fail on error

2. **internal/service/user/module.go**
   - Added `"github.com/jackc/pgx/v5/pgxpool"` import
   - Updated Service provider to inject pool parameter

### Implementation Details

**IP Address Handling**:
```go
// Parse IP address if provided
if metadata.UploadedFromIP != nil {
    addr, err := netip.ParseAddr(*metadata.UploadedFromIP)
    if err != nil {
        _ = s.storage.Delete(ctx, storedKey)
        return nil, fmt.Errorf("failed to parse IP address: %w", err)
    }
    createParams.UploadedFromIp = addr
}
```

**UpdateUser Fix**:
Changed from silently ignoring errors:
```go
_, err = s.repo.UpdateUser(ctx, userID, UpdateUserParams{
    AvatarURL: &avatarURL,
})
if err != nil {
    // Log error but don't fail
    _ = err
}
```

To properly failing within transaction:
```go
_, err = txQueries.UpdateUser(ctx, db.UpdateUserParams{
    UserID:    userID,
    AvatarUrl: &avatarURL,
})
if err != nil {
    _ = s.storage.Delete(ctx, storedKey)
    return nil, fmt.Errorf("failed to update user avatar: %w", err)
}
```

### Testing

**Integration Test Required**: Transaction rollback test for avatar upload
- Verify all operations succeed together or all fail
- Verify stored file is cleaned up on DB failure
- Verify no orphaned avatar records or inconsistent states

**Test Status**: ⚠️ PENDING (requires Docker/testcontainers)

### Impact

**Before**:
- Race conditions could create multiple "current" avatars
- Avatar records could exist without user.avatar_url update
- Inconsistent database state on partial failures
- Silent failures on UpdateUser

**After**:
- All DB operations are atomic - succeed together or fail together
- Stored files are cleaned up on any DB error
- UpdateUser failures now properly fail the entire operation
- No orphaned records or inconsistent states

**Related Security Issues**: A7.1.1 (User Registration), A7.1.3 (Session Refresh)

### References

- Source: [TODO_A7_SECURITY_FIXES.md](TODO_A7_SECURITY_FIXES.md) lines 79-133
- Pattern Reference: internal/service/rbac/adapter.go, internal/service/auth/service.go Register method
- Report: [REPORT_2_IMPLEMENTATION_VERIFICATION.md](REPORT_2_IMPLEMENTATION_VERIFICATION.md)

---

## A7.1.3: Incorrect Operation Order in Session Refresh

**Status**: ✅ FIXED
**Date**: 2026-02-05
**Priority**: P1 (High - User Experience / Security)
**Category**: Security / Data Integrity

### Bug Description

Session refresh had operations in wrong order: it revoked the old session BEFORE creating the new one. If CreateSession failed after RevokeSession succeeded, the user would be left without any valid session, forcing them to log in again.

**Operation Flow (Before Fix)**:
1. Get old session by refresh token
2. Generate new tokens
3. **Revoke old session** ← Done first
4. **Create new session** ← Done second
5. If step 4 fails, user has no session!

### Root Cause

The RefreshSession method (lines 149-181) performed operations in non-resilient order. The old session was revoked (line 150) before the new session was created (line 169), violating the principle of "create before destroy" for critical resources.

### Fix Implemented

**Approach**: Reorder operations to create new session first, then revoke old session only if new one succeeded.

**New Operation Flow**:
1. Get old session by refresh token
2. Generate new tokens
3. **Create new session first** ← Done first
4. **Revoke old session only if creation succeeded** ← Done second
5. If revocation fails, log error but don't fail (new session is valid)

**Changes**:
1. Moved CreateSession block before RevokeSession block
2. Changed RevokeSession error handling from fail to warn
3. Added comments explaining the importance of operation order
4. Added error logging with zap.Warn for revocation failures

### Implementation Details

**Before (Problematic Order)**:
```go
// Revoke old session first
if err := s.repo.RevokeSession(ctx, session.ID, &reason); err != nil {
    return "", "", fmt.Errorf("failed to revoke old session: %w", err)
}

// Create new session second
_, err = s.repo.CreateSession(ctx, CreateSessionParams{...})
if err != nil {
    return "", "", fmt.Errorf("failed to create refreshed session: %w", err)
}
```

**After (Resilient Order)**:
```go
// Create new session first
_, err = s.repo.CreateSession(ctx, CreateSessionParams{...})
if err != nil {
    return "", "", fmt.Errorf("failed to create refreshed session: %w", err)
}

// Revoke old session only after new one exists
if err := s.repo.RevokeSession(ctx, session.ID, &reason); err != nil {
    // Log but don't fail - new session is valid
    s.logger.Warn("failed to revoke old session during refresh",
        zap.Error(err), zap.String("session_id", session.ID.String()))
}
```

### Files Changed

1. **internal/service/session/service.go**
   - Reordered RefreshSession method operations (lines 146-185)
   - Moved CreateSession before RevokeSession
   - Changed RevokeSession error handling from return error to log warning
   - Added explanatory comments about operation order

### Testing

**Manual Testing**:
- Verify session refresh succeeds normally
- Verify user keeps new session even if revocation fails
- Verify error is logged when revocation fails

**Integration Test Required**: ⚠️ PENDING
- Test CreateSession failure doesn't leave user without session (no longer possible with fix)
- Test RevokeSession failure doesn't prevent refresh success
- Verify old session is revoked in normal case

### Impact

**Before**:
- CreateSession failure left user without any valid session
- User forced to log in again after refresh failure
- Poor user experience during transient errors
- Unnecessary authentication loops

**After**:
- User always gets new session or keeps old session (atomic from user perspective)
- CreateSession failure: old session still valid, user not logged out
- RevokeSession failure: new session valid, old session eventually expires, logged for cleanup
- Resilient to transient database errors

**Edge Case Handling**:
- If RevokeSession fails, old refresh token remains in database
- Old session will expire naturally based on expiry time
- Logged warning allows monitoring and manual cleanup if needed
- User security not compromised (new session is valid and secure)

**Related Security Issues**: A7.1.1 (User Registration), A7.1.2 (Avatar Upload)

### References

- Source: [TODO_A7_SECURITY_FIXES.md](TODO_A7_SECURITY_FIXES.md) lines 136-171
- Report: [REPORT_2_IMPLEMENTATION_VERIFICATION.md](REPORT_2_IMPLEMENTATION_VERIFICATION.md)

---

## A7.2: Username Enumeration via Timing Attack in Login

**Status**: ✅ FIXED
**Date**: 2026-02-05
**Priority**: P0 (Critical - Security)
**Category**: Security / Authentication

### Bug Description

Login timing differed based on whether user exists, allowing attackers to enumerate valid usernames via timing analysis. This is a classic timing attack vulnerability.

**Timing Difference**:
- User not found: Fast path (~1ms) - returns error immediately without password hash comparison
- User found: Slow path (~50-100ms) - performs Argon2id password hash comparison

**Attack Scenario**:
1. Attacker sends login requests with various usernames
2. Measures response time for each request
3. Fast responses = username doesn't exist
4. Slow responses = username exists (regardless of password correctness)
5. Attacker builds list of valid usernames for targeted attacks

### Root Cause

The Login method (lines 256-297) had two execution paths with significantly different timing:

**Fast Path** (user not found):
```go
user, err := s.repo.GetUserByUsername(ctx, username)
if err != nil {
    user, err = s.repo.GetUserByEmail(ctx, username)
    if err != nil {
        return nil, errors.New("invalid username or password") // NO hash comparison
    }
}
```

**Slow Path** (user found):
```go
// Verify password using Argon2id (takes ~50-100ms)
match, err := s.hasher.VerifyPassword(password, user.PasswordHash)
```

### Fix Implemented

**Approach**: Always perform password hash comparison, even if user doesn't exist, using a precomputed dummy hash to ensure constant-time behavior.

**Implementation**:
1. Generate precomputed Argon2id dummy hash
2. Track whether user was found (but don't return early)
3. Select hash to compare: dummy hash if user not found, real hash if found
4. ALWAYS perform password verification
5. Check both conditions (user found AND password matched) after hash comparison

**Key Changes**:
```go
// Determine which hash to compare
hashToCompare := dummyPasswordHash
if userFound {
    hashToCompare = user.PasswordHash
}

// ALWAYS verify password (even if user not found)
match, err := s.hasher.VerifyPassword(password, hashToCompare)

// Check both conditions AFTER hash comparison
if !userFound || !match {
    return nil, errors.New("invalid username or password")
}
```

### Files Changed

1. **internal/service/auth/service.go**
   - Added `dummyPasswordHash` constant (precomputed Argon2id hash)
   - Added security comment explaining timing attack mitigation
   - Refactored Login method to use constant-time pattern (lines 256-302)
   - Track `userFound` boolean instead of returning early
   - Always perform password hash comparison
   - Check both conditions after hash comparison

**Dummy Hash**:
```go
const dummyPasswordHash = "$argon2id$v=19$m=65536,t=1,p=24$tQMNjFt979tvL7ho1P6xXw$DXkAY76TwLxFcMyqpMQQowtoWwhHfcs5Da9lFIid0Bg"
```

### Implementation Details

**Before (Vulnerable)**:
```go
user, err := s.repo.GetUserByUsername(ctx, username)
if err != nil {
    user, err = s.repo.GetUserByEmail(ctx, username)
    if err != nil {
        // FAST PATH: Return immediately without hash comparison
        return nil, errors.New("invalid username or password")
    }
}

// SLOW PATH: Only reached if user found
match, err := s.hasher.VerifyPassword(password, user.PasswordHash)
```

**After (Secure)**:
```go
user, err := s.repo.GetUserByUsername(ctx, username)
userFound := (err == nil)
if err != nil {
    user, err = s.repo.GetUserByEmail(ctx, username)
    userFound = (err == nil)
}

// Always compare hash (constant-time)
hashToCompare := dummyPasswordHash
if userFound {
    hashToCompare = user.PasswordHash
}

// ALWAYS verify password (takes ~50-100ms regardless of username validity)
match, err := s.hasher.VerifyPassword(password, hashToCompare)

// Check both conditions AFTER hash comparison
if !userFound || !match {
    return nil, errors.New("invalid username or password")
}
```

### Testing

**Manual Testing**:
- Measure login timing for valid username
- Measure login timing for invalid username
- Verify timing difference is minimal (within noise margin)
- Verify error messages are identical

**Integration Test Required**: ⚠️ PENDING
```go
func TestLogin_ConstantTime(t *testing.T) {
    // Test that login timing is similar for valid vs invalid usernames
    // Both should take ~50-100ms due to Argon2id comparison
}
```

**Benchmark Test**: Document expected timing behavior

### Impact

**Before**:
- Attackers could enumerate valid usernames via timing analysis
- Fast response (~1ms) = username doesn't exist
- Slow response (~50-100ms) = username exists
- Enables targeted brute-force attacks on known usernames

**After**:
- Login timing is constant (~50-100ms) regardless of username validity
- Attackers cannot determine username validity from response time
- Username enumeration via timing attacks is prevented
- Same error message for all failure cases

**Security Improvement**:
- Prevents reconnaissance phase of targeted attacks
- Forces attackers to guess both username AND password blindly
- Significantly increases attack difficulty and cost
- Protects user privacy (usernames not leaked)

**Performance Note**:
- Invalid username logins now take ~50-100ms instead of ~1ms
- This is acceptable security trade-off
- Legitimate users don't notice (they provide valid usernames)
- Rate limiting (A7.5) will prevent brute-force attempts

**Related Security Issues**: A7.5 (Account Lockout / Rate Limiting)

### References

- Source: [TODO_A7_SECURITY_FIXES.md](TODO_A7_SECURITY_FIXES.md) lines 174-223
- Report: [REPORT_2_IMPLEMENTATION_VERIFICATION.md](REPORT_2_IMPLEMENTATION_VERIFICATION.md)
- OWASP: [Authentication Timing Attacks](https://owasp.org/www-community/attacks/Timing_attack)

---

**Next Bug**: A7.3 - Goroutine Leaks in Notification Dispatcher

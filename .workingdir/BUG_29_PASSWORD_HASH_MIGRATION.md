# BUG #29: Password Hash Migration (bcrypt → argon2id)

**Discovered**: 2026-02-03 17:20 CET
**Severity**: HIGH
**Status**: IDENTIFIED
**Type**: Data Migration Issue

---

## Problem Description

Users with existing bcrypt passwords cannot log in after the crypto package was migrated from bcrypt to argon2id.

### Error Observed
```
level="warn" time="2026-02-03T16:16:36.858Z" logger="api" msg="Login failed"
error="password verification failed: failed to verify password: argon2id: hash is not in the correct format"
username="testuser"
```

### Root Cause
Database contains passwords hashed with bcrypt (`$2a$12$...`):
```sql
SELECT username, substring(password_hash, 1, 50) FROM shared.users WHERE username='testuser';
 username |                    hash_preview
----------+----------------------------------------------------
 testuser | $2a$12$6ee8jJHoXzrWnRUSfcpxkuLnhQegt2gNqybYQ8EqA4U
```

But `internal/crypto/password.go` now uses argon2id format (`$argon2id$v=19$...`).

---

## Impact

### Affected Components
- **Authentication**: All existing users cannot log in
- **Login endpoint**: `/api/v1/auth/login` returns 401
- **User migration**: Existing database entries incompatible

### User Impact
- **Severity**: HIGH - breaks authentication for all existing users
- **Workaround**: None (users locked out)
- **Data Loss**: No (passwords intact, just wrong format)

---

## Technical Analysis

### Code Changes
1. **Before** (implied from database): Used bcrypt
2. **After**: `internal/crypto/password.go` uses argon2id
   ```go
   func (h *PasswordHasher) VerifyPassword(password, hash string) (bool, error) {
       match, err := argon2id.ComparePasswordAndHash(password, hash)
       // This fails on bcrypt hashes
   }
   ```

### Why This Happened
- New crypto package introduced with argon2id
- No migration script provided
- Test data still uses old bcrypt format
- No backward compatibility check

---

## Solutions

### Option 1: Hybrid Password Verifier (Recommended)
Add backward compatibility to support both formats during migration.

```go
// VerifyPassword verifies a password against a hash (supports bcrypt and argon2id)
func (h *PasswordHasher) VerifyPassword(password, hash string) (bool, error) {
    if password == "" {
        return false, fmt.Errorf("password cannot be empty")
    }
    if hash == "" {
        return false, fmt.Errorf("hash cannot be empty")
    }

    // Check hash format
    if strings.HasPrefix(hash, "$2a$") || strings.HasPrefix(hash, "$2b$") || strings.HasPrefix(hash, "$2y$") {
        // bcrypt format - use bcrypt verification
        err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
        if err == bcrypt.ErrMismatchedHashAndPassword {
            return false, nil
        }
        if err != nil {
            return false, fmt.Errorf("bcrypt verification failed: %w", err)
        }
        return true, nil
    }

    // argon2id format - use argon2id verification
    match, err := argon2id.ComparePasswordAndHash(password, hash)
    if err != nil {
        return false, fmt.Errorf("argon2id verification failed: %w", err)
    }

    return match, nil
}
```

**Pros**:
- Immediate fix for existing users
- Graceful migration path
- No database changes needed
- Users auto-migrate on next login

**Cons**:
- Code complexity
- Need to maintain both algorithms temporarily
- Need bcrypt dependency

### Option 2: Database Migration Script
Force rehash all passwords (requires password reset).

```sql
-- Mark all users for password reset
UPDATE shared.users
SET password_hash = NULL,
    email_verified = false
WHERE password_hash LIKE '$2a$%';
```

**Pros**:
- Clean break from bcrypt
- Forces security upgrade
- Single algorithm in code

**Cons**:
- **BREAKS ALL LOGINS** - users must reset passwords
- User friction
- Requires email system working
- Not acceptable for production

### Option 3: Lazy Migration (Recommended for Production)
Combine Option 1 with automatic rehashing on successful login.

```go
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
    // ... existing code ...

    match, err := s.hasher.VerifyPassword(req.Password, user.PasswordHash)
    if err != nil {
        return nil, err
    }
    if !match {
        return nil, ErrInvalidCredentials
    }

    // Check if password needs rehashing (from bcrypt to argon2id)
    if strings.HasPrefix(user.PasswordHash, "$2a$") {
        newHash, err := s.hasher.HashPassword(req.Password)
        if err != nil {
            // Log error but don't fail login
            log.Warn().Err(err).Msg("Failed to rehash password")
        } else {
            // Update password hash in background
            go func() {
                err := s.repo.UpdatePasswordHash(context.Background(), user.ID, newHash)
                if err != nil {
                    log.Error().Err(err).Msg("Failed to update password hash")
                }
            }()
        }
    }

    // ... rest of login logic ...
}
```

**Pros**:
- Zero user impact
- Automatic migration over time
- No forced resets
- Gradual security upgrade

**Cons**:
- Takes time to migrate all users
- Need to monitor migration progress
- Some users might never migrate (inactive accounts)

---

## Recommended Implementation

### Phase 1: Immediate Fix (Option 1)
1. Add bcrypt support to `VerifyPassword`
2. Test with existing bcrypt and new argon2id hashes
3. Deploy immediately

### Phase 2: Lazy Migration (Option 3)
1. Add rehashing logic to login flow
2. Add repository method `UpdatePasswordHash`
3. Monitor migration progress with metrics

### Phase 3: Cleanup (Future)
1. After 3-6 months, check unmigrated users
2. Send migration emails to inactive users
3. Eventually remove bcrypt support

---

## Test Cases Needed

```go
func TestPasswordHasher_VerifyBcryptPassword(t *testing.T) {
    hasher := NewPasswordHasher()

    // bcrypt hash of "TestPass123!"
    bcryptHash := "$2a$12$6ee8jJHoXzrWnRUSfcpxkuLnhQegt2gNqybYQ8EqA4U"

    match, err := hasher.VerifyPassword("TestPass123!", bcryptHash)
    assert.NoError(t, err)
    assert.True(t, match)
}

func TestPasswordHasher_VerifyArgon2idPassword(t *testing.T) {
    hasher := NewPasswordHasher()

    // Create new argon2id hash
    hash, err := hasher.HashPassword("TestPass123!")
    require.NoError(t, err)

    match, err := hasher.VerifyPassword("TestPass123!", hash)
    assert.NoError(t, err)
    assert.True(t, match)
}

func TestAuthService_LazyPasswordMigration(t *testing.T) {
    // Test that bcrypt passwords get rehashed to argon2id on login
    // ...
}
```

---

## Metrics to Track

```go
prometheus.NewCounterVec(
    prometheus.CounterOpts{
        Name: "auth_password_hash_format_total",
        Help: "Password verifications by hash format",
    },
    []string{"format"}, // "bcrypt", "argon2id"
)

prometheus.NewCounter(
    prometheus.CounterOpts{
        Name: "auth_password_migrations_total",
        Help: "Number of passwords migrated from bcrypt to argon2id",
    },
)
```

---

## Prevention

1. **Migration Scripts**: Always provide migration scripts when changing data formats
2. **Backward Compatibility**: Maintain compatibility during transition periods
3. **Feature Flags**: Use feature flags for cryptographic changes
4. **Monitoring**: Track algorithm usage and migration progress
5. **Documentation**: Document breaking changes clearly

---

## Related Issues
- BUG_28: Password hashing inconsistency (documented earlier)
- This is the actual root cause of BUG_28

---

## Action Items

- [ ] Implement hybrid VerifyPassword (Option 1)
- [ ] Add bcrypt dependency to go.mod
- [ ] Write tests for both hash formats
- [ ] Implement lazy migration (Option 3)
- [ ] Add UpdatePasswordHash repository method
- [ ] Add migration metrics
- [ ] Update documentation
- [ ] Monitor migration progress

---

**Priority**: HIGH (blocks authentication)
**Estimated Effort**: 3-4 hours
**Risk**: Low (backward compatible fix)

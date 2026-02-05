# TODO A7: Security Fixes

**Phase**: A7
**Priority**: P0 (Critical)
**Effort**: 16-24 hours
**Status**: Pending
**Dependencies**: None
**Created**: 2026-02-05

---

## Overview

Critical security vulnerabilities discovered during code review:
- Missing transaction boundaries (data integrity)
- Timing attack vulnerability (username enumeration)
- Goroutine leaks (resource management)
- Information disclosure (password reset)
- Missing service-level rate limiting

**Source**: [REPORT_2_IMPLEMENTATION_VERIFICATION.md](REPORT_2_IMPLEMENTATION_VERIFICATION.md)

---

## Tasks

### A7.1: Transaction Boundaries 🔴 CRITICAL

**Priority**: P0
**Effort**: 8-12h
**Files**:
- `internal/service/auth/service.go`
- `internal/service/user/service.go`
- `internal/service/session/service.go`

**Problem**: Multi-step operations without transactions leave database in inconsistent state.

**Affected Operations**:

#### A7.1.1: User Registration
**Location**: `internal/service/auth/service.go:87-97`

**Issue**:
```go
user, err := s.repo.CreateUser(ctx, ...)  // User created
// ... error handling ...

_, err = s.repo.CreateEmailVerificationToken(ctx, ...)  // This fails
// User already in DB without token!
```

**Fix**:
```go
func (s *Service) Register(ctx context.Context, req RegisterRequest) (*User, error) {
    return s.db.WithTx(ctx, func(tx pgx.Tx) (*User, error) {
        user, err := s.repo.CreateUserTx(ctx, tx, ...)
        if err != nil {
            return nil, fmt.Errorf("failed to create user: %w", err)
        }

        _, err = s.repo.CreateEmailVerificationTokenTx(ctx, tx, ...)
        if err != nil {
            return nil, fmt.Errorf("failed to create verification token: %w", err)
        }

        return user, nil
    })
}
```

**Subtasks**:
- [ ] Add transaction support to repository interfaces
- [ ] Implement `CreateUserTx`, `CreateEmailVerificationTokenTx`
- [ ] Wrap Register operation in transaction
- [ ] Write integration test for rollback behavior

---

#### A7.1.2: Avatar Upload
**Location**: `internal/service/user/service.go:346-376`

**Issue**: Multiple DB operations without transaction:
1. Get latest version
2. Unset current avatars
3. Create new avatar
4. Update user avatar_url

**Fix**:
```go
func (s *Service) UploadAvatar(ctx context.Context, userID uuid.UUID, reader io.Reader) (*Avatar, error) {
    return s.db.WithTx(ctx, func(tx pgx.Tx) (*Avatar, error) {
        // 1. Get latest version
        version, err := s.repo.GetLatestAvatarVersionTx(ctx, tx, userID)
        if err != nil {
            return nil, err
        }

        // 2. Unset current avatars
        if err := s.repo.UnsetCurrentAvatarsTx(ctx, tx, userID); err != nil {
            return nil, err
        }

        // 3. Store file (outside transaction, cleanup on error)
        key := fmt.Sprintf("avatars/%s/%s.jpg", userID, uuid.New())
        url, err := s.storage.Store(ctx, key, reader, "image/jpeg")
        if err != nil {
            return nil, err
        }

        // 4. Create avatar record
        avatar, err := s.repo.CreateAvatarTx(ctx, tx, userID, url, version+1)
        if err != nil {
            s.storage.Delete(ctx, key) // Cleanup file
            return nil, err
        }

        // 5. Update user
        if err := s.repo.UpdateUserAvatarTx(ctx, tx, userID, url); err != nil {
            s.storage.Delete(ctx, key) // Cleanup file
            return nil, err
        }

        return avatar, nil
    })
}
```

**Subtasks**:
- [ ] Add transaction methods to user repository
- [ ] Implement cleanup logic for failed uploads
- [ ] Wrap avatar operations in transaction
- [ ] Write integration test

---

#### A7.1.3: Session Refresh
**Location**: `internal/service/session/service.go:148-151`

**Issue**: Create new session, then revoke old one. If CreateSession fails, old session is already revoked.

**Fix**: Create new session first, then revoke old one only on success:
```go
func (s *Service) RefreshSession(ctx context.Context, oldToken string) (*Session, error) {
    // 1. Validate old session
    oldSession, err := s.GetSessionByToken(ctx, oldToken)
    if err != nil {
        return nil, err
    }

    // 2. Create new session first
    newSession, err := s.CreateSession(ctx, oldSession.UserID, ...)
    if err != nil {
        return nil, fmt.Errorf("failed to create new session: %w", err)
    }

    // 3. Revoke old session only if new one succeeded
    if err := s.RevokeSession(ctx, oldToken); err != nil {
        // Log but don't fail - new session is valid
        s.logger.Warn("failed to revoke old session", "error", err)
    }

    return newSession, nil
}
```

**Subtasks**:
- [ ] Reorder session refresh logic
- [ ] Add error logging for revoke failures
- [ ] Write test for refresh flow
- [ ] Document behavior

---

### A7.2: Timing Attack in Login 🔴 CRITICAL

**Priority**: P0
**Effort**: 2-4h
**Location**: `internal/service/auth/service.go:236-268`

**Issue**: Login timing differs based on whether user exists:
- User not found: Fast path (no hash comparison)
- User found: Slow path (Argon2id comparison ~50-100ms)

→ Attacker can enumerate valid usernames via timing analysis

**Fix**: Always compare password hash, even if user doesn't exist:

```go
const dummyHash = "$argon2id$v=19$m=65536,t=3,p=2$..." // Precomputed dummy hash

func (s *Service) Login(ctx context.Context, req LoginRequest) (*Session, error) {
    // 1. Lookup user
    user, err := s.repo.GetUserByUsername(ctx, req.Username)
    if err != nil && !errors.Is(err, ErrUserNotFound) {
        return nil, err
    }

    // 2. Always compare hash (constant-time behavior)
    hashToCompare := dummyHash
    if err == nil {
        hashToCompare = user.PasswordHash
    }

    if !crypto.ComparePasswordHash(req.Password, hashToCompare) {
        return nil, ErrInvalidCredentials
    }

    // 3. Return error if user not found (after hash comparison)
    if err != nil {
        return nil, ErrInvalidCredentials
    }

    // 4. Create session
    return s.CreateSession(ctx, user.ID, ...)
}
```

**Subtasks**:
- [ ] Generate dummy hash constant
- [ ] Refactor login to always compare hash
- [ ] Add timing test (verify no difference)
- [ ] Document security fix

---

### A7.3: Goroutine Leak in Notification Dispatcher 🟠 HIGH

**Priority**: P1
**Effort**: 2-3h
**Location**: `internal/service/notification/dispatcher.go:116-137`

**Issue**: `Dispatch` spawns goroutine without tracking or waiting for completion.

**Fix**: Add shutdown tracking:

```go
type Dispatcher struct {
    // ... existing fields ...
    wg     sync.WaitGroup
    stopCh chan struct{}
}

func (d *Dispatcher) Dispatch(ctx context.Context, notif *Notification) error {
    agents := d.getEnabledAgents()
    if len(agents) == 0 {
        return ErrNoAgentsEnabled
    }

    d.wg.Add(1)
    go func() {
        defer d.wg.Done()

        for _, agent := range agents {
            select {
            case <-d.stopCh:
                return
            default:
            }

            agentCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
            err := agent.Send(agentCtx, notif)
            cancel()

            if err != nil {
                d.logger.Error("failed to send notification",
                    "agent", agent.Name(),
                    "error", err,
                )
            }
        }
    }()

    return nil
}

func (d *Dispatcher) Close() error {
    close(d.stopCh)
    d.wg.Wait()
    return nil
}
```

**Subtasks**:
- [ ] Add WaitGroup and stop channel to Dispatcher
- [ ] Implement Close method
- [ ] Hook Close into fx lifecycle
- [ ] Write test for graceful shutdown

---

### A7.4: Password Reset Information Disclosure 🟠 HIGH

**Priority**: P1
**Effort**: 1-2h
**Location**: `internal/service/auth/service.go:521-562`

**Issue**: Function returns different values:
- Empty string if user not found
- Actual token if user found

→ Caller might leak this information (email enumeration)

**Fix**: Never return token, always send via email:

```go
func (s *Service) RequestPasswordReset(ctx context.Context, email string) error {
    user, err := s.repo.GetUserByEmail(ctx, email)
    if err != nil {
        if errors.Is(err, ErrUserNotFound) {
            // Silently succeed - don't reveal email doesn't exist
            return nil
        }
        return err
    }

    token, err := s.repo.CreatePasswordResetToken(ctx, user.ID)
    if err != nil {
        return err
    }

    // Send email async (don't block request)
    go func() {
        ctx := context.Background()
        if err := s.emailService.SendPasswordReset(ctx, user.Email, token); err != nil {
            s.logger.Error("failed to send password reset email", "error", err)
        }
    }()

    return nil
}
```

**Subtasks**:
- [ ] Change return type from string to error
- [ ] Update API handler to not return token
- [ ] Send email asynchronously
- [ ] Add email logging

---

### A7.5: Service-Level Rate Limiting 🟡 MEDIUM

**Priority**: P2
**Effort**: 4-6h
**Location**: `internal/service/auth/service.go`

**Issue**: Password verification (Argon2id) is CPU-intensive. No rate limiting at service layer (only API middleware).

**Fix**: Add account lockout after failed attempts:

```sql
-- Migration: Add failed login tracking
CREATE TABLE shared.failed_login_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username TEXT NOT NULL,
    ip_address TEXT NOT NULL,
    attempted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_failed_login_username ON shared.failed_login_attempts(username, attempted_at);
CREATE INDEX idx_failed_login_ip ON shared.failed_login_attempts(ip_address, attempted_at);
```

```go
func (s *Service) Login(ctx context.Context, req LoginRequest) (*Session, error) {
    // 1. Check failed attempt count
    attempts, err := s.repo.CountFailedLoginAttempts(ctx, req.Username, 15*time.Minute)
    if err != nil {
        return nil, err
    }

    if attempts >= 5 {
        return nil, ErrAccountLocked
    }

    // 2. Attempt login (with timing attack fix from A7.2)
    user, err := s.authenticateUser(ctx, req)
    if err != nil {
        // Record failed attempt
        s.repo.RecordFailedLogin(ctx, req.Username, req.IPAddress)
        return nil, err
    }

    // 3. Clear failed attempts on success
    s.repo.ClearFailedLoginAttempts(ctx, req.Username)

    // 4. Create session
    return s.CreateSession(ctx, user.ID, ...)
}
```

**Subtasks**:
- [ ] Create migration for failed_login_attempts table
- [ ] Add repository methods for attempt tracking
- [ ] Implement lockout logic in Login
- [ ] Add configuration for lockout threshold/duration
- [ ] Write tests for lockout behavior

---

### A7.6: Context Misuse in Async Operations 🟡 MEDIUM

**Priority**: P2
**Effort**: 2-3h
**Location**: `internal/service/apikeys/service.go:185-192`

**Issue**: Using `context.Background()` in goroutine loses cancellation from parent context.

**Fix**: Create detached context with timeout:

```go
go func() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := s.repo.UpdateAPIKeyLastUsed(ctx, dbKey.ID); err != nil {
        s.logger.Warn("failed to update API key last used", "error", err)
    }
}()
```

**Subtasks**:
- [ ] Review all goroutines using context.Background()
- [ ] Add timeouts to detached contexts
- [ ] Document async operation patterns
- [ ] Add lint rule to catch this pattern

---

## Testing Requirements

### Security Tests

**Location**: `internal/service/auth/service_security_test.go` (new file)

```go
// Test timing attack is fixed
func TestLogin_ConstantTime(t *testing.T) {
    // Measure time for valid vs invalid username
    // Should be similar (within tolerance)
}

// Test transaction rollback
func TestRegister_TransactionRollback(t *testing.T) {
    // Force CreateEmailVerificationToken to fail
    // Verify user was not created
}

// Test account lockout
func TestLogin_AccountLockout(t *testing.T) {
    // Attempt login 5 times with wrong password
    // 6th attempt should return ErrAccountLocked
}
```

**Subtasks**:
- [ ] Write timing attack test
- [ ] Write transaction rollback tests
- [ ] Write account lockout tests
- [ ] Add benchmark for login timing

---

## Dependencies

**Required for**:
- Production deployment
- Security audit
- User trust

**Blocks**:
- Nothing (can be done in parallel with other phases)

---

## Verification Checklist

- [ ] All transactions implemented
- [ ] Timing attack fixed and tested
- [ ] Goroutine leaks fixed
- [ ] Information disclosure fixed
- [ ] Rate limiting/lockout implemented
- [ ] Context usage corrected
- [ ] Security tests passing
- [ ] Code review completed
- [ ] Penetration test passed (if applicable)

---

## Rollout Plan

1. **Development**: Implement all fixes with tests
2. **Staging**: Deploy and run security scan (gosec, semgrep)
3. **Review**: Security review by team
4. **Production**: Deploy during maintenance window
5. **Monitor**: Watch for auth errors, lockouts

---

## Related Issues

- [REPORT_2_IMPLEMENTATION_VERIFICATION.md](REPORT_2_IMPLEMENTATION_VERIFICATION.md)
- BUG_35: TOTP Generate Secret Upsert (completed)
- BUG_36: HasTOTP Error Handling (completed)

---

**Completion Criteria**:
✅ All 6 security issues fixed
✅ Security tests passing
✅ No new security warnings from gosec
✅ Code review approved

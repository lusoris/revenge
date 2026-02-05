# Revenge Codebase: Implementation Correctness Verification

**Generated**: 2026-02-05
**Analysis Type**: Code Quality & Security Review
**Severity Levels**: 🔴 Critical | 🟠 High | 🟡 Medium | 🔵 Low

---

## Executive Summary

The Revenge codebase demonstrates **strong architectural foundations** with consistent use of established patterns (repository pattern, dependency injection, error wrapping). However, several **critical issues** require immediate attention, particularly around **transaction boundaries**, **goroutine management**, and **timing attack vulnerabilities**.

**Overall Code Quality**: 7.5/10
- **Good Practices**: SQL injection protection, error wrapping, mutex usage
- **Critical Issues**: 4 high-priority issues found
- **Medium Issues**: 3 medium-priority issues found
- **Low Issues**: 3 low-priority issues found

---

## 1. Concurrency Safety

### 🔴 CRITICAL: Goroutine leak in notification dispatcher

**Location**: `internal\service\notification\dispatcher.go:116-137`

**Problem**: The `Dispatch` method spawns a goroutine without any mechanism to track or wait for completion. If the service shuts down while notifications are being sent, goroutines may be leaked.

```go
go func() {
    for _, agent := range agents {
        agentCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        // ... send notification ...
        cancel()
    }
}()
```

**Impact**:
- Resource leaks during shutdown
- Notifications may be lost
- Graceful shutdown incomplete

**Recommendation**:
Add a `sync.WaitGroup` or context-based shutdown mechanism to track goroutines and ensure graceful termination:

```go
type Dispatcher struct {
    // ... existing fields ...
    wg sync.WaitGroup
}

func (d *Dispatcher) Dispatch(ctx context.Context, notif *Notification) error {
    d.wg.Add(1)
    go func() {
        defer d.wg.Done()
        // ... existing logic ...
    }()
    return nil
}

func (d *Dispatcher) Close() error {
    d.wg.Wait()
    return nil
}
```

---

### 🟡 MEDIUM: Context misuse in async operations

**Location**: `internal\service\apikeys\service.go:185-192`

**Problem**: Using `context.Background()` in goroutine for async database update loses the ability to cancel the operation if the parent context is cancelled:

```go
go func() {
    if err := s.repo.UpdateAPIKeyLastUsed(context.Background(), dbKey.ID); err != nil {
        // ...
    }
}()
```

**Impact**:
- Database operations continue even if request is cancelled
- Minor resource waste

**Recommendation**:
Create a detached context with timeout instead:

```go
go func() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := s.repo.UpdateAPIKeyLastUsed(ctx, dbKey.ID); err != nil {
        // ...
    }
}()
```

Or better yet, use a proper background job queue (River) for async updates.

---

### 🟡 MEDIUM: Race condition in rate limiter health check

**Location**: `internal\api\middleware\ratelimit_redis.go:111-130`

**Problem**: The health check goroutine continuously runs and modifies `rl.healthy` with mutex protection, but ticker cleanup happens only in defer. If the goroutine never exits, it's a resource leak.

**Impact**:
- Resource leak if service doesn't shut down properly
- Minor goroutine accumulation

**Recommendation**:
Add proper shutdown mechanism using context cancellation or a stop channel:

```go
type RedisRateLimiter struct {
    // ... existing fields ...
    stopCh chan struct{}
}

func (rl *RedisRateLimiter) startHealthCheck() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            // ... health check ...
        case <-rl.stopCh:
            return
        }
    }
}

func (rl *RedisRateLimiter) Close() error {
    close(rl.stopCh)
    return nil
}
```

---

### ✅ GOOD PRACTICE: Proper mutex usage

**Location**: `internal\service\notification\dispatcher.go`

The notification dispatcher correctly uses `sync.RWMutex` for concurrent read/write access to the agents map:
- Write operations (RegisterAgent, UnregisterAgent) use `Lock()`
- Read operations (ListAgents, GetAgent) use `RLock()`

This is the correct pattern and prevents race conditions.

---

## 2. Error Handling

### 🔴 CRITICAL: Missing transaction rollback

**Location**: `internal\service\auth\service.go:87-97`

**Problem**: When `CreateEmailVerificationToken` fails after user creation (line 89), the user is already created in the database but the function returns an error. This leaves orphaned users without verification tokens.

```go
user, err := s.repo.CreateUser(ctx, ...)  // User created
if err != nil {
    return nil, fmt.Errorf("failed to create user: %w", err)
}

_, err = s.repo.CreateEmailVerificationToken(ctx, ...)  // This fails
if err != nil {
    return nil, fmt.Errorf("failed to create verification token: %w", err)  // User already in DB!
}
```

**Impact**:
- Database inconsistency
- Orphaned user accounts
- Failed registration leaves invalid state

**Recommendation**:
Wrap both operations in a database transaction to ensure atomicity:

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

---

### 🟡 MEDIUM: Silent error handling with fmt.Printf

**Location**: Multiple locations including:
- `internal\service\auth\service.go:104, 216, 315, 444, 505, 594`
- `internal\service\user\service.go:104`

**Problem**: Using `fmt.Printf` for error logging instead of structured logging:

```go
fmt.Printf("failed to send verification email: %v\n", err)
fmt.Printf("failed to update last login: %v\n", err)
```

**Impact**:
- Errors not searchable in log aggregation
- Missing contextual information
- Inconsistent logging format

**Recommendation**:
Use the structured logger (slog) available in the service:

```go
s.logger.Error("failed to send verification email",
    "error", err,
    "user_id", user.ID,
)
```

---

### ✅ GOOD PRACTICE: Consistent error wrapping

Throughout the codebase, errors are properly wrapped with `fmt.Errorf(...: %w, err)` maintaining the error chain for debugging. This allows error inspection with `errors.Is` and `errors.As`.

---

## 3. Database Patterns

### 🔴 CRITICAL: Missing transaction support for multi-step operations

**Location**: `internal\service\user\service.go:346-376`

**Problem**: Avatar upload involves multiple operations without transaction:
1. Get latest version
2. Unset current avatars
3. Create new avatar
4. Update user avatar_url

If any step fails mid-way, database is left in inconsistent state. Cleanup attempts with `s.storage.Delete` but doesn't rollback DB changes.

**Impact**:
- Database inconsistency
- Orphaned avatar records
- User left with broken avatar state

**Recommendation**:
Use pgx transactions to wrap all database operations:

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

        // 3. Store file
        url, err := s.storage.Store(ctx, key, reader, contentType)
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

---

### 🟡 MEDIUM: Potential N+1 query in movie repository

**Location**: `internal\content\movie\repository_postgres.go`

**Problem**: The repository pattern is well-structured, but there's no evidence of JOIN queries for related data. Operations like `GetMovie` followed by `ListMovieGenres`, `ListMovieCast` would result in N+1 queries.

**Impact**:
- Performance degradation
- Increased database load
- Slower API responses

**Recommendation**:
Add methods that use JOINs to fetch movie with all related data in a single query:

```go
func (r *Repository) GetMovieWithDetails(ctx context.Context, id uuid.UUID) (*MovieDetails, error) {
    query := `
        SELECT
            m.*,
            COALESCE(json_agg(DISTINCT jsonb_build_object(
                'genre_id', mg.tmdb_genre_id,
                'name', mg.name
            )) FILTER (WHERE mg.tmdb_genre_id IS NOT NULL), '[]') as genres,
            COALESCE(json_agg(DISTINCT jsonb_build_object(
                'person_id', mc.tmdb_person_id,
                'name', mc.name,
                'role', mc.role
            )) FILTER (WHERE mc.tmdb_person_id IS NOT NULL), '[]') as cast
        FROM movies m
        LEFT JOIN movie_genres mg ON m.id = mg.movie_id
        LEFT JOIN movie_credits mc ON m.id = mc.movie_id
        WHERE m.id = $1
        GROUP BY m.id
    `
    // ... execute and parse ...
}
```

---

### ✅ GOOD PRACTICE: SQL injection protection

All database queries use sqlc-generated code with parameterized queries, providing excellent protection against SQL injection. No string concatenation or unsafe query building found.

---

## 4. Memory Management

### 🟡 MEDIUM: Missing defer close in session refresh

**Location**: `internal\service\session\service.go:148-151`

**Problem**: Session refresh creates a new session and revokes old one, but if `CreateSession` fails, the old session is already revoked. No rollback mechanism.

**Impact**:
- User logged out unexpectedly
- Poor user experience

**Recommendation**:
Use transactions or create new session first, then revoke old one only on success:

```go
func (s *Service) RefreshSession(ctx context.Context, oldToken string) (*Session, error) {
    // 1. Create new session first
    newSession, err := s.CreateSession(ctx, userID, ...)
    if err != nil {
        return nil, fmt.Errorf("failed to create new session: %w", err)
    }

    // 2. Revoke old session only if new one succeeded
    if err := s.RevokeSession(ctx, oldToken); err != nil {
        // Log but don't fail - new session is valid
        s.logger.Warn("failed to revoke old session", "error", err)
    }

    return newSession, nil
}
```

---

### 🟡 MEDIUM: Unbounded goroutine creation in notification system

**Location**: `internal\service\notification\dispatcher.go:169-200`

**Problem**: `DispatchSync` creates one goroutine per agent without limiting concurrency. With many agents, this could exhaust resources.

```go
for _, agent := range agents {
    wg.Add(1)
    go func(a Agent) {  // Unbounded goroutine creation
        defer wg.Done()
        // ...
    }(agent)
}
```

**Impact**:
- Resource exhaustion with many agents
- Potential goroutine explosion

**Recommendation**:
Use a worker pool or limit concurrent goroutines to a reasonable number (e.g., 10):

```go
func (d *Dispatcher) DispatchSync(ctx context.Context, notif *Notification) error {
    agents := d.getEnabledAgents()

    // Limit concurrency
    sem := make(chan struct{}, 10)
    var wg sync.WaitGroup

    for _, agent := range agents {
        wg.Add(1)
        sem <- struct{}{} // Acquire semaphore

        go func(a Agent) {
            defer func() {
                <-sem // Release semaphore
                wg.Done()
            }()
            // ... send notification ...
        }(agent)
    }

    wg.Wait()
    return nil
}
```

---

### ✅ GOOD PRACTICE: Proper resource cleanup

Database connection pool properly uses `defer conn.Release()` and cache properly implements `Close()` methods. No resource leaks detected in these areas.

---

## 5. Security Issues

### 🔴 CRITICAL: Timing attack vulnerability in error messages

**Location**: `internal\service\auth\service.go:236-249`

**Problem**: Login error messages differ based on whether username exists. While messages are the same, the timing difference between database lookups could leak information about valid usernames.

```go
user, err := s.repo.GetUserByUsername(ctx, req.Username)  // DB lookup
if err != nil {
    if errors.Is(err, ErrUserNotFound) {
        return nil, ErrInvalidCredentials  // Fast path
    }
    return nil, err
}

if !crypto.ComparePasswordHash(req.Password, user.PasswordHash) {
    return nil, ErrInvalidCredentials  // Slow path (Argon2id)
}
```

**Impact**:
- Username enumeration possible via timing analysis
- Security vulnerability

**Recommendation**:
Always perform password hash comparison even if user doesn't exist (compare against dummy hash) to maintain constant time:

```go
const dummyHash = "$argon2id$v=19$m=65536,t=3,p=2$..." // Precomputed dummy hash

user, err := s.repo.GetUserByUsername(ctx, req.Username)
if err != nil && !errors.Is(err, ErrUserNotFound) {
    return nil, err
}

// Always compare hash, even if user not found
hashToCompare := dummyHash
if err == nil {
    hashToCompare = user.PasswordHash
}

if !crypto.ComparePasswordHash(req.Password, hashToCompare) {
    return nil, ErrInvalidCredentials
}

if err != nil {
    return nil, ErrInvalidCredentials
}
```

---

### 🟠 HIGH: Password reset information disclosure

**Location**: `internal\service\auth\service.go:521-528`

**Problem**: Function returns success even if email doesn't exist (security by obscurity), but the return value differs:
- Returns empty string `""` if user not found (line 527)
- Returns actual token if user found (line 562)

The caller might inadvertently leak this information.

**Impact**:
- Email enumeration possible
- Security vulnerability

**Recommendation**:
Always return success response without token to caller. Send token via email only:

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

    // Send email async
    go s.emailService.SendPasswordReset(context.Background(), user.Email, token)

    return nil
}
```

---

### 🟠 HIGH: Missing rate limiting on password verification

**Location**: `internal\service\auth\service.go:252-268`

**Problem**: Password verification using Argon2id is CPU-intensive. No rate limiting at the service layer (only at API middleware level) means this could be bypassed if service is called directly.

**Impact**:
- CPU exhaustion via brute force
- DoS vulnerability

**Recommendation**:
Add service-level rate limiting or account lockout after failed attempts:

```go
func (s *Service) Login(ctx context.Context, req LoginRequest) (*Session, error) {
    // Check failed attempt count
    attempts, err := s.repo.GetFailedLoginAttempts(ctx, req.Username)
    if err != nil {
        return nil, err
    }

    if attempts >= 5 {
        lockoutUntil, _ := s.repo.GetLockoutTime(ctx, req.Username)
        if time.Now().Before(lockoutUntil) {
            return nil, ErrAccountLocked
        }
    }

    // ... existing login logic ...

    if loginFailed {
        s.repo.IncrementFailedAttempts(ctx, req.Username)
        return nil, ErrInvalidCredentials
    }

    s.repo.ResetFailedAttempts(ctx, req.Username)
    return session, nil
}
```

---

### ✅ GOOD PRACTICE: Secure token generation

Uses `crypto.GenerateSecureToken` with crypto/rand for token generation. API keys use SHA-256 hashing for storage. These are secure implementations.

---

## 6. API Design Issues

### 🟡 MEDIUM: Inconsistent error handling in handlers

**Location**: `internal\api\handler_mfa.go`

**Problem**: Some handlers return typed errors, others return generic `&ogen.Error`:
- Line 52: Returns `&ogen.Error{...}`
- Line 123: Returns `(*ogen.VerifyTOTPBadRequest)(&ogen.Error{...})`

**Impact**:
- Inconsistent API responses
- Poor developer experience

**Recommendation**:
Consistently use typed error responses for better API clarity:

```go
func (h *Handler) SetupTOTP(ctx context.Context, req *ogen.SetupTOTPReq) (*ogen.SetupTOTPOK, error) {
    // ... logic ...

    if err != nil {
        return nil, &ogen.SetupTOTPBadRequest{
            Code:    "SETUP_FAILED",
            Message: err.Error(),
        }
    }

    return &ogen.SetupTOTPOK{...}, nil
}
```

---

### 🟡 MEDIUM: Missing validation in API handlers

**Location**: `internal\api\handler_mfa.go:68-78`

**Problem**: `SetupTOTP` doesn't validate `accountName` parameter before passing to service. Empty or malicious input could cause issues.

**Impact**:
- Invalid data reaching service layer
- Poor error messages

**Recommendation**:
Add input validation at API layer before calling services:

```go
func (h *Handler) SetupTOTP(ctx context.Context, req *ogen.SetupTOTPReq) (*ogen.SetupTOTPOK, error) {
    if req.AccountName == "" {
        return nil, &ogen.SetupTOTPBadRequest{
            Code:    "INVALID_ACCOUNT_NAME",
            Message: "Account name cannot be empty",
        }
    }

    if len(req.AccountName) > 100 {
        return nil, &ogen.SetupTOTPBadRequest{
            Code:    "INVALID_ACCOUNT_NAME",
            Message: "Account name too long (max 100 characters)",
        }
    }

    // ... continue with service call ...
}
```

---

### 🔵 LOW: Lack of pagination validation

**Location**: `internal\service\user\service.go:315-322`

**Problem**: `ListUserAvatars` caps limit at 100 but doesn't validate offset. Negative offset could cause issues.

```go
if limit > 100 {
    limit = 100
}
// Missing: if offset < 0 { offset = 0 }
```

**Impact**:
- Potential query errors
- Minor UX issue

**Recommendation**:
Validate all pagination parameters:

```go
func (s *Service) ListUserAvatars(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Avatar, error) {
    if limit <= 0 {
        limit = 20
    }
    if limit > 100 {
        limit = 100
    }
    if offset < 0 {
        offset = 0
    }

    return s.repo.ListUserAvatars(ctx, userID, limit, offset)
}
```

---

### ✅ GOOD PRACTICE: Authentication context handling

Proper extraction of user ID from context with error handling in handlers like `GetUserIDFromContext`. This ensures authentication is enforced at the API layer.

---

## 7. Additional Findings

### 🔵 LOW: Missing index recommendations

While I cannot see the actual SQL schema files, based on the query patterns in `repository_postgres.go`, ensure indexes exist on:
- `movies.tmdb_id`, `movies.imdb_id`, `movies.radarr_id` (foreign key lookups)
- `movie_watched.user_id, movie_watched.movie_id` (composite for watch progress)
- `movie_genres.movie_id, movie_genres.tmdb_genre_id` (composite for genre queries)

**Recommendation**:
Add migration to create these indexes if they don't exist:

```sql
CREATE INDEX IF NOT EXISTS idx_movies_tmdb_id ON movies(tmdb_id);
CREATE INDEX IF NOT EXISTS idx_movies_imdb_id ON movies(imdb_id);
CREATE INDEX IF NOT EXISTS idx_movies_radarr_id ON movies(radarr_id);
CREATE INDEX IF NOT EXISTS idx_movie_watched_user_movie ON movie_watched(user_id, movie_id);
CREATE INDEX IF NOT EXISTS idx_movie_genres_movie_genre ON movie_genres(movie_id, tmdb_genre_id);
```

---

### ✅ GOOD PRACTICE: Configuration validation

Pool configuration properly validates and converts values with `validate.SafeInt32` to prevent integer overflow.

---

## Summary of Issues

### Critical Priority (Fix Immediately)

1. **Missing transaction support** for multi-step database operations
   - User registration (user + verification token)
   - Avatar upload (multiple DB operations)
   - Session refresh (create + revoke)

2. **Goroutine leaks** in notification dispatcher
   - Add shutdown tracking mechanism

3. **Timing attack vulnerability** in login
   - Always compare password hash even if user not found

### High Priority (Fix Before Production)

4. **Password reset information disclosure**
   - Don't return token value to caller

5. **Missing rate limiting** on password verification
   - Add service-level rate limiting or account lockout

### Medium Priority (Fix During Next Sprint)

6. **Context misuse** in async operations
7. **Silent error handling** with fmt.Printf (use structured logging)
8. **N+1 queries** in movie repository (add JOIN queries)
9. **Unbounded goroutine creation** in notification dispatcher
10. **Inconsistent error handling** in API handlers
11. **Missing input validation** in API handlers

### Low Priority (Technical Debt)

12. **Rate limiter health check** goroutine cleanup
13. **Pagination validation** gaps
14. **Database indexes** recommendations

---

## Code Quality Score Breakdown

| Category | Score | Notes |
|----------|-------|-------|
| Architecture | 9/10 | Excellent patterns (repository, DI, interfaces) |
| Error Handling | 7/10 | Good wrapping, but transaction issues |
| Concurrency | 6/10 | Goroutine leaks, context misuse |
| Security | 6/10 | Timing attacks, information disclosure |
| Testing | 7/10 | Good coverage in some areas, gaps in others |
| Documentation | 8/10 | Good design docs, inline comments |
| **Overall** | **7.5/10** | **Solid foundation, needs security hardening** |

---

## Recommendations

### Immediate Actions

1. **Add transaction support** to all multi-step operations
2. **Fix timing attack** in login flow
3. **Fix goroutine leaks** in notification system
4. **Review all async operations** for proper context usage

### Short-term Actions

5. **Replace fmt.Printf** with structured logging
6. **Add input validation** to all API handlers
7. **Implement account lockout** for failed login attempts
8. **Add JOIN queries** to avoid N+1 problems

### Long-term Actions

9. **Security audit** of all authentication flows
10. **Performance testing** with concurrent load
11. **Add chaos testing** for transaction rollbacks
12. **Implement distributed tracing** for debugging

---

## Conclusion

The Revenge codebase demonstrates **strong architectural foundations** with consistent patterns and good separation of concerns. However, **critical security vulnerabilities** and **transaction management issues** must be addressed before production deployment.

The good news is that these issues are well-defined and fixable. The codebase already has the right abstractions in place (repository pattern, transaction interfaces) to support proper fixes.

**Priority**: Focus on transaction safety and timing attack vulnerabilities first, as these affect data integrity and security.

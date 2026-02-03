# Session Service Bugs

## Bug 1: ValidateSession/RefreshSession don't return ErrUnauthorized for invalid tokens

**Location**: `internal/service/session/service.go`

**Problem**:
- `ValidateSession()` und `RefreshSession()` geben "failed to get session: no rows in result set" zurück statt `errors.ErrUnauthorized`
- GetSessionByTokenHash/GetSessionByRefreshTokenHash returnen bei pgx "no rows in result set" statt `sql.ErrNoRows`
- Repository-Code prüft nur auf `sql.ErrNoRows`, aber pgx gibt einen anderen Error zurück

**Expected**:
```go
_, err := service.ValidateSession(ctx, "invalid")
// Should return: errors.ErrUnauthorized
```

**Actual**:
```go
// Returns: "failed to get session: no rows in result set"
```

**Fix**: Repository muss auch pgx.ErrNoRows checken, nicht nur sql.ErrNoRows

## Bug 2: RefreshSession duplicate key violation on refresh_token_hash

**Location**: `internal/service/session/service.go:176`

**Problem**:
RefreshSession() versucht den alten refresh_token_hash wiederzuverwenden nach dem die alte Session revoked wurde:

```go
_, err = s.repo.CreateSession(ctx, CreateSessionParams{
    UserID:           session.UserID,
    TokenHash:        newTokenHash,
    RefreshTokenHash: &refreshTokenHash, // ❌ Reuse refresh token - DUPLICATE KEY!
    ...
})
```

Die alte Session wird revoked (nicht deleted), daher bleibt der refresh_token_hash in der DB mit unique constraint.

**Error**:
```
ERROR: duplicate key value violates unique constraint "sessions_refresh_token_hash_key"
```

**Fix**: Generate new refresh token:

```go
// Generate new refresh token
newRefreshToken, newRefreshTokenHash, err := s.generateToken()
if err != nil {
    return "", fmt.Errorf("failed to generate new refresh token: %w", err)
}

_, err = s.repo.CreateSession(ctx, CreateSessionParams{
    UserID:           session.UserID,
    TokenHash:        newTokenHash,
    RefreshTokenHash: &newRefreshTokenHash,  // ✅ New refresh token
    ...
})

// Return both new tokens
return newToken, newRefreshToken, nil  // Signature needs to change!
```

**Impact**: RefreshSession ist aktuell komplett broken und kann nicht verwendet werden.

**Alternative Fix**: Delete old session instead of revoking, oder unique constraint nur auf non-revoked sessions.

# Bug #36: HasTOTP Swallows Database Errors

**Status**: OPEN
**Date**: 2026-02-04
**Severity**: LOW
**Component**: internal/service/mfa/totp.go
**Found by**: Code Review

## Summary

The `HasTOTP` function returns `false, nil` for ALL errors instead of distinguishing between "no rows found" and actual database errors.

## Problem Code

`totp.go:199-207`:
```go
// HasTOTP checks if a user has TOTP configured
func (s *TOTPService) HasTOTP(ctx context.Context, userID uuid.UUID) (bool, error) {
    _, err := s.queries.GetUserTOTPSecret(ctx, userID)
    if err != nil {
        // SQLC doesn't define a specific error, check for no rows
        return false, nil  // BUG: All errors are swallowed!
    }
    return true, nil
}
```

## Expected Behavior

```go
func (s *TOTPService) HasTOTP(ctx context.Context, userID uuid.UUID) (bool, error) {
    _, err := s.queries.GetUserTOTPSecret(ctx, userID)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return false, nil
        }
        return false, fmt.Errorf("failed to check TOTP: %w", err)
    }
    return true, nil
}
```

## Impact

- Database errors (connection timeout, etc.) are interpreted as "no TOTP"
- Can lead to incorrect security decisions
- Low severity since real DB errors would cause other operations to fail too

## Test Coverage

Current test `TestTOTPService_HasTOTP` only tests the happy path.

New test needed:
```go
t.Run("database error", func(t *testing.T) {
    // Simulate database error
    // Verify error is propagated, not swallowed
})
```

## Files to Modify

- `internal/service/mfa/totp.go` - HasTOTP function
- `internal/service/mfa/totp_test.go` - Add error handling test

## Related

- Use `pgx.ErrNoRows` for "no rows found" check
- Use `errors.Is()` for correct error comparison (consistent with other services)

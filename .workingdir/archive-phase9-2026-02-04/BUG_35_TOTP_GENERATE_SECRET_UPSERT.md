# Bug #35: TOTP GenerateSecret Missing Upsert Logic

**Status**: OPEN
**Date**: 2026-02-04
**Severity**: MEDIUM
**Component**: internal/service/mfa/totp.go
**Found by**: Test `TestTOTPService_MultipleSecrets`

## Summary

When a user already has a TOTP secret and wants to generate a new one (e.g., re-enrollment), `GenerateSecret` fails with a unique constraint error.

## Symptom

```
Error: failed to store TOTP secret: ERROR: duplicate key value violates unique constraint "user_totp_secrets_pkey" (SQLSTATE 23505)
```

## Root Cause

`totp.go:GenerateSecret` only calls `CreateTOTPSecret` (INSERT) without checking if a secret already exists.

**Problem Code** (`totp.go:92-100`):
```go
// Store encrypted secret in database
_, err = s.queries.CreateTOTPSecret(ctx, db.CreateTOTPSecretParams{
    UserID:          userID,
    EncryptedSecret: encryptedSecret,
    Nonce:           encryptedSecret[:12], // First 12 bytes are the nonce
})
```

There's already an `UpdateTOTPSecret` query in `mfa.sql:21-29` that should be used for re-enrollment.

## Expected Behavior

1. Check if user already has a TOTP secret
2. If yes: use `UpdateTOTPSecret` (invalidates old secret)
3. If no: use `CreateTOTPSecret`

## Fix Required

```go
// Check if user already has TOTP secret
_, err = s.queries.GetUserTOTPSecret(ctx, userID)
if err == nil {
    // User has existing secret, update it
    err = s.queries.UpdateTOTPSecret(ctx, db.UpdateTOTPSecretParams{
        UserID:          userID,
        EncryptedSecret: encryptedSecret,
        Nonce:           encryptedSecret[:12],
    })
} else {
    // No existing secret, create new one
    _, err = s.queries.CreateTOTPSecret(ctx, db.CreateTOTPSecretParams{
        UserID:          userID,
        EncryptedSecret: encryptedSecret,
        Nonce:           encryptedSecret[:12],
    })
}
```

## Test Case

`TestTOTPService_MultipleSecrets` - Tests re-enrollment scenario

## Impact

- Users cannot re-enroll TOTP without manual deletion first
- Poor UX when switching devices or losing authenticator
- Medium severity since workaround exists (DeleteTOTP → GenerateSecret)

## Files to Modify

- `internal/service/mfa/totp.go` - GenerateSecret function

## Related

- SQL Query `UpdateTOTPSecret` already exists
- Test `TestTOTPService_MultipleSecrets` verifies fix

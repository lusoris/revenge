# Bug #34: gosec Security Issues

**Status**: RESOLVED
**Date**: 2026-02-04
**Severity**: HIGH
**Component**: Multiple

## Summary

gosec security scanner found 9 HIGH severity issues:
- 6x Integer overflow conversions (webauthn.go)
- 1x Integer overflow conversion (webhook.go)
- 2x TLS InsecureSkipVerify (email.go) - INTENTIONAL CONFIG

## Issues Found

### 1. WebAuthn SignCount Integer Overflows (6 instances)

**File**: `internal/service/mfa/webauthn.go`
**Lines**: 138, 198, 304, 362, 387, 409
**Issue**: G115 (CWE-190) - Integer overflow conversion between int32 and uint32
**Risk**: Potential integer overflow when converting sign counters

**Instances**:
- L138: `SignCount: uint32(cred.SignCount)` (int32 → uint32)
- L198: `SignCount: uint32(cred.SignCount)` (int32 → uint32)
- L304: `SignCount: uint32(cred.SignCount)` (int32 → uint32)
- L362: `SignCount: uint32(cred.SignCount)` (int32 → uint32)
- L387: `oldCounter := uint32(dbCred.SignCount)` (int32 → uint32)
- L409: `SignCount: int32(newCounter)` (uint32 → int32)

### 2. Webhook Exponential Backoff Integer Overflow

**File**: `internal/service/notification/agents/webhook.go`
**Line**: 166
**Issue**: G115 (CWE-190) - Integer overflow conversion int → uint
**Code**: `backoff := time.Duration(1<<uint(attempt-1)) * time.Second`
**Risk**: Potential overflow if attempt is large

### 3. TLS InsecureSkipVerify (INTENTIONAL)

**File**: `internal/service/notification/agents/email.go`
**Lines**: 325, 348
**Issue**: G402 (CWE-295) - TLS InsecureSkipVerify may be true
**Status**: INTENTIONAL - User-configurable option for self-signed certificates
**Action**: Add #nosec comment with justification

## Fixes Required

1. Add safe conversion helpers for SignCount (validate range before conversion)
2. Cap webhook retry attempts to prevent overflow
3. Add #nosec comments to TLS skip verify with justification

## Resolution

**All 9 issues RESOLVED:**

1. **WebAuthn SignCount overflows** (6 instances): Added safe conversion helpers
   - `safeUint32ToInt32()` - caps at max int32
   - `safeInt32ToUint32()` - treats negative as 0
   - Applied to all conversions in webauthn.go

2. **Webhook backoff overflow**: Capped attempt value before conversion
   - Added `safeAttempt := min(attempt-1, 6)` before shift operation
   - Added #nosec G115 comment with justification

3. **TLS InsecureSkipVerify** (2 instances): Added #nosec G402 comments
   - Documented as user-configurable for self-signed certificates
   - email.go lines 325, 348

**Verification**: `gosec` scan passes with 0 HIGH severity issues in fixed files

# Gosec G115 Helper Functions - Duplication Report

**Date**: 2026-02-04
**Purpose**: Identify duplicated gosec G115 safe conversion helpers that should use the shared library

---

## Summary

Found **1 file** with duplicated local helpers that should use the shared `internal/util/safeconv.go` library.

---

## Shared Libraries

### 1. `internal/util/safeconv.go` - Capping/Saturation Approach

The primary shared library for safe integer conversions. Uses **capping** (saturation) - values exceeding bounds are capped to max/min.

```go
SafeIntToInt32(v int) int32
SafeInt64ToInt32(v int64) int32
SafeUint64ToInt32(v uint64) int32
SafeInt32ToUint32(v int32) uint32    // negatives → 0
SafeUint32ToInt32(v uint32) int32
SafeIntToUint(v int) uint            // negatives → 0
```

**Use case**: When overflow should be handled silently by capping to bounds.

### 2. `internal/validate/convert.go` - Error-Returning Approach

A different library for validation contexts. Returns **errors** on overflow.

```go
SafeInt32(value int) (int32, error)
SafeUint32(value int) (uint32, error)
SafeUint(value int) (uint, error)
MustInt32(value int) int32           // panics on overflow
MustUint32(value int) uint32         // panics on overflow
MustUint(value int) uint             // panics on overflow
```

**Use case**: When overflow indicates a bug/invalid input that should be reported.

---

## Duplicates Found

### `internal/service/mfa/webauthn.go` (Lines 38-52)

**Status**: ⚠️ DUPLICATE - Should use `util/safeconv.go`

Local functions that duplicate the shared library:

```go
// Line 38-45: Duplicate of util.SafeUint32ToInt32
func safeUint32ToInt32(val uint32) int32 {
    const maxInt32 = 2147483647
    if val > maxInt32 {
        return maxInt32
    }
    return int32(val) // #nosec G115 -- validated above
}

// Line 47-52: Duplicate of util.SafeInt32ToUint32
func safeInt32ToUint32(val int32) uint32 {
    if val < 0 {
        return 0
    }
    return uint32(val) // #nosec G115 -- validated above
}
```

**Used at**:
- Line 295: `int32(safeUint32ToInt32(cred.ID[i]))`
- Line 311: `int32(safeUint32ToInt32(cred.PublicKey[i]))`
- Line 347: `byte(safeInt32ToUint32(b))`
- Line 360: `byte(safeInt32ToUint32(b))`

**Recommendation**: Replace with:
```go
import "github.com/lusoris/revenge/internal/util"

// Replace safeUint32ToInt32 calls with:
util.SafeUint32ToInt32(cred.ID[i])

// Replace safeInt32ToUint32 calls with:
util.SafeInt32ToUint32(b)
```

---

## Acceptable Inline Handling

### `internal/service/notification/agents/webhook.go` (Line 167-168)

**Status**: ✅ OK - Inline context-specific capping

```go
safeAttempt := min(attempt-1, 6) // Max 2^6 = 64s
backoffSeconds := 1 << uint(safeAttempt) // #nosec G115 -- safeAttempt is capped at 6
```

This is acceptable because:
1. It's a simple one-liner specific to this context
2. The capping value (6) is business logic specific to retry backoff
3. Creating a generic helper for this wouldn't add value

---

## Action Items

| Priority | File | Action | Status |
|----------|------|--------|--------|
| P1 | `internal/service/mfa/webauthn.go` | Remove local `safeUint32ToInt32` and `safeInt32ToUint32`, use `util.SafeUint32ToInt32` and `util.SafeInt32ToUint32` | ✅ FIXED |

---

## When to Use Which Library

| Scenario | Library | Example |
|----------|---------|---------|
| Silent capping is acceptable | `util/safeconv.go` | Database field conversion, API response mapping |
| Overflow is a bug/invalid input | `validate/convert.go` | User input validation, config parsing |
| Context-specific inline logic | Inline with `#nosec G115` | Retry backoff calculation |

---

## Files Checked

All files containing `#nosec G115` or `func safe[A-Z]` patterns were reviewed:

- `internal/util/safeconv.go` - ✅ Shared library
- `internal/validate/convert.go` - ✅ Different purpose (error-returning)
- `internal/service/mfa/webauthn.go` - ⚠️ Has duplicates
- `internal/service/notification/agents/webhook.go` - ✅ Inline OK
- `.workingdir/archive-*` - Historical documentation only

# Bug #32: mediainfo.go DispositionFlag Compile Error

**Date**: 2026-02-04
**Found During**: mockery mock generation
**Severity**: HIGH
**Status**: 🔍 INVESTIGATING

## Description
mockery fails to parse `internal/content/movie/mediainfo.go` due to invalid `DispositionFlagDefault` constant usage:

```
error="C:\\Users\\ms\\dev\\revenge\\internal\\content\\movie\\mediainfo.go:214:54:
cannot use astiav.DispositionFlagDefault (constant unknown with invalid type) as astiav.DispositionFlag value
in argument to stream.DispositionFlags().Has"
```

## Affected Code
**File**: `internal/content/movie/mediainfo.go:214, 244`

```go
// Line 214
audioInfo.IsDefault = stream.DispositionFlags().Has(astiav.DispositionFlagDefault)

// Line 244
subInfo.IsDefault = disposition.Has(astiav.DispositionFlagDefault)
```

## Investigation
- `go build ./internal/content/movie/` succeeds (no error)
- mockery fails to parse the file
- go-astiav version: v0.40.0
- Possible causes:
  1. mockery uses different parser than go compiler
  2. CGO constants not available to mockery
  3. `DispositionFlagDefault` doesn't exist or has wrong type in go-astiav v0.40.0

## Potential Fixes

### Option A: Check FFmpeg/astiav documentation
Look up correct constant name in go-astiav v0.40.0

### Option B: Fallback to bitwise check
```go
// Instead of:
audioInfo.IsDefault = stream.DispositionFlags().Has(astiav.DispositionFlagDefault)

// Use:
audioInfo.IsDefault = (stream.DispositionFlags() & astiav.AV_DISPOSITION_DEFAULT) != 0
```

### Option C: Check if constant exists
```go
const (
    // Check if this exists in go-astiav
    DispositionFlagDefault DispositionFlag = C.AV_DISPOSITION_DEFAULT
)
```

## Next Steps
1. Check go-astiav v0.40.0 source for correct constant name
2. Check FFmpeg 8.0 documentation for AV_DISPOSITION_DEFAULT
3. Apply fix
4. Verify mockery can parse
5. Re-run mock generation

---

## Root Cause
mockery cannot parse CGO constants. `DispositionFlagDefault` is defined as:
```go
DispositionFlagDefault = DispositionFlag(C.AV_DISPOSITION_DEFAULT)
```

The `C.AV_DISPOSITION_DEFAULT` requires CGO to be executed, but mockery only does static parsing.

## Resolution
**Decision**: Exclude movie package from mockery generation for now.
- Phase 1 focuses on Session/Auth/User/RBAC Services
- Movie package tests will use real implementations or manual mocks
- Not a blocker for Phase 1 testing

**Alternative for later**:
- Run mockery with CGO_ENABLED=1 and proper build tags
- Create manual mocks for movie interfaces
- Or use real implementations in tests (integration testing)

**Status**: ✅ RESOLVED (excluded from mockery, not a Phase 1 blocker)

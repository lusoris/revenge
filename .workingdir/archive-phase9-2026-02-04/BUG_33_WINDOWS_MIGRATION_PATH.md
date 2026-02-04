# Bug #33: Windows Migration File Path Parsing Error

**Status**: RESOLVED
**Date**: 2026-02-04
**Severity**: HIGH
**Component**: `internal/testutil/testdb_migrate.go`

## Problem

Integration tests failing on Windows with file path URL parsing error:

```
failed to create template database: failed to run migrations on template:
failed to open source, "file://c:\\Users\\ms\\dev\\revenge\\migrations":
parse "file://c:\\Users\\ms\\dev\\revenge\\migrations":
invalid port ":\\Users\\ms\\dev\\revenge\\migrations" after host
```

## Root Cause

The migrations path on Windows (e.g., `c:\Users\ms\dev\revenge\migrations`) contains a colon after the drive letter. When prefixed with `file://`, the URL parser interprets the colon as a port separator, causing the parsing to fail.

Expected: `file:///c:/Users/ms/dev/revenge/migrations` (note the triple slash and forward slashes)
Actual: `file://c:\Users\ms\dev\revenge\migrations` (double slash and backslashes)

## Affected Tests

All integration tests using `testutil.NewTestDB()`:
- `TestService_CreateSession_MaxPerUser`
- `TestService_ValidateSession`
- `TestService_RefreshSession`
- `TestRepositoryPG_*` (all repository integration tests)
- And more...

## Fix Required

In `internal/testutil/testdb.go`, the migration path construction needs to:
1. Convert Windows backslashes to forward slashes
2. Use `file:///` (triple slash) for absolute paths on Windows
3. Properly handle the drive letter in the file:// URL scheme

## Impact

- **Unit tests with mocks**: ✅ PASS (20/20 exhaustive tests passed)
- **Integration tests**: ❌ FAIL (all tests using testDB)
- **Coverage reporting**: Blocked until fixed

## Resolution

Fixed in `internal/testutil/testdb_migrate.go` by using Go's `net/url.URL` to properly construct file:// URLs:

```go
func pathToFileURL(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	u := &url.URL{
		Scheme: "file",
		Path:   filepath.ToSlash(absPath),
	}

	return u.String()
}
```

This correctly handles:
- Windows absolute paths (C:\Users\... → file:///C:/Users/...)
- Unix absolute paths (/home/... → file:///home/...)
- Proper URL encoding via net/url package

## Verification

Test `TestService_CreateSession_MaxPerUser` now passes: ✅ PASS (11.26s)

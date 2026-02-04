# Bug Tracker

**Purpose**: Track bugs discovered during implementation

---

## Open Bugs

_None yet_

---

## Fixed Bugs

_None yet_

---

## Format

```markdown
### BUG-XXX: Short Description
**Severity**: Critical/High/Medium/Low
**Found**: YYYY-MM-DD
**File**: path/to/file.go:line
**Description**: What's wrong
**Fix**: How it was fixed (when resolved)
**Status**: Open/Fixed
```

## A0.4: Avatar Upload Implementation - COMPLETED

**Status**: Fixed
**Date**: 2026-02-04

### Changes Made:
1. Created `internal/service/storage/` package with:
   - `Storage` interface for abstraction (clustering-ready)
   - `LocalStorage` for local filesystem storage
   - `MockStorage` for testing
   - Path sanitization to prevent traversal attacks
   - Unique key generation for avatars

2. Updated `internal/service/user/service.go`:
   - Added storage and avatar config dependencies
   - `UploadAvatar` now actually stores files via storage interface
   - Cleanup on error (delete stored file if DB operation fails)

3. Created `internal/api/image_utils.go`:
   - `detectImageInfoWithReader` to detect image type and dimensions
   - Supports JPEG, PNG, GIF, WebP
   - Returns new reader since original is consumed

4. Updated `internal/api/handler.go` `UploadAvatar`:
   - Validates file size against config
   - Detects content type from file bytes (not trusting client)
   - Validates against allowed types
   - Extracts image dimensions
   - Returns proper `*ogen.Avatar` response

5. Updated test files to use new `NewService` signature

### Security:
- Path sanitization prevents directory traversal
- Content-type detection from file bytes (not trusting headers)
- File size validation before processing
- MIME type allowlist validation


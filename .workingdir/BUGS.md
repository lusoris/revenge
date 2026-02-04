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


## A0.5: Request Metadata Extraction [P1] - COMPLETED

**Status**: Fixed
**Date**: 2026-02-04

### Changes Made:
1. Created `internal/api/middleware/request_metadata.go`:
   - `RequestMetadata` struct with IP, UserAgent, AcceptLanguage
   - `RequestMetadataMiddleware()` ogen middleware
   - `extractClientIP()` with X-Forwarded-For, X-Real-IP, RemoteAddr support
   - `stripPort()` handling IPv4 and IPv6 correctly
   - Context helpers: `WithRequestMetadata`, `GetRequestMetadata`, `GetIPAddress`, `GetUserAgent`

2. Added comprehensive tests for all cases:
   - X-Forwarded-For chain parsing
   - X-Real-IP fallback
   - IPv4 with/without port
   - IPv6 with brackets and port: `[::1]:8080`
   - IPv6 without port: `::1`

### Usage:
```go
// In handler
meta := middleware.GetRequestMetadata(ctx)
ip := meta.IPAddress
ua := meta.UserAgent
```


## A0.6: WebAuthn Session Cache [P1] - COMPLETED

**Status**: Fixed
**Date**: 2026-02-04

### Changes Made:
1. Updated `internal/service/mfa/webauthn.go`:
   - Added `cache *cache.Cache` field to `WebAuthnService`
   - Added session storage constants: `webAuthnSessionTTL`, key prefixes
   - Added internal helpers: `storeSession`, `getSession`, `deleteSession`
   - Updated `BeginRegistration` to store session in cache
   - Updated `BeginLogin` to store session in cache
   - Added public methods:
     - `GetRegistrationSession(ctx, userID)` - retrieve cached registration session
     - `GetLoginSession(ctx, userID)` - retrieve cached login session
     - `DeleteRegistrationSession(ctx, userID)` - cleanup after finish
     - `DeleteLoginSession(ctx, userID)` - cleanup after finish
     - `HasCache()` - check if cache is configured

2. Updated `internal/service/mfa/module.go`:
   - `NewWebAuthnServiceFromConfig` now accepts `*cache.Client`
   - Creates dedicated named cache "webauthn" with 5-minute TTL
   - Gracefully handles missing cache (logs warning, continues without)

3. Added comprehensive tests in `webauthn_test.go`:
   - `TestWebAuthnService_HasCache` - cache availability check
   - `TestWebAuthnService_SessionCache` - store/retrieve/delete sessions

### Architecture:
- Uses L1 (otter in-memory) + L2 (Dragonfly via rueidis) cache layers
- Sessions expire after 5 minutes (WebAuthn timeout)
- Graceful degradation: works without cache (client provides session)
- Cache key pattern: `webauthn:registration:{userID}`, `webauthn:login:{userID}`

### Usage:
```go
// Handler flow:
options, _ := webauthnService.BeginRegistration(ctx, userID, ...)
// ... client does WebAuthn ceremony ...

// Retrieve session from cache:
session, err := webauthnService.GetRegistrationSession(ctx, userID)
if err != nil {
    // Fallback: use session from client
}

// Finish and cleanup:
_ = webauthnService.FinishRegistration(ctx, userID, ..., *session, ...)
webauthnService.DeleteRegistrationSession(ctx, userID)
```


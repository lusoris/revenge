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


## A0.7: OIDC New User Creation [P1] - COMPLETED

**Status**: Fixed
**Date**: 2026-02-04

### Changes Made:
1. Added to `internal/service/auth/service.go`:
   - `RegisterFromOIDCRequest` struct for OIDC user registration data
   - `RegisterFromOIDC(ctx, req)` - creates user with random unusable password, email already verified
   - `CreateSessionForUser(ctx, userID, ipAddress, userAgent, deviceName)` - creates session for authenticated user (OIDC login flow)

2. Updated `internal/service/oidc/service.go`:
   - Added `ProviderID` field to `CallbackResult` struct
   - Set `ProviderID` in `HandleCallback` when returning new user info

3. Updated `internal/api/handler_oidc.go`:
   - Full implementation of `OidcCallback` handler:
     - When `IsNewUser=true`: creates user via `authService.RegisterFromOIDC`, links to OIDC provider
     - Creates session via `authService.CreateSessionForUser`
     - Returns proper JWT access/refresh tokens (not OIDC provider tokens)
   - Username generation from OIDC: uses `preferred_username` claim, falls back to email prefix

### Security:
- OIDC users get random 64-byte unusable passwords (cannot be used for login)
- Email marked as verified (OIDC provider already verified it)
- Proper JWT tokens generated (not returning OIDC provider tokens)
- Activity logging with `oidc_login: true` metadata

### Usage:
```go
// OIDC callback flow:
// 1. User authenticates with OIDC provider
// 2. Provider redirects to callback with code
// 3. Handler validates code, gets user info
// 4. If new user: create account + link to provider
// 5. Create session with proper JWT tokens
// 6. Return access_token + refresh_token
```


## A0.8: MFA Remember Device Setting [P2] - COMPLETED

**Status**: Fixed
**Date**: 2026-02-04

### Changes Made:
1. Fixed `GetStatus()` in `internal/service/mfa/manager.go`:
   - Now properly fetches `remember_device_enabled` from `user_mfa_settings` table
   - Removed TODO comment, implemented actual fetch

2. Added new methods to `MFAManager`:
   - `SetRememberDevice(ctx, userID, enabled, durationDays)` - enable/disable remember device
   - `GetRememberDeviceSettings(ctx, userID)` - get remember device settings
   - Both methods handle case when MFA settings don't exist (creates them)

### Database:
- Uses existing `UpdateMFASettingsRememberDevice` query
- Uses existing `GetUserMFASettings` query
- Settings created with defaults if they don't exist

### Usage:
```go
// Enable remember device for 30 days
err := mfaManager.SetRememberDevice(ctx, userID, true, 30)

// Disable remember device
err := mfaManager.SetRememberDevice(ctx, userID, false, 0)

// Get current settings
enabled, days, err := mfaManager.GetRememberDeviceSettings(ctx, userID)
```


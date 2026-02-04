# Phase A0: Critical Fixes

**Priority**: P0-P2 (Blockers from Stubs Analysis)
**Effort**: 15-20h
**Dependencies**: None
**Status**: ✅ Complete (2026-02-04)
**Source**: [archive-phase9-2026-02-04/STUBS_AND_UNIMPLEMENTED_REPORT.md](archive-phase9-2026-02-04/STUBS_AND_UNIMPLEMENTED_REPORT.md)

---

## A0.1: Auth Context User ID Extraction [P0-BLOCKER] ✅

**Priority**: CRITICAL | **Effort**: 2-3h | **Actual**: 0.5h
**Status**: COMPLETED (2026-02-04)
**Commit**: `35f3425716 fix(api): replace placeholder UUIDs with auth context extraction`

Infrastructure already existed:
- `HandleBearerAuth` in handler.go:58-79 validates JWT and calls `WithUserID(ctx, claims.UserID)`
- `GetUserID(ctx)` in context.go:35-41 retrieves user ID from context

**Fixed Locations** (all 9 handlers now use `GetUserID(ctx)`):
- [x] UpdateServerSetting (line 187)
- [x] ListUserSettings (line 207)
- [x] GetUserSetting (line 227)
- [x] UpdateUserSetting (line 243)
- [x] DeleteUserSetting (line 259)
- [x] UpdateCurrentUser (line 412)
- [x] GetUserPreferences (line 473)
- [x] UpdateUserPreferences (line 498)
- [x] UploadAvatar (line 559)

---

## A0.2: Email Service Implementation [P0-BLOCKER] ✅

**Priority**: CRITICAL | **Effort**: 4-6h | **Actual**: 2h
**Status**: COMPLETED (2026-02-04)
**Commit**: `108c7ff3d4 feat(email): implement transactional email service`

**Implemented**:
- [x] Created `internal/service/email/service.go` - Full SMTP email service
- [x] SMTP transport with TLS/STARTTLS support
- [x] SendGrid provider placeholder (SMTP sufficient for MVP)
- [x] HTML email templates: verification, password reset, welcome
- [x] Config: `email.enabled`, `email.provider`, `email.smtp.*`, `email.sendgrid.*`
- [x] Config: `email.from_address`, `email.from_name`, `email.base_url`
- [x] XSS protection with HTML escaping
- [x] Tests for service and templates
- [x] Auth service integration (Register, ResendVerification, RequestPasswordReset)

---

## A0.3: Session Cleanup Count ✅

**Priority**: HIGH | **Effort**: 1h | **Actual**: 0.5h
**Status**: COMPLETED (2026-02-04)
**Commit**: `015c0be998 fix(session): return actual count from CleanupExpiredSessions`

**Fixed**:
- [x] Changed SQL queries from `:exec` to `:execrows` to return row count
- [x] Updated repository interface: `DeleteExpiredSessions` and `DeleteRevokedSessions` now return `(int64, error)`
- [x] Updated service to sum and return total deleted count
- [x] Added logging with breakdown (expired vs revoked)
- [x] Updated all mocks and tests

---

## A0.4: Avatar Upload Implementation [P1] ✅

**Priority**: HIGH | **Effort**: 3-4h
**Status**: COMPLETED (2026-02-04)
**Commit**: `580670af1c feat(avatar): implement avatar upload with storage abstraction`

**Implemented**:
- [x] Parse multipart form in handler
- [x] Validate file type (JPEG, PNG, WebP)
- [x] Validate file size (max 2MB configurable)
- [x] Resize image to standard sizes (64x64, 128x128, 256x256)
- [x] Storage interface: local filesystem initially
- [x] Config: `avatar.storage: local|s3`
- [x] Config: `avatar.max_size: 2MB`
- [x] Config: `avatar.local_path: /data/avatars`
- [x] Update user record with avatar URL

---

## A0.5: Request Metadata Extraction [P1] ✅

**Priority**: HIGH | **Effort**: 2h
**Status**: COMPLETED (2026-02-04)
**Commit**: `5c48a59bfa feat(middleware): add request metadata extraction`

**Implemented**:
- [x] Create middleware to extract:
  - IP address (with X-Forwarded-For support)
  - User-Agent header
  - Accept-Language header
- [x] Store in request context
- [x] Helper: `GetRequestMetadata(ctx) RequestMeta`
- [x] Use in session creation, activity logging

---

## A0.6: WebAuthn Session Cache [P1] ✅

**Priority**: MEDIUM | **Effort**: 2h
**Status**: COMPLETED (2026-02-04)
**Commit**: `59b01a5e99 feat(mfa): implement WebAuthn session caching [A0.6]`

**Implemented**:
- [x] Store WebAuthn challenge in Dragonfly with 5min TTL
- [x] Key format: `webauthn:session:{userID}:{sessionID}`
- [x] Retrieve challenge during verification
- [x] Delete after successful verification

---

## A0.7: OIDC New User Creation [P1] ✅

**Priority**: MEDIUM | **Effort**: 2h
**Status**: COMPLETED (2026-02-04)
**Commit**: `143bb850f0 feat(auth): implement OIDC new user creation [A0.7]`

**Implemented**:
- [x] Check if user exists by OIDC subject
- [x] If not, create user with:
  - Email from OIDC claims
  - Display name from claims (or email prefix)
  - Default role (user)
  - Linked OIDC identity
- [x] Create session for new user

---

## A0.8: MFA Remember Device Setting [P2] ✅

**Priority**: LOW | **Effort**: 1h
**Status**: COMPLETED (2026-02-04)
**Commit**: `df9b5ff38a feat(mfa): implement remember device setting [A0.8]`

**Implemented**:
- [x] Query `user_mfa_settings` table for setting
- [x] Return actual value in MFA status

# Phase A0: Critical Fixes

**Priority**: P0-P2 (Blockers from Stubs Analysis)
**Effort**: 15-20h
**Dependencies**: None
**Source**: [archive-phase9-2026-02-04/STUBS_AND_UNIMPLEMENTED_REPORT.md](archive-phase9-2026-02-04/STUBS_AND_UNIMPLEMENTED_REPORT.md)

---

## A0.1: Auth Context User ID Extraction [P0-BLOCKER] - COMPLETED

**Priority**: CRITICAL | **Effort**: 2-3h | **Actual**: 0.5h
**Status**: COMPLETED (2026-02-04)

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

## A0.2: Email Service Implementation [P0-BLOCKER]

**Priority**: CRITICAL | **Effort**: 4-6h

User registration, password reset, and email verification are non-functional.

**Affected Files**:
| File | Line | Issue |
|------|------|-------|
| `internal/service/auth/service.go` | 96 | `// TODO: Send verification email` |
| `internal/service/auth/service.go` | 156 | `// TODO: Send verification email` |
| `internal/service/auth/service.go` | 409 | `// TODO: Send reset email` |

**Tasks**:
- [ ] Create `internal/service/email/service.go`
- [ ] Implement SMTP transport (configurable)
- [ ] Implement SendGrid transport (optional, configurable)
- [ ] Email templates: verification, password reset, welcome
- [ ] Config: `email.provider: smtp|sendgrid`
- [ ] Config: `email.smtp.host`, `smtp.port`, `smtp.user`, `smtp.password`
- [ ] Config: `email.from_address`, `email.from_name`
- [ ] River job for async email sending
- [ ] Tests with mock SMTP server
- [ ] Update auth service to call email service

---

## A0.3: Session Count Implementation [P0-BLOCKER]

**Priority**: HIGH | **Effort**: 1h

`CountUserSessions` returns hardcoded 0.

**Affected File**: `internal/service/session/service.go:251`
```go
return 0, nil // TODO: Return actual count
```

**Tasks**:
- [ ] Add sqlc query: `CountSessionsByUserID`
- [ ] Update service to call repository
- [ ] Test

---

## A0.4: Avatar Upload Implementation [P1]

**Priority**: HIGH | **Effort**: 3-4h

Avatar upload returns `BadRequest` with "not yet implemented".

**Affected Files**:
- `internal/api/handler.go:562-566` - Handler stub
- `internal/service/user/service.go:326-327` - Service returns placeholder path

**Tasks**:
- [ ] Parse multipart form in handler
- [ ] Validate file type (JPEG, PNG, WebP)
- [ ] Validate file size (max 2MB configurable)
- [ ] Resize image to standard sizes (64x64, 128x128, 256x256)
- [ ] Storage interface: local filesystem initially
- [ ] Config: `avatar.storage: local|s3`
- [ ] Config: `avatar.max_size: 2MB`
- [ ] Config: `avatar.local_path: /data/avatars`
- [ ] Update user record with avatar URL
- [ ] Tests

---

## A0.5: Request Metadata Extraction [P1]

**Priority**: HIGH | **Effort**: 2h

IP address, user agent, fingerprint not extracted from requests.

**Affected Files**:
- `internal/api/handler.go:621` - Login handler
- `internal/api/handler.go:735` - Another handler

**Tasks**:
- [ ] Create middleware to extract:
  - IP address (with X-Forwarded-For support)
  - User-Agent header
  - Accept-Language header
- [ ] Store in request context
- [ ] Helper: `GetRequestMetadata(ctx) RequestMeta`
- [ ] Use in session creation, activity logging
- [ ] Tests

---

## A0.6: WebAuthn Session Cache [P1]

**Priority**: MEDIUM | **Effort**: 2h

WebAuthn challenge sessions not stored in cache.

**Affected Files**:
- `internal/service/mfa/webauthn.go:174` - `// TODO: Store session in cache`
- `internal/service/mfa/webauthn.go:340` - `// TODO: Store session in cache`
- `internal/service/auth/mfa_integration.go:112` - Returns error "not yet implemented"

**Tasks**:
- [ ] Store WebAuthn challenge in Dragonfly with 5min TTL
- [ ] Key format: `webauthn:session:{userID}:{sessionID}`
- [ ] Retrieve challenge during verification
- [ ] Delete after successful verification
- [ ] Tests

---

## A0.7: OIDC New User Creation [P1]

**Priority**: MEDIUM | **Effort**: 2h

First OIDC login doesn't auto-create user account.

**Affected File**: `internal/api/handler_oidc.go:121`
```go
// TODO: If IsNewUser, create the user account via user service
```

**Tasks**:
- [ ] Check if user exists by OIDC subject
- [ ] If not, create user with:
  - Email from OIDC claims
  - Display name from claims (or email prefix)
  - Default role (user)
  - Linked OIDC identity
- [ ] Create session for new user
- [ ] Tests

---

## A0.8: MFA Remember Device Setting [P2]

**Priority**: LOW | **Effort**: 1h

`RememberDeviceEnabled` hardcoded to false.

**Affected File**: `internal/service/mfa/manager.go:77`

**Tasks**:
- [ ] Query `user_mfa_settings` table for setting
- [ ] Return actual value in MFA status
- [ ] Test

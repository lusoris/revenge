# Authentication & Authorization Flow

> Written from code as of 2026-02-06. Source of truth: the Go files linked below.

## Overview

```
Client Request
    |
    v
+-------------------+     +------------------+     +------------------+
| HandleBearerAuth  |---->| TokenManager     |---->| Context with     |
| (handler.go)      |     | .ValidateAccess  |     | UserID, Username |
+-------------------+     | Token (jwt.go)   |     +--------+---------+
                          +------------------+              |
                                                            v
                                                   +------------------+
                                                   | Handler method   |
                                                   | GetUserID(ctx)   |
                                                   +--------+---------+
                                                            |
                                                            v
                                                   +------------------+
                                                   | rbac.Service     |
                                                   | .Enforce(sub,    |
                                                   |  resource, act)  |
                                                   +------------------+

Alternative auth (API keys):
  rv_<hex> --> apikeys.ValidateKey --> CheckScope --> handler
```

Key source files:

| Component | File |
|---|---|
| JWT tokens | `internal/service/auth/jwt.go` |
| Bearer middleware | `internal/api/handler.go` (HandleBearerAuth) |
| Context helpers | `internal/api/context.go` |
| Auth service | `internal/service/auth/service.go` |
| MFA integration | `internal/service/auth/mfa_integration.go` |
| MFA manager | `internal/service/mfa/manager.go` |
| TOTP | `internal/service/mfa/totp.go` |
| Sessions | `internal/service/session/service.go` |
| RBAC | `internal/service/rbac/service.go` |
| RBAC permissions | `internal/service/rbac/permissions.go` |
| API keys | `internal/service/apikeys/service.go` |

---

## Login Flow

`auth.Service.Login` (`internal/service/auth/service.go`):

1. **Account lockout check** -- if enabled, count failed attempts within `lockoutWindow`. If `>= lockoutThreshold`, reject.
2. **User lookup** -- by username, then by email if not found.
3. **Constant-time password check** -- always runs `hasher.VerifyPassword`, even for non-existent users (compares against `dummyPasswordHash` to prevent timing attacks).
4. **Active check** -- reject disabled accounts.
5. **Generate JWT access token** -- `tokenManager.GenerateAccessToken(userID, username)`.
6. **Generate refresh token** -- 32 bytes from `crypto/rand`, hex-encoded.
7. **Store refresh token** -- SHA-256 hash in DB with device metadata and expiry.
8. **Clear failed attempts** on success, update `last_login_at`.

```go
// From auth/service.go -- the critical anti-timing-attack pattern:
hashToCompare := dummyPasswordHash
if userFound {
    hashToCompare = user.PasswordHash
}
match, err := s.hasher.VerifyPassword(password, hashToCompare)
if !userFound || !match {
    return nil, errors.New("invalid username or password")
}
```

Login returns `LoginResponse{User, AccessToken, RefreshToken, ExpiresIn}`.

---

## JWT Structure and Validation

**Algorithm**: HS256 (HMAC-SHA256), implemented in stdlib (`crypto/hmac`). No third-party JWT library.

**Claims** (`internal/service/auth/jwt.go`):

```go
type Claims struct {
    UserID    uuid.UUID `json:"user_id"`
    Username  string    `json:"username"`
    IssuedAt  int64     `json:"iat"`      // milliseconds
    ExpiresAt int64     `json:"exp"`      // milliseconds
}
```

The JWT contains only `user_id` and `username`. No roles, no scopes. This keeps tokens small and avoids stale permission data.

**Validation** (`HandleBearerAuth` in `internal/api/handler.go`):

```go
func (h *Handler) HandleBearerAuth(ctx context.Context, operationName ogen.OperationName, t ogen.BearerAuth) (context.Context, error) {
    claims, err := h.tokenManager.ValidateAccessToken(t.Token)
    if err != nil {
        return nil, errors.Wrap(err, "invalid token")
    }
    ctx = WithUserID(ctx, claims.UserID)
    ctx = WithUsername(ctx, claims.Username)
    return ctx, nil
}
```

Validation steps in `ValidateAccessToken`:
1. Split token into 3 parts (header.payload.signature).
2. Recompute HMAC-SHA256 of `header.payload` with server secret.
3. Constant-time compare signatures (`hmac.Equal`).
4. Decode payload, check `exp` against current time (millisecond precision).

Context extraction in handlers:

```go
userID, err := api.GetUserID(ctx)   // internal/api/context.go
username, err := api.GetUsername(ctx)
```

---

## Token Refresh Flow

`auth.Service.RefreshToken` (`internal/service/auth/service.go`):

1. Hash incoming refresh token with SHA-256.
2. Look up hash in DB -- reject if not found or expired.
3. Look up user by `authToken.UserID`.
4. Generate new JWT access token.
5. Update `last_used_at` on the refresh token record.
6. Return **same refresh token** (no rotation at auth service level).

```go
// RefreshToken does NOT rotate the refresh token:
return &LoginResponse{
    User:         user,
    AccessToken:  accessToken,
    RefreshToken: refreshToken, // same token returned
    ExpiresIn:    int64(s.jwtExpiry.Seconds()),
}, nil
```

The session service (`internal/service/session/service.go`) **does** rotate tokens on refresh -- `RefreshSession` creates a new session+refresh token pair and revokes the old session.

**Logout** revokes refresh tokens:
- Single device: `Logout(ctx, refreshToken)` -- revokes by token hash.
- All devices: `LogoutAll(ctx, userID)` -- revokes all user tokens.

---

## MFA Integration

MFA is checked **after** password verification. Source: `auth.Service.LoginWithMFA` (`internal/service/auth/mfa_integration.go`).

### Flow

```
1. LoginWithMFA(username, password, ..., mfaAuthenticator)
       |
       v
2. s.Login(...)  -- password check, returns LoginResponse
       |
       v
3. mfaAuthenticator.CheckMFARequired(userID)
       |
       +-- MFA not required --> return LoginResponse (tokens issued)
       |
       +-- MFA required --> return MFALoginResponse + ErrMFARequired
                            (NO tokens issued yet)
       |
       v
4. Client sends MFA code
       |
       v
5. mfaAuthenticator.VerifyMFA(MFAVerifyRequest{Method, Code})
       |
       +-- "totp"        --> mfaManager.VerifyTOTP(userID, code)
       +-- "backup_code" --> mfaManager.VerifyBackupCode(userID, code, ip)
       +-- "webauthn"    --> not yet implemented
       |
       v
6. s.CompleteMFALogin(sessionID, result)
       --> marks session as mfa_verified in DB
```

### Supported methods

| Method | Implementation | Storage |
|---|---|---|
| TOTP | `mfa.TOTPService` -- SHA1, 6 digits, 30s period | Secret encrypted with AES-256-GCM in DB |
| Backup codes | `mfa.BackupCodesService` | Hashed in DB, single-use |
| WebAuthn | `mfa.WebAuthnService` (stubbed) | Not yet implemented |

### MFA status check

```go
// From mfa/manager.go
type MFAStatus struct {
    HasTOTP               bool
    WebAuthnCount         int64
    UnusedBackupCodes     int64
    RequireMFA            bool   // gates login
    RememberDeviceEnabled bool
}
```

`RequireMFA` must be explicitly enabled per user via `MFAManager.EnableMFA`. Requires at least one method to be configured.

---

## RBAC Permission Checking

**Engine**: Casbin with a custom PostgreSQL adapter (`internal/service/rbac/adapter.go`). Policies stored in `shared.casbin_rule`.

**Model**: subject (userID) -> role -> (resource, action).

### Built-in roles

| Role | Key permissions |
|---|---|
| `admin` | `admin:*` (wildcard, full access) |
| `moderator` | Full content CRUD, user list/get, integrations, audit read |
| `user` | Own profile, view content, stream, create requests |
| `guest` | Read-only content, stream only |

### Checking permissions

```go
// Direct enforcement -- internal/service/rbac/service.go
allowed, err := rbacService.Enforce(ctx, userID.String(), "movies", "create")

// Convenience wrapper
allowed, err := rbacService.EnforceWithContext(ctx, userID, "movies", "create")

// Check specific permission
allowed, err := rbacService.CheckUserPermission(ctx, userID, "libraries", "scan")
```

Permission format is `resource:action`. Defined constants in `internal/service/rbac/permissions.go`:

```go
const (
    PermMoviesList   = "movies:list"
    PermMoviesCreate = "movies:create"
    PermLibrariesScan = "libraries:scan"
    PermAdminAll     = "admin:*"
    // ... etc
)
```

### Admin check in handlers

```go
// internal/api/handler.go
func (h *Handler) isAdmin(ctx context.Context) bool {
    userID, ok := h.getUserID(ctx)
    if !ok { return false }
    roles, err := h.rbacService.GetUserRoles(ctx, userID)
    for _, r := range roles {
        if r == "admin" { return true }
    }
    return false
}
```

---

## API Key Authentication

API keys provide scoped, long-lived access. Source: `internal/service/apikeys/service.go`.

**Format**: `rv_<64 hex chars>` (prefix `rv_` + 32 random bytes hex-encoded).

**Storage**: SHA-256 hash in DB. First 8 chars stored as `key_prefix` for identification. Raw key shown only at creation.

### Validation flow

```go
func (s *Service) ValidateKey(ctx context.Context, rawKey string) (*APIKey, error) {
    if !s.isValidKeyFormat(rawKey) { return nil, ErrInvalidKeyFormat }
    keyHash := s.hashKey(rawKey)
    dbKey, err := s.repo.GetAPIKeyByHash(ctx, keyHash)
    // check IsActive, check ExpiresAt
    // async update LastUsedAt
    return &key, nil
}
```

### Scope checking

```go
func (s *Service) CheckScope(ctx context.Context, keyID uuid.UUID, requiredScope string) (bool, error) {
    key, err := s.GetKey(ctx, keyID)
    for _, scope := range key.Scopes {
        if scope == requiredScope || scope == "admin" {
            return true, nil
        }
    }
    return false, nil
}
```

Valid scopes: `read`, `write`, `admin`. Max 10 keys per user (`DefaultMaxKeysPerUser`).

**Note**: As of this writing, there is no `HandleApiKeyAuth` middleware in the ogen-generated handler. API key auth is available at the service layer but not yet wired into the HTTP security handler.

---

## How to Add a New Protected Endpoint

1. **Define the endpoint in the OpenAPI spec** (`api/openapi/`) with `security: [BearerAuth: []]`. Run `make generate` to regenerate the ogen handler interface.

2. **Implement the handler method** on `api.Handler`. Extract user from context:

```go
func (h *Handler) MyNewEndpoint(ctx context.Context, params ogen.MyNewEndpointParams) (ogen.MyNewEndpointRes, error) {
    userID, err := GetUserID(ctx)
    if err != nil {
        return &ogen.Error{Code: 401, Message: "Unauthorized"}, nil
    }

    // Permission check via RBAC
    allowed, err := h.rbacService.CheckUserPermission(ctx, userID, "myresource", "myaction")
    if err != nil || !allowed {
        return &ogen.Error{Code: 403, Message: "Forbidden"}, nil
    }

    // Business logic...
}
```

3. **Add the permission constant** to `internal/service/rbac/permissions.go` if it is a new resource/action:

```go
const PermMyResourceAction = "myresource:myaction"
```

4. **Assign the permission to roles** in `DefaultRolePermissions` (same file) so built-in roles get access on policy initialization.

5. **JWT validation is automatic** -- `HandleBearerAuth` runs before your handler for any endpoint with `BearerAuth` security in the OpenAPI spec. No additional middleware registration needed.

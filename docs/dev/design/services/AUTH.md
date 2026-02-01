---
sources:
  - name: Uber fx
    url: ../../sources/tooling/fx.md
    note: Auto-resolved from fx
  - name: ogen OpenAPI Generator
    url: ../../sources/tooling/ogen.md
    note: Auto-resolved from ogen
design_refs:
  - title: services
    path: INDEX.md
  - title: 01_ARCHITECTURE
    path: ../architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: ../architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: ../architecture/03_METADATA_SYSTEM.md
---

## Table of Contents

- [Auth Service](#auth-service)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Service Structure](#service-structure)
    - [Dependencies](#dependencies)
    - [Provides](#provides)
    - [Component Diagram](#component-diagram)
  - [Implementation](#implementation)
    - [File Structure](#file-structure)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
    - [Config Keys](#config-keys)
  - [API Endpoints](#api-endpoints)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)


# Auth Service


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: service


> > Authentication, registration, and password management

**Package**: `internal/service/auth`
**fx Module**: `auth.Module`

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | ✅ | - |
| Instructions | ✅ | - |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---


## Architecture

### Service Structure

```
internal/service/auth/
├── module.go              # fx module definition
├── service.go             # Service implementation
├── repository.go          # Data access (if needed)
├── handler.go             # HTTP handlers (if exposed)
├── middleware.go          # Middleware (if needed)
├── types.go               # Domain types
└── service_test.go        # Tests
```

### Dependencies
**Go Packages**:
- `github.com/google/uuid`
- `github.com/jackc/pgx/v5`
- `golang.org/x/crypto/argon2` - Password hashing
- `github.com/alexedwards/argon2id` - Argon2 helper
- `crypto/rand` - Token generation
- `crypto/sha256` - Token hashing
- `go.uber.org/fx`

**External Services**:
- Email service (SMTP) for verification and password reset emails


### Provides
<!-- Service provides -->

### Component Diagram

<!-- Component diagram -->


## Implementation

### File Structure

<!-- File structure -->

### Key Interfaces

```go
type AuthService interface {
  // Registration
  Register(ctx context.Context, req RegisterRequest) (*User, error)
  VerifyEmail(ctx context.Context, token string) error
  ResendVerification(ctx context.Context, userID uuid.UUID) error

  // Login
  Login(ctx context.Context, username, password string) (*User, error)
  Logout(ctx context.Context, sessionID uuid.UUID) error

  // Password management
  ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
  RequestPasswordReset(ctx context.Context, email string) error
  ResetPassword(ctx context.Context, token, newPassword string) error

  // Account security
  CheckRateLimit(ctx context.Context, identifier string) error
  RecordLoginAttempt(ctx context.Context, userID uuid.UUID, success bool, ip net.IP) error
}

type RegisterRequest struct {
  Username    string `json:"username"`
  Email       string `json:"email"`
  Password    string `json:"password"`
  DisplayName string `json:"display_name,omitempty"`
}

type User struct {
  ID            uuid.UUID  `db:"id" json:"id"`
  Username      string     `db:"username" json:"username"`
  Email         string     `db:"email" json:"email"`
  DisplayName   *string    `db:"display_name" json:"display_name,omitempty"`
  EmailVerified bool       `db:"email_verified" json:"email_verified"`
  IsActive      bool       `db:"is_active" json:"is_active"`
  IsAdmin       bool       `db:"is_admin" json:"is_admin"`
  CreatedAt     time.Time  `db:"created_at" json:"created_at"`
}
```


### Dependencies

**Go Packages**:
- `github.com/google/uuid`
- `github.com/jackc/pgx/v5`
- `golang.org/x/crypto/argon2` - Password hashing
- `github.com/alexedwards/argon2id` - Argon2 helper
- `crypto/rand` - Token generation
- `crypto/sha256` - Token hashing
- `go.uber.org/fx`

**External Services**:
- Email service (SMTP) for verification and password reset emails






## Configuration
### Environment Variables

```bash
AUTH_PASSWORD_MIN_LENGTH=8
AUTH_ARGON2_TIME=1
AUTH_ARGON2_MEMORY=64
AUTH_ARGON2_THREADS=4
AUTH_RESET_TOKEN_EXPIRY=1h
AUTH_VERIFICATION_TOKEN_EXPIRY=24h
AUTH_MAX_LOGIN_ATTEMPTS=5
AUTH_LOCKOUT_DURATION=15m
```


### Config Keys

```yaml
auth:
  password:
    min_length: 8
    require_uppercase: true
    require_lowercase: true
    require_number: true
    require_special: false
  argon2:
    time: 1
    memory: 64  # MiB
    threads: 4
  tokens:
    reset_expiry: 1h
    verification_expiry: 24h
  security:
    max_login_attempts: 5
    lockout_duration: 15m
```



## API Endpoints
```
POST   /api/v1/auth/register          # Register new user
POST   /api/v1/auth/login             # Login with credentials
POST   /api/v1/auth/logout            # Logout (invalidate session)
POST   /api/v1/auth/verify-email      # Verify email with token
POST   /api/v1/auth/resend-verification # Resend verification email
POST   /api/v1/auth/forgot-password   # Request password reset
POST   /api/v1/auth/reset-password    # Reset password with token
POST   /api/v1/auth/change-password   # Change password (authenticated)
```

**Example Register Request**:
```json
{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "SecurePass123!",
  "display_name": "John Doe"
}
```

**Example Register Response**:
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "username": "johndoe",
  "email": "john@example.com",
  "display_name": "John Doe",
  "email_verified": false,
  "is_active": true,
  "created_at": "2026-02-01T10:00:00Z"
}
```

**Example Login Request**:
```json
{
  "username": "johndoe",
  "password": "SecurePass123!"
}
```

**Example Login Response**:
```json
{
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "username": "johndoe",
    "email": "john@example.com"
  },
  "session_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```



## Testing Strategy

### Unit Tests

<!-- Unit test strategy -->

### Integration Tests

<!-- Integration test strategy -->

### Test Coverage

Target: **80% minimum**







## Related Documentation
### Design Documents
- [services](INDEX.md)
- [01_ARCHITECTURE](../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Uber fx](../../sources/tooling/fx.md) - Auto-resolved from fx
- [ogen OpenAPI Generator](../../sources/tooling/ogen.md) - Auto-resolved from ogen


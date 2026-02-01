---
sources:
  - name: Uber fx
    url: ../../sources/tooling/fx.md
    note: Auto-resolved from fx
  - name: golang.org/x/crypto
    url: ../../sources/go/x/crypto.md
    note: Auto-resolved from golang-x-crypto
  - name: ogen OpenAPI Generator
    url: ../../sources/tooling/ogen.md
    note: Auto-resolved from ogen
  - name: pgx PostgreSQL Driver
    url: ../../sources/database/pgx.md
    note: Auto-resolved from pgx
  - name: PostgreSQL Arrays
    url: ../../sources/database/postgresql-arrays.md
    note: Auto-resolved from postgresql-arrays
  - name: PostgreSQL JSON Functions
    url: ../../sources/database/postgresql-json.md
    note: Auto-resolved from postgresql-json
  - name: River Job Queue
    url: ../../sources/tooling/river.md
    note: Auto-resolved from river
  - name: sqlc
    url: ../../sources/database/sqlc.md
    note: Auto-resolved from sqlc
  - name: sqlc Configuration
    url: ../../sources/database/sqlc-config.md
    note: Auto-resolved from sqlc-config
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

- [User Service](#user-service)
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
- [User management (admin)](#user-management-admin)
- [Current user (self)](#current-user-self)
- [Profile](#profile)
- [Notifications](#notifications)
- [GDPR](#gdpr)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)


# User Service


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: service


> > User account management and authentication

**Package**: `internal/service/user`
**fx Module**: `user.Module`

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
internal/service/user/
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
- `github.com/riverqueue/river` - Background data export jobs
- `io` - File handling
- `archive/zip` - Data export ZIP creation
- `go.uber.org/fx`


### Provides
<!-- Service provides -->

### Component Diagram

<!-- Component diagram -->


## Implementation

### File Structure

<!-- File structure -->

### Key Interfaces

```go
type UserService interface {
  // User management
  GetUser(ctx context.Context, userID uuid.UUID) (*User, error)
  ListUsers(ctx context.Context, filters UserFilters) ([]User, error)
  UpdateUser(ctx context.Context, userID uuid.UUID, update UserUpdate) (*User, error)
  DeleteUser(ctx context.Context, userID uuid.UUID) error

  // Profile
  GetProfile(ctx context.Context, userID uuid.UUID) (*UserProfile, error)
  UpdateProfile(ctx context.Context, userID uuid.UUID, profile ProfileUpdate) (*UserProfile, error)
  UploadAvatar(ctx context.Context, userID uuid.UUID, file io.Reader) (string, error)

  // Notifications
  GetNotificationPreferences(ctx context.Context, userID uuid.UUID) (*NotificationPreferences, error)
  UpdateNotificationPreferences(ctx context.Context, userID uuid.UUID, prefs NotificationPreferences) error

  // GDPR
  RequestDataExport(ctx context.Context, userID uuid.UUID) (*DataExport, error)
  GetDataExport(ctx context.Context, exportID uuid.UUID) (*DataExport, error)
  RequestAccountDeletion(ctx context.Context, userID uuid.UUID, reason string) error
  CancelDeletion(ctx context.Context, requestID uuid.UUID) error
}

type UserProfile struct {
  UserID            uuid.UUID  `db:"user_id" json:"user_id"`
  Bio               *string    `db:"bio" json:"bio,omitempty"`
  AvatarURL         *string    `db:"avatar_url" json:"avatar_url,omitempty"`
  BannerURL         *string    `db:"banner_url" json:"banner_url,omitempty"`
  Timezone          *string    `db:"timezone" json:"timezone,omitempty"`
  Language          string     `db:"language" json:"language"`
  ProfileVisibility string     `db:"profile_visibility" json:"profile_visibility"`
}
```


### Dependencies

**Go Packages**:
- `github.com/google/uuid`
- `github.com/jackc/pgx/v5`
- `github.com/riverqueue/river` - Background data export jobs
- `io` - File handling
- `archive/zip` - Data export ZIP creation
- `go.uber.org/fx`






## Configuration
### Environment Variables

```bash
USER_DEFAULT_STORAGE_QUOTA_MB=100
USER_MAX_AVATAR_SIZE_MB=5
USER_DATA_EXPORT_EXPIRY=168h  # 7 days
USER_DELETION_GRACE_PERIOD=720h  # 30 days
```


### Config Keys

```yaml
user:
  storage:
    default_quota_mb: 100
    max_avatar_size_mb: 5
    upload_path: /data/uploads/avatars
  gdpr:
    data_export_expiry: 168h
    deletion_grace_period: 720h
  profile:
    default_visibility: private
```



## API Endpoints
```
# User management (admin)
GET    /api/v1/users               # List users
GET    /api/v1/users/:id           # Get user
PUT    /api/v1/users/:id           # Update user
DELETE /api/v1/users/:id           # Delete user

# Current user (self)
GET    /api/v1/users/me            # Get current user
PUT    /api/v1/users/me            # Update current user

# Profile
GET    /api/v1/users/:id/profile   # Get profile
PUT    /api/v1/users/me/profile    # Update profile
POST   /api/v1/users/me/avatar     # Upload avatar

# Notifications
GET    /api/v1/users/me/notifications/preferences
PUT    /api/v1/users/me/notifications/preferences

# GDPR
POST   /api/v1/users/me/data-export           # Request export
GET    /api/v1/users/me/data-export/:id       # Get export status
GET    /api/v1/users/me/data-export/:id/download  # Download export
POST   /api/v1/users/me/delete                # Request deletion
DELETE /api/v1/users/me/delete/:id            # Cancel deletion
```

**Example Profile Update Request**:
```json
{
  "bio": "Movie enthusiast and avid reader",
  "timezone": "America/New_York",
  "language": "en",
  "profile_visibility": "friends"
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
- [golang.org/x/crypto](../../sources/go/x/crypto.md) - Auto-resolved from golang-x-crypto
- [ogen OpenAPI Generator](../../sources/tooling/ogen.md) - Auto-resolved from ogen
- [pgx PostgreSQL Driver](../../sources/database/pgx.md) - Auto-resolved from pgx
- [PostgreSQL Arrays](../../sources/database/postgresql-arrays.md) - Auto-resolved from postgresql-arrays
- [PostgreSQL JSON Functions](../../sources/database/postgresql-json.md) - Auto-resolved from postgresql-json
- [River Job Queue](../../sources/tooling/river.md) - Auto-resolved from river
- [sqlc](../../sources/database/sqlc.md) - Auto-resolved from sqlc
- [sqlc Configuration](../../sources/database/sqlc-config.md) - Auto-resolved from sqlc-config


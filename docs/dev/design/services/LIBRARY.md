# Library Service

> Library management and access control

**Location**: `internal/service/library/`

---

## Overview

The Library service manages media libraries including:

- Library CRUD operations
- User access control (grants/revokes)
- Library type validation
- Access-aware queries

---

## Library Types

```go
const (
    LibraryTypeMovie      = "movie"
    LibraryTypeTvshow     = "tvshow"
    LibraryTypeMusic      = "music"
    LibraryTypeAudiobook  = "audiobook"
    LibraryTypeBook       = "book"
    LibraryTypePodcast    = "podcast"
    LibraryTypePhoto      = "photo"
    LibraryTypeLivetv     = "livetv"
    LibraryTypeComics     = "comics"
    LibraryTypeAdultMovie = "adult_movie"
    LibraryTypeAdultScene = "adult_scene"
)
```

---

## Operations

### Create Library

```go
type CreateParams struct {
    Name              string
    LibraryType       string
    Paths             []string
    ScanEnabled       bool
    ScanIntervalHours int32
    PreferredLanguage *string
    DownloadImages    bool
    DownloadNfo       bool
    GenerateChapters  bool
    IsPrivate         bool
    OwnerUserID       pgtype.UUID
    SortOrder         int32
    Icon              *string
}

func (s *Service) Create(ctx context.Context, params CreateParams) (*db.Library, error)
```

### Get Library

```go
// Get by ID (no access check)
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*db.Library, error)

// Get by ID with user access check
func (s *Service) GetByIDWithAccess(ctx context.Context, id, userID uuid.UUID, isAdmin bool) (*db.Library, error)
```

### List Libraries

```go
// List all (admin only)
func (s *Service) List(ctx context.Context) ([]db.Library, error)

// List accessible to user
func (s *Service) ListAccessible(ctx context.Context, userID uuid.UUID) ([]db.Library, error)

// List by type
func (s *Service) ListByType(ctx context.Context, libraryType db.LibraryType) ([]db.Library, error)

// List with access info (admin)
func (s *Service) ListAll(ctx context.Context) ([]WithAccess, error)

// List with access info (user)
func (s *Service) ListForUser(ctx context.Context, userID uuid.UUID) ([]WithAccess, error)
```

### Update Library

```go
type UpdateParams struct {
    ID                uuid.UUID
    Name              *string
    Paths             []string
    ScanEnabled       *bool
    ScanIntervalHours *int32
    PreferredLanguage *string
    DownloadImages    *bool
    DownloadNfo       *bool
    GenerateChapters  *bool
    IsPrivate         *bool
    SortOrder         *int32
    Icon              *string
}

func (s *Service) Update(ctx context.Context, params UpdateParams) (*db.Library, error)
```

### Delete Library

```go
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error
```

---

## Access Control

### Grant Access

```go
func (s *Service) GrantAccess(ctx context.Context, libraryID, userID uuid.UUID, canManage bool) error
```

### Revoke Access

```go
func (s *Service) RevokeAccess(ctx context.Context, libraryID, userID uuid.UUID) error
```

### Check Access

```go
func (s *Service) UserCanAccess(ctx context.Context, libraryID, userID uuid.UUID) (bool, error)
```

### List Users

```go
func (s *Service) ListUsers(ctx context.Context, libraryID uuid.UUID) ([]db.ListLibraryUsersRow, error)
```

---

## Access Model

```go
type WithAccess struct {
    Library   db.Library
    CanManage bool
}
```

Access rules:
- **Admins**: Access all libraries, can manage all
- **Owners**: Access and manage own private libraries
- **Granted users**: Access based on `library_access` table

---

## Errors

| Error | Description |
|-------|-------------|
| `ErrLibraryNotFound` | Library does not exist |
| `ErrAccessDenied` | User lacks access to library |
| `ErrInvalidLibraryType` | Invalid library type specified |

---

## Related

- [Library Types](../features/shared/LIBRARY_TYPES.md) - Per-module library design
- [RBAC](../features/shared/RBAC_CASBIN.md) - Permission model

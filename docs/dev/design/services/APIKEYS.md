# API Keys Service

> API key generation, validation, and management

**Location**: `internal/service/apikeys/`

---

## Overview

The API Keys service provides programmatic access management:

- Key generation with secure random values
- Hash-based storage (raw key never stored)
- Scope-based permissions
- Expiration support
- Usage tracking

---

## Operations

### Create API Key

```go
type CreateParams struct {
    UserID    uuid.UUID
    Name      string
    Scopes    []string
    ExpiresAt *time.Time
}

type CreateResult struct {
    Key    *db.ApiKey
    RawKey string  // Only returned once!
}

func (s *Service) Create(ctx context.Context, params CreateParams) (*CreateResult, error)
```

**Important**: The raw key is only returned on creation. It cannot be retrieved later.

### Validate API Key

```go
func (s *Service) Validate(ctx context.Context, rawKey string) (*db.ApiKey, error)
```

Returns the API key record if valid, updates usage statistics.

### Retrieve Keys

```go
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*db.ApiKey, error)
func (s *Service) ListByUser(ctx context.Context, userID uuid.UUID) ([]db.ApiKey, error)
```

### Delete Keys

```go
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error
func (s *Service) DeleteExpired(ctx context.Context) error
```

---

## Key Security

```go
// Generation: 32 random bytes
keyBytes := make([]byte, 32)
rand.Read(keyBytes)
rawKey := base64.URLEncoding.EncodeToString(keyBytes)

// Storage: SHA-256 hash only
hash := sha256.Sum256([]byte(rawKey))
keyHash := hex.EncodeToString(hash[:])

// Prefix: first 8 chars for identification
keyPrefix := rawKey[:8]
```

**Storage Model**:
- `key_hash`: SHA-256 hash for lookup
- `key_prefix`: First 8 chars for display
- Raw key: **never stored**

---

## Scopes

```go
func HasScope(apiKey *db.ApiKey, scope string) bool
```

Common scopes:
- `*` - Full access
- `read` - Read-only access
- `library:read` - Read library data
- `library:write` - Modify library data
- `playback` - Playback operations

---

## Errors

| Error | Description |
|-------|-------------|
| `ErrKeyNotFound` | API key does not exist |
| `ErrKeyExpired` | API key has expired |
| `ErrInvalidKey` | Invalid API key format |

---

## Related

- [Auth Service](AUTH.md) - Session-based auth
- [RBAC Service](RBAC.md) - Permission checking

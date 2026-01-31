# API Keys Service

<!-- SOURCES: fx, ogen, pgx, postgresql-arrays, postgresql-json, river, sqlc, sqlc-config -->

<!-- DESIGN: services, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> API key generation, validation, and management


<!-- TOC-START -->

## Table of Contents

- [Developer Resources](#developer-resources)
- [Status](#status)
- [Overview](#overview)
- [Operations](#operations)
  - [Create API Key](#create-api-key)
  - [Validate API Key](#validate-api-key)
  - [Retrieve Keys](#retrieve-keys)
  - [Delete Keys](#delete-keys)
- [Key Security](#key-security)
- [Scopes](#scopes)
- [Errors](#errors)
- [Implementation Checklist](#implementation-checklist)
  - [Phase 1: Core Infrastructure](#phase-1-core-infrastructure)
  - [Phase 2: Database](#phase-2-database)
  - [Phase 3: Service Layer](#phase-3-service-layer)
  - [Phase 4: API Integration](#phase-4-api-integration)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documents](#related-documents)

<!-- TOC-END -->

**Module**: `internal/service/apikeys`

## Developer Resources

> Package versions: [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#go-dependencies-core)

| Package | Purpose |
|---------|---------|
| crypto/rand | Secure random key generation |
| crypto/sha256 | Key hash storage |
| encoding/base64 | URL-safe key encoding |
| pgx | PostgreSQL driver |

## Status

| Dimension | Status |
|-----------|--------|
| Design | ✅ |
| Sources | ✅ |
| Instructions | ✅ |
| Code | 🔴 |
| Linting | 🔴 |
| Unit Testing | 🔴 |
| Integration Testing | 🔴 |---

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

## Implementation Checklist

### Phase 1: Core Infrastructure
- [ ] Create `internal/service/apikeys/` package structure
- [ ] Define entity types in `entity.go`
- [ ] Create repository interface in `repository.go`
- [ ] Add fx module wiring in `module.go`

### Phase 2: Database
- [ ] Create migration for `api_keys` table
- [ ] Add indexes (user_id, key_hash, expires_at)
- [ ] Write sqlc queries

### Phase 3: Service Layer
- [ ] Implement key generation (32 bytes random)
- [ ] Implement SHA-256 hash storage
- [ ] Implement key validation
- [ ] Add scope checking

### Phase 4: API Integration
- [ ] Define OpenAPI endpoints
- [ ] Generate ogen handlers
- [ ] Wire handlers to service
- [ ] Add authentication middleware

---


## Related Documents

- [Auth Service](AUTH.md) - Session-based auth
- [RBAC Service](RBAC.md) - Permission checking
- [Session Service](SESSION.md) - Token management patterns
- [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) - Service inventory

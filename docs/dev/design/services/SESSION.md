# Session Service

<!-- SOURCES: fx, ogen, pgx, postgresql-arrays, postgresql-json, river, sqlc, sqlc-config -->

<!-- DESIGN: services, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> Session token management and device tracking


<!-- TOC-START -->

## Table of Contents

- [Developer Resources](#developer-resources)
- [Status](#status)
- [Overview](#overview)
- [Configuration](#configuration)
- [Operations](#operations)
  - [Create Session](#create-session)
  - [Validate Token](#validate-token)
  - [Deactivate Sessions](#deactivate-sessions)
  - [Update Activity](#update-activity)
- [Token Security](#token-security)
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

**Module**: `internal/service/session`

## Developer Resources

> Package versions: [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#go-dependencies-core)

| Package | Purpose |
|---------|---------|
| crypto/rand | Secure token generation |
| crypto/sha256 | Token hash storage |
| netip | IP address handling |
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

The Session service manages user sessions:

- Token generation and validation
- Device tracking (name, type, client info)
- Session expiration
- Activity updates
- Session limits per user

---

## Configuration

```go
type Service struct {
    queries            *db.Queries
    logger             *slog.Logger
    sessionDuration    time.Duration  // Default: 24h
    maxSessionsPerUser int            // 0 = unlimited
}

func (s *Service) SetSessionDuration(d time.Duration)
func (s *Service) SetMaxSessionsPerUser(maxSessions int)
```

---

## Operations

### Create Session

```go
type CreateParams struct {
    UserID        uuid.UUID
    ProfileID     *uuid.UUID
    DeviceName    *string
    DeviceType    *string
    ClientName    *string
    ClientVersion *string
    IPAddress     netip.Addr
    UserAgent     *string
}

type CreateResult struct {
    Session *db.Session
    Token   string  // Raw token - only returned on creation
}

func (s *Service) Create(ctx context.Context, params CreateParams) (*CreateResult, error)
```

### Validate Token

```go
func (s *Service) ValidateToken(ctx context.Context, token string) (*db.Session, error)
```

Checks:
1. Token hash exists in database
2. Session is active
3. Session not expired

### Deactivate Sessions

```go
// Single session
func (s *Service) Deactivate(ctx context.Context, sessionID uuid.UUID) error

// All sessions for user
func (s *Service) DeactivateAllForUser(ctx context.Context, userID uuid.UUID) error
```

### Update Activity

```go
func (s *Service) UpdateActivity(ctx context.Context, sessionID uuid.UUID, ipAddress *netip.Addr) error
```

---

## Token Security

- **Generation**: 32 bytes random via `crypto/rand`
- **Encoding**: Base64 URL-safe
- **Storage**: SHA-256 hash only (raw token never stored)
- **Lookup**: Hash-based lookup for O(1) validation

```go
// Token generation
tokenBytes := make([]byte, 32)
rand.Read(tokenBytes)
token := base64.URLEncoding.EncodeToString(tokenBytes)

// Storage: hash only
hash := sha256.Sum256([]byte(token))
tokenHash := base64.URLEncoding.EncodeToString(hash[:])
```

---

## Errors

| Error | Description |
|-------|-------------|
| `ErrSessionNotFound` | Session not found or invalid token |
| `ErrSessionExpired` | Session has expired |
| `ErrSessionInactive` | Session was deactivated |
| `ErrTooManySessions` | Max sessions per user exceeded |

---

## Implementation Checklist

### Phase 1: Core Infrastructure
- [ ] Create `internal/service/session/` package structure
- [ ] Define entity types in `entity.go`
- [ ] Create repository interface
- [ ] Add fx module wiring

### Phase 2: Database
- [ ] Create migration for `sessions` table
- [ ] Add indexes (user_id, token_hash, expires_at)
- [ ] Write sqlc queries

### Phase 3: Service Layer
- [ ] Implement token generation (32 bytes random)
- [ ] Implement SHA-256 hash storage
- [ ] Implement session validation
- [ ] Add device tracking
- [ ] Implement session limits per user

### Phase 4: API Integration
- [ ] Define OpenAPI endpoints
- [ ] Generate ogen handlers
- [ ] Wire handlers to service
- [ ] Add activity update middleware

---


## Related Documents

- [Auth Service](AUTH.md) - Login/logout flows
- [User Service](USER.md) - User accounts
- [Activity Service](ACTIVITY.md) - Session activity logging
- [API Keys Service](APIKEYS.md) - Alternative authentication
- [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) - Service inventory

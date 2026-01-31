# Grants Service

> Polymorphic resource access grants for fine-grained sharing


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Developer Resources](#developer-resources)
- [Overview](#overview)
- [Goals](#goals)
- [Non-Goals](#non-goals)
- [Technical Design](#technical-design)
  - [Grant Types](#grant-types)
  - [Permission Levels](#permission-levels)
  - [Repository Interface](#repository-interface)
  - [Service Layer](#service-layer)
  - [Integration with Casbin](#integration-with-casbin)
- [Database Schema](#database-schema)
- [API Endpoints](#api-endpoints)
- [Configuration](#configuration)
- [Implementation Files](#implementation-files)
- [Checklist](#checklist)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documents](#related-documents)

<!-- TOC-END -->

**Module**: `internal/service/grants`
**Dependencies**: [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#go-dependencies-security--rbac)

## Status

| Dimension | Status |
|-----------|--------|
| Design | ✅ |
| Sources | ✅ |
| Instructions | ✅ |
| Code | 🔴 |
| Linting | 🔴 |
| Unit Testing | 🔴 |
| Integration Testing | 🔴 |
## Developer Resources

> Package versions: [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#go-dependencies-security--rbac)

| Package | Purpose |
|---------|---------|
| Casbin | RBAC integration |
| pgx | PostgreSQL driver |
| otter | Grant caching |
| fx | Dependency injection |

---

## Overview

The Grants service provides fine-grained access control beyond RBAC roles. It enables:
- Sharing specific resources (libraries, playlists, items) with specific users
- Time-limited access grants
- Granular permissions (read, write, share, admin)
- Cascading grants (library access includes all content)

This complements Casbin RBAC (role-based) with resource-specific grants (ABAC-like).

## Goals

- Enable users to share specific content with other users
- Support time-limited sharing (expiring grants)
- Allow granular permission levels per grant
- Integrate with RBAC for permission checks

## Non-Goals

- Replace RBAC for role-based permissions
- Public/anonymous sharing (use API keys instead)
- Cross-server federation sharing

---

## Technical Design

### Grant Types

| Grant Type | Scope | Example |
|------------|-------|---------|
| Library | All content in library | Share "Movies" library with user |
| Collection | Specific collection | Share "Marvel Movies" collection |
| Playlist | Specific playlist | Share "Party Mix" playlist |
| Item | Single media item | Share specific movie |

### Permission Levels

| Level | Can Read | Can Edit | Can Share | Can Admin |
|-------|----------|----------|-----------|-----------|
| `view` | Yes | No | No | No |
| `edit` | Yes | Yes | No | No |
| `share` | Yes | Yes | Yes | No |
| `admin` | Yes | Yes | Yes | Yes |

### Repository Interface

```go
type GrantsRepository interface {
    // Create/manage grants
    CreateGrant(ctx context.Context, grant *Grant) error
    RevokeGrant(ctx context.Context, grantID uuid.UUID) error
    UpdateGrant(ctx context.Context, grantID uuid.UUID, updates GrantUpdates) error

    // Query grants
    GetGrant(ctx context.Context, grantID uuid.UUID) (*Grant, error)
    ListGrantsForResource(ctx context.Context, resourceType string, resourceID uuid.UUID) ([]Grant, error)
    ListGrantsForUser(ctx context.Context, userID uuid.UUID) ([]Grant, error)
    ListGrantsByGrantor(ctx context.Context, grantorID uuid.UUID) ([]Grant, error)

    // Check access
    HasAccess(ctx context.Context, userID uuid.UUID, resourceType string, resourceID uuid.UUID, permission string) (bool, error)
    GetEffectivePermission(ctx context.Context, userID uuid.UUID, resourceType string, resourceID uuid.UUID) (string, error)
}

type Grant struct {
    ID           uuid.UUID
    GrantorID    uuid.UUID  // Who created the grant
    GranteeID    uuid.UUID  // Who receives access
    ResourceType string     // "library", "collection", "playlist", "movie", etc.
    ResourceID   uuid.UUID
    Permission   string     // "view", "edit", "share", "admin"
    ExpiresAt    *time.Time // Optional expiration
    CreatedAt    time.Time
}
```

### Service Layer

```go
type GrantsService struct {
    repo   GrantsRepository
    casbin *casbin.Enforcer
}

func (s *GrantsService) GrantAccess(ctx context.Context, grantorID, granteeID uuid.UUID, resourceType string, resourceID uuid.UUID, permission string, expiresAt *time.Time) (*Grant, error)
func (s *GrantsService) RevokeAccess(ctx context.Context, grantorID, grantID uuid.UUID) error
func (s *GrantsService) CanAccess(ctx context.Context, userID uuid.UUID, resourceType string, resourceID uuid.UUID, action string) (bool, error)
func (s *GrantsService) ListSharedWithMe(ctx context.Context, userID uuid.UUID) ([]SharedResource, error)
func (s *GrantsService) ListMyShares(ctx context.Context, userID uuid.UUID) ([]Grant, error)
```

### Integration with Casbin

```go
// Permission check flow:
// 1. Check RBAC role permissions via Casbin
// 2. If denied, check resource-specific grants
// 3. Return highest permission level found

func (s *GrantsService) CanAccess(ctx context.Context, userID uuid.UUID, resourceType string, resourceID uuid.UUID, action string) (bool, error) {
    // First check RBAC
    if allowed, _ := s.casbin.Enforce(userID.String(), resourceType, action); allowed {
        return true, nil
    }

    // Then check grants
    return s.repo.HasAccess(ctx, userID, resourceType, resourceID, action)
}
```

---

## Database Schema

```sql
CREATE TABLE resource_grants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    grantor_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    grantee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    resource_type VARCHAR(50) NOT NULL,
    resource_id UUID NOT NULL,
    permission VARCHAR(20) NOT NULL CHECK (permission IN ('view', 'edit', 'share', 'admin')),
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(grantee_id, resource_type, resource_id)
);

CREATE INDEX idx_grants_grantee ON resource_grants (grantee_id);
CREATE INDEX idx_grants_resource ON resource_grants (resource_type, resource_id);
CREATE INDEX idx_grants_expires ON resource_grants (expires_at) WHERE expires_at IS NOT NULL;
```

---

## API Endpoints

```
POST   /api/v1/grants
       { "grantee_id": "uuid", "resource_type": "library", "resource_id": "uuid", "permission": "view", "expires_at": "2024-12-31T23:59:59Z" }

GET    /api/v1/grants                    # List grants I created
GET    /api/v1/grants/shared-with-me     # List resources shared with me
DELETE /api/v1/grants/{id}               # Revoke a grant
PATCH  /api/v1/grants/{id}               # Update permission or expiration
```

---

## Configuration

```yaml
grants:
  enabled: true
  max_grants_per_user: 100      # Limit grants a user can create
  default_expiration: 0         # 0 = no default expiration
  cascade_library_access: true  # Library grant includes all content
```

---

## Implementation Files

| File | Action | Description |
|------|--------|-------------|
| `internal/service/grants/service.go` | CREATE | Core grants service |
| `internal/service/grants/repository.go` | CREATE | Repository interface |
| `internal/service/grants/repository_pg.go` | CREATE | PostgreSQL implementation |
| `internal/service/grants/module.go` | CREATE | fx module |
| `migrations/shared/000019_resource_grants.up.sql` | EXISTS | Check if migration exists |

---

## Checklist

- [ ] Database migration created/verified
- [ ] Repository interface defined
- [ ] PostgreSQL repository implemented
- [ ] Service layer with Casbin integration
- [ ] API handlers created
- [ ] Expiration cleanup job (River)
- [ ] Tests written
- [ ] Documentation updated

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [Casbin](https://pkg.go.dev/github.com/casbin/casbin/v2) | [Local](../../sources/security/casbin.md) |
| [PostgreSQL Arrays](https://www.postgresql.org/docs/current/arrays.html) | [Local](../../sources/database/postgresql-arrays.md) |
| [PostgreSQL JSON Functions](https://www.postgresql.org/docs/current/functions-json.html) | [Local](../../sources/database/postgresql-json.md) |
| [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) | [Local](../../sources/tooling/river.md) |
| [Uber fx](https://pkg.go.dev/go.uber.org/fx) | [Local](../../sources/tooling/fx.md) |
| [pgx PostgreSQL Driver](https://pkg.go.dev/github.com/jackc/pgx/v5) | [Local](../../sources/database/pgx.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Services](INDEX.md)

### In This Section

- [Activity Service](ACTIVITY.md)
- [Analytics Service](ANALYTICS.md)
- [API Keys Service](APIKEYS.md)
- [Auth Service](AUTH.md)
- [Fingerprint Service](FINGERPRINT.md)
- [Library Service](LIBRARY.md)
- [Metadata Service](METADATA.md)
- [Notification Service](NOTIFICATION.md)

### Related Topics

- [Revenge - Architecture v2](../architecture/01_ARCHITECTURE.md) _Architecture_
- [Revenge - Design Principles](../architecture/02_DESIGN_PRINCIPLES.md) _Architecture_
- [Revenge - Metadata System](../architecture/03_METADATA_SYSTEM.md) _Architecture_
- [Revenge - Player Architecture](../architecture/04_PLAYER_ARCHITECTURE.md) _Architecture_
- [Plugin Architecture Decision](../architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) _Architecture_

### Indexes

- [Design Index](../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

## Related Documents

- [RBAC Service](RBAC.md) - Role-based permissions (Casbin)
- [Library Service](LIBRARY.md) - Library access control
- [User Service](USER.md) - User management
- [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) - Service inventory

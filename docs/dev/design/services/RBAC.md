# RBAC Service

> Role-based access control with Casbin


<!-- TOC-START -->

## Table of Contents

- [Developer Resources](#developer-resources)
- [Status](#status)
- [Overview](#overview)
- [Model](#model)
- [Operations](#operations)
  - [Check Permission](#check-permission)
  - [Role Management](#role-management)
  - [Policy Management](#policy-management)
- [Default Roles](#default-roles)
- [Resource Types](#resource-types)
- [Implementation Checklist](#implementation-checklist)
  - [Phase 1: Core Infrastructure](#phase-1-core-infrastructure)
  - [Phase 2: Casbin Setup](#phase-2-casbin-setup)
  - [Phase 3: Service Layer](#phase-3-service-layer)
  - [Phase 4: Middleware](#phase-4-middleware)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documents](#related-documents)

<!-- TOC-END -->

**Module**: `internal/service/rbac`

## Developer Resources

> See [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#backend-services) for service inventory and status.
> See [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md#go-dependencies-security--rbac) for Casbin package versions.

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
---

## Overview

The RBAC service provides permission management using Casbin:

- Policy enforcement
- Role assignment
- Permission checking
- Dynamic policy updates

---

## Model

Uses RBAC with resource/action model:

```ini
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
```

---

## Operations

### Check Permission

```go
func (s *Service) Enforce(ctx context.Context, sub, obj, act string) (bool, error)
```

Example:
```go
allowed, err := rbac.Enforce(ctx, "user:123", "library:456", "read")
```

### Role Management

```go
// Add role to user
func (s *Service) AddRoleForUser(ctx context.Context, user, role string) error

// Remove role from user
func (s *Service) DeleteRoleForUser(ctx context.Context, user, role string) error

// Get user's roles
func (s *Service) GetRolesForUser(ctx context.Context, user string) ([]string, error)
```

### Policy Management

```go
// Add policy
func (s *Service) AddPolicy(ctx context.Context, sub, obj, act string) error

// Remove policy
func (s *Service) RemovePolicy(ctx context.Context, sub, obj, act string) error

// Get policies for subject
func (s *Service) GetPoliciesForSubject(ctx context.Context, sub string) ([][]string, error)
```

---

## Default Roles

| Role | Permissions |
|------|-------------|
| `admin` | Full access to all resources |
| `moderator` | Manage content, limited user management |
| `user` | Access own libraries, standard features |
| `guest` | Read-only access to public content |

---

## Resource Types

| Resource | Actions |
|----------|---------|
| `library:*` | read, write, delete, manage |
| `user:*` | read, write, delete |
| `settings` | read, write |
| `admin` | access |

---

## Implementation Checklist

### Phase 1: Core Infrastructure
- [ ] Create `internal/service/rbac/` package structure
- [ ] Define role and permission types
- [ ] Configure Casbin model
- [ ] Add fx module wiring

### Phase 2: Casbin Setup
- [ ] Define RBAC model file
- [ ] Create PostgreSQL adapter
- [ ] Load default policies

### Phase 3: Service Layer
- [ ] Implement role assignment
- [ ] Implement permission checking
- [ ] Implement policy management
- [ ] Add caching for hot paths

### Phase 4: Middleware
- [ ] Implement RBAC middleware
- [ ] Add permission enforcement
- [ ] Wire into API router

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
- [Grants Service](GRANTS.md)
- [Library Service](LIBRARY.md)
- [Metadata Service](METADATA.md)

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

- [RBAC_CASBIN.md](../features/shared/RBAC_CASBIN.md) - Full design doc
- [User Service](USER.md) - User management
- [Library Service](LIBRARY.md) - Library access
- [Grants Service](GRANTS.md) - Fine-grained resource grants
- [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) - Service inventory

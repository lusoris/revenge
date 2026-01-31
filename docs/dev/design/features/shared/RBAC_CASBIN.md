# Dynamic RBAC with Casbin

<!-- SOURCES: casbin, casbin-docs, casbin-pgx-adapter, fx, ogen, pgx, postgresql-arrays, postgresql-json, river, sqlc, sqlc-config -->

<!-- DESIGN: features/shared, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> Role-Based Access Control using Casbin for dynamic permission management.


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
- [Architecture](#architecture)
  - [Components](#components)
  - [Data Flow](#data-flow)
- [Database Schema](#database-schema)
  - [roles Table](#roles-table)
  - [System Roles (Cannot Be Deleted)](#system-roles-cannot-be-deleted)
  - [casbin_rules Table (Managed by Casbin)](#casbin-rules-table-managed-by-casbin)
  - [permission_definitions Table](#permission-definitions-table)
- [Permission Categories](#permission-categories)
- [Available Permissions](#available-permissions)
  - [System Permissions](#system-permissions)
  - [User Permissions](#user-permissions)
  - [Library Permissions](#library-permissions)
  - [Content Permissions](#content-permissions)
  - [Playback Permissions](#playback-permissions)
  - [Social Permissions](#social-permissions)
  - [Request Permissions](#request-permissions)
  - [Adult Request Permissions (Schema `c`)](#adult-request-permissions-schema-c)
  - [Adult Permissions](#adult-permissions)
- [Usage](#usage)
  - [Checking Permissions in Handlers](#checking-permissions-in-handlers)
  - [Using Middleware](#using-middleware)
  - [Managing Roles via API](#managing-roles-via-api)
- [Casbin Model](#casbin-model)
- [Performance](#performance)
- [Migration from Static RBAC](#migration-from-static-rbac)
- [Dependencies](#dependencies)
- [Metadata Auditing](#metadata-auditing)
  - [Design Principles](#design-principles)
  - [Audit Log Schema](#audit-log-schema)
  - [Async Write Pattern](#async-write-pattern)
  - [Retention & Cleanup](#retention-cleanup)
  - [Configuration](#configuration)
  - [Audited Actions](#audited-actions)
  - [Metadata Locking](#metadata-locking)
  - [Viewing Edit History](#viewing-edit-history)
  - [Rollback Support](#rollback-support)
- [Resource-Level Permissions (Polymorphic)](#resource-level-permissions-polymorphic)
  - [Two Permission Types](#two-permission-types)
  - [Polymorphic Resource Grants](#polymorphic-resource-grants)
  - [Why Polymorphic?](#why-polymorphic)
  - [Resource Types](#resource-types)
  - [Grant Types](#grant-types)
  - [Usage](#usage)
  - [Cleanup on Delete](#cleanup-on-delete)
- [Implementation Checklist](#implementation-checklist)
  - [Phase 1: Core Infrastructure](#phase-1-core-infrastructure)
  - [Phase 2: Database](#phase-2-database)
  - [Phase 3: Casbin Integration](#phase-3-casbin-integration)
  - [Phase 4: Service Layer](#phase-4-service-layer)
  - [Phase 5: Resource Grants](#phase-5-resource-grants)
  - [Phase 6: Audit Logging](#phase-6-audit-logging)
  - [Phase 7: Metadata Locking](#phase-7-metadata-locking)
  - [Phase 8: Middleware & API Integration](#phase-8-middleware-api-integration)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related](#related)

<!-- TOC-END -->

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Full design with Casbin model, DB schema, audit logging |
| Sources | ✅ | casbin/casbin, pckhoi/casbin-pgx-adapter documented |
| Instructions | ✅ | Implementation checklist added |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |---

## Overview

Revenge uses [Casbin](https://casbin.org/) for dynamic Role-Based Access Control (RBAC). This allows administrators to:

- Create custom roles via the admin interface
- Assign granular permissions to roles at runtime
- Modify role-permission mappings without code changes or restarts

## Architecture

### Components

| Component | Description |
|-----------|-------------|
| `CasbinService` | Core RBAC service in `internal/service/rbac/casbin.go` |
| `casbin_rules` | PostgreSQL table storing Casbin policies |
| `roles` | Custom roles table for role metadata (name, description, color, etc.) |
| `permission_definitions` | Reference table of available permissions for UI |
| `RBACMiddleware` | HTTP middleware for permission checks |

### Data Flow

```
Request → Auth Middleware → RBAC Middleware → Handler
                              ↓
                         CasbinService
                              ↓
                    Casbin Enforcer (cached)
                              ↓
                    casbin_rules (PostgreSQL)
```

## Database Schema

### roles Table

Stores role definitions for admin management:

```sql
CREATE TABLE roles (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(100) NOT NULL UNIQUE,  -- e.g., "editor", "viewer"
    display_name    VARCHAR(255) NOT NULL,          -- e.g., "Content Editor"
    description     TEXT,
    color           VARCHAR(7),                      -- Hex color for UI
    icon            VARCHAR(50),                     -- Icon name for UI
    is_system       BOOLEAN NOT NULL DEFAULT false,  -- Cannot delete system roles
    is_default      BOOLEAN NOT NULL DEFAULT false,  -- Default for new users
    priority        INT NOT NULL DEFAULT 0,          -- UI sorting
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by      UUID REFERENCES users(id)
);
```

### System Roles (Cannot Be Deleted)

| Role | Description | Priority |
|------|-------------|----------|
| `admin` | Full access to all features | 1000 |
| `moderator` | Manage libraries, metadata, content | 500 |
| `user` | Browse, play, rate, playlists | 100 |
| `guest` | Browse only | 0 |

### casbin_rules Table (Managed by Casbin)

Stores permission policies:

```sql
-- Automatically created by Casbin adapter
CREATE TABLE casbin_rules (
    id SERIAL PRIMARY KEY,
    ptype VARCHAR(10),  -- "p" for policy, "g" for grouping
    v0 VARCHAR(256),    -- role name
    v1 VARCHAR(256),    -- permission name
    v2 VARCHAR(256),    -- action (always "allow")
    v3 VARCHAR(256),
    v4 VARCHAR(256),
    v5 VARCHAR(256)
);
```

### permission_definitions Table

Reference table for UI to display available permissions:

```sql
CREATE TABLE permission_definitions (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(100) NOT NULL UNIQUE,    -- e.g., "content.browse"
    display_name    VARCHAR(255) NOT NULL,           -- e.g., "Browse Content"
    description     TEXT NOT NULL,
    category        VARCHAR(50) NOT NULL,            -- e.g., "Content"
    is_dangerous    BOOLEAN NOT NULL DEFAULT false,  -- Needs extra confirmation
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

## Permission Categories

| Category | Description |
|----------|-------------|
| `System` | Server settings, logs, jobs, API keys, roles |
| `Users` | User management |
| `Libraries` | Library management |
| `Content` | Content browsing and metadata |
| `Playback` | Streaming, downloading, transcoding |
| `Social` | Ratings, playlists, collections, history |
| `Adult` | QAR adult content access (schema `qar`) |

## Available Permissions

### System Permissions

| Permission | Description | Dangerous |
|------------|-------------|-----------|
| `system.settings.read` | View server settings | No |
| `system.settings.write` | Modify server settings | Yes |
| `system.logs.read` | View activity logs | No |
| `system.jobs.read` | View background jobs | No |
| `system.jobs.manage` | Manage background jobs | Yes |
| `system.apikeys.manage` | Manage API keys | Yes |
| `system.roles.read` | View roles | No |
| `system.roles.manage` | Manage roles | Yes |

### User Permissions

| Permission | Description | Dangerous |
|------------|-------------|-----------|
| `users.read` | View user list | No |
| `users.create` | Create users | Yes |
| `users.update` | Update users | Yes |
| `users.delete` | Delete users | Yes |
| `users.sessions.manage` | Force logout | Yes |

### Library Permissions

| Permission | Description | Dangerous |
|------------|-------------|-----------|
| `libraries.read` | View libraries | No |
| `libraries.create` | Create libraries | Yes |
| `libraries.update` | Update libraries | Yes |
| `libraries.delete` | Delete libraries | Yes |
| `libraries.scan` | Trigger scans | No |

### Content Permissions

| Permission | Description | Dangerous |
|------------|-------------|-----------|
| `content.browse` | Browse content | No |
| `content.metadata.read` | View metadata | No |
| `content.metadata.write` | Edit metadata (audited) | No |
| `content.metadata.lock` | Lock metadata from auto-updates | No |
| `content.metadata.audit` | View metadata edit history | No |
| `content.images.manage` | Manage images | No |
| `content.delete` | Delete content | Yes |

### Playback Permissions

| Permission | Description | Dangerous |
|------------|-------------|-----------|
| `playback.stream` | Stream content | No |
| `playback.download` | Download files | No |
| `playback.transcode` | Request transcoding | No |

### Social Permissions

| Permission | Description | Dangerous |
|------------|-------------|-----------|
| `social.rate` | Rate content | No |
| `social.playlists.create` | Create playlists | No |
| `social.playlists.manage` | Manage playlists | No |
| `social.collections.create` | Create collections | No |
| `social.collections.manage` | Manage collections | No |
| `social.history.read` | View history | No |
| `social.favorites.manage` | Manage favorites | No |

### Request Permissions

| Permission | Description | Dangerous |
|------------|-------------|-----------|
| `requests.submit` | Submit content requests | No |
| `requests.view.own` | View own requests | No |
| `requests.vote` | Vote on requests | No |
| `requests.comment` | Comment on requests | No |
| `requests.cancel.own` | Cancel own pending requests | No |
| `requests.view.all` | View all users' requests | No |
| `requests.approve` | Approve requests | No |
| `requests.decline` | Decline requests | No |
| `requests.priority` | Set request priority | No |
| `requests.rules.read` | View request rules | No |
| `requests.rules.manage` | Create/edit request rules | Yes |
| `requests.quotas.read` | View user quotas | No |
| `requests.quotas.manage` | Manage user quotas | Yes |
| `requests.polls.vote` | Vote in polls | No |
| `requests.polls.create` | Create polls | No |
| `requests.polls.manage` | Manage polls | Yes |

### Adult Request Permissions (Schema `c`)

| Permission | Description | Dangerous |
|------------|-------------|-----------|
| `adult.requests.submit` | Submit adult content requests | No |
| `adult.requests.view.own` | View own adult requests | No |
| `adult.requests.vote` | Vote on adult requests | No |
| `adult.requests.approve` | Approve adult requests | No |
| `adult.requests.decline` | Decline adult requests | No |
| `adult.requests.rules.manage` | Manage adult request rules | Yes |
| `adult.requests.polls.create` | Create adult content polls | No |
| `adult.requests.polls.manage` | Manage adult polls | Yes |

### Adult Permissions

| Permission | Description | Dangerous |
|------------|-------------|-----------|
| `adult.browse` | Browse adult content | No |
| `adult.stream` | Stream adult content | No |
| `adult.metadata.write` | Edit adult metadata | No |

## Usage

### Checking Permissions in Handlers

```go
func (h *Handler) DeleteMovie(ctx context.Context, params api.DeleteMovieParams) error {
    user, ok := middleware.UserFromContext(ctx)
    if !ok {
        return ErrUnauthorized
    }

    // Check permission
    if err := h.rbac.RequirePermission(ctx, user.ID, rbac.PermContentDelete); err != nil {
        return ErrForbidden
    }

    // Proceed with deletion
    return h.service.DeleteMovie(ctx, params.ID)
}
```

### Using Middleware

```go
// Single permission
router.Handle("/api/v1/admin/settings",
    rbacMiddleware.RequirePermission(rbac.PermSystemSettingsWrite)(settingsHandler))

// Any of multiple permissions
router.Handle("/api/v1/libraries/{id}/scan",
    rbacMiddleware.RequireAnyPermission(rbac.PermLibrariesScan, rbac.PermSystemJobsManage)(scanHandler))

// Role-based (for backwards compatibility)
router.Handle("/api/v1/admin/users",
    rbacMiddleware.RequireRole("admin", "moderator")(usersHandler))
```

### Managing Roles via API

```go
// Create custom role
role, err := rbacService.CreateRole(ctx, rbac.CreateRoleParams{
    Name:        "editor",
    DisplayName: "Content Editor",
    Description: "Can edit content metadata and images",
    Color:       "#10B981",
    Permissions: []string{
        rbac.PermContentBrowse,
        rbac.PermContentMetadataRead,
        rbac.PermContentMetadataWrite,
        rbac.PermContentImagesManage,
    },
    CreatedBy: adminUserID,
})

// Update role permissions
err = rbacService.SetRolePermissions(ctx, "editor", []string{
    rbac.PermContentBrowse,
    rbac.PermContentMetadataRead,
    rbac.PermContentMetadataWrite,
    rbac.PermContentImagesManage,
    rbac.PermLibrariesScan, // Added new permission
})

// Delete custom role (only non-system roles)
err = rbacService.DeleteRole(ctx, "editor")
```

## Casbin Model

Located at `internal/service/rbac/model.conf`:

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

## Performance

- Policies are cached in memory by Casbin
- Database is only hit on policy changes or explicit reload
- Policy changes auto-save to database
- Thread-safe with read-write mutex

## Migration from Static RBAC

The migration `000018_dynamic_rbac.up.sql`:

1. Creates `roles` table with system roles (admin, moderator, user, guest)
2. Creates `permission_definitions` table
3. Adds `role_id` column to `users` table
4. Migrates existing users from ENUM to role_id
5. Casbin adapter creates `casbin_rules` table on first run
6. Default policies are seeded by `CasbinService` on startup

## Dependencies

```go
github.com/casbin/casbin/v2           // Casbin core
github.com/pckhoi/casbin-pgx-adapter/v3  // PostgreSQL adapter (pgx native)
```

## Metadata Auditing

All metadata edits by users with `content.metadata.write` permission are logged to the `activity_log` table for accountability and rollback capability.

### Design Principles

1. **Async writes** - Audit entries written via River job queue, never blocking
2. **Retention policy** - Auto-cleanup after configurable period (default: 90 days)
3. **Partitioned** - Monthly partitions for fast queries and efficient cleanup

### Audit Log Schema

```sql
-- activity_log with monthly partitioning for performance
CREATE TABLE activity_log (
    id              UUID NOT NULL DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,  -- No FK to avoid write contention
    action          VARCHAR(50) NOT NULL,
    module          VARCHAR(50) NOT NULL,
    entity_id       UUID NOT NULL,
    entity_type     VARCHAR(50) NOT NULL,
    changes         JSONB NOT NULL,
    ip_address      INET,
    user_agent      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id, created_at)  -- Partition key included
) PARTITION BY RANGE (created_at);

-- Create partitions (managed by pg_partman or River job)
CREATE TABLE activity_log_2026_01 PARTITION OF activity_log
    FOR VALUES FROM ('2026-01-01') TO ('2026-02-01');

-- Indexes on each partition (auto-created)
CREATE INDEX idx_activity_log_entity ON activity_log(module, entity_type, entity_id);
CREATE INDEX idx_activity_log_user ON activity_log(user_id, created_at DESC);
```

### Async Write Pattern

```go
// Never write directly - queue via River
func (s *AuditService) LogAction(ctx context.Context, entry AuditEntry) error {
    // Fire-and-forget via River job queue
    _, err := s.river.Insert(ctx, AuditLogArgs{Entry: entry}, nil)
    return err  // Only fails if queue insert fails, not the actual log write
}

// River worker processes audit entries async
type AuditLogWorker struct {
    river.WorkerDefaults[AuditLogArgs]
    db *pgxpool.Pool
}

func (w *AuditLogWorker) Work(ctx context.Context, job *river.Job[AuditLogArgs]) error {
    _, err := w.db.Exec(ctx, `INSERT INTO activity_log ...`, job.Args.Entry)
    return err
}
```

### Retention & Cleanup

```go
// River scheduled job - runs daily at 03:00
type AuditCleanupArgs struct{}

func (AuditCleanupArgs) Kind() string { return "audit_cleanup" }

func (w *AuditCleanupWorker) Work(ctx context.Context, job *river.Job[AuditCleanupArgs]) error {
    retention := w.config.AuditRetentionDays  // Default: 90
    cutoff := time.Now().AddDate(0, 0, -retention)

    // Drop old partitions (fast, O(1))
    // Or DELETE for non-partitioned: DELETE FROM activity_log WHERE created_at < $1
    _, err := w.db.Exec(ctx, `DROP TABLE IF EXISTS activity_log_`+cutoff.Format("2006_01"))
    return err
}
```

### Configuration

```yaml
audit:
  enabled: true
  retention_days: 90          # Auto-delete after 90 days
  async: true                 # Always async via River (recommended)
  partition_by: month         # monthly partitions
```

### Audited Actions

| Action | Description |
|--------|-------------|
| `metadata.edit` | Field values changed (title, overview, year, etc.) |
| `metadata.lock` | Metadata locked from automatic provider updates |
| `metadata.unlock` | Metadata unlocked, allows auto-updates again |
| `metadata.refresh` | Manual metadata refresh triggered |
| `image.upload` | Custom image uploaded |
| `image.select` | Different image selected from provider |
| `image.delete` | Image removed |

### Metadata Locking

Editors can lock specific items from automatic metadata updates:

```go
// Lock metadata - prevents automatic updates
err := metadataService.LockMetadata(ctx, rbac.LockMetadataParams{
    EntityType: "movie",
    EntityID:   movieID,
    LockedBy:   userID,
    Reason:     "Manual corrections applied",
})

// Locked metadata is skipped during scheduled refreshes
func (s *MetadataService) RefreshMetadata(ctx context.Context, entityType string, entityID uuid.UUID) error {
    if s.IsLocked(ctx, entityType, entityID) {
        slog.Debug("skipping locked metadata", "type", entityType, "id", entityID)
        return nil
    }
    // ... fetch from provider
}
```

### Viewing Edit History

Users with `content.metadata.audit` permission can view edit history:

```go
// Get edit history for an item
history, err := activityService.GetEntityHistory(ctx, rbac.EntityHistoryParams{
    Module:     "movie",
    EntityType: "movie",
    EntityID:   movieID,
    Limit:      50,
})

// Returns list of changes with diffs
for _, entry := range history {
    fmt.Printf("%s by %s at %s\n", entry.Action, entry.UserName, entry.CreatedAt)
    for field, change := range entry.Changes {
        fmt.Printf("  %s: %q → %q\n", field, change.Old, change.New)
    }
}
```

### Rollback Support

Admins can revert to previous values using the audit log:

```go
// Revert a specific change
err := metadataService.RevertChange(ctx, rbac.RevertChangeParams{
    ActivityLogID: logEntryID,
    RevertedBy:    adminUserID,
    Reason:        "Incorrect edit by user",
})
```

---

## Resource-Level Permissions (Polymorphic)

Role permissions (above) define **what actions a role can perform**. Resource permissions define **which specific resources a user can access**.

### Two Permission Types

| Type | Question | Example |
|------|----------|---------|
| **Role** (Casbin) | "Can this role do X?" | Can `user` role stream content? |
| **Resource** (Polymorphic) | "Can this user access Y?" | Can user 123 view movie library 456? |

### Polymorphic Resource Grants

```sql
-- shared/000019_resource_grants.up.sql
-- Polymorphic resource access (no central registry needed)
CREATE TABLE resource_grants (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Polymorphic reference (grant knows what it's for)
    resource_type   VARCHAR(50) NOT NULL,   -- 'movie_library', 'playlist', 'collection'
    resource_id     UUID NOT NULL,          -- UUID of the actual resource

    -- Grant level
    grant_type      VARCHAR(20) NOT NULL DEFAULT 'view',  -- 'view', 'edit', 'manage', 'owner'

    -- Audit
    granted_by      UUID REFERENCES users(id),
    granted_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at      TIMESTAMPTZ,            -- Optional expiration

    UNIQUE (user_id, resource_type, resource_id)
);

CREATE INDEX idx_resource_grants_user ON resource_grants(user_id);
CREATE INDEX idx_resource_grants_resource ON resource_grants(resource_type, resource_id);
CREATE INDEX idx_resource_grants_expires ON resource_grants(expires_at) WHERE expires_at IS NOT NULL;
```

### Why Polymorphic?

1. **No registry to sync** - No central table tracking all resources
2. **Grant owns the reference** - Self-describing, module validates existence
3. **Works for any resource** - Libraries, playlists, collections, items
4. **Clean deletes** - Module deletes grants when resource deleted
5. **No FK constraints to modules** - Loose coupling

### Resource Types

| Type | Description | Module |
|------|-------------|--------|
| `movie_library` | Movie library | movie |
| `tv_library` | TV library | tvshow |
| `music_library` | Music library | music |
| `adult_library` | Adult library | c (adult) |
| `playlist` | User playlist | shared |
| `collection` | User collection | shared |

### Grant Types

| Grant | Can View | Can Edit | Can Manage | Is Owner |
|-------|----------|----------|------------|----------|
| `view` | ✅ | ❌ | ❌ | ❌ |
| `edit` | ✅ | ✅ | ❌ | ❌ |
| `manage` | ✅ | ✅ | ✅ | ❌ |
| `owner` | ✅ | ✅ | ✅ | ✅ |

### Usage

```go
// Check resource access
func (s *MovieModule) CanAccessLibrary(ctx context.Context, userID, libraryID uuid.UUID) (bool, error) {
    // 1. Check if user owns the library (stored in library table)
    lib, err := s.repo.GetLibrary(ctx, libraryID)
    if err != nil {
        return false, err
    }
    if lib.OwnerUserID == userID {
        return true, nil
    }

    // 2. Check polymorphic resource grant
    return s.grants.HasGrant(ctx, userID, "movie_library", libraryID, "view")
}

// Grant access to another user
func (s *MovieModule) ShareLibrary(ctx context.Context, ownerID, targetUserID, libraryID uuid.UUID) error {
    // Verify ownership
    lib, err := s.repo.GetLibrary(ctx, libraryID)
    if err != nil {
        return err
    }
    if lib.OwnerUserID != ownerID {
        return ErrNotOwner
    }

    // Create polymorphic grant
    return s.grants.CreateGrant(ctx, grants.CreateParams{
        UserID:       targetUserID,
        ResourceType: "movie_library",
        ResourceID:   libraryID,
        GrantType:    "view",
        GrantedBy:    ownerID,
    })
}
```

### Cleanup on Delete

Modules delete grants when resource is deleted:

```go
func (s *MovieModule) DeleteLibrary(ctx context.Context, libraryID uuid.UUID) error {
    // Delete the library
    if err := s.repo.DeleteLibrary(ctx, libraryID); err != nil {
        return err
    }

    // Clean up grants (no FK, so must be explicit)
    return s.grants.DeleteByResource(ctx, "movie_library", libraryID)
}
```

---

## Implementation Checklist

**Location**: `internal/service/rbac/`

### Phase 1: Core Infrastructure
- [ ] Create package structure `internal/service/rbac/`
- [ ] Define entities: `Role`, `PermissionDefinition`, `ResourceGrant`, `AuditEntry`
- [ ] Create Casbin model file `model.conf` with RBAC rules
- [ ] Create repository interface `RBACRepository`
- [ ] Implement fx module `rbac.Module`
- [ ] Define permission constants (e.g., `PermContentDelete`, `PermSystemSettingsWrite`)

### Phase 2: Database
- [ ] Create migration `xxx_dynamic_rbac.up.sql`
- [ ] Create `roles` table with system role flag
- [ ] Create `permission_definitions` table for UI reference
- [ ] Create `resource_grants` table for polymorphic grants
- [ ] Create `activity_log` table with monthly partitioning
- [ ] Add `role_id` column to `users` table
- [ ] Seed system roles (admin, moderator, user, guest)
- [ ] Seed permission definitions for all categories
- [ ] Generate sqlc queries for role and permission CRUD
- [ ] Configure `casbin_rules` table (auto-created by adapter)

### Phase 3: Casbin Integration
- [ ] Install `github.com/casbin/casbin/v2`
- [ ] Install `github.com/pckhoi/casbin-pgx-adapter/v3`
- [ ] Initialize Casbin enforcer with PostgreSQL adapter
- [ ] Load model configuration on startup
- [ ] Implement policy caching for performance
- [ ] Add thread-safe read-write mutex for policy updates
- [ ] Seed default policies on first startup

### Phase 4: Service Layer
- [ ] Implement `CasbinService` with permission checks
- [ ] Implement `RequirePermission(ctx, userID, permission)` method
- [ ] Implement `HasPermission(ctx, userID, permission)` method
- [ ] Implement `CreateRole(ctx, params)` for custom roles
- [ ] Implement `SetRolePermissions(ctx, roleName, permissions)` method
- [ ] Implement `DeleteRole(ctx, roleName)` (block system roles)
- [ ] Implement `AssignRoleToUser(ctx, userID, roleName)` method
- [ ] Add role validation (prevent deleting system roles)

### Phase 5: Resource Grants
- [ ] Implement `GrantsService` for polymorphic resource access
- [ ] Implement `CreateGrant(ctx, params)` for sharing resources
- [ ] Implement `HasGrant(ctx, userID, resourceType, resourceID, grantType)` check
- [ ] Implement `DeleteByResource(ctx, resourceType, resourceID)` for cleanup
- [ ] Implement grant expiration handling
- [ ] Add grant type hierarchy (owner > manage > edit > view)

### Phase 6: Audit Logging
- [ ] Implement `AuditService` with async River job writes
- [ ] Create `AuditLogWorker` for processing audit entries
- [ ] Create `AuditCleanupWorker` for retention policy (drop old partitions)
- [ ] Implement `LogAction(ctx, entry)` fire-and-forget method
- [ ] Implement `GetEntityHistory(ctx, params)` for viewing edit history
- [ ] Implement `RevertChange(ctx, params)` for rollback support
- [ ] Configure daily cleanup job (default: 90 day retention)

### Phase 7: Metadata Locking
- [ ] Implement `MetadataLockService` for editor protections
- [ ] Implement `LockMetadata(ctx, params)` method
- [ ] Implement `UnlockMetadata(ctx, params)` method
- [ ] Implement `IsLocked(ctx, entityType, entityID)` check
- [ ] Update metadata refresh jobs to skip locked items

### Phase 8: Middleware & API Integration
- [ ] Implement `RBACMiddleware` for HTTP routes
- [ ] Implement `RequirePermission(permission)` middleware
- [ ] Implement `RequireAnyPermission(permissions...)` middleware
- [ ] Implement `RequireRole(roles...)` middleware (backwards compat)
- [ ] Define OpenAPI spec for role management endpoints
- [ ] Generate ogen handlers for RBAC admin API
- [ ] Implement `GET /api/v1/admin/roles` - list roles
- [ ] Implement `POST /api/v1/admin/roles` - create role
- [ ] Implement `PUT /api/v1/admin/roles/:id` - update role permissions
- [ ] Implement `DELETE /api/v1/admin/roles/:id` - delete custom role
- [ ] Implement `GET /api/v1/admin/permissions` - list permission definitions
- [ ] Add RBAC checks to all admin handlers

---


## Related

- [RBAC Service](../../services/RBAC.md) - Service implementation (code, API, middleware)
- [ARCHITECTURE.md](../../architecture/01_ARCHITECTURE.md) - Overall architecture
- [DESIGN_PRINCIPLES.md](../../architecture/02_DESIGN_PRINCIPLES.md) - Design principles
- [LIBRARY_TYPES.md](LIBRARY_TYPES.md) - Per-module library architecture

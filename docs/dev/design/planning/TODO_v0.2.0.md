# TODO v0.2.0 - Core

<!-- DESIGN: planning, README, test_output_claude, test_output_wiki -->


<!-- TOC-START -->

## Table of Contents

- [Overview](#overview)
- [Deliverables](#deliverables)
  - [Auth Service](#auth-service)
  - [User Service](#user-service)
  - [Session Service](#session-service)
  - [RBAC Service (Casbin)](#rbac-service-casbin)
  - [API Keys Service](#api-keys-service)
  - [OIDC Service](#oidc-service)
  - [Settings Service](#settings-service)
  - [Activity Service](#activity-service)
  - [Library Service](#library-service)
  - [Health Service (Enhancement)](#health-service-enhancement)
  - [PostgreSQL Integration](#postgresql-integration)
  - [Dragonfly/Redis Integration](#dragonflyredis-integration)
  - [River Job Queue Setup](#river-job-queue-setup)
- [Verification Checklist](#verification-checklist)
- [Dependencies from SOURCE_OF_TRUTH](#dependencies-from-source-of-truth)
- [Related Documentation](#related-documentation)

<!-- TOC-END -->


> Backend Services

**Status**: 🔴 Not Started
**Tag**: `v0.2.0`
**Focus**: Core Backend Services

**Depends On**: [v0.1.0](TODO_v0.1.0.md) (Project structure must exist)

---

## Overview

This milestone implements all core backend services: authentication, user management, sessions, RBAC, API keys, OIDC, settings, activity logging, and library management. These services power all content modules.

---

## Deliverables

### Auth Service

- [ ] **Database Schema** (`migrations/`)
  - [ ] `auth_tokens` table (refresh tokens)
  - [ ] `password_reset_tokens` table
  - [ ] `email_verification_tokens` table

- [ ] **Repository** (`internal/service/auth/`)
  - [ ] `repository.go` - Interface definition
  - [ ] `repository_pg.go` - PostgreSQL implementation
  - [ ] Token CRUD operations
  - [ ] Token validation queries

- [ ] **Service** (`internal/service/auth/service.go`)
  - [ ] Login (email/password)
  - [ ] Logout (invalidate tokens)
  - [ ] Register (create user + send verification)
  - [ ] Password reset flow
  - [ ] Email verification flow
  - [ ] JWT generation (access + refresh)
  - [ ] Token refresh logic

- [ ] **Handler** (`internal/api/auth_handler.go`)
  - [ ] `POST /api/v1/auth/login`
  - [ ] `POST /api/v1/auth/logout`
  - [ ] `POST /api/v1/auth/register`
  - [ ] `POST /api/v1/auth/refresh`
  - [ ] `POST /api/v1/auth/forgot-password`
  - [ ] `POST /api/v1/auth/reset-password`
  - [ ] `POST /api/v1/auth/verify-email`

- [ ] **Middleware** (`internal/api/middleware/auth.go`)
  - [ ] JWT validation middleware
  - [ ] Token extraction from header
  - [ ] User context injection

- [ ] **Tests**
  - [ ] Unit tests with mockery mocks
  - [ ] Integration tests with embedded-postgres

### User Service

- [ ] **Database Schema** (`migrations/`)
  - [ ] `shared.users` table
  - [ ] `shared.user_preferences` table
  - [ ] `shared.user_avatars` table
  - [ ] Indexes on email, username

- [ ] **Repository** (`internal/service/user/`)
  - [ ] `repository.go` - Interface
  - [ ] `repository_pg.go` - Implementation
  - [ ] User CRUD operations
  - [ ] Preference operations
  - [ ] Avatar operations

- [ ] **Service** (`internal/service/user/service.go`)
  - [ ] Get user profile
  - [ ] Update user profile
  - [ ] Change password
  - [ ] Upload avatar
  - [ ] Get/Set preferences
  - [ ] Delete account

- [ ] **Handler** (`internal/api/user_handler.go`)
  - [ ] `GET /api/v1/users/me`
  - [ ] `PATCH /api/v1/users/me`
  - [ ] `POST /api/v1/users/me/avatar`
  - [ ] `DELETE /api/v1/users/me/avatar`
  - [ ] `GET /api/v1/users/me/preferences`
  - [ ] `PATCH /api/v1/users/me/preferences`
  - [ ] `POST /api/v1/users/me/change-password`
  - [ ] `DELETE /api/v1/users/me` (soft delete)

- [ ] **Admin Endpoints**
  - [ ] `GET /api/v1/admin/users` (list all)
  - [ ] `GET /api/v1/admin/users/:id`
  - [ ] `PATCH /api/v1/admin/users/:id`
  - [ ] `DELETE /api/v1/admin/users/:id`
  - [ ] `POST /api/v1/admin/users/:id/disable`
  - [ ] `POST /api/v1/admin/users/:id/enable`

- [ ] **Tests**
  - [ ] Unit tests
  - [ ] Integration tests

### Session Service

- [ ] **Database Schema** (`migrations/`)
  - [ ] `shared.sessions` table
  - [ ] `shared.session_devices` table
  - [ ] Indexes on user_id, token

- [ ] **Repository** (`internal/service/session/`)
  - [ ] `repository.go` - Interface
  - [ ] `repository_pg.go` - PostgreSQL (persistent)
  - [ ] `repository_cache.go` - Dragonfly (fast lookup)
  - [ ] Session CRUD
  - [ ] Device tracking

- [ ] **Service** (`internal/service/session/service.go`)
  - [ ] Create session (on login)
  - [ ] Validate session
  - [ ] Extend session (on activity)
  - [ ] Revoke session (on logout)
  - [ ] Revoke all sessions (security)
  - [ ] List active sessions
  - [ ] Device fingerprinting

- [ ] **Handler** (`internal/api/session_handler.go`)
  - [ ] `GET /api/v1/sessions` (list user's sessions)
  - [ ] `DELETE /api/v1/sessions/:id` (revoke specific)
  - [ ] `DELETE /api/v1/sessions` (revoke all)

- [ ] **Tests**
  - [ ] Unit tests
  - [ ] Integration tests with Dragonfly

### RBAC Service (Casbin)

- [ ] **Database Schema** (`migrations/`)
  - [ ] `shared.casbin_rules` table
  - [ ] Casbin policy tables

- [ ] **Model Configuration** (`config/casbin_model.conf`)
  - [ ] RBAC with resource permissions
  - [ ] Role inheritance
  - [ ] Domain support (for libraries)

- [ ] **Repository** (`internal/service/rbac/`)
  - [ ] `repository.go` - Interface
  - [ ] `repository_pg.go` - casbin-pgx-adapter

- [ ] **Service** (`internal/service/rbac/service.go`)
  - [ ] Initialize Casbin enforcer
  - [ ] Check permission
  - [ ] Add policy
  - [ ] Remove policy
  - [ ] Get roles for user
  - [ ] Get users for role
  - [ ] Assign role to user
  - [ ] Remove role from user

- [ ] **Middleware** (`internal/api/middleware/rbac.go`)
  - [ ] Permission check middleware
  - [ ] Route-based enforcement

- [ ] **Handler** (`internal/api/rbac_handler.go`)
  - [ ] `GET /api/v1/admin/roles`
  - [ ] `POST /api/v1/admin/roles`
  - [ ] `DELETE /api/v1/admin/roles/:name`
  - [ ] `GET /api/v1/admin/roles/:name/users`
  - [ ] `POST /api/v1/admin/users/:id/roles`
  - [ ] `DELETE /api/v1/admin/users/:id/roles/:role`

- [ ] **Default Roles**
  - [ ] `admin` - Full access
  - [ ] `user` - Standard user
  - [ ] `guest` - Limited access
  - [ ] `legacy:read` - QAR access

- [ ] **Tests**
  - [ ] Unit tests with embedded Casbin
  - [ ] Integration tests

### API Keys Service

- [ ] **Database Schema** (`migrations/`)
  - [ ] `shared.api_keys` table (hashed keys)
  - [ ] `shared.api_key_permissions` table

- [ ] **Repository** (`internal/service/apikeys/`)
  - [ ] `repository.go` - Interface
  - [ ] `repository_pg.go` - Implementation

- [ ] **Service** (`internal/service/apikeys/service.go`)
  - [ ] Generate API key
  - [ ] Validate API key
  - [ ] List API keys (without secret)
  - [ ] Revoke API key
  - [ ] Set key permissions
  - [ ] Rate limiting per key

- [ ] **Middleware** (`internal/api/middleware/apikey.go`)
  - [ ] API key extraction (header/query)
  - [ ] Key validation
  - [ ] Permission injection

- [ ] **Handler** (`internal/api/apikey_handler.go`)
  - [ ] `GET /api/v1/users/me/api-keys`
  - [ ] `POST /api/v1/users/me/api-keys`
  - [ ] `DELETE /api/v1/users/me/api-keys/:id`
  - [ ] `PATCH /api/v1/users/me/api-keys/:id`

- [ ] **Tests**
  - [ ] Unit tests
  - [ ] Integration tests

### OIDC Service

- [ ] **Database Schema** (`migrations/`)
  - [ ] `shared.oidc_providers` table
  - [ ] `shared.oidc_user_links` table

- [ ] **Repository** (`internal/service/oidc/`)
  - [ ] `repository.go` - Interface
  - [ ] `repository_pg.go` - Implementation

- [ ] **Service** (`internal/service/oidc/service.go`)
  - [ ] Provider configuration (Authentik, Keycloak, etc.)
  - [ ] OAuth2 flow (authorize URL, callback)
  - [ ] Token exchange
  - [ ] User info fetching
  - [ ] Account linking/unlinking
  - [ ] Auto-registration

- [ ] **Handler** (`internal/api/oidc_handler.go`)
  - [ ] `GET /api/v1/auth/oidc/providers` (list enabled)
  - [ ] `GET /api/v1/auth/oidc/:provider/authorize`
  - [ ] `GET /api/v1/auth/oidc/:provider/callback`
  - [ ] `POST /api/v1/users/me/oidc/link/:provider`
  - [ ] `DELETE /api/v1/users/me/oidc/unlink/:provider`

- [ ] **Admin Endpoints**
  - [ ] `GET /api/v1/admin/oidc/providers`
  - [ ] `POST /api/v1/admin/oidc/providers`
  - [ ] `PATCH /api/v1/admin/oidc/providers/:id`
  - [ ] `DELETE /api/v1/admin/oidc/providers/:id`

- [ ] **Tests**
  - [ ] Unit tests with mock provider
  - [ ] Integration tests

### Settings Service

- [ ] **Database Schema** (`migrations/`)
  - [ ] `shared.server_settings` table (key-value)
  - [ ] `shared.user_settings` table

- [ ] **Repository** (`internal/service/settings/`)
  - [ ] `repository.go` - Interface
  - [ ] `repository_pg.go` - Implementation

- [ ] **Service** (`internal/service/settings/service.go`)
  - [ ] Get server setting
  - [ ] Set server setting
  - [ ] Get user setting
  - [ ] Set user setting
  - [ ] Setting validation
  - [ ] Setting encryption (secrets)

- [ ] **Handler** (`internal/api/settings_handler.go`)
  - [ ] `GET /api/v1/settings` (public server settings)
  - [ ] `GET /api/v1/admin/settings`
  - [ ] `PATCH /api/v1/admin/settings`

- [ ] **Tests**
  - [ ] Unit tests
  - [ ] Integration tests

### Activity Service

- [ ] **Database Schema** (`migrations/`)
  - [ ] `shared.activity_logs` table
  - [ ] Indexes on user_id, action, timestamp
  - [ ] Partitioning by month

- [ ] **Repository** (`internal/service/activity/`)
  - [ ] `repository.go` - Interface
  - [ ] `repository_pg.go` - Implementation

- [ ] **Service** (`internal/service/activity/service.go`)
  - [ ] Log activity (async via River)
  - [ ] Query activity (with filters)
  - [ ] Activity aggregation
  - [ ] Retention policy (cleanup old)

- [ ] **Handler** (`internal/api/activity_handler.go`)
  - [ ] `GET /api/v1/users/me/activity`
  - [ ] `GET /api/v1/admin/activity` (all users)
  - [ ] `GET /api/v1/admin/activity/stats`

- [ ] **River Job** (`internal/service/activity/jobs.go`)
  - [ ] ActivityLogJob - Async activity logging

- [ ] **Tests**
  - [ ] Unit tests
  - [ ] Integration tests

### Library Service

- [ ] **Database Schema** (`migrations/`)
  - [ ] `shared.libraries` table
  - [ ] `shared.library_paths` table
  - [ ] `shared.library_access` table (user ↔ library)

- [ ] **Repository** (`internal/service/library/`)
  - [ ] `repository.go` - Interface
  - [ ] `repository_pg.go` - Implementation

- [ ] **Service** (`internal/service/library/service.go`)
  - [ ] Create library
  - [ ] Update library
  - [ ] Delete library
  - [ ] Add/Remove paths
  - [ ] Grant/Revoke user access
  - [ ] Get libraries for user
  - [ ] Check library access

- [ ] **Handler** (`internal/api/library_handler.go`)
  - [ ] `GET /api/v1/libraries`
  - [ ] `GET /api/v1/libraries/:id`
  - [ ] `POST /api/v1/admin/libraries`
  - [ ] `PATCH /api/v1/admin/libraries/:id`
  - [ ] `DELETE /api/v1/admin/libraries/:id`
  - [ ] `POST /api/v1/admin/libraries/:id/paths`
  - [ ] `DELETE /api/v1/admin/libraries/:id/paths/:pathId`
  - [ ] `POST /api/v1/admin/libraries/:id/access`
  - [ ] `DELETE /api/v1/admin/libraries/:id/access/:userId`

- [ ] **Tests**
  - [ ] Unit tests
  - [ ] Integration tests

### Health Service (Enhancement)

- [ ] **Enhanced Checks** (`internal/infra/health/`)
  - [ ] Add service-level health checks
  - [ ] River job queue health
  - [ ] External API health (metadata providers)

- [ ] **Prometheus Metrics** (`internal/infra/metrics/`)
  - [ ] Request duration histogram
  - [ ] Request count by endpoint
  - [ ] Active sessions gauge
  - [ ] Database pool metrics
  - [ ] Cache hit ratio

### PostgreSQL Integration

- [ ] **Pool Enhancement** (`internal/infra/database/`)
  - [ ] Connection pool metrics
  - [ ] Query logging (debug mode)
  - [ ] Slow query detection
  - [ ] Self-healing (Reset on error)

- [ ] **sqlc Queries** (`internal/infra/database/queries/`)
  - [ ] User queries
  - [ ] Session queries
  - [ ] Library queries
  - [ ] Settings queries
  - [ ] Activity queries

### Dragonfly/Redis Integration

- [ ] **Client Setup** (`internal/infra/cache/rueidis.go`)
  - [ ] rueidis client configuration
  - [ ] Client-side caching
  - [ ] Connection pooling

- [ ] **Cache Operations** (`internal/infra/cache/cache.go`)
  - [ ] Session caching
  - [ ] Rate limiting
  - [ ] Distributed locks
  - [ ] Pub/Sub for invalidation

- [ ] **otter L1 Cache** (`internal/infra/cache/otter.go`)
  - [ ] In-memory W-TinyLFU cache
  - [ ] TTL configuration
  - [ ] Size limits

### River Job Queue Setup

- [ ] **Client Setup** (`internal/infra/jobs/river.go`)
  - [ ] River client configuration
  - [ ] Worker pool setup
  - [ ] Graceful shutdown

- [ ] **Queue Configuration** (`internal/infra/jobs/queues.go`)
  - [ ] Define queue priorities
  - [ ] Worker count per queue
  - [ ] Retry policies

- [ ] **Base Jobs**
  - [ ] ActivityLogJob
  - [ ] EmailSendJob
  - [ ] CleanupJob (expired tokens, sessions)

---

## Verification Checklist

- [ ] All services have 80%+ test coverage
- [ ] All endpoints documented in OpenAPI spec
- [ ] Authentication flow works end-to-end
- [ ] RBAC permissions enforced correctly
- [ ] Session management works with Dragonfly
- [ ] River jobs process correctly
- [ ] Health endpoints report all dependencies
- [ ] CI pipeline passes

---

## Dependencies from SOURCE_OF_TRUTH

| Package | Version | Purpose |
|---------|---------|---------|
| github.com/casbin/casbin/v2 | v2.135.0 | RBAC framework |
| github.com/pckhoi/casbin-pgx-adapter/v3 | v3.2.0 | Casbin PostgreSQL |
| github.com/redis/rueidis | v1.0.49 | Redis/Dragonfly |
| github.com/maypok86/otter/v2 | v2.x | In-memory cache |
| github.com/riverqueue/river | v0.26.0 | Job queue |
| github.com/wneessen/go-mail | v0.6.2 | Email sending |
| golang.org/x/crypto | v0.47.0 | Password hashing |

---

## Design Documentation

> **Note**: Design work for v0.2.0 scope is **COMPLETE**. The following design documents exist and should be referenced during implementation:

### Core Service Designs (All 19 services are designed)
- [AUTH.md](../services/AUTH.md) - Authentication service (login, register, JWT, tokens)
- [USER.md](../services/USER.md) - User management service (profiles, preferences, avatars)
- [SESSION.md](../services/SESSION.md) - Session service (persistent sessions, device tracking)
- [RBAC.md](../services/RBAC.md) - RBAC service with Casbin (roles, permissions, policies)
- [APIKEYS.md](../services/APIKEYS.md) - API key management service
- [OIDC.md](../services/OIDC.md) - OIDC authentication integration
- [SETTINGS.md](../services/SETTINGS.md) - Settings service (user + server settings)
- [ACTIVITY.md](../services/ACTIVITY.md) - Activity logging service (audit trail)
- [LIBRARY.md](../services/LIBRARY.md) - Library scanner service (file watching, providers)

### Integration Designs
- [DRAGONFLY.md](../integrations/infrastructure/DRAGONFLY.md) - Session storage, caching patterns
- [RIVER.md](../integrations/infrastructure/RIVER.md) - Background job processing
- [POSTGRESQL.md](../integrations/infrastructure/POSTGRESQL.md) - Database schema patterns

### Auth Provider Integrations (Designed, can add post-v0.2.0)
- [AUTHELIA.md](../integrations/auth/AUTHELIA.md) - Authelia OIDC integration
- [AUTHENTIK.md](../integrations/auth/AUTHENTIK.md) - Authentik OIDC integration
- [KEYCLOAK.md](../integrations/auth/KEYCLOAK.md) - Keycloak OIDC integration
- [GENERIC_OIDC.md](../integrations/auth/GENERIC_OIDC.md) - Generic OIDC provider support

### Technical Architecture
- [EMAIL.md](../technical/EMAIL.md) - Email sending system (verification, password reset)
- [RBAC.md](../features/shared/RBAC.md) - Dynamic RBAC architecture with Casbin

---

## Related Documentation

- [ROADMAP.md](ROADMAP.md) - Full roadmap overview
- [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) - Authoritative versions

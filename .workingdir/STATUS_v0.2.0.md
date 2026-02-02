# v0.2.0 Status

**Version**: v0.2.0 - Core Backend Services
**Start Date**: 2026-02-02
**Target**: TBD
**Current Status**: � Infrastructure Complete (Database, Cache, Jobs 100%)

## Overview

Backend services implementation: Auth, User, Session, RBAC, API Keys, OIDC, Settings, Activity, Library, Health, PostgreSQL, Dragonfly, River.

## Progress Tracker

### Services

| Service | Status | Progress | Notes |
|---------|--------|----------|-------|
| Auth | � Complete | 100% | ✓ DB ✓ Repo ✓ JWT ✓ Service ✓ Middleware ✓ API (Commits 20-25) |
| User | 🟢 Complete | 100% | ✓ DB ✓ Repo ✓ Service ✓ API (Commits 17-19) |
| Session | 🟡 In Progress | 70% | ✓ DB ✓ Repo ✓ Service ⏳ API (Commit 26) |
| RBAC | 🔴 Not Started | 0% | Casbin integration |
| API Keys | 🔴 Not Started | 0% | Key generation, validation |
| OIDC | 🔴 Not Started | 0% | SSO providers |
| Settings | 🟢 Complete | 100% | ✓ DB ✓ Service ✓ API (Commits 11-16) |
| Activity | 🔴 Not Started | 0% | Audit logging |
| Library | 🔴 Not Started | 0% | Library CRUD |
| Health | 🔴 Not Started | 0% | Enhanced checks |

### Infrastructure

| Component | Status | Progress | Notes |
|-----------|--------|----------|-------|
| PostgreSQL | 🟢 Complete | 100% | ✓ Migrations ✓ sqlc ✓ Metrics ✓ Query Logging (4/4) |
| Dragonfly | 🟢 Complete | 100% | ✓ Rueidis client ✓ Otter L1 ✓ Cache Ops (3/3) |
| River | � Complete | 100% | ✓ River client ✓ Queue config ✓ Cleanup job (3/3) |

### Testing

| Category | Coverage | Target | Status |
|----------|----------|--------|--------|
| Auth | ~30% | 80% | 🟡 In Progress (DB/Repo/JWT tests pending) |
| User | ~40% | 80% | 🟡 Partial (Service tests exist) |
| Session | 0% | 80% | 🔴 Not Started |
| RBAC | 0% | 80% | 🔴 Not Started |
| API Keys | 0% | 80% | 🔴 Not Started |
| OIDC | 0% | 80% | 🔴 Not Started |
| Settings | ~50% | 80% | 🟡 Partial (DB tests exist) |
| Activity | 0% | 80% | 🔴 Not Started |
| Library | 0% | 80% | 🔴 Not Started |

## Current Sprint

**Sprint**: Not Started
**Focus**: TBD

### Active Tasks

- None yet

## Completed Milestones

- ✅ **2026-02-02**: Database Layer (Migrations, sqlc, Metrics, Query Logging) - 4/4 steps
- ✅ **2026-02-02**: Dragonfly/Redis Cache (Rueidis, Otter L1, Cache Ops) - 3/3 steps
- ✅ **2026-02-02**: River Client Setup - Step 3.1 complete
- ✅ **2026-02-02**: River Queue Configuration - Step 3.2 complete (3 queues, 2 backoff strategies)
- ✅ **2026-02-02**: River Cleanup Job - Step 3.3 complete (validation, dry-run, 8 tests)
- ✅ **2026-02-02**: River Job Queue - FULLY COMPLETE (all 3 steps, 31 tests, 65.6% coverage)
- ✅ **2026-02-02**: Settings Service (Commits 11-16) - DB, Service layer, API, 6 commits
- ✅ **2026-02-02**: User Service (Commits 17-19) - DB, Repository, Service, 3 commits
- ✅ **2026-02-02**: Auth Service Step 6.1 (Commit 20) - 3 token tables, SHA-256 hashing
- ✅ **2026-02-02**: Auth Service Step 6.2 (Commit 21) - 27 sqlc queries, PostgreSQL repo
- ✅ **2026-02-02**: Auth Service Step 6.3 (Commit 22) - JWT manager (stdlib crypto only)
- ✅ **2026-02-02**: Auth Service Step 6.4 (Commit 23) - Service layer (9 methods, Argon2id)
- ✅ **2026-02-02**: Auth Service Step 6.5 (Commit 24) - Middleware (JWT validation, context)
- ✅ **2026-02-02**: Auth Service Step 6.6 (Commit 25) - API Handler (8 endpoints, 0 lint)
- ✅ **2026-02-02**: Session Service Step 7 (Commit 26) - Repository + Service (17 queries, 0 lint)

## Next Steps

1. **Step 7**: Session Service (Active sessions, device management, session history)
2. **Step 8**: RBAC Service (Casbin integration, roles, permissions)
3. **Step 9+**: API Keys, OIDC, Activity, Library services

## Reference

- **Design Doc**: [TODO_v0.2.0.md](/docs/dev/design/planning/TODO_v0.2.0.md)
- **Source of Truth**: [00_SOURCE_OF_TRUTH.md](/docs/dev/design/00_SOURCE_OF_TRUTH.md)
- **Roadmap**: [ROADMAP.md](/docs/dev/design/planning/ROADMAP.md)

## Updates Log

| Date | Update |
|------|--------|
| 2026-02-02 | Created initial status file |
| 2026-02-02 | Completed Database Layer (4/4): Migrations, sqlc, Metrics, Query Logging |
| 2026-02-02 | Completed Cache Layer (3/3): Rueidis client, Otter L1, Cache Operations |
| 2026-02-02 | Completed Step 3.1: River client setup (36% coverage, 0 lint issues) |
| 2026-02-02 | Completed Step 3.2: Queue config (3 queues, 2 backoff, coverage 55.6%) |
| 2026-02-02 | Completed Step 3.3: Cleanup job (validation, dry-run, coverage 65.6%) |
| 2026-02-02 | ✅ INFRASTRUCTURE COMPLETE: PostgreSQL + Dragonfly + River (100%) |
| 2026-02-02 | Completed Settings Service (Commits 11-16): Database + Service + API |
| 2026-02-02 | Completed User Service (Commits 17-19): Users + Preferences + Avatars |
| 2026-02-02 | Completed Auth Step 6.1 (Commit 20): 3 token tables with SHA-256 hashing |
| 2026-02-02 | Completed Auth Step 6.2 (Commit 21): 27 repository methods + sqlc |
| 2026-02-02 | Completed Auth Step 6.3 (Commit 22): JWT manager (HMAC-SHA256, stdlib) || 2026-02-02 | Completed Auth Step 6.4 (Commit 23): Service layer (9 methods, Argon2id passwords) |
| 2026-02-02 | Completed Auth Step 6.5 (Commit 24): Middleware (HandleBearerAuth, context injection) |
| 2026-02-02 | ✅ Auth Service COMPLETE (Commit 25): API Handler (8 endpoints, 9 schemas, 0 lint) |
| 2026-02-02 | Session Service 70% (Commit 26): Repository + Service layer (API deferred) |

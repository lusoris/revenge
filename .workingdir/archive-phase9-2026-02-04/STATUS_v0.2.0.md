# v0.2.0 Status

**Version**: v0.2.0 - Core Backend Services
**Start Date**: 2026-02-02
**Target**: TBD
**Current Status**: 🟢 COMPLETE - All 10 services implemented
**Testing Phase**: 🟡 IN PROGRESS - Notification 97.6%, Session 59.6%, others pending

## Overview

Backend services implementation: Auth, User, Session, RBAC, API Keys, OIDC, Settings, Activity, Library, Health, PostgreSQL, Dragonfly, River.

## Progress Tracker

### Services

| Service | Status | Progress | Notes |
|---------|--------|----------|-------|
| Auth | 🟢 Complete | 100% | ✓ DB ✓ Repo ✓ JWT ✓ Service ✓ Middleware ✓ API (Commits 20-25) |
| User | 🟢 Complete | 100% | ✓ DB ✓ Repo ✓ Service ✓ API (Commits 17-19) |
| Session | � Complete | 100% | ✓ DB ✓ Repo ✓ Service ✓ API (Commits 26, 28) |
| RBAC | 🟢 Complete | 100% | ✓ DB ✓ Adapter ✓ Service ✓ API (Commits 27, 28) |
| API Keys | � Complete | 100% | ✓ DB ✓ Repo ✓ Service ✓ API (Commit 29) |
| OIDC | 🟢 Complete | 100% | ✓ DB ✓ Repo ✓ Service ✓ API (Commit 30) |
| Settings | 🟢 Complete | 100% | ✓ DB ✓ Service ✓ API (Commits 11-16) |
| Activity | 🟢 Complete | 100% | ✓ DB ✓ Repo ✓ Service ✓ Cleanup Job ✓ API |
| Library | � Complete | 100% | ✓ DB ✓ Repo ✓ Service ✓ API (Step 12) |
| Health | 🟢 Complete | 100% | ✓ Real checks for Cache/Jobs/DB (Step 13) |

### Infrastructure

| Component | Status | Progress | Notes |
|-----------|--------|----------|-------|
| PostgreSQL | 🟢 Complete | 100% | ✓ Migrations ✓ sqlc ✓ Metrics ✓ Query Logging (4/4) |
| Dragonfly | 🟢 Complete | 100% | ✓ Rueidis client ✓ Otter L1 ✓ Cache Ops (3/3) |
| River | � Complete | 100% | ✓ River client ✓ Queue config ✓ Cleanup job (3/3) |

### Testing (Updated 2026-02-04)

| Category | Coverage | Target | Status |
|----------|----------|--------|--------|
| Auth | 29.9% | 80% | 🟡 In Progress |
| Session | 59.6% | 80% | 🟡 Good Progress |
| Notification | 97.6% | 80% | 🟢 **Complete** |
| Notification Agents | 83.9% | 80% | 🟢 **Complete** |
| Search | 37.0% | 80% | 🟡 In Progress |
| MFA | 12.7% | 80% | 🟡 Started |
| RBAC | 1.3% | 80% | 🔴 Needs Work |
| OIDC | 1.7% | 80% | 🔴 Needs Work |
| Activity | 1.2% | 80% | 🔴 Needs Work |
| User | 0% | 80% | 🔴 Needs Work (integration tests exist) |
| Settings | 0% | 80% | 🔴 Needs Work |
| API Keys | 0% | 80% | 🔴 Needs Work |
| Library | 0% | 80% | 🔴 Needs Work |

## Current Sprint

**Sprint**: v0.2.0 Complete
**Focus**: Core Backend Services

### Active Tasks

- ✅ Step 13: Health Service - Complete

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
- ✅ **2026-02-02**: RBAC Service Step 8 (Commit 27) - Casbin integration (12 methods, 0 lint)
- ✅ **2026-02-02**: OIDC Service Step 10 - SSO providers (11 endpoints, OAuth2 flows, token exchange)
- ✅ **2026-02-02**: Activity Service Step 11 - Audit logging (5 admin endpoints, River cleanup job)
- ✅ **2026-02-02**: Library Service Step 12 - Library CRUD (10 endpoints, scans, permissions)
- ✅ **2026-02-02**: Health Service Step 13 - Enhanced checks (cache/jobs/db real checks, 18 tests)

## Next Steps

### Phase 2: Testing & Quality Assurance 🟡 IN PROGRESS

**TestDB Pattern Implementation ✅**
- **Location**: `internal/testutil/testdb.go`, `testdb_migrate.go`
- **Pattern**: PostgreSQL Template Database für instant cloning
- **Performance**: ~90ms pro Test (vs. 3-5s vorher)
- **Features**:
  - `sync.Once` für shared PostgreSQL instance
  - Template DB mit allen Migrations pre-applied
  - `NewTestDB(t)` cloned Template instant
  - `t.Cleanup()` dropped Test-DB automatisch
  - `TestMain` mit `StopSharedPostgres()` für sauberes Cleanup
- **Dokumentation**: `.workingdir/RESEARCH_parallel_db_testing.md`

**Current Focus**: Test Coverage Expansion
- [ ] Messe aktuelle Coverage
- [ ] Migriere langsame Tests zu TestDB Pattern
- [ ] Implementiere fehlende Service Tests

v0.2.0 Core Backend Services is complete. Remaining work deferred to v0.3.0+:
- Content Services (Movies, Shows, Music, Collections)
- Search Integration (Meilisearch)
- Transcoding Integration

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
| 2026-02-02 | Completed Auth Step 6.5 (Commit 24): Middleware (HandleBearerAuth, context injection) |
| 2026-02-02 | ✅ Auth Service COMPLETE (Commit 25): API Handler (8 endpoints, 9 schemas, 0 lint) |
| 2026-02-02 | Session Service 70% (Commit 26): Repository + Service layer (API deferred) |
| 2026-02-02 | RBAC Service 80% (Commit 27): Casbin adapter + Service (API deferred) |
| 2026-02-02 | Completed Settings Service (Commits 11-16): Database + Service + API |
| 2026-02-02 | Completed User Service (Commits 17-19): Users + Preferences + Avatars |
| 2026-02-02 | Completed Auth Step 6.1 (Commit 20): 3 token tables with SHA-256 hashing |
| 2026-02-02 | Completed Auth Step 6.2 (Commit 21): 27 repository methods + sqlc |
| 2026-02-02 | Completed Auth Step 6.3 (Commit 22): JWT manager (HMAC-SHA256, stdlib) || 2026-02-02 | Completed Auth Step 6.4 (Commit 23): Service layer (9 methods, Argon2id passwords) |
| 2026-02-02 | Completed Auth Step 6.5 (Commit 24): Middleware (HandleBearerAuth, context injection) |
| 2026-02-02 | ✅ Auth Service COMPLETE (Commit 25): API Handler (8 endpoints, 9 schemas, 0 lint) |
| 2026-02-02 | Session Service 70% (Commit 26): Repository + Service layer (API deferred) |
| 2026-02-02 | RBAC Service 80% (Commit 27): Casbin adapter + Service (API deferred) |
| 2026-02-02 | ✅ Session API COMPLETE (Commit 28): 6 endpoints, SessionInfo schema, Error type pattern |
| 2026-02-02 | ✅ RBAC API COMPLETE (Commit 28): 6 endpoints (admin only), dedicated type aliases for 403 |
| 2026-02-02 | ✅ API Keys Service COMPLETE (Commit 29): 4 endpoints, SHA-256 hashing, rv_ key format |
| 2026-02-02 | ✅ OIDC Service COMPLETE (Commit 30): 11 endpoints, OAuth2 flows, JWT token exchange |
| 2026-02-02 | ✅ Activity Service COMPLETE (Step 11): 5 admin endpoints, cleanup job, 15 queries |
| 2026-02-02 | ✅ Library Service COMPLETE (Step 12): 10 endpoints, CRUD, scans, permissions |
| 2026-02-02 | ✅ Health Service COMPLETE (Step 13): Real checks for Cache/Jobs/DB, 18 tests |
| 2026-02-02 | ✅ v0.2.0 CORE BACKEND SERVICES COMPLETE - All 10 services implemented |
| 2026-02-03 | 🧪 Testing Phase Started: TestDB pattern + User Service 81.3% coverage |

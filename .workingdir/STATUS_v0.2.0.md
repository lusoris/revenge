# v0.2.0 Status

**Version**: v0.2.0 - Core Backend Services
**Start Date**: 2026-02-02
**Target**: TBD
**Current Status**: 🟡 In Progress (Database: 80%, Cache: 100%, Jobs: 30%)

## Overview

Backend services implementation: Auth, User, Session, RBAC, API Keys, OIDC, Settings, Activity, Library, Health, PostgreSQL, Dragonfly, River.

## Progress Tracker

### Services

| Service | Status | Progress | Notes |
|---------|--------|----------|-------|
| Auth | 🔴 Not Started | 0% | Login, JWT, tokens, password reset |
| User | 🔴 Not Started | 0% | Profile, preferences, avatar |
| Session | 🔴 Not Started | 0% | Token management, devices |
| RBAC | 🔴 Not Started | 0% | Casbin integration |
| API Keys | 🔴 Not Started | 0% | Key generation, validation |
| OIDC | 🔴 Not Started | 0% | SSO providers |
| Settings | 🟡 In Progress | 30% | ✓ Database layer ⏳ Service layer |
| Activity | 🔴 Not Started | 0% | Audit logging |
| Library | 🔴 Not Started | 0% | Library CRUD |
| Health | 🔴 Not Started | 0% | Enhanced checks |

### Infrastructure

| Component | Status | Progress | Notes |
|-----------|--------|----------|-------|
| PostgreSQL | 🟢 Complete | 100% | ✓ Migrations ✓ sqlc ✓ Metrics ✓ Query Logging (4/4) |
| Dragonfly | 🟢 Complete | 100% | ✓ Rueidis client ✓ Otter L1 ✓ Cache Ops (3/3) |
| River | 🟡 In Progress | 30% | ✓ River client ⏳ Queue config ⏳ Job types (1/3) |

### Testing

| Category | Coverage | Target | Status |
|----------|----------|--------|--------|
| Auth | 0% | 80% | 🔴 Not Started |
| User | 0% | 80% | 🔴 Not Started |
| Session | 0% | 80% | 🔴 Not Started |
| RBAC | 0% | 80% | 🔴 Not Started |
| API Keys | 0% | 80% | 🔴 Not Started |
| OIDC | 0% | 80% | 🔴 Not Started |
| Settings | 0% | 80% | 🔴 Not Started |
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

## Next Steps

1. Complete Step 3.2: Queue Configuration (priorities, retry policies)
2. Complete Step 3.3: Base Job Types (cleanup job)
3. Start Step 4: Settings Service (migrations, service layer)
4. Start Step 5: Auth Service (core authentication)

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

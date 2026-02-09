# Build Readiness Analysis

**Date**: 2026-02-06
**Branch**: develop (5 commits ahead of origin)
**Go**: 1.25.6 with `GOEXPERIMENT=greenteagc,jsonv2`

---

## TL;DR

The project **compiles and is substantially built**. 72K lines of Go source across 262 files, with 63K lines of tests (144 files). All 15 core services are implemented, both content modules (movie, tvshow) are complete, and infrastructure is solid. There are **4 broken test files** from recent refactoring and **1 doc inaccuracy** (JOBS.md) that should be fixed before continuing development.

---

## 1. Code Health

### Compilation: PASS

`go build ./...` succeeds with zero errors.

### Go Vet: 3 issues in 2 packages

| Package | Issue |
|---------|-------|
| `internal/service/user` | `NewService` signature changed (added `*pgxpool.Pool` param), tests not updated |
| `internal/integration/radarr` | `decimal.Decimal.InexactFloat64()` no longer exists, mapper test broken |

### Unit Tests: 4 packages failing

| Package | Failure | Root Cause |
|---------|---------|------------|
| `service/user` | Build fail | `NewService` takes 5 args now, tests pass 4 |
| `service/auth` | Panic (nil pool) | `Register` now calls `pool.Begin()`, test passes nil pool |
| `service/session` | 2 test failures | Mock expectations stale after `RefreshSession` refactor |
| `integration/radarr` | Build fail | `InexactFloat64()` method removed from decimal.Decimal |

**All other packages pass.** These 4 failures share a root cause: production code was refactored but tests weren't updated.

### Code Metrics

| Metric | Value |
|--------|-------|
| Source files (excl. ogen) | 262 |
| Test files (excl. ogen) | 144 |
| Source LOC | 72,017 |
| Test LOC | 62,906 |
| Test:source ratio | 0.87 |
| Dependencies (go.mod) | 154 |
| Migrations | 32 (64 files) |
| OpenAPI spec | 8,889 lines, 124 endpoints |

---

## 2. Implementation State

### Fully Implemented (27 components)

| Component | Files | LOC | Tests | fx Module |
|-----------|-------|-----|-------|-----------|
| **Entry & Wiring** |
| cmd/revenge | 2 | 171 | - | - |
| internal/app | 1 | ~100 | - | app.Module |
| **Config & API** |
| config | 6 | 1,332 | 3 | config.Module |
| api (handlers) | 41 | 14,383 | 14 | api.Module |
| **Core Services** |
| auth | 7 | 1,784 | 6 | auth.Module |
| session | 6 | 762 | 5 | session.Module |
| mfa | 5 | 1,423 | 4 | mfa.Module |
| oidc | 4 | 1,521 | 4 | oidc.Module |
| user | 5 | 1,188 | 5 | user.Module |
| rbac | 6 | 1,374 | 6 | rbac.Module |
| apikeys | 4 | 497 | 4 | apikeys.Module |
| metadata | 18 | 8,087 | 0* | metadatafx.Module |
| library | 6 | 1,589 | 3 | library.Module |
| search | 4 | 982 | 1 | search.Module |
| settings | 5 | 748 | 4 | settings.Module |
| activity | 6 | 1,072 | 7 | activity.Module |
| email | 2 | 414 | 1 | email.Module |
| notification | 6 | 2,134 | 4 | notification.Module |
| storage | 4 | 518 | 1 | storage.Module |
| **Content Modules** |
| movie | 27 | 7,616 | 14 | movie.Module |
| tvshow | 22 | 9,247 | 5 | tvshow.Module |
| shared | 13 | 2,301 | 11 | - |
| **Infrastructure** |
| database | 19 | 7,869 | 8 | database.Module |
| cache | 4 | 870 | 8 | cache.Module |
| jobs | 6 | 763 | 5 | jobs.Module |
| health | 4 | 418 | 3 | health.Module |
| image | 2 | 432 | 1 | image.Module |
| logging | 2 | 194 | 1 | logging.Module |
| observability | 5 | 548 | 1 | observability.Module |
| search infra | 1 | 240 | 1 | search.Module |
| **Integrations** |
| radarr | 9 | 2,182 | 3 | radarr.Module |
| sonarr | 9 | 2,613 | 1 | sonarr.Module |
| **Integration Tests** |
| tests/integration | 13 | 3,785 | 13 | - |

*metadata service has 0 unit tests but is covered by integration tests

### Partially Implemented (2 components)

| Component | State | Details |
|-----------|-------|---------|
| qar (adult content) | Schema only | 4 files (424 LOC), DB models generated, no service logic |
| raft (clustering) | Scaffold | 2 files (325 LOC), module exists, not operational |

### Not Started

Nothing critical is missing. All planned v0.2.0 services are built.

---

## 3. Documentation Alignment

Spot-checked 14 docs against code:

| Check | Result |
|-------|--------|
| Perfect match | 10/14 (71%) |
| Minor deviation | 2/14 (14%) |
| Major mismatch | 1/14 (7%) |
| Medium issue | 1/14 (7%) |

### Issues Found

#### HIGH: JOBS.md describes wrong architecture
- **Doc claims**: 17 workers across 5 priority queues (critical/high/default/low/bulk)
- **Reality**: 11 workers (4 movie + 5 tvshow + 2 infra) on a single default queue
- **Action**: Rewrite JOBS.md queue/worker section

#### MEDIUM: ARCHITECTURE.md missing 3 fx modules
- **Missing**: health.Module, observability.Module, raft.Module
- **Action**: Add to module list

#### LOW: METADATA_SYSTEM.md method count
- **Doc claims**: 31 methods on Service interface
- **Reality**: 27 unique methods (4 overcounted)
- **Action**: Correct count

---

## 4. What's Needed Before Continuing Development

### Must Fix (blocks development confidence)

1. **Fix 4 broken test files** (~30 min work)
   - `service/user/service_test.go` + `cached_service_test.go` — add pool param
   - `service/auth/service_exhaustive_test.go` — mock pool or restructure
   - `service/session/service_exhaustive_test.go` — update mock expectations
   - `integration/radarr/mapper_test.go` — fix decimal method call

2. **Fix JOBS.md** (~15 min)
   - Correct worker count (11, not 17)
   - Remove fictional 5-tier queue system, document single-queue reality

### Should Fix (improves accuracy)

3. **Add 3 missing modules to ARCHITECTURE.md** (health, observability, raft)
4. **Correct method count in METADATA_SYSTEM.md** (27, not 31)
5. **Add metadata service unit tests** (8,087 LOC with 0 unit tests is a risk)

### Nice to Have

6. **Run `make lint`** to check for linter issues (requires golangci-lint v2.8.0)
7. **Run integration tests** (`make test-integration`) to verify full stack

---

## 5. What Can You Build Next?

The project is at **v0.2.0 (Core)** milestone. Based on the roadmap, the next logical steps are:

### Option A: Continue to v0.3.0 (MVP - Movies)
The movie module is already implemented. What's left for a real MVP:
- Frontend (SvelteKit) — no frontend exists yet
- Movie playback/streaming — transcoding not built
- Library scanning refinement — exists but may need polish
- User onboarding flow — auth works, needs setup wizard

### Option B: Harden What Exists
- Fix the 4 broken tests
- Add metadata service unit tests
- Run full CI pipeline (`make ci`)
- Load testing / performance baseline
- Security audit (all auth flows)

### Option C: Build QAR Module
- Service layer for adult content (schema exists)
- Whisparr integration
- StashDB provider

### Option D: Start Frontend
- SvelteKit skeleton
- API client generation from OpenAPI spec
- Auth flow UI
- Movie browsing UI

---

## 6. Architecture Strengths

- **Clean DI**: 32 fx modules, proper lifecycle management
- **Repository pattern**: Every service has interface + postgres impl
- **Caching layers**: sync.Map (L0) + otter (L1) + rueidis/Dragonfly (L2)
- **Test infrastructure**: testcontainers for integration, mockery for unit
- **CI/CD ready**: GitHub Actions, Docker multi-stage, Helm chart
- **OpenAPI-first**: 8,889-line spec → ogen codegen → type-safe handlers

## 7. Architecture Risks

- **Metadata service**: 8K LOC with 0 unit tests — any refactor is risky
- **sync.Map caching**: Used in TMDb/TVDb/Radarr/Sonarr clients instead of the proper L1/L2 cache infra (tracked in CODE_ISSUES.md)
- **No frontend**: Entire backend exists but no UI — can't demo anything
- **RAFT not integrated**: Module exists but isn't tested or operational

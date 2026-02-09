# Deep Scan Report

**Date**: 2026-02-07
**Branch**: develop
**Scope**: Full Go codebase (excluding generated ogen code, vendor, working dirs)

## Executive Summary

- **API Coverage**: 160/160 endpoints implemented (100%)
- **Real Stubs**: 0 blocking stubs found
- **Non-functional Code**: 3 config fields never used, 2 unused deps
- **Known Security Issues**: 6 pending (documented in `.workingdir/TODO_A7_SECURITY_FIXES.md`)
- **TODOs in Code**: Only 3 non-test files have TODO comments

---

## 1. Stubs & Placeholders

### NONE BLOCKING

No function stubs or placeholder implementations that block functionality were found. All handler methods have real implementations.

### Intentional Placeholders (OK - by design)

| File | What | Why |
|------|------|-----|
| `internal/content/movie/db/placeholder.sql.go` | `SELECT 1` query | sqlc requires non-empty query dir, real queries in shared schema |
| `internal/content/tvshow/db/placeholder.sql.go` | `SELECT 1` query | Same reason |
| `internal/content/qar/db/placeholder.sql.go` | `SELECT 1` query | QAR module planned for v0.3.0 |
| `internal/content/movie/mediainfo_windows.go` | Media prober returns error | FFmpeg/CGO not available on Windows - platform stub |

---

## 2. TODO/FIXME Comments in Code

Only 3 non-test Go files have TODO/FIXME:

| File | Line | Comment | Severity |
|------|------|---------|----------|
| `internal/api/localization.go` | 15 | `// TODO: Implement when user settings are available` - user language from context | LOW (fallback to Accept-Language works) |
| `internal/config/config.go` | 194 | `// Optional for v0.1.0 (auth not implemented yet)` - outdated comment, auth IS implemented | LOW (misleading comment) |

Test files with TODOs (not actionable bugs):

| File | Line | Comment |
|------|------|---------|
| `tests/integration/health_test.go` | 81 | `t.Skip("Health service doesn't currently detect DB failures")` |
| `tests/integration/search/search_test.go` | 47 | `t.Skip("Health check endpoint has issues with typesense-go v2 client")` |

---

## 3. Non-functional / Dead Code

### Config Fields Defined but Never Used by Runtime Code

These `JobsConfig` fields are defined, have defaults, and are tested, but are never read by the actual River worker infrastructure:

| Field | Defined In | Used In | Problem |
|-------|-----------|---------|---------|
| `Jobs.FetchPollInterval` | `config.go:170` | Only tests + defaults | River not configured with this value |
| `Jobs.RescueStuckJobsAfter` | `config.go:173` | Only tests + defaults | River not configured with this value |
| `Jobs.MaxWorkers` | `config.go:161` | Only tests + defaults | River not configured with this value |
| `Jobs.FetchCooldown` | `config.go:167` | Only tests + defaults | River not configured with this value |

### Unused Dependencies

From `TODO.md`, these deps are in `go.mod` but unused in code:

| Dependency | Purpose | Status |
|-----------|---------|--------|
| `gobreaker` | Circuit breaker | In go.mod, never imported |
| `sturdyc` | Request coalescing/caching | In go.mod, never imported |

---

## 4. Known Security Issues (Previously Documented)

These are tracked in `.workingdir/TODO_A7_SECURITY_FIXES.md` and remain **unfixed**:

| ID | Issue | Severity | Location |
|----|-------|----------|----------|
| A7.1 | Missing transaction boundaries (registration, avatar upload, session refresh) | CRITICAL | `auth/service.go`, `user/service.go`, `session/service.go` |
| A7.2 | Login timing attack enables username enumeration | CRITICAL | `auth/service.go:236-268` |
| A7.3 | Goroutine leak in notification dispatcher | HIGH | `notification/dispatcher.go:116-137` |
| A7.4 | Password reset token returned (info disclosure) | HIGH | `auth/service.go:521-562` |
| A7.5 | No service-level rate limiting for Argon2id | MEDIUM | `auth/service.go` |
| A7.6 | `context.Background()` in goroutines loses cancellation | MEDIUM | `apikeys/service.go:185-192` |

---

## 5. API Endpoint Coverage

**160/160 endpoints implemented (100%)**

All endpoints in the OpenAPI spec at `api/openapi/openapi.yaml` have real handler implementations. No endpoint falls through to the `UnimplementedHandler` stubs.

Breakdown by domain:
- Health: 3
- Auth & MFA: 20 (incl. 7 new WebAuthn)
- Movies & Collections: 22
- TV Shows & Episodes: 26
- Metadata & Images: 8
- Search: 4
- Settings: 7
- Users: 6
- Sessions: 6
- RBAC: 12
- API Keys: 4
- OIDC: 14
- Activity Logs: 5
- Libraries: 10
- Integrations: 10
- WebAuthn: 7

---

## 6. Bugs Found

### BUG-1: Health check does not detect database failures

**File**: `internal/infra/health/service.go`
**Impact**: Kubernetes readiness probe won't fail when DB is down
**Evidence**: Integration test explicitly skipped with note "Health service doesn't currently detect DB failures"

### BUG-2: Typesense health check broken with v2 client

**File**: `tests/integration/search/search_test.go`
**Impact**: Search service health not monitored
**Evidence**: Test skipped: "Health check endpoint has issues with typesense-go v2 client"

### BUG-3: CleanupWorker not registered with River

**File**: `internal/infra/jobs/cleanup_job.go`
**Impact**: Auth token/session cleanup job is defined but never scheduled
**Evidence**: Found in previous session (.workingdir3), `CleanupWorker` exists but `river.AddWorker()` is never called for it

### BUG-4: Outdated auth config comment

**File**: `internal/config/config.go:194`
**Impact**: Misleading - says "auth not implemented yet" but auth is fully implemented
**Fix**: Remove the comment "Optional for v0.1.0 (auth not implemented yet)"

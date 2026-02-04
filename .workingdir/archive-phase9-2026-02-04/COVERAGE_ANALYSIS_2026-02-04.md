# Coverage Analysis Report

**Date**: 2026-02-04
**Total Coverage**: 4.13%
**Target**: 80%
**Gap**: 75.87 percentage points

---

## Setup Summary

✅ **CGO + FFmpeg Setup Complete**:
- MSYS2 installed via winget
- MinGW-w64 toolchain (GCC 15.2.0)
- FFmpeg 8.0.1 + all dev libraries
- Go 1.25.6 (upgraded from 1.25.5)
- CGO_ENABLED=1
- go-astiav compilation working

✅ **Test Suite Execution**:
- All tests passed (exit code 0)
- Coverage profile generated: `coverage.out`

---

## Coverage by Category (Initial Assessment)

### 🔴 0% Coverage (No Tests)
**API Handlers** (all 0%):
- `handler.go` - Main API handlers (auth, settings, users)
- `handler_activity.go` - Activity logging handlers
- `handler_apikeys.go` - API key management
- `handler_library.go` - Library management
- `handler_metadata.go` - TMDb metadata proxy
- `handler_mfa.go` - MFA management
- `handler_oidc.go` - OIDC authentication
- `handler_radarr.go` - Radarr integration
- `handler_rbac.go` - RBAC management
- `handler_search.go` - Search endpoints

**Services** (mostly 0%):
- `cmd/revenge/main.go` - Main entry point
- `cmd/revenge/migrate.go` - Migration commands
- Context helpers (`context.go`)
- Most service implementations

### 🟢 100% Coverage (Well Tested)
- `internal/api/errors.go` - Error handling functions
- `internal/errors/` package

### 🟡 Partial Coverage
- `internal/api/middleware` - 82.6% ✅
- `internal/config` - 74.4% ✅
- `internal/crypto` - 84.7% ✅
- `internal/errors` - 100% ✅
- `internal/infra/cache` - 40.9%

---

## Priority Test Writing Plan

### Phase 1: Core Services (Target: 80%+ each)
**High Value, Foundation Layer**

1. **Session Service** (`internal/service/session/`)
   - ValidateSession, CreateSession, RevokeSession
   - Token generation and validation
   - Expiry handling

2. **Auth Service** (`internal/service/auth/`)
   - Login, Register, Logout
   - Password verification
   - MFA integration

3. **User Service** (`internal/service/user/`)
   - CRUD operations
   - Password changes
   - Email verification

4. **RBAC Service** (`internal/service/rbac/`)
   - Enforce, GetUserRoles
   - Policy management
   - Custom roles

5. **Settings Service** (`internal/service/settings/`)
   - Get/Set server settings
   - Get/Set user settings

---

### Phase 2: Content Services (Target: 70%+)
**Business Logic Layer**

6. **Movie Service** (`internal/content/movie/`)
   - CRUD operations
   - Watch progress
   - Search integration

7. **Library Service** (`internal/service/library/`)
   - Scan operations
   - Permission checks

8. **Search Service** (`internal/service/search/`)
   - Index operations
   - Query handling

9. **Activity Service** (`internal/service/activity/`)
   - Audit logging
   - Query operations

---

### Phase 3: Integration Services (Target: 60%+)
**External Integrations**

10. **Radarr Integration** (`internal/integration/radarr/`)
    - Sync operations
    - Webhook handling

11. **TMDb Service** (`internal/content/movie/metadata_service.go`)
    - Metadata fetching
    - Image proxying

12. **Notification Service** (`internal/service/notification/`)
    - Event dispatch
    - Agent delivery

---

### Phase 4: API Handlers (Target: 50%+)
**HTTP Layer** (after services are tested)

13. Handler tests for:
    - Auth endpoints
    - User endpoints
    - Library endpoints
    - Movie endpoints
    - RBAC endpoints

---

### Phase 5: Infrastructure (Target: 70%+)
**Supporting Components**

14. **Cache Integration** (`internal/infra/cache/`)
    - Current: 40.9% → Target: 70%
    - L1/L2 invalidation
    - Pattern matching

15. **Database Integration** (`internal/infra/database/`)
    - Connection pooling
    - Transaction handling
    - Migration system

16. **Job Queue** (`internal/infra/jobs/`)
    - Worker registration
    - Job execution
    - Error handling

---

## Test Writing Strategy

### Test Types Needed

**Unit Tests**:
- Mock dependencies (repository, cache, external APIs)
- Focus on business logic
- Fast execution (<100ms per test)

**Integration Tests**:
- Use testcontainers for PostgreSQL, Dragonfly, Typesense
- Test real database operations
- Test cache behavior
- Test job execution

**API Tests**:
- Use httptest for handlers
- Mock service layer
- Test error cases
- Test auth/authz

---

## Estimated Effort

| Phase | Packages | Est. Tests | Est. Hours |
|-------|----------|------------|------------|
| 1. Core Services | 5 | 150-200 | 12-16 |
| 2. Content Services | 4 | 100-120 | 8-12 |
| 3. Integration Services | 3 | 60-80 | 6-8 |
| 4. API Handlers | 10+ | 100-150 | 10-15 |
| 5. Infrastructure | 3 | 40-60 | 4-6 |
| **Total** | **25+** | **450-610** | **40-57h** |

---

## Known Issues / Blockers

None - CGO + FFmpeg setup complete, all tests passing.

---

## Next Actions

1. ✅ Start with **Session Service Tests**
   - Most fundamental service
   - Used by all authenticated endpoints
   - Clear test cases

2. ✅ Follow with **Auth Service Tests**
   - Depends on Session Service
   - Critical for security
   - MFA integration tests

3. ✅ Continue systematically through Phase 1-5

4. ✅ Monitor coverage after each package:
   ```bash
   go test -coverprofile=coverage.out ./internal/service/session/
   go tool cover -func=coverage.out | grep total
   ```

5. ✅ Re-run full suite after each phase to verify progress

---

## Coverage Tracking

**Current State**:
```
Total: 4.13%
├─ Middleware: 82.6% ✅
├─ Config: 74.4% ✅
├─ Crypto: 84.7% ✅
├─ Errors: 100% ✅
├─ Cache: 40.9%
└─ Everything else: ~0%
```

**Target After Each Phase**:
- Phase 1 complete: ~25% total
- Phase 2 complete: ~40% total
- Phase 3 complete: ~55% total
- Phase 4 complete: ~70% total
- Phase 5 complete: ~80%+ total ✅

---

**Status**: Ready to begin test writing
**Starting Point**: Session Service

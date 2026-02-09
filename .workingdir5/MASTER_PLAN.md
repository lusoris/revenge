# Plan: Integration Tests + Local Container Setup + Live Testing

## Context

Unit test coverage push is complete — 20 packages at 80%+, but 21 packages remain below 80%. Most are blocked by `repository_pg.go` / direct DB dependencies that need real PostgreSQL. The user wants three phases:

1. **Integration tests** with real PostgreSQL to push remaining packages past 80%
2. **Local container configuration** — fix Docker setup, get full stack running locally
3. **Real live end-to-end testing** — smoke tests against running containers

## Phase 1: Integration Tests (embedded-postgres)

### Strategy

Use the existing `testutil.NewTestDB(t)` embedded-postgres infrastructure (fast, no Docker needed, ~10ms per test). Focus on service-level integration tests that exercise both service.go AND repository_pg.go code paths with real SQL.

### Target Packages (ordered by ROI)

| Package | Current | Realistic Target | Approach |
|---------|---------|-----------------|----------|
| auth | 50.4% | 75-80% | Test failed-login tracking (5 untested repo methods), service error paths |
| oidc | 41.0% | 70-80% | Service-level tests with real DB (encryption, state, callback flow) |
| library | 41.3% | 70-80% | Service-level scan lifecycle + permission edge cases |
| session | 71.8% | 80%+ | Edge cases: expired sessions, concurrent revocations |
| rbac | 61.9% | 75-80% | Add adapter.go integration tests (Casbin + PostgreSQL) |
| apikeys | 74.3% | 80%+ | Error paths: expired keys, revoked keys, duplicate prefixes |
| user | 53.7% | 60-65% | Test HardDeleteAvatar, service integration paths |
| mfa | 19.9% | 30-40% | Service tests using real db.Queries (no repo interface to mock) |

**Skip**: content/movie (33%, cgo/FFmpeg ceiling), content/tvshow (44.7%, no repo_pg in package)

### Files to Create/Modify

#### 1. `internal/service/auth/service_integration_test.go` (expand existing)
- Test `RecordFailedLoginAttempt` + `CountFailedLoginAttemptsByUsername` + `CountFailedLoginAttemptsByIP`
- Test `ClearFailedLoginAttemptsByUsername` + `DeleteOldFailedLoginAttempts`
- Test `InvalidateEmailVerificationTokensByEmail` + `DeleteVerifiedEmailTokens`
- Test full Login flow with lockout (threshold exceeded → lockout → clear)
- Test Register → VerifyEmail → Login transaction boundary
- **Setup**: Uses existing `TestMain` with `testutil.StopSharedPostgres()`

#### 2. `internal/service/oidc/service_integration_test.go` (new)
- Add `TestMain` with embedded-postgres cleanup
- Test `AddProvider` → `GetProvider` → `UpdateProvider` → `DeleteProvider` lifecycle
- Test `encryptSecret`/`decryptSecret` with real provider storage
- Test `GetAuthURL` → state creation → `HandleCallback` with mock HTTP (state validation)
- Test `LinkUser` → `ListUserLinks` → `UnlinkUser` lifecycle
- Test `CleanupExpiredStates`
- Test `SetDefaultProvider` → `GetDefaultProvider`
- Test `EnableProvider`/`DisableProvider` toggle
- **Constructor**: `NewService(repo, logger, callbackURL, encryptKey)` — encryptKey is 32-byte AES key

#### 3. `internal/service/library/service_integration_test.go` (new)
- Add `TestMain` with embedded-postgres cleanup
- Test full scan lifecycle: `TriggerScan` → `StartScan` → `UpdateScanProgress` → `CompleteScan`
- Test scan failure: `TriggerScan` → `StartScan` → `FailScan`
- Test scan cancellation: `TriggerScan` → `CancelScan`
- Test permission model: `GrantPermission` → `CheckPermission` → `CanAccess`/`CanDownload`/`CanManage`
- Test `ListAccessible` with mixed permissions
- Test `RevokePermission` → verify `CheckPermission` returns false
- Test `Delete` library cascades to scans and permissions
- **Constructor**: `NewService(repo, logger, activityLogger)` — activityLogger can be nop

#### 4. `internal/service/session/service_integration_test.go` (expand or new)
- Test `RevokeAllUserSessionsExcept` — current session survives
- Test `CountActiveUserSessions` after create/revoke cycle
- Test `DeleteExpiredSessions` + `DeleteRevokedSessions` cleanup
- Test `GetInactiveSessions` + `RevokeInactiveSessions` idle timeout
- Test `UpdateSessionActivity` updates last_active timestamp
- Test concurrent session creation (max_per_user boundary)

#### 5. `internal/service/rbac/adapter_test.go` (new)
- Add `TestMain` with embedded-postgres cleanup
- Create `Adapter` with real pool: `NewAdapter(pool)`
- Test `SavePolicy` → `LoadPolicy` roundtrip with Casbin model
- Test `AddPolicy` → verify in DB → `RemovePolicy` → verify removed
- Test `RemoveFilteredPolicy` with field index matching
- Test `SavePolicy` transaction atomicity (replaces all rules)
- Test concurrent policy modifications

#### 6. `internal/service/apikeys/service_integration_test.go` (new or expand)
- Test `ValidateKey` with expired key → error
- Test `ValidateKey` with revoked key → error
- Test `CreateAPIKey` → `ValidateKey` → `RevokeKey` → `ValidateKey` fails
- Test `UpdateScopes` → `ValidateKey` returns new scopes
- Test `DeleteExpiredAPIKeys` cleanup
- Test `CountUserAPIKeys` accuracy after create/delete cycles

#### 7. `internal/service/mfa/manager_integration_test.go` (new)
- Add `TestMain` with embedded-postgres cleanup
- Create real `db.Queries` from test pool
- Test `GetStatus` with no MFA configured → returns disabled
- Test `EnableMFA` → `RequiresMFA` returns true → `DisableMFA`
- Test `SetRememberDevice` → `GetRememberDeviceSettings` roundtrip
- Test `HasAnyMethod` with/without TOTP configured
- Test `RemoveAllMethods` clears everything
- **Note**: MFA uses `*db.Queries` directly, no Repository interface

### Test Pattern (reuse everywhere)

```go
// TestMain for each new test file
func TestMain(m *testing.M) {
    code := m.Run()
    testutil.StopSharedPostgres()
    os.Exit(code)
}

// Setup helper
func setupIntegrationService(t *testing.T) (*Service, testutil.DB) {
    t.Helper()
    testDB := testutil.NewTestDB(t)
    queries := db.New(testDB.Pool())
    repo := NewRepositoryPg(queries)
    svc := NewService(repo, zap.NewNop(), /* other deps */)
    return svc, testDB
}
```

### Execution: 5 Parallel Agents

| Agent | Packages | Est. Tests |
|-------|----------|-----------|
| 1 | auth (failed login + email verification + lockout flow) | 15-20 |
| 2 | oidc (full provider lifecycle + auth flow) | 15-20 |
| 3 | library (scan lifecycle + permissions) + apikeys (error paths) | 20-25 |
| 4 | session (edge cases) + rbac adapter (Casbin + PG) | 15-20 |
| 5 | mfa (manager with real DB) + user (remaining gaps) | 10-15 |

---

## Phase 2: Local Container Configuration

### Issues to Fix

1. **Missing `deploy/typesense.Dockerfile`** — `docker-compose.dev.yml:76` references it but doesn't exist. Fix: use official `typesense/typesense:27.1` image directly instead of custom Dockerfile.

2. **Port inconsistency** — config.go defaults to 8080, Docker uses 8096. Standardize.

3. **`docker-compose.dev.yml` verification** — ensure all services start, health checks pass, migrations run.

### New/Modified Files

#### 1. Fix `docker-compose.dev.yml`
- Replace `build: deploy/typesense.Dockerfile` with `image: typesense/typesense:27.1`
- Verify all `depends_on` health checks work
- Ensure `REVENGE_DATABASE_URL`, `REVENGE_CACHE_URL`, `REVENGE_SEARCH_URL` env vars are correct

#### 2. Add `make docker-local` target to Makefile
```makefile
docker-local: docker-build  ## Build and run full local stack
    docker compose -f docker-compose.dev.yml up -d
    @echo "Waiting for services..."
    @sleep 5
    @curl -sf http://localhost:8096/health/live && echo "Revenge is healthy!" || echo "Startup failed"
```

#### 3. Verify `docker-compose.dev.yml` services
- PostgreSQL 18 starts and accepts connections
- Dragonfly starts and responds to PING
- Typesense starts and responds to health check
- Revenge starts, runs migrations, passes health checks
- All volumes properly mounted

### Quick Verification
```bash
make docker-local
curl http://localhost:8096/health/live   # → 200
curl http://localhost:8096/health/ready  # → 200
make docker-down
```

---

## Phase 3: Live End-to-End Testing

### Approach

Create a Go test file with `-tags=live` that runs against a real running stack. Tests use the generated ogen client or raw HTTP.

### File: `tests/live/smoke_test.go`

```
//go:build live
```

**Test Categories:**

#### Health & Infrastructure
- `TestLive_HealthLive` — GET /health/live → 200
- `TestLive_HealthReady` — GET /health/ready → 200 (all deps up)

#### Auth Flow
- `TestLive_Register` — POST /api/v1/auth/register → 201
- `TestLive_Login` — POST /api/v1/auth/login → 200 + tokens
- `TestLive_RefreshToken` — POST /api/v1/auth/refresh → 200 + new tokens
- `TestLive_Logout` — POST /api/v1/auth/logout → 204
- `TestLive_LoginInvalidPassword` — POST /api/v1/auth/login → 401

#### User CRUD
- `TestLive_GetCurrentUser` — GET /api/v1/users/me → 200
- `TestLive_UpdateUser` — PUT /api/v1/users/me → 200
- `TestLive_GetUserPreferences` — GET /api/v1/users/me/preferences → 200

#### Library CRUD (admin)
- `TestLive_CreateLibrary` — POST /api/v1/libraries → 201
- `TestLive_ListLibraries` — GET /api/v1/libraries → 200
- `TestLive_TriggerScan` — POST /api/v1/libraries/{id}/scan → 202
- `TestLive_DeleteLibrary` — DELETE /api/v1/libraries/{id} → 204

#### Settings
- `TestLive_GetServerSettings` — GET /api/v1/admin/settings → 200
- `TestLive_UpdateServerSetting` — PUT /api/v1/admin/settings/{key} → 200

#### Search (if Typesense is up)
- `TestLive_SearchMovies` — GET /api/v1/search/movies?q=test → 200

### Configuration
```go
var (
    baseURL  = envOr("REVENGE_TEST_URL", "http://localhost:8096")
    adminUser = envOr("REVENGE_TEST_ADMIN_USER", "admin")
    adminPass = envOr("REVENGE_TEST_ADMIN_PASS", "changeme123!")
)
```

### Run Command
```bash
# Start stack
make docker-local

# Run live tests
GOEXPERIMENT=greenteagc,jsonv2 go test -tags=live -v ./tests/live/...

# Tear down
make docker-down
```

### Add Makefile target
```makefile
test-live: docker-local  ## Run live end-to-end tests against running stack
    @echo "Running live smoke tests..."
    GOEXPERIMENT=greenteagc,jsonv2 go test -tags=live -v -count=1 ./tests/live/...
```

---

## Execution Order

1. **Phase 1** — Integration tests (5 parallel agents, ~30 min)
   - Run full test suite after to verify
   - Measure final coverage
2. **Phase 2** — Fix docker-compose, add make targets (~15 min)
   - Verify `make docker-local` works
3. **Phase 3** — Write live smoke tests (~20 min)
   - Run `make test-live` end-to-end

## Verification

1. `make test` — all unit + integration tests pass
2. `make docker-local` — full stack starts successfully
3. `curl localhost:8096/health/ready` — returns 200
4. `make test-live` — all smoke tests pass
5. Coverage report shows improvement in auth, oidc, library, session, rbac, apikeys packages

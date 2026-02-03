# Service Layer Integration Testing Session

**Date**: 2026-02-03
**Phase**: Application Core Testing (Service Layer)
**Previous Phase**: Infrastructure Testing (Database, Cache, Search) - COMPLETED ✅

## Session Overview

After completing infrastructure testing (41/42 tests passing), user realized we only tested infrastructure components, not the actual **application business logic**. This session focuses on testing the 9 service modules with real database/cache/search dependencies.

## Service Modules to Test

1. **User Service** (`/internal/service/user/`) - ✅ COMPLETED (11/11 tests)
2. **Auth Service** (`/internal/service/auth/`) - ✅ COMPLETED (12/12 tests)
3. **Session Service** (`/internal/service/session/`) - ⏳ PENDING
4. **Settings Service** (`/internal/service/settings/`) - ⏳ PENDING
5. **RBAC Service** (`/internal/service/rbac/`) - ⏳ PENDING
6. **API Keys Service** (`/internal/service/apikeys/`) - ⏳ PENDING
7. **Library Service** (`/internal/service/library/`) - ⏳ PENDING
8. **Activity Service** (`/internal/service/activity/`) - ⏳ PENDING
9. **OIDC Service** (`/internal/service/oidc/`) - ⏳ PENDING

---

## Test Results Summary

### User Service (`tests/integration/service/user_service_test.go`)

**Status**: ✅ ALL PASSING (11/11 tests)
**Test File**: `/home/kilian/dev/revenge/tests/integration/service/user_service_test.go`
**Execution Time**: 0.667s
**Last Run**: 2026-02-03

#### Test Coverage

| Test Name | Status | Duration | Description |
|-----------|--------|----------|-------------|
| TestUserService_CreateUser | ✅ PASS | 0.22s | Create user with password hashing |
| TestUserService_CreateUserDuplicateUsername | ✅ PASS | 0.21s | Duplicate username validation |
| TestUserService_CreateUserDuplicateEmail | ✅ PASS | 0.22s | Duplicate email validation |
| TestUserService_GetUser | ✅ PASS | 0.22s | Get user by ID |
| TestUserService_GetUserByUsername | ✅ PASS | 0.22s | Get user by username |
| TestUserService_GetUserByEmail | ✅ PASS | 0.22s | Get user by email |
| TestUserService_UpdateUser | ✅ PASS | 0.23s | Update user fields |
| TestUserService_DeleteUser | ✅ PASS | 0.23s | Soft delete user |
| TestUserService_ListUsers | ✅ PASS | 1.05s | List users with pagination |
| TestUserService_PasswordHashing | ✅ PASS | 0.59s | Bcrypt password hashing |
| TestUserService_ConcurrentCreation | ✅ PASS | 0.23s | Concurrent user creation race |

#### Implementation Details

**Setup**:
- Real PostgreSQL database (not mocked)
- Connection URL: `postgres://revenge:revenge_dev_pass@localhost:5432/revenge`
- Uses `db.New(pool)` for queries
- Repository: `user.NewPostgresRepository(queries)`

**Key Findings**:
- Password hashing uses **bcrypt** with cost factor 12
- Username and email uniqueness enforced at service layer
- Concurrent creation properly handles race conditions
- Timestamps added to test data to prevent pollution between runs

---

### Auth Service (`tests/integration/service/auth_service_test.go`)

**Status**: ✅ ALL PASSING (12/12 tests)
**Test File**: `/home/kilian/dev/revenge/tests/integration/service/auth_service_test.go`
**Execution Time**: 2.675s
**Last Run**: 2026-02-03

#### Test Coverage

| Test Name | Status | Duration | Description |
|-----------|--------|----------|-------------|
| TestAuthService_Register | ✅ PASS | 0.03s | Register new user with Argon2id |
| TestAuthService_Login | ✅ PASS | 0.07s | Login with username |
| TestAuthService_LoginWithEmail | ✅ PASS | 0.05s | Login with email |
| TestAuthService_LoginWrongPassword | ✅ PASS | 0.04s | Wrong password rejection |
| TestAuthService_RefreshToken | ✅ PASS | 2.07s | JWT refresh token flow |
| TestAuthService_Logout | ✅ PASS | 0.08s | Logout and revoke token |
| TestAuthService_ChangePassword | ✅ PASS | 0.11s | Change password flow |
| TestAuthService_ChangePasswordWrongOldPassword | ✅ PASS | 0.05s | Wrong old password rejection |
| **TestPasswordCompatibility_UserServiceToAuthService** | ✅ PASS | 0.07s | **Critical: Cross-service compatibility** |
| **TestPasswordCompatibility_AuthServiceToUserService** | ✅ PASS | 0.04s | **Critical: Reverse compatibility** |
| TestAuthService_MultipleDeviceLogin | ✅ PASS | 0.09s | Multiple concurrent sessions |
| TestAuthService_LogoutAll | ✅ PASS | 0.13s | Revoke all user tokens |

#### Implementation Details

**Setup**:
- Real PostgreSQL database (shared with User Service)
- JWT TokenManager with HS256 signing
- 15-minute access token expiry
- 7-day refresh token expiry

**Key Findings**:
- **Bug #28 VERIFIED FIXED**: Cross-service password compatibility working ✅
- Argon2id hashing consistent across both services
- JWT tokens generated with username and user ID claims
- Refresh tokens stored as SHA-256 hashes
- Multi-device support via device fingerprints
- Token revocation working correctly

**Critical Tests**:
- `TestPasswordCompatibility_UserServiceToAuthService`: User created by User Service can login via Auth Service ✅
- `TestPasswordCompatibility_AuthServiceToUserService`: User registered via Auth Service can be verified by User Service ✅

These tests confirm that the shared crypto service successfully resolved the password hashing inconsistency!

---

## Bugs Found

### Bug #28: Password Hashing Inconsistency (CRITICAL) 🔴

**Status**: 🔍 DISCOVERED - NOT YET FIXED
**Severity**: HIGH
**Impact**: Security vulnerability

**Description**:
Two different password hashing algorithms are used in the codebase:

1. **User Service** (`internal/service/user/service.go` line 164):
   - Uses **bcrypt** with cost factor 12
   - Code: `bcrypt.GenerateFromPassword([]byte(password), 12)`

2. **Auth Service** (`internal/service/auth/service.go` line 56):
   - Uses **Argon2id** with default params
   - Code: `argon2id.CreateHash(req.Password, argon2id.DefaultParams)`

**Problem**:
- User created via User Service will have bcrypt hash
- User created via Auth Service (registration) will have Argon2id hash
- Login flow uses Argon2id verification only
- Users created by User Service **cannot log in**!

**Evidence**:
```go
// User Service (service/user/service.go:164)
func (s *Service) HashPassword(password string) (string, error) {
    hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    // ...
}

// Auth Service (service/auth/service.go:56)
func (s *Service) Register(ctx context.Context, req RegisterRequest) (*db.SharedUser, error) {
    passwordHash, err := argon2id.CreateHash(req.Password, argon2id.DefaultParams)
    // ...
}

// Auth Service Login (service/auth/service.go:171)
match, err := argon2id.ComparePasswordAndHash(password, user.PasswordHash)
```

**Affected Files**:
- `/internal/service/user/service.go` (lines 164-170, 175-180)
- `/internal/service/auth/service.go` (lines 56-60, 171-176)

**Recommended Fix**:
1. **Option A**: Standardize on Argon2id (per AUTH.md documentation)
   - Change User Service to use Argon2id
   - Update all password hashing/verification to use Argon2id

2. **Option B**: Support both algorithms with migration
   - Detect hash type by prefix
   - Migrate bcrypt users to Argon2id on next login

3. **Option C**: Use only bcrypt everywhere
   - Change Auth Service to use bcrypt
   - Update documentation

**Next Steps**:
- [ ] Verify AUTH.md specification (which algorithm is required?)
- [ ] Check if any existing users have bcrypt hashes
- [ ] Implement fix based on specification
- [ ] Add integration test for cross-service password compatibility
- [ ] Update documentation

---

## Questions & Issues

### Q1: Which Password Hashing Algorithm Should Be Used?

**Status**: 🔍 NEEDS CLARIFICATION
**Context**: Bug #28 discovered inconsistency

**Question**: Should the system use:
- A) Argon2id (as specified in AUTH.md)?
- B) Bcrypt (simpler, widely supported)?
- C) Support both with migration path?

**Blocker**: Cannot proceed with Auth Service testing until this is resolved

---

### Q2: Database Cleanup Between Tests

**Status**: ✅ RESOLVED
**Solution**: Add timestamps to all test data

**Original Issue**: Test data pollution caused failures on repeat runs
**Fix Applied**: All test data now includes `time.Now().UnixNano()` suffix for uniqueness

---

### Q3: Docker Services Not Running

**Status**: ✅ RESOLVED
**Solution**: Started services manually with `docker start`

**Original Issue**: Tests failed with "connection refused"
**Fix Applied**:
```bash
docker start revenge-postgres-dev revenge-dragonfly-dev revenge-typesense-dev
```

---

## Test Infrastructure

### Database Setup

**Direct Connection** (not using testutil.TestDB):
```go
const testDatabaseURL = "postgres://revenge:revenge_dev_pass@localhost:5432/revenge?sslmode=disable"

pool, err := pgxpool.New(ctx, testDatabaseURL)
queries := db.New(pool)
repo := user.NewPostgresRepository(queries)
svc := user.NewService(repo)
```

**Why not embedded Postgres?**
- Service tests need full integration environment
- Testing against real production-like database
- All migrations already applied
- Faster than creating isolated databases per test

### Data Isolation Strategy

**Timestamps**: Every test creates unique data with nanosecond timestamp suffix
```go
timestamp := time.Now().UnixNano()
username := fmt.Sprintf("testuser_%d", timestamp)
email := fmt.Sprintf("test_%d@example.com", timestamp)
```

**Pros**:
- No cleanup needed
- Tests can run in parallel
- No shared state between tests
- Fast execution

**Cons**:
- Database grows with test data
- Need periodic cleanup
- Cannot test exact duplicates

---

## Next Steps

### Immediate (Today)

1. **BLOCKER**: Resolve Bug #28 (password hashing inconsistency)
   - Check AUTH.md specification
   - Decide on standardization approach
   - Implement fix
   - Verify with tests

2. **Auth Service Tests**: Create integration tests for:
   - Registration flow
   - Login flow
   - Token generation/verification
   - Password reset
   - Email verification
   - Refresh token rotation

3. **Cross-Service Tests**: Test interactions between services
   - Auth Service creating user → User Service retrieving
   - Password compatibility between services

### Short-term (This Week)

4. **Session Service Tests**
5. **Settings Service Tests**
6. **RBAC Service Tests**

### Medium-term

7. **API Keys Service Tests**
8. **Library Service Tests**
9. **Activity Service Tests**
10. **OIDC Service Tests**

---

## Test Execution Commands

### Run All Service Tests
```bash
go test -tags=integration -v ./tests/integration/service/... -timeout 30s
```

### Run User Service Tests Only
```bash
go test -tags=integration -v ./tests/integration/service/... -run TestUserService -timeout 30s
```

### Run Specific Test
```bash
go test -tags=integration -v ./tests/integration/service/... -run TestUserService_CreateUser -timeout 30s
```

### Check Services Running
```bash
docker ps | grep -E '(postgres|dragonfly|typesense)'
```

### Start Services
```bash
docker start revenge-postgres-dev revenge-dragonfly-dev revenge-typesense-dev
```

---

## Session Timeline

1. **15:30** - Started User Service test creation
2. **15:45** - Fixed compilation errors (database setup, type names)
3. **16:00** - Started docker services (were stopped)
4. **16:05** - First test passed ✅
5. **16:10** - Fixed data pollution with timestamps
6. **16:15** - All 11 User Service tests passing ✅
7. **16:20** - **DISCOVERED BUG #28**: Password hashing inconsistency 🔴
8. **16:25** - Started Auth Service code review
9. **16:30** - Creating documentation (this file)

---

## Statistics

### Infrastructure Testing (Previous Session)
- Tests Created: 42
- Tests Passing: 41
- Pass Rate: 97.6%
- Bugs Found: 8 (5 fixed immediately)

### Service Layer Testing (Current Session)
- Services Tested: 1/9 (11%)
- Tests Created: 11
- Tests Passing: 11
- Pass Rate: 100%
- **Bugs Found: 1 (CRITICAL security issue)** 🔴

### Combined Session
- Total Tests: 53
- Total Passing: 52
- Total Pass Rate: 98.1%
- Total Bugs Found: 9
- **Critical Bugs**: 1 (password hashing)

---

## Notes

- User Service tests provide good template for other service tests
- Need to check password hashing in Auth Service before creating those tests
- Consider creating shared test helpers for common patterns
- Database connection string hardcoded - should come from env/config?
- Test data cleanup strategy needed for long-term maintenance

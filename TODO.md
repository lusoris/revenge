# Project TODO

## High Priority

### 🔴 BUG #29: Password Hash Migration (bcrypt → argon2id)
**Status**: Identified - blocks authentication  
**Impact**: All existing users cannot log in  
**Effort**: 3-4 hours

**Tasks**:
- [ ] Add bcrypt support to `VerifyPassword` for backward compatibility
- [ ] Add bcrypt dependency to go.mod
- [ ] Write tests for both bcrypt and argon2id formats
- [ ] Implement lazy migration (rehash on login)
- [ ] Add `UpdatePasswordHash` repository method
- [ ] Add migration metrics
- [ ] Test with existing bcrypt passwords from database

**See**: `.workingdir/BUG_29_PASSWORD_HASH_MIGRATION.md`

---

### 🟡 River Job Worker Implementation
**Status**: Currently a stub  
**Impact**: Background jobs not processing  
**Effort**: 8-16 hours

**Current State**: `internal/infra/jobs/` has stub implementation:
```
16:14:40 INF job workers started (stub)
```

**Tasks**:
- [ ] Design job worker architecture
- [ ] Implement River job worker integration
- [ ] Add job queue persistence (PostgreSQL)
- [ ] Implement job types:
  - [ ] Email sending
  - [ ] Media processing
  - [ ] Library scanning
  - [ ] Thumbnail generation
  - [ ] Metadata fetching
- [ ] Add job monitoring and metrics
- [ ] Add job retry logic and dead letter queue
- [ ] Write integration tests
- [ ] Document job system architecture

**Research Needed**:
- River library API and best practices
- Job priority and scheduling
- Worker pool sizing
- Error handling strategies

---

## Medium Priority

### 🔵 Security: Fix G602 Slice Bounds Issues
**Status**: Identified  
**Count**: 10 issues  
**Location**: `internal/api/ogen/oas_router_gen.go` (generated code)  
**Effort**: 2-4 hours

**Tasks**:
- [ ] Review Ogen router generated code
- [ ] Add slice bounds checking with `validate.ValidateSliceIndex`
- [ ] Test with edge cases
- [ ] Re-run gosec to verify fixes

---

### 🔵 Security: Suppress G101 False Positives
**Status**: Identified  
**Count**: 43 issues  
**Location**: `internal/infra/database/db/*.sql.go` (SQLC-generated)  
**Effort**: 1 hour

**Tasks**:
- [ ] Add `#nosec G101` comments to SQLC-generated code
- [ ] Add justifications (e.g., "SQL query name, not actual credentials")
- [ ] Document suppression rationale
- [ ] Re-run gosec to verify clean scan

---

### 🔵 Security: Fix G204 Subprocess Issue
**Status**: Identified  
**Count**: 1 issue  
**Location**: `internal/testutil/testdb.go:292`  
**Effort**: 30 minutes

**Tasks**:
- [ ] Review subprocess call in test utilities
- [ ] Sanitize input if needed
- [ ] Add input validation
- [ ] Test with edge cases

---

## Completed ✅

### ✅ Security: Integer Overflow Fixes (G115)
**Completed**: 2026-02-03  
**Result**: 14/14 vulnerabilities fixed  

- ✅ Created `internal/validate` package with safe conversions
- ✅ Fixed database pool configuration (3 issues)
- ✅ Fixed API handler pagination (11 issues)
- ✅ Added comprehensive tests (100% coverage)
- ✅ Security scan: 68 → 54 issues (-20.6%)
- ✅ All tests passing (27 test files)
- ✅ Live system tested and verified

**See**: `.workingdir/security-fixes-summary.md`

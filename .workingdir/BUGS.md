# Bug Tracker

## Fixed Bugs

### BUG-002: User service tests expected bcrypt instead of Argon2id (2026-02-03)

**Problem**: 
User service tests failed with `Should be true` when checking password hash prefix.

**Root Cause**:
Tests in `internal/service/user/service_test.go` expected password hashes to start with 
`$2a$` (bcrypt format), but the crypto package was changed to use Argon2id which 
produces hashes starting with `$argon2id$`.

**Solution**:
Updated test assertions to check for `$argon2id$` prefix instead of `$2a$`.

**Files Changed**:
- `internal/service/user/service_test.go` - Lines 49, 373

---

### BUG-001: Testutil used duplicated/outdated migrations (2026-02-03)

**Problem**: 
Tests in `internal/api/` failed with `ERROR: column "mfa_verified" does not exist`.

**Root Cause**:
`internal/testutil/testdb_migrate.go` used `//go:embed migrations/*.sql` which embedded 
a local copy of migrations from `internal/testutil/migrations/`. This copy was outdated
(only 15 migrations) while the real migrations in `migrations/` had 26 files including
the MFA migration (000020).

**Why it happened**:
`go:embed` can only embed files from the same package directory or below. It cannot access
files from parent directories like `../../migrations/`. So a copy was created, which
got out of sync.

**Solution**:
Changed `testdb_migrate.go` to use `runtime.Caller()` to find the project root at runtime,
then load migrations from `migrations/` using `file://` source instead of embedded FS.

**Files Changed**:
- `internal/testutil/testdb_migrate.go` - Use dynamic path resolution
- `internal/testutil/migrations/` - **DELETED** (no longer needed)

**Lesson Learned**:
Never duplicate files that need to stay in sync. Always use a single source of truth.

---

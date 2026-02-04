# Bug Tracker

## Fixed Bugs

### BUG-004: Migration test hardcoded to check wrong table (2026-02-03)

**Problem**:
`TestMigrationsUpDown/MigrateDown` failed with `Should be false` / `Should be true`.

**Root Cause**:
Test was hardcoded to check for `user_avatars` table (migration 000007) as the newest
migration, but the project now has 26 migrations with `movie_watched` (000026) being
the latest.

**Solution**:
Updated test to check for `movie_watched` table in `public` schema (000026) and verify
`movie_genres` (000025) still exists after stepping down one migration.

**Files Changed**:
- `internal/infra/database/migrations_test.go` - Updated MigrateDown subtest

---

### BUG-003: Config required validation on optional Movie fields (2026-02-03)

**Problem**:
Config tests failed with validation errors on `Movie.TMDb.APIKey` and `Movie.Library.Paths`.

**Root Cause**:
These fields had `validate:"required"` and `validate:"required,min=1"` tags, but they
are optional features. The Movie service is only used when explicitly enabled, so having
an API key shouldn't be required at startup.

**Solution**:
Removed `required` validation tags from `TMDbConfig.APIKey` and `LibraryConfig.Paths`.

**Files Changed**:
- `internal/config/config.go`

---

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

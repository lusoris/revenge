# RBAC Service Bug

## Bug: Casbin policies persist across test runs

**Location**: `internal/service/rbac/adapter.go`

**Issue**: The Casbin adapter uses a shared table (`shared.casbin_rule`) that persists policies across test runs. When running tests in parallel, policies from previous tests appear in later tests, causing assertion failures.

**Evidence**:
- `TestService_GetPolicies` expects 3 policies but gets 10 (includes 7 from previous tests)
- `TestService_SavePolicy` expects 2 policies but gets 9 (includes 7 from previous tests)
- Tests pass if run in isolation but fail when run together

**Root Cause**: The TestDB creates separate databases for each test, but Casbin's LoadPolicy is called during enforcer initialization in `setupTestService`, which loads ALL policies from the shared table across all test databases created from the same template.

**Solution**:
1. Clear the casbin_rule table in setupTestService before creating the enforcer, OR
2. Don't call LoadPolicy automatically during enforcer creation (Casbin loads on first use), OR
3. Use a unique table name per test

**Fix Applied**: Clear table in setupTestService before enforcer initialization.

**Status**: FIXED

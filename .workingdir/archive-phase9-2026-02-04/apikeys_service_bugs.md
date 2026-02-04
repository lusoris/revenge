# API Keys Service Bugs

## Bug 1: ValidateKey cannot distinguish between "not found" and "inactive"

**Discovered:** 2026-02-03 during test creation

**Problem:**
The `ValidateKey` method uses `GetAPIKeyByHash` which has `AND is_active = true` in the SQL query. When a revoked/inactive key is validated, the query returns no rows, causing the service to return `ErrKeyNotFound` instead of `ErrKeyInactive`.

This means:
- Users cannot tell the difference between an invalid key and a revoked key
- Error messages are misleading
- API consumers cannot provide better UX (e.g., "This key was revoked" vs "Invalid key")

**Expected behavior:**
- `ValidateKey` should return `ErrKeyInactive` for revoked keys
- `ValidateKey` should return `ErrKeyNotFound` for non-existent keys

**Actual behavior:**
- Both cases return `ErrKeyNotFound`

**Fix:**
Remove the `AND is_active = true` filter from `GetAPIKeyByHash` SQL query, and let the service layer check `is_active` and return the appropriate error.

### SQL Changes Required:

**File:** `internal/infra/database/queries/shared/apikeys.sql`

```sql
-- BEFORE:
-- name: GetAPIKeyByHash :one
SELECT * FROM shared.api_keys
WHERE key_hash = $1 AND is_active = true;

-- AFTER:
-- name: GetAPIKeyByHash :one
SELECT * FROM shared.api_keys
WHERE key_hash = $1;
```

Then regenerate sqlc: `go run github.com/sqlc-dev/sqlc/cmd/sqlc generate`

The service code already has the check, so it will work correctly once the SQL is fixed.

**Impact:**
- Better error messages for API consumers
- Ability to distinguish revoked keys from invalid keys in logs/metrics
- More accurate error handling

**Status:** FIXING NOW

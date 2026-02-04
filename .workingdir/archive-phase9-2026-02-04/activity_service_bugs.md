# Activity Service Bugs

## Bug 1: Repository tests fail due to foreign key constraint

**Status:** FOUND - Test Design Issue

**Description:**
The activity_log table has a foreign key constraint `user_id REFERENCES shared.users(id)`, but tests were using random UUIDs that don't exist in the users table.

**Error:**
```
ERROR: insert or update on table "activity_log" violates foreign key constraint "activity_log_user_id_fkey" (SQLSTATE 23503)
```

**Fix:**
Tests need to either:
1. Create actual users first, OR
2. Use NULL for user_id (field is nullable), OR
3. Test with system activities that don't require a user

**Solution:** Use NULL for user_id in tests, or create actual users when testing user-related activities.

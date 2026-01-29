-- Rollback RBAC: Remove role column from users table

-- Remove index
DROP INDEX IF EXISTS idx_users_role;

-- Remove role column
ALTER TABLE users DROP COLUMN IF EXISTS role;

-- Drop enum type
DROP TYPE IF EXISTS user_role;

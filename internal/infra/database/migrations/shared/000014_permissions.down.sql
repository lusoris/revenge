-- Rollback permissions tables

-- Drop indexes
DROP INDEX IF EXISTS idx_role_permissions_role;
DROP INDEX IF EXISTS idx_permissions_category;

-- Drop tables
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;

-- Drop enum type
DROP TYPE IF EXISTS permission_category;

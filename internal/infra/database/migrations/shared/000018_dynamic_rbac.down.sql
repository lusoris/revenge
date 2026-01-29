-- Dynamic RBAC: Rollback to static ENUM-based roles
BEGIN;

-- Remove role_id from users
ALTER TABLE users DROP COLUMN IF EXISTS role_id;

-- Drop new tables
DROP TABLE IF EXISTS permission_definitions;
DROP TABLE IF EXISTS roles;

-- Note: Old role_permissions and permissions tables from 000013/000014 are kept
-- They will be restored as the active RBAC system after rollback

COMMIT;

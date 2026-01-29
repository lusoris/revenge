-- name: ListRoles :many
SELECT id, name, display_name, description, color, icon, is_system, is_default, priority, created_at, updated_at, created_by
FROM roles
ORDER BY priority DESC, name ASC;

-- name: GetRoleByID :one
SELECT id, name, display_name, description, color, icon, is_system, is_default, priority, created_at, updated_at, created_by
FROM roles
WHERE id = $1;

-- name: GetRoleByName :one
SELECT id, name, display_name, description, color, icon, is_system, is_default, priority, created_at, updated_at, created_by
FROM roles
WHERE name = $1;

-- name: GetDefaultRole :one
SELECT id, name, display_name, description, color, icon, is_system, is_default, priority, created_at, updated_at, created_by
FROM roles
WHERE is_default = true
LIMIT 1;

-- name: CreateRole :one
INSERT INTO roles (name, display_name, description, color, icon, priority, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, name, display_name, description, color, icon, is_system, is_default, priority, created_at, updated_at, created_by;

-- name: UpdateRole :one
UPDATE roles
SET display_name = $2, description = $3, color = $4, icon = $5, priority = $6, updated_at = NOW()
WHERE name = $1
RETURNING id, name, display_name, description, color, icon, is_system, is_default, priority, created_at, updated_at, created_by;

-- name: DeleteRole :exec
DELETE FROM roles WHERE name = $1 AND is_system = false;

-- name: SetDefaultRole :exec
UPDATE roles SET is_default = (name = $1);

-- name: CountUsersWithRole :one
SELECT COUNT(*) FROM users WHERE role_id = $1;

-- name: GetUserRoleName :one
SELECT r.name
FROM users u
JOIN roles r ON r.id = u.role_id
WHERE u.id = $1;

-- name: SetUserRole :exec
UPDATE users SET role_id = (SELECT id FROM roles WHERE name = $2)
WHERE users.id = $1;

-- name: ListPermissionDefinitions :many
SELECT id, name, display_name, description, category, is_dangerous, created_at
FROM permission_definitions
ORDER BY category, name;

-- name: GetPermissionDefinitionsByCategory :many
SELECT id, name, display_name, description, category, is_dangerous, created_at
FROM permission_definitions
WHERE category = $1
ORDER BY name;

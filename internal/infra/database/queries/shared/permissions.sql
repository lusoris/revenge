-- RBAC Permission queries

-- name: GetAllPermissions :many
SELECT * FROM permissions ORDER BY category, name;

-- name: GetPermissionByID :one
SELECT * FROM permissions WHERE id = $1;

-- name: GetPermissionByName :one
SELECT * FROM permissions WHERE name = $1;

-- name: GetPermissionsByCategory :many
SELECT * FROM permissions WHERE category = $1 ORDER BY name;

-- name: GetPermissionsForRole :many
SELECT p.* FROM permissions p
INNER JOIN role_permissions rp ON p.id = rp.permission_id
WHERE rp.role = $1
ORDER BY p.category, p.name;

-- name: GetPermissionNamesForRole :many
SELECT p.name FROM permissions p
INNER JOIN role_permissions rp ON p.id = rp.permission_id
WHERE rp.role = $1
ORDER BY p.name;

-- name: UserHasPermission :one
SELECT EXISTS(
    SELECT 1 FROM role_permissions rp
    INNER JOIN permissions p ON p.id = rp.permission_id
    INNER JOIN users u ON u.role = rp.role
    WHERE u.id = $1 AND p.name = $2
);

-- name: UserHasAnyPermission :one
SELECT EXISTS(
    SELECT 1 FROM role_permissions rp
    INNER JOIN permissions p ON p.id = rp.permission_id
    INNER JOIN users u ON u.role = rp.role
    WHERE u.id = $1 AND p.name = ANY($2::text[])
);

-- name: GetUserPermissions :many
SELECT p.* FROM permissions p
INNER JOIN role_permissions rp ON p.id = rp.permission_id
INNER JOIN users u ON u.role = rp.role
WHERE u.id = $1
ORDER BY p.category, p.name;

-- name: GetUserPermissionNames :many
SELECT p.name FROM permissions p
INNER JOIN role_permissions rp ON p.id = rp.permission_id
INNER JOIN users u ON u.role = rp.role
WHERE u.id = $1
ORDER BY p.name;

-- name: AddRolePermission :exec
INSERT INTO role_permissions (role, permission_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RemoveRolePermission :exec
DELETE FROM role_permissions
WHERE role = $1 AND permission_id = $2;

-- name: GetRolesForPermission :many
SELECT role FROM role_permissions
WHERE permission_id = $1
ORDER BY role;

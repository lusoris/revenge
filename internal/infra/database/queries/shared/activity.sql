-- Activity Log (uses new partitioned schema from migration 000021)
-- Note: Prefer using audit_log.sql queries for audit-specific operations

-- name: CreateActivityLog :one
INSERT INTO activity_log (
    user_id, action, module, entity_id, entity_type, changes, ip_address, user_agent
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: ListActivityLogByUser :many
SELECT * FROM activity_log
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListActivityLogByAction :many
SELECT * FROM activity_log
WHERE action = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListActivityLogByModule :many
SELECT * FROM activity_log
WHERE module = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListRecentActivity :many
SELECT * FROM activity_log
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: DeleteOldActivityLogs :exec
DELETE FROM activity_log
WHERE created_at < $1;

-- Server Settings
-- name: GetServerSetting :one
SELECT * FROM server_settings WHERE key = $1;

-- name: GetServerSettingsByCategory :many
SELECT * FROM server_settings WHERE category = $1 ORDER BY key ASC;

-- name: GetPublicServerSettings :many
SELECT * FROM server_settings WHERE is_public = true ORDER BY key ASC;

-- name: ListServerSettings :many
SELECT * FROM server_settings ORDER BY category, key ASC;

-- name: UpsertServerSetting :one
INSERT INTO server_settings (key, value, category, description, is_public)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (key) DO UPDATE SET
    value = $2,
    category = COALESCE($3, server_settings.category),
    description = COALESCE($4, server_settings.description),
    is_public = COALESCE($5, server_settings.is_public),
    updated_at = NOW()
RETURNING *;

-- name: DeleteServerSetting :exec
DELETE FROM server_settings WHERE key = $1;

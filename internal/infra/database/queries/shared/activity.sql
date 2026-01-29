-- Activity Log
-- name: CreateActivityLog :one
INSERT INTO activity_log (
    user_id, profile_id, action, module, item_id, item_type,
    details, ip_address, user_agent, severity
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
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

-- name: ListServerSettings :many
SELECT * FROM server_settings ORDER BY key ASC;

-- name: UpsertServerSetting :one
INSERT INTO server_settings (key, value, description, updated_by)
VALUES ($1, $2, $3, $4)
ON CONFLICT (key) DO UPDATE SET
    value = $2,
    description = COALESCE($3, server_settings.description),
    updated_by = $4,
    updated_at = NOW()
RETURNING *;

-- name: DeleteServerSetting :exec
DELETE FROM server_settings WHERE key = $1;

-- name: CreateAuditLog :one
INSERT INTO activity_log (user_id, action, module, entity_id, entity_type, changes, ip_address, user_agent)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetAuditLogByID :one
SELECT * FROM activity_log WHERE id = $1 AND created_at = $2;

-- name: ListAuditLogsByEntity :many
SELECT * FROM activity_log
WHERE module = $1 AND entity_type = $2 AND entity_id = $3
ORDER BY created_at DESC
LIMIT $4 OFFSET $5;

-- name: ListAuditLogsByUser :many
SELECT * FROM activity_log
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListAuditLogsByAction :many
SELECT * FROM activity_log
WHERE action = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListAuditLogsByModule :many
SELECT * FROM activity_log
WHERE module = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountAuditLogsByEntity :one
SELECT COUNT(*) FROM activity_log
WHERE module = $1 AND entity_type = $2 AND entity_id = $3;

-- name: ListRecentAuditLogs :many
SELECT * FROM activity_log
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: DeleteOldAuditLogs :exec
DELETE FROM activity_log WHERE created_at < $1;

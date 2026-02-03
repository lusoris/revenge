-- name: CreateActivityLog :one
-- Insert a new activity log entry
INSERT INTO public.activity_log (
    user_id,
    username,
    action,
    resource_type,
    resource_id,
    changes,
    metadata,
    ip_address,
    user_agent,
    success,
    error_message
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetActivityLog :one
-- Get a single activity log entry by ID
SELECT * FROM public.activity_log WHERE id = $1;

-- name: ListActivityLogs :many
-- List activity logs with pagination, ordered by created_at DESC
SELECT * FROM public.activity_log
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountActivityLogs :one
-- Count total activity logs
SELECT COUNT(*) FROM public.activity_log;

-- name: SearchActivityLogs :many
-- Search activity logs with optional filters
SELECT * FROM public.activity_log
WHERE
    (sqlc.narg('user_id')::UUID IS NULL OR user_id = sqlc.narg('user_id'))
    AND (sqlc.narg('action')::VARCHAR IS NULL OR action = sqlc.narg('action'))
    AND (sqlc.narg('resource_type')::VARCHAR IS NULL OR resource_type = sqlc.narg('resource_type'))
    AND (sqlc.narg('resource_id')::UUID IS NULL OR resource_id = sqlc.narg('resource_id'))
    AND (sqlc.narg('success')::BOOLEAN IS NULL OR success = sqlc.narg('success'))
    AND (sqlc.narg('start_time')::TIMESTAMPTZ IS NULL OR created_at >= sqlc.narg('start_time'))
    AND (sqlc.narg('end_time')::TIMESTAMPTZ IS NULL OR created_at <= sqlc.narg('end_time'))
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountSearchActivityLogs :one
-- Count activity logs matching search filters
SELECT COUNT(*) FROM public.activity_log
WHERE
    (sqlc.narg('user_id')::UUID IS NULL OR user_id = sqlc.narg('user_id'))
    AND (sqlc.narg('action')::VARCHAR IS NULL OR action = sqlc.narg('action'))
    AND (sqlc.narg('resource_type')::VARCHAR IS NULL OR resource_type = sqlc.narg('resource_type'))
    AND (sqlc.narg('resource_id')::UUID IS NULL OR resource_id = sqlc.narg('resource_id'))
    AND (sqlc.narg('success')::BOOLEAN IS NULL OR success = sqlc.narg('success'))
    AND (sqlc.narg('start_time')::TIMESTAMPTZ IS NULL OR created_at >= sqlc.narg('start_time'))
    AND (sqlc.narg('end_time')::TIMESTAMPTZ IS NULL OR created_at <= sqlc.narg('end_time'));

-- name: GetUserActivityLogs :many
-- Get activity logs for a specific user
SELECT * FROM public.activity_log
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountUserActivityLogs :one
-- Count activity logs for a specific user
SELECT COUNT(*) FROM public.activity_log WHERE user_id = $1;

-- name: GetResourceActivityLogs :many
-- Get activity logs for a specific resource
SELECT * FROM public.activity_log
WHERE resource_type = $1 AND resource_id = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountResourceActivityLogs :one
-- Count activity logs for a specific resource
SELECT COUNT(*) FROM public.activity_log
WHERE resource_type = $1 AND resource_id = $2;

-- name: GetFailedActivityLogs :many
-- Get failed activity logs (for monitoring)
SELECT * FROM public.activity_log
WHERE success = false
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetActivityLogsByAction :many
-- Get activity logs by action type
SELECT * FROM public.activity_log
WHERE action = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetActivityLogsByIP :many
-- Get activity logs from a specific IP address
SELECT * FROM public.activity_log
WHERE ip_address = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: DeleteOldActivityLogs :execrows
-- Delete activity logs older than a specific date (for cleanup job)
DELETE FROM public.activity_log
WHERE created_at < $1;

-- name: GetOldActivityLogsCount :one
-- Count activity logs older than a specific date (for dry-run)
SELECT COUNT(*) FROM public.activity_log
WHERE created_at < $1;

-- name: GetActivityLogStats :one
-- Get activity log statistics
SELECT
    COUNT(*) as total_count,
    COUNT(*) FILTER (WHERE success = true) as success_count,
    COUNT(*) FILTER (WHERE success = false) as failed_count,
    MIN(created_at) as oldest_entry,
    MAX(created_at) as newest_entry
FROM public.activity_log;

-- name: GetRecentActions :many
-- Get recent distinct actions (for autocomplete/filtering)
SELECT DISTINCT action, COUNT(*) as count
FROM public.activity_log
WHERE created_at > now() - interval '7 days'
GROUP BY action
ORDER BY count DESC
LIMIT $1;

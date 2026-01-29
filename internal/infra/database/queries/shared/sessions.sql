-- name: GetSessionByID :one
SELECT * FROM sessions WHERE id = $1;

-- name: GetSessionByTokenHash :one
SELECT * FROM sessions
WHERE token_hash = $1 AND is_active = true AND expires_at > NOW();

-- name: ListSessionsByUser :many
SELECT * FROM sessions
WHERE user_id = $1
ORDER BY last_activity DESC;

-- name: ListActiveSessions :many
SELECT * FROM sessions
WHERE is_active = true AND expires_at > NOW()
ORDER BY last_activity DESC
LIMIT $1 OFFSET $2;

-- name: CreateSession :one
INSERT INTO sessions (
    user_id, profile_id, token_hash,
    device_name, device_type, client_name, client_version,
    ip_address, user_agent, expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: UpdateSessionActivity :exec
UPDATE sessions SET
    last_activity = NOW(),
    ip_address = COALESCE(sqlc.narg('ip_address'), ip_address)
WHERE id = sqlc.arg('id');

-- name: DeactivateSession :exec
UPDATE sessions SET is_active = false WHERE id = $1;

-- name: DeactivateUserSessions :exec
UPDATE sessions SET is_active = false WHERE user_id = $1;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE id = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions WHERE expires_at < NOW();

-- name: CountActiveSessionsByUser :one
SELECT COUNT(*) FROM sessions
WHERE user_id = $1 AND is_active = true AND expires_at > NOW();

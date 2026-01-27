-- Session queries for Jellyfin Go

-- =============================================================================
-- BASIC CRUD
-- =============================================================================

-- name: GetSessionByID :one
SELECT * FROM sessions WHERE id = $1 LIMIT 1;

-- name: GetSessionByTokenHash :one
SELECT * FROM sessions WHERE token_hash = $1 LIMIT 1;

-- name: GetSessionByRefreshTokenHash :one
SELECT * FROM sessions WHERE refresh_token_hash = $1 LIMIT 1;

-- name: ListUserSessions :many
SELECT * FROM sessions
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: ListActiveSessions :many
SELECT * FROM sessions
WHERE expires_at > NOW()
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountUserSessions :one
SELECT count(*) FROM sessions WHERE user_id = $1;

-- name: CreateSession :one
INSERT INTO sessions (
    user_id, token_hash, refresh_token_hash, device_id, device_name,
    client_name, client_version, ip_address, expires_at, refresh_expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: UpdateSessionRefreshToken :exec
UPDATE sessions
SET
    refresh_token_hash = $2,
    refresh_expires_at = $3
WHERE id = $1;

-- name: DeleteSession :exec
DELETE FROM sessions WHERE id = $1;

-- name: DeleteSessionByTokenHash :exec
DELETE FROM sessions WHERE token_hash = $1;

-- name: DeleteUserSessions :exec
DELETE FROM sessions WHERE user_id = $1;

-- name: DeleteExpiredSessions :execrows
DELETE FROM sessions WHERE expires_at < NOW();

-- name: SessionExists :one
SELECT EXISTS(SELECT 1 FROM sessions WHERE token_hash = $1 AND expires_at > NOW());

-- =============================================================================
-- JOIN QUERIES
-- =============================================================================

-- name: GetSessionWithUser :one
SELECT
    s.*,
    u.id AS user_id,
    u.username,
    u.email,
    u.display_name,
    u.is_admin,
    u.is_disabled
FROM sessions s
JOIN users u ON s.user_id = u.id
WHERE s.token_hash = $1 AND s.expires_at > NOW()
LIMIT 1;

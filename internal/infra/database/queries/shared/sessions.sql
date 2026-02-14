-- Session Management Queries
-- Persistent session tracking with device information

-- name: CreateSession :one
INSERT INTO
    shared.sessions (
        user_id,
        token_hash,
        refresh_token_hash,
        ip_address,
        user_agent,
        device_name,
        scopes,
        expires_at,
        last_activity_at
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        NOW()
    ) RETURNING *;

-- name: GetSessionByTokenHash :one
SELECT *
FROM shared.sessions
WHERE
    token_hash = $1
    AND revoked_at IS NULL
    AND expires_at > NOW()
LIMIT 1;

-- name: GetSessionByID :one
SELECT *
FROM shared.sessions
WHERE
    id = $1
    AND revoked_at IS NULL
LIMIT 1;

-- name: GetSessionByRefreshTokenHash :one
SELECT *
FROM shared.sessions
WHERE
    refresh_token_hash = $1
    AND revoked_at IS NULL
    AND expires_at > NOW()
LIMIT 1;

-- name: ListUserSessions :many
SELECT *
FROM shared.sessions
WHERE
    user_id = $1
    AND revoked_at IS NULL
    AND expires_at > NOW()
ORDER BY last_activity_at DESC;

-- name: ListAllUserSessions :many
-- Includes expired but not revoked sessions (for user to see full history)
SELECT *
FROM shared.sessions
WHERE
    user_id = $1
    AND revoked_at IS NULL
ORDER BY last_activity_at DESC;

-- name: CountActiveUserSessions :one
SELECT COUNT(*)
FROM shared.sessions
WHERE
    user_id = $1
    AND revoked_at IS NULL
    AND expires_at > NOW();

-- name: CountAllActiveSessions :one
SELECT COUNT(*)
FROM shared.sessions
WHERE
    revoked_at IS NULL
    AND expires_at > NOW();

-- name: UpdateSessionActivity :exec
UPDATE shared.sessions
SET
    last_activity_at = NOW()
WHERE
    id = $1
    AND revoked_at IS NULL;

-- name: UpdateSessionActivityByTokenHash :exec
UPDATE shared.sessions
SET
    last_activity_at = NOW()
WHERE
    token_hash = $1
    AND revoked_at IS NULL;

-- name: RevokeSession :exec
UPDATE shared.sessions
SET
    revoked_at = NOW(),
    revoke_reason = sqlc.narg ('reason')
WHERE
    id = $1;

-- name: RevokeSessionByTokenHash :exec
UPDATE shared.sessions
SET
    revoked_at = NOW(),
    revoke_reason = sqlc.narg ('reason')
WHERE
    token_hash = $1;

-- name: RevokeAllUserSessions :exec
UPDATE shared.sessions
SET
    revoked_at = NOW(),
    revoke_reason = sqlc.narg ('reason')
WHERE
    user_id = $1
    AND revoked_at IS NULL;

-- name: RevokeAllUserSessionsExcept :exec
UPDATE shared.sessions
SET
    revoked_at = NOW(),
    revoke_reason = sqlc.narg ('reason')
WHERE
    user_id = $1
    AND id != $2
    AND revoked_at IS NULL;

-- name: DeleteExpiredSessions :execrows
DELETE FROM shared.sessions
WHERE
    expires_at < NOW() - INTERVAL '30 days';

-- name: DeleteRevokedSessions :execrows
DELETE FROM shared.sessions
WHERE
    revoked_at < NOW() - INTERVAL '30 days';

-- name: GetInactiveSessions :many
-- Sessions that haven't been active for N hours
SELECT * FROM shared.sessions
WHERE last_activity_at < sqlc.arg('inactive_since')::TIMESTAMPTZ
  AND revoked_at IS NULL
  AND expires_at > NOW();

-- name: RevokeInactiveSessions :exec
UPDATE shared.sessions
SET revoked_at = NOW(),
    revoke_reason = 'Inactivity timeout'
WHERE last_activity_at < sqlc.arg('inactive_since')::TIMESTAMPTZ
  AND revoked_at IS NULL;

-- MFA Session Tracking

-- name: MarkSessionMFAVerified :exec
UPDATE shared.sessions
SET
    mfa_verified = TRUE,
    mfa_verified_at = NOW()
WHERE
    id = $1
    AND revoked_at IS NULL;

-- name: MarkSessionMFAVerifiedByTokenHash :exec
UPDATE shared.sessions
SET
    mfa_verified = TRUE,
    mfa_verified_at = NOW()
WHERE
    token_hash = $1
    AND revoked_at IS NULL;

-- name: GetMFAVerifiedSessions :many
SELECT *
FROM shared.sessions
WHERE
    user_id = $1
    AND mfa_verified = TRUE
    AND revoked_at IS NULL
    AND expires_at > NOW()
ORDER BY mfa_verified_at DESC;

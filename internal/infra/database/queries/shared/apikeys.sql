-- name: CreateAPIKey :one
INSERT INTO shared.api_keys (
    user_id,
    name,
    description,
    key_hash,
    key_prefix,
    scopes,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetAPIKey :one
SELECT * FROM shared.api_keys
WHERE id = $1;

-- name: GetAPIKeyByHash :one
SELECT * FROM shared.api_keys
WHERE key_hash = $1;

-- name: GetAPIKeyByPrefix :one
SELECT * FROM shared.api_keys
WHERE key_prefix = $1 AND is_active = true
LIMIT 1;

-- name: ListUserAPIKeys :many
SELECT * FROM shared.api_keys
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: ListActiveUserAPIKeys :many
SELECT * FROM shared.api_keys
WHERE user_id = $1 AND is_active = true
ORDER BY created_at DESC;

-- name: CountUserAPIKeys :one
SELECT COUNT(*) FROM shared.api_keys
WHERE user_id = $1 AND is_active = true;

-- name: RevokeAPIKey :exec
UPDATE shared.api_keys
SET is_active = false, updated_at = NOW()
WHERE id = $1;

-- name: UpdateAPIKeyLastUsed :exec
UPDATE shared.api_keys
SET last_used_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: UpdateAPIKeyScopes :exec
UPDATE shared.api_keys
SET scopes = $2, updated_at = NOW()
WHERE id = $1;

-- name: DeleteAPIKey :exec
DELETE FROM shared.api_keys
WHERE id = $1;

-- name: DeleteExpiredAPIKeys :exec
DELETE FROM shared.api_keys
WHERE expires_at IS NOT NULL
  AND expires_at < NOW()
  AND is_active = false;

-- name: GetAPIKeyLastUsedAt :one
SELECT last_used_at FROM shared.api_keys
WHERE id = $1;

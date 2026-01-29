-- name: GetAPIKeyByID :one
SELECT * FROM api_keys WHERE id = $1;

-- name: GetAPIKeyByHash :one
SELECT * FROM api_keys
WHERE key_hash = $1 AND (expires_at IS NULL OR expires_at > NOW());

-- name: GetAPIKeyByPrefix :many
SELECT * FROM api_keys WHERE key_prefix = $1;

-- name: ListAPIKeysByUser :many
SELECT * FROM api_keys
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: CreateAPIKey :one
INSERT INTO api_keys (
    user_id, name, key_hash, key_prefix, scopes, expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: UpdateAPIKeyUsage :exec
UPDATE api_keys SET
    last_used_at = NOW(),
    use_count = use_count + 1
WHERE id = $1;

-- name: DeleteAPIKey :exec
DELETE FROM api_keys WHERE id = $1;

-- name: DeleteExpiredAPIKeys :exec
DELETE FROM api_keys WHERE expires_at IS NOT NULL AND expires_at < NOW();

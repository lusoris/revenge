-- name: GetResourceGrant :one
SELECT * FROM resource_grants
WHERE user_id = $1 AND resource_type = $2 AND resource_id = $3;

-- name: HasResourceGrant :one
SELECT EXISTS(
    SELECT 1 FROM resource_grants
    WHERE user_id = $1
    AND resource_type = $2
    AND resource_id = $3
    AND grant_type = ANY($4::text[])
    AND (expires_at IS NULL OR expires_at > NOW())
);

-- name: ListUserResourceGrants :many
SELECT * FROM resource_grants
WHERE user_id = $1
AND (expires_at IS NULL OR expires_at > NOW())
ORDER BY granted_at DESC
LIMIT $2 OFFSET $3;

-- name: ListUserResourceGrantsByType :many
SELECT * FROM resource_grants
WHERE user_id = $1
AND resource_type = $2
AND (expires_at IS NULL OR expires_at > NOW())
ORDER BY granted_at DESC
LIMIT $3 OFFSET $4;

-- name: ListResourceGrants :many
SELECT * FROM resource_grants
WHERE resource_type = $1 AND resource_id = $2
ORDER BY granted_at DESC;

-- name: CountResourceGrants :one
SELECT COUNT(*) FROM resource_grants
WHERE resource_type = $1 AND resource_id = $2;

-- name: CreateResourceGrant :one
INSERT INTO resource_grants (user_id, resource_type, resource_id, grant_type, granted_by, expires_at)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (user_id, resource_type, resource_id)
DO UPDATE SET grant_type = EXCLUDED.grant_type, granted_by = EXCLUDED.granted_by, granted_at = NOW(), expires_at = EXCLUDED.expires_at
RETURNING *;

-- name: UpdateResourceGrant :one
UPDATE resource_grants
SET grant_type = $4, expires_at = $5
WHERE user_id = $1 AND resource_type = $2 AND resource_id = $3
RETURNING *;

-- name: DeleteResourceGrant :exec
DELETE FROM resource_grants
WHERE user_id = $1 AND resource_type = $2 AND resource_id = $3;

-- name: DeleteResourceGrantByID :exec
DELETE FROM resource_grants WHERE id = $1;

-- name: DeleteResourceGrantsByResource :exec
DELETE FROM resource_grants
WHERE resource_type = $1 AND resource_id = $2;

-- name: DeleteExpiredResourceGrants :exec
DELETE FROM resource_grants
WHERE expires_at IS NOT NULL AND expires_at < NOW();

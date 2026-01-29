-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: CreateUser :one
INSERT INTO users (
    username, email, password_hash, is_admin,
    max_rating_level, adult_enabled,
    preferred_language, preferred_rating_system
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: UpdateUser :one
UPDATE users SET
    username = COALESCE(sqlc.narg('username'), username),
    email = COALESCE(sqlc.narg('email'), email),
    password_hash = COALESCE(sqlc.narg('password_hash'), password_hash),
    is_admin = COALESCE(sqlc.narg('is_admin'), is_admin),
    is_disabled = COALESCE(sqlc.narg('is_disabled'), is_disabled),
    max_rating_level = COALESCE(sqlc.narg('max_rating_level'), max_rating_level),
    adult_enabled = COALESCE(sqlc.narg('adult_enabled'), adult_enabled),
    preferred_language = COALESCE(sqlc.narg('preferred_language'), preferred_language),
    preferred_rating_system = COALESCE(sqlc.narg('preferred_rating_system'), preferred_rating_system)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateUserLastLogin :exec
UPDATE users SET last_login_at = NOW() WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: UserExistsByUsername :one
SELECT EXISTS(SELECT 1 FROM users WHERE username = $1);

-- name: UserExistsByEmail :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);

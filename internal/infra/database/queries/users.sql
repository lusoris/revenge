-- User queries for Jellyfin Go
-- sqlc generates type-safe Go code from these queries

-- =============================================================================
-- BASIC CRUD
-- =============================================================================

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY username ASC
LIMIT $1 OFFSET $2;

-- name: ListAdminUsers :many
SELECT * FROM users
WHERE is_admin = true
ORDER BY username ASC;

-- name: CountUsers :one
SELECT count(*) FROM users;

-- name: CountAdminUsers :one
SELECT count(*) FROM users WHERE is_admin = true;

-- name: CreateUser :one
INSERT INTO users (
    username, email, password_hash, display_name, is_admin
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
    email = COALESCE($2, email),
    display_name = COALESCE($3, display_name),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :exec
UPDATE users
SET password_hash = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserAdmin :exec
UPDATE users
SET is_admin = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserDisabled :exec
UPDATE users
SET is_disabled = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserLastLogin :exec
UPDATE users
SET last_login_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserLastActivity :exec
UPDATE users
SET last_activity_at = NOW()
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: UserExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE id = $1);

-- name: UsernameExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE username = $1);

-- name: EmailExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);

-- name: GetServerSetting :one
-- Get a server setting by key
SELECT * FROM shared.server_settings
WHERE key = $1
LIMIT 1;

-- name: ListServerSettings :many
-- Get all server settings
SELECT * FROM shared.server_settings
ORDER BY category, key;

-- name: ListServerSettingsByCategory :many
-- Get settings by category
SELECT * FROM shared.server_settings
WHERE category = $1
ORDER BY key;

-- name: ListPublicServerSettings :many
-- Get public settings (exposed in API)
SELECT * FROM shared.server_settings
WHERE is_public = TRUE
ORDER BY category, key;

-- name: CreateServerSetting :one
-- Create a new server setting
INSERT INTO shared.server_settings (
    key,
    value,
    description,
    category,
    data_type,
    is_secret,
    is_public,
    allowed_values,
    min_value,
    max_value,
    pattern,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
RETURNING *;

-- name: UpdateServerSetting :one
-- Update a server setting value
UPDATE shared.server_settings
SET
    value = $2,
    updated_at = NOW(),
    updated_by = $3
WHERE key = $1
RETURNING *;

-- name: UpsertServerSetting :one
-- Insert or update a server setting
INSERT INTO shared.server_settings (
    key,
    value,
    description,
    category,
    data_type,
    is_secret,
    is_public,
    updated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
ON CONFLICT (key) DO UPDATE
SET
    value = EXCLUDED.value,
    updated_at = NOW(),
    updated_by = EXCLUDED.updated_by
RETURNING *;

-- name: DeleteServerSetting :exec
-- Delete a server setting
DELETE FROM shared.server_settings
WHERE key = $1;

-- ============================================================================
-- User Settings Queries
-- ============================================================================

-- name: GetUserSetting :one
-- Get a user setting by user_id and key
SELECT * FROM shared.user_settings
WHERE user_id = $1 AND key = $2
LIMIT 1;

-- name: ListUserSettings :many
-- Get all settings for a user
SELECT * FROM shared.user_settings
WHERE user_id = $1
ORDER BY category, key;

-- name: ListUserSettingsByCategory :many
-- Get user settings by category
SELECT * FROM shared.user_settings
WHERE user_id = $1 AND category = $2
ORDER BY key;

-- name: CreateUserSetting :one
-- Create a new user setting
INSERT INTO shared.user_settings (
    user_id,
    key,
    value,
    description,
    category,
    data_type
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateUserSetting :one
-- Update a user setting value
UPDATE shared.user_settings
SET
    value = $2,
    updated_at = NOW()
WHERE user_id = $1 AND key = $3
RETURNING *;

-- name: UpsertUserSetting :one
-- Insert or update a user setting
INSERT INTO shared.user_settings (
    user_id,
    key,
    value,
    description,
    category,
    data_type
) VALUES (
    $1, $2, $3, $4, $5, $6
)
ON CONFLICT (user_id, key) DO UPDATE
SET
    value = EXCLUDED.value,
    updated_at = NOW()
RETURNING *;

-- name: DeleteUserSetting :exec
-- Delete a user setting
DELETE FROM shared.user_settings
WHERE user_id = $1 AND key = $2;

-- name: DeleteAllUserSettings :exec
-- Delete all settings for a user (used when user is deleted)
DELETE FROM shared.user_settings
WHERE user_id = $1;

-- name: GetUserByID :one
-- Get a user by their UUID
SELECT * FROM shared.users WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserByUsername :one
-- Get a user by username
SELECT *
FROM shared.users
WHERE
    username = $1
    AND deleted_at IS NULL;

-- name: GetUserByEmail :one
-- Get a user by email
SELECT * FROM shared.users WHERE email = $1 AND deleted_at IS NULL;

-- name: ListUsers :many
-- List all active users with optional filters
-- Use sqlc.narg for nullable parameters
SELECT * FROM shared.users
WHERE deleted_at IS NULL
  AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active'))
  AND (sqlc.narg('is_admin')::boolean IS NULL OR is_admin = sqlc.narg('is_admin'))
  AND (sqlc.narg('query')::text IS NULL OR (
    username ILIKE '%' || sqlc.narg('query') || '%'
    OR email ILIKE '%' || sqlc.narg('query') || '%'
    OR display_name ILIKE '%' || sqlc.narg('query') || '%'
  ))
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountUsers :one
-- Count users matching filters
SELECT COUNT(*) FROM shared.users
WHERE deleted_at IS NULL
  AND (sqlc.narg('is_active')::boolean IS NULL OR is_active = sqlc.narg('is_active'))
  AND (sqlc.narg('is_admin')::boolean IS NULL OR is_admin = sqlc.narg('is_admin'))
  AND (sqlc.narg('query')::text IS NULL OR (
    username ILIKE '%' || sqlc.narg('query') || '%'
    OR email ILIKE '%' || sqlc.narg('query') || '%'
    OR display_name ILIKE '%' || sqlc.narg('query') || '%'
  ));

-- name: CreateUser :one
-- Create a new user
INSERT INTO
    shared.users (
        username,
        email,
        password_hash,
        display_name,
        timezone,
        qar_enabled,
        is_active,
        is_admin
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8
    ) RETURNING *;

-- name: UpdateUser :one
-- Update user fields
UPDATE shared.users
SET
    email = COALESCE(sqlc.narg ('email'), email),
    display_name = COALESCE(
        sqlc.narg ('display_name'),
        display_name
    ),
    avatar_url = COALESCE(
        sqlc.narg ('avatar_url'),
        avatar_url
    ),
    timezone = COALESCE(
        sqlc.narg ('timezone'),
        timezone
    ),
    qar_enabled = COALESCE(
        sqlc.narg ('qar_enabled'),
        qar_enabled
    ),
    is_active = COALESCE(
        sqlc.narg ('is_active'),
        is_active
    ),
    is_admin = COALESCE(
        sqlc.narg ('is_admin'),
        is_admin
    ),
    updated_at = NOW()
WHERE
    id = sqlc.arg ('user_id')
    AND deleted_at IS NULL RETURNING *;

-- name: UpdatePassword :exec
-- Update user password hash
UPDATE shared.users
SET
    password_hash = $2,
    updated_at = NOW()
WHERE
    id = $1
    AND deleted_at IS NULL;

-- name: UpdateLastLogin :exec
-- Update user last login timestamp
UPDATE shared.users
SET
    last_login_at = NOW(),
    updated_at = NOW()
WHERE
    id = $1
    AND deleted_at IS NULL;

-- name: VerifyEmail :exec
-- Mark email as verified
UPDATE shared.users
SET
    email_verified = TRUE,
    email_verified_at = NOW(),
    updated_at = NOW()
WHERE
    id = $1
    AND deleted_at IS NULL;

-- name: DeleteUser :exec
-- Soft delete a user
UPDATE shared.users
SET
    deleted_at = NOW(),
    updated_at = NOW()
WHERE
    id = $1
    AND deleted_at IS NULL;

-- name: HardDeleteUser :exec
-- Permanently delete a user (GDPR compliance)
DELETE FROM shared.users WHERE id = $1;

-- ============================================================================
-- User Preferences Queries
-- ============================================================================

-- name: GetUserPreferences :one
-- Get user preferences
SELECT * FROM shared.user_preferences WHERE user_id = $1;

-- name: UpsertUserPreferences :one
-- Create or update user preferences
INSERT INTO
    shared.user_preferences (
        user_id,
        email_notifications,
        push_notifications,
        digest_notifications,
        profile_visibility,
        show_email,
        show_activity,
        theme,
        display_language,
        content_language,
        metadata_language,
        show_adult_content,
        show_spoilers,
        auto_play_videos
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
        $9,
        $10,
        $11,
        $12,
        $13,
        $14
    ) ON CONFLICT (user_id) DO
UPDATE
SET
    email_notifications = COALESCE(
        EXCLUDED.email_notifications,
        user_preferences.email_notifications
    ),
    push_notifications = COALESCE(
        EXCLUDED.push_notifications,
        user_preferences.push_notifications
    ),
    digest_notifications = COALESCE(
        EXCLUDED.digest_notifications,
        user_preferences.digest_notifications
    ),
    profile_visibility = COALESCE(
        EXCLUDED.profile_visibility,
        user_preferences.profile_visibility
    ),
    show_email = COALESCE(
        EXCLUDED.show_email,
        user_preferences.show_email
    ),
    show_activity = COALESCE(
        EXCLUDED.show_activity,
        user_preferences.show_activity
    ),
    theme = COALESCE(
        EXCLUDED.theme,
        user_preferences.theme
    ),
    display_language = COALESCE(
        EXCLUDED.display_language,
        user_preferences.display_language
    ),
    content_language = COALESCE(
        EXCLUDED.content_language,
        user_preferences.content_language
    ),
    metadata_language = COALESCE(
        EXCLUDED.metadata_language,
        user_preferences.metadata_language
    ),
    show_adult_content = COALESCE(
        EXCLUDED.show_adult_content,
        user_preferences.show_adult_content
    ),
    show_spoilers = COALESCE(
        EXCLUDED.show_spoilers,
        user_preferences.show_spoilers
    ),
    auto_play_videos = COALESCE(
        EXCLUDED.auto_play_videos,
        user_preferences.auto_play_videos
    ),
    updated_at = NOW() RETURNING *;

-- name: DeleteUserPreferences :exec
-- Delete user preferences (cleanup on user deletion)
DELETE FROM shared.user_preferences WHERE user_id = $1;

-- ============================================================================
-- User Avatars Queries
-- ============================================================================

-- name: GetCurrentAvatar :one
-- Get the current avatar for a user
SELECT *
FROM shared.user_avatars
WHERE
    user_id = $1
    AND is_current = TRUE
    AND deleted_at IS NULL;

-- name: GetAvatarByID :one
-- Get a specific avatar by ID
SELECT *
FROM shared.user_avatars
WHERE
    id = $1
    AND deleted_at IS NULL;

-- name: ListUserAvatars :many
-- List all avatars for a user (for history)
SELECT *
FROM shared.user_avatars
WHERE
    user_id = $1
    AND deleted_at IS NULL
ORDER BY version DESC
LIMIT $2
OFFSET
    $3;

-- name: CreateAvatar :one
-- Upload a new avatar (sets it as current)
INSERT INTO
    shared.user_avatars (
        user_id,
        file_path,
        file_size_bytes,
        mime_type,
        width,
        height,
        is_animated,
        version,
        is_current,
        uploaded_from_ip,
        uploaded_from_user_agent
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
        TRUE,
        $9,
        $10
    ) RETURNING *;

-- name: UnsetCurrentAvatars :exec
-- Mark all user's avatars as not current (before setting a new current)
UPDATE shared.user_avatars
SET
    is_current = FALSE,
    updated_at = NOW()
WHERE
    user_id = $1
    AND is_current = TRUE
    AND deleted_at IS NULL;

-- name: SetCurrentAvatar :exec
-- Set an existing avatar as current
UPDATE shared.user_avatars
SET
    is_current = TRUE,
    updated_at = NOW()
WHERE
    id = $1
    AND deleted_at IS NULL;

-- name: DeleteAvatar :exec
-- Soft delete an avatar
UPDATE shared.user_avatars
SET
    deleted_at = NOW(),
    updated_at = NOW()
WHERE
    id = $1
    AND deleted_at IS NULL;

-- name: HardDeleteAvatar :exec
-- Permanently delete an avatar file
DELETE FROM shared.user_avatars WHERE id = $1;

-- name: GetLatestAvatarVersion :one
-- Get the latest avatar version number for a user
SELECT COALESCE(MAX(version), 0)::integer as max_version
FROM shared.user_avatars
WHERE user_id = $1;

-- name: GetProfileByID :one
SELECT * FROM profiles WHERE id = $1;

-- name: ListProfilesByUser :many
SELECT * FROM profiles
WHERE user_id = $1
ORDER BY is_default DESC, created_at ASC;

-- name: GetDefaultProfile :one
SELECT * FROM profiles
WHERE user_id = $1 AND is_default = true;

-- name: CreateProfile :one
INSERT INTO profiles (
    user_id, name, avatar_url, is_default, is_kids,
    max_rating_level, adult_enabled,
    preferred_language, preferred_audio_language, preferred_subtitle_language,
    autoplay_next, autoplay_previews
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
) RETURNING *;

-- name: UpdateProfile :one
UPDATE profiles SET
    name = COALESCE(sqlc.narg('name'), name),
    avatar_url = COALESCE(sqlc.narg('avatar_url'), avatar_url),
    is_kids = COALESCE(sqlc.narg('is_kids'), is_kids),
    max_rating_level = COALESCE(sqlc.narg('max_rating_level'), max_rating_level),
    adult_enabled = COALESCE(sqlc.narg('adult_enabled'), adult_enabled),
    preferred_language = COALESCE(sqlc.narg('preferred_language'), preferred_language),
    preferred_audio_language = COALESCE(sqlc.narg('preferred_audio_language'), preferred_audio_language),
    preferred_subtitle_language = COALESCE(sqlc.narg('preferred_subtitle_language'), preferred_subtitle_language),
    autoplay_next = COALESCE(sqlc.narg('autoplay_next'), autoplay_next),
    autoplay_previews = COALESCE(sqlc.narg('autoplay_previews'), autoplay_previews)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: SetDefaultProfile :exec
UPDATE profiles SET is_default = (id = $2)
WHERE user_id = $1;

-- name: DeleteProfile :exec
DELETE FROM profiles WHERE id = $1 AND is_default = false;

-- name: CountProfilesByUser :one
SELECT COUNT(*) FROM profiles WHERE user_id = $1;

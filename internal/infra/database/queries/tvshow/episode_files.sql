-- name: GetEpisodeFile :one
SELECT * FROM tvshow.episode_files WHERE id = $1;

-- name: GetEpisodeFileByPath :one
SELECT * FROM tvshow.episode_files WHERE file_path = $1;

-- name: GetEpisodeFileBySonarrID :one
SELECT * FROM tvshow.episode_files WHERE sonarr_file_id = $1;

-- name: ListEpisodeFilesByEpisode :many
SELECT * FROM tvshow.episode_files
WHERE episode_id = $1
ORDER BY created_at ASC;

-- name: CreateEpisodeFile :one
INSERT INTO tvshow.episode_files (
    episode_id, file_path, file_name, file_size,
    container, resolution, quality_profile,
    video_codec, audio_codec, bitrate_kbps, duration_seconds,
    audio_languages, subtitle_languages,
    sonarr_file_id
) VALUES (
    $1, $2, $3, $4,
    $5, $6, $7,
    $8, $9, $10, $11,
    $12, $13,
    $14
)
RETURNING *;

-- name: UpdateEpisodeFile :one
UPDATE tvshow.episode_files SET
    file_path = COALESCE(sqlc.narg('file_path'), file_path),
    file_name = COALESCE(sqlc.narg('file_name'), file_name),
    file_size = COALESCE(sqlc.narg('file_size'), file_size),
    container = COALESCE(sqlc.narg('container'), container),
    resolution = COALESCE(sqlc.narg('resolution'), resolution),
    quality_profile = COALESCE(sqlc.narg('quality_profile'), quality_profile),
    video_codec = COALESCE(sqlc.narg('video_codec'), video_codec),
    audio_codec = COALESCE(sqlc.narg('audio_codec'), audio_codec),
    bitrate_kbps = COALESCE(sqlc.narg('bitrate_kbps'), bitrate_kbps),
    duration_seconds = COALESCE(sqlc.narg('duration_seconds'), duration_seconds),
    audio_languages = COALESCE(sqlc.narg('audio_languages'), audio_languages),
    subtitle_languages = COALESCE(sqlc.narg('subtitle_languages'), subtitle_languages),
    sonarr_file_id = COALESCE(sqlc.narg('sonarr_file_id'), sonarr_file_id)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteEpisodeFile :exec
DELETE FROM tvshow.episode_files WHERE id = $1;

-- name: DeleteEpisodeFilesByEpisode :exec
DELETE FROM tvshow.episode_files WHERE episode_id = $1;

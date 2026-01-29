-- TV Show Images and Videos Queries

-- Series Images

-- name: GetSeriesImages :many
SELECT * FROM series_images
WHERE series_id = $1
ORDER BY is_primary DESC, vote_average DESC NULLS LAST;

-- name: GetSeriesImagesByType :many
SELECT * FROM series_images
WHERE series_id = $1 AND image_type = $2
ORDER BY is_primary DESC, vote_average DESC NULLS LAST;

-- name: GetPrimarySeriesImage :one
SELECT * FROM series_images
WHERE series_id = $1 AND image_type = $2 AND is_primary = true
LIMIT 1;

-- name: CreateSeriesImage :one
INSERT INTO series_images (
    series_id, image_type, url, local_path,
    width, height, aspect_ratio,
    language, vote_average, vote_count, blurhash,
    provider, provider_id, is_primary
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING *;

-- name: SetPrimarySeriesImage :exec
UPDATE series_images SET is_primary = (id = $2)
WHERE series_id = $1 AND image_type = $3;

-- name: DeleteSeriesImages :exec
DELETE FROM series_images WHERE series_id = $1;

-- name: DeleteSeriesImagesByType :exec
DELETE FROM series_images WHERE series_id = $1 AND image_type = $2;

-- Season Images

-- name: GetSeasonImages :many
SELECT * FROM season_images
WHERE season_id = $1
ORDER BY is_primary DESC, vote_average DESC NULLS LAST;

-- name: GetSeasonImagesByType :many
SELECT * FROM season_images
WHERE season_id = $1 AND image_type = $2
ORDER BY is_primary DESC, vote_average DESC NULLS LAST;

-- name: GetPrimarySeasonImage :one
SELECT * FROM season_images
WHERE season_id = $1 AND image_type = $2 AND is_primary = true
LIMIT 1;

-- name: CreateSeasonImage :one
INSERT INTO season_images (
    season_id, image_type, url, local_path,
    width, height, aspect_ratio,
    language, vote_average, vote_count, blurhash,
    provider, provider_id, is_primary
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING *;

-- name: SetPrimarySeasonImage :exec
UPDATE season_images SET is_primary = (id = $2)
WHERE season_id = $1 AND image_type = $3;

-- name: DeleteSeasonImages :exec
DELETE FROM season_images WHERE season_id = $1;

-- name: DeleteSeasonImagesBySeries :exec
DELETE FROM season_images WHERE season_id IN (SELECT id FROM seasons WHERE series_id = $1);

-- Episode Images

-- name: GetEpisodeImages :many
SELECT * FROM episode_images
WHERE episode_id = $1
ORDER BY is_primary DESC, vote_average DESC NULLS LAST;

-- name: GetPrimaryEpisodeImage :one
SELECT * FROM episode_images
WHERE episode_id = $1 AND is_primary = true
LIMIT 1;

-- name: CreateEpisodeImage :one
INSERT INTO episode_images (
    episode_id, image_type, url, local_path,
    width, height, aspect_ratio,
    vote_average, vote_count, blurhash,
    provider, provider_id, is_primary
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING *;

-- name: SetPrimaryEpisodeImage :exec
UPDATE episode_images SET is_primary = (id = $2)
WHERE episode_id = $1;

-- name: DeleteEpisodeImages :exec
DELETE FROM episode_images WHERE episode_id = $1;

-- name: DeleteEpisodeImagesBySeries :exec
DELETE FROM episode_images WHERE episode_id IN (SELECT id FROM episodes WHERE series_id = $1);

-- Series Videos (trailers, etc.)

-- name: GetSeriesVideos :many
SELECT * FROM series_videos
WHERE series_id = $1
ORDER BY
    CASE video_type
        WHEN 'trailer' THEN 1
        WHEN 'teaser' THEN 2
        WHEN 'featurette' THEN 3
        WHEN 'behind_the_scenes' THEN 4
        ELSE 10
    END,
    size DESC;

-- name: GetSeriesVideosByType :many
SELECT * FROM series_videos
WHERE series_id = $1 AND video_type = $2
ORDER BY size DESC;

-- name: CreateSeriesVideo :one
INSERT INTO series_videos (
    series_id, video_type, name, key, site, size, language, is_official, tmdb_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: DeleteSeriesVideos :exec
DELETE FROM series_videos WHERE series_id = $1;

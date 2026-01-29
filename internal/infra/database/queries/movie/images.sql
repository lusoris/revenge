-- Movie Images and Videos Queries

-- Images

-- name: GetMovieImages :many
SELECT * FROM movie_images
WHERE movie_id = $1
ORDER BY is_primary DESC, vote_average DESC NULLS LAST;

-- name: GetMovieImagesByType :many
SELECT * FROM movie_images
WHERE movie_id = $1 AND image_type = $2
ORDER BY is_primary DESC, vote_average DESC NULLS LAST;

-- name: GetPrimaryMovieImage :one
SELECT * FROM movie_images
WHERE movie_id = $1 AND image_type = $2 AND is_primary = true
LIMIT 1;

-- name: CreateMovieImage :one
INSERT INTO movie_images (
    movie_id, image_type, path, blurhash, width, height, aspect_ratio, language, vote_average, vote_count, is_primary, source
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: SetPrimaryMovieImage :exec
UPDATE movie_images SET is_primary = (id = $2)
WHERE movie_id = $1 AND image_type = $3;

-- name: DeleteMovieImages :exec
DELETE FROM movie_images WHERE movie_id = $1;

-- name: DeleteMovieImagesByType :exec
DELETE FROM movie_images WHERE movie_id = $1 AND image_type = $2;

-- Videos

-- name: GetMovieVideos :many
SELECT * FROM movie_videos
WHERE movie_id = $1
ORDER BY
    CASE video_type
        WHEN 'trailer' THEN 1
        WHEN 'teaser' THEN 2
        WHEN 'clip' THEN 3
        ELSE 10
    END,
    size DESC;

-- name: GetMovieVideosByType :many
SELECT * FROM movie_videos
WHERE movie_id = $1 AND video_type = $2
ORDER BY size DESC;

-- name: CreateMovieVideo :one
INSERT INTO movie_videos (movie_id, video_type, site, key, name, language, size)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: DeleteMovieVideos :exec
DELETE FROM movie_videos WHERE movie_id = $1;

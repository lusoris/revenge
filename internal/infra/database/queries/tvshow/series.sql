-- TV Series Core Queries

-- name: GetSeriesByID :one
SELECT * FROM series WHERE id = $1;

-- name: GetSeriesByPath :one
SELECT s.* FROM series s
JOIN episodes e ON e.series_id = s.id
WHERE e.path = $1
LIMIT 1;

-- name: GetSeriesByTmdbID :one
SELECT * FROM series WHERE tmdb_id = $1;

-- name: GetSeriesByImdbID :one
SELECT * FROM series WHERE imdb_id = $1;

-- name: GetSeriesByTvdbID :one
SELECT * FROM series WHERE tvdb_id = $1;

-- name: ListSeries :many
SELECT * FROM series
ORDER BY
    CASE WHEN @sort_by::text = 'title' AND @sort_order::text = 'asc' THEN sort_title END ASC,
    CASE WHEN @sort_by::text = 'title' AND @sort_order::text = 'desc' THEN sort_title END DESC,
    CASE WHEN @sort_by::text = 'date_added' AND @sort_order::text = 'asc' THEN date_added END ASC,
    CASE WHEN @sort_by::text = 'date_added' AND @sort_order::text = 'desc' THEN date_added END DESC,
    CASE WHEN @sort_by::text = 'first_air_date' AND @sort_order::text = 'asc' THEN first_air_date END ASC,
    CASE WHEN @sort_by::text = 'first_air_date' AND @sort_order::text = 'desc' THEN first_air_date END DESC,
    CASE WHEN @sort_by::text = 'rating' AND @sort_order::text = 'asc' THEN community_rating END ASC NULLS LAST,
    CASE WHEN @sort_by::text = 'rating' AND @sort_order::text = 'desc' THEN community_rating END DESC NULLS LAST,
    sort_title ASC
LIMIT $1 OFFSET $2;

-- name: ListSeriesByLibrary :many
SELECT * FROM series
WHERE library_id = $1
ORDER BY
    CASE WHEN @sort_by::text = 'title' AND @sort_order::text = 'asc' THEN sort_title END ASC,
    CASE WHEN @sort_by::text = 'title' AND @sort_order::text = 'desc' THEN sort_title END DESC,
    CASE WHEN @sort_by::text = 'date_added' AND @sort_order::text = 'asc' THEN date_added END ASC,
    CASE WHEN @sort_by::text = 'date_added' AND @sort_order::text = 'desc' THEN date_added END DESC,
    CASE WHEN @sort_by::text = 'first_air_date' AND @sort_order::text = 'asc' THEN first_air_date END ASC,
    CASE WHEN @sort_by::text = 'first_air_date' AND @sort_order::text = 'desc' THEN first_air_date END DESC,
    CASE WHEN @sort_by::text = 'rating' AND @sort_order::text = 'asc' THEN community_rating END ASC NULLS LAST,
    CASE WHEN @sort_by::text = 'rating' AND @sort_order::text = 'desc' THEN community_rating END DESC NULLS LAST,
    sort_title ASC
LIMIT $2 OFFSET $3;

-- name: ListRecentlyAddedSeries :many
SELECT * FROM series
WHERE library_id = ANY(@library_ids::uuid[])
ORDER BY date_added DESC
LIMIT $1;

-- name: ListRecentlyPlayedSeries :many
SELECT * FROM series
WHERE library_id = ANY(@library_ids::uuid[])
  AND last_played_at IS NOT NULL
ORDER BY last_played_at DESC
LIMIT $1;

-- name: ListCurrentlyAiringSeries :many
SELECT * FROM series
WHERE library_id = ANY(@library_ids::uuid[])
  AND status IN ('Returning Series', 'In Production')
ORDER BY last_air_date DESC NULLS LAST
LIMIT $1;

-- name: SearchSeries :many
SELECT * FROM series
WHERE to_tsvector('english', COALESCE(title, '') || ' ' || COALESCE(original_title, '') || ' ' || COALESCE(overview, ''))
      @@ plainto_tsquery('english', $1)
ORDER BY ts_rank(
    to_tsvector('english', COALESCE(title, '') || ' ' || COALESCE(original_title, '') || ' ' || COALESCE(overview, '')),
    plainto_tsquery('english', $1)
) DESC
LIMIT $2 OFFSET $3;

-- name: CountSeries :one
SELECT COUNT(*) FROM series;

-- name: CountSeriesByLibrary :one
SELECT COUNT(*) FROM series WHERE library_id = $1;

-- name: CreateSeries :one
INSERT INTO series (
    library_id, title, sort_title, original_title, tagline, overview,
    first_air_date, last_air_date, year, status, type,
    content_rating, rating_level, community_rating, vote_count,
    poster_path, poster_blurhash, backdrop_path, backdrop_blurhash, logo_path,
    tmdb_id, imdb_id, tvdb_id,
    network_name, network_logo_path
) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8, $9, $10, $11,
    $12, $13, $14, $15,
    $16, $17, $18, $19, $20,
    $21, $22, $23,
    $24, $25
) RETURNING *;

-- name: UpdateSeries :one
UPDATE series SET
    title = COALESCE(sqlc.narg('title'), title),
    sort_title = COALESCE(sqlc.narg('sort_title'), sort_title),
    original_title = COALESCE(sqlc.narg('original_title'), original_title),
    tagline = COALESCE(sqlc.narg('tagline'), tagline),
    overview = COALESCE(sqlc.narg('overview'), overview),
    first_air_date = COALESCE(sqlc.narg('first_air_date'), first_air_date),
    last_air_date = COALESCE(sqlc.narg('last_air_date'), last_air_date),
    year = COALESCE(sqlc.narg('year'), year),
    status = COALESCE(sqlc.narg('status'), status),
    type = COALESCE(sqlc.narg('type'), type),
    content_rating = COALESCE(sqlc.narg('content_rating'), content_rating),
    rating_level = COALESCE(sqlc.narg('rating_level'), rating_level),
    community_rating = COALESCE(sqlc.narg('community_rating'), community_rating),
    vote_count = COALESCE(sqlc.narg('vote_count'), vote_count),
    poster_path = COALESCE(sqlc.narg('poster_path'), poster_path),
    poster_blurhash = COALESCE(sqlc.narg('poster_blurhash'), poster_blurhash),
    backdrop_path = COALESCE(sqlc.narg('backdrop_path'), backdrop_path),
    backdrop_blurhash = COALESCE(sqlc.narg('backdrop_blurhash'), backdrop_blurhash),
    logo_path = COALESCE(sqlc.narg('logo_path'), logo_path),
    tmdb_id = COALESCE(sqlc.narg('tmdb_id'), tmdb_id),
    imdb_id = COALESCE(sqlc.narg('imdb_id'), imdb_id),
    tvdb_id = COALESCE(sqlc.narg('tvdb_id'), tvdb_id),
    network_name = COALESCE(sqlc.narg('network_name'), network_name),
    network_logo_path = COALESCE(sqlc.narg('network_logo_path'), network_logo_path),
    is_locked = COALESCE(sqlc.narg('is_locked'), is_locked)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateSeriesPlaybackStats :exec
UPDATE series SET last_played_at = NOW() WHERE id = $1;

-- name: DeleteSeries :exec
DELETE FROM series WHERE id = $1;

-- name: DeleteSeriesByLibrary :exec
DELETE FROM series WHERE library_id = $1;

-- name: SeriesExistsByTmdbID :one
SELECT EXISTS(SELECT 1 FROM series WHERE tmdb_id = $1);

-- name: SeriesExistsByTvdbID :one
SELECT EXISTS(SELECT 1 FROM series WHERE tvdb_id = $1);

-- name: ListSeriesPaths :many
SELECT DISTINCT s.id, e.path FROM series s
JOIN episodes e ON e.series_id = s.id
WHERE s.library_id = $1;

-- Movie Core Queries

-- name: GetMovieByID :one
SELECT * FROM movies WHERE id = $1;

-- name: GetMovieByPath :one
SELECT * FROM movies WHERE path = $1;

-- name: GetMovieByTmdbID :one
SELECT * FROM movies WHERE tmdb_id = $1;

-- name: GetMovieByImdbID :one
SELECT * FROM movies WHERE imdb_id = $1;

-- name: ListMovies :many
SELECT * FROM movies
ORDER BY
    CASE WHEN @sort_by::text = 'title' AND @sort_order::text = 'asc' THEN sort_title END ASC,
    CASE WHEN @sort_by::text = 'title' AND @sort_order::text = 'desc' THEN sort_title END DESC,
    CASE WHEN @sort_by::text = 'date_added' AND @sort_order::text = 'asc' THEN date_added END ASC,
    CASE WHEN @sort_by::text = 'date_added' AND @sort_order::text = 'desc' THEN date_added END DESC,
    CASE WHEN @sort_by::text = 'release_date' AND @sort_order::text = 'asc' THEN release_date END ASC,
    CASE WHEN @sort_by::text = 'release_date' AND @sort_order::text = 'desc' THEN release_date END DESC,
    CASE WHEN @sort_by::text = 'rating' AND @sort_order::text = 'asc' THEN community_rating END ASC NULLS LAST,
    CASE WHEN @sort_by::text = 'rating' AND @sort_order::text = 'desc' THEN community_rating END DESC NULLS LAST,
    sort_title ASC
LIMIT $1 OFFSET $2;

-- name: ListMoviesByLibrary :many
SELECT * FROM movies
WHERE movie_library_id = $1
ORDER BY
    CASE WHEN @sort_by::text = 'title' AND @sort_order::text = 'asc' THEN sort_title END ASC,
    CASE WHEN @sort_by::text = 'title' AND @sort_order::text = 'desc' THEN sort_title END DESC,
    CASE WHEN @sort_by::text = 'date_added' AND @sort_order::text = 'asc' THEN date_added END ASC,
    CASE WHEN @sort_by::text = 'date_added' AND @sort_order::text = 'desc' THEN date_added END DESC,
    CASE WHEN @sort_by::text = 'release_date' AND @sort_order::text = 'asc' THEN release_date END ASC,
    CASE WHEN @sort_by::text = 'release_date' AND @sort_order::text = 'desc' THEN release_date END DESC,
    CASE WHEN @sort_by::text = 'rating' AND @sort_order::text = 'asc' THEN community_rating END ASC NULLS LAST,
    CASE WHEN @sort_by::text = 'rating' AND @sort_order::text = 'desc' THEN community_rating END DESC NULLS LAST,
    sort_title ASC
LIMIT $2 OFFSET $3;

-- name: ListMoviesByCollection :many
SELECT * FROM movies
WHERE collection_id = $1
ORDER BY collection_order ASC, release_date ASC;

-- name: ListRecentlyAddedMovies :many
SELECT * FROM movies
WHERE movie_library_id = ANY(@library_ids::uuid[])
ORDER BY date_added DESC
LIMIT $1;

-- name: ListRecentlyPlayedMovies :many
SELECT * FROM movies
WHERE movie_library_id = ANY(@library_ids::uuid[])
  AND last_played_at IS NOT NULL
ORDER BY last_played_at DESC
LIMIT $1;

-- name: SearchMovies :many
SELECT * FROM movies
WHERE to_tsvector('english', COALESCE(title, '') || ' ' || COALESCE(original_title, '') || ' ' || COALESCE(overview, ''))
      @@ plainto_tsquery('english', $1)
ORDER BY ts_rank(
    to_tsvector('english', COALESCE(title, '') || ' ' || COALESCE(original_title, '') || ' ' || COALESCE(overview, '')),
    plainto_tsquery('english', $1)
) DESC
LIMIT $2 OFFSET $3;

-- name: CountMovies :one
SELECT COUNT(*) FROM movies;

-- name: CountMoviesByLibrary :one
SELECT COUNT(*) FROM movies WHERE movie_library_id = $1;

-- name: CreateMovie :one
INSERT INTO movies (
    movie_library_id, path, container, size_bytes, runtime_ticks,
    title, sort_title, original_title, tagline, overview,
    release_date, year, content_rating, rating_level,
    budget, revenue, community_rating, vote_count,
    poster_path, poster_blurhash, backdrop_path, backdrop_blurhash, logo_path,
    tmdb_id, imdb_id, tvdb_id,
    collection_id, collection_order
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9, $10,
    $11, $12, $13, $14,
    $15, $16, $17, $18,
    $19, $20, $21, $22, $23,
    $24, $25, $26,
    $27, $28
) RETURNING *;

-- name: UpdateMovie :one
UPDATE movies SET
    title = COALESCE(sqlc.narg('title'), title),
    sort_title = COALESCE(sqlc.narg('sort_title'), sort_title),
    original_title = COALESCE(sqlc.narg('original_title'), original_title),
    tagline = COALESCE(sqlc.narg('tagline'), tagline),
    overview = COALESCE(sqlc.narg('overview'), overview),
    release_date = COALESCE(sqlc.narg('release_date'), release_date),
    year = COALESCE(sqlc.narg('year'), year),
    content_rating = COALESCE(sqlc.narg('content_rating'), content_rating),
    rating_level = COALESCE(sqlc.narg('rating_level'), rating_level),
    budget = COALESCE(sqlc.narg('budget'), budget),
    revenue = COALESCE(sqlc.narg('revenue'), revenue),
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
    collection_id = sqlc.narg('collection_id'),
    collection_order = COALESCE(sqlc.narg('collection_order'), collection_order),
    is_locked = COALESCE(sqlc.narg('is_locked'), is_locked)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateMoviePlaybackStats :exec
UPDATE movies SET
    last_played_at = NOW(),
    play_count = play_count + 1
WHERE id = $1;

-- name: DeleteMovie :exec
DELETE FROM movies WHERE id = $1;

-- name: DeleteMoviesByLibrary :exec
DELETE FROM movies WHERE movie_library_id = $1;

-- name: MovieExistsByPath :one
SELECT EXISTS(SELECT 1 FROM movies WHERE path = $1);

-- name: MovieExistsByTmdbID :one
SELECT EXISTS(SELECT 1 FROM movies WHERE tmdb_id = $1);

-- name: ListMoviePaths :many
SELECT id, path FROM movies WHERE movie_library_id = $1;

-- Collections

-- name: GetCollectionByID :one
SELECT * FROM movie_collections WHERE id = $1;

-- name: GetCollectionByTmdbID :one
SELECT * FROM movie_collections WHERE tmdb_id = $1;

-- name: ListCollections :many
SELECT * FROM movie_collections
ORDER BY name ASC
LIMIT $1 OFFSET $2;

-- name: CountCollections :one
SELECT COUNT(*) FROM movie_collections;

-- name: CreateCollection :one
INSERT INTO movie_collections (name, sort_name, overview, poster_path, poster_blurhash, backdrop_path, backdrop_blurhash, tmdb_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateCollection :one
UPDATE movie_collections SET
    name = COALESCE(sqlc.narg('name'), name),
    sort_name = COALESCE(sqlc.narg('sort_name'), sort_name),
    overview = COALESCE(sqlc.narg('overview'), overview),
    poster_path = COALESCE(sqlc.narg('poster_path'), poster_path),
    poster_blurhash = COALESCE(sqlc.narg('poster_blurhash'), poster_blurhash),
    backdrop_path = COALESCE(sqlc.narg('backdrop_path'), backdrop_path),
    backdrop_blurhash = COALESCE(sqlc.narg('backdrop_blurhash'), backdrop_blurhash)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteCollection :exec
DELETE FROM movie_collections WHERE id = $1;

-- Studios

-- name: GetStudioByID :one
SELECT * FROM movie_studios WHERE id = $1;

-- name: GetStudioByTmdbID :one
SELECT * FROM movie_studios WHERE tmdb_id = $1;

-- name: ListStudios :many
SELECT * FROM movie_studios ORDER BY name ASC LIMIT $1 OFFSET $2;

-- name: CreateStudio :one
INSERT INTO movie_studios (name, logo_path, tmdb_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: LinkMovieStudio :exec
INSERT INTO movie_studio_link (movie_id, studio_id, display_order)
VALUES ($1, $2, $3)
ON CONFLICT (movie_id, studio_id) DO UPDATE SET display_order = $3;

-- name: UnlinkMovieStudios :exec
DELETE FROM movie_studio_link WHERE movie_id = $1;

-- name: ListMovieStudios :many
SELECT s.* FROM movie_studios s
JOIN movie_studio_link msl ON s.id = msl.studio_id
WHERE msl.movie_id = $1
ORDER BY msl.display_order ASC;

-- Movie CRUD Operations
-- name: CreateMovie :one
INSERT INTO
    movie.movies (
        tmdb_id,
        imdb_id,
        title,
        original_title,
        year,
        release_date,
        runtime,
        overview,
        tagline,
        status,
        original_language,
        titles_i18n,
        taglines_i18n,
        overviews_i18n,
        age_ratings,
        external_ratings,
        poster_path,
        backdrop_path,
        trailer_url,
        vote_average,
        vote_count,
        popularity,
        budget,
        revenue,
        radarr_id,
        metadata_updated_at
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
        $14,
        $15,
        $16,
        $17,
        $18,
        $19,
        $20,
        $21,
        $22,
        $23,
        $24,
        $25,
        $26
    ) RETURNING *;

-- name: GetMovie :one
SELECT * FROM movie.movies WHERE id = $1 AND deleted_at IS NULL;

-- name: GetMovieByTMDbID :one
SELECT *
FROM movie.movies
WHERE
    tmdb_id = $1
    AND deleted_at IS NULL;

-- name: GetMovieByIMDbID :one
SELECT *
FROM movie.movies
WHERE
    imdb_id = $1
    AND deleted_at IS NULL;

-- name: GetMovieByRadarrID :one
SELECT *
FROM movie.movies
WHERE
    radarr_id = $1
    AND deleted_at IS NULL;

-- name: UpdateMovie :one
UPDATE movie.movies
SET
    tmdb_id = COALESCE(
        sqlc.narg ('tmdb_id'),
        tmdb_id
    ),
    imdb_id = COALESCE(
        sqlc.narg ('imdb_id'),
        imdb_id
    ),
    title = COALESCE(sqlc.narg ('title'), title),
    original_title = COALESCE(
        sqlc.narg ('original_title'),
        original_title
    ),
    year = COALESCE(sqlc.narg ('year'), year),
    release_date = COALESCE(
        sqlc.narg ('release_date'),
        release_date
    ),
    runtime = COALESCE(
        sqlc.narg ('runtime'),
        runtime
    ),
    overview = COALESCE(
        sqlc.narg ('overview'),
        overview
    ),
    tagline = COALESCE(
        sqlc.narg ('tagline'),
        tagline
    ),
    status = COALESCE(sqlc.narg ('status'), status),
    original_language = COALESCE(
        sqlc.narg ('original_language'),
        original_language
    ),
    titles_i18n = COALESCE(
        sqlc.narg ('titles_i18n'),
        titles_i18n
    ),
    taglines_i18n = COALESCE(
        sqlc.narg ('taglines_i18n'),
        taglines_i18n
    ),
    overviews_i18n = COALESCE(
        sqlc.narg ('overviews_i18n'),
        overviews_i18n
    ),
    age_ratings = COALESCE(
        sqlc.narg ('age_ratings'),
        age_ratings
    ),
    external_ratings = COALESCE(
        sqlc.narg ('external_ratings'),
        external_ratings
    ),
    poster_path = COALESCE(
        sqlc.narg ('poster_path'),
        poster_path
    ),
    backdrop_path = COALESCE(
        sqlc.narg ('backdrop_path'),
        backdrop_path
    ),
    trailer_url = COALESCE(
        sqlc.narg ('trailer_url'),
        trailer_url
    ),
    vote_average = COALESCE(
        sqlc.narg ('vote_average'),
        vote_average
    ),
    vote_count = COALESCE(
        sqlc.narg ('vote_count'),
        vote_count
    ),
    popularity = COALESCE(
        sqlc.narg ('popularity'),
        popularity
    ),
    budget = COALESCE(sqlc.narg ('budget'), budget),
    revenue = COALESCE(
        sqlc.narg ('revenue'),
        revenue
    ),
    radarr_id = COALESCE(
        sqlc.narg ('radarr_id'),
        radarr_id
    ),
    metadata_updated_at = COALESCE(
        sqlc.narg ('metadata_updated_at'),
        metadata_updated_at
    )
WHERE
    id = sqlc.arg ('id')
    AND deleted_at IS NULL RETURNING *;

-- name: DeleteMovie :exec
UPDATE movie.movies SET deleted_at = NOW() WHERE id = $1;

-- name: ListMovies :many
SELECT * FROM movie.movies
WHERE deleted_at IS NULL
ORDER BY
    CASE WHEN sqlc.narg('order_by')::text = 'title' THEN title END ASC,
    CASE WHEN sqlc.narg('order_by')::text = 'year' THEN year END DESC,
    CASE WHEN sqlc.narg('order_by')::text = 'added' THEN library_added_at END DESC,
    CASE WHEN sqlc.narg('order_by')::text = 'rating' THEN vote_average END DESC,
    library_added_at DESC
LIMIT $1 OFFSET $2;

-- name: CountMovies :one
SELECT COUNT(*) FROM movie.movies WHERE deleted_at IS NULL;

-- name: SearchMoviesByTitle :many
SELECT *
FROM movie.movies
WHERE
    deleted_at IS NULL
    AND (
        title % $1
        OR original_title % $1
    )
ORDER BY similarity (title, $1) DESC, similarity (original_title, $1) DESC
LIMIT $2
OFFSET
    $3;

-- name: SearchMoviesByTitleAnyLanguage :many
SELECT *
FROM movie.movies
WHERE
    deleted_at IS NULL
    AND (
        title ILIKE '%' || $1 || '%'
        OR original_title ILIKE '%' || $1 || '%'
        OR EXISTS (
            SELECT 1
            FROM jsonb_each_text (titles_i18n)
            WHERE
                value ILIKE '%' || $1 || '%'
        )
    )
ORDER BY
    CASE
        WHEN title ILIKE $1 THEN 1
        WHEN original_title ILIKE $1 THEN 2
        WHEN title ILIKE $1 || '%' THEN 3
        ELSE 4
    END,
    vote_average DESC NULLS LAST
LIMIT $2
OFFSET
    $3;

-- name: ListMoviesByYear :many
SELECT *
FROM movie.movies
WHERE
    deleted_at IS NULL
    AND year = $1
ORDER BY vote_average DESC NULLS LAST, title ASC
LIMIT $2
OFFSET
    $3;

-- name: ListRecentlyAdded :many
SELECT *
FROM movie.movies
WHERE
    deleted_at IS NULL
ORDER BY library_added_at DESC
LIMIT $1
OFFSET
    $2;

-- name: ListTopRated :many
SELECT *
FROM movie.movies
WHERE
    deleted_at IS NULL
    AND vote_average IS NOT NULL
    AND vote_count > $1
ORDER BY vote_average DESC, vote_count DESC
LIMIT $2
OFFSET
    $3;

-- name: CountTopRated :one
SELECT COUNT(*)
FROM movie.movies
WHERE
    deleted_at IS NULL
    AND vote_average IS NOT NULL
    AND vote_count > $1;

-- Movie Files Operations
-- name: CreateMovieFile :one
INSERT INTO
    movie.movie_files (
        movie_id,
        file_path,
        file_size,
        resolution,
        quality_profile,
        video_codec,
        audio_codec,
        container,
        bitrate_kbps,
        audio_languages,
        subtitle_languages,
        radarr_file_id
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
        $12
    ) RETURNING *;

-- name: GetMovieFile :one
SELECT *
FROM movie.movie_files
WHERE
    id = $1
    AND deleted_at IS NULL;

-- name: ListMovieFilesByMovieID :many
SELECT *
FROM movie.movie_files
WHERE
    movie_id = $1
    AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: GetMovieFileByPath :one
SELECT *
FROM movie.movie_files
WHERE
    file_path = $1
    AND deleted_at IS NULL;

-- name: GetMovieFileByRadarrID :one
SELECT *
FROM movie.movie_files
WHERE
    radarr_file_id = $1
    AND deleted_at IS NULL;

-- name: UpdateMovieFile :one
UPDATE movie.movie_files
SET
    file_path = COALESCE(
        sqlc.narg ('file_path'),
        file_path
    ),
    file_size = COALESCE(
        sqlc.narg ('file_size'),
        file_size
    ),
    resolution = COALESCE(
        sqlc.narg ('resolution'),
        resolution
    ),
    quality_profile = COALESCE(
        sqlc.narg ('quality_profile'),
        quality_profile
    ),
    video_codec = COALESCE(
        sqlc.narg ('video_codec'),
        video_codec
    ),
    audio_codec = COALESCE(
        sqlc.narg ('audio_codec'),
        audio_codec
    ),
    container = COALESCE(
        sqlc.narg ('container'),
        container
    ),
    bitrate_kbps = COALESCE(
        sqlc.narg ('bitrate_kbps'),
        bitrate_kbps
    ),
    audio_languages = COALESCE(
        sqlc.narg ('audio_languages'),
        audio_languages
    ),
    subtitle_languages = COALESCE(
        sqlc.narg ('subtitle_languages'),
        subtitle_languages
    ),
    radarr_file_id = COALESCE(
        sqlc.narg ('radarr_file_id'),
        radarr_file_id
    )
WHERE
    id = sqlc.arg ('id')
    AND deleted_at IS NULL RETURNING *;

-- name: DeleteMovieFile :exec
UPDATE movie.movie_files SET deleted_at = NOW() WHERE id = $1;

-- Movie Credits Operations
-- name: CreateMovieCredit :one
INSERT INTO
    movie.movie_credits (
        movie_id,
        tmdb_person_id,
        name,
        credit_type,
        character,
        job,
        department,
        cast_order,
        profile_path
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
        $9
    ) RETURNING *;

-- name: ListMovieCast :many
SELECT *
FROM movie.movie_credits
WHERE
    movie_id = $1
    AND credit_type = 'cast'
    AND deleted_at IS NULL
ORDER BY cast_order ASC NULLS LAST
LIMIT $2
OFFSET
    $3;

-- name: CountMovieCast :one
SELECT COUNT(*)
FROM movie.movie_credits
WHERE
    movie_id = $1
    AND credit_type = 'cast'
    AND deleted_at IS NULL;

-- name: ListMovieCrew :many
SELECT *
FROM movie.movie_credits
WHERE
    movie_id = $1
    AND credit_type = 'crew'
    AND deleted_at IS NULL
ORDER BY
    CASE department
        WHEN 'Directing' THEN 1
        WHEN 'Writing' THEN 2
        WHEN 'Production' THEN 3
        ELSE 99
    END,
    name ASC
LIMIT $2
OFFSET
    $3;

-- name: CountMovieCrew :one
SELECT COUNT(*)
FROM movie.movie_credits
WHERE
    movie_id = $1
    AND credit_type = 'crew'
    AND deleted_at IS NULL;

-- name: DeleteMovieCredits :exec
UPDATE movie.movie_credits
SET
    deleted_at = NOW()
WHERE
    movie_id = $1;

-- Movie Collections Operations
-- name: CreateMovieCollection :one
INSERT INTO
    movie.movie_collections (
        tmdb_collection_id,
        name,
        overview,
        poster_path,
        backdrop_path
    )
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetMovieCollection :one
SELECT *
FROM movie.movie_collections
WHERE
    id = $1
    AND deleted_at IS NULL;

-- name: GetMovieCollectionByTMDbID :one
SELECT *
FROM movie.movie_collections
WHERE
    tmdb_collection_id = $1
    AND deleted_at IS NULL;

-- name: UpdateMovieCollection :one
UPDATE movie.movie_collections
SET
    tmdb_collection_id = COALESCE(
        sqlc.narg ('tmdb_collection_id'),
        tmdb_collection_id
    ),
    name = COALESCE(sqlc.narg ('name'), name),
    overview = COALESCE(
        sqlc.narg ('overview'),
        overview
    ),
    poster_path = COALESCE(
        sqlc.narg ('poster_path'),
        poster_path
    ),
    backdrop_path = COALESCE(
        sqlc.narg ('backdrop_path'),
        backdrop_path
    )
WHERE
    id = sqlc.arg ('id')
    AND deleted_at IS NULL RETURNING *;

-- name: AddMovieToCollection :exec
INSERT INTO
    movie.movie_collection_members (
        collection_id,
        movie_id,
        collection_order
    )
VALUES ($1, $2, $3) ON CONFLICT (collection_id, movie_id) DO
UPDATE
SET
    collection_order = EXCLUDED.collection_order;

-- name: RemoveMovieFromCollection :exec
DELETE FROM movie.movie_collection_members
WHERE
    collection_id = $1
    AND movie_id = $2;

-- name: ListMoviesByCollection :many
SELECT m.*
FROM movie.movies m
    JOIN movie.movie_collection_members mcm ON m.id = mcm.movie_id
WHERE
    mcm.collection_id = $1
    AND m.deleted_at IS NULL
ORDER BY mcm.collection_order ASC NULLS LAST, m.year ASC;

-- name: GetCollectionForMovie :one
SELECT c.*
FROM movie.movie_collections c
    JOIN movie.movie_collection_members mcm ON c.id = mcm.collection_id
WHERE
    mcm.movie_id = $1
    AND c.deleted_at IS NULL
LIMIT 1;

-- Movie Genres Operations
-- name: AddMovieGenre :exec
INSERT INTO
    movie.movie_genres (movie_id, tmdb_genre_id, name)
VALUES ($1, $2, $3) ON CONFLICT (movie_id, tmdb_genre_id) DO NOTHING;

-- name: ListMovieGenres :many
SELECT *
FROM movie.movie_genres
WHERE
    movie_id = $1
ORDER BY name ASC;

-- name: ListDistinctMovieGenres :many
SELECT tmdb_genre_id, name, COUNT(DISTINCT movie_id)::bigint AS item_count
FROM movie.movie_genres
GROUP BY tmdb_genre_id, name
ORDER BY name ASC;

-- name: DeleteMovieGenres :exec
DELETE FROM movie.movie_genres WHERE movie_id = $1;

-- name: ListMoviesByGenre :many
SELECT m.*
FROM movie.movies m
    JOIN movie.movie_genres mg ON m.id = mg.movie_id
WHERE
    mg.tmdb_genre_id = $1
    AND m.deleted_at IS NULL
ORDER BY m.vote_average DESC NULLS LAST, m.title ASC
LIMIT $2
OFFSET
    $3;

-- Movie Watch Progress Operations
-- name: CreateOrUpdateWatchProgress :one
INSERT INTO
    movie.movie_watched (
        user_id,
        movie_id,
        progress_seconds,
        duration_seconds,
        is_completed
    )
VALUES ($1, $2, $3, $4, $5) ON CONFLICT (user_id, movie_id) DO
UPDATE
SET
    progress_seconds = EXCLUDED.progress_seconds,
    duration_seconds = EXCLUDED.duration_seconds,
    is_completed = EXCLUDED.is_completed,
    watch_count = CASE
        WHEN EXCLUDED.is_completed
        AND NOT movie_watched.is_completed THEN movie_watched.watch_count + 1
        ELSE movie_watched.watch_count
    END,
    last_watched_at = NOW() RETURNING *;

-- name: GetWatchProgress :one
SELECT *
FROM movie.movie_watched
WHERE
    user_id = $1
    AND movie_id = $2;

-- name: DeleteWatchProgress :exec
DELETE FROM movie.movie_watched
WHERE
    user_id = $1
    AND movie_id = $2;

-- name: ListContinueWatching :many
SELECT m.*, mw.progress_seconds, mw.duration_seconds, mw.progress_percent, mw.last_watched_at
FROM movie.movies m
    JOIN movie.movie_watched mw ON m.id = mw.movie_id
WHERE
    mw.user_id = $1
    AND mw.is_completed = FALSE
    AND mw.progress_percent > 5
    AND m.deleted_at IS NULL
ORDER BY mw.last_watched_at DESC
LIMIT $2;

-- name: ListWatchedMovies :many
SELECT m.*, mw.watch_count, mw.last_watched_at
FROM movie.movies m
    JOIN movie.movie_watched mw ON m.id = mw.movie_id
WHERE
    mw.user_id = $1
    AND mw.is_completed = TRUE
    AND m.deleted_at IS NULL
ORDER BY mw.last_watched_at DESC
LIMIT $2
OFFSET
    $3;

-- name: GetUserMovieStats :one
SELECT
    COUNT(*) FILTER (WHERE is_completed) as watched_count,
    COUNT(*) FILTER (WHERE NOT is_completed AND progress_percent > 5) as in_progress_count,
    COALESCE(SUM(watch_count), 0)::bigint as total_watches
FROM movie.movie_watched
WHERE user_id = $1;

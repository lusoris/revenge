-- name: GetSeries :one
SELECT * FROM tvshow.series WHERE id = $1;

-- name: GetSeriesByTMDbID :one
SELECT * FROM tvshow.series WHERE tmdb_id = $1;

-- name: GetSeriesByTVDbID :one
SELECT * FROM tvshow.series WHERE tvdb_id = $1;

-- name: GetSeriesBySonarrID :one
SELECT * FROM tvshow.series WHERE sonarr_id = $1;

-- name: ListSeries :many
SELECT * FROM tvshow.series
ORDER BY
    CASE WHEN @sort_by::text = 'title' AND @sort_order::text = 'asc' THEN title END ASC,
    CASE WHEN @sort_by::text = 'title' AND @sort_order::text = 'desc' THEN title END DESC,
    CASE WHEN @sort_by::text = 'first_air_date' AND @sort_order::text = 'asc' THEN first_air_date END ASC,
    CASE WHEN @sort_by::text = 'first_air_date' AND @sort_order::text = 'desc' THEN first_air_date END DESC,
    CASE WHEN @sort_by::text = 'created_at' AND @sort_order::text = 'asc' THEN created_at END ASC,
    CASE WHEN @sort_by::text = 'created_at' AND @sort_order::text = 'desc' THEN created_at END DESC,
    created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountSeries :one
SELECT COUNT(*) FROM tvshow.series;

-- name: SearchSeriesByTitle :many
SELECT *
FROM tvshow.series
WHERE
    title ILIKE '%' || $1 || '%'
    OR original_title ILIKE '%' || $1 || '%'
ORDER BY
    CASE
        WHEN LOWER(title) = LOWER($1) THEN 0
        WHEN LOWER(title) LIKE LOWER($1) || '%' THEN 1
        ELSE 2
    END,
    popularity DESC NULLS LAST
LIMIT $2
OFFSET
    $3;

-- name: SearchSeriesByTitleAnyLanguage :many
SELECT * FROM tvshow.series
WHERE title ILIKE '%' || $1 || '%'
   OR original_title ILIKE '%' || $1 || '%'
   OR titles_i18n::text ILIKE '%' || $1 || '%'
ORDER BY
    CASE
        WHEN LOWER(title) = LOWER($1) THEN 0
        WHEN LOWER(title) LIKE LOWER($1) || '%' THEN 1
        ELSE 2
    END,
    popularity DESC NULLS LAST
LIMIT $2 OFFSET $3;

-- name: ListRecentlyAddedSeries :many
SELECT *
FROM tvshow.series
ORDER BY created_at DESC
LIMIT $1
OFFSET
    $2;

-- name: ListSeriesByStatus :many
SELECT *
FROM tvshow.series
WHERE
    status = $1
ORDER BY first_air_date DESC
LIMIT $2
OFFSET
    $3;

-- name: CreateSeries :one
INSERT INTO
    tvshow.series (
        tmdb_id,
        tvdb_id,
        imdb_id,
        sonarr_id,
        title,
        tagline,
        overview,
        titles_i18n,
        taglines_i18n,
        overviews_i18n,
        age_ratings,
        external_ratings,
        original_language,
        original_title,
        status,
        type,
        first_air_date,
        last_air_date,
        vote_average,
        vote_count,
        popularity,
        poster_path,
        backdrop_path,
        total_seasons,
        total_episodes,
        trailer_url,
        homepage,
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
        $26,
        $27,
        $28
    ) RETURNING *;

-- name: UpdateSeries :one
UPDATE tvshow.series
SET
    tmdb_id = COALESCE(
        sqlc.narg ('tmdb_id'),
        tmdb_id
    ),
    tvdb_id = COALESCE(
        sqlc.narg ('tvdb_id'),
        tvdb_id
    ),
    imdb_id = COALESCE(
        sqlc.narg ('imdb_id'),
        imdb_id
    ),
    sonarr_id = COALESCE(
        sqlc.narg ('sonarr_id'),
        sonarr_id
    ),
    title = COALESCE(sqlc.narg ('title'), title),
    tagline = COALESCE(
        sqlc.narg ('tagline'),
        tagline
    ),
    overview = COALESCE(
        sqlc.narg ('overview'),
        overview
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
    original_language = COALESCE(
        sqlc.narg ('original_language'),
        original_language
    ),
    original_title = COALESCE(
        sqlc.narg ('original_title'),
        original_title
    ),
    status = COALESCE(sqlc.narg ('status'), status),
    type = COALESCE(sqlc.narg ('type'), type),
    first_air_date = COALESCE(
        sqlc.narg ('first_air_date'),
        first_air_date
    ),
    last_air_date = COALESCE(
        sqlc.narg ('last_air_date'),
        last_air_date
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
    poster_path = COALESCE(
        sqlc.narg ('poster_path'),
        poster_path
    ),
    backdrop_path = COALESCE(
        sqlc.narg ('backdrop_path'),
        backdrop_path
    ),
    total_seasons = COALESCE(
        sqlc.narg ('total_seasons'),
        total_seasons
    ),
    total_episodes = COALESCE(
        sqlc.narg ('total_episodes'),
        total_episodes
    ),
    trailer_url = COALESCE(
        sqlc.narg ('trailer_url'),
        trailer_url
    ),
    homepage = COALESCE(
        sqlc.narg ('homepage'),
        homepage
    ),
    metadata_updated_at = COALESCE(
        sqlc.narg ('metadata_updated_at'),
        metadata_updated_at
    )
WHERE
    id = sqlc.arg ('id') RETURNING *;

-- name: DeleteSeries :exec
DELETE FROM tvshow.series WHERE id = $1;

-- name: UpdateSeriesStats :exec
UPDATE tvshow.series
SET
    total_seasons = (
        SELECT COUNT(*)
        FROM tvshow.seasons se
        WHERE
            se.series_id = $1
    ),
    total_episodes = (
        SELECT COUNT(*)
        FROM tvshow.episodes ep
        WHERE
            ep.series_id = $1
    )
WHERE
    id = $1;

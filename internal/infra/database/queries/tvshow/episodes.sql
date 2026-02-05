-- name: GetEpisode :one
SELECT * FROM tvshow.episodes WHERE id = $1;

-- name: GetEpisodeByTMDbID :one
SELECT * FROM tvshow.episodes WHERE tmdb_id = $1;

-- name: GetEpisodeByNumber :one
SELECT * FROM tvshow.episodes
WHERE series_id = $1 AND season_number = $2 AND episode_number = $3;

-- name: ListEpisodesBySeries :many
SELECT * FROM tvshow.episodes
WHERE series_id = $1
ORDER BY season_number ASC, episode_number ASC;

-- name: ListEpisodesBySeason :many
SELECT * FROM tvshow.episodes
WHERE season_id = $1
ORDER BY episode_number ASC;

-- name: ListEpisodesBySeasonNumber :many
SELECT * FROM tvshow.episodes
WHERE series_id = $1 AND season_number = $2
ORDER BY episode_number ASC;

-- name: ListRecentEpisodes :many
SELECT e.*, s.title as series_title, s.poster_path as series_poster_path
FROM tvshow.episodes e
JOIN tvshow.series s ON e.series_id = s.id
WHERE e.air_date IS NOT NULL AND e.air_date <= CURRENT_DATE
ORDER BY e.air_date DESC
LIMIT $1 OFFSET $2;

-- name: ListUpcomingEpisodes :many
SELECT e.*, s.title as series_title, s.poster_path as series_poster_path
FROM tvshow.episodes e
JOIN tvshow.series s ON e.series_id = s.id
WHERE e.air_date IS NOT NULL AND e.air_date > CURRENT_DATE
ORDER BY e.air_date ASC
LIMIT $1 OFFSET $2;

-- name: CountEpisodesBySeries :one
SELECT COUNT(*) FROM tvshow.episodes WHERE series_id = $1;

-- name: CountEpisodesBySeason :one
SELECT COUNT(*) FROM tvshow.episodes WHERE season_id = $1;

-- name: CreateEpisode :one
INSERT INTO tvshow.episodes (
    series_id, season_id, tmdb_id, tvdb_id, imdb_id,
    season_number, episode_number,
    title, overview,
    titles_i18n, overviews_i18n,
    air_date, runtime,
    vote_average, vote_count,
    still_path, production_code
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7,
    $8, $9,
    $10, $11,
    $12, $13,
    $14, $15,
    $16, $17
)
RETURNING *;

-- name: UpdateEpisode :one
UPDATE tvshow.episodes SET
    tmdb_id = COALESCE(sqlc.narg('tmdb_id'), tmdb_id),
    tvdb_id = COALESCE(sqlc.narg('tvdb_id'), tvdb_id),
    imdb_id = COALESCE(sqlc.narg('imdb_id'), imdb_id),
    title = COALESCE(sqlc.narg('title'), title),
    overview = COALESCE(sqlc.narg('overview'), overview),
    titles_i18n = COALESCE(sqlc.narg('titles_i18n'), titles_i18n),
    overviews_i18n = COALESCE(sqlc.narg('overviews_i18n'), overviews_i18n),
    air_date = COALESCE(sqlc.narg('air_date'), air_date),
    runtime = COALESCE(sqlc.narg('runtime'), runtime),
    vote_average = COALESCE(sqlc.narg('vote_average'), vote_average),
    vote_count = COALESCE(sqlc.narg('vote_count'), vote_count),
    still_path = COALESCE(sqlc.narg('still_path'), still_path),
    production_code = COALESCE(sqlc.narg('production_code'), production_code)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteEpisode :exec
DELETE FROM tvshow.episodes WHERE id = $1;

-- name: DeleteEpisodesBySeries :exec
DELETE FROM tvshow.episodes WHERE series_id = $1;

-- name: DeleteEpisodesBySeason :exec
DELETE FROM tvshow.episodes WHERE season_id = $1;

-- name: UpsertEpisode :one
INSERT INTO tvshow.episodes (
    series_id, season_id, tmdb_id, tvdb_id, imdb_id,
    season_number, episode_number,
    title, overview,
    titles_i18n, overviews_i18n,
    air_date, runtime,
    vote_average, vote_count,
    still_path, production_code
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7,
    $8, $9,
    $10, $11,
    $12, $13,
    $14, $15,
    $16, $17
)
ON CONFLICT (series_id, season_number, episode_number) DO UPDATE SET
    tmdb_id = EXCLUDED.tmdb_id,
    tvdb_id = EXCLUDED.tvdb_id,
    imdb_id = EXCLUDED.imdb_id,
    title = EXCLUDED.title,
    overview = EXCLUDED.overview,
    titles_i18n = EXCLUDED.titles_i18n,
    overviews_i18n = EXCLUDED.overviews_i18n,
    air_date = EXCLUDED.air_date,
    runtime = EXCLUDED.runtime,
    vote_average = EXCLUDED.vote_average,
    vote_count = EXCLUDED.vote_count,
    still_path = EXCLUDED.still_path,
    production_code = EXCLUDED.production_code
RETURNING *;

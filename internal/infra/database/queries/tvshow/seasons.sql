-- name: GetSeason :one
SELECT * FROM tvshow.seasons WHERE id = $1;

-- name: GetSeasonByNumber :one
SELECT * FROM tvshow.seasons
WHERE series_id = $1 AND season_number = $2;

-- name: ListSeasonsBySeries :many
SELECT * FROM tvshow.seasons
WHERE series_id = $1
ORDER BY season_number ASC;

-- name: ListSeasonsBySeriesWithEpisodeCount :many
SELECT
    s.*,
    (SELECT COUNT(*) FROM tvshow.episodes e WHERE e.season_id = s.id) as actual_episode_count
FROM tvshow.seasons s
WHERE s.series_id = $1
ORDER BY s.season_number ASC;

-- name: CreateSeason :one
INSERT INTO tvshow.seasons (
    series_id, tmdb_id, season_number,
    name, overview,
    names_i18n, overviews_i18n,
    poster_path, episode_count, air_date, vote_average
) VALUES (
    $1, $2, $3,
    $4, $5,
    $6, $7,
    $8, $9, $10, $11
)
RETURNING *;

-- name: UpdateSeason :one
UPDATE tvshow.seasons SET
    tmdb_id = COALESCE(sqlc.narg('tmdb_id'), tmdb_id),
    name = COALESCE(sqlc.narg('name'), name),
    overview = COALESCE(sqlc.narg('overview'), overview),
    names_i18n = COALESCE(sqlc.narg('names_i18n'), names_i18n),
    overviews_i18n = COALESCE(sqlc.narg('overviews_i18n'), overviews_i18n),
    poster_path = COALESCE(sqlc.narg('poster_path'), poster_path),
    episode_count = COALESCE(sqlc.narg('episode_count'), episode_count),
    air_date = COALESCE(sqlc.narg('air_date'), air_date),
    vote_average = COALESCE(sqlc.narg('vote_average'), vote_average)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteSeason :exec
DELETE FROM tvshow.seasons WHERE id = $1;

-- name: DeleteSeasonsBySeries :exec
DELETE FROM tvshow.seasons WHERE series_id = $1;

-- name: UpsertSeason :one
INSERT INTO tvshow.seasons (
    series_id, tmdb_id, season_number,
    name, overview,
    names_i18n, overviews_i18n,
    poster_path, episode_count, air_date, vote_average
) VALUES (
    $1, $2, $3,
    $4, $5,
    $6, $7,
    $8, $9, $10, $11
)
ON CONFLICT (series_id, season_number) DO UPDATE SET
    tmdb_id = EXCLUDED.tmdb_id,
    name = EXCLUDED.name,
    overview = EXCLUDED.overview,
    names_i18n = EXCLUDED.names_i18n,
    overviews_i18n = EXCLUDED.overviews_i18n,
    poster_path = EXCLUDED.poster_path,
    episode_count = EXCLUDED.episode_count,
    air_date = EXCLUDED.air_date,
    vote_average = EXCLUDED.vote_average
RETURNING *;

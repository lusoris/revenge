-- Season Queries

-- name: GetSeasonByID :one
SELECT * FROM seasons WHERE id = $1;

-- name: GetSeasonByNumber :one
SELECT * FROM seasons WHERE series_id = $1 AND season_number = $2;

-- name: GetSeasonByTmdbID :one
SELECT * FROM seasons WHERE series_id = $1 AND tmdb_id = $2;

-- name: ListSeasons :many
SELECT * FROM seasons
WHERE series_id = $1
ORDER BY season_number ASC;

-- name: ListSeasonsWithSpecials :many
SELECT * FROM seasons
WHERE series_id = $1
ORDER BY
    CASE WHEN season_number = 0 THEN 1 ELSE 0 END,
    season_number ASC;

-- name: CountSeasons :one
SELECT COUNT(*) FROM seasons WHERE series_id = $1 AND season_number > 0;

-- name: CountSpecialSeasons :one
SELECT COUNT(*) FROM seasons WHERE series_id = $1 AND season_number = 0;

-- name: CreateSeason :one
INSERT INTO seasons (
    series_id, season_number, name, overview,
    air_date, year,
    poster_path, poster_blurhash,
    tmdb_id, tvdb_id
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: UpdateSeason :one
UPDATE seasons SET
    name = COALESCE(sqlc.narg('name'), name),
    overview = COALESCE(sqlc.narg('overview'), overview),
    air_date = COALESCE(sqlc.narg('air_date'), air_date),
    year = COALESCE(sqlc.narg('year'), year),
    poster_path = COALESCE(sqlc.narg('poster_path'), poster_path),
    poster_blurhash = COALESCE(sqlc.narg('poster_blurhash'), poster_blurhash),
    tmdb_id = COALESCE(sqlc.narg('tmdb_id'), tmdb_id),
    tvdb_id = COALESCE(sqlc.narg('tvdb_id'), tvdb_id)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteSeason :exec
DELETE FROM seasons WHERE id = $1;

-- name: DeleteSeasonsBySeries :exec
DELETE FROM seasons WHERE series_id = $1;

-- name: SeasonExistsByNumber :one
SELECT EXISTS(SELECT 1 FROM seasons WHERE series_id = $1 AND season_number = $2);

-- name: GetOrCreateSeason :one
INSERT INTO seasons (series_id, season_number, name)
VALUES ($1, $2, COALESCE($3, 'Season ' || $2::text))
ON CONFLICT (series_id, season_number) DO UPDATE SET updated_at = NOW()
RETURNING *;

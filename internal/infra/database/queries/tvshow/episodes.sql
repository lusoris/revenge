-- Episode Queries

-- name: GetEpisodeByID :one
SELECT * FROM episodes WHERE id = $1;

-- name: GetEpisodeByPath :one
SELECT * FROM episodes WHERE path = $1;

-- name: GetEpisodeByNumber :one
SELECT * FROM episodes
WHERE series_id = $1 AND season_number = $2 AND episode_number = $3;

-- name: GetEpisodeByAbsoluteNumber :one
SELECT * FROM episodes
WHERE series_id = $1 AND absolute_number = $2;

-- name: GetEpisodeByTmdbID :one
SELECT * FROM episodes WHERE series_id = $1 AND tmdb_id = $2;

-- name: ListEpisodes :many
SELECT * FROM episodes
WHERE series_id = $1
ORDER BY season_number ASC, episode_number ASC;

-- name: ListEpisodesBySeason :many
SELECT * FROM episodes
WHERE season_id = $1
ORDER BY episode_number ASC;

-- name: ListEpisodesBySeasonNumber :many
SELECT * FROM episodes
WHERE series_id = $1 AND season_number = $2
ORDER BY episode_number ASC;

-- name: ListRecentlyAddedEpisodes :many
SELECT e.*, s.title as series_title, s.poster_path as series_poster
FROM episodes e
JOIN series s ON e.series_id = s.id
WHERE s.library_id = ANY(@library_ids::uuid[])
ORDER BY e.date_added DESC
LIMIT $1;

-- name: ListRecentlyAiredEpisodes :many
SELECT e.*, s.title as series_title, s.poster_path as series_poster
FROM episodes e
JOIN series s ON e.series_id = s.id
WHERE s.library_id = ANY(@library_ids::uuid[])
  AND e.air_date IS NOT NULL
  AND e.air_date <= CURRENT_DATE
ORDER BY e.air_date DESC
LIMIT $1;

-- name: ListUpcomingEpisodes :many
SELECT e.*, s.title as series_title, s.poster_path as series_poster
FROM episodes e
JOIN series s ON e.series_id = s.id
WHERE s.library_id = ANY(@library_ids::uuid[])
  AND e.air_date IS NOT NULL
  AND e.air_date > CURRENT_DATE
ORDER BY e.air_date ASC
LIMIT $1;

-- name: CountEpisodes :one
SELECT COUNT(*) FROM episodes WHERE series_id = $1;

-- name: CountEpisodesBySeason :one
SELECT COUNT(*) FROM episodes WHERE season_id = $1;

-- name: CountSpecialEpisodes :one
SELECT COUNT(*) FROM episodes WHERE series_id = $1 AND season_number = 0;

-- name: CreateEpisode :one
INSERT INTO episodes (
    series_id, season_id, path, container, size_bytes, runtime_ticks,
    season_number, episode_number, absolute_number,
    title, overview, production_code,
    air_date, air_date_utc,
    community_rating, vote_count,
    still_path, still_blurhash,
    tmdb_id, imdb_id, tvdb_id
) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8, $9,
    $10, $11, $12,
    $13, $14,
    $15, $16,
    $17, $18,
    $19, $20, $21
) RETURNING *;

-- name: UpdateEpisode :one
UPDATE episodes SET
    season_id = COALESCE(sqlc.narg('season_id'), season_id),
    container = COALESCE(sqlc.narg('container'), container),
    size_bytes = COALESCE(sqlc.narg('size_bytes'), size_bytes),
    runtime_ticks = COALESCE(sqlc.narg('runtime_ticks'), runtime_ticks),
    season_number = COALESCE(sqlc.narg('season_number'), season_number),
    episode_number = COALESCE(sqlc.narg('episode_number'), episode_number),
    absolute_number = COALESCE(sqlc.narg('absolute_number'), absolute_number),
    title = COALESCE(sqlc.narg('title'), title),
    overview = COALESCE(sqlc.narg('overview'), overview),
    production_code = COALESCE(sqlc.narg('production_code'), production_code),
    air_date = COALESCE(sqlc.narg('air_date'), air_date),
    air_date_utc = COALESCE(sqlc.narg('air_date_utc'), air_date_utc),
    community_rating = COALESCE(sqlc.narg('community_rating'), community_rating),
    vote_count = COALESCE(sqlc.narg('vote_count'), vote_count),
    still_path = COALESCE(sqlc.narg('still_path'), still_path),
    still_blurhash = COALESCE(sqlc.narg('still_blurhash'), still_blurhash),
    tmdb_id = COALESCE(sqlc.narg('tmdb_id'), tmdb_id),
    imdb_id = COALESCE(sqlc.narg('imdb_id'), imdb_id),
    tvdb_id = COALESCE(sqlc.narg('tvdb_id'), tvdb_id),
    is_locked = COALESCE(sqlc.narg('is_locked'), is_locked)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateEpisodePlaybackStats :exec
UPDATE episodes SET
    last_played_at = NOW(),
    play_count = play_count + 1
WHERE id = $1;

-- name: DeleteEpisode :exec
DELETE FROM episodes WHERE id = $1;

-- name: DeleteEpisodesBySeries :exec
DELETE FROM episodes WHERE series_id = $1;

-- name: DeleteEpisodesBySeason :exec
DELETE FROM episodes WHERE season_id = $1;

-- name: EpisodeExistsByPath :one
SELECT EXISTS(SELECT 1 FROM episodes WHERE path = $1);

-- name: EpisodeExistsByNumber :one
SELECT EXISTS(SELECT 1 FROM episodes WHERE series_id = $1 AND season_number = $2 AND episode_number = $3);

-- name: ListEpisodePaths :many
SELECT id, path FROM episodes
WHERE series_id IN (SELECT id FROM series WHERE tv_library_id = $1);

-- name: GetNextEpisode :one
SELECT * FROM episodes
WHERE series_id = $1
  AND (season_number > $2 OR (season_number = $2 AND episode_number > $3))
  AND season_number > 0
ORDER BY season_number ASC, episode_number ASC
LIMIT 1;

-- name: GetPreviousEpisode :one
SELECT * FROM episodes
WHERE series_id = $1
  AND (season_number < $2 OR (season_number = $2 AND episode_number < $3))
  AND season_number > 0
ORDER BY season_number DESC, episode_number DESC
LIMIT 1;

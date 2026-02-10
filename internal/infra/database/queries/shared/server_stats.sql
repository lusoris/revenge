-- name: UpsertServerStat :exec
-- Upsert a single server statistic
INSERT INTO
    shared.server_stats (
        stat_key,
        stat_value,
        computed_at
    )
VALUES ($1, $2, NOW()) ON CONFLICT (stat_key) DO
UPDATE
SET
    stat_value = EXCLUDED.stat_value,
    computed_at = NOW();

-- name: GetServerStat :one
-- Get a single server statistic by key
SELECT
    stat_key,
    stat_value,
    computed_at
FROM shared.server_stats
WHERE
    stat_key = $1;

-- name: GetAllServerStats :many
-- Get all server statistics
SELECT
    stat_key,
    stat_value,
    computed_at
FROM shared.server_stats
ORDER BY stat_key;

-- =============================================================================
-- Aggregate queries used by the stats aggregation worker
-- =============================================================================

-- name: CountActiveUsers :one
-- Count total active (non-deleted) users
SELECT COUNT(*) FROM shared.users WHERE deleted_at IS NULL;

-- name: CountActiveUsersLast24h :one
-- Count users active in the last 24 hours (by session activity)
SELECT COUNT(DISTINCT user_id)
FROM shared.sessions
WHERE
    last_activity_at > NOW() - INTERVAL '24 hours';

-- name: CountTotalLibraries :one
-- Count total enabled libraries
SELECT COUNT(*) FROM public.libraries WHERE enabled = true;

-- name: CountTotalMovies :one
-- Count total movies across all libraries
SELECT COUNT(*) FROM movie.movies;

-- name: CountTotalSeries :one
-- Count total TV series
SELECT COUNT(*) FROM tvshow.series;

-- name: CountTotalEpisodes :one
-- Count total TV episodes
SELECT COUNT(*) FROM tvshow.episodes;

-- name: CountTotalMovieWatches :one
-- Count total movie watch records
SELECT COALESCE(SUM(watch_count), 0)::bigint FROM movie.movie_watched;

-- name: CountTotalEpisodeWatches :one
-- Count total episode watch records
SELECT COALESCE(SUM(watch_count), 0)::bigint FROM tvshow.episode_watched;

-- name: SumMovieWatchDurationSeconds :one
-- Sum total movie watch duration in seconds
SELECT COALESCE(SUM(duration_seconds::bigint * watch_count::bigint), 0)::bigint
FROM movie.movie_watched
WHERE is_completed = true;

-- name: SumEpisodeWatchDurationSeconds :one
-- Sum total episode watch duration in seconds
SELECT COALESCE(SUM(duration_seconds::bigint * watch_count::bigint), 0)::bigint
FROM tvshow.episode_watched
WHERE is_completed = true;

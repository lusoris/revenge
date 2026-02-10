-- name: GetEpisodeWatchProgress :one
SELECT *
FROM tvshow.episode_watched
WHERE
    user_id = $1
    AND episode_id = $2;

-- name: ListWatchedEpisodesBySeries :many
SELECT ew.*, e.season_number, e.episode_number, e.title as episode_title
FROM tvshow.episode_watched ew
    JOIN tvshow.episodes e ON ew.episode_id = e.id
WHERE
    ew.user_id = $1
    AND e.series_id = $2
ORDER BY e.season_number ASC, e.episode_number ASC;

-- name: ListContinueWatchingSeries :many
SELECT DISTINCT
    ON (s.id) s.*,
    e.id as last_episode_id,
    e.season_number as last_season_number,
    e.episode_number as last_episode_number,
    e.title as last_episode_title,
    ew.progress_seconds,
    ew.duration_seconds,
    ew.last_watched_at
FROM tvshow.episode_watched ew
    JOIN tvshow.episodes e ON ew.episode_id = e.id
    JOIN tvshow.series s ON e.series_id = s.id
WHERE
    ew.user_id = $1
    AND NOT ew.is_completed
    AND ew.progress_seconds > 0
ORDER BY s.id, ew.last_watched_at DESC
LIMIT $2;

-- name: ListWatchedEpisodesByUser :many
SELECT
    ew.*,
    e.season_number,
    e.episode_number,
    e.title as episode_title,
    e.still_path as episode_still_path,
    s.id as series_id,
    s.title as series_title,
    s.poster_path as series_poster_path
FROM tvshow.episode_watched ew
    JOIN tvshow.episodes e ON ew.episode_id = e.id
    JOIN tvshow.series s ON e.series_id = s.id
WHERE
    ew.user_id = $1
    AND ew.is_completed = TRUE
ORDER BY ew.last_watched_at DESC
LIMIT $2
OFFSET
    $3;

-- name: GetSeriesWatchStats :one
SELECT
    COUNT(*) FILTER (
        WHERE
            ew.is_completed = TRUE
    ) as watched_count,
    COUNT(*) FILTER (
        WHERE
            ew.is_completed = FALSE
            AND ew.progress_seconds > 0
    ) as in_progress_count,
    SUM(ew.watch_count) as total_watches,
    (
        SELECT COUNT(*)
        FROM tvshow.episodes ep
        WHERE
            ep.series_id = $2
    ) as total_episodes
FROM tvshow.episode_watched ew
    JOIN tvshow.episodes e ON ew.episode_id = e.id
WHERE
    ew.user_id = $1
    AND e.series_id = $2;

-- name: GetUserTVStats :one
SELECT
    COUNT(DISTINCT e.series_id) as series_count,
    COUNT(*) FILTER (
        WHERE
            ew.is_completed = TRUE
    ) as episodes_watched,
    COUNT(*) FILTER (
        WHERE
            ew.is_completed = FALSE
            AND ew.progress_seconds > 0
    ) as episodes_in_progress,
    COALESCE(SUM(ew.watch_count), 0) as total_watches
FROM tvshow.episode_watched ew
    JOIN tvshow.episodes e ON ew.episode_id = e.id
WHERE
    ew.user_id = $1;

-- name: CreateOrUpdateWatchProgress :one
INSERT INTO
    tvshow.episode_watched (
        user_id,
        episode_id,
        progress_seconds,
        duration_seconds,
        is_completed,
        watch_count,
        last_watched_at
    )
VALUES ($1, $2, $3, $4, $5, 1, NOW()) ON CONFLICT (user_id, episode_id) DO
UPDATE
SET
    progress_seconds = $3,
    duration_seconds = $4,
    is_completed = $5,
    watch_count = CASE
        WHEN $5
        AND NOT tvshow.episode_watched.is_completed THEN tvshow.episode_watched.watch_count + 1
        ELSE tvshow.episode_watched.watch_count
    END,
    last_watched_at = NOW() RETURNING *;

-- name: MarkEpisodeWatched :one
INSERT INTO
    tvshow.episode_watched (
        user_id,
        episode_id,
        progress_seconds,
        duration_seconds,
        is_completed,
        watch_count,
        last_watched_at
    )
VALUES (
        $1,
        $2,
        $3,
        $3,
        TRUE,
        1,
        NOW()
    ) ON CONFLICT (user_id, episode_id) DO
UPDATE
SET
    progress_seconds = $3,
    duration_seconds = $3,
    is_completed = TRUE,
    watch_count = tvshow.episode_watched.watch_count + 1,
    last_watched_at = NOW() RETURNING *;

-- name: DeleteWatchProgress :exec
DELETE FROM tvshow.episode_watched
WHERE
    user_id = $1
    AND episode_id = $2;

-- name: DeleteSeriesWatchProgress :exec
DELETE FROM tvshow.episode_watched
WHERE
    user_id = $1
    AND episode_id IN (
        SELECT id
        FROM tvshow.episodes
        WHERE
            series_id = $2
    );

-- name: MarkEpisodesWatchedBulk :execrows
WITH episode_durations AS (
    SELECT id, COALESCE(runtime * 60, 2700) AS duration_secs
    FROM tvshow.episodes
    WHERE id = ANY(@episode_ids::uuid[])
)
INSERT INTO tvshow.episode_watched (
    user_id, episode_id, progress_seconds, duration_seconds, is_completed, watch_count, last_watched_at
)
SELECT @user_id, ed.id, ed.duration_secs, ed.duration_secs, TRUE, 1, NOW()
FROM episode_durations ed
ON CONFLICT (user_id, episode_id) DO UPDATE SET
    progress_seconds = EXCLUDED.progress_seconds,
    duration_seconds = EXCLUDED.duration_seconds,
    is_completed = TRUE,
    watch_count = tvshow.episode_watched.watch_count + 1,
    last_watched_at = NOW();

-- name: GetNextUnwatchedEpisode :one
SELECT e.*
FROM tvshow.episodes e
    LEFT JOIN tvshow.episode_watched ew ON e.id = ew.episode_id
    AND ew.user_id = $1
WHERE
    e.series_id = $2
    AND (
        ew.is_completed IS NULL
        OR ew.is_completed = FALSE
    )
ORDER BY e.season_number ASC, e.episode_number ASC
LIMIT 1;

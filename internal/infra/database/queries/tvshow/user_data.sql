-- TV Show User Data Queries

-- Series User Ratings

-- name: GetSeriesUserRating :one
SELECT * FROM series_user_ratings
WHERE user_id = $1 AND series_id = $2;

-- name: ListUserSeriesRatings :many
SELECT sur.*, s.title, s.poster_path
FROM series_user_ratings sur
JOIN series s ON sur.series_id = s.id
WHERE sur.user_id = $1
ORDER BY sur.updated_at DESC
LIMIT $2 OFFSET $3;

-- name: SetSeriesUserRating :one
INSERT INTO series_user_ratings (user_id, series_id, rating, review)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id, series_id) DO UPDATE SET
    rating = $3,
    review = COALESCE($4, series_user_ratings.review),
    updated_at = NOW()
RETURNING *;

-- name: DeleteSeriesUserRating :exec
DELETE FROM series_user_ratings
WHERE user_id = $1 AND series_id = $2;

-- name: GetAverageSeriesUserRating :one
SELECT AVG(rating) as avg_rating, COUNT(*) as rating_count
FROM series_user_ratings
WHERE series_id = $1;

-- Episode User Ratings

-- name: GetEpisodeUserRating :one
SELECT * FROM episode_user_ratings
WHERE user_id = $1 AND episode_id = $2;

-- name: SetEpisodeUserRating :one
INSERT INTO episode_user_ratings (user_id, episode_id, rating)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, episode_id) DO UPDATE SET
    rating = $3,
    updated_at = NOW()
RETURNING *;

-- name: DeleteEpisodeUserRating :exec
DELETE FROM episode_user_ratings
WHERE user_id = $1 AND episode_id = $2;

-- Series Favorites

-- name: IsSeriesFavorite :one
SELECT EXISTS(
    SELECT 1 FROM series_favorites
    WHERE user_id = $1 AND series_id = $2
);

-- name: AddSeriesFavorite :exec
INSERT INTO series_favorites (user_id, series_id)
VALUES ($1, $2)
ON CONFLICT (user_id, series_id) DO NOTHING;

-- name: RemoveSeriesFavorite :exec
DELETE FROM series_favorites
WHERE user_id = $1 AND series_id = $2;

-- name: ListUserFavoriteSeries :many
SELECT s.* FROM series s
JOIN series_favorites sf ON s.id = sf.series_id
WHERE sf.user_id = $1
ORDER BY sf.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountUserFavoriteSeries :one
SELECT COUNT(*) FROM series_favorites WHERE user_id = $1;

-- Episode Watch History

-- name: GetEpisodeWatchHistory :one
SELECT * FROM episode_watch_history
WHERE user_id = $1 AND episode_id = $2 AND completed = false
ORDER BY last_updated_at DESC
LIMIT 1;

-- name: GetCompletedEpisodeWatchHistory :many
SELECT * FROM episode_watch_history
WHERE user_id = $1 AND episode_id = $2 AND completed = true
ORDER BY completed_at DESC;

-- name: ListUserEpisodeWatchHistory :many
SELECT ewh.*, e.title as episode_title, e.season_number, e.episode_number, e.runtime_ticks,
       s.id as series_id, s.title as series_title, s.poster_path as series_poster
FROM episode_watch_history ewh
JOIN episodes e ON ewh.episode_id = e.id
JOIN series s ON e.series_id = s.id
WHERE ewh.user_id = $1
ORDER BY ewh.last_updated_at DESC
LIMIT $2 OFFSET $3;

-- name: ListResumeableEpisodes :many
SELECT ewh.*, e.title as episode_title, e.season_number, e.episode_number, e.runtime_ticks, e.still_path,
       s.id as series_id, s.title as series_title, s.poster_path as series_poster
FROM episode_watch_history ewh
JOIN episodes e ON ewh.episode_id = e.id
JOIN series s ON e.series_id = s.id
WHERE ewh.user_id = $1
  AND ewh.completed = false
  AND ewh.played_percentage > 5
  AND ewh.played_percentage < 90
ORDER BY ewh.last_updated_at DESC
LIMIT $2;

-- name: CreateEpisodeWatchHistory :one
INSERT INTO episode_watch_history (
    user_id, profile_id, episode_id, position_ticks, duration_ticks,
    device_name, device_type, client_name, play_method
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateEpisodeWatchHistory :one
UPDATE episode_watch_history SET
    position_ticks = $2,
    duration_ticks = COALESCE($3, duration_ticks),
    last_updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: MarkEpisodeWatchHistoryCompleted :exec
UPDATE episode_watch_history SET
    completed = true,
    completed_at = NOW(),
    position_ticks = COALESCE(duration_ticks, position_ticks)
WHERE id = $1;

-- name: DeleteEpisodeWatchHistory :exec
DELETE FROM episode_watch_history WHERE id = $1;

-- name: IsEpisodeWatched :one
SELECT EXISTS(
    SELECT 1 FROM episode_watch_history
    WHERE user_id = $1 AND episode_id = $2 AND completed = true
);

-- name: CountUserWatchedEpisodes :one
SELECT COUNT(DISTINCT episode_id) FROM episode_watch_history
WHERE user_id = $1 AND completed = true;

-- name: CountUserWatchedEpisodesBySeries :one
SELECT COUNT(DISTINCT ewh.episode_id) FROM episode_watch_history ewh
JOIN episodes e ON ewh.episode_id = e.id
WHERE ewh.user_id = $1 AND e.series_id = $2 AND ewh.completed = true;

-- Series Watchlist

-- name: IsSeriesInWatchlist :one
SELECT EXISTS(
    SELECT 1 FROM series_watchlist
    WHERE user_id = $1 AND series_id = $2
);

-- name: AddSeriesToWatchlist :exec
INSERT INTO series_watchlist (user_id, series_id, sort_order)
VALUES ($1, $2, COALESCE($3, (SELECT COALESCE(MAX(sort_order), 0) + 1 FROM series_watchlist WHERE user_id = $1)))
ON CONFLICT (user_id, series_id) DO NOTHING;

-- name: RemoveSeriesFromWatchlist :exec
DELETE FROM series_watchlist
WHERE user_id = $1 AND series_id = $2;

-- name: ListUserSeriesWatchlist :many
SELECT s.* FROM series s
JOIN series_watchlist sw ON s.id = sw.series_id
WHERE sw.user_id = $1
ORDER BY sw.sort_order ASC, sw.added_at DESC
LIMIT $2 OFFSET $3;

-- name: CountUserSeriesWatchlist :one
SELECT COUNT(*) FROM series_watchlist WHERE user_id = $1;

-- name: ReorderSeriesWatchlist :exec
UPDATE series_watchlist SET sort_order = $3
WHERE user_id = $1 AND series_id = $2;

-- Series Watch Progress (Continue Watching)

-- name: GetSeriesWatchProgress :one
SELECT * FROM series_watch_progress
WHERE user_id = $1 AND series_id = $2;

-- name: ListContinueWatchingSeries :many
SELECT swp.*, s.title, s.poster_path, s.backdrop_path,
       e.id as next_episode_id, e.title as next_episode_title,
       e.season_number as next_season, e.episode_number as next_episode,
       e.still_path as next_episode_still
FROM series_watch_progress swp
JOIN series s ON swp.series_id = s.id
LEFT JOIN episodes e ON e.id = swp.last_episode_id
WHERE swp.user_id = $1
  AND swp.is_watching = true
ORDER BY swp.last_watched_at DESC
LIMIT $2;

-- name: ListCompletedSeries :many
SELECT swp.*, s.title, s.poster_path
FROM series_watch_progress swp
JOIN series s ON swp.series_id = s.id
WHERE swp.user_id = $1
  AND swp.completed_at IS NOT NULL
ORDER BY swp.completed_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateSeriesWatchProgress :one
INSERT INTO series_watch_progress (
    user_id, series_id, last_episode_id, last_season_number, last_episode_number,
    total_episodes, watched_episodes, is_watching, started_at, last_watched_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
ON CONFLICT (user_id, series_id) DO UPDATE SET
    last_episode_id = EXCLUDED.last_episode_id,
    last_season_number = EXCLUDED.last_season_number,
    last_episode_number = EXCLUDED.last_episode_number,
    total_episodes = EXCLUDED.total_episodes,
    watched_episodes = EXCLUDED.watched_episodes,
    is_watching = EXCLUDED.is_watching,
    last_watched_at = NOW()
RETURNING *;

-- name: DeleteSeriesWatchProgress :exec
DELETE FROM series_watch_progress
WHERE user_id = $1 AND series_id = $2;

-- Series External Ratings

-- name: GetSeriesExternalRatings :many
SELECT * FROM series_external_ratings
WHERE series_id = $1;

-- name: UpsertSeriesExternalRating :exec
INSERT INTO series_external_ratings (series_id, source, rating, vote_count, certified)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (series_id, source) DO UPDATE SET
    rating = $3,
    vote_count = $4,
    certified = $5,
    last_updated = NOW();

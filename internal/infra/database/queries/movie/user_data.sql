-- Movie User Data Queries

-- User Ratings

-- name: GetMovieUserRating :one
SELECT * FROM movie_user_ratings
WHERE user_id = $1 AND movie_id = $2;

-- name: ListUserMovieRatings :many
SELECT mur.*, m.title, m.poster_path
FROM movie_user_ratings mur
JOIN movies m ON mur.movie_id = m.id
WHERE mur.user_id = $1
ORDER BY mur.updated_at DESC
LIMIT $2 OFFSET $3;

-- name: SetMovieUserRating :one
INSERT INTO movie_user_ratings (user_id, movie_id, rating, review)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id, movie_id) DO UPDATE SET
    rating = $3,
    review = COALESCE($4, movie_user_ratings.review),
    updated_at = NOW()
RETURNING *;

-- name: DeleteMovieUserRating :exec
DELETE FROM movie_user_ratings
WHERE user_id = $1 AND movie_id = $2;

-- name: GetAverageUserRating :one
SELECT AVG(rating) as avg_rating, COUNT(*) as rating_count
FROM movie_user_ratings
WHERE movie_id = $1;

-- Favorites

-- name: IsMovieFavorite :one
SELECT EXISTS(
    SELECT 1 FROM movie_favorites
    WHERE user_id = $1 AND movie_id = $2
);

-- name: AddMovieFavorite :exec
INSERT INTO movie_favorites (user_id, movie_id)
VALUES ($1, $2)
ON CONFLICT (user_id, movie_id) DO NOTHING;

-- name: RemoveMovieFavorite :exec
DELETE FROM movie_favorites
WHERE user_id = $1 AND movie_id = $2;

-- name: ListUserFavoriteMovies :many
SELECT m.* FROM movies m
JOIN movie_favorites mf ON m.id = mf.movie_id
WHERE mf.user_id = $1
ORDER BY mf.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountUserFavoriteMovies :one
SELECT COUNT(*) FROM movie_favorites WHERE user_id = $1;

-- Watch History

-- name: GetMovieWatchHistory :one
SELECT * FROM movie_watch_history
WHERE user_id = $1 AND movie_id = $2 AND completed = false
ORDER BY last_updated_at DESC
LIMIT 1;

-- name: GetCompletedMovieWatchHistory :many
SELECT * FROM movie_watch_history
WHERE user_id = $1 AND movie_id = $2 AND completed = true
ORDER BY completed_at DESC;

-- name: ListUserWatchHistory :many
SELECT mwh.*, m.title, m.poster_path, m.runtime_ticks
FROM movie_watch_history mwh
JOIN movies m ON mwh.movie_id = m.id
WHERE mwh.user_id = $1
ORDER BY mwh.last_updated_at DESC
LIMIT $2 OFFSET $3;

-- name: ListResumeableMovies :many
SELECT mwh.*, m.title, m.poster_path, m.runtime_ticks
FROM movie_watch_history mwh
JOIN movies m ON mwh.movie_id = m.id
WHERE mwh.user_id = $1
  AND mwh.completed = false
  AND mwh.played_percentage > 5   -- At least 5% watched
  AND mwh.played_percentage < 90  -- Less than 90% watched
ORDER BY mwh.last_updated_at DESC
LIMIT $2;

-- name: CreateWatchHistory :one
INSERT INTO movie_watch_history (
    user_id, profile_id, movie_id, position_ticks, duration_ticks,
    device_name, device_type, client_name, play_method
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: UpdateWatchHistory :one
UPDATE movie_watch_history SET
    position_ticks = $2,
    duration_ticks = COALESCE($3, duration_ticks),
    last_updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: MarkWatchHistoryCompleted :exec
UPDATE movie_watch_history SET
    completed = true,
    completed_at = NOW(),
    position_ticks = COALESCE(duration_ticks, position_ticks)
WHERE id = $1;

-- name: DeleteWatchHistory :exec
DELETE FROM movie_watch_history WHERE id = $1;

-- name: IsMovieWatched :one
SELECT EXISTS(
    SELECT 1 FROM movie_watch_history
    WHERE user_id = $1 AND movie_id = $2 AND completed = true
);

-- name: CountUserWatchedMovies :one
SELECT COUNT(DISTINCT movie_id) FROM movie_watch_history
WHERE user_id = $1 AND completed = true;

-- Watchlist

-- name: IsMovieInWatchlist :one
SELECT EXISTS(
    SELECT 1 FROM movie_watchlist
    WHERE user_id = $1 AND movie_id = $2
);

-- name: AddMovieToWatchlist :exec
INSERT INTO movie_watchlist (user_id, movie_id, sort_order)
VALUES ($1, $2, COALESCE($3, (SELECT COALESCE(MAX(sort_order), 0) + 1 FROM movie_watchlist WHERE user_id = $1)))
ON CONFLICT (user_id, movie_id) DO NOTHING;

-- name: RemoveMovieFromWatchlist :exec
DELETE FROM movie_watchlist
WHERE user_id = $1 AND movie_id = $2;

-- name: ListUserWatchlist :many
SELECT m.* FROM movies m
JOIN movie_watchlist mw ON m.id = mw.movie_id
WHERE mw.user_id = $1
ORDER BY mw.sort_order ASC, mw.added_at DESC
LIMIT $2 OFFSET $3;

-- name: CountUserWatchlist :one
SELECT COUNT(*) FROM movie_watchlist WHERE user_id = $1;

-- name: ReorderWatchlist :exec
UPDATE movie_watchlist SET sort_order = $3
WHERE user_id = $1 AND movie_id = $2;

-- External Ratings

-- name: GetMovieExternalRatings :many
SELECT * FROM movie_external_ratings
WHERE movie_id = $1;

-- name: UpsertMovieExternalRating :exec
INSERT INTO movie_external_ratings (movie_id, source, rating, vote_count, certified)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (movie_id, source) DO UPDATE SET
    rating = $3,
    vote_count = $4,
    certified = $5,
    last_updated = NOW();

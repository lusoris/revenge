-- Movie Genre Queries

-- name: GetMovieGenres :many
SELECT g.* FROM genres g
JOIN movie_genres mg ON g.id = mg.genre_id
WHERE mg.movie_id = $1
ORDER BY g.name ASC;

-- name: LinkMovieGenre :exec
INSERT INTO movie_genres (movie_id, genre_id)
VALUES ($1, $2)
ON CONFLICT (movie_id, genre_id) DO NOTHING;

-- name: UnlinkMovieGenres :exec
DELETE FROM movie_genres WHERE movie_id = $1;

-- name: ListMoviesByGenre :many
SELECT m.* FROM movies m
JOIN movie_genres mg ON m.id = mg.movie_id
WHERE mg.genre_id = $1
ORDER BY m.sort_title ASC
LIMIT $2 OFFSET $3;

-- name: CountMoviesByGenre :one
SELECT COUNT(*) FROM movie_genres WHERE genre_id = $1;

-- name: ListGenresWithMovieCounts :many
SELECT g.*, COUNT(mg.movie_id) as movie_count
FROM genres g
LEFT JOIN movie_genres mg ON g.id = mg.genre_id
WHERE g.domain = 'movie'
GROUP BY g.id
ORDER BY movie_count DESC, g.name ASC;

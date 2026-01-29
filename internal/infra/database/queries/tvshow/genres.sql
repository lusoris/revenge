-- TV Show Genre Queries

-- name: GetSeriesGenres :many
SELECT g.* FROM tvshow_genres g
JOIN series_genre_link sg ON g.id = sg.genre_id
WHERE sg.series_id = $1
ORDER BY g.name ASC;

-- name: LinkSeriesGenre :exec
INSERT INTO series_genre_link (series_id, genre_id)
VALUES ($1, $2)
ON CONFLICT (series_id, genre_id) DO NOTHING;

-- name: UnlinkSeriesGenres :exec
DELETE FROM series_genre_link WHERE series_id = $1;

-- name: ListSeriesByGenre :many
SELECT s.* FROM series s
JOIN series_genre_link sg ON s.id = sg.series_id
WHERE sg.genre_id = $1
ORDER BY s.sort_title ASC
LIMIT $2 OFFSET $3;

-- name: CountSeriesByGenre :one
SELECT COUNT(*) FROM series_genre_link WHERE genre_id = $1;

-- name: ListTVShowGenresWithCounts :many
SELECT g.*, COUNT(sg.series_id) as series_count
FROM tvshow_genres g
LEFT JOIN series_genre_link sg ON g.id = sg.genre_id
GROUP BY g.id
ORDER BY series_count DESC, g.name ASC;

-- name: GetTVShowGenreByID :one
SELECT * FROM tvshow_genres WHERE id = $1;

-- name: GetTVShowGenreByName :one
SELECT * FROM tvshow_genres WHERE name = $1;

-- name: GetTVShowGenreByTmdbID :one
SELECT * FROM tvshow_genres WHERE tmdb_id = $1;

-- name: ListTVShowGenres :many
SELECT * FROM tvshow_genres ORDER BY name ASC;

-- name: CreateTVShowGenre :one
INSERT INTO tvshow_genres (name, tmdb_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetOrCreateTVShowGenre :one
INSERT INTO tvshow_genres (name, tmdb_id)
VALUES ($1, $2)
ON CONFLICT (name) DO UPDATE SET tmdb_id = COALESCE(EXCLUDED.tmdb_id, tvshow_genres.tmdb_id)
RETURNING *;

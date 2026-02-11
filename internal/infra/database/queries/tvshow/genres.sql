-- name: ListSeriesGenres :many
SELECT *
FROM tvshow.series_genres
WHERE
    series_id = $1
ORDER BY name ASC;

-- name: ListDistinctSeriesGenres :many
SELECT slug, name, COUNT(DISTINCT series_id)::bigint AS item_count
FROM tvshow.series_genres
GROUP BY slug, name
ORDER BY name ASC;

-- name: AddSeriesGenre :exec
INSERT INTO
    tvshow.series_genres (series_id, slug, name)
VALUES ($1, $2, $3) ON CONFLICT (series_id, slug) DO NOTHING;

-- name: DeleteSeriesGenres :exec
DELETE FROM tvshow.series_genres WHERE series_id = $1;

-- name: ListSeriesByGenre :many
SELECT s.*
FROM tvshow.series s
    JOIN tvshow.series_genres sg ON s.id = sg.series_id
WHERE
    sg.slug = $1
ORDER BY s.popularity DESC NULLS LAST
LIMIT $2
OFFSET
    $3;

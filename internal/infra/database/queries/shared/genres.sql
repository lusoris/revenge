-- name: GetGenreByID :one
SELECT * FROM genres WHERE id = $1;

-- name: GetGenreBySlug :one
SELECT * FROM genres WHERE domain = $1 AND slug = $2;

-- name: ListGenresByDomain :many
SELECT * FROM genres
WHERE domain = $1
ORDER BY name ASC;

-- name: ListTopLevelGenresByDomain :many
SELECT * FROM genres
WHERE domain = $1 AND parent_id IS NULL
ORDER BY name ASC;

-- name: ListChildGenres :many
SELECT * FROM genres
WHERE parent_id = $1
ORDER BY name ASC;

-- name: SearchGenres :many
SELECT * FROM genres
WHERE domain = $1 AND name ILIKE '%' || $2 || '%'
ORDER BY
    CASE WHEN name ILIKE $2 THEN 0 ELSE 1 END,
    name ASC
LIMIT $3;

-- name: CreateGenre :one
INSERT INTO genres (domain, name, slug, description, parent_id, external_ids)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateGenre :one
UPDATE genres SET
    name = COALESCE(sqlc.narg('name'), name),
    slug = COALESCE(sqlc.narg('slug'), slug),
    description = COALESCE(sqlc.narg('description'), description),
    parent_id = sqlc.narg('parent_id'),
    external_ids = COALESCE(sqlc.narg('external_ids'), external_ids)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteGenre :exec
DELETE FROM genres WHERE id = $1;

-- name: GenreExists :one
SELECT EXISTS(SELECT 1 FROM genres WHERE id = $1);

-- name: GenreSlugExists :one
SELECT EXISTS(SELECT 1 FROM genres WHERE domain = $1 AND slug = $2);

-- name: CountGenresByDomain :one
SELECT COUNT(*) FROM genres WHERE domain = $1;

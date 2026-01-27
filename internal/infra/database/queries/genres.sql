-- name: GetGenreByID :one
SELECT id, domain, name, slug, description, parent_id, external_ids, created_at, updated_at
FROM genres
WHERE id = $1;

-- name: GetGenreBySlug :one
SELECT id, domain, name, slug, description, parent_id, external_ids, created_at, updated_at
FROM genres
WHERE domain = $1 AND slug = $2;

-- name: ListGenresByDomain :many
SELECT id, domain, name, slug, description, parent_id, external_ids, created_at, updated_at
FROM genres
WHERE domain = $1
ORDER BY name ASC;

-- name: ListTopLevelGenresByDomain :many
SELECT id, domain, name, slug, description, parent_id, external_ids, created_at, updated_at
FROM genres
WHERE domain = $1 AND parent_id IS NULL
ORDER BY name ASC;

-- name: ListChildGenres :many
SELECT id, domain, name, slug, description, parent_id, external_ids, created_at, updated_at
FROM genres
WHERE parent_id = $1
ORDER BY name ASC;

-- name: ListGenresForMediaItem :many
SELECT g.id, g.domain, g.name, g.slug, g.description, g.parent_id, g.external_ids, g.created_at, g.updated_at
FROM genres g
INNER JOIN media_item_genres mig ON mig.genre_id = g.id
WHERE mig.media_item_id = $1
ORDER BY g.name ASC;

-- name: SearchGenres :many
SELECT id, domain, name, slug, description, parent_id, external_ids, created_at, updated_at
FROM genres
WHERE domain = $1 AND name ILIKE '%' || $2 || '%'
ORDER BY
    CASE WHEN name ILIKE $2 THEN 0 ELSE 1 END,  -- Exact matches first
    name ASC
LIMIT $3;

-- name: SearchGenresAllDomains :many
SELECT id, domain, name, slug, description, parent_id, external_ids, created_at, updated_at
FROM genres
WHERE name ILIKE '%' || $1 || '%'
ORDER BY
    CASE WHEN name ILIKE $1 THEN 0 ELSE 1 END,
    domain ASC,
    name ASC
LIMIT $2;

-- name: CreateGenre :one
INSERT INTO genres (domain, name, slug, description, parent_id, external_ids)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, domain, name, slug, description, parent_id, external_ids, created_at, updated_at;

-- name: UpdateGenre :one
UPDATE genres
SET
    name = COALESCE(sqlc.narg('name'), name),
    slug = COALESCE(sqlc.narg('slug'), slug),
    description = COALESCE(sqlc.narg('description'), description),
    parent_id = sqlc.narg('parent_id'),
    external_ids = COALESCE(sqlc.narg('external_ids'), external_ids),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING id, domain, name, slug, description, parent_id, external_ids, created_at, updated_at;

-- name: DeleteGenre :exec
DELETE FROM genres WHERE id = $1;

-- name: GenreExists :one
SELECT EXISTS(SELECT 1 FROM genres WHERE id = $1);

-- name: GenreSlugExists :one
SELECT EXISTS(SELECT 1 FROM genres WHERE domain = $1 AND slug = $2);

-- name: CountGenresByDomain :one
SELECT COUNT(*) FROM genres WHERE domain = $1;

-- Media Item Genre Associations

-- name: AssignGenreToMediaItem :exec
INSERT INTO media_item_genres (media_item_id, genre_id, source, confidence)
VALUES ($1, $2, $3, $4)
ON CONFLICT (media_item_id, genre_id) DO UPDATE
SET source = EXCLUDED.source, confidence = EXCLUDED.confidence;

-- name: RemoveGenreFromMediaItem :exec
DELETE FROM media_item_genres
WHERE media_item_id = $1 AND genre_id = $2;

-- name: RemoveAllGenresFromMediaItem :exec
DELETE FROM media_item_genres
WHERE media_item_id = $1;

-- name: GetMediaItemGenres :many
SELECT mig.media_item_id, mig.genre_id, mig.source, mig.confidence, mig.created_at
FROM media_item_genres mig
WHERE mig.media_item_id = $1;

-- name: CountMediaItemsWithGenre :one
SELECT COUNT(DISTINCT media_item_id)
FROM media_item_genres
WHERE genre_id = $1;

-- name: ListGenresWithCounts :many
SELECT
    g.id, g.domain, g.name, g.slug, g.description, g.parent_id, g.external_ids, g.created_at, g.updated_at,
    COUNT(mig.media_item_id) as item_count
FROM genres g
LEFT JOIN media_item_genres mig ON mig.genre_id = g.id
WHERE g.domain = $1
GROUP BY g.id
ORDER BY item_count DESC, g.name ASC
LIMIT $2;

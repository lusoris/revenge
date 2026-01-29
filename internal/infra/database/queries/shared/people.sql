-- name: GetPersonByID :one
SELECT * FROM people WHERE id = $1;

-- name: GetPersonByTmdbID :one
SELECT * FROM people WHERE tmdb_id = $1;

-- name: GetPersonByImdbID :one
SELECT * FROM people WHERE imdb_id = $1;

-- name: SearchPeople :many
SELECT * FROM people
WHERE name ILIKE '%' || $1 || '%'
ORDER BY
    CASE WHEN name ILIKE $1 THEN 0 ELSE 1 END,
    name ASC
LIMIT $2;

-- name: ListPeople :many
SELECT * FROM people
ORDER BY name ASC
LIMIT $1 OFFSET $2;

-- name: CreatePerson :one
INSERT INTO people (
    name, sort_name, original_name,
    biography, birthdate, deathdate, birthplace, gender,
    primary_image_url, primary_image_blurhash,
    tmdb_id, imdb_id, tvdb_id, musicbrainz_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
) RETURNING *;

-- name: UpdatePerson :one
UPDATE people SET
    name = COALESCE(sqlc.narg('name'), name),
    sort_name = COALESCE(sqlc.narg('sort_name'), sort_name),
    original_name = COALESCE(sqlc.narg('original_name'), original_name),
    biography = COALESCE(sqlc.narg('biography'), biography),
    birthdate = COALESCE(sqlc.narg('birthdate'), birthdate),
    deathdate = COALESCE(sqlc.narg('deathdate'), deathdate),
    birthplace = COALESCE(sqlc.narg('birthplace'), birthplace),
    gender = COALESCE(sqlc.narg('gender'), gender),
    primary_image_url = COALESCE(sqlc.narg('primary_image_url'), primary_image_url),
    primary_image_blurhash = COALESCE(sqlc.narg('primary_image_blurhash'), primary_image_blurhash),
    tmdb_id = COALESCE(sqlc.narg('tmdb_id'), tmdb_id),
    imdb_id = COALESCE(sqlc.narg('imdb_id'), imdb_id),
    tvdb_id = COALESCE(sqlc.narg('tvdb_id'), tvdb_id),
    musicbrainz_id = COALESCE(sqlc.narg('musicbrainz_id'), musicbrainz_id)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeletePerson :exec
DELETE FROM people WHERE id = $1;

-- name: CountPeople :one
SELECT COUNT(*) FROM people;

-- name: PersonExistsByTmdbID :one
SELECT EXISTS(SELECT 1 FROM people WHERE tmdb_id = $1);

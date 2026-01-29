-- Movie Libraries queries

-- name: GetMovieLibraryByID :one
SELECT * FROM movie_libraries WHERE id = $1;

-- name: ListMovieLibraries :many
SELECT * FROM movie_libraries
ORDER BY sort_order, name
LIMIT $1 OFFSET $2;

-- name: ListMovieLibrariesByOwner :many
SELECT * FROM movie_libraries
WHERE owner_user_id = $1
ORDER BY sort_order, name;

-- name: ListAccessibleMovieLibraries :many
SELECT ml.* FROM movie_libraries ml
LEFT JOIN movie_library_access mla ON mla.library_id = ml.id AND mla.user_id = $1
WHERE ml.is_private = false
   OR ml.owner_user_id = $1
   OR mla.user_id IS NOT NULL
ORDER BY ml.sort_order, ml.name;

-- name: CreateMovieLibrary :one
INSERT INTO movie_libraries (
    name,
    paths,
    scan_enabled,
    scan_interval_hours,
    preferred_language,
    tmdb_enabled,
    imdb_enabled,
    download_trailers,
    download_backdrops,
    download_nfo,
    generate_chapters,
    is_private,
    owner_user_id,
    sort_order,
    icon
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
)
RETURNING *;

-- name: UpdateMovieLibrary :one
UPDATE movie_libraries
SET
    name = $2,
    paths = $3,
    scan_enabled = $4,
    scan_interval_hours = $5,
    preferred_language = $6,
    tmdb_enabled = $7,
    imdb_enabled = $8,
    download_trailers = $9,
    download_backdrops = $10,
    download_nfo = $11,
    generate_chapters = $12,
    is_private = $13,
    owner_user_id = $14,
    sort_order = $15,
    icon = $16
WHERE id = $1
RETURNING *;

-- name: DeleteMovieLibrary :exec
DELETE FROM movie_libraries WHERE id = $1;

-- name: UpdateMovieLibraryScanStatus :exec
UPDATE movie_libraries
SET
    last_scan_at = $2,
    last_scan_duration = $3
WHERE id = $1;

-- name: CountMoviesByMovieLibrary :one
SELECT COUNT(*) FROM movies WHERE movie_library_id = $1;

-- name: GrantMovieLibraryAccess :exec
INSERT INTO movie_library_access (library_id, user_id, can_manage)
VALUES ($1, $2, $3)
ON CONFLICT (library_id, user_id)
DO UPDATE SET can_manage = EXCLUDED.can_manage;

-- name: RevokeMovieLibraryAccess :exec
DELETE FROM movie_library_access WHERE library_id = $1 AND user_id = $2;

-- name: ListMovieLibraryAccess :many
SELECT * FROM movie_library_access WHERE library_id = $1;

-- name: HasMovieLibraryAccess :one
SELECT EXISTS(
    SELECT 1 FROM movie_library_access
    WHERE library_id = $1 AND user_id = $2
) AS has_access;

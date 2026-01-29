-- TV Libraries queries

-- name: GetTVLibraryByID :one
SELECT * FROM tv_libraries WHERE id = $1;

-- name: ListTVLibraries :many
SELECT * FROM tv_libraries
ORDER BY sort_order, name
LIMIT $1 OFFSET $2;

-- name: ListTVLibrariesByOwner :many
SELECT * FROM tv_libraries
WHERE owner_user_id = $1
ORDER BY sort_order, name;

-- name: ListAccessibleTVLibraries :many
SELECT tl.* FROM tv_libraries tl
LEFT JOIN tv_library_access tla ON tla.library_id = tl.id AND tla.user_id = $1
WHERE tl.is_private = false
   OR tl.owner_user_id = $1
   OR tla.user_id IS NOT NULL
ORDER BY tl.sort_order, tl.name;

-- name: CreateTVLibrary :one
INSERT INTO tv_libraries (
    name,
    paths,
    scan_enabled,
    scan_interval_hours,
    preferred_language,
    tmdb_enabled,
    tvdb_enabled,
    download_backdrops,
    download_nfo,
    generate_chapters,
    season_folder_format,
    episode_naming_format,
    auto_add_missing,
    is_private,
    owner_user_id,
    sort_order,
    icon
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
)
RETURNING *;

-- name: UpdateTVLibrary :one
UPDATE tv_libraries
SET
    name = $2,
    paths = $3,
    scan_enabled = $4,
    scan_interval_hours = $5,
    preferred_language = $6,
    tmdb_enabled = $7,
    tvdb_enabled = $8,
    download_backdrops = $9,
    download_nfo = $10,
    generate_chapters = $11,
    season_folder_format = $12,
    episode_naming_format = $13,
    auto_add_missing = $14,
    is_private = $15,
    owner_user_id = $16,
    sort_order = $17,
    icon = $18
WHERE id = $1
RETURNING *;

-- name: DeleteTVLibrary :exec
DELETE FROM tv_libraries WHERE id = $1;

-- name: UpdateTVLibraryScanStatus :exec
UPDATE tv_libraries
SET
    last_scan_at = $2,
    last_scan_duration = $3
WHERE id = $1;

-- name: CountSeriesByTVLibrary :one
SELECT COUNT(*) FROM series WHERE tv_library_id = $1;

-- name: GrantTVLibraryAccess :exec
INSERT INTO tv_library_access (library_id, user_id, can_manage)
VALUES ($1, $2, $3)
ON CONFLICT (library_id, user_id)
DO UPDATE SET can_manage = EXCLUDED.can_manage;

-- name: RevokeTVLibraryAccess :exec
DELETE FROM tv_library_access WHERE library_id = $1 AND user_id = $2;

-- name: ListTVLibraryAccess :many
SELECT * FROM tv_library_access WHERE library_id = $1;

-- name: HasTVLibraryAccess :one
SELECT EXISTS(
    SELECT 1 FROM tv_library_access
    WHERE library_id = $1 AND user_id = $2
) AS has_access;

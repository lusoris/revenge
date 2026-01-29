-- name: GetLibraryByID :one
SELECT * FROM libraries WHERE id = $1;

-- name: ListLibraries :many
SELECT * FROM libraries
ORDER BY sort_order ASC, name ASC;

-- name: ListLibrariesByType :many
SELECT * FROM libraries
WHERE type = $1
ORDER BY sort_order ASC, name ASC;

-- name: ListAccessibleLibraries :many
SELECT l.* FROM libraries l
LEFT JOIN library_user_access lua ON l.id = lua.library_id
WHERE l.is_private = false
   OR l.owner_user_id = $1
   OR lua.user_id = $1
ORDER BY l.sort_order ASC, l.name ASC;

-- name: CreateLibrary :one
INSERT INTO libraries (
    name, type, paths,
    scan_enabled, scan_interval_hours,
    preferred_language, download_images, download_nfo, generate_chapters,
    is_private, owner_user_id, sort_order, icon
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
) RETURNING *;

-- name: UpdateLibrary :one
UPDATE libraries SET
    name = COALESCE(sqlc.narg('name'), name),
    paths = COALESCE(sqlc.narg('paths'), paths),
    scan_enabled = COALESCE(sqlc.narg('scan_enabled'), scan_enabled),
    scan_interval_hours = COALESCE(sqlc.narg('scan_interval_hours'), scan_interval_hours),
    preferred_language = COALESCE(sqlc.narg('preferred_language'), preferred_language),
    download_images = COALESCE(sqlc.narg('download_images'), download_images),
    download_nfo = COALESCE(sqlc.narg('download_nfo'), download_nfo),
    generate_chapters = COALESCE(sqlc.narg('generate_chapters'), generate_chapters),
    is_private = COALESCE(sqlc.narg('is_private'), is_private),
    sort_order = COALESCE(sqlc.narg('sort_order'), sort_order),
    icon = COALESCE(sqlc.narg('icon'), icon)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateLibraryScanStatus :exec
UPDATE libraries SET
    last_scan_at = NOW(),
    last_scan_duration = $2
WHERE id = $1;

-- name: DeleteLibrary :exec
DELETE FROM libraries WHERE id = $1;

-- name: CountLibraries :one
SELECT COUNT(*) FROM libraries;

-- Library User Access
-- name: GrantLibraryAccess :exec
INSERT INTO library_user_access (library_id, user_id, can_manage)
VALUES ($1, $2, $3)
ON CONFLICT (library_id, user_id) DO UPDATE SET can_manage = $3;

-- name: RevokeLibraryAccess :exec
DELETE FROM library_user_access
WHERE library_id = $1 AND user_id = $2;

-- name: ListLibraryUsers :many
SELECT lua.*, u.username, u.email
FROM library_user_access lua
JOIN users u ON lua.user_id = u.id
WHERE lua.library_id = $1;

-- name: UserCanAccessLibrary :one
SELECT EXISTS(
    SELECT 1 FROM libraries l
    LEFT JOIN library_user_access lua ON l.id = lua.library_id
    WHERE l.id = $1
      AND (l.is_private = false OR l.owner_user_id = $2 OR lua.user_id = $2)
);

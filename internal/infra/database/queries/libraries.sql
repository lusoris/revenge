-- name: GetLibraryByID :one
SELECT id, name, type, paths, settings, is_visible, is_adult,
       scan_interval_hours, last_scan_at, created_at, updated_at
FROM libraries
WHERE id = $1;

-- name: GetLibraryByName :one
SELECT id, name, type, paths, settings, is_visible, is_adult,
       scan_interval_hours, last_scan_at, created_at, updated_at
FROM libraries
WHERE name = $1;

-- name: ListLibraries :many
SELECT id, name, type, paths, settings, is_visible, is_adult,
       scan_interval_hours, last_scan_at, created_at, updated_at
FROM libraries
ORDER BY name;

-- name: ListLibrariesByType :many
SELECT id, name, type, paths, settings, is_visible, is_adult,
       scan_interval_hours, last_scan_at, created_at, updated_at
FROM libraries
WHERE type = $1
ORDER BY name;

-- name: ListVisibleLibraries :many
SELECT id, name, type, paths, settings, is_visible, is_adult,
       scan_interval_hours, last_scan_at, created_at, updated_at
FROM libraries
WHERE is_visible = true
ORDER BY name;

-- name: ListNonAdultLibraries :many
SELECT id, name, type, paths, settings, is_visible, is_adult,
       scan_interval_hours, last_scan_at, created_at, updated_at
FROM libraries
WHERE is_adult = false
ORDER BY name;

-- name: ListLibrariesForUser :many
-- List libraries accessible to a user (considers adult content settings)
SELECT l.id, l.name, l.type, l.paths, l.settings, l.is_visible, l.is_adult,
       l.scan_interval_hours, l.last_scan_at, l.created_at, l.updated_at
FROM libraries l
JOIN users u ON u.id = $1
WHERE l.is_visible = true
  AND (NOT l.is_adult OR u.adult_content_enabled = true)
ORDER BY l.name;

-- name: CreateLibrary :one
INSERT INTO libraries (name, type, paths, settings, is_visible, is_adult, scan_interval_hours)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, name, type, paths, settings, is_visible, is_adult,
          scan_interval_hours, last_scan_at, created_at, updated_at;

-- name: UpdateLibrary :exec
UPDATE libraries
SET name = COALESCE(sqlc.narg('name'), name),
    paths = COALESCE(sqlc.narg('paths'), paths),
    settings = COALESCE(sqlc.narg('settings'), settings),
    is_visible = COALESCE(sqlc.narg('is_visible'), is_visible),
    scan_interval_hours = COALESCE(sqlc.narg('scan_interval_hours'), scan_interval_hours),
    updated_at = NOW()
WHERE id = @id;

-- name: DeleteLibrary :exec
DELETE FROM libraries WHERE id = $1;

-- name: UpdateLibraryLastScan :exec
UPDATE libraries
SET last_scan_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: CountLibraries :one
SELECT COUNT(*) FROM libraries;

-- name: LibraryNameExists :one
SELECT EXISTS (SELECT 1 FROM libraries WHERE name = $1) AS exists;

-- name: LibraryNameExistsExcluding :one
SELECT EXISTS (SELECT 1 FROM libraries WHERE name = $1 AND id != $2) AS exists;

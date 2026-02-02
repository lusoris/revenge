-- name: CreateLibrary :one
-- Creates a new library
INSERT INTO public.libraries (
    name,
    type,
    paths,
    enabled,
    scan_on_startup,
    realtime_monitoring,
    metadata_provider,
    preferred_language,
    scanner_config
) VALUES (
    @name,
    @type,
    @paths,
    @enabled,
    @scan_on_startup,
    @realtime_monitoring,
    @metadata_provider,
    @preferred_language,
    @scanner_config
)
RETURNING *;

-- name: GetLibrary :one
-- Gets a library by ID
SELECT * FROM public.libraries
WHERE id = @id;

-- name: GetLibraryByName :one
-- Gets a library by name
SELECT * FROM public.libraries
WHERE name = @name;

-- name: ListLibraries :many
-- Lists all libraries
SELECT * FROM public.libraries
ORDER BY name ASC;

-- name: ListEnabledLibraries :many
-- Lists all enabled libraries
SELECT * FROM public.libraries
WHERE enabled = true
ORDER BY name ASC;

-- name: ListLibrariesByType :many
-- Lists libraries by type
SELECT * FROM public.libraries
WHERE type = @type
ORDER BY name ASC;

-- name: UpdateLibrary :one
-- Updates a library
UPDATE public.libraries SET
    name = COALESCE(sqlc.narg('name'), name),
    type = COALESCE(sqlc.narg('type'), type),
    paths = COALESCE(sqlc.narg('paths'), paths),
    enabled = COALESCE(sqlc.narg('enabled'), enabled),
    scan_on_startup = COALESCE(sqlc.narg('scan_on_startup'), scan_on_startup),
    realtime_monitoring = COALESCE(sqlc.narg('realtime_monitoring'), realtime_monitoring),
    metadata_provider = COALESCE(sqlc.narg('metadata_provider'), metadata_provider),
    preferred_language = COALESCE(sqlc.narg('preferred_language'), preferred_language),
    scanner_config = COALESCE(sqlc.narg('scanner_config'), scanner_config)
WHERE id = @id
RETURNING *;

-- name: DeleteLibrary :exec
-- Deletes a library by ID
DELETE FROM public.libraries
WHERE id = @id;

-- name: CountLibraries :one
-- Counts total libraries
SELECT COUNT(*) FROM public.libraries;

-- name: CountLibrariesByType :one
-- Counts libraries by type
SELECT COUNT(*) FROM public.libraries
WHERE type = @type;

-- ============================================================================
-- Library Scans
-- ============================================================================

-- name: CreateLibraryScan :one
-- Creates a new library scan record
INSERT INTO public.library_scans (
    library_id,
    scan_type,
    status
) VALUES (
    @library_id,
    @scan_type,
    @status
)
RETURNING *;

-- name: GetLibraryScan :one
-- Gets a library scan by ID
SELECT * FROM public.library_scans
WHERE id = @id;

-- name: ListLibraryScans :many
-- Lists scans for a library
SELECT * FROM public.library_scans
WHERE library_id = @library_id
ORDER BY created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: CountLibraryScans :one
-- Counts scans for a library
SELECT COUNT(*) FROM public.library_scans
WHERE library_id = @library_id;

-- name: GetLatestLibraryScan :one
-- Gets the most recent scan for a library
SELECT * FROM public.library_scans
WHERE library_id = @library_id
ORDER BY created_at DESC
LIMIT 1;

-- name: GetRunningScans :many
-- Gets all currently running scans
SELECT * FROM public.library_scans
WHERE status = 'running'
ORDER BY started_at ASC;

-- name: UpdateLibraryScanStatus :one
-- Updates scan status
UPDATE public.library_scans SET
    status = @status,
    started_at = COALESCE(sqlc.narg('started_at'), started_at),
    completed_at = sqlc.narg('completed_at'),
    duration_seconds = sqlc.narg('duration_seconds'),
    error_message = sqlc.narg('error_message')
WHERE id = @id
RETURNING *;

-- name: UpdateLibraryScanProgress :one
-- Updates scan progress
UPDATE public.library_scans SET
    items_scanned = @items_scanned,
    items_added = @items_added,
    items_updated = @items_updated,
    items_removed = @items_removed,
    errors_count = @errors_count
WHERE id = @id
RETURNING *;

-- name: DeleteOldLibraryScans :execrows
-- Deletes library scans older than a given time
DELETE FROM public.library_scans
WHERE created_at < @older_than
AND status IN ('completed', 'failed', 'cancelled');

-- ============================================================================
-- Library Permissions
-- ============================================================================

-- name: CreateLibraryPermission :one
-- Grants a permission to a user for a library
INSERT INTO public.library_permissions (
    library_id,
    user_id,
    permission
) VALUES (
    @library_id,
    @user_id,
    @permission
)
ON CONFLICT (library_id, user_id, permission) DO NOTHING
RETURNING *;

-- name: GetLibraryPermission :one
-- Gets a specific permission
SELECT * FROM public.library_permissions
WHERE library_id = @library_id
AND user_id = @user_id
AND permission = @permission;

-- name: ListLibraryPermissions :many
-- Lists all permissions for a library
SELECT * FROM public.library_permissions
WHERE library_id = @library_id
ORDER BY created_at ASC;

-- name: ListUserLibraryPermissions :many
-- Lists all library permissions for a user
SELECT * FROM public.library_permissions
WHERE user_id = @user_id
ORDER BY created_at ASC;

-- name: CheckLibraryPermission :one
-- Checks if a user has a specific permission for a library
SELECT EXISTS(
    SELECT 1 FROM public.library_permissions
    WHERE library_id = @library_id
    AND user_id = @user_id
    AND permission = @permission
) AS has_permission;

-- name: GetUserAccessibleLibraries :many
-- Gets all libraries a user has view access to
SELECT DISTINCT l.* FROM public.libraries l
JOIN public.library_permissions lp ON l.id = lp.library_id
WHERE lp.user_id = @user_id
AND lp.permission = 'view'
AND l.enabled = true
ORDER BY l.name ASC;

-- name: DeleteLibraryPermission :exec
-- Revokes a permission from a user for a library
DELETE FROM public.library_permissions
WHERE library_id = @library_id
AND user_id = @user_id
AND permission = @permission;

-- name: DeleteAllLibraryPermissions :exec
-- Revokes all permissions for a library (used when deleting library)
DELETE FROM public.library_permissions
WHERE library_id = @library_id;

-- name: DeleteUserLibraryPermissions :exec
-- Revokes all library permissions for a user
DELETE FROM public.library_permissions
WHERE user_id = @user_id;

-- name: CountLibraryPermissions :one
-- Counts permissions for a library
SELECT COUNT(*) FROM public.library_permissions
WHERE library_id = @library_id;

-- Migrate from shared libraries to per-module tv_libraries
-- Phase 2.5.2: Update Foreign Keys
BEGIN;

-- Step 1: Create default TV library for migration if none exists
INSERT INTO tv_libraries (id, name, paths, scan_enabled)
SELECT
    gen_random_uuid(),
    'Default TV Library',
    ARRAY['/tv'],
    true
WHERE NOT EXISTS (SELECT 1 FROM tv_libraries LIMIT 1);

-- Step 2: Migrate series from old library_id to tv_library_id
-- Creates a corresponding tv_library for each old library entry
INSERT INTO tv_libraries (id, name, paths, scan_enabled, preferred_language, is_private, owner_user_id, sort_order, created_at, updated_at)
SELECT
    l.id,  -- Use same UUID for easy mapping
    l.name,
    l.paths,
    l.scan_enabled,
    l.preferred_language,
    l.is_private,
    l.owner_user_id,
    l.sort_order,
    l.created_at,
    l.updated_at
FROM libraries l
WHERE l.type = 'tvshow'
AND NOT EXISTS (SELECT 1 FROM tv_libraries tl WHERE tl.id = l.id)
ON CONFLICT (id) DO NOTHING;

-- Step 3: Update series to use the new tv_library_id
UPDATE series
SET tv_library_id = library_id
WHERE tv_library_id IS NULL
AND library_id IS NOT NULL
AND EXISTS (SELECT 1 FROM tv_libraries tl WHERE tl.id = series.library_id);

-- Step 4: For any remaining series without tv_library_id, assign to default library
UPDATE series
SET tv_library_id = (SELECT id FROM tv_libraries ORDER BY created_at LIMIT 1)
WHERE tv_library_id IS NULL;

-- Step 5: Make tv_library_id NOT NULL now that all data is migrated
ALTER TABLE series ALTER COLUMN tv_library_id SET NOT NULL;

-- Step 6: Drop the old library_id column and its index
DROP INDEX IF EXISTS idx_series_library;
ALTER TABLE series DROP COLUMN library_id;

-- Step 7: Rename the new index for clarity
DROP INDEX IF EXISTS idx_series_tv_library;
CREATE INDEX idx_series_library ON series(tv_library_id);

COMMIT;

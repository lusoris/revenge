-- Migrate from shared libraries to per-module movie_libraries
-- Phase 2.5.2: Update Foreign Keys
BEGIN;

-- Step 1: Create default movie library for migration if none exists
INSERT INTO movie_libraries (id, name, paths, scan_enabled)
SELECT
    gen_random_uuid(),
    'Default Movie Library',
    ARRAY['/movies'],
    true
WHERE NOT EXISTS (SELECT 1 FROM movie_libraries LIMIT 1);

-- Step 2: Migrate movies from old library_id to movie_library_id
-- Creates a corresponding movie_library for each old library entry
INSERT INTO movie_libraries (id, name, paths, scan_enabled, preferred_language, is_private, owner_user_id, sort_order, created_at, updated_at)
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
WHERE l.type = 'movie'
AND NOT EXISTS (SELECT 1 FROM movie_libraries ml WHERE ml.id = l.id)
ON CONFLICT (id) DO NOTHING;

-- Step 3: Update movies to use the new movie_library_id
UPDATE movies
SET movie_library_id = library_id
WHERE movie_library_id IS NULL
AND library_id IS NOT NULL
AND EXISTS (SELECT 1 FROM movie_libraries ml WHERE ml.id = movies.library_id);

-- Step 4: For any remaining movies without movie_library_id, assign to default library
UPDATE movies
SET movie_library_id = (SELECT id FROM movie_libraries ORDER BY created_at LIMIT 1)
WHERE movie_library_id IS NULL;

-- Step 5: Make movie_library_id NOT NULL now that all data is migrated
ALTER TABLE movies ALTER COLUMN movie_library_id SET NOT NULL;

-- Step 6: Drop the old library_id column and its index
DROP INDEX IF EXISTS idx_movies_library;
ALTER TABLE movies DROP COLUMN library_id;

-- Step 7: Rename the new index for clarity
DROP INDEX IF EXISTS idx_movies_movie_library;
CREATE INDEX idx_movies_library ON movies(movie_library_id);

COMMIT;

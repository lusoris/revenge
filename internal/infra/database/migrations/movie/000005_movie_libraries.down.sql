-- Rollback movie_libraries migration
BEGIN;

-- Remove new FK column from movies
ALTER TABLE movies DROP COLUMN IF EXISTS movie_library_id;

-- Drop user access table
DROP TABLE IF EXISTS movie_library_access;

-- Drop movie_libraries table
DROP TABLE IF EXISTS movie_libraries;

COMMIT;

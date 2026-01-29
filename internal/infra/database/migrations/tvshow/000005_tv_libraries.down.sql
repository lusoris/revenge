-- Rollback tv_libraries migration
BEGIN;

-- Remove new FK column from series
ALTER TABLE series DROP COLUMN IF EXISTS tv_library_id;

-- Drop user access table
DROP TABLE IF EXISTS tv_library_access;

-- Drop tv_libraries table
DROP TABLE IF EXISTS tv_libraries;

COMMIT;

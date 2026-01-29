-- Deprecate shared libraries table
-- Phase 2.5.3: Mark shared library infrastructure as deprecated
--
-- The shared 'libraries' table is being replaced by per-module library tables:
-- - movie_libraries (for movies)
-- - tv_libraries (for TV shows)
-- - qar.fleets (for adult content)
-- - music_libraries (planned)
-- - etc.
--
-- This migration does NOT drop the table yet - it adds deprecation notices
-- and removes any remaining dependencies where possible.
BEGIN;

-- Add deprecation comment to the libraries table
COMMENT ON TABLE libraries IS 'DEPRECATED: Use per-module library tables instead (movie_libraries, tv_libraries, qar.fleets, etc.). This table will be removed in a future version.';

-- Add deprecation comment to the library_user_access table
COMMENT ON TABLE library_user_access IS 'DEPRECATED: Use per-module library access tables instead (movie_library_access, tv_library_access, etc.). This table will be removed in a future version.';

-- Add deprecation comment to the enum
COMMENT ON TYPE library_type IS 'DEPRECATED: Library types are now implicit in per-module tables. This enum will be removed in a future version.';

-- NOTE: A unified all_libraries view is not created here because sqlc processes
-- schemas separately. The view would require cross-module table access.
-- If needed, implement library aggregation at the application level.

COMMIT;

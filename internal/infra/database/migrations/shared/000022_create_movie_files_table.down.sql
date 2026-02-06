-- Drop movie_files table
DROP TRIGGER IF EXISTS update_movie_files_updated_at ON public.movie_files;
DROP INDEX IF EXISTS idx_movie_files_last_scanned;
DROP INDEX IF EXISTS idx_movie_files_quality;
DROP INDEX IF EXISTS idx_movie_files_resolution;
DROP INDEX IF EXISTS idx_movie_files_radarr_id;
DROP INDEX IF EXISTS idx_movie_files_movie_id;
DROP TABLE IF EXISTS public.movie_files;

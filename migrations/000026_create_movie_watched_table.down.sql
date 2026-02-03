-- Drop movie_watched table
DROP TRIGGER IF EXISTS update_movie_watched_updated_at ON public.movie_watched;
DROP INDEX IF EXISTS idx_movie_watched_in_progress;
DROP INDEX IF EXISTS idx_movie_watched_completed;
DROP INDEX IF EXISTS idx_movie_watched_last_watched;
DROP INDEX IF EXISTS idx_movie_watched_movie_id;
DROP INDEX IF EXISTS idx_movie_watched_user_id;
DROP TABLE IF EXISTS public.movie_watched;

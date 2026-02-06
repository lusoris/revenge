-- Drop movies table
DROP TRIGGER IF EXISTS update_movies_updated_at ON public.movies;
DROP INDEX IF EXISTS idx_movies_original_title_trgm;
DROP INDEX IF EXISTS idx_movies_title_trgm;
DROP INDEX IF EXISTS idx_movies_vote_average;
DROP INDEX IF EXISTS idx_movies_library_added_at;
DROP INDEX IF EXISTS idx_movies_release_date;
DROP INDEX IF EXISTS idx_movies_year;
DROP INDEX IF EXISTS idx_movies_radarr_id;
DROP INDEX IF EXISTS idx_movies_imdb_id;
DROP INDEX IF EXISTS idx_movies_tmdb_id;
DROP TABLE IF EXISTS public.movies;

-- Drop movie_credits table
DROP TRIGGER IF EXISTS update_movie_credits_updated_at ON public.movie_credits;
DROP INDEX IF EXISTS idx_movie_credits_unique;
DROP INDEX IF EXISTS idx_movie_credits_job;
DROP INDEX IF EXISTS idx_movie_credits_cast_order;
DROP INDEX IF EXISTS idx_movie_credits_type;
DROP INDEX IF EXISTS idx_movie_credits_person_id;
DROP INDEX IF EXISTS idx_movie_credits_movie_id;
DROP TABLE IF EXISTS public.movie_credits;

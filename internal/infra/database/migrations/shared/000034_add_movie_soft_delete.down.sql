-- Remove soft-delete columns from movie tables

DROP INDEX IF EXISTS idx_movie_collections_deleted_at;
DROP INDEX IF EXISTS idx_movie_credits_deleted_at;
DROP INDEX IF EXISTS idx_movie_files_deleted_at;
DROP INDEX IF EXISTS idx_movies_deleted_at;

ALTER TABLE public.movie_collections DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE public.movie_credits DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE public.movie_files DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE public.movies DROP COLUMN IF EXISTS deleted_at;

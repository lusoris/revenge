-- Add soft-delete support to movie tables.
-- Queries already filter by deleted_at IS NULL; this migration adds the column.

ALTER TABLE public.movies ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;
ALTER TABLE public.movie_files ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;
ALTER TABLE public.movie_credits ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;
ALTER TABLE public.movie_collections ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

-- Partial indexes for efficient filtering of non-deleted rows
CREATE INDEX IF NOT EXISTS idx_movies_deleted_at ON public.movies(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_movie_files_deleted_at ON public.movie_files(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_movie_credits_deleted_at ON public.movie_credits(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_movie_collections_deleted_at ON public.movie_collections(deleted_at) WHERE deleted_at IS NULL;

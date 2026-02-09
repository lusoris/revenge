-- Remove external_ratings JSONB column from movies and series tables
ALTER TABLE public.movies DROP COLUMN IF EXISTS external_ratings;

ALTER TABLE tvshow.series DROP COLUMN IF EXISTS external_ratings;

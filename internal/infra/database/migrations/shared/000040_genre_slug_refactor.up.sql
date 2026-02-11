-- Migration 000040: Replace tmdb_genre_id with slug in genre tables
-- This decouples the genre system from TMDb-specific integer IDs,
-- using a canonical slug as the universal deduplication key.

-- ============================================================
-- Movie Genres
-- ============================================================

-- Add slug column
ALTER TABLE movie.movie_genres ADD COLUMN slug TEXT;

-- Populate slug from existing name (lowercase, spaces/special chars → hyphens)
UPDATE movie.movie_genres
SET slug = LOWER(REGEXP_REPLACE(TRIM(name), '[^a-zA-Z0-9]+', '-', 'g'));

-- Remove trailing hyphens
UPDATE movie.movie_genres
SET slug = REGEXP_REPLACE(slug, '-+$', '');

-- Make slug NOT NULL
ALTER TABLE movie.movie_genres ALTER COLUMN slug SET NOT NULL;

-- Drop old unique constraint and index
ALTER TABLE movie.movie_genres DROP CONSTRAINT IF EXISTS uq_movie_genre;
DROP INDEX IF EXISTS movie.idx_movie_genres_genre_id;

-- Drop tmdb_genre_id column
ALTER TABLE movie.movie_genres DROP COLUMN tmdb_genre_id;

-- Add new unique constraint and index on slug
ALTER TABLE movie.movie_genres ADD CONSTRAINT uq_movie_genre UNIQUE(movie_id, slug);
CREATE INDEX idx_movie_genres_slug ON movie.movie_genres(slug);

-- ============================================================
-- Series Genres
-- ============================================================

-- Add slug column
ALTER TABLE tvshow.series_genres ADD COLUMN slug TEXT;

-- Populate slug from existing name
UPDATE tvshow.series_genres
SET slug = LOWER(REGEXP_REPLACE(TRIM(name), '[^a-zA-Z0-9]+', '-', 'g'));

-- Remove trailing hyphens
UPDATE tvshow.series_genres
SET slug = REGEXP_REPLACE(slug, '-+$', '');

-- Make slug NOT NULL
ALTER TABLE tvshow.series_genres ALTER COLUMN slug SET NOT NULL;

-- Drop old unique constraint and index
ALTER TABLE tvshow.series_genres DROP CONSTRAINT IF EXISTS series_genres_series_id_tmdb_genre_id_key;
DROP INDEX IF EXISTS tvshow.idx_series_genres_genre;

-- Drop tmdb_genre_id column
ALTER TABLE tvshow.series_genres DROP COLUMN tmdb_genre_id;

-- Add new unique constraint and index on slug
ALTER TABLE tvshow.series_genres ADD CONSTRAINT uq_series_genre UNIQUE(series_id, slug);
CREATE INDEX idx_series_genres_slug ON tvshow.series_genres(slug);

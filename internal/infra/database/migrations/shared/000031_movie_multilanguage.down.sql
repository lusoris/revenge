-- Rollback: Remove multi-language support from movies
-- Phase: A9.1.1
-- Date: 2026-02-05
--
-- This migration rolls back multi-language support, removing JSONB columns.
-- The original title, tagline, and overview columns remain unchanged.

-- Drop indexes
DROP INDEX IF EXISTS idx_movies_titles_i18n;
DROP INDEX IF EXISTS idx_movies_taglines_i18n;
DROP INDEX IF EXISTS idx_movies_overviews_i18n;
DROP INDEX IF EXISTS idx_movies_age_ratings;
DROP INDEX IF EXISTS idx_movies_original_language;

-- Drop columns (except original_language which existed before this migration)
ALTER TABLE movies DROP COLUMN IF EXISTS titles_i18n;
ALTER TABLE movies DROP COLUMN IF EXISTS taglines_i18n;
ALTER TABLE movies DROP COLUMN IF EXISTS overviews_i18n;
ALTER TABLE movies DROP COLUMN IF EXISTS age_ratings;
-- Note: original_language is NOT dropped as it existed in the original table schema (migration 000021)

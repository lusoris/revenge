-- Migration: Add multi-language support to movies
-- Phase: A9.1.1
-- Date: 2026-02-05
--
-- This migration adds JSONB columns for multi-language metadata support.
-- Enables storing titles, taglines, overviews, and age ratings in multiple languages.
--
-- Key features:
-- - Hybrid approach: Keep default fields + add i18n JSONB for performance
-- - Support for 40+ languages from TMDb/TheTVDB
-- - Age ratings per country (MPAA, FSK, BBFC, PEGI, etc.)
-- - Original language tracking
-- - Migrates existing English data to new structure

-- Add multi-language JSONB columns
ALTER TABLE movies ADD COLUMN IF NOT EXISTS titles_i18n JSONB DEFAULT '{}';
ALTER TABLE movies ADD COLUMN IF NOT EXISTS taglines_i18n JSONB DEFAULT '{}';
ALTER TABLE movies ADD COLUMN IF NOT EXISTS overviews_i18n JSONB DEFAULT '{}';
ALTER TABLE movies ADD COLUMN IF NOT EXISTS age_ratings JSONB DEFAULT '{}';

-- Add original language tracking (ISO 639-1 codes: en, de, fr, es, ja, etc.)
ALTER TABLE movies ADD COLUMN IF NOT EXISTS original_language TEXT DEFAULT 'en';

-- Migrate existing data to default language (en)
-- Only migrate rows that haven't been migrated yet (titles_i18n is empty)
UPDATE movies
SET
    original_language = COALESCE(original_language, 'en'),
    titles_i18n = jsonb_build_object('en', title),
    taglines_i18n = CASE
        WHEN tagline IS NOT NULL AND tagline != '' THEN jsonb_build_object('en', tagline)
        ELSE '{}'::jsonb
    END,
    overviews_i18n = CASE
        WHEN overview IS NOT NULL AND overview != '' THEN jsonb_build_object('en', overview)
        ELSE '{}'::jsonb
    END
WHERE titles_i18n = '{}'::jsonb;

-- Create GIN indexes for JSONB lookups (efficient for key existence and text search)
CREATE INDEX IF NOT EXISTS idx_movies_titles_i18n ON movies USING GIN (titles_i18n);
CREATE INDEX IF NOT EXISTS idx_movies_taglines_i18n ON movies USING GIN (taglines_i18n);
CREATE INDEX IF NOT EXISTS idx_movies_overviews_i18n ON movies USING GIN (overviews_i18n);
CREATE INDEX IF NOT EXISTS idx_movies_age_ratings ON movies USING GIN (age_ratings);

-- Create B-tree index for original_language (for filtering by language)
CREATE INDEX IF NOT EXISTS idx_movies_original_language ON movies(original_language);

-- Add column comments for documentation
COMMENT ON COLUMN movies.titles_i18n IS 'Movie titles by ISO 639-1 language code: {"en": "The Shawshank Redemption", "de": "Die Verurteilten", "fr": "Les Évadés"}';
COMMENT ON COLUMN movies.taglines_i18n IS 'Taglines by language code: {"en": "Fear can hold you prisoner. Hope can set you free.", "de": "Angst kann dich gefangen halten. Hoffnung kann dich befreien."}';
COMMENT ON COLUMN movies.overviews_i18n IS 'Plot overviews by language code: {"en": "Imprisoned in the 1940s...", "de": "In den 1940er Jahren eingesperrt..."}';
COMMENT ON COLUMN movies.age_ratings IS 'Age ratings by country code and rating system: {"US": {"MPAA": "R"}, "DE": {"FSK": "12"}, "GB": {"BBFC": "15"}}';
COMMENT ON COLUMN movies.original_language IS 'ISO 639-1 language code of the original movie language (en, de, fr, es, ja, ko, etc.)';

-- Add metadata_language column to user_preferences
-- Controls which language is used for displaying titles, overviews, and taglines.
-- ISO 639-1 code (e.g., "en", "de", "fr"). Default: "en" (British English).

ALTER TABLE shared.user_preferences
ADD COLUMN metadata_language VARCHAR(10) DEFAULT 'en';

COMMENT ON COLUMN shared.user_preferences.metadata_language
IS 'Language for metadata display (titles, overviews, taglines). ISO 639-1 code.';

-- 000015_genres.down.sql
-- Revert genre domain separation

BEGIN;

-- Drop trigger first
DROP TRIGGER IF EXISTS genres_updated_at ON genres;
DROP FUNCTION IF EXISTS update_genre_timestamp();

-- Drop tables (junction first due to FK)
DROP TABLE IF EXISTS media_item_genres;
DROP TABLE IF EXISTS genres;

-- Drop enum
DROP TYPE IF EXISTS genre_domain;

COMMIT;

BEGIN;

DROP TABLE IF EXISTS movie_studio_link;
DROP TABLE IF EXISTS movie_studios;
DROP INDEX IF EXISTS idx_movies_collection;
ALTER TABLE movies DROP COLUMN IF EXISTS collection_order;
ALTER TABLE movies DROP COLUMN IF EXISTS collection_id;
DROP TABLE IF EXISTS movie_collections;
DROP TABLE IF EXISTS movies;

COMMIT;

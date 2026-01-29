-- 000009_content_ratings.down.sql

BEGIN;

DROP TABLE IF EXISTS content_ratings CASCADE;
DROP TYPE IF EXISTS rating_system;

COMMIT;

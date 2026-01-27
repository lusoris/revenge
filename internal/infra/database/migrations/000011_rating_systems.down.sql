-- 000011_rating_systems.down.sql

DROP TRIGGER IF EXISTS trigger_refresh_content_min_rating_levels ON content_ratings;
DROP FUNCTION IF EXISTS refresh_content_min_rating_levels();
DROP MATERIALIZED VIEW IF EXISTS content_min_rating_levels;
DROP TABLE IF EXISTS content_ratings;
DROP TABLE IF EXISTS rating_equivalents;
DROP TABLE IF EXISTS ratings;
DROP TABLE IF EXISTS rating_systems;

-- 000014_seed_rating_systems.down.sql
-- Remove seeded rating systems data

-- Delete in reverse order of foreign key dependencies
DELETE FROM rating_equivalents;
DELETE FROM ratings;
DELETE FROM rating_systems;

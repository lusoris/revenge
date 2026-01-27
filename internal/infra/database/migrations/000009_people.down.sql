-- 000009_people.down.sql
DROP TABLE IF EXISTS media_people;
DROP TRIGGER IF EXISTS update_people_updated_at ON people;
DROP TABLE IF EXISTS people;
DROP TYPE IF EXISTS person_type;

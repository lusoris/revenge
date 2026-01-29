BEGIN;

DROP TABLE IF EXISTS movie_credits;
DROP TYPE IF EXISTS movie_credit_role;
DROP TRIGGER IF EXISTS movie_people_updated_at ON movie_people;
DROP TABLE IF EXISTS movie_people;

COMMIT;

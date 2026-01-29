-- TV Show Credits: Rollback
BEGIN;

DROP TABLE IF EXISTS episode_credits;
DROP TABLE IF EXISTS series_credits;
DROP TYPE IF EXISTS tvshow_credit_role;
DROP TABLE IF EXISTS tvshow_people;

COMMIT;

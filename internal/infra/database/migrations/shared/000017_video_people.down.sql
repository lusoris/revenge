-- Video People: Rollback
BEGIN;

DROP TYPE IF EXISTS video_credit_role;
DROP TABLE IF EXISTS video_people;

COMMIT;

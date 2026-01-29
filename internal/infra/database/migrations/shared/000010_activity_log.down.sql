-- 000008_activity_log.down.sql

BEGIN;

DROP TABLE IF EXISTS activity_log CASCADE;
DROP TYPE IF EXISTS activity_type;
DROP TYPE IF EXISTS activity_severity;

COMMIT;

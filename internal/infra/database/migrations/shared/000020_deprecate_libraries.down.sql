-- Rollback: Remove deprecation notices
BEGIN;

-- Remove deprecation comments
COMMENT ON TABLE libraries IS NULL;
COMMENT ON TABLE library_user_access IS NULL;
COMMENT ON TYPE library_type IS NULL;

COMMIT;

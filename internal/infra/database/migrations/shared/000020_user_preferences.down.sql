-- Rollback user preferences columns
BEGIN;

ALTER TABLE profiles
    DROP CONSTRAINT IF EXISTS check_auto_play_delay,
    DROP CONSTRAINT IF EXISTS check_continue_watching_days,
    DROP CONSTRAINT IF EXISTS check_mark_watched_percent;

ALTER TABLE profiles
    DROP COLUMN IF EXISTS auto_play_enabled,
    DROP COLUMN IF EXISTS auto_play_delay_seconds,
    DROP COLUMN IF EXISTS continue_watching_days,
    DROP COLUMN IF EXISTS mark_watched_percent,
    DROP COLUMN IF EXISTS adult_pin_hash;

COMMIT;

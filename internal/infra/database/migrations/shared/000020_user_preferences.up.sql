-- User Preferences: Playback and content settings
BEGIN;

-- Add preference columns to profiles table
ALTER TABLE profiles
    ADD COLUMN IF NOT EXISTS auto_play_enabled BOOLEAN NOT NULL DEFAULT true,
    ADD COLUMN IF NOT EXISTS auto_play_delay_seconds INT NOT NULL DEFAULT 10,
    ADD COLUMN IF NOT EXISTS continue_watching_days INT NOT NULL DEFAULT 30,
    ADD COLUMN IF NOT EXISTS mark_watched_percent INT NOT NULL DEFAULT 90,
    ADD COLUMN IF NOT EXISTS adult_pin_hash TEXT;

-- Constraints
ALTER TABLE profiles
    ADD CONSTRAINT check_auto_play_delay CHECK (auto_play_delay_seconds >= 0 AND auto_play_delay_seconds <= 60),
    ADD CONSTRAINT check_continue_watching_days CHECK (continue_watching_days >= 1 AND continue_watching_days <= 365),
    ADD CONSTRAINT check_mark_watched_percent CHECK (mark_watched_percent >= 50 AND mark_watched_percent <= 100);

COMMENT ON COLUMN profiles.auto_play_enabled IS 'Whether to auto-play next episode/item';
COMMENT ON COLUMN profiles.auto_play_delay_seconds IS 'Countdown seconds before auto-play (0-60)';
COMMENT ON COLUMN profiles.continue_watching_days IS 'Days to keep items in continue watching (1-365)';
COMMENT ON COLUMN profiles.mark_watched_percent IS 'Percentage of media watched to mark as complete (50-100)';
COMMENT ON COLUMN profiles.adult_pin_hash IS 'Hashed PIN for adult content access (optional)';

COMMIT;

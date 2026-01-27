-- 000012_user_rating_settings.up.sql
-- User extensions for content rating and age restriction

ALTER TABLE users ADD COLUMN birthdate DATE;
ALTER TABLE users ADD COLUMN max_rating_level INT NOT NULL DEFAULT 100;
ALTER TABLE users ADD COLUMN adult_content_enabled BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE users ADD COLUMN preferred_rating_system VARCHAR(20);  -- 'fsk', 'mpaa', NULL for auto
ALTER TABLE users ADD COLUMN parental_pin_hash VARCHAR(255);       -- PIN for unlocking restricted content
ALTER TABLE users ADD COLUMN hide_restricted BOOLEAN NOT NULL DEFAULT true;  -- Hide vs show locked items

-- Constraint to ensure max_rating_level is valid
ALTER TABLE users ADD CONSTRAINT users_max_rating_level_check
    CHECK (max_rating_level >= 0 AND max_rating_level <= 100);

-- Index for filtering users by adult content setting
CREATE INDEX idx_users_adult_content ON users(adult_content_enabled) WHERE adult_content_enabled = true;

-- Comment on columns
COMMENT ON COLUMN users.birthdate IS 'User birthdate for age-based content filtering';
COMMENT ON COLUMN users.max_rating_level IS 'Maximum normalized rating level (0-100) user can view. Parental override.';
COMMENT ON COLUMN users.adult_content_enabled IS 'Whether user can access adult (XXX) content libraries';
COMMENT ON COLUMN users.preferred_rating_system IS 'Preferred rating system code for display (fsk, mpaa, bbfc)';
COMMENT ON COLUMN users.parental_pin_hash IS 'Hashed PIN for temporary unlock of restricted content';
COMMENT ON COLUMN users.hide_restricted IS 'If true, hide restricted content entirely. If false, show as locked.';

-- 000012_user_rating_settings.down.sql

ALTER TABLE users DROP CONSTRAINT IF EXISTS users_max_rating_level_check;
DROP INDEX IF EXISTS idx_users_adult_content;

ALTER TABLE users DROP COLUMN IF EXISTS hide_restricted;
ALTER TABLE users DROP COLUMN IF EXISTS parental_pin_hash;
ALTER TABLE users DROP COLUMN IF EXISTS preferred_rating_system;
ALTER TABLE users DROP COLUMN IF EXISTS adult_content_enabled;
ALTER TABLE users DROP COLUMN IF EXISTS max_rating_level;
ALTER TABLE users DROP COLUMN IF EXISTS birthdate;

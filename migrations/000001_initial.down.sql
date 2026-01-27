-- 001_initial.down.sql
-- Rollback initial schema

-- Drop triggers first
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_oidc_providers_updated_at ON oidc_providers;
DROP TRIGGER IF EXISTS update_libraries_updated_at ON libraries;
DROP TRIGGER IF EXISTS update_media_items_updated_at ON media_items;
DROP TRIGGER IF EXISTS update_playback_progress_updated_at ON playback_progress;
DROP TRIGGER IF EXISTS update_people_updated_at ON people;

DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS activity_log;
DROP TABLE IF EXISTS media_people;
DROP TABLE IF EXISTS people;
DROP TABLE IF EXISTS playback_progress;
DROP TABLE IF EXISTS images;
DROP TABLE IF EXISTS media_items;
DROP TABLE IF EXISTS libraries;
DROP TABLE IF EXISTS oidc_user_links;
DROP TABLE IF EXISTS oidc_providers;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;

-- Drop types
DROP TYPE IF EXISTS activity_type;
DROP TYPE IF EXISTS person_type;
DROP TYPE IF EXISTS image_type;
DROP TYPE IF EXISTS media_type;
DROP TYPE IF EXISTS library_type;

-- Note: We don't drop the uuid-ossp and pgcrypto extensions
-- as they may be used by other databases

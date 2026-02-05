-- Drop triggers first
DROP TRIGGER IF EXISTS update_episode_credits_updated_at ON tvshow.episode_credits;
DROP TRIGGER IF EXISTS update_series_credits_updated_at ON tvshow.series_credits;
DROP TRIGGER IF EXISTS update_episode_watched_updated_at ON tvshow.episode_watched;
DROP TRIGGER IF EXISTS update_episode_files_updated_at ON tvshow.episode_files;
DROP TRIGGER IF EXISTS update_episodes_updated_at ON tvshow.episodes;
DROP TRIGGER IF EXISTS update_seasons_updated_at ON tvshow.seasons;
DROP TRIGGER IF EXISTS update_series_updated_at ON tvshow.series;

-- Drop trigger function
DROP FUNCTION IF EXISTS tvshow.update_updated_at_column();

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS tvshow.series_networks;
DROP TABLE IF EXISTS tvshow.networks;
DROP TABLE IF EXISTS tvshow.episode_credits;
DROP TABLE IF EXISTS tvshow.series_credits;
DROP TABLE IF EXISTS tvshow.series_genres;
DROP TABLE IF EXISTS tvshow.episode_watched;
DROP TABLE IF EXISTS tvshow.episode_files;
DROP TABLE IF EXISTS tvshow.episodes;
DROP TABLE IF EXISTS tvshow.seasons;
DROP TABLE IF EXISTS tvshow.series;

-- Drop schema
DROP SCHEMA IF EXISTS tvshow;

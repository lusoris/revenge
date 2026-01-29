-- TV Show User Data: Rollback
BEGIN;

DROP TRIGGER IF EXISTS episode_watch_history_update_progress ON episode_watch_history;
DROP FUNCTION IF EXISTS update_series_watch_progress();

DROP TABLE IF EXISTS series_external_ratings;
DROP TABLE IF EXISTS series_watch_progress;
DROP TABLE IF EXISTS series_watchlist;
DROP TABLE IF EXISTS episode_watch_history;
DROP TABLE IF EXISTS series_favorites;
DROP TABLE IF EXISTS episode_user_ratings;
DROP TABLE IF EXISTS series_user_ratings;

COMMIT;

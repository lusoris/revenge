-- 000008_playback_progress.down.sql
DROP TRIGGER IF EXISTS update_playback_progress_updated_at ON playback_progress;
DROP TABLE IF EXISTS playback_progress;

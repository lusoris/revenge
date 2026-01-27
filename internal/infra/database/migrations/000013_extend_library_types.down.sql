-- 000013_extend_library_types.down.sql
-- Note: PostgreSQL does not support removing enum values
-- The enum values will remain, but we can remove the is_adult column

DROP INDEX IF EXISTS idx_libraries_adult;
ALTER TABLE libraries DROP COLUMN IF EXISTS is_adult;

-- IMPORTANT: Enum values cannot be removed in PostgreSQL
-- The following types will remain in the enum even after rollback:
-- library_type: musicvideos, homevideos, boxsets, livetv, playlists, books, audiobooks, podcasts, adult_movies, adult_shows
-- media_type: musicvideo, trailer, homevideo, audiobook_chapter, podcast_episode, book, audiobook, podcast, boxset, playlist, channel, program, recording

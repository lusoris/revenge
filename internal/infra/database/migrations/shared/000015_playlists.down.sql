-- Drop playlists tables and types

DROP TRIGGER IF EXISTS trg_playlist_items_update_metadata ON playlist_items;
DROP FUNCTION IF EXISTS update_playlist_metadata();

DROP TABLE IF EXISTS playlist_collaborators;
DROP TABLE IF EXISTS playlist_items;
DROP TABLE IF EXISTS playlists;

DROP TYPE IF EXISTS playlist_visibility;
DROP TYPE IF EXISTS playlist_type;

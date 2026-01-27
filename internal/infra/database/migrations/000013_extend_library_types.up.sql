-- 000013_extend_library_types.up.sql
-- Extend library_type and media_type enums with additional types

-- Add new library types
-- Note: PostgreSQL doesn't allow removing enum values, only adding
ALTER TYPE library_type ADD VALUE IF NOT EXISTS 'musicvideos';
ALTER TYPE library_type ADD VALUE IF NOT EXISTS 'homevideos';
ALTER TYPE library_type ADD VALUE IF NOT EXISTS 'boxsets';
ALTER TYPE library_type ADD VALUE IF NOT EXISTS 'livetv';
ALTER TYPE library_type ADD VALUE IF NOT EXISTS 'playlists';
ALTER TYPE library_type ADD VALUE IF NOT EXISTS 'books';
ALTER TYPE library_type ADD VALUE IF NOT EXISTS 'audiobooks';
ALTER TYPE library_type ADD VALUE IF NOT EXISTS 'podcasts';
ALTER TYPE library_type ADD VALUE IF NOT EXISTS 'adult_movies';
ALTER TYPE library_type ADD VALUE IF NOT EXISTS 'adult_shows';

-- Add new media types
ALTER TYPE media_type ADD VALUE IF NOT EXISTS 'musicvideo';
ALTER TYPE media_type ADD VALUE IF NOT EXISTS 'trailer';
ALTER TYPE media_type ADD VALUE IF NOT EXISTS 'homevideo';
ALTER TYPE media_type ADD VALUE IF NOT EXISTS 'audiobook_chapter';
ALTER TYPE media_type ADD VALUE IF NOT EXISTS 'podcast_episode';
ALTER TYPE media_type ADD VALUE IF NOT EXISTS 'book';
ALTER TYPE media_type ADD VALUE IF NOT EXISTS 'audiobook';
ALTER TYPE media_type ADD VALUE IF NOT EXISTS 'podcast';
ALTER TYPE media_type ADD VALUE IF NOT EXISTS 'boxset';
ALTER TYPE media_type ADD VALUE IF NOT EXISTS 'playlist';
ALTER TYPE media_type ADD VALUE IF NOT EXISTS 'channel';
ALTER TYPE media_type ADD VALUE IF NOT EXISTS 'program';
ALTER TYPE media_type ADD VALUE IF NOT EXISTS 'recording';

-- Add is_adult flag to libraries for quick filtering
ALTER TABLE libraries ADD COLUMN is_adult BOOLEAN NOT NULL DEFAULT false;

-- Index for adult library filtering
CREATE INDEX idx_libraries_adult ON libraries(is_adult) WHERE is_adult = true;

-- Update existing adult libraries (if any exist with the new types)
-- This will be a no-op if no adult libraries exist yet
UPDATE libraries SET is_adult = true
WHERE type IN ('adult_movies', 'adult_shows');

COMMENT ON COLUMN libraries.is_adult IS 'Whether this library contains adult (XXX) content. Auto-set for adult_movies/adult_shows types.';

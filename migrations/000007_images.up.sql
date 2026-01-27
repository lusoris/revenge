-- 000007_images.up.sql
-- Images for media items (posters, backdrops, etc.)

CREATE TYPE image_type AS ENUM (
    'primary', 'backdrop', 'logo', 'thumb', 'banner',
    'art', 'disc', 'box', 'screenshot', 'chapter'
);

CREATE TABLE images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID NOT NULL REFERENCES media_items(id) ON DELETE CASCADE,
    type image_type NOT NULL,
    index INT NOT NULL DEFAULT 0,             -- For multiple images of same type
    path TEXT NOT NULL,                       -- Filesystem path or URL
    width INT,
    height INT,
    blurhash VARCHAR(100),                    -- BlurHash for placeholder
    provider VARCHAR(50),                     -- Source provider (local, tmdb, tvdb, etc.)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(item_id, type, index)
);

CREATE INDEX idx_images_item_id ON images(item_id);
CREATE INDEX idx_images_type ON images(item_id, type);

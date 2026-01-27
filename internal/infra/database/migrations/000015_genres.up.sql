-- 000015_genres.up.sql
-- Domain-scoped genres with proper relationships
-- Replaces the TEXT[] genres column on media_items with normalized tables

BEGIN;

-- Genre domain enum - scopes genres to content types
CREATE TYPE genre_domain AS ENUM (
    'movie',    -- Movies, trailers
    'tv',       -- TV shows, episodes, seasons
    'music',    -- Audio, albums, artists
    'book',     -- Books, audiobooks
    'podcast',  -- Podcasts
    'game'      -- Future: games
);

-- Master genres table with domain scoping
CREATE TABLE genres (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    domain genre_domain NOT NULL,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL,               -- URL-safe: "sci-fi", "hip-hop"
    description TEXT,
    parent_id UUID REFERENCES genres(id) ON DELETE SET NULL, -- Hierarchical: Rock → Alternative Rock
    external_ids JSONB NOT NULL DEFAULT '{}', -- {"tmdb": "123", "musicbrainz": "abc"}
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(domain, slug)
);

CREATE INDEX idx_genres_domain ON genres(domain);
CREATE INDEX idx_genres_parent ON genres(parent_id) WHERE parent_id IS NOT NULL;
CREATE INDEX idx_genres_name ON genres(name);
CREATE INDEX idx_genres_slug ON genres(domain, slug);

-- Many-to-many: media items to genres
CREATE TABLE media_item_genres (
    media_item_id UUID NOT NULL REFERENCES media_items(id) ON DELETE CASCADE,
    genre_id UUID NOT NULL REFERENCES genres(id) ON DELETE CASCADE,
    source VARCHAR(50) NOT NULL DEFAULT 'manual', -- 'tmdb', 'musicbrainz', 'manual', 'nfo'
    confidence DECIMAL(3,2) DEFAULT 1.00,         -- 0.00-1.00 for auto-tagged
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (media_item_id, genre_id)
);

CREATE INDEX idx_media_item_genres_genre ON media_item_genres(genre_id);
CREATE INDEX idx_media_item_genres_source ON media_item_genres(source);

-- Trigger to update genres.updated_at
CREATE OR REPLACE FUNCTION update_genre_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER genres_updated_at
    BEFORE UPDATE ON genres
    FOR EACH ROW
    EXECUTE FUNCTION update_genre_timestamp();

-- Seed common genres per domain
-- Movie genres
INSERT INTO genres (domain, name, slug) VALUES
    ('movie', 'Action', 'action'),
    ('movie', 'Adventure', 'adventure'),
    ('movie', 'Animation', 'animation'),
    ('movie', 'Comedy', 'comedy'),
    ('movie', 'Crime', 'crime'),
    ('movie', 'Documentary', 'documentary'),
    ('movie', 'Drama', 'drama'),
    ('movie', 'Family', 'family'),
    ('movie', 'Fantasy', 'fantasy'),
    ('movie', 'History', 'history'),
    ('movie', 'Horror', 'horror'),
    ('movie', 'Music', 'music'),
    ('movie', 'Mystery', 'mystery'),
    ('movie', 'Romance', 'romance'),
    ('movie', 'Science Fiction', 'science-fiction'),
    ('movie', 'Thriller', 'thriller'),
    ('movie', 'War', 'war'),
    ('movie', 'Western', 'western'),
    ('movie', 'Adult', 'adult');

-- TV genres
INSERT INTO genres (domain, name, slug) VALUES
    ('tv', 'Action & Adventure', 'action-adventure'),
    ('tv', 'Animation', 'animation'),
    ('tv', 'Comedy', 'comedy'),
    ('tv', 'Crime', 'crime'),
    ('tv', 'Documentary', 'documentary'),
    ('tv', 'Drama', 'drama'),
    ('tv', 'Family', 'family'),
    ('tv', 'Kids', 'kids'),
    ('tv', 'Mystery', 'mystery'),
    ('tv', 'News', 'news'),
    ('tv', 'Reality', 'reality'),
    ('tv', 'Sci-Fi & Fantasy', 'sci-fi-fantasy'),
    ('tv', 'Soap', 'soap'),
    ('tv', 'Talk', 'talk'),
    ('tv', 'War & Politics', 'war-politics'),
    ('tv', 'Western', 'western');

-- Music genres (top-level)
INSERT INTO genres (domain, name, slug) VALUES
    ('music', 'Rock', 'rock'),
    ('music', 'Pop', 'pop'),
    ('music', 'Hip Hop', 'hip-hop'),
    ('music', 'R&B', 'rnb'),
    ('music', 'Electronic', 'electronic'),
    ('music', 'Jazz', 'jazz'),
    ('music', 'Classical', 'classical'),
    ('music', 'Country', 'country'),
    ('music', 'Folk', 'folk'),
    ('music', 'Blues', 'blues'),
    ('music', 'Metal', 'metal'),
    ('music', 'Punk', 'punk'),
    ('music', 'Reggae', 'reggae'),
    ('music', 'Latin', 'latin'),
    ('music', 'World', 'world'),
    ('music', 'Soul', 'soul'),
    ('music', 'Funk', 'funk'),
    ('music', 'Disco', 'disco'),
    ('music', 'House', 'house'),
    ('music', 'Techno', 'techno'),
    ('music', 'Ambient', 'ambient'),
    ('music', 'Soundtrack', 'soundtrack'),
    ('music', 'Spoken Word', 'spoken-word');

-- Music sub-genres (examples with parent references)
INSERT INTO genres (domain, name, slug, parent_id) 
SELECT 'music', 'Alternative Rock', 'alternative-rock', id FROM genres WHERE domain = 'music' AND slug = 'rock';

INSERT INTO genres (domain, name, slug, parent_id) 
SELECT 'music', 'Indie Rock', 'indie-rock', id FROM genres WHERE domain = 'music' AND slug = 'rock';

INSERT INTO genres (domain, name, slug, parent_id) 
SELECT 'music', 'Hard Rock', 'hard-rock', id FROM genres WHERE domain = 'music' AND slug = 'rock';

INSERT INTO genres (domain, name, slug, parent_id) 
SELECT 'music', 'Progressive Rock', 'progressive-rock', id FROM genres WHERE domain = 'music' AND slug = 'rock';

INSERT INTO genres (domain, name, slug, parent_id) 
SELECT 'music', 'Trap', 'trap', id FROM genres WHERE domain = 'music' AND slug = 'hip-hop';

INSERT INTO genres (domain, name, slug, parent_id) 
SELECT 'music', 'Boom Bap', 'boom-bap', id FROM genres WHERE domain = 'music' AND slug = 'hip-hop';

INSERT INTO genres (domain, name, slug, parent_id) 
SELECT 'music', 'Deep House', 'deep-house', id FROM genres WHERE domain = 'music' AND slug = 'house';

INSERT INTO genres (domain, name, slug, parent_id) 
SELECT 'music', 'Tech House', 'tech-house', id FROM genres WHERE domain = 'music' AND slug = 'house';

-- Book genres
INSERT INTO genres (domain, name, slug) VALUES
    ('book', 'Fiction', 'fiction'),
    ('book', 'Non-Fiction', 'non-fiction'),
    ('book', 'Mystery', 'mystery'),
    ('book', 'Thriller', 'thriller'),
    ('book', 'Romance', 'romance'),
    ('book', 'Science Fiction', 'science-fiction'),
    ('book', 'Fantasy', 'fantasy'),
    ('book', 'Horror', 'horror'),
    ('book', 'Biography', 'biography'),
    ('book', 'History', 'history'),
    ('book', 'Self-Help', 'self-help'),
    ('book', 'Business', 'business'),
    ('book', 'Children', 'children'),
    ('book', 'Young Adult', 'young-adult');

-- Podcast genres
INSERT INTO genres (domain, name, slug) VALUES
    ('podcast', 'News', 'news'),
    ('podcast', 'Comedy', 'comedy'),
    ('podcast', 'True Crime', 'true-crime'),
    ('podcast', 'Technology', 'technology'),
    ('podcast', 'Business', 'business'),
    ('podcast', 'Education', 'education'),
    ('podcast', 'Health', 'health'),
    ('podcast', 'Sports', 'sports'),
    ('podcast', 'Music', 'music'),
    ('podcast', 'Society & Culture', 'society-culture'),
    ('podcast', 'Science', 'science'),
    ('podcast', 'History', 'history'),
    ('podcast', 'Fiction', 'fiction');

COMMIT;

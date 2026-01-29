-- 000006_genres.up.sql
-- Global genres table for cross-module genre management

BEGIN;

-- Genres table
CREATE TABLE genres (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_genres_slug ON genres(slug);

-- Trigger for updated_at
CREATE TRIGGER update_genres_updated_at
    BEFORE UPDATE ON genres
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert common genres (movies, TV, music)
INSERT INTO genres (name, slug) VALUES
    ('Action', 'action'),
    ('Adventure', 'adventure'),
    ('Animation', 'animation'),
    ('Comedy', 'comedy'),
    ('Crime', 'crime'),
    ('Documentary', 'documentary'),
    ('Drama', 'drama'),
    ('Family', 'family'),
    ('Fantasy', 'fantasy'),
    ('History', 'history'),
    ('Horror', 'horror'),
    ('Music', 'music'),
    ('Mystery', 'mystery'),
    ('Romance', 'romance'),
    ('Science Fiction', 'science-fiction'),
    ('Thriller', 'thriller'),
    ('War', 'war'),
    ('Western', 'western'),
    -- Music genres
    ('Rock', 'rock'),
    ('Pop', 'pop'),
    ('Jazz', 'jazz'),
    ('Classical', 'classical'),
    ('Hip Hop', 'hip-hop'),
    ('Electronic', 'electronic'),
    ('Metal', 'metal'),
    ('Blues', 'blues'),
    ('Country', 'country'),
    ('R&B', 'rnb')
ON CONFLICT (slug) DO NOTHING;

COMMIT;

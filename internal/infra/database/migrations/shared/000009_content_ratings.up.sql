-- 000009_content_ratings.up.sql
-- Content rating systems (MPAA, FSK, PEGI, BBFC, etc.)

BEGIN;

-- Rating systems enum
CREATE TYPE rating_system AS ENUM (
    'MPAA',      -- US Motion Picture Association
    'FSK',       -- German Freiwillige Selbstkontrolle
    'PEGI',      -- Pan European Game Information
    'BBFC',      -- British Board of Film Classification
    'ACB',       -- Australian Classification Board
    'KIJKWIJZER' -- Dutch age rating system
);

-- Content ratings table
CREATE TABLE content_ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    system rating_system NOT NULL,
    code VARCHAR(20) NOT NULL,  -- e.g., 'PG-13', 'FSK 12', 'PEGI 18'
    display_name VARCHAR(50) NOT NULL,
    description TEXT,
    min_age INT,  -- Minimum recommended age
    icon_url TEXT,  -- Optional icon/image URL
    sort_order INT NOT NULL DEFAULT 0,  -- For UI ordering
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (system, code)
);

-- Indexes
CREATE INDEX idx_content_ratings_system ON content_ratings(system);
CREATE INDEX idx_content_ratings_min_age ON content_ratings(min_age);
CREATE INDEX idx_content_ratings_sort ON content_ratings(system, sort_order);

-- Trigger for updated_at
CREATE TRIGGER update_content_ratings_updated_at
    BEFORE UPDATE ON content_ratings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert MPAA ratings
INSERT INTO content_ratings (system, code, display_name, description, min_age, sort_order) VALUES
    ('MPAA', 'G', 'G - General Audiences', 'All ages admitted', 0, 1),
    ('MPAA', 'PG', 'PG - Parental Guidance', 'Some material may not be suitable for children', 0, 2),
    ('MPAA', 'PG-13', 'PG-13', 'Some material may be inappropriate for children under 13', 13, 3),
    ('MPAA', 'R', 'R - Restricted', 'Under 17 requires accompanying parent or adult guardian', 17, 4),
    ('MPAA', 'NC-17', 'NC-17', 'No one 17 and under admitted', 18, 5)
ON CONFLICT (system, code) DO NOTHING;

-- Insert FSK ratings (Germany)
INSERT INTO content_ratings (system, code, display_name, description, min_age, sort_order) VALUES
    ('FSK', '0', 'FSK 0', 'Approved for all ages', 0, 1),
    ('FSK', '6', 'FSK 6', 'Approved for ages 6 and up', 6, 2),
    ('FSK', '12', 'FSK 12', 'Approved for ages 12 and up', 12, 3),
    ('FSK', '16', 'FSK 16', 'Approved for ages 16 and up', 16, 4),
    ('FSK', '18', 'FSK 18', 'Adults only', 18, 5)
ON CONFLICT (system, code) DO NOTHING;

-- Insert PEGI ratings (Europe)
INSERT INTO content_ratings (system, code, display_name, description, min_age, sort_order) VALUES
    ('PEGI', '3', 'PEGI 3', 'Suitable for all ages', 3, 1),
    ('PEGI', '7', 'PEGI 7', 'Suitable for ages 7 and up', 7, 2),
    ('PEGI', '12', 'PEGI 12', 'Suitable for ages 12 and up', 12, 3),
    ('PEGI', '16', 'PEGI 16', 'Suitable for ages 16 and up', 16, 4),
    ('PEGI', '18', 'PEGI 18', 'Adults only', 18, 5)
ON CONFLICT (system, code) DO NOTHING;

-- Insert BBFC ratings (UK)
INSERT INTO content_ratings (system, code, display_name, description, min_age, sort_order) VALUES
    ('BBFC', 'U', 'U - Universal', 'Suitable for all', 0, 1),
    ('BBFC', 'PG', 'PG', 'Parental guidance', 0, 2),
    ('BBFC', '12A', '12A', 'Cinema release for 12 and over', 12, 3),
    ('BBFC', '12', '12', 'Video release for 12 and over', 12, 4),
    ('BBFC', '15', '15', 'Suitable for 15 and over', 15, 5),
    ('BBFC', '18', '18', 'Adults only', 18, 6),
    ('BBFC', 'R18', 'R18', 'Restricted to licensed premises', 18, 7)
ON CONFLICT (system, code) DO NOTHING;

COMMIT;

-- 000009_people.up.sql
-- People (actors, directors, etc.) and their links to media

CREATE TYPE person_type AS ENUM ('actor', 'director', 'writer', 'producer', 'composer', 'guest_star');

CREATE TABLE people (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    sort_name VARCHAR(255),
    overview TEXT,
    birth_date DATE,
    death_date DATE,
    birth_place VARCHAR(255),
    provider_ids JSONB NOT NULL DEFAULT '{}', -- {"imdb": "nm123", "tmdb": "456"}
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_people_name ON people(name);
CREATE INDEX idx_people_sort_name ON people(sort_name);

CREATE TRIGGER update_people_updated_at
    BEFORE UPDATE ON people
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Junction table linking media items to people
CREATE TABLE media_people (
    item_id UUID NOT NULL REFERENCES media_items(id) ON DELETE CASCADE,
    person_id UUID NOT NULL REFERENCES people(id) ON DELETE CASCADE,
    type person_type NOT NULL,
    role VARCHAR(255),                        -- Character name for actors
    sort_order INT NOT NULL DEFAULT 0,
    PRIMARY KEY (item_id, person_id, type)
);

CREATE INDEX idx_media_people_item_id ON media_people(item_id);
CREATE INDEX idx_media_people_person_id ON media_people(person_id);
CREATE INDEX idx_media_people_type ON media_people(item_id, type);

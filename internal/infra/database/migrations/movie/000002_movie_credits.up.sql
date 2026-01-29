-- Movie Credits: Cast and crew relationships
BEGIN;

-- Movie People: actors, directors, writers, etc.
CREATE TABLE movie_people (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                VARCHAR(255) NOT NULL,
    sort_name           VARCHAR(255),
    original_name       VARCHAR(255),

    -- Bio
    biography           TEXT,
    birthdate           DATE,
    deathdate           DATE,
    birthplace          VARCHAR(255),
    gender              VARCHAR(20),

    -- Images
    primary_image_url   TEXT,
    primary_image_blurhash VARCHAR(50),

    -- External IDs
    tmdb_id             INT,
    imdb_id             VARCHAR(20),
    tvdb_id             INT,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_movie_people_name ON movie_people(name);
CREATE INDEX idx_movie_people_tmdb ON movie_people(tmdb_id) WHERE tmdb_id IS NOT NULL;
CREATE INDEX idx_movie_people_imdb ON movie_people(imdb_id) WHERE imdb_id IS NOT NULL;

CREATE TRIGGER movie_people_updated_at
    BEFORE UPDATE ON movie_people
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Credit role enum
CREATE TYPE movie_credit_role AS ENUM (
    'actor',
    'director',
    'writer',
    'producer',
    'executive_producer',
    'composer',
    'cinematographer',
    'editor',
    'production_designer',
    'costume_designer',
    'makeup_artist',
    'visual_effects',
    'stunt_coordinator',
    'sound_designer'
);

-- Movie Credits: Link movies to people with roles
CREATE TABLE movie_credits (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    movie_id            UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    person_id           UUID NOT NULL REFERENCES movie_people(id) ON DELETE CASCADE,

    role                movie_credit_role NOT NULL,
    character_name      VARCHAR(255),        -- For actors
    department          VARCHAR(100),        -- For crew (e.g., "Art", "Sound")
    job                 VARCHAR(100),        -- Specific job title

    billing_order       INT NOT NULL DEFAULT 0,
    is_guest            BOOLEAN NOT NULL DEFAULT false,

    -- TMDb credit ID for updates
    tmdb_credit_id      VARCHAR(50),

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_movie_credits_movie ON movie_credits(movie_id);
CREATE INDEX idx_movie_credits_person ON movie_credits(person_id);
CREATE INDEX idx_movie_credits_role ON movie_credits(role);
CREATE INDEX idx_movie_credits_billing ON movie_credits(movie_id, billing_order);

-- Unique constraint to prevent duplicate credits
CREATE UNIQUE INDEX idx_movie_credits_unique ON movie_credits(movie_id, person_id, role, COALESCE(character_name, ''));

COMMIT;

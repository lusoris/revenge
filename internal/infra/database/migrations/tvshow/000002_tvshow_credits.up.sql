-- TV Show Credits: Cast and crew relationships
BEGIN;

-- TV Show People: actors, directors, writers, etc.
CREATE TABLE tvshow_people (
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

CREATE INDEX idx_tvshow_people_name ON tvshow_people(name);
CREATE INDEX idx_tvshow_people_tmdb ON tvshow_people(tmdb_id) WHERE tmdb_id IS NOT NULL;
CREATE INDEX idx_tvshow_people_imdb ON tvshow_people(imdb_id) WHERE imdb_id IS NOT NULL;

CREATE TRIGGER tvshow_people_updated_at
    BEFORE UPDATE ON tvshow_people
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Credit role enum for TV
CREATE TYPE tvshow_credit_role AS ENUM (
    'actor',
    'director',
    'writer',
    'creator',
    'showrunner',
    'producer',
    'executive_producer',
    'composer',
    'cinematographer',
    'editor',
    'guest_star'
);

-- Series Credits: Regular cast and crew for the whole series
CREATE TABLE series_credits (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    series_id           UUID NOT NULL REFERENCES series(id) ON DELETE CASCADE,
    person_id           UUID NOT NULL REFERENCES tvshow_people(id) ON DELETE CASCADE,

    role                tvshow_credit_role NOT NULL,
    character_name      VARCHAR(255),
    department          VARCHAR(100),
    job                 VARCHAR(100),

    billing_order       INT NOT NULL DEFAULT 0,

    tmdb_credit_id      VARCHAR(50),

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_series_credits_series ON series_credits(series_id);
CREATE INDEX idx_series_credits_person ON series_credits(person_id);
CREATE INDEX idx_series_credits_role ON series_credits(role);
CREATE INDEX idx_series_credits_billing ON series_credits(series_id, billing_order);

CREATE UNIQUE INDEX idx_series_credits_unique ON series_credits(series_id, person_id, role, COALESCE(character_name, ''));

-- Episode Credits: Guest stars and episode-specific crew
CREATE TABLE episode_credits (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    episode_id          UUID NOT NULL REFERENCES episodes(id) ON DELETE CASCADE,
    person_id           UUID NOT NULL REFERENCES tvshow_people(id) ON DELETE CASCADE,

    role                tvshow_credit_role NOT NULL,
    character_name      VARCHAR(255),
    department          VARCHAR(100),
    job                 VARCHAR(100),

    billing_order       INT NOT NULL DEFAULT 0,
    is_guest            BOOLEAN NOT NULL DEFAULT false,

    tmdb_credit_id      VARCHAR(50),

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_episode_credits_episode ON episode_credits(episode_id);
CREATE INDEX idx_episode_credits_person ON episode_credits(person_id);
CREATE INDEX idx_episode_credits_role ON episode_credits(role);
CREATE INDEX idx_episode_credits_guest ON episode_credits(episode_id) WHERE is_guest = true;

CREATE UNIQUE INDEX idx_episode_credits_unique ON episode_credits(episode_id, person_id, role, COALESCE(character_name, ''));

COMMIT;

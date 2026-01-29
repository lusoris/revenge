-- Movie Credits: Cast and crew relationships
-- Uses shared video_people table from shared/000017_video_people
BEGIN;

-- Movie Credits: Link movies to people with roles
CREATE TABLE movie_credits (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    movie_id            UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    person_id           UUID NOT NULL REFERENCES video_people(id) ON DELETE CASCADE,

    role                video_credit_role NOT NULL,
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

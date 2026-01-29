-- Adult movies and shared adult metadata (schema c)
CREATE SCHEMA IF NOT EXISTS c;

-- Studios
CREATE TABLE c.studios (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    parent_id       UUID REFERENCES c.studios(id),
    stashdb_id      VARCHAR(100),
    tpdb_id         VARCHAR(100),
    url             TEXT,
    logo_path       TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(stashdb_id),
    UNIQUE(tpdb_id)
);

CREATE TRIGGER c_studios_updated_at
    BEFORE UPDATE ON c.studios
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Performers
CREATE TABLE c.performers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL,
    disambiguation  VARCHAR(255),
    gender          VARCHAR(50),
    birthdate       DATE,
    death_date      DATE,
    birth_city      VARCHAR(255),
    ethnicity       VARCHAR(100),
    nationality     VARCHAR(100),
    hair_color      VARCHAR(50),
    eye_color       VARCHAR(50),
    height_cm       INT,
    weight_kg       INT,
    measurements    VARCHAR(50),
    cup_size        VARCHAR(10),
    breast_type     VARCHAR(50),
    tattoos         TEXT,
    piercings       TEXT,
    career_start    INT,
    career_end      INT,
    bio             TEXT,

    stash_id        VARCHAR(100),
    stashdb_id      VARCHAR(100),
    tpdb_id         VARCHAR(100),
    freeones_id     VARCHAR(100),

    twitter         TEXT,
    instagram       TEXT,

    image_path      TEXT,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(stashdb_id),
    UNIQUE(tpdb_id)
);

CREATE INDEX idx_c_performers_name ON c.performers(name);
CREATE INDEX idx_c_performers_stashdb ON c.performers(stashdb_id);

CREATE TRIGGER c_performers_updated_at
    BEFORE UPDATE ON c.performers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Performer aliases
CREATE TABLE c.performer_aliases (
    performer_id    UUID REFERENCES c.performers(id) ON DELETE CASCADE,
    alias           VARCHAR(255) NOT NULL,
    PRIMARY KEY (performer_id, alias)
);

-- Performer images
CREATE TABLE c.performer_images (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    performer_id    UUID REFERENCES c.performers(id) ON DELETE CASCADE,
    path            TEXT NOT NULL,
    type            VARCHAR(50) DEFAULT 'photo',
    source          VARCHAR(50),
    primary_image   BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Tags
CREATE TABLE c.tags (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name            VARCHAR(255) NOT NULL UNIQUE,
    description     TEXT,
    parent_id       UUID REFERENCES c.tags(id),
    stashdb_id      VARCHAR(100),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Movies
CREATE TABLE c.movies (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    library_id      UUID NOT NULL,
    whisparr_id     INT,
    stashdb_id      VARCHAR(100),
    tpdb_id         VARCHAR(100),

    title           VARCHAR(500) NOT NULL,
    sort_title      VARCHAR(500),
    original_title  VARCHAR(500),
    overview        TEXT,
    release_date    DATE,
    runtime_ticks   BIGINT,
    studio_id       UUID REFERENCES c.studios(id),
    director        VARCHAR(255),
    series          VARCHAR(255),

    path            TEXT NOT NULL,
    size_bytes      BIGINT,
    container       VARCHAR(50),
    video_codec     VARCHAR(50),
    audio_codec     VARCHAR(50),
    resolution      VARCHAR(50),

    phash           VARCHAR(64),
    oshash          VARCHAR(64),

    has_file        BOOLEAN DEFAULT TRUE,
    is_hdr          BOOLEAN DEFAULT FALSE,
    is_3d           BOOLEAN DEFAULT FALSE,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(path),
    UNIQUE(stashdb_id)
);

CREATE INDEX idx_c_movies_library ON c.movies(library_id);
CREATE INDEX idx_c_movies_studio ON c.movies(studio_id);
CREATE INDEX idx_c_movies_phash ON c.movies(phash);
CREATE INDEX idx_c_movies_oshash ON c.movies(oshash);

CREATE TRIGGER c_movies_updated_at
    BEFORE UPDATE ON c.movies
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Movie performers
CREATE TABLE c.movie_performers (
    movie_id        UUID REFERENCES c.movies(id) ON DELETE CASCADE,
    performer_id    UUID REFERENCES c.performers(id) ON DELETE CASCADE,
    character_name  VARCHAR(255),
    PRIMARY KEY (movie_id, performer_id)
);

-- Movie tags
CREATE TABLE c.movie_tags (
    movie_id        UUID REFERENCES c.movies(id) ON DELETE CASCADE,
    tag_id          UUID REFERENCES c.tags(id) ON DELETE CASCADE,
    PRIMARY KEY (movie_id, tag_id)
);

-- Movie images
CREATE TABLE c.movie_images (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    movie_id        UUID REFERENCES c.movies(id) ON DELETE CASCADE,
    type            VARCHAR(50) NOT NULL,
    path            TEXT NOT NULL,
    source          VARCHAR(50),
    primary_image   BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Galleries
CREATE TABLE c.galleries (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    movie_id        UUID REFERENCES c.movies(id) ON DELETE SET NULL,
    title           VARCHAR(500) NOT NULL,
    path            TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE c.gallery_images (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    gallery_id      UUID REFERENCES c.galleries(id) ON DELETE CASCADE,
    path            TEXT NOT NULL,
    sort_order      INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- User ratings
CREATE TABLE c.user_ratings (
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
    movie_id        UUID REFERENCES c.movies(id) ON DELETE CASCADE,
    rating          SMALLINT CHECK (rating >= 1 AND rating <= 10),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, movie_id)
);

CREATE TRIGGER c_user_ratings_updated_at
    BEFORE UPDATE ON c.user_ratings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- User favorites
CREATE TABLE c.user_favorites (
    user_id         UUID NOT NULL,
    movie_id        UUID REFERENCES c.movies(id) ON DELETE CASCADE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, movie_id)
);

-- Watch history
CREATE TABLE c.watch_history (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    movie_id        UUID REFERENCES c.movies(id) ON DELETE CASCADE,
    watched_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    position_ticks  BIGINT,
    completed       BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_c_watch_history_user ON c.watch_history(user_id, watched_at DESC);

-- User performer favorites
CREATE TABLE c.user_performer_favorites (
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
    performer_id    UUID REFERENCES c.performers(id) ON DELETE CASCADE,
    added_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, performer_id)
);

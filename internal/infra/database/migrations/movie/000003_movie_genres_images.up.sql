-- Movie Genres and Images
BEGIN;

-- Movie Genres: Per-module genre definitions
CREATE TABLE movie_genres (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                VARCHAR(100) NOT NULL,
    slug                VARCHAR(100) NOT NULL,
    description         TEXT,
    parent_id           UUID REFERENCES movie_genres(id) ON DELETE SET NULL,
    external_ids        JSONB NOT NULL DEFAULT '{}'::jsonb,  -- {tmdb: "28"}

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (slug)
);

CREATE INDEX idx_movie_genres_parent ON movie_genres(parent_id) WHERE parent_id IS NOT NULL;

CREATE TRIGGER movie_genres_updated_at
    BEFORE UPDATE ON movie_genres
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Junction: Movie <-> Genre
CREATE TABLE movie_genre_link (
    movie_id            UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    genre_id            UUID NOT NULL REFERENCES movie_genres(id) ON DELETE CASCADE,

    PRIMARY KEY (movie_id, genre_id)
);

CREATE INDEX idx_movie_genre_link_genre ON movie_genre_link(genre_id);

-- Movie Images: Additional artwork beyond primary poster/backdrop
CREATE TYPE movie_image_type AS ENUM (
    'poster',
    'backdrop',
    'logo',
    'thumb',
    'banner',
    'disc',
    'clearart',
    'clearlogo',
    'keyart'
);

CREATE TABLE movie_images (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    movie_id            UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,

    image_type          movie_image_type NOT NULL,
    path                TEXT NOT NULL,
    blurhash            VARCHAR(50),
    width               INT,
    height              INT,
    aspect_ratio        DECIMAL(5,3),
    language            VARCHAR(10),
    vote_average        DECIMAL(3,1),
    vote_count          INT,
    is_primary          BOOLEAN NOT NULL DEFAULT false,

    -- Source info
    source              VARCHAR(50) NOT NULL DEFAULT 'tmdb',  -- tmdb, fanart, local

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_movie_images_movie ON movie_images(movie_id);
CREATE INDEX idx_movie_images_type ON movie_images(movie_id, image_type);
CREATE INDEX idx_movie_images_primary ON movie_images(movie_id, image_type, is_primary) WHERE is_primary = true;

-- Movie Videos: Trailers, clips, featurettes
CREATE TYPE movie_video_type AS ENUM (
    'trailer',
    'teaser',
    'clip',
    'featurette',
    'behind_the_scenes',
    'bloopers'
);

CREATE TYPE movie_video_site AS ENUM (
    'youtube',
    'vimeo'
);

CREATE TABLE movie_videos (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    movie_id            UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,

    video_type          movie_video_type NOT NULL,
    site                movie_video_site NOT NULL,
    key                 VARCHAR(50) NOT NULL,           -- YouTube/Vimeo video ID
    name                VARCHAR(255),
    language            VARCHAR(10),
    size                INT,                            -- 360, 480, 720, 1080

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_movie_videos_movie ON movie_videos(movie_id);
CREATE INDEX idx_movie_videos_type ON movie_videos(movie_id, video_type);

COMMIT;

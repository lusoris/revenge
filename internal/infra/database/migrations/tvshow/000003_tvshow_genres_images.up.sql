-- TV Show Genres and Images
BEGIN;

-- Genres (shared with movies but module-isolated for flexibility)
CREATE TABLE tvshow_genres (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                VARCHAR(100) NOT NULL UNIQUE,
    tmdb_id             INT UNIQUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tvshow_genres_name ON tvshow_genres(name);

-- Series <-> Genre junction
CREATE TABLE series_genre_link (
    series_id           UUID NOT NULL REFERENCES series(id) ON DELETE CASCADE,
    genre_id            UUID NOT NULL REFERENCES tvshow_genres(id) ON DELETE CASCADE,
    PRIMARY KEY (series_id, genre_id)
);

CREATE INDEX idx_series_genre_link_genre ON series_genre_link(genre_id);

-- Series Images (posters, backdrops, logos, banners)
CREATE TABLE series_images (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    series_id           UUID NOT NULL REFERENCES series(id) ON DELETE CASCADE,

    image_type          VARCHAR(50) NOT NULL,    -- poster, backdrop, logo, banner, thumb
    url                 TEXT NOT NULL,
    local_path          TEXT,                    -- Cached locally

    -- Dimensions
    width               INT,
    height              INT,
    aspect_ratio        DECIMAL(4,2),

    -- Quality info
    language            VARCHAR(10),
    vote_average        DECIMAL(3,1),
    vote_count          INT,

    -- Blurhash for placeholders
    blurhash            VARCHAR(50),

    -- Provider info
    provider            VARCHAR(50),             -- tmdb, tvdb, fanart
    provider_id         VARCHAR(100),

    is_primary          BOOLEAN NOT NULL DEFAULT false,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_series_images_series ON series_images(series_id);
CREATE INDEX idx_series_images_type ON series_images(series_id, image_type);
CREATE INDEX idx_series_images_primary ON series_images(series_id, image_type) WHERE is_primary = true;

-- Season Images
CREATE TABLE season_images (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    season_id           UUID NOT NULL REFERENCES seasons(id) ON DELETE CASCADE,

    image_type          VARCHAR(50) NOT NULL,    -- poster, backdrop
    url                 TEXT NOT NULL,
    local_path          TEXT,

    width               INT,
    height              INT,
    aspect_ratio        DECIMAL(4,2),

    language            VARCHAR(10),
    vote_average        DECIMAL(3,1),
    vote_count          INT,
    blurhash            VARCHAR(50),

    provider            VARCHAR(50),
    provider_id         VARCHAR(100),

    is_primary          BOOLEAN NOT NULL DEFAULT false,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_season_images_season ON season_images(season_id);
CREATE INDEX idx_season_images_type ON season_images(season_id, image_type);

-- Episode Images (stills/thumbnails)
CREATE TABLE episode_images (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    episode_id          UUID NOT NULL REFERENCES episodes(id) ON DELETE CASCADE,

    image_type          VARCHAR(50) NOT NULL,    -- still, thumb
    url                 TEXT NOT NULL,
    local_path          TEXT,

    width               INT,
    height              INT,
    aspect_ratio        DECIMAL(4,2),

    vote_average        DECIMAL(3,1),
    vote_count          INT,
    blurhash            VARCHAR(50),

    provider            VARCHAR(50),
    provider_id         VARCHAR(100),

    is_primary          BOOLEAN NOT NULL DEFAULT false,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_episode_images_episode ON episode_images(episode_id);
CREATE INDEX idx_episode_images_primary ON episode_images(episode_id) WHERE is_primary = true;

-- Series Videos (trailers, teasers, behind the scenes)
CREATE TABLE series_videos (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    series_id           UUID NOT NULL REFERENCES series(id) ON DELETE CASCADE,

    video_type          VARCHAR(50) NOT NULL,    -- trailer, teaser, featurette, behind_the_scenes
    name                VARCHAR(255),
    key                 VARCHAR(100),            -- YouTube/Vimeo key
    site                VARCHAR(50),             -- youtube, vimeo
    size                INT,                     -- Quality: 360, 480, 720, 1080

    language            VARCHAR(10),
    is_official         BOOLEAN NOT NULL DEFAULT true,

    tmdb_id             VARCHAR(50),

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_series_videos_series ON series_videos(series_id);
CREATE INDEX idx_series_videos_type ON series_videos(series_id, video_type);

COMMIT;

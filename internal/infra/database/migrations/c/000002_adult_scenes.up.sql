-- Adult scenes (schema c)
CREATE SCHEMA IF NOT EXISTS c;

CREATE TABLE c.scenes (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    library_id      UUID NOT NULL,
    title           VARCHAR(500) NOT NULL,
    sort_title      VARCHAR(500),
    overview        TEXT,
    release_date    DATE,
    runtime_minutes INT,
    studio_id       UUID REFERENCES c.studios(id),

    whisparr_id     INT,
    stash_id        VARCHAR(100),
    stashdb_id      VARCHAR(100),
    tpdb_id         VARCHAR(100),

    path            TEXT NOT NULL,
    size_bytes      BIGINT,
    video_codec     VARCHAR(50),
    audio_codec     VARCHAR(50),
    resolution      VARCHAR(20),

    oshash          VARCHAR(32),
    phash           VARCHAR(32),
    md5             VARCHAR(64),

    cover_path      TEXT,

    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(path)
);

CREATE INDEX idx_c_scenes_library ON c.scenes(library_id);
CREATE INDEX idx_c_scenes_studio ON c.scenes(studio_id);
CREATE INDEX idx_c_scenes_oshash ON c.scenes(oshash);
CREATE INDEX idx_c_scenes_stashdb ON c.scenes(stashdb_id);

CREATE TRIGGER c_scenes_updated_at
    BEFORE UPDATE ON c.scenes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Scene-performer relationship
CREATE TABLE c.scene_performers (
    scene_id        UUID REFERENCES c.scenes(id) ON DELETE CASCADE,
    performer_id    UUID REFERENCES c.performers(id) ON DELETE CASCADE,
    role            VARCHAR(100),
    PRIMARY KEY (scene_id, performer_id)
);

-- Scene tags
CREATE TABLE c.scene_tags (
    scene_id        UUID REFERENCES c.scenes(id) ON DELETE CASCADE,
    tag_id          UUID REFERENCES c.tags(id) ON DELETE CASCADE,
    PRIMARY KEY (scene_id, tag_id)
);

-- Scene markers (chapters/positions)
CREATE TABLE c.scene_markers (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scene_id        UUID REFERENCES c.scenes(id) ON DELETE CASCADE,
    title           VARCHAR(255),
    start_seconds   FLOAT NOT NULL,
    end_seconds     FLOAT,
    tag_id          UUID REFERENCES c.tags(id),
    stash_marker_id VARCHAR(100),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_c_markers_scene ON c.scene_markers(scene_id);

-- User data (per-module, in c schema)
CREATE TABLE c.user_scene_data (
    user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
    scene_id        UUID REFERENCES c.scenes(id) ON DELETE CASCADE,
    position_ms     BIGINT DEFAULT 0,
    watch_count     INT DEFAULT 0,
    last_watched    TIMESTAMPTZ,
    rating          SMALLINT CHECK (rating >= 1 AND rating <= 10),
    o_counter       INT DEFAULT 0,
    is_favorite     BOOLEAN DEFAULT FALSE,
    is_organized    BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (user_id, scene_id)
);

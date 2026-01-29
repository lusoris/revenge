-- Movie User Data: Ratings, favorites, watch history
BEGIN;

-- Movie User Ratings
CREATE TABLE movie_user_ratings (
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    movie_id            UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,

    rating              DECIMAL(3,1) NOT NULL CHECK (rating >= 0 AND rating <= 10),
    review              TEXT,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_id, movie_id)
);

CREATE INDEX idx_movie_user_ratings_movie ON movie_user_ratings(movie_id);

CREATE TRIGGER movie_user_ratings_updated_at
    BEFORE UPDATE ON movie_user_ratings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Movie Favorites
CREATE TABLE movie_favorites (
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    movie_id            UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_id, movie_id)
);

CREATE INDEX idx_movie_favorites_movie ON movie_favorites(movie_id);
CREATE INDEX idx_movie_favorites_user_date ON movie_favorites(user_id, created_at DESC);

-- Movie Watch History: Track playback progress
CREATE TABLE movie_watch_history (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    profile_id          UUID REFERENCES profiles(id) ON DELETE SET NULL,
    movie_id            UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,

    -- Playback state
    position_ticks      BIGINT NOT NULL DEFAULT 0,       -- Current position
    duration_ticks      BIGINT,                          -- Total duration at time of play
    played_percentage   DECIMAL(5,2) GENERATED ALWAYS AS (
        CASE WHEN duration_ticks > 0
        THEN (position_ticks::DECIMAL / duration_ticks * 100)
        ELSE 0 END
    ) STORED,

    -- Completion tracking
    completed           BOOLEAN NOT NULL DEFAULT false,
    completed_at        TIMESTAMPTZ,

    -- Session info
    device_name         VARCHAR(100),
    device_type         VARCHAR(50),
    client_name         VARCHAR(100),
    play_method         VARCHAR(50),                     -- direct, transcode

    -- Timestamps
    started_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_movie_watch_history_user ON movie_watch_history(user_id);
CREATE INDEX idx_movie_watch_history_profile ON movie_watch_history(profile_id) WHERE profile_id IS NOT NULL;
CREATE INDEX idx_movie_watch_history_movie ON movie_watch_history(movie_id);
CREATE INDEX idx_movie_watch_history_recent ON movie_watch_history(user_id, last_updated_at DESC);

-- Unique constraint: one active playback session per user/movie
CREATE UNIQUE INDEX idx_movie_watch_history_active ON movie_watch_history(user_id, movie_id)
    WHERE completed = false;

-- Movie Watchlist: Movies users want to watch
CREATE TABLE movie_watchlist (
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    movie_id            UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,

    added_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    sort_order          INT,

    PRIMARY KEY (user_id, movie_id)
);

CREATE INDEX idx_movie_watchlist_user ON movie_watchlist(user_id, added_at DESC);

-- External Ratings: From other services (Rotten Tomatoes, etc.)
CREATE TABLE movie_external_ratings (
    movie_id            UUID NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    source              VARCHAR(50) NOT NULL,            -- rotten_tomatoes, metacritic, imdb

    rating              DECIMAL(5,2),                    -- Normalized to 0-100
    vote_count          INT,
    certified           BOOLEAN DEFAULT false,           -- e.g., "Certified Fresh"

    last_updated        TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (movie_id, source)
);

COMMIT;

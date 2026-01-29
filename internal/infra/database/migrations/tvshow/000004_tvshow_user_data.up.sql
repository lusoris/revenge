-- TV Show User Data: Ratings, favorites, watch history
BEGIN;

-- Series User Ratings
CREATE TABLE series_user_ratings (
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    series_id           UUID NOT NULL REFERENCES series(id) ON DELETE CASCADE,

    rating              DECIMAL(3,1) NOT NULL CHECK (rating >= 0 AND rating <= 10),
    review              TEXT,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_id, series_id)
);

CREATE INDEX idx_series_user_ratings_series ON series_user_ratings(series_id);

CREATE TRIGGER series_user_ratings_updated_at
    BEFORE UPDATE ON series_user_ratings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Episode User Ratings (optional, for per-episode ratings)
CREATE TABLE episode_user_ratings (
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    episode_id          UUID NOT NULL REFERENCES episodes(id) ON DELETE CASCADE,

    rating              DECIMAL(3,1) NOT NULL CHECK (rating >= 0 AND rating <= 10),

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_id, episode_id)
);

CREATE INDEX idx_episode_user_ratings_episode ON episode_user_ratings(episode_id);

CREATE TRIGGER episode_user_ratings_updated_at
    BEFORE UPDATE ON episode_user_ratings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Series Favorites
CREATE TABLE series_favorites (
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    series_id           UUID NOT NULL REFERENCES series(id) ON DELETE CASCADE,

    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_id, series_id)
);

CREATE INDEX idx_series_favorites_series ON series_favorites(series_id);
CREATE INDEX idx_series_favorites_user_date ON series_favorites(user_id, created_at DESC);

-- Episode Watch History: Track playback progress per episode
CREATE TABLE episode_watch_history (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    profile_id          UUID REFERENCES profiles(id) ON DELETE SET NULL,
    episode_id          UUID NOT NULL REFERENCES episodes(id) ON DELETE CASCADE,

    -- Playback state
    position_ticks      BIGINT NOT NULL DEFAULT 0,
    duration_ticks      BIGINT,
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
    play_method         VARCHAR(50),

    -- Timestamps
    started_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_episode_watch_history_user ON episode_watch_history(user_id);
CREATE INDEX idx_episode_watch_history_profile ON episode_watch_history(profile_id) WHERE profile_id IS NOT NULL;
CREATE INDEX idx_episode_watch_history_episode ON episode_watch_history(episode_id);
CREATE INDEX idx_episode_watch_history_recent ON episode_watch_history(user_id, last_updated_at DESC);

-- Unique: one active playback per user/episode
CREATE UNIQUE INDEX idx_episode_watch_history_active ON episode_watch_history(user_id, episode_id)
    WHERE completed = false;

-- Series Watchlist: Shows users want to watch
CREATE TABLE series_watchlist (
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    series_id           UUID NOT NULL REFERENCES series(id) ON DELETE CASCADE,

    added_at            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    sort_order          INT,

    PRIMARY KEY (user_id, series_id)
);

CREATE INDEX idx_series_watchlist_user ON series_watchlist(user_id, added_at DESC);

-- Series Watch Progress: Aggregated progress per series
-- This is a materialized view/cache for "Continue Watching" queries
CREATE TABLE series_watch_progress (
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    series_id           UUID NOT NULL REFERENCES series(id) ON DELETE CASCADE,

    -- Last watched episode
    last_episode_id     UUID REFERENCES episodes(id) ON DELETE SET NULL,
    last_season_number  INT,
    last_episode_number INT,

    -- Progress
    total_episodes      INT NOT NULL DEFAULT 0,
    watched_episodes    INT NOT NULL DEFAULT 0,
    progress_percentage DECIMAL(5,2) GENERATED ALWAYS AS (
        CASE WHEN total_episodes > 0
        THEN (watched_episodes::DECIMAL / total_episodes * 100)
        ELSE 0 END
    ) STORED,

    -- Is currently watching (has incomplete episode)
    is_watching         BOOLEAN NOT NULL DEFAULT false,

    -- Timestamps
    started_at          TIMESTAMPTZ,
    last_watched_at     TIMESTAMPTZ,
    completed_at        TIMESTAMPTZ,

    PRIMARY KEY (user_id, series_id)
);

CREATE INDEX idx_series_watch_progress_watching ON series_watch_progress(user_id, last_watched_at DESC)
    WHERE is_watching = true;

-- Function to update series watch progress
CREATE OR REPLACE FUNCTION update_series_watch_progress()
RETURNS TRIGGER AS $$
DECLARE
    v_series_id UUID;
    v_total_episodes INT;
    v_watched_episodes INT;
    v_last_episode_id UUID;
    v_last_season INT;
    v_last_episode INT;
    v_is_watching BOOLEAN;
    v_started_at TIMESTAMPTZ;
    v_completed_at TIMESTAMPTZ;
BEGIN
    -- Get series_id from episode
    SELECT series_id INTO v_series_id FROM episodes WHERE id = COALESCE(NEW.episode_id, OLD.episode_id);

    -- Count total and watched episodes
    SELECT COUNT(*) INTO v_total_episodes
    FROM episodes WHERE series_id = v_series_id AND season_number > 0;

    SELECT COUNT(DISTINCT e.id) INTO v_watched_episodes
    FROM episodes e
    JOIN episode_watch_history ewh ON ewh.episode_id = e.id
    WHERE e.series_id = v_series_id
      AND e.season_number > 0
      AND ewh.user_id = COALESCE(NEW.user_id, OLD.user_id)
      AND ewh.completed = true;

    -- Get last watched episode
    SELECT ewh.episode_id, e.season_number, e.episode_number
    INTO v_last_episode_id, v_last_season, v_last_episode
    FROM episode_watch_history ewh
    JOIN episodes e ON e.id = ewh.episode_id
    WHERE e.series_id = v_series_id
      AND ewh.user_id = COALESCE(NEW.user_id, OLD.user_id)
    ORDER BY ewh.last_updated_at DESC
    LIMIT 1;

    -- Check if currently watching (has incomplete episode)
    SELECT EXISTS(
        SELECT 1 FROM episode_watch_history ewh
        JOIN episodes e ON e.id = ewh.episode_id
        WHERE e.series_id = v_series_id
          AND ewh.user_id = COALESCE(NEW.user_id, OLD.user_id)
          AND ewh.completed = false
    ) INTO v_is_watching;

    -- Get started_at (first watch)
    SELECT MIN(started_at) INTO v_started_at
    FROM episode_watch_history ewh
    JOIN episodes e ON e.id = ewh.episode_id
    WHERE e.series_id = v_series_id
      AND ewh.user_id = COALESCE(NEW.user_id, OLD.user_id);

    -- Check if series completed
    IF v_watched_episodes >= v_total_episodes AND v_total_episodes > 0 THEN
        v_completed_at := NOW();
    ELSE
        v_completed_at := NULL;
    END IF;

    -- Upsert progress
    INSERT INTO series_watch_progress (
        user_id, series_id, last_episode_id, last_season_number, last_episode_number,
        total_episodes, watched_episodes, is_watching, started_at, last_watched_at, completed_at
    ) VALUES (
        COALESCE(NEW.user_id, OLD.user_id), v_series_id, v_last_episode_id, v_last_season, v_last_episode,
        v_total_episodes, v_watched_episodes, v_is_watching, v_started_at, NOW(), v_completed_at
    )
    ON CONFLICT (user_id, series_id) DO UPDATE SET
        last_episode_id = EXCLUDED.last_episode_id,
        last_season_number = EXCLUDED.last_season_number,
        last_episode_number = EXCLUDED.last_episode_number,
        total_episodes = EXCLUDED.total_episodes,
        watched_episodes = EXCLUDED.watched_episodes,
        is_watching = EXCLUDED.is_watching,
        last_watched_at = EXCLUDED.last_watched_at,
        completed_at = EXCLUDED.completed_at;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER episode_watch_history_update_progress
    AFTER INSERT OR UPDATE ON episode_watch_history
    FOR EACH ROW EXECUTE FUNCTION update_series_watch_progress();

-- External Ratings for series
CREATE TABLE series_external_ratings (
    series_id           UUID NOT NULL REFERENCES series(id) ON DELETE CASCADE,
    source              VARCHAR(50) NOT NULL,

    rating              DECIMAL(5,2),
    vote_count          INT,
    certified           BOOLEAN DEFAULT false,

    last_updated        TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (series_id, source)
);

COMMIT;

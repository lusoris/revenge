-- Create tvshow schema
-- Phase A11: TV Shows Module
CREATE SCHEMA IF NOT EXISTS tvshow;

-- Series table
CREATE TABLE tvshow.series (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- External IDs
    tmdb_id INTEGER UNIQUE,
    tvdb_id INTEGER UNIQUE,
    imdb_id TEXT,
    sonarr_id INTEGER,

    -- Default language fields
    title TEXT NOT NULL,
    tagline TEXT,
    overview TEXT,

    -- Multi-language support (A9)
    titles_i18n JSONB NOT NULL DEFAULT '{}',
    taglines_i18n JSONB NOT NULL DEFAULT '{}',
    overviews_i18n JSONB NOT NULL DEFAULT '{}',
    age_ratings JSONB NOT NULL DEFAULT '{}',

    -- Original language info
    original_language TEXT NOT NULL DEFAULT 'en',
    original_title TEXT,

    -- Series metadata
    status TEXT, -- 'Returning Series', 'Ended', 'Canceled', 'In Production'
    type TEXT,   -- 'Scripted', 'Documentary', 'Reality', 'Animation', 'Talk Show', 'News'
    first_air_date DATE,
    last_air_date DATE,

    -- Ratings
    vote_average DECIMAL(4, 2),
    vote_count INTEGER,
    popularity DECIMAL(10, 3),

    -- Media paths
    poster_path TEXT,
    backdrop_path TEXT,

    -- Stats (denormalized for performance)
    total_seasons INTEGER NOT NULL DEFAULT 0,
    total_episodes INTEGER NOT NULL DEFAULT 0,

    -- External integrations
    trailer_url TEXT,
    homepage TEXT,

    -- Timestamps
    metadata_updated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Seasons table
CREATE TABLE tvshow.seasons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    series_id UUID NOT NULL REFERENCES tvshow.series(id) ON DELETE CASCADE,

    -- External IDs
    tmdb_id INTEGER,

    -- Season info
    season_number INTEGER NOT NULL,

    -- Multi-language content
    name TEXT NOT NULL,
    overview TEXT,
    names_i18n JSONB NOT NULL DEFAULT '{}',
    overviews_i18n JSONB NOT NULL DEFAULT '{}',

    -- Media
    poster_path TEXT,

    -- Stats
    episode_count INTEGER NOT NULL DEFAULT 0,
    air_date DATE,

    -- Ratings
    vote_average DECIMAL(4, 2),

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    UNIQUE(series_id, season_number)
);

-- Episodes table
CREATE TABLE tvshow.episodes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    series_id UUID NOT NULL REFERENCES tvshow.series(id) ON DELETE CASCADE,
    season_id UUID NOT NULL REFERENCES tvshow.seasons(id) ON DELETE CASCADE,

    -- External IDs
    tmdb_id INTEGER,
    tvdb_id INTEGER,
    imdb_id TEXT,

    -- Episode info
    season_number INTEGER NOT NULL,
    episode_number INTEGER NOT NULL,

    -- Multi-language content
    title TEXT NOT NULL,
    overview TEXT,
    titles_i18n JSONB NOT NULL DEFAULT '{}',
    overviews_i18n JSONB NOT NULL DEFAULT '{}',

    -- Episode metadata
    air_date DATE,
    runtime INTEGER, -- minutes

    -- Ratings
    vote_average DECIMAL(4, 2),
    vote_count INTEGER,

    -- Media
    still_path TEXT,

    -- Production code (often used by TV networks)
    production_code TEXT,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    UNIQUE(series_id, season_number, episode_number)
);

-- Episode files table
CREATE TABLE tvshow.episode_files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    episode_id UUID NOT NULL REFERENCES tvshow.episodes(id) ON DELETE CASCADE,

    -- File info
    file_path TEXT NOT NULL UNIQUE,
    file_name TEXT NOT NULL,
    file_size BIGINT NOT NULL,

    -- Media info (from FFmpeg/MediaInfo)
    container TEXT,
    resolution TEXT,
    quality_profile TEXT,
    video_codec TEXT,
    audio_codec TEXT,
    bitrate_kbps INTEGER,
    duration_seconds DECIMAL(10, 3),

    -- Audio/subtitle languages
    audio_languages TEXT[] NOT NULL DEFAULT '{}',
    subtitle_languages TEXT[] NOT NULL DEFAULT '{}',

    -- Sonarr integration
    sonarr_file_id INTEGER,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Watch progress for episodes
CREATE TABLE tvshow.episode_watched (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    episode_id UUID NOT NULL REFERENCES tvshow.episodes(id) ON DELETE CASCADE,

    -- Progress tracking
    progress_seconds INTEGER NOT NULL DEFAULT 0,
    duration_seconds INTEGER NOT NULL DEFAULT 0,
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,

    -- Watch history
    watch_count INTEGER NOT NULL DEFAULT 0,
    last_watched_at TIMESTAMPTZ,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    UNIQUE(user_id, episode_id)
);

-- Series genres junction table
CREATE TABLE tvshow.series_genres (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    series_id UUID NOT NULL REFERENCES tvshow.series(id) ON DELETE CASCADE,
    tmdb_genre_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(series_id, tmdb_genre_id)
);

-- Series credits (cast and crew)
CREATE TABLE tvshow.series_credits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    series_id UUID NOT NULL REFERENCES tvshow.series(id) ON DELETE CASCADE,

    -- Person info
    tmdb_person_id INTEGER NOT NULL,
    name TEXT NOT NULL,

    -- Credit type
    credit_type TEXT NOT NULL CHECK (credit_type IN ('cast', 'crew')),

    -- Cast specific
    character TEXT,
    cast_order INTEGER,

    -- Crew specific
    job TEXT,
    department TEXT,

    -- Media
    profile_path TEXT,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Episode credits (guest stars)
CREATE TABLE tvshow.episode_credits (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    episode_id UUID NOT NULL REFERENCES tvshow.episodes(id) ON DELETE CASCADE,

    -- Person info
    tmdb_person_id INTEGER NOT NULL,
    name TEXT NOT NULL,

    -- Credit type
    credit_type TEXT NOT NULL CHECK (credit_type IN ('guest_star', 'crew')),

    -- Guest star specific
    character TEXT,
    cast_order INTEGER,

    -- Crew specific
    job TEXT,
    department TEXT,

    -- Media
    profile_path TEXT,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Networks (e.g., HBO, Netflix, ABC)
CREATE TABLE tvshow.networks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tmdb_id INTEGER UNIQUE NOT NULL,
    name TEXT NOT NULL,
    logo_path TEXT,
    origin_country TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Series-Network junction
CREATE TABLE tvshow.series_networks (
    series_id UUID NOT NULL REFERENCES tvshow.series(id) ON DELETE CASCADE,
    network_id UUID NOT NULL REFERENCES tvshow.networks(id) ON DELETE CASCADE,
    PRIMARY KEY (series_id, network_id)
);

-- Indexes for performance
CREATE INDEX idx_series_tmdb_id ON tvshow.series(tmdb_id);
CREATE INDEX idx_series_tvdb_id ON tvshow.series(tvdb_id);
CREATE INDEX idx_series_sonarr_id ON tvshow.series(sonarr_id) WHERE sonarr_id IS NOT NULL;
CREATE INDEX idx_series_status ON tvshow.series(status);
CREATE INDEX idx_series_first_air_date ON tvshow.series(first_air_date);
CREATE INDEX idx_series_titles_i18n ON tvshow.series USING GIN (titles_i18n);
CREATE INDEX idx_series_created_at ON tvshow.series(created_at DESC);

CREATE INDEX idx_seasons_series ON tvshow.seasons(series_id);
CREATE INDEX idx_seasons_number ON tvshow.seasons(series_id, season_number);

CREATE INDEX idx_episodes_series ON tvshow.episodes(series_id);
CREATE INDEX idx_episodes_season ON tvshow.episodes(season_id);
CREATE INDEX idx_episodes_number ON tvshow.episodes(series_id, season_number, episode_number);
CREATE INDEX idx_episodes_air_date ON tvshow.episodes(air_date);
CREATE INDEX idx_episodes_tmdb_id ON tvshow.episodes(tmdb_id);

CREATE INDEX idx_episode_files_episode ON tvshow.episode_files(episode_id);
CREATE INDEX idx_episode_files_path ON tvshow.episode_files(file_path);
CREATE INDEX idx_episode_files_sonarr ON tvshow.episode_files(sonarr_file_id) WHERE sonarr_file_id IS NOT NULL;

CREATE INDEX idx_episode_watched_user ON tvshow.episode_watched(user_id);
CREATE INDEX idx_episode_watched_episode ON tvshow.episode_watched(episode_id);
CREATE INDEX idx_episode_watched_user_episode ON tvshow.episode_watched(user_id, episode_id);
CREATE INDEX idx_episode_watched_last ON tvshow.episode_watched(user_id, last_watched_at DESC) WHERE NOT is_completed;

CREATE INDEX idx_series_genres_series ON tvshow.series_genres(series_id);
CREATE INDEX idx_series_genres_genre ON tvshow.series_genres(tmdb_genre_id);

CREATE INDEX idx_series_credits_series ON tvshow.series_credits(series_id);
CREATE INDEX idx_series_credits_person ON tvshow.series_credits(tmdb_person_id);

CREATE INDEX idx_episode_credits_episode ON tvshow.episode_credits(episode_id);

-- Trigger for updated_at timestamps
CREATE OR REPLACE FUNCTION tvshow.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_series_updated_at
    BEFORE UPDATE ON tvshow.series
    FOR EACH ROW EXECUTE FUNCTION tvshow.update_updated_at_column();

CREATE TRIGGER update_seasons_updated_at
    BEFORE UPDATE ON tvshow.seasons
    FOR EACH ROW EXECUTE FUNCTION tvshow.update_updated_at_column();

CREATE TRIGGER update_episodes_updated_at
    BEFORE UPDATE ON tvshow.episodes
    FOR EACH ROW EXECUTE FUNCTION tvshow.update_updated_at_column();

CREATE TRIGGER update_episode_files_updated_at
    BEFORE UPDATE ON tvshow.episode_files
    FOR EACH ROW EXECUTE FUNCTION tvshow.update_updated_at_column();

CREATE TRIGGER update_episode_watched_updated_at
    BEFORE UPDATE ON tvshow.episode_watched
    FOR EACH ROW EXECUTE FUNCTION tvshow.update_updated_at_column();

CREATE TRIGGER update_series_credits_updated_at
    BEFORE UPDATE ON tvshow.series_credits
    FOR EACH ROW EXECUTE FUNCTION tvshow.update_updated_at_column();

CREATE TRIGGER update_episode_credits_updated_at
    BEFORE UPDATE ON tvshow.episode_credits
    FOR EACH ROW EXECUTE FUNCTION tvshow.update_updated_at_column();

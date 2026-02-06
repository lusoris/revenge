-- Create movies table in public schema
-- Core movie metadata from TMDb/Radarr

CREATE TABLE IF NOT EXISTS public.movies (
    -- Primary key
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- External IDs
    tmdb_id INTEGER UNIQUE,
    imdb_id TEXT UNIQUE,

    -- Basic Information
    title TEXT NOT NULL,
    original_title TEXT,
    year INTEGER,
    release_date DATE,
    runtime INTEGER, -- in minutes
    overview TEXT,
    tagline TEXT,
    status TEXT, -- released, post-production, in-production, etc.
    original_language TEXT, -- ISO 639-1 code (en, de, fr, etc.)

    -- Images (relative paths from metadata service)
    poster_path TEXT,
    backdrop_path TEXT,
    trailer_url TEXT,

    -- Ratings & Popularity
    vote_average NUMERIC(3, 1), -- e.g., 7.5 (TMDb rating)
    vote_count INTEGER,
    popularity NUMERIC(10, 3),

    -- Financial (TMDb data)
    budget BIGINT,
    revenue BIGINT,

    -- Library Management
    library_added_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata_updated_at TIMESTAMPTZ,
    radarr_id INTEGER, -- Radarr movie ID for sync

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for efficient queries
CREATE UNIQUE INDEX idx_movies_tmdb_id ON public.movies(tmdb_id) WHERE tmdb_id IS NOT NULL;
CREATE UNIQUE INDEX idx_movies_imdb_id ON public.movies(imdb_id) WHERE imdb_id IS NOT NULL;
CREATE INDEX idx_movies_radarr_id ON public.movies(radarr_id) WHERE radarr_id IS NOT NULL;
CREATE INDEX idx_movies_year ON public.movies(year) WHERE year IS NOT NULL;
CREATE INDEX idx_movies_release_date ON public.movies(release_date DESC) WHERE release_date IS NOT NULL;
CREATE INDEX idx_movies_library_added_at ON public.movies(library_added_at DESC);
CREATE INDEX idx_movies_vote_average ON public.movies(vote_average DESC) WHERE vote_average IS NOT NULL;

-- Full-text search index (for title search)
CREATE INDEX idx_movies_title_trgm ON public.movies USING gin(title gin_trgm_ops);
CREATE INDEX idx_movies_original_title_trgm ON public.movies USING gin(original_title gin_trgm_ops) WHERE original_title IS NOT NULL;

-- Trigger for updated_at
CREATE TRIGGER update_movies_updated_at
    BEFORE UPDATE ON public.movies
    FOR EACH ROW
    EXECUTE FUNCTION shared.update_updated_at_column();

-- Comments
COMMENT ON TABLE public.movies IS 'Movie metadata from TMDb/Radarr with library tracking';
COMMENT ON COLUMN public.movies.id IS 'UUID v7 primary key (time-ordered)';
COMMENT ON COLUMN public.movies.tmdb_id IS 'The Movie Database (TMDb) ID - primary metadata source';
COMMENT ON COLUMN public.movies.imdb_id IS 'Internet Movie Database ID (tt1234567 format)';
COMMENT ON COLUMN public.movies.radarr_id IS 'Radarr movie ID for PRIMARY metadata sync';
COMMENT ON COLUMN public.movies.poster_path IS 'Relative path to poster image (from TMDb/Radarr)';
COMMENT ON COLUMN public.movies.backdrop_path IS 'Relative path to backdrop image (from TMDb/Radarr)';
COMMENT ON COLUMN public.movies.library_added_at IS 'When the movie was added to library';
COMMENT ON COLUMN public.movies.metadata_updated_at IS 'Last metadata refresh from TMDb/Radarr';
COMMENT ON COLUMN public.movies.vote_average IS 'TMDb average rating (0-10 scale)';
COMMENT ON COLUMN public.movies.status IS 'Release status (released, post-production, in-production, etc.)';

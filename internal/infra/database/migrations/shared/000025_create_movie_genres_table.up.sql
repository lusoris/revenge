-- Create movie_genres junction table
-- Links movies to genres from TMDb

CREATE TABLE IF NOT EXISTS public.movie_genres (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Foreign key
    movie_id UUID NOT NULL REFERENCES public.movies(id) ON DELETE CASCADE,

    -- Genre information (from TMDb)
    tmdb_genre_id INTEGER NOT NULL,
    name TEXT NOT NULL, -- Action, Comedy, Drama, etc.

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Unique constraint
    CONSTRAINT uq_movie_genre UNIQUE(movie_id, tmdb_genre_id)
);

-- Indexes
CREATE INDEX idx_movie_genres_movie_id ON public.movie_genres(movie_id);
CREATE INDEX idx_movie_genres_genre_id ON public.movie_genres(tmdb_genre_id);
CREATE INDEX idx_movie_genres_name ON public.movie_genres(name);

-- Comments
COMMENT ON TABLE public.movie_genres IS 'Junction table linking movies to TMDb genres';
COMMENT ON COLUMN public.movie_genres.tmdb_genre_id IS 'TMDb genre ID (28=Action, 35=Comedy, etc.)';
COMMENT ON COLUMN public.movie_genres.name IS 'Genre name for display (Action, Comedy, Drama, etc.)';

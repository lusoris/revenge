-- Create movie collections tables
-- Collections group related movies (e.g., Marvel Cinematic Universe, Star Wars)

-- Collections (e.g., "Marvel Cinematic Universe", "The Lord of the Rings")
CREATE TABLE IF NOT EXISTS public.movie_collections (
    -- Primary key
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- External IDs
    tmdb_collection_id INTEGER UNIQUE,

    -- Collection Information
    name TEXT NOT NULL,
    overview TEXT,
    poster_path TEXT,
    backdrop_path TEXT,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Junction table for movies in collections
CREATE TABLE IF NOT EXISTS public.movie_collection_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Foreign keys
    collection_id UUID NOT NULL REFERENCES public.movie_collections(id) ON DELETE CASCADE,
    movie_id UUID NOT NULL REFERENCES public.movies(id) ON DELETE CASCADE,

    -- Order in collection (1, 2, 3, etc. for sequels)
    collection_order INTEGER,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Unique constraint
    CONSTRAINT uq_collection_movie UNIQUE(collection_id, movie_id)
);

-- Indexes for movie_collections
CREATE INDEX idx_collections_tmdb_id ON public.movie_collections(tmdb_collection_id)
    WHERE tmdb_collection_id IS NOT NULL;
CREATE INDEX idx_collections_name_trgm ON public.movie_collections USING gin(name gin_trgm_ops);

-- Indexes for movie_collection_members
CREATE INDEX idx_collection_members_collection ON public.movie_collection_members(collection_id);
CREATE INDEX idx_collection_members_movie ON public.movie_collection_members(movie_id);
CREATE INDEX idx_collection_members_order ON public.movie_collection_members(collection_id, collection_order)
    WHERE collection_order IS NOT NULL;

-- Triggers
CREATE TRIGGER update_movie_collections_updated_at
    BEFORE UPDATE ON public.movie_collections
    FOR EACH ROW
    EXECUTE FUNCTION shared.update_updated_at_column();

-- Comments
COMMENT ON TABLE public.movie_collections IS 'Movie collections from TMDb (e.g., MCU, Star Wars)';
COMMENT ON TABLE public.movie_collection_members IS 'Junction table linking movies to collections';
COMMENT ON COLUMN public.movie_collections.tmdb_collection_id IS 'TMDb collection ID for metadata sync';
COMMENT ON COLUMN public.movie_collection_members.collection_order IS 'Order in collection (1 = first movie, 2 = sequel, etc.)';

-- Create movie_credits table for cast and crew information
-- Stores actors, directors, writers, etc. from TMDb

CREATE TABLE IF NOT EXISTS public.movie_credits (
    -- Primary key
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Foreign key
    movie_id UUID NOT NULL REFERENCES public.movies(id) ON DELETE CASCADE,

    -- Person Information (from TMDb)
    tmdb_person_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    profile_path TEXT, -- Relative path to person's photo

    -- Credit Details
    credit_type TEXT NOT NULL, -- 'cast' or 'crew'

    -- Cast-specific fields
    character TEXT, -- Character name for actors
    cast_order INTEGER, -- Order in cast list (0 = lead)

    -- Crew-specific fields
    job TEXT, -- Director, Writer, Producer, Cinematographer, etc.
    department TEXT, -- Directing, Writing, Production, Camera, etc.

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT chk_credit_type CHECK (credit_type IN ('cast', 'crew')),
    CONSTRAINT chk_cast_fields CHECK (
        (credit_type = 'cast' AND character IS NOT NULL) OR
        (credit_type = 'crew' AND job IS NOT NULL AND department IS NOT NULL)
    )
);

-- Indexes
CREATE INDEX idx_movie_credits_movie_id ON public.movie_credits(movie_id);
CREATE INDEX idx_movie_credits_person_id ON public.movie_credits(tmdb_person_id);
CREATE INDEX idx_movie_credits_type ON public.movie_credits(credit_type);
CREATE INDEX idx_movie_credits_cast_order ON public.movie_credits(movie_id, cast_order)
    WHERE credit_type = 'cast' AND cast_order IS NOT NULL;
CREATE INDEX idx_movie_credits_job ON public.movie_credits(job)
    WHERE credit_type = 'crew' AND job IS NOT NULL;

-- Composite index for efficient lookups
CREATE UNIQUE INDEX idx_movie_credits_unique ON public.movie_credits(
    movie_id, tmdb_person_id, credit_type, COALESCE(character, ''), COALESCE(job, '')
);

-- Trigger for updated_at
CREATE TRIGGER update_movie_credits_updated_at
    BEFORE UPDATE ON public.movie_credits
    FOR EACH ROW
    EXECUTE FUNCTION shared.update_updated_at_column();

-- Comments
COMMENT ON TABLE public.movie_credits IS 'Cast and crew information from TMDb';
COMMENT ON COLUMN public.movie_credits.credit_type IS 'Either cast (actor) or crew (director, writer, etc.)';
COMMENT ON COLUMN public.movie_credits.character IS 'Character name for cast members';
COMMENT ON COLUMN public.movie_credits.cast_order IS 'Order in cast list (0 = lead actor)';
COMMENT ON COLUMN public.movie_credits.job IS 'Job title for crew members (Director, Writer, etc.)';
COMMENT ON COLUMN public.movie_credits.department IS 'Department for crew members (Directing, Writing, etc.)';
COMMENT ON COLUMN public.movie_credits.tmdb_person_id IS 'TMDb person ID for linking to person data';

-- Create movie_watched table for user watch tracking
-- Tracks which users have watched which movies and their progress

CREATE TABLE IF NOT EXISTS public.movie_watched (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Foreign keys
    user_id UUID NOT NULL REFERENCES shared.users(id) ON DELETE CASCADE,
    movie_id UUID NOT NULL REFERENCES public.movies(id) ON DELETE CASCADE,

    -- Watch Progress
    progress_seconds INTEGER NOT NULL DEFAULT 0, -- Current position in seconds
    duration_seconds INTEGER, -- Total duration (from file)
    progress_percent NUMERIC(5, 2) GENERATED ALWAYS AS (
        CASE
            WHEN duration_seconds > 0
            THEN (progress_seconds::NUMERIC / duration_seconds * 100)
            ELSE 0
        END
    ) STORED,

    -- Watch Status
    is_completed BOOLEAN DEFAULT FALSE, -- Watched to end (>90% progress)
    last_watched_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ, -- When user finished watching
    watch_count INTEGER DEFAULT 1, -- Number of times watched

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Unique constraint - one watch record per user per movie
    CONSTRAINT uq_user_movie_watched UNIQUE(user_id, movie_id)
);

-- Indexes
CREATE INDEX idx_movie_watched_user_id ON public.movie_watched(user_id);
CREATE INDEX idx_movie_watched_movie_id ON public.movie_watched(movie_id);
CREATE INDEX idx_movie_watched_last_watched ON public.movie_watched(user_id, last_watched_at DESC);
CREATE INDEX idx_movie_watched_completed ON public.movie_watched(user_id, is_completed, last_watched_at DESC);
CREATE INDEX idx_movie_watched_in_progress ON public.movie_watched(user_id, progress_percent)
    WHERE is_completed = FALSE AND progress_percent > 0;

-- Trigger for updated_at
CREATE TRIGGER update_movie_watched_updated_at
    BEFORE UPDATE ON public.movie_watched
    FOR EACH ROW
    EXECUTE FUNCTION shared.update_updated_at_column();

-- Comments
COMMENT ON TABLE public.movie_watched IS 'User watch history and progress tracking for movies';
COMMENT ON COLUMN public.movie_watched.progress_seconds IS 'Current playback position in seconds';
COMMENT ON COLUMN public.movie_watched.progress_percent IS 'Calculated progress percentage (0-100)';
COMMENT ON COLUMN public.movie_watched.is_completed IS 'TRUE when user has watched >90% of the movie';
COMMENT ON COLUMN public.movie_watched.watch_count IS 'Number of times user has watched this movie';
COMMENT ON COLUMN public.movie_watched.last_watched_at IS 'Most recent watch time (for continue watching)';
COMMENT ON COLUMN public.movie_watched.completed_at IS 'Timestamp when user completed watching';

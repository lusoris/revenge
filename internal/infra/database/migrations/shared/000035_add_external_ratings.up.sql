-- Add external_ratings JSONB column to movies and series tables
-- Stores ratings from external providers (IMDb, Rotten Tomatoes, Metacritic, TMDb, etc.)
-- Format: [{"source": "Internet Movie Database", "value": "8.8/10", "score": 88.0}, ...]

-- Movies
ALTER TABLE public.movies
ADD COLUMN IF NOT EXISTS external_ratings JSONB NOT NULL DEFAULT '[]';

-- TV Shows
ALTER TABLE tvshow.series
ADD COLUMN IF NOT EXISTS external_ratings JSONB NOT NULL DEFAULT '[]';

-- Comments
COMMENT ON COLUMN public.movies.external_ratings IS 'External ratings from various providers (IMDb, RT, Metacritic, etc.) as JSON array';

COMMENT ON COLUMN tvshow.series.external_ratings IS 'External ratings from various providers (IMDb, RT, Metacritic, etc.) as JSON array';

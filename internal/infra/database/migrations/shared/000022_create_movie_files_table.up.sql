-- Create movie_files table for physical file tracking
-- Links movies to actual media files on disk

CREATE TABLE IF NOT EXISTS public.movie_files (
    -- Primary key
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Foreign key
    movie_id UUID NOT NULL REFERENCES public.movies(id) ON DELETE CASCADE,

    -- File Information
    file_path TEXT NOT NULL UNIQUE, -- Full path on disk
    file_size BIGINT NOT NULL, -- Size in bytes
    file_name TEXT NOT NULL,

    -- Media Information
    resolution TEXT, -- 1080p, 2160p, 720p, etc.
    quality_profile TEXT, -- Bluray-1080p, WEB-DL-2160p, etc.
    video_codec TEXT, -- h264, h265, av1, etc.
    audio_codec TEXT, -- aac, ac3, dts, etc.
    container TEXT, -- mkv, mp4, avi, etc.

    -- Technical Details
    duration_seconds INTEGER, -- Parsed from file
    bitrate_kbps INTEGER,
    framerate NUMERIC(5, 2), -- 23.976, 24.000, 25.000, etc.

    -- HDR/Color
    dynamic_range TEXT, -- SDR, HDR, HDR10, HDR10+, Dolby Vision
    color_space TEXT, -- BT.709, BT.2020, etc.

    -- Audio Tracks
    audio_channels TEXT, -- 2.0, 5.1, 7.1, etc.
    audio_languages TEXT[], -- Array of ISO 639-1 codes

    -- Subtitles
    subtitle_languages TEXT[], -- Array of ISO 639-1 codes

    -- Library Management
    radarr_file_id INTEGER, -- Radarr's file ID for sync
    last_scanned_at TIMESTAMPTZ,
    is_monitored BOOLEAN DEFAULT TRUE,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_movie_files_movie_id ON public.movie_files(movie_id);
CREATE INDEX idx_movie_files_radarr_id ON public.movie_files(radarr_file_id) WHERE radarr_file_id IS NOT NULL;
CREATE INDEX idx_movie_files_resolution ON public.movie_files(resolution) WHERE resolution IS NOT NULL;
CREATE INDEX idx_movie_files_quality ON public.movie_files(quality_profile) WHERE quality_profile IS NOT NULL;
CREATE INDEX idx_movie_files_last_scanned ON public.movie_files(last_scanned_at DESC);

-- Trigger for updated_at
CREATE TRIGGER update_movie_files_updated_at
    BEFORE UPDATE ON public.movie_files
    FOR EACH ROW
    EXECUTE FUNCTION shared.update_updated_at_column();

-- Comments
COMMENT ON TABLE public.movie_files IS 'Physical media files associated with movies';
COMMENT ON COLUMN public.movie_files.file_path IS 'Absolute path to file on disk';
COMMENT ON COLUMN public.movie_files.radarr_file_id IS 'Radarr file ID for sync with PRIMARY metadata source';
COMMENT ON COLUMN public.movie_files.quality_profile IS 'Quality profile from Radarr (e.g., Bluray-1080p)';
COMMENT ON COLUMN public.movie_files.duration_seconds IS 'Actual file duration (may differ from movie runtime)';
COMMENT ON COLUMN public.movie_files.audio_languages IS 'Array of audio track language codes';
COMMENT ON COLUMN public.movie_files.subtitle_languages IS 'Array of subtitle language codes';

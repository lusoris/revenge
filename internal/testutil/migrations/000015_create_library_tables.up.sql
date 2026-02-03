-- Migration: 000015_create_library_tables
-- Description: Create libraries, library_scans, and library_permissions tables
-- Schema: public

-- Libraries (top-level organizational unit for media)
CREATE TABLE IF NOT EXISTS public.libraries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Library details
    name TEXT NOT NULL,
    type VARCHAR(50) NOT NULL,                     -- 'movie', 'tvshow', 'music', 'photo', 'book', 'audiobook', 'comic', 'podcast', 'adult'

    -- File system paths (array of paths)
    paths TEXT[] NOT NULL DEFAULT '{}',

    -- Settings
    enabled BOOLEAN NOT NULL DEFAULT true,
    scan_on_startup BOOLEAN NOT NULL DEFAULT false,
    realtime_monitoring BOOLEAN NOT NULL DEFAULT true,

    -- Metadata settings
    metadata_provider VARCHAR(50),                 -- 'tmdb', 'tvdb', 'musicbrainz', 'openlib', etc.
    preferred_language VARCHAR(10) NOT NULL DEFAULT 'en',

    -- Scanner settings (type-specific configuration)
    scanner_config JSONB,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Indexes for libraries
CREATE INDEX IF NOT EXISTS idx_libraries_type ON public.libraries(type);
CREATE INDEX IF NOT EXISTS idx_libraries_enabled ON public.libraries(enabled) WHERE enabled = true;
CREATE INDEX IF NOT EXISTS idx_libraries_name ON public.libraries(name);

-- Comments
COMMENT ON TABLE public.libraries IS 'Media libraries organizing content by type and location';
COMMENT ON COLUMN public.libraries.type IS 'Library content type: movie, tvshow, music, photo, book, audiobook, comic, podcast, adult';
COMMENT ON COLUMN public.libraries.paths IS 'Array of file system paths to scan for this library';
COMMENT ON COLUMN public.libraries.metadata_provider IS 'Primary metadata provider: tmdb, tvdb, musicbrainz, openlib, etc.';
COMMENT ON COLUMN public.libraries.scanner_config IS 'Type-specific scanner configuration as JSONB';

-- Library scan jobs (track scanning progress)
CREATE TABLE IF NOT EXISTS public.library_scans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    library_id UUID NOT NULL REFERENCES public.libraries(id) ON DELETE CASCADE,

    -- Scan details
    scan_type VARCHAR(20) NOT NULL,                -- 'full', 'incremental', 'metadata'
    status VARCHAR(20) NOT NULL DEFAULT 'pending', -- 'pending', 'running', 'completed', 'failed', 'cancelled'

    -- Progress tracking
    items_scanned INTEGER NOT NULL DEFAULT 0,
    items_added INTEGER NOT NULL DEFAULT 0,
    items_updated INTEGER NOT NULL DEFAULT 0,
    items_removed INTEGER NOT NULL DEFAULT 0,
    errors_count INTEGER NOT NULL DEFAULT 0,

    -- Error details
    error_message TEXT,

    -- Timing
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    duration_seconds INTEGER,

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Indexes for library_scans
CREATE INDEX IF NOT EXISTS idx_library_scans_library ON public.library_scans(library_id, started_at DESC);
CREATE INDEX IF NOT EXISTS idx_library_scans_status ON public.library_scans(status);
CREATE INDEX IF NOT EXISTS idx_library_scans_created ON public.library_scans(created_at DESC);

-- Comments
COMMENT ON TABLE public.library_scans IS 'Track library scan jobs and their progress';
COMMENT ON COLUMN public.library_scans.scan_type IS 'Type of scan: full, incremental, metadata';
COMMENT ON COLUMN public.library_scans.status IS 'Scan status: pending, running, completed, failed, cancelled';

-- Library permissions (user access control to libraries)
CREATE TABLE IF NOT EXISTS public.library_permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    library_id UUID NOT NULL REFERENCES public.libraries(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES shared.users(id) ON DELETE CASCADE,

    -- Permission type
    permission VARCHAR(50) NOT NULL,               -- 'view', 'download', 'manage'

    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Indexes for library_permissions
CREATE UNIQUE INDEX IF NOT EXISTS idx_library_perms_unique ON public.library_permissions(library_id, user_id, permission);
CREATE INDEX IF NOT EXISTS idx_library_perms_user ON public.library_permissions(user_id);
CREATE INDEX IF NOT EXISTS idx_library_perms_library ON public.library_permissions(library_id);

-- Comments
COMMENT ON TABLE public.library_permissions IS 'Per-user access permissions to libraries';
COMMENT ON COLUMN public.library_permissions.permission IS 'Permission type: view, download, manage';

-- Trigger to update updated_at on libraries
CREATE OR REPLACE FUNCTION update_libraries_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_libraries_updated_at
    BEFORE UPDATE ON public.libraries
    FOR EACH ROW
    EXECUTE FUNCTION update_libraries_updated_at();

-- Migration: 000015_create_library_tables (down)
-- Description: Drop libraries, library_scans, and library_permissions tables

-- Drop trigger first
DROP TRIGGER IF EXISTS trigger_libraries_updated_at ON public.libraries;
DROP FUNCTION IF EXISTS update_libraries_updated_at();

-- Drop tables in reverse order (dependencies first)
DROP TABLE IF EXISTS public.library_permissions;
DROP TABLE IF EXISTS public.library_scans;
DROP TABLE IF EXISTS public.libraries;

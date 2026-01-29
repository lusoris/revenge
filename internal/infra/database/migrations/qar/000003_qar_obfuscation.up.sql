-- QAR Full Obfuscation: Queen Anne's Revenge themed naming
-- Renames all adult content tables to pirate-themed obfuscated names
BEGIN;

-- ============================================================================
-- FLEETS (Libraries)
-- ============================================================================

CREATE TABLE qar.fleets (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                    VARCHAR(255) NOT NULL,
    fleet_type              VARCHAR(20) NOT NULL CHECK (fleet_type IN ('expedition', 'voyage')),
    paths                   TEXT[] NOT NULL,

    -- Adult-specific settings
    stashdb_endpoint        TEXT DEFAULT 'https://stashdb.org/graphql',
    tpdb_enabled            BOOLEAN NOT NULL DEFAULT true,
    whisparr_sync           BOOLEAN NOT NULL DEFAULT false,
    auto_tag_crew           BOOLEAN NOT NULL DEFAULT true,
    fingerprint_on_scan     BOOLEAN NOT NULL DEFAULT true,

    -- Access control
    owner_user_id           UUID REFERENCES users(id) ON DELETE SET NULL,

    -- Timestamps
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_qar_fleets_owner ON qar.fleets(owner_user_id) WHERE owner_user_id IS NOT NULL;

CREATE TRIGGER qar_fleets_updated_at
    BEFORE UPDATE ON qar.fleets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- ============================================================================
-- PORTS (Studios)
-- ============================================================================

-- Rename studios -> ports
ALTER TABLE qar.studios RENAME TO ports;

-- ============================================================================
-- FLAGS (Tags)
-- ============================================================================

-- Rename tags -> flags
ALTER TABLE qar.tags RENAME TO flags;

-- Add waters (category) column if not exists
ALTER TABLE qar.flags ADD COLUMN IF NOT EXISTS waters VARCHAR(50);

-- ============================================================================
-- CREW (Performers)
-- ============================================================================

-- Rename performers -> crew with obfuscated columns
ALTER TABLE qar.performers RENAME TO crew;

-- Rename columns to obfuscated versions
ALTER TABLE qar.crew RENAME COLUMN birthdate TO christening;
ALTER TABLE qar.crew RENAME COLUMN ethnicity TO origin;
ALTER TABLE qar.crew RENAME COLUMN career_start TO maiden_voyage;
ALTER TABLE qar.crew RENAME COLUMN career_end TO last_port;
ALTER TABLE qar.crew RENAME COLUMN hair_color TO rigging;
ALTER TABLE qar.crew RENAME COLUMN eye_color TO compass;
ALTER TABLE qar.crew RENAME COLUMN tattoos TO markings;
ALTER TABLE qar.crew RENAME COLUMN piercings TO anchors;
ALTER TABLE qar.crew RENAME COLUMN stashdb_id TO charter;
ALTER TABLE qar.crew RENAME COLUMN tpdb_id TO registry;
ALTER TABLE qar.crew RENAME COLUMN freeones_id TO manifest;

-- Rename performer_aliases -> crew_names
ALTER TABLE qar.performer_aliases RENAME TO crew_names;
ALTER TABLE qar.crew_names RENAME COLUMN performer_id TO crew_id;
ALTER TABLE qar.crew_names RENAME COLUMN alias TO name;

-- Rename performer_images -> crew_portraits
ALTER TABLE qar.performer_images RENAME TO crew_portraits;
ALTER TABLE qar.crew_portraits RENAME COLUMN performer_id TO crew_id;

-- ============================================================================
-- EXPEDITIONS (Movies/Full-length adult content)
-- ============================================================================

-- Rename movies -> expeditions
ALTER TABLE qar.movies RENAME TO expeditions;

-- Update to use fleet_id instead of library_id
ALTER TABLE qar.expeditions ADD COLUMN fleet_id UUID REFERENCES qar.fleets(id) ON DELETE CASCADE;
ALTER TABLE qar.expeditions RENAME COLUMN studio_id TO port_id;
ALTER TABLE qar.expeditions RENAME COLUMN stashdb_id TO charter;
ALTER TABLE qar.expeditions RENAME COLUMN tpdb_id TO registry;
ALTER TABLE qar.expeditions RENAME COLUMN release_date TO launch_date;
ALTER TABLE qar.expeditions RENAME COLUMN phash TO coordinates;

-- Rename movie_performers -> expedition_crew
ALTER TABLE qar.movie_performers RENAME TO expedition_crew;
ALTER TABLE qar.expedition_crew RENAME COLUMN movie_id TO expedition_id;
ALTER TABLE qar.expedition_crew RENAME COLUMN performer_id TO crew_id;

-- Rename movie_tags -> expedition_flags
ALTER TABLE qar.movie_tags RENAME TO expedition_flags;
ALTER TABLE qar.expedition_flags RENAME COLUMN movie_id TO expedition_id;
ALTER TABLE qar.expedition_flags RENAME COLUMN tag_id TO flag_id;

-- Rename movie_images -> expedition_charts (images)
ALTER TABLE qar.movie_images RENAME TO expedition_charts;
ALTER TABLE qar.expedition_charts RENAME COLUMN movie_id TO expedition_id;

-- ============================================================================
-- VOYAGES (Scenes)
-- ============================================================================

-- Rename scenes -> voyages
ALTER TABLE qar.scenes RENAME TO voyages;

-- Update to use fleet_id instead of library_id
ALTER TABLE qar.voyages ADD COLUMN fleet_id UUID REFERENCES qar.fleets(id) ON DELETE CASCADE;
ALTER TABLE qar.voyages RENAME COLUMN studio_id TO port_id;
ALTER TABLE qar.voyages RENAME COLUMN stashdb_id TO charter;
ALTER TABLE qar.voyages RENAME COLUMN tpdb_id TO registry;
ALTER TABLE qar.voyages RENAME COLUMN release_date TO launch_date;
ALTER TABLE qar.voyages RENAME COLUMN runtime_minutes TO distance;
ALTER TABLE qar.voyages RENAME COLUMN phash TO coordinates;

-- Rename scene_performers -> voyage_crew
ALTER TABLE qar.scene_performers RENAME TO voyage_crew;
ALTER TABLE qar.voyage_crew RENAME COLUMN scene_id TO voyage_id;
ALTER TABLE qar.voyage_crew RENAME COLUMN performer_id TO crew_id;

-- Rename scene_tags -> voyage_flags
ALTER TABLE qar.scene_tags RENAME TO voyage_flags;
ALTER TABLE qar.voyage_flags RENAME COLUMN scene_id TO voyage_id;
ALTER TABLE qar.voyage_flags RENAME COLUMN tag_id TO flag_id;

-- Rename scene_markers -> voyage_waypoints
ALTER TABLE qar.scene_markers RENAME TO voyage_waypoints;
ALTER TABLE qar.voyage_waypoints RENAME COLUMN scene_id TO voyage_id;
ALTER TABLE qar.voyage_waypoints RENAME COLUMN tag_id TO flag_id;

-- ============================================================================
-- TREASURES (Galleries/Images)
-- ============================================================================

-- Rename galleries -> treasures
ALTER TABLE qar.galleries RENAME TO treasures;
ALTER TABLE qar.treasures RENAME COLUMN movie_id TO expedition_id;

-- Rename gallery_images -> treasure_maps
ALTER TABLE qar.gallery_images RENAME TO treasure_maps;
ALTER TABLE qar.treasure_maps RENAME COLUMN gallery_id TO treasure_id;
ALTER TABLE qar.treasure_maps RENAME COLUMN sort_order TO bearing;

-- ============================================================================
-- USER DATA (Obfuscated)
-- ============================================================================

-- Rename user_ratings -> crew_bounties (user ratings)
ALTER TABLE qar.user_ratings RENAME TO bounties;
ALTER TABLE qar.bounties RENAME COLUMN movie_id TO expedition_id;
ALTER TABLE qar.bounties RENAME COLUMN rating TO reward;

-- Rename user_favorites -> marked_charts
ALTER TABLE qar.user_favorites RENAME TO marked_charts;
ALTER TABLE qar.marked_charts RENAME COLUMN movie_id TO expedition_id;

-- Rename watch_history -> voyage_log
ALTER TABLE qar.watch_history RENAME TO voyage_log;
ALTER TABLE qar.voyage_log RENAME COLUMN movie_id TO expedition_id;
ALTER TABLE qar.voyage_log RENAME COLUMN position_ticks TO bearing_ticks;

-- Rename user_scene_data -> voyage_records
ALTER TABLE qar.user_scene_data RENAME TO voyage_records;
ALTER TABLE qar.voyage_records RENAME COLUMN scene_id TO voyage_id;
ALTER TABLE qar.voyage_records RENAME COLUMN position_ms TO bearing_ms;
ALTER TABLE qar.voyage_records RENAME COLUMN watch_count TO crossings;
ALTER TABLE qar.voyage_records RENAME COLUMN rating TO bounty;
ALTER TABLE qar.voyage_records RENAME COLUMN is_favorite TO marked;

-- Rename user_performer_favorites -> crew_roster (favorite performers)
ALTER TABLE qar.user_performer_favorites RENAME TO crew_roster;
ALTER TABLE qar.crew_roster RENAME COLUMN performer_id TO crew_id;
ALTER TABLE qar.crew_roster RENAME COLUMN added_at TO enlisted_at;

-- ============================================================================
-- UPDATE INDEXES (Rename to match new table names)
-- ============================================================================

-- Crew indexes
ALTER INDEX IF EXISTS idx_qar_performers_name RENAME TO idx_qar_crew_name;
ALTER INDEX IF EXISTS idx_qar_performers_stashdb RENAME TO idx_qar_crew_charter;

-- Expedition indexes
ALTER INDEX IF EXISTS idx_qar_movies_library RENAME TO idx_qar_expeditions_fleet;
ALTER INDEX IF EXISTS idx_qar_movies_studio RENAME TO idx_qar_expeditions_port;
ALTER INDEX IF EXISTS idx_qar_movies_phash RENAME TO idx_qar_expeditions_coordinates;
ALTER INDEX IF EXISTS idx_qar_movies_oshash RENAME TO idx_qar_expeditions_oshash;

-- Voyage indexes
ALTER INDEX IF EXISTS idx_qar_scenes_library RENAME TO idx_qar_voyages_fleet;
ALTER INDEX IF EXISTS idx_qar_scenes_studio RENAME TO idx_qar_voyages_port;
ALTER INDEX IF EXISTS idx_qar_scenes_oshash RENAME TO idx_qar_voyages_oshash;
ALTER INDEX IF EXISTS idx_qar_scenes_stashdb RENAME TO idx_qar_voyages_charter;

-- Voyage waypoints
ALTER INDEX IF EXISTS idx_qar_markers_scene RENAME TO idx_qar_waypoints_voyage;

-- User data
ALTER INDEX IF EXISTS idx_qar_watch_history_user RENAME TO idx_qar_voyage_log_user;

-- ============================================================================
-- UPDATE TRIGGERS (Rename to match new table names)
-- ============================================================================

-- Drop old triggers and recreate with new names
DROP TRIGGER IF EXISTS c_studios_updated_at ON qar.ports;
CREATE TRIGGER qar_ports_updated_at
    BEFORE UPDATE ON qar.ports
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

DROP TRIGGER IF EXISTS c_performers_updated_at ON qar.crew;
CREATE TRIGGER qar_crew_updated_at
    BEFORE UPDATE ON qar.crew
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

DROP TRIGGER IF EXISTS c_movies_updated_at ON qar.expeditions;
CREATE TRIGGER qar_expeditions_updated_at
    BEFORE UPDATE ON qar.expeditions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

DROP TRIGGER IF EXISTS c_scenes_updated_at ON qar.voyages;
CREATE TRIGGER qar_voyages_updated_at
    BEFORE UPDATE ON qar.voyages
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

DROP TRIGGER IF EXISTS c_user_ratings_updated_at ON qar.bounties;
CREATE TRIGGER qar_bounties_updated_at
    BEFORE UPDATE ON qar.bounties
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

COMMIT;

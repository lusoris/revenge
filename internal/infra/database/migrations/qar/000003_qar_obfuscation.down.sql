-- Rollback QAR Full Obfuscation
-- Reverts all pirate-themed naming back to original names
BEGIN;

-- ============================================================================
-- REVERT TRIGGERS
-- ============================================================================

DROP TRIGGER IF EXISTS qar_ports_updated_at ON qar.ports;
DROP TRIGGER IF EXISTS qar_crew_updated_at ON qar.crew;
DROP TRIGGER IF EXISTS qar_expeditions_updated_at ON qar.expeditions;
DROP TRIGGER IF EXISTS qar_voyages_updated_at ON qar.voyages;
DROP TRIGGER IF EXISTS qar_bounties_updated_at ON qar.bounties;
DROP TRIGGER IF EXISTS qar_fleets_updated_at ON qar.fleets;

-- ============================================================================
-- REVERT USER DATA
-- ============================================================================

ALTER TABLE qar.crew_roster RENAME COLUMN enlisted_at TO added_at;
ALTER TABLE qar.crew_roster RENAME COLUMN crew_id TO performer_id;
ALTER TABLE qar.crew_roster RENAME TO user_performer_favorites;

ALTER TABLE qar.voyage_records RENAME COLUMN marked TO is_favorite;
ALTER TABLE qar.voyage_records RENAME COLUMN bounty TO rating;
ALTER TABLE qar.voyage_records RENAME COLUMN crossings TO watch_count;
ALTER TABLE qar.voyage_records RENAME COLUMN bearing_ms TO position_ms;
ALTER TABLE qar.voyage_records RENAME COLUMN voyage_id TO scene_id;
ALTER TABLE qar.voyage_records RENAME TO user_scene_data;

ALTER TABLE qar.voyage_log RENAME COLUMN bearing_ticks TO position_ticks;
ALTER TABLE qar.voyage_log RENAME COLUMN expedition_id TO movie_id;
ALTER TABLE qar.voyage_log RENAME TO watch_history;

ALTER TABLE qar.marked_charts RENAME COLUMN expedition_id TO movie_id;
ALTER TABLE qar.marked_charts RENAME TO user_favorites;

ALTER TABLE qar.bounties RENAME COLUMN reward TO rating;
ALTER TABLE qar.bounties RENAME COLUMN expedition_id TO movie_id;
ALTER TABLE qar.bounties RENAME TO user_ratings;

-- ============================================================================
-- REVERT TREASURES
-- ============================================================================

ALTER TABLE qar.treasure_maps RENAME COLUMN bearing TO sort_order;
ALTER TABLE qar.treasure_maps RENAME COLUMN treasure_id TO gallery_id;
ALTER TABLE qar.treasure_maps RENAME TO gallery_images;

ALTER TABLE qar.treasures RENAME COLUMN expedition_id TO movie_id;
ALTER TABLE qar.treasures RENAME TO galleries;

-- ============================================================================
-- REVERT VOYAGES (Scenes)
-- ============================================================================

ALTER TABLE qar.voyage_waypoints RENAME COLUMN flag_id TO tag_id;
ALTER TABLE qar.voyage_waypoints RENAME COLUMN voyage_id TO scene_id;
ALTER TABLE qar.voyage_waypoints RENAME TO scene_markers;

ALTER TABLE qar.voyage_flags RENAME COLUMN flag_id TO tag_id;
ALTER TABLE qar.voyage_flags RENAME COLUMN voyage_id TO scene_id;
ALTER TABLE qar.voyage_flags RENAME TO scene_tags;

ALTER TABLE qar.voyage_crew RENAME COLUMN crew_id TO performer_id;
ALTER TABLE qar.voyage_crew RENAME COLUMN voyage_id TO scene_id;
ALTER TABLE qar.voyage_crew RENAME TO scene_performers;

ALTER TABLE qar.voyages RENAME COLUMN coordinates TO phash;
ALTER TABLE qar.voyages RENAME COLUMN distance TO runtime_minutes;
ALTER TABLE qar.voyages RENAME COLUMN launch_date TO release_date;
ALTER TABLE qar.voyages RENAME COLUMN registry TO tpdb_id;
ALTER TABLE qar.voyages RENAME COLUMN charter TO stashdb_id;
ALTER TABLE qar.voyages RENAME COLUMN port_id TO studio_id;
ALTER TABLE qar.voyages DROP COLUMN IF EXISTS fleet_id;
ALTER TABLE qar.voyages RENAME TO scenes;

-- ============================================================================
-- REVERT EXPEDITIONS (Movies)
-- ============================================================================

ALTER TABLE qar.expedition_charts RENAME COLUMN expedition_id TO movie_id;
ALTER TABLE qar.expedition_charts RENAME TO movie_images;

ALTER TABLE qar.expedition_flags RENAME COLUMN flag_id TO tag_id;
ALTER TABLE qar.expedition_flags RENAME COLUMN expedition_id TO movie_id;
ALTER TABLE qar.expedition_flags RENAME TO movie_tags;

ALTER TABLE qar.expedition_crew RENAME COLUMN crew_id TO performer_id;
ALTER TABLE qar.expedition_crew RENAME COLUMN expedition_id TO movie_id;
ALTER TABLE qar.expedition_crew RENAME TO movie_performers;

ALTER TABLE qar.expeditions RENAME COLUMN coordinates TO phash;
ALTER TABLE qar.expeditions RENAME COLUMN launch_date TO release_date;
ALTER TABLE qar.expeditions RENAME COLUMN registry TO tpdb_id;
ALTER TABLE qar.expeditions RENAME COLUMN charter TO stashdb_id;
ALTER TABLE qar.expeditions RENAME COLUMN port_id TO studio_id;
ALTER TABLE qar.expeditions DROP COLUMN IF EXISTS fleet_id;
ALTER TABLE qar.expeditions RENAME TO movies;

-- ============================================================================
-- REVERT CREW (Performers)
-- ============================================================================

ALTER TABLE qar.crew_portraits RENAME COLUMN crew_id TO performer_id;
ALTER TABLE qar.crew_portraits RENAME TO performer_images;

ALTER TABLE qar.crew_names RENAME COLUMN name TO alias;
ALTER TABLE qar.crew_names RENAME COLUMN crew_id TO performer_id;
ALTER TABLE qar.crew_names RENAME TO performer_aliases;

ALTER TABLE qar.crew RENAME COLUMN manifest TO freeones_id;
ALTER TABLE qar.crew RENAME COLUMN registry TO tpdb_id;
ALTER TABLE qar.crew RENAME COLUMN charter TO stashdb_id;
ALTER TABLE qar.crew RENAME COLUMN anchors TO piercings;
ALTER TABLE qar.crew RENAME COLUMN markings TO tattoos;
ALTER TABLE qar.crew RENAME COLUMN compass TO eye_color;
ALTER TABLE qar.crew RENAME COLUMN rigging TO hair_color;
ALTER TABLE qar.crew RENAME COLUMN last_port TO career_end;
ALTER TABLE qar.crew RENAME COLUMN maiden_voyage TO career_start;
ALTER TABLE qar.crew RENAME COLUMN origin TO ethnicity;
ALTER TABLE qar.crew RENAME COLUMN christening TO birthdate;
ALTER TABLE qar.crew RENAME TO performers;

-- ============================================================================
-- REVERT FLAGS (Tags)
-- ============================================================================

ALTER TABLE qar.flags DROP COLUMN IF EXISTS waters;
ALTER TABLE qar.flags RENAME TO tags;

-- ============================================================================
-- REVERT PORTS (Studios)
-- ============================================================================

ALTER TABLE qar.ports RENAME TO studios;

-- ============================================================================
-- DROP FLEETS
-- ============================================================================

DROP TABLE IF EXISTS qar.fleets;

-- ============================================================================
-- RECREATE ORIGINAL TRIGGERS
-- ============================================================================

CREATE TRIGGER c_studios_updated_at
    BEFORE UPDATE ON qar.studios
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER c_performers_updated_at
    BEFORE UPDATE ON qar.performers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER c_movies_updated_at
    BEFORE UPDATE ON qar.movies
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER c_scenes_updated_at
    BEFORE UPDATE ON qar.scenes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER c_user_ratings_updated_at
    BEFORE UPDATE ON qar.user_ratings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

COMMIT;

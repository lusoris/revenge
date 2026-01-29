-- TV Shows Core: Rollback
BEGIN;

DROP TRIGGER IF EXISTS seasons_update_counts ON seasons;
DROP TRIGGER IF EXISTS episodes_update_counts ON episodes;
DROP FUNCTION IF EXISTS update_series_counts();

DROP TABLE IF EXISTS series_network_link;
DROP TABLE IF EXISTS tv_networks;
DROP TABLE IF EXISTS episodes;
DROP TABLE IF EXISTS seasons;
DROP TABLE IF EXISTS series;

COMMIT;

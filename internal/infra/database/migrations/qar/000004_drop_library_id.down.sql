-- Rollback: Restore library_id columns
-- Note: Data cannot be restored as it was dropped

BEGIN;

-- Make fleet_id nullable again
ALTER TABLE qar.expeditions ALTER COLUMN fleet_id DROP NOT NULL;
ALTER TABLE qar.voyages ALTER COLUMN fleet_id DROP NOT NULL;

-- Restore library_id columns (empty - data is lost)
ALTER TABLE qar.expeditions ADD COLUMN library_id UUID;
ALTER TABLE qar.voyages ADD COLUMN library_id UUID;

COMMIT;

-- Drop deprecated library_id columns from expeditions and voyages
-- These tables now use fleet_id to reference qar.fleets

BEGIN;

-- Drop library_id from expeditions (replaced by fleet_id)
ALTER TABLE qar.expeditions DROP COLUMN IF EXISTS library_id;

-- Drop library_id from voyages (replaced by fleet_id)
ALTER TABLE qar.voyages DROP COLUMN IF EXISTS library_id;

-- Make fleet_id NOT NULL (it was added as nullable in 000003)
ALTER TABLE qar.expeditions ALTER COLUMN fleet_id SET NOT NULL;
ALTER TABLE qar.voyages ALTER COLUMN fleet_id SET NOT NULL;

COMMIT;

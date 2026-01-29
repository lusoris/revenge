-- Drop collections tables and types

DROP TRIGGER IF EXISTS trg_collection_items_update_metadata ON collection_items;
DROP FUNCTION IF EXISTS update_collection_metadata();

DROP TABLE IF EXISTS collection_subscriptions;
DROP TABLE IF EXISTS collection_tags;
DROP TABLE IF EXISTS collection_items;
DROP TABLE IF EXISTS collections;

DROP TYPE IF EXISTS collection_type;

-- 000006_media_items.down.sql
DROP TRIGGER IF EXISTS update_media_items_updated_at ON media_items;
DROP TABLE IF EXISTS media_items;
DROP TYPE IF EXISTS media_type;

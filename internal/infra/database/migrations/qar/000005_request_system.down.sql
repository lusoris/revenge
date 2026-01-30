-- QAR Request System: Rollback

BEGIN;

-- Drop triggers first
DROP TRIGGER IF EXISTS qar_provision_ayes_count_trigger ON qar.provision_ayes;
DROP FUNCTION IF EXISTS qar.update_provision_ayes_count();

DROP TRIGGER IF EXISTS qar_cargo_hold_updated_at ON qar.cargo_hold;
DROP TRIGGER IF EXISTS qar_articles_updated_at ON qar.articles;
DROP TRIGGER IF EXISTS qar_rations_updated_at ON qar.rations;
DROP TRIGGER IF EXISTS qar_provisions_updated_at ON qar.provisions;

-- Drop tables in reverse order (respecting FK constraints)
DROP TABLE IF EXISTS qar.cargo_hold;
DROP TABLE IF EXISTS qar.articles;
DROP TABLE IF EXISTS qar.rations;
DROP TABLE IF EXISTS qar.provision_missives;
DROP TABLE IF EXISTS qar.provision_ayes;
DROP TABLE IF EXISTS qar.provisions;

COMMIT;

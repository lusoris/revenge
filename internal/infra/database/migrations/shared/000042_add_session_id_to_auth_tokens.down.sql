DROP INDEX IF EXISTS shared.idx_auth_tokens_session_id;

ALTER TABLE shared.auth_tokens DROP COLUMN IF EXISTS session_id;

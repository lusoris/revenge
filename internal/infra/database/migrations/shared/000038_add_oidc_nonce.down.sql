-- Remove nonce column from oidc_states
ALTER TABLE shared.oidc_states DROP COLUMN IF EXISTS nonce;

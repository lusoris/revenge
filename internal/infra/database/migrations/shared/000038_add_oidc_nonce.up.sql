-- Add nonce column to oidc_states for OIDC nonce validation
ALTER TABLE shared.oidc_states ADD COLUMN nonce VARCHAR(128);

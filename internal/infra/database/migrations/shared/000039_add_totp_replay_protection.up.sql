-- Add last_used_code column for TOTP replay protection
ALTER TABLE public.user_totp_secrets
ADD COLUMN last_used_code VARCHAR(10);

COMMENT ON COLUMN public.user_totp_secrets.last_used_code IS 'Last successfully used TOTP code to prevent replay attacks';

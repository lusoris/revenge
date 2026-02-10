-- Remove last_used_code column
ALTER TABLE public.user_totp_secrets DROP COLUMN IF EXISTS last_used_code;

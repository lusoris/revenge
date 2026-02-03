-- Revert activity log table
DROP INDEX IF EXISTS public.idx_activity_log_success;
DROP INDEX IF EXISTS public.idx_activity_log_ip;
DROP INDEX IF EXISTS public.idx_activity_log_created;
DROP INDEX IF EXISTS public.idx_activity_log_resource;
DROP INDEX IF EXISTS public.idx_activity_log_action;
DROP INDEX IF EXISTS public.idx_activity_log_user;
DROP TABLE IF EXISTS public.activity_log;

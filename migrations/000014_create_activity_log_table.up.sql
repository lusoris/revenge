-- Activity audit log for tracking user actions
CREATE TABLE IF NOT EXISTS public.activity_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    -- Actor (who performed the action)
    user_id UUID REFERENCES shared.users(id) ON DELETE SET NULL,
    username VARCHAR(255),  -- Denormalized for historical records

    -- Action details
    action VARCHAR(100) NOT NULL,      -- 'user.login', 'user.logout', 'library.create', 'settings.update'
    resource_type VARCHAR(50),         -- 'user', 'library', 'movie', 'setting'
    resource_id UUID,

    -- Change data
    changes JSONB,                     -- {"field": {"old": "...", "new": "..."}}
    metadata JSONB,                    -- Additional context

    -- Request info
    ip_address INET,
    user_agent TEXT,

    -- Result
    success BOOLEAN DEFAULT true,
    error_message TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Performance indexes
CREATE INDEX idx_activity_log_user ON public.activity_log(user_id, created_at DESC);
CREATE INDEX idx_activity_log_action ON public.activity_log(action, created_at DESC);
CREATE INDEX idx_activity_log_resource ON public.activity_log(resource_type, resource_id);
CREATE INDEX idx_activity_log_created ON public.activity_log(created_at DESC);
CREATE INDEX idx_activity_log_ip ON public.activity_log(ip_address, created_at DESC);
CREATE INDEX idx_activity_log_success ON public.activity_log(success, created_at DESC) WHERE success = false;

COMMENT ON TABLE public.activity_log IS 'Audit log for tracking user actions and system events';
COMMENT ON COLUMN public.activity_log.action IS 'Action type: user.login, user.logout, library.create, settings.update, etc.';
COMMENT ON COLUMN public.activity_log.changes IS 'JSON object with field changes: {"field": {"old": "...", "new": "..."}}';
COMMENT ON COLUMN public.activity_log.metadata IS 'Additional context data as JSON';

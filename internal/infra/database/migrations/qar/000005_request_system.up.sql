-- QAR Request System: Community content requests for adult content
-- Isolated from public request system, uses QAR schema and terminology

BEGIN;

-- ============================================================================
-- PROVISIONS (Content Requests)
-- "Crew members can request provisions for the ship"
-- ============================================================================

CREATE TABLE qar.provisions (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id                 UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    -- Request type
    content_type            VARCHAR(50) NOT NULL CHECK (content_type IN ('expedition', 'voyage')),
    request_subtype         VARCHAR(50),  -- 'scene', 'port', 'crew', 'flag_combo'

    -- External reference
    external_id             VARCHAR(200),  -- StashDB ID, TPDB ID, etc.
    external_source         VARCHAR(50),   -- 'stashdb', 'tpdb', 'whisparr'
    title                   VARCHAR(500) NOT NULL,
    release_year            INT,

    -- Metadata (JSON with type-specific data)
    manifest                JSONB,         -- flags, crew_ids, port_id, etc.

    -- Status workflow
    status                  VARCHAR(50) NOT NULL DEFAULT 'pending'
                            CHECK (status IN ('pending', 'approved', 'processing', 'available', 'declined', 'on_hold')),
    auto_approved           BOOLEAN NOT NULL DEFAULT FALSE,
    auto_article_id         UUID,          -- Which rule triggered auto-approval

    -- Approval info
    approved_by_user_id     UUID REFERENCES users(id) ON DELETE SET NULL,
    approved_at             TIMESTAMPTZ,
    declined_reason         TEXT,

    -- Priority & voting
    priority                INT NOT NULL DEFAULT 0,
    ayes_count              INT NOT NULL DEFAULT 0,  -- Vote count (ayes = pirate votes)

    -- Integration
    integration_id          VARCHAR(200),  -- Whisparr ID after approval
    integration_status      VARCHAR(100),

    -- Storage estimation
    estimated_cargo_gb      DECIMAL(10,2), -- Estimated size
    actual_cargo_gb         DECIMAL(10,2), -- Actual size after download

    -- Timestamps
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    available_at            TIMESTAMPTZ,

    -- Automation tracking
    triggered_by_automation BOOLEAN NOT NULL DEFAULT FALSE,
    parent_provision_id     UUID REFERENCES qar.provisions(id) ON DELETE SET NULL
);

-- Indexes
CREATE INDEX idx_qar_provisions_user ON qar.provisions(user_id);
CREATE INDEX idx_qar_provisions_status ON qar.provisions(status);
CREATE INDEX idx_qar_provisions_content_type ON qar.provisions(content_type);
CREATE INDEX idx_qar_provisions_created ON qar.provisions(created_at DESC);
CREATE INDEX idx_qar_provisions_priority ON qar.provisions(priority DESC, ayes_count DESC);
CREATE INDEX idx_qar_provisions_external ON qar.provisions(external_source, external_id) WHERE external_id IS NOT NULL;

-- Updated at trigger
CREATE TRIGGER qar_provisions_updated_at
    BEFORE UPDATE ON qar.provisions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- ============================================================================
-- PROVISION AYES (Request Votes)
-- "The crew votes aye on provisions they want"
-- ============================================================================

CREATE TABLE qar.provision_ayes (
    provision_id            UUID NOT NULL REFERENCES qar.provisions(id) ON DELETE CASCADE,
    user_id                 UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    voted_at                TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (provision_id, user_id)
);

CREATE INDEX idx_qar_provision_ayes_user ON qar.provision_ayes(user_id);

-- ============================================================================
-- PROVISION MISSIVES (Request Comments)
-- "Messages about the provision request"
-- ============================================================================

CREATE TABLE qar.provision_missives (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provision_id            UUID NOT NULL REFERENCES qar.provisions(id) ON DELETE CASCADE,
    user_id                 UUID REFERENCES users(id) ON DELETE SET NULL,
    message                 TEXT NOT NULL,
    is_captain_order        BOOLEAN NOT NULL DEFAULT FALSE,  -- Admin/mod comment
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_qar_provision_missives_provision ON qar.provision_missives(provision_id);

-- ============================================================================
-- RATIONS (Request Quotas per User)
-- "Each crew member has limited rations they can request"
-- ============================================================================

CREATE TABLE qar.rations (
    user_id                 UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,

    -- Request limits
    daily_limit             INT NOT NULL DEFAULT 5,
    weekly_limit            INT NOT NULL DEFAULT 20,
    monthly_limit           INT NOT NULL DEFAULT 50,
    daily_used              INT NOT NULL DEFAULT 0,
    weekly_used             INT NOT NULL DEFAULT 0,
    monthly_used            INT NOT NULL DEFAULT 0,

    -- Storage quota (GB)
    cargo_quota_gb          DECIMAL(10,2) NOT NULL DEFAULT 500,
    cargo_used_gb           DECIMAL(10,2) NOT NULL DEFAULT 0,

    -- Reset timestamps
    last_reset_daily        DATE NOT NULL DEFAULT CURRENT_DATE,
    last_reset_weekly       DATE NOT NULL DEFAULT CURRENT_DATE,
    last_reset_monthly      DATE NOT NULL DEFAULT CURRENT_DATE,

    -- Timestamps
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER qar_rations_updated_at
    BEFORE UPDATE ON qar.rations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- ============================================================================
-- ARTICLES (Request Rules / Auto-approval rules)
-- "The Pirate's Code - articles governing provisions"
-- ============================================================================

CREATE TABLE qar.articles (
    id                      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name                    VARCHAR(200) NOT NULL,
    description             TEXT,

    -- Targeting
    content_type            VARCHAR(50),  -- NULL = all types, 'expedition' or 'voyage'

    -- Condition
    condition_type          VARCHAR(50) NOT NULL,
    -- Types: 'user_role', 'trust_score', 'crew_preference', 'port_preference',
    --        'flag_preference', 'storage_available', 'release_year'
    condition_value         JSONB NOT NULL,

    -- Action
    action                  VARCHAR(50) NOT NULL DEFAULT 'auto_approve'
                            CHECK (action IN ('auto_approve', 'require_approval', 'decline', 'on_hold')),

    -- Automation trigger
    automation_trigger      VARCHAR(50),
    -- Triggers: NULL (manual), 'watch_completed', 'storage_low'

    -- State
    enabled                 BOOLEAN NOT NULL DEFAULT TRUE,
    priority                INT NOT NULL DEFAULT 0,  -- Higher = checked first

    -- Timestamps
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_qar_articles_enabled ON qar.articles(enabled, priority DESC) WHERE enabled = TRUE;
CREATE INDEX idx_qar_articles_automation ON qar.articles(automation_trigger) WHERE automation_trigger IS NOT NULL;

CREATE TRIGGER qar_articles_updated_at
    BEFORE UPDATE ON qar.articles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- ============================================================================
-- GLOBAL CARGO QUOTA (Server-wide storage limits for QAR content)
-- ============================================================================

CREATE TABLE qar.cargo_hold (
    id                      INT PRIMARY KEY DEFAULT 1 CHECK (id = 1),  -- Single row
    total_quota_gb          DECIMAL(10,2) NOT NULL DEFAULT 2000,
    total_used_gb           DECIMAL(10,2) NOT NULL DEFAULT 0,
    expedition_quota_gb     DECIMAL(10,2) NOT NULL DEFAULT 1000,
    expedition_used_gb      DECIMAL(10,2) NOT NULL DEFAULT 0,
    voyage_quota_gb         DECIMAL(10,2) NOT NULL DEFAULT 1000,
    voyage_used_gb          DECIMAL(10,2) NOT NULL DEFAULT 0,
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Insert default global quota
INSERT INTO qar.cargo_hold (id) VALUES (1);

CREATE TRIGGER qar_cargo_hold_updated_at
    BEFORE UPDATE ON qar.cargo_hold
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- ============================================================================
-- FUNCTIONS: Aye counting
-- ============================================================================

-- Function to update ayes_count when votes change
CREATE OR REPLACE FUNCTION qar.update_provision_ayes_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE qar.provisions
        SET ayes_count = ayes_count + 1
        WHERE id = NEW.provision_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE qar.provisions
        SET ayes_count = ayes_count - 1
        WHERE id = OLD.provision_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER qar_provision_ayes_count_trigger
    AFTER INSERT OR DELETE ON qar.provision_ayes
    FOR EACH ROW EXECUTE FUNCTION qar.update_provision_ayes_count();

COMMIT;

-- Profiles table: Netflix-style user profiles
-- Each user can have multiple profiles with separate watch history, preferences
CREATE TABLE profiles (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name            VARCHAR(100) NOT NULL,
    avatar_url      TEXT,                            -- Profile picture
    is_default      BOOLEAN NOT NULL DEFAULT false,  -- Primary profile for user
    is_kids         BOOLEAN NOT NULL DEFAULT false,  -- Kids profile (restricted content)

    -- Content access (can only be MORE restrictive than parent user)
    max_rating_level    INT NOT NULL DEFAULT 100,
    adult_enabled       BOOLEAN NOT NULL DEFAULT false,

    -- Preferences (override user defaults)
    preferred_language          VARCHAR(10),
    preferred_audio_language    VARCHAR(10),
    preferred_subtitle_language VARCHAR(10),
    autoplay_next               BOOLEAN DEFAULT true,
    autoplay_previews           BOOLEAN DEFAULT true,

    -- Timestamps
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Constraints
    CONSTRAINT profile_rating_not_exceeds_user CHECK (max_rating_level <= 100)
);

-- Indexes
CREATE INDEX idx_profiles_user_id ON profiles(user_id);
CREATE UNIQUE INDEX idx_profiles_user_default ON profiles(user_id) WHERE is_default = true;

-- Trigger for updated_at
CREATE TRIGGER profiles_updated_at
    BEFORE UPDATE ON profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- Create default profile for each new user
CREATE OR REPLACE FUNCTION create_default_profile()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO profiles (user_id, name, is_default, max_rating_level, adult_enabled)
    VALUES (NEW.id, NEW.username, true, NEW.max_rating_level, NEW.adult_enabled);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_create_default_profile
    AFTER INSERT ON users
    FOR EACH ROW EXECUTE FUNCTION create_default_profile();

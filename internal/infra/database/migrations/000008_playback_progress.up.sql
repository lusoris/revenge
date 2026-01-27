-- 000008_playback_progress.up.sql
-- User playback progress tracking

CREATE TABLE playback_progress (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_id UUID NOT NULL REFERENCES media_items(id) ON DELETE CASCADE,
    position_ticks BIGINT NOT NULL DEFAULT 0, -- Playback position in ticks (100ns)
    played BOOLEAN NOT NULL DEFAULT false,    -- Has been fully watched
    play_count INT NOT NULL DEFAULT 0,
    last_played_at TIMESTAMPTZ,
    audio_stream_index INT,                   -- Selected audio track
    subtitle_stream_index INT,                -- Selected subtitle track
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, item_id)
);

CREATE INDEX idx_playback_progress_user_id ON playback_progress(user_id);
CREATE INDEX idx_playback_progress_item_id ON playback_progress(item_id);
CREATE INDEX idx_playback_progress_last_played ON playback_progress(user_id, last_played_at DESC);
CREATE INDEX idx_playback_progress_played ON playback_progress(user_id, played) WHERE played = true;

CREATE TRIGGER update_playback_progress_updated_at
    BEFORE UPDATE ON playback_progress
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

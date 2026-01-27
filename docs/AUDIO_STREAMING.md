# Revenge - Audio Streaming & Progress Tracking

> Complete audio streaming with progress persistence and session management.

## Overview

Audio streaming in Revenge covers multiple content types with unique requirements:

| Content Type | Progress Unit | Sync Frequency | Special Features |
|--------------|---------------|----------------|------------------|
| Music | Track position | On pause/next | Scrobbling, queue |
| Podcast | Episode position | Every 30s | Resume, skip segments |
| Audiobook | Chapter position | Every 30s | Chapter navigation, bookmarks |
| Live Radio | N/A | N/A | Now playing metadata |

---

## Audio Streaming Architecture

### Flow Diagram

```
┌──────────┐     ┌──────────┐     ┌──────────────┐     ┌──────────────┐
│  Client  │ ──→ │ Revenge  │ ──→ │   Storage    │     │  Blackbeard  │
│  (Web)   │     │  Server  │     │   (Files)    │     │ (Transcode)  │
└──────────┘     └──────────┘     └──────────────┘     └──────────────┘
     │                │                   │                    │
     │   1. Play      │                   │                    │
     │───────────────→│                   │                    │
     │                │  2. Check format  │                    │
     │                │───────────────────│                    │
     │                │                   │                    │
     │                │  3. Direct/Trans? │                    │
     │                │                   │                    │
     │                ├── Direct ─────────┼── 4a. Stream ─────→│
     │                │                   │                    │
     │                └── Transcode ──────┼── 4b. Request ────→│
     │                                    │                    │
     │   5. HLS/Audio stream              │    Transcoded      │
     │←───────────────────────────────────┼────────────────────┘
     │                                    │
     │   6. Progress updates (periodic)   │
     │───────────────────────────────────→│
```

### Direct Play vs Transcoding

```go
type AudioStreamDecision struct {
    CanDirectPlay bool
    Reason        string
    TranscodeArgs *TranscodeRequest
}

func (s *AudioService) DecideStreamMethod(ctx context.Context, track *AudioTrack, client *ClientCapabilities) AudioStreamDecision {
    // Check codec support
    codecSupported := slices.Contains(client.SupportedAudioCodecs, track.Codec)

    // Check container support
    containerSupported := slices.Contains(client.SupportedContainers, track.Container)

    // Check bitrate limits
    bitrateOK := client.MaxAudioBitrate == 0 || track.Bitrate <= client.MaxAudioBitrate

    // Check channel support
    channelsOK := client.MaxAudioChannels == 0 || track.Channels <= client.MaxAudioChannels

    if codecSupported && containerSupported && bitrateOK && channelsOK {
        return AudioStreamDecision{
            CanDirectPlay: true,
            Reason:        "direct_play",
        }
    }

    // Need transcoding
    reason := []string{}
    if !codecSupported {
        reason = append(reason, fmt.Sprintf("codec %s not supported", track.Codec))
    }
    if !bitrateOK {
        reason = append(reason, fmt.Sprintf("bitrate %d exceeds max %d", track.Bitrate, client.MaxAudioBitrate))
    }

    return AudioStreamDecision{
        CanDirectPlay: false,
        Reason:        strings.Join(reason, ", "),
        TranscodeArgs: &TranscodeRequest{
            MediaID:        track.ID,
            TargetCodec:    preferredAudioCodec(client),
            TargetBitrate:  min(track.Bitrate, client.MaxAudioBitrate),
            TargetChannels: min(track.Channels, client.MaxAudioChannels),
        },
    }
}

func preferredAudioCodec(client *ClientCapabilities) string {
    // Preference order: opus > aac > mp3
    if slices.Contains(client.SupportedAudioCodecs, "opus") {
        return "opus"
    }
    if slices.Contains(client.SupportedAudioCodecs, "aac") {
        return "aac"
    }
    return "mp3"
}
```

### Audio Formats Support

| Format | Container | Codec | Direct Play | Transcode Target |
|--------|-----------|-------|-------------|------------------|
| FLAC | flac | flac | Supported clients | AAC/Opus |
| MP3 | mp3 | mp3 | Universal | - |
| AAC | m4a, mp4 | aac | Most clients | MP3 |
| Opus | ogg, webm | opus | Modern browsers | AAC |
| ALAC | m4a | alac | Apple devices | AAC |
| WAV | wav | pcm | Limited | AAC |
| OGG Vorbis | ogg | vorbis | Firefox, Chrome | AAC |

---

## Bandwidth-Aware Audio Streaming

### Overview

External clients (outside local network) need adaptive quality based on measured bandwidth:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Bandwidth Adaptation Flow                            │
└─────────────────────────────────────────────────────────────────────────┘

┌──────────┐  1. Request  ┌──────────┐  2. Check  ┌──────────────────────┐
│  Client  │ ───────────→ │ Revenge  │ ─────────→ │ Is External Client?  │
└──────────┘              └──────────┘            └──────────────────────┘
                                                           │
                              ┌─────────────────────────────┤
                              │                             │
                              ▼ YES                         ▼ NO
                    ┌─────────────────────┐      ┌─────────────────────┐
                    │ Measure Bandwidth   │      │ Use Full Quality    │
                    │ (rolling average)   │      │ (direct play/FLAC)  │
                    └─────────────────────┘      └─────────────────────┘
                              │
                              ▼
                    ┌─────────────────────┐
                    │ Select Quality Tier │
                    │ Based on Bandwidth  │
                    └─────────────────────┘
```

### Audio Quality Tiers

| Bandwidth | Quality Tier | Codec | Bitrate | Use Case |
|-----------|--------------|-------|---------|----------|
| > 5 Mbps | Lossless | FLAC | ~1400 kbps | Local/WiFi |
| 1-5 Mbps | High | AAC/Opus | 320 kbps | Stable connection |
| 500kbps-1Mbps | Medium | AAC/Opus | 192 kbps | Mobile (good) |
| 250-500 kbps | Low | AAC/Opus | 128 kbps | Mobile (fair) |
| < 250 kbps | Minimal | AAC | 64 kbps | Poor connection |

### Bandwidth Measurement

```go
type AudioBandwidthAdapter struct {
    samples       []BandwidthSample
    windowSize    int           // Number of samples to keep
    minSamples    int           // Minimum samples for reliable estimate
    jitterWeight  float64       // How much to penalize unstable connections
}

type BandwidthSample struct {
    Timestamp time.Time
    Kbps      int
    Latency   time.Duration
}

func (a *AudioBandwidthAdapter) RecordSample(bytesReceived int64, duration time.Duration) {
    if duration == 0 {
        return
    }
    
    kbps := int(float64(bytesReceived*8) / duration.Seconds() / 1000)
    a.samples = append(a.samples, BandwidthSample{
        Timestamp: time.Now(),
        Kbps:      kbps,
    })
    
    // Keep rolling window
    if len(a.samples) > a.windowSize {
        a.samples = a.samples[len(a.samples)-a.windowSize:]
    }
}

func (a *AudioBandwidthAdapter) RecommendedAudioBitrate() int {
    if len(a.samples) < a.minSamples {
        return 192 // Default to medium quality until we have data
    }
    
    // Calculate average
    var sum int
    for _, s := range a.samples {
        sum += s.Kbps
    }
    avg := sum / len(a.samples)
    
    // Calculate jitter (standard deviation)
    var variance float64
    for _, s := range a.samples {
        diff := float64(s.Kbps - avg)
        variance += diff * diff
    }
    jitter := int(math.Sqrt(variance / float64(len(a.samples))))
    
    // Conservative estimate: 70% of average minus jitter penalty
    safeBandwidth := int(float64(avg)*0.7) - int(float64(jitter)*a.jitterWeight)
    
    return a.mapToAudioTier(safeBandwidth)
}

func (a *AudioBandwidthAdapter) mapToAudioTier(safeBandwidth int) int {
    switch {
    case safeBandwidth >= 1500:
        return 0       // Lossless (no transcode)
    case safeBandwidth >= 400:
        return 320     // High quality
    case safeBandwidth >= 250:
        return 192     // Medium quality
    case safeBandwidth >= 150:
        return 128     // Low quality
    default:
        return 64      // Minimal (voice-quality)
    }
}
```

### Client Bandwidth Reporting

```go
// Client reports bandwidth in session update
type SessionBandwidthUpdate struct {
    MeasuredKbps    int    `json:"measured_kbps"`
    NetworkType     string `json:"network_type"`  // "wifi", "cellular", "ethernet"
    SaveDataEnabled bool   `json:"save_data"`     // Browser Save-Data header
}

// Session handler updates bandwidth
func (h *SessionHandler) HandleBandwidthUpdate(w http.ResponseWriter, r *http.Request) {
    var update SessionBandwidthUpdate
    if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
        http.Error(w, "invalid request", http.StatusBadRequest)
        return
    }
    
    session := getSessionFromContext(r.Context())
    
    // Update bandwidth adapter
    session.BandwidthAdapter.RecordSample(
        int64(update.MeasuredKbps*1000/8), // Convert to bytes
        time.Second,
    )
    
    // Check if quality should change
    newBitrate := session.BandwidthAdapter.RecommendedAudioBitrate()
    if newBitrate != session.CurrentAudioBitrate {
        session.CurrentAudioBitrate = newBitrate
        
        // Notify client to switch quality
        h.notifyQualityChange(session, newBitrate)
    }
    
    w.WriteHeader(http.StatusOK)
}
```

### Frontend Bandwidth Monitoring

```typescript
// lib/audio/bandwidth.ts
class AudioBandwidthMonitor {
    private samples: number[] = [];
    private readonly windowSize = 10;
    private measurementInterval: number | null = null;
    
    start() {
        // Measure every 30 seconds during playback
        this.measurementInterval = setInterval(() => {
            this.measureAndReport();
        }, 30000);
    }
    
    stop() {
        if (this.measurementInterval) {
            clearInterval(this.measurementInterval);
        }
    }
    
    async measureAndReport() {
        const bandwidth = await this.measureBandwidth();
        this.samples.push(bandwidth);
        
        if (this.samples.length > this.windowSize) {
            this.samples.shift();
        }
        
        // Report to server
        await fetch('/api/v1/session/bandwidth', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                measured_kbps: Math.round(bandwidth),
                network_type: this.getNetworkType(),
                save_data: navigator.connection?.saveData ?? false,
            }),
        });
    }
    
    private async measureBandwidth(): Promise<number> {
        const testUrl = '/api/v1/bandwidth-test?size=50000'; // 50KB test
        const start = performance.now();
        
        try {
            const response = await fetch(testUrl, { cache: 'no-store' });
            const data = await response.arrayBuffer();
            const duration = performance.now() - start;
            
            // Calculate kbps
            return (data.byteLength * 8) / duration; // kbps
        } catch {
            return this.getAverageKbps(); // Use cached average on error
        }
    }
    
    private getNetworkType(): string {
        const conn = (navigator as any).connection;
        return conn?.effectiveType || 'unknown';
    }
    
    getAverageKbps(): number {
        if (this.samples.length === 0) return 0;
        return this.samples.reduce((a, b) => a + b, 0) / this.samples.length;
    }
}

export const bandwidthMonitor = new AudioBandwidthMonitor();
```

### Quality Switch Handling

```typescript
// Listen for server-initiated quality changes
websocket.on('quality_change', (data: { bitrate: number }) => {
    const player = getAudioPlayer();
    
    // Store current position
    const position = player.currentTime;
    
    // Get new stream URL with updated bitrate
    const newUrl = `/api/v1/audio/${currentTrackId}/stream?bitrate=${data.bitrate}`;
    
    // Seamless switch
    player.src = newUrl;
    player.currentTime = position;
    player.play();
    
    // Show notification
    showToast(`Quality adjusted to ${data.bitrate}kbps due to network conditions`);
});
```

---

## Progress Tracking

### Database Schema

```sql
-- Universal progress tracking (per module)
CREATE TABLE music_playback_progress (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    track_id        UUID NOT NULL REFERENCES tracks(id) ON DELETE CASCADE,

    -- Progress
    position_ms     BIGINT NOT NULL DEFAULT 0,
    duration_ms     BIGINT NOT NULL,
    percentage      FLOAT GENERATED ALWAYS AS (position_ms::float / NULLIF(duration_ms, 0) * 100) STORED,

    -- Play count
    play_count      INT NOT NULL DEFAULT 0,
    completed_count INT NOT NULL DEFAULT 0,  -- Full listens

    -- Timestamps
    last_played_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (user_id, track_id)
);

CREATE TABLE podcast_playback_progress (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    episode_id      UUID NOT NULL REFERENCES podcast_episodes(id) ON DELETE CASCADE,

    -- Progress
    position_ms     BIGINT NOT NULL DEFAULT 0,
    duration_ms     BIGINT NOT NULL,
    percentage      FLOAT GENERATED ALWAYS AS (position_ms::float / NULLIF(duration_ms, 0) * 100) STORED,

    -- State
    is_completed    BOOLEAN NOT NULL DEFAULT false,

    -- Timestamps
    last_played_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (user_id, episode_id)
);

CREATE TABLE audiobook_playback_progress (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    audiobook_id    UUID NOT NULL REFERENCES audiobooks(id) ON DELETE CASCADE,

    -- Overall progress
    current_chapter INT NOT NULL DEFAULT 0,
    position_ms     BIGINT NOT NULL DEFAULT 0,  -- Within current chapter
    total_position_ms BIGINT NOT NULL DEFAULT 0, -- Overall position

    -- Completion
    total_duration_ms BIGINT NOT NULL,
    percentage      FLOAT GENERATED ALWAYS AS (total_position_ms::float / NULLIF(total_duration_ms, 0) * 100) STORED,
    is_completed    BOOLEAN NOT NULL DEFAULT false,

    -- Bookmarks
    bookmarks       JSONB DEFAULT '[]',
    -- [{"chapter": 5, "position_ms": 123456, "note": "Important quote", "created_at": "..."}]

    -- Timestamps
    last_played_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (user_id, audiobook_id)
);

-- Indexes for efficient queries
CREATE INDEX idx_music_progress_user ON music_playback_progress(user_id, last_played_at DESC);
CREATE INDEX idx_podcast_progress_user ON podcast_playback_progress(user_id, last_played_at DESC);
CREATE INDEX idx_audiobook_progress_user ON audiobook_playback_progress(user_id, last_played_at DESC);
```

### Progress Service

```go
type ProgressService struct {
    db     *pgxpool.Pool
    cache  *cache.Client
    logger *slog.Logger
}

type PlaybackProgress struct {
    ItemID        uuid.UUID     `json:"item_id"`
    ItemType      string        `json:"item_type"`  // "track", "episode", "audiobook"
    PositionMs    int64         `json:"position_ms"`
    DurationMs    int64         `json:"duration_ms"`
    Percentage    float64       `json:"percentage"`
    IsCompleted   bool          `json:"is_completed"`
    LastPlayedAt  time.Time     `json:"last_played_at"`
}

// Update progress (called periodically from client)
func (s *ProgressService) UpdateProgress(ctx context.Context, userID uuid.UUID, update ProgressUpdate) error {
    // Validate
    if update.PositionMs < 0 || update.PositionMs > update.DurationMs {
        return errors.New("invalid position")
    }

    // Check completion threshold (e.g., 90% for podcasts/audiobooks)
    isCompleted := false
    if update.ItemType == "episode" || update.ItemType == "audiobook" {
        isCompleted = float64(update.PositionMs)/float64(update.DurationMs) >= 0.9
    }

    // Music: count as "complete" if played >50%
    if update.ItemType == "track" {
        isCompleted = float64(update.PositionMs)/float64(update.DurationMs) >= 0.5
    }

    // Upsert to database
    switch update.ItemType {
    case "track":
        return s.updateMusicProgress(ctx, userID, update, isCompleted)
    case "episode":
        return s.updatePodcastProgress(ctx, userID, update, isCompleted)
    case "audiobook":
        return s.updateAudiobookProgress(ctx, userID, update, isCompleted)
    default:
        return fmt.Errorf("unknown item type: %s", update.ItemType)
    }
}

func (s *ProgressService) updateMusicProgress(ctx context.Context, userID uuid.UUID, update ProgressUpdate, isCompleted bool) error {
    _, err := s.db.Exec(ctx, `
        INSERT INTO music_playback_progress (user_id, track_id, position_ms, duration_ms, play_count, completed_count, last_played_at)
        VALUES ($1, $2, $3, $4, 1, CASE WHEN $5 THEN 1 ELSE 0 END, NOW())
        ON CONFLICT (user_id, track_id) DO UPDATE SET
            position_ms = EXCLUDED.position_ms,
            duration_ms = EXCLUDED.duration_ms,
            play_count = music_playback_progress.play_count + 1,
            completed_count = music_playback_progress.completed_count + CASE WHEN $5 AND music_playback_progress.position_ms < $3 THEN 1 ELSE 0 END,
            last_played_at = NOW()
    `, userID, update.ItemID, update.PositionMs, update.DurationMs, isCompleted)
    return err
}

// Get resume position
func (s *ProgressService) GetResumePosition(ctx context.Context, userID uuid.UUID, itemID uuid.UUID, itemType string) (*PlaybackProgress, error) {
    // Try cache first
    cacheKey := fmt.Sprintf("progress:%s:%s:%s", userID, itemType, itemID)
    if cached, err := s.cache.Get(ctx, cacheKey); err == nil {
        var progress PlaybackProgress
        if json.Unmarshal(cached, &progress) == nil {
            return &progress, nil
        }
    }

    // Query database based on type
    var progress PlaybackProgress
    var err error

    switch itemType {
    case "track":
        err = s.db.QueryRow(ctx, `
            SELECT track_id, position_ms, duration_ms, percentage,
                   play_count >= 1 AND completed_count >= 1, last_played_at
            FROM music_playback_progress
            WHERE user_id = $1 AND track_id = $2
        `, userID, itemID).Scan(
            &progress.ItemID, &progress.PositionMs, &progress.DurationMs,
            &progress.Percentage, &progress.IsCompleted, &progress.LastPlayedAt,
        )
        progress.ItemType = "track"

    case "episode":
        err = s.db.QueryRow(ctx, `
            SELECT episode_id, position_ms, duration_ms, percentage,
                   is_completed, last_played_at
            FROM podcast_playback_progress
            WHERE user_id = $1 AND episode_id = $2
        `, userID, itemID).Scan(
            &progress.ItemID, &progress.PositionMs, &progress.DurationMs,
            &progress.Percentage, &progress.IsCompleted, &progress.LastPlayedAt,
        )
        progress.ItemType = "episode"

    case "audiobook":
        err = s.db.QueryRow(ctx, `
            SELECT audiobook_id, total_position_ms, total_duration_ms, percentage,
                   is_completed, last_played_at
            FROM audiobook_playback_progress
            WHERE user_id = $1 AND audiobook_id = $2
        `, userID, itemID).Scan(
            &progress.ItemID, &progress.PositionMs, &progress.DurationMs,
            &progress.Percentage, &progress.IsCompleted, &progress.LastPlayedAt,
        )
        progress.ItemType = "audiobook"
    }

    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, nil // No progress yet
        }
        return nil, err
    }

    // Cache for quick access
    data, _ := json.Marshal(progress)
    s.cache.Set(ctx, cacheKey, data, 5*time.Minute)

    return &progress, nil
}
```

### Progress Update Batching

```go
// Batch progress updates to reduce database writes
type ProgressBatcher struct {
    mu       sync.Mutex
    pending  map[string]ProgressUpdate
    service  *ProgressService
    interval time.Duration
}

func NewProgressBatcher(service *ProgressService, interval time.Duration) *ProgressBatcher {
    b := &ProgressBatcher{
        pending:  make(map[string]ProgressUpdate),
        service:  service,
        interval: interval,
    }
    go b.run()
    return b
}

func (b *ProgressBatcher) Queue(userID uuid.UUID, update ProgressUpdate) {
    b.mu.Lock()
    defer b.mu.Unlock()

    key := fmt.Sprintf("%s:%s:%s", userID, update.ItemType, update.ItemID)
    b.pending[key] = update
}

func (b *ProgressBatcher) run() {
    ticker := time.NewTicker(b.interval)
    defer ticker.Stop()

    for range ticker.C {
        b.flush()
    }
}

func (b *ProgressBatcher) flush() {
    b.mu.Lock()
    updates := b.pending
    b.pending = make(map[string]ProgressUpdate)
    b.mu.Unlock()

    if len(updates) == 0 {
        return
    }

    ctx := context.Background()
    for key, update := range updates {
        parts := strings.SplitN(key, ":", 3)
        userID, _ := uuid.Parse(parts[0])

        if err := b.service.UpdateProgress(ctx, userID, update); err != nil {
            slog.Error("failed to update progress", "key", key, "error", err)
        }
    }
}
```

---

## Playback Session

### Session State

```go
type AudioPlaybackSession struct {
    ID           uuid.UUID          `json:"id"`
    UserID       uuid.UUID          `json:"user_id"`
    DeviceID     string             `json:"device_id"`

    // Current item
    ItemID       uuid.UUID          `json:"item_id"`
    ItemType     string             `json:"item_type"`  // "track", "episode", "audiobook_chapter"

    // State
    State        PlaybackState      `json:"state"`      // "playing", "paused", "stopped"
    PositionMs   int64              `json:"position_ms"`
    DurationMs   int64              `json:"duration_ms"`

    // Queue (for music)
    Queue        []uuid.UUID        `json:"queue,omitempty"`
    QueueIndex   int                `json:"queue_index,omitempty"`
    ShuffleMode  bool               `json:"shuffle_mode,omitempty"`
    RepeatMode   RepeatMode         `json:"repeat_mode,omitempty"` // "off", "one", "all"

    // Metadata (for real-time updates)
    NowPlaying   *NowPlayingInfo    `json:"now_playing,omitempty"`

    // Timing
    StartedAt    time.Time          `json:"started_at"`
    LastUpdateAt time.Time          `json:"last_update_at"`
}

type NowPlayingInfo struct {
    Title       string   `json:"title"`
    Artist      string   `json:"artist,omitempty"`
    Album       string   `json:"album,omitempty"`
    ArtworkURL  string   `json:"artwork_url,omitempty"`
    Duration    int64    `json:"duration"`
}

type RepeatMode string

const (
    RepeatOff RepeatMode = "off"
    RepeatOne RepeatMode = "one"
    RepeatAll RepeatMode = "all"
)
```

### Session Manager

```go
type AudioSessionManager struct {
    sessions sync.Map // map[userID+deviceID]*AudioPlaybackSession
    cache    *cache.Client
    progress *ProgressService
}

func (m *AudioSessionManager) StartSession(ctx context.Context, req StartSessionRequest) (*AudioPlaybackSession, error) {
    session := &AudioPlaybackSession{
        ID:         uuid.New(),
        UserID:     req.UserID,
        DeviceID:   req.DeviceID,
        ItemID:     req.ItemID,
        ItemType:   req.ItemType,
        State:      PlaybackStatePlaying,
        PositionMs: req.StartPosition,
        StartedAt:  time.Now(),
    }

    // Get resume position if not specified
    if req.StartPosition == 0 && req.Resume {
        if progress, _ := m.progress.GetResumePosition(ctx, req.UserID, req.ItemID, req.ItemType); progress != nil {
            // Don't resume if nearly complete
            if progress.Percentage < 95 {
                session.PositionMs = progress.PositionMs
            }
        }
    }

    // Store session
    key := fmt.Sprintf("%s:%s", req.UserID, req.DeviceID)
    m.sessions.Store(key, session)

    // Cache for cross-device sync
    m.cacheSession(ctx, session)

    return session, nil
}

func (m *AudioSessionManager) UpdatePosition(ctx context.Context, userID uuid.UUID, deviceID string, positionMs int64) error {
    key := fmt.Sprintf("%s:%s", userID, deviceID)
    val, ok := m.sessions.Load(key)
    if !ok {
        return ErrSessionNotFound
    }

    session := val.(*AudioPlaybackSession)
    session.PositionMs = positionMs
    session.LastUpdateAt = time.Now()

    // Update cache
    m.cacheSession(ctx, session)

    return nil
}

func (m *AudioSessionManager) cacheSession(ctx context.Context, session *AudioPlaybackSession) {
    // Store in Dragonfly for cross-device access
    key := fmt.Sprintf("audio_session:%s", session.UserID)
    data, _ := json.Marshal(session)
    m.cache.Set(ctx, key, data, 24*time.Hour)
}
```

---

## Music Queue

### Queue Management

```go
type MusicQueueService struct {
    cache *cache.Client
}

type MusicQueue struct {
    UserID       uuid.UUID   `json:"user_id"`
    DeviceID     string      `json:"device_id"`
    Tracks       []uuid.UUID `json:"tracks"`
    CurrentIndex int         `json:"current_index"`
    ShuffleMode  bool        `json:"shuffle_mode"`
    RepeatMode   RepeatMode  `json:"repeat_mode"`
    OriginalOrder []uuid.UUID `json:"original_order,omitempty"` // For unshuffle
}

func (s *MusicQueueService) CreateQueue(ctx context.Context, userID uuid.UUID, deviceID string, tracks []uuid.UUID, startIndex int) (*MusicQueue, error) {
    queue := &MusicQueue{
        UserID:       userID,
        DeviceID:     deviceID,
        Tracks:       tracks,
        CurrentIndex: startIndex,
        ShuffleMode:  false,
        RepeatMode:   RepeatOff,
    }

    return queue, s.saveQueue(ctx, queue)
}

func (s *MusicQueueService) AddToQueue(ctx context.Context, userID uuid.UUID, deviceID string, tracks []uuid.UUID, position string) error {
    queue, err := s.GetQueue(ctx, userID, deviceID)
    if err != nil {
        return err
    }

    switch position {
    case "next":
        // Insert after current
        insertAt := queue.CurrentIndex + 1
        queue.Tracks = slices.Insert(queue.Tracks, insertAt, tracks...)
    case "last":
        // Append to end
        queue.Tracks = append(queue.Tracks, tracks...)
    default:
        return fmt.Errorf("invalid position: %s", position)
    }

    return s.saveQueue(ctx, queue)
}

func (s *MusicQueueService) Shuffle(ctx context.Context, userID uuid.UUID, deviceID string, enable bool) error {
    queue, err := s.GetQueue(ctx, userID, deviceID)
    if err != nil {
        return err
    }

    if enable && !queue.ShuffleMode {
        // Save original order
        queue.OriginalOrder = slices.Clone(queue.Tracks)

        // Shuffle everything except current track
        currentTrack := queue.Tracks[queue.CurrentIndex]
        remaining := slices.Delete(slices.Clone(queue.Tracks), queue.CurrentIndex, queue.CurrentIndex+1)

        rand.Shuffle(len(remaining), func(i, j int) {
            remaining[i], remaining[j] = remaining[j], remaining[i]
        })

        // Put current track first
        queue.Tracks = append([]uuid.UUID{currentTrack}, remaining...)
        queue.CurrentIndex = 0
        queue.ShuffleMode = true

    } else if !enable && queue.ShuffleMode {
        // Restore original order
        currentTrack := queue.Tracks[queue.CurrentIndex]
        queue.Tracks = queue.OriginalOrder
        queue.CurrentIndex = slices.Index(queue.Tracks, currentTrack)
        queue.OriginalOrder = nil
        queue.ShuffleMode = false
    }

    return s.saveQueue(ctx, queue)
}

func (s *MusicQueueService) Next(ctx context.Context, userID uuid.UUID, deviceID string) (*uuid.UUID, error) {
    queue, err := s.GetQueue(ctx, userID, deviceID)
    if err != nil {
        return nil, err
    }

    switch queue.RepeatMode {
    case RepeatOne:
        // Stay on current track
        return &queue.Tracks[queue.CurrentIndex], nil

    case RepeatAll:
        // Move to next, wrap around
        queue.CurrentIndex = (queue.CurrentIndex + 1) % len(queue.Tracks)

    case RepeatOff:
        // Move to next, stop at end
        if queue.CurrentIndex >= len(queue.Tracks)-1 {
            return nil, nil // End of queue
        }
        queue.CurrentIndex++
    }

    if err := s.saveQueue(ctx, queue); err != nil {
        return nil, err
    }

    return &queue.Tracks[queue.CurrentIndex], nil
}

func (s *MusicQueueService) saveQueue(ctx context.Context, queue *MusicQueue) error {
    key := fmt.Sprintf("music_queue:%s:%s", queue.UserID, queue.DeviceID)
    data, _ := json.Marshal(queue)
    return s.cache.Set(ctx, key, data, 24*time.Hour)
}
```

---

## Podcast Features

### Skip Segments

```go
// Skip intro/outro/ads
type SkipSegment struct {
    Type     string `json:"type"` // "intro", "outro", "sponsor", "silence"
    StartMs  int64  `json:"start_ms"`
    EndMs    int64  `json:"end_ms"`
    Source   string `json:"source"` // "user", "community", "ai"
}

func (s *PodcastService) GetSkipSegments(ctx context.Context, episodeID uuid.UUID) ([]SkipSegment, error) {
    // Check database for user/community segments
    segments, err := s.repo.GetSkipSegments(ctx, episodeID)
    if err != nil {
        return nil, err
    }

    // Could also integrate with SponsorBlock-like service
    return segments, nil
}
```

### Episode Played Status

```go
func (s *PodcastService) MarkPlayed(ctx context.Context, userID, episodeID uuid.UUID) error {
    return s.db.Exec(ctx, `
        INSERT INTO podcast_playback_progress (user_id, episode_id, is_completed, position_ms, duration_ms)
        SELECT $1, $2, true, duration_ms, duration_ms
        FROM podcast_episodes WHERE id = $2
        ON CONFLICT (user_id, episode_id) DO UPDATE SET
            is_completed = true,
            position_ms = podcast_playback_progress.duration_ms
    `, userID, episodeID)
}

func (s *PodcastService) MarkUnplayed(ctx context.Context, userID, episodeID uuid.UUID) error {
    return s.db.Exec(ctx, `
        UPDATE podcast_playback_progress
        SET is_completed = false, position_ms = 0
        WHERE user_id = $1 AND episode_id = $2
    `, userID, episodeID)
}
```

---

## Audiobook Features

### Chapter Navigation

```go
type AudiobookChapter struct {
    Index    int    `json:"index"`
    Title    string `json:"title"`
    StartMs  int64  `json:"start_ms"`
    EndMs    int64  `json:"end_ms"`
    Duration int64  `json:"duration"`
}

func (s *AudiobookService) GetChapters(ctx context.Context, audiobookID uuid.UUID) ([]AudiobookChapter, error) {
    // From database (parsed from file metadata)
    return s.repo.GetChapters(ctx, audiobookID)
}

func (s *AudiobookService) SeekToChapter(ctx context.Context, userID, audiobookID uuid.UUID, chapterIndex int) error {
    chapters, err := s.GetChapters(ctx, audiobookID)
    if err != nil {
        return err
    }

    if chapterIndex < 0 || chapterIndex >= len(chapters) {
        return errors.New("invalid chapter index")
    }

    chapter := chapters[chapterIndex]

    // Update progress
    return s.db.Exec(ctx, `
        UPDATE audiobook_playback_progress
        SET current_chapter = $3, position_ms = $4, total_position_ms = $4
        WHERE user_id = $1 AND audiobook_id = $2
    `, userID, audiobookID, chapterIndex, chapter.StartMs)
}
```

### Bookmarks

```go
type Bookmark struct {
    ID        uuid.UUID `json:"id"`
    Chapter   int       `json:"chapter"`
    PositionMs int64    `json:"position_ms"`
    Note      string    `json:"note,omitempty"`
    CreatedAt time.Time `json:"created_at"`
}

func (s *AudiobookService) AddBookmark(ctx context.Context, userID, audiobookID uuid.UUID, chapter int, positionMs int64, note string) error {
    bookmark := Bookmark{
        ID:        uuid.New(),
        Chapter:   chapter,
        PositionMs: positionMs,
        Note:      note,
        CreatedAt: time.Now(),
    }

    _, err := s.db.Exec(ctx, `
        UPDATE audiobook_playback_progress
        SET bookmarks = bookmarks || $3::jsonb
        WHERE user_id = $1 AND audiobook_id = $2
    `, userID, audiobookID, bookmark)

    return err
}

func (s *AudiobookService) GetBookmarks(ctx context.Context, userID, audiobookID uuid.UUID) ([]Bookmark, error) {
    var bookmarks []Bookmark
    err := s.db.QueryRow(ctx, `
        SELECT bookmarks FROM audiobook_playback_progress
        WHERE user_id = $1 AND audiobook_id = $2
    `, userID, audiobookID).Scan(&bookmarks)
    return bookmarks, err
}
```

### Sleep Timer

```go
type SleepTimer struct {
    UserID     uuid.UUID     `json:"user_id"`
    DeviceID   string        `json:"device_id"`
    Mode       SleepTimerMode `json:"mode"`
    EndTime    *time.Time    `json:"end_time,omitempty"`     // For duration mode
    EndChapter *int          `json:"end_chapter,omitempty"`  // For chapter mode
    CreatedAt  time.Time     `json:"created_at"`
}

type SleepTimerMode string

const (
    SleepTimerDuration   SleepTimerMode = "duration"    // Stop after X minutes
    SleepTimerChapter    SleepTimerMode = "chapter"     // Stop at end of chapter
    SleepTimerChapters   SleepTimerMode = "chapters"    // Stop after N chapters
)

func (s *AudiobookService) SetSleepTimer(ctx context.Context, userID uuid.UUID, deviceID string, mode SleepTimerMode, value int) error {
    timer := SleepTimer{
        UserID:    userID,
        DeviceID:  deviceID,
        Mode:      mode,
        CreatedAt: time.Now(),
    }

    switch mode {
    case SleepTimerDuration:
        endTime := time.Now().Add(time.Duration(value) * time.Minute)
        timer.EndTime = &endTime
    case SleepTimerChapter:
        timer.EndChapter = &value
    }

    key := fmt.Sprintf("sleep_timer:%s:%s", userID, deviceID)
    data, _ := json.Marshal(timer)
    return s.cache.Set(ctx, key, data, 24*time.Hour)
}
```

---

## API Endpoints

```yaml
# api/openapi/audio.yaml
paths:
  # Streaming
  /api/v1/audio/{itemType}/{id}/stream:
    get:
      summary: Stream audio content
      parameters:
        - name: itemType
          in: path
          schema:
            type: string
            enum: [track, episode, audiobook]
        - name: id
          in: path
          schema:
            type: string
            format: uuid
        - name: position
          in: query
          schema:
            type: integer
            description: Start position in milliseconds

  # Progress
  /api/v1/audio/progress:
    post:
      summary: Update playback progress
      requestBody:
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/ProgressUpdate'

  /api/v1/audio/{itemType}/{id}/progress:
    get:
      summary: Get resume position

  # Music queue
  /api/v1/music/queue:
    get:
      summary: Get current queue
    post:
      summary: Create/replace queue
    patch:
      summary: Modify queue (add, remove, reorder)

  /api/v1/music/queue/next:
    post:
      summary: Skip to next track

  /api/v1/music/queue/previous:
    post:
      summary: Go to previous track

  /api/v1/music/queue/shuffle:
    post:
      summary: Toggle shuffle mode

  /api/v1/music/queue/repeat:
    post:
      summary: Set repeat mode

  # Audiobook
  /api/v1/audiobooks/{id}/chapters:
    get:
      summary: Get chapter list

  /api/v1/audiobooks/{id}/bookmarks:
    get:
      summary: Get bookmarks
    post:
      summary: Add bookmark

  /api/v1/audiobooks/{id}/bookmarks/{bookmarkId}:
    delete:
      summary: Remove bookmark

  /api/v1/audiobooks/{id}/sleep-timer:
    post:
      summary: Set sleep timer
    delete:
      summary: Cancel sleep timer

  # Podcast
  /api/v1/podcasts/{podcastId}/episodes/{episodeId}/played:
    post:
      summary: Mark episode as played
    delete:
      summary: Mark episode as unplayed
```

---

## Configuration

```yaml
# configs/config.yaml
audio:
  # Progress tracking
  progress:
    sync_interval: 30s            # How often clients should sync
    batch_interval: 10s           # Server-side batching interval
    completion_threshold: 0.9     # 90% = completed (podcasts/audiobooks)
    music_scrobble_threshold: 0.5 # 50% = scrobbled (music)

  # Transcoding
  transcoding:
    enabled: true
    default_codec: aac
    default_bitrate: 256000       # 256 kbps
    max_bitrate: 320000           # 320 kbps

  # Queue
  queue:
    max_size: 10000               # Max tracks in queue
    history_size: 50              # Tracks to keep in history
```

---

## Summary

| Feature | Implementation |
|---------|----------------|
| Audio Streaming | Direct play + Blackbeard transcode |
| Progress Tracking | Per-module tables with batching |
| Music Queue | Dragonfly-backed with shuffle/repeat |
| Podcast Features | Skip segments, played/unplayed |
| Audiobook Features | Chapters, bookmarks, sleep timer |
| Session Management | Cross-device sync via cache |

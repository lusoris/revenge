# Revenge - User Experience Features

> User-facing features inspired by modern streaming services.

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Full design with feature matrix, DB schemas, UI patterns |
| Sources | ✅ | Netflix, Spotify, Disney+, YouTube features documented |
| Instructions | ✅ | Implementation checklist complete |
| Code | 🔴 | |
| Linting | 🔴 | |
| Unit Testing | 🔴 | |
| Integration Testing | 🔴 | |

**Last Updated**: 2026-01-30
**Location**: `internal/ux/`

---

## Overview

This document catalogs **user experience features** for Revenge. Focus is on UI/UX - all transcoding is handled by Blackbeard.

**Design Principle:** See [DESIGN_PRINCIPLES.md](DESIGN_PRINCIPLES.md) for architectural decisions.

---

## Feature Matrix

| Feature | Netflix | Spotify | Disney+ | YouTube | Revenge |
|---------|---------|---------|---------|---------|---------|
| Skip Intro/Outro | ✅ | — | ✅ | ✅ | ✅ |
| Continue Watching | ✅ | ✅ | ✅ | ✅ | ✅ |
| "Still Watching?" | ✅ | — | ✅ | — | ✅ |
| User Profiles | ✅ | ✅ | ✅ | ✅ | ✅ |
| Kids Mode | ✅ | ✅ | ✅ | ✅ | ✅ |
| Watch Party | ✅ | — | ✅ | ✅ | ✅ |
| Downloads | ✅ | ✅ | ✅ | ✅ | ✅ |
| Playback Speed | — | — | — | ✅ | ✅ |
| Sleep Timer | ✅ | ✅ | — | — | ✅ |
| Random Episode | — | ✅ | — | — | ✅ |
| Crossfade | — | ✅ | — | — | ✅ |
| Gapless Playback | — | ✅ | — | — | ✅ |
| Lyrics Sync | — | ✅ | — | — | ✅ |
| Chapters | — | ✅ | — | ✅ | ✅ |
| Picture-in-Picture | ✅ | — | ✅ | ✅ | ✅ |

---

## Video Features

### Skip Intro / Recap / Outro

**Source:** Netflix, Disney+

Automatically detect and offer skip buttons for repetitive content.

| Marker | Description | Skip Button |
|--------|-------------|-------------|
| Intro | Opening credits/theme | "Skip Intro" |
| Recap | "Previously on..." | "Skip Recap" |
| Outro | End credits | "Skip to Next" |
| Post-credits | After credits scene | Alert: "Post-credits scene!" |

```sql
-- Content markers for skip functionality
CREATE TABLE content_markers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    media_item_id UUID NOT NULL,
    marker_type VARCHAR(50) NOT NULL,  -- intro, recap, outro, post_credits
    start_ms BIGINT NOT NULL,
    end_ms BIGINT NOT NULL,
    detection_method VARCHAR(50),       -- manual, audio_fingerprint, community
    confidence FLOAT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

**Detection Methods:**
1. **Audio fingerprinting** - Compare audio across episodes to find recurring intros
2. **Community contributions** - Users submit markers, verified by votes
3. **Manual entry** - Admin/user can set markers manually

---

### Continue Watching

**Source:** All streaming services

Smart resume with context awareness.

```go
type ContinueWatching struct {
    MediaID       uuid.UUID
    Position      time.Duration
    Duration      time.Duration
    Progress      float64        // 0.0 - 1.0
    LastWatched   time.Time

    // Context
    NextEpisode   *MediaItem     // For TV shows
    SeasonNumber  int
    EpisodeNumber int

    // Display
    ThumbnailURL  string         // Frame at position
    TimeRemaining time.Duration
}
```

**Rules:**
- Show if progress > 5% and < 95%
- Sort by last watched
- Remove after 30 days of inactivity
- Auto-advance to next episode at 95%

---

### "Are You Still Watching?"

**Source:** Netflix

Pause playback after inactivity to save bandwidth and track accurate watch history.

```yaml
playback:
  still_watching:
    enabled: true
    trigger_after_episodes: 3      # Ask after X consecutive episodes
    trigger_after_hours: 4         # Or after X hours
    inactivity_minutes: 90         # No input for 90 min

    # Exceptions
    skip_for_music: true
    skip_for_audiobooks: true
```

---

### Playback Speed Control

**Source:** YouTube, Podcasts apps

Variable playback speed for all content types.

| Content Type | Speed Options |
|--------------|---------------|
| Movies/TV | 0.5x, 0.75x, 1x, 1.25x, 1.5x |
| Podcasts | 0.5x, 0.75x, 1x, 1.25x, 1.5x, 1.75x, 2x, 2.5x, 3x |
| Audiobooks | 0.5x, 0.75x, 1x, 1.1x, 1.25x, 1.5x, 2x |
| Music | 1x only (pitch preservation) |

**Features:**
- Remember speed per content type
- Pitch correction for speech
- Silence trimming for podcasts (smart speed)

---

### Sleep Timer

**Source:** Netflix, Spotify, Podcast apps

Auto-stop playback after duration or content completion.

```go
type SleepTimer struct {
    Mode     string        // "duration", "end_of_episode", "end_of_album"
    Duration time.Duration // For duration mode
    FadeOut  bool          // Gradually lower volume
    FadeTime time.Duration // 30 seconds default
}
```

**Options:**
- 15 min, 30 min, 45 min, 1 hour, 2 hours
- End of current episode/track
- End of current season/album
- Custom duration

---

### Picture-in-Picture (PiP)

**Source:** YouTube, Netflix

Continue watching while browsing.

```typescript
// Client-side PiP support
interface PiPConfig {
  enabled: boolean;
  position: 'bottom-right' | 'bottom-left' | 'top-right' | 'top-left';
  size: 'small' | 'medium' | 'large';
  autoEnable: boolean;  // Enable when navigating away
  showControls: boolean;
}
```

---

### Random Episode / Shuffle Play

**Source:** Netflix "Surprise Me", Spotify Shuffle

Play random content from a series or playlist.

| Mode | Description |
|------|-------------|
| Shuffle Series | Random episode from selected series |
| Shuffle Watched | Random from already-watched episodes |
| Shuffle Unwatched | Random unwatched episode |
| True Shuffle | Mathematically random |
| Smart Shuffle | Avoids recently played (Spotify-style) |

```go
type ShuffleConfig struct {
    Mode           string   // true_shuffle, smart_shuffle
    AvoidRecent    int      // Don't repeat last N items
    WeightUnplayed bool     // Prefer unplayed content
    Seed           int64    // For reproducible shuffles
}
```

---

## Audio Features (Spotify-inspired)

### Gapless Playback

**Source:** Spotify, Apple Music

Seamless transitions between tracks, especially for live albums and classical music.

```go
type GaplessConfig struct {
    Enabled           bool
    PreloadNextTrack  time.Duration  // Preload 10s before end
    CrossfadeOverride bool           // Disable crossfade for gapless albums
}
```

**Album types that need gapless:**
- Live albums
- Classical music (movements)
- Concept albums (Pink Floyd, etc.)
- DJ mixes

---

### Crossfade

**Source:** Spotify

Smooth transitions between tracks with configurable overlap.

```yaml
audio:
  crossfade:
    enabled: true
    duration_ms: 5000        # 5 second crossfade
    curve: equal_power       # linear, equal_power, logarithmic

    # Smart crossfade
    disable_for_gapless: true
    disable_for_same_album: true
    auto_adjust_by_genre: true  # Shorter for classical, longer for party playlists
```

---

### Volume Normalization

**Source:** Spotify, Apple Music

Consistent volume across tracks to avoid jarring level changes.

| Mode | Description |
|------|-------------|
| Off | No normalization |
| Normal | Target -14 LUFS (Spotify default) |
| Quiet | Target -23 LUFS (quiet environments) |
| Loud | Target -11 LUFS (noisy environments) |

```go
type NormalizationConfig struct {
    Enabled    bool
    TargetLUFS float64  // -14 LUFS default
    Mode       string   // album, track
    // Album mode: normalize album as unit
    // Track mode: normalize each track individually
}
```

---

### Synced Lyrics

**Source:** Spotify, Apple Music

Display lyrics synchronized to playback.

```sql
CREATE TABLE track_lyrics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    track_id UUID NOT NULL REFERENCES tracks(id),
    language VARCHAR(10) DEFAULT 'en',

    -- Sync types
    sync_type VARCHAR(20) NOT NULL,  -- line, word, unsynced

    -- Content
    lyrics_json JSONB NOT NULL,  -- [{time_ms, text}, ...]

    -- Source
    source VARCHAR(50),  -- lrclib, musixmatch, genius, manual

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

**Lyrics sources:**
1. LRCLib (free, community)
2. Embedded in file (ID3 USLT/SYLT)
3. Manual upload (.lrc files)

---

### Audio Quality Selection

**Source:** Spotify, Tidal

Let users choose streaming quality.

| Quality | Bitrate | Use Case |
|---------|---------|----------|
| Low | 96 kbps | Save data |
| Normal | 160 kbps | Default mobile |
| High | 320 kbps | Default desktop |
| Lossless | FLAC 16/44.1 | Audiophile |
| Hi-Res | FLAC 24/96+ | Premium |

```yaml
audio:
  quality:
    wifi: high
    cellular: normal
    download: lossless
```

---

## Discovery & Recommendations

### Smart Home Screen

**Source:** Netflix, Spotify

Personalized home with dynamic rows.

| Row Type | Content |
|----------|---------|
| Continue Watching | In-progress media |
| Because You Watched X | Similar content |
| New in Library | Recently added |
| Trending | Popular with all users |
| Top Picks for You | ML recommendations |
| Genre: Action | Genre-based row |
| Mood: Chill | Mood-based (music) |

---

### "More Like This"

**Source:** All services

Related content recommendations.

```go
type SimilarityFactors struct {
    Genre      float64  // Weight: 0.3
    Cast       float64  // Weight: 0.2
    Director   float64  // Weight: 0.1
    Year       float64  // Weight: 0.1
    Rating     float64  // Weight: 0.1
    UserTaste  float64  // Weight: 0.2 (collaborative filtering)
}
```

---

### Collections & Lists

**Source:** Spotify playlists, Netflix "My List"

```sql
-- User collections (watchlist, favorites, custom lists)
CREATE TABLE user_collections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    collection_type VARCHAR(50) NOT NULL,  -- watchlist, favorites, custom
    is_public BOOLEAN DEFAULT false,
    cover_image_url VARCHAR(1024),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE user_collection_items (
    collection_id UUID NOT NULL REFERENCES user_collections(id),
    media_item_id UUID NOT NULL,
    position INT NOT NULL,
    added_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (collection_id, media_item_id)
);
```

---

## Social Features

### User Profiles

**Source:** Netflix, Disney+

Multiple profiles per account with separate preferences.

| Feature | Description |
|---------|-------------|
| Avatar | Custom upload, Gravatar, or generated |
| Name | Display name |
| Language | UI and content language |
| Maturity | Content rating filter |
| Autoplay | Next episode autoplay |
| History | Separate watch history |
| Recommendations | Personalized per profile |

#### Avatar Sources

Avatars use a cascading fallback system:

| Priority | Source | Description |
|----------|--------|-------------|
| 1 | Custom Upload | User-uploaded image stored in `/media/avatars/{user_id}.{ext}` |
| 2 | Gravatar | Email hash lookup: `https://gravatar.com/avatar/{md5(email)}?d=404` |
| 3 | DiceBear | Generated avatar: `https://api.dicebear.com/7.x/bottts/svg?seed={username}` |

**DiceBear Styles** (configurable per server):
- `bottts` - Robot avatars (default)
- `avataaars` - Cartoon people
- `identicon` - Geometric patterns
- `pixel-art` - Retro pixel art
- `shapes` - Abstract shapes

**Implementation:**
```go
func (u *User) AvatarURL() string {
    // 1. Custom upload
    if u.AvatarPath != "" {
        return fmt.Sprintf("/media/avatars/%s", u.AvatarPath)
    }
    // 2. Gravatar (if email exists and Gravatar returns 200)
    if u.Email != "" {
        hash := md5.Sum([]byte(strings.ToLower(strings.TrimSpace(u.Email))))
        return fmt.Sprintf("https://gravatar.com/avatar/%x?d=404", hash)
    }
    // 3. DiceBear fallback
    return fmt.Sprintf("https://api.dicebear.com/7.x/bottts/svg?seed=%s", u.Username)
}
```

**Frontend Handling:**
```typescript
function getAvatarUrl(user: User): string {
  // Try custom avatar first
  if (user.avatar_path) return `/media/avatars/${user.avatar_path}`;
  // Then Gravatar with DiceBear fallback in URL
  const emailHash = md5(user.email?.toLowerCase().trim() || '');
  const dicebear = encodeURIComponent(`https://api.dicebear.com/7.x/bottts/svg?seed=${user.username}`);
  return `https://gravatar.com/avatar/${emailHash}?d=${dicebear}`;
}
```

```sql
CREATE TABLE profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    name VARCHAR(100) NOT NULL,
    avatar_path VARCHAR(255),        -- Custom upload path (NULL = use fallback)
    avatar_url VARCHAR(1024),        -- Computed/cached avatar URL
    is_kids BOOLEAN DEFAULT false,
    max_maturity_rating VARCHAR(20),  -- G, PG, PG-13, R, NC-17
    language VARCHAR(10) DEFAULT 'en',
    autoplay_next BOOLEAN DEFAULT true,
    autoplay_previews BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

### Kids Mode

**Source:** Netflix Kids, Disney+, YouTube Kids

Safe viewing environment for children.

| Feature | Description |
|---------|-------------|
| Content Filter | Only G/PG rated content |
| Simplified UI | Big tiles, simple navigation |
| No Search | Curated content only |
| PIN Exit | Require PIN to exit kids mode |
| Time Limits | Optional viewing time limits |
| No Ads | Remove any promotional content |

---

### Watch Party

**Source:** Disney+ GroupWatch, Amazon Watch Party, Teleparty

Synchronized viewing with friends.

```go
type WatchParty struct {
    ID           uuid.UUID
    HostUserID   uuid.UUID
    MediaItemID  uuid.UUID

    // Sync state
    PlaybackState  string        // playing, paused
    Position       time.Duration
    LastSyncAt     time.Time

    // Participants
    Participants   []PartyMember
    MaxParticipants int          // Default 8

    // Features
    ChatEnabled    bool
    ReactionsEnabled bool
}

type PartyMember struct {
    UserID    uuid.UUID
    Name      string
    AvatarURL string
    IsHost    bool
    JoinedAt  time.Time
}
```

**Sync Protocol (WebSocket):**
```json
{
  "type": "sync",
  "position_ms": 125000,
  "state": "playing",
  "timestamp": "2026-01-28T12:00:00Z"
}
```

---

### Activity Feed

**Source:** Spotify Friend Activity, Letterboxd

See what friends are watching/listening to.

```sql
CREATE TABLE activity_feed (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID NOT NULL REFERENCES profiles(id),
    activity_type VARCHAR(50) NOT NULL,  -- watching, finished, rated, added_to_list
    media_item_id UUID NOT NULL,
    -- ENCRYPTED activity metadata (AES-256-GCM)
    activity_data BYTEA NOT NULL,  -- Encrypted JSON: {rating, list_name, etc.}
    is_public BOOLEAN DEFAULT false,  -- Privacy by default
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Encryption key derived from user's master key
-- Master key encrypted with user password (see DESIGN_PRINCIPLES.md)
```

**Privacy controls:**
- Share nothing (default)
- Share with friends only
- Share publicly
- Exclude specific content types (adult, specific genres)

---

### Year in Review / Wrapped

**Source:** Spotify Wrapped, Apple Music Replay

Annual statistics and highlights.

| Stat | Description |
|------|-------------|
| Total Watch Time | Hours watched this year |
| Top Genres | Most watched genres |
| Top Shows | Most binged series |
| Top Artists | Most listened artists |
| Listening Streak | Longest consecutive days |
| Night Owl | Late night viewing stats |
| Binge Sessions | Marathon viewing counts |

---

## Downloads & Offline

### Smart Downloads

**Source:** Netflix, Spotify

Automatic download management.

| Feature | Description |
|---------|-------------|
| Auto-download | Download next episodes automatically |
| Auto-delete | Remove watched downloads |
| WiFi-only | Download only on WiFi |
| Storage limit | Max GB for downloads |
| Quality | Download quality setting |
| Expiry | Downloads expire after X days |

```yaml
downloads:
  enabled: true
  wifi_only: true
  storage_limit_gb: 10
  quality: high

  smart_downloads:
    enabled: true
    max_episodes: 3      # Keep 3 episodes ahead
    auto_delete_watched: true

  expiry:
    enabled: true
    days: 30             # Downloads expire after 30 days
    renew_on_connect: true  # Refresh expiry when online
```

---

## Accessibility

### Audio Descriptions

**Source:** Netflix, Disney+

Narrated descriptions of visual elements for blind/low-vision users.

```sql
CREATE TABLE audio_descriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    media_item_id UUID NOT NULL,
    language VARCHAR(10) NOT NULL,
    audio_track_index INT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

### Subtitle Customization

**Source:** Netflix, YouTube

Fully customizable subtitle appearance.

| Setting | Options |
|---------|---------|
| Font | System, Netflix Sans, Arial, etc. |
| Size | 50% - 200% |
| Color | White, Yellow, Cyan, Green, etc. |
| Background | Transparent, Semi-transparent, Opaque |
| Edge Style | None, Drop Shadow, Raised, Depressed |
| Position | Bottom, Top, Custom |

```sql
CREATE TABLE user_subtitle_preferences (
    user_id UUID PRIMARY KEY REFERENCES users(id),
    font_family VARCHAR(100) DEFAULT 'system',
    font_size INT DEFAULT 100,  -- percentage
    font_color VARCHAR(20) DEFAULT '#FFFFFF',
    background_color VARCHAR(20) DEFAULT '#000000',
    background_opacity FLOAT DEFAULT 0.75,
    edge_style VARCHAR(20) DEFAULT 'drop_shadow',
    position VARCHAR(20) DEFAULT 'bottom'
);
```

---

### Reduced Motion

**Source:** iOS, modern web

Option to reduce animations for vestibular disorders.

```yaml
accessibility:
  reduced_motion: false
  high_contrast: false
  screen_reader_optimized: false
```

---

## Quality of Life

### Intro Preview / Autoplay Previews

**Source:** Netflix (on hover), Disney+

Preview content without clicking.

| Feature | Behavior |
|---------|----------|
| Hover preview | Play preview after 2s hover |
| Sound | Muted by default |
| Duration | 30s max |
| Disable option | User preference to disable |

---

### Content Warnings

**Source:** Netflix, Disney+, HBO Max

Content-specific warnings before playback.

| Warning | Description |
|---------|-------------|
| Flashing Lights | Photosensitive epilepsy |
| Violence | Graphic violence |
| Language | Strong language |
| Smoking | Tobacco depiction |
| Outdated Content | "Dated cultural depictions" |

---

### Parental Controls

**Source:** All services

Comprehensive content restrictions.

```sql
CREATE TABLE parental_controls (
    user_id UUID PRIMARY KEY REFERENCES users(id),
    pin_hash VARCHAR(255) NOT NULL,
    max_rating VARCHAR(20),  -- G, PG, PG-13, R
    blocked_titles UUID[],   -- Specific blocked content
    viewing_restrictions JSONB,  -- Time limits, etc.
    require_pin_for_changes BOOLEAN DEFAULT true
);
```

---

## Implementation Priority

### Phase 1: Core Playback
- [x] Continue watching
- [x] Skip intro/outro markers
- [x] Playback speed control
- [x] Sleep timer

### Phase 2: Profiles & Personalization
- [ ] User profiles
- [ ] Kids mode
- [ ] Watch history per profile
- [ ] Basic recommendations

### Phase 3: Audio
- [ ] Gapless playback
- [ ] Crossfade
- [ ] Volume normalization
- [ ] Synced lyrics

### Phase 4: Social
- [ ] Watch party (WebSocket)
- [ ] Activity feed
- [ ] Year in review

### Phase 5: Downloads
- [ ] Offline downloads
- [ ] Smart downloads
- [ ] Download management

---

## Summary

| Category | Features |
|----------|----------|
| **Playback** | Skip intro, continue watching, speed control, sleep timer, PiP |
| **Audio** | Gapless, crossfade, normalization, lyrics |
| **Discovery** | Home rows, recommendations, "more like this" |
| **Social** | Profiles, watch party, activity feed, wrapped |
| **Downloads** | Offline viewing, smart downloads |
| **Accessibility** | Audio descriptions, subtitle styling, reduced motion |
| **Parental** | Kids mode, PIN, content filters |

All streaming/transcoding is delegated to Blackbeard. Revenge focuses on the user experience layer.

---

## Implementation Checklist

### Phase 1: Core Infrastructure
- [ ] Create package structure at `internal/ux/`
- [ ] Create sub-packages: `profiles/`, `playback/`, `preferences/`, `social/`
- [ ] Define profile entity (`profiles/entity.go`)
- [ ] Define user preferences entity
- [ ] Define content marker entity (intro/outro)
- [ ] Define collection entity (watchlist, favorites)
- [ ] Create repository interfaces
- [ ] Implement PostgreSQL repositories
- [ ] Create fx module (`module.go`)
- [ ] Add configuration structs for all features

### Phase 2: Database
- [ ] Create migration for `profiles` table
- [ ] Create migration for `content_markers` table
- [ ] Create migration for `track_lyrics` table
- [ ] Create migration for `user_collections` table
- [ ] Create migration for `user_collection_items` table
- [ ] Create migration for `user_subtitle_preferences` table
- [ ] Create migration for `parental_controls` table
- [ ] Create migration for `activity_feed` table (encrypted)
- [ ] Add indexes for profile lookups, collection queries
- [ ] Write sqlc queries for profile CRUD
- [ ] Write sqlc queries for preferences management
- [ ] Write sqlc queries for collections
- [ ] Write sqlc queries for activity feed

### Phase 3: Service Layer
- [ ] Implement ProfileService
  - [ ] Profile CRUD operations
  - [ ] Avatar handling (upload, Gravatar, DiceBear fallback)
  - [ ] Kids mode restrictions
  - [ ] Maturity rating enforcement
- [ ] Implement PlaybackService
  - [ ] Continue watching logic (5%-95% rules)
  - [ ] Skip intro/outro detection and UI triggers
  - [ ] Playback speed per content type
  - [ ] Sleep timer functionality
  - [ ] "Are you still watching?" prompts
- [ ] Implement PreferencesService
  - [ ] Subtitle customization
  - [ ] Audio quality selection
  - [ ] Autoplay settings
  - [ ] Reduced motion / accessibility
- [ ] Implement SocialService
  - [ ] Watch party sync (WebSocket)
  - [ ] Activity feed (with encryption)
  - [ ] Year-in-review stats calculation
- [ ] Add caching for frequently accessed preferences

### Phase 4: Background Jobs
- [ ] Create River job for content marker detection (audio fingerprinting)
- [ ] Create River job for "continue watching" cleanup (30-day inactivity)
- [ ] Create River job for year-in-review generation
- [ ] Create River job for activity feed pruning
- [ ] Create River job for smart download management
- [ ] Register jobs in fx module

### Phase 5: API Integration
- [ ] Add OpenAPI spec for profile endpoints
- [ ] Add OpenAPI spec for playback preference endpoints
- [ ] Add OpenAPI spec for collection endpoints
- [ ] Add OpenAPI spec for social feature endpoints
- [ ] Generate ogen handlers
- [ ] Implement profile endpoints:
  - [ ] GET/POST/PUT/DELETE /api/v1/profiles
  - [ ] POST /api/v1/profiles/{id}/avatar
- [ ] Implement playback endpoints:
  - [ ] GET /api/v1/continue-watching
  - [ ] PUT /api/v1/playback/progress
  - [ ] GET /api/v1/content/{id}/markers
- [ ] Implement collection endpoints:
  - [ ] GET/POST /api/v1/collections
  - [ ] POST/DELETE /api/v1/collections/{id}/items
- [ ] Implement social endpoints:
  - [ ] WebSocket /api/v1/watch-party/{id}
  - [ ] GET /api/v1/activity-feed
  - [ ] GET /api/v1/wrapped/{year}
- [ ] Add authentication middleware
- [ ] Add profile-based authorization

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Shared](INDEX.md)

### In This Section

- [Time-Based Access Controls](ACCESS_CONTROLS.md)
- [Tracearr Analytics Service](ANALYTICS_SERVICE.md)
- [Revenge - Client Support & Device Capabilities](CLIENT_SUPPORT.md)
- [Content Rating System](CONTENT_RATING.md)
- [Revenge - Internationalization (i18n)](I18N.md)
- [Library Types](LIBRARY_TYPES.md)
- [News System](NEWS_SYSTEM.md)
- [Revenge - NSFW Toggle](NSFW_TOGGLE.md)

### Related Topics

- [Revenge - Architecture v2](../../architecture/01_ARCHITECTURE.md) _Architecture_
- [Revenge - Design Principles](../../architecture/02_DESIGN_PRINCIPLES.md) _Architecture_
- [Revenge - Metadata System](../../architecture/03_METADATA_SYSTEM.md) _Architecture_
- [Revenge - Player Architecture](../../architecture/04_PLAYER_ARCHITECTURE.md) _Architecture_
- [Plugin Architecture Decision](../../architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) _Architecture_

### Indexes

- [Design Index](../../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

## Related

- [Design Principles](DESIGN_PRINCIPLES.md) - Architectural decisions
- [Client Support](CLIENT_SUPPORT.md) - Device capabilities and streaming
- [Voice Control](VOICE_CONTROL.md) - Voice assistant integration

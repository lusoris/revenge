# Revenge - Design Principles

> Unverhandelbare Architektur-Prinzipien für das gesamte Projekt.

## Core Principles

### 1. Performance First

**UX darf niemals blockiert werden.**

| Rule | Implementation |
|------|----------------|
| No blocking I/O in request handlers | Async via River Jobs |
| No heavy computation in hot path | Offload to Blackbeard / background |
| Database queries optimized | Indexes, prepared statements, connection pooling |
| Caching aggressive | Dragonfly for hot data |

```go
// ❌ WRONG - Blocks request
func (h *Handler) GetMovie(w http.ResponseWriter, r *http.Request) {
    metadata := h.tmdb.FetchMetadata(id)  // Blocks!
    // ...
}

// ✅ RIGHT - Returns cached, triggers background refresh
func (h *Handler) GetMovie(w http.ResponseWriter, r *http.Request) {
    movie := h.cache.GetMovie(id)  // Fast
    if movie.NeedsRefresh() {
        h.jobs.Enqueue(RefreshMetadataJob{ID: id})  // Async
    }
    // ...
}
```

### 2. Client Agnostic

**Keine eigenen Clients entwickeln. Fremdclients unterstützen.**

| Client Type | Support Strategy |
|-------------|------------------|
| Web | SvelteKit WebUI (einzige Eigenentwicklung) |
| Mobile | Jellyfin/Infuse/VLC via kompatible API |
| TV | Jellyfin TV Apps, Kodi, Plex-kompatibel |
| Desktop | VLC, mpv, IINA via Direct Play |

**API Compatibility Layers:**
- Jellyfin-compatible endpoints für existierende Apps
- Subsonic API für Music Apps (DSub, Ultrasonic)
- DLNA/UPnP für Smart TVs
- Chromecast/AirPlay Support

```yaml
api:
  compatibility:
    jellyfin: true      # Jellyfin client support
    subsonic: true      # Music app support
    dlna: true          # Smart TV support
```

### 3. Privacy by Default, Features by Choice

**Minimales Tracking, maximale Kontrolle.**

| Data | Storage | Encryption | Purpose |
|------|---------|------------|---------|
| Watch History | Local DB | At-rest (AES-256) | Continue Watching, Statistics |
| Play Events | Local DB | At-rest | Year in Review, Recommendations |
| User Preferences | Local DB | At-rest | Personalization |

**Rules:**
- Alle User-Daten verschlüsselt at-rest
- Keine externe Telemetrie ohne explizites Opt-in
- Keine Cloud-Calls ohne User-Consent
- Export aller eigenen Daten jederzeit möglich (GDPR)

```sql
-- Encrypted user activity storage
CREATE TABLE user_activity (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    -- Encrypted blob containing activity details
    activity_data BYTEA NOT NULL,  -- AES-256-GCM encrypted JSON
    activity_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Encryption key derived from user's master key
-- Master key encrypted with user password
```

### 4. Bleeding Edge, Stable Core

**Neueste stabile Versionen, keine Experimente in Production.**

| Component | Version Policy | Current |
|-----------|---------------|---------|
| Go | Latest stable | 1.25 |
| PostgreSQL | Latest stable | 18+ || Frontend | SvelteKit (latest stable) | 2 |
| UI Library | Tailwind CSS (latest stable) | 4 |
| Search | Typesense (latest stable) | 0.25+ || Dependencies | Latest stable, actively maintained | See go.mod |

**Forbidden:**
- Alpha/Beta releases in production
- Deprecated packages
- Unmaintained dependencies (>1 year no commits)
- Packages with known CVEs

```go
// go.mod - Only stable, maintained packages
require (
    github.com/jackc/pgx/v5 v5.7.0      // Active, stable
    github.com/riverqueue/river v0.15.0  // Active, stable
    go.uber.org/fx v1.24.0               // Active, stable
)
```

### 5. Optional ML Integration

**Self-hosted ML nur wenn explizit konfiguriert.**

| Feature | Without ML | With ML (Ollama/etc.) |
|---------|------------|----------------------|
| Recommendations | Genre/Cast/Director matching | Collaborative filtering, embeddings |
| Search | Typesense full-text | Semantic search |
| Intro Detection | Community markers, chapters | Audio fingerprinting |

```yaml
ml:
  enabled: false
  provider: ollama  # ollama, localai, etc.
  endpoint: http://localhost:11434

  features:
    recommendations: true
    semantic_search: true
    intro_detection: false  # CPU-intensive
```

**No ML = Full functionality.** ML enhances, never requires.

### 6. Resource-Aware Background Tasks

**Heavy tasks nur bei verfügbaren Ressourcen.**

| Task Priority | Condition | Examples |
|---------------|-----------|----------|
| Critical | Always run | Session cleanup, webhook delivery |
| High | Low load | Metadata refresh, image download |
| Normal | Idle | Library scan, search reindex |
| Low | Very idle + opt-in | Audio fingerprinting, ML training |

```go
// Resource-aware job scheduling
type JobPriority int

const (
    PriorityCritical JobPriority = iota  // Always
    PriorityHigh                          // Load < 70%
    PriorityNormal                        // Load < 50%
    PriorityLow                           // Load < 20%, user opt-in
)

func (s *Scheduler) ShouldRun(priority JobPriority) bool {
    load := s.GetSystemLoad()
    switch priority {
    case PriorityCritical:
        return true
    case PriorityHigh:
        return load < 0.7
    case PriorityNormal:
        return load < 0.5
    case PriorityLow:
        return load < 0.2 && s.config.LowPriorityEnabled
    }
    return false
}
```

### 7. Profile-Based Multi-User

**Ein Account, mehrere Profile (Netflix-Modell).**

| Concept | Description |
|---------|-------------|
| User | Login credentials, admin rights, billing (if any) |
| Profile | Watch history, preferences, recommendations, restrictions |

```sql
-- Parent user account
CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    is_admin BOOLEAN DEFAULT false,
    max_profiles INT DEFAULT 5,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Child profiles under user
CREATE TABLE profiles (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    avatar_url VARCHAR(1024),
    is_kids BOOLEAN DEFAULT false,
    is_default BOOLEAN DEFAULT false,
    pin_hash VARCHAR(255),  -- Optional PIN for profile
    max_maturity VARCHAR(20),  -- G, PG, PG-13, R, NC-17
    language VARCHAR(10) DEFAULT 'en',

    -- Preferences
    autoplay_next BOOLEAN DEFAULT true,
    autoplay_previews BOOLEAN DEFAULT false,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, name)
);

-- All user data tied to profile, not user
CREATE TABLE watch_history (
    profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    -- ...
);
```

### 8. WebUI Player Capabilities

**Unsere WebUI muss vollständige Player-Features haben.**

| Feature | Implementation | Priority |
|---------|----------------|----------|
| Gapless Playback | Web Audio API, preload next track | High |
| Crossfade | Web Audio API, gain nodes | High |
| Volume Normalization | ReplayGain tags + client-side gain | High |
| Picture-in-Picture | Browser PiP API | Medium |
| Chromecast | Google Cast SDK | Medium |
| AirPlay | Safari only (native) | Low |

```typescript
// WebUI Audio Player with gapless support
class GaplessPlayer {
  private audioContext: AudioContext;
  private currentSource: AudioBufferSourceNode | null = null;
  private nextBuffer: AudioBuffer | null = null;

  async preloadNext(url: string): Promise<void> {
    const response = await fetch(url);
    const arrayBuffer = await response.arrayBuffer();
    this.nextBuffer = await this.audioContext.decodeAudioData(arrayBuffer);
  }

  crossfadeTo(nextTrack: AudioBuffer, duration: number): void {
    // Fade out current, fade in next using gain nodes
  }
}
```

**Frontend Stack:**
- Framework: **SvelteKit 2**
- UI Library: **Tailwind CSS 4** + shadcn-svelte
- Player: Unified (Shaka/hls.js for video, Web Audio API for audio)
- State Management: Svelte Stores (client) + TanStack Query (server)

---

### 9. External Transcoding

**Revenge NEVER transcodes internally. All transcoding via Blackbeard.**

| Rule | Implementation |
|------|----------------|
| No FFmpeg in Revenge | Blackbeard handles all transcoding |
| Stream proxy only | Revenge proxies for access control & tracking |
| Scalable | Multiple Blackbeard instances possible |
| Replaceable | Swap transcoder without touching Revenge |

**Why:**
- **Revenge stays lean** - No heavy codec dependencies (FFmpeg ~200MB)
- **Scalable transcoding** - Spin up Blackbeard instances as needed
- **Regional deployment** - Blackbeard near storage, Revenge near users
- **GPU optimization** - Blackbeard uses hardware acceleration without affecting Revenge

**Architecture:**
```
Client → Revenge (Auth, Session, Proxy) → Blackbeard (Transcode) → Storage
         ↑                                                            ↓
         └────────────── Stream flows through Revenge ───────────────┘
```

**Blackbeard APIs (internal):**
- `POST /transcode/start` - Request transcoded stream
- `GET /transcode/{id}/master.m3u8` - HLS manifest
- `GET /transcode/{id}/{segment}.ts` - Video segment
- `DELETE /transcode/{id}` - Stop transcoding, cleanup

**Revenge APIs (client-facing):**
- `GET /stream/{sessionId}/master.m3u8` - Proxied HLS manifest
- `GET /stream/{sessionId}/{segment}.ts` - Proxied video segment
- `WebSocket /playback/{sessionId}` - Quality switching, position tracking

**Benefits:**
- Centralized access control (all streams via Revenge)
- Progress tracking (Revenge knows what client watches)
- Bandwidth monitoring (measure actual throughput)
- Session management (pause, seek, stop)

---

## Anti-Patterns (FORBIDDEN)

### ❌ Never Do This

| Anti-Pattern | Why | Alternative |
|--------------|-----|-------------|
| Sync external API in request | Blocks UX | Background job + cache |
| Store plaintext sensitive data | Security | Encrypt at-rest |
| Global state | Testing nightmare | Dependency injection |
| Panic for errors | Crashes server | Return error, handle gracefully |
| Build native mobile apps | Maintenance burden | Support existing clients |
| Require ML for basic features | Not everyone has GPU | ML enhances, core works without |
| Track without consent | Privacy violation | Opt-in only for non-essential |

---

## Decision Log

| Date | Decision | Rationale |
|------|----------|-----------|
| 2026-01-28 | No native mobile apps | Support Jellyfin/Infuse instead |
| 2026-01-28 | Profiles under Users | Family sharing like Netflix |
| 2026-01-28 | Optional ML via Ollama | Self-hosted, not required |
| 2026-01-28 | Encrypted activity tracking | Privacy + features (Wrapped) |
| 2026-01-28 | Resource-aware background jobs | Don't overload home servers |
| 2026-01-28 | WebUI with full player features | Primary interface |
| 2026-01-28 | External transcoding (Blackbeard) | Keep Revenge lean, scalable |
| 2026-01-28 | SvelteKit 2 + Tailwind CSS 4 | Modern, fast, accessible WebUI |

---

## Summary

```
┌─────────────────────────────────────────────────────────────────┐
│                    REVENGE DESIGN PRINCIPLES                    │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  1. Performance First     - Never block UX                      │
│  2. Client Agnostic       - Support Jellyfin/Subsonic clients   │
│  3. Privacy by Default    - Encrypted, local, opt-in tracking  │
│  4. Bleeding Edge Stable  - Latest stable, no alpha/beta       │
│  5. Optional ML           - Ollama integration, not required   │
│  6. Resource Aware        - Background tasks respect load      │
│  7. Profile Multi-User    - Netflix model (User → Profiles)    │
│  8. Full WebUI Player     - Gapless, crossfade, PiP, Cast      │
│  9. External Transcoding  - Delegate to Blackbeard service     │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

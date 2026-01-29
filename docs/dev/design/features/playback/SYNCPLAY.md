# SyncPlay (Watch Together)

> Synchronized playback for multiple users watching together

**Status**: 🔴 PLANNING
**Priority**: 🟢 HIGH (Critical Gap - Jellyfin has this)
**Inspired By**: Jellyfin SyncPlay

---

## Overview

SyncPlay allows multiple users to watch the same content simultaneously with synchronized playback. When one user pauses, seeks, or plays, all connected users are affected.

---

## Features

### Core Functionality

| Feature | Description |
|---------|-------------|
| Create Session | Host creates a watch party |
| Join Session | Users join via invite link/code |
| Synchronized Play/Pause | All users play/pause together |
| Synchronized Seeking | Seeking syncs across all clients |
| Latency Compensation | Adjusts for network delays |
| Chat Integration | Optional text chat during playback |

### Session Types

| Type | Description |
|------|-------------|
| Public | Anyone with link can join |
| Private | Invite-only, host approval |
| Friends Only | Only friends can join |

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    SyncPlay Server                          │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │  Session    │  │  Playback   │  │   Chat      │        │
│  │  Manager    │  │  Sync       │  │   Service   │        │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘        │
│         │                │                │                │
│         └────────────────┴────────────────┘                │
│                          │                                  │
│              ┌───────────▼───────────┐                     │
│              │    WebSocket Hub      │                     │
│              │   (gorilla/websocket) │                     │
│              └───────────┬───────────┘                     │
└──────────────────────────┼──────────────────────────────────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
   ┌────▼────┐       ┌────▼────┐       ┌────▼────┐
   │ Client 1│       │ Client 2│       │ Client 3│
   │  (Host) │       │ (Guest) │       │ (Guest) │
   └─────────┘       └─────────┘       └─────────┘
```

---

## Go Packages

| Package | Purpose | URL |
|---------|---------|-----|
| **gorilla/websocket** | WebSocket connections | github.com/gorilla/websocket |
| **centrifugal/centrifuge** | Real-time pub/sub (alternative) | github.com/centrifugal/centrifuge |

---

## Database Schema

```sql
CREATE TABLE syncplay_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(8) NOT NULL UNIQUE, -- Join code (e.g., "ABC123")
    host_user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    content_type VARCHAR(50) NOT NULL,
    content_id UUID NOT NULL,

    -- Session settings
    visibility VARCHAR(20) DEFAULT 'private', -- public, private, friends
    max_participants INT DEFAULT 10,
    chat_enabled BOOLEAN DEFAULT true,

    -- Playback state
    is_playing BOOLEAN DEFAULT false,
    position_ms BIGINT DEFAULT 0,
    playback_speed DECIMAL(3,2) DEFAULT 1.0,
    last_sync_at TIMESTAMPTZ DEFAULT NOW(),

    -- Lifecycle
    started_at TIMESTAMPTZ DEFAULT NOW(),
    ended_at TIMESTAMPTZ,

    CONSTRAINT valid_visibility CHECK (visibility IN ('public', 'private', 'friends'))
);

CREATE TABLE syncplay_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID REFERENCES syncplay_sessions(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,

    -- Connection state
    is_connected BOOLEAN DEFAULT true,
    is_ready BOOLEAN DEFAULT false, -- Has buffered enough to play
    latency_ms INT DEFAULT 0,

    joined_at TIMESTAMPTZ DEFAULT NOW(),
    left_at TIMESTAMPTZ,

    UNIQUE(session_id, user_id)
);

CREATE TABLE syncplay_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID REFERENCES syncplay_sessions(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_syncplay_sessions_code ON syncplay_sessions(code);
CREATE INDEX idx_syncplay_sessions_host ON syncplay_sessions(host_user_id);
CREATE INDEX idx_syncplay_participants_session ON syncplay_participants(session_id);
```

---

## WebSocket Protocol

### Client → Server Messages

```typescript
// Join session
{ "type": "join", "session_code": "ABC123" }

// Leave session
{ "type": "leave" }

// Playback commands (host only or all, based on settings)
{ "type": "play" }
{ "type": "pause" }
{ "type": "seek", "position_ms": 123456 }

// Ready state
{ "type": "ready", "buffered_ms": 5000 }

// Chat message
{ "type": "chat", "message": "Hello!" }

// Ping for latency measurement
{ "type": "ping", "timestamp": 1234567890 }
```

### Server → Client Messages

```typescript
// Session state update
{
    "type": "state",
    "is_playing": true,
    "position_ms": 123456,
    "server_time": 1234567890,
    "participants": [...]
}

// Playback command
{ "type": "command", "action": "play" | "pause" | "seek", "position_ms": 123456 }

// Participant update
{ "type": "participant", "action": "joined" | "left" | "ready", "user": {...} }

// Chat message
{ "type": "chat", "user": {...}, "message": "Hello!", "timestamp": 1234567890 }

// Pong for latency measurement
{ "type": "pong", "client_timestamp": 1234567890, "server_timestamp": 1234567891 }
```

---

## Latency Compensation

```go
type LatencyCompensator struct {
    participants map[uuid.UUID]*ParticipantLatency
    mu           sync.RWMutex
}

type ParticipantLatency struct {
    UserID      uuid.UUID
    Samples     []int // Last N latency samples (ms)
    AvgLatency  int   // Rolling average
    MaxLatency  int   // Max observed
}

func (c *LatencyCompensator) CalculateSyncOffset(participants []Participant) int {
    // Find the participant with highest latency
    maxLatency := 0
    for _, p := range participants {
        if p.AvgLatency > maxLatency {
            maxLatency = p.AvgLatency
        }
    }

    // All participants wait for slowest + buffer
    return maxLatency + 100 // 100ms buffer
}

// When issuing play command:
// - Calculate sync offset
// - Send "play at server_time + offset" to all clients
// - Clients with lower latency wait, others play immediately
```

---

## Go Service Implementation

```go
// internal/service/syncplay/

type Service struct {
    hub        *WebSocketHub
    sessions   SessionRepository
    messages   MessageRepository
}

type WebSocketHub struct {
    sessions   map[string]*Session
    register   chan *Client
    unregister chan *Client
    broadcast  chan *Message
    mu         sync.RWMutex
}

type Session struct {
    ID           uuid.UUID
    Code         string
    Host         *Client
    Participants map[uuid.UUID]*Client
    State        *PlaybackState
    Compensator  *LatencyCompensator
}

func (h *WebSocketHub) Run() {
    for {
        select {
        case client := <-h.register:
            h.addClient(client)
        case client := <-h.unregister:
            h.removeClient(client)
        case msg := <-h.broadcast:
            h.broadcastToSession(msg)
        }
    }
}
```

---

## API Endpoints

```
# Sessions
POST /api/v1/syncplay/sessions           # Create session
GET  /api/v1/syncplay/sessions/:code     # Get session info
DELETE /api/v1/syncplay/sessions/:code   # End session (host only)

# Participants
POST /api/v1/syncplay/sessions/:code/join   # Join session
POST /api/v1/syncplay/sessions/:code/leave  # Leave session
GET  /api/v1/syncplay/sessions/:code/participants

# WebSocket
GET  /api/v1/syncplay/ws?code=ABC123     # WebSocket connection
```

---

## RBAC Permissions

| Permission | Description |
|------------|-------------|
| `syncplay.create` | Create watch parties |
| `syncplay.join` | Join watch parties |
| `syncplay.host_controls` | Control playback as non-host |

---

## Configuration

```yaml
syncplay:
  enabled: true
  max_participants: 20
  session_timeout: 4h
  latency_samples: 10
  max_latency_ms: 5000  # Kick if latency too high
  chat:
    enabled: true
    max_message_length: 500
    rate_limit: 10/minute
```

---

## Related Documentation

- [Client Support](CLIENT_SUPPORT.md)
- [Go Packages](../architecture/GO_PACKAGES.md)
- [WebSocket Integration](../technical/WEBSOCKETS.md)

## Table of Contents

- [SyncPlay (Watch Together)](#syncplay-watch-together)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Database Schema](#database-schema)
    - [Module Structure](#module-structure)
    - [Component Interaction](#component-interaction)
  - [Implementation](#implementation)
    - [File Structure](#file-structure)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
- [SyncPlay settings](#syncplay-settings)
- [Sync settings](#sync-settings)
- [Chat settings](#chat-settings)
- [Room settings](#room-settings)
- [WebSocket](#websocket)
    - [Config Keys](#config-keys)
  - [API Endpoints](#api-endpoints)
    - [Content Management](#content-management)
      - [GET /api/v1/syncplay/rooms](#get-apiv1syncplayrooms)
      - [POST /api/v1/syncplay/rooms](#post-apiv1syncplayrooms)
      - [GET /api/v1/syncplay/rooms/:id](#get-apiv1syncplayroomsid)
      - [PUT /api/v1/syncplay/rooms/:id](#put-apiv1syncplayroomsid)
      - [DELETE /api/v1/syncplay/rooms/:id](#delete-apiv1syncplayroomsid)
      - [POST /api/v1/syncplay/rooms/:id/join](#post-apiv1syncplayroomsidjoin)
      - [POST /api/v1/syncplay/rooms/:id/leave](#post-apiv1syncplayroomsidleave)
      - [GET /api/v1/syncplay/rooms/:id/participants](#get-apiv1syncplayroomsidparticipants)
      - [POST /api/v1/syncplay/rooms/:id/chat](#post-apiv1syncplayroomsidchat)
      - [GET /api/v1/syncplay/rooms/:id/chat](#get-apiv1syncplayroomsidchat)
      - [POST /api/v1/syncplay/rooms/:id/invitations](#post-apiv1syncplayroomsidinvitations)
      - [GET /api/v1/syncplay/invitations/:code](#get-apiv1syncplayinvitationscode)
      - [WS /api/v1/syncplay/ws](#ws-apiv1syncplayws)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# SyncPlay (Watch Together)


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: feature


> Content module for 

> Synchronized playback for multiple users watching together

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | - |
| Sources | ✅ | - |
| Instructions | ✅ | - |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete



---


## Architecture

```
  ┌─────────────┐     ┌──────────────┐     ┌─────────────┐
  │   Client    │────▶│  API Handler │────▶│   Service   │
  │  (Web/App)  │◀────│   (ogen)     │◀────│   (Logic)   │
  └─────────────┘     └──────────────┘     └──────┬──────┘
                                                   │
                            ┌──────────────────────┼────────────┐
                            ▼                      ▼            ▼
                      ┌──────────┐          ┌───────────┐  ┌────────┐
                      │Repository│          │ Metadata  │  │  Cache │
                      │  (sqlc)  │          │  Service  │  │(otter) │
                      └────┬─────┘          └─────┬─────┘  └────────┘
                           │                      │
                           ▼                      ▼
                    ┌─────────────┐        ┌──────────┐
                    │ PostgreSQL  │        │ External │
                    │   (pgx)     │        │   APIs   │
                    └─────────────┘        └──────────┘
  ```

### Database Schema

**Schema**: `public`

<!-- Schema diagram -->

### Module Structure

```
internal/content/syncplay_(watch_together)/
├── module.go              # fx module definition
├── repository.go          # Database operations
├── service.go             # Business logic
├── handler.go             # HTTP handlers (ogen)
├── types.go               # Domain types
└── syncplay_(watch_together)_test.go
```

### Component Interaction

<!-- Component interaction diagram -->


## Implementation

### File Structure

```
internal/playback/syncplay/
├── module.go                    # fx module registration
├── repository.go                # Database operations (sqlc)
├── queries.sql                  # SQL queries for sqlc
├── service.go                   # Business logic
├── handler.go                   # HTTP/WebSocket handlers
├── types.go                     # Domain types
├── room_manager.go              # Room lifecycle management
├── sync_engine.go               # Playback synchronization logic
├── websocket.go                 # WebSocket connection handling
├── chat.go                      # Chat functionality
├── events.go                    # Event logging and broadcasting
└── cache.go                     # Caching layer (otter)

cmd/server/
└── main.go                      # Server entry point with fx

migrations/
├── 033_syncplay.up.sql          # SyncPlay tables
└── 033_syncplay.down.sql        # Rollback

api/openapi/
└── syncplay.yaml                # OpenAPI spec (HTTP endpoints)

web/src/lib/components/syncplay/
├── SyncPlayRoom.svelte          # Room UI
├── ParticipantsList.svelte      # Participants sidebar
├── SyncPlayChat.svelte          # Chat panel
├── SyncPlayControls.svelte      # Playback controls
└── SyncPlayInvite.svelte        # Invite dialog
```


### Key Interfaces

```go
// Repository interface for syncplay database operations
type Repository interface {
    // Rooms
    CreateRoom(ctx context.Context, params CreateRoomParams) (*SyncPlayRoom, error)
    GetRoom(ctx context.Context, id uuid.UUID) (*SyncPlayRoom, error)
    ListRooms(ctx context.Context, filters RoomFilters, limit, offset int) ([]*SyncPlayRoom, int64, error)
    UpdateRoom(ctx context.Context, id uuid.UUID, params UpdateRoomParams) (*SyncPlayRoom, error)
    UpdatePlaybackState(ctx context.Context, roomID uuid.UUID, position float64, isPlaying bool, userID uuid.UUID) error
    EndRoom(ctx context.Context, id uuid.UUID) error
    DeleteRoom(ctx context.Context, id uuid.UUID) error

    // Participants
    AddParticipant(ctx context.Context, roomID, userID uuid.UUID) (*SyncPlayParticipant, error)
    GetParticipant(ctx context.Context, roomID, userID uuid.UUID) (*SyncPlayParticipant, error)
    ListParticipants(ctx context.Context, roomID uuid.UUID, connectedOnly bool) ([]*SyncPlayParticipant, error)
    UpdateParticipantState(ctx context.Context, roomID, userID uuid.UUID, params UpdateParticipantParams) error
    RemoveParticipant(ctx context.Context, roomID, userID uuid.UUID) error
    UpdateParticipantConnection(ctx context.Context, roomID, userID uuid.UUID, isConnected bool) error

    // Chat
    CreateChatMessage(ctx context.Context, roomID, userID uuid.UUID, message string) (*SyncPlayChatMessage, error)
    GetChatMessages(ctx context.Context, roomID uuid.UUID, limit, offset int) ([]*SyncPlayChatMessage, error)
    DeleteChatMessage(ctx context.Context, id uuid.UUID) error

    // Events
    LogEvent(ctx context.Context, params LogEventParams) error
    GetEvents(ctx context.Context, roomID uuid.UUID, limit, offset int) ([]*SyncPlayEvent, error)

    // Invitations
    CreateInvitation(ctx context.Context, params CreateInvitationParams) (*SyncPlayInvitation, error)
    GetInvitationByCode(ctx context.Context, code string) (*SyncPlayInvitation, error)
    DeleteInvitation(ctx context.Context, id uuid.UUID) error
    CleanupExpiredInvitations(ctx context.Context) (int64, error)
}

// Service interface for syncplay operations
type Service interface {
    // Room management
    CreateRoom(ctx context.Context, userID uuid.UUID, req CreateRoomRequest) (*SyncPlayRoom, error)
    GetRoom(ctx context.Context, roomID uuid.UUID) (*SyncPlayRoomDetail, error)
    ListRooms(ctx context.Context, filters RoomFilters) (*RoomListResponse, error)
    EndRoom(ctx context.Context, roomID, userID uuid.UUID) error
    UpdateRoomSettings(ctx context.Context, roomID, userID uuid.UUID, updates RoomSettingsUpdate) (*SyncPlayRoom, error)

    // Joining/leaving
    JoinRoom(ctx context.Context, userID uuid.UUID, req JoinRoomRequest) (*JoinRoomResponse, error)
    LeaveRoom(ctx context.Context, roomID, userID uuid.UUID) error

    // Invitations
    CreateInvitation(ctx context.Context, roomID, userID uuid.UUID, invitedUserID *uuid.UUID, expiresIn *time.Duration) (*SyncPlayInvitation, error)
    ValidateInvitation(ctx context.Context, code string) (*SyncPlayRoom, error)

    // Chat
    SendChatMessage(ctx context.Context, roomID, userID uuid.UUID, message string) error
    GetChatHistory(ctx context.Context, roomID uuid.UUID, pagination Pagination) (*ChatHistoryResponse, error)
}

// RoomManager interface for room lifecycle
type RoomManager interface {
    CreateRoom(roomID uuid.UUID, contentType string, contentID uuid.UUID, ownerID uuid.UUID) error
    GetRoom(roomID uuid.UUID) (*Room, error)
    CloseRoom(roomID uuid.UUID) error
    BroadcastToRoom(roomID uuid.UUID, event RoomEvent) error
}

// SyncEngine interface for playback synchronization
type SyncEngine interface {
    HandlePlaybackCommand(ctx context.Context, roomID, userID uuid.UUID, cmd PlaybackCommand) error
    SyncParticipant(ctx context.Context, roomID, userID uuid.UUID, clientState ClientState) (*SyncResponse, error)
    CalculateSyncOffset(roomState *RoomState, clientState *ClientState) (offset float64, shouldSync bool)
}

// WebSocketManager interface for WebSocket connections
type WebSocketManager interface {
    HandleConnection(conn *websocket.Conn, userID uuid.UUID, roomID uuid.UUID) error
    SendToUser(roomID, userID uuid.UUID, message interface{}) error
    SendToRoom(roomID uuid.UUID, message interface{}, excludeUserID *uuid.UUID) error
    DisconnectUser(roomID, userID uuid.UUID) error
}
```


### Dependencies
**Go Packages**:
```go
require (
    // Core
    github.com/google/uuid v1.6.0
    go.uber.org/fx v1.23.0

    // Database
    github.com/jackc/pgx/v5 v5.7.2
    github.com/sqlc-dev/sqlc v1.28.0

    // API
    github.com/ogen-go/ogen v1.7.0

    // WebSocket
    github.com/gorilla/websocket v1.5.3

    // Caching
    github.com/maypok86/otter v1.2.4

    // Password hashing
    golang.org/x/crypto v0.31.0

    // Testing
    github.com/stretchr/testify v1.10.0
    github.com/testcontainers/testcontainers-go v0.35.0
)
```

**External Dependencies**:
- **PostgreSQL 18+**: Database
- **WebSocket**: Real-time communication







## Configuration

### Environment Variables

```bash
# SyncPlay settings
SYNCPLAY_ENABLED=true
SYNCPLAY_DEFAULT_MAX_PARTICIPANTS=50
SYNCPLAY_MAX_ROOM_AGE_HOURS=24         # Auto-end inactive rooms

# Sync settings
SYNCPLAY_SYNC_THRESHOLD_SECONDS=2      # Sync if offset > 2 seconds
SYNCPLAY_HEARTBEAT_INTERVAL_SECONDS=5  # Client heartbeat interval
SYNCPLAY_DISCONNECT_TIMEOUT_SECONDS=30 # Mark disconnected after X seconds

# Chat settings
SYNCPLAY_CHAT_ENABLED=true
SYNCPLAY_CHAT_MAX_MESSAGE_LENGTH=1000
SYNCPLAY_CHAT_HISTORY_LIMIT=100

# Room settings
SYNCPLAY_ALLOW_PUBLIC_ROOMS=true
SYNCPLAY_ALLOW_PASSWORD_ROOMS=true
SYNCPLAY_INVITE_EXPIRY_HOURS=24

# WebSocket
SYNCPLAY_WS_READ_BUFFER_SIZE=4096
SYNCPLAY_WS_WRITE_BUFFER_SIZE=4096
SYNCPLAY_WS_PING_INTERVAL_SECONDS=30
SYNCPLAY_WS_PONG_TIMEOUT_SECONDS=60
```


### Config Keys
```yaml
syncplay:
  # Feature toggle
  enabled: true

  # Room settings
  rooms:
    default_max_participants: 50
    max_room_age_hours: 24
    allow_public: true
    allow_password_protected: true

  # Synchronization
  sync:
    threshold_seconds: 2          # Sync if offset > 2 sec
    heartbeat_interval_seconds: 5
    disconnect_timeout_seconds: 30
    max_latency_ms: 500           # Warn if latency > 500ms

  # Chat
  chat:
    enabled: true
    max_message_length: 1000
    history_limit: 100
    rate_limit_messages_per_minute: 60

  # Invitations
  invitations:
    default_expiry_hours: 24
    code_length: 8

  # WebSocket
  websocket:
    read_buffer_size: 4096
    write_buffer_size: 4096
    ping_interval_seconds: 30
    pong_timeout_seconds: 60
    max_connections_per_user: 5

  # Cleanup
  cleanup:
    inactive_rooms_hours: 24
    expired_invitations_hours: 48
```



## API Endpoints

### Content Management
#### GET /api/v1/syncplay/rooms

List SyncPlay rooms

---
#### POST /api/v1/syncplay/rooms

Create a new room

---
#### GET /api/v1/syncplay/rooms/:id

Get room details

---
#### PUT /api/v1/syncplay/rooms/:id

Update room settings

---
#### DELETE /api/v1/syncplay/rooms/:id

End/delete room

---
#### POST /api/v1/syncplay/rooms/:id/join

Join a room

---
#### POST /api/v1/syncplay/rooms/:id/leave

Leave a room

---
#### GET /api/v1/syncplay/rooms/:id/participants

List room participants

---
#### POST /api/v1/syncplay/rooms/:id/chat

Send chat message

---
#### GET /api/v1/syncplay/rooms/:id/chat

Get chat history

---
#### POST /api/v1/syncplay/rooms/:id/invitations

Create invitation

---
#### GET /api/v1/syncplay/invitations/:code

Get invitation details

---
#### WS /api/v1/syncplay/ws

WebSocket connection for real-time sync

---







## Related Documentation
### Design Documents
- [01_ARCHITECTURE](../../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../../architecture/02_DESIGN_PRINCIPLES.md)
- [03_METADATA_SYSTEM](../../architecture/03_METADATA_SYSTEM.md)

### External Sources
- [Go sync](../../../sources/go/stdlib/sync.md) - Auto-resolved from go-sync
- [Jellyfin SyncPlay](../../../sources/apis/jellyfin-syncplay.md) - Auto-resolved from jellyfin-syncplay


# Jellyfin SyncPlay API

> Source: https://jellyfin.org/docs/general/server/syncplay.html
> Type: stub
> Status: Placeholder - needs content fetch

---

## Overview

SyncPlay allows multiple users to watch content together in sync. The API manages group creation, membership, and playback synchronization.

## WebSocket Connection

SyncPlay uses WebSocket for real-time synchronization:
```
wss://{server}/socket?api_key={api_key}
```

## Core Endpoints

### Groups

| Endpoint | Description |
|----------|-------------|
| `POST /SyncPlay/New` | Create new group |
| `POST /SyncPlay/Join` | Join existing group |
| `POST /SyncPlay/Leave` | Leave group |
| `GET /SyncPlay/List` | List available groups |

### Playback Control

| Endpoint | Description |
|----------|-------------|
| `POST /SyncPlay/Play` | Play/resume |
| `POST /SyncPlay/Pause` | Pause playback |
| `POST /SyncPlay/Seek` | Seek to position |
| `POST /SyncPlay/SetNewQueue` | Set playback queue |

### State

| Endpoint | Description |
|----------|-------------|
| `POST /SyncPlay/Ping` | Send ping for sync |
| `POST /SyncPlay/Ready` | Mark as ready |
| `POST /SyncPlay/Buffering` | Report buffering |

## WebSocket Messages

```json
{
  "MessageType": "SyncPlayGroupUpdate",
  "Data": {
    "GroupId": "...",
    "PlaylistItemId": "...",
    "PositionTicks": 123456789,
    "State": "Playing"
  }
}
```

## Related

- [SyncPlay Feature](../../design/features/playback/SYNCPLAY.md)
- [WebSocket Protocol](../protocols/websocket.md)

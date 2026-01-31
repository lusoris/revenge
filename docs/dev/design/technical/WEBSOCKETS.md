# WebSockets

> Real-time bidirectional communication

**Source of Truth**: [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md)

---

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | 🟡 | Scaffold |
| Sources | 🔴 |  |
| Instructions | 🔴 |  |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |

---

## Overview

Revenge uses WebSockets for real-time features:
- Playback progress sync
- SyncPlay coordination
- Live notifications
- Remote control commands
- Admin dashboard updates

---

## Connection

```
wss://{server}/socket?token={jwt}
```

Authentication via JWT token in query parameter or header.

---

## Message Types

### Client → Server

| Type | Purpose |
|------|---------|
| `ping` | Keepalive |
| `playback.progress` | Report playback position |
| `syncplay.command` | SyncPlay actions |
| `remote.command` | Remote control |

### Server → Client

| Type | Purpose |
|------|---------|
| `pong` | Keepalive response |
| `notification` | Push notification |
| `syncplay.state` | SyncPlay sync |
| `library.update` | Library changes |
| `session.update` | Session changes |

---

## Message Format

```json
{
  "type": "playback.progress",
  "id": "msg-123",
  "data": {
    "itemId": "uuid",
    "positionTicks": 123456789,
    "isPaused": false
  }
}
```

---

## Implementation

Using `github.com/coder/websocket` for WebSocket handling:

```go
func (h *SocketHandler) HandleConnection(w http.ResponseWriter, r *http.Request) {
    conn, err := websocket.Accept(w, r, nil)
    if err != nil {
        return
    }
    defer conn.Close(websocket.StatusNormalClosure, "")

    // Handle messages
    for {
        _, data, err := conn.Read(ctx)
        if err != nil {
            break
        }
        h.processMessage(conn, data)
    }
}
```

---

## Related

- [SyncPlay Feature](../features/playback/SYNCPLAY.md)
- [Session Service](../services/SESSION.md)
- [WebSocket Source](../../sources/tooling/websocket.md)

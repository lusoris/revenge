---
applyTo: "**/internal/api/**/*.go,**/internal/service/playback/**/*.go"
---

# WebSocket - coder/websocket

> Modern WebSocket library for Watch Party, live updates, quality switching

## Overview

Use `coder/websocket` (formerly nhooyr.io/websocket) for all WebSocket communication. It's the modern replacement for gorilla/websocket.

**Package**: `github.com/coder/websocket`

## Features

- Zero external dependencies
- Full RFC 6455 compliance
- Context-based cancellation
- Graceful close handshake
- Compression support (permessage-deflate)
- Concurrent-safe writes

## Installation

```bash
go get github.com/coder/websocket
```

## Server-Side Usage

### Accept Connection

```go
import "github.com/coder/websocket"

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
        Subprotocols:         []string{"revenge.v1"},
        InsecureSkipVerify:   false,
        OriginPatterns:       []string{"*.revenge.local", "localhost:*"},
        CompressionMode:      websocket.CompressionContextTakeover,
        CompressionThreshold: 256, // Compress messages > 256 bytes
    })
    if err != nil {
        log.Error("websocket accept failed", "error", err)
        return
    }
    defer conn.CloseNow()

    // Handle connection
    ctx := r.Context()
    if err := handleConnection(ctx, conn); err != nil {
        log.Error("websocket error", "error", err)
    }
}
```

### Read/Write Messages

```go
func handleConnection(ctx context.Context, conn *websocket.Conn) error {
    for {
        // Read message
        msgType, data, err := conn.Read(ctx)
        if err != nil {
            if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
                return nil // Normal close
            }
            return err
        }

        // Process message
        if msgType == websocket.MessageText {
            var msg Message
            if err := json.Unmarshal(data, &msg); err != nil {
                continue
            }

            // Handle message...
            response := processMessage(msg)

            // Send response
            responseData, _ := json.Marshal(response)
            if err := conn.Write(ctx, websocket.MessageText, responseData); err != nil {
                return err
            }
        }
    }
}
```

### Graceful Close

```go
// Close with status code and reason
err := conn.Close(websocket.StatusNormalClosure, "session ended")

// Close immediately (on error)
conn.CloseNow()
```

## Watch Party Implementation

```go
package watchparty

import (
    "context"
    "encoding/json"
    "sync"
    "github.com/coder/websocket"
    "github.com/google/uuid"
)

type Party struct {
    ID      uuid.UUID
    Host    *Participant
    Members map[uuid.UUID]*Participant
    mu      sync.RWMutex
}

type Participant struct {
    ID       uuid.UUID
    UserID   uuid.UUID
    Conn     *websocket.Conn
    SendChan chan []byte
}

type Message struct {
    Type    string          `json:"type"`
    Payload json.RawMessage `json:"payload"`
}

type SyncPayload struct {
    Position float64 `json:"position"`
    Playing  bool    `json:"playing"`
    Speed    float64 `json:"speed"`
}

func (p *Party) HandleParticipant(ctx context.Context, participant *Participant) error {
    defer p.removeParticipant(participant.ID)

    // Start writer goroutine
    go p.writeLoop(ctx, participant)

    // Read loop
    for {
        _, data, err := participant.Conn.Read(ctx)
        if err != nil {
            return err
        }

        var msg Message
        if err := json.Unmarshal(data, &msg); err != nil {
            continue
        }

        switch msg.Type {
        case "sync":
            // Only host can sync
            if participant.ID == p.Host.ID {
                p.broadcastExcept(data, participant.ID)
            }
        case "chat":
            p.broadcastAll(data)
        case "reaction":
            p.broadcastAll(data)
        }
    }
}

func (p *Party) writeLoop(ctx context.Context, participant *Participant) {
    for {
        select {
        case <-ctx.Done():
            return
        case msg := <-participant.SendChan:
            ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
            err := participant.Conn.Write(ctx, websocket.MessageText, msg)
            cancel()
            if err != nil {
                return
            }
        }
    }
}

func (p *Party) broadcastAll(data []byte) {
    p.mu.RLock()
    defer p.mu.RUnlock()

    for _, member := range p.Members {
        select {
        case member.SendChan <- data:
        default:
            // Channel full, skip
        }
    }
}

func (p *Party) broadcastExcept(data []byte, excludeID uuid.UUID) {
    p.mu.RLock()
    defer p.mu.RUnlock()

    for id, member := range p.Members {
        if id == excludeID {
            continue
        }
        select {
        case member.SendChan <- data:
        default:
        }
    }
}
```

## Playback Progress Updates

```go
package playback

type ProgressUpdate struct {
    SessionID uuid.UUID `json:"session_id"`
    MediaID   uuid.UUID `json:"media_id"`
    Position  float64   `json:"position"`   // seconds
    Duration  float64   `json:"duration"`   // seconds
    Playing   bool      `json:"playing"`
}

func (h *Handler) HandleProgressWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := websocket.Accept(w, r, nil)
    if err != nil {
        return
    }
    defer conn.CloseNow()

    ctx := r.Context()
    sessionID := r.PathValue("sessionID")

    // Send progress updates every 5 seconds
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            progress, err := h.svc.GetProgress(ctx, sessionID)
            if err != nil {
                continue
            }

            data, _ := json.Marshal(progress)
            if err := conn.Write(ctx, websocket.MessageText, data); err != nil {
                return
            }
        }
    }
}
```

## Quality Switching

```go
type QualitySwitch struct {
    Profile    string `json:"profile"`    // "1080p", "720p", etc.
    Bitrate    int    `json:"bitrate"`    // kbps
    FromServer bool   `json:"from_server"` // true if server-initiated
}

func (h *Handler) HandleStreamWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, err := websocket.Accept(w, r, nil)
    if err != nil {
        return
    }
    defer conn.CloseNow()

    ctx := r.Context()
    sessionID := r.PathValue("sessionID")

    for {
        _, data, err := conn.Read(ctx)
        if err != nil {
            return
        }

        var msg Message
        json.Unmarshal(data, &msg)

        switch msg.Type {
        case "quality_request":
            var req QualitySwitch
            json.Unmarshal(msg.Payload, &req)

            // Update transcode profile
            if err := h.svc.SetQuality(ctx, sessionID, req.Profile); err != nil {
                // Send error response
                continue
            }

            // Confirm quality change
            response, _ := json.Marshal(Message{
                Type:    "quality_changed",
                Payload: msg.Payload,
            })
            conn.Write(ctx, websocket.MessageText, response)

        case "bandwidth_report":
            var report BandwidthReport
            json.Unmarshal(msg.Payload, &report)

            // Server may initiate quality switch
            if newProfile := h.svc.SuggestQuality(report.Kbps); newProfile != "" {
                suggestion, _ := json.Marshal(Message{
                    Type: "quality_suggestion",
                    Payload: json.RawMessage(fmt.Sprintf(`{"profile":"%s","from_server":true}`, newProfile)),
                })
                conn.Write(ctx, websocket.MessageText, suggestion)
            }
        }
    }
}
```

## Ping/Pong Heartbeat

```go
func handleWithHeartbeat(ctx context.Context, conn *websocket.Conn) error {
    // Set read deadline
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    // Start ping goroutine
    go func() {
        ticker := time.NewTicker(30 * time.Second)
        defer ticker.Stop()

        for {
            select {
            case <-ctx.Done():
                return
            case <-ticker.C:
                ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
                err := conn.Ping(ctx)
                cancel()
                if err != nil {
                    return
                }
            }
        }
    }()

    // Normal message handling...
    return handleMessages(ctx, conn)
}
```

## DO's and DON'Ts

### DO

- ✅ Use context for cancellation and timeouts
- ✅ Use `CloseNow()` in defer for cleanup
- ✅ Use separate goroutine for writes
- ✅ Buffer send channels to prevent blocking
- ✅ Use compression for large messages
- ✅ Implement heartbeat (ping/pong)

### DON'T

- ❌ Use gorilla/websocket - use coder/websocket
- ❌ Write from multiple goroutines without coordination
- ❌ Forget to close connections on error
- ❌ Block on full send channels
- ❌ Use unbounded message queues
---

## Related

- [INDEX.instructions.md](INDEX.instructions.md) - Main instruction index with all cross-references
- [ARCHITECTURE_V2.md](../../docs/dev/design/architecture/ARCHITECTURE_V2.md) - System architecture
- [BEST_PRACTICES.md](../../docs/dev/design/operations/BEST_PRACTICES.md) - Best practices

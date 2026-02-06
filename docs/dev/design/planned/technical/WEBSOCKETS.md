## Table of Contents

- [WebSockets](#websockets)
  - [Status](#status)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)

# WebSockets

<!-- DESIGN: technical, README, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES -->


**Created**: 2026-01-31
**Status**: ✅ Complete
**Category**: technical


> > Real-time bidirectional communication via WebSockets

Complete WebSocket infrastructure:
- **Library**: nhooyr.io/websocket (coder/websocket)
- **Channels**: Playback sync, notifications, live updates, transcoding progress
- **Authentication**: JWT token validation on connect
- **Protocol**: JSON message format with type-based routing
- **Features**: Automatic reconnection, heartbeat/ping, message acknowledgment

---


## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Complete WebSocket system design |
| Sources | ✅ | coder/websocket documentation included |
| Instructions | ✅ | Generated from design |
| Code | 🔴 | - |
| Linting | 🔴 | - |
| Unit Testing | 🔴 | - |
| Integration Testing | 🔴 | - |

**Overall**: ✅ Complete


## Related Documentation
### Design Documents
- [technical](INDEX.md)
- [01_ARCHITECTURE](../architecture/01_ARCHITECTURE.md)
- [02_DESIGN_PRINCIPLES](../architecture/02_DESIGN_PRINCIPLES.md)
- [WEBSOCKETS (Syncplay)](../features/playback/SYNCPLAY.md)

### External Sources
- [coder/websocket GitHub README](../../sources/tooling/websocket-guide.md) - Auto-resolved from coder-websocket-docs
- [nhooyr.io/websocket](https://pkg.go.dev/nhooyr.io/websocket) - WebSocket library (coder/websocket)
- [Uber fx](../../sources/tooling/fx.md) - Auto-resolved from fx


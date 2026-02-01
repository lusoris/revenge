

---
sources:
  - name: coder/websocket GitHub README
    url: ../../sources/tooling/websocket-guide.md
    note: Auto-resolved from coder-websocket-docs
  - name: nhooyr.io/websocket
    url: https://pkg.go.dev/nhooyr.io/websocket
    note: WebSocket library (coder/websocket)
  - name: Uber fx
    url: ../../sources/tooling/fx.md
    note: Auto-resolved from fx
design_refs:
  - title: technical
    path: INDEX.md
  - title: 01_ARCHITECTURE
    path: ../architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: ../architecture/02_DESIGN_PRINCIPLES.md
  - title: WEBSOCKETS (Syncplay)
    path: ../features/playback/SYNCPLAY.md
---

## Table of Contents

- [WebSockets](#websockets)
  - [Status](#status)
  - [Architecture](#architecture)
    - [Components](#components)
  - [Implementation](#implementation)
    - [File Structure](#file-structure)
    - [Key Interfaces](#key-interfaces)
    - [Dependencies](#dependencies)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
    - [Config Keys](#config-keys)
  - [Testing Strategy](#testing-strategy)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [Test Coverage](#test-coverage)
  - [Related Documentation](#related-documentation)
    - [Design Documents](#design-documents)
    - [External Sources](#external-sources)


# WebSockets


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



---


## Architecture

<!-- Architecture diagram placeholder -->

### Components

<!-- Component description -->


## Implementation

### File Structure

<!-- File structure -->

### Key Interfaces

<!-- Interface definitions -->

### Dependencies

<!-- Dependency list -->





## Configuration
### Environment Variables

<!-- Environment variables -->

### Config Keys

<!-- Configuration keys -->




## Testing Strategy

### Unit Tests

<!-- Unit test strategy -->

### Integration Tests

<!-- Integration test strategy -->

### Test Coverage

Target: **80% minimum**







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


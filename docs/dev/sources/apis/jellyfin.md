# Jellyfin API

> Source: https://jellyfin.org/docs/general/server/api.html
> Type: stub
> Status: Placeholder - needs content fetch

---

## Overview

Jellyfin is a free, open-source media server. The API provides access to library management, playback, user management, and system administration.

## API Base URL

```
https://{server}/api
```

## Authentication

API key or user authentication via headers:
```
X-Emby-Token: {api_key}
```

Or Authorization header with bearer token after login.

## Core Endpoints

### Items

| Endpoint | Description |
|----------|-------------|
| `GET /Items` | Get library items |
| `GET /Items/{id}` | Get item details |
| `GET /Items/{id}/Images` | Get item images |
| `POST /Items/{id}/Refresh` | Refresh item metadata |

### Users

| Endpoint | Description |
|----------|-------------|
| `GET /Users` | Get users |
| `GET /Users/{id}` | Get user details |
| `POST /Users/AuthenticateByName` | Login |

### Playback

| Endpoint | Description |
|----------|-------------|
| `GET /Videos/{id}/stream` | Stream video |
| `POST /Sessions/Playing` | Report playback start |
| `POST /Sessions/Playing/Progress` | Report progress |
| `POST /Sessions/Playing/Stopped` | Report playback stop |

### Library

| Endpoint | Description |
|----------|-------------|
| `GET /Library/VirtualFolders` | Get libraries |
| `POST /Library/Refresh` | Refresh all libraries |

## Related

- [Player Architecture](../../design/architecture/04_PLAYER_ARCHITECTURE.md)
- [Client Support](../../design/features/shared/CLIENT_SUPPORT.md)

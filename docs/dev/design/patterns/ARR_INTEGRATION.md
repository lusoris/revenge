# Arr Integration Pattern

> Common patterns for integrating with *arr services (Radarr, Sonarr, Lidarr, Whisparr, etc.)

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

All *arr services (Radarr, Sonarr, Lidarr, Whisparr, Chaptarr) share common API patterns. This document defines reusable patterns for:

- API client configuration
- Authentication handling
- Webhook processing
- Event synchronization
- Error handling

---

## Common API Pattern

### Client Structure

```go
type ArrClient struct {
    baseURL    string
    apiKey     string
    httpClient *resty.Client
}

func NewArrClient(baseURL, apiKey string) *ArrClient {
    client := resty.New().
        SetBaseURL(baseURL).
        SetHeader("X-Api-Key", apiKey).
        SetTimeout(30 * time.Second)

    return &ArrClient{
        baseURL:    baseURL,
        apiKey:     apiKey,
        httpClient: client,
    }
}
```

### Standard Endpoints

| Endpoint | Purpose |
|----------|---------|
| `GET /api/v3/system/status` | Health check |
| `GET /api/v3/{resource}` | List items |
| `GET /api/v3/{resource}/{id}` | Get single item |
| `POST /api/v3/{resource}` | Create item |
| `PUT /api/v3/{resource}/{id}` | Update item |
| `DELETE /api/v3/{resource}/{id}` | Delete item |

---

## Webhook Pattern

### Event Types

| Event | Trigger |
|-------|---------|
| `Grab` | Item grabbed from indexer |
| `Download` | Download completed |
| `Rename` | File renamed |
| `Delete` | Item deleted |
| `Health` | Health status changed |

### Handler Pattern

```go
func (h *WebhookHandler) HandleArrWebhook(w http.ResponseWriter, r *http.Request) {
    var event ArrWebhookEvent
    if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
        http.Error(w, "Invalid payload", http.StatusBadRequest)
        return
    }

    switch event.EventType {
    case "Download":
        h.handleDownload(event)
    case "Grab":
        h.handleGrab(event)
    // ...
    }

    w.WriteHeader(http.StatusOK)
}
```

---

## Related

- [Radarr Integration](../integrations/servarr/RADARR.md)
- [Sonarr Integration](../integrations/servarr/SONARR.md)
- [Lidarr Integration](../integrations/servarr/LIDARR.md)
- [Whisparr Integration](../integrations/servarr/WHISPARR.md)
- [Chaptarr Integration](../integrations/servarr/CHAPTARR.md)
- [Webhook Patterns](WEBHOOK_PATTERNS.md)

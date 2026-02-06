# Webhook Patterns

> Incoming webhook handling: receive events from external services, convert types, queue async jobs. Written from code as of 2026-02-06.

---

## Current State

Webhooks are **incoming only** — Radarr and Sonarr send events to Revenge. No outgoing webhooks exist yet.

| Source | Endpoint | Handler |
|--------|----------|---------|
| Radarr | `POST /api/v1/webhooks/radarr` | `HandleRadarrWebhook` |
| Sonarr | `POST /api/v1/webhooks/sonarr` | `HandleSonarrWebhook` |

---

## Incoming Webhook Flow

```
Radarr/Sonarr fires event
    ↓
POST /api/v1/webhooks/{arr}
    ↓
API Handler receives ogen-generated request type
    ↓
convertWebhookPayload() -> internal types
    ↓
River client available?
    ├── Yes: Insert job (QueueHigh, 1min timeout) -> 202 Accepted
    └── No:  Log warning -> 202 Accepted (event dropped)
    ↓
Worker picks up job
    ↓
WebhookHandler.HandleWebhook() routes by EventType
    ↓
Action (sync, log, or no-op)
```

### Type Conversion

Webhook payloads arrive as ogen-generated types (from OpenAPI spec) and must be converted to internal types before processing:

```go
func convertWebhookPayload(req *ogen.RadarrWebhookPayload) *radarr.WebhookPayload {
    payload := &radarr.WebhookPayload{
        EventType:      string(req.EventType),
        InstanceName:   req.InstanceName.Value,
        DownloadClient: req.DownloadClient.Value,
        IsUpgrade:      req.IsUpgrade.Value,
    }

    // Optional nested objects use .Set checks
    if movie := req.Movie; movie.Set {
        payload.Movie = &radarr.WebhookMovie{
            ID:     int(movie.Value.ID.Value),
            Title:  movie.Value.Title.Value,
            TMDbID: int(movie.Value.TmdbId.Value),
        }
    }
    return payload
}
```

Ogen optional fields use `req.Field.Set` (bool) and `req.Field.Value` (actual data).

### Job Queuing

```go
if h.riverClient != nil {
    _, err := h.riverClient.Insert(ctx, &radarr.RadarrWebhookJobArgs{
        Payload: *payload,
    }, nil)
    if err != nil {
        return &ogen.Error{Code: 400, Message: "Failed to process webhook"}, nil
    }
    return &ogen.HandleRadarrWebhookAccepted{}, nil
}

// Fallback: no River client
h.logger.Warn("Webhook received but no handler configured")
return &ogen.HandleRadarrWebhookAccepted{}, nil
```

Always returns 202 immediately. Processing happens asynchronously via River workers.

### Event Routing

The `WebhookHandler` in each integration package dispatches by event type:

| Event | Action | Details |
|-------|--------|---------|
| Download | **Sync** | Main trigger — sync content + files from arr |
| Rename | **Sync** | Update file paths |
| Delete | **Log** | Don't auto-delete from local DB |
| File Delete | **Sync** | Reflect removed files |
| Grab | **No-op** | Wait for download to complete |
| Health | **Log** | Warning-level logging |
| Test | **Ack** | Confirms webhook connectivity |

---

## Authentication

**Current**: No webhook-specific authentication. Webhook endpoints are public — assumed to be behind a firewall or reverse proxy.

**Planned**: Build a proper webhook service with:
- HMAC signature validation
- Per-source API key authentication
- Request deduplication (event ID tracking)
- IP whitelisting

---

## Payload Differences

### Radarr Webhook

```go
type WebhookPayload struct {
    EventType      string
    InstanceName   string
    Movie          *WebhookMovie     // TMDbID, IMDbID, Title, Year
    MovieFile      *WebhookMovieFile // Path, Quality, Size
    Release        *WebhookRelease   // ReleaseGroup, Indexer, Size
    DownloadClient string
    DownloadID     string
    IsUpgrade      bool
}
```

### Sonarr Webhook

```go
type WebhookPayload struct {
    EventType          string
    InstanceName       string
    Series             *WebhookSeries       // TVDbID, TVMazeID, IMDbID
    Episodes           []WebhookEpisode     // Array (season packs)
    EpisodeFile        *WebhookEpisodeFile  // Path, Quality, Size, DateAdded
    DeletedFiles       []WebhookEpisodeFile // For cleanup events
    Release            *WebhookRelease      // Same as Radarr
    DownloadClient     string
    DownloadClientType string               // Radarr doesn't have this
    DownloadID         string
    IsUpgrade          bool
}
```

Key difference: Sonarr has **episodes array** (multiple episodes per download for season packs) and **deleted files array** (for cleanup events). Radarr is single-movie.

---

## Planned: Outgoing Webhooks

Not yet implemented. The vision is a standalone webhook service that:

1. Listens for internal domain events (movie added, sync complete, health change)
2. Dispatches to user-configured webhook URLs
3. Supports retry with exponential backoff
4. Provides webhook delivery logs and status

This would live in `internal/service/webhooks/` as a new service module.

---

## Adding a New Webhook Source

1. Define ogen types in `api/openapi/` for the webhook payload
2. Add handler method in `internal/api/handler_{source}.go`
3. Implement `convert{Source}WebhookPayload()` for type conversion
4. Create job args type in the integration package
5. Create worker in the integration package
6. Create `WebhookHandler` with event routing switch

See [Arr Integration Pattern](SERVARR.md) for the full 8-layer template.

---

## Related Documentation

- [Arr Integration Pattern](SERVARR.md) — Radarr/Sonarr webhook handlers in context
- [River Workers](RIVER_WORKERS.md) — Async job processing for webhooks
- [Error Handling](ERROR_HANDLING.md) — How webhook errors propagate

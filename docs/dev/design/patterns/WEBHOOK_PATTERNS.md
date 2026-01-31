# Webhook Patterns

<!-- SOURCES: river -->

<!-- DESIGN: patterns, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> Patterns for processing webhooks and handling events

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

Revenge receives webhooks from external services (*arr stack, scrobbling services, etc.). This document defines patterns for:

- Webhook registration
- Payload validation
- Event processing
- Error handling
- Retry logic

---

## Webhook Handler Pattern

### Router Setup

```go
func SetupWebhookRoutes(r chi.Router, h *WebhookHandler) {
    r.Route("/webhooks", func(r chi.Router) {
        r.Post("/radarr", h.HandleRadarr)
        r.Post("/sonarr", h.HandleSonarr)
        r.Post("/lidarr", h.HandleLidarr)
        r.Post("/whisparr", h.HandleWhisparr)
        r.Post("/trakt", h.HandleTrakt)
    })
}
```

### Validation Middleware

```go
func ValidateWebhookSignature(secret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            signature := r.Header.Get("X-Webhook-Signature")
            if !verifySignature(r.Body, signature, secret) {
                http.Error(w, "Invalid signature", http.StatusUnauthorized)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

---

## Event Processing

### Async Processing with River

```go
type WebhookEventArgs struct {
    Source    string          `json:"source"`
    EventType string          `json:"event_type"`
    Payload   json.RawMessage `json:"payload"`
}

func (w *WebhookWorker) Work(ctx context.Context, job *river.Job[WebhookEventArgs]) error {
    switch job.Args.Source {
    case "radarr":
        return w.processRadarrEvent(ctx, job.Args)
    case "sonarr":
        return w.processSonarrEvent(ctx, job.Args)
    }
    return nil
}
```

---

## Idempotency

### Deduplication Pattern

```go
func (h *WebhookHandler) isDuplicate(eventID string) bool {
    key := fmt.Sprintf("webhook:event:%s", eventID)
    exists, _ := h.cache.Exists(ctx, key)
    if exists {
        return true
    }
    h.cache.Set(ctx, key, "1", 24*time.Hour)
    return false
}
```

---

## Related

- [Arr Integration Pattern](ARR_INTEGRATION.md)
- [Servarr Integrations](../integrations/servarr/)
- [Scrobbling](../features/shared/SCROBBLING.md)

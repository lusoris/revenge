# Phase A4: Webhook & Notification Enhancements

**Priority**: P2
**Effort**: 2-3h
**Dependencies**: A0

---

## A4.1: Custom Webhook Templates

**Affected File**: `internal/service/notification/agents/webhook.go:190`

**Current State**:
```go
// TODO: Support custom PayloadTemplate with Go templates
```

**Tasks**:
- [ ] Parse `PayloadTemplate` as Go template
- [ ] Provide template data: event type, payload, timestamp, etc.
- [ ] Validate template on agent creation
- [ ] Tests

**Template Data Structure**:
```go
type WebhookTemplateData struct {
    EventType   string
    Timestamp   time.Time
    Payload     any
    User        *UserInfo  // optional
    Movie       *MovieInfo // optional, if movie-related event
}
```

**Example Template**:
```json
{
  "event": "{{.EventType}}",
  "time": "{{.Timestamp.Format \"2006-01-02T15:04:05Z07:00\"}}",
  "data": {{.Payload | json}}
}
```

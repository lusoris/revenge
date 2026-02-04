# Phase A4: Webhook & Notification Enhancements

**Priority**: P2
**Effort**: 2-3h
**Dependencies**: A0
**Status**: ✅ Complete (2026-02-04)

---

## A4.1: Custom Webhook Templates ✅

**Affected File**: `internal/service/notification/agents/webhook.go`

**Completed Tasks**:
- [x] Parse `PayloadTemplate` as Go template at agent creation
- [x] Create `WebhookTemplateData` struct with all event fields
- [x] Validate template on agent creation (returns error if invalid)
- [x] Add `json` and `jsonIndent` template functions
- [x] Execute template in `buildTemplatedPayload` method

**Implementation Details**:

Template data structure:
```go
type WebhookTemplateData struct {
    EventID   string
    EventType string
    Timestamp time.Time
    UserID    string
    TargetID  string
    Data      map[string]any
    Metadata  map[string]string
    Source    string
}
```

Template functions available:
- `json` - Encodes value as compact JSON
- `jsonIndent` - Encodes value as indented JSON

Example template:
```json
{
  "event": "{{.EventType}}",
  "time": "{{.Timestamp.Format "2006-01-02T15:04:05Z07:00"}}",
  "data": {{.Data | json}}
}
```

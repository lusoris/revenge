# Notification System Split - Design Notes

**Date**: 2026-02-06
**Context**: During Step 6 doc alignment, decided to split notification into two distinct systems.
**Decisions**: See [QUESTIONS_INBOX_HELPDESK.md](QUESTIONS_INBOX_HELPDESK.md) for all 8 decisions.
**Planned doc**: `docs/dev/design/planned/services/ANNOUNCEMENTS.md` (was INBOX.md, renamed per Q1 decision)

---

## Decision

Split what was a single "notification" concept into **four independent systems**:

### 1. External Notification Service (EXISTS - `internal/service/notification/`)

Current event dispatcher with pluggable agents. Outbound-only, admin-configured.

- **What it does**: Fire-and-forget notifications to external channels on system events
- **Agents**: webhook, discord, email, gotify, ntfy
- **27 event types** across 7 categories (content, requests, library, user, auth, playback, system)
- **Doc**: `docs/dev/design/services/NOTIFICATION.md` (aligned with code)

### 2. In-App Notification + Helpdesk System (PLANNED - new module)

User-facing notification center inside the app. Much broader scope than the event dispatcher.

**Notification sources** (needs adapters to each):
- Wiki/helpdesk integration - contextual help, guided assistance, "how do I..."
- App news - release notes, changelog highlights
- Admin announcements - system-wide messages from admins
- Mod messages - moderation notices, content warnings
- Content module news - release calendars, new episodes, movie premieres (adapters to movie, tvshow, etc.)
- System alerts - storage warnings, scan failures, update available

**Helpdesk component**:
- Connected to wiki system
- Contextual help (page-aware suggestions)
- User self-service troubleshooting
- FAQ / knowledge base integration

**Dependencies this will need**:
- Adapters to every content module (movie, tvshow, future modules)
- User settings (notification preferences per category, per channel)
- Admin settings (which notification types are enabled, default preferences)
- RBAC integration (admin messages vs user messages vs mod messages)
- Database layer (notification storage, read/unread state, history)
- Possibly WebSocket/SSE for real-time delivery

**User settings expansion needed**:
- Per-category notification preferences (content news, admin messages, helpdesk, etc.)
- Delivery preferences (in-app only, email digest, push, etc.)
- Quiet hours / DND
- Per-library notification settings (only notify for libraries I care about)

**RBAC considerations**:
- `notifications:read` - view own notifications
- `notifications:manage` - manage notification preferences
- `notifications:admin` - send admin announcements
- `notifications:mod` - send mod messages
- `helpdesk:read` - access helpdesk/wiki
- `helpdesk:manage` - manage helpdesk content (admin)

---

## Naming

- External dispatcher: keep as `notification` service (or rename to `alerting`/`dispatch`?)
- In-app system: new package, maybe `internal/service/inbox` or `internal/service/announcements`?
- Helpdesk: could be part of in-app system or separate `internal/service/helpdesk`

## Open Questions

- Should the external dispatcher also feed into the in-app notification center? (e.g., "library scan completed" shows both as Discord message AND in-app notification)
- How tightly coupled is helpdesk to the wiki? Separate service or wiki adapter?
- Push notifications (mobile/browser) - part of in-app system or separate?

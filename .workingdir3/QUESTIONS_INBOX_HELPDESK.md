# Open Questions: In-App Announcements + Helpdesk System

**Date**: 2026-02-06
**Context**: Splitting notification into two systems. Decisions captured below.

---

## Decisions (2026-02-06)

| # | Question | Decision |
|---|----------|----------|
| Q1 | Package naming | `internal/service/announcements/` |
| Q2 | Dispatcher → inbox bridge | **Keep independent** - no bridge, content modules create entries directly |
| Q3 | Helpdesk scope | **Full wiki integration** - helpdesk IS the wiki, bidirectional sync, user-editable |
| Q4 | Real-time delivery | **WebSocket for all** - single WS for notifications + future features (playback, chat) |
| Q5 | Notification preferences | **Extend settings service** - add notification prefs as UserSettings |
| Q6 | Content module adapter | **Event-based** - modules emit events, announcements service subscribes |
| Q7 | Push notifications | **Separate service** - own complexity (FCM/APNs, device management) |
| Q8 | Moderation messages | **Direct messages** - mods send directly, RBAC controls who can send |

---

## Architecture Summary (from decisions)

```
Three independent systems:

1. External Dispatcher (EXISTS)
   Package: internal/service/notification/
   Purpose: Outbound alerts to external channels (Discord, webhook, gotify, ntfy, email)
   No database, no user preferences

2. Announcements (PLANNED)
   Package: internal/service/announcements/
   Purpose: In-app user-facing notifications, admin/mod messages, content news
   Database: shared.announcements, shared.user_announcements
   Preferences: Via existing settings service (UserSettings)
   Delivery: WebSocket (real-time) + email digest (scheduled)
   Content adapter: Event-based (subscribes to content module events)
   Mod messages: Direct to users, RBAC-controlled

3. Helpdesk / Wiki (PLANNED)
   Scope: Full wiki integration - content IS the wiki
   Bidirectional sync, user-editable
   Contextual help per page
   Separate from announcements but can create announcement entries

4. Push Notifications (PLANNED, separate)
   Package: internal/service/push/ (future)
   Purpose: FCM/APNs, browser push, device token management
   Separate service due to complexity
```

---

## Original Questions (for reference)

### Q1: Naming / Package Structure

The external dispatcher is `internal/service/notification/`. What should the in-app system be called?

- Option A: `internal/service/inbox/`
- **Option B: `internal/service/announcements/`** ✅
- Option C: Two packages: `inbox/` + `helpdesk/`
- Option D: Single `inbox/` with helpdesk inside

### Q2: External Dispatcher → Inbox Bridge

Should the existing event dispatcher also create in-app notifications?

- Option A: Yes, bridge them
- **Option B: Keep independent** ✅
- Option C: Hybrid

### Q3: Helpdesk Scope

How deep should the helpdesk/wiki integration go?

- Option A: Simple FAQ/knowledge base
- Option B: Interactive helpdesk
- **Option C: Full wiki integration** ✅
- Option D: Start simple, extend later

### Q4: Real-time Delivery

How should in-app notifications reach the user in real-time?

- **Option A/D: WebSocket for all** ✅
- Option B: SSE
- Option C: Polling

### Q5: User Settings Expansion

Where should notification preferences live?

- **Option A: Extend existing settings service** ✅
- Option B: Own table in announcements service
- Option C: New user_preferences service

### Q6: Content Module Adapter Pattern

How should content modules feed notifications?

- Option A: Direct dependency
- **Option B: Event-based** ✅
- Option C: Interface/adapter
- Option D: River jobs

### Q7: Push Notifications

Part of this system or separate?

- Option A: Part of inbox/announcements
- **Option B: Separate service** ✅
- Option C: Defer

### Q8: Moderation Messages

How should mod messages work?

- Option A: Same as admin announcements
- Option B: With approval workflow
- **Option C: Direct messages, RBAC-controlled** ✅

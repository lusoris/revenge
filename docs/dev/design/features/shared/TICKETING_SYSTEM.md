# Ticketing System

> User portal for bug reports, feature requests, and support tickets


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Executive Summary](#executive-summary)
- [User Portal](#user-portal)
  - [Ticket Submission Form](#ticket-submission-form)
  - [Ticket Viewing](#ticket-viewing)
  - [User Actions](#user-actions)
- [Admin Interface](#admin-interface)
  - [Triage Dashboard](#triage-dashboard)
  - [Admin Actions](#admin-actions)
- [PostgreSQL Schema](#postgresql-schema)
- [API Endpoints](#api-endpoints)
  - [User Endpoints](#user-endpoints)
  - [Admin Endpoints](#admin-endpoints)
- [Email Notifications](#email-notifications)
  - [User Notifications](#user-notifications)
  - [Admin Notifications](#admin-notifications)
- [Webhook Triggers](#webhook-triggers)
  - [Events](#events)
  - [Payload Example](#payload-example)
  - [Integrations](#integrations)
- [GitHub Issues Sync (Optional)](#github-issues-sync-optional)
  - [Two-Way Sync](#two-way-sync)
  - [GitHub Issue Metadata](#github-issue-metadata)
  - [Sync Rules](#sync-rules)
- [UI/UX Design](#uiux-design)
  - [User Ticket Submission Form](#user-ticket-submission-form)
  - [Admin Dashboard](#admin-dashboard)
  - [Ticket Detail Page](#ticket-detail-page)
- [Implementation Phases](#implementation-phases)
  - [Phase 1: Backend (Week 1)](#phase-1-backend-week-1)
  - [Phase 2: User Portal (Week 2)](#phase-2-user-portal-week-2)
  - [Phase 3: Admin Interface (Week 3)](#phase-3-admin-interface-week-3)
  - [Phase 4: Notifications (Week 4)](#phase-4-notifications-week-4)
  - [Phase 5: GitHub Sync (Optional, Week 5)](#phase-5-github-sync-optional-week-5)
- [Implementation Checklist](#implementation-checklist)
  - [Phase 1: Core Infrastructure](#phase-1-core-infrastructure)
  - [Phase 2: Database](#phase-2-database)
  - [Phase 3: Service Layer](#phase-3-service-layer)
  - [Phase 4: Background Jobs](#phase-4-background-jobs)
  - [Phase 5: API Integration](#phase-5-api-integration)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related](#related)

<!-- TOC-END -->

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Full design with DB schema, API endpoints, GitHub sync |
| Sources | 🟡 |  |
| Instructions | ✅ | Implementation checklist complete |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |
**Last Updated**: 2026-01-30
**Location**: `internal/ticketing/`
**Dependencies**: User authentication, email notifications, webhook triggers

---

## Executive Summary

**Problem**: Users need a structured way to report bugs, request features, and get support. Admins need to triage, assign, and track resolution of issues.

**Solution**: Built-in ticketing system with user portal (submit tickets), admin interface (manage tickets), and integration with GitHub Issues (optional two-way sync).

**Key Features**:
- User-friendly ticket submission (no GitHub account required)
- Admin triage & assignment
- Priority & status tracking
- Email notifications
- Webhook triggers (e.g., Discord/Slack)
- Optional GitHub Issues sync

---

## User Portal

### Ticket Submission Form

**URL**: `/support/new`

**Fields**:
- **Type**: Bug Report | Feature Request | Support Question
- **Title**: One-line summary (max 200 chars)
- **Description**: Detailed description (Markdown supported)
- **Severity** (Bug reports only): Critical | High | Medium | Low
- **Category**: Movies | TV Shows | Music | Books | Adult | Other
- **Attachments**: Screenshots, logs (max 10MB per file, 5 files max)
- **Environment** (Auto-detected):
  - Revenge version
  - Browser/Client (user-agent)
  - OS (from user-agent)
  - Server OS (from server logs)

**Validation**:
- Title: Required, 10-200 chars
- Description: Required, min 20 chars
- Attachments: Optional, max 10MB per file, allowed types: `.png`, `.jpg`, `.log`, `.txt`

### Ticket Viewing

**URL**: `/support/tickets` (user's own tickets)

**List View**:
- Ticket ID (e.g., `#1234`)
- Type (icon + label)
- Title
- Status (Open, In Progress, Resolved, Closed)
- Created date
- Last updated date

**Detail View** (`/support/tickets/{id}`):
- All form fields (read-only for users after submission)
- Comments/replies (threaded conversation)
- Status history (timeline)
- Assigned admin (if any)

### User Actions
- **Create ticket**: Submit new ticket
- **Add comment**: Reply to own ticket
- **Mark resolved**: User accepts resolution (moves to "Closed")
- **Reopen**: If issue persists after resolution

---

## Admin Interface

### Triage Dashboard

**URL**: `/admin/support`

**Views**:
1. **Unassigned** (default): New tickets awaiting triage
2. **Assigned to Me**: Tickets assigned to current admin
3. **All Open**: All open tickets
4. **Resolved**: Tickets marked resolved (awaiting user confirmation)
5. **Closed**: User-confirmed resolutions

**Filters**:
- Type (Bug, Feature, Support)
- Status (Open, In Progress, Resolved, Closed)
- Priority (Critical → Low)
- Category (Movies, TV, Music, etc.)
- Date range (created/updated)

**Sorting**:
- Priority (Critical first)
- Created date (newest first)
- Updated date (most recent activity)

### Admin Actions

**Per-Ticket**:
- **Assign**: Assign to self or another admin
- **Change Priority**: Critical | High | Medium | Low
- **Change Status**: Open → In Progress → Resolved → Closed
- **Add Labels**: Bug, Wontfix, Duplicate, Enhancement, etc.
- **Add Comment**: Internal notes (private) OR public replies (visible to user)
- **Link to GitHub Issue**: Create/link GitHub issue for tracking
- **Close**: Mark as resolved with resolution notes

**Bulk Actions** (multi-select):
- Assign to admin
- Change priority
- Change status
- Add label
- Close with bulk resolution note

---

## PostgreSQL Schema

```sql
-- Tickets table
CREATE TABLE tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_number SERIAL UNIQUE NOT NULL,    -- Human-readable ID (e.g., #1234)
    user_id UUID REFERENCES users(id) ON DELETE SET NULL, -- NULL if user deleted
    type VARCHAR(50) NOT NULL CHECK (type IN ('bug', 'feature', 'support')),
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    severity VARCHAR(20) CHECK (severity IN ('critical', 'high', 'medium', 'low')), -- Bug reports only
    category VARCHAR(50),                    -- Movies, TV, Music, etc.
    status VARCHAR(20) NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'in_progress', 'resolved', 'closed')),
    priority INT DEFAULT 3 CHECK (priority BETWEEN 1 AND 5), -- 1=Critical, 5=Low
    assigned_admin_id UUID REFERENCES users(id) ON DELETE SET NULL, -- Admin assigned to ticket
    github_issue_url VARCHAR(500),           -- Optional GitHub issue link
    metadata_json JSONB,                     -- Environment info (version, browser, OS)
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMPTZ,
    closed_at TIMESTAMPTZ
);

-- Ticket comments/replies
CREATE TABLE ticket_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID REFERENCES tickets(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL, -- Commenter (user or admin)
    comment TEXT NOT NULL,
    is_internal BOOLEAN DEFAULT FALSE,       -- TRUE = admin-only note, FALSE = public reply
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Ticket attachments
CREATE TABLE ticket_attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID REFERENCES tickets(id) ON DELETE CASCADE,
    filename VARCHAR(500) NOT NULL,
    file_path VARCHAR(1000) NOT NULL,        -- Storage path (local or S3)
    file_size_bytes INT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    uploaded_at TIMESTAMPTZ DEFAULT NOW()
);

-- Ticket labels (tags)
CREATE TABLE ticket_labels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    color VARCHAR(7),                        -- Hex color (e.g., #ff0000)
    description TEXT
);

-- Ticket-Label relationships (many-to-many)
CREATE TABLE ticket_label_assignments (
    ticket_id UUID REFERENCES tickets(id) ON DELETE CASCADE,
    label_id UUID REFERENCES ticket_labels(id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (ticket_id, label_id)
);

-- Status history (audit trail)
CREATE TABLE ticket_status_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID REFERENCES tickets(id) ON DELETE CASCADE,
    old_status VARCHAR(20),
    new_status VARCHAR(20) NOT NULL,
    changed_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    changed_at TIMESTAMPTZ DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_tickets_user_id ON tickets(user_id);
CREATE INDEX idx_tickets_assigned_admin_id ON tickets(assigned_admin_id);
CREATE INDEX idx_tickets_status ON tickets(status);
CREATE INDEX idx_tickets_priority ON tickets(priority);
CREATE INDEX idx_tickets_created_at ON tickets(created_at DESC);
CREATE INDEX idx_ticket_comments_ticket_id ON ticket_comments(ticket_id);
```

---

## API Endpoints

### User Endpoints

```bash
# Create ticket
POST /api/v1/support/tickets
Request:
{
  "type": "bug",
  "title": "Video playback stutters on Chrome",
  "description": "When playing...",
  "severity": "high",
  "category": "movies",
  "attachments": ["file1.png", "file2.log"]
}
Response: { "id": "uuid", "ticket_number": 1234, "status": "open" }

# Get user's tickets
GET /api/v1/support/tickets?status=open&type=bug
Response: [{ "id": "uuid", "ticket_number": 1234, ... }]

# Get ticket detail
GET /api/v1/support/tickets/{id}
Response: { "id": "uuid", "title": "...", "comments": [...], ... }

# Add comment to ticket
POST /api/v1/support/tickets/{id}/comments
Request: { "comment": "Still happening after update" }

# Mark ticket as resolved (user accepts)
PUT /api/v1/support/tickets/{id}/resolve
```

### Admin Endpoints

```bash
# Get all tickets (admin view)
GET /api/v1/admin/support/tickets?status=open&assigned_to=me
Response: [{ "id": "uuid", "ticket_number": 1234, ... }]

# Assign ticket
PUT /api/v1/admin/support/tickets/{id}/assign
Request: { "admin_id": "uuid" }

# Change status
PUT /api/v1/admin/support/tickets/{id}/status
Request: { "status": "in_progress" }

# Change priority
PUT /api/v1/admin/support/tickets/{id}/priority
Request: { "priority": 1 }

# Add labels
POST /api/v1/admin/support/tickets/{id}/labels
Request: { "label_ids": ["uuid1", "uuid2"] }

# Add admin comment (internal note)
POST /api/v1/admin/support/tickets/{id}/comments
Request: { "comment": "Reproduced locally", "is_internal": true }

# Close ticket with resolution
PUT /api/v1/admin/support/tickets/{id}/close
Request: { "resolution": "Fixed in v1.2.3" }

# Bulk actions
POST /api/v1/admin/support/tickets/bulk
Request: { "ticket_ids": ["uuid1", "uuid2"], "action": "assign", "admin_id": "uuid" }
```

---

## Email Notifications

### User Notifications
- **Ticket created**: Confirmation email with ticket number
- **Admin reply**: Email when admin adds public comment
- **Status change**: Notify when status changes (Open → In Progress → Resolved)
- **Ticket closed**: Notify when admin closes ticket with resolution

### Admin Notifications
- **New ticket**: Notify all admins (or specific team)
- **Assigned**: Notify when ticket assigned to admin
- **User reply**: Notify assigned admin when user adds comment
- **Escalation**: Notify when Critical ticket is open for >4 hours

**Email Template Example** (User - Admin Reply):
```
Subject: [Revenge Support #1234] Admin replied to your ticket

Hi {user_name},

An admin has replied to your support ticket:

Ticket #1234: Video playback stutters on Chrome
Status: In Progress
Priority: High

Admin Reply:
---
We've reproduced the issue and identified the cause. A fix will be included in the next release (v1.2.3).
---

View ticket: https://revenge.local/support/tickets/1234

Thanks,
Revenge Support Team
```

---

## Webhook Triggers

### Events
- `ticket.created`: New ticket submitted
- `ticket.assigned`: Ticket assigned to admin
- `ticket.status_changed`: Status changed
- `ticket.commented`: New comment added
- `ticket.resolved`: Ticket marked resolved
- `ticket.closed`: Ticket closed

### Payload Example
```json
{
  "event": "ticket.created",
  "timestamp": "2026-01-28T12:00:00Z",
  "ticket": {
    "id": "uuid",
    "ticket_number": 1234,
    "type": "bug",
    "title": "Video playback stutters",
    "severity": "high",
    "status": "open",
    "user": { "id": "uuid", "username": "johndoe" },
    "url": "https://revenge.local/support/tickets/1234"
  }
}
```

### Integrations
- **Discord**: Post to `#support` channel
- **Slack**: Post to `#tickets` channel
- **GitHub**: Create GitHub issue automatically
- **Custom**: POST to user-defined webhook URL

---

## GitHub Issues Sync (Optional)

### Two-Way Sync
- **Revenge → GitHub**: Create GitHub issue when ticket created
- **GitHub → Revenge**: Update ticket status when issue closed

### GitHub Issue Metadata
```markdown
<!-- Revenge Ticket Metadata -->
revenge_ticket_id: uuid
revenge_ticket_number: 1234
revenge_user: johndoe
revenge_url: https://revenge.local/support/tickets/1234
```

### Sync Rules
1. **Create**: Revenge ticket → GitHub issue (with label `revenge-ticket`)
2. **Update**: Comment on Revenge ticket → Comment on GitHub issue
3. **Close**: GitHub issue closed → Revenge ticket marked "Resolved"
4. **Reopen**: Revenge ticket reopened → Reopen GitHub issue

---

## UI/UX Design

### User Ticket Submission Form
- **Layout**: Single-page form with sections
- **Validation**: Real-time validation (red borders on errors)
- **Markdown Preview**: Live preview of description
- **Attachments**: Drag-and-drop file upload
- **Auto-Save**: Draft saved every 30 seconds (localStorage)

### Admin Dashboard
- **Kanban Board**: Columns for Open, In Progress, Resolved, Closed
- **Drag-and-Drop**: Move tickets between columns
- **Quick Actions**: Hover menu for Assign, Priority, Labels
- **Filters**: Sticky sidebar with filter checkboxes
- **Bulk Select**: Checkbox column for bulk actions

### Ticket Detail Page
- **Split Layout**: Left = Ticket info, Right = Comments/Activity
- **Comment Thread**: Threaded replies (admin + user)
- **Timeline**: Visual timeline of status changes
- **Admin Toolbox**: Floating action buttons (Assign, Priority, Labels, Close)

---

## Implementation Phases

### Phase 1: Backend (Week 1)
- [ ] Create PostgreSQL schema (tickets, ticket_comments, ticket_attachments, ticket_labels, ticket_status_history)
- [ ] Implement API endpoints (user + admin)
- [ ] File upload service (attachments)

### Phase 2: User Portal (Week 2)
- [ ] Ticket submission form (Svelte 5 + shadcn-svelte)
- [ ] Ticket list view (user's own tickets)
- [ ] Ticket detail view (with comments)

### Phase 3: Admin Interface (Week 3)
- [ ] Admin dashboard (Kanban board)
- [ ] Triage view (filters, sorting, bulk actions)
- [ ] Admin toolbox (assign, priority, labels, close)

### Phase 4: Notifications (Week 4)
- [ ] Email templates
- [ ] SMTP integration
- [ ] Webhook triggers (Discord, Slack)

### Phase 5: GitHub Sync (Optional, Week 5)
- [ ] GitHub API integration
- [ ] Two-way sync logic
- [ ] Metadata embedding

**Total Estimated Time**: 4-5 weeks

---

## Implementation Checklist

### Phase 1: Core Infrastructure
- [ ] Create package structure at `internal/ticketing/`
- [ ] Define ticket entity (`entity.go`)
- [ ] Define comment entity
- [ ] Define attachment entity
- [ ] Define label entity
- [ ] Create repository interface (`repository.go`)
- [ ] Implement PostgreSQL repository (`postgres_repository.go`)
- [ ] Create fx module (`module.go`)
- [ ] Add configuration structs

### Phase 2: Database
- [ ] Create migration for `tickets` table
- [ ] Create migration for `ticket_comments` table
- [ ] Create migration for `ticket_attachments` table
- [ ] Create migration for `ticket_labels` table
- [ ] Create migration for `ticket_label_assignments` table
- [ ] Create migration for `ticket_status_history` table
- [ ] Add indexes for user_id, assigned_admin_id, status, priority, created_at
- [ ] Write sqlc queries for ticket CRUD
- [ ] Write sqlc queries for comment operations
- [ ] Write sqlc queries for attachment operations
- [ ] Write sqlc queries for label management
- [ ] Write sqlc queries for status history tracking

### Phase 3: Service Layer
- [ ] Implement TicketService (`service.go`)
- [ ] Add ticket creation logic with validation
- [ ] Add ticket update logic (status, priority, assignment)
- [ ] Add comment management (public/internal)
- [ ] Add label assignment logic
- [ ] Add attachment upload/delete logic
- [ ] Implement status transition validation
- [ ] Add caching for ticket lists (Redis)
- [ ] Implement bulk actions (assign, priority, status, close)
- [ ] Add ticket search/filter functionality

### Phase 4: Background Jobs
- [ ] Create River job for escalation alerts (Critical tickets >4h)
- [ ] Create River job for auto-close stale resolved tickets
- [ ] Create River job for attachment cleanup
- [ ] Create River job for GitHub sync (if enabled)
- [ ] Register jobs in fx module

### Phase 5: API Integration
- [ ] Add OpenAPI spec for user ticket endpoints
- [ ] Add OpenAPI spec for admin ticket endpoints
- [ ] Generate ogen handlers
- [ ] Implement user endpoints:
  - [ ] POST /api/v1/support/tickets (create)
  - [ ] GET /api/v1/support/tickets (list own)
  - [ ] GET /api/v1/support/tickets/{id} (detail)
  - [ ] POST /api/v1/support/tickets/{id}/comments (add comment)
  - [ ] PUT /api/v1/support/tickets/{id}/resolve (mark resolved)
- [ ] Implement admin endpoints:
  - [ ] GET /api/v1/admin/support/tickets (list all)
  - [ ] PUT /api/v1/admin/support/tickets/{id}/assign
  - [ ] PUT /api/v1/admin/support/tickets/{id}/status
  - [ ] PUT /api/v1/admin/support/tickets/{id}/priority
  - [ ] POST /api/v1/admin/support/tickets/{id}/labels
  - [ ] POST /api/v1/admin/support/tickets/{id}/comments (internal)
  - [ ] PUT /api/v1/admin/support/tickets/{id}/close
  - [ ] POST /api/v1/admin/support/tickets/bulk
- [ ] Add authentication middleware
- [ ] Add RBAC permission checks
- [ ] Integrate email notifications
- [ ] Integrate webhook triggers

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [PostgreSQL Arrays](https://www.postgresql.org/docs/current/arrays.html) | [Local](../../../sources/database/postgresql-arrays.md) |
| [PostgreSQL JSON Functions](https://www.postgresql.org/docs/current/functions-json.html) | [Local](../../../sources/database/postgresql-json.md) |
| [River Job Queue](https://pkg.go.dev/github.com/riverqueue/river) | [Local](../../../sources/tooling/river.md) |
| [Svelte 5 Documentation](https://svelte.dev/docs/svelte/overview) | [Local](../../../sources/frontend/svelte5.md) |
| [Svelte 5 Runes](https://svelte.dev/docs/svelte/$state) | [Local](../../../sources/frontend/svelte-runes.md) |
| [SvelteKit Documentation](https://svelte.dev/docs/kit/introduction) | [Local](../../../sources/frontend/sveltekit.md) |
| [Uber fx](https://pkg.go.dev/go.uber.org/fx) | [Local](../../../sources/tooling/fx.md) |
| [ogen OpenAPI Generator](https://pkg.go.dev/github.com/ogen-go/ogen) | [Local](../../../sources/tooling/ogen.md) |
| [pgx PostgreSQL Driver](https://pkg.go.dev/github.com/jackc/pgx/v5) | [Local](../../../sources/database/pgx.md) |
| [rueidis](https://pkg.go.dev/github.com/redis/rueidis) | [Local](../../../sources/tooling/rueidis.md) |
| [shadcn-svelte](https://www.shadcn-svelte.com/docs) | [Local](../../../sources/frontend/shadcn-svelte.md) |
| [sqlc](https://docs.sqlc.dev/en/stable/) | [Local](../../../sources/database/sqlc.md) |
| [sqlc Configuration](https://docs.sqlc.dev/en/stable/reference/config.html) | [Local](../../../sources/database/sqlc-config.md) |

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Shared](INDEX.md)

### In This Section

- [Time-Based Access Controls](ACCESS_CONTROLS.md)
- [Tracearr Analytics Service](ANALYTICS_SERVICE.md)
- [Revenge - Client Support & Device Capabilities](CLIENT_SUPPORT.md)
- [Content Rating System](CONTENT_RATING.md)
- [Revenge - Internationalization (i18n)](I18N.md)
- [Library Types](LIBRARY_TYPES.md)
- [News System](NEWS_SYSTEM.md)
- [Revenge - NSFW Toggle](NSFW_TOGGLE.md)

### Related Topics

- [Revenge - Architecture v2](../../architecture/01_ARCHITECTURE.md) _Architecture_
- [Revenge - Design Principles](../../architecture/02_DESIGN_PRINCIPLES.md) _Architecture_
- [Revenge - Metadata System](../../architecture/03_METADATA_SYSTEM.md) _Architecture_
- [Revenge - Player Architecture](../../architecture/04_PLAYER_ARCHITECTURE.md) _Architecture_
- [Plugin Architecture Decision](../../architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) _Architecture_

### Indexes

- [Design Index](../../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

## Related

- [Auditing System](AUDITING_SYSTEM.md) - Error logs, metadata conflicts, moderation
- [User Management](../operations/USER_MANAGEMENT.md) - User roles & permissions
- [Email Notifications](../technical/EMAIL.md) - SMTP configuration
- [Webhook Configuration](../technical/WEBHOOKS.md) - Webhook setup

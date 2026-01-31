# Time-Based Access Controls

> User access restrictions based on time, limits, and schedules

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | ✅ | Full design with DB schema, Go implementation, API endpoints |
| Sources | 🟡 | Inspired by Emby Parental Controls |
| Instructions | ✅ | Implementation checklist complete |
| Code | 🔴 | |
| Linting | 🔴 | |
| Unit Testing | 🔴 | |
| Integration Testing | 🔴 | |

**Priority**: 🟡 MEDIUM (Emby has this)
**Inspired By**: Emby Parental Controls
**Location**: `internal/service/access/`

---

## Overview

Time-based access controls allow administrators and parents to restrict when users can access content, set viewing time limits, and enforce access schedules.

---

## Features

| Feature | Description |
|---------|-------------|
| Time Limits | Daily/weekly viewing time limits |
| Access Schedules | Allowed viewing hours (e.g., 6pm-9pm) |
| Automatic Logout | Force logout after time expires |
| Bedtime Mode | Block access after specific time |
| Device Limits | Concurrent stream limits |
| Content Lockout | Time-based content restrictions |

---

## Use Cases

### Parental Controls

```
Child Account "Timmy":
- Max 2 hours/day viewing
- Only allowed 4:00 PM - 8:00 PM weekdays
- Only allowed 10:00 AM - 9:00 PM weekends
- No access after 8:30 PM (bedtime)
- Max 1 concurrent stream
```

### Guest Accounts

```
Guest Account:
- Max 4 hours total (lifetime limit)
- Expires after 7 days
- Max 1 concurrent stream
```

### Household Rules

```
Global Rule:
- No streaming after midnight
- Max 5 concurrent streams household-wide
```

---

## Database Schema

```sql
-- User access rules
CREATE TABLE access_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,

    -- Rule type
    rule_type VARCHAR(50) NOT NULL, -- time_limit, schedule, bedtime, device_limit, expiry

    -- Time limits (daily/weekly)
    daily_limit_minutes INT,
    weekly_limit_minutes INT,
    reset_time TIME DEFAULT '00:00:00',

    -- Schedule (allowed hours)
    schedule_enabled BOOLEAN DEFAULT false,
    weekday_start TIME, -- e.g., 16:00
    weekday_end TIME,   -- e.g., 20:00
    weekend_start TIME,
    weekend_end TIME,

    -- Bedtime
    bedtime_enabled BOOLEAN DEFAULT false,
    bedtime_weekday TIME,
    bedtime_weekend TIME,
    bedtime_warning_minutes INT DEFAULT 15,

    -- Device/stream limits
    max_concurrent_streams INT,

    -- Account expiry
    expires_at TIMESTAMPTZ,
    max_total_minutes INT, -- Lifetime limit

    -- Settings
    is_enabled BOOLEAN DEFAULT true,
    enforcement VARCHAR(20) DEFAULT 'soft', -- soft (warning), hard (block)
    notify_admin BOOLEAN DEFAULT false,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Usage tracking
CREATE TABLE access_usage (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,

    -- Daily tracking
    watch_seconds BIGINT DEFAULT 0,
    session_count INT DEFAULT 0,

    -- Running totals
    total_watch_seconds BIGINT DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(user_id, date)
);

-- Access violations
CREATE TABLE access_violations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    rule_id UUID REFERENCES access_rules(id) ON DELETE SET NULL,

    violation_type VARCHAR(50) NOT NULL, -- time_exceeded, schedule_violation, bedtime_violation
    details JSONB,
    action_taken VARCHAR(50), -- warned, blocked, logged_out

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_access_rules_user ON access_rules(user_id);
CREATE INDEX idx_access_usage_user_date ON access_usage(user_id, date DESC);
CREATE INDEX idx_access_violations_user ON access_violations(user_id);
```

---

## Go Implementation

```go
// internal/service/access/

type Service struct {
    rules    RuleRepository
    usage    UsageRepository
    sessions SessionService
}

type AccessChecker struct {
    rules []AccessRule
    usage *DailyUsage
}

func (c *AccessChecker) CanAccess(now time.Time) (bool, *Violation) {
    for _, rule := range c.rules {
        if !rule.IsEnabled {
            continue
        }

        // Check schedule
        if rule.ScheduleEnabled {
            if !c.isWithinSchedule(now, rule) {
                return false, &Violation{
                    Type:    "schedule_violation",
                    Message: "Outside allowed hours",
                }
            }
        }

        // Check bedtime
        if rule.BedtimeEnabled {
            if c.isPastBedtime(now, rule) {
                return false, &Violation{
                    Type:    "bedtime_violation",
                    Message: "Past bedtime",
                }
            }
        }

        // Check time limit
        if rule.DailyLimitMinutes > 0 {
            usedMinutes := c.usage.WatchSeconds / 60
            if int(usedMinutes) >= rule.DailyLimitMinutes {
                return false, &Violation{
                    Type:    "time_exceeded",
                    Message: "Daily limit reached",
                }
            }
        }

        // Check expiry
        if rule.ExpiresAt != nil && now.After(*rule.ExpiresAt) {
            return false, &Violation{
                Type:    "account_expired",
                Message: "Account has expired",
            }
        }
    }

    return true, nil
}

func (c *AccessChecker) isWithinSchedule(now time.Time, rule AccessRule) bool {
    weekday := now.Weekday()
    currentTime := now.Format("15:04:05")

    var start, end string
    if weekday == time.Saturday || weekday == time.Sunday {
        start = rule.WeekendStart
        end = rule.WeekendEnd
    } else {
        start = rule.WeekdayStart
        end = rule.WeekdayEnd
    }

    return currentTime >= start && currentTime <= end
}

func (c *AccessChecker) GetRemainingTime() time.Duration {
    for _, rule := range c.rules {
        if rule.DailyLimitMinutes > 0 {
            usedMinutes := c.usage.WatchSeconds / 60
            remaining := rule.DailyLimitMinutes - int(usedMinutes)
            if remaining < 0 {
                remaining = 0
            }
            return time.Duration(remaining) * time.Minute
        }
    }
    return -1 // No limit
}
```

---

## API Endpoints

```
# Rules (admin/parent)
GET  /api/v1/access/rules                    # List all rules
GET  /api/v1/access/users/:user_id/rules     # Get user's rules
POST /api/v1/access/users/:user_id/rules     # Create rule
PUT  /api/v1/access/rules/:id                # Update rule
DELETE /api/v1/access/rules/:id              # Delete rule

# Usage
GET  /api/v1/access/users/:user_id/usage     # Get user's usage
GET  /api/v1/access/users/:user_id/usage/today # Today's usage
GET  /api/v1/access/users/:user_id/usage/week  # This week's usage

# Violations
GET  /api/v1/access/users/:user_id/violations # Get user's violations

# Check access (called before playback)
GET  /api/v1/access/check                    # Check current user's access
```

---

## Client Integration

### Playback Enforcement

```typescript
// Check access before playback
const checkAccess = async () => {
    const response = await fetch('/api/v1/access/check');
    const result = await response.json();

    if (!result.canAccess) {
        showAccessDeniedModal(result.violation);
        return false;
    }

    if (result.warningMinutes > 0 && result.warningMinutes <= 15) {
        showTimeWarning(`${result.warningMinutes} minutes remaining`);
    }

    return true;
};

// Periodic check during playback
setInterval(async () => {
    const result = await checkAccess();
    if (!result) {
        player.pause();
    }
}, 60000); // Check every minute
```

### Warning UI

```
┌─────────────────────────────────────────┐
│  ⚠️ Time Warning                        │
│                                         │
│  You have 15 minutes of viewing time    │
│  remaining today.                       │
│                                         │
│            [ OK ]                       │
└─────────────────────────────────────────┘
```

### Blocked UI

```
┌─────────────────────────────────────────┐
│  🚫 Access Restricted                   │
│                                         │
│  You have reached your daily viewing    │
│  limit.                                 │
│                                         │
│  Your limit will reset at 12:00 AM.     │
│                                         │
│        [ Contact Parent ]               │
└─────────────────────────────────────────┘
```

---

## RBAC Permissions

| Permission | Description |
|------------|-------------|
| `access.rules.view` | View access rules |
| `access.rules.manage` | Create/edit rules |
| `access.usage.view` | View usage stats |
| `access.bypass` | Bypass all access restrictions |

---

## Configuration

```yaml
access_controls:
  enabled: true
  default_enforcement: soft  # soft or hard

  warnings:
    time_remaining_minutes: [30, 15, 5, 1]
    bedtime_minutes: 15

  tracking:
    update_interval_seconds: 60
    include_paused_time: false

  notifications:
    admin_on_violation: false
    parent_on_limit_reached: true
```

---

## Implementation Checklist

### Phase 1: Core Infrastructure
- [ ] Create package structure `internal/service/access/`
- [ ] Define `AccessRule` entity with all rule types (time_limit, schedule, bedtime, device_limit, expiry)
- [ ] Define `DailyUsage` entity for tracking
- [ ] Define `Violation` entity for violations log
- [ ] Create repository interfaces for rules, usage, and violations
- [ ] Register fx module `internal/service/access/module.go`

### Phase 2: Database
- [ ] Create migration `shared/000XXX_access_rules.up.sql`
- [ ] Create `access_rules` table with schedule, bedtime, limit columns
- [ ] Create `access_usage` table with daily/weekly tracking
- [ ] Create `access_violations` table for audit log
- [ ] Add indexes on user_id, date for efficient queries
- [ ] Generate sqlc queries for CRUD operations
- [ ] Generate queries for usage aggregation

### Phase 3: Service Layer
- [ ] Implement `AccessChecker` with rule evaluation logic
- [ ] Implement schedule checking (weekday vs weekend hours)
- [ ] Implement bedtime enforcement with warning minutes
- [ ] Implement daily/weekly time limit tracking
- [ ] Implement `GetRemainingTime()` for client warnings
- [ ] Add caching for user rules (Redis with invalidation)
- [ ] Implement concurrent stream limit checking

### Phase 4: Background Jobs
- [ ] Create River job for daily usage reset at `reset_time`
- [ ] Create River job for weekly usage rollup
- [ ] Create River job for expired account cleanup
- [ ] Create notification job for limit warnings
- [ ] Schedule periodic access check during active sessions

### Phase 5: API Integration
- [ ] Add OpenAPI schema for access rules endpoints
- [ ] Implement `GET /api/v1/access/rules` (admin view all)
- [ ] Implement `GET /api/v1/access/users/:user_id/rules`
- [ ] Implement `POST /api/v1/access/users/:user_id/rules`
- [ ] Implement `PUT /api/v1/access/rules/:id`
- [ ] Implement `DELETE /api/v1/access/rules/:id`
- [ ] Implement `GET /api/v1/access/check` (current user access check)
- [ ] Implement usage endpoints (today, week, history)
- [ ] Add middleware for playback access enforcement
- [ ] Add RBAC permissions (`access.rules.view`, `access.rules.manage`, `access.bypass`)

---


<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

<!-- SOURCE-BREADCRUMBS-END -->

<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Shared](INDEX.md)

### In This Section

- [Tracearr Analytics Service](ANALYTICS_SERVICE.md)
- [Revenge - Client Support & Device Capabilities](CLIENT_SUPPORT.md)
- [Content Rating System](CONTENT_RATING.md)
- [Revenge - Internationalization (i18n)](I18N.md)
- [Library Types](LIBRARY_TYPES.md)
- [News System](NEWS_SYSTEM.md)
- [Revenge - NSFW Toggle](NSFW_TOGGLE.md)
- [Dynamic RBAC with Casbin](RBAC_CASBIN.md)

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

- [RBAC with Casbin](RBAC_CASBIN.md)
- [User Experience Features](USER_EXPERIENCE_FEATURES.md)
- [Analytics Service](ANALYTICS_SERVICE.md)

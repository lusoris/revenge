# Time-Based Access Controls

> User access restrictions based on time, limits, and schedules

**Status**: 🔴 PLANNING
**Priority**: 🟡 MEDIUM (Emby has this)
**Inspired By**: Emby Parental Controls

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

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
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

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),

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

    created_at TIMESTAMPTZ DEFAULT NOW()
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

## Related Documentation

- [RBAC with Casbin](RBAC_CASBIN.md)
- [User Experience Features](USER_EXPERIENCE_FEATURES.md)
- [Analytics Service](ANALYTICS_SERVICE.md)

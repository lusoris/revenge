# Revenge - NSFW Toggle

> User preference component for adult content visibility.
> Referenced by [WHISPARR_STASHDB_SCHEMA.md](WHISPARR_STASHDB_SCHEMA.md) and [ADULT_CONTENT_SYSTEM.md](ADULT_CONTENT_SYSTEM.md).


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Overview](#overview)
- [Architecture](#architecture)
- [User Preference Storage](#user-preference-storage)
  - [Database Schema](#database-schema)
  - [Cache Layer](#cache-layer)
- [API Middleware](#api-middleware)
  - [Route Protection](#route-protection)
  - [Search Filtering](#search-filtering)
- [Frontend Components](#frontend-components)
  - [Toggle Component (Svelte)](#toggle-component-svelte)
  - [Sidebar Integration](#sidebar-integration)
  - [Settings Page](#settings-page)
- [Auto-Lock Feature](#auto-lock-feature)
  - [Backend Implementation](#backend-implementation)
  - [Activity Tracking](#activity-tracking)
- [Security Considerations](#security-considerations)
  - [PIN Storage](#pin-storage)
  - [Audit Logging](#audit-logging)
- [API Endpoints](#api-endpoints)
- [Behavior Summary](#behavior-summary)
- [Implementation Checklist](#implementation-checklist)
  - [Phase 1: Core Infrastructure](#phase-1-core-infrastructure)
  - [Phase 2: Database](#phase-2-database)
  - [Phase 3: Service Layer](#phase-3-service-layer)
  - [Phase 4: Background Jobs](#phase-4-background-jobs)
  - [Phase 5: Middleware](#phase-5-middleware)
  - [Phase 6: API Integration](#phase-6-api-integration)
  - [Phase 7: Search Integration](#phase-7-search-integration)
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
| Design | ✅ | Full design with DB schema, middleware, Svelte components |
| Sources | 🟡 |  |
| Instructions | ✅ | Implementation checklist added |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |
---

## Overview

The NSFW toggle controls visibility and access to adult content modules (`adult_movie`, `adult_scene`) stored in PostgreSQL schema `qar`.

**Default State:** OFF (explicit opt-in required)

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     NSFW Toggle Flow                         │
└─────────────────────────────────────────────────────────────┘

┌──────────┐     ┌─────────────────┐     ┌─────────────────────┐
│  Toggle  │ ──→ │  User Settings  │ ──→ │  API Middleware     │
│   UI     │     │  (Database)     │     │  (Access Control)   │
└──────────┘     └─────────────────┘     └─────────────────────┘
                        │                          │
                        ▼                          ▼
              ┌─────────────────┐       ┌──────────────────────┐
              │  Session Cache  │       │  Route Filtering     │
              │  (Dragonfly)    │       │  (/qar/* visibility) │
              └─────────────────┘       └──────────────────────┘
                                                   │
                                                   ▼
                                        ┌──────────────────────┐
                                        │  Content Filtering   │
                                        │  (Search, Dashboard) │
                                        └──────────────────────┘
```

---

## User Preference Storage

### Database Schema

```sql
-- User preferences table (in public schema)
CREATE TABLE user_preferences (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,

    -- NSFW Settings
    nsfw_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    nsfw_pin VARCHAR(6),  -- Optional 4-6 digit PIN for additional security
    nsfw_timeout_minutes INTEGER DEFAULT 30,  -- Auto-lock after inactivity

    -- Other preferences
    theme VARCHAR(20) DEFAULT 'system',
    language VARCHAR(10) DEFAULT 'en',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Audit log for NSFW toggles (security/compliance)
CREATE TABLE nsfw_toggle_audit (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    action VARCHAR(20) NOT NULL,  -- 'enabled', 'disabled', 'pin_set', 'pin_verified'
    ip_address INET,
    user_agent TEXT,
    timestamp TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_nsfw_audit_user ON nsfw_toggle_audit(user_id, timestamp DESC);
```

### Cache Layer

```go
// Cache NSFW state in session for fast middleware checks
const (
    nsfwCacheKey    = "user:%s:nsfw"
    nsfwCacheTTL    = 5 * time.Minute
)

func (s *SessionService) GetNSFWEnabled(ctx context.Context, userID uuid.UUID) (bool, error) {
    key := fmt.Sprintf(nsfwCacheKey, userID)

    // Check cache first
    val, err := s.cache.Get(ctx, key).Bool()
    if err == nil {
        return val, nil
    }

    // Fallback to database
    enabled, err := s.repo.GetNSFWEnabled(ctx, userID)
    if err != nil {
        return false, err
    }

    // Update cache
    s.cache.Set(ctx, key, enabled, nsfwCacheTTL)
    return enabled, nil
}
```

---

## API Middleware

### Route Protection

```go
// NSFWMiddleware blocks /qar/* routes when NSFW is disabled
func NSFWMiddleware(sessions *SessionService) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Only apply to /qar/* routes
            if !strings.HasPrefix(r.URL.Path, "/api/v1/qar/") {
                next.ServeHTTP(w, r)
                return
            }

            userID := auth.UserIDFromContext(r.Context())
            if userID == uuid.Nil {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }

            enabled, err := sessions.GetNSFWEnabled(r.Context(), userID)
            if err != nil || !enabled {
                // Return 404 to obscure existence of adult content
                http.NotFound(w, r)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

### Search Filtering

```go
// SearchService filters adult content based on NSFW setting
func (s *SearchService) Search(ctx context.Context, query string, opts SearchOpts) (*SearchResults, error) {
    userID := auth.UserIDFromContext(ctx)
    nsfwEnabled, _ := s.sessions.GetNSFWEnabled(ctx, userID)

    results := &SearchResults{}

    // Search non-adult modules
    results.Movies = s.searchMovies(ctx, query, opts)
    results.Shows = s.searchShows(ctx, query, opts)
    results.Music = s.searchMusic(ctx, query, opts)
    // ... other modules

    // Only include adult results if enabled
    if nsfwEnabled {
        results.AdultMovies = s.searchAdultMovies(ctx, query, opts)
        results.AdultScenes = s.searchAdultScenes(ctx, query, opts)
    }

    return results, nil
}
```

---

## Frontend Components

### Toggle Component (Svelte)

```svelte
<!-- NSFWToggle.svelte -->
<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    import { slide } from 'svelte/transition';
    import { userPreferences } from '$lib/stores/preferences';
    import { api } from '$lib/api';

    const dispatch = createEventDispatcher();

    let showPinModal = false;
    let pin = '';
    let error = '';
    let loading = false;

    $: nsfwEnabled = $userPreferences.nsfw_enabled;
    $: hasPin = $userPreferences.nsfw_pin_set;

    async function toggle() {
        if (nsfwEnabled) {
            // Disabling - no PIN required
            await updateNSFW(false);
        } else {
            // Enabling - check if PIN required
            if (hasPin) {
                showPinModal = true;
            } else {
                await updateNSFW(true);
            }
        }
    }

    async function verifyPin() {
        loading = true;
        error = '';

        try {
            const verified = await api.verifyNSFWPin(pin);
            if (verified) {
                await updateNSFW(true);
                showPinModal = false;
                pin = '';
            } else {
                error = 'Incorrect PIN';
            }
        } catch (e) {
            error = 'Verification failed';
        } finally {
            loading = false;
        }
    }

    async function updateNSFW(enabled: boolean) {
        loading = true;
        try {
            await api.updatePreferences({ nsfw_enabled: enabled });
            userPreferences.update(p => ({ ...p, nsfw_enabled: enabled }));
            dispatch('change', { enabled });
        } finally {
            loading = false;
        }
    }
</script>

<div class="nsfw-toggle">
    <button
        on:click={toggle}
        class:enabled={nsfwEnabled}
        disabled={loading}
        aria-label={nsfwEnabled ? 'Disable adult content' : 'Enable adult content'}
    >
        <span class="icon">{nsfwEnabled ? '🔓' : '🔒'}</span>
        <span class="label">{nsfwEnabled ? 'Adult: On' : 'Adult: Off'}</span>
    </button>
</div>

{#if showPinModal}
    <div class="modal-backdrop" transition:slide>
        <div class="modal">
            <h3>Enter PIN</h3>
            <input
                type="password"
                inputmode="numeric"
                pattern="[0-9]*"
                maxlength="6"
                bind:value={pin}
                🔴 Not implemented="Enter PIN"
                class:error={!!error}
            />
            {#if error}
                <p class="error-text">{error}</p>
            {/if}
            <div class="modal-actions">
                <button on:click={() => { showPinModal = false; pin = ''; }}>
                    Cancel
                </button>
                <button on:click={verifyPin} disabled={loading || pin.length < 4}>
                    Verify
                </button>
            </div>
        </div>
    </div>
{/if}

<style>
    .nsfw-toggle button {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.5rem 1rem;
        border-radius: 9999px;
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        cursor: pointer;
        transition: all 0.2s;
    }

    .nsfw-toggle button.enabled {
        background: var(--color-danger-subtle);
        border-color: var(--color-danger);
    }

    .modal-backdrop {
        position: fixed;
        inset: 0;
        background: rgba(0, 0, 0, 0.5);
        display: grid;
        place-items: center;
        z-index: 100;
    }

    .modal {
        background: var(--color-surface);
        padding: 1.5rem;
        border-radius: 0.5rem;
        min-width: 280px;
    }

    input.error {
        border-color: var(--color-danger);
    }

    .error-text {
        color: var(--color-danger);
        font-size: 0.875rem;
    }
</style>
```

### Sidebar Integration

```svelte
<!-- Sidebar.svelte (partial) -->
<script lang="ts">
    import { userPreferences } from '$lib/stores/preferences';

    $: nsfwEnabled = $userPreferences.nsfw_enabled;
</script>

<nav class="sidebar">
    <!-- Regular modules -->
    <NavItem href="/movies" icon="film">Movies</NavItem>
    <NavItem href="/shows" icon="tv">TV Shows</NavItem>
    <NavItem href="/music" icon="music">Music</NavItem>
    <!-- ... other modules ... -->

    <!-- Adult section (only visible when enabled) -->
    {#if nsfwEnabled}
        <div class="sidebar-divider" />
        <NavItem href="/qar/expeditions" icon="lock">Adult Movies</NavItem>
        <NavItem href="/qar/voyages" icon="lock">Adult Scenes</NavItem>
    {/if}
</nav>
```

### Settings Page

```svelte
<!-- SettingsNSFW.svelte -->
<script lang="ts">
    import { userPreferences } from '$lib/stores/preferences';
    import { api } from '$lib/api';

    let newPin = '';
    let confirmPin = '';
    let currentPin = '';
    let showSetPin = false;
    let timeout = $userPreferences.nsfw_timeout_minutes || 30;

    $: hasPin = $userPreferences.nsfw_pin_set;

    async function setPin() {
        if (newPin !== confirmPin) {
            return; // Show error
        }
        if (newPin.length < 4 || newPin.length > 6) {
            return; // Show error
        }

        await api.setNSFWPin(newPin, hasPin ? currentPin : undefined);
        userPreferences.update(p => ({ ...p, nsfw_pin_set: true }));
        newPin = '';
        confirmPin = '';
        currentPin = '';
        showSetPin = false;
    }

    async function removePin() {
        await api.removeNSFWPin(currentPin);
        userPreferences.update(p => ({ ...p, nsfw_pin_set: false }));
        currentPin = '';
    }

    async function updateTimeout() {
        await api.updatePreferences({ nsfw_timeout_minutes: timeout });
    }
</script>

<section class="settings-section">
    <h2>Adult Content</h2>

    <div class="setting-row">
        <div>
            <h3>PIN Protection</h3>
            <p class="description">
                Require a PIN to enable adult content. PIN is 4-6 digits.
            </p>
        </div>
        {#if hasPin}
            <button on:click={() => showSetPin = true}>Change PIN</button>
            <button class="danger" on:click={removePin}>Remove PIN</button>
        {:else}
            <button on:click={() => showSetPin = true}>Set PIN</button>
        {/if}
    </div>

    <div class="setting-row">
        <div>
            <h3>Auto-Lock Timeout</h3>
            <p class="description">
                Automatically disable adult content after inactivity.
            </p>
        </div>
        <select bind:value={timeout} on:change={updateTimeout}>
            <option value={0}>Never</option>
            <option value={15}>15 minutes</option>
            <option value={30}>30 minutes</option>
            <option value={60}>1 hour</option>
            <option value={240}>4 hours</option>
        </select>
    </div>

    {#if showSetPin}
        <!-- PIN setup modal -->
        <div class="modal">
            {#if hasPin}
                <label>
                    Current PIN
                    <input type="password" bind:value={currentPin} maxlength="6" />
                </label>
            {/if}
            <label>
                New PIN
                <input type="password" bind:value={newPin} maxlength="6" />
            </label>
            <label>
                Confirm PIN
                <input type="password" bind:value={confirmPin} maxlength="6" />
            </label>
            <button on:click={setPin}>Save</button>
            <button on:click={() => showSetPin = false}>Cancel</button>
        </div>
    {/if}
</section>
```

---

## Auto-Lock Feature

### Backend Implementation

```go
// Auto-lock NSFW after inactivity
type NSFWLockJob struct {
    river.WorkerDefaults[NSFWLockArgs]
    repo    PreferencesRepository
    cache   *redis.Client
}

type NSFWLockArgs struct {
    UserID uuid.UUID `json:"user_id"`
}

func (NSFWLockArgs) Kind() string { return "nsfw.auto_lock" }

func (w *NSFWLockJob) Work(ctx context.Context, job *river.Job[NSFWLockArgs]) error {
    // Check if user is still active
    lastActivity, err := w.cache.Get(ctx,
        fmt.Sprintf("user:%s:last_activity", job.Args.UserID)).Time()
    if err != nil || time.Since(lastActivity) > 30*time.Minute {
        // Lock NSFW
        return w.repo.SetNSFWEnabled(ctx, job.Args.UserID, false)
    }

    // Still active, reschedule
    return river.ErrSnooze(10 * time.Minute)
}
```

### Activity Tracking

```go
// Track user activity to prevent premature auto-lock
func ActivityMiddleware(cache *redis.Client) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if userID := auth.UserIDFromContext(r.Context()); userID != uuid.Nil {
                key := fmt.Sprintf("user:%s:last_activity", userID)
                cache.Set(r.Context(), key, time.Now(), 1*time.Hour)
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

---

## Security Considerations

### PIN Storage

```go
// PIN is hashed, not stored in plain text
func (s *PreferencesService) SetNSFWPin(ctx context.Context, userID uuid.UUID, pin string) error {
    hash, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    return s.repo.UpdateNSFWPin(ctx, userID, string(hash))
}

func (s *PreferencesService) VerifyNSFWPin(ctx context.Context, userID uuid.UUID, pin string) (bool, error) {
    hash, err := s.repo.GetNSFWPinHash(ctx, userID)
    if err != nil {
        return false, err
    }
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pin)) == nil, nil
}
```

### Audit Logging

All NSFW toggle actions are logged for compliance:

```go
func (s *PreferencesService) ToggleNSFW(ctx context.Context, userID uuid.UUID, enabled bool, r *http.Request) error {
    // Update preference
    if err := s.repo.SetNSFWEnabled(ctx, userID, enabled); err != nil {
        return err
    }

    // Invalidate cache
    s.cache.Del(ctx, fmt.Sprintf("user:%s:nsfw", userID))

    // Audit log
    action := "disabled"
    if enabled {
        action = "enabled"
    }
    return s.repo.LogNSFWAction(ctx, NSFWAuditEntry{
        UserID:    userID,
        Action:    action,
        IPAddress: r.RemoteAddr,
        UserAgent: r.UserAgent(),
        Timestamp: time.Now(),
    })
}
```

---

## API Endpoints

```yaml
# OpenAPI spec fragment
paths:
  /api/v1/user/preferences/nsfw:
    get:
      summary: Get NSFW preference status
      responses:
        200:
          content:
            application/json:
              schema:
                type: object
                properties:
                  enabled:
                    type: boolean
                  pin_set:
                    type: boolean
                  timeout_minutes:
                    type: integer

    put:
      summary: Update NSFW enabled state
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                enabled:
                  type: boolean
      responses:
        200:
          description: Updated successfully
        403:
          description: PIN verification required

  /api/v1/user/preferences/nsfw/pin:
    post:
      summary: Set or update NSFW PIN
      requestBody:
        content:
          application/json:
            schema:
              type: object
              required: [pin]
              properties:
                pin:
                  type: string
                  minLength: 4
                  maxLength: 6
                current_pin:
                  type: string
                  description: Required when changing existing PIN

    delete:
      summary: Remove NSFW PIN
      requestBody:
        content:
          application/json:
            schema:
              type: object
              required: [current_pin]
              properties:
                current_pin:
                  type: string

  /api/v1/user/preferences/nsfw/verify:
    post:
      summary: Verify NSFW PIN
      requestBody:
        content:
          application/json:
            schema:
              type: object
              required: [pin]
              properties:
                pin:
                  type: string
      responses:
        200:
          description: PIN verified, NSFW enabled
        401:
          description: Invalid PIN
```

---

## Behavior Summary

| State | `/qar/*` Routes | Search | Sidebar | Dashboard |
|-------|-----------------|--------|---------|-----------|
| NSFW OFF | 404 | Hidden | Hidden | Excluded |
| NSFW ON | Accessible | Included | Visible | Included |
| NSFW ON + Timeout | Auto-locks | — | — | — |
| NSFW ON + PIN | Requires PIN to enable | — | — | — |

---

## Implementation Checklist

**Location**: `internal/service/preferences/`

### Phase 1: Core Infrastructure
- [ ] Create package structure `internal/service/preferences/`
- [ ] Define entities: `UserPreferences`, `NSFWAuditEntry`
- [ ] Create repository interface `PreferencesRepository`
- [ ] Implement fx module `preferences.Module`
- [ ] Add configuration struct for NSFW settings (timeout defaults)

### Phase 2: Database
- [ ] Create migration `xxx_nsfw_toggle.up.sql`
- [ ] Create `user_preferences` table with NSFW fields
- [ ] Create `nsfw_toggle_audit` table for compliance logging
- [ ] Add index on `nsfw_toggle_audit(user_id, timestamp DESC)`
- [ ] Generate sqlc queries for preference CRUD
- [ ] Generate sqlc queries for audit log insertion

### Phase 3: Service Layer
- [ ] Implement `PreferencesService` for NSFW state management
- [ ] Implement PIN hashing with `bcrypt` (never store plain text)
- [ ] Implement PIN verification with timing-safe comparison
- [ ] Add Dragonfly cache layer for fast middleware checks (`user:{id}:nsfw`)
- [ ] Implement cache invalidation on preference changes
- [ ] Add audit logging for all toggle actions (enabled, disabled, pin_set, pin_verified)
- [ ] Implement auto-lock timeout logic with activity tracking

### Phase 4: Background Jobs
- [ ] Create `NSFWAutoLockWorker` for inactivity timeout
- [ ] Implement activity tracking via cache (`user:{id}:last_activity`)
- [ ] Configure River job snoozing for active users
- [ ] Add job scheduling on NSFW enable (when timeout configured)

### Phase 5: Middleware
- [ ] Implement `NSFWMiddleware` for `/api/v1/qar/*` route protection
- [ ] Return 404 (not 403) to obscure adult content existence
- [ ] Implement `ActivityMiddleware` for tracking user activity
- [ ] Integrate middleware into router chain

### Phase 6: API Integration
- [ ] Define OpenAPI spec for NSFW preference endpoints
- [ ] Generate ogen handlers for preferences
- [ ] Implement `GET /api/v1/user/preferences/nsfw` - get status
- [ ] Implement `PUT /api/v1/user/preferences/nsfw` - update enabled state
- [ ] Implement `POST /api/v1/user/preferences/nsfw/pin` - set/update PIN
- [ ] Implement `DELETE /api/v1/user/preferences/nsfw/pin` - remove PIN
- [ ] Implement `POST /api/v1/user/preferences/nsfw/verify` - verify PIN
- [ ] Add auth middleware to all endpoints
- [ ] Implement rate limiting on PIN verification (prevent brute force)

### Phase 7: Search Integration
- [ ] Update `SearchService` to check NSFW state before including adult results
- [ ] Filter adult modules (`AdultMovies`, `AdultScenes`) based on user preference
- [ ] Ensure search respects middleware protection

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
| [Casbin](https://pkg.go.dev/github.com/casbin/casbin/v2) | [Local](../../../sources/security/casbin.md) |
| [Dragonfly Documentation](https://www.dragonflydb.io/docs) | [Local](../../../sources/infrastructure/dragonfly.md) |
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

- [ADULT_CONTENT_SYSTEM.md](ADULT_CONTENT_SYSTEM.md) - Full adult content architecture
- [WHISPARR_STASHDB_SCHEMA.md](WHISPARR_STASHDB_SCHEMA.md) - Adult content schema and metadata
- [USER_EXPERIENCE_FEATURES.md](USER_EXPERIENCE_FEATURES.md) - General UX features

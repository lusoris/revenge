# Revenge - NSFW Toggle

> User preference component for adult content visibility.
> Referenced by [WHISPARR_STASHDB_SCHEMA.md](WHISPARR_STASHDB_SCHEMA.md) and [ADULT_CONTENT_SYSTEM.md](ADULT_CONTENT_SYSTEM.md).

## Overview

The NSFW toggle controls visibility and access to adult content modules (`adult_movie`, `adult_scene`) stored in PostgreSQL schema `c`.

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
              │  (Dragonfly)    │       │  (/c/* visibility)   │
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

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
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
// NSFWMiddleware blocks /c/* routes when NSFW is disabled
func NSFWMiddleware(sessions *SessionService) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Only apply to /c/* routes
            if !strings.HasPrefix(r.URL.Path, "/api/v1/c/") {
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
                placeholder="Enter PIN"
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
        <NavItem href="/c/movies" icon="lock">Adult Movies</NavItem>
        <NavItem href="/c/scenes" icon="lock">Adult Scenes</NavItem>
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

| State | `/c/*` Routes | Search | Sidebar | Dashboard |
|-------|---------------|--------|---------|-----------|
| NSFW OFF | 404 | Hidden | Hidden | Excluded |
| NSFW ON | Accessible | Included | Visible | Included |
| NSFW ON + Timeout | Auto-locks | — | — | — |
| NSFW ON + PIN | Requires PIN to enable | — | — | — |

---

## Related Documents

- [ADULT_CONTENT_SYSTEM.md](ADULT_CONTENT_SYSTEM.md) - Full adult content architecture
- [WHISPARR_STASHDB_SCHEMA.md](WHISPARR_STASHDB_SCHEMA.md) - Adult content schema and metadata
- [USER_EXPERIENCE_FEATURES.md](USER_EXPERIENCE_FEATURES.md) - General UX features

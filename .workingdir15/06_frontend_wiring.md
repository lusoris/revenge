# Frontend Wiring — Implementation Plan

## Overview
Wire the SvelteKit SPA to the backend's ~184 API endpoints. Build incrementally:
API client → Auth → Layout → Library → Detail → Playback → Search → Settings.

## Phase 1: API Client Layer
- **`$lib/api/client.ts`** — Typed fetch wrapper with auth header injection, token refresh, error mapping
- **`$lib/api/types.ts`** — TypeScript interfaces generated/hand-written from OpenAPI schemas
- **`$lib/api/endpoints/`** — Per-domain modules (auth, movies, tvshows, search, playback, etc.)
- Pattern: All endpoints return typed promises, use `@tanstack/svelte-query` for caching/dedup

## Phase 2: Auth Flow
- **`$lib/stores/auth.ts`** — Svelte 5 runes-based auth state ($state for user, tokens, loading)
- **`$lib/api/endpoints/auth.ts`** — login, register, logout, refresh, verify-email
- **`/login`** route — Email/password form, OIDC buttons, MFA challenge
- **`/register`** route — Registration form with validation (zod)
- **Route guard** — `+layout.ts` load function checks auth, redirects to /login
- Token refresh — Automatic via interceptor in client.ts (401 → refresh → retry)

## Phase 3: Core Layout
- **`$components/layout/Sidebar.svelte`** — Navigation: Home, Movies, TV Shows, Search, Settings
- **`$components/layout/Header.svelte`** — User avatar, search bar, theme toggle
- **`$components/layout/AppShell.svelte`** — Sidebar + header + main content area
- **`/(app)/+layout.svelte`** — Authenticated layout group wrapping all app routes
- Dark theme by default (mode-watcher), responsive sidebar (mobile: bottom nav or drawer)

## Phase 4: Library Browsing
- **`/(app)/movies/+page.svelte`** — Movie grid with infinite scroll, sort/filter
- **`/(app)/tvshows/+page.svelte`** — TV show grid
- **`/(app)/+page.svelte`** — Home/dashboard: continue watching, recently added, stats
- **`$components/media/MediaCard.svelte`** — Poster card with title, year, rating
- **`$components/media/MediaGrid.svelte`** — Responsive grid layout
- Use `@tanstack/svelte-query` infinite queries for pagination

## Phase 5: Detail Pages
- **`/(app)/movies/[id]/+page.svelte`** — Movie detail: backdrop, metadata, cast, files, play button
- **`/(app)/tvshows/[id]/+page.svelte`** — TV show detail: seasons, episodes, cast
- **`/(app)/tvshows/[id]/season/[seasonNum]/+page.svelte`** — Season episodes list
- **`$components/media/CastList.svelte`** — Horizontal scroll cast cards
- Watch progress tracking (progress bar on cards, resume button)

## Phase 6: Playback
- **`/(app)/play/[sessionId]/+page.svelte`** — Full-screen vidstack player
- **`$lib/api/endpoints/playback.ts`** — startPlaybackSession, heartbeat, stop
- HLS.js integration via vidstack (master playlist → adaptive bitrate)
- Audio/subtitle track selection, seek, quality selection
- Heartbeat interval (30s) to keep session alive + report position

## Phase 7: Search
- **`/(app)/search/+page.svelte`** — Multi-search with autocomplete
- Debounced input → autocomplete API → full search on enter
- Faceted filtering (genre, year, rating)

## Phase 8: Settings & Admin (later)
- User settings, MFA setup, API keys, session management
- Admin: users, RBAC, libraries, integrations, activity logs

## Technical Decisions
- **No OpenAPI codegen** — Hand-written types (simpler, smaller bundle)
- **`$lib/api/client.ts`** — Thin wrapper around fetch, not axios
- **Auth tokens** — Stored in memory ($state) + refresh token in httpOnly cookie (if backend supports) or localStorage
- **Image proxy** — Use backend's `/api/v1/images/{type}/{size}/{path}` for all TMDb images
- **Routing** — SvelteKit route groups: `(auth)` for login/register, `(app)` for authenticated routes

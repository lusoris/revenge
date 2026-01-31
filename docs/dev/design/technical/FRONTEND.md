# Revenge - Frontend Architecture

<!-- SOURCES: gohlslib, lastfm-api, shadcn-svelte, svelte-runes, svelte5, sveltekit, tanstack-query -->

<!-- DESIGN: technical, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


> Modern, responsive web interface with full RBAC and theme support.


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Technology Stack](#technology-stack)
- [Project Structure](#project-structure)
- [Role-Based Access Control (RBAC)](#role-based-access-control-rbac)
  - [Roles](#roles)
  - [Permissions](#permissions)
  - [Permission Check](#permission-check)
- [Theme System](#theme-system)
  - [CSS Variables](#css-variables)
  - [Theme Store](#theme-store)
  - [Custom Themes (Admin)](#custom-themes-admin)
- [API Integration](#api-integration)
  - [OpenAPI Client Generation](#openapi-client-generation)
  - [TanStack Query Integration](#tanstack-query-integration)
  - [Usage in Components](#usage-in-components)
- [Video Player](#video-player)
  - [Features](#features)
  - [Player Component](#player-component)
- [PWA Support](#pwa-support)
  - [Service Worker](#service-worker)
  - [Manifest](#manifest)
- [Internationalization (i18n)](#internationalization-i18n)
  - [Setup](#setup)
  - [Usage](#usage)
- [Responsive Design](#responsive-design)
  - [Breakpoints](#breakpoints)
  - [Mobile-First Components](#mobile-first-components)
- [Admin Panel Features](#admin-panel-features)
  - [Dashboard](#dashboard)
  - [User Management](#user-management)
  - [Library Management](#library-management)
  - [Server Settings](#server-settings)
  - [Activity Logs](#activity-logs)
- [Development](#development)
  - [Setup](#setup)
  - [Build](#build)
  - [Component Development](#component-development)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Related Documentation](#related-documentation)

<!-- TOC-END -->

## Status

| Dimension | Status |
|-----------|--------|
| Design | 🔴 |
| Sources | 🔴 |
| Instructions | 🔴 |
| Code | 🔴 |
| Linting | 🔴 |
| Unit Testing | 🔴 |
| Integration Testing | 🔴 |
---

## Technology Stack

| Component | Technology | Version | Purpose |
|-----------|------------|---------|---------|
| Framework | SvelteKit | 2.x | SSR, routing, API integration |
| Language | TypeScript | 5.x | Type safety |
| UI Framework | Tailwind CSS | 4.x | Utility-first styling |
| Components | shadcn-svelte | latest | Accessible UI components |
| State | Svelte Stores | built-in | Client-side state |
| Server State | TanStack Query | 5.x | Caching, background refresh |
| Forms | Superforms | latest | Form validation |
| Icons | Lucide | latest | Icon library |
| Charts | Chart.js | 4.x | Analytics/dashboards |

## Project Structure

```
web/
├── src/
│   ├── lib/
│   │   ├── components/
│   │   │   ├── ui/                 # shadcn-svelte components
│   │   │   │   ├── button/
│   │   │   │   ├── card/
│   │   │   │   ├── dialog/
│   │   │   │   └── ...
│   │   │   ├── media/              # Media-specific components
│   │   │   │   ├── MediaCard.svelte
│   │   │   │   ├── MediaGrid.svelte
│   │   │   │   ├── VideoPlayer.svelte
│   │   │   │   ├── AudioPlayer.svelte
│   │   │   │   └── ImageGallery.svelte
│   │   │   ├── admin/              # Admin panel components
│   │   │   │   ├── UserTable.svelte
│   │   │   │   ├── LibraryManager.svelte
│   │   │   │   └── ActivityLog.svelte
│   │   │   └── layout/             # Layout components
│   │   │       ├── Navbar.svelte
│   │   │       ├── Sidebar.svelte
│   │   │       └── Footer.svelte
│   │   ├── stores/
│   │   │   ├── auth.ts             # Authentication state
│   │   │   ├── user.ts             # User preferences
│   │   │   ├── theme.ts            # Theme state
│   │   │   └── playback.ts         # Playback state
│   │   ├── api/
│   │   │   ├── client.ts           # API client (from OpenAPI)
│   │   │   ├── movies.ts           # Movie API
│   │   │   ├── shows.ts            # TV show API
│   │   │   └── ...
│   │   └── utils/
│   │       ├── format.ts           # Formatters (date, duration)
│   │       ├── auth.ts             # Auth helpers
│   │       └── media.ts            # Media helpers
│   ├── routes/
│   │   ├── (app)/                  # Main app (authenticated)
│   │   │   ├── +layout.svelte
│   │   │   ├── +page.svelte        # Home/dashboard
│   │   │   ├── (admin)/            # Admin routes
│   │   │   │   ├── +layout.svelte  # Admin layout
│   │   │   │   ├── users/
│   │   │   │   ├── libraries/
│   │   │   │   ├── settings/
│   │   │   │   ├── logs/
│   │   │   │   └── tasks/
│   │   │   ├── (media)/            # Media browsing
│   │   │   │   ├── movies/
│   │   │   │   │   ├── +page.svelte        # Movie list
│   │   │   │   │   └── [id]/+page.svelte   # Movie detail
│   │   │   │   ├── shows/
│   │   │   │   ├── music/
│   │   │   │   ├── audiobooks/
│   │   │   │   ├── photos/
│   │   │   │   └── search/
│   │   │   ├── (player)/           # Full-screen players
│   │   │   │   ├── video/[id]/
│   │   │   │   └── audio/[id]/
│   │   │   └── profile/            # User profile
│   │   ├── (auth)/                 # Authentication routes
│   │   │   ├── login/
│   │   │   ├── register/
│   │   │   ├── forgot-password/
│   │   │   └── oidc/
│   │   │       └── callback/
│   │   └── api/                    # SvelteKit API routes (BFF)
│   │       └── [...path]/
│   ├── app.css                     # Global styles + Tailwind
│   ├── app.html                    # HTML template
│   └── hooks.server.ts             # Server hooks (auth)
├── static/
│   ├── favicon.ico
│   └── manifest.json               # PWA manifest
├── tailwind.config.ts
├── svelte.config.js
├── vite.config.ts
├── tsconfig.json
└── package.json
```

## Role-Based Access Control (RBAC)

### Roles

| Role | Level | Description |
|------|-------|-------------|
| `admin` | 100 | Full system access |
| `moderator` | 50 | Content and library management |
| `user` | 10 | Standard user |
| `guest` | 0 | Read-only (if enabled) |

### Permissions

```typescript
const permissions = {
  // User management
  'users:read': ['admin', 'moderator'],
  'users:create': ['admin'],
  'users:update': ['admin'],
  'users:delete': ['admin'],

  // Library management
  'libraries:read': ['admin', 'moderator', 'user'],
  'libraries:create': ['admin'],
  'libraries:update': ['admin', 'moderator'],
  'libraries:delete': ['admin'],
  'libraries:scan': ['admin', 'moderator'],

  // Content
  'content:read': ['admin', 'moderator', 'user', 'guest'],
  'content:update_metadata': ['admin', 'moderator'],
  'content:delete': ['admin'],

  // Playback
  'playback:stream': ['admin', 'moderator', 'user'],
  'playback:download': ['admin', 'moderator', 'user'], // if enabled

  // Server settings
  'settings:read': ['admin'],
  'settings:update': ['admin'],

  // Activity logs
  'logs:read': ['admin', 'moderator'],

  // Adult content (requires explicit scope)
  'adult:read': ['admin', 'user'], // only if adult_enabled
}
```

### Permission Check

```svelte
<script>
  import { hasPermission } from '$lib/stores/auth';
</script>

{#if $hasPermission('users:create')}
  <Button on:click={createUser}>Create User</Button>
{/if}
```

## Theme System

### CSS Variables

```css
/* app.css */
:root {
  /* Colors */
  --color-background: 0 0% 100%;
  --color-foreground: 222.2 84% 4.9%;
  --color-primary: 222.2 47.4% 11.2%;
  --color-primary-foreground: 210 40% 98%;
  --color-secondary: 210 40% 96.1%;
  --color-accent: 210 40% 96.1%;
  --color-destructive: 0 84.2% 60.2%;

  /* Spacing */
  --radius: 0.5rem;

  /* Media player */
  --player-background: 0 0% 0%;
  --player-controls: 0 0% 100%;
}

.dark {
  --color-background: 222.2 84% 4.9%;
  --color-foreground: 210 40% 98%;
  /* ... dark theme values */
}
```

### Theme Store

```typescript
// lib/stores/theme.ts
import { writable } from 'svelte/store';
import { browser } from '$app/environment';

type Theme = 'light' | 'dark' | 'system';

function createThemeStore() {
  const stored = browser ? localStorage.getItem('theme') as Theme : 'system';
  const { subscribe, set } = writable<Theme>(stored || 'system');

  return {
    subscribe,
    set: (theme: Theme) => {
      if (browser) {
        localStorage.setItem('theme', theme);
        applyTheme(theme);
      }
      set(theme);
    }
  };
}

export const theme = createThemeStore();
```

### Custom Themes (Admin)

Admins can create custom themes via settings:

```typescript
interface CustomTheme {
  id: string;
  name: string;
  base: 'light' | 'dark';
  colors: {
    primary: string;
    secondary: string;
    accent: string;
    background: string;
    foreground: string;
  };
  isDefault: boolean;
}
```

## API Integration

### OpenAPI Client Generation

```bash
# Generate TypeScript client from OpenAPI spec
npx openapi-typescript-codegen \
  --input ../api/openapi/revenge.yaml \
  --output src/lib/api/generated \
  --client fetch
```

### TanStack Query Integration

```typescript
// lib/api/movies.ts
import { createQuery, createMutation } from '@tanstack/svelte-query';
import { api } from './client';

export function useMovies(libraryId?: string) {
  return createQuery({
    queryKey: ['movies', libraryId],
    queryFn: () => api.movies.list({ libraryId }),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

export function useMovie(id: string) {
  return createQuery({
    queryKey: ['movies', id],
    queryFn: () => api.movies.get(id),
  });
}

export function useUpdateMovie() {
  return createMutation({
    mutationFn: (data: UpdateMovieRequest) => api.movies.update(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['movies'] });
    },
  });
}
```

### Usage in Components

```svelte
<script lang="ts">
  import { useMovies } from '$lib/api/movies';

  const movies = useMovies();
</script>

{#if $movies.isLoading}
  <Spinner />
{:else if $movies.error}
  <Error message={$movies.error.message} />
{:else}
  <MediaGrid items={$movies.data} />
{/if}
```

## Video Player

### Features

- HLS.js for adaptive streaming
- Keyboard shortcuts (space, arrows, f, m, etc.)
- Picture-in-Picture support
- Chromecast support
- Subtitle selection and styling
- Audio track selection
- Quality selection (auto + manual)
- Progress persistence
- Chapter markers
- Skip intro/outro
- Binge mode (auto-play next)

### Player Component

```svelte
<!-- lib/components/media/VideoPlayer.svelte -->
<script lang="ts">
  import Hls from 'hls.js';
  import { onMount, onDestroy } from 'svelte';
  import { playbackStore } from '$lib/stores/playback';

  export let src: string;
  export let poster: string;
  export let startPosition: number = 0;

  let video: HTMLVideoElement;
  let hls: Hls;

  onMount(() => {
    if (Hls.isSupported()) {
      hls = new Hls({
        enableWorker: true,
        lowLatencyMode: false,
      });
      hls.loadSource(src);
      hls.attachMedia(video);
      hls.on(Hls.Events.MANIFEST_PARSED, () => {
        video.currentTime = startPosition;
        video.play();
      });
    } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
      // Safari native HLS
      video.src = src;
    }
  });

  onDestroy(() => {
    hls?.destroy();
  });
</script>

<div class="player-container">
  <video
    bind:this={video}
    {poster}
    playsinline
    on:timeupdate={handleTimeUpdate}
    on:ended={handleEnded}
  />
  <PlayerControls {video} />
</div>
```

## PWA Support

### Service Worker

```typescript
// service-worker.ts
import { build, files, version } from '$service-worker';

const CACHE_NAME = `revenge-${version}`;

self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => {
      return cache.addAll([...build, ...files]);
    })
  );
});

self.addEventListener('fetch', (event) => {
  // Network-first for API, cache-first for static
  if (event.request.url.includes('/api/')) {
    event.respondWith(networkFirst(event.request));
  } else {
    event.respondWith(cacheFirst(event.request));
  }
});
```

### Manifest

```json
{
  "name": "Revenge Media Server",
  "short_name": "Revenge",
  "start_url": "/",
  "display": "standalone",
  "background_color": "#000000",
  "theme_color": "#7c3aed",
  "icons": [
    { "src": "/icon-192.png", "sizes": "192x192", "type": "image/png" },
    { "src": "/icon-512.png", "sizes": "512x512", "type": "image/png" }
  ]
}
```

## Internationalization (i18n)

### Setup

```typescript
// lib/i18n/index.ts
import { init, register, getLocaleFromNavigator } from 'svelte-i18n';

register('en', () => import('./locales/en.json'));
register('de', () => import('./locales/de.json'));
register('fr', () => import('./locales/fr.json'));
register('es', () => import('./locales/es.json'));

init({
  fallbackLocale: 'en',
  initialLocale: getLocaleFromNavigator(),
});
```

### Usage

```svelte
<script>
  import { t } from 'svelte-i18n';
</script>

<h1>{$t('movies.title')}</h1>
<Button>{$t('common.save')}</Button>
```

## Responsive Design

### Breakpoints

```typescript
// Tailwind breakpoints
const screens = {
  'sm': '640px',   // Mobile landscape
  'md': '768px',   // Tablet
  'lg': '1024px',  // Desktop
  'xl': '1280px',  // Large desktop
  '2xl': '1536px', // Extra large
};
```

### Mobile-First Components

```svelte
<!-- Responsive media grid -->
<div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
  {#each items as item}
    <MediaCard {item} />
  {/each}
</div>

<!-- Responsive sidebar -->
<aside class="
  fixed inset-y-0 left-0 z-50 w-64
  -translate-x-full lg:translate-x-0
  transition-transform duration-300
">
  <Sidebar />
</aside>
```

## Admin Panel Features

### Dashboard

- Active sessions count
- Library statistics
- Recent activity
- System health (CPU, RAM, disk)
- Transcoding queue status

### User Management

- User list with search/filter
- Create/edit/delete users
- Role assignment
- Password reset
- Session management (force logout)
- Activity history per user

### Library Management

- Add/remove libraries
- Configure library settings
- Manual/scheduled scanning
- Metadata refresh
- Library permissions

### Server Settings

- General settings (server name, URL)
- Authentication settings (OIDC, registration)
- Transcoding settings (Blackbeard URL)
- Module enable/disable
- Adult content toggle
- Default user preferences

### Activity Logs

- Filterable activity log
- User actions
- Playback history
- Login history
- Admin actions

## Development

### Setup

```bash
cd web
pnpm install
pnpm dev
```

### Build

```bash
pnpm build
pnpm preview
```

### Component Development

```bash
# Add new shadcn-svelte component
npx shadcn-svelte@latest add button
npx shadcn-svelte@latest add dialog
```

---


## Related Documentation

| Document | Description |
|----------|-------------|
| [I18N.md](I18N.md) | Multi-language support (UI, metadata, audio/subtitle) |
| [AUDIO_STREAMING.md](AUDIO_STREAMING.md) | Audio player integration, progress tracking |
| [SCROBBLING.md](SCROBBLING.md) | External service connections (Trakt, Last.fm) |
| [ARCHITECTURE.md](../architecture/01_ARCHITECTURE.md) | Backend architecture overview |

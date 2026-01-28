---
applyTo: "**/web/**/*.svelte,**/web/**/*.ts,**/web/src/**/*"
alwaysApply: false
---

# Frontend Architecture Instructions

> SvelteKit 2 + Tailwind CSS 4 + shadcn-svelte WebUI.

## Stack

| Component  | Technology                     | Version |
| ---------- | ------------------------------ | ------- |
| Framework  | SvelteKit                      | 2.x     |
| UI Library | shadcn-svelte                  | Latest  |
| Styling    | Tailwind CSS                   | 4.x     |
| State      | Svelte Stores + TanStack Query | -       |
| API Client | OpenAPI Generated              | -       |
| Auth       | JWT + OIDC                     | -       |
| i18n       | Built-in                       | -       |
| PWA        | Service Worker                 | -       |

## Project Structure

```
web/
  src/
    lib/
      components/
        ui/              # shadcn-svelte components
        media/           # Media cards, grids, players
        admin/           # Admin panel components
      stores/            # Svelte stores (auth, theme, playback)
      api/               # Generated API client from OpenAPI
      utils/             # Helpers, formatters
    routes/
      (app)/             # Main app layout
        (admin)/         # Admin routes (/admin/...)
          users/
          libraries/
          settings/
        (media)/         # Media browsing
          movies/
          shows/
          music/
        (player)/        # Player overlay/page
      (auth)/            # Auth routes (login, register, oidc)
      api/               # API routes (BFF pattern if needed)
    app.css
    app.html
  static/
  tailwind.config.ts
  svelte.config.js
```

## Design Principles

### 1. Performance First

- Code splitting per route
- Lazy load heavy components (player, visualizer)
- Prefetch data on hover (media cards)
- Virtual scrolling for large lists
- Image lazy loading with blurhash placeholders

### 2. Accessibility

- Semantic HTML (`<nav>`, `<main>`, `<article>`)
- ARIA labels on all interactive elements
- Keyboard navigation (Tab, Arrow keys, Space, Enter)
- Focus visible indicators
- Screen reader friendly

### 3. Responsive Design

- Mobile-first approach
- Breakpoints: `sm: 640px`, `md: 768px`, `lg: 1024px`, `xl: 1280px`, `2xl: 1536px`
- Touch-friendly targets (min 44x44px)
- Swipe gestures on mobile

### 4. Offline Support

- Service Worker for static assets
- IndexedDB for cached media metadata
- Offline indicator in UI
- Retry failed requests automatically

## Code Patterns

### DO

- ✅ Use `$:` reactive statements for derived state
- ✅ Use `onMount()` for side effects, cleanup in `onDestroy()`
- ✅ Use `bind:this` for DOM references
- ✅ Use TypeScript for all components
- ✅ Use TanStack Query for server state
- ✅ Use Svelte stores for client state
- ✅ Prefetch routes with `data-sveltekit-preload-data`
- ✅ Use `progressive enhancement` (works without JS)

### DON'T

- ❌ Use `window` or `document` outside `onMount()`
- ❌ Mutate props directly (use events)
- ❌ Use `any` type in TypeScript
- ❌ Fetch data in components (use load functions)
- ❌ Use inline styles (use Tailwind classes)
- ❌ Hardcode text (use i18n)
- ❌ Use global CSS (scope with Svelte's `<style>`)

## Routing

### Route Groups

```
(app)/        # Main layout with nav, sidebar
(admin)/      # Admin layout with admin nav
(auth)/       # Centered layout, no nav
(player)/     # Fullscreen player layout
```

### Load Functions

```typescript
// +page.ts - runs on both server and client
export async function load({ fetch, params }) {
  const movie = await fetch(`/api/movies/${params.id}`).then((r) => r.json());
  return { movie };
}

// +page.server.ts - runs only on server
export async function load({ locals, params }) {
  const user = locals.user; // From handle hook
  // Sensitive data here
}
```

### Preloading

```svelte
<!-- Prefetch on hover -->
<a href="/movies/123" data-sveltekit-preload-data="hover">
    Watch Movie
</a>

<!-- Prefetch on viewport -->
<a href="/shows/456" data-sveltekit-preload-data="viewport">
    Browse Shows
</a>
```

## State Management

### Client State (Svelte Stores)

```typescript
// lib/stores/auth.ts
import { writable } from "svelte/store";

interface User {
  id: string;
  username: string;
  role: "admin" | "user";
}

export const user = writable<User | null>(null);
export const isAuthenticated = derived(user, ($user) => $user !== null);
```

### Server State (TanStack Query)

```typescript
// lib/api/movies.ts
import { createQuery } from "@tanstack/svelte-query";

export function movieQuery(id: string) {
  return createQuery({
    queryKey: ["movie", id],
    queryFn: () => fetch(`/api/movies/${id}`).then((r) => r.json()),
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

// In component
const movieQuery = movieQuery(movieId);
$: movie = $movieQuery.data;
$: loading = $movieQuery.isLoading;
```

### Playback State (Global Store)

```typescript
// lib/stores/playback.ts
import { writable } from "svelte/store";

interface PlaybackState {
  mediaId: string | null;
  playing: boolean;
  position: number;
  duration: number;
  volume: number;
  muted: boolean;
  quality: string;
  subtitles: string | null;
}

export const playback = writable<PlaybackState>({
  mediaId: null,
  playing: false,
  position: 0,
  duration: 0,
  volume: 1.0,
  muted: false,
  quality: "auto",
  subtitles: null,
});
```

## Component Patterns

### Media Card

```svelte
<!-- lib/components/media/MediaCard.svelte -->
<script lang="ts">
    import { fade } from 'svelte/transition';
    import { playback } from '$lib/stores/playback';

    export let media: Media;
    export let showProgress = false;

    let imageLoaded = false;

    function handlePlay() {
        playback.set({ mediaId: media.id, playing: true, ... });
    }
</script>

<article class="media-card" transition:fade>
    <a href="/movies/{media.id}" data-sveltekit-preload-data="hover">
        <!-- Blurhash placeholder -->
        {#if !imageLoaded}
            <div class="blurhash" style="background: {media.blurhash}" />
        {/if}

        <!-- Image -->
        <img
            src={media.poster}
            alt={media.title}
            loading="lazy"
            on:load={() => imageLoaded = true}
        />

        <!-- Progress bar (if watched) -->
        {#if showProgress && media.progress > 0}
            <div class="progress">
                <div class="progress-bar" style="width: {media.progress}%" />
            </div>
        {/if}
    </a>

    <!-- Quick actions -->
    <button on:click={handlePlay} aria-label="Play {media.title}">
        <PlayIcon />
    </button>
</article>

<style>
    .media-card {
        position: relative;
        aspect-ratio: 2/3;
        border-radius: 0.5rem;
        overflow: hidden;
        transition: transform 0.2s;
    }

    .media-card:hover {
        transform: scale(1.05);
    }
</style>
```

### Virtual Scroller (Large Lists)

```svelte
<!-- lib/components/media/VirtualGrid.svelte -->
<script lang="ts">
    import { onMount, onDestroy } from 'svelte';

    export let items: any[];
    export let itemHeight = 300;
    export let itemWidth = 200;

    let container: HTMLElement;
    let visibleRange = { start: 0, end: 20 };

    function updateVisibleRange() {
        const scrollTop = container.scrollTop;
        const containerHeight = container.clientHeight;

        const start = Math.floor(scrollTop / itemHeight);
        const end = Math.ceil((scrollTop + containerHeight) / itemHeight);

        visibleRange = { start, end };
    }

    onMount(() => {
        container.addEventListener('scroll', updateVisibleRange);
    });

    onDestroy(() => {
        container?.removeEventListener('scroll', updateVisibleRange);
    });

    $: visibleItems = items.slice(visibleRange.start, visibleRange.end);
</script>

<div bind:this={container} class="virtual-scroller">
    <div class="spacer" style="height: {items.length * itemHeight}px">
        <div class="items" style="transform: translateY({visibleRange.start * itemHeight}px)">
            {#each visibleItems as item (item.id)}
                <MediaCard {item} />
            {/each}
        </div>
    </div>
</div>
```

## Theming

### Theme Store

```typescript
// lib/stores/theme.ts
import { writable } from "svelte/store";
import { browser } from "$app/environment";

type Theme = "light" | "dark" | "system";

function createThemeStore() {
  const { subscribe, set } = writable<Theme>("system");

  return {
    subscribe,
    set: (theme: Theme) => {
      if (browser) {
        localStorage.setItem("theme", theme);
        applyTheme(theme);
      }
      set(theme);
    },
    init: () => {
      if (browser) {
        const saved = localStorage.getItem("theme") as Theme;
        const theme = saved || "system";
        applyTheme(theme);
        set(theme);
      }
    },
  };
}

function applyTheme(theme: Theme) {
  const isDark =
    theme === "dark" ||
    (theme === "system" &&
      window.matchMedia("(prefers-color-scheme: dark)").matches);

  document.documentElement.classList.toggle("dark", isDark);
}

export const theme = createThemeStore();
```

### CSS Variables

```css
/* app.css */
@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  :root {
    /* Light mode */
    --background: 0 0% 100%;
    --foreground: 222.2 84% 4.9%;
    --primary: 221.2 83.2% 53.3%;
    --accent: 210 40% 96.1%;
  }

  .dark {
    /* Dark mode */
    --background: 222.2 84% 4.9%;
    --foreground: 210 40% 98%;
    --primary: 217.2 91.2% 59.8%;
    --accent: 217.2 32.6% 17.5%;
  }
}
```

## API Integration

### Generated Client

```typescript
// lib/api/client.ts
import { Configuration, DefaultApi } from "./generated";

const config = new Configuration({
  basePath: import.meta.env.VITE_API_URL || "http://localhost:8096",
  accessToken: () => localStorage.getItem("token") || "",
});

export const api = new DefaultApi(config);
```

### TanStack Query Integration

```typescript
// lib/api/movies.ts
import { createQuery, createMutation } from "@tanstack/svelte-query";
import { api } from "./client";

export function moviesQuery(page = 1, limit = 20) {
  return createQuery({
    queryKey: ["movies", page, limit],
    queryFn: () => api.getMovies({ page, limit }),
  });
}

export function rateMovieMutation() {
  return createMutation({
    mutationFn: ({ movieId, rating }: { movieId: string; rating: number }) =>
      api.rateMovie({ movieId, rating }),
    onSuccess: () => {
      // Invalidate movies query to refetch
      queryClient.invalidateQueries({ queryKey: ["movies"] });
    },
  });
}
```

## Authentication

### Auth Store

```typescript
// lib/stores/auth.ts
import { writable, derived } from "svelte/store";
import { goto } from "$app/navigation";
import { api } from "$lib/api/client";

function createAuthStore() {
  const { subscribe, set } = writable<User | null>(null);

  return {
    subscribe,
    login: async (username: string, password: string) => {
      const response = await api.login({ username, password });
      localStorage.setItem("token", response.token);
      set(response.user);
      goto("/");
    },
    logout: () => {
      localStorage.removeItem("token");
      set(null);
      goto("/login");
    },
    init: async () => {
      const token = localStorage.getItem("token");
      if (token) {
        try {
          const user = await api.getMe();
          set(user);
        } catch (e) {
          localStorage.removeItem("token");
        }
      }
    },
  };
}

export const auth = createAuthStore();
export const isAdmin = derived(auth, ($auth) => $auth?.role === "admin");
```

### Protected Routes

```typescript
// hooks.server.ts
export async function handle({ event, resolve }) {
  const token = event.cookies.get("token");

  if (token) {
    try {
      const user = await verifyToken(token);
      event.locals.user = user;
    } catch (e) {
      // Invalid token
    }
  }

  // Protect admin routes
  if (event.url.pathname.startsWith("/admin") && !event.locals.user?.isAdmin) {
    return new Response("Unauthorized", { status: 401 });
  }

  return resolve(event);
}
```

## i18n

### Translation Files

```typescript
// lib/i18n/translations/en.json
{
    "nav": {
        "home": "Home",
        "movies": "Movies",
        "shows": "TV Shows",
        "music": "Music"
    },
    "player": {
        "play": "Play",
        "pause": "Pause",
        "volume": "Volume"
    }
}

// lib/i18n/index.ts
import { derived, writable } from 'svelte/store';
import en from './translations/en.json';
import de from './translations/de.json';

const translations = { en, de };

export const locale = writable('en');
export const t = derived(locale, $locale => (key: string) => {
    const keys = key.split('.');
    let value: any = translations[$locale];

    for (const k of keys) {
        value = value[k];
    }

    return value || key;
});
```

### Usage

```svelte
<script>
    import { t } from '$lib/i18n';
</script>

<button>{$t('player.play')}</button>
```

## PWA

### Service Worker

```typescript
// service-worker.ts
const CACHE_NAME = "revenge-v1";
const STATIC_ASSETS = [
  "/",
  "/app.css",
  "/app.js",
  // ... fonts, icons, etc.
];

self.addEventListener("install", (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => cache.addAll(STATIC_ASSETS)),
  );
});

self.addEventListener("fetch", (event) => {
  event.respondWith(
    caches.match(event.request).then((response) => {
      return response || fetch(event.request);
    }),
  );
});
```

### Manifest

```json
// static/manifest.json
{
  "name": "Revenge Media Server",
  "short_name": "Revenge",
  "start_url": "/",
  "display": "standalone",
  "background_color": "#000000",
  "theme_color": "#3b82f6",
  "icons": [
    {
      "src": "/icon-192.png",
      "sizes": "192x192",
      "type": "image/png"
    },
    {
      "src": "/icon-512.png",
      "sizes": "512x512",
      "type": "image/png"
    }
  ]
}
```

## Testing

### Unit Tests (Vitest)

```typescript
// lib/utils/format.test.ts
import { describe, it, expect } from "vitest";
import { formatDuration } from "./format";

describe("formatDuration", () => {
  it("formats seconds correctly", () => {
    expect(formatDuration(90)).toBe("1:30");
    expect(formatDuration(3661)).toBe("1:01:01");
  });
});
```

### Component Tests (Testing Library)

```typescript
// lib/components/media/MediaCard.test.ts
import { render, fireEvent } from "@testing-library/svelte";
import MediaCard from "./MediaCard.svelte";

describe("MediaCard", () => {
  it("renders media title", () => {
    const { getByText } = render(MediaCard, {
      props: { media: { id: "1", title: "Test Movie" } },
    });

    expect(getByText("Test Movie")).toBeInTheDocument();
  });

  it("calls play on button click", async () => {
    const { getByLabelText } = render(MediaCard, {
      props: { media: { id: "1", title: "Test" } },
    });

    const playButton = getByLabelText("Play Test");
    await fireEvent.click(playButton);

    // Assert playback store updated
  });
});
```

## Performance Optimization

### Code Splitting

```typescript
// Lazy load heavy components
const Player = () => import('$lib/components/Player.svelte');
const Visualizer = () => import('$lib/components/Visualizer.svelte');

// In component
{#await Player() then Component}
    <Component.default />
{/await}
```

### Image Optimization

```svelte
<!-- Use blurhash placeholders -->
<img
    src={`/api/images/${media.posterId}?w=300&h=450`}
    alt={media.title}
    loading="lazy"
    decoding="async"
    style="background: {media.blurhash}"
/>
```

### Prefetching

```typescript
// Prefetch on hover
function handleMouseEnter(mediaId: string) {
  queryClient.prefetchQuery({
    queryKey: ["movie", mediaId],
    queryFn: () => api.getMovie(mediaId),
  });
}
```

## Summary

```yaml
Frontend Stack:
  Framework: SvelteKit 2
  UI: shadcn-svelte + Tailwind CSS 4
  State: Svelte Stores (client) + TanStack Query (server)
  Auth: JWT + OIDC
  i18n: Built-in translation system
  PWA: Service Worker + Manifest

  Key Features:
    - Code splitting per route
    - Virtual scrolling for large lists
    - Lazy image loading with blurhash
    - Theme system (light/dark/system)
    - Prefetching on hover
    - Offline support
    - Accessibility (ARIA, keyboard nav)
```

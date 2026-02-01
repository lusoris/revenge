# TODO v0.3.0 - MVP (Movies)

<!-- DESIGN: planning, README, SCAFFOLD_TEMPLATE, test_output_claude -->

> Movie Module + Basic Frontend

**Status**: 🔴 Not Started
**Tag**: `v0.3.0`
**Focus**: Movie Module + Full Backend + Basic UI

**Depends On**: [v0.2.0](TODO_v0.2.0.md) (Core services required)

---

## Overview

This is the **MVP milestone**. It delivers a fully functional movie library with metadata, search, Radarr integration, and a basic but usable web frontend. Users can browse, search, and play movies.

---

## Deliverables

### Movie Module (Backend)

- [ ] **Database Schema** (`migrations/`)
  - [ ] `public.movies` table
    - [ ] id (UUID v7)
    - [ ] title, original_title
    - [ ] year, release_date
    - [ ] runtime_minutes
    - [ ] overview, tagline
    - [ ] status (released, announced, etc.)
    - [ ] poster_path, backdrop_path
    - [ ] tmdb_id, imdb_id
    - [ ] added_at, updated_at
  - [ ] `public.movie_genres` table
  - [ ] `public.movie_cast` table
  - [ ] `public.movie_crew` table
  - [ ] `public.movie_files` table
  - [ ] `public.movie_watch_progress` table
  - [ ] Indexes on tmdb_id, imdb_id, title

- [ ] **Entity** (`internal/content/movie/entity.go`)
  - [ ] Movie struct
  - [ ] MovieFile struct
  - [ ] MovieGenre struct
  - [ ] MovieCast struct
  - [ ] MovieCrew struct
  - [ ] WatchProgress struct

- [ ] **Repository** (`internal/content/movie/`)
  - [ ] `repository.go` - Interface
  - [ ] `repository_pg.go` - PostgreSQL implementation
  - [ ] CRUD operations
  - [ ] List with filters (genre, year, etc.)
  - [ ] Search by title
  - [ ] Watch progress operations

- [ ] **Service** (`internal/content/movie/service.go`)
  - [ ] Get movie by ID
  - [ ] List movies (paginated)
  - [ ] Search movies
  - [ ] Update watch progress
  - [ ] Get continue watching
  - [ ] Get recently added
  - [ ] Trigger metadata refresh

- [ ] **Library Provider** (`internal/content/movie/library_service.go`)
  - [ ] Implement LibraryProvider interface
  - [ ] Scan library path
  - [ ] Match files to movies
  - [ ] Handle file changes (add/remove/modify)

- [ ] **Handler** (`internal/api/movie_handler.go`)
  - [ ] `GET /api/v1/movies` (list, paginated)
  - [ ] `GET /api/v1/movies/:id`
  - [ ] `GET /api/v1/movies/:id/files`
  - [ ] `GET /api/v1/movies/:id/cast`
  - [ ] `GET /api/v1/movies/:id/crew`
  - [ ] `GET /api/v1/movies/:id/similar`
  - [ ] `POST /api/v1/movies/:id/progress`
  - [ ] `GET /api/v1/movies/:id/progress`
  - [ ] `DELETE /api/v1/movies/:id/progress`
  - [ ] `POST /api/v1/movies/:id/refresh`

- [ ] **River Jobs** (`internal/content/movie/jobs.go`)
  - [ ] MovieMetadataRefreshJob
  - [ ] MovieLibraryScanJob
  - [ ] MovieFileMatchJob

- [ ] **fx Module** (`internal/content/movie/module.go`)
  - [ ] Provide repository
  - [ ] Provide service
  - [ ] Provide library provider
  - [ ] Register River workers

- [ ] **Tests**
  - [ ] Unit tests (80%+ coverage)
  - [ ] Integration tests

### Collection Support

- [ ] **Database Schema** (`migrations/`)
  - [ ] `public.collections` table
  - [ ] `public.collection_movies` table

- [ ] **Service** (`internal/content/movie/collection_service.go`)
  - [ ] Get collection by ID
  - [ ] List collections
  - [ ] Get movies in collection
  - [ ] Auto-detect collections from TMDb

- [ ] **Handler** (`internal/api/collection_handler.go`)
  - [ ] `GET /api/v1/collections`
  - [ ] `GET /api/v1/collections/:id`
  - [ ] `GET /api/v1/collections/:id/movies`

### Metadata Service (TMDb)

- [ ] **TMDb Client** (`internal/service/metadata/tmdb/client.go`)
  - [ ] API key configuration
  - [ ] Rate limiting (50 req/s)
  - [ ] Retry with backoff
  - [ ] Response caching

- [ ] **TMDb Service** (`internal/service/metadata/tmdb/service.go`)
  - [ ] Search movie
  - [ ] Get movie details
  - [ ] Get movie credits (cast/crew)
  - [ ] Get movie images
  - [ ] Get similar movies
  - [ ] Get collection details

- [ ] **Image Handler** (`internal/service/metadata/tmdb/images.go`)
  - [ ] Poster download/cache
  - [ ] Backdrop download/cache
  - [ ] Profile image download/cache
  - [ ] Image proxy endpoint

- [ ] **Handler** (`internal/api/metadata_handler.go`)
  - [ ] `GET /api/v1/metadata/search/movie?q=`
  - [ ] `GET /api/v1/metadata/movie/:tmdbId`
  - [ ] `GET /api/v1/images/:type/:path` (proxy)

- [ ] **Tests**
  - [ ] Unit tests with mock API
  - [ ] Integration tests (optional, needs API key)

### Search Service (Typesense)

- [ ] **Typesense Setup** (`internal/service/search/typesense.go`)
  - [ ] Client configuration
  - [ ] Collection schemas
  - [ ] Index management

- [ ] **Movie Collection Schema**
  ```json
  {
    "name": "movies",
    "fields": [
      {"name": "id", "type": "string"},
      {"name": "title", "type": "string"},
      {"name": "original_title", "type": "string"},
      {"name": "overview", "type": "string"},
      {"name": "year", "type": "int32"},
      {"name": "genres", "type": "string[]"},
      {"name": "cast", "type": "string[]"},
      {"name": "director", "type": "string"},
      {"name": "rating", "type": "float"},
      {"name": "added_at", "type": "int64"}
    ]
  }
  ```

- [ ] **Search Service** (`internal/service/search/service.go`)
  - [ ] Index movie
  - [ ] Remove from index
  - [ ] Search movies (full-text)
  - [ ] Faceted search (genre, year)
  - [ ] Autocomplete

- [ ] **Handler** (`internal/api/search_handler.go`)
  - [ ] `GET /api/v1/search?q=&type=movie`
  - [ ] `GET /api/v1/search/autocomplete?q=`

- [ ] **River Jobs** (`internal/service/search/jobs.go`)
  - [ ] SearchIndexJob - Index single item
  - [ ] SearchReindexJob - Full reindex

- [ ] **Tests**
  - [ ] Unit tests
  - [ ] Integration tests with Typesense container

### Radarr Integration

- [ ] **Radarr Client** (`internal/service/metadata/radarr/client.go`)
  - [ ] API v3 implementation
  - [ ] Authentication (API key)
  - [ ] Error handling

- [ ] **Radarr Service** (`internal/service/metadata/radarr/service.go`)
  - [ ] Get all movies
  - [ ] Get movie by ID
  - [ ] Get movie files
  - [ ] Sync library (Radarr → Revenge)
  - [ ] Trigger refresh in Radarr
  - [ ] Get quality profiles
  - [ ] Get root folders

- [ ] **Sync Logic** (`internal/service/metadata/radarr/sync.go`)
  - [ ] Full sync (initial)
  - [ ] Incremental sync (changes only)
  - [ ] File path mapping
  - [ ] Conflict resolution

- [ ] **Webhook Handler** (`internal/api/webhook_handler.go`)
  - [ ] `POST /api/v1/webhooks/radarr`
  - [ ] Handle: Grab, Download, Rename, Delete events

- [ ] **Handler** (`internal/api/radarr_handler.go`)
  - [ ] `GET /api/v1/admin/integrations/radarr/status`
  - [ ] `POST /api/v1/admin/integrations/radarr/sync`
  - [ ] `GET /api/v1/admin/integrations/radarr/quality-profiles`

- [ ] **River Jobs** (`internal/service/metadata/radarr/jobs.go`)
  - [ ] RadarrSyncJob - Full library sync
  - [ ] RadarrWebhookJob - Process webhook events

- [ ] **Tests**
  - [ ] Unit tests with mock API
  - [ ] Integration tests (optional)

### Frontend (Basic SvelteKit)

- [ ] **Project Setup** (`frontend/`)
  - [ ] SvelteKit 2 initialization
  - [ ] Svelte 5 configuration
  - [ ] TypeScript setup
  - [ ] Tailwind CSS 4 setup
  - [ ] shadcn-svelte components

- [ ] **Authentication Flow**
  - [ ] Login page (`/login`)
  - [ ] Registration page (`/register`)
  - [ ] Password reset flow
  - [ ] JWT storage (httpOnly cookie)
  - [ ] Auth store (Svelte store)
  - [ ] Protected routes

- [ ] **Layout** (`frontend/src/routes/+layout.svelte`)
  - [ ] Navigation sidebar
  - [ ] Header with user menu
  - [ ] Responsive design
  - [ ] Dark mode (default)

- [ ] **Library Browser**
  - [ ] Movies grid view (`/movies`)
  - [ ] Movie card component
  - [ ] Sorting (title, year, added)
  - [ ] Filtering (genre, year)
  - [ ] Pagination/infinite scroll
  - [ ] Search integration

- [ ] **Movie Detail Page** (`/movies/[id]`)
  - [ ] Hero backdrop
  - [ ] Poster image
  - [ ] Title, year, runtime
  - [ ] Overview
  - [ ] Cast carousel
  - [ ] Crew list
  - [ ] Similar movies
  - [ ] Play button
  - [ ] Watch progress

- [ ] **Search** (`/search`)
  - [ ] Global search bar
  - [ ] Search results page
  - [ ] Autocomplete dropdown

- [ ] **Basic Player Integration**
  - [ ] Player page (`/play/[id]`)
  - [ ] HLS.js integration
  - [ ] Basic controls (play, pause, seek)
  - [ ] Progress tracking
  - [ ] Quality selection (if available)
  - [ ] Subtitle selection (if available)

- [ ] **Settings** (`/settings`)
  - [ ] Profile settings
  - [ ] Playback preferences
  - [ ] Language preference

- [ ] **Admin Pages** (`/admin/*`)
  - [ ] Dashboard overview
  - [ ] Library management
  - [ ] User management
  - [ ] Integration settings (Radarr)

- [ ] **Components** (shadcn-svelte based)
  - [ ] Button, Input, Card
  - [ ] Dialog, Sheet
  - [ ] Select, Dropdown
  - [ ] Avatar, Badge
  - [ ] Skeleton loaders
  - [ ] Toast notifications

- [ ] **API Client**
  - [ ] Type-safe API client (generated or manual)
  - [ ] Error handling
  - [ ] Token refresh logic
  - [ ] TanStack Query integration

### Infrastructure

- [ ] **Typesense Deployment**
  - [ ] Docker Compose service
  - [ ] Helm chart subchart
  - [ ] Environment variables

- [ ] **Full Docker Compose Stack**
  - [ ] revenge (backend)
  - [ ] revenge-frontend
  - [ ] postgresql
  - [ ] dragonfly
  - [ ] typesense
  - [ ] traefik (reverse proxy)

- [ ] **Docker Images**
  - [ ] Backend multi-stage Dockerfile
  - [ ] Frontend multi-stage Dockerfile
  - [ ] Combined nginx config

### Documentation

- [ ] **User Documentation**
  - [ ] Getting started guide
  - [ ] Installation guide (Docker)
  - [ ] Configuration reference
  - [ ] Radarr setup guide

- [ ] **API Documentation**
  - [ ] Complete OpenAPI spec
  - [ ] Swagger UI endpoint
  - [ ] API authentication guide

---

## Verification Checklist

- [ ] Movies display in frontend
- [ ] Search works end-to-end
- [ ] Radarr sync imports movies
- [ ] Watch progress saves and restores
- [ ] Player plays video files
- [ ] Authentication works (login/logout)
- [ ] RBAC enforced on admin pages
- [ ] All tests pass (80%+ coverage)
- [ ] CI pipeline passes
- [ ] Docker Compose stack works

---

## MVP Definition of Done

The MVP is complete when a user can:

1. ✅ Register and login
2. ✅ Browse their movie library
3. ✅ Search for movies
4. ✅ View movie details (metadata, cast, crew)
5. ✅ Play a movie
6. ✅ Resume watching from where they left off
7. ✅ Admin can add Radarr integration
8. ✅ Movies sync from Radarr automatically

---

## Dependencies from SOURCE_OF_TRUTH

### Backend
| Package | Version | Purpose |
|---------|---------|---------|
| github.com/typesense/typesense-go/v4 | v4.x | Typesense client |
| github.com/go-resty/resty/v2 | v2.17.1 | HTTP client (Radarr) |

### Frontend
| Package | Version | Purpose |
|---------|---------|---------|
| SvelteKit | 2.x | Framework |
| Svelte | 5.x | UI library |
| Tailwind CSS | 4.x | Styling |
| shadcn-svelte | latest | Components |
| TanStack Query | latest | Data fetching |

---

## Related Documentation

- [ROADMAP.md](ROADMAP.md) - Full roadmap overview
- [00_SOURCE_OF_TRUTH.md](../00_SOURCE_OF_TRUTH.md) - Authoritative versions
- [MOVIE_MODULE.md](../features/video/MOVIE_MODULE.md) - Movie module design
- [TMDB.md](../integrations/metadata/video/TMDB.md) - TMDb integration
- [RADARR.md](../integrations/servarr/RADARR.md) - Radarr integration
- [TYPESENSE.md](../integrations/infrastructure/TYPESENSE.md) - Search setup

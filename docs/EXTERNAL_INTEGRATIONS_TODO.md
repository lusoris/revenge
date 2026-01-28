# External Integrations TODO

> Comprehensive list of ALL external services/APIs that Revenge integrates with.
> For each service: Dev docs, API reference, authentication, rate limits, and implementation needs.

**Status Legend:**
- ✅ Documentation exists
- ⚠️ Partially documented
- ❌ Not documented
- 📚 Dev docs available
- 🔗 API reference available

---

## 1. Content Management (Servarr Ecosystem)

### 1.1 Radarr - Movie Management

**Purpose:** Automatic movie downloading, metadata, library management
**Integration:** Webhook listener + API client for metadata sync
**Priority:** 🔴 CRITICAL (Phase 2 - Movie Module)

**Developer Resources:**
- 📚 API Docs: https://radarr.video/docs/api/
- 🔗 OpenAPI Spec: https://github.com/Radarr/Radarr/blob/develop/src/Radarr.Api.V3/openapi.json
- 🔗 GitHub: https://github.com/Radarr/Radarr
- 🔗 Wiki: https://wiki.servarr.com/radarr

**Implementation Needs:**
- [ ] `docs/integrations/RADARR.md` - Integration guide
- [ ] `.github/instructions/radarr-client.instructions.md` - Client patterns
- [ ] `internal/service/metadata/provider_radarr.go` - API client
- [ ] Webhook handler for import events
- [ ] Metadata sync (title, overview, posters, cast, crew)
- [ ] Quality profile mapping
- [ ] Root folder management

**API Details:**
- Authentication: API Key header `X-Api-Key`
- Rate Limit: None documented (self-hosted)
- Base Path: `/api/v3/`
- Key Endpoints: `/movie`, `/importlist`, `/metadata`, `/mediamanagement`

---

### 1.2 Sonarr - TV Show Management

**Purpose:** Automatic TV show downloading, metadata, library management
**Integration:** Webhook listener + API client for metadata sync
**Priority:** 🔴 CRITICAL (Phase 3 - TV Show Module)

**Developer Resources:**
- 📚 API Docs: https://sonarr.tv/docs/api/
- 🔗 OpenAPI Spec: https://github.com/Sonarr/Sonarr/blob/develop/src/Sonarr.Api.V3/openapi.json
- 🔗 GitHub: https://github.com/Sonarr/Sonarr
- 🔗 Wiki: https://wiki.servarr.com/sonarr

**Implementation Needs:**
- [ ] `docs/integrations/SONARR.md`
- [ ] `.github/instructions/sonarr-client.instructions.md`
- [ ] `internal/service/metadata/provider_sonarr.go`
- [ ] Webhook handler for import events
- [ ] Series/season/episode metadata sync
- [ ] Quality profile mapping
- [ ] Season pack handling

**API Details:**
- Authentication: API Key header `X-Api-Key`
- Rate Limit: None documented (self-hosted)
- Base Path: `/api/v3/`
- Key Endpoints: `/series`, `/episode`, `/episodefile`, `/importlist`

---

### 1.3 Lidarr - Music Management

**Purpose:** Automatic music downloading, metadata, library management
**Integration:** Webhook listener + API client for metadata sync
**Priority:** 🟡 HIGH (Phase 4 - Music Module)

**Developer Resources:**
- 📚 API Docs: https://lidarr.audio/docs/api/
- 🔗 OpenAPI Spec: https://github.com/Lidarr/Lidarr/blob/develop/src/Lidarr.Api.V3/openapi.json
- 🔗 GitHub: https://github.com/Lidarr/Lidarr
- 🔗 Wiki: https://wiki.servarr.com/lidarr

**Implementation Needs:**
- [ ] `docs/integrations/LIDARR.md`
- [ ] `.github/instructions/lidarr-client.instructions.md`
- [ ] `internal/service/metadata/provider_lidarr.go`
- [ ] Artist/album/track metadata sync
- [ ] MusicBrainz ID mapping
- [ ] Quality profile mapping
- [ ] Multi-artist album handling

**API Details:**
- Authentication: API Key header `X-Api-Key`
- Rate Limit: None documented (self-hosted)
- Base Path: `/api/v1/`
- Key Endpoints: `/artist`, `/album`, `/track`, `/importlist`

---

### 1.4 Whisparr - Adult Content Management

**Purpose:** Automatic adult content downloading, metadata, library management
**Integration:** Webhook listener + API client for metadata sync
**Priority:** 🟡 MEDIUM (Phase 7 - Adult Modules)

**Developer Resources:**
- 📚 API Docs: https://whisparr.com/docs/api/ (if available)
- 🔗 GitHub: https://github.com/Whisparr/Whisparr
- 🔗 Based on Radarr/Sonarr API structure

**Implementation Needs:**
- [ ] `docs/integrations/WHISPARR.md`
- [ ] `.github/instructions/whisparr-client.instructions.md`
- [ ] `internal/content/c/movie/provider_whisparr.go`
- [ ] Performer metadata sync
- [ ] Studio/tag handling
- [ ] Scene extraction
- [ ] Privacy controls (schema `c` isolation)

**API Details:**
- Authentication: API Key header `X-Api-Key`
- Rate Limit: None documented (self-hosted)
- Base Path: `/api/v3/` (assumed, Radarr-based)
- Key Endpoints: TBD (similar to Radarr)

---

### 1.5 Readarr - Book Management

**Purpose:** Book downloading, metadata, library management
**Integration:** Webhook listener + API client for metadata sync
**Priority:** 🟡 MEDIUM (Phase 6 - Book Module)

**Developer Resources:**
- 📚 API Docs: https://readarr.com/docs/api/
- 🔗 OpenAPI Spec: https://github.com/Readarr/Readarr/blob/develop/src/Readarr.Api.V1/openapi.json
- 🔗 GitHub: https://github.com/Readarr/Readarr
- 🔗 Wiki: https://wiki.servarr.com/readarr

**Implementation Needs:**
- [ ] `docs/integrations/READARR.md`
- [ ] `.github/instructions/readarr-client.instructions.md`
- [ ] `internal/content/book/provider_readarr.go`
- [ ] Book/author metadata sync
- [ ] Edition handling
- [ ] Quality profile mapping
- [ ] Audiobook vs ebook separation

**API Details:**
- Authentication: API Key header `X-Api-Key`
- Rate Limit: None documented (self-hosted)
- Base Path: `/api/v1/`
- Key Endpoints: `/book`, `/author`, `/bookfile`, `/importlist`

---

## 2. Audiobook Management

### 2.1 Audiobookshelf

**Purpose:** Audiobook, book, podcast library management
**Integration:** API client for metadata + library sync
**Priority:** 🟡 HIGH (Phase 6 - Audiobook/Book/Podcast Modules)

**Developer Resources:**
- 📚 API Docs: https://api.audiobookshelf.org/
- 🔗 GitHub: https://github.com/advplyr/audiobookshelf
- 🔗 OpenAPI Spec: https://github.com/advplyr/audiobookshelf/blob/master/server/openapi.yaml

**Implementation Needs:**
- [ ] `docs/integrations/AUDIOBOOKSHELF.md`
- [ ] `.github/instructions/audiobookshelf-client.instructions.md`
- [ ] `internal/content/audiobook/provider_audiobookshelf.go`
- [ ] `internal/content/book/provider_audiobookshelf.go`
- [ ] `internal/content/podcast/provider_audiobookshelf.go`
- [ ] Library sync (audiobooks, books, podcasts)
- [ ] Chapter marker sync
- [ ] Progress sync (bidirectional)
- [ ] Narrator metadata

**API Details:**
- Authentication: Bearer token header `Authorization: Bearer {token}`
- Rate Limit: None documented (self-hosted)
- Base Path: `/api/`
- Key Endpoints: `/libraries`, `/items`, `/me/progress`, `/podcasts`

---

## 3. Request Management

### 3.1 Overseerr

**Purpose:** Media request management (movies, TV shows)
**Integration:** API client for request handling + webhooks
**Priority:** 🟡 MEDIUM (Phase 9 - External Services)

**Developer Resources:**
- 📚 API Docs: https://api-docs.overseerr.dev/
- 🔗 GitHub: https://github.com/sct/overseerr
- 🔗 OpenAPI Spec: https://github.com/sct/overseerr/blob/develop/overseerr-api.yml

**Implementation Needs:**
- [ ] `docs/integrations/OVERSEERR.md`
- [ ] `.github/instructions/overseerr-adapter.instructions.md`
- [ ] `internal/service/requests/adapter_overseerr.go`
- [ ] Request status sync
- [ ] Approval workflow
- [ ] User quota management
- [ ] Webhook events

**API Details:**
- Authentication: API Key header `X-Api-Key`
- Rate Limit: None documented (self-hosted)
- Base Path: `/api/v1/`
- Key Endpoints: `/request`, `/movie`, `/tv`, `/user`

---

### 3.2 Jellyseerr

**Purpose:** Jellyfin-focused media request management
**Integration:** API client for request handling + webhooks
**Priority:** 🟡 MEDIUM (Phase 9 - External Services)

**Developer Resources:**
- 📚 API Docs: https://github.com/Fallenbagel/jellyseerr (fork of Overseerr)
- 🔗 GitHub: https://github.com/Fallenbagel/jellyseerr
- 🔗 Compatible with Overseerr API

**Implementation Needs:**
- [ ] `docs/integrations/JELLYSEERR.md`
- [ ] Can reuse Overseerr adapter (compatible API)
- [ ] Jellyfin-specific authentication flow
- [ ] User mapping from Jellyfin

**API Details:**
- Same as Overseerr (forked)
- Authentication: API Key header `X-Api-Key`

---

## 4. Metadata Providers

### 4.1 TMDb (The Movie Database)

**Purpose:** Primary metadata for movies and TV shows
**Integration:** HTTP API client with rate limiting
**Priority:** 🔴 CRITICAL (Phase 2/3 - Movie/TV Modules)

**Developer Resources:**
- 📚 API Docs: https://developers.themoviedb.org/3
- 🔗 API v3: https://api.themoviedb.org/3/
- 🔗 API v4 (experimental): https://developers.themoviedb.org/4
- 🔗 Image CDN: https://image.tmdb.org/t/p/
- 🔗 Status Page: https://status.themoviedb.org/

**Implementation Needs:**
- [ ] `docs/integrations/TMDB.md`
- [ ] `.github/instructions/tmdb-provider.instructions.md`
- [ ] `internal/service/metadata/provider_tmdb.go`
- [ ] Movie metadata (title, overview, cast, crew, images)
- [ ] TV show metadata (series, seasons, episodes)
- [ ] Image download (posters, backdrops, logos)
- [ ] Blurhash generation
- [ ] External ID mapping (IMDb, TheTVDB)

**API Details:**
- Authentication: API Key query param `?api_key={key}`
- Rate Limit: 40 requests per 10 seconds
- Free tier: Available
- Regions: i18n support with `language` param
- Key Endpoints: `/movie/{id}`, `/tv/{id}`, `/person/{id}`, `/search`

---

### 4.2 TheTVDB

**Purpose:** Primary metadata for TV shows
**Integration:** HTTP API client with JWT authentication
**Priority:** 🔴 CRITICAL (Phase 3 - TV Show Module)

**Developer Resources:**
- 📚 API Docs: https://thetvdb.github.io/v4-api/
- 🔗 API v4: https://api4.thetvdb.com/v4/
- 🔗 GitHub: https://github.com/thetvdb/v4-api

**Implementation Needs:**
- [ ] `docs/integrations/THETVDB.md`
- [ ] `.github/instructions/thetvdb-provider.instructions.md`
- [ ] `internal/service/metadata/provider_thetvdb.go`
- [ ] Series metadata
- [ ] Season/episode metadata
- [ ] Artwork (posters, fanart, banners)
- [ ] JWT token refresh
- [ ] External ID mapping

**API Details:**
- Authentication: JWT Bearer token (POST /login first)
- Rate Limit: Varies by subscription tier
- Free tier: Limited requests
- Key Endpoints: `/series/{id}`, `/seasons/{id}`, `/episodes/{id}`, `/artwork`

---

### 4.3 OMDb (Open Movie Database)

**Purpose:** Fallback metadata, IMDb ratings
**Integration:** HTTP API client
**Priority:** 🟡 HIGH (Phase 2 - Movie Module)

**Developer Resources:**
- 📚 API Docs: https://www.omdbapi.com/
- 🔗 API: http://www.omdbapi.com/

**Implementation Needs:**
- [ ] `docs/integrations/OMDB.md`
- [ ] `.github/instructions/omdb-provider.instructions.md`
- [ ] `internal/service/metadata/provider_omdb.go`
- [ ] Basic movie metadata
- [ ] IMDb ratings
- [ ] Poster fallback

**API Details:**
- Authentication: API Key query param `?apikey={key}`
- Rate Limit: 1,000 requests/day (free tier)
- Key Endpoints: `/?i={imdb_id}`, `/?t={title}&y={year}`

---

### 4.4 Fanart.tv

**Purpose:** High-quality artwork (posters, logos, clearart)
**Integration:** HTTP API client
**Priority:** 🟡 MEDIUM (Phase 8 - Media Enhancements)

**Developer Resources:**
- 📚 API Docs: https://fanart.tv/api-docs/
- 🔗 API: http://webservice.fanart.tv/v3/

**Implementation Needs:**
- [ ] `docs/integrations/FANART_TV.md`
- [ ] `.github/instructions/fanart-provider.instructions.md`
- [ ] `internal/service/metadata/provider_fanart.go`
- [ ] HD posters, backdrops
- [ ] ClearLogo, ClearArt
- [ ] Banner, thumb images

**API Details:**
- Authentication: API Key query param `?api_key={key}`
- Rate Limit: None documented
- Key Endpoints: `/movies/{tmdb_id}`, `/tv/{tvdb_id}`, `/music/{musicbrainz_id}`

---

### 4.5 ThePosterDB

**Purpose:** Curated posters for movies/TV
**Integration:** HTTP scraper (no official API)
**Priority:** 🟢 LOW (Phase 8 - Media Enhancements)

**Developer Resources:**
- 🔗 Website: https://theposterdb.com/
- ❌ No official API
- 🔗 Community API: https://github.com/jarulsamy/ThePosterDB-API (unofficial)

**Implementation Needs:**
- [ ] `docs/integrations/THEPOSTERDB.md`
- [ ] `.github/instructions/posterdb-scraper.instructions.md`
- [ ] `internal/service/metadata/provider_posterdb.go`
- [ ] Web scraping with rate limiting
- [ ] Poster downloads
- [ ] Set collections (4K, IMAX, etc.)

**API Details:**
- ❌ No official API
- Unofficial API available (community-maintained)
- Respect robots.txt and rate limits

---

### 4.6 MusicBrainz

**Purpose:** Primary metadata for music (artists, albums, tracks)
**Integration:** HTTP API client with User-Agent requirement
**Priority:** 🟡 HIGH (Phase 4 - Music Module)

**Developer Resources:**
- 📚 API Docs: https://musicbrainz.org/doc/MusicBrainz_API
- 🔗 API: https://musicbrainz.org/ws/2/
- 🔗 Cover Art Archive: https://coverartarchive.org/
- 🔗 Rate Limiting: https://musicbrainz.org/doc/MusicBrainz_API/Rate_Limiting

**Implementation Needs:**
- [ ] `docs/integrations/MUSICBRAINZ.md`
- [ ] `.github/instructions/musicbrainz-provider.instructions.md`
- [ ] `internal/content/music/provider_musicbrainz.go`
- [ ] Artist metadata
- [ ] Album/release metadata
- [ ] Track metadata
- [ ] Cover art via Cover Art Archive
- [ ] User-Agent compliance

**API Details:**
- Authentication: None (User-Agent required)
- Rate Limit: 1 request per second (50/minute burst)
- User-Agent: `{AppName}/{Version} ( {ContactEmail} )`
- Key Endpoints: `/artist/{mbid}`, `/release/{mbid}`, `/recording/{mbid}`

---

### 4.7 Last.fm

**Purpose:** Music metadata, artist bio, tags, similar artists
**Integration:** HTTP API client
**Priority:** 🟡 HIGH (Phase 4 - Music Module + Scrobbling)

**Developer Resources:**
- 📚 API Docs: https://www.last.fm/api
- 🔗 API: https://ws.audioscrobbler.com/2.0/
- 🔗 Scrobbling: https://www.last.fm/api/scrobbling

**Implementation Needs:**
- [ ] `docs/integrations/LASTFM.md`
- [ ] `.github/instructions/lastfm-client.instructions.md`
- [ ] `internal/content/music/provider_lastfm.go` (metadata)
- [ ] `internal/service/scrobble/client_lastfm.go` (scrobbling)
- [ ] Artist bio, tags, similar artists
- [ ] Album/track info
- [ ] Scrobble "Now Playing" + track completion
- [ ] OAuth authentication for scrobbling

**API Details:**
- Authentication: API Key + Secret (metadata) / OAuth (scrobbling)
- Rate Limit: Not strictly enforced
- Key Endpoints: `/artist.getInfo`, `/track.scrobble`, `/auth.getSession`

---

### 4.8 Spotify

**Purpose:** Music cover art, popularity scores
**Integration:** HTTP API client with OAuth
**Priority:** 🟡 MEDIUM (Phase 4 - Music Module)

**Developer Resources:**
- 📚 API Docs: https://developer.spotify.com/documentation/web-api
- 🔗 API: https://api.spotify.com/v1/
- 🔗 Auth: https://accounts.spotify.com/api/token

**Implementation Needs:**
- [ ] `docs/integrations/SPOTIFY.md`
- [ ] `.github/instructions/spotify-provider.instructions.md`
- [ ] `internal/content/music/provider_spotify.go`
- [ ] Album cover art (high quality)
- [ ] Popularity scores
- [ ] Artist images
- [ ] Client Credentials OAuth flow

**API Details:**
- Authentication: OAuth 2.0 (Client Credentials for metadata)
- Rate Limit: Varies (typically lenient)
- Key Endpoints: `/search`, `/albums/{id}`, `/artists/{id}`, `/tracks/{id}`

---

### 4.9 Discogs

**Purpose:** Music release metadata, labels, formats
**Integration:** HTTP API client
**Priority:** 🟢 LOW (Phase 4 - Music Module)

**Developer Resources:**
- 📚 API Docs: https://www.discogs.com/developers
- 🔗 API: https://api.discogs.com/
- 🔗 Rate Limiting: https://www.discogs.com/developers#page:home,header:home-rate-limiting

**Implementation Needs:**
- [ ] `docs/integrations/DISCOGS.md`
- [ ] `.github/instructions/discogs-provider.instructions.md`
- [ ] `internal/content/music/provider_discogs.go`
- [ ] Release metadata
- [ ] Label information
- [ ] Format details (vinyl, CD, etc.)

**API Details:**
- Authentication: OAuth or Personal Access Token
- Rate Limit: 60 requests/minute (authenticated)
- User-Agent required
- Key Endpoints: `/releases/{id}`, `/artists/{id}`, `/labels/{id}`

---

### 4.10 Goodreads

**Purpose:** Book ratings, reviews, metadata
**Integration:** Web scraping (API retired 2020)
**Priority:** 🟡 MEDIUM (Phase 6 - Book Module)

**Developer Resources:**
- ❌ API retired December 2020
- 🔗 Website: https://www.goodreads.com/
- 🔗 Alternative: OpenLibrary

**Implementation Needs:**
- [ ] `docs/integrations/GOODREADS.md`
- [ ] `.github/instructions/goodreads-scraper.instructions.md`
- [ ] `internal/content/book/provider_goodreads.go`
- [ ] Web scraping with caching
- [ ] Book ratings/reviews
- [ ] Author info
- [ ] Respectful scraping (robots.txt, rate limiting)

**API Details:**
- ❌ No API (scraping only)
- Rate limiting required (manual)
- Consider OpenLibrary as alternative

---

### 4.11 OpenLibrary

**Purpose:** Open book metadata database
**Integration:** HTTP API client
**Priority:** 🟡 HIGH (Phase 6 - Book Module)

**Developer Resources:**
- 📚 API Docs: https://openlibrary.org/developers/api
- 🔗 API: https://openlibrary.org/api/
- 🔗 Covers API: https://covers.openlibrary.org/

**Implementation Needs:**
- [ ] `docs/integrations/OPENLIBRARY.md`
- [ ] `.github/instructions/openlibrary-provider.instructions.md`
- [ ] `internal/content/book/provider_openlibrary.go`
- [ ] Book metadata (title, author, ISBN)
- [ ] Cover images
- [ ] Author bio
- [ ] Edition information

**API Details:**
- Authentication: None (public API)
- Rate Limit: Not strictly enforced (be respectful)
- Key Endpoints: `/books/{isbn}.json`, `/authors/{olid}.json`, `/search.json`

---

### 4.12 Audible

**Purpose:** Audiobook metadata, narrators
**Integration:** Web scraping (no public API)
**Priority:** 🟡 MEDIUM (Phase 6 - Audiobook Module)

**Developer Resources:**
- ❌ No public API
- 🔗 Website: https://www.audible.com/
- 🔗 Community lib: https://github.com/mkb79/audible (Python)

**Implementation Needs:**
- [ ] `docs/integrations/AUDIBLE.md`
- [ ] `.github/instructions/audible-scraper.instructions.md`
- [ ] `internal/content/audiobook/provider_audible.go`
- [ ] Web scraping for audiobook metadata
- [ ] Narrator information
- [ ] Series information
- [ ] Duration, ratings

**API Details:**
- ❌ No public API
- Scraping only (respect rate limits)
- Alternative: Audiobookshelf metadata

---

### 4.13 Hardcover

**Purpose:** Book metadata, reading tracking (Goodreads alternative)
**Integration:** GraphQL API
**Priority:** 🟡 MEDIUM (Phase 6 - Book Module)

**Developer Resources:**
- 📚 API Docs: https://hardcover.app/docs/api
- 🔗 GraphQL API: https://api.hardcover.app/v1/graphql
- 🔗 Website: https://hardcover.app/

**Implementation Needs:**
- [ ] `docs/integrations/HARDCOVER.md`
- [ ] `.github/instructions/hardcover-provider.instructions.md`
- [ ] `internal/content/book/provider_hardcover.go`
- [ ] GraphQL client for book metadata
- [ ] Book ratings/reviews
- [ ] Author information
- [ ] Reading lists/progress sync (optional)

**API Details:**
- Authentication: API Key (header `Authorization: Bearer {token}`)
- Rate Limit: TBD (check docs)
- GraphQL endpoint
- Key Queries: `book`, `author`, `search`

---

### 4.14 Stash / StashDB

**Purpose:** Adult content metadata (scenes, performers, studios)
**Integration:** GraphQL API (Stash), HTTP API (StashDB)
**Priority:** 🟡 MEDIUM (Phase 7 - Adult Modules)

**Developer Resources:**
- 📚 Stash API Docs: https://github.com/stashapp/stash/blob/develop/graphql/schema/schema.graphql
- 📚 StashDB API: https://stashdb.org/graphql (GraphQL)
- 🔗 GitHub: https://github.com/stashapp/stash

**Implementation Needs:**
- [ ] `docs/integrations/STASH.md` (update existing ADULT_METADATA.md)
- [ ] `.github/instructions/stash-provider.instructions.md`
- [ ] `internal/content/c/movie/provider_stash.go`
- [ ] GraphQL client for scene metadata
- [ ] Performer metadata (bio, measurements, images)
- [ ] Studio information
- [ ] Tag/category mapping
- [ ] Fingerprinting (phash)

**API Details:**
- Authentication: API Key (Stash local) / StashDB OAuth
- Rate Limit: None (self-hosted Stash) / StashDB varies
- GraphQL endpoints
- Key Queries: `findScene`, `findPerformer`, `findStudio`

---

### 4.15 ThePornDB

**Purpose:** Adult content metadata (movies, scenes, performers)
**Integration:** HTTP API
**Priority:** 🟡 MEDIUM (Phase 7 - Adult Modules)

**Developer Resources:**
- 📚 API Docs: https://theporndb.net/docs/api
- 🔗 API: https://api.theporndb.net/
- 🔗 Requires account + API key

**Implementation Needs:**
- [ ] `docs/integrations/THEPORNDB.md` (update existing ADULT_METADATA.md)
- [ ] `.github/instructions/theporndb-provider.instructions.md`
- [ ] `internal/content/c/movie/provider_theporndb.go`
- [ ] Movie/scene metadata
- [ ] Performer metadata
- [ ] Studio information
- [ ] Image downloads

**API Details:**
- Authentication: API Key (header `Authorization: Bearer {key}`)
- Rate Limit: TBD (check API docs)
- Key Endpoints: `/scenes/{id}`, `/performers/{id}`, `/studios/{id}`

---

## 5. Scrobbling / Rating Sync

### 5.1 Trakt

**Purpose:** Watch history, ratings sync for movies/TV
**Integration:** OAuth + HTTP API client
**Priority:** 🟡 HIGH (Phase 9 - External Services)

**Developer Resources:**
- 📚 API Docs: https://trakt.docs.apiary.io/
- 🔗 API: https://api.trakt.tv/
- 🔗 OAuth: https://trakt.tv/oauth/applications

**Implementation Needs:**
- [ ] `docs/integrations/TRAKT.md`
- [ ] `.github/instructions/trakt-scrobbler.instructions.md`
- [ ] `internal/service/scrobble/client_trakt.go`
- [ ] OAuth 2.0 flow (user authorization)
- [ ] Watch history sync (movies, episodes)
- [ ] Ratings sync (bidirectional)
- [ ] Collection/watchlist sync
- [ ] "Now Watching" scrobble
- [ ] Progress tracking

**API Details:**
- Authentication: OAuth 2.0 (Bearer token)
- Rate Limit: 1,000 requests per 5 minutes (per user)
- Client ID + Client Secret required
- Key Endpoints: `/sync/history`, `/sync/ratings`, `/scrobble/start`, `/scrobble/stop`

---

### 5.2 Last.fm (Scrobbling)

**Purpose:** Music scrobbling, play count tracking
**Integration:** OAuth + HTTP API client
**Priority:** 🟡 HIGH (Phase 4 - Music Module + Phase 9)

**Developer Resources:**
- (See 4.7 Last.fm Metadata above - same API)

**Implementation Needs:**
- [ ] `internal/service/scrobble/client_lastfm.go`
- [ ] OAuth authentication
- [ ] Scrobble track plays
- [ ] "Now Playing" updates
- [ ] Love/unlove tracks

**API Details:**
- (See 4.7 above)

---

### 5.3 ListenBrainz

**Purpose:** Open-source music scrobbling (Last.fm alternative)
**Integration:** HTTP API client with token auth
**Priority:** 🟡 MEDIUM (Phase 4 - Music Module + Phase 9)

**Developer Resources:**
- 📚 API Docs: https://listenbrainz.readthedocs.io/
- 🔗 API: https://api.listenbrainz.org/
- 🔗 GitHub: https://github.com/metabrainz/listenbrainz-server

**Implementation Needs:**
- [ ] `docs/integrations/LISTENBRAINZ.md`
- [ ] `.github/instructions/listenbrainz-scrobbler.instructions.md`
- [ ] `internal/service/scrobble/client_listenbrainz.go`
- [ ] Token authentication
- [ ] Submit listens (scrobbles)
- [ ] "Now Playing" updates
- [ ] MusicBrainz ID mapping

**API Details:**
- Authentication: User token (header `Authorization: Token {token}`)
- Rate Limit: Not strictly enforced
- Key Endpoints: `/1/submit-listens`, `/1/playing-now`

---

### 5.4 Letterboxd

**Purpose:** Movie diary, ratings, reviews
**Integration:** Web scraping (no public API) OR RSS export
**Priority:** 🟢 LOW (Phase 9 - External Services)

**Developer Resources:**
- ❌ No public API
- 🔗 Website: https://letterboxd.com/
- 🔗 RSS Exports: User diary/watchlist feeds

**Implementation Needs:**
- [ ] `docs/integrations/LETTERBOXD.md`
- [ ] `.github/instructions/letterboxd-sync.instructions.md`
- [ ] `internal/service/scrobble/sync_letterboxd.go`
- [ ] RSS feed parsing for diary exports
- [ ] One-way sync only (export from Revenge)
- [ ] CSV export format

**API Details:**
- ❌ No API
- RSS feeds available per-user
- Export only (no import via API)

---

### 5.5 Simkl

**Purpose:** Multi-platform watch tracking (anime, shows, movies)
**Integration:** OAuth + HTTP API client
**Priority:** 🟢 LOW (Phase 9 - External Services)

**Developer Resources:**
- 📚 API Docs: https://simkl.docs.apiary.io/
- 🔗 API: https://api.simkl.com/
- 🔗 OAuth: https://simkl.com/apps/

**Implementation Needs:**
- [ ] `docs/integrations/SIMKL.md`
- [ ] `.github/instructions/simkl-scrobbler.instructions.md`
- [ ] `internal/service/scrobble/client_simkl.go`
- [ ] OAuth 2.0 flow
- [ ] Watch history sync
- [ ] Ratings sync
- [ ] Anime support

**API Details:**
- Authentication: OAuth 2.0 (Bearer token)
- Rate Limit: Not documented
- Key Endpoints: `/sync/history`, `/sync/ratings`, `/anime/episodes/watched`

---

## 5A. Anime Tracking

### 5A.1 AniList

**Purpose:** Anime/manga tracking, metadata, ratings (GraphQL API)
**Integration:** GraphQL client with OAuth
**Priority:** 🟡 HIGH (Anime module + scrobbling)

**Developer Resources:**
- 📚 API Docs: https://docs.anilist.co/
- 🔗 GraphQL Endpoint: https://graphql.anilist.co
- 🔗 OAuth: https://anilist.co/api/v2/oauth/authorize
- 🔗 Interactive Explorer: https://studio.apollographql.com/sandbox/explorer?endpoint=https://graphql.anilist.co
- 🔗 GitHub: https://github.com/AniList/docs (VitePress, TypeScript, 1.2k stars)

**Implementation Needs:**
- [ ] `docs/integrations/ANILIST.md`
- [ ] `.github/instructions/anilist-graphql.instructions.md`
- [ ] `internal/content/anime/provider_anilist.go` (metadata)
- [ ] `internal/service/scrobble/client_anilist.go` (tracking)
- [ ] GraphQL client (queries, mutations)
- [ ] OAuth 2.0 flow (implicit + authorization code)
- [ ] Anime/manga metadata (500k+ entries)
- [ ] Character/staff data with relations
- [ ] Airing schedule tracking
- [ ] User list sync (watching/completed/dropped)
- [ ] Ratings/reviews sync

**API Details:**
- Authentication: OAuth 2.0 (implicit/authorization code) → Bearer token
- Rate Limit: Not specified (90 req/min recommended by community)
- GraphQL schema with queries/mutations
- Data model: titles (multiple languages), ratings, airing status, relations, characters, staff, streaming links
- Apollo Studio recommended (GraphiQL deprecated)
- Static docs at https://docs.anilist.co/reference/query (manually updated, may be outdated)

---

### 5A.2 MyAnimeList (MAL)

**Purpose:** Largest anime/manga database, tracking, forums
**Integration:** REST API v2 (beta) with OAuth
**Priority:** 🟡 HIGH (Anime module + scrobbling)

**Developer Resources:**
- 📚 API Docs: https://myanimelist.net/apiconfig/references/api/v2
- 🔗 Base API: https://api.myanimelist.net/v2
- 🔗 OAuth: https://myanimelist.net/v1/oauth2/authorize
- 🔗 Official Club: https://myanimelist.net/clubs.php?cid=13727 (1,839 members)

**Implementation Needs:**
- [ ] `docs/integrations/MYANIMELIST.md`
- [ ] `.github/instructions/mal-api-client.instructions.md`
- [ ] `internal/content/anime/provider_mal.go` (metadata)
- [ ] `internal/service/scrobble/client_mal.go` (tracking)
- [ ] OAuth 2.0 flow (implicit flow, `write:users` scope) OR X-MAL-CLIENT-ID header
- [ ] Anime search/details/ranking/seasonal
- [ ] Manga search/details/ranking
- [ ] User library (animelist/mangalist) sync
- [ ] Status updates (watching/completed/dropped, episode progress)
- [ ] Ratings sync (0-10 score, tags, comments)
- [ ] Forum integration (optional)

**API Details:**
- Authentication: OAuth 2.0 (`write:users` scope) OR `X-MAL-CLIENT-ID` header (non-user endpoints)
- Rate Limit: Not specified (reasonable usage)
- Versioning: v0.x during beta, increments on breaking changes
- Base Path: `/api/myanimelist.net/v2`
- Key Endpoints: `/anime` (search limit 100), `/anime/{id}`, `/anime/ranking` (limit 500), `/anime/season/{year}/{season}`, `/manga`, `/users/{user}/animelist` (limit 1000), PATCH `/anime/{id}/my_list_status`
- Sparse fieldsets via `fields` param
- Date formats: ISO 8601, YYYY-MM-DD, YYYY-MM, YYYY
- Error codes: 400 (invalid params), 401 (invalid_token), 403 (DoS), 404

---

### 5A.3 Kitsu

**Purpose:** Modern anime/manga API (JSON:API spec)
**Integration:** REST API with JSON:API conventions
**Priority:** 🟡 MEDIUM (Anime module alternative)

**Developer Resources:**
- 📚 API Docs: https://kitsu.docs.apiary.io/
- 🔗 Base API: https://kitsu.io/api/edge
- 🔗 OAuth: https://kitsu.io/api/oauth
- 🔗 GitHub Tools: https://github.com/hummingbird-me/kitsu-tools (2.2k stars, Apache-2.0, docker-compose)

**Implementation Needs:**
- [ ] `docs/integrations/KITSU.md`
- [ ] `.github/instructions/kitsu-jsonapi.instructions.md`
- [ ] `internal/content/anime/provider_kitsu.go` (metadata)
- [ ] `internal/service/scrobble/client_kitsu.go` (tracking)
- [ ] JSON:API client (filtering, pagination, includes, sparse fieldsets)
- [ ] OAuth 2.0 flow
- [ ] Anime/manga metadata with relationships
- [ ] Episodes/chapters data
- [ ] Trending anime/manga
- [ ] Streaming links integration
- [ ] User library sync (extensive filters)
- [ ] Cross-reference mappings (MAL integration)

**API Details:**
- Authentication: OAuth 2.0 → Bearer token
- Rate Limit: Not specified (be respectful)
- JSON:API spec compliance (https://jsonapi.org/format/)
- Required headers: `Accept: application/vnd.api+json`, `Content-Type: application/vnd.api+json`
- Filtering: `filter[attribute]=value`, `filter[text]=query`
- Pagination: `page[limit]=10` (default 10, max 20), `page[offset]=0`
- Sorting: `sort=-field` (descending), comma-delimited
- Includes: `include=relationship.nested`
- Sparse fieldsets: `fields[type]=field1,field2`
- Key Endpoints: `/anime`, `/manga`, `/trending/anime`, `/library-entries` (max limit 500), `/streaming-links`, `/media-relationships`, `/mappings` (externalSite: myanimelist/anime)
- NSFW content hidden for unauthenticated/disabled users

---

## 6. Authentication (OIDC/SSO)

### 6.1 Authelia

**Purpose:** SSO authentication provider
**Integration:** OIDC client
**Priority:** 🟡 MEDIUM (Phase 1 - Core Infrastructure)

**Developer Resources:**
- 📚 OIDC Docs: https://www.authelia.com/integration/openid-connect/introduction/
- 🔗 GitHub: https://github.com/authelia/authelia

**Implementation Needs:**
- [ ] `docs/integrations/AUTHELIA.md`
- [ ] OIDC discovery
- [ ] User profile mapping
- [ ] Group/role mapping

**API Details:**
- Standard OIDC endpoints
- Discovery: `/.well-known/openid-configuration`

---

### 6.2 Authentik

**Purpose:** SSO authentication provider
**Integration:** OIDC client
**Priority:** 🟡 MEDIUM (Phase 1 - Core Infrastructure)

**Developer Resources:**
- 📚 OIDC Docs: https://goauthentik.io/docs/providers/oauth2
- 🔗 GitHub: https://github.com/goauthentik/authentik

**Implementation Needs:**
- [ ] `docs/integrations/AUTHENTIK.md`
- [ ] OIDC discovery
- [ ] User profile mapping
- [ ] Group/role mapping

**API Details:**
- Standard OIDC endpoints
- Discovery: `/application/o/{slug}/.well-known/openid-configuration`

---

### 6.3 Keycloak

**Purpose:** Enterprise SSO authentication
**Integration:** OIDC client
**Priority:** 🟡 MEDIUM (Phase 1 - Core Infrastructure)

**Developer Resources:**
- 📚 OIDC Docs: https://www.keycloak.org/docs/latest/securing_apps/
- 🔗 GitHub: https://github.com/keycloak/keycloak

**Implementation Needs:**
- [ ] `docs/integrations/KEYCLOAK.md`
- [ ] OIDC discovery
- [ ] User profile mapping
- [ ] Realm/client configuration

**API Details:**
- Standard OIDC endpoints
- Discovery: `/realms/{realm}/.well-known/openid-configuration`

---

### 6.4 Generic OIDC

**Purpose:** Support any OIDC-compliant provider
**Integration:** OIDC client library
**Priority:** 🟡 MEDIUM (Phase 1 - Core Infrastructure)

**Developer Resources:**
- 📚 OIDC Spec: https://openid.net/specs/openid-connect-core-1_0.html

**Implementation Needs:**
- [ ] Generic OIDC discovery
- [ ] Configurable claim mapping
- [ ] Support for various providers (Okta, Auth0, Google, etc.)

---

## 7. Transcoding

### 7.1 Blackbeard (External Service)

**Purpose:** External transcoding service (HLS/DASH generation)
**Integration:** Internal HTTP API (Blackbeard → Revenge for raw files)
**Priority:** 🔴 CRITICAL (Phase 5 - Playback Service)

**Developer Resources:**
- 📚 Internal API (to be documented)
- 🔗 Repository: TBD (separate Revenge project)

**Implementation Needs:**
- [ ] `docs/integrations/BLACKBEARD.md` (update OFFLOADING.md)
- [ ] `.github/instructions/blackbeard-client.instructions.md`
- [ ] `internal/service/playback/transcoder.go` (already exists!)
- [ ] Profile selection API
- [ ] Stream request/response
- [ ] Progress tracking
- [ ] Internal file serving endpoint (Revenge → Blackbeard)

**API Details:**
- Authentication: Internal token
- Endpoints: `/transcode/start`, `/transcode/{id}/master.m3u8`
- Reverse endpoint: `GET /internal/stream/{mediaID}` (Revenge side)

---

## 8. Live TV (PVR Backends)

### 8.1 TVHeadend

**Purpose:** DVR backend, EPG, channel management
**Integration:** HTTP API + HTSP protocol
**Priority:** 🟢 LOW (Phase 6 - LiveTV Module)

**Developer Resources:**
- 📚 API Docs: https://tvheadend.org/projects/tvheadend/wiki/Webui
- 🔗 HTSP Protocol: https://tvheadend.org/projects/tvheadend/wiki/Htsp
- 🔗 GitHub: https://github.com/tvheadend/tvheadend

**Implementation Needs:**
- [ ] `docs/integrations/TVHEADEND.md`
- [ ] `.github/instructions/tvheadend-client.instructions.md`
- [ ] `internal/content/livetv/backend_tvheadend.go`
- [ ] Channel list sync
- [ ] EPG (Electronic Program Guide) sync
- [ ] Recording management
- [ ] Live stream URLs

**API Details:**
- Authentication: Digest or Basic Auth
- HTTP API + HTSP binary protocol
- Key Endpoints: `/api/channel/grid`, `/api/dvr/entry/grid_upcoming`

---

### 8.2 NextPVR

**Purpose:** Windows-based DVR backend
**Integration:** HTTP API
**Priority:** 🟢 LOW (Phase 6 - LiveTV Module)

**Developer Resources:**
- 📚 API Docs: https://github.com/sub3suite/NextPVR/wiki/NextPVR-API
- 🔗 GitHub: https://github.com/sub3suite/NextPVR

**Implementation Needs:**
- [ ] `docs/integrations/NEXTPVR.md`
- [ ] `.github/instructions/nextpvr-client.instructions.md`
- [ ] `internal/content/livetv/backend_nextpvr.go`
- [ ] Channel list
- [ ] EPG sync
- [ ] Recording management

**API Details:**
- Authentication: PIN-based
- Base path: `/service`
- Key Methods: `channel.list`, `recording.list`, `guide.get`

---

## 9. Casting / Streaming Protocols

### 9.1 Google Chromecast

**Purpose:** Cast media to Chromecast devices
**Integration:** Google Cast SDK
**Priority:** 🟡 MEDIUM (Phase 5 - Playback Service)

**Developer Resources:**
- 📚 Cast Docs: https://developers.google.com/cast
- 🔗 Web Receiver: https://developers.google.com/cast/docs/web_receiver
- 🔗 Cast Application Framework (CAF): https://developers.google.com/cast/docs/caf_receiver

**Implementation Needs:**
- [ ] `docs/integrations/CHROMECAST.md` (update CLIENT_SUPPORT.md)
- [ ] `.github/instructions/chromecast-integration.instructions.md`
- [ ] Custom Web Receiver app
- [ ] Media session handling
- [ ] Subtitle support
- [ ] Queue management

**API Details:**
- Google Cast SDK (JavaScript for Web Receiver)
- Load Media API
- Playback control messages

---

### 9.2 DLNA / UPnP

**Purpose:** Stream media to DLNA-capable devices
**Integration:** DLNA/UPnP server implementation
**Priority:** 🟡 MEDIUM (Phase 5 - Playback Service)

**Developer Resources:**
- 📚 DLNA Spec: https://www.dlna.org/
- 🔗 UPnP AV: https://openconnectivity.org/developer/specifications/upnp-resources/upnp
- 🔗 Go Library: https://github.com/anacrolix/dms

**Implementation Needs:**
- [ ] `docs/integrations/DLNA.md` (update CLIENT_SUPPORT.md)
- [ ] `.github/instructions/dlna-server.instructions.md`
- [ ] `internal/service/dlna/server.go`
- [ ] UPnP device discovery (SSDP)
- [ ] Content directory service
- [ ] Media server implementation
- [ ] Transcoding profile selection

**API Details:**
- SOAP-based UPnP services
- SSDP discovery (multicast UDP)
- HTTP streaming

---

## 10. Infrastructure (Already Implemented)

### 10.1 PostgreSQL

**Purpose:** Primary database
**Status:** ✅ Implemented
**Docs:** Needs performance tuning guide

**Implementation Needs:**
- [ ] `docs/PERFORMANCE_TUNING.md` - Indexing strategies, query optimization

---

### 10.2 Dragonfly (Redis-compatible)

**Purpose:** Cache & sessions
**Status:** ✅ Implemented (not registered in main.go)
**Docs:** Needs cache patterns

**Implementation Needs:**
- [ ] `.github/instructions/dragonfly-cache-patterns.instructions.md`

---

### 10.3 Typesense

**Purpose:** Full-text search
**Status:** ✅ Implemented (not registered in main.go)
**Docs:** Needs search architecture

**Implementation Needs:**
- [ ] `docs/SEARCH.md`
- [ ] `.github/instructions/typesense-integration.instructions.md`

---

### 10.4 River

**Purpose:** Job queue (PostgreSQL-native)
**Status:** ✅ Implemented (not registered in main.go)
**Docs:** Needs job queue architecture

**Implementation Needs:**
- [ ] `docs/JOB_QUEUE.md`
- [ ] `.github/instructions/river-job-queue.instructions.md`

---

## 11. Frontend/UX/UI Resources

### 11.1 Design Systems (5 sources)

**Purpose:** Component patterns, design tokens, accessibility standards
**Integration:** Reference documentation (not API)
**Priority:** 🟡 HIGH (Phase 5 - Frontend Architecture)

**Resources:**
- **Material Design 3 (M3 Expressive):** Motion physics, 35 shapes, 14 components, Figma UI Kit
  - https://m3.material.io/
- **Apple Human Interface Guidelines:** iOS/macOS/watchOS/tvOS/visionOS, Generative AI, Hierarchy
  - https://developer.apple.com/design/human-interface-guidelines/
- **Microsoft Fluent 2:** Web (React), iOS, Android, Windows, AI Chat v1
  - https://fluent2.microsoft.design/
- **Atlassian Design System:** Rovo AI Patterns, token-based theming
  - https://atlassian.design/
- **IBM Carbon:** React/Angular/Vue/**Svelte**, Web Components, AI Chat v1
  - https://carbondesignsystem.com/

**Implementation Needs:**
- [ ] `docs/UX_DESIGN_PRINCIPLES.md` - Revenge-specific UX guidelines
- [ ] `.github/instructions/frontend-ux-guidelines.instructions.md`
- [ ] Component naming conventions (from UI Guideline synthesis)
- [ ] Token-based theming system (CSS custom properties)

---

### 11.2 Accessibility & Standards (4 sources)

**Purpose:** WCAG compliance, usability heuristics, international standards
**Integration:** Compliance framework (not API)
**Priority:** 🔴 CRITICAL (Phase 5 - Frontend)

**Resources:**
- **W3C WCAG 2.2:** 4 principles (Perceivable, Operable, Understandable, Robust), 13 guidelines, 3 levels (A/AA/AAA)
  - https://www.w3.org/WAI/WCAG22/quickref/
  - ISO/IEC 40500:2025
- **ISO 9241-11:2018:** Usability as Outcome of Use, Effectiveness/Efficiency/Satisfaction
  - https://www.iso.org/standard/63500.html
- **Nielsen Norman 10 Usability Heuristics:** Jakob Nielsen 1994 (updated 2020), free posters
  - https://www.nngroup.com/articles/ten-usability-heuristics/
- **Laws of UX:** 26 laws (Fitts's, Hick's, Miller's, Jakob's), O'Reilly book
  - https://lawsofux.com/

**Implementation Needs:**
- [ ] `docs/ACCESSIBILITY_IMPLEMENTATION.md`
- [ ] `.github/instructions/wcag-compliance.instructions.md`
- [ ] `.github/instructions/accessibility-standards.instructions.md`
- [ ] Automated testing setup (axe-core, pa11y)
- [ ] Manual testing checklist (keyboard nav, screen readers)

---

### 11.3 Component Libraries & Patterns (3 sources)

**Purpose:** Standardized UI components, naming, anatomy, best practices
**Integration:** Reference documentation
**Priority:** 🟡 HIGH (Phase 5 - Frontend)

**Resources:**
- **UI Guideline:** 60+ components synthesized from top 20 systems
  - https://uiguideline.com/
  - Methodology: Annual top 20 selection → manual review → consolidate patterns → synthesize definition
  - Components: Atomic (Badge, Avatar), Buttons, Controls (Slider, Switch), Data Display (Table, Card), Inputs (Date Picker, Select), Loaders, Navigation (Breadcrumbs, Tabs), Notifications, Overlays (Modal, Tooltip), States
  - Each component: Possible names, properties, HTML structure, anatomy, best practices from 20 systems
  - 6,700+ users, trusted by Amsterdam Uni, Asana, Atlassian, Feedly, Zendesk
- **Open UI (W3C Community Group):** UI standardization, component research
  - https://open-ui.org/
  - GitHub: https://github.com/openui/open-ui (4.4k stars, 98 contributors)
  - Purpose: Allow web devs to style/extend built-in web UI controls (<select>, checkboxes, radio, date pickers)
  - Scope: Research component patterns from 3rd-party frameworks, capture common language, define developer needs, recommend to WHATWG/CSSWG/W3C/TC39
  - Components researched: Customizable Select, Exclusive Accordion, Invoker Commands, Popover, Combobox, Switch, Breadcrumb, Menu, Tabs, Tooltip, and more
  - Design Principles, Component Matrix, Test suites
- **web.dev Patterns:** Modern CSS/JS patterns with accessibility
  - https://web.dev/patterns/
  - Categories: Animation, Clipboard, Components, Files, Layout (CSS Grid/Flexbox), Media, Theming, Web Apps
  - Focus: Modern CSS APIs, prefers-reduced-motion, user preferences

**Implementation Needs:**
- [ ] `docs/FRONTEND_COMPONENTS.md` - Component library structure
- [ ] `.github/instructions/svelte-component-patterns.instructions.md`
- [ ] `.github/instructions/shadcn-svelte-customization.instructions.md`
- [ ] Component naming standard (based on UI Guideline synthesis)
- [ ] Anatomy definitions (props, slots, events)

---

### 11.4 UX Research & Best Practices (5 sources)

**Purpose:** UX methodology, design thinking, e-commerce patterns
**Integration:** Reference documentation
**Priority:** 🟡 MEDIUM (Phase 5 - Frontend)

**Resources:**
- **A List Apart (est. 1998):** Web standards authority, progressive enhancement
  - https://alistapart.com/
  - ISSN: 1534-0295
  - Latest: "Design for Amiability" (Vienna Circle CS history), "Design Dialects" (systems as living languages), "Shared Design Leadership", "Beta to Bedrock" (long-term stability), "User Research as Storytelling", "Mobile-First CSS Rethink"
  - Topics: Web standards, accessibility, design thinking, career
- **Smashing Magazine:** UX design articles, patterns, accessibility
  - https://www.smashingmagazine.com/category/ux-design/
  - 57+ UX articles (homepage accessible, category page 404)
  - Topics: Design patterns, accessibility, mobile-first, infinite scroll
- **Baymard Institute:** E-commerce UX research (200k+ hours)
  - https://baymard.com/
  - 443 articles, 18k+ annotated examples, UX Benchmark (326 sites, 275k scores)
  - Latest: Mobile UX 2025 (9 pitfalls), Checkout UX 2025 (10 best practices), Homepage/Category Nav (67% mediocre)
  - Topics: Homepage/Category (45), Search (27), Product List (52), Product Page (51), Cart/Checkout (70), Accessibility (8), Mobile (128)
- **Interaction Design Foundation:** Largest UX community (1.2M+ enrollments)
  - https://www.interaction-design.org/
  - Courses: UX Beginner's Guide, Get First Job as UX Designer, User Research Methods
  - Topics: What is UX, Graphic→UX career, Wireframing tools, User-Centered Design
- **UK GDS 11 Principles:** Government design standards
  - https://www.gov.uk/guidance/government-design-principles
  - Principles: Start with user needs, Do less, Design with data, Do hard work to make simple, Iterate, This is for everyone, Understand context, Build digital services not websites, Be consistent not uniform, Make things open, **NEW 2025**: Minimise environmental impact

**Implementation Needs:**
- [ ] `docs/UX_DESIGN_PRINCIPLES.md` (update with research findings)
- [ ] `.github/instructions/frontend-ux-guidelines.instructions.md` (update)
- [ ] Progressive enhancement strategy
- [ ] Mobile-first responsive design
- [ ] User research methodology (optional)

---

### 11.5 Media Player Best Practices (3 sources)

**Purpose:** Video/audio player UX, accessibility, performance
**Integration:** Reference + potential library adoption
**Priority:** 🟡 HIGH (Phase 5 - Player Architecture)

**Resources:**
- **Plyr:** Simple HTML5/YouTube/Vimeo player
  - https://plyr.io/ (GitHub: 29.5k stars)
  - Features: Accessible (VTT captions, screen readers), customizable, clean HTML (<input type="range">, <progress>, <button>), responsive, monetization, streaming (hls.js, Shaka, dash.js), API, events, fullscreen, PiP, keyboard shortcuts, speed controls, multiple captions, i18n, preview thumbnails
  - CSS custom properties (40+ tokens)
  - Used by: Selz, Peugeot, TomTom, BitChute, Koel
  - Supports: HTML5 Video/Audio, YouTube, Vimeo
- **Video.js:** Industry standard video player
  - https://videojs.com/
  - Advanced plugin ecosystem, HLS/DASH support
- **W3C Media Events:** HTML5 media element API reference
  - https://www.w3.org/2010/05/video/mediaevents.html
  - Event model: loadstart, loadeddata, canplay, canplaythrough, playing, pause, timeupdate, volumechange, seeking, seeked, ended, etc.

**Implementation Needs:**
- [ ] `docs/PLAYER_UX.md` - Video/audio player architecture
- [ ] `.github/instructions/video-player-patterns.instructions.md`
- [ ] `.github/instructions/audio-player-patterns.instructions.md`
- [ ] Controls UX (play/pause, seek, volume, fullscreen, quality, subtitles)
- [ ] Accessibility (WCAG 2.2 for media, keyboard nav, screen reader support)
- [ ] Performance (adaptive bitrate, preloading, buffering)
- [ ] Format support (HLS, DASH, WebM, MP4)

---

### 11.6 External Client Integration (2 sources)

**Purpose:** Integration with existing media center clients
**Integration:** API compatibility layers
**Priority:** 🟢 LOW (Phase 9 - External Services)

**Resources:**
- **Kodi JSON-RPC API:** JSON-RPC interface for Kodi control
  - https://kodi.wiki/view/JSON-RPC_API
  - Transports: Python (addon), HTTP (POST/GET), TCP (port 9090), WebSocket
  - API Version: v12 (Matrix), v13 (Nexus), v14 (Omega)
  - Functionalities: Response, Notifications (server/client), File download
  - Methods: Player, Playlist, Files, AudioLibrary, VideoLibrary, GUI, System, JSONRPC
  - Authentication: API key (HTTP), none (local Python/TCP)
  - Self-documented via JSONRPC.Introspect (JSON schema)
- **Jellyfin API:** (Documentation 404 errors, check GitHub)
  - https://api.jellyfin.org/ (Failed to extract)
  - Alternative: Check https://github.com/jellyfin/jellyfin for OpenAPI spec
- **Plex API:** (No public documentation found)
  - https://www.plex.tv/api/ (404)
  - Community libraries exist (unofficial)

**Implementation Needs:**
- [ ] `docs/EXTERNAL_CLIENTS.md` - Client integration patterns
- [ ] `.github/instructions/kodi-jsonrpc-compatibility.instructions.md`
- [ ] Kodi JSON-RPC compatibility layer (optional)
- [ ] Jellyfin API compatibility (optional)
- [ ] Device detection & client capabilities
- [ ] Remote control API (play/pause/seek from external apps)

---

### 11.7 Database Design Best Practices

**Purpose:** Advanced PostgreSQL patterns, indexing, performance
**Integration:** Implementation guidelines
**Priority:** 🔴 CRITICAL (Phase 1-2 - Database Architecture)

**Resources:**
- **PostgreSQL 18 DDL Documentation:**
  - https://www.postgresql.org/docs/current/ddl.html
  - Chapters: Table Basics, Default Values, Identity Columns, Generated Columns, Constraints (Check, Not-Null, Unique, Primary Keys, Foreign Keys, Exclusion), System Columns, Modifying Tables, Privileges, Row Security Policies, Schemas, Inheritance, **Table Partitioning** (Declarative, Partition Pruning, Best Practices), Foreign Data, Dependency Tracking

**Implementation Needs:**
- [ ] `docs/DATABASE_DESIGN_PRINCIPLES.md` - Advanced patterns
- [ ] `.github/instructions/postgresql-performance.instructions.md`
- [ ] Indexing strategies (B-tree, GiST, GIN for JSONB, partial indexes)
- [ ] Query optimization (EXPLAIN ANALYZE)
- [ ] Partitioning strategies (range, list, hash)
- [ ] Migration best practices (zero-downtime, rollback)
- [ ] Connection pooling (pgBouncer, pgx pool)
- [ ] Replication & backup
- [ ] Performance monitoring (pg_stat_statements)
- [ ] Schema isolation (adult content `c` schema)

---

## Summary Counts

| Category | Services | Docs | Instructions | Priority |
|----------|----------|------|--------------|----------|
| Content Management | 5 | 0/5 | 0/5 | 🔴 CRITICAL |
| Audiobook Management | 1 | 0/1 | 0/1 | 🟡 HIGH |
| Request Management | 2 | 0/2 | 0/2 | 🟡 MEDIUM |
| Metadata Providers | 15 | 1/15 | 0/15 | 🔴 CRITICAL |
| Scrobbling/Sync | 5 | 0/5 | 0/5 | 🟡 HIGH |
| **Anime Tracking** | **3** | **0/3** | **0/3** | **🟡 HIGH** |
| Authentication (OIDC) | 4 | 0/4 | 0/4 | 🟡 MEDIUM |
| Transcoding | 1 | 1/1 | 0/1 | 🔴 CRITICAL |
| Live TV | 2 | 0/2 | 0/2 | 🟢 LOW |
| Casting | 2 | 0/2 | 0/2 | 🟡 MEDIUM |
| Infrastructure | 4 | 0/4 | 1/4 | 🔴 CRITICAL |
| **Frontend/UX/UI** | **22** | **0/22** | **0/11** | **🔴 CRITICAL** |

**TOTAL:** 66 External Integrations (41 APIs + 3 Anime + 22 UX/UI Resources)
**Documentation Needed:** 62 docs, 50 instructions
**Existing:** ADULT_METADATA.md (partial), OFFLOADING.md (partial)

---

## Next Steps

1. **Create base patterns first:**
   - `external-api-client-pattern.instructions.md`
   - `rate-limiting.instructions.md`
   - `error-handling-patterns.instructions.md`

2. **Critical integrations for MVP:**
   - Servarr (Radarr, Sonarr, Lidarr)
   - TMDb, TheTVDB, MusicBrainz
   - Blackbeard (transcoding)
   - Infrastructure (River, Typesense, Dragonfly)
   - **Anime Tracking (AniList, MyAnimeList, Kitsu)**

3. **Frontend/UX resources:**
   - Design systems reference (Material 3, Fluent 2, Carbon Svelte)
   - Accessibility standards (WCAG 2.2 AA, Nielsen heuristics)
   - Component patterns (UI Guideline, Open UI)
   - UX research (A List Apart, Baymard)
   - Media player UX (Plyr patterns, W3C events)
   - Database design (PostgreSQL 18 DDL, partitioning)

4. **Extended integrations:**
   - Scrobbling services
   - Audiobookshelf
   - Adult metadata (Stash, ThePornDB)
   - OIDC providers
   - External clients (Kodi, Jellyfin compatibility)


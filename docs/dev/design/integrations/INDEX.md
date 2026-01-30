# External Integrations

> Third-party services and APIs

---

## Overview

Revenge integrates with numerous external services for:
- **Metadata** - Media information enrichment
- **Scrobbling** - Activity tracking
- **Authentication** - SSO/OIDC providers
- **Infrastructure** - Core stack components
- **Media Management** - Servarr stack

---

## Categories

### 🎬 [Metadata Providers](metadata/INDEX.md)
External sources for media metadata.

| Category | Providers |
|----------|-----------|
| [Video](metadata/video/INDEX.md) | TMDB, TVDB, OMDB, Fanart.tv |
| [Music](metadata/music/INDEX.md) | MusicBrainz, Last.fm, Spotify, Discogs |
| [Books](metadata/books/INDEX.md) | Open Library, Google Books, Goodreads |
| [Comics](metadata/comics/INDEX.md) | ComicVine, Marvel API, GCD |
| [Adult](metadata/adult/INDEX.md) | StashDB, TPDB, FreeOnes |

### 📺 [Anime](anime/INDEX.md)
Anime-specific metadata and tracking.

| Provider | Purpose |
|----------|---------|
| AniList | Primary anime database |
| MyAnimeList | Secondary, largest community |
| Kitsu | Alternative with good API |

### 📊 [Scrobbling](scrobbling/INDEX.md)
Activity tracking and sync services.

| Provider | Content |
|----------|---------|
| Trakt | Movies, TV |
| Last.fm | Music |
| ListenBrainz | Music (open source) |
| Letterboxd | Movies |
| Simkl | Movies, TV, Anime |

### 🔑 [Authentication](auth/INDEX.md)
OIDC/SSO providers.

| Provider | Type |
|----------|------|
| Authelia | Self-hosted SSO |
| Authentik | Self-hosted IdP |
| Keycloak | Enterprise IdP |
| Generic OIDC | Any provider |

### 📚 [Wiki](wiki/INDEX.md)
Supplementary information sources.

| Provider | Purpose |
|----------|---------|
| Wikipedia | General info |
| Wikidata | Structured data |
| [Adult Wiki](wiki/adult/INDEX.md) | Performer info |

### 🗄️ [Servarr](servarr/INDEX.md)
Media management automation.

| Application | Content |
|-------------|---------|
| Radarr | Movies |
| Sonarr | TV Shows |
| Lidarr | Music |
| Readarr | Books |
| Whisparr | Adult |

### 🔗 [External Services](external/INDEX.md)
Third-party integrations.

| Category | Services |
|----------|----------|
| [Adult](external/adult/INDEX.md) | Twitter/X, Instagram |

### 🎧 [Audiobook & Podcast](audiobook/INDEX.md)
Native audiobook and podcast management.

| Feature | Implementation |
|---------|----------------|
| Audiobooks | Native library + metadata |
| Podcasts | Native RSS + downloads |

### 🔄 [Transcoding](transcoding/INDEX.md)
External transcoding services.

| Service | Purpose |
|---------|---------|
| Blackbeard | Video transcoding |

### 📺 [Live TV](livetv/INDEX.md)
PVR backend integration.

| Provider | Purpose |
|----------|---------|
| TVHeadend | Full PVR |
| NextPVR | Windows PVR |

### 📡 [Casting](casting/INDEX.md)
Device casting protocols.

| Protocol | Devices |
|----------|---------|
| Chromecast | Google devices |
| DLNA | Smart TVs, consoles |

### 🏗️ [Infrastructure](infrastructure/INDEX.md)
Core stack components.

| Component | Purpose |
|-----------|---------|
| PostgreSQL | Primary database |
| Dragonfly | Cache & sessions |
| Typesense | Search engine |
| River | Job queue |

---

## Integration Status

| Status | Meaning |
|--------|---------|
| 🟢 | Active / Supported |
| 🟡 | Planned / In Progress |
| 🔴 | Low Priority |

---

## Configuration

All integrations configured in `configs/config.yaml`:

```yaml
# Enable/disable categories
metadata:
  video:
    tmdb:
      enabled: true
      api_key: "${TMDB_API_KEY}"

scrobbling:
  trakt:
    enabled: true
    client_id: "${TRAKT_CLIENT_ID}"

auth:
  oidc:
    enabled: true
```

---

## Related Documentation

- [Architecture](../architecture/ARCHITECTURE.md)
- [Tech Stack](../technical/TECH_STACK.md)
- [Setup Guide](../operations/SETUP.md)

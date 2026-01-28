# Scrobbling Services

> Track and share playback activity

---

## Overview

Scrobbling services allow users to:
- Track watch/listen history externally
- Sync progress across platforms
- Share activity with friends
- Get recommendations based on history
- Maintain unified viewing statistics

---

## Providers

| Provider | Content | API | Status |
|----------|---------|-----|--------|
| [Trakt](TRAKT.md) | Movies, TV | REST/OAuth | 🟢 Primary |
| [Last.fm](LASTFM_SCROBBLE.md) | Music | REST | 🟢 Primary |
| [ListenBrainz](LISTENBRAINZ.md) | Music | REST | 🟡 Alternative |
| [Letterboxd](LETTERBOXD.md) | Movies | CSV Export | 🟡 Limited |
| [Simkl](SIMKL.md) | Movies, TV, Anime | REST/OAuth | 🟡 Secondary |

---

## Provider Details

### Trakt
**Primary for video content**

- ✅ Movies and TV shows
- ✅ Watch history sync
- ✅ Ratings sync
- ✅ Lists and collections
- ✅ VIP features (calendar, progress)

### Last.fm
**Primary for music scrobbling**

- ✅ Track scrobbling
- ✅ Now playing
- ✅ Love/unlove tracks
- ✅ Long history (since 2002)

### ListenBrainz
**Open-source Last.fm alternative**

- ✅ Track scrobbling
- ✅ Open data
- ✅ MusicBrainz integration
- ✅ No account limits

### Letterboxd
**Film diary and social**

- ✅ Film diary entries
- ✅ Reviews and ratings
- ⚠️ No write API (export only)

### Simkl
**Multi-content tracker**

- ✅ Movies, TV, Anime
- ✅ Watch history
- ✅ Ratings
- ✅ Anime-specific features

---

## Scrobble Flow

```
User plays content
    ↓
Playback service detects progress
    ↓
At threshold (e.g., 80%), trigger scrobble
    ↓
Queue scrobble job (River)
    ↓
Send to enabled services (parallel)
    ↓
Handle failures with retry
```

---

## Configuration

```yaml
scrobbling:
  # Global settings
  threshold: 0.8  # 80% completion

  # Per-service
  trakt:
    enabled: true
    client_id: "${TRAKT_CLIENT_ID}"
    client_secret: "${TRAKT_CLIENT_SECRET}"

  lastfm:
    enabled: true
    api_key: "${LASTFM_API_KEY}"
    api_secret: "${LASTFM_API_SECRET}"

  listenbrainz:
    enabled: false
    user_token: "${LISTENBRAINZ_TOKEN}"
```

---

## User Authentication

Most services require per-user OAuth:

```
User → Settings → Connect Trakt
    ↓
Redirect to Trakt OAuth
    ↓
User authorizes
    ↓
Callback with code
    ↓
Exchange for tokens
    ↓
Store encrypted tokens
```

---

## Related Documentation

- [Metadata Providers](../metadata/INDEX.md)
- [External Services](../external/INDEX.md)

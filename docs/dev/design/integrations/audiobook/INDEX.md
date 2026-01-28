# Audiobook Providers

> Audiobook library management and metadata

---

## Overview

Audiobook integration provides:
- Audiobookshelf library sync
- Chapter information
- Narrator metadata
- Progress tracking
- Podcast support (via Audiobookshelf)

---

## Providers

| Provider | Type | Status |
|----------|------|--------|
| [Audiobookshelf](AUDIOBOOKSHELF.md) | Library Server | 🟢 Primary |

---

## Provider Details

### Audiobookshelf
**Self-hosted audiobook server**

- ✅ Audiobook management
- ✅ Podcast support
- ✅ Progress sync
- ✅ Multi-user support
- ✅ Mobile apps (iOS/Android)
- ✅ REST API + Socket.io

---

## Integration Modes

### Standalone
Audiobookshelf manages audiobook library independently, Revenge links for playback.

### Sync Mode
Progress and metadata synced between Revenge and Audiobookshelf.

```yaml
audiobook:
  provider: audiobookshelf
  url: "http://audiobookshelf:13378"
  api_key: "${AUDIOBOOKSHELF_API_KEY}"

  sync:
    progress: true
    metadata: true
```

---

## Data Flow

```
Audiobookshelf manages audiobooks
    ↓
Revenge discovers via API
    ↓
User plays in Revenge
    ↓
Progress synced back to Audiobookshelf
    ↓
Mobile app shows updated progress
```

---

## Configuration

```yaml
integrations:
  audiobook:
    enabled: true
    provider: audiobookshelf

    audiobookshelf:
      url: "http://audiobookshelf:13378"
      api_key: "${AUDIOBOOKSHELF_API_KEY}"

      # Library mapping
      libraries:
        - name: "Audiobooks"
          abs_library_id: "lib_xyz"

      # Sync settings
      sync:
        progress: true
        interval: "5m"
```

---

## Related Documentation

- [Book Metadata](../metadata/books/INDEX.md)
- [Readarr Integration](../servarr/READARR.md)

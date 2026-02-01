## Table of Contents

- [Metadata Service](#metadata-service)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [Related Pages](#related-pages)
    - [Learn More](#learn-more)

# Metadata Service




> Unified metadata service with Arr-first priority chain

Metadata service orchestrates content enrichment using PRIMARY sources (Arr services - Radarr, Sonarr, Lidarr, Chaptarr, Whisparr) that aggregate metadata locally, with SUPPLEMENTARY external APIs (TMDb, TheTVDB, MusicBrainz, etc.) as fallback and enrichment. Priority chain: L1 cache → L2 cache → Arr services (PRIMARY) → External APIs (SUPPLEMENTARY via optional proxy/VPN). Automatic filename matching, manual override, and scheduled background refresh jobs keep metadata current.

---





---


## How It Works

<!-- How it works -->




## Features
<!-- Feature list placeholder -->



## Configuration
<!-- User-friendly configuration guide -->









## Related Documentation
### Related Pages
<!-- Related wiki pages -->

### Learn More

Official documentation and guides:
- [Uber fx](../../sources/tooling/fx.md)
- [Last.fm API](../../sources/apis/lastfm.md)
- [River Job Queue](../../sources/tooling/river.md)



---

**Need Help?** [Open an issue](https://github.com/revenge-project/revenge/issues) or [Join the discussion](https://github.com/revenge-project/revenge/discussions)
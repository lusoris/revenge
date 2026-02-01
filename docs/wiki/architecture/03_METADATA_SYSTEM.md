## Table of Contents

- [Revenge - Metadata System](#revenge-metadata-system)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [Related Pages](#related-pages)
    - [Learn More](#learn-more)



---
---

# Revenge - Metadata System




> How Revenge finds and stores information about your media

The metadata system gathers information about your media from multiple sources. It always checks local cache first for instant display, then queries Arr services (Radarr, Sonarr) which already have metadata, then external APIs like TMDb or MusicBrainz. Background jobs enrich media with additional data like cast info, thumbnails, and blurhash previews. Two-tier caching (memory + distributed) ensures fast lookups even for large libraries.

---




## Contents

<!-- TOC will be auto-generated here by markdown-toc -->

---


## How It Works

<!-- User-friendly explanation -->




## Features
<!-- Feature list placeholder -->



## Configuration
<!-- User-friendly configuration guide -->









## Related Documentation
### Related Pages
<!-- Related wiki pages -->

### Learn More

Official documentation and guides:
- [Dragonfly Documentation](../../sources/infrastructure/dragonfly.md)
- [go-blurhash](../../sources/media/go-blurhash.md)
- [Last.fm API](../../sources/apis/lastfm.md)
- [pgx PostgreSQL Driver](../../sources/database/pgx.md)
- [PostgreSQL Arrays](../../sources/database/postgresql-arrays.md)
- [PostgreSQL JSON Functions](../../sources/database/postgresql-json.md)
- [River Job Queue](../../sources/tooling/river.md)
- [rueidis](../../sources/tooling/rueidis.md)
- [rueidis GitHub README](../../sources/tooling/rueidis-guide.md)



---

**Need Help?** [Open an issue](https://github.com/revenge-project/revenge/issues) or [Join the discussion](https://github.com/revenge-project/revenge/discussions)
## Table of Contents

- [EPG Service](#epg-service)
- [EPG Service](#epg-service)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [See Also](#see-also)



---
sources:
  - name: XMLTV Format
    url: http://wiki.xmltv.org/index.php/XMLTVFormat
    note: EPG data format standard
  - name: Typesense Go Client
    url: ../../sources/infrastructure/typesense-go.md
    note: Full-text search for programs
  - name: River Job Queue
    url: ../../sources/tooling/river.md
    note: Scheduled EPG refresh jobs
  - name: Uber fx
    url: ../../sources/tooling/fx.md
    note: Dependency injection
design_refs:
  - title: services
    path: INDEX.md
  - title: 01_ARCHITECTURE
    path: ../architecture/01_ARCHITECTURE.md
  - title: LIVE_TV_DVR
    path: ../features/livetv/LIVE_TV_DVR.md
  - title: TVHEADEND
    path: ../integrations/livetv/TVHEADEND.md
  - title: NEXTPVR
    path: ../integrations/livetv/NEXTPVR.md
  - title: ERSATZTV
    path: ../integrations/livetv/ERSATZTV.md
---

# EPG Service




# EPG Service

> TV program schedules and guide data for Live TV and DVR

The Electronic Program Guide (EPG) Service provides comprehensive TV program schedule information for Live TV and DVR functionality. Fetches guide data from multiple sources in XMLTV format, indexes programs for fast search, and delivers schedule information via REST API. Automatic refresh keeps guide data current. Supports filtering by channel, genre, time range, and full-text search across program titles and descriptions.

---




## Contents

<!-- TOC will be auto-generated here by markdown-toc -->

---


## How It Works

<!-- How it works -->




## Features
<!-- Feature list placeholder -->



## Configuration
<!-- User-friendly configuration guide -->









## Related Documentation
### See Also
<!-- Related wiki pages -->



---

**Need Help?** [Open an issue](https://github.com/revenge-project/revenge/issues) or [Join the discussion](https://github.com/revenge-project/revenge/discussions)
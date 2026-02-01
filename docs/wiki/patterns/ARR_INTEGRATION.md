## Table of Contents

- [Arr Integration Pattern](#arr-integration-pattern)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [See Also](#see-also)



---
sources:
  - name: Radarr API Docs
    url: ../../sources/apis/radarr-docs.md
    note: Radarr webhook events
  - name: Sonarr API Docs
    url: ../../sources/apis/sonarr-docs.md
    note: Sonarr webhook events
  - name: Lidarr API Docs
    url: ../../sources/apis/lidarr-docs.md
    note: Lidarr webhook events
  - name: Servarr Wiki
    url: ../../sources/apis/servarr-wiki.md
    note: Shared Arr stack documentation
design_refs:
  - title: patterns
    path: patterns.md
  - title: 01_ARCHITECTURE
    path: ../architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: ../architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: ../architecture/03_METADATA_SYSTEM.md
---

# Arr Integration Pattern




> Seamless integration with your Arr stack for automated media management

The Arr Integration Pattern provides a consistent approach for integrating with Radarr, Sonarr, Lidarr, and Whisparr. Webhook events automatically sync new downloads, metadata updates, and deletions. The system prioritizes Arr metadata (which is already cached from upstream sources) to minimize external API calls, while background jobs enrich content with additional metadata when needed.

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
### See Also
<!-- Related wiki pages -->



---

**Need Help?** [Open an issue](https://github.com/revenge-project/revenge/issues) or [Join the discussion](https://github.com/revenge-project/revenge/discussions)
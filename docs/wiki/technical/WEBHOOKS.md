## Table of Contents

- [Webhooks](#webhooks)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [See Also](#see-also)



---
sources:
  - name: Radarr API Docs
    url: ../sources/apis/radarr-docs.md
    note: Auto-resolved from radarr-docs
  - name: Sonarr API Docs
    url: ../sources/apis/sonarr-docs.md
    note: Auto-resolved from sonarr-docs
  - name: Lidarr API Docs
    url: ../sources/apis/lidarr-docs.md
    note: Auto-resolved from lidarr-docs
  - name: River Job Queue
    url: ../sources/tooling/river.md
    note: Auto-resolved from river
  - name: Uber fx
    url: ../sources/tooling/fx.md
    note: Auto-resolved from fx
design_refs:
  - title: technical
    path: technical.md
  - title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
  - title: WEBHOOK_PATTERNS
    path: patterns/WEBHOOK_PATTERNS.md
  - title: ARR_INTEGRATION
    path: patterns/ARR_INTEGRATION.md
  - title: NOTIFICATIONS
    path: technical/NOTIFICATIONS.md
---

# Webhooks




> Real-time event integration with external services

The Webhook System enables bidirectional event-driven integration. Receive webhooks from Arr services (Radarr, Sonarr, Lidarr, Whisparr) to automatically sync new downloads and metadata updates. Send webhooks to Discord, Slack, or custom endpoints for notifications and automation. All webhooks use HMAC signatures for security and are queued via River for reliable delivery with automatic retries.

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
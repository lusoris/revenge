## Table of Contents

- [Notifications System](#notifications-system)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [See Also](#see-also)



---
sources:
  - name: go-mail GitHub README
    url: ../../sources/tooling/go-mail-guide.md
    note: Auto-resolved from go-mail
  - name: go-fcm
    url: ../../sources/tooling/go-fcm.md
    note: FCM push notifications
  - name: River Job Queue
    url: ../../sources/tooling/river.md
    note: Auto-resolved from river
  - name: Uber fx
    url: ../../sources/tooling/fx.md
    note: Auto-resolved from fx
design_refs:
  - title: technical
    path: technical.md
  - title: 01_ARCHITECTURE
    path: ../architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: ../architecture/02_DESIGN_PRINCIPLES.md
  - title: EMAIL
    path: ../technical/EMAIL.md
  - title: WEBHOOKS
    path: ../technical/WEBHOOKS.md
---

# Notifications System




> Stay informed with email, push, and webhook notifications

The Notifications System delivers real-time alerts across multiple channels. Users receive email notifications for account updates, push notifications on mobile devices via Firebase Cloud Messaging, and developers can subscribe to webhooks for automation. The system queues all notifications via River for reliable delivery with automatic retries and respects user preferences for each notification type.

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
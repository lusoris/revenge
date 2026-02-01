## Table of Contents

- [Email System](#email-system)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [See Also](#see-also)

---
sources:
- name: go-mail GitHub README
    url: ../sources/tooling/go-mail-guide.md
    note: Auto-resolved from go-mail-docs
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
- title: NOTIFICATIONS
    path: technical/NOTIFICATIONS.md
---

# Email System

> Reliable email delivery for notifications, alerts, and user communications

The Email System provides reliable SMTP email delivery with support for HTML templates, TLS encryption, and async sending. Built on go-mail library with connection pooling and automatic retries. All emails are queued via River for resilient delivery, with bounce detection and unsubscribe management built-in.

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

## Table of Contents

- [Fingerprint Service](#fingerprint-service)
- [Fingerprint Service](#fingerprint-service)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [See Also](#see-also)

---
sources:
- name: FFmpeg Documentation
    url: ../sources/media/ffmpeg.md
    note: Auto-resolved from ffmpeg
- name: FFmpeg Codecs
    url: ../sources/media/ffmpeg-codecs.md
    note: Auto-resolved from ffmpeg-codecs
- name: FFmpeg Formats
    url: ../sources/media/ffmpeg-formats.md
    note: Auto-resolved from ffmpeg-formats
- name: Uber fx
    url: ../sources/tooling/fx.md
    note: Auto-resolved from fx
- name: go-astiav (FFmpeg bindings)
    url: ../sources/media/go-astiav.md
    note: Auto-resolved from go-astiav
- name: go-astiav GitHub README
    url: ../sources/media/go-astiav-guide.md
    note: Auto-resolved from go-astiav-docs
- name: pgx PostgreSQL Driver
    url: ../sources/database/pgx.md
    note: Auto-resolved from pgx
- name: PostgreSQL Arrays
    url: ../sources/database/postgresql-arrays.md
    note: Auto-resolved from postgresql-arrays
- name: PostgreSQL JSON Functions
    url: ../sources/database/postgresql-json.md
    note: Auto-resolved from postgresql-json
- name: River Job Queue
    url: ../sources/tooling/river.md
    note: Auto-resolved from river
design_refs:
- title: services
    path: services/INDEX.md
- title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
- title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
- title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# Fingerprint Service

# Fingerprint Service

> Media file identification via perceptual hashing and acoustic fingerprinting

Fingerprinting identifies media files by their content, not just filenames. Perceptual hashes match images and video frames even after re-encoding. Acoustic fingerprints identify music tracks. This enables duplicate detection, automatic metadata matching, and skip intro detection. Fingerprints are computed on library scan and stored for fast lookups.

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

## Table of Contents

- [Advanced Offloading Architecture](#advanced-offloading-architecture)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [Related Pages](#related-pages)
    - [Learn More](#learn-more)



---
---

# Advanced Offloading Architecture




> Keep your API fast by offloading heavy work to specialized services

The Advanced Offloading Architecture ensures fast API response times by delegating heavy operations to background workers and specialized services. Transcoding, metadata enrichment, and image processing run asynchronously via River queue. Session storage and rate limiting use Dragonfly cache. Full-text search queries route to Typesense. Metrics collection offloads to Prometheus. The result is a responsive API that scales horizontally while handling resource-intensive tasks efficiently.

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
- [Uber fx](../../sources/tooling/fx.md)
- [koanf](../../sources/tooling/koanf.md)
- [Prometheus Go Client](../../sources/observability/prometheus.md)
- [Prometheus Metric Types](../../sources/observability/prometheus-metrics.md)
- [River Job Queue](../../sources/tooling/river.md)
- [rueidis](../../sources/tooling/rueidis.md)
- [rueidis GitHub README](../../sources/tooling/rueidis-guide.md)
- [Typesense API](../../sources/infrastructure/typesense.md)
- [Typesense Go Client](../../sources/infrastructure/typesense-go.md)



---

**Need Help?** [Open an issue](https://github.com/revenge-project/revenge/issues) or [Join the discussion](https://github.com/revenge-project/revenge/discussions)
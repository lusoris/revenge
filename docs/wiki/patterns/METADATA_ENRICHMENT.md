## Table of Contents

- [Metadata Enrichment Pattern](#metadata-enrichment-pattern)
  - [Features](#features)
  - [Related Documentation](#related-documentation)
    - [Related Pages](#related-pages)
    - [Learn More](#learn-more)

# Metadata Enrichment Pattern




> Fast, efficient metadata with intelligent caching and background enrichment

The Metadata Enrichment Pattern ensures fast UI response times while maintaining rich, up-to-date metadata. The system uses a five-tier priority chain starting with local cache for instant results, falling back to Arr services (which cache upstream data), then internal sources, external APIs, and finally background enrichment jobs. Multi-tier caching with request coalescing prevents duplicate API calls and reduces external API quota consumption.

---






## Features
<!-- Feature list placeholder -->
## Related Documentation
### Related Pages
<!-- Related wiki pages -->

### Learn More

Official documentation and guides:
- [River Job Queue](../../sources/tooling/river.md)
- [rueidis](../../sources/tooling/rueidis.md)
- [Otter](https://pkg.go.dev/github.com/maypok86/otter)
- [Sturdyc](../../sources/tooling/sturdyc-guide.md)



---

**Need Help?** [Open an issue](https://github.com/revenge-project/revenge/issues) or [Join the discussion](https://github.com/revenge-project/revenge/discussions)
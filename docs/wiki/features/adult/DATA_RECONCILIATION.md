## Table of Contents

- [Adult Data Reconciliation](#adult-data-reconciliation)
  - [How It Works](#how-it-works)
    - [Content Flow](#content-flow)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [Related Pages](#related-pages)
    - [Learn More](#learn-more)

# Adult Data Reconciliation




> Smart matching when metadata sources disagree

When multiple metadata sources provide conflicting information, the reconciliation system resolves differences. Uses fuzzy string matching and confidence scoring to pick the best data. Manually override any field when automatic matching is wrong. Track which source provided each piece of metadata. Particularly useful for adult content where naming is inconsistent.

---





---


## How It Works

<!-- Wiki how it works description pending -->

### Content Flow

1. **Add Content**: Import or add manually
2. **Metadata**: Automatically fetched
3. **Library**: Organized by folder structure
4. **Playback**: Stream directly or use external player




## Features
<!-- Feature list placeholder -->




## Configuration









## Related Documentation
### Related Pages
<!-- Related wiki pages -->

### Learn More

Official documentation and guides:
- [Casbin](../../../sources/security/casbin.md)
- [River Job Queue](../../../sources/tooling/river.md)
- [sqlc](../../../sources/database/sqlc.md)
- [sqlc Configuration](../../../sources/database/sqlc-config.md)



---

**Need Help?** [Open an issue](https://github.com/revenge-project/revenge/issues) or [Join the discussion](https://github.com/revenge-project/revenge/discussions)
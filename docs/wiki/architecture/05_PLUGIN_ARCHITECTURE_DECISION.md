## Table of Contents

- [Plugin Architecture Decision](#plugin-architecture-decision)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [Related Pages](#related-pages)
    - [Learn More](#learn-more)

# Plugin Architecture Decision




> Why Revenge uses integrations instead of plugins

Revenge deliberately chose not to implement a plugin system. Instead, common integrations (Radarr, Sonarr, TMDb, etc.) are built directly into the codebase with first-class support. This means faster development, better security (no arbitrary code execution), and simpler maintenance. External systems can still integrate via webhooks and the REST API. For power users who need custom automation, scripting support (Lua or Starlark) may be added in the future.

---





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
- [pgx PostgreSQL Driver](../../sources/database/pgx.md)
- [PostgreSQL Arrays](../../sources/database/postgresql-arrays.md)
- [PostgreSQL JSON Functions](../../sources/database/postgresql-json.md)
- [River Job Queue](../../sources/tooling/river.md)
- [rueidis](../../sources/tooling/rueidis.md)
- [rueidis GitHub README](../../sources/tooling/rueidis-guide.md)
- [Typesense API](../../sources/infrastructure/typesense.md)
- [Typesense Go Client](../../sources/infrastructure/typesense-go.md)



---

**Need Help?** [Open an issue](https://github.com/revenge-project/revenge/issues) or [Join the discussion](https://github.com/revenge-project/revenge/discussions)
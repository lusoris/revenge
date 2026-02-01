## Table of Contents

- [Database Auto-Healing & Recovery](#database-auto-healing-recovery)
  - [Features](#features)
  - [Related Documentation](#related-documentation)
    - [Related Pages](#related-pages)
    - [Learn More](#learn-more)

# Database Auto-Healing & Recovery




> How Revenge automatically recovers from database issues

Revenge includes built-in database recovery mechanisms. The connection pool (pgxpool) automatically reconnects after network interruptions. Health checks run periodically with exponential backoff on failures. If the database becomes unavailable, Revenge enters a degraded read-only mode using cached data. Database migrations run automatically on startup, and daily backups can be configured for disaster recovery.

---





---






## Features
<!-- Feature list placeholder -->













## Related Documentation
### Related Pages
<!-- Related wiki pages -->

### Learn More

Official documentation and guides:
- [pgx PostgreSQL Driver](../../sources/database/pgx.md)
- [PostgreSQL Arrays](../../sources/database/postgresql-arrays.md)
- [PostgreSQL JSON Functions](../../sources/database/postgresql-json.md)
- [River Job Queue](../../sources/tooling/river.md)



---

**Need Help?** [Open an issue](https://github.com/revenge-project/revenge/issues) or [Join the discussion](https://github.com/revenge-project/revenge/discussions)
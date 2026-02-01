## Table of Contents

- [RBAC Service](#rbac-service)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [Related Pages](#related-pages)
    - [Learn More](#learn-more)

# RBAC Service




> Role-based access control with Casbin

The RBAC service controls who can do what in Revenge. Built on Casbin, it supports roles (Admin, User, Guest) and fine-grained permissions. Admins define policies like "Users can view movies" or "Guests cannot access adult content". Policies are stored in PostgreSQL and cached for fast authorization checks on every request.

---





---


## How It Works

<!-- How it works -->
## Features
<!-- Feature list placeholder -->
## Configuration
## Related Documentation
### Related Pages
<!-- Related wiki pages -->

### Learn More

Official documentation and guides:
- [Casbin](../../sources/security/casbin.md)
- [Uber fx](../../sources/tooling/fx.md)
- [pgx PostgreSQL Driver](../../sources/database/pgx.md)
- [PostgreSQL Arrays](../../sources/database/postgresql-arrays.md)
- [PostgreSQL JSON Functions](../../sources/database/postgresql-json.md)
- [River Job Queue](../../sources/tooling/river.md)



---

**Need Help?** [Open an issue](https://github.com/revenge-project/revenge/issues) or [Join the discussion](https://github.com/revenge-project/revenge/discussions)
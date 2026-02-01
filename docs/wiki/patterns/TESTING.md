

---
---

## Table of Contents

- [Testing Patterns](#testing-patterns)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [Related Pages](#related-pages)
    - [Learn More](#learn-more)


# Testing Patterns




> Testing standards and patterns for code quality


Revenge maintains high code quality through comprehensive testing. All new code requires 80% test coverage minimum. Unit tests use table-driven patterns for clarity and completeness, with mocks generated automatically via Mockery. Integration tests run against real PostgreSQL using embedded-postgres for speed or testcontainers for full fidelity. The CI pipeline runs all tests on every pull request, blocking merges that reduce coverage or introduce failing tests.


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
- [Go Testing](../../sources/go/stdlib/testing.md)
- [Testify](https://github.com/stretchr/testify)
- [Mockery](../../sources/testing/mockery-guide.md)
- [Testcontainers Go](../../sources/testing/testcontainers.md)
- [Embedded Postgres](../../sources/testing/embedded-postgres-guide.md)



---

**Need Help?** [Open an issue](https://github.com/revenge-project/revenge/issues) or [Join the discussion](https://github.com/revenge-project/revenge/discussions)
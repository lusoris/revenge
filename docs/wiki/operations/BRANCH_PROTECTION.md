## Table of Contents

- [Branch Protection Rules](#branch-protection-rules)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [See Also](#see-also)



---
sources:
  - name: Conventional Commits
    url: ../../sources/standards/conventional-commits.md
    note: Auto-resolved from conventional-commits
  - name: Git Flow
    url: ../../sources/standards/gitflow.md
    note: Auto-resolved from gitflow
design_refs:
  - title: operations
    path: operations/INDEX.md
  - title: 01_ARCHITECTURE
    path: architecture/01_ARCHITECTURE.md
  - title: 02_DESIGN_PRINCIPLES
    path: architecture/02_DESIGN_PRINCIPLES.md
  - title: 03_METADATA_SYSTEM
    path: architecture/03_METADATA_SYSTEM.md
---

# Branch Protection Rules




> Rules that protect main and develop branches from accidental changes

Branch protection rules prevent direct pushes to main and develop branches. All changes must go through pull requests with at least one approval. CI checks (tests, linting, coverage) must pass before merging. Force pushes are disabled to preserve commit history. These rules are configured in GitHub repository settings and ensure code quality and review for all changes.

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
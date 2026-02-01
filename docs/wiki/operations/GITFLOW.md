## Table of Contents

- [Git Workflow & Branching Strategy](#git-workflow-branching-strategy)
  - [Contents](#contents)
  - [How It Works](#how-it-works)
  - [Features](#features)
  - [Configuration](#configuration)
  - [Related Documentation](#related-documentation)
    - [See Also](#see-also)



---
sources:
  - name: Conventional Commits
    url: ../sources/standards/conventional-commits.md
    note: Auto-resolved from conventional-commits
  - name: Git Flow
    url: ../sources/standards/gitflow.md
    note: Auto-resolved from gitflow
  - name: Go io
    url: ../sources/go/stdlib/io.md
    note: Auto-resolved from go-io
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

# Git Workflow & Branching Strategy




> How we use Git branches for features, releases, and hotfixes

Revenge uses a Gitflow-style branching strategy. The main branch contains production-ready code, while develop is the integration branch for the next release. Create feature/* branches for new work, fix/* for bug fixes. All merges require pull request review. Commits follow Conventional Commits format (feat:, fix:, docs:) for automatic changelog generation. Git hooks enforce commit message format and run linters before push.

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
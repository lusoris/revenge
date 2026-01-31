# Branch Protection Rules



<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Main Branch (`main`)](#main-branch-main)
  - [Active Protection Rules:](#active-protection-rules)
- [Develop Branch (`develop`)](#develop-branch-develop)
  - [Active Protection Rules:](#active-protection-rules)
- [Feature Branches (`feature/*`, `fix/*`)](#feature-branches-feature-fix)
  - [Naming Conventions:](#naming-conventions)
  - [Best Practices:](#best-practices)
- [Release Branches (`release/*`)](#release-branches-release)
  - [Protection Rules:](#protection-rules)
  - [Release Process:](#release-process)
- [Hotfix Branches (`hotfix/*`)](#hotfix-branches-hotfix)
  - [Protection Rules:](#protection-rules)
  - [Hotfix Process:](#hotfix-process)
- [Tag Protection](#tag-protection)
  - [Protected Tag Patterns:](#protected-tag-patterns)
  - [Rules:](#rules)
- [Repository Settings](#repository-settings)
  - [General:](#general)
  - [Security:](#security)
- [Setting Up Protection Rules](#setting-up-protection-rules)
  - [Via GitHub UI:](#via-github-ui)
  - [Via GitHub CLI:](#via-github-cli)
  - [Via Terraform (Infrastructure as Code):](#via-terraform-infrastructure-as-code)
- [Workflow Summary](#workflow-summary)
- [Additional Resources](#additional-resources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)

<!-- TOC-END -->

## Status

| Dimension | Status | Notes |
|-----------|--------|-------|
| Design | 🔴 |  |
| Sources | 🔴 |  |
| Instructions | 🔴 |  |
| Code | 🔴 |  |
| Linting | 🔴 |  |
| Unit Testing | 🔴 |  |
| Integration Testing | 🔴 |  |

---

This document outlines the branch protection rules for the revenge repository.

> **Status**: ✅ Branch protection is **ACTIVE** on `main` and `develop`

## Main Branch (`main`)

The `main` branch contains production-ready code.

### Active Protection Rules:

- ✅ **Require pull request reviews before merging**
  - Required approving reviews: 1
  - Dismiss stale pull request approvals when new commits are pushed

- ✅ **Require status checks to pass before merging**
  - Require branches to be up to date before merging
  - Required status checks:
    - `test`
    - `lint`
    - `build`

- ✅ **Require linear history** (no merge commits)

- ✅ **Allow force pushes**: Disabled

- ✅ **Allow deletions**: Disabled

## Develop Branch (`develop`)

The `develop` branch is the integration branch for ongoing development.

### Active Protection Rules:

- ✅ **Require status checks to pass before merging**
  - Require branches to be up to date before merging
  - Required status checks:
    - `test`
    - `lint`
    - `build`

- ✅ **Allow force pushes**: Disabled

- ✅ **Allow deletions**: Disabled

## Feature Branches (`feature/*`, `fix/*`)

Feature branches don't have strict protection rules but should follow naming conventions.

### Naming Conventions:

- `feature/<short-description>` - For new features
- `fix/<short-description>` - For bug fixes
- `docs/<short-description>` - For documentation changes
- `refactor/<short-description>` - For code refactoring
- `test/<short-description>` - For test additions/changes
- `chore/<short-description>` - For maintenance tasks

### Best Practices:

- Keep branches short-lived (< 1 week)
- Rebase on `develop` regularly
- Delete after merging
- Include issue number in branch name when applicable: `feature/123-add-user-auth`

## Release Branches (`release/*`)

Release branches are created from `develop` when preparing a new release.

### Protection Rules:

- ✅ **Require pull request reviews before merging to main**
  - Required approving reviews: 2

- ✅ **Only allow merges to `main` and `develop`**

### Release Process:

1. Create release branch from `develop`: `release/v1.0.0`
2. Perform final testing and bug fixes
3. Update version numbers and CHANGELOG
4. Merge to `main` via PR
5. Tag the release on `main`
6. Merge back to `develop`
7. Delete release branch

## Hotfix Branches (`hotfix/*`)

Hotfix branches are created from `main` for critical production fixes.

### Protection Rules:

- ✅ **Require pull request reviews before merging**
  - Required approving reviews: 1 (can be expedited for critical fixes)

### Hotfix Process:

1. Create hotfix branch from `main`: `hotfix/v1.0.1-critical-bug`
2. Fix the issue
3. Update version and CHANGELOG
4. Create PR to `main`
5. After merge to `main`, tag the release
6. Cherry-pick or merge back to `develop`
7. Delete hotfix branch

## Tag Protection

### Protected Tag Patterns:

- `v*.*.*` - All version tags

### Rules:

- ✅ **Only allow repository admins to create/delete tags**
- ✅ **Require signed tags**

## Repository Settings

### General:

- ✅ **Disable merge commits** - Use squash or rebase only
- ✅ **Automatically delete head branches** after PR merge
- ✅ **Disable wiki** (use docs/ instead)
- ✅ **Enable issue templates**
- ✅ **Enable PR templates**

### Security:

- ✅ **Enable Dependabot alerts**
- ✅ **Enable Dependabot security updates**
- ✅ **Enable CodeQL analysis**
- ✅ **Enable secret scanning**
- ✅ **Enable push protection for secrets**

## Setting Up Protection Rules

### Via GitHub UI:

1. Go to Repository Settings
2. Click on "Branches" in the left sidebar
3. Click "Add rule" under "Branch protection rules"
4. Enter the branch name pattern (e.g., `main`)
5. Configure the protection rules as listed above
6. Click "Create" or "Save changes"

### Via GitHub CLI:

```bash
# Install GitHub CLI first
gh repo view --web

# Then configure in the UI, or use the API
```

### Via Terraform (Infrastructure as Code):

```hcl
resource "github_branch_protection" "main" {
  repository_id = github_repository.revenge_go.node_id
  pattern       = "main"

  required_pull_request_reviews {
    required_approving_review_count = 2
    dismiss_stale_reviews          = true
    require_code_owner_reviews     = true
  }

  required_status_checks {
    strict = true
    contexts = [
      "Lint & Code Quality",
      "Test (Go 1.22 / ubuntu-latest)",
      # ... other checks
    ]
  }

  enforce_admins        = true
  require_signed_commits = true
  require_linear_history = true

  restrict_pushes {
    blocks_creations = true
  }
}
```

## Workflow Summary

```
feature/xyz ──┐
              ├──> develop ──> release/v1.x ──> main (tagged v1.x)
fix/abc ──────┘                                   │
                                                  │
hotfix/v1.x.y ────────────────────────────────────┘
```

## Additional Resources

- [GitHub Branch Protection Documentation](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/defining-the-mergeability-of-pull-requests/about-protected-branches)
- [Git Flow Workflow](https://nvie.com/posts/a-successful-git-branching-model/)
- [Conventional Commits](https://www.conventionalcommits.org/)


<!-- DESIGN-BREADCRUMBS-START -->

## Related Design Docs

> Auto-generated cross-references to related design documentation

**Category**: [Operations](INDEX.md)

### In This Section

- [Advanced Patterns & Best Practices](BEST_PRACTICES.md)
- [Database Auto-Healing & Consistency Restoration](DATABASE_AUTO_HEALING.md)
- [Clone repository](DEVELOPMENT.md)
- [GitFlow Workflow Guide](GITFLOW.md)
- [Revenge - Reverse Proxy & Deployment Best Practices](REVERSE_PROXY.md)
- [revenge - Setup Guide](SETUP.md)

### Related Topics

- [Revenge - Architecture v2](../architecture/01_ARCHITECTURE.md) _Architecture_
- [Revenge - Design Principles](../architecture/02_DESIGN_PRINCIPLES.md) _Architecture_
- [Revenge - Metadata System](../architecture/03_METADATA_SYSTEM.md) _Architecture_
- [Revenge - Player Architecture](../architecture/04_PLAYER_ARCHITECTURE.md) _Architecture_
- [Plugin Architecture Decision](../architecture/05_PLUGIN_ARCHITECTURE_DECISION.md) _Architecture_

### Indexes

- [Design Index](../DESIGN_INDEX.md) - All design docs by category/topic
- [Source of Truth](../00_SOURCE_OF_TRUTH.md) - Package versions and status

<!-- DESIGN-BREADCRUMBS-END -->

---

<!-- SOURCE-BREADCRUMBS-START -->

## Sources & Cross-References

> Auto-generated section linking to external documentation sources

### Cross-Reference Indexes

- [All Sources Index](../../sources/SOURCES_INDEX.md) - Complete list of external documentation
- [Design ↔ Sources Map](../../sources/DESIGN_CROSSREF.md) - Which docs reference which sources

### Referenced Sources

| Source | Documentation |
|--------|---------------|
| [Conventional Commits](https://www.conventionalcommits.org/) | [Local](../../sources/standards/conventional-commits.md) |
| [Git Flow](https://nvie.com/posts/a-successful-git-branching-model/) | [Local](../../sources/standards/gitflow.md) |

<!-- SOURCE-BREADCRUMBS-END -->
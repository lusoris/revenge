# GitHub Documentation

> Reference documentation for GitHub configuration

## Documents

| Document | Description |
|----------|-------------|
| [ACTIONS.md](ACTIONS.md) | Project workflow patterns and conventions |
| [DEPENDABOT.md](DEPENDABOT.md) | Automated dependency updates configuration |

## External Documentation

Comprehensive GitHub Actions documentation (auto-fetched from GitHub):

| Document | Description |
|----------|-------------|
| [github-actions/main.md](github-actions/main.md) | GitHub Actions overview and getting started |
| [github-actions/syntax.md](github-actions/syntax.md) | Complete workflow syntax reference |
| [github-actions/contexts.md](github-actions/contexts.md) | Contexts and expressions |
| [github-actions/variables.md](github-actions/variables.md) | Environment variables and secrets |
| [github-actions/secrets.md](github-actions/secrets.md) | Encrypted secrets management |
| [github-actions/caching.md](github-actions/caching.md) | Caching dependencies |
| [github-actions/matrix.md](github-actions/matrix.md) | Matrix build strategies |
| [github-actions/security.md](github-actions/security.md) | Security best practices |

## Quick Links

- [GitHub Actions Docs](https://docs.github.com/en/actions)
- [Dependabot Docs](https://docs.github.com/en/code-security/dependabot)
- [Workflow Syntax](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions)

## Project Configuration

| File | Purpose |
|------|---------|
| [workflows/](../workflows/) | GitHub Actions workflows |
| [dependabot.yml](../dependabot.yml) | Dependency update config |
| [labeler.yml](../labeler.yml) | Auto-labeling rules |
| [CODEOWNERS](../CODEOWNERS) | Code ownership rules |

## Workflows

### Build & Test
| Workflow | Trigger | Purpose |
|----------|---------|---------|
| [ci.yml](../workflows/ci.yml) | Push/PR to main | CI pipeline for main branch |
| [dev.yml](../workflows/dev.yml) | Push/PR to develop | Development build and tests |
| [coverage.yml](../workflows/coverage.yml) | Push/PR | Code coverage reporting |
| [validate-sot.yml](../workflows/validate-sot.yml) | PR/Weekly | Validate SOURCE_OF_TRUTH consistency |

### Release Management
| Workflow | Trigger | Purpose |
|----------|---------|---------|
| [release.yml](../workflows/release.yml) | Tag push (v*.*.*) | Build and publish releases |
| [release-please.yml](../workflows/release-please.yml) | Push to main | Automated release PR creation |

### Security
| Workflow | Trigger | Purpose |
|----------|---------|---------|
| [security.yml](../workflows/security.yml) | PR/Push/Weekly | Security scanning (CodeQL, Trivy, gosec) |

### Automation
| Workflow | Trigger | Purpose |
|----------|---------|---------|
| [auto-label.yml](../workflows/auto-label.yml) | PR/Issue events | Auto-label PRs and issues |
| [labels.yml](../workflows/labels.yml) | Push to main | Sync repository labels |
| [pr-checks.yml](../workflows/pr-checks.yml) | PR events | PR validation and checks |
| [build-status.yml](../workflows/build-status.yml) | Workflow completion | Report build status |
| [stale.yml](../workflows/stale.yml) | Daily | Mark stale issues/PRs |

### Documentation
| Workflow | Trigger | Purpose |
|----------|---------|---------|
| [doc-validation.yml](../workflows/doc-validation.yml) | PR/Push | Validate documentation structure |
| [fetch-sources.yml](../workflows/fetch-sources.yml) | Weekly/Manual | Fetch external documentation |
| [source-refresh.yml](../workflows/source-refresh.yml) | Weekly/Manual | Refresh external source docs |

### Dependencies
| Workflow | Trigger | Purpose |
|----------|---------|---------|
| [dependency-update.yml](../workflows/dependency-update.yml) | Weekly | Automated dependency updates |

### Utilities
| Workflow | Trigger | Purpose |
|----------|---------|---------|
| [_versions.yml](../workflows/_versions.yml) | Reusable | Extract versions from SOURCE_OF_TRUTH |

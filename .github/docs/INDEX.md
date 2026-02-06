# GitHub Documentation

> Reference documentation for GitHub configuration

## Documents

| Document | Description |
|----------|-------------|
| [ACTIONS.md](ACTIONS.md) | Project workflow patterns and conventions |
| [DEPENDABOT.md](DEPENDABOT.md) | Automated dependency updates configuration |

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
| [ci.yml](../workflows/ci.yml) | Push/PR to main, develop | Lint, unit tests, Docker build + Trivy, govulncheck |
| [coverage.yml](../workflows/coverage.yml) | PR | Coverage report as PR comment |

### Release

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| [develop.yml](../workflows/develop.yml) | Push to develop | Auto-build `:develop` Docker image + Helm chart |
| [release-please.yml](../workflows/release-please.yml) | Push to main | Automated release PR creation |

### Security

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| [security.yml](../workflows/security.yml) | Schedule + PR | CodeQL, Trivy, govulncheck, dependency review |

### Automation

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| [labels.yml](../workflows/labels.yml) | Push to main | Sync repository labels |
| [pr-checks.yml](../workflows/pr-checks.yml) | PR events | Title format, branch name, merge conflicts |
| [stale.yml](../workflows/stale.yml) | Daily | Mark stale issues/PRs |
| [wiki-sync.yml](../workflows/wiki-sync.yml) | Push to develop (docs changes) | Auto-sync docs to GitHub wiki |

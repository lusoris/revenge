# CI/CD

<!-- DESIGN: operations -->

**Platform**: GitHub Actions
**Registry**: `ghcr.io/lusoris/revenge` (Docker) + `ghcr.io/lusoris/charts/revenge` (Helm)
**Release**: Release Please (automated semantic versioning)

> Continuous integration, security scanning, and automated releases

---

## Workflows

### ci.yml — Primary CI Pipeline

**Trigger**: Push/PR to main, develop

| Job | Purpose | Details |
|-----|---------|---------|
| lint | Code quality | golangci-lint v2.8.0 + go vet |
| test | Unit tests | Race detection, coverage to Codecov |
| docker | Build + scan | Docker build + Trivy SARIF scan (CRITICAL, HIGH) |
| integration | Integration tests | Runs after lint+test pass |
| vuln | Vulnerability check | govulncheck |

### develop.yml — Dev Builds

**Trigger**: Push to develop

| Job | Purpose | Details |
|-----|---------|---------|
| docker | Multi-arch build | amd64/arm64, pushes `:develop` and `:dev-{sha}` tags |
| helm-dev | Dev chart | Packages chart version `0.0.0-dev.{sha}` |
| cleanup | Prune images | Keeps 5 tagged builds, removes orphaned manifests |

### release-please.yml — Automated Releases

**Trigger**: Push to main

| Job | Purpose | Details |
|-----|---------|---------|
| release-please | Detect changes | Creates release PR with semver bump |
| docker | Release image | Multi-arch, tags: `{version}`, `{major}.{minor}`, `latest`, SBOM (CycloneDX) |
| helm | Release chart | Updates Chart.yaml, packages and pushes |

### security.yml — Security Scanning

**Trigger**: Push/PR to main/develop, weekly schedule (Mon 00:00 UTC)

| Job | Purpose | Details |
|-----|---------|---------|
| codeql | Static analysis | security-extended + security-and-quality queries |
| trivy-repo | Filesystem scan | CRITICAL, HIGH severity |
| trivy-docker | Image scan | CRITICAL, HIGH severity |
| govulncheck | Go vulnerabilities | Go vulnerability database |
| dependency-review | PR dependencies | Fail on high severity (PR-only) |

### coverage.yml — Coverage Reporting

**Trigger**: PR to main/develop

Generates coverage report, uploads to Codecov, posts PR comment with summary.

### pr-checks.yml — PR Validation

**Trigger**: PR opened/sync/edit/ready

| Job | Purpose |
|-----|---------|
| pr-title | Enforces conventional commit format in PR title |
| pr-branch | Validates branch naming: `{type}/{name}` (feature/, fix/, docs/, etc.) |
| pr-conflicts | Auto-labels PRs with merge conflicts |

### stale.yml — Stale Cleanup

**Trigger**: Daily (00:00 UTC)

Issues stale after 60d (closed +14d), PRs stale after 30d (closed +14d). Exempts: pinned, security, bug, in-progress, do-not-close.

### wiki-sync.yml — Documentation Sync

**Trigger**: Push to develop (docs/dev/design/, root MDs)

Auto-syncs design documentation to GitHub Wiki. Generates Home.md + sidebar from design docs.

### labels.yml — Label Sync

**Trigger**: Push to main (when `.github/labels.yml` changes)

Syncs GitHub label definitions from config file.

## Environment

All workflows use:

```
GO_VERSION: 1.25
GOEXPERIMENT: greenteagc,jsonv2
```

## Release Strategy

- **Release Please** handles semver bumps from conventional commits
- Current version: `0.2.0` (from `.release-please-manifest.json`)
- Changelog sections: feat, fix, perf, docs, refactor (visible); test, build, ci, chore, style (hidden)
- Version stamped in: `cmd/revenge/main.go`, `internal/version/version.go`
- Branch strategy: `develop` (working) → `main` (stable/release)

## Current Status

CI/CD workflows were recently simplified and corrected. Some workflows may still need further fixes — tracked in `.workingdir3/CODE_ISSUES.md`.

## Related Documentation

- [DEVELOPMENT.md](DEVELOPMENT.md) - Development setup and Makefile targets
- [../technical/TESTING.md](../technical/TESTING.md) - Testing infrastructure

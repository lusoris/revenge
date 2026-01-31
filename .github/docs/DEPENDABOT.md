# Dependabot Configuration

> Source: https://docs.github.com/en/code-security/dependabot/dependabot-version-updates/configuration-options-for-the-dependabot.yml-file

## Overview

The `dependabot.yml` file configures automatic dependency updates. Place it in `.github/dependabot.yml`.

## Basic Structure

```yaml
version: 2
updates:
  - package-ecosystem: "npm"
    directory: "/"
    schedule:
      interval: "weekly"
```

## Required Fields

| Field | Description |
|-------|-------------|
| `version` | Must be `2` |
| `package-ecosystem` | Package manager type |
| `directory` | Location of manifest files |
| `schedule.interval` | Update frequency |

## Supported Ecosystems

| Ecosystem | Manifest Files |
|-----------|----------------|
| `npm` | package.json, package-lock.json |
| `pip` | requirements.txt, Pipfile |
| `gomod` | go.mod |
| `docker` | Dockerfile |
| `github-actions` | .github/workflows/*.yml |
| `terraform` | *.tf |
| `composer` | composer.json |
| `cargo` | Cargo.toml |
| `bundler` | Gemfile |
| `maven` | pom.xml |
| `gradle` | build.gradle |

## Schedule Options

```yaml
schedule:
  interval: "daily"    # daily, weekly, monthly
  day: "monday"        # For weekly
  time: "09:00"        # UTC by default
  timezone: "Europe/Berlin"
```

## Dependency Grouping

```yaml
groups:
  development-dependencies:
    patterns:
      - "eslint*"
      - "prettier*"
    update-types:
      - "minor"
      - "patch"

  production-dependencies:
    dependency-type: "production"
```

## Allow/Ignore Rules

### Allow Specific Dependencies

```yaml
allow:
  - dependency-name: "lodash"
  - dependency-name: "express"
    dependency-type: "production"
```

### Ignore Dependencies

```yaml
ignore:
  - dependency-name: "aws-sdk"
  - dependency-name: "typescript"
    versions: [">=5.0.0"]
  - dependency-name: "*"
    update-types: ["version-update:semver-major"]
```

## PR Configuration

```yaml
updates:
  - package-ecosystem: "npm"
    directory: "/"
    schedule:
      interval: "weekly"

    # PR settings
    assignees:
      - "username"
    reviewers:
      - "team:reviewers"
    labels:
      - "dependencies"
      - "automated"
    milestone: 4

    # Commit message
    commit-message:
      prefix: "deps"
      prefix-development: "deps(dev)"
      include: "scope"
```

## Private Registries

```yaml
registries:
  npm-private:
    type: npm-registry
    url: https://npm.pkg.github.com
    token: ${{ secrets.NPM_TOKEN }}

updates:
  - package-ecosystem: "npm"
    directory: "/"
    registries:
      - npm-private
```

## Advanced Options

### Target Branch

```yaml
target-branch: "develop"
```

### Open PR Limit

```yaml
open-pull-requests-limit: 10
```

### Rebase Strategy

```yaml
rebase-strategy: "auto"  # auto, disabled
```

### Versioning Strategy

```yaml
versioning-strategy: "increase"  # auto, increase, increase-if-necessary, lockfile-only, widen
```

### Vendor Dependencies

```yaml
vendor: true
```

## Complete Example

```yaml
version: 2

registries:
  npm-github:
    type: npm-registry
    url: https://npm.pkg.github.com
    token: ${{ secrets.GITHUB_TOKEN }}

updates:
  # JavaScript/TypeScript
  - package-ecosystem: "npm"
    directory: "/"
    schedule:
      interval: "weekly"
      day: "monday"
      time: "09:00"
      timezone: "Europe/Berlin"
    registries:
      - npm-github
    groups:
      dev-deps:
        patterns:
          - "@types/*"
          - "eslint*"
          - "prettier*"
          - "vitest*"
        update-types: ["minor", "patch"]
    ignore:
      - dependency-name: "*"
        update-types: ["version-update:semver-major"]
    labels:
      - "dependencies"
      - "javascript"
    commit-message:
      prefix: "deps(npm)"

  # Go modules
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"
    labels:
      - "dependencies"
      - "go"
    commit-message:
      prefix: "deps(go)"

  # GitHub Actions
  - package-ecosystem: "github-actions"
    directory: "/"
    schedule:
      interval: "weekly"
    labels:
      - "dependencies"
      - "ci"
    commit-message:
      prefix: "ci"

  # Docker
  - package-ecosystem: "docker"
    directory: "/"
    schedule:
      interval: "monthly"
    labels:
      - "dependencies"
      - "docker"
```

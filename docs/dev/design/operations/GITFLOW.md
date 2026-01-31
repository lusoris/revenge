# GitFlow Workflow Guide

<!-- SOURCES: conventional-commits, gitflow, go-io -->

<!-- DESIGN: operations, 01_ARCHITECTURE, 02_DESIGN_PRINCIPLES, 03_METADATA_SYSTEM -->


<!-- TOC-START -->

## Table of Contents

- [Status](#status)
- [Branch Structure](#branch-structure)
- [Main Branches](#main-branches)
  - [`main`](#main)
  - [`develop`](#develop)
- [Supporting Branches](#supporting-branches)
  - [Feature Branches (`feature/*`)](#feature-branches-feature)
  - [Fix Branches (`fix/*`)](#fix-branches-fix)
  - [Release Branches (`release/*`)](#release-branches-release)
  - [Hotfix Branches (`hotfix/*`)](#hotfix-branches-hotfix)
- [Commit Message Format](#commit-message-format)
  - [Types:](#types)
  - [Examples:](#examples)
- [Pull Request Process](#pull-request-process)
  - [Creating a PR](#creating-a-pr)
  - [Reviewing a PR](#reviewing-a-pr)
  - [Merging a PR](#merging-a-pr)
- [Common Scenarios](#common-scenarios)
  - [Updating your branch with latest develop](#updating-your-branch-with-latest-develop)
  - [Fixing a mistake in your last commit](#fixing-a-mistake-in-your-last-commit)
  - [Cherry-picking a commit](#cherry-picking-a-commit)
  - [Reverting a commit](#reverting-a-commit)
- [Repository Setup](#repository-setup)
  - [Initial Setup](#initial-setup)
  - [Developer Setup](#developer-setup)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)
  - [Merge Conflicts](#merge-conflicts)
  - [Accidentally committed to wrong branch](#accidentally-committed-to-wrong-branch)
  - [Need to sync fork](#need-to-sync-fork)
- [Resources](#resources)
- [Related Design Docs](#related-design-docs)
  - [In This Section](#in-this-section)
  - [Related Topics](#related-topics)
  - [Indexes](#indexes)
- [Sources & Cross-References](#sources-cross-references)
  - [Cross-Reference Indexes](#cross-reference-indexes)
  - [Referenced Sources](#referenced-sources)

<!-- TOC-END -->

## Status

| Dimension | Status |
|-----------|--------|
| Design | 🔴 |
| Sources | 🔴 |
| Instructions | 🔴 |
| Code | 🔴 |
| Linting | 🔴 |
| Unit Testing | 🔴 |
| Integration Testing | 🔴 |
---

This document describes the Git branching strategy for revenge.

## Branch Structure

```
main (production)
  │
  ├── develop (integration)
  │     │
  │     ├── feature/user-auth
  │     ├── feature/media-scanner
  │     ├── fix/database-leak
  │     └── ...
  │
  ├── release/v0.1.0
  └── hotfix/v0.0.2-critical-bug
```

## Main Branches

### `main`
- **Purpose**: Production-ready code
- **Protection**: Highest level
- **Updates**: Only via PR from `release/*` or `hotfix/*`
- **Tags**: All releases are tagged here

### `develop`
- **Purpose**: Integration branch for features
- **Protection**: Medium level
- **Updates**: Via PR from `feature/*`, `fix/*`, etc.
- **Tags**: None

## Supporting Branches

### Feature Branches (`feature/*`)

**Purpose**: Develop new features

**Naming**: `feature/<issue-number>-<short-description>`
- Example: `feature/123-add-jwt-auth`

**Branch from**: `develop`
**Merge to**: `develop`

**Workflow**:
```bash
# Create feature branch
git checkout develop
git pull origin develop
git checkout -b feature/123-add-jwt-auth

# Work on feature
git add .
git commit -m "feat(auth): implement JWT authentication"

# Keep up to date with develop
git checkout develop
git pull origin develop
git checkout feature/123-add-jwt-auth
git rebase develop

# Push and create PR
git push origin feature/123-add-jwt-auth
# Create PR: feature/123-add-jwt-auth -> develop
```

### Fix Branches (`fix/*`)

**Purpose**: Fix bugs in development

**Naming**: `fix/<issue-number>-<short-description>`
- Example: `fix/456-database-connection-leak`

**Branch from**: `develop`
**Merge to**: `develop`

**Workflow**: Same as feature branches

### Release Branches (`release/*`)

**Purpose**: Prepare for production release

**Naming**: `release/v<major>.<minor>.<patch>`
- Example: `release/v0.1.0`

**Branch from**: `develop`
**Merge to**: `main` AND `develop`

**Workflow**:
```bash
# Create release branch
git checkout develop
git pull origin develop
git checkout -b release/v0.1.0

# Prepare release
# - Update version numbers
# - Update CHANGELOG.md
# - Final bug fixes only

# Update version
echo "v0.1.0" > VERSION

# Commit changes
git add .
git commit -m "chore(release): prepare v0.1.0"

# Push release branch
git push origin release/v0.1.0

# Create PR to main
# After approval and merge:

# Tag the release on main
git checkout main
git pull origin main
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0

# Merge back to develop
git checkout develop
git merge release/v0.1.0
git push origin develop

# Delete release branch
git branch -d release/v0.1.0
git push origin --delete release/v0.1.0
```

### Hotfix Branches (`hotfix/*`)

**Purpose**: Emergency fixes for production

**Naming**: `hotfix/v<major>.<minor>.<patch>-<description>`
- Example: `hotfix/v0.1.1-security-patch`

**Branch from**: `main`
**Merge to**: `main` AND `develop`

**Workflow**:
```bash
# Create hotfix branch from main
git checkout main
git pull origin main
git checkout -b hotfix/v0.1.1-security-patch

# Fix the issue
git add .
git commit -m "fix(security): patch critical vulnerability"

# Update version
echo "v0.1.1" > VERSION

# Commit version update
git add VERSION
git commit -m "chore(release): bump to v0.1.1"

# Push and create PR to main
git push origin hotfix/v0.1.1-security-patch
# Create PR: hotfix/v0.1.1-security-patch -> main

# After merge to main, tag the release
git checkout main
git pull origin main
git tag -a v0.1.1 -m "Hotfix v0.1.1"
git push origin v0.1.1

# Merge to develop
git checkout develop
git merge hotfix/v0.1.1-security-patch
git push origin develop

# Delete hotfix branch
git branch -d hotfix/v0.1.1-security-patch
git push origin --delete hotfix/v0.1.1-security-patch
```

## Commit Message Format

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `style`: Formatting
- `refactor`: Code restructuring
- `perf`: Performance
- `test`: Tests
- `chore`: Maintenance
- `ci`: CI/CD
- `build`: Build system

### Examples:
```bash
feat(auth): add JWT token generation
fix(database): resolve connection pool leak
docs(readme): update installation instructions
test(user): add unit tests for user service
chore(deps): update dependencies
```

## Pull Request Process

### Creating a PR

1. **Push your branch**:
   ```bash
   git push origin feature/123-my-feature
   ```

2. **Create PR on GitHub**:
   - Use the PR template
   - Link related issues
   - Add appropriate labels
   - Request reviewers
   - Assign to yourself

3. **PR Checks**:
   - CI must pass
   - Code coverage maintained
   - No conflicts with base branch
   - Required reviews obtained

### Reviewing a PR

1. **Code Review Checklist**:
   - [ ] Code follows project standards
   - [ ] Tests included and passing
   - [ ] Documentation updated
   - [ ] No security issues
   - [ ] Performance acceptable
   - [ ] API compatibility maintained

2. **Approval**:
   - Approve if all checks pass
   - Request changes if issues found
   - Comment for questions or suggestions

### Merging a PR

1. **Squash and Merge** (default for most PRs):
   - Combines all commits into one
   - Clean commit history
   - Use for feature branches

2. **Rebase and Merge**:
   - Maintains individual commits
   - Linear history
   - Use when commits are well-organized

3. **Merge Commit** (disabled):
   - Not used to maintain linear history

## Common Scenarios

### Updating your branch with latest develop

```bash
git checkout develop
git pull origin develop
git checkout feature/my-feature
git rebase develop

# If conflicts occur
git rebase --continue
# or
git rebase --abort
```

### Fixing a mistake in your last commit

```bash
# Amend the last commit
git add .
git commit --amend --no-edit

# Force push (only if not yet reviewed!)
git push origin feature/my-feature --force-with-lease
```

### Cherry-picking a commit

```bash
git checkout target-branch
git cherry-pick <commit-hash>
```

### Reverting a commit

```bash
# Create a new commit that undoes changes
git revert <commit-hash>
git push origin <branch>
```

## Repository Setup

### Initial Setup

```bash
# Clone repository
git clone https://github.com/revenge/revenge.git
cd revenge

# Install Git hooks
./scripts/install-hooks.sh  # Linux/macOS
# or
.\scripts\install-hooks.ps1  # Windows

# Create develop branch (maintainers only)
git checkout -b develop
git push origin develop

# Set develop as default branch (GitHub settings)
```

### Developer Setup

```bash
# Clone and setup
git clone https://github.com/revenge/revenge.git
cd revenge
./scripts/install-hooks.sh

# Checkout develop
git checkout develop
git pull origin develop

# Start working on feature
git checkout -b feature/my-awesome-feature
```

## Best Practices

1. **Keep branches short-lived**: Merge within 1 week
2. **One feature per branch**: Don't mix unrelated changes
3. **Rebase regularly**: Keep up-to-date with base branch
4. **Write good commits**: Clear, descriptive messages
5. **Test before pushing**: Run tests locally
6. **Review thoroughly**: Both as author and reviewer
7. **Delete merged branches**: Keep repository clean
8. **Use draft PRs**: For work-in-progress
9. **Link issues**: Connect PRs to issues
10. **Update documentation**: Keep docs in sync with code

## Troubleshooting

### Merge Conflicts

```bash
# During rebase
git rebase develop
# Fix conflicts in files
git add <resolved-files>
git rebase --continue

# If you want to start over
git rebase --abort
```

### Accidentally committed to wrong branch

```bash
# Move commit to correct branch
git checkout correct-branch
git cherry-pick <commit-hash>

# Remove from wrong branch
git checkout wrong-branch
git reset --hard HEAD~1
```

### Need to sync fork

```bash
# Add upstream remote
git remote add upstream https://github.com/revenge/revenge.git

# Fetch and merge
git fetch upstream
git checkout develop
git merge upstream/develop
git push origin develop
```

## Resources

- [Git Flow Original Article](https://nvie.com/posts/a-successful-git-branching-model/)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [GitHub Flow](https://guides.github.com/introduction/flow/)
- [Git Best Practices](https://sethrobertson.github.io/GitBestPractices/)


---


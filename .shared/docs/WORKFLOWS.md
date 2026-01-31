# Development Workflows

**Purpose**: Step-by-step workflows for common development tasks in the Revenge project

**Last Updated**: 2026-01-31

---

## Overview

This guide provides detailed workflows for common development tasks, with specific commands and tool recommendations for each step.

---

## Table of Contents

- [Local Development Workflow](#local-development-workflow)
- [Remote Development Workflow (Coder)](#remote-development-workflow-coder)
- [Code Review Workflow](#code-review-workflow)
- [Release Workflow](#release-workflow)
- [Testing Workflow](#testing-workflow)
- [Documentation Workflow](#documentation-workflow)
- [Emergency Hotfix Workflow](#emergency-hotfix-workflow)

---

## Local Development Workflow

### Prerequisites

- Go 1.25.6+
- Node.js 20+
- Python 3.12+
- Docker & Docker Compose (for services)
- Git configured

### Step 1: Initial Setup

```bash
# Clone repository
git clone https://github.com/lusoris/revenge.git
cd revenge

# Install Go dependencies
go mod download

# Install Python tools
pip install ruff pytest

# Install frontend dependencies
cd web
npm install
cd ..

# Start development services
docker-compose -f docker-compose.dev.yml up -d

# Verify services are running
docker-compose -f docker-compose.dev.yml ps
```

### Step 2: Create Feature Branch

```bash
# Ensure you're on develop
git checkout develop
git pull origin develop

# Create feature branch
git checkout -b feature/movie-search-api

# Verify branch
git branch --show-current
```

### Step 3: Development Cycle

**Backend (Go)**:
```bash
# Start development server with hot reload
air

# In another terminal, run tests on file save
go test -v ./...

# Access API at http://localhost:8096
```

**Frontend (Svelte)**:
```bash
cd web
npm run dev

# Access at http://localhost:5173
```

**Both**:
- Edit code in VS Code or Zed
- Format on save (configured in IDE)
- Tests run automatically (or manually trigger)

### Step 4: Testing

```bash
# Run all Go tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test -v ./internal/api/...

# Run with race detector
go test -race ./...

# Frontend tests
cd web
npm run test
npm run test:coverage
```

### Step 5: Linting

```bash
# Go linting
golangci-lint run

# Python linting
ruff check scripts/

# Frontend linting
cd web
npm run lint
```

### Step 6: Commit Changes

```bash
# Check status
git status

# Stage changes
git add internal/api/handlers/search.go
git add internal/api/handlers/search_test.go

# Commit with conventional commit format
git commit -m "feat(api): add movie search endpoint

- Implement /api/v1/movies/search handler
- Add query parameter validation
- Include unit tests with 85% coverage
- Update OpenAPI spec

Closes #123"

# Pre-commit hooks run automatically
```

### Step 7: Push and Create PR

```bash
# Push branch to remote
git push origin feature/movie-search-api

# Create PR using GitHub CLI
gh pr create \
  --title "feat(api): add movie search endpoint" \
  --body "$(cat <<EOF
## Summary
Implements movie search API endpoint with query parameter support.

## Changes
- New `/api/v1/movies/search` endpoint
- Query validation and sanitization
- Unit tests (85% coverage)
- OpenAPI spec updated

## Testing
- [x] Unit tests pass
- [x] Integration tests pass
- [x] Manual testing completed
- [x] API documentation updated

## Related Issues
Closes #123
EOF
)"

# Or open browser to create PR manually
gh pr create --web
```

### Step 8: Address PR Feedback

```bash
# Pull latest changes from develop
git fetch origin
git merge origin/develop

# Make requested changes
# ... edit files ...

# Amend existing commit or create new one
git add .
git commit -m "refactor: address PR review comments"

# Force push (if amending)
git push --force-with-lease

# Or normal push (if new commit)
git push
```

### Step 9: Merge and Cleanup

```bash
# After approval, merge via GitHub UI or CLI
gh pr merge --squash

# Switch back to develop and pull
git checkout develop
git pull origin develop

# Delete feature branch
git branch -D feature/movie-search-api
git push origin --delete feature/movie-search-api
```

---

## Remote Development Workflow (Coder)

### Prerequisites

- Coder CLI installed
- Access to https://coder.ancilla.lol
- SSH client

### Step 1: Login to Coder

```bash
# Login (opens browser)
coder login https://coder.ancilla.lol

# Verify login
coder list
```

### Step 2: Create Workspace

```bash
# Create workspace from revenge template
coder create revenge-dev --template revenge

# Or use custom parameters
coder create revenge-dev \
  --template revenge \
  --parameter cpu=8 \
  --parameter memory=16 \
  --parameter ide=vscode-browser

# Wait for workspace to be ready
coder list
```

### Step 3: Connect to Workspace

**Option A: VS Code (Browser)**:
```bash
# Open browser-based VS Code
coder open revenge-dev

# Access at https://coder.ancilla.lol/workspaces/revenge-dev/code
```

**Option B: VS Code (Desktop)**:
```bash
# Install Coder VS Code extension first
# Then connect
coder code revenge-dev

# VS Code desktop will open and connect via SSH
```

**Option C: Zed (SSH)**:
```bash
# Get SSH config
coder ssh revenge-dev

# In Zed, connect to remote
# File → Open Remote → SSH → coder-revenge-dev
```

**Option D: Terminal Only**:
```bash
# SSH into workspace
coder ssh revenge-dev

# Now use vim/nano or any terminal editor
```

### Step 4: Development in Remote Workspace

```bash
# Inside workspace, clone repo if not already done
cd /workspace
git clone https://github.com/lusoris/revenge.git
cd revenge

# All services are already running in workspace:
# - PostgreSQL (localhost:5432)
# - Dragonfly (localhost:6379)
# - Typesense (localhost:8108)

# Start development server
air

# Run tests
go test ./...

# Access forwarded ports from local browser
# Coder automatically forwards ports
```

### Step 5: Port Forwarding

```bash
# Ports are automatically forwarded by Coder
# Access from local browser:
# - Backend: http://localhost:8096
# - Frontend: http://localhost:5173

# Manual port forward if needed
coder port-forward revenge-dev --tcp 8080:8080
```

### Step 6: Stop/Restart Workspace

```bash
# Stop workspace when done
coder stop revenge-dev

# Start again later
coder start revenge-dev

# Check resource usage
coder stat revenge-dev
```

### Step 7: Delete Workspace

```bash
# Delete when no longer needed
coder delete revenge-dev --yes
```

---

## Code Review Workflow

### As Reviewer

#### Step 1: Review PR on GitHub

```bash
# List open PRs
gh pr list

# View specific PR
gh pr view 456

# Checkout PR locally to test
gh pr checkout 456
```

#### Step 2: Review Code Changes

**On GitHub Web**:
1. Go to PR page
2. Click "Files changed" tab
3. Leave inline comments on specific lines
4. Use "Start a review" to batch comments

**In VS Code**:
```bash
# After checking out PR
git diff develop...HEAD

# Or use GitHub Pull Requests extension
# View → Extensions → Search "GitHub Pull Requests"
# Review PRs directly in VS Code
```

#### Step 3: Test Changes Locally

```bash
# Ensure PR branch is checked out
git branch --show-current

# Run tests
go test ./...

# Run linters
golangci-lint run

# Build and run
go run ./cmd/revenge

# Test manually
curl http://localhost:8096/api/v1/movies/search?q=matrix
```

#### Step 4: Leave Review

**Approve**:
```bash
# Via GitHub CLI
gh pr review 456 --approve --body "LGTM! Great work on the tests."
```

**Request Changes**:
```bash
# Via GitHub CLI
gh pr review 456 --request-changes --body "Please address the following:
- Add error handling for empty query
- Update API documentation"
```

**Comment Only**:
```bash
gh pr review 456 --comment --body "Looks good overall. Minor suggestion: consider caching results."
```

### As Author

#### Responding to Feedback

```bash
# Pull latest comments
gh pr view 456

# Make requested changes
# ... edit files ...

# Commit changes
git add .
git commit -m "refactor: address review comments

- Add error handling for empty query
- Update API documentation
- Add caching for search results"

# Push updates
git push

# Reply to review comments
gh pr comment 456 --body "Thanks for the review! All feedback addressed."
```

---

## Release Workflow

### Prerequisites

- Maintainer access to repository
- GitHub CLI configured
- Changelog updated

### Step 1: Prepare Release

```bash
# Ensure on develop branch
git checkout develop
git pull origin develop

# Create release branch
git checkout -b release/v0.2.0

# Update version in relevant files
# - go.mod (if module version changed)
# - package.json (frontend version)
# - VERSION file (if exists)

# Update CHANGELOG.md
# Add new section for v0.2.0 with all changes since last release
```

### Step 2: Final Testing

```bash
# Run full test suite
go test ./...

# Run integration tests
go test -tags=integration ./...

# Build release binary
GOEXPERIMENT=greenteagc,jsonv2 go build -o bin/revenge ./cmd/revenge

# Test release binary
./bin/revenge --version
```

### Step 3: Create Release PR

```bash
# Push release branch
git push origin release/v0.2.0

# Create PR to main
gh pr create \
  --base main \
  --title "Release v0.2.0" \
  --body "$(cat <<EOF
## Release v0.2.0

### New Features
- Movie search API endpoint
- Metadata enrichment improvements

### Bug Fixes
- Fixed connection pool leak
- Resolved cache invalidation issue

### Breaking Changes
None

### Migration Guide
No migration needed.
EOF
)"
```

### Step 4: Merge and Tag

```bash
# After approval, merge to main
gh pr merge --merge

# Switch to main and pull
git checkout main
git pull origin main

# Create annotated tag
git tag -a v0.2.0 -m "Release v0.2.0

New Features:
- Movie search API
- Enhanced metadata enrichment

Bug Fixes:
- Connection pool leak fix
- Cache invalidation improvements"

# Push tag
git push origin v0.2.0
```

### Step 5: Create GitHub Release

```bash
# Create release with notes
gh release create v0.2.0 \
  --title "Revenge v0.2.0" \
  --notes "$(cat CHANGELOG.md | sed -n '/## v0.2.0/,/## v0.1.0/p' | sed '$ d')"

# Upload release artifacts if needed
gh release upload v0.2.0 bin/revenge-linux-amd64
gh release upload v0.2.0 bin/revenge-darwin-amd64
```

### Step 6: Merge Back to Develop

```bash
# Merge main back to develop
git checkout develop
git pull origin develop
git merge main
git push origin develop
```

---

## Testing Workflow

### Unit Testing

```bash
# Run all tests
go test ./...

# Run tests for specific package
go test ./internal/api/handlers

# Run tests with verbose output
go test -v ./...

# Run specific test function
go test -v -run TestMovieSearchHandler ./internal/api/handlers

# Run with coverage
go test -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out

# Run with race detector
go test -race ./...
```

### Integration Testing

```bash
# Run integration tests (requires full stack)
go test -tags=integration ./...

# Run integration tests with cleanup
go test -tags=integration -cleanup ./...

# Run specific integration test
go test -tags=integration -run TestDatabaseMigration ./...
```

### Frontend Testing

```bash
cd web

# Run unit tests
npm run test

# Run tests in watch mode
npm run test:watch

# Run with coverage
npm run test:coverage

# Run E2E tests
npm run test:e2e

# Run specific test file
npm run test -- src/components/MovieCard.test.ts
```

### Debugging Failed Tests

**Go**:
```bash
# Run with verbose output
go test -v ./internal/api/handlers

# Run with increased timeout
go test -timeout 30s ./...

# Print test logs
go test -v ./... 2>&1 | tee test.log

# Debug with delve
dlv test ./internal/api/handlers -- -test.run TestMovieSearchHandler
```

**Frontend**:
```bash
# Debug in VS Code
# Add breakpoint in test file
# F5 to start debugging

# Or use browser dev tools
npm run test:debug
```

---

## Documentation Workflow

### Finding Documentation

```bash
# Search across all docs
grep -r "metadata enrichment" docs/

# Find specific design doc
find docs/dev/design -name "*METADATA*"

# Use doc search (if available)
/doc-search metadata
```

### Updating Documentation

#### Step 1: Locate Relevant Document

```bash
# Design docs are in docs/dev/design/
# Find the right category:
ls docs/dev/design/
# - architecture/
# - features/
# - integrations/
# - services/
# - technical/

# For API changes, update:
docs/api/openapi.yaml
```

#### Step 2: Edit Document

```bash
# Use Zed for fast editing
zed docs/dev/design/services/METADATA.md

# Or VS Code
code docs/dev/design/services/METADATA.md
```

#### Step 3: Validate Changes

```bash
# Validate document structure
python scripts/validate-doc-structure.py

# Check for broken links
python scripts/validate-links.py

# Generate/update indexes
python scripts/generate-design-indexes.py

# Or use doc pipeline
./scripts/doc-pipeline.sh --apply
```

#### Step 4: Commit Documentation

```bash
# Stage changes
git add docs/dev/design/services/METADATA.md

# Commit
git commit -m "docs(services): update metadata enrichment process

- Add section on provider priority
- Document fallback behavior
- Include examples for each provider"

# Push
git push origin feature/update-metadata-docs
```

### Creating New Documentation

```bash
# Use Claude Code skill
/add-design-doc services/NOTIFICATIONS

# Or manually create using template
cp docs/dev/design/01_DESIGN_DOC_TEMPLATE.md \
   docs/dev/design/services/NOTIFICATIONS.md

# Edit the new document
zed docs/dev/design/services/NOTIFICATIONS.md

# Run doc pipeline to update indexes
./scripts/doc-pipeline.sh --apply
```

---

## Emergency Hotfix Workflow

### When to Use

- Critical production bug
- Security vulnerability
- Data corruption issue

### Step 1: Assess Severity

```bash
# Verify the issue
# - Can it wait for normal release cycle?
# - Does it affect production users?
# - Is data at risk?

# If YES to above, proceed with hotfix
# If NO, use normal feature workflow
```

### Step 2: Create Hotfix Branch from Main

```bash
# Checkout main
git checkout main
git pull origin main

# Create hotfix branch
git checkout -b hotfix/v0.1.1-api-crash

# Verify you're on hotfix branch
git branch --show-current
```

### Step 3: Fix the Issue

```bash
# Identify the bug
# Use minimal changes - only fix the critical issue

# Example: Add nil check
# Edit file
zed internal/api/handlers/movie.go

# Test the fix
go test ./internal/api/handlers
go run ./cmd/revenge

# Verify fix works
curl http://localhost:8096/api/v1/movies/12345
```

### Step 4: Commit and Push

```bash
# Commit fix
git add internal/api/handlers/movie.go
git commit -m "fix(api): prevent crash on missing movie

Adds nil check before accessing movie data to prevent panic.

Fixes #789"

# Push hotfix branch
git push origin hotfix/v0.1.1-api-crash
```

### Step 5: Create Emergency PR to Main

```bash
# Create PR with URGENT label
gh pr create \
  --base main \
  --title "HOTFIX: Prevent API crash on missing movie" \
  --label "urgent,hotfix" \
  --body "## Emergency Hotfix

### Issue
API crashes with panic when movie ID doesn't exist.

### Fix
Add nil check before accessing movie data.

### Testing
- [x] Reproduced issue locally
- [x] Verified fix resolves crash
- [x] Tested with missing IDs
- [x] All tests pass

### Impact
Critical - affects all production users.

Fixes #789"
```

### Step 6: Fast-Track Review

```bash
# Request immediate review via Slack/Discord
# Once approved, merge immediately
gh pr merge --squash

# Create hotfix tag
git checkout main
git pull origin main
git tag -a v0.1.1 -m "Hotfix: Prevent API crash on missing movie"
git push origin v0.1.1
```

### Step 7: Deploy Hotfix

```bash
# Deploy to production (follow your deployment process)
# Monitor for any issues

# Create GitHub release
gh release create v0.1.1 \
  --title "Hotfix v0.1.1" \
  --notes "Emergency fix for API crash on missing movie ID."
```

### Step 8: Merge Back to Develop

```bash
# Merge hotfix to develop
git checkout develop
git pull origin develop
git merge main
git push origin develop

# Delete hotfix branch
git branch -D hotfix/v0.1.1-api-crash
git push origin --delete hotfix/v0.1.1-api-crash
```

---

## Best Practices

### General

1. **Always work on feature branches** - Never commit directly to develop or main
2. **Write descriptive commit messages** - Follow conventional commit format
3. **Test before committing** - Run tests locally first
4. **Keep commits focused** - One logical change per commit
5. **Pull before pushing** - Avoid merge conflicts

### Code Review

1. **Review promptly** - Don't let PRs sit for days
2. **Be constructive** - Suggest improvements, don't just criticize
3. **Test locally** - Don't just review code, test it
4. **Approve explicitly** - Use GitHub's review feature

### Testing

1. **Write tests first** - TDD when possible
2. **Aim for 80%+ coverage** - Don't sacrifice quality for coverage
3. **Test edge cases** - Not just happy path
4. **Integration tests matter** - Unit tests aren't enough

### Documentation

1. **Update docs with code** - Don't delay documentation
2. **Keep SOURCE_OF_TRUTH current** - It's the single source of truth
3. **Run validation** - Use doc pipeline before committing
4. **Link related docs** - Use cross-references

---

**Last Updated**: 2026-01-31
**Maintained By**: Development Team

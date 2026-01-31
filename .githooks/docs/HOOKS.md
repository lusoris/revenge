# Git Hooks Reference

> Source: https://git-scm.com/docs/githooks

## Overview

Git hooks are scripts that run automatically at certain points in Git's execution. They must be executable and placed in `.git/hooks/` (or a custom path via `core.hooksPath`).

## Hook Location

```bash
# Default
.git/hooks/

# Custom (recommended for shared hooks)
git config core.hooksPath .githooks
```

## Local Hooks (Client-Side)

### pre-commit

Runs before commit message is created.

```bash
#!/bin/bash
# Prevent commits with TODO/FIXME
if git diff --cached | grep -i "TODO\|FIXME"; then
    echo "Error: Commit contains TODO/FIXME"
    exit 1
fi

# Run linters
npm run lint-staged
```

**Bypass:** `git commit --no-verify`

### prepare-commit-msg

Edit commit message before editor opens.

```bash
#!/bin/bash
# $1 = message file, $2 = source (message|template|merge|squash|commit)
BRANCH=$(git branch --show-current)
if [[ $BRANCH =~ ^([A-Z]+-[0-9]+) ]]; then
    TICKET="${BASH_REMATCH[1]}"
    sed -i "1s/^/[$TICKET] /" "$1"
fi
```

### commit-msg

Validate commit message format.

```bash
#!/bin/bash
MESSAGE=$(cat "$1")

# Conventional commits
if ! echo "$MESSAGE" | grep -qE "^(feat|fix|docs|style|refactor|test|chore)(\(.+\))?: .+"; then
    echo "Error: Commit message must follow conventional commits format"
    echo "Example: feat(auth): add login functionality"
    exit 1
fi
```

**Bypass:** `git commit --no-verify`

### post-commit

Runs after commit is created. Cannot affect outcome.

```bash
#!/bin/bash
echo "Commit created: $(git rev-parse HEAD)"
# Notify external service, update tracking, etc.
```

### pre-push

Runs before push, receives refs via stdin.

```bash
#!/bin/bash
# $1 = remote name, $2 = remote URL

# Prevent push to main
PROTECTED_BRANCH="main"
while read local_ref local_sha remote_ref remote_sha; do
    if [[ "$remote_ref" == "refs/heads/$PROTECTED_BRANCH" ]]; then
        echo "Error: Cannot push directly to $PROTECTED_BRANCH"
        exit 1
    fi
done

# Run tests before push
npm test
```

**Bypass:** `git push --no-verify`

### pre-rebase

Runs before rebase starts.

```bash
#!/bin/bash
# $1 = upstream, $2 = branch (optional)

# Prevent rebase of main
BRANCH=${2:-$(git branch --show-current)}
if [[ "$BRANCH" == "main" ]]; then
    echo "Error: Cannot rebase main branch"
    exit 1
fi
```

### post-checkout

Runs after checkout/switch completes.

```bash
#!/bin/bash
# $1 = prev HEAD, $2 = new HEAD, $3 = branch checkout (1) or file checkout (0)

if [[ "$3" == "1" ]]; then
    # Branch checkout - reinstall dependencies if lockfile changed
    if git diff --name-only "$1" "$2" | grep -q "package-lock.json"; then
        echo "package-lock.json changed, running npm ci..."
        npm ci
    fi
fi
```

### post-merge

Runs after merge completes.

```bash
#!/bin/bash
# $1 = squash flag (1 = squash merge)

# Reinstall dependencies if lockfile changed
if git diff --name-only HEAD@{1} HEAD | grep -q "package-lock.json"; then
    npm ci
fi
```

### post-rewrite

Runs after commits are rewritten (amend, rebase).

```bash
#!/bin/bash
# $1 = command (amend|rebase)
# stdin: old-sha new-sha [extra-info]

echo "Commits rewritten by $1"
```

## Server-Side Hooks

### pre-receive

Runs on remote before any refs are updated.

```bash
#!/bin/bash
# stdin: old-sha new-sha ref-name

while read old_sha new_sha ref; do
    # Reject force push to main
    if [[ "$ref" == "refs/heads/main" ]]; then
        if ! git merge-base --is-ancestor "$old_sha" "$new_sha"; then
            echo "Error: Force push to main not allowed"
            exit 1
        fi
    fi
done
```

### update

Runs once per ref being updated.

```bash
#!/bin/bash
# $1 = ref, $2 = old-sha, $3 = new-sha

REF=$1
OLD=$2
NEW=$3

# Enforce fast-forward only
if [[ "$REF" == "refs/heads/main" ]]; then
    if ! git merge-base --is-ancestor "$OLD" "$NEW"; then
        echo "Error: Non-fast-forward to main rejected"
        exit 1
    fi
fi
```

### post-receive

Runs after all refs updated. Cannot affect outcome.

```bash
#!/bin/bash
# stdin: old-sha new-sha ref-name

while read old_sha new_sha ref; do
    if [[ "$ref" == "refs/heads/main" ]]; then
        # Trigger deployment
        curl -X POST https://api.example.com/deploy
    fi
done
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `GIT_DIR` | Repository directory |
| `GIT_WORK_TREE` | Working tree root |
| `GIT_AUTHOR_NAME` | Commit author name |
| `GIT_AUTHOR_EMAIL` | Commit author email |
| `GIT_COMMITTER_NAME` | Committer name |
| `GIT_COMMITTER_EMAIL` | Committer email |

## Common Patterns

### Run Tests Before Push

```bash
#!/bin/bash
# .githooks/pre-push

echo "Running tests before push..."
npm test || exit 1
go test ./... || exit 1
```

### Lint Staged Files

```bash
#!/bin/bash
# .githooks/pre-commit

# Use lint-staged for efficiency
npx lint-staged

# Or manual approach
FILES=$(git diff --cached --name-only --diff-filter=ACM | grep -E '\.(js|ts)$')
if [[ -n "$FILES" ]]; then
    echo "$FILES" | xargs npx eslint --fix
    echo "$FILES" | xargs git add
fi
```

### Enforce Branch Naming

```bash
#!/bin/bash
# .githooks/pre-push

BRANCH=$(git branch --show-current)
PATTERN="^(feature|bugfix|hotfix|release)/[a-z0-9-]+$"

if ! [[ "$BRANCH" =~ $PATTERN ]]; then
    echo "Error: Branch name must match pattern: $PATTERN"
    echo "Examples: feature/add-auth, bugfix/fix-login"
    exit 1
fi
```

### Block Sensitive Files

```bash
#!/bin/bash
# .githooks/pre-commit

SENSITIVE_FILES=".env credentials.json *.pem *.key"
for pattern in $SENSITIVE_FILES; do
    if git diff --cached --name-only | grep -q "$pattern"; then
        echo "Error: Attempting to commit sensitive file: $pattern"
        exit 1
    fi
done
```

## Setup for Team

```bash
# In project root
mkdir -p .githooks
chmod +x .githooks/*

# Add to README or setup script
git config core.hooksPath .githooks
```

Or use a setup script:

```bash
#!/bin/bash
# scripts/setup-hooks.sh

git config core.hooksPath .githooks
echo "Git hooks configured"
```

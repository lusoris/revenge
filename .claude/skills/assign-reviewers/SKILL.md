---
name: assign-reviewers
description: Configure and assign PR reviewers based on CODEOWNERS and expertise
argument-hint: "[--init|--view|--update-codeowners|--suggest PR_NUMBER] [--auto-assign]"
disable-model-invocation: false
allowed-tools: Bash(cat .github/CODEOWNERS *, python scripts/automation/github_security.py *)
---

# Assign Reviewers

Manage code ownership and reviewer assignments using CODEOWNERS configuration and intelligent suggestions.

## Usage

```
/assign-reviewers --view                       # View current CODEOWNERS
/assign-reviewers --init                       # Initialize CODEOWNERS file
/assign-reviewers --update-codeowners          # Update based on team structure
/assign-reviewers --suggest 123                # Suggest reviewers for PR 123
/assign-reviewers --view --file "api/auth.go"  # Show owners for file
```

## Arguments

- `$0`: Action (--init, --view, --update-codeowners, --suggest)
- `$1+`: Options (PR number, file path, --auto-assign for automatic assignment)

## CODEOWNERS File

**Location**: `.github/CODEOWNERS`

**Purpose**:
- Define code ownership
- Require specific reviewers for areas
- Automate reviewer assignment
- Document responsibility

## Code Ownership Structure

| Area | Owner | Expertise |
|------|-------|-----------|
| Backend (Go) | @lusoris | Go, APIs, architecture |
| Frontend | @lusoris | SvelteKit, TypeScript |
| Database | @lusoris | PostgreSQL, migrations |
| Docs | @lusoris | Documentation, guides |
| CI/CD | @lusoris | GitHub Actions, automation |
| Security | @lusoris | Auth, security features |

## Prerequisites

- GitHub repository
- Admin access to manage CODEOWNERS
- GitHub CLI (`gh`) installed and authenticated
- `.github/CODEOWNERS` file

## Task

Set up and manage code ownership for automatic reviewer assignment.

### Step 1: View Current CODEOWNERS

**View entire CODEOWNERS file**:
```bash
cat .github/CODEOWNERS
```

**Check ownership for file**:
```bash
grep "api/auth" .github/CODEOWNERS
```

**View by area**:
```bash
# Backend ownership
grep "internal/" .github/CODEOWNERS

# Frontend ownership
grep "frontend/" .github/CODEOWNERS

# Database
grep "migrations/" .github/CODEOWNERS
```

### Step 2: Initialize CODEOWNERS

**Create standard CODEOWNERS structure**:
```bash
python scripts/automation/github_security.py --setup-codeowners --init
```

**Creates**:
- Default owners
- By-area ownership
- By-file patterns
- Comment documentation

### Step 3: Update CODEOWNERS

**Add new area**:
1. Edit `.github/CODEOWNERS`
2. Add pattern: `path/to/code/ @owner`
3. Commit and push
4. GitHub auto-uses new rules

**Add new owner**:
```
# API changes require API lead
/internal/api/ @lusoris

# Database changes require DB expert
/migrations/ @lusoris
```

### Step 4: Suggest Reviewers for PR

**Get reviewer suggestions**:
```bash
python scripts/automation/github_security.py --suggest-reviewers --pr-number 123
```

**Shows**:
- Required reviewers (CODEOWNERS)
- Suggested reviewers
- Reviewer availability
- Recent reviewer activity

### Step 5: Auto-Assign Reviewers

**Automatic assignment on PR open**:
```bash
# In GitHub Actions workflow
gh pr edit $PR_NUMBER --add-reviewer @suggested-reviewer
```

**Manual assignment**:
```bash
gh pr edit 123 --add-reviewer @lusoris
```

## CODEOWNERS File Structure

**Current file** (`.github/CODEOWNERS`):

```
# Code Owners
# This file defines who is responsible for code in this repository
# More info: https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-code-owners

# Default owners for everything in the repo
* @lusoris

# Documentation
/docs/ @lusoris
*.md @lusoris

# Protected design documentation (manual changes only)
# Auto-fetcher MUST NOT modify these files
/docs/dev/design/ @lusoris

# Auto-fetched sources (bot can create PRs, requires review)
/docs/dev/sources/ @lusoris

# CI/CD and Infrastructure
/.github/ @lusoris
/Dockerfile @lusoris
/docker-compose*.yml @lusoris
/Makefile @lusoris
/.goreleaser.yml @lusoris

# Configuration
/configs/ @lusoris
/pkg/config/ @lusoris

# Database and migrations
/migrations/ @lusoris
/internal/infra/database/ @lusoris

# API layer
/internal/api/ @lusoris

# Core services
/internal/service/ @lusoris

# Security sensitive files
/SECURITY.md @lusoris
/internal/api/middleware/auth*.go @lusoris
```

## Ownership Patterns

### Exact Path

```
/internal/api/auth.go @lusoris
```
- Matches specific file
- Only one file

### Directory

```
/internal/api/ @lusoris
```
- All files in directory
- Recursive (subdirectories)

### Pattern

```
*.md @lusoris
```
- All markdown files
- Any directory

### Multiple Patterns

```
/internal/api/**/*.go @lusoris
/frontend/**/*.ts @lusoris
```
- Complex glob patterns
- Multiple file types

## Reviewer Assignment Flow

### Automatic (Branch Protection)

1. **PR created** with changes
2. **GitHub checks CODEOWNERS**
3. **Finds matching patterns**
4. **Adds required reviewers**
5. **Blocks merge until reviewed**

### Manual Assignment

1. **Open PR**
2. **Click "Reviewers"**
3. **Select from suggestions**
4. **Assign reviewer**
5. **Notification sent**

## Examples

**View CODEOWNERS**:
```bash
/assign-reviewers --view
```

**Initialize CODEOWNERS**:
```bash
/assign-reviewers --init
```

**Check ownership for file**:
```bash
/assign-reviewers --view --file "internal/api/auth.go"
```

**Suggest reviewers for PR**:
```bash
/assign-reviewers --suggest 123
```

**Update CODEOWNERS**:
```bash
/assign-reviewers --update-codeowners
```

## Managing Reviewers

### Adding Reviewer to PR

**Via GitHub UI**:
1. Click "Reviewers" (right sidebar)
2. Search for user
3. Click to add
4. Notification sent

**Via CLI**:
```bash
gh pr edit 123 --add-reviewer @username
```

### Removing Reviewer

**Via GitHub UI**:
1. Click "Reviewers"
2. Click X next to reviewer

**Via CLI**:
```bash
gh pr edit 123 --remove-reviewer @username
```

### Requesting Re-Review

**After updates**:
1. Make changes based on feedback
2. Click "Re-request review"
3. Notification sent to reviewer

## Reviewer Responsibilities

### Code Review Checklist

- [ ] Understand the changes
- [ ] Check code quality
- [ ] Verify tests pass
- [ ] Check for security issues
- [ ] Verify documentation updated
- [ ] Look for performance issues
- [ ] Suggest improvements
- [ ] Approve or request changes

### Comments

**Request changes**:
- Be specific
- Explain why
- Suggest alternatives
- Be respectful

**Approve**:
- Code looks good
- Tests pass
- No concerns

**Comment**:
- Questions
- Suggestions
- FYI notes

## Expertise Matching

### By Technology

**Go Backend**:
- @lusoris - Lead developer
- Expertise: Architecture, APIs, database

**Frontend**:
- @lusoris - Full stack
- Expertise: SvelteKit, TypeScript, CSS

**Database**:
- @lusoris - DB specialist
- Expertise: PostgreSQL, migrations, optimization

**DevOps**:
- @lusoris - Infrastructure
- Expertise: Docker, CI/CD, Kubernetes

### By Domain

**Authentication**:
- @lusoris - Auth specialist
- Files: `/internal/api/middleware/auth*.go`

**Metadata**:
- @lusoris - Metadata expert
- Files: `/internal/service/metadata/`

**Search**:
- @lusoris - Search expert
- Files: `/internal/service/search/`

## Adding New Team Members

**Update CODEOWNERS**:
1. Edit `.github/CODEOWNERS`
2. Add area: `/path @newuser`
3. Commit: `chore: update CODEOWNERS`
4. Create PR for approval

**Example**:
```
# New team member for frontend
/frontend/ @newuser @existing-reviewer
```

**Multiple reviewers**:
```
# Requires all listed reviewers
/internal/api/middleware/ @lusoris @security-lead
```

## Troubleshooting

**"CODEOWNERS not working"**:
1. Check file path: `.github/CODEOWNERS`
2. Verify syntax (no quotes needed)
3. Ensure branch protection rule enabled
4. Check user exists and has access

**"Reviewer not notified"**:
1. Check reviewer username is correct
2. Verify user is team member
3. Check GitHub notification settings
4. Try re-requesting review

**"Can't require reviewer"**:
1. Need admin access
2. Branch protection rule needed
3. Enable "Require review from Code Owners"
4. CODEOWNERS file must exist

**"Multiple reviewers not working"**:
1. List owners space-separated: `@user1 @user2`
2. Branch protection can require all
3. Only one approval needed by default

## Tips

1. **Keep CODEOWNERS updated**:
   - When adding features
   - When team structure changes
   - Quarterly reviews

2. **Use specific patterns**:
   - More specific = better targeting
   - Auth-related files for security expert
   - API files for API lead

3. **Automate where possible**:
   - CODEOWNERS auto-assigns
   - Branch protection enforces
   - Reduce manual overhead

4. **Document ownership**:
   - Comments in CODEOWNERS
   - Team wiki
   - Decision record

5. **Balance expertise**:
   - Don't overload one person
   - Spread responsibility
   - Build team knowledge

## Best Practices

1. **Clear ownership**:
   - Every area has owner
   - No ambiguity
   - Document rationale

2. **Appropriate experts**:
   - Right person for review
   - Actual expertise needed
   - Rotate reviewers sometimes

3. **Timely reviews**:
   - Set expectations (24-48 hours)
   - Monitor turnaround time
   - Escalate if blocked

4. **Clear feedback**:
   - Explain reasoning
   - Suggest improvements
   - Be constructive

5. **Trust reviewers**:
   - Empower decisions
   - Document disagreements
   - Learn from reviews

## Integration with Branch Protection

**Enforce CODEOWNERS review**:
1. Settings → Branches
2. Edit protection rule
3. Check "Require review from Code Owners"
4. Set required approvals
5. Save

**Combined with other checks**:
```
Required checks:
✓ Code Owner review
✓ CI tests passing
✓ Coverage > 80%
```

## Reviewer Analytics

**Track review times**:
```bash
# Recent PRs and their reviewers
gh pr list --state closed --limit 10 --json number,reviewers,createdAt,closedAt

# Calculate average review time
gh pr list --state closed --limit 30 --json number,reviewers,createdAt,closedAt \
  --jq '.[] | {number, reviewers: .reviewers[].login, days: ((.closedAt | tonumber) - (.createdAt | tonumber)) / 86400}'
```

## Related Skills

- `/configure-branch-protection` - Branch protection rules
- `/manage-labels` - Issue labeling
- `/manage-ci-workflows` - CI/CD workflows
- `/setup-github-projects` - Project management

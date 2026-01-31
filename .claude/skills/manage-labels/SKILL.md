---
name: manage-labels
description: Manage GitHub issue and PR labels for categorization and automation
argument-hint: "[--init|--list|--create|--delete|--sync] [--label NAME] [--color HEX]"
disable-model-invocation: false
allowed-tools: Bash(python scripts/automation/github_labels.py *)
---

# Manage Labels

Create, organize, and manage GitHub issue and PR labels for consistent categorization and workflow automation.

## Usage

```
/manage-labels --list                            # List all labels
/manage-labels --init                            # Create default label set
/manage-labels --create --label "bug" --color ff0000   # Create label
/manage-labels --delete --label "deprecated"     # Delete label
/manage-labels --sync                            # Sync labels with config
```

## Arguments

- `$0`: Action (--init, --list, --create, --delete, --sync)
- `$1+`: Options (--label for label name, --color for hex color)

## Default Label Categories

| Category | Labels | Usage |
|----------|--------|-------|
| Type | task, bug, feature, docs, refactor, chore | Issue type |
| Priority | priority:low, priority:medium, priority:high, priority:critical | Urgency |
| Status | status:blocked, status:in-progress, status:ready | Workflow |
| Area | go, frontend, database, api, infra, docs | Component |
| Effort | effort:small, effort:medium, effort:large | Size |
| QA | testing, tested, qa-blocked | Quality |
| Release | release:breaking, release:enhancement, release:bugfix | Changelog |

## Label Colors

| Category | Color | Hex |
|----------|-------|-----|
| Bug | Red | `d73a49` |
| Feature | Green | `28a745` |
| Documentation | Blue | `0075ca` |
| Priority High | Orange | `ff6b45` |
| In Progress | Purple | `6f42c1` |
| Blocked | Red | `d73a49` |

## Prerequisites

- GitHub CLI (`gh`) installed and authenticated
- Repository admin access
- Python 3.10+ installed

## Task

Create and manage labels for issue/PR organization and automation.

### Step 1: View Current Labels

**List all labels**:
```bash
python scripts/automation/github_labels.py --list
```

**Shows**:
- Label name
- Color
- Usage count
- Description

### Step 2: Initialize Default Labels

**Create standard label set**:
```bash
python scripts/automation/github_labels.py --init
```

**Creates**:
- Type labels (task, bug, feature, docs, etc.)
- Priority labels (low, medium, high, critical)
- Status labels (blocked, in-progress, ready)
- Area labels (go, frontend, database, etc.)
- Effort labels (small, medium, large)
- QA labels (testing, tested, qa-blocked)
- Release labels (breaking, enhancement, bugfix)

### Step 3: Create Custom Label

**Create new label**:
```bash
python scripts/automation/github_labels.py --create --label "experimental" --color "ffaa00"
```

**Create with description**:
```bash
python scripts/automation/github_labels.py --create --label "needs-review" --color "ffd700"
```

### Step 4: Delete Unused Label

**Remove label**:
```bash
python scripts/automation/github_labels.py --delete --label "deprecated"
```

### Step 5: Sync Labels

**Update label definitions**:
```bash
python scripts/automation/github_labels.py --sync
```

## Default Label Set

### Type Labels

**task**:
- Color: `87ceeb` (Sky Blue)
- Description: General task or work item
- Usage: Backlog items

**bug**:
- Color: `d73a49` (Red)
- Description: Something isn't working
- Usage: Bug reports

**feature**:
- Color: `28a745` (Green)
- Description: New feature or enhancement
- Usage: Feature requests

**docs**:
- Color: `0075ca` (Blue)
- Description: Documentation
- Usage: Documentation tasks

**refactor**:
- Color: `fdbcb4` (Light Red)
- Description: Code refactoring
- Usage: Technical debt

**chore**:
- Color: `e4e669` (Yellow)
- Description: Maintenance tasks
- Usage: Regular maintenance

### Priority Labels

**priority:critical**:
- Color: `ff0000` (Red)
- Description: Critical priority
- Usage: Must fix immediately

**priority:high**:
- Color: `ff6b45` (Orange)
- Description: High priority
- Usage: Do soon

**priority:medium**:
- Color: `ffcc00` (Yellow)
- Description: Medium priority
- Usage: Normal priority

**priority:low**:
- Color: `90ee90` (Light Green)
- Description: Low priority
- Usage: Nice to have

### Status Labels

**status:blocked**:
- Color: `d73a49` (Red)
- Description: Work is blocked
- Usage: Waiting on something

**status:in-progress**:
- Color: `6f42c1` (Purple)
- Description: Work in progress
- Usage: Currently being worked on

**status:ready**:
- Color: `28a745` (Green)
- Description: Ready to start
- Usage: No blockers

### Area Labels

**go**:
- Color: `0599a8` (Cyan)
- Description: Go backend code
- Usage: Backend changes

**frontend**:
- Color: `1f6feb` (Blue)
- Description: Frontend code
- Usage: Frontend changes

**database**:
- Color: `6f42c1` (Purple)
- Description: Database/migrations
- Usage: Database changes

**api**:
- Color: `1f6feb` (Blue)
- Description: API changes
- Usage: API design/changes

**infra**:
- Color: `808080` (Gray)
- Description: Infrastructure
- Usage: Deployment/infrastructure

**docs**:
- Color: `0075ca` (Blue)
- Description: Documentation area
- Usage: Doc updates

### Effort Labels

**effort:small**:
- Color: `90ee90` (Light Green)
- Description: Small effort
- Usage: Quick fixes (1-2 hours)

**effort:medium**:
- Color: `ffcc00` (Yellow)
- Description: Medium effort
- Usage: 1-2 days

**effort:large**:
- Color: `ff6b45` (Orange)
- Description: Large effort
- Usage: 3+ days

### QA Labels

**testing**:
- Color: `fbca04` (Yellow)
- Description: Needs testing
- Usage: Testing required

**tested**:
- Color: `28a745` (Green)
- Description: Tested and verified
- Usage: QA approved

**qa-blocked**:
- Color: `d73a49` (Red)
- Description: QA blocked
- Usage: Testing issue

### Release Labels

**release:breaking**:
- Color: `d73a49` (Red)
- Description: Breaking change
- Usage: Major version bump

**release:enhancement**:
- Color: `28a745` (Green)
- Description: Feature enhancement
- Usage: Minor version bump

**release:bugfix**:
- Color: `fdbcb4` (Light Red)
- Description: Bug fix
- Usage: Patch version bump

## Examples

**List all labels**:
```bash
/manage-labels --list
```

**Create standard label set**:
```bash
/manage-labels --init
```

**Create custom label**:
```bash
/manage-labels --create --label "windows-only" --color "0078d4"
```

**Delete unused label**:
```bash
/manage-labels --delete --label "old-label"
```

**Sync label configuration**:
```bash
/manage-labels --sync
```

## Using Labels in Workflows

### Auto-Labeling PRs

**Based on file changes**:
```yaml
# .github/workflows/auto-label.yml
on: pull_request

jobs:
  auto-label:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/labeler@v4
        with:
          configuration-path: .github/labeler.yml
```

**Label configuration** (`.github/labeler.yml`):
```yaml
frontend:
  - changed-files:
      - any-glob-to-any-file:
          - 'frontend/**'

go:
  - changed-files:
      - any-glob-to-any-file:
          - 'internal/**'
          - 'cmd/**'

database:
  - changed-files:
      - any-glob-to-any-file:
          - 'migrations/**'
```

### Label-Based Automation

**Add to project based on label**:
```yaml
on:
  issues:
    types: [opened, labeled]

jobs:
  add-to-project:
    runs-on: ubuntu-latest
    steps:
      - uses: github/project-automation@v2
        with:
          project-name: Development
          github-token: ${{ secrets.GITHUB_TOKEN }}
```

## Label Naming Conventions

### Naming Scheme

**Format**: `category:subcategory` or `category`

**Examples**:
- `priority:high`
- `type:bug`
- `area:database`
- `status:blocked`

**Benefits**:
- Namespace labels logically
- Easier to find related labels
- Supports grouping in UI
- Better for filtering

### Consistency

**Maintain consistency**:
- Use lowercase only
- Use hyphens for spaces
- Consistent prefixes
- Clear descriptions

## Common Label Workflows

### Triaging New Issues

1. **Add type label**: bug, feature, task, docs
2. **Add priority**: critical, high, medium, low
3. **Add area**: go, frontend, database, api
4. **Set status**: ready, blocked (if applicable)

**Example combination**:
- `type:bug` + `priority:high` + `area:api` + `status:ready`

### Tracking Work

1. **In backlog**: `status:ready`
2. **Starting work**: `status:in-progress`
3. **Under review**: `status:in-review`
4. **Complete**: `status:done`

### Release Planning

1. **Mark changes**: `release:breaking`, `release:enhancement`, `release:bugfix`
2. **Group by release**: Add to milestone
3. **Generate changelog**: From labels and milestone

## Bulk Operations

### Update Multiple Issues

**Add label to issues**:
```bash
gh issue list --label "type:bug" --json number | jq '.[] | .number' | \
  xargs -I {} gh issue edit {} --add-label "priority:high"
```

**Remove label**:
```bash
gh issue list --label "old-label" --json number | jq '.[] | .number' | \
  xargs -I {} gh issue edit {} --remove-label "old-label"
```

### Export Label Usage

**See which issues use which labels**:
```bash
gh issue list --json number,labels --jq '.[] | "\(.number): \(.labels | map(.name) | join(", "))"'
```

## Troubleshooting

**"Label not appearing on issue"**:
1. Verify label name spelling
2. Check user permissions
3. Try re-adding label
4. Refresh browser

**"Can't create label"**:
1. Check permissions (admin required)
2. Verify label doesn't exist
3. Check name validity
4. Try different name if taken

**"Label color not correct"**:
1. Edit label in Settings
2. Update hex color
3. Verify valid hex format (6 digits)
4. Save changes

**"Too many labels causing clutter"**:
1. Archive unused labels
2. Consolidate similar labels
3. Use label filtering
4. Create label standards

## Tips

1. **Keep label set focused**:
   - 30-50 labels maximum
   - Clear naming conventions
   - Regular cleanup of unused

2. **Use colors strategically**:
   - Type: Blue/red
   - Priority: Red/yellow
   - Status: Purple/green
   - Area: Varied

3. **Combine labels for richness**:
   - One type label
   - One priority label
   - One area label
   - Additional as needed

4. **Auto-label when possible**:
   - File-based labeling
   - Branch-based labeling
   - Reduces manual work

5. **Regular review**:
   - Monthly label audit
   - Remove unused labels
   - Update descriptions
   - Maintain consistency

## Best Practices

1. **Consistent naming**:
   - `type:`, `priority:`, `area:`, `status:`, `effort:` prefixes
   - Lowercase only
   - Use hyphens not spaces

2. **Clear descriptions**:
   - Each label has purpose
   - Visible when hovering
   - Helps new contributors

3. **Limit label count per issue**:
   - Max 5-7 labels
   - One type + one priority minimum
   - Others as needed

4. **Keep updated**:
   - Remove obsolete labels
   - Update descriptions
   - Maintain color scheme consistency

5. **Use for automation**:
   - Label triggers project additions
   - Label triggers workflow actions
   - Label affects CI/CD behavior

## Related Skills

- `/setup-github-projects` - Project management
- `/manage-milestones` - Milestone management
- `/manage-ci-workflows` - CI/CD workflows

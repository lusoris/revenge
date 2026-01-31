---
name: setup-github-projects
description: Set up GitHub Projects for task/issue tracking and workflow automation
argument-hint: "[--init|--list|--view PROJECT_NAME] [--template kanban|table|roadmap]"
disable-model-invocation: false
allowed-tools: Bash(python scripts/automation/github_projects.py *)
---

# Setup GitHub Projects

Initialize and manage GitHub Projects (v2) for task tracking, sprint planning, and workflow automation.

## Usage

```
/setup-github-projects --list                      # List all projects
/setup-github-projects --init                      # Create default projects
/setup-github-projects --init --template kanban    # Create with kanban template
/setup-github-projects --view "Development"        # View project details
/setup-github-projects --view "Releases"           # View release project
```

## Arguments

- `$0`: Action (--init to create, --list to show, --view to inspect)
- `$1+`: Options (--template for template type, project name for view)

## Default Projects

| Project | Purpose | Template | Automation |
|---------|---------|----------|-----------|
| Development | Current sprint tasks | Table | Auto-add PRs/issues |
| Releases | Release planning | Roadmap | Milestone linked |
| Bug Triage | Bug tracking | Kanban | Auto-add bugs |
| Documentation | Docs tasks | Table | By label |

## Prerequisites

- GitHub CLI (`gh`) installed and authenticated
- Admin access to repository
- GitHub Projects v2 enabled

## Task

Set up GitHub Projects for organized task tracking and workflow management.

### Step 1: List Current Projects

**View all projects**:
```bash
python scripts/automation/github_projects.py --list
```

**Output shows**:
- Project names
- Number of issues
- Status fields
- Automation configured

### Step 2: Initialize Default Projects

**Create standard project structure**:
```bash
python scripts/automation/github_projects.py --init
```

**Creates**:
1. Development project (table view)
2. Releases project (roadmap view)
3. Bug Triage project (kanban view)
4. Documentation project (table view)

### Step 3: View Project Details

**View specific project**:
```bash
python scripts/automation/github_projects.py --view "Development"
```

**Shows**:
- Columns/statuses
- Custom fields
- Automation rules
- Number of items

### Step 4: Configure Automation

**Auto-add new issues to Development**:
```bash
# Edit project automation in GitHub UI
# Settings → Automation → Add items
```

**Auto-link with milestones**:
```bash
# Issues with milestone added to Releases project
```

## Project Details

### Development Project

**Type**: Table view

**Columns**:
- Backlog
- Todo
- In Progress
- In Review
- Done

**Fields**:
- Status (select)
- Priority (High/Medium/Low)
- Assignee (person)
- Due date
- Labels

**Automation**:
- Auto-add issues with label `type:task`
- Auto-add PRs to "In Review"

### Releases Project

**Type**: Roadmap view

**Columns**:
- Backlog
- Planning
- In Development
- Testing
- Released

**Fields**:
- Status
- Release version
- Milestone
- Due date

**Automation**:
- Link to milestones automatically

### Bug Triage Project

**Type**: Kanban view

**Columns**:
- New
- Triaging
- Backlog
- In Progress
- Done

**Fields**:
- Status
- Priority
- Severity
- Assignee
- Due date

**Automation**:
- Auto-add issues with label `type:bug`

### Documentation Project

**Type**: Table view

**Columns**:
- Backlog
- Writing
- Review
- Published

**Fields**:
- Status
- Type (doc/guide/api-ref)
- Assignee
- Due date

**Automation**:
- Auto-add issues with label `type:docs`

## Examples

**List all projects**:
```bash
/setup-github-projects --list
```

**Set up default projects**:
```bash
/setup-github-projects --init
```

**View Development project**:
```bash
/setup-github-projects --view "Development"
```

**Create with kanban template**:
```bash
/setup-github-projects --init --template kanban
```

## Managing Projects

### Add Item to Project

**Via GitHub CLI**:
```bash
# Add issue to project
gh project item-add 1 --id ISSUE_ID

# Add PR to project
gh project item-add 1 --id PR_ID
```

**Via automation**:
- Label-based (label triggers auto-add)
- Milestone-based (milestone triggers auto-add)
- PR-based (PR automatically added to "In Review")

### Move Item Between Columns

**Via CLI**:
```bash
gh project item-edit ITEM_ID --field "Status" --single-select-option "In Progress"
```

**Via GitHub UI**:
- Drag and drop in project view
- Right-click context menu
- Column dropdown

### Update Item Fields

**Priority**:
```bash
gh project item-edit ITEM_ID --field "Priority" --single-select-option "High"
```

**Due date**:
```bash
gh project item-edit ITEM_ID --field "Due date" --date "2026-02-28"
```

**Assignee**:
```bash
gh project item-edit ITEM_ID --field "Assignee" --text "@username"
```

## Automation Rules

### Label-Based Automation

**Add issues with specific labels to Development**:
1. Go to project → Automation
2. Create rule: "Add items when they're labeled"
3. Select label: `type:task`
4. Target: Development project

### Milestone-Based Automation

**Add to Releases project by milestone**:
1. Create milestone: `v1.1.0`
2. Issues with this milestone auto-add to Releases
3. Status field shows release phase

### PR Automation

**PRs automatically added to Development**:
- PR opened → Added to "In Review" column
- PR merged → Moved to "Done" column

## Views and Reports

### Table View (Development)

**Best for**:
- Sprint tracking
- Task management
- Custom fields display

**Columns shown**:
- Issue number
- Title
- Status
- Priority
- Assignee
- Due date

### Kanban View (Bug Triage)

**Best for**:
- Bug tracking
- Workflow visualization
- Column-based progress

**Columns**:
- New → Triaging → Backlog → In Progress → Done

### Roadmap View (Releases)

**Best for**:
- Release planning
- Timeline visualization
- Milestone tracking

**Timeline shows**:
- Release dates
- Feature grouping
- Progress

## Filtering and Searching

### Filter by Status

```bash
# Show items in "In Progress"
# In project: Status = "In Progress"
```

### Filter by Assignee

```bash
# Show items assigned to me
# In project: Assignee = "@myusername"
```

### Filter by Priority

```bash
# Show high priority items
# In project: Priority = "High"
```

### Sort Options

- Sort by due date
- Sort by priority
- Sort by assignee
- Sort by creation date

## Integration with Issues

### Issue Labels

**Development labels**:
- `type:task` → Development project
- `type:bug` → Bug Triage project
- `type:docs` → Documentation project
- `priority:high` → High priority field
- `priority:low` → Low priority field

### Issue Milestones

**Milestone integration**:
- Milestone set → Issue added to Releases project
- Release notes generated from milestone
- Status shows release phase

### Issue Templates

**Link to projects**:
- Tasks created from template → Development
- Bugs created from template → Bug Triage
- Docs from template → Documentation

## Troubleshooting

**"Project not found"**:
1. Verify project exists: `/setup-github-projects --list`
2. Check project name spelling
3. Ensure you have access
4. Try re-creating: `/setup-github-projects --init`

**"Can't add items to project"**:
1. Check automation rules enabled
2. Verify labels configured
3. Check issue meets automation criteria
4. Try manual add via CLI

**"Automation not working"**:
1. Review automation rules
2. Check label names match
3. Verify issue has label when created
4. Re-configure automation

**"Columns not showing**:
1. Check project type (table vs kanban)
2. Verify field configuration
3. Try refreshing browser
4. Check custom field setup

## Tips

1. **Use projects for sprint planning**:
   - Development project for current sprint
   - Move items through columns
   - Track progress

2. **Link issues to projects**:
   ```bash
   # When creating issue
   gh issue create --title "New task" --label type:task
   ```

3. **Regular triage**:
   - Weekly triage of bug project
   - Move bugs to backlog or in progress
   - Update priority and assignee

4. **View for stakeholders**:
   - Roadmap project for high-level planning
   - Share status with team
   - Track release progress

5. **Automated filtering**:
   - Use saved filters
   - Quick access to high priority
   - Track personal assignments

## Performance Notes

- Projects v2 is fast even with 1000+ items
- Filters apply in real-time
- Automation triggers within minutes
- Search indexes all items

## Best Practices

1. **Use consistent labels**:
   - Standard types (task, bug, docs, feature)
   - Standard priorities (high, medium, low)
   - Consistent naming

2. **Keep projects updated**:
   - Move items through columns
   - Close/archive done items
   - Update fields regularly

3. **Regular reviews**:
   - Weekly sprint reviews
   - Monthly roadmap updates
   - Quarterly backlog grooming

4. **Automate where possible**:
   - Labels trigger additions
   - PRs auto-move through review
   - Milestones link to releases

## Related Skills

- `/manage-labels` - Manage issue labels
- `/manage-milestones` - Manage milestones
- `/manage-ci-workflows` - CI/CD workflows

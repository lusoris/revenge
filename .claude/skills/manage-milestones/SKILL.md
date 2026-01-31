---
name: manage-milestones
description: Manage GitHub milestones for release planning and version tracking
argument-hint: "[--init|--list|--create VERSION] [--close VERSION] [--view VERSION]"
disable-model-invocation: false
allowed-tools: Bash(python scripts/automation/github_milestones.py *)
---

# Manage Milestones

Create and manage GitHub milestones for version planning, release coordination, and progress tracking.

## Usage

```
/manage-milestones --list                       # List all milestones
/manage-milestones --init                       # Create default milestones
/manage-milestones --create v1.1.0              # Create new milestone
/manage-milestones --view v1.0.0                # View milestone details
/manage-milestones --close v1.0.0               # Close completed milestone
```

## Arguments

- `$0`: Action (--init, --list, --create, --view, --close)
- `$1+`: Options (version for specific milestone)

## Default Milestones

| Milestone | Purpose | Status |
|-----------|---------|--------|
| v1.0.0 | Initial release | Planning |
| v1.1.0 | Next minor release | Planning |
| v2.0.0 | Major version | Planning |
| Backlog | Future work | Planning |
| In Progress | Current sprint | In Progress |

## Milestone Workflow

| Phase | Status | Activity |
|-------|--------|----------|
| Planning | Open | Create issues, plan work |
| Active | Open | Assign issues, track progress |
| Release | Open | Bug fixes, release prep |
| Released | Closed | Archive, document |

## Prerequisites

- GitHub repository
- Admin or write access
- GitHub CLI (`gh`) installed and authenticated
- Python 3.10+ installed

## Task

Create and manage milestones for release planning and tracking.

### Step 1: View Current Milestones

**List all milestones**:
```bash
python scripts/automation/github_milestones.py --list
```

**Shows**:
- Milestone name
- Description
- Due date
- Open/closed issue count
- Progress percentage

### Step 2: Initialize Default Milestones

**Create standard milestone structure**:
```bash
python scripts/automation/github_milestones.py --init
```

**Creates**:
- v1.0.0 - Initial release
- v1.1.0 - Next minor release
- v2.0.0 - Major version
- Backlog - Future work
- In Progress - Current sprint

### Step 3: Create New Milestone

**Create new version milestone**:
```bash
python scripts/automation/github_milestones.py --create v1.2.0
```

**With description and due date**:
```bash
python scripts/automation/github_milestones.py --create v1.2.0 --description "Minor feature release" --due-date "2026-03-31"
```

### Step 4: View Milestone Details

**Show milestone information**:
```bash
python scripts/automation/github_milestones.py --view v1.0.0
```

**Shows**:
- All issues in milestone
- Issue status breakdown
- Progress percentage
- Due date
- Description

### Step 5: Close Completed Milestone

**Mark milestone as complete**:
```bash
python scripts/automation/github_milestones.py --close v1.0.0
```

**Creates release notes and archives**

## Milestone Details

### Release Milestones

**v1.0.0 (Initial Release)**:
- Description: "Initial product release"
- Due date: TBD
- Type: Major release
- Status: Planning/Active
- Features: Core functionality

**v1.1.0 (Next Minor)**:
- Description: "Minor feature release"
- Due date: TBD
- Type: Minor release
- Status: Planning
- Features: Enhancements to v1.0

**v2.0.0 (Major Version)**:
- Description: "Major redesign/upgrade"
- Due date: TBD
- Type: Major release
- Status: Planning
- Features: Breaking changes

### Work Milestones

**In Progress (Sprint)**:
- Description: "Current sprint"
- Target duration: 2 weeks
- Issues: Current work
- Review: Weekly

**Backlog (Future)**:
- Description: "Future work items"
- Type: Collection
- Status: Planning
- Action: Move to release milestone when ready

## Using Milestones

### Link Issues to Milestone

**Set milestone when creating issue**:
```bash
gh issue create --title "New feature" --milestone v1.1.0 --label "type:feature"
```

**Add existing issue to milestone**:
```bash
gh issue edit 123 --milestone v1.1.0
```

**Bulk add issues**:
```bash
gh issue list --label "priority:high" --json number | jq '.[] | .number' | \
  xargs -I {} gh issue edit {} --milestone v1.1.0
```

### Track Progress

**View milestone progress**:
1. Repository → Projects → Milestones
2. Click milestone name
3. See progress bar
4. View open/closed issues

**Percentage complete**:
- Calculated: Closed / (Closed + Open)
- Shows progress toward release

### Release from Milestone

**Generate release notes**:
```bash
gh release create v1.0.0 \
  --title "Version 1.0.0" \
  --notes "$(gh api repos/:owner/:repo/milestones/1 --jq '.description')"
```

## Milestone Planning

### Planning Phase

1. **Create milestone**:
   - Version number (semantic)
   - Description
   - Due date

2. **Define scope**:
   - What features?
   - What fixes?
   - What's excluded?

3. **Create issues**:
   - Break down into tasks
   - Assign to milestone
   - Set priority and effort

### Active Phase

1. **Assign issues**:
   - Team members pick up work
   - Move to "In Progress" status
   - Update regularly

2. **Monitor progress**:
   - Track closure rate
   - Identify blockers
   - Adjust scope if needed

3. **Manage risks**:
   - Watch for slipping
   - Re-prioritize if needed
   - Communicate timeline

### Release Phase

1. **Code freeze**:
   - Only bug fixes
   - No new features
   - Testing focus

2. **Testing**:
   - QA testing
   - Integration testing
   - Regression testing

3. **Release**:
   - Create release tag
   - Build release artifacts
   - Publish release notes
   - Close milestone

## Integration with Release Please

**Automatic milestone creation**:
- Release Please creates release PRs
- Links to milestone
- Generates changelog from issues

**Workflow**:
1. Issue created → assigned to milestone
2. PR created → references issue
3. PR merged → closes issue
4. Milestone closed → triggers release

## Example Release Timeline

### v1.0.0 Release Planning

**Week 1-2**: Planning
```
- Create v1.0.0 milestone
- Create 20-30 issues
- Estimate effort for each
- Assign to team members
```

**Week 3-6**: Development
```
- Team works on assigned issues
- Daily standup review
- Move issues through statuses
- Manage blockers
```

**Week 7**: Code Freeze
```
- No new features
- Bug fixes only
- Testing focus
```

**Week 8**: Release
```
- Final testing
- Release documentation
- Create GitHub release
- Close milestone
```

## Tracking Progress

### Issue Counts

**Total issues**:
```
Open: 5
Closed: 15
Total: 20
Progress: 75% complete
```

### Burndown Chart

**Track closure over time**:
- Plot remaining issues
- X-axis: Time
- Y-axis: Issues remaining
- Trend line shows slope

**Manual tracking**:
```bash
# Check progress daily
gh api repos/:owner/:repo/milestones/:milestone_number \
  --jq '{open: .open_issues, closed: .closed_issues}'
```

## Troubleshooting

**"Milestone not showing issues"**:
1. Verify issues assigned to milestone
2. Check issue is in correct repo
3. Refresh browser
4. Try filtering

**"Can't add issue to milestone"**:
1. Verify write access
2. Check milestone exists
3. Ensure issue is in same repo
4. Try via API

**"Due date not working"**:
1. Use valid date format (YYYY-MM-DD)
2. Date must be future
3. Check timezone handling
4. Verify date in GitHub UI

**"Can't close milestone"**:
1. Close all issues first (or keep open)
2. Verify you have permission
3. Create release first
4. Try deleting if needed

## Best Practices

1. **Use semantic versioning**:
   - v1.0.0 = major.minor.patch
   - v1.0.0-beta.1 = prerelease
   - Matches Release Please output

2. **Set realistic due dates**:
   - Buffer for unknowns
   - Account for testing
   - Leave room for issues

3. **Regular scope reviews**:
   - Weekly progress check
   - Identify slipping
   - Re-prioritize as needed
   - Cut low-priority items if needed

4. **Link releases to milestones**:
   - Track what went into each version
   - Generate comprehensive release notes
   - Identify what's in production

5. **Archive old milestones**:
   - Keep recent (2-3) visible
   - Archive older releases
   - Use for reference

## Tips

1. **Create next milestone early**:
   - Plan as current milestone nears completion
   - Gather feedback
   - Prioritize items
   - Set due date

2. **Use milestone for project communication**:
   - Public view of roadmap
   - Transparency on progress
   - Shows commitment
   - Sets expectations

3. **Combine with labels**:
   - `release:breaking` issues go in major milestone
   - `release:feature` goes in minor milestone
   - Easy filtering and categorization

4. **Milestone checklist**:
   ```
   - [ ] Milestone created
   - [ ] Description written
   - [ ] Due date set
   - [ ] Issues created
   - [ ] Issues prioritized
   - [ ] Work assigned
   - [ ] Tracking started
   - [ ] Regular reviews scheduled
   - [ ] Release date confirmed
   - [ ] Release notes prepared
   ```

5. **Monitor velocity**:
   - How many issues per sprint?
   - What's the trend?
   - Plan future milestones accordingly

## Milestone Dashboard

**Create tracking spreadsheet**:
```
| Milestone | Target Date | Issues | Closed | % Done | On Track |
|-----------|-------------|--------|--------|--------|----------|
| v1.0.0    | 2026-03-01  | 20     | 15     | 75%    | Yes      |
| v1.1.0    | 2026-04-15  | 15     | 0      | 0%     | Planning |
```

## Integration with GitHub API

**Query milestone info**:
```bash
gh api repos/:owner/:repo/milestones --jq '.[] | {title: .title, open: .open_issues, closed: .closed_issues}'
```

**Create milestone with API**:
```bash
gh api repos/:owner/:repo/milestones -f title="v1.1.0" -f description="Next release" -f due_on="2026-04-30T23:59:59Z"
```

## Related Skills

- `/configure-release-please` - Automated versioning
- `/manage-labels` - Issue labeling
- `/setup-github-projects` - Project tracking
- `/manage-ci-workflows` - CI/CD workflows

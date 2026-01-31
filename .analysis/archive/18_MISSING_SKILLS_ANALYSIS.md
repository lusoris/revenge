# Missing Skills Analysis - Complete System Coverage

**Created**: 2026-01-31
**Purpose**: Identify all missing Claude Code skills needed for complete automation system coverage
**Based On**: 17_CLAUDE_SKILLS_SPECIFICATION.md + 15_FINAL_ANSWERS_SUMMARY.md

---

## Overview

**Skills Specified**: 6 (doc automation only)
**Skills Missing**: ~15-20 (GitHub management, release, config, monitoring)

**Gap**: We have comprehensive doc automation skills, but missing skills for:
1. GitHub project management (Projects, Discussions, Issues, PRs)
2. Dependency & release management (Dependabot, Release Please)
3. Code quality automation (linting, formatting, testing)
4. Infrastructure management (Coder, Docker, CI/CD)
5. Monitoring & observability (logs, metrics, alerts)

---

## Category 1: GitHub Project Management (7 skills needed)

### Skill 7: setup-github-projects

**Purpose**: Configure GitHub Projects integration with automation rules

**Usage**:
```bash
/setup-github-projects
/setup-github-projects --project "Revenge v1.0"
```

**Features**:
- Create project board with columns (Backlog, Todo, In Progress, Review, Done)
- Configure automation rules:
  - Auto-add issues to project
  - Auto-move cards on PR open/merge
  - Auto-assign to milestones
- Set up project views (Board, Table, Roadmap)
- Configure custom fields (Priority, Effort, Module)

**Interactive Flow**:
1. Ask project name
2. Ask board structure (default or custom columns)
3. Ask automation rules to enable
4. Create project via GitHub API
5. Configure automation
6. Show project URL

---

### Skill 8: setup-github-discussions

**Purpose**: Configure GitHub Discussions with categories and templates

**Usage**:
```bash
/setup-github-discussions
```

**Features**:
- Enable Discussions on repo
- Create categories:
  - 💡 Ideas (for feature requests)
  - ❓ Q&A (for questions)
  - 📢 Announcements (for updates)
  - 🐛 Bugs (auto-convert to issues)
- Create discussion templates
- Configure auto-convert rules (discussion → issue when labeled)

---

### Skill 9: configure-branch-protection

**Purpose**: Set up branch protection rules from SOT configuration

**Usage**:
```bash
/configure-branch-protection
/configure-branch-protection --branch develop --strict
```

**Features**:
- Configure protection for `develop` and `main`
- Rules from automation-config.yml:
  - Require pull request reviews (1+ approvals)
  - Require status checks pass (CI, lint, tests)
  - Require linear history
  - No force push
  - Include administrators
- Show current protection status
- Validate protection active

---

### Skill 10: setup-codeql

**Purpose**: Enable and configure GitHub Advanced Security with CodeQL

**Usage**:
```bash
/setup-codeql
/setup-codeql --languages go,javascript
```

**Features**:
- Enable GitHub Advanced Security (if available)
- Configure CodeQL analysis for Go + JavaScript
- Set up automated scanning on push/PR
- Configure security alerts
- Set up dependency review
- Show security overview

---

### Skill 11: manage-labels

**Purpose**: Synchronize GitHub labels from configuration

**Usage**:
```bash
/manage-labels
/manage-labels --sync
/manage-labels --add "type: feature" --color "#0E8A16"
```

**Features**:
- Load label config from `.github/labels.yml`
- Create missing labels
- Update existing labels (name, color, description)
- Delete labels not in config
- Show label usage stats

**Label Categories**:
- Type: `type: feature`, `type: bug`, `type: docs`
- Priority: `priority: critical`, `priority: high`, `priority: low`
- Size: `size: S`, `size: M`, `size: L`, `size: XL`
- Status: `status: blocked`, `status: wip`
- Module: `module: video`, `module: music`, `module: auth`

---

### Skill 12: assign-reviewers

**Purpose**: Auto-assign reviewers based on CODEOWNERS and PR changes

**Usage**:
```bash
/assign-reviewers --pr 123
/assign-reviewers --auto  # Set up GitHub Action
```

**Features**:
- Parse CODEOWNERS file
- Match changed files to owners
- Assign reviewers via GitHub API
- Configure auto-assignment GitHub Action
- Handle team assignments

---

### Skill 13: manage-milestones

**Purpose**: Create and manage GitHub milestones with automation

**Usage**:
```bash
/manage-milestones
/manage-milestones --create "v0.3.0" --due 2026-03-01
/manage-milestones --close v0.2.0 --move-open-to v0.3.0
```

**Features**:
- Create milestone with due date
- Auto-assign issues based on labels/priority
- Close milestone when all issues done
- Move open issues to next milestone
- Generate milestone report (completed, open, closed issues)

---

## Category 2: Dependency & Release Management (3 skills needed)

### Skill 14: configure-dependabot

**Purpose**: Set up and manage Dependabot configuration

**Usage**:
```bash
/configure-dependabot
/configure-dependabot --check
```

**Features**:
- Generate `.github/dependabot.yml` from SOT
- Configure package ecosystems:
  - Go modules (go.mod)
  - npm (package.json)
  - Python (requirements.txt)
  - GitHub Actions (.github/workflows/)
- Set update schedule (weekly, daily)
- Configure auto-merge rules
- Set version bump limits (patch, minor, major)
- Configure ignore conditions

**Example Config**:
```yaml
version: 2
updates:
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 10
    reviewers:
      - "revenge-bot"
    labels:
      - "dependencies"
      - "go"

  - package-ecosystem: "npm"
    directory: "/"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 5
    reviewers:
      - "revenge-bot"
    labels:
      - "dependencies"
      - "javascript"
```

---

### Skill 15: configure-release-please

**Purpose**: Set up Release Please for automated releases

**Usage**:
```bash
/configure-release-please
/configure-release-please --check
```

**Features**:
- Generate `.github/release-please-config.json`
- Configure release types (Go, Node monorepo)
- Set up changelog generation
- Configure release branches (main, develop)
- Set up GitHub Actions workflow
- Test release please locally

**Example Config**:
```json
{
  "packages": {
    ".": {
      "release-type": "go",
      "package-name": "revenge",
      "changelog-path": "CHANGELOG.md",
      "bump-minor-pre-major": true,
      "bump-patch-for-minor-pre-major": true,
      "extra-files": [
        "docs/dev/design/00_SOURCE_OF_TRUTH.md"
      ]
    },
    "web": {
      "release-type": "node",
      "package-name": "@revenge/web",
      "changelog-path": "web/CHANGELOG.md"
    }
  }
}
```

---

### Skill 16: update-dependencies

**Purpose**: Manual dependency update with validation

**Usage**:
```bash
/update-dependencies
/update-dependencies --check  # Check for updates
/update-dependencies --package fx --version v1.24.0
```

**Features**:
- Check for available updates (Go, npm, Python)
- Show current vs latest versions
- Update specific package
- Update SOT when package updated
- Trigger doc regeneration
- Run tests after update
- Create PR with update

---

## Category 3: Code Quality Automation (4 skills needed)

### Skill 17: run-linters

**Purpose**: Run all linters with auto-fix option

**Usage**:
```bash
/run-linters
/run-linters --fix
/run-linters --type go,python
```

**Features**:
- Run golangci-lint (Go)
- Run ruff (Python)
- Run markdownlint (Docs)
- Run prettier (JavaScript/TypeScript)
- Parallel execution
- Auto-fix when possible
- Show summary report
- Create GitHub check if in CI

**Output**:
```
🔍 Running linters...

[1/4] golangci-lint (Go)
✅ 0 issues found

[2/4] ruff (Python)
⚠️  3 issues found (2 auto-fixed)
  scripts/automation/sot_parser.py:45:1: E501 Line too long (92 > 88)

[3/4] markdownlint (Docs)
✅ 0 issues found

[4/4] prettier (TypeScript)
✅ 0 issues found

Summary:
- Total issues: 3
- Auto-fixed: 2
- Remaining: 1

Next steps:
1. Fix remaining issue: scripts/automation/sot_parser.py:45
2. Commit fixes: git add . && git commit -m "lint: auto-fix issues"
```

---

### Skill 18: run-tests

**Purpose**: Run all test suites with coverage reporting

**Usage**:
```bash
/run-tests
/run-tests --type unit
/run-tests --coverage
/run-tests --watch
```

**Features**:
- Run Go tests (unit + integration)
- Run Python tests
- Run frontend tests (Vitest)
- Parallel execution
- Coverage reporting (with 80% threshold check)
- Watch mode for development
- Generate coverage reports (HTML, JSON)

**Output**:
```
🧪 Running tests...

[1/3] Go Tests
  Unit tests: 245 passed, 0 failed
  Integration tests: 67 passed, 0 failed
  Coverage: 84.2% (threshold: 80%)
  Time: 12.3s

[2/3] Python Tests
  Tests: 42 passed, 0 failed
  Coverage: 91.5%
  Time: 2.1s

[3/3] Frontend Tests
  Tests: 156 passed, 0 failed
  Coverage: 78.3% ⚠️  (below threshold)
  Time: 8.7s

Summary:
- Total tests: 510
- Passed: 510
- Failed: 0
- Coverage: 82.1% (overall)
- Time: 23.1s

⚠️  Frontend coverage below 80% threshold
Recommendation: Add tests for uncovered components
```

---

### Skill 19: format-code

**Purpose**: Auto-format all code with standard formatters

**Usage**:
```bash
/format-code
/format-code --check  # Dry run
/format-code --type go
```

**Features**:
- Run gofmt (Go)
- Run ruff format (Python)
- Run prettier (JavaScript/TypeScript/JSON/YAML/Markdown)
- Show files changed
- Commit changes (optional)

---

### Skill 20: check-licenses

**Purpose**: Check dependency licenses for compliance

**Usage**:
```bash
/check-licenses
/check-licenses --policy strict
```

**Features**:
- Scan Go dependencies
- Scan npm dependencies
- Check against allow/deny list (from SOT)
- Show license summary
- Flag incompatible licenses
- Generate license report

**Allow List** (from SOT):
- MIT, Apache-2.0, BSD-3-Clause, BSD-2-Clause
- ISC, MPL-2.0

**Deny List**:
- GPL-3.0, AGPL-3.0 (viral licenses)

---

## Category 4: Infrastructure Management (3 skills needed)

### Skill 21: manage-coder-workspace

**Purpose**: Manage Coder workspace templates and instances

**Usage**:
```bash
/manage-coder-workspace
/manage-coder-workspace --create
/manage-coder-workspace --update
```

**Features**:
- Update Coder template (.coder/template.tf)
- Sync tool versions from SOT
- Create workspace
- Start/stop workspace
- Show workspace status
- SSH into workspace

---

### Skill 22: manage-docker-config

**Purpose**: Manage Docker and docker-compose configurations

**Usage**:
```bash
/manage-docker-config
/manage-docker-config --sync
/manage-docker-config --build
```

**Features**:
- Sync Dockerfile base images from SOT
- Sync docker-compose.yml service versions
- Build images locally
- Push to registry
- Validate configs

---

### Skill 23: manage-ci-workflows

**Purpose**: Manage GitHub Actions workflows

**Usage**:
```bash
/manage-ci-workflows
/manage-ci-workflows --sync
/manage-ci-workflows --validate
/manage-ci-workflows --run ci
```

**Features**:
- List all workflows
- Sync tool versions from SOT
- Validate workflow syntax
- Trigger workflow run
- Show workflow status
- Download workflow logs

---

## Category 5: Monitoring & Observability (2 skills needed)

### Skill 24: check-health

**Purpose**: Check health of all systems (extends check-automation)

**Usage**:
```bash
/check-health
/check-health --verbose
```

**Features**:
- Check automation system (existing check-automation)
- Check backend services (database, cache, search)
- Check frontend build
- Check external integrations (TMDb, etc.)
- Show resource usage (CPU, memory, disk)
- Check for errors in logs

---

### Skill 25: view-logs

**Purpose**: View and search logs from automation runs

**Usage**:
```bash
/view-logs
/view-logs --workflow doc-generation
/view-logs --failed
/view-logs --search "error"
```

**Features**:
- List recent workflow runs
- Show logs for specific run
- Search logs
- Filter by success/failure
- Download logs locally

---

## Summary by Category

### Documentation Automation (6 skills) ✅
1. scaffold-doc ✅
2. generate-docs ✅
3. validate-doc ✅
4. migrate-doc ✅
5. sync-configs ✅
6. check-automation ✅

### GitHub Project Management (7 skills) ❌
7. setup-github-projects
8. setup-github-discussions
9. configure-branch-protection
10. setup-codeql
11. manage-labels
12. assign-reviewers
13. manage-milestones

### Dependency & Release (3 skills) ❌
14. configure-dependabot
15. configure-release-please
16. update-dependencies

### Code Quality (4 skills) ❌
17. run-linters
18. run-tests
19. format-code
20. check-licenses

### Infrastructure (3 skills) ❌
21. manage-coder-workspace
22. manage-docker-config
23. manage-ci-workflows

### Monitoring (2 skills) ❌
24. check-health
25. view-logs

**Total**: 25 skills (6 implemented, 19 missing)

---

## Priority Ranking

### P0: Critical (implement before Day 1)
- check-automation (already specified)
- run-linters
- run-tests
- check-health

### P1: High (implement Phase 1-2)
- setup-github-projects
- setup-github-discussions
- configure-branch-protection
- setup-codeql
- configure-dependabot
- configure-release-please
- manage-labels

### P2: Medium (implement Phase 3-6)
- assign-reviewers
- manage-milestones
- update-dependencies
- format-code
- check-licenses
- manage-coder-workspace
- manage-docker-config

### P3: Low (implement Phase 7-8 or post-MVP)
- manage-ci-workflows
- view-logs

---

## Implementation Strategy

### Option A: Comprehensive (recommended)
Implement all 25 skills in parallel with automation system.

**Timeline**: +10 days (total 35-40 days)
**Effort**: 2-3 days per category
**Benefit**: Complete system coverage from Day 1

### Option B: Phased
Implement in phases:
- Phase 1: Doc automation (6 skills) - Days 1-7
- Phase 2: GitHub management (7 skills) - Days 8-14
- Phase 3: Code quality (4 skills) - Days 15-18
- Phase 4: Infrastructure + monitoring (5 skills) - Days 19-23
- Phase 5: Dependency/release (3 skills) - Days 24-27

**Timeline**: 27 days
**Benefit**: Gradual rollout, test as you go

### Option C: Minimal Viable Skills (MVS)
Implement P0+P1 only (13 skills), defer P2/P3.

**Timeline**: +5 days (total 30 days)
**Benefit**: Faster to MVP, sufficient for daily use

---

## Recommendation

**Go with Option A: Comprehensive**

Reasons:
1. User selected "EVERYTHING" for automation - be consistent
2. Many skills are simple (1-2 hours each)
3. Implementing together ensures integration works
4. Skills are force multipliers - worth the upfront investment
5. Having all skills from Day 1 prevents workflow gaps

**Revised Timeline**: 35-40 days (was 25 days)

**Breakdown**:
- Days 1-7: Foundation + Doc automation (6 skills)
- Days 8-14: GitHub management (7 skills)
- Days 15-21: Code quality + Infrastructure (7 skills)
- Days 22-27: Dependency/Release + Monitoring (5 skills)
- Days 28-35: Testing, refinement, documentation
- Days 36-40: Buffer for issues

---

## Next Steps

1. **User decision**: Which option (A, B, or C)?
2. **Create detailed specs**: Extend 17_CLAUDE_SKILLS_SPECIFICATION.md with skills 7-25
3. **Update implementation plan**: Revise 16_IMPLEMENTATION_PLAN.md with skill phases
4. **Implement skills**: Build Python scripts in `.claude/skills/`

---

**Status**: Analysis complete
**Missing skills identified**: 19 (across 5 categories)
**Recommendation**: Option A (Comprehensive, 35-40 days)


# Final Implementation Plan - Documentation Automation System

**Created**: 2026-01-31
**Status**: ✅ FINALIZED - Ready for implementation
**Scope**: COMPREHENSIVE - Everything automated, all 25 skills, full GitHub integration

---

## Finalized Decisions

1. ✅ **Content Writing**: Build automation + templates NOW. Write doc content LATER (before coding phase)
2. ✅ **Wiki Generation**: INCLUDED - Dual output (Claude + Wiki) from single template
3. ✅ **GitHub Integration**: ALL FEATURES - Projects, Discussions, branch protection, CodeQL, labels, reviewers, milestones
4. ✅ **Skills**: ALL 25 SKILLS - 6 doc automation + 19 management
5. ✅ **Timeline**: No estimates - Build it right, ship when ready

---

## System Overview

### Scope
**Everything Selected, Everything Automated**:

**Documentation Generation**:
- ✅ Design docs (Claude-optimized in `docs/dev/design/`)
- ✅ Wiki docs (User-friendly in `docs/wiki/`)
- ✅ User docs (`docs/user/`)
- ✅ API reference docs (`docs/api/`)
- ✅ Project files (README, CONTRIBUTING, etc.)
- ✅ GitHub templates (issue, PR templates)

**Configuration Synchronization**:
- ✅ IDE settings (VS Code, Zed, JetBrains)
- ✅ Coder workspace templates
- ✅ Language version files (.tool-versions, .nvmrc, .python-version, go.mod)
- ✅ CI/CD workflows (GitHub Actions)
- ✅ Linter configs (golangci-lint, ruff, markdownlint)
- ✅ Docker configs (Dockerfile, docker-compose.yml)

**GitHub Integration** (10+ features):
- ✅ GitHub Projects (boards, automation rules)
- ✅ GitHub Discussions (categories, templates)
- ✅ Branch protection rules (develop, main)
- ✅ CodeQL security scanning (Go + JavaScript)
- ✅ Repository settings sync
- ✅ Label management
- ✅ Auto-assign reviewers (CODEOWNERS)
- ✅ Milestone automation
- ✅ Dependabot configuration
- ✅ Release Please (auto-versioning, changelog)

**Code Quality Automation**:
- ✅ Auto-format on commit (gofmt, prettier, ruff)
- ✅ PR lint reports
- ✅ Coverage tracking (80% threshold)
- ✅ License compliance checking
- ✅ Secret scanning (gitleaks)

**Issue/PR Management**:
- ✅ Auto-labeling (type, size, area)
- ✅ Auto-assign reviewers
- ✅ Auto-close on merge
- ✅ Stale bot

**Monitoring & Observability**:
- ✅ Automation failure alerts (GitHub issues)
- ✅ Health checks
- ✅ Log viewing

---

## Phase Breakdown

### Phase 1: Foundation (Core Infrastructure)
**Goal**: Set up project structure, dependencies, bot account, fetch sources

**Deliverables**:
- [ ] Project structure created (templates/, schemas/, data/, scripts/automation/)
- [ ] Dependencies installed (Python packages, npm tools, gitleaks)
- [ ] Bot user account created (`revenge-bot`)
- [ ] 17 new sources fetched (GitHub docs, style guides, API standards)
- [ ] SOT parser built and tested
- [ ] Shared data extraction working (data/shared-sot.yaml)

**Key Files**:
- `scripts/automation/sot_parser.py`
- `scripts/requirements.txt` (updated)
- `.github/` (automation-config.yml)

---

### Phase 2: Template System
**Goal**: Create Jinja2 templates for all doc types with Claude/Wiki dual output

**Deliverables**:
- [ ] Base template with blocks (`templates/base.md.jinja2`)
- [ ] Feature template (`templates/feature.md.jinja2`)
- [ ] Service template (`templates/service.md.jinja2`)
- [ ] Integration template (`templates/integration.md.jinja2`)
- [ ] **Wiki templates** (`templates/wiki/feature.md.jinja2`, etc.)
- [ ] User docs template (`templates/user.md.jinja2`)
- [ ] API docs template (`templates/api.md.jinja2`)
- [ ] Project files templates (README, CONTRIBUTING, etc.)
- [ ] JSON schemas for all types (feature.schema.json, etc.)
- [ ] Template testing framework
- [ ] **Pilot migration** (MOVIE_MODULE, MUSIC_MODULE, TMDB)

**Key Files**:
- `templates/*.jinja2`
- `schemas/*.schema.json`
- `tests/test_templates.py`

---

### Phase 3: Data Extraction & Migration
**Goal**: Convert all 136+ existing docs to template-based system

**Deliverables**:
- [ ] Markdown parser built (`scripts/automation/markdown_parser.py`)
- [ ] Extract data from existing docs → YAML files
- [ ] Validate extracted data (YAML schema)
- [ ] Multi-stage migration:
  - [ ] Pilot (3 docs): Validate template + parser
  - [ ] 10% (13 docs): Validate at scale
  - [ ] 50% (68 docs): Validate performance
  - [ ] 100% (136 docs): Full migration
- [ ] All docs migrated to `data/{category}/{DOC_NAME}.yaml`
- [ ] Category-level shared data files created

**Key Files**:
- `scripts/automation/markdown_parser.py`
- `data/features/*.yaml`
- `data/services/*.yaml`
- `data/integrations/*.yaml`
- `data/shared-features.yaml`
- `data/shared-services.yaml`
- `data/shared-integrations.yaml`

---

### Phase 4: Validation Pipeline
**Goal**: Build comprehensive validation for all generated content

**Deliverables**:
- [ ] YAML schema validation (yamale + JSON schemas)
- [ ] Markdown linting (markdownlint-cli)
- [ ] Link validation (markdown-link-check)
- [ ] SOT reference validator (versions match)
- [ ] Secret scanning (gitleaks)
- [ ] Frontmatter validation
- [ ] Full validation pipeline script
- [ ] Validation report generation

**Key Files**:
- `scripts/automation/validator.py`
- `.markdownlint.json`
- `.markdown-link-check.json`
- `.gitleaksignore`

---

### Phase 5: Generation Pipeline
**Goal**: Build doc generation with atomic operations, loop prevention, PR automation

**Deliverables**:
- [ ] Generation script (`scripts/automation/doc_generator.py`)
- [ ] Atomic operations (temp → validate → swap)
- [ ] Loop prevention:
  - [ ] Bot user check (skip if author = revenge-bot)
  - [ ] Cooldown lock (.automation-lock with 1hr timeout)
  - [ ] No automatic SOT update (human review gate)
- [ ] PR creation automation (batched by trigger type)
- [ ] Auto-merge for docs-only PRs
- [ ] TOC generation (markdown-toc post-processing)
- [ ] **Wiki generation** (dual output from same data)
- [ ] End-to-end generation tested

**Key Files**:
- `scripts/automation/doc_generator.py`
- `scripts/automation/pr_creator.py`
- `.github/workflows/doc-generation.yml`

**Wiki Generation**:
```
Single data source → Two outputs:
  1. docs/dev/design/{category}/{DOC}.md (Claude-optimized)
  2. docs/wiki/{category}/{DOC}.md (User-friendly)

Template conditionals:
  {% if claude %}...technical details...{% endif %}
  {% if wiki %}...user guide, screenshots...{% endif %}
```

---

### Phase 6: Config Synchronization
**Goal**: Auto-sync all config files from SOURCE_OF_TRUTH.md

**Deliverables**:
- [ ] Config sync script (`scripts/automation/config_sync.py`)
- [ ] IDE settings sync (VS Code, Zed, JetBrains)
- [ ] Language version files sync (.tool-versions, .nvmrc, .python-version, go.mod)
- [ ] CI/CD workflows sync (GitHub Actions - tool versions)
- [ ] Linter configs sync (golangci-lint, ruff, markdownlint)
- [ ] Docker configs sync (Dockerfile, docker-compose.yml)
- [ ] Coder template sync (.coder/template.tf)
- [ ] Validation for all configs
- [ ] Development Tools table in SOT (tool → version → config sync paths)

**Key Files**:
- `scripts/automation/config_sync.py`
- `docs/dev/design/00_SOURCE_OF_TRUTH.md` (updated with Development Tools table)

---

### Phase 7: GitHub Projects & Discussions
**Goal**: Set up GitHub project management features

**Deliverables**:
- [ ] GitHub Projects configured:
  - [ ] Project board created (Backlog, Todo, In Progress, Review, Done)
  - [ ] Automation rules (auto-add issues, auto-move on PR)
  - [ ] Custom fields (Priority, Effort, Module)
  - [ ] Project views (Board, Table, Roadmap)
- [ ] GitHub Discussions configured:
  - [ ] Discussions enabled
  - [ ] Categories created (Ideas, Q&A, Announcements, Bugs)
  - [ ] Discussion templates
  - [ ] Auto-convert rules (discussion → issue)
- [ ] Integration scripts (add to project, create discussion)

**Key Files**:
- `scripts/automation/github_projects.py`
- `scripts/automation/github_discussions.py`
- `.github/DISCUSSION_TEMPLATE/*.md`

---

### Phase 8: GitHub Security & Branch Protection
**Goal**: Set up GitHub Advanced Security and branch protection

**Deliverables**:
- [ ] Branch protection rules configured:
  - [ ] develop: Require PR reviews, status checks, linear history
  - [ ] main: Require PR reviews, status checks, linear history
  - [ ] Include administrators
  - [ ] No force push
- [ ] CodeQL configured:
  - [ ] GitHub Advanced Security enabled
  - [ ] CodeQL analysis for Go + JavaScript
  - [ ] Automated scanning on push/PR
  - [ ] Security alerts configured
  - [ ] Dependency review enabled
- [ ] Secret scanning enabled
- [ ] Configuration scripts

**Key Files**:
- `scripts/automation/github_security.py`
- `.github/workflows/codeql.yml`

---

### Phase 9: GitHub Automation (Labels, Reviewers, Milestones)
**Goal**: Automate issue/PR management workflows

**Deliverables**:
- [ ] Label management:
  - [ ] Label config file (`.github/labels.yml`)
  - [ ] Auto-sync labels from config
  - [ ] Auto-label PRs (by type, size, area)
- [ ] Auto-assign reviewers:
  - [ ] CODEOWNERS file configured
  - [ ] Auto-assignment on PR creation
  - [ ] Team assignments
- [ ] Milestone automation:
  - [ ] Auto-create milestones
  - [ ] Auto-assign issues to milestones
  - [ ] Auto-close milestones when done
  - [ ] Move open issues to next milestone
- [ ] Stale bot configured

**Key Files**:
- `.github/labels.yml`
- `CODEOWNERS`
- `scripts/automation/github_labels.py`
- `scripts/automation/github_milestones.py`
- `.github/workflows/stale.yml`

---

### Phase 10: Dependency & Release Automation
**Goal**: Set up Dependabot and Release Please

**Deliverables**:
- [ ] Dependabot configured:
  - [ ] `.github/dependabot.yml` generated from SOT
  - [ ] Package ecosystems: Go, npm, Python, GitHub Actions
  - [ ] Update schedule (weekly)
  - [ ] Auto-merge rules for patch updates
  - [ ] Reviewers assigned
- [ ] Release Please configured:
  - [ ] `.github/release-please-config.json` created
  - [ ] Release types: Go (backend), Node (frontend)
  - [ ] Changelog generation
  - [ ] Monorepo handling
  - [ ] GitHub Actions workflow
- [ ] Manual dependency update script

**Key Files**:
- `.github/dependabot.yml`
- `.github/release-please-config.json`
- `.github/workflows/release-please.yml`
- `scripts/automation/update_dependencies.py`

---

### Phase 11: Code Quality Automation
**Goal**: Set up linting, testing, formatting, license checking

**Deliverables**:
- [ ] Linter runner script:
  - [ ] Run golangci-lint (Go)
  - [ ] Run ruff (Python)
  - [ ] Run markdownlint (Docs)
  - [ ] Run prettier (TypeScript/JSON/YAML)
  - [ ] Parallel execution
  - [ ] Auto-fix option
- [ ] Test runner script:
  - [ ] Run Go tests (unit + integration)
  - [ ] Run Python tests
  - [ ] Run frontend tests (Vitest)
  - [ ] Coverage reporting (80% threshold)
  - [ ] Watch mode
- [ ] Code formatter script
- [ ] License checker:
  - [ ] Scan Go + npm dependencies
  - [ ] Check against allow/deny lists (from SOT)
  - [ ] Generate license report

**Key Files**:
- `scripts/automation/run_linters.py`
- `scripts/automation/run_tests.py`
- `scripts/automation/format_code.py`
- `scripts/automation/check_licenses.py`

---

### Phase 12: Infrastructure Management
**Goal**: Manage Coder, Docker, CI/CD configs

**Deliverables**:
- [ ] Coder workspace management:
  - [ ] Update template from SOT
  - [ ] Create/start/stop workspace scripts
  - [ ] SSH helper
- [ ] Docker config management:
  - [ ] Sync Dockerfile base images from SOT
  - [ ] Sync docker-compose.yml versions
  - [ ] Build/push scripts
- [ ] CI/CD workflow management:
  - [ ] List/validate workflows
  - [ ] Trigger workflow runs
  - [ ] Download logs

**Key Files**:
- `scripts/automation/manage_coder.py`
- `scripts/automation/manage_docker.py`
- `scripts/automation/manage_ci.py`

---

### Phase 13: Monitoring & Health Checks
**Goal**: Set up monitoring, health checks, log viewing

**Deliverables**:
- [ ] System health checker:
  - [ ] Check automation system (dependencies, templates, schemas)
  - [ ] Check backend services (database, cache, search)
  - [ ] Check frontend build
  - [ ] Check external integrations
  - [ ] Resource usage monitoring
- [ ] Log viewer:
  - [ ] View workflow runs
  - [ ] Search logs
  - [ ] Filter by success/failure
  - [ ] Download logs
- [ ] Failure alerting (GitHub issues on failure)

**Key Files**:
- `scripts/automation/check_health.py`
- `scripts/automation/view_logs.py`

---

### Phase 14: Claude Code Skills
**Goal**: Implement all 25 skills for user interaction

**Skills to Implement**:

**Documentation (6 skills)**:
1. scaffold-doc
2. generate-docs
3. validate-doc
4. migrate-doc
5. sync-configs
6. check-automation

**GitHub Management (7 skills)**:
7. setup-github-projects
8. setup-github-discussions
9. configure-branch-protection
10. setup-codeql
11. manage-labels
12. assign-reviewers
13. manage-milestones

**Dependency & Release (3 skills)**:
14. configure-dependabot
15. configure-release-please
16. update-dependencies

**Code Quality (4 skills)**:
17. run-linters
18. run-tests
19. format-code
20. check-licenses

**Infrastructure (3 skills)**:
21. manage-coder-workspace
22. manage-docker-config
23. manage-ci-workflows

**Monitoring (2 skills)**:
24. check-health
25. view-logs

**Deliverables**:
- [ ] All 25 skill Python scripts in `.claude/skills/`
- [ ] Common utilities (`common/sot_parser.py`, `common/template_renderer.py`, etc.)
- [ ] Skill tests
- [ ] Skill documentation in `.claude/CLAUDE.md`

**Key Files**:
- `.claude/skills/*.py` (25 files)
- `.claude/skills/common/*.py`
- `tests/skills/*.py`

---

### Phase 15: Testing & Validation
**Goal**: Comprehensive testing of entire system

**Deliverables**:
- [ ] Unit tests for all scripts (80%+ coverage)
- [ ] Integration tests for pipelines
- [ ] E2E tests for full automation workflow
- [ ] Load tests (generate all 136 docs, check performance)
- [ ] Failure scenario tests (partial failures, network errors, race conditions)
- [ ] Rollback tests (atomic operations, git revert)
- [ ] Security tests (secret scanning, YAML injection, Jinja2 sandboxing)
- [ ] All tests passing

**Test Scenarios**:
- Generate single doc
- Generate category (20 docs)
- Generate all docs (136 docs)
- Concurrent generation (multiple triggers)
- Validation failures
- Network failures (external sources)
- Template errors
- Data extraction errors
- Config sync failures
- GitHub API failures

**Key Files**:
- `tests/test_sot_parser.py`
- `tests/test_doc_generator.py`
- `tests/test_validator.py`
- `tests/test_config_sync.py`
- `tests/integration/test_full_pipeline.py`

---

### Phase 16: Documentation & Refinement
**Goal**: Document the automation system itself

**Deliverables**:
- [ ] Architecture documentation (`.claude/docs/automation/`)
- [ ] Troubleshooting guide
- [ ] Maintenance playbook
- [ ] Onboarding documentation
- [ ] Update `.claude/CLAUDE.md` with automation instructions
- [ ] Update `CONTRIBUTING.md` with doc contribution guidelines
- [ ] Create runbooks for common operations
- [ ] Performance optimization based on testing results

**Documentation Topics**:
- How the automation works (data flow, triggers, pipelines)
- How to add a new doc type
- How to update templates
- How to fix broken automation
- Common errors and solutions
- Performance tuning
- Security considerations

**Key Files**:
- `.claude/docs/automation/ARCHITECTURE.md`
- `.claude/docs/automation/TROUBLESHOOTING.md`
- `.claude/docs/automation/MAINTENANCE.md`
- `.claude/docs/automation/ONBOARDING.md`

---

## Success Criteria

### Functional Requirements
✅ All 136+ docs migrated to template system
✅ SOT → docs generation working (Claude + Wiki)
✅ Validation pipeline passing (YAML, lint, links, SOT refs, secrets)
✅ Loop prevention working (no infinite loops)
✅ PR automation working (batched by trigger type)
✅ Auto-merge for docs-only PRs working
✅ Config sync working (IDE, CI/CD, language files, Docker, Coder)
✅ All 10+ GitHub features configured
✅ All 25 skills implemented and tested
✅ Dependabot configured and running
✅ Release Please configured and running
✅ CodeQL scanning active
✅ Branch protection enforced
✅ Failure alerting working (GitHub issues created)
✅ 17 new sources fetched and documented
✅ Zero regressions (all existing docs still render correctly)

### Quality Requirements
✅ 80%+ test coverage for all automation scripts
✅ All validation checks passing
✅ No security vulnerabilities (gitleaks, CodeQL clean)
✅ Performance: <30s to generate all 136 docs
✅ Performance: <5s to generate single doc
✅ All configs synced from SOT (no drift)
✅ All GitHub features operational

### Documentation Requirements
✅ All automation documented
✅ All skills documented
✅ Troubleshooting guide complete
✅ Onboarding guide complete
✅ `.claude/CLAUDE.md` updated

---

## Key Technologies

### Core Stack
- **Python 3.12+**: Automation scripts, skills
- **Jinja2 3.1.5+**: Template engine (SandboxedEnvironment)
- **PyYAML 6.0+**: YAML parsing (safe_load only)
- **yamale 4.0+**: YAML schema validation
- **Node 20.x**: npm tools (markdown-link-check, markdown-toc)
- **gitleaks 8.18+**: Secret scanning
- **GitHub API**: Projects, Discussions, branch protection, CodeQL

### Development Tools
- **Go 1.25+**: Backend (future)
- **golangci-lint v1.61.0**: Go linting
- **ruff 0.4+**: Python linting/formatting
- **markdownlint**: Doc linting
- **pytest 8.0+**: Python testing

### Infrastructure
- **GitHub Actions**: CI/CD
- **Coder v2.17.2+**: Dev environments
- **Docker 27+**: Containers

---

## Implementation Approach

### Build Order
**Sequential phases, parallel within phases**:
- Phases 1-5: **Core automation system** (foundation → generation pipeline)
- Phases 6-10: **GitHub integration** (all features in parallel)
- Phases 11-13: **Code quality & infrastructure** (can parallelize)
- Phase 14: **Skills** (can build alongside other phases)
- Phases 15-16: **Testing & documentation** (final polish)

### Development Workflow
1. **Build feature**
2. **Write tests** (80%+ coverage)
3. **Validate** (run linters, tests)
4. **Document** (code comments, user docs, troubleshooting)
5. **Integration test** (with full system)
6. **Commit** (conventional commits)

### Testing Strategy
- **Unit tests**: Every function, 80%+ coverage
- **Integration tests**: Full pipelines
- **E2E tests**: Real-world scenarios
- **Failure tests**: Edge cases, errors
- **Load tests**: Performance validation

### Quality Gates
- All tests pass
- 80%+ coverage
- No lint errors
- No security vulnerabilities
- All validation checks pass
- Documentation complete

---

## Risk Mitigation

### Loop Prevention
- Bot user authorship check
- Cooldown lock (1hr timeout)
- No automatic SOT updates
- Commit message markers

### Data Safety
- Atomic operations (temp → validate → swap)
- Backups before migration
- Git version control
- Rollback tested

### Security
- yaml.safe_load() only (no RCE)
- SandboxedEnvironment for Jinja2
- gitleaks secret scanning
- CodeQL security scanning
- Input validation everywhere

### Performance
- Multiprocessing for bulk generation
- Caching (SOT parsing, shared data)
- Incremental generation (only changed docs)
- Performance tests (load testing)

---

## Deliverables Checklist

### Code
- [ ] `scripts/automation/` - 15+ automation scripts
- [ ] `.claude/skills/` - 25 skill scripts
- [ ] `templates/` - 10+ Jinja2 templates
- [ ] `schemas/` - 5+ JSON schemas
- [ ] `data/` - 136+ YAML data files

### Configuration
- [ ] `.github/automation-config.yml`
- [ ] `.github/dependabot.yml`
- [ ] `.github/release-please-config.json`
- [ ] `.github/labels.yml`
- [ ] `.github/workflows/*.yml` (updated)
- [ ] `.markdownlint.json`
- [ ] `.markdown-link-check.json`
- [ ] `.gitleaksignore`
- [ ] `CODEOWNERS`

### Documentation
- [ ] `.claude/docs/automation/` - 4+ docs
- [ ] `.claude/CLAUDE.md` (updated)
- [ ] `CONTRIBUTING.md` (updated)
- [ ] `docs/dev/design/00_SOURCE_OF_TRUTH.md` (updated with Development Tools table)
- [ ] 17 new sources in `docs/dev/sources/`

### GitHub
- [ ] Projects configured
- [ ] Discussions configured
- [ ] Branch protection active
- [ ] CodeQL active
- [ ] Labels synced
- [ ] Milestones configured

---

## Post-Implementation

### Content Writing Phase (Future)
**After automation complete, before coding**:
- Write missing design docs using scaffold-doc skill
- Complete partial docs (MUSIC_MODULE, etc.)
- Write critical docs (API.md, FRONTEND.md)
- Achieve 99% design completion goal
- Use automation to maintain docs as designs evolve

**Estimated effort**: 170-230 hours (from earlier analysis)
**But**: Automation makes this faster and easier

### Continuous Improvement
- Monitor automation performance
- Add more skills as needs arise
- Optimize generation speed
- Enhance validation checks
- Add more GitHub integrations
- Improve error messages
- Add more test coverage

---

## Status

**Plan Status**: ✅ FINALIZED
**Decisions Made**: 4/4 critical decisions resolved
**Scope**: COMPREHENSIVE - Everything automated, all features, all skills
**Ready**: YES - Begin implementation immediately

**Next Step**: Begin Phase 1 - Foundation

---

**Last Updated**: 2026-01-31
**Supersedes**: 16_IMPLEMENTATION_PLAN.md


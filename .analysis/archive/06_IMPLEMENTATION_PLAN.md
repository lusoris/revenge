# Implementation Plan - Documentation Restructuring & MVP Planning

**Date**: 2026-01-31
**Status**: DRAFT - Awaiting validation
**Estimated Duration**: 15-25 days

---

## Overview

This document provides the detailed implementation plan for:
1. Template-based documentation split (Claude vs Wiki)
2. MVP definition and roadmap (v0.1.x → v0.3.x → v1.0.0)
3. Automation (settings sync, wiki generation, code status verification)
4. New Claude skills (4 skills)

**Approach**: Sequential 7-phase implementation with validation gates

---

## Phase 0: Planning & Approval ✅ COMPLETE

**Duration**: 1 day (2026-01-31)

### Tasks Completed
- [x] Analyze current documentation system
- [x] Create comprehensive questions
- [x] Gather all user decisions
- [x] Create versioning guidance
- [x] Document all answers
- [x] Create summary document
- [x] Add required sources to SOURCES.yaml
- [x] Fetch external documentation

### Outputs
- `.analysis/00_ANALYSIS_REPORT.md`
- `.analysis/01_ORIGINAL_TODO_BACKUP.md`
- `.analysis/02_CRITICAL_QUESTIONS.md`
- `.analysis/03_ANSWERS.md`
- `.analysis/04_VERSIONING_GUIDANCE.md`
- `.analysis/05_SUMMARY.md`
- `.analysis/06_IMPLEMENTATION_PLAN.md` (this file)

---

## Phase 1: MVP Definition & Roadmap

**Duration**: 1-2 days
**Dependencies**: Phase 0 complete

### Objectives
- Define MVP scope document (v0.3.x)
- Create detailed roadmap (v0.1.x → v1.0.0)
- Create milestone TODOs
- Update SOURCE_OF_TRUTH.md

### Tasks

#### 1.1: Create MVP Definition Document
**File**: `docs/dev/design/planning/MVP_DEFINITION.md`

**Content Structure**:
```markdown
# MVP Definition (v0.3.x)

## What is MVP?
- Version: v0.3.x
- Timeline: [Based on development progress]
- Goal: Usable Movies + TV media server

## Core Features
### Content Types
- Movies (complete)
- TV Shows (complete)

### Services
- Auth, User, Session, RBAC
- Metadata & Search (Typesense)
- Library Management
- Playback (Direct play + HLS/DASH)

### Integrations
- TMDB, TheTVDB
- Radarr, Sonarr
- OIDC/OAuth
- Typesense

### Frontend
- Full SvelteKit 2 + Svelte 5 UI
- Tailwind CSS 4
- Browse, search, play

## Explicitly NOT in MVP
- QAR/Adult Content
- LiveTV/DVR
- Photos, Comics, Books
- Advanced playback features
- Transcoding (waiting for Blackbeard)

## Success Criteria
- Can add movies and TV shows
- Metadata fetches correctly
- Browse and search works
- Direct playback works
- HLS/DASH streaming works
- User authentication works
- Can replace Jellyfin for basic use
```

**Validation**:
- [ ] Document follows design template
- [ ] All MVP features clearly defined
- [ ] All excluded features listed
- [ ] Success criteria measurable

---

#### 1.2: Create Implementation Roadmap
**File**: `docs/dev/design/planning/IMPLEMENTATION_ROADMAP.md`

**Content Structure**:
```markdown
# Implementation Roadmap (v0.1.x → v1.0.0)

## Overview
Phased implementation following SemVer pre-1.0 strategy

## Milestone Structure

### v0.1.x - Core Foundation
**Timeline**: [TBD]
**Goal**: Backend infrastructure complete

**Deliverables**:
- PostgreSQL + migrations
- Dragonfly cache
- River job queue
- Auth, User, Session, RBAC services
- Library scanner
- Health checks, logging, metrics
- API framework (ogen)

**Exit Criteria**:
- All backend services running
- API responds to health checks
- Tests pass (80%+ coverage)
- No UI yet

**TODO**: `planning/v0.1.x_TODO.md`

---

### v0.2.x - First Content Type
**Timeline**: [TBD]
**Goal**: One content type end-to-end

**Deliverables**:
- Movie module (OR TV module)
- TMDB integration (OR TheTVDB)
- Radarr integration (OR Sonarr)
- Metadata fetching & storage
- Basic SvelteKit UI
- Direct play only

**Exit Criteria**:
- Can scan movie library
- Metadata fetches automatically
- Can browse via web UI
- Can play movies directly

**TODO**: `planning/v0.2.x_TODO.md`

---

### v0.3.x - MVP Complete ⭐
**Timeline**: [TBD]
**Goal**: Production-ready for basic use

**Deliverables**:
- Movie + TV modules complete
- TMDB + TheTVDB + Radarr + Sonarr
- Full SvelteKit UI
- Direct play + HLS/DASH streaming
- Typesense search
- OIDC authentication
- User management UI

**Exit Criteria**:
- All MVP features working
- 80%+ test coverage
- Documentation complete
- Can replace Jellyfin for basic use

**TODO**: `planning/v0.3.x_TODO.md`

---

### v0.4.x - Third Content Type
**Timeline**: [TBD]
**Goal**: Add music support

**Deliverables**:
- Music module
- MusicBrainz/Spotify/Last.fm integration
- Lidarr integration
- Audio playback

**TODO**: `planning/v0.4.x_TODO.md`

---

### v0.5.x - Transcoding
**Timeline**: [TBD, depends on Blackbeard]
**Goal**: Format conversion support

**Deliverables**:
- Blackbeard integration
- Transcoding profiles
- Quality selection

**TODO**: `planning/v0.5.x_TODO.md`

---

### v0.6-0.9.x - Feature Expansion
**Timeline**: [TBD]
**Goal**: Advanced features + more content types

**Deliverables**:
- Skip intro, trickplay, syncplay
- Photos, Comics, Books modules
- LiveTV/DVR
- QAR (adult content)
- Plugin system
- Performance optimization

**TODO**: Multiple milestone TODOs

---

### v1.0.0 - Stable Release
**Timeline**: [TBD]
**Goal**: Full design spec, production-ready

**Criteria**:
- All modules from design complete
- 80%+ test coverage across all code
- All integrations working
- Full documentation (Claude + Wiki)
- API stable (no more breaking changes)
- Performance benchmarks met
```

**Validation**:
- [ ] All milestones defined
- [ ] Exit criteria clear
- [ ] Dependencies documented
- [ ] Links to TODO files

---

#### 1.3: Create Milestone TODO Files
**Files**:
- `docs/dev/design/planning/v0.1.x_TODO.md`
- `docs/dev/design/planning/v0.2.x_TODO.md`
- `docs/dev/design/planning/v0.3.x_TODO.md`
- `docs/dev/design/planning/v0.4.x_TODO.md`
- `docs/dev/design/planning/v0.5.x_TODO.md`

**Format** (example for v0.1.x):
```markdown
# v0.1.x - Core Foundation TODO

## Infrastructure
- [ ] Set up PostgreSQL 18+ with connection pooling
- [ ] Implement migration system (golang-migrate)
- [ ] Set up Dragonfly cache
- [ ] Set up River job queue

## Core Services
- [ ] Implement Auth service (interface + implementation)
- [ ] Implement User service
- [ ] Implement Session service
- [ ] Implement RBAC service with Casbin

## Library Management
- [ ] Implement file scanner (fsnotify)
- [ ] Implement library database schema
- [ ] Implement library service

## Observability
- [ ] Set up structured logging (slog)
- [ ] Set up metrics (Prometheus)
- [ ] Set up tracing (OpenTelemetry)
- [ ] Implement health checks

## API
- [ ] Generate ogen API from OpenAPI spec
- [ ] Implement API handlers
- [ ] Add authentication middleware
- [ ] Add RBAC middleware

## Testing
- [ ] Unit tests for all services (80%+ coverage)
- [ ] Integration tests with testcontainers
- [ ] API tests

## Documentation
- [ ] Update design docs with implementation details
- [ ] Document API endpoints
```

**Validation**:
- [ ] All TODO files created
- [ ] Tasks are specific and actionable
- [ ] Dependencies between tasks noted
- [ ] Linked from roadmap document

---

#### 1.4: Update SOURCE_OF_TRUTH.md
**File**: `docs/dev/design/00_SOURCE_OF_TRUTH.md`

**Changes to make**:

1. **Documentation Map section** - Add planning category:
```markdown
## Documentation Map

### Planning & Roadmap
- [MVP Definition](planning/MVP_DEFINITION.md) - What constitutes v0.3.x MVP
- [Implementation Roadmap](planning/IMPLEMENTATION_ROADMAP.md) - v0.1.x → v1.0.0
- Milestone TODOs: [v0.1.x](planning/v0.1.x_TODO.md), [v0.2.x](planning/v0.2.x_TODO.md), [v0.3.x](planning/v0.3.x_TODO.md)
```

2. **Automation section** - Document pipelines:
```markdown
## Automation & Pipelines

### Documentation Pipeline (6 stages)
Located: `scripts/doc-pipeline/`

1. Generate indexes (DESIGN_INDEX.md, category INDEXes)
2. Add breadcrumbs (navigation links)
3. Update status tables (multi-dimensional tracking)
4. Validate structure (required sections, formatting)
5. Fix broken links (detect and repair)
6. Generate metadata (TOCs, cross-references)

### Source Pipeline (3 stages)
Located: `scripts/source-pipeline/`

1. Fetch sources (from SOURCES.yaml)
2. Generate source index
3. Add design↔source cross-references

### Settings Sync (new)
Located: `scripts/sync-tool-settings.py`

- Reads SOURCE_OF_TRUTH.md
- Updates all tool configs (.vscode, .zed, .jetbrains, .coder)
- Syncs: language versions, formatters, linters, LSP settings

### Wiki Generation (new)
Located: `scripts/generate-wiki.py`

- Reads design doc templates
- Renders with Jinja2
- Generates Claude + Wiki versions
- Auto-syncs to GitHub Wiki
```

3. **Template System section** - Document templates:
```markdown
## Documentation Templates

### Template Types
1. **Claude Template** (`.claude.template.md`)
   - Includes: file paths, internal APIs, DB schemas, code patterns
   - Variables: `{{ package_path }}`, `{{ service_interface }}`
   - Conditionals: `{{ if claude }}...{{ endif }}`

2. **Wiki Template** (generated from Claude template)
   - Includes: screenshots, tutorials, user guides, external API docs
   - Omits: implementation details, internal references
   - Conditionals: `{{ if wiki }}...{{ endif }}`

### Template Variables
Sourced from SOURCE_OF_TRUTH.md:
- `{{ go_version }}` - Go 1.25.6
- `{{ postgres_version }}` - PostgreSQL 18+
- `{{ python_version }}` - Python 3.12+
- (And all dependencies...)
```

4. **Status Tracking section** - Clarify meanings:
```markdown
## Status Tracking

### Status Dimensions
1. **Design** ✅ - Architecture, API, and schema documented
2. **Sources** ✅ - External docs fetched and linked
3. **Instructions** ✅ - Claude Code instructions written
4. **Code** ✅ - Go implementation exists (auto-verified)
5. **Linting** ✅ - golangci-lint passes
6. **Unit Testing** ✅ - 80%+ coverage achieved
7. **Integration Testing** ✅ - Integration tests pass

### Status Verification
- **Design**: Manual review
- **Code**: Automated (via `code-status` skill)
- **Testing**: Automated (via test runners + coverage tools)
- **Last Verified**: Timestamp added to status tables
```

**Validation**:
- [ ] All sections added to SOT
- [ ] Links resolve correctly
- [ ] No duplicate information
- [ ] Linting passes

---

### Phase 1 Deliverables
- [ ] MVP_DEFINITION.md created and validated
- [ ] IMPLEMENTATION_ROADMAP.md created and validated
- [ ] 5 milestone TODO files created
- [ ] SOURCE_OF_TRUTH.md updated
- [ ] All documents linted (no errors)
- [ ] User reviews and approves documents

### Phase 1 Validation Gates
- All documents follow design template
- All links resolve correctly
- No markdown linting errors
- SOT remains single source of truth
- User explicitly approves before Phase 2

---

## Phase 2: Template Creation

**Duration**: 2-3 days
**Dependencies**: Phase 1 complete, Jinja2 docs fetched

### Objectives
- Create Claude doc template
- Create template rendering system
- Test with sample documents
- Validate output quality

### Tasks

#### 2.1: Study Jinja2 Documentation
**Prerequisites**: Jinja2 docs fetched in Phase 0

**Study**:
- `docs/dev/sources/python_tools/jinja2/main.md`
- `docs/dev/sources/python_tools/jinja2/templates.md`
- `docs/dev/sources/python_tools/jinja2/api.md`

**Understanding Needed**:
- Template syntax (`{{ }}`, `{% %}`, `{# #}`)
- Conditionals (`{% if %}...{% endif %}`)
- Variables and filters
- Template inheritance
- Whitespace control

---

#### 2.2: Create Base Template
**File**: `docs/dev/design/.templates/BASE_TEMPLATE.template.md`

**Content** (example structure):
```markdown
# {{ title }}

<!-- SOURCES: {{ sources|join(', ') }} -->
<!-- DESIGN: {{ related_design|join(', ') }} -->

{{ '<!-- TOC-START -->' if generate_toc }}
{{ toc_content }}
{{ '<!-- TOC-END -->' if generate_toc }}

## Status

| Dimension | Status |
|-----------|--------|
| Design | {{ status.design }} |
| Sources | {{ status.sources }} |
| Instructions | {{ status.instructions }} |
| Code | {{ status.code }} |
| Linting | {{ status.linting }} |
| Unit Testing | {{ status.unit_testing }} |
| Integration Testing | {{ status.integration_testing }} |

---

## Overview

{{ overview }}

{% if claude %}
## Implementation Details

**Package**: `{{ package_path }}`
**Interface**: `{{ interface_name }}`

### Service Structure
{{ service_structure }}

### Database Schema
{{ database_schema }}

### Code Patterns
{{ code_patterns }}
{% endif %}

{% if wiki %}
## Getting Started

{{ getting_started_guide }}

### Screenshots
{{ screenshots }}

### Step-by-Step Tutorial
{{ tutorial_steps }}
{% endif %}

## {{ content_sections }}

---

{% if claude %}
## Related Design Docs
{{ related_design_docs }}

## Sources & Cross-References
{{ cross_references }}
{% endif %}

{% if wiki %}
## External Resources
{{ external_resources }}

## Community
- GitHub: {{ github_url }}
- Discord: {{ discord_url }}
{% endif %}
```

**Validation**:
- [ ] Template syntax valid
- [ ] All variables defined
- [ ] Conditionals work correctly
- [ ] Renders without errors

---

#### 2.3: Create Template Rendering Script
**File**: `scripts/generate-wiki.py`

**Features**:
- Load template from file
- Read variables from SOURCE_OF_TRUTH.md
- Read content from design docs
- Render both Claude and Wiki versions
- Write output files
- Validate output

**Pseudocode**:
```python
import jinja2
import yaml
from pathlib import Path

def load_sot_variables():
    """Extract variables from SOURCE_OF_TRUTH.md"""
    # Parse SOT markdown
    # Extract versions, packages, etc.
    # Return as dict

def render_template(template_path, variables, output_mode='claude'):
    """Render template with Jinja2"""
    env = jinja2.Environment(
        loader=jinja2.FileSystemLoader('docs/dev/design/.templates'),
        trim_blocks=True,
        lstrip_blocks=True
    )

    template = env.get_template(template_path)
    variables['claude'] = (output_mode == 'claude')
    variables['wiki'] = (output_mode == 'wiki')

    return template.render(**variables)

def generate_docs(template_file, output_claude, output_wiki):
    """Generate both versions from template"""
    variables = load_sot_variables()

    # Render Claude version
    claude_content = render_template(template_file, variables, 'claude')
    Path(output_claude).write_text(claude_content)

    # Render Wiki version
    wiki_content = render_template(template_file, variables, 'wiki')
    Path(output_wiki).write_text(wiki_content)

    print(f"Generated: {output_claude}, {output_wiki}")
```

**Validation**:
- [ ] Script runs without errors
- [ ] Both versions generate correctly
- [ ] Variables substitute properly
- [ ] Conditionals filter correctly

---

#### 2.4: Create Sample Templates
**Files**:
- `docs/dev/design/features/video/MOVIE_MODULE.template.md`
- `docs/dev/design/services/AUTH.template.md`

**Process**:
1. Take existing MOVIE_MODULE.md
2. Convert to template format
3. Add variables for dynamic content
4. Add claude/wiki conditionals
5. Test rendering

**Validation**:
- [ ] Sample renders match originals (Claude mode)
- [ ] Wiki version omits implementation details
- [ ] Wiki version includes user-friendly content

---

#### 2.5: Test Template System
**Test Cases**:
1. **Variable Substitution**
   - Verify `{{ go_version }}` → "1.25.6"
   - Verify package paths correct

2. **Conditional Rendering**
   - Claude version has implementation details
   - Wiki version has screenshots/tutorials
   - No overlap

3. **Template Inheritance**
   - Base template properties inherited
   - Overrides work correctly

4. **Error Handling**
   - Missing variables caught
   - Invalid syntax reported
   - Graceful failures

**Validation**:
- [ ] All tests pass
- [ ] No rendering errors
- [ ] Output quality acceptable

---

### Phase 2 Deliverables
- [ ] BASE_TEMPLATE.template.md created
- [ ] generate-wiki.py script working
- [ ] 2+ sample templates converted
- [ ] Both Claude and Wiki versions generate correctly
- [ ] All tests pass
- [ ] Documentation for template system

### Phase 2 Validation Gates
- Templates render without errors
- Output matches quality standards
- Variables substitute correctly
- Conditionals work as expected
- User approves template approach

---

---

## Phase 3: Automation Scripts

**Duration**: 3-4 days
**Dependencies**: Phase 2 complete, Go AST docs fetched, GitHub Wiki docs fetched

### Objectives
- Create wiki generation automation
- Create settings sync script
- Create code status verification
- Test all scripts thoroughly

### Key Tasks

**3.1: Wiki Generation Script** (`scripts/generate-wiki.py`)
- Integrate with template system from Phase 2
- Auto-generate on commit (git hook or GitHub Action)
- Sync to GitHub Wiki via API
- Performance: < 30s for all 179 docs

**3.2: Settings Sync Script** (`scripts/sync-tool-settings.py`)
- Read SOURCE_OF_TRUTH.md
- Extract: Go version, Python version, Node version, formatter/linter configs
- Update: `.vscode/settings.json`, `.zed/settings.json`, `.jetbrains/`, `.coder/`
- Validation: Ensure all configs match

**3.3: Code Status Verification** (`scripts/verify-code-status.py`)
- Use Go AST (go/ast, go/parser, golang.org/x/tools/go/packages)
- Parse Go codebase to detect implemented services
- Compare with design doc status
- Update status tables automatically
- Add "last verified" timestamps

**3.4: Integration with Doc Pipeline**
- Add wiki generation as stage 7 of doc-pipeline
- Add settings validation to CI/CD
- Add code status check to status update script

### Deliverables
- [ ] 3 new automation scripts working
- [ ] Integrated with existing pipelines
- [ ] Full test coverage on scripts
- [ ] Documentation for each script

### Validation Gates
- Scripts run without errors
- Output quality acceptable
- Performance requirements met (< 30s)
- CI/CD integration working

---

## Phase 4: Pipeline Integration

**Duration**: 2-3 days
**Dependencies**: Phase 3 complete

### Objectives
- Extend doc-pipeline for wiki generation
- Add settings sync to CI/CD
- Add MVP filtering to status generation
- End-to-end testing

### Key Tasks

**4.1: Extend Doc Pipeline**
- Add stage 7: Wiki generation
- Modify stage 3: Add MVP filtering to status tables
- Update pipeline runner script

**4.2: CI/CD Workflows**
- Modify `.github/workflows/doc-validation.yml` - add wiki validation
- Create `.github/workflows/settings-sync.yml` - validate settings match SOT
- Modify `.github/workflows/fetch-sources.yml` - include new sources

**4.3: End-to-End Testing**
- Test full pipeline: design change → wiki generation → GitHub Wiki sync
- Test settings sync: SOT change → all tool configs update
- Test code status: code change → status table update

### Deliverables
- [ ] Doc pipeline extended (7 stages)
- [ ] 2 new CI/CD workflows
- [ ] Full end-to-end tests passing

### Validation Gates
- Full pipeline runs without errors
- Wiki syncs to GitHub correctly
- Settings stay in sync
- All tests pass

---

## Phase 5: Documentation Migration

**Duration**: 3-5 days
**Dependencies**: Phase 4 complete

### Objectives
- Convert existing docs to template format
- Generate both Claude and Wiki versions
- Fix all broken links
- Full validation

### Key Tasks

**5.1: Convert Design Docs** (179 files)
Priority order:
1. Feature docs (41) - Start with video/MOVIE_MODULE.md, video/TVSHOW_MODULE.md
2. Service docs (15) - AUTH.md, USER.md, etc.
3. Integration docs (65) - Highest count, do in batches
4. Operations docs (7)
5. Technical docs (6)

**5.2: Generate Output**
- Run generate-wiki.py on all converted templates
- Verify Claude versions match originals
- Verify Wiki versions are user-friendly

**5.3: Link Validation**
- Run link checker on all generated docs
- Fix broken internal links
- Fix broken external links

**5.4: Quality Pass**
- Lint all generated docs
- Check formatting
- Verify consistency

### Deliverables
- [ ] All 179 design docs converted to templates
- [ ] Both Claude and Wiki versions generated
- [ ] Zero broken links
- [ ] All linting clean

### Validation Gates
- All links resolve correctly
- Markdown linting passes
- Claude versions match originals
- Wiki versions are user-friendly
- User spot-checks and approves

---

## Phase 6: Skills & Tooling

**Duration**: 2-3 days
**Dependencies**: Phase 5 complete

### Objectives
- Create 4 new Claude Code skills
- Update .claude/CLAUDE.md
- Test all skills

### Key Tasks

**6.1: Create Skills**

1. **mvp-status** (`.claude/skills/mvp-status/SKILL.md`)
   - Show MVP completion percentage
   - List what's done vs pending
   - Suggest next task from roadmap

2. **generate-wiki** (`.claude/skills/generate-wiki/SKILL.md`)
   - Wrapper for generate-wiki.py
   - Regenerate specific docs or all
   - Deploy to GitHub Wiki

3. **verify-settings** (`.claude/skills/verify-settings/SKILL.md`)
   - Run sync-tool-settings.py in check mode
   - Report mismatches
   - Optionally auto-fix

4. **code-status** (`.claude/skills/code-status/SKILL.md`)
   - Run verify-code-status.py
   - Show code implementation status
   - Compare with design status

**6.2: Update .claude/CLAUDE.md**
- Document new skills
- Update workflow section
- Add MVP/roadmap references

**6.3: Test Skills**
- Test each skill individually
- Test skill integration
- Verify output quality

### Deliverables
- [ ] 4 new skills created and working
- [ ] .claude/CLAUDE.md updated
- [ ] Skills documentation complete

### Validation Gates
- All skills work correctly
- Skills integrate with existing automation
- Documentation clear
- User can successfully use all skills

---

## Phase 7: Validation & Rollout

**Duration**: 1-2 days
**Dependencies**: Phase 6 complete

### Objectives
- Full validation of all changes
- Create git tag backup
- Commit and push
- Archive .analysis/ directory

### Key Tasks

**7.1: Full Validation**
- Run full doc pipeline on all docs
- Run full source pipeline
- Lint all files (markdown, Python, YAML)
- Validate all links
- Check all workflows
- Test all skills
- Verify settings sync

**7.2: Quality Checks**
- Zero markdown linting errors
- Zero Python linting errors (ruff)
- Zero broken links
- All CI/CD workflows pass
- All tests pass

**7.3: User Review**
- User reviews all changes
- User tests wiki generation
- User tests new skills
- User approves rollout

**7.4: Git Operations**
- Create tag: `git tag pre-restructure-2026-01-31`
- Push tag: `git push origin pre-restructure-2026-01-31`
- Commit all changes with detailed message
- Push to develop branch

**7.5: Cleanup**
- Archive .analysis/ directory to `.backup/analysis-2026-01-31/`
- Update TODO.md with post-implementation tasks
- Create GitHub issue for any follow-up work

### Deliverables
- [ ] All validation passed
- [ ] Git tag created
- [ ] All changes committed and pushed
- [ ] .analysis/ archived

### Validation Gates
- Zero errors in all validation
- User explicitly approves
- All changes in git
- Backup tag created

---

## Risk Management

### High-Risk Areas
1. **Link breakage** from doc restructuring
   - Mitigation: Automated link checker, incremental changes
2. **Template rendering errors**
   - Mitigation: Extensive testing on samples first
3. **CI/CD pipeline failures**
   - Mitigation: Test locally first, feature branch testing
4. **Wiki sync failures**
   - Mitigation: Manual sync fallback, API error handling

### Rollback Procedure
If anything breaks:
1. `git reset --hard pre-restructure-2026-01-31`
2. Review .analysis/ backups
3. Identify root cause
4. Fix issue
5. Re-test and retry

---

## Success Metrics

### Quantitative
- ✅ 179 design docs converted to templates
- ✅ Both Claude and Wiki versions generate
- ✅ Zero broken links
- ✅ Zero linting errors
- ✅ All 4 skills working
- ✅ Settings sync 100% accurate
- ✅ Wiki auto-generates in < 30s

### Qualitative
- ✅ Claude docs have all implementation details
- ✅ Wiki docs are beginner-friendly
- ✅ MVP roadmap is clear
- ✅ Automation reduces manual work
- ✅ User satisfaction with system

---

## Timeline Summary

| Phase | Duration | Start | Dependencies |
|-------|----------|-------|--------------|
| Phase 0: Planning | 1 day | Day 1 | None |
| Phase 1: MVP Definition | 1-2 days | Day 2 | Phase 0 |
| Phase 2: Templates | 2-3 days | Day 4 | Phase 1, Jinja2 docs |
| Phase 3: Automation | 3-4 days | Day 7 | Phase 2, AST/Wiki docs |
| Phase 4: Pipelines | 2-3 days | Day 11 | Phase 3 |
| Phase 5: Migration | 3-5 days | Day 14 | Phase 4 |
| Phase 6: Skills | 2-3 days | Day 19 | Phase 5 |
| Phase 7: Validation | 1-2 days | Day 22 | Phase 6 |
| **Total** | **15-25 days** | - | - |

---

## Next Steps

**Immediate** (After approval of this plan):
1. Wait for fetch-sources.py to complete
2. Verify all sources fetched correctly
3. Begin Phase 1: Create MVP_DEFINITION.md
4. Continue sequentially through phases

**Before Each Phase**:
- Review phase objectives
- Check dependencies met
- Prepare validation criteria
- User approval for complex changes

---

**STATUS**: DRAFT - Awaiting user review and approval

**Ready to proceed with Phase 1?**

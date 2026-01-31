# Final Answers Summary - Documentation Automation System

**Created**: 2026-01-31
**Purpose**: Master reference of all decisions for implementation
**Status**: ✅ COMPLETE - Ready for implementation plan

**Questions Answered**: 32 total (17 P0 Critical + 12 P1 High + 3 GitHub Management)

---

## 1. DATA FLOW & SOURCE OF TRUTH

### Data Extraction Strategy
**Decision**: **Auto-extract from SOT parsing**
- Parse SOURCE_OF_TRUTH.md markdown tables/sections
- Auto-generate data/*.yaml files
- Single source, automated
- Build markdown parser for extraction

### Shared Data Organization
**Decision**: **Category-level sharing**
- `shared-features.yaml` - Common data for all feature docs
- `shared-services.yaml` - Common data for all service docs
- `shared-integrations.yaml` - Common data for all integration docs
- Doc-specific YAML: Unique content only
- Pro: Smaller shared files, scoped, manageable

### Bootstrap Procedure
**Decision**: **Hybrid SOT sections**
- Some SOT sections manual (principles, architecture, design philosophy)
- Some SOT sections auto-generated (dependency tables, version lists, tool inventories)
- Best of both: human-readable narrative + automated data tables
- Clear boundaries between manual and generated sections

### SOT is Master
**Confirmed**: "sot is master!!!!!"
- SOURCE_OF_TRUTH.md is PRIMARY source
- Everything else is DERIVED:
  - data/*.yaml files → auto-generated from SOT
  - Design docs → generated from templates + data
  - Config files → synced from SOT
  - Tool versions → stored in SOT, synced to configs

---

## 2. LOOP PREVENTION

### Commit Authorship Detection
**Decision**: **Dedicated bot user account**
- Create GitHub bot user: `revenge-bot`
- All automation commits authored by bot
- Git hooks check `git log --format='%an'`
- If author == "revenge-bot", skip automation triggers
- Standard practice, clear, tamper-resistant

### Regeneration Cooldown
**Decision**: **File-based lock with timeout**
- Create `.automation-lock` file when regeneration starts
- Lock file contains: timestamp, trigger reason, PID
- If lock exists and < 1 hour old, skip trigger
- Delete lock when regeneration completes
- Handles crashes: stale locks > 1hr ignored
- Simple, prevents concurrent regeneration

### Dependabot Loop Prevention
**Decision**: **No automatic SOT update on dependabot merge**
- Flow:
  1. Dependabot merges dependency update to `go.mod`
  2. Automation creates PR to update SOURCE_OF_TRUTH.md (doesn't auto-merge)
  3. Human reviews SOT PR, merges manually
  4. SOT merge triggers metadata sync
  5. Metadata sync updates `go.mod`
  6. Dependabot ignores update (version matches original dependabot PR)
- Human in the loop prevents auto-loop
- Double review: dependabot PR + SOT PR

---

## 3. PR WORKFLOW & BATCHING

### PR Batching Strategy
**Decision**: **Batch by trigger type**
- **Dependabot updates**: Individual PRs (standard dependabot behavior)
- **SOT changes**: Batch all doc regenerations into **one PR**
- **Template changes**: Batch all affected docs into **one PR**
- **Scheduled runs**: Batch everything in that run into **one PR**
- Balanced, reviewable chunks
- Reduces PR noise from 30 PRs/week to ~5-10 PRs/week

### Auto-Approval for Trusted Changes
**Decision**: **Auto-merge for docs-only changes**
- If PR only touches `docs/` directory: auto-approve after CI passes
- If PR touches code, configs, scripts: require human review
- Reduces review burden for low-risk documentation updates
- Speeds up doc updates while maintaining safety for critical files

---

## 4. MIGRATION STRATEGY

### Pilot Testing
**Decision**: **Yes, pilot with 2-3 representative docs**
- **Pilot docs**:
  1. `MOVIE_MODULE.md` (complete feature doc)
  2. `MUSIC_MODULE.md` (partial/scaffold doc)
  3. `TMDB.md` (integration doc)
- **Pilot validates**:
  - Template design works
  - Data extraction parser works
  - YAML schema validation works
  - Link checking works
  - TOC generation works
  - Rendering produces expected output
- **Then**: Big bang migrate remaining 133 docs
- **Timeline**: 1-2 days for pilot, then 1 day for big bang

### Data Extraction Tooling
**Decision**: **Build markdown parser for existing docs**
- **Parser extracts**:
  - YAML frontmatter (new format)
  - HTML comment sources: `<!-- SOURCES: ... -->` → `sources_list`
  - HTML comment design refs: `<!-- DESIGN: ... -->` → `design_references`
  - Status table → 7-dimension status + notes
  - Architecture diagrams → multiline string
  - Database schemas → SQL code blocks
  - Implementation checklists → structured phases list
  - Related documents → links list
- **Output**: `data/{category}/{DOC_NAME}.yaml`
- **One-time migration tool**, then maintain YAML files going forward

### Migration Validation Checkpoints
**Decision**: **Multi-stage validation**
- **After pilot (3 docs)**:
  - Validate template design works
  - Validate data extraction works
  - Validate rendering produces expected output
- **After 10% (13 docs)**:
  - Validate data extraction works at scale
  - Validate performance acceptable (<30s per doc)
- **After 50% (68 docs)**:
  - Validate performance still acceptable
  - Check for memory issues, parser edge cases
- **After 100% (136 docs)**:
  - Full validation suite:
    - YAML schema validation (all docs)
    - Markdown linting (all docs)
    - Link checking (all docs)
    - SOT reference checking (all docs)
    - Secret scanning (all docs)
- Catch issues at each stage, stop if critical failure

---

## 5. VALIDATION PIPELINE

### YAML Schema Validation
**Decision**: **JSON Schema with yamale**
- Create JSON Schema for each doc type:
  - `schemas/feature.schema.json`
  - `schemas/service.schema.json`
  - `schemas/integration.schema.json`
- Use `yamale` library to validate
- Example validation:
  ```python
  import yamale
  schema = yamale.make_schema('schemas/feature.schema.json')
  data = yamale.make_data('data/features/video/MOVIE_MODULE.yaml')
  yamale.validate(schema, data)  # Raises exception if invalid
  ```
- Industry standard, clear error messages, tooling support

### Link Validation
**Decision**: **markdown-link-check with custom config**
- Use `markdown-link-check` npm package
- Check internal links (relative paths within repo)
- Check external links (with retries, 30s timeout)
- Custom config `.markdown-link-check.json`:
  - Ignore known-slow sites (with retry exhausted)
  - Set timeout: 30s
  - Set retries: 3
- Fail validation if broken links found
- Run post-generation as validation step

### SOT Reference Validation
**Decision**: **Parse SOT + docs, compare versions**
- **Extract from SOT**:
  - Parse "Go Dependencies" table → `{package: version}`
  - Parse "Development Tools" table → `{tool: version}`
  - Parse "Infrastructure Components" → `{component: version}`
- **Extract from docs**:
  - Parse generated docs for version references
  - Example: "uses fx v1.23.0" → extract package + version
- **Compare**:
  - Flag mismatches: `fx v1.22.0` in doc but SOT has `v1.23.0`
  - Report all mismatches
- Catches drift, ensures docs stay in sync with SOT

---

## 6. TEMPLATE DESIGN

### Template Inheritance
**Decision**: **Base template with blocks**
- **Base template**: `templates/base.md.jinja2`
  ```jinja2
  {% block frontmatter %}{% endblock %}
  {% block title %}{% endblock %}
  {% block status_table %}{% endblock %}
  {% block toc %}{% endblock %}
  {% block overview %}{% endblock %}
  {% block architecture %}{% endblock %}
  {% block database_schema %}{% endblock %}
  {% block implementation_checklist %}{% endblock %}
  {% block related_documents %}{% endblock %}
  ```
- **Feature template**: `templates/feature.md.jinja2`
  ```jinja2
  {% extends "base.md.jinja2" %}
  {% block title %}# {{ feature_name }}{% endblock %}
  {% block overview %}{{ claude_overview }}{% endblock %}
  ```
- DRY, consistent structure, easy to update
- All docs follow same base structure
- Doc-type-specific blocks override base

### Metadata Format
**Decision**: **YAML frontmatter (modern standard)**
- Replace HTML comments with YAML frontmatter
- Example:
  ```yaml
  ---
  sources:
    - fx
    - pgx
    - sqlc
  design_refs:
    - 01_ARCHITECTURE
    - 02_DESIGN_PRINCIPLES
  category: features/video
  last_updated: 2026-01-31
  ---
  ```
- Modern standard (Jekyll, Hugo compatible)
- Parseable by tools
- Future-proof for static site generation
- Optimized, not preserving legacy for preservation's sake

### TOC Generation
**Decision**: **markdown-toc tool post-generation**
- Generate doc without TOC initially
- Run `markdown-toc` as post-processing step
- Insert TOC at `<!-- toc -->` marker
- Standard tool, works well, widely used
- No need to build custom TOC generator

### Status Table Structure
**Decision**: **7-dimension exact match**
- Match existing pattern exactly:
  ```markdown
  | Dimension | Status | Notes |
  |-----------|--------|-------|
  | Design | {{ status_design }} | {{ status_design_notes }} |
  | Sources | {{ status_sources }} | {{ status_sources_notes }} |
  | Instructions | {{ status_instructions }} | {{ status_instructions_notes }} |
  | Code | {{ status_code }} | {{ status_code_notes }} |
  | Linting | {{ status_linting }} | {{ status_linting_notes }} |
  | Unit Testing | {{ status_unit_testing }} | {{ status_unit_testing_notes }} |
  | Integration Testing | {{ status_integration_testing }} | {{ status_integration_testing_notes }} |
  ```
- Consistency with existing docs
- Full tracking granularity for 99% completion goal
- Supports "heavy testing" requirement (explicit test dimensions)
- Data file has 14 fields: 7 status + 7 notes

### Implementation Checklist Structure
**Decision**: **Generate from structured data**
- Data YAML:
  ```yaml
  implementation_phases:
    - phase: 1
      name: "Core Infrastructure"
      tasks:
        - "Create package structure"
        - "Define interfaces"
    - phase: 2
      name: "Database"
      tasks:
        - "Create migration"
        - "Create tables"
  ```
- Template:
  ```jinja2
  {% for phase in implementation_phases %}
  ### Phase {{ phase.phase }}: {{ phase.name }}
  {% for task in phase.tasks %}
  - [ ] {{ task }}
  {% endfor %}
  {% endfor %}
  ```
- Structured, reusable, can programmatically track progress
- Matches existing multi-phase pattern (Phase 1-6)

---

## 7. SOT SETTINGS & CONFIGURATION

### Automation Settings in SOT
**Decision**: **Minimal settings in SOT, rest in separate config**
- **In SOURCE_OF_TRUTH.md**:
  ```markdown
  ## Automation

  | Feature | Enabled | Config |
  |---------|---------|--------|
  | Doc Generation | ✅ | `.github/automation-config.yml` |
  | Dependency Management | ✅ | `.github/dependabot.yml` |
  | Release Automation | ✅ | `.github/release-please-config.json` |
  ```
- **Details in** `.github/automation-config.yml`:
  ```yaml
  doc_generation:
    triggers:
      on_sot_change: true
      on_template_change: true
      on_data_change: true
      on_commit: false  # Prevent loops
    batching:
      enabled: true
      timeout_minutes: 60
    validation:
      yaml_schema: true
      markdown_lint: true
      link_checking: true
      sot_reference_check: true
      secret_scanning: true
  ```
- Keeps SOT clean (high-level overview)
- Detailed config in appropriate location (.github/)
- Config file is auto-synced from SOT-defined schema

### Tool Versions in SOT
**Decision**: **New Development Tools table in SOT**
- Add to SOURCE_OF_TRUTH.md:
  ```markdown
  ## Development Tools

  | Tool | Version | Purpose | Status | Config Sync |
  |------|---------|---------|--------|-------------|
  | Go | 1.25+ | Backend | ✅ | .tool-versions, Dockerfile, .github/workflows/*.yml |
  | Node | 20.x | Frontend | ✅ | .nvmrc, package.json, .github/workflows/*.yml |
  | Python | 3.12+ | Scripts | ✅ | .python-version, .github/workflows/*.yml |
  | gopls | latest | Go LSP | ✅ | .vscode/settings.json, .zed/settings.json |
  | golangci-lint | v1.61.0 | Go Linter | ✅ | .golangci.yml, .github/workflows/lint.yml |
  | ruff | 0.4+ | Python Lint | ✅ | ruff.toml, .github/workflows/lint.yml |
  | markdownlint | latest | Docs Lint | ✅ | .markdownlint.json, .github/workflows/lint.yml |
  | Docker | 27+ | Containers | ✅ | Dockerfile, docker-compose.yml |
  | Coder | v2.17.2 | Dev Env | ✅ | .coder/template.tf |
  ```
- Clear inventory of ALL tools used
- Version tracking
- Config sync mapping (which files to update when version changes)
- Status tracking (which tools are integrated)

---

## 8. SECURITY

### YAML Safe Loading
**Decision**: **yaml.safe_load() only**
- Always use `yaml.safe_load(f)` (NEVER `yaml.load()`)
- Prevents remote code execution (RCE)
- Standard practice, industry recommendation
- Cannot use advanced YAML features (Python objects), but we don't need them

### Jinja2 Template Sandboxing
**Decision**: **SandboxedEnvironment + StrictUndefined**
- Use Jinja2's `SandboxedEnvironment`:
  ```python
  from jinja2 import SandboxedEnvironment, StrictUndefined, FileSystemLoader

  env = SandboxedEnvironment(
      loader=FileSystemLoader('templates/'),
      undefined=StrictUndefined,  # Fail on undefined variables
      autoescape=False,  # Markdown, not HTML
      trim_blocks=True,
      lstrip_blocks=True,
  )
  ```
- Prevents code execution in templates
- Catches undefined variables early
- Safe, production-ready

### Secret Scanning
**Decision**: **gitleaks in validation pipeline**
- Run `gitleaks` on all generated files before commit
- Fail validation if secrets found
- Example patterns detected:
  - API keys
  - AWS credentials
  - GitHub tokens
  - Private keys
  - Database passwords
- Catches accidental leaks before they reach git history
- May have false positives (review needed)

---

## 9. MONITORING & OBSERVABILITY

### Automation Failure Alerting
**Decision**: **Create GitHub issue on failure**
- When automation fails:
  1. Create GitHub issue
  2. Add label: `automation-failure`
  3. Title: `Automation Failure: {reason}`
  4. Body includes:
     - Error message
     - Stack trace
     - Logs (last 100 lines)
     - Trigger reason
     - Timestamp
     - Environment (commit SHA, branch)
  5. Assign to: automation maintainers team
- Persistent, traceable, searchable
- Could create spam if systemic failure (implement deduplication)

---

## 10. MISSING SOURCES TO FETCH

### GitHub Documentation (ALL selected)
**Fetch these sources** and add to `docs/dev/sources/SOURCES.yaml`:

1. **GitHub README best practices**
   - URL: https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-readmes
   - Output: `docs/dev/sources/devops/github-readme.md`

2. **GitHub CONTRIBUTING.md guide**
   - URL: https://docs.github.com/en/communities/setting-up-your-project-for-healthy-contributions/setting-guidelines-for-repository-contributors
   - Output: `docs/dev/sources/devops/github-contributing.md`

3. **GitHub issue/PR template guide**
   - URL: https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests
   - Output: `docs/dev/sources/devops/github-templates.md`

4. **GitHub repo metadata guide**
   - URL: https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-repository-languages
   - Output: `docs/dev/sources/devops/github-metadata.md`

5. **"everything related to those topics..."**
   - GitHub collaboration docs
   - GitHub branch protection docs
   - GitHub Actions best practices
   - (comprehensive GitHub docs coverage)

### Documentation Style Guides
**Fetch these sources**:

1. **Google Developer Documentation Style Guide**
   - URL: https://developers.google.com/style
   - Output: `docs/dev/sources/devops/google-style-guide.md`

2. **Write the Docs best practices**
   - URL: https://www.writethedocs.org/guide/
   - Output: `docs/dev/sources/devops/writethedocs.md`

3. **MarkdownGuide**
   - URL: https://www.markdownguide.org/basic-syntax/
   - Output: `docs/dev/sources/devops/markdown-guide.md`

### API Documentation Standards (ALL selected)
**Fetch these sources**:

1. **OpenAPI 3.1 Specification**
   - URL: https://spec.openapis.org/oas/v3.1.0
   - Output: `docs/dev/sources/apis/openapi-spec.md`

2. **Stoplight API Design Guide**
   - URL: https://stoplight.io/api-design-guide
   - Output: `docs/dev/sources/apis/stoplight-guide.md`

3. **API Documentation Best Practices (Readme.com)**
   - URL: https://readme.com/blog/api-documentation-best-practices
   - Output: `docs/dev/sources/apis/readme-api-docs.md`

4. **Swagger/OpenAPI docs**
   - (User selected even though already have ogen docs - add explicit reference)

### GitHub Project Management Documentation (ALL selected)
**Fetch these sources**:

1. **GitHub Issues documentation**
   - URL: https://docs.github.com/en/issues
   - Output: `docs/dev/sources/devops/github-issues.md`

2. **GitHub Projects documentation**
   - URL: https://docs.github.com/en/issues/planning-and-tracking-with-projects
   - Output: `docs/dev/sources/devops/github-projects.md`

3. **GitHub Discussions documentation**
   - URL: https://docs.github.com/en/discussions
   - Output: `docs/dev/sources/devops/github-discussions.md`

4. **GitHub Advanced Security documentation**
   - URL: https://docs.github.com/en/code-security
   - Output: `docs/dev/sources/devops/github-security.md`

**Total new sources**: ~17 sources to fetch and add to SOURCES.yaml

---

## 11. GITHUB PROJECT MANAGEMENT INTEGRATION (ALL SELECTED!)

User selected **EVERYTHING** for comprehensive GitHub integration:

### Automation Features to Implement

1. **Issue Tracking Automation** (already selected earlier)
   - Auto-label by files changed
   - Auto-assign based on CODEOWNERS
   - Auto-close on PR merge
   - Stale issue bot

2. **GitHub Projects Integration** ✅
   - Auto-add issues/PRs to project boards
   - Auto-move cards based on status changes
   - Project templates for new features
   - Automation rules: "When issue labeled 'bug', add to Bugs column"

3. **Milestones Automation** (already selected earlier)
   - Auto-assign issues to milestones based on labels
   - Auto-close milestones when all issues complete
   - Auto-create next milestone

4. **GitHub Discussions Integration** ✅
   - Auto-create discussion for major features
   - Link discussions to related issues
   - Discussion templates (Q&A, Ideas, Announcements)
   - Auto-convert discussions to issues when actionable

5. **Branch Protection Rules** ✅
   - Auto-configure branch protection for `develop` and `main`
   - Require pull request reviews (1+ approvals)
   - Require status checks pass (CI, linting, tests)
   - Require linear history
   - Include administrators in restrictions

6. **Code Scanning (CodeQL)** ✅
   - Auto-enable GitHub Advanced Security
   - Configure CodeQL analysis for Go + JavaScript
   - Security advisories integration
   - Dependency review enforcement

7. **Repository Settings Sync** ✅
   - Auto-configure repo settings from SOT:
     - Repository description
     - Topics/tags
     - Features enabled (Issues, Projects, Wiki, Discussions)
     - Merge button settings (squash, rebase, merge commit)
     - Auto-delete head branches after merge
   - Sync script runs on SOT changes

8. **Release Automation** (already covered by Release Please)
   - Auto-create releases with assets
   - Generate changelog from commits
   - Release notes with contributors

9. **Dependency Graph Monitoring** (already covered by Dependabot)
   - Monitor dependency graph
   - Alerts for vulnerable dependencies
   - Auto-create PRs for security updates

10. **All Core Features Already Covered** ✅
    - Issue tracking ✅
    - Milestones ✅
    - GitHub Actions ✅
    - (User confirmed these are included)

---

## 12. IMPLEMENTATION SUMMARY

### Scope
**Comprehensive professional automation** covering:
- ✅ Documentation generation (Design, User, API, Project files)
- ✅ GitHub templates (Issues, PRs)
- ✅ Config synchronization (IDE, Coder, language files, CI/CD)
- ✅ Dependency management (Dependabot, vulnerability scanning, version matrix)
- ✅ Release automation (Release Please, changelog, milestones)
- ✅ Code quality (auto-format, lint reports, coverage tracking, license compliance)
- ✅ Issue/PR management (auto-label, auto-assign, auto-close, stale bot)
- ✅ Documentation automation (API docs from OpenAPI, CLI help, badges)
- ✅ Project metadata sync (versions across all files, copyright, tool versions)
- ✅ **GitHub Projects integration** (boards, automation rules)
- ✅ **GitHub Discussions integration** (community Q&A, feature discussions)
- ✅ **Branch protection rules** (enforce review, CI checks)
- ✅ **CodeQL security scanning** (Go + JavaScript)
- ✅ **Repository settings sync** (description, topics, features)

### Key Technical Decisions
1. **SOT is Master**: SOURCE_OF_TRUTH.md is primary, everything else derived
2. **Category-level sharing**: shared-features.yaml, shared-services.yaml, shared-integrations.yaml
3. **Hybrid SOT**: Manual narrative + auto-generated tables
4. **Loop prevention**: Bot user + file-based lock + no automatic SOT update
5. **Batch by trigger type**: Minimize PR noise while maintaining clarity
6. **Pilot before big bang**: Validate with 3 docs before migrating 133
7. **Build markdown parser**: Auto-extract from existing docs
8. **Multi-stage validation**: Checkpoints at pilot, 10%, 50%, 100%
9. **JSON Schema + yamale**: Strict YAML validation
10. **markdown-link-check**: Standard link validation
11. **Base template with blocks**: Jinja2 inheritance for DRY
12. **YAML frontmatter**: Modern standard, not legacy HTML comments
13. **markdown-toc post-processing**: Standard TOC generation
14. **7-dimension status table**: Exact match for 99% completion tracking
15. **Structured implementation checklists**: Programmatically trackable
16. **Minimal SOT settings**: High-level in SOT, details in .github/automation-config.yml
17. **Development Tools table**: Track ALL tools with version + config sync mapping
18. **yaml.safe_load() + SandboxedEnvironment**: Security hardened
19. **gitleaks secret scanning**: Prevent accidental leaks
20. **GitHub issue alerts**: Persistent, traceable failure tracking
21. **Auto-merge docs-only PRs**: Reduce review burden for low-risk changes

### Migration Strategy
1. **Pilot** (1-2 days):
   - Migrate MOVIE_MODULE.md, MUSIC_MODULE.md, TMDB.md
   - Validate template + parser + validation pipeline
2. **Multi-stage migration** (1 day):
   - 10% (13 docs) → validate
   - 50% (68 docs) → validate
   - 100% (136 docs) → full validation
3. **Source fetching** (parallel):
   - Fetch 17 new sources (GitHub docs, style guides, API standards)
   - Add to SOURCES.yaml
   - Run fetch-sources.py

### Automation Workflow
```
Trigger (SOT change, template change, dependabot merge)
  ↓
Check cooldown lock (skip if < 1hr)
  ↓
Check commit author (skip if revenge-bot)
  ↓
Create .automation-lock
  ↓
Generate docs (temp directory)
  ↓
Validate (YAML schema, lint, links, SOT refs, secrets)
  ↓
If validation passes: atomic swap (temp → docs/)
  ↓
Create batched PR (by trigger type)
  ↓
Auto-merge if docs-only, else require review
  ↓
Delete .automation-lock
  ↓
If failure: create GitHub issue with logs
```

---

## 13. NEXT STEPS

1. ✅ **All questions answered** (32 total)
2. ⏭️ **Create detailed implementation plan**
   - Phase-by-phase breakdown
   - Scripts to build
   - Templates to create
   - Config files to add
   - Validation pipeline
   - Migration procedure
   - Testing strategy
3. ⏭️ **Begin implementation with heavy testing**
   - Build markdown parser
   - Create templates
   - Build generation pipeline
   - Build validation pipeline
   - Pilot migration
   - Full migration
   - Source fetching
   - GitHub integration setup

---

**Status**: ✅ COMPLETE - Ready for implementation
**Total Decisions**: 32
**Automation Scope**: Comprehensive (everything selected)
**Timeline Estimate**: 16-25 days for full implementation + testing


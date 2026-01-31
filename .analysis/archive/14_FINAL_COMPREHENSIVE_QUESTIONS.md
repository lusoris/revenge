# Final Comprehensive Questions - Doc Automation System

**Created**: 2026-01-31
**Purpose**: Resolve all critical gaps and finalize automation design before implementation
**Based On**:
- Gap Analysis (13_CRITICAL_GAP_ANALYSIS.md)
- Pattern Analysis (Explore agent deep dive)
- P0 answers already collected
- User requirement: "SOT is master!!!!!"

---

## Question Categories

1. **🔴 P0 Critical** - Blockers (must answer before ANY implementation)
2. **🟡 P1 High** - Important (should answer before detailed design)
3. **🟢 P2 Medium** - Nice to have (can refine during implementation)

---

## SECTION 1: Data Flow & Source of Truth

### ✅ RESOLVED: SOT is Master

**User confirmed**: "sot is master!!!!!" - SOURCE_OF_TRUTH.md is primary source, everything else is derived.

### 🔴 Q1.1: Data Extraction from SOT (CRITICAL)

**Context**: We have 136+ existing docs. SOT is master. How do we extract data for templates?

**Question**: How should we extract template data (YAML files) from SOURCE_OF_TRUTH.md?

**Options**:
A) **Auto-extract from SOT parsing** (Recommended)
   - Parse SOT markdown tables/sections
   - Generate data/*.yaml files automatically
   - Example: Parse "Go Dependencies" table → extract versions
   - Pro: Single source, automated
   - Con: Parsing complexity

B) **Manual data file creation**
   - Developers manually create data/*.yaml per doc
   - Reference SOT manually for versions
   - Pro: Simple, explicit
   - Con: Risk of version drift, manual effort

C) **Hybrid: Auto-extract versions + manual content**
   - Auto-extract from SOT: versions, dependencies, config keys
   - Manual: architecture diagrams, business logic, implementation notes
   - Pro: Automated for objective data, manual for subjective
   - Con: Two systems to maintain

D) **Extract from existing docs**
   - Parse existing design docs → YAML
   - One-time migration process
   - Then maintain YAML files going forward
   - Pro: Preserves existing content
   - Con: Parser must handle inconsistencies

---

### 🔴 Q1.2: Shared Data Organization (CRITICAL)

**Context**: Hierarchical data organization selected. Need to define what's shared vs specific.

**Question**: What data should live in shared.yaml vs doc-specific YAML files?

**Options**:
A) **Maximize sharing** (Recommended)
   - shared.yaml: All SOT data (versions, dependencies, config keys, API namespaces)
   - doc-specific: Only unique content (diagrams, implementation notes, business rules)
   - Pro: DRY, single update point
   - Con: Large shared file

B) **Category-level sharing**
   - shared-features.yaml, shared-services.yaml, shared-integrations.yaml
   - Doc-specific: Everything unique to that doc
   - Pro: Smaller shared files, scoped
   - Con: Multiple shared files to sync

C) **Minimal sharing**
   - shared.yaml: Only common metadata (project name, repo URL)
   - Doc-specific: Everything else (including versions)
   - Pro: Simple, isolated
   - Con: Massive duplication, drift risk

D) **SOT-only (no shared.yaml)**
   - All shared data comes directly from SOT parsing
   - No intermediate shared.yaml file
   - Pro: True single source
   - Con: Every doc generation parses SOT

---

### 🔴 Q1.3: Bootstrap Procedure (CRITICAL)

**Context**: Chicken-and-egg problem - need SOT to generate docs, but SOT is a doc.

**Question**: How do we bootstrap the automation system initially?

**Options**:
A) **SOT is hand-maintained, generates everything else** (Recommended)
   - SOURCE_OF_TRUTH.md is the ONLY manually maintained doc
   - Everything else (design docs, user docs, API docs, configs) is generated
   - Pro: Clear boundary, single manual touchpoint
   - Con: SOT must be perfect before generating

B) **Gradual conversion**
   - Start with SOT + config files as manual
   - Generate design docs first
   - Later add user docs, API docs
   - Pro: Incremental, testable
   - Con: Long transition period

C) **Template SOT too**
   - Even SOT is generated from templates
   - Master data lives in raw YAML files
   - Pro: Consistent generation for everything
   - Con: Loses human-readability of SOT as primary doc

D) **Hybrid: SOT sections**
   - Some SOT sections manual (principles, architecture)
   - Some auto-generated (dependency tables, version lists)
   - Pro: Best of both worlds
   - Con: Complexity in SOT structure

---

## SECTION 2: Loop Prevention

### 🔴 Q2.1: Commit Authorship Detection (CRITICAL)

**Context**: Need to prevent automation from triggering itself infinitely.

**Question**: How do we detect if a commit is from automation vs human?

**Options**:
A) **Dedicated bot user account** (Recommended)
   - Create GitHub bot user "revenge-bot"
   - All automation commits authored by bot
   - Git hook checks author, skips triggers if bot
   - Pro: Clear, standard practice
   - Con: Requires bot account setup

B) **Commit message markers**
   - All automation commits include `[skip ci]` or `[bot]`
   - Git hooks parse commit message, skip if marker present
   - Pro: Simple, no extra account needed
   - Con: Can be bypassed, relies on discipline

C) **Git notes metadata**
   - Automation adds git notes: `git notes add -m "automation=true"`
   - Hooks check notes before triggering
   - Pro: Separate from commit message, harder to bypass
   - Con: Git notes less commonly used, tooling needed

D) **Commit signing**
   - Bot commits GPG signed with bot key
   - Human commits signed with human keys
   - Hooks check signature, skip if bot key
   - Pro: Secure, tamper-proof
   - Con: GPG key management overhead

---

### 🔴 Q2.2: Regeneration Cooldown (CRITICAL)

**Context**: Multiple triggers could cause regeneration storm. Need rate limiting.

**Question**: Should we implement regeneration cooldown to prevent loops?

**Options**:
A) **Yes, file-based lock with timeout** (Recommended)
   - Create `.automation-lock` file when regen starts
   - If lock exists and < 1 hour old, skip trigger
   - Delete lock when regen completes
   - Pro: Simple, prevents concurrent regen
   - Con: Stale lock file if process crashes

B) **Yes, GitHub API check**
   - Check if automation PR already open
   - If yes, skip new trigger
   - Pro: No lock files, uses GitHub state
   - Con: API calls, rate limiting

C) **No cooldown, trust commit authorship**
   - Rely on Q2.1 to prevent loops
   - Pro: Simpler
   - Con: Vulnerable to edge cases

D) **Yes, Redis/Dragonfly distributed lock**
   - Use Dragonfly (already in stack) for distributed locking
   - SET NX with TTL for regen lock
   - Pro: Robust, works in distributed environment
   - Con: Requires Dragonfly running for automation

---

### 🔴 Q2.3: Dependabot Loop Prevention (CRITICAL)

**Context**: Dependabot updates go.mod → SOT PR → metadata sync updates go.mod → loop?

**Question**: How do we prevent Dependabot from detecting automation's go.mod changes?

**Options**:
A) **No automatic SOT update on dependabot merge** (From P0 answers: "PR for SOT update")
   - Dependabot merges dependency update
   - Automation creates PR to update SOT (doesn't auto-merge)
   - Human reviews SOT PR, merges manually
   - SOT merge triggers metadata sync
   - Metadata sync updates go.mod
   - **Dependabot ignores update because version matches original dependabot PR**
   - Pro: Human in the loop, no auto-loop
   - Con: Manual SOT PR review required

B) **Commit message triggers dependabot ignore**
   - Metadata sync commits include `[dependabot skip]`
   - Configure .github/dependabot.yml to ignore these commits
   - Pro: Automated, no manual steps
   - Con: Dependabot config may not support this

C) **Separate branch for automation**
   - Automation works on `automation` branch
   - Never auto-merges to `develop`
   - Always creates PR
   - Dependabot only watches `develop`, ignores `automation`
   - Pro: Complete isolation
   - Con: Extra branch management

D) **Timestamp-based deduplication**
   - Track when dependabot last updated package
   - Metadata sync checks: "did I just update this 5 min ago?"
   - If yes, skip
   - Pro: Smart deduplication
   - Con: State tracking complexity

---

## SECTION 3: PR Workflow & Batching

### 🔴 Q3.1: PR Batching Strategy (CRITICAL)

**Context**: User selected "Dependabot-style auto PRs" but review burden could be 30 PRs/week.

**Question**: When should we batch multiple changes into one PR vs create individual PRs?

**Options**:
A) **Batch by trigger type** (Recommended)
   - Dependabot updates: Individual PRs (standard dependabot behavior)
   - SOT changes: Batch all doc regenerations into one PR
   - Template changes: Batch all affected docs into one PR
   - Scheduled runs: Batch everything in that run into one PR
   - Pro: Balanced, reviewable chunks
   - Con: Logic to implement batching

B) **Always batch within time window**
   - Collect all triggers within 1 hour
   - Create single PR with all changes
   - Pro: Minimal PR noise
   - Con: Delayed updates, complex state management

C) **Never batch, always individual**
   - Every trigger = separate PR
   - Pro: Clear attribution, easy rollback
   - Con: PR noise, review burden

D) **Batch by affected docs**
   - If > 5 docs affected: single PR
   - If ≤ 5 docs: individual PRs per doc
   - Pro: Adaptive based on scale
   - Con: Arbitrary threshold, still complex

---

### 🟡 Q3.2: Auto-Approval for Trusted Changes (HIGH)

**Context**: Some changes are low-risk (version bumps in docs). Auto-approve?

**Question**: Should certain automation PRs be auto-approved and merged?

**Options**:
A) **No auto-merge, always require human review** (Recommended for safety)
   - Every PR waits for approval
   - Pro: Human oversight, catch unexpected changes
   - Con: Review burden, slower updates

B) **Auto-merge for docs-only changes**
   - If PR only touches `docs/` directory, auto-approve
   - If touches code/config, require review
   - Pro: Reduces review for low-risk docs
   - Con: Risk of bad doc generation slipping through

C) **Auto-merge with post-merge validation**
   - Auto-merge after CI passes
   - Run additional validation post-merge
   - If validation fails, auto-revert
   - Pro: Fast updates, safety net
   - Con: Reverts create noise, potential for brief broken state

D) **Auto-merge on schedule, manual for immediate**
   - Scheduled regenerations: auto-merge
   - Triggered regenerations (SOT change): require review
   - Pro: Balanced automation vs control
   - Con: Two different workflows

---

## SECTION 4: Migration Strategy

### 🔴 Q4.1: Pilot Scope (CRITICAL)

**Context**: User selected "Big bang migration" but gap analysis flagged high risk.

**Question**: Should we pilot test the system before big bang migration?

**Options**:
A) **Yes, pilot with 2-3 representative docs** (Recommended)
   - Select 1 complete doc (MOVIE_MODULE.md)
   - Select 1 partial doc (MUSIC_MODULE.md)
   - Select 1 integration doc (TMDB.md)
   - Validate: template works, data extraction works, validation works
   - Then big bang the rest
   - Pro: Catch template bugs early, validate approach
   - Con: Delays full migration by 1-2 days

B) **No pilot, trust big bang**
   - Convert all 136+ docs at once
   - Extensive testing in branch before merge
   - Pro: Fastest completion
   - Con: High risk, difficult rollback

C) **Pilot per category**
   - Pilot features category (1 doc)
   - Then migrate all features
   - Pilot services category (1 doc)
   - Then migrate all services
   - Repeat for each category
   - Pro: Incremental validation
   - Con: Long migration period

D) **Pilot with synthetic test doc**
   - Create fake test doc with all patterns
   - Validate against test doc
   - Then big bang real docs
   - Pro: Safe testing environment
   - Con: Test doc may not cover real edge cases

---

### 🔴 Q4.2: Data Extraction Tooling (CRITICAL)

**Context**: 136+ docs need data/*.yaml files. Manual = 170-230 hours. Need automation.

**Question**: How do we extract data from existing docs into YAML files?

**Options**:
A) **Build markdown parser for existing docs** (Recommended)
   - Parse markdown sections (## Status, ## Architecture, ## Database Schema)
   - Extract tables, code blocks, lists
   - Generate YAML from extracted data
   - Pro: Automated, fast, consistent
   - Con: Parser complexity, edge cases in inconsistent docs

B) **Manual extraction with validation**
   - Humans copy-paste from docs to YAML
   - Automated validation ensures no missing fields
   - Pro: Handles inconsistencies, ensures quality
   - Con: 170-230 hours effort, error-prone

C) **LLM-assisted extraction**
   - Feed each doc to Claude
   - Claude extracts structured data → YAML
   - Human reviews and approves
   - Pro: Handles inconsistencies, faster than manual
   - Con: Requires API access, cost, potential hallucinations

D) **Hybrid: Auto-extract structured data + manual for content**
   - Auto-extract: Status tables, package lists (from HTML comments)
   - Manual: Architecture diagrams, business logic, implementation notes
   - Pro: Automated for mechanical parts, manual for creative parts
   - Con: Two processes to coordinate

---

### 🟡 Q4.3: Migration Validation Checkpoints (HIGH)

**Context**: Big bang migration with pilot. Need validation stages.

**Question**: What validation checkpoints should we have during migration?

**Options**:
A) **Multi-stage validation** (Recommended)
   - After pilot (3 docs): Validate template design works
   - After 10% (13 docs): Validate data extraction works at scale
   - After 50% (68 docs): Validate performance acceptable
   - After 100% (136 docs): Full validation suite
   - Pro: Catch issues at each stage
   - Con: More manual checkpoints

B) **Pre-migration and post-migration only**
   - Validate templates before starting
   - Migrate everything
   - Validate all outputs
   - Pro: Simple, two-stage
   - Con: Late discovery of issues

C) **Continuous validation during migration**
   - Every doc migrated is immediately validated
   - Stop on first failure
   - Pro: Immediate feedback
   - Con: Could stop frequently, slow progress

D) **Category-based validation**
   - Migrate all features docs, validate
   - Migrate all services docs, validate
   - Migrate all integration docs, validate
   - Pro: Natural boundaries
   - Con: Could waste effort if template broken

---

## SECTION 5: Validation Pipeline

### 🔴 Q5.1: YAML Schema Definition (CRITICAL)

**Context**: "Heavy testing" includes YAML schema validation. Need to define schema.

**Question**: How should we define and validate YAML schema for data files?

**Options**:
A) **JSON Schema with yamale** (Recommended)
   - Create JSON Schema for each doc type
   - Use yamale library to validate YAML against schema
   - Pro: Industry standard, tooling support, clear errors
   - Con: Need to write schemas for all doc types

B) **Pydantic models**
   - Define Python Pydantic models for data structure
   - Load YAML, validate against models
   - Pro: Type hints, excellent errors, Python-native
   - Con: Requires Python 3.7+, more complex setup

C) **Custom validation functions**
   - Write Python functions to check required fields
   - Pro: Simple, flexible
   - Con: Ad-hoc, harder to maintain, poor error messages

D) **Jinja2 StrictUndefined only**
   - Rely on Jinja2's StrictUndefined to catch missing vars
   - No separate schema validation
   - Pro: Built-in, no extra tooling
   - Con: Only catches issues during rendering, not before

---

### 🔴 Q5.2: Link Validation Strategy (CRITICAL)

**Context**: Generated docs have many cross-references. Need to validate links don't break.

**Question**: How should we validate links in generated documentation?

**Options**:
A) **markdown-link-check with custom config** (Recommended)
   - Use markdown-link-check npm package
   - Check internal links (relative paths)
   - Check external links (with retries, timeout)
   - Fail validation if broken links found
   - Pro: Standard tool, fast, configurable
   - Con: External link checking slow, false positives

B) **Custom Python link checker**
   - Parse markdown, extract all links
   - Check file:// links exist locally
   - Check http(s):// links return 200
   - Pro: Customizable, can skip known-slow sites
   - Con: Need to write parser, maintain checker

C) **GitHub Actions marketplace action**
   - Use existing action like "lychee-action"
   - Pro: Maintained by community, feature-rich
   - Con: Another dependency, less control

D) **Internal links only, skip external**
   - Only validate relative links within repo
   - Trust external URLs are correct
   - Pro: Fast, no false positives from external sites
   - Con: External links could break unnoticed

---

### 🟡 Q5.3: SOT Reference Validation (HIGH)

**Context**: Docs reference versions from SOT. Need to ensure they match.

**Question**: How do we validate that doc references to SOT are correct?

**Options**:
A) **Parse SOT + docs, compare versions** (Recommended)
   - Extract all versions from SOURCE_OF_TRUTH.md
   - Extract all version references from generated docs
   - Compare: flag mismatches
   - Pro: Catches drift, automated
   - Con: Parser complexity

B) **No validation, trust generation**
   - If docs generated from SOT, versions must match
   - Pro: Simple
   - Con: Doesn't catch manual edits or bugs

C) **Checksum-based validation**
   - Generate checksum of SOT dependency section
   - Store in each doc as HTML comment
   - Validate checksum matches on doc read
   - Pro: Fast, cryptographically sound
   - Con: Doesn't show which version is wrong

D) **Linting rule**
   - Custom linter rule: "no hardcoded versions in docs"
   - Force all version references to use variables/links
   - Pro: Prevents problem at source
   - Con: Requires custom linter

---

## SECTION 6: Template Design

### 🔴 Q6.1: Template Inheritance (CRITICAL)

**Context**: Pattern analysis shows docs have common sections. Should templates use inheritance?

**Question**: Should we use Jinja2 template inheritance to avoid duplication?

**Options**:
A) **Yes, base template with blocks** (Recommended)
   ```jinja2
   # base.md.jinja2
   {% block title %}{% endblock %}
   {% block status_table %}...{% endblock %}
   {% block overview %}{% endblock %}
   {% block architecture %}{% endblock %}

   # feature.md.jinja2
   {% extends "base.md.jinja2" %}
   {% block title %}{{ feature_name }}{% endblock %}
   ```
   - Pro: DRY, consistent structure, easy updates
   - Con: More complex to understand, debugging harder

B) **Yes, include for common sections**
   ```jinja2
   # feature.md.jinja2
   {% include "partials/status_table.jinja2" %}
   {% include "partials/toc.jinja2" %}
   ```
   - Pro: Modular, clear what's included
   - Con: Many partial files to manage

C) **No, separate templates with duplication**
   - Each doc type has full template
   - Accept duplication for clarity
   - Pro: Simple, self-contained, easy to understand
   - Con: Changes must be replicated across templates

D) **Hybrid: base + includes for small sections**
   - Base template for overall structure
   - Includes for tiny snippets (breadcrumbs, footer)
   - Pro: Balanced
   - Con: Two systems to learn

---

### 🔴 Q6.2: HTML Comment Metadata (CRITICAL)

**Context**: Existing docs use `<!-- SOURCES: ... -->` and `<!-- DESIGN: ... -->` comments.

**Question**: Should generated docs preserve these HTML comment metadata patterns?

**Options**:
A) **Yes, auto-generate from data** (Recommended)
   - Template includes:
   ```jinja2
   <!-- SOURCES: {{ sources_list }} -->
   <!-- DESIGN: {{ design_references }} -->
   ```
   - Data YAML has: `sources_list: "fx, pgx, sqlc"`
   - Pro: Preserves existing pattern, automation-compatible
   - Con: Need to extract from existing docs

B) **No, replace with YAML frontmatter**
   - Use Jekyll/Hugo style frontmatter:
   ```yaml
   ---
   sources: [fx, pgx, sqlc]
   design_refs: [01_ARCHITECTURE, 02_DESIGN_PRINCIPLES]
   ---
   ```
   - Pro: Standard, parseable, better tooling
   - Con: Changes existing pattern, breaks automation

C) **Both: frontmatter + HTML comments**
   - Frontmatter for tooling
   - HTML comments for human readability
   - Pro: Best of both worlds
   - Con: Duplication, must stay in sync

D) **No metadata in docs, only in data files**
   - Metadata lives only in data/*.yaml
   - Docs are pure output
   - Pro: Clean separation
   - Con: Harder to trace doc to sources

---

### 🟡 Q6.3: TOC Auto-Generation (HIGH)

**Context**: Existing docs use `<!-- TOC-START -->` markers for table of contents.

**Question**: Should we auto-generate TOC in templates or expect it in data?

**Options**:
A) **Auto-generate from headings** (Recommended)
   - Template parses its own headings
   - Generates TOC dynamically
   - Pro: Always correct, no manual maintenance
   - Con: Jinja2 doesn't have built-in TOC generation

B) **Use markdown-toc tool post-generation**
   - Generate doc without TOC
   - Run markdown-toc as post-processing step
   - Pro: Standard tool, works well
   - Con: Extra step, markers needed

C) **Manual TOC in data files**
   - Data YAML includes: `toc: ["Overview", "Architecture", ...]`
   - Template renders from data
   - Pro: Explicit control
   - Con: Manual maintenance, can drift

D) **No TOC in generated docs**
   - Rely on GitHub's automatic TOC rendering
   - Pro: Zero maintenance
   - Con: Less control, not visible in raw markdown

---

## SECTION 7: Missing Documentation Sources

### 🟡 Q7.1: GitHub Documentation Best Practices (HIGH)

**Context**: We have GitHub Wiki docs but missing broader GitHub documentation guides.

**Question**: Which GitHub documentation sources should we fetch and add to SOURCES.yaml?

**Options** (Multi-select):
A) **GitHub README best practices**
   - https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-readmes
   - Covers: What makes a good README, badges, structure

B) **GitHub CONTRIBUTING.md guide**
   - https://docs.github.com/en/communities/setting-up-your-project-for-healthy-contributions/setting-guidelines-for-repository-contributors
   - Covers: Contribution guidelines, code of conduct

C) **GitHub issue/PR template guide**
   - https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests
   - Covers: Issue templates, PR templates, forms

D) **GitHub repo metadata guide**
   - https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/customizing-your-repository/about-repository-languages
   - Covers: .gitattributes, LICENSE, CODEOWNERS

**Recommended**: Fetch A, B, C (all except D)

---

### 🟡 Q7.2: Documentation Style Guides (HIGH)

**Question**: Should we fetch industry documentation style guides?

**Options** (Multi-select):
A) **Google Developer Documentation Style Guide**
   - https://developers.google.com/style
   - Industry standard for technical writing

B) **Microsoft Writing Style Guide**
   - https://learn.microsoft.com/en-us/style-guide/welcome/
   - Another major industry standard

C) **Write the Docs best practices**
   - https://www.writethedocs.org/guide/
   - Community best practices

D) **MarkdownGuide**
   - https://www.markdownguide.org/basic-syntax/
   - Markdown formatting reference

**Recommended**: Fetch A, D (Google + MarkdownGuide, skip B/C to avoid style conflicts)

---

### 🟡 Q7.3: API Documentation Standards (HIGH)

**Question**: Should we fetch OpenAPI and API documentation best practices?

**Options** (Multi-select):
A) **OpenAPI 3.1 Specification**
   - https://spec.openapis.org/oas/v3.1.0
   - We use ogen which generates OpenAPI

B) **Stoplight API Design Guide**
   - https://stoplight.io/api-design-guide
   - Best practices for REST API design

C) **API Documentation Best Practices (Readme.com)**
   - https://readme.com/blog/api-documentation-best-practices
   - User-facing API docs guide

D) **Swagger/OpenAPI docs**
   - Already have via ogen docs

**Recommended**: Fetch A, B (OpenAPI spec + design guide)

---

## SECTION 8: SOT Settings Table

### 🔴 Q8.1: Automation Settings in SOT (CRITICAL)

**Context**: User said "sot needs a settings table to let us setup this automation system!"

**Question**: What should the SOT automation settings section include?

**Options**:
A) **Comprehensive settings table** (Recommended)
   ```yaml
   automation:
     doc_generation:
       enabled: true
       trigger_on_commit: false  # Prevent loops
       trigger_on_sot_change: true
       trigger_on_template_change: true
       batch_prs: true
       batch_timeout_minutes: 60

     dependency_management:
       dependabot_enabled: true
       auto_update_sot: false  # Manual SOT updates
       create_pr_for_sot_update: true

     validation:
       yaml_schema: true
       markdown_lint: true
       link_checking: true
       sot_reference_check: true
       coverage_threshold: 80

     rollback:
       strategy: "atomic"
       backup_before_gen: true

     commit_settings:
       author_name: "Revenge Bot"
       author_email: "bot@revenge.dev"
       gpg_sign: false
       conventional_commits: true

     pr_settings:
       auto_assign_reviewers: true
       auto_label: true
       auto_close_on_merge: true
       stale_days: 60
   ```
   - Pro: Complete configuration, one place
   - Con: Large YAML section in SOT

B) **Minimal settings, rest in separate config**
   - SOT only has: enabled/disabled flags
   - Details in `.github/automation-config.yml`
   - Pro: Keeps SOT clean
   - Con: Two config files

C) **No settings table, hardcoded in scripts**
   - Automation behavior hardcoded
   - Pro: No config overhead
   - Con: Can't customize without code changes

D) **Settings as code (Python config file)**
   - `scripts/automation_config.py` with settings
   - Pro: Validation via Python types
   - Con: Not in SOT as user requested

**Recommended**: A (Comprehensive settings table in SOT)

---

### 🟡 Q8.2: Tool Versions in SOT (HIGH)

**Context**: User said "tools the project uses or wants to use should be set in the sot too"

**Question**: How should we document external tools in SOT?

**Options**:
A) **New "Development Tools" table** (Recommended)
   ```markdown
   ## Development Tools

   | Tool | Version | Purpose | Status | Config Sync |
   |------|---------|---------|--------|-------------|
   | Go | 1.25+ | Backend | ✅ | .tool-versions, Dockerfile, CI |
   | Node | 20.x | Frontend | ✅ | .nvmrc, package.json, CI |
   | Python | 3.12+ | Scripts | ✅ | .python-version, CI |
   | gopls | latest | LSP | ✅ | IDE settings |
   | golangci-lint | v1.61.0 | Linter | ✅ | .golangci.yml, CI |
   | ruff | 0.4+ | Python lint | ✅ | pyproject.toml, CI |
   | markdownlint | latest | Doc lint | 🟡 | .markdownlint.json |
   | Docker | 27+ | Containers | ✅ | Compose, K8s |
   | Coder | v2.17.2 | Dev env | ✅ | .coder/template.tf |
   ```
   - Pro: Clear inventory, status tracking, config sync mapping
   - Con: Another table to maintain

B) **Extend existing tables**
   - Add to "Infrastructure Components" table
   - Pro: Reuses existing structure
   - Con: Mixes runtime with dev tools

C) **Separate TOOLS.md doc**
   - New doc: `docs/dev/design/operations/TOOLS.md`
   - Keep SOT for runtime only
   - Pro: Separation of concerns
   - Con: Not in SOT as user requested

D) **YAML section instead of table**
   ```yaml
   tools:
     required:
       go: "1.25+"
       node: "20.x"
     optional:
       docker: "27+"
   ```
   - Pro: Parseable
   - Con: Less readable than table

**Recommended**: A (Development Tools table in SOT)

---

## SECTION 9: Template Validation Against Patterns

### 🟡 Q9.1: Status Table Structure (HIGH)

**Context**: Existing docs use 7-dimension status table. Current template uses simplified version.

**Question**: Should template exactly match existing 7-dimension pattern?

**Options**:
A) **Yes, use exact 7-dimension pattern** (Recommended)
   ```jinja2
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
   - Pro: Consistency with existing, full tracking
   - Con: 14 YAML fields per doc (status + notes for each dimension)

B) **Simplified 3-dimension for generated docs**
   - Design, Code, Testing only
   - Pro: Simpler data files
   - Con: Loses granularity

C) **Make dimensions configurable per doc type**
   - Features: 7 dimensions
   - Services: 5 dimensions
   - Integrations: 4 dimensions
   - Pro: Flexible
   - Con: Inconsistent

D) **Status in SOT only, not in individual docs**
   - SOT has master status table
   - Individual docs reference SOT
   - Pro: Single source for status
   - Con: Docs not self-contained

**Recommended**: A (7-dimension exact match)

---

### 🟡 Q9.2: Implementation Checklist Structure (HIGH)

**Context**: Existing docs use multi-phase checklists (Phase 1-6). Template should match.

**Question**: How should template handle implementation checklists?

**Options**:
A) **Generate from structured data** (Recommended)
   ```yaml
   # In data file:
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

   # In template:
   {% for phase in implementation_phases %}
   ### Phase {{ phase.phase }}: {{ phase.name }}
   {% for task in phase.tasks %}
   - [ ] {{ task }}
   {% endfor %}
   {% endfor %}
   ```
   - Pro: Structured, reusable, can track progress
   - Con: Verbose data files

B) **Markdown string in data**
   ```yaml
   implementation_checklist: |
     ### Phase 1: Core
     - [ ] Create package
     ### Phase 2: Database
     - [ ] Create migration
   ```
   - Pro: Simple, flexible formatting
   - Con: Loses structure, can't programmatically track

C) **Template includes standard phases, data fills tasks**
   - Template hardcodes: Phase 1-6 names
   - Data only provides tasks per phase
   - Pro: Consistent phase structure
   - Con: Not flexible if different phase structure needed

D) **No checklist in template, link to GitHub issues**
   - Template generates link to GitHub project board
   - Pro: Integration with issue tracking
   - Con: Checklist not in doc

**Recommended**: A (Generate from structured data)

---

## SECTION 10: Security

### 🔴 Q10.1: YAML Safe Loading (CRITICAL)

**Context**: Gap analysis flagged PyYAML unsafe load as RCE risk.

**Question**: How should we safely parse YAML files?

**Options**:
A) **yaml.safe_load() only** (Recommended)
   ```python
   with open(data_file) as f:
       data = yaml.safe_load(f)  # SAFE - no code exec
   ```
   - Pro: Prevents RCE, standard practice
   - Con: Can't use advanced YAML features (Python objects)

B) **yaml.load() with warnings disabled**
   ```python
   data = yaml.load(f, Loader=yaml.FullLoader)  # UNSAFE
   ```
   - Pro: Full YAML features
   - Con: RCE vulnerability

C) **ruamel.yaml for safer full loading**
   - Use ruamel.yaml instead of PyYAML
   - Pro: Safer than PyYAML, more features
   - Con: Extra dependency

D) **Validate YAML structure post-load**
   - Use safe_load, then validate structure
   - Pro: Defense in depth
   - Con: Extra validation step

**Recommended**: A (safe_load only)

---

### 🔴 Q10.2: Jinja2 Template Sandboxing (CRITICAL)

**Context**: Jinja2 can execute Python. Need to sandbox.

**Question**: How should we sandbox Jinja2 template execution?

**Options**:
A) **SandboxedEnvironment with StrictUndefined** (Recommended)
   ```python
   from jinja2 import SandboxedEnvironment, StrictUndefined

   env = SandboxedEnvironment(
       loader=FileSystemLoader(...),
       undefined=StrictUndefined,  # Fail on undefined vars
       autoescape=False,  # Markdown, not HTML
   )
   ```
   - Pro: Prevents code execution, catches errors
   - Con: Slightly more restrictive

B) **Regular Environment, trust templates**
   ```python
   env = Environment(...)
   ```
   - Pro: Full Jinja2 features
   - Con: Potential code execution

C) **Disable dangerous filters/functions**
   ```python
   env = Environment(...)
   env.filters.pop('attr', None)  # Remove dangerous filters
   ```
   - Pro: Custom control
   - Con: Manual maintenance, easy to miss

D) **Read-only file system for template execution**
   - Run template rendering in container with RO filesystem
   - Pro: OS-level isolation
   - Con: Overhead, slower

**Recommended**: A (SandboxedEnvironment + StrictUndefined)

---

### 🟡 Q10.3: Secret Scanning (HIGH)

**Context**: Auto-generated docs could leak secrets if SOT contains them.

**Question**: Should we add secret scanning to validation pipeline?

**Options**:
A) **Yes, use gitleaks or truffleHog** (Recommended)
   - Run secret scanner on generated files before commit
   - Fail validation if secrets found
   - Pro: Catches accidental leaks
   - Con: False positives, slower validation

B) **Yes, regex patterns for common secrets**
   - Custom regex for API keys, tokens, passwords
   - Pro: Fast, customizable
   - Con: Incomplete, false positives/negatives

C) **No, trust manual review**
   - Human reviews PRs for secrets
   - Pro: No tooling overhead
   - Con: Humans miss things

D) **Yes, GitHub secret scanning**
   - Enable GitHub's built-in secret scanning
   - Pro: Automatic, maintained by GitHub
   - Con: Only scans after push (reactive)

**Recommended**: A (gitleaks in pre-commit validation)

---

## SECTION 11: Monitoring & Observability

### 🟡 Q11.1: Automation Failure Alerting (HIGH)

**Question**: How should we be notified if automation fails?

**Options**:
A) **Create GitHub issue on failure** (Recommended)
   - Automation creates issue with label `automation-failure`
   - Issue includes: logs, error message, stack trace
   - Pro: Persistent, traceable, searchable
   - Con: Could create issue spam if systemic failure

B) **GitHub Actions workflow failure notification**
   - Built-in email notification
   - Pro: Standard, zero setup
   - Con: Email noise, not actionable

C) **Slack webhook**
   - Post to #automation channel on failure
   - Pro: Real-time, team visible
   - Con: Requires Slack setup, ephemeral

D) **No alerting, check manually**
   - Review automation PRs to see if any failures
   - Pro: Zero overhead
   - Con: Delayed detection

**Recommended**: A (GitHub issue on failure)

---

### 🟡 Q11.2: Automation Metrics Dashboard (MEDIUM)

**Question**: Should we track automation metrics (success rate, time, etc.)?

**Options**:
A) **Yes, GitHub Actions metrics only**
   - Use built-in workflow run metrics
   - Pro: Zero setup, basic visibility
   - Con: Limited metrics

B) **Yes, custom dashboard (Grafana)**
   - Export metrics to Grafana
   - Track: success rate, regen time, PR count, etc.
   - Pro: Rich visualization, trends
   - Con: Infrastructure overhead

C) **Yes, CSV log file in repo**
   - Append to `.automation/metrics.csv` on each run
   - Pro: Simple, version controlled
   - Con: Manual analysis

D) **No metrics tracking**
   - Pro: Zero overhead
   - Con: No visibility into trends

**Recommended**: A or C (GitHub Actions metrics, or CSV log if we need more)

---

## Summary by Priority

### 🔴 P0 Critical (MUST Answer Before Any Code)

**Data Flow & SOT:**
- Q1.1: Data extraction from SOT
- Q1.2: Shared data organization
- Q1.3: Bootstrap procedure

**Loop Prevention:**
- Q2.1: Commit authorship detection
- Q2.2: Regeneration cooldown
- Q2.3: Dependabot loop prevention

**PR Workflow:**
- Q3.1: PR batching strategy

**Migration:**
- Q4.1: Pilot scope
- Q4.2: Data extraction tooling

**Validation:**
- Q5.1: YAML schema definition
- Q5.2: Link validation strategy

**Templates:**
- Q6.1: Template inheritance
- Q6.2: HTML comment metadata

**SOT Settings:**
- Q8.1: Automation settings in SOT

**Security:**
- Q10.1: YAML safe loading
- Q10.2: Jinja2 sandboxing

**Total: 17 P0 questions**

---

### 🟡 P1 High (Should Answer Before Detailed Design)

- Q3.2: Auto-approval for trusted changes
- Q4.3: Migration validation checkpoints
- Q5.3: SOT reference validation
- Q6.3: TOC auto-generation
- Q7.1: GitHub documentation sources
- Q7.2: Documentation style guides
- Q7.3: API documentation standards
- Q8.2: Tool versions in SOT
- Q9.1: Status table structure
- Q9.2: Implementation checklist structure
- Q10.3: Secret scanning
- Q11.1: Automation failure alerting

**Total: 12 P1 questions**

---

### 🟢 P2 Medium (Can Refine During Implementation)

- Q11.2: Automation metrics dashboard

**Total: 1 P2 question**

---

## Recommended Question Round Order

1. **Ask all P0 Critical (17 questions)** - Blocking for implementation plan
2. **Ask all P1 High (12 questions)** - Important for detailed design
3. **Defer P2 Medium (1 question)** - Can decide during implementation

---

**Next Step**: Begin P0 question round with user using AskUserQuestion tool.


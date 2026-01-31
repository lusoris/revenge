# User Answers - Documentation Restructuring

**Date**: 2026-01-31
**Status**: In Progress - Iterating

---

## Important Clarifications

**Versioning Strategy**:
- **v0.1.0** = MVP (Minimal Viable Product) ← First usable version
- **v0.x.x** = Incremental features
- **v1.0.0** = Full product per complete design specs

**External Dependencies**:
- Blackbeard (transcoding) is in development by Lawrence, NOT ready for MVP
- MVP will not include transcoding

---

## Section 1: Documentation Split Strategy ✅

### Q1.1: Doc Split Approach
**Answer**: Option D: Template-Based (Jinja2/variables)

**Implications**:
- Need to implement Jinja2 template engine
- Create `.template.md` files with variables and conditionals
- Generate both Claude and Wiki versions from templates
- Most flexible but requires template system setup

### Q1.2: Claude Docs Content (ALL selected)
**Include**:
- ✅ File paths (internal/content/movie/service.go)
- ✅ Internal API details
- ✅ Database schema specifics
- ✅ Code patterns and hints

### Q1.3: Wiki Docs Content (ALL selected)
**Include**:
- ✅ Screenshots and diagrams
- ✅ Step-by-step tutorials
- ✅ User perspective content
- ✅ External API documentation

### Q1.4: Wiki Location
**Answer**: Both (repo + auto-sync)

**Implications**:
- Store generated wiki in `docs/wiki/` (in git)
- Auto-sync to GitHub Wiki on commit/PR
- Need GitHub Wiki sync workflow
- Redundancy for backup

---

## Section 2: MVP Definition ✅ (pending version clarification)

### Q2.1: MVP Content Modules
**Answer**: Movies + TV Shows

### Q2.2: Frontend Scope
**Answer**: Full SvelteKit UI

### Q2.3: Playback Capabilities (without Blackbeard)
**Answer**: Direct play + HLS/DASH streaming

### Q2.4: Core Services (ALL)
**Confirmed**:
- ✅ Auth, User, Session, RBAC
- ✅ Metadata & Search
- ✅ Library Management
- ✅ Playback Service (no transcoding)

### Q2.5: External Integrations (ALL)
**Confirmed**:
- ✅ TMDB + TheTVDB
- ✅ Radarr + Sonarr
- ✅ OIDC/OAuth Providers
- ✅ Typesense Search

### Q2.6: Deferred to Post-MVP (ALL)
**Confirmed**:
- ⏸️ QAR/Adult Content
- ⏸️ LiveTV/DVR
- ⏸️ Photos + Comics + Books
- ⏸️ Advanced Playback Features
- ⏸️ Transcoding (waiting for Blackbeard)

### Q2.7: Milestone Structure
**Answer**: "we have a versioning schema... dunno what actually fits"

**CONFLICT FOUND**: VERSIONING.md says:
```
v0.1.x  Phase 1: Core Foundation (Auth, Library, Direct Play)
v0.2.x  Phase 1: Media Management
v0.3.x  Phase 1: MVP Complete  ← MVP at v0.3.x!
```

**Need to clarify**: MVP at v0.1.0 or v0.3.x?

---

## Section 3: Automation & Settings Sync ✅

### Q3.1: Settings Auto-Sync from SOT (ALL selected)
**Answer**: Sync all tool settings
- ✅ Language versions (Go 1.25.6, Python 3.12+, Node 20+)
- ✅ Formatter configs (gofmt, ruff, prettier)
- ✅ Linter configs (golangci-lint, eslint)
- ✅ LSP settings (gopls, typescript-language-server)

**Implications**:
- Need script: `scripts/sync-tool-settings.py`
- Reads SOT → updates all tool configs (.vscode, .zed, .jetbrains, .coder)
- CI/CD validation to ensure sync
- Prevents config drift across tools

---

## Section 4: Claude Skills ✅

### Q4.1: New Skills to Create (ALL selected)
**Answer**: Create all proposed skills
- ✅ `mvp-status` - Show MVP completion %, suggest next task
- ✅ `generate-wiki` - Generate wiki from design templates
- ✅ `verify-settings` - Validate tool configs match SOT
- ✅ `code-status` - Verify code implementation matches design

**Implications**:
- 4 new skills to implement
- Each needs SKILL.md with implementation
- Integration with existing automation

---

## Section 5: Wiki Generation ✅

### Q5.1: Wiki Trigger
**Answer**: On every commit (automated)

**Implications**:
- Git hook or GitHub Action on push
- Auto-generate wiki when design docs change
- Auto-sync to GitHub Wiki
- Need fast generation (< 30s for CI/CD)

---

## Section 6: Versioning Strategy ✅

### Q6.1: MVP Version Number
**Answer**: Keep VERSIONING.md plan (MVP = v0.3.x)

**Rationale** (from guidance document):
- Media servers need quality (not "barely works")
- Infrastructure-heavy projects need foundation first
- v0.1.x = Core Foundation (Auth, Library, Direct Play)
- v0.2.x = First Content Type
- v0.3.x = MVP Complete (Movies + TV + Full UI + Direct play/HLS)
- v1.0.0 = Full design spec implemented

**This matches industry best practices** for complex systems (PostgreSQL, Kubernetes, etc.)

---

## Summary of All Answers ✅

### Documentation Split
- **Method**: Template-based (Jinja2)
- **Claude content**: Implementation details, file paths, DB schema, code patterns
- **Wiki content**: Screenshots, tutorials, user perspective, external API docs
- **Location**: Both repo (docs/wiki/) + auto-sync to GitHub Wiki
- **Trigger**: On every commit (automated)

### MVP Definition (v0.3.x)
- **Content**: Movies + TV Shows
- **Frontend**: Full SvelteKit UI
- **Playback**: Direct play + HLS/DASH (no transcoding until Blackbeard ready)
- **Services**: Auth, User, Session, RBAC, Metadata, Search, Library, Playback
- **Integrations**: TMDB, TheTVDB, Radarr, Sonarr, OIDC, Typesense
- **Deferred**: QAR, LiveTV, Photos, Comics, Books, Advanced Playback, Transcoding

### Automation
- **Settings Sync**: ALL (language versions, formatters, linters, LSP) from SOT
- **New Skills**: ALL 4 (mvp-status, generate-wiki, verify-settings, code-status)
- **Wiki Generation**: Automated on every commit

### Versioning
- **MVP**: v0.3.x (not v0.1.0)
- **v1.0.0**: Full design spec
- **Follow**: Existing VERSIONING.md plan

---

---

## Section 7: Testing Requirements ✅

### Q7.1: Testing Before Implementation (ALL selected)
**Answer**: Full testing suite
- ✅ Test all pipelines on sample docs
- ✅ Test template generation with samples
- ✅ Lint all generated docs
- ✅ End-to-end validation (template → generate → validate → deploy)

**Implications**:
- Create test dataset before implementation
- Validate every stage before moving to next
- No production changes until all tests pass

---

## Section 8: Risk Mitigation ✅

### Q8.1: Handling Breaking Changes
**Answer**: Multi-layer safety
- ✅ Keep .analysis/ backups until confirmed working
- ✅ Test on feature branch first
- ✅ Git tag before major changes
- ❌ NOT using symlinks (clean break)

**Implications**:
- .analysis/ dir stays until full validation
- All work on feature branch, merge after validation
- Create git tag: `pre-restructure-2026-01-31`

---

## Section 9: Implementation Sequence ✅

### Q9.1: Approval
**Answer**: Approved as-is

**7-Phase Sequence**:
1. Phase 0: Planning & Approval (DONE)
2. Phase 1: MVP Definition
3. Phase 2: Template Creation
4. Phase 3: Automation Scripts
5. Phase 4: Pipeline Integration
6. Phase 5: Documentation Migration
7. Phase 6: Skills & Tooling
8. Phase 7: Validation & Rollout

---

## Section 10: External Sources ✅

### Q10.1: Sources to Fetch (ALL selected)
**Answer**: Fetch all required sources
- ✅ Jinja2 documentation (Required for templates)
- ✅ GitHub Wiki format/API (Required for auto-sync)
- ✅ Go AST parsing (For code-status skill)
- ✅ Markdown linting rules (For validation)

**Implications**:
- Add to SOURCES.yaml
- Fetch before Phase 2 begins
- Save in docs/dev/sources/

---

## ✅ ALL QUESTIONS ANSWERED

**Status**: Ready to proceed

**Next Steps**:
1. Fetch external sources (Jinja2, GitHub Wiki, Go AST, markdownlint)
2. Create detailed implementation plan (Phase 1-7 breakdown)
3. Create MVP definition document
4. Create roadmap document (v0.1.x → v0.3.x MVP → v1.0.0)
5. Begin Phase 1 implementation

**Awaiting**: User approval to start fetching sources and creating implementation plan

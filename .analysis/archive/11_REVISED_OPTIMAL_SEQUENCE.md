# REVISED OPTIMAL SEQUENCE

**Date**: 2026-01-31
**Based on user feedback**: "scaffold everything first -> then do the automation stuff rework to build all docs the way we want them to be and then in the end we will create the contents for all gaps, step by step..."

---

## WHY THIS SEQUENCE IS BETTER

**Original Problem**:
- Design writing (Phase 1) blocks everything for weeks/months
- Automation can't be built without knowing structure
- Risk of changing design format mid-automation

**New Approach**:
1. **Scaffold everything** → see full scope, establish structure
2. **Build automation** → templates, pipelines, wiki generation working
3. **Fill in content** → systematic design writing with automation already working

**Benefits**:
- Immediate visibility of full scope
- Automation ready before heavy content work
- Can test automation with scaffolds
- Content work is last step, not blocking step

---

## NEW PHASE SEQUENCE

### Phase 0: Planning & Approval ✅ COMPLETE
- [x] Analysis, questions, decisions, source fetching

---

### Phase 1: SCAFFOLDING (Quick Structural Work)

**Goal**: Create structure for ALL missing/incomplete docs

**Duration**: 1-2 days

#### 1.1: Create Scaffold Template
- [ ] Create `docs/dev/design/.templates/SCAFFOLD_TEMPLATE.md`
- [ ] Test template with one example
- [ ] Validate structure

#### 1.2: Scaffold ALL Missing Documents
**Scaffold these quickly (structure only, no content yet)**:

##### Critical Technical (2 docs)
- [ ] API.md - add detailed structure if needed
- [ ] FRONTEND.md - add detailed structure if needed

##### Content Modules (3 docs)
- [ ] MUSIC_MODULE.md - enhance scaffold
- [ ] AUDIOBOOK_MODULE.md - enhance scaffold
- [ ] BOOK_MODULE.md - enhance scaffold

##### Wiki Integrations (6 docs)
- [ ] WIKIPEDIA.md - scaffold
- [ ] FANDOM.md - scaffold
- [ ] TVTROPES.md - scaffold
- [ ] BABEPEDIA.md - scaffold
- [ ] BOOBPEDIA.md - scaffold
- [ ] IAFD.md - scaffold

##### Scrobbling Services (4 docs)
- [ ] TRAKT.md - enhance scaffold (20% → scaffold)
- [ ] LETTERBOXD.md - enhance scaffold
- [ ] LISTENBRAINZ.md - enhance scaffold
- [ ] SIMKL.md - enhance scaffold

##### Auth Integrations (2 docs)
- [ ] AUTHENTIK.md - enhance scaffold (50% → scaffold)
- [ ] KEYCLOAK.md - enhance scaffold

##### Adult Metadata (4 docs)
- [ ] STASH.md - enhance scaffold
- [ ] THEPORNDB.md - enhance scaffold
- [ ] PORNHUB.md - enhance scaffold
- [ ] THENUDE.md - enhance scaffold

##### Technical Documentation (3 docs)
- [ ] NOTIFICATIONS.md - enhance scaffold
- [ ] EMAIL.md - enhance scaffold
- [ ] WEBHOOKS.md - enhance scaffold

##### Operations Documentation (4 docs)
- [ ] SETUP.md - enhance scaffold
- [ ] DEVELOPMENT.md - enhance scaffold
- [ ] GITFLOW.md - enhance scaffold
- [ ] BRANCH_PROTECTION.md - enhance scaffold

**Total**: ~30 scaffolds to create/enhance

#### 1.3: Link Scaffolds
- [ ] Link all scaffolds from category INDEX.md files
- [ ] Link from DESIGN_INDEX.md
- [ ] Add to status tracking tables

#### 1.4: Commit Scaffolds
- [ ] Create feature branch: `design-scaffolding`
- [ ] Commit all scaffolds
- [ ] Push to remote
- [ ] Keep branch (will merge after content complete)

**Exit Criteria**:
- ✅ All ~30 docs scaffolded
- ✅ All linked properly
- ✅ Structure validated
- ✅ Committed to feature branch

---

### Phase 2: MVP Definition & Coding Roadmap

**Goal**: Define implementation sequence

**Duration**: 3-5 days

#### 2.1: Create MVP Definition
- [ ] Create `planning/MVP_DEFINITION.md`
- [ ] Define v0.3.x scope
- [ ] List included/excluded features

#### 2.2: Create Implementation Roadmap
- [ ] Create `planning/IMPLEMENTATION_ROADMAP.md`
- [ ] Define milestones (v0.1.x → v1.0.0)
- [ ] Dependency graph

#### 2.3: Create Detailed Phase TODOs
- [ ] `planning/v0.1.x_TODO.md` (100+ tasks - Core Foundation)
- [ ] `planning/v0.2.x_TODO.md` (100+ tasks - First Content Type)
- [ ] `planning/v0.3.x_TODO.md` (100+ tasks - MVP Complete)
- [ ] `planning/v0.4.x_TODO.md` (Music)
- [ ] `planning/v0.5.x_TODO.md` (Transcoding)
- [ ] `planning/v0.6.x_TODO.md` (Audiobooks + Books)
- [ ] `planning/v0.7.x_TODO.md` (Photos + Comics)
- [ ] `planning/v0.8.x_TODO.md` (LiveTV)
- [ ] `planning/v0.9.x_TODO.md` (QAR + Advanced Playback)
- [ ] `planning/v1.0.0_TODO.md` (Polish)

#### 2.4: Update SOURCE_OF_TRUTH
- [ ] Add Planning & Roadmap section
- [ ] Link all planning docs

**Exit Criteria**:
- ✅ MVP defined
- ✅ 8+ milestones with 100+ tasks each
- ✅ Linked from SOT
- ✅ User approves

---

### Phase 3: Template System

**Goal**: Create Jinja2 templates for dual doc generation

**Duration**: 2-3 days

#### 3.1: Design Template Structure
- [ ] Research Jinja2
- [ ] Define variables (from SOT)
- [ ] Define conditionals ({{ if claude }}, {{ if wiki }})

#### 3.2: Create Templates
- [ ] `CLAUDE_DOC_TEMPLATE.template.md`
- [ ] `WIKI_DOC_TEMPLATE.template.md`

#### 3.3: Test with Scaffolds
- [ ] Convert 2-3 scaffolds to templates
- [ ] Generate both versions
- [ ] Validate quality

**Exit Criteria**:
- ✅ Templates work correctly
- ✅ Can generate Claude + Wiki versions
- ✅ User approves

---

### Phase 4: Automation Scripts

**Goal**: Build all automation infrastructure

**Duration**: 3-4 days

#### 4.1: Wiki Generation
- [ ] Create `scripts/generate-wiki.py`
- [ ] Integrate Jinja2
- [ ] Optimize for speed (< 30s)

#### 4.2: GitHub Wiki Sync
- [ ] Implement sync mechanism
- [ ] Test with actual wiki

#### 4.3: Settings Sync
- [ ] Create `scripts/sync-tool-settings.py`
- [ ] Read SOT → update all tool configs
- [ ] Validation mode

#### 4.4: Code Status Verification
- [ ] Create `scripts/verify-code-status.py`
- [ ] Use Go AST to parse codebase
- [ ] Compare with design status
- [ ] Auto-update status tables

**Exit Criteria**:
- ✅ All scripts working
- ✅ Wiki gen < 30s
- ✅ Settings sync accurate
- ✅ Code verification works

---

### Phase 5: Pipeline Integration

**Goal**: Integrate automation into pipelines

**Duration**: 2-3 days

#### 5.1: Extend Doc Pipeline
- [ ] Add Stage 7: Wiki Generation
- [ ] Modify Stage 3: Add MVP filtering, code status

#### 5.2: Extend Source Pipeline
- [ ] Verify new sources included

#### 5.3: CI/CD Workflows
- [ ] Create `wiki-sync.yml`
- [ ] Create `settings-sync-validation.yml`
- [ ] Update `doc-validation.yml`
- [ ] Update `fetch-sources.yml`

#### 5.4: End-to-End Testing
- [ ] Test full pipeline with scaffolds
- [ ] Fix any issues

**Exit Criteria**:
- ✅ Pipelines extended
- ✅ CI/CD working
- ✅ E2E tests pass

---

### Phase 6: Documentation Migration

**Goal**: Convert existing docs to template format

**Duration**: 3-5 days

#### 6.1: Convert to Templates
**Convert ALL existing docs** (including scaffolds):
- [ ] Feature docs (41)
- [ ] Service docs (15)
- [ ] Integration docs (65)
- [ ] Operations docs (7)
- [ ] Technical docs (6)
- [ ] Architecture docs (5)

**For each**:
- Identify Claude vs Wiki content
- Extract variables
- Convert to .template.md

#### 6.2: Generate Both Versions
- [ ] Run generate-wiki.py on all
- [ ] Verify Claude versions
- [ ] Verify Wiki versions

#### 6.3: Link Validation
- [ ] Check all links
- [ ] Fix broken links

#### 6.4: Quality Pass
- [ ] Lint all docs
- [ ] Format check
- [ ] User approval

**Exit Criteria**:
- ✅ All 136+ docs converted
- ✅ Both versions generated
- ✅ Zero broken links
- ✅ Linting clean

---

### Phase 7: Skills & Tooling

**Goal**: Create Claude Code skills

**Duration**: 2-3 days

#### 7.1: Create Skills
- [ ] `mvp-status` - Show completion %
- [ ] `generate-wiki` - Generate wiki docs
- [ ] `verify-settings` - Validate settings
- [ ] `code-status` - Show implementation status

#### 7.2: Update Documentation
- [ ] Update `.claude/CLAUDE.md`
- [ ] Add usage examples

#### 7.3: Test Skills
- [ ] Test individually
- [ ] Test integration
- [ ] User acceptance

**Exit Criteria**:
- ✅ 4 skills working
- ✅ Documentation updated
- ✅ User can use skills

---

### Phase 8: DESIGN WRITING (Content Creation)

**Goal**: Fill in ALL scaffolds with 99% perfection content

**Duration**: Variable (170-230 hours collaborative)

**This is where the heavy design work happens**

#### 8.1: Priority Tier 1 - FOUNDATION (CRITICAL)
**~40-60 hours**

1. **API.md** - Complete endpoint specifications
   - ALL endpoints for ALL services
   - Request/Response schemas
   - Error codes
   - Rate limiting
   - Auth flow
   - Pagination
   - Filtering/Sorting

2. **FRONTEND.md** - Complete component architecture
   - SvelteKit structure
   - Component patterns
   - State management
   - API integration
   - Testing strategies
   - Performance
   - Accessibility
   - Theming
   - Forms

#### 8.2: Priority Tier 2 - CONTENT MODULES (MVP)
**~30-40 hours**

3. **MUSIC_MODULE.md** - Complete design
4. **AUDIOBOOK_MODULE.md** - Complete design
5. **BOOK_MODULE.md** - Complete design

#### 8.3: Priority Tier 3 - HIGH-PRIORITY INTEGRATIONS
**~40-50 hours**

6. **WIKIPEDIA.md** - Complete design
7. **AUTHENTIK.md** - 50% → 100%
8. **KEYCLOAK.md** - 50% → 100%
9. **TRAKT.md** - Complete design
10. **NOTIFICATIONS.md** - Complete design
11. **SETUP.md** - Complete operations guide

#### 8.4: Priority Tier 4 - MEDIUM-PRIORITY
**~40-50 hours**

12-25. All remaining scrobbling, wiki, adult metadata, technical docs

#### 8.5: Priority Tier 5 - POLISH
**~20-30 hours**

26-32. Operations docs, source references, final review

#### 8.6: Design Validation
- [ ] All designs follow template
- [ ] All have architecture diagrams
- [ ] All link to sources
- [ ] All cross-references valid
- [ ] Markdown linting passes
- [ ] User review and approval

**Exit Criteria**:
- ✅ Zero "🔴 PLANNED" status
- ✅ All v1.0 features ✅ or 🟡
- ✅ API.md ✅ (foundational)
- ✅ FRONTEND.md ✅ (foundational)
- ✅ All 3 content modules ✅
- ✅ User approves 99% perfection achieved

---

### Phase 9: Validation & Rollout

**Goal**: Final validation and deployment

**Duration**: 1-2 days

#### 9.1: Full Validation
- [ ] Run all pipelines
- [ ] Lint all files
- [ ] Validate all links
- [ ] Test all workflows
- [ ] Test all skills
- [ ] Verify settings sync
- [ ] Run code status check

#### 9.2: Quality Checks
- [ ] Zero linting errors
- [ ] Zero broken links
- [ ] All CI/CD passes
- [ ] Wiki renders correctly
- [ ] Settings match SOT 100%

#### 9.3: User Review
- [ ] Review all changes
- [ ] Test wiki generation
- [ ] Test new skills
- [ ] Approve rollout

#### 9.4: Git Operations
- [ ] Merge design-scaffolding branch
- [ ] Commit automation work
- [ ] Create tag: `pre-v1-design-complete-2026-01-31`
- [ ] Push to develop
- [ ] Create PR
- [ ] Merge

#### 9.5: Cleanup
- [ ] Archive .analysis/ to `.backup/`
- [ ] Update TODO.md
- [ ] Create follow-up issues

**Exit Criteria**:
- ✅ All validation passed
- ✅ Git tags created
- ✅ Changes merged
- ✅ No regressions
- ✅ Ready for implementation

---

## NEW DEPENDENCY GRAPH

```
Phase 0 (Planning) ✅
    ↓
Phase 1 (Scaffolding) ← FAST (1-2 days)
    ↓
Phase 2 (MVP & Roadmap)
    ↓
Phase 3 (Templates)
    ↓
Phase 4 (Automation Scripts)
    ↓
Phase 5 (Pipeline Integration)
    ↓
Phase 6 (Migration to Templates)
    ↓
Phase 7 (Skills & Tooling)
    ↓
Phase 8 (Design Writing) ← HEAVY WORK, but automation ready
    ↓
Phase 9 (Validation & Rollout)
```

---

## TIMELINE COMPARISON

### Old Sequence
```
Phase 1: Design Writing (170-230 hours) ← BLOCKS EVERYTHING
    ↓ (weeks/months)
Phase 2-8: Automation & rollout (16-25 days)
```

### New Sequence
```
Phase 1: Scaffolding (1-2 days) ← FAST
    ↓
Phases 2-7: Automation (16-25 days) ← BUILD INFRASTRUCTURE
    ↓
Phase 8: Design Writing (170-230 hours) ← HEAVY WORK, automation ready
    ↓
Phase 9: Validation (1-2 days)
```

**Key Advantage**: Automation is ready BEFORE heavy content work begins!

---

## NEXT IMMEDIATE STEPS

1. ✅ Master plan created
2. ✅ Revised sequence created
3. → Ask user if any questions about this sequence
4. → Verify source fetch completed
5. → Start Phase 1.1: Create scaffold template
6. → Begin scaffolding all docs (Phase 1.2)

---

**STATUS**: Optimal sequence defined - ready to proceed with Phase 1 (Scaffolding)

# REVISED Implementation Sequence

**Date**: 2026-01-31
**Critical Change**: Design completion BEFORE coding roadmap

---

## User Clarification

> "the roadmap is only about coding so our goal is to first finish the full design spec for the current full roadmap goal of the version 1 release..."

**This means**:
1. ✅ Complete ALL design documentation for v1.0 scope FIRST
2. ✅ Then create coding roadmap (v0.1.x → v1.0.0)
3. ✅ Then implement templates, automation, etc.
4. ✅ Then start actual coding

---

## NEW Phase Sequence

### Phase 0: Planning & Approval ✅ COMPLETE
- Analysis, questions, decisions, source fetching

### **NEW Phase 1: Design Completion** (MOST IMPORTANT)
**Goal**: Complete ALL design docs for v1.0 scope

**Sub-phases**:

#### 1.1: Gap Analysis
- Identify ALL missing design docs
- Categorize by feature area
- Prioritize by MVP vs post-MVP
- Create complete list

#### 1.2: Design Scaffolding
- Create placeholder docs for ALL gaps
- Use scaffold template
- Link from DESIGN_INDEX.md
- Add to status tracking

#### 1.3: Design Writing (Iterative)
**For each missing/incomplete design**:
1. Research (check sources, similar systems)
2. Draft architecture
3. Define database schema
4. Define API endpoints
5. Define integrations
6. Define testing strategy
7. Review and validate
8. Mark as complete

**Priority Order**:
1. **MVP-Critical** (for v0.3.x)
   - Any incomplete service docs
   - Any incomplete integration docs
   - Any incomplete feature docs for Movies/TV

2. **v0.4-v0.9 Features**
   - Music, Audiobooks, Podcasts modules
   - Photos, Comics, Books modules
   - LiveTV/DVR
   - QAR (adult content)
   - Advanced playback features
   - Plugin architecture

3. **v1.0 Polish**
   - Performance optimization docs
   - Deployment docs
   - Operations docs

#### 1.4: Design Validation
- All designs follow template
- All designs have clear architecture
- All designs link to sources
- All cross-references valid
- Status tables updated

**Exit Criteria**:
- ✅ Zero "🔴 PLANNED" status (all at least 🟡 or ✅)
- ✅ All v1.0 features have design docs
- ✅ All designs validated and reviewed
- ✅ User approves design completeness

**Estimated Duration**: 2-4 weeks (depends on gaps)

---

### Phase 2: MVP Definition & Coding Roadmap
**Goal**: Define what to code and in what order

**Tasks**:
1. Create MVP_DEFINITION.md (v0.3.x scope)
2. Create IMPLEMENTATION_ROADMAP.md (v0.1.x → v1.0.0)
3. Create VERY DETAILED phase TODOs (100+ tasks each)
4. Update SOURCE_OF_TRUTH.md with all links

**Exit Criteria**:
- ✅ MVP clearly defined
- ✅ Roadmap has 8+ milestones
- ✅ Each milestone has detailed TODO (100+ tasks)
- ✅ All linked from SOT
- ✅ User approves roadmap

**Estimated Duration**: 3-5 days

---

### Phase 3: Template System
(Same as before - create Jinja2 templates)

### Phase 4: Automation Scripts
(Same as before - wiki gen, settings sync, code status)

### Phase 5: Pipeline Integration
(Same as before - extend doc pipeline)

### Phase 6: Documentation Migration
(Same as before - convert to templates)

### Phase 7: Skills & Tooling
(Same as before - 4 new skills)

### Phase 8: Validation & Rollout
(Same as before - final validation)

---

## Phase 1 Detailed Plan: Design Completion

### Step 1: Gap Analysis

**Objective**: Find ALL missing design docs

**Method**:
1. Read VERSIONING.md - list all v1.0 features
2. Check existing design docs - what exists?
3. Compare lists - what's missing?
4. Create gap report

**Expected Gaps** (preliminary):

**Services**:
- ❓ Playback Service - exists but might need detail
- ❓ Scrobbling Service - might be incomplete
- ❓ Analytics Service - might be incomplete
- ❓ Notification Service - might be incomplete

**Features**:
- ❌ Music Module - missing or incomplete?
- ❌ Audiobook Module - scaffold only?
- ❌ Podcast Module - scaffold only?
- ❌ Photos Module - planned but not designed?
- ❌ Comics Module - planned but not designed?
- ❌ Books Module - planned but not designed?
- ❌ LiveTV/DVR - planned but not designed?
- ❌ Skip Intro feature - missing?
- ❌ Trickplay feature - missing?
- ❌ Syncplay feature - missing?
- ❌ Plugin Architecture - missing?

**Integrations**:
- ❌ Music metadata (MusicBrainz, Spotify, Last.fm, Discogs)
- ❌ Book metadata (OpenLibrary, Goodreads, Audible)
- ❌ Comic metadata (ComicVine)
- ❌ Photo integrations (if any)
- ❌ Scrobbling integrations (Trakt, Last.fm, Letterboxd, Simkl, ListenBrainz)
- ❌ Wiki integrations (Wikipedia, Fandom, TVTropes, etc.)

**Operations**:
- ❓ Deployment guides - complete?
- ❓ Performance tuning - documented?
- ❓ Backup/restore - documented?
- ❓ Monitoring - documented?

**Deliverable**: `.analysis/09_DESIGN_GAPS.md`

---

### Step 2: Create Scaffold Template

**File**: `docs/dev/design/.templates/SCAFFOLD_TEMPLATE.md`

```markdown
# {{ feature_name }}

<!-- STATUS: 🔴 DESIGN INCOMPLETE - Scaffold only -->

<!-- SOURCES: [to be added] -->

<!-- DESIGN: [to be linked] -->

## Status

| Dimension | Status |
|-----------|--------|
| Design | 🔴 |
| Sources | 🔴 |
| Instructions | 🔴 |
| Code | 🔴 |
| Linting | 🔴 |
| Unit Testing | 🔴 |
| Integration Testing | 🔴 |

---

## Overview

> **Purpose**: [One sentence description]

**Scope**: {{ scope_description }}

**Target Version**: {{ target_version }}

**Priority**: {{ priority_level }}

---

## Design TODO

This is a scaffold document. The following sections need to be completed:

### Architecture
- [ ] Define overall architecture
- [ ] Define component interactions
- [ ] Define data flow
- [ ] Create architecture diagram

### Database Schema
- [ ] Define tables and columns
- [ ] Define relationships
- [ ] Define indexes
- [ ] Define constraints
- [ ] Create ER diagram

### API Endpoints
- [ ] List all endpoints
- [ ] Define request/response formats
- [ ] Define authentication requirements
- [ ] Define error responses

### External Integrations
- [ ] List required external services
- [ ] Define integration points
- [ ] Define data synchronization
- [ ] Define error handling

### Testing Strategy
- [ ] Define unit test approach
- [ ] Define integration test approach
- [ ] Define test data requirements
- [ ] Define performance test criteria

### Implementation Notes
- [ ] List Go packages needed
- [ ] List dependencies
- [ ] List potential challenges
- [ ] List open questions

---

## Dependencies

**Depends on**:
- {{ dependency_list }}

**Blocks**:
- {{ blocks_list }}

**Related to**:
- {{ related_list }}

---

## Timeline

- **Design Start**: [TBD]
- **Design Complete**: [TBD]
- **Implementation Start**: [TBD]
- **Target Release**: {{ target_version }}

---

## Open Questions

1. {{ question_1 }}
2. {{ question_2 }}
...

---

**Last Updated**: {{ date }}
**Status**: 🔴 SCAFFOLD - Awaiting detailed design
```

---

### Step 3: Scaffold All Missing Docs

**For each gap identified**:
1. Copy scaffold template
2. Fill in basic metadata
3. Save to appropriate category
4. Link from category INDEX.md
5. Link from DESIGN_INDEX.md
6. Add to status tracking table

**Output**: 20-40 new scaffold files (depending on gaps)

---

### Step 4: Prioritize Design Work

**Create priority list**:

**Tier 1 - MVP Critical** (must design first):
- Anything needed for v0.1.x, v0.2.x, v0.3.x
- Core services
- Movie/TV modules
- Essential integrations

**Tier 2 - v0.4-v0.6** (design next):
- Additional content types
- Advanced features
- Optional integrations

**Tier 3 - v0.7-v0.9** (design last):
- Plugin system
- Nice-to-have features
- Performance optimizations

**Tier 4 - v1.0 Polish**:
- Documentation polish
- Operations guides
- Deployment automation

---

### Step 5: Design Writing Process

**For EACH missing/incomplete design** (iterative):

1. **Research** (30 min - 2 hours)
   - Read related designs
   - Check external sources
   - Study similar systems
   - Research best practices

2. **Architecture** (1-3 hours)
   - Define components
   - Define interactions
   - Create diagrams
   - Document patterns

3. **Database Schema** (1-2 hours)
   - Define tables
   - Define relationships
   - Define indexes
   - Create ER diagram

4. **API Design** (1-2 hours)
   - List endpoints
   - Define request/response
   - Define authentication
   - Define errors

5. **Integrations** (1-3 hours)
   - List external services
   - Define integration points
   - Define data sync
   - Define error handling

6. **Testing Strategy** (30 min - 1 hour)
   - Unit test approach
   - Integration test approach
   - Test data needs
   - Performance criteria

7. **Review** (30 min - 1 hour)
   - Self-review
   - Check template compliance
   - Validate links
   - Run linter

8. **Finalize** (15 min)
   - Update status to ✅ or 🟡
   - Link to related docs
   - Add to DESIGN_INDEX
   - Commit changes

**Estimated time per design**: 4-12 hours

**Total time** (if 30 missing designs): 120-360 hours = 15-45 days

---

### Step 6: Track Progress

**Create tracking document**: `.analysis/10_DESIGN_PROGRESS.md`

```markdown
# Design Completion Progress

**Goal**: Complete all design docs for v1.0

**Total Designs Needed**: [count]
**Currently Complete**: [count]
**In Progress**: [count]
**Not Started**: [count]

## Progress by Category

### Features
- [ ] Movie Module ✅ (already complete)
- [ ] TV Show Module ✅ (already complete)
- [ ] Music Module 🔴 (not started)
- [ ] Audiobook Module 🔴 (scaffold only)
...

### Services
- [ ] Auth Service ✅ (already complete)
- [ ] User Service ✅ (already complete)
...

### Integrations
- [ ] TMDB ✅ (already complete)
...

## This Week's Goals
- [ ] Complete Music Module design
- [ ] Complete Audiobook Module design
- [ ] Start Podcast Module design

## Blockers
- Need clarification on [topic]
- Waiting for [dependency]
```

---

## Phase 1 Estimated Timeline

**Optimistic**: 2-3 weeks (if few gaps, fast writing)
**Realistic**: 3-4 weeks (moderate gaps, thorough design)
**Pessimistic**: 5-6 weeks (many gaps, complex designs)

**Recommendation**: Plan for 4 weeks (1 month)

---

## Updated SOURCE_OF_TRUTH.md Structure

Add these sections:

```markdown
## Documentation Status

### Design Completion
- **Total Designs**: 179+ docs
- **Complete (✅)**: [count]
- **Partial (🟡)**: [count]
- **Scaffold (🔴)**: [count]
- **Progress Tracker**: [.analysis/10_DESIGN_PROGRESS.md](.analysis/10_DESIGN_PROGRESS.md)

### Design Gaps
- **Gap Analysis**: [.analysis/09_DESIGN_GAPS.md](.analysis/09_DESIGN_GAPS.md)
- **Missing Designs**: [list major gaps]

## Planning & Roadmap

### Phase 1: Design Completion (Current)
- **Goal**: Complete ALL v1.0 design docs
- **Status**: In Progress
- **Progress**: [link to tracker]

### Phase 2: MVP Definition (Next)
- **File**: [planning/MVP_DEFINITION.md](planning/MVP_DEFINITION.md)
- **Scope**: v0.3.x features

### Coding Roadmap (After Design Complete)
- **File**: [planning/IMPLEMENTATION_ROADMAP.md](planning/IMPLEMENTATION_ROADMAP.md)
- **Milestones**: v0.1.x → v1.0.0
- **Phase TODOs**: [list]
  - [v0.1.x TODO](planning/v0.1.x_TODO.md) - Core Foundation
  - [v0.2.x TODO](planning/v0.2.x_TODO.md) - First Content Type
  - [v0.3.x TODO](planning/v0.3.x_TODO.md) - MVP Complete
  - [v0.4.x TODO](planning/v0.4.x_TODO.md) - Third Content Type
  - [v0.5.x TODO](planning/v0.5.x_TODO.md) - Transcoding
  - [v0.6-v0.9.x TODOs](planning/) - Feature Expansion
```

---

## Next Immediate Steps

1. ✅ Wait for source fetch to complete
2. ✅ Verify all sources fetched correctly
3. **→ Start Phase 1: Gap Analysis**
   - Create `.analysis/09_DESIGN_GAPS.md`
   - Identify ALL missing designs
   - Create priority list
4. **→ Create Scaffolds**
   - Scaffold ALL missing docs
   - Link from indexes
5. **→ Begin Design Writing**
   - Start with Tier 1 (MVP-critical)
   - Iterate through priority list

---

**Question for User**: Should I start Phase 1 (gap analysis) now, or wait for you to review this revised plan first?

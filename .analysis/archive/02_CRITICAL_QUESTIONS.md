# CRITICAL QUESTIONS - Documentation Restructuring & MVP Planning

**Date**: 2026-01-31 (Updated)
**Status**: PARTIALLY ANSWERED - Phase 1 questions needed

---

## ✅ ANSWERED QUESTIONS (Sections 1-10)

All questions from initial round have been answered. See [03_ANSWERS.md](03_ANSWERS.md) for details.

**Summary**:
- ✅ Doc split: Template-based (Jinja2)
- ✅ MVP version: v0.3.x (not v0.1.0)
- ✅ Settings sync: ALL
- ✅ New skills: ALL 4
- ✅ Wiki trigger: On every commit
- ✅ Testing: Full suite
- ✅ Risk mitigation: Branch + backups + git tag
- ✅ Sequence: Approved (but now revised)
- ✅ External sources: ALL fetched

---

## 🔄 REVISED UNDERSTANDING

**Critical Change**: User clarified that design completion comes BEFORE coding roadmap.

**New Sequence**:
1. **Phase 1**: Complete ALL v1.0 design docs (CURRENT FOCUS)
2. **Phase 2**: MVP definition + coding roadmap
3. **Phase 3+**: Templates, automation, etc.

This creates NEW questions for Phase 1...

---

## ⏸️ NEW QUESTIONS - Phase 1: Design Completion

### Q11: Design Gap Analysis Scope

**Q11.1**: When analyzing design gaps, what categories should we check?

Options:
- [ ] Only features explicitly listed in VERSIONING.md
- [ ] All features from current design docs + VERSIONING.md
- [ ] All features PLUS anything mentioned in TODOs/discussions
- [ ] Comprehensive (features + services + integrations + operations)

**Q11.2**: How should we categorize missing designs?

Proposed categories:
- MVP-Critical (v0.1-v0.3.x)
- Post-MVP Features (v0.4-v0.9.x)
- v1.0 Only
- Future/Deferred

Is this categorization correct?

**Q11.3**: Should gap analysis include checking EXISTING docs for completeness?

Some docs might exist but be incomplete (missing sections, TODOs, etc.). Should we:
- [ ] Just find completely missing docs
- [ ] Also flag incomplete/partial docs
- [ ] Rate existing docs (Complete, Partial, Scaffold)

---

### Q12: Design Scaffolding Strategy

**Q12.1**: What level of detail for scaffolds?

Options:
- [ ] Minimal (just title + "TODO: Design this")
- [ ] Basic (title + overview + empty sections)
- [ ] Structured (template with all sections, each marked TODO)
- [ ] Rich (template + initial research/notes)

**Recommendation**: Structured (template with TODOs in each section)

**Q12.2**: Should scaffolds include research links?

When creating scaffold, should we:
- [ ] Add links to relevant external sources
- [ ] Add links to similar existing designs
- [ ] Add questions/open issues
- [ ] Just structure, no content

**Q12.3**: Priority for scaffold creation?

- [ ] Create ALL scaffolds at once (complete picture)
- [ ] Create Tier 1 (MVP) scaffolds first
- [ ] Create scaffolds as we analyze each category

---

### Q13: Design Writing Process

**Q13.1**: Who writes the designs?

Options:
- [ ] User writes all designs (I assist with research/structure)
- [ ] I draft designs, user reviews/approves
- [ ] Collaborative (I draft, user refines, iterate)
- [ ] Split by complexity (I do simple, user does complex)

**Q13.2**: How much detail in each design?

For each design doc, should it include:
- [ ] High-level architecture only
- [ ] Architecture + database schema
- [ ] Architecture + schema + API endpoints
- [ ] Architecture + schema + API + testing strategy
- [ ] Full detail (all sections from template)

**Recommendation**: Full detail (ensures nothing is forgotten later)

**Q13.3**: Design review process?

After each design is written:
- [ ] No review, just mark complete
- [ ] Self-review against template
- [ ] User reviews every design
- [ ] User reviews only complex/critical designs
- [ ] Batch review (review 5-10 at once)

**Q13.4**: Should designs include implementation estimates?

For each design, should we estimate:
- [ ] Lines of code?
- [ ] Time to implement?
- [ ] Complexity level (simple/medium/complex)?
- [ ] Dependencies/blockers?
- [ ] None (estimates come later in coding roadmap)

---

### Q14: Design Dependencies & Ordering

**Q14.1**: Should we map design dependencies?

Some designs depend on others (e.g., Music module depends on Metadata service). Should we:
- [ ] Create dependency graph
- [ ] Just note dependencies in each doc
- [ ] No formal tracking
- [ ] Track in separate dependency map file

**Q14.2**: Design writing order?

- [ ] Alphabetical (easy to track)
- [ ] By priority (MVP first)
- [ ] By dependency (foundations first)
- [ ] By category (all features, then services, etc.)
- [ ] Opportunistic (whatever is easiest/most interesting)

**Recommendation**: By dependency (build from foundation up)

---

### Q15: External Research Requirements

**Q15.1**: For designs needing external research, what sources are acceptable?

- [ ] Only already-fetched sources (docs/dev/sources/)
- [ ] Fetch new sources as needed (add to SOURCES.yaml)
- [ ] One-time web research (not added to sources)
- [ ] No external research (only internal knowledge)

**Q15.2**: Should all external research be documented?

If we research external systems for design:
- [ ] Add all sources to SOURCES.yaml
- [ ] Add only frequently-referenced sources
- [ ] Just link in design doc
- [ ] No requirement to track

---

### Q16: Design Validation Criteria

**Q16.1**: What makes a design "complete"?

Checklist for marking design as ✅:
- [ ] All template sections filled
- [ ] Architecture diagram included
- [ ] Database schema defined
- [ ] API endpoints listed
- [ ] Testing strategy outlined
- [ ] Cross-references added
- [ ] Linting passes
- [ ] User approval

Which of these are required?

**Q16.2**: How to handle uncertain designs?

If we're unsure about architecture/approach:
- [ ] Document multiple options, decide later
- [ ] Pick best option now, refine during implementation
- [ ] Ask user for decision
- [ ] Mark as partial (🟡) until resolved

---

### Q17: Timeline & Progress Tracking

**Q17.1**: Estimated time for design completion?

Based on gap analysis, how long should Phase 1 take:
- [ ] 1-2 weeks (few gaps, fast writing)
- [ ] 3-4 weeks (moderate gaps, thorough design)
- [ ] 1-2 months (many gaps, complex designs)
- [ ] Ongoing (design continuously, no deadline)

**Q17.2**: Progress tracking frequency?

Update design progress tracker:
- [ ] After each design completed
- [ ] Daily summary
- [ ] Weekly summary
- [ ] When category is complete
- [ ] No formal tracking

**Q17.3**: Milestones within Phase 1?

Should we break Phase 1 into sub-milestones:
- [ ] Week 1: Gap analysis + scaffolding
- [ ] Week 2: Tier 1 (MVP) designs
- [ ] Week 3: Tier 2 designs
- [ ] Week 4: Tier 3 designs + validation
- [ ] No sub-milestones (just "design completion")

---

### Q18: Integration with Existing Workflow

**Q18.1**: Should design work use the doc pipeline?

As we write designs:
- [ ] Run doc-pipeline after each design
- [ ] Run pipeline at end of each day
- [ ] Run pipeline at end of Phase 1
- [ ] Don't run until all designs complete

**Q18.2**: Should designs be committed incrementally?

- [ ] Commit after each design completed
- [ ] Commit after each category completed
- [ ] Commit at phase milestones
- [ ] Single commit at end of Phase 1

**Q18.3**: Feature branch or develop?

- [ ] Work on feature branch (phase-1-design-completion)
- [ ] Work directly on develop
- [ ] Work in .analysis/ only (no commits yet)

---

### Q19: Design Format & Style

**Q19.1**: Should all designs follow exact template?

- [ ] Yes, every design uses exact same template
- [ ] No, adapt template per type (service vs feature vs integration)
- [ ] Base template + category-specific additions
- [ ] Flexible (whatever makes sense for that design)

**Q19.2**: Diagram requirements?

For architecture diagrams:
- [ ] Required for all designs
- [ ] Required for complex designs only
- [ ] Optional but recommended
- [ ] Not needed (text description enough)

**Q19.3**: Code examples in design docs?

Should designs include:
- [ ] Pseudocode examples
- [ ] Go interface definitions
- [ ] Example API requests/responses
- [ ] No code (just descriptions)

**Recommendation**: Go interfaces + API examples (concrete specs)

---

### Q20: Open Questions & Decisions

**Q20.1**: How to handle design decisions that need user input?

When writing design, if architectural choice is unclear:
- [ ] Document options, ask user immediately
- [ ] Pick best option, note as "pending review"
- [ ] Add to "open questions" section, batch ask later
- [ ] Use general-purpose agent to research and recommend

**Q20.2**: Should we create a central decision log?

Track all major architectural decisions:
- [ ] Yes, create DESIGN_DECISIONS.md
- [ ] Note in each design doc
- [ ] Add to SOURCE_OF_TRUTH.md
- [ ] No formal tracking

---

## PRIORITY QUESTIONS (Must answer before starting Phase 1)

These are the MOST CRITICAL questions to answer before beginning gap analysis:

### ✅ P0 - ANSWERED:

1. **Q11.1**: What scope for gap analysis?
   - **ANSWER**: ✅ Comprehensive (features + services + integrations + operations)

2. **Q12.1**: What level of detail for scaffolds?
   - **ANSWER**: ✅ Structured (template with TODOs in each section)

3. **Q13.1**: Who writes designs?
   - **ANSWER**: ✅ Collaborative (I draft, you review/refine)

4. **Q13.2**: How much detail per design?
   - **ANSWER**: ✅ Full detail (all sections from template) - targeting 99% perfection

5. **Q17.1**: Timeline?
   - **ANSWER**: ✅ No fixed timeline - complete when ready

### ✅ P1 - ANSWERED (Proceeding with recommendations):

6. **Q11.3**: Check existing docs for completeness?
   - **ANSWER**: Yes, flag incomplete/partial docs

7. **Q12.2**: Include research links in scaffolds?
   - **ANSWER**: Yes, add research links

8. **Q14.1**: Map design dependencies?
   - **ANSWER**: Yes, create dependency graph

9. **Q18.3**: Feature branch or develop?
   - **ANSWER**: Feature branch (phase-1-design-completion)

### 🟢 P2 - Can answer during design writing:

10. **Q13.3**: Design review process?
11. **Q13.4**: Include implementation estimates?
12. **Q15.1**: External research sources?
13. **Q16.2**: How to handle uncertain designs?

---

## Target: 99% Perfection

Per user requirement: "if possible we want about 99% perfection before we start meaning: only things we will likely not have finished by then are the actuall final branding files (placeholder graphics for now)"

This means:
- Full detail in all design sections
- Complete architecture diagrams
- Complete database schemas
- Complete API endpoint definitions
- Complete integration specifications
- Complete testing strategies
- Only placeholder graphics acceptable as incomplete

---

## Next Steps - STARTING NOW

1. ✅ P0 questions answered (using recommendations)
2. → Start gap analysis (create `.analysis/09_DESIGN_GAPS.md`)
3. → Create scaffolds for ALL missing docs
4. → Begin design writing (in dependency order)
5. → Track progress in `.analysis/10_DESIGN_PROGRESS.md`

---

**STATUS**: ✅ READY TO START PHASE 1 - Beginning gap analysis now

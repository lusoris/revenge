# File Renaming Strategy

> Proposed numbered file convention for design documentation

**Status**: PROPOSED - Awaiting approval

---

## Current Structure Analysis

### Root Design Folder
```
docs/dev/design/
├── 00_SOURCE_OF_TRUTH.md      # Canonical reference
├── 01_DESIGN_DOC_TEMPLATE.md  # Template for new docs
├── 02_QUESTIONS_TO_DISCUSS.md # Open questions
└── [subfolders]/
```

### Subfolders (13 total)
- `architecture/` (6 files) - Core architecture docs
- `technical/` (7 files) - Tech stack, API, configuration
- `services/` (16 files) - Service design docs
- `features/` (25+ files) - Feature specifications
- `integrations/` (40+ files) - External integrations
- `operations/` (6 files) - DevOps and operations
- `planning/` (3 files) - Versioning, roadmap
- `research/` (3 files) - UX research
- `.archive/` - Obsolete documents

---

## Proposed Strategy: Selective Numbering

### Option A: Number Core Docs Only (Recommended)

Only add numeric prefixes to **top-level** and **architecture** docs to indicate reading order. Subfolders keep current naming with INDEX.md for navigation.

**Root (reading order):**
```
00_SOURCE_OF_TRUTH.md       # Read first - the canonical reference
01_DESIGN_DOC_TEMPLATE.md   # How to write new docs
02_QUESTIONS_TO_DISCUSS.md  # Open architectural questions
```

**Architecture (reading order):**
```
00_INDEX.md                        # Navigation
01_ARCHITECTURE.md                 # Main architecture overview
02_DESIGN_PRINCIPLES.md            # Core principles
03_METADATA_SYSTEM.md              # Metadata architecture
04_PLAYER_ARCHITECTURE.md          # Player design
05_PLUGIN_ARCHITECTURE_DECISION.md # Plugin decisions
```

**Other folders:** Keep current naming (already have INDEX.md files)

### Option B: Number Everything

Add numeric prefixes to all files based on a hierarchy:
```
1.0_ARCHITECTURE.md
1.1_DESIGN_PRINCIPLES.md
2.0_TECH_STACK.md
2.1_API.md
```

**Pros:** Clear reading order everywhere
**Cons:** Complex to maintain, harder to add new docs

### Option C: No Renaming

Keep current naming. Use INDEX.md files for navigation.

**Pros:** No changes needed, URLs stay stable
**Cons:** No clear reading order for newcomers

---

## Recommendation

**Option A** provides the best balance:
- Clear entry point (00_SOURCE_OF_TRUTH.md)
- Logical reading order for architecture
- Minimal changes to existing links
- Easy to maintain

### Files to Rename (Option A)

| Current | Proposed |
|---------|----------|
| `00_SOURCE_OF_TRUTH.md` | `00_SOURCE_OF_TRUTH.md` |
| `01_DESIGN_DOC_TEMPLATE.md` | `01_DESIGN_DOC_TEMPLATE.md` |
| `02_QUESTIONS_TO_DISCUSS.md` | `02_QUESTIONS_TO_DISCUSS.md` |
| `architecture/00_INDEX.md` | `architecture/00_INDEX.md` |
| `architecture/01_ARCHITECTURE.md` | `architecture/01_ARCHITECTURE.md` |
| `architecture/02_DESIGN_PRINCIPLES.md` | `architecture/02_DESIGN_PRINCIPLES.md` |
| `architecture/03_METADATA_SYSTEM.md` | `architecture/03_METADATA_SYSTEM.md` |
| `architecture/04_PLAYER_ARCHITECTURE.md` | `architecture/04_PLAYER_ARCHITECTURE.md` |
| `architecture/05_PLUGIN_ARCHITECTURE_DECISION.md` | `architecture/05_PLUGIN_ARCHITECTURE_DECISION.md` |

### Link Updates Required

After renaming, update all cross-references:
- Links in INDEX.md files
- Links in all design docs referencing core docs
- Any external references

---

## Decision Required

Choose one:
- [ ] **Option A**: Number core docs only (recommended)
- [ ] **Option B**: Number everything
- [ ] **Option C**: Keep current naming

---

## Execution Plan (if Option A approved)

1. Create backup of current state (git commit)
2. Rename root files with git mv
3. Rename architecture files with git mv
4. Update all cross-references with grep + sed
5. Verify no broken links
6. Commit changes

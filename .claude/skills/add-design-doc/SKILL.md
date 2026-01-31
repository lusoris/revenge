---
name: add-design-doc
description: Create a new design document following project conventions
argument-hint: <category/name>
disable-model-invocation: true
---

# Add Design Document

Create a new design document in the specified category following project conventions.

## Usage

```
/add-design-doc services/NOTIFICATIONS
/add-design-doc features/video/TRAILER_SUPPORT
/add-design-doc integrations/metadata/comics/MANGA_API
```

## Arguments

- `$ARGUMENTS`: The path for the new document (e.g., `services/AUTH` or `features/adult/GALLERIES`)

## Task

Create a new design document at `docs/dev/design/$ARGUMENTS.md` with:

1. **Standard structure** following the template:
   - Title (# heading)
   - Description blockquote
   - Status table with all 7 dimensions
   - Overview section
   - Source of Truth reference
   - Source breadcrumb placeholder

2. **Proper status table**:
   ```markdown
   | Dimension | Status | Notes |
   |-----------|--------|-------|
   | Design | 🔴 | |
   | Sources | 🔴 | |
   | Instructions | 🔴 | |
   | Code | 🔴 | |
   | Linting | 🔴 | |
   | Unit Testing | 🔴 | |
   | Integration Testing | 🔴 | |
   ```

3. **Category-specific sections**:
   - `services/`: Include Module path, Dependencies
   - `integrations/`: Include API reference, Authentication
   - `features/`: Include Implementation, Design

4. **Update the INDEX.md** in the parent directory if needed

5. **Add to 03_DESIGN_DOCS_STATUS.md** if it exists

## Template Reference

See `docs/dev/design/01_DESIGN_DOC_TEMPLATE.md` for the full template.

## After Creation

Run the doc pipeline to update indexes:
```bash
./scripts/doc-pipeline.sh --apply
```

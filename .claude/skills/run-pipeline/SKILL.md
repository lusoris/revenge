---
name: run-pipeline
description: Run documentation pipelines (source or doc)
argument-hint: <source|doc|all> [--apply]
disable-model-invocation: true
allowed-tools: Bash(./scripts/*.sh *)
---

# Run Pipeline

Run the documentation pipelines to update indexes, breadcrumbs, and cross-references.

## Usage

```
/run-pipeline doc                  # Dry run doc pipeline
/run-pipeline doc --apply          # Apply doc pipeline changes
/run-pipeline source               # Dry run source pipeline
/run-pipeline source --apply       # Apply source pipeline changes
/run-pipeline all --apply          # Run both pipelines
```

## Arguments

- `$0`: Pipeline to run (`source`, `doc`, or `all`)
- `$1`: Optional `--apply` to write changes (default: dry-run)

## Pipelines

### Source Pipeline (`source`)

Updates external documentation sources:

1. **01-fetch.py**: Fetch external documentation
2. **02-index.py**: Generate SOURCES.md index
3. **03-breadcrumbs.py**: Add source breadcrumbs to design docs

```bash
./scripts/source-pipeline.sh [--apply]
```

### Doc Pipeline (`doc`)

Updates internal documentation structure:

1. **01-indexes.py**: Generate INDEX.md files
2. **02-breadcrumbs.py**: Add design breadcrumbs
3. **03-status.py**: Sync status tables
4. **04-validate.py**: Validate document structure
5. **05-fix.py**: Fix broken links
6. **06-meta.py**: Generate DESIGN_INDEX.md

```bash
./scripts/doc-pipeline.sh [--apply]
```

## Options

| Option | Description |
|--------|-------------|
| `--apply` | Actually write changes (default: dry-run) |
| `--step N` | Run only step N |
| `--force` | Force update unchanged sources (source pipeline only) |
| `--validate` | Only run validation (doc pipeline only) |

## Examples

Check what would change:
```bash
./scripts/doc-pipeline.sh
```

Apply all changes:
```bash
./scripts/doc-pipeline.sh --apply
```

Run only validation:
```bash
./scripts/doc-pipeline.sh --validate
```

Refresh all sources:
```bash
./scripts/source-pipeline.sh --force --apply
```

## Safety

- Both pipelines default to **dry-run** mode
- They check for uncommitted changes before applying
- Use git to review and revert changes if needed

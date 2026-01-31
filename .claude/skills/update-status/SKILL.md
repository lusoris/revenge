---
name: update-status
description: Update status tables in design documents
argument-hint: [category]
disable-model-invocation: true
---

# Update Status

Update status tables in design documents, either for a specific document or category.

## Usage

```
/update-status                          # Check all status tables
/update-status services                  # Check services category
/update-status services/AUTH.md Design ✅  # Update specific dimension
```

## Arguments

- `$0`: Category or file path (optional)
- `$1`: Dimension name (optional, e.g., Design, Sources, Code)
- `$2`: New status (optional, e.g., ✅, 🟡, 🔴)

## Task

### If no arguments: Audit all status tables

1. Run the status sync script in check mode:
   ```bash
   python3 scripts/doc-pipeline/03-status.py --verbose
   ```

2. Report:
   - Documents without status tables
   - Documents with invalid status values
   - Status coverage percentage

### If category specified: Check category

1. Run for specific category:
   ```bash
   python3 scripts/doc-pipeline/03-status.py --category $0
   ```

### If file and status specified: Update status

1. Read the specified document
2. Find the status table
3. Update the specified dimension with the new status
4. Write the changes

## Valid Statuses

| Status | Meaning |
|--------|---------|
| ✅ | Complete |
| 🟡 | Partial / In Progress |
| 🔴 | Not Started |
| ⚪ | N/A |

## Valid Dimensions

- Design
- Sources
- Instructions
- Code
- Linting
- Unit Testing
- Integration Testing

## Example

Update AUTH.md Design status to complete:
```
/update-status services/AUTH.md Design ✅
```

This will change the status table row from:
```
| Design | 🟡 | In progress |
```
to:
```
| Design | ✅ | Complete |
```

# v0.2.0 Bugs & Issues

**Version**: v0.2.0 - Core Backend Services
**Last Updated**: 2026-02-02

## Active Bugs

### BUG-001: Terminal heredoc corruption with file creation tools
**Severity**: Medium
**Reported**: 2026-02-02
**Status**: Workaround implemented

**Description**: When using `cat << 'EOF'` or `dd << 'EOF'` in terminal, the shell's prompt hook executes `ls` which injects directory listings into the heredoc content, corrupting the files being created.

**Impact**: Files created via heredoc in terminal get corrupted with ls output interspersed in the code, causing syntax errors.

**Workaround**: Use `replace_string_in_file` with empty oldString on touched files, or use Python/base64 for file creation.

**Root Cause**: Terminal prompt configuration executes commands on prompt display.

**Resolution**: Use alternative file creation methods. This is an environment issue, not a code bug.

---

## Resolved Bugs

No bugs resolved yet.

## Known Issues

No known issues yet.

## Technical Debt

No technical debt yet.

## Future Considerations

No items yet.

---

## Bug Template

When adding bugs, use this format:

```markdown
### [BUG-XXX] Short Title

**Severity**: Critical / High / Medium / Low
**Status**: Open / In Progress / Resolved / Won't Fix
**Component**: Auth / User / Session / etc.
**Reported**: YYYY-MM-DD
**Resolved**: YYYY-MM-DD (if resolved)

**Description**:
Brief description of the bug.

**Steps to Reproduce**:
1. Step 1
2. Step 2
3. Step 3

**Expected Behavior**:
What should happen.

**Actual Behavior**:
What actually happens.

**Workaround** (if any):
Temporary solution.

**Fix** (if resolved):
How it was fixed.

**Related**:
- Issue #XXX
- Commit: abc123
```

---

## Issue Categories

- **Security**: Security vulnerabilities
- **Performance**: Performance issues
- **Functionality**: Feature not working as designed
- **UX**: User experience issues
- **API**: API contract violations
- **Database**: Database-related issues
- **Cache**: Caching issues
- **Jobs**: Background job issues

---

## Updates Log

| Date | Update |
|------|--------|
| 2026-02-02 | Created initial bugs file |

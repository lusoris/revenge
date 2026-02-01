# Bugfixes - v0.1.0 Implementation

> Issues encountered and their fixes (for future test creation)

---

## Issues

### [ISSUE-001] TestDefault fails - Database.URL has no default value
**Problem**: TestDefault test failure in CI: "Database.URL should have a default value"
**Cause**: Both `Defaults()` function (config.go:170) and `Default()` function (module.go:37) returned empty string for Database.URL
**Fix**: Changed both locations to use placeholder: `"postgres://revenge:changeme@localhost:5432/revenge?sslmode=disable"`
**Test Hint**: Test that default config values match expected placeholder values
**Files Changed**:
- [internal/config/config.go:170](internal/config/config.go#L170)
- [internal/config/module.go:37](internal/config/module.go#L37)

---

# Testing Decisions - Final

**Datum**: 2026-02-04
**Status**: ✅ Alle Entscheidungen getroffen

---

## ✅ Entscheidungen

### 1. Mock-Strategie
**Entscheidung**: **mockery** (automatische Mock-Generierung)
- Moderner als gomock
- Type-safe
- Einfach zu verwenden

### 2. Integration Tests
**Entscheidung**: **testcontainers-go** (Docker Container)
- Isoliert
- Reale PostgreSQL-Instanz
- Bereits in testutil implementiert

### 3. Test-Tiefe
**Entscheidung**: **Exhaustive** (alle Edge Cases)
- Happy paths
- Error cases
- Edge cases
- Boundary conditions
- **Hinweis**: Zeitaufwendiger, aber maximale Qualität

### 4. Test Helpers
**Entscheidung**: **Existing testutil Package verwenden**
- `internal/testutil/` existiert bereits
- Fast parallel DB testing (embedded-postgres)
- Testcontainers integration
- Fixtures, assertions, database helpers

### 5. Execution
**Entscheidung**: **Sequenziell** (ein Package nach dem anderen)
- Klarer Progress
- Einfacher zu debuggen

### 6. Test Style
**Entscheidung**: **Immer table-driven**
- Go Convention
- Konsistent
- Einfach zu erweitern

### 7. Bug Handling
**Entscheidung**: **Sofort fixen + dokumentieren**
- Test schreiben → Bug finden → sofort fixen
- Bug in `.workingdir/BUG_XX_*.md` dokumentieren
- Keine technische Schulden

### 8. Coverage Tracking
**Entscheidung**: **Nach jeder Phase**
- Report nach Phase 1, 2, 3, 4, 5
- Weniger Overhead als per-Package
- Klare Milestones

### 9. Package-Priorität
**Entscheidung**: ✅ **Approved**
- **CRITICAL**: Session, Auth, RBAC, User (>80%)
- **HIGH**: Movie, Library, Settings (>70%)
- **MEDIUM**: Activity, Search, Notification (>50%)
- **LOW**: Handlers, Integrations (>30%)

### 10. Start-Approach
**Entscheidung**: **mockery setup → Mock-Gen → Tests**
1. mockery installieren & konfigurieren (15min)
2. Alle Interfaces finden & Mocks generieren (15min)
3. Session Service Tests schreiben (2-3h)

---

## 🚀 Execution Plan

### Phase 0: Mock Setup (30min)
1. ✅ Install mockery: `go install github.com/vektra/mockery/v2@latest`
2. ✅ Create `.mockery.yaml` config
3. ✅ Identify all interfaces to mock:
   - Repository interfaces (User, Session, Auth, etc.)
   - Service interfaces
   - External dependencies (Cache, Queue, etc.)
4. ✅ Generate mocks: `mockery`
5. ✅ Verify mocks compile

### Phase 1: Core Services (12-16h)
**Target**: 80%+ coverage per package

1. **Session Service** (2-3h)
   - ValidateSession (happy + expired + invalid + edge cases)
   - CreateSession (happy + duplicate + validation errors)
   - RevokeSession (happy + not found + already revoked)
   - RevokeAllUserSessions (happy + no sessions + partial failures)
   - CleanupExpiredSessions (happy + empty + partial cleanup)

2. **Auth Service** (3-4h)
   - Login (happy + wrong password + user not found + account locked + MFA required)
   - Register (happy + duplicate email + duplicate username + validation errors)
   - VerifyPassword (happy + wrong + hash errors)
   - GenerateTokens (happy + signing errors)
   - RefreshToken (happy + expired + revoked + invalid)

3. **User Service** (2-3h)
   - Create, Get, Update, Delete (happy + errors)
   - GetByUsername, GetByEmail (happy + not found)
   - ChangePassword (happy + wrong old password + validation)
   - Edge cases (nil values, empty strings, etc.)

4. **RBAC Service** (3-4h)
   - Enforce (happy + denied + no policy + edge cases)
   - AssignRole, RemoveRole (happy + errors)
   - GetUserRoles (happy + no roles + multiple roles)
   - CreateRole, DeleteRole, UpdateRolePermissions (happy + errors + conflicts)
   - ListPermissions (happy + empty)

5. **Settings Service** (2-3h)
   - Get, Set server settings (happy + not found + validation)
   - Get, Set user settings (happy + not found + validation)
   - Cache invalidation tests

**Milestone**: Coverage-Report nach Phase 1

### Phase 2-5: Continue as planned

---

## 📝 Bug Documentation Template

When a test finds a bug:

```markdown
# Bug #XX: [Short Description]

**Date**: 2026-02-04
**Found During**: [Package] test writing
**Severity**: [CRITICAL/HIGH/MEDIUM/LOW]

## Description
[What is the bug?]

## Reproduction
[Steps to reproduce or test case that triggered it]

## Root Cause
[Why does this happen?]

## Fix
[How was it fixed?]

## Files Changed
- [List of files]

## Tests Added
- [List of tests that now cover this case]

---

**Status**: ✅ FIXED
```

---

## 📊 Coverage Expectations

**Current**: 4.13%
**After Phase 0**: 4.13% (no change, just setup)
**After Phase 1**: ~25-30% (Core Services)
**After Phase 2**: ~40-45% (Content Services)
**After Phase 3**: ~55-60% (Integrations)
**After Phase 4**: ~70-75% (Handlers)
**After Phase 5**: ~80-85% ✅ (Infrastructure)

---

**Status**: Ready to start with mockery setup
**Next Action**: Install & configure mockery

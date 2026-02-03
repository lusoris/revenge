# Design Documentation Audit - 2026-02-03

**Created**: 2026-02-03
**Purpose**: Identify deviations between implementation and design documentation
**Status**: 🟡 In Progress - Awaiting User Decisions

---

## Executive Summary

Audit identified **3 categories of deviations** between implementation and design documentation:

1. **Missing Design Document**: User Settings/Preferences system fully implemented but no design doc
2. **Token Generation Changes**: Session token generation implemented inline, not as separate service + token format uses SHA-256 not documented
3. **Package Version Mismatches**: Several packages in go.mod differ from SOURCE_OF_TRUTH.md versions

---

## 1. Missing Design Documentation

### User Settings & Preferences System

**Status**: ✅ **FULLY IMPLEMENTED** but 🔴 **NO DESIGN DOC**

**Evidence**:
- Database migrations exist:
  - `000005_create_user_settings_table.up.sql` (shared schema)
  - `000006_create_user_preferences_table.up.sql` (shared schema)
  
- Tables implemented:
  ```sql
  shared.user_settings (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES shared.users(id),
    key TEXT NOT NULL,
    value JSONB NOT NULL,
    data_type TEXT NOT NULL,
    category TEXT,
    description TEXT,
    is_public BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
  )
  
  shared.user_preferences (
    user_id UUID PRIMARY KEY REFERENCES shared.users(id),
    email_notifications JSONB DEFAULT '{}',
    push_notifications JSONB DEFAULT '{}',
    digest_notifications JSONB DEFAULT '{}',
    profile_visibility TEXT DEFAULT 'friends',
    activity_privacy TEXT DEFAULT 'friends',
    library_privacy TEXT DEFAULT 'private',
    playback_settings JSONB DEFAULT '{}',
    subtitle_settings JSONB DEFAULT '{}',
    theme_settings JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
  )
  ```

**Current Design Doc** (`docs/dev/design/services/SETTINGS.md`):
- Exists but labeled as "Server settings persistence and retrieval"
- **Does NOT document user_settings or user_preferences tables**
- Only references server-level settings
- Status shows "✅ Complete" but implementation differs

**Gap Analysis**:
- Design doc focuses on SERVER settings (server_settings table)
- USER settings and USER preferences are completely undocumented
- Two distinct systems exist:
  1. **Server Settings**: Global configuration (config.yaml, environment variables)
  2. **User Settings**: Per-user key-value configuration
  3. **User Preferences**: Per-user UI/notification/privacy preferences

**Required Action**:
- Option A: Create new `docs/dev/design/services/USER_SETTINGS.md`
- Option B: Expand `SETTINGS.md` to cover both server and user settings
- Option C: Merge into `USER.md` as user configuration subsystem

**Questions for User**:
1. Should user_settings and user_preferences be documented separately or as part of USER.md?
2. Are user_settings (flexible key-value) and user_preferences (structured) intentionally separate?
3. What's the intended use case differentiation between settings and preferences?

---

## 2. Token Generation Implementation Deviation

### Issue: Token Hashing Helper Inline vs Separate Service

**Documented Approach** (`docs/dev/design/services/SESSION.md`):
```go
// Session service dependencies (from design doc)
**Go Packages**:
- `crypto/rand` - Token generation
- `crypto/sha256` - Token hashing
```

**Actual Implementation** (`internal/service/session/service.go`):
```go
// Inline helper methods (NOT a separate service/package)
func (s *Service) generateToken() (string, string, error) {
    token := make([]byte, s.tokenLength)
    if _, err := rand.Read(token); err != nil {
        return "", "", err
    }
    
    tokenStr := hex.EncodeToString(token)
    tokenHash := s.hashToken(tokenStr)
    
    return tokenStr, tokenHash, nil
}

func (s *Service) hashToken(token string) string {
    hash := sha256.Sum256([]byte(token))
    return hex.EncodeToString(hash[:])
}
```

**Analysis**:
- ✅ Tokens are generated using `crypto/rand` (secure)
- ✅ Tokens are hashed using SHA-256 (as documented)
- ⚠️ Implementation is inline, not a separate "HashToken service"
- ⚠️ Token format not documented in SESSION.md

**Token Format Details**:
- **Token Generation**: 32 bytes random → hex encoded = 64-character hex string
- **Token Storage**: SHA-256 hash of token (64-character hex string)
- **Database Column**: `token_hash TEXT NOT NULL`
- **Security**: Tokens never stored in plaintext, only hashes

**SESSION.md Documentation Gap**:
```yaml
# Missing from SESSION.md:
session:
  token_length: 32              # bytes (results in 64 hex chars)
  token_hash_algorithm: sha256  # NOT documented
  token_format: hex             # NOT documented
```

**Current Config** (`config/config.yaml`):
```yaml
session:
  token_length: 32           # ✅ Documented
  expiry: 720h               # ✅ Documented
  # Missing: hash algorithm, format
```

**Deviation Assessment**:
- **Severity**: LOW - Implementation is correct and secure
- **Issue**: Documentation doesn't specify:
  - Token hash algorithm (SHA-256)
  - Token encoding format (hex)
  - Final token length (64 chars, not 32)
  - Storage format (hash only, not plaintext)

**Required Action**:
Update `docs/dev/design/services/SESSION.md`:
1. Add token format specification section
2. Document SHA-256 hashing
3. Clarify token_length is pre-encoding bytes
4. Add example token and hash

**Questions for User**:
1. Should token generation be extracted to a separate `internal/crypto/token.go` package?
2. Is inline implementation acceptable for this use case?
3. Should we support other hash algorithms (SHA-512, BLAKE2b) or is SHA-256 sufficient?

---

## 3. Package Version Mismatches

### SOURCE_OF_TRUTH.md vs go.mod Discrepancies

**Mismatched Packages**:

| Package | SOURCE_OF_TRUTH.md | go.mod | Severity | Notes |
|---------|-------------------|--------|----------|-------|
| `github.com/maypok86/otter` | v1.2.4 | v2.3.0 | **HIGH** | Major version upgrade (v1→v2), API changes likely |
| `github.com/redis/rueidis` | v1.0.49 | v1.0.71 | LOW | Minor patch, backward compatible |
| `github.com/ogen-go/ogen` | v1.18.0 | v1.18.0 | ✅ MATCH | - |
| `github.com/jackc/pgx/v5` | v5.8.0 | v5.8.0 | ✅ MATCH | - |
| `github.com/riverqueue/river` | v0.26.0 | v0.26.0 | ✅ MATCH | - |
| `go.uber.org/fx` | v1.24.0 | Not in go.mod direct | ⚠️ | Likely indirect dependency |
| `github.com/knadh/koanf/v2` | v2.3.0 | v2.3.2 | LOW | Patch update |

**Critical Issue: otter v1.2.4 → v2.3.0**

SOURCE_OF_TRUTH.md says:
```markdown
| `github.com/maypok86/otter` | v1.2.4 | In-memory cache | W-TinyLFU, faster than Ristretto |
```

go.mod says:
```go
github.com/maypok86/otter/v2 v2.3.0
```

**Impact**:
- Otter v2 is a MAJOR version upgrade with breaking API changes
- Import path changed: `github.com/maypok86/otter` → `github.com/maypok86/otter/v2`
- Code imports v2: `"github.com/maypok86/otter/v2"` (verified in cache package)
- Documentation is outdated

**Actual Implementation** (`internal/infra/cache/otter.go`):
```go
import (
    "github.com/maypok86/otter/v2"  // v2 API
)
```

**Action Required**: Update SOURCE_OF_TRUTH.md to reflect otter v2.3.0

---

### Missing Packages in SOURCE_OF_TRUTH.md

**Packages in go.mod NOT listed in SOURCE_OF_TRUTH.md**:

| Package | Version | Purpose | Used In |
|---------|---------|---------|---------|
| `github.com/alexedwards/argon2id` | v1.0.0 | Password hashing | `internal/crypto/password.go` |
| `golang.org/x/crypto` | (indirect) | Cryptography | Via argon2id |
| `github.com/coreos/go-oidc/v3` | v3.17.0 | OIDC client | OIDC service |
| `github.com/fergusstrange/embedded-postgres` | v1.33.0 | Testing | `internal/testutil/testdb.go` |

**Critical Missing: argon2id**

Password hashing implementation uses `github.com/alexedwards/argon2id` but SOURCE_OF_TRUTH.md says:
```markdown
| `golang.org/x/crypto` | v0.47.0 | Cryptography | Argon2, bcrypt, AES-GCM |
```

**Reality**:
- We use `github.com/alexedwards/argon2id` (wrapper library)
- NOT direct `golang.org/x/crypto/argon2`
- Benefits: Simplified API, secure defaults, timing attack protection

**Action Required**: Add argon2id to SOURCE_OF_TRUTH.md Security section

---

## 4. Additional Findings from .workingdir Scan

### Bug #29: Password Hash Migration

**File**: `.workingdir/BUG_29_PASSWORD_HASH_MIGRATION.md`

**Deviation Identified**:
- Database contains bcrypt hashes (`$2a$12$...`)
- Code expects argon2id hashes (`$argon2id$v=19$...`)
- No migration documented in design docs

**Resolution Status**: ✅ RESOLVED
- Implemented hybrid verifier in `internal/crypto/password.go`
- Supports both bcrypt (legacy) and argon2id (current)
- Automatic migration on next login

**Design Doc Impact**: None (temporary migration code, will be removed)

---

### Token Hash Storage Format

**Observation**: Session tokens use SHA-256 hashing

**Design Doc Gap**: `SESSION.md` doesn't document:
- Hash algorithm (SHA-256)
- Hash format (hex string, 64 chars)
- Token format (hex string, 64 chars)
- Why SHA-256 vs bcrypt/argon2id (speed, not password hashing)

**Recommendation**: Add "Token Security Model" section to SESSION.md

---

## Summary of Required Actions

### High Priority (Blocks v1.0)

1. **Create User Settings Design Doc** (NEW DOCUMENT)
   - Document user_settings table (flexible key-value)
   - Document user_preferences table (structured preferences)
   - Explain separation of concerns
   - API endpoints, configuration, use cases

2. **Update SOURCE_OF_TRUTH.md Packages**
   - Change otter v1.2.4 → v2.3.0 ✅ CRITICAL
   - Add argon2id v1.0.0 (currently missing)
   - Update rueidis v1.0.49 → v1.0.71
   - Update koanf v2.3.0 → v2.3.2

3. **Update SESSION.md**
   - Add "Token Security Model" section
   - Document SHA-256 hashing
   - Document token format (hex, 64 chars)
   - Clarify token_length config (bytes, not final length)

### Medium Priority (Documentation Quality)

4. **Expand SETTINGS.md**
   - Clarify scope (server-level only)
   - Add references to USER_SETTINGS.md
   - Document server_settings table schema

5. **Review All Design Docs**
   - Check for other version mismatches
   - Verify all migrations have corresponding design docs
   - Cross-reference SOURCE_OF_TRUTH.md

### Low Priority (Nice to Have)

6. **Consider Token Helper Extraction**
   - Extract `generateToken()` and `hashToken()` to `internal/crypto/token.go`?
   - Reusable across auth_tokens, password_reset_tokens, email_verification_tokens
   - Consistent token format across services

---

## Questions for User Decision

### 1. User Settings Documentation Strategy

**Question**: How should we document user_settings and user_preferences?

**Option A**: Create new `docs/dev/design/services/USER_SETTINGS.md`
- ✅ Focused, dedicated document
- ✅ Clear separation from server settings
- ❌ Adds another design doc

**Option B**: Expand existing `docs/dev/design/services/SETTINGS.md`
- ✅ All settings in one place
- ❌ Mixes server and user concerns
- ❌ Document becomes large

**Option C**: Merge into `docs/dev/design/services/USER.md`
- ✅ User-related features together
- ❌ USER.md already comprehensive
- ❌ Settings buried in larger doc

**Recommendation**: **Option A** - Create USER_SETTINGS.md

---

### 2. user_settings vs user_preferences Distinction

**Question**: Why two separate tables? Should we consolidate?

**Current Implementation**:
- `user_settings`: Flexible JSONB key-value store (any setting)
- `user_preferences`: Structured columns (email_notifications, theme_settings, etc.)

**Possible Reasons**:
1. **Performance**: Structured columns faster than JSONB queries
2. **Type Safety**: Known columns vs dynamic keys
3. **Use Cases**: 
   - Settings = custom app-specific config
   - Preferences = UI/UX/notification settings

**Questions**:
- Is this intentional separation or accidental duplication?
- Should we consolidate into one table?
- Document the distinction clearly?

---

### 3. Token Helper Service Extraction

**Question**: Should token generation be a separate service/package?

**Current**: Inline methods in `internal/service/session/service.go`

**Proposed**: Extract to `internal/crypto/token.go`
```go
package crypto

type TokenService interface {
    GenerateToken(length int) (token, hash string, err error)
    HashToken(token string) string
    ValidateToken(token, hash string) bool
}
```

**Benefits**:
- ✅ Reusable across auth_tokens, password_reset, email_verification
- ✅ Consistent token format
- ✅ Testable in isolation
- ✅ Single source of truth for token security

**Drawbacks**:
- ❌ Over-engineering for simple SHA-256 hash
- ❌ Adds DI complexity
- ❌ Session service becomes more coupled

**Recommendation**: **KEEP INLINE** (simple, localized, not worth extraction)

---

### 4. Package Version Update Strategy

**Question**: Should we update packages immediately or wait for Dependabot?

**Identified Updates**:
- otter v2.3.0 (CRITICAL - already in go.mod, just fix docs)
- rueidis v1.0.71 (LOW - patch update)
- koanf v2.3.2 (LOW - patch update)

**Recommendation**:
1. Update SOURCE_OF_TRUTH.md to match go.mod (documentation fix, not code change)
2. Let Dependabot handle future updates
3. Review Dependabot PRs before merging

---

## Next Steps

**Immediate Actions** (Before MFA Implementation):

1. **User Decision Required**: Choose User Settings documentation strategy (Option A/B/C)
2. **Update SOURCE_OF_TRUTH.md**: Fix package versions (otter, add argon2id)
3. **Create/Update Design Docs**: Based on user decision
4. **Commit Documentation**: Sync design docs with implementation

**After Documentation Sync**:
5. Resume MFA implementation (Phase 1: Migrations)

---

## Appendix: Scan Results

### Design Docs Found (20 total)
```
✅ ACTIVITY.md      ✅ ANALYTICS.md     ✅ AUTH.md         ✅ APIKEYS.md
✅ EPG.md           ✅ OIDC.md          ✅ MFA.md          ✅ USER.md
✅ TRANSCODING.md   ✅ SETTINGS.md      ✅ SESSION.md      ✅ SEARCH.md
✅ RBAC.md          ✅ NOTIFICATION.md  ✅ METADATA.md     ✅ LIBRARY.md
✅ INDEX.md         ✅ HTTP_CLIENT.md   ✅ GRANTS.md       ✅ FINGERPRINT.md
```

### Migrations Found (15 pairs)
```
000001 schemas          000002 users            000003 sessions
000004 server_settings  000005 user_settings    000006 user_preferences
000007 user_avatars     000008 auth_tokens      000009 password_reset
000010 email_verify     000011 casbin_rules     000012 api_keys
000013 oidc_tables      000014 activity_log     000015 library_tables
```

**Gap**: user_settings (000005) and user_preferences (000006) have no design doc

### .workingdir Files Scanned (50 files)
- No additional deviations found beyond those documented above
- Bug #29 (password migration) already resolved
- Security scan results show 0 issues (all fixed)

---

**Audit Complete**: Awaiting user decisions on documentation strategy before proceeding.

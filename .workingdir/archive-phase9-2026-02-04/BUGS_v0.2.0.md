# v0.2.0 Bugs & Issues

**Version**: v0.2.0 - Core Backend Services
**Last Updated**: 2026-02-02

## Active Bugs

### BUG-002: Session Table Schema Mismatch
**Severity**: High
**Reported**: 2026-02-02
**Status**: ✅ Resolved
**Resolved**: 2026-02-02 (Commit 26)
**Component**: Session, Database

**Description**: The `shared.sessions` table schema in migration `000003_create_sessions_table.up.sql` doesn't match the schema that sqlc generated models from. This prevents session service implementation.

**Migration Schema** (what we want):
- `token_hash` TEXT (SHA256 of JWT access token)
- `refresh_token_hash` TEXT (SHA256 of refresh token)
- `last_activity_at` TIMESTAMPTZ
- `scopes` TEXT[] (OAuth2-style scopes)
- `revoked_at` TIMESTAMPTZ
- `revoke_reason` TEXT

**Current DB Schema** (what exists - from models.go):
- `refresh_token` TEXT (plain token, not hash)
- `access_token_hash` TEXT (exists)
- `last_used_at` TIMESTAMPTZ (not `last_activity_at`)
- No `scopes` field
- `is_active` BOOLEAN (instead of `revoked_at`)
- `revoked_at` pgtype.Timestamptz (exists)

**Impact**: Cannot implement session service without aligned schema. Session queries fail sqlc validation.

**Steps to Reproduce**:
1. Create sessions.sql queries matching SESSION.md design
2. Run `make sqlc`
3. Error: "column token_hash does not exist", "column scopes does not exist"

**Root Cause**: Migration 000003 was created/modified but database wasn't migrated, or models.go is from older schema

**Fix Required**:
1. Either update migration 000003 to match current DB schema
2. OR run migrations to update DB to match migration file
3. OR create new migration to ALTER table structure
4. Then regenerate sqlc models

**Workaround**: None - session service blocked

---

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

### BUG-002: Session Table Schema Mismatch ✅
**Resolved**: 2026-02-02 (Commit 26)
**Fix**:
1. Started PostgreSQL database via docker-compose
2. Created revenge_dev database
3. Ran all migrations (000001-000010)
4. Fixed sqlc.yaml schema path (internal/infra/database/migrations/shared/ → migrations/)
5. Regenerated sqlc models from actual database schema
6. Fixed repository_pg.go to match generated types (netip.Addr, time.Time, *string)
7. All 17 session queries working
8. Build successful with 0 lint issues

## Known Issues

No known issues yet.

## Technical Debt

### DEBT-001: HTTP Request Metadata Extraction (Auth Service)
**Priority**: Medium
**Component**: Auth
**Status**: Pending
**Reported**: 2026-02-02

**Description**: Auth handlers need to extract IP address, user agent, and device fingerprint from HTTP requests for security/audit purposes.

**Locations**:
- `internal/api/handler.go:560` - Login handler (IP, user agent, fingerprint)
- `internal/api/handler.go:674` - ForgotPassword handler (IP, user agent)

**Current State**: Handlers pass `nil` for these parameters

**Required**:
1. Add middleware or helper to extract IP from `X-Forwarded-For` / `X-Real-IP` / `RemoteAddr`
2. Extract user agent from `User-Agent` header
3. Implement device fingerprinting (browser fingerprint, TLS fingerprint, or similar)

**Impact**: Reduced security audit trail, can't track login locations/devices properly

---

### DEBT-002: Email Service Integration (Auth Service)
**Priority**: High
**Component**: Auth, Email
**Status**: Pending
**Reported**: 2026-02-02

**Description**: Auth service generates email verification and password reset tokens but doesn't send emails. Email service integration needed.

**Current State**: Tokens generated but not delivered to users

**Required**:
1. Implement email service with SMTP configuration
2. Create email templates (HTML + plain text):
   - Welcome email with verification link
   - Password reset email with reset link
   - Verification resend email
3. Wire email service into auth handlers

**Affected Endpoints**:
- `/auth/register` - Should send verification email
- `/auth/resend-verification` - Should send verification email
- `/auth/forgot-password` - Should send reset email

**Impact**: Users cannot verify emails or reset passwords without manual intervention

---

### DEBT-003: Settings Handlers Missing Auth Context
**Priority**: Low
**Component**: Settings
**Status**: Pending
**Reported**: 2026-02-02

**Description**: User settings endpoints have TODOs for getting user ID from auth context (5 handlers affected).

**Locations**:
- `internal/api/handler.go:179` - ListUserSettings
- `internal/api/handler.go:199` - GetUserSetting
- `internal/api/handler.go:215` - UpdateUserSetting
- `internal/api/handler.go:231` - DeleteUserSetting
- `internal/api/handler.go:159` - UpdateServerSetting (needs admin check)

**Current State**: Hardcoded userID = 1

**Required**: Use `GetUserID(ctx)` from auth middleware (Step 6.5)

**Impact**: Settings endpoints not properly scoped to authenticated user

---

### DEBT-004: Notification Settings JSONB Handling
**Priority**: Low
**Component**: User
**Status**: Pending
**Reported**: 2026-02-02

**Description**: User update endpoint has TODO for handling notification settings stored as JSONB.

**Location**: `internal/api/handler.go:444`

**Current State**: Notification settings fields ignored in updates

**Required**: Add JSONB unmarshaling/validation for notification_settings field

**Impact**: Cannot update notification preferences via API

---

### DEBT-005: Avatar Upload Implementation
**Priority**: Low
**Component**: User
**Status**: Pending
**Reported**: 2026-02-02

**Description**: Avatar upload endpoint exists but doesn't process multipart form data.

**Location**: `internal/api/handler.go:501`

**Current State**: Returns 501 Not Implemented

**Required**:
1. Parse multipart form data
2. Validate image file (type, size, dimensions)
3. Store file (local filesystem / S3 / CDN)
4. Update user.avatar_url
5. Generate thumbnails if needed

**Impact**: Users cannot upload profile pictures

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
| 2026-02-02 | Added 5 technical debt items from auth/user/settings TODOs |

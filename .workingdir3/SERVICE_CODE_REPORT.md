# Service Code Analysis Report

Generated: 2026-02-06
Purpose: Reference for aligning service design docs with actual code (Step 6)

---

## 1. AUTH SERVICE (`internal/service/auth/`)

**Files**: module.go, service.go, repository.go, repository_pg.go, jwt.go, mfa_integration.go, service_testing.go, + test files
**fx Module**: `auth.Module` (fx.Module("auth"))

**Service struct** (not interface):
```go
type Service struct {
  pool, repo, tokenManager, hasher, activityLogger, emailService
  jwtExpiry, refreshExpiry time.Duration
  lockoutThreshold int, lockoutWindow time.Duration, lockoutEnabled bool
}
```

**Methods** (16):
- Register, VerifyEmail, ResendVerification, RegisterFromOIDC
- Login (with ipAddress, userAgent, deviceName, deviceFingerprint), Logout (by refreshToken), LogoutAll, RefreshToken, CreateSessionForUser
- ChangePassword, RequestPasswordReset (with ipAddress, userAgent), ResetPassword
- LoginWithMFA, CompleteMFALogin, GetSessionMFAInfo

**Returns**: `*db.SharedUser` (not local User type), `*LoginResponse` (AccessToken, RefreshToken, ExpiresIn)

**TokenManager interface** (jwt.go): GenerateAccessToken, GenerateRefreshToken, ValidateAccessToken, HashRefreshToken, ExtractClaims

**Config** (from config.go AuthConfig): jwt_secret, jwt_expiry (24h), refresh_expiry (168h), lockout_threshold (5), lockout_window (15m), lockout_enabled (true)

---

## 2. SESSION SERVICE (`internal/service/session/`)

**Files**: module.go, service.go, repository.go, repository_pg.go, cached_service.go, service_testing.go
**fx Module**: `session.Module` (fx.Module("session"))

**Service struct**:
```go
type Service struct {
  repo Repository, logger *zap.Logger
  tokenLength int, expiry, refreshExpiry time.Duration, maxPerUser int
}
```

**Methods** (8):
- CreateSession(ctx, userID, DeviceInfo, scopes) → (token, refreshToken, error)
- ValidateSession(ctx, token) → (*db.SharedSession, error)
- RefreshSession(ctx, refreshToken) → (newToken, newRefreshToken, error)
- RevokeSession, RevokeAllUserSessions, RevokeAllUserSessionsExcept
- ListUserSessions → []SessionInfo
- CleanupExpiredSessions → (int, error)

**Types**: DeviceInfo{DeviceName *string, UserAgent *string, IPAddress *netip.Addr}, SessionInfo{..., IsCurrent bool}
**CachedService**: wraps Service with cache.Cache

**Config** (SessionConfig): cache_enabled (true), cache_ttl (5m), max_per_user (10), token_length (32)

---

## 3. MFA SERVICE (`internal/service/mfa/`)

**Files**: module.go, manager.go, totp.go, webauthn.go, backup_codes.go, + test files
**fx Module**: `mfa.Module` (fx.Module("mfa"))

**MFAManager struct**: queries, totp, webauthn, backupCodes, logger
**Methods** (11): GetStatus, HasAnyMethod, RequiresMFA, EnableMFA, DisableMFA, SetRememberDevice, GetRememberDeviceSettings, VerifyTOTP, VerifyBackupCode, RemoveAllMethods

**TOTPService**: GenerateSecret→TOTPSetup, VerifyCode, EnableTOTP, DisableTOTP, DeleteTOTP, HasTOTP
**WebAuthnService**: BeginRegistration, FinishRegistration, BeginLogin, FinishLogin, ListCredentials, DeleteCredential, RenameCredential, HasWebAuthn, + session management
**BackupCodesService**: GenerateCodes, RegenerateCodes, VerifyCode, GetRemainingCount, HasBackupCodes, DeleteAllCodes

**Types**: MFAStatus, VerificationResult, TOTPSetup, BackupCode, WebAuthnUser
**Dependencies**: crypto.Encryptor (TOTP secrets), crypto.PasswordHasher (backup codes), cache.Client (WebAuthn sessions)

---

## 4. OIDC SERVICE (`internal/service/oidc/`)

**Files**: service.go, repository.go, repository_pg.go, module.go
**fx Module**: `oidc.Module` (fx.Module("oidc"))

**Service struct**: repo, logger, callbackURL, encryptKey
**Methods** (17):
- Provider mgmt: AddProvider, GetProvider, GetProviderByName, GetDefaultProvider, ListProviders, ListEnabledProviders, UpdateProvider, DeleteProvider, EnableProvider, DisableProvider, SetDefaultProvider
- OAuth: GetAuthURL, HandleCallback, LinkUser, UnlinkUser, ListUserLinks
- Cleanup: CleanupExpiredStates

**Repository** (25 methods): Provider CRUD (11), UserLink CRUD (10), State management (5)

**Types**: Provider (full OIDC provider with ClaimMappings, RoleMappings, endpoints), UserLink, UserLinkWithProvider, State, AuthURLResult, CallbackResult, UserInfo
**Config**: callbackURL derived from server config, encryptKey from auth.jwt_secret

---

## 5. USER SERVICE (`internal/service/user/`)

**Files**: service.go, repository.go, repository_pg.go, module.go, cached_service.go
**fx Module**: `user.Module` (fx.Module("user"))

**Service struct**: pool, repo, hasher, activityLogger, storage, avatarConfig
**Methods** (21):
- User CRUD: GetUser, GetUserByUsername, GetUserByEmail, ListUsers, CreateUser, UpdateUser, DeleteUser, HardDeleteUser, VerifyEmail, RecordLogin
- Password: HashPassword, VerifyPassword, UpdatePassword
- Preferences: GetUserPreferences, UpdateUserPreferences, UpdateNotificationPreferences
- Avatars: GetCurrentAvatar, ListUserAvatars, UploadAvatar, SetCurrentAvatar, DeleteAvatar

**Repository** (24 methods): User (11), Preferences (3), Avatars (10)
**Types**: UserFilters, CreateUserParams, UpdateUserParams, UpsertPreferencesParams, AvatarMetadata, NotificationSettings
**CachedService**: wraps with cache.Cache
**Config**: AvatarConfig (storage_path, max_size_bytes, allowed_types)

---

## 6. RBAC SERVICE (`internal/service/rbac/`)

**Files**: service.go, module.go, permissions.go, roles.go, adapter.go, cached_service.go
**fx Module**: `rbac.Module` (fx.Module("rbac"))

**Service struct**: enforcer (*casbin.Enforcer), logger, activityLogger
**Methods** (23):
- Enforce: Enforce, EnforceWithContext
- Policy: AddPolicy, RemovePolicy, GetPolicies
- Roles: AssignRole, RemoveRole, GetUserRoles, GetUsersForRole, HasRole, ListRoles, GetRole, CreateRole, DeleteRole, UpdateRolePermissions, GetRolePermissions, AddPermissionToRole, RemovePermissionFromRole, GetAllRoleNames, CheckUserPermission
- Lifecycle: LoadPolicy, SavePolicy, ListPermissions

**Types**: Role, Permission (Resource+Action)
**Permissions**: 40+ permission constants (users, profile, movies, libraries, playback, requests, settings, audit, integrations, notifications, admin)
**Adapter**: Custom Casbin pgx adapter (shared.casbin_rule table)
**CachedService**: wraps with cache.Cache
**Config** (RBACConfig): model_path, policy_reload_interval (5m)

---

## 7. API KEYS SERVICE (`internal/service/apikeys/`)

**Files**: service.go, repository.go, repository_pg.go, module.go
**fx Module**: `apikeys.Module` (fx.Module("apikeys"))

**Service struct**: repo, logger, maxKeysPerUser, defaultExpiry
**Methods** (8): CreateKey, GetKey, ListUserKeys, ValidateKey, RevokeKey, CheckScope, UpdateScopes, CleanupExpiredKeys

**Repository** (12 methods): CreateAPIKey, GetAPIKey, GetAPIKeyByHash, GetAPIKeyByPrefix, ListUserAPIKeys, ListActiveUserAPIKeys, CountUserAPIKeys, RevokeAPIKey, UpdateAPIKeyLastUsed, UpdateAPIKeyScopes, DeleteAPIKey, DeleteExpiredAPIKeys
**Types**: CreateKeyRequest, APIKey, CreateKeyResponse (includes RawKey)
**Constants**: KeyPrefix="rv_", KeyLength=32, DefaultMaxKeysPerUser=10

---

## 8. LIBRARY SERVICE (`internal/service/library/`)

**Files**: module.go, service.go, repository.go, repository_pg.go, cached_service.go, cleanup.go
**fx Module**: `library.Module` (fx.Module("library"))

**Service struct**: repo, logger, activityLogger
**Methods** (29): Create, Get, GetByName, List, ListEnabled, ListByType, ListAccessible, Update, Delete, Count, TriggerScan, GetScan, ListScans, GetLatestScan, GetRunningScans, StartScan, CompleteScan, FailScan, CancelScan, UpdateScanProgress, GrantPermission, RevokePermission, CheckPermission, ListPermissions, ListUserPermissions, GetPermission, CanAccess, CanDownload, CanManage

**Repository** (30 methods): Library CRUD (10), Scans (9), Permissions (11)
**Types**: Library, LibraryUpdate, LibraryScan, ScanStatusUpdate, ScanProgress, Permission, CreateLibraryRequest
**Cleanup Worker**: LibraryScanCleanupWorker (River job, uses raft.LeaderElection)
**CachedService**: wraps with cache.Cache

---

## 9. SEARCH SERVICE (`internal/service/search/`)

**Files**: module.go, movie_service.go, movie_schema.go, cached_service.go
**fx Module**: `search_service.Module` (fx.Module("search_service"))

**MovieSearchService struct**: client (*search.Client), logger (*slog.Logger)
**Methods** (10): IsEnabled, InitializeCollection, IndexMovie, UpdateMovie, RemoveMovie, BulkIndexMovies, Search, Autocomplete, GetFacets, ReindexAll

**Types**: MovieDocument (full Typesense document), SearchResult, MovieHit, FacetValue, SearchParams, MovieWithRelations
**CachedService**: CachedMovieSearchService wraps with cache
**Config** (SearchConfig): url, api_key, enabled

---

## 10. SETTINGS SERVICE (`internal/service/settings/`)

**Files**: module.go, service.go, repository.go, repository_pg.go, cached_service.go
**fx Module**: `settings.Module` (fx.Module("settings"))

**Service interface** (not struct!):
- Server: GetServerSetting, ListServerSettings, ListServerSettingsByCategory, ListPublicServerSettings, SetServerSetting, DeleteServerSetting
- User: GetUserSetting, ListUserSettings, ListUserSettingsByCategory, SetUserSetting, SetUserSettingsBulk, DeleteUserSetting

**Repository** (14 methods): Server settings (7), User settings (7)
**Types**: ServerSetting (Key, Value, Description, Category, DataType, IsSecret, IsPublic, AllowedValues), UserSetting
**CachedService**: wraps with cache.Cache
**No dedicated config** - uses general database/pool config

---

## 11. ACTIVITY SERVICE (`internal/service/activity/`)

**Files**: module.go, service.go, repository.go, repository_pg.go, logger.go
**fx Module**: `activity.Module` (fx.Module("activity"))

**Service struct**: repo, logger
**Methods** (12): Log, LogWithContext, LogFailure, Get, List, Search, GetUserActivity, GetResourceActivity, GetFailedActivity, GetStats, GetRecentActions, CleanupOldLogs, CountOldLogs

**Logger interface**: LogAction, LogFailure (used by other services as dependency)
**Implementations**: ServiceLogger (wraps Service), NoopLogger (for testing)

**Repository** (14 methods): CRUD + Search + Stats + Cleanup
**Types**: Entry, SearchFilters, Stats, ActionCount, LogRequest, LogActionRequest, LogFailureRequest
**Action constants**: 27 predefined actions (user.login, session.create, library.scan, etc.)
**Resource types**: user, session, apikey, oidc, setting, library, movie, tvshow, episode
**Config** (ActivityConfig): retention_days (90)

---

## Config Summary (from config.go)

| Service | Config Struct | koanf namespace | Key fields |
|---------|--------------|-----------------|------------|
| Auth | AuthConfig | auth.* | jwt_secret, jwt_expiry, refresh_expiry, lockout_* |
| Session | SessionConfig | session.* | cache_enabled, cache_ttl, max_per_user, token_length |
| RBAC | RBACConfig | rbac.* | model_path, policy_reload_interval |
| Search | SearchConfig | search.* | url, api_key, enabled |
| Activity | ActivityConfig | activity.* | retention_days |
| Storage | StorageConfig | storage.* | backend, local.path, s3.* |
| Email | EmailConfig | email.* | enabled, provider, smtp.*, sendgrid.* |
| Avatar | AvatarConfig | avatar.* | storage_path, max_size_bytes, allowed_types |
| Movie/TMDb | MovieConfig.TMDb | movie.tmdb.* | api_key, rate_limit, cache_ttl, proxy_url |
| Integrations | IntegrationsConfig | integrations.* | radarr.*, sonarr.* |

**No dedicated config in config.go for**: MFA, OIDC, User, API Keys, Library, Settings, Notification

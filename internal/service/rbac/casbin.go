// Package rbac provides role-based access control using Casbin.
package rbac

import (
	"context"
	"embed"
	_ "embed"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxadapter "github.com/pckhoi/casbin-pgx-adapter/v3"

	"github.com/lusoris/revenge/internal/infra/database/db"
)

//go:embed model.conf
var modelConf embed.FS

var (
	// ErrPermissionDenied indicates the user lacks required permissions.
	ErrPermissionDenied = errors.New("permission denied")
	// ErrRoleNotFound indicates the role was not found.
	ErrRoleNotFound = errors.New("role not found")
	// ErrRoleInUse indicates the role cannot be deleted because it's assigned to users.
	ErrRoleInUse = errors.New("role is in use")
	// ErrSystemRole indicates the operation is not allowed on system roles.
	ErrSystemRole = errors.New("cannot modify system role")
)

// Permission constants for common permission names.
const (
	// System permissions
	PermSystemSettingsRead  = "system.settings.read"
	PermSystemSettingsWrite = "system.settings.write"
	PermSystemLogsRead      = "system.logs.read"
	PermSystemJobsRead      = "system.jobs.read"
	PermSystemJobsManage    = "system.jobs.manage"
	PermSystemAPIKeysManage = "system.apikeys.manage"
	PermSystemRolesRead     = "system.roles.read"
	PermSystemRolesManage   = "system.roles.manage"

	// User permissions
	PermUsersRead           = "users.read"
	PermUsersCreate         = "users.create"
	PermUsersUpdate         = "users.update"
	PermUsersDelete         = "users.delete"
	PermUsersSessionsManage = "users.sessions.manage"

	// Library permissions
	PermLibrariesRead   = "libraries.read"
	PermLibrariesCreate = "libraries.create"
	PermLibrariesUpdate = "libraries.update"
	PermLibrariesDelete = "libraries.delete"
	PermLibrariesScan   = "libraries.scan"

	// Content permissions
	PermContentBrowse        = "content.browse"
	PermContentMetadataRead  = "content.metadata.read"
	PermContentMetadataWrite = "content.metadata.write"
	PermContentMetadataLock  = "content.metadata.lock"
	PermContentMetadataAudit = "content.metadata.audit"
	PermContentImagesManage  = "content.images.manage"
	PermContentDelete        = "content.delete"

	// Playback permissions
	PermPlaybackStream    = "playback.stream"
	PermPlaybackDownload  = "playback.download"
	PermPlaybackTranscode = "playback.transcode"

	// Social permissions
	PermSocialRate              = "social.rate"
	PermSocialPlaylistsCreate   = "social.playlists.create"
	PermSocialPlaylistsManage   = "social.playlists.manage"
	PermSocialCollectionsCreate = "social.collections.create"
	PermSocialCollectionsManage = "social.collections.manage"
	PermSocialHistoryRead       = "social.history.read"
	PermSocialFavoritesManage   = "social.favorites.manage"

	// Request permissions
	PermRequestsSubmit       = "requests.submit"
	PermRequestsViewOwn      = "requests.view.own"
	PermRequestsVote         = "requests.vote"
	PermRequestsComment      = "requests.comment"
	PermRequestsCancelOwn    = "requests.cancel.own"
	PermRequestsViewAll      = "requests.view.all"
	PermRequestsApprove      = "requests.approve"
	PermRequestsDecline      = "requests.decline"
	PermRequestsPriority     = "requests.priority"
	PermRequestsRulesRead    = "requests.rules.read"
	PermRequestsRulesManage  = "requests.rules.manage"
	PermRequestsQuotasRead   = "requests.quotas.read"
	PermRequestsQuotasManage = "requests.quotas.manage"
	PermRequestsPollsVote    = "requests.polls.vote"
	PermRequestsPollsCreate  = "requests.polls.create"
	PermRequestsPollsManage  = "requests.polls.manage"

	// Adult content permissions
	PermAdultBrowse               = "adult.browse"
	PermAdultStream               = "adult.stream"
	PermAdultMetadataWrite        = "adult.metadata.write"
	PermAdultRequestsSubmit       = "adult.requests.submit"
	PermAdultRequestsViewOwn      = "adult.requests.view.own"
	PermAdultRequestsVote         = "adult.requests.vote"
	PermAdultRequestsApprove      = "adult.requests.approve"
	PermAdultRequestsDecline      = "adult.requests.decline"
	PermAdultRequestsRulesManage  = "adult.requests.rules.manage"
	PermAdultRequestsPollsCreate  = "adult.requests.polls.create"
	PermAdultRequestsPollsManage  = "adult.requests.polls.manage"
)

// AllPermissions returns all available permission names.
func AllPermissions() []string {
	return []string{
		// System
		PermSystemSettingsRead, PermSystemSettingsWrite, PermSystemLogsRead,
		PermSystemJobsRead, PermSystemJobsManage, PermSystemAPIKeysManage,
		PermSystemRolesRead, PermSystemRolesManage,
		// Users
		PermUsersRead, PermUsersCreate, PermUsersUpdate, PermUsersDelete, PermUsersSessionsManage,
		// Libraries
		PermLibrariesRead, PermLibrariesCreate, PermLibrariesUpdate, PermLibrariesDelete, PermLibrariesScan,
		// Content
		PermContentBrowse, PermContentMetadataRead, PermContentMetadataWrite,
		PermContentMetadataLock, PermContentMetadataAudit,
		PermContentImagesManage, PermContentDelete,
		// Playback
		PermPlaybackStream, PermPlaybackDownload, PermPlaybackTranscode,
		// Social
		PermSocialRate, PermSocialPlaylistsCreate, PermSocialPlaylistsManage,
		PermSocialCollectionsCreate, PermSocialCollectionsManage, PermSocialHistoryRead, PermSocialFavoritesManage,
		// Requests
		PermRequestsSubmit, PermRequestsViewOwn, PermRequestsVote, PermRequestsComment, PermRequestsCancelOwn,
		PermRequestsViewAll, PermRequestsApprove, PermRequestsDecline, PermRequestsPriority,
		PermRequestsRulesRead, PermRequestsRulesManage, PermRequestsQuotasRead, PermRequestsQuotasManage,
		PermRequestsPollsVote, PermRequestsPollsCreate, PermRequestsPollsManage,
		// Adult
		PermAdultBrowse, PermAdultStream, PermAdultMetadataWrite,
		PermAdultRequestsSubmit, PermAdultRequestsViewOwn, PermAdultRequestsVote,
		PermAdultRequestsApprove, PermAdultRequestsDecline, PermAdultRequestsRulesManage,
		PermAdultRequestsPollsCreate, PermAdultRequestsPollsManage,
	}
}

// Role represents a role in the system.
type Role struct {
	ID          uuid.UUID
	Name        string
	DisplayName string
	Description string
	Color       string
	Icon        string
	IsSystem    bool
	IsDefault   bool
	Priority    int
	Permissions []string
}

// CasbinService provides dynamic RBAC using Casbin.
type CasbinService struct {
	enforcer *casbin.Enforcer
	queries  *db.Queries
	pool     *pgxpool.Pool
	logger   *slog.Logger
	mu       sync.RWMutex
}

// NewCasbinService creates a new Casbin-based RBAC service.
func NewCasbinService(pool *pgxpool.Pool, queries *db.Queries, logger *slog.Logger) (*CasbinService, error) {
	// Load model from embedded file
	modelData, err := modelConf.ReadFile("model.conf")
	if err != nil {
		return nil, fmt.Errorf("failed to read model.conf: %w", err)
	}

	m, err := model.NewModelFromString(string(modelData))
	if err != nil {
		return nil, fmt.Errorf("failed to parse casbin model: %w", err)
	}

	// Create PostgreSQL adapter using pgx pool
	adapter, err := pgxadapter.NewAdapter(pool, pgxadapter.WithTableName("casbin_rules"))
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin adapter: %w", err)
	}

	// Create enforcer
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin enforcer: %w", err)
	}

	// Enable auto-save
	enforcer.EnableAutoSave(true)

	svc := &CasbinService{
		enforcer: enforcer,
		queries:  queries,
		pool:     pool,
		logger:   logger.With(slog.String("service", "rbac")),
	}

	// Seed default policies if empty
	if err := svc.seedDefaultPolicies(); err != nil {
		return nil, fmt.Errorf("failed to seed default policies: %w", err)
	}

	return svc, nil
}

// seedDefaultPolicies creates default role-permission mappings if they don't exist.
func (s *CasbinService) seedDefaultPolicies() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if policies already exist
	policies, err := s.enforcer.GetPolicy()
	if err != nil {
		return err
	}
	if len(policies) > 0 {
		s.logger.Info("Casbin policies already exist, skipping seed")
		return nil
	}

	s.logger.Info("Seeding default Casbin policies")

	// Admin: All permissions
	for _, perm := range AllPermissions() {
		if _, err := s.enforcer.AddPolicy("admin", perm, "allow"); err != nil {
			return fmt.Errorf("failed to add admin policy for %s: %w", perm, err)
		}
	}

	// Moderator permissions
	moderatorPerms := []string{
		PermSystemLogsRead, PermSystemJobsRead, PermSystemRolesRead,
		PermUsersRead,
		PermLibrariesRead, PermLibrariesCreate, PermLibrariesUpdate, PermLibrariesScan,
		PermContentBrowse, PermContentMetadataRead, PermContentMetadataWrite,
		PermContentMetadataLock, PermContentMetadataAudit, PermContentImagesManage,
		PermPlaybackStream, PermPlaybackDownload, PermPlaybackTranscode,
		PermSocialRate, PermSocialPlaylistsCreate, PermSocialPlaylistsManage,
		PermSocialCollectionsCreate, PermSocialCollectionsManage, PermSocialHistoryRead, PermSocialFavoritesManage,
		PermRequestsSubmit, PermRequestsViewOwn, PermRequestsVote, PermRequestsComment, PermRequestsCancelOwn,
		PermRequestsViewAll, PermRequestsApprove, PermRequestsDecline, PermRequestsPriority,
		PermRequestsRulesRead, PermRequestsPollsVote, PermRequestsPollsCreate,
		// Adult content - moderators have full adult access
		PermAdultBrowse, PermAdultStream, PermAdultMetadataWrite,
		PermAdultRequestsSubmit, PermAdultRequestsViewOwn, PermAdultRequestsVote,
		PermAdultRequestsApprove, PermAdultRequestsDecline, PermAdultRequestsRulesManage,
		PermAdultRequestsPollsCreate, PermAdultRequestsPollsManage,
	}
	for _, perm := range moderatorPerms {
		if _, err := s.enforcer.AddPolicy("moderator", perm, "allow"); err != nil {
			return fmt.Errorf("failed to add moderator policy for %s: %w", perm, err)
		}
	}

	// User permissions
	userPerms := []string{
		PermLibrariesRead,
		PermContentBrowse, PermContentMetadataRead,
		PermPlaybackStream, PermPlaybackDownload, PermPlaybackTranscode,
		PermSocialRate, PermSocialPlaylistsCreate, PermSocialPlaylistsManage,
		PermSocialCollectionsCreate, PermSocialCollectionsManage, PermSocialHistoryRead, PermSocialFavoritesManage,
		PermRequestsSubmit, PermRequestsViewOwn, PermRequestsVote, PermRequestsComment, PermRequestsCancelOwn,
		PermRequestsPollsVote,
	}
	for _, perm := range userPerms {
		if _, err := s.enforcer.AddPolicy("user", perm, "allow"); err != nil {
			return fmt.Errorf("failed to add user policy for %s: %w", perm, err)
		}
	}

	// Guest permissions
	guestPerms := []string{
		PermLibrariesRead,
		PermContentBrowse, PermContentMetadataRead,
	}
	for _, perm := range guestPerms {
		if _, err := s.enforcer.AddPolicy("guest", perm, "allow"); err != nil {
			return fmt.Errorf("failed to add guest policy for %s: %w", perm, err)
		}
	}

	s.logger.Info("Default Casbin policies seeded successfully")
	return nil
}

// HasPermission checks if a user has a specific permission.
func (s *CasbinService) HasPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error) {
	// Get user's role
	role, err := s.getUserRole(ctx, userID)
	if err != nil {
		return false, err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	allowed, err := s.enforcer.Enforce(role, permission, "allow")
	if err != nil {
		return false, fmt.Errorf("casbin enforce error: %w", err)
	}

	return allowed, nil
}

// HasAnyPermission checks if a user has any of the specified permissions.
func (s *CasbinService) HasAnyPermission(ctx context.Context, userID uuid.UUID, permissions []string) (bool, error) {
	for _, perm := range permissions {
		has, err := s.HasPermission(ctx, userID, perm)
		if err != nil {
			return false, err
		}
		if has {
			return true, nil
		}
	}
	return false, nil
}

// RequirePermission checks if a user has the required permission and returns an error if not.
func (s *CasbinService) RequirePermission(ctx context.Context, userID uuid.UUID, permission string) error {
	has, err := s.HasPermission(ctx, userID, permission)
	if err != nil {
		return err
	}
	if !has {
		s.logger.Debug("Permission denied",
			slog.String("user_id", userID.String()),
			slog.String("permission", permission),
		)
		return ErrPermissionDenied
	}
	return nil
}

// RequireAnyPermission checks if a user has any of the required permissions.
func (s *CasbinService) RequireAnyPermission(ctx context.Context, userID uuid.UUID, permissions []string) error {
	has, err := s.HasAnyPermission(ctx, userID, permissions)
	if err != nil {
		return err
	}
	if !has {
		s.logger.Debug("Permission denied",
			slog.String("user_id", userID.String()),
			slog.Any("permissions", permissions),
		)
		return ErrPermissionDenied
	}
	return nil
}

// GetUserPermissions returns all permissions for a user based on their role.
func (s *CasbinService) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, error) {
	role, err := s.getUserRole(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.GetRolePermissions(role)
}

// GetRolePermissions returns all permissions assigned to a role.
func (s *CasbinService) GetRolePermissions(roleName string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	policies, err := s.enforcer.GetFilteredPolicy(0, roleName)
	if err != nil {
		return nil, err
	}

	permissions := make([]string, 0, len(policies))
	for _, policy := range policies {
		if len(policy) >= 2 {
			permissions = append(permissions, policy[1])
		}
	}

	return permissions, nil
}

// AddRolePermission adds a permission to a role.
func (s *CasbinService) AddRolePermission(ctx context.Context, roleName, permission string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.enforcer.AddPolicy(roleName, permission, "allow")
	if err != nil {
		return fmt.Errorf("failed to add permission: %w", err)
	}

	s.logger.Info("Added permission to role",
		slog.String("role", roleName),
		slog.String("permission", permission),
	)

	return nil
}

// RemoveRolePermission removes a permission from a role.
func (s *CasbinService) RemoveRolePermission(ctx context.Context, roleName, permission string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.enforcer.RemovePolicy(roleName, permission, "allow")
	if err != nil {
		return fmt.Errorf("failed to remove permission: %w", err)
	}

	s.logger.Info("Removed permission from role",
		slog.String("role", roleName),
		slog.String("permission", permission),
	)

	return nil
}

// SetRolePermissions replaces all permissions for a role.
func (s *CasbinService) SetRolePermissions(ctx context.Context, roleName string, permissions []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove all existing permissions for this role
	_, err := s.enforcer.RemoveFilteredPolicy(0, roleName)
	if err != nil {
		return fmt.Errorf("failed to clear role permissions: %w", err)
	}

	// Add new permissions
	for _, perm := range permissions {
		if _, err := s.enforcer.AddPolicy(roleName, perm, "allow"); err != nil {
			return fmt.Errorf("failed to add permission %s: %w", perm, err)
		}
	}

	s.logger.Info("Set role permissions",
		slog.String("role", roleName),
		slog.Int("count", len(permissions)),
	)

	return nil
}

// ListRoles returns all roles from the database.
func (s *CasbinService) ListRoles(ctx context.Context) ([]Role, error) {
	dbRoles, err := s.queries.ListRoles(ctx)
	if err != nil {
		return nil, err
	}

	roles := make([]Role, len(dbRoles))
	for i, r := range dbRoles {
		perms, _ := s.GetRolePermissions(r.Name)
		roles[i] = Role{
			ID:          r.ID,
			Name:        r.Name,
			DisplayName: r.DisplayName,
			Description: stringFromPtr(r.Description),
			Color:       stringFromPtr(r.Color),
			Icon:        stringFromPtr(r.Icon),
			IsSystem:    r.IsSystem,
			IsDefault:   r.IsDefault,
			Priority:    int(r.Priority),
			Permissions: perms,
		}
	}

	return roles, nil
}

// GetRole returns a role by name.
func (s *CasbinService) GetRole(ctx context.Context, name string) (*Role, error) {
	r, err := s.queries.GetRoleByName(ctx, name)
	if err != nil {
		return nil, ErrRoleNotFound
	}

	perms, _ := s.GetRolePermissions(r.Name)

	return &Role{
		ID:          r.ID,
		Name:        r.Name,
		DisplayName: r.DisplayName,
		Description: stringFromPtr(r.Description),
		Color:       stringFromPtr(r.Color),
		Icon:        stringFromPtr(r.Icon),
		IsSystem:    r.IsSystem,
		IsDefault:   r.IsDefault,
		Priority:    int(r.Priority),
		Permissions: perms,
	}, nil
}

// CreateRole creates a new custom role.
func (s *CasbinService) CreateRole(ctx context.Context, params CreateRoleParams) (*Role, error) {
	role, err := s.queries.CreateRole(ctx, db.CreateRoleParams{
		Name:        params.Name,
		DisplayName: params.DisplayName,
		Description: ptrFromString(params.Description),
		Color:       ptrFromString(params.Color),
		Icon:        ptrFromString(params.Icon),
		Priority:    int32(params.Priority),
		CreatedBy:   pgtypeUUIDFromUUID(params.CreatedBy),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	// Add permissions via Casbin
	if len(params.Permissions) > 0 {
		if err := s.SetRolePermissions(ctx, params.Name, params.Permissions); err != nil {
			return nil, err
		}
	}

	s.logger.Info("Created role",
		slog.String("name", params.Name),
		slog.String("created_by", params.CreatedBy.String()),
	)

	return &Role{
		ID:          role.ID,
		Name:        role.Name,
		DisplayName: role.DisplayName,
		Description: params.Description,
		Color:       params.Color,
		Icon:        params.Icon,
		IsSystem:    false,
		IsDefault:   false,
		Priority:    params.Priority,
		Permissions: params.Permissions,
	}, nil
}

// UpdateRole updates an existing role.
func (s *CasbinService) UpdateRole(ctx context.Context, name string, params UpdateRoleParams) (*Role, error) {
	// Check if it's a system role
	existing, err := s.queries.GetRoleByName(ctx, name)
	if err != nil {
		return nil, ErrRoleNotFound
	}

	if existing.IsSystem && (params.Name != "" && params.Name != name) {
		return nil, ErrSystemRole
	}

	// Update role in database
	role, err := s.queries.UpdateRole(ctx, db.UpdateRoleParams{
		Name:        name,
		DisplayName: params.DisplayName,
		Description: ptrFromString(params.Description),
		Color:       ptrFromString(params.Color),
		Icon:        ptrFromString(params.Icon),
		Priority:    int32(params.Priority),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	// Update permissions if provided
	if params.Permissions != nil {
		if err := s.SetRolePermissions(ctx, name, params.Permissions); err != nil {
			return nil, err
		}
	}

	perms, _ := s.GetRolePermissions(name)

	s.logger.Info("Updated role", slog.String("name", name))

	return &Role{
		ID:          role.ID,
		Name:        role.Name,
		DisplayName: role.DisplayName,
		Description: stringFromPtr(role.Description),
		Color:       stringFromPtr(role.Color),
		Icon:        stringFromPtr(role.Icon),
		IsSystem:    role.IsSystem,
		IsDefault:   role.IsDefault,
		Priority:    int(role.Priority),
		Permissions: perms,
	}, nil
}

// DeleteRole deletes a custom role.
func (s *CasbinService) DeleteRole(ctx context.Context, name string) error {
	// Check if it's a system role
	existing, err := s.queries.GetRoleByName(ctx, name)
	if err != nil {
		return ErrRoleNotFound
	}

	if existing.IsSystem {
		return ErrSystemRole
	}

	// Check if role is in use
	count, err := s.queries.CountUsersWithRole(ctx, pgtype.UUID{Bytes: existing.ID, Valid: true})
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrRoleInUse
	}

	// Remove all Casbin policies for this role
	s.mu.Lock()
	_, err = s.enforcer.RemoveFilteredPolicy(0, name)
	s.mu.Unlock()
	if err != nil {
		return fmt.Errorf("failed to remove role policies: %w", err)
	}

	// Delete from database
	if err := s.queries.DeleteRole(ctx, name); err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	s.logger.Info("Deleted role", slog.String("name", name))

	return nil
}

// CanAccessAdultContent checks if a user can access adult content.
func (s *CasbinService) CanAccessAdultContent(ctx context.Context, user *db.User) (bool, error) {
	if !user.AdultEnabled {
		return false, nil
	}
	return s.HasPermission(ctx, user.ID, PermAdultBrowse)
}

// ReloadPolicies reloads all policies from the database.
func (s *CasbinService) ReloadPolicies() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.enforcer.LoadPolicy()
}

// getUserRole gets the role name for a user.
func (s *CasbinService) getUserRole(ctx context.Context, userID uuid.UUID) (string, error) {
	role, err := s.queries.GetUserRoleName(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user role: %w", err)
	}
	return role, nil
}

// SetUserRole assigns a role to a user by role name.
func (s *CasbinService) SetUserRole(ctx context.Context, userID uuid.UUID, roleName string) error {
	// Verify role exists
	_, err := s.queries.GetRoleByName(ctx, roleName)
	if err != nil {
		return ErrRoleNotFound
	}

	if err := s.queries.SetUserRole(ctx, db.SetUserRoleParams{
		ID:   userID,
		Name: roleName,
	}); err != nil {
		return fmt.Errorf("failed to set user role: %w", err)
	}

	s.logger.Info("User role changed",
		slog.String("user_id", userID.String()),
		slog.String("role", roleName),
	)

	return nil
}

// GetUserRole returns the full Role struct for a user.
func (s *CasbinService) GetUserRole(ctx context.Context, userID uuid.UUID) (*Role, error) {
	roleName, err := s.getUserRole(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.GetRole(ctx, roleName)
}

// GetUsersForRole returns all users with a specific role.
func (s *CasbinService) GetUsersForRole(ctx context.Context, roleName string, limit, offset int) ([]uuid.UUID, error) {
	users, err := s.queries.ListUsersByRoleName(ctx, db.ListUsersByRoleNameParams{
		Name:   roleName,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list users for role: %w", err)
	}

	userIDs := make([]uuid.UUID, len(users))
	for i, u := range users {
		userIDs[i] = u.ID
	}

	return userIDs, nil
}

// CountUsersForRole returns the number of users with a specific role.
func (s *CasbinService) CountUsersForRole(ctx context.Context, roleName string) (int64, error) {
	return s.queries.CountUsersByRoleName(ctx, roleName)
}

// AddRoleForUser is an alias for SetUserRole (single role per user model).
func (s *CasbinService) AddRoleForUser(ctx context.Context, userID uuid.UUID, roleName string) error {
	return s.SetUserRole(ctx, userID, roleName)
}

// RemoveRoleForUser removes a user's role by setting them to the default role.
func (s *CasbinService) RemoveRoleForUser(ctx context.Context, userID uuid.UUID) error {
	defaultRole, err := s.queries.GetDefaultRole(ctx)
	if err != nil {
		return fmt.Errorf("failed to get default role: %w", err)
	}

	return s.SetUserRole(ctx, userID, defaultRole.Name)
}

// CreateRoleParams contains parameters for creating a role.
type CreateRoleParams struct {
	Name        string
	DisplayName string
	Description string
	Color       string
	Icon        string
	Priority    int
	Permissions []string
	CreatedBy   uuid.UUID
}

// UpdateRoleParams contains parameters for updating a role.
type UpdateRoleParams struct {
	Name        string
	DisplayName string
	Description string
	Color       string
	Icon        string
	Priority    int
	Permissions []string
}

// Helper functions
func stringFromPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ptrFromString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func pgtypeUUIDFromUUID(id uuid.UUID) pgtype.UUID {
	if id == uuid.Nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: id, Valid: true}
}

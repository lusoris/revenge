package rbac

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"github.com/google/uuid"
)

// Errors for role management
var (
	ErrRoleNotFound      = errors.New("role not found")
	ErrRoleAlreadyExists = errors.New("role already exists")
	ErrBuiltInRole       = errors.New("cannot modify built-in role")
	ErrRoleInUse         = errors.New("role is assigned to users")
)

// Role represents a role with its permissions.
type Role struct {
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Permissions []Permission `json:"permissions"`
	IsBuiltIn   bool         `json:"is_built_in"`
	UserCount   int          `json:"user_count"`
}

// Permission represents a single permission (object + action).
type Permission struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

// String returns the permission as a string (resource:action).
func (p Permission) String() string {
	return fmt.Sprintf("%s:%s", p.Resource, p.Action)
}

// ParsePermission parses a permission string (resource:action).
func ParsePermission(s string) (Permission, error) {
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return Permission{}, fmt.Errorf("invalid permission format: %s", s)
	}
	return Permission{Resource: parts[0], Action: parts[1]}, nil
}

// BuiltInRoles defines the system's built-in roles that cannot be deleted.
var BuiltInRoles = map[string]string{
	"admin":     "Full system access",
	"moderator": "Content moderation and user management",
	"user":      "Standard user access",
	"guest":     "Read-only access",
}

// AvailableResources defines all available resources in the system.
// Uses FineGrainedResources from permissions.go for consistency.
var AvailableResources = FineGrainedResources

// AvailableActions defines all available actions.
// Uses FineGrainedActions from permissions.go for consistency.
var AvailableActions = FineGrainedActions

// ListRoles returns all roles with their permissions.
func (s *Service) ListRoles(ctx context.Context) ([]Role, error) {
	policies, err := s.enforcer.GetPolicy()
	if err != nil {
		return nil, fmt.Errorf("failed to get policies: %w", err)
	}

	// Group policies by role
	rolePerms := make(map[string][]Permission)
	for _, policy := range policies {
		if len(policy) >= 3 {
			roleName := policy[0]
			perm := Permission{Resource: policy[1], Action: policy[2]}
			rolePerms[roleName] = append(rolePerms[roleName], perm)
		}
	}

	// Build role list
	roles := make([]Role, 0, len(rolePerms))
	for roleName, perms := range rolePerms {
		desc, isBuiltIn := BuiltInRoles[roleName]

		// Get user count
		users, _ := s.enforcer.GetUsersForRole(roleName)

		roles = append(roles, Role{
			Name:        roleName,
			Description: desc,
			Permissions: perms,
			IsBuiltIn:   isBuiltIn,
			UserCount:   len(users),
		})
	}

	// Sort by name
	sort.Slice(roles, func(i, j int) bool {
		return roles[i].Name < roles[j].Name
	})

	return roles, nil
}

// GetRole returns a single role by name.
func (s *Service) GetRole(ctx context.Context, name string) (*Role, error) {
	policies, err := s.enforcer.GetFilteredPolicy(0, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get policies for role: %w", err)
	}

	if len(policies) == 0 {
		return nil, ErrRoleNotFound
	}

	perms := make([]Permission, 0, len(policies))
	for _, policy := range policies {
		if len(policy) >= 3 {
			perms = append(perms, Permission{Resource: policy[1], Action: policy[2]})
		}
	}

	desc, isBuiltIn := BuiltInRoles[name]
	users, _ := s.enforcer.GetUsersForRole(name)

	return &Role{
		Name:        name,
		Description: desc,
		Permissions: perms,
		IsBuiltIn:   isBuiltIn,
		UserCount:   len(users),
	}, nil
}

// CreateRole creates a new role with the given permissions.
func (s *Service) CreateRole(ctx context.Context, name, description string, permissions []Permission) (*Role, error) {
	// Check if role already exists
	existing, err := s.enforcer.GetFilteredPolicy(0, name)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing role: %w", err)
	}
	if len(existing) > 0 {
		return nil, ErrRoleAlreadyExists
	}

	// Validate role name
	if name == "" {
		return nil, fmt.Errorf("role name cannot be empty")
	}
	if strings.ContainsAny(name, " \t\n") {
		return nil, fmt.Errorf("role name cannot contain whitespace")
	}

	// Add policies for each permission
	for _, perm := range permissions {
		if _, err := s.enforcer.AddPolicy(name, perm.Resource, perm.Action); err != nil {
			// Rollback on failure
			if _, rollbackErr := s.enforcer.RemoveFilteredPolicy(0, name); rollbackErr != nil {
				s.logger.Error("failed to rollback policies during role creation failure",
					slog.String("role", name),
					slog.Any("error", rollbackErr))
			}
			return nil, fmt.Errorf("failed to add permission %s: %w", perm.String(), err)
		}
	}

	s.logger.Info("role created",
		slog.String("name", name),
		slog.Int("permissions", len(permissions)),
	)

	return &Role{
		Name:        name,
		Description: description,
		Permissions: permissions,
		IsBuiltIn:   false,
		UserCount:   0,
	}, nil
}

// DeleteRole deletes a role if it's not built-in and has no users.
func (s *Service) DeleteRole(ctx context.Context, name string) error {
	// Check if built-in
	if _, isBuiltIn := BuiltInRoles[name]; isBuiltIn {
		return ErrBuiltInRole
	}

	// Check if role exists
	existing, err := s.enforcer.GetFilteredPolicy(0, name)
	if err != nil {
		return fmt.Errorf("failed to check existing role: %w", err)
	}
	if len(existing) == 0 {
		return ErrRoleNotFound
	}

	// Check if role has users
	users, err := s.enforcer.GetUsersForRole(name)
	if err != nil {
		return fmt.Errorf("failed to get users for role: %w", err)
	}
	if len(users) > 0 {
		return ErrRoleInUse
	}

	// Remove all policies for this role
	if _, err := s.enforcer.RemoveFilteredPolicy(0, name); err != nil {
		return fmt.Errorf("failed to remove role policies: %w", err)
	}

	s.logger.Info("role deleted", slog.String("name", name))
	return nil
}

// UpdateRolePermissions updates the permissions for a role.
func (s *Service) UpdateRolePermissions(ctx context.Context, name string, permissions []Permission) (*Role, error) {
	// Check if role exists
	existing, err := s.enforcer.GetFilteredPolicy(0, name)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing role: %w", err)
	}
	if len(existing) == 0 {
		return nil, ErrRoleNotFound
	}

	// Get current description
	desc, isBuiltIn := BuiltInRoles[name]

	// Remove old permissions
	if _, err := s.enforcer.RemoveFilteredPolicy(0, name); err != nil {
		return nil, fmt.Errorf("failed to remove old permissions: %w", err)
	}

	// Add new permissions
	for _, perm := range permissions {
		if _, err := s.enforcer.AddPolicy(name, perm.Resource, perm.Action); err != nil {
			return nil, fmt.Errorf("failed to add permission %s: %w", perm.String(), err)
		}
	}

	users, _ := s.enforcer.GetUsersForRole(name)

	s.logger.Info("role permissions updated",
		slog.String("name", name),
		slog.Int("permissions", len(permissions)),
	)

	return &Role{
		Name:        name,
		Description: desc,
		Permissions: permissions,
		IsBuiltIn:   isBuiltIn,
		UserCount:   len(users),
	}, nil
}

// ListPermissions returns all available permissions.
func (s *Service) ListPermissions(ctx context.Context) []Permission {
	perms := make([]Permission, 0, len(AvailableResources)*len(AvailableActions))

	for _, resource := range AvailableResources {
		for _, action := range AvailableActions {
			perms = append(perms, Permission{
				Resource: resource,
				Action:   action,
			})
		}
	}

	return perms
}

// GetRolePermissions returns permissions for a specific role.
func (s *Service) GetRolePermissions(ctx context.Context, role string) ([]Permission, error) {
	policies, err := s.enforcer.GetFilteredPolicy(0, role)
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}

	perms := make([]Permission, 0, len(policies))
	for _, policy := range policies {
		if len(policy) >= 3 {
			perms = append(perms, Permission{
				Resource: policy[1],
				Action:   policy[2],
			})
		}
	}

	return perms, nil
}

// AddPermissionToRole adds a single permission to a role.
func (s *Service) AddPermissionToRole(ctx context.Context, role string, perm Permission) error {
	added, err := s.enforcer.AddPolicy(role, perm.Resource, perm.Action)
	if err != nil {
		return fmt.Errorf("failed to add permission: %w", err)
	}
	if !added {
		return fmt.Errorf("permission already exists for role")
	}

	s.logger.Info("permission added to role",
		slog.String("role", role),
		slog.String("permission", perm.String()),
	)

	return nil
}

// RemovePermissionFromRole removes a single permission from a role.
func (s *Service) RemovePermissionFromRole(ctx context.Context, role string, perm Permission) error {
	removed, err := s.enforcer.RemovePolicy(role, perm.Resource, perm.Action)
	if err != nil {
		return fmt.Errorf("failed to remove permission: %w", err)
	}
	if !removed {
		return fmt.Errorf("permission not found for role")
	}

	s.logger.Info("permission removed from role",
		slog.String("role", role),
		slog.String("permission", perm.String()),
	)

	return nil
}

// GetAllRoleNames returns just the names of all roles.
func (s *Service) GetAllRoleNames(ctx context.Context) ([]string, error) {
	policies, err := s.enforcer.GetPolicy()
	if err != nil {
		return nil, fmt.Errorf("failed to get policies: %w", err)
	}

	roleSet := make(map[string]struct{})
	for _, policy := range policies {
		if len(policy) >= 1 {
			roleSet[policy[0]] = struct{}{}
		}
	}

	roles := make([]string, 0, len(roleSet))
	for role := range roleSet {
		roles = append(roles, role)
	}
	sort.Strings(roles)

	return roles, nil
}

// CheckUserPermission checks if a user has a specific permission (via any of their roles).
func (s *Service) CheckUserPermission(ctx context.Context, userID uuid.UUID, resource, action string) (bool, error) {
	return s.Enforce(ctx, userID.String(), resource, action)
}

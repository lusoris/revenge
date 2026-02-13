package rbac

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	"github.com/lusoris/revenge/internal/service/activity"
)

// Service provides RBAC functionality using Casbin.
type Service struct {
	enforcer       *casbin.SyncedEnforcer
	logger         *slog.Logger
	activityLogger activity.Logger
}

// NewService creates a new RBAC service.
func NewService(enforcer *casbin.SyncedEnforcer, logger *slog.Logger, activityLogger activity.Logger) *Service {
	return &Service{
		enforcer:       enforcer,
		logger:         logger.With("component", "rbac"),
		activityLogger: activityLogger,
	}
}

// Enforce checks if a subject has permission to perform an action on an object.
func (s *Service) Enforce(ctx context.Context, sub, obj, act string) (bool, error) {
	allowed, err := s.enforcer.Enforce(sub, obj, act)
	if err != nil {
		s.logger.Error("failed to enforce policy",
			slog.String("subject", sub),
			slog.String("object", obj),
			slog.String("action", act),
			slog.Any("error", err),
		)
		return false, fmt.Errorf("failed to enforce policy: %w", err)
	}

	s.logger.Debug("policy enforcement",
		slog.String("subject", sub),
		slog.String("object", obj),
		slog.String("action", act),
		slog.Bool("allowed", allowed),
	)

	return allowed, nil
}

// EnforceWithContext checks if a user has permission to perform an action on a resource.
// This is a convenience method that converts userID to string.
func (s *Service) EnforceWithContext(ctx context.Context, userID uuid.UUID, resource, action string) (bool, error) {
	return s.Enforce(ctx, userID.String(), resource, action)
}

// AddPolicy adds a policy rule.
func (s *Service) AddPolicy(ctx context.Context, sub, obj, act string) error {
	added, err := s.enforcer.AddPolicy(sub, obj, act)
	if err != nil {
		s.logger.Error("failed to add policy",
			slog.String("subject", sub),
			slog.String("object", obj),
			slog.String("action", act),
			slog.Any("error", err),
		)
		return fmt.Errorf("failed to add policy: %w", err)
	}

	if !added {
		s.logger.Warn("policy already exists",
			slog.String("subject", sub),
			slog.String("object", obj),
			slog.String("action", act),
		)
	}

	s.logger.Info("policy added",
		slog.String("subject", sub),
		slog.String("object", obj),
		slog.String("action", act),
	)

	return nil
}

// RemovePolicy removes a policy rule.
func (s *Service) RemovePolicy(ctx context.Context, sub, obj, act string) error {
	removed, err := s.enforcer.RemovePolicy(sub, obj, act)
	if err != nil {
		s.logger.Error("failed to remove policy",
			slog.String("subject", sub),
			slog.String("object", obj),
			slog.String("action", act),
			slog.Any("error", err),
		)
		return fmt.Errorf("failed to remove policy: %w", err)
	}

	if !removed {
		s.logger.Warn("policy not found",
			slog.String("subject", sub),
			slog.String("object", obj),
			slog.String("action", act),
		)
		return fmt.Errorf("policy not found")
	}

	s.logger.Info("policy removed",
		slog.String("subject", sub),
		slog.String("object", obj),
		slog.String("action", act),
	)

	return nil
}

// GetPolicies returns all policy rules.
func (s *Service) GetPolicies(ctx context.Context) ([][]string, error) {
	policies, err := s.enforcer.GetPolicy()
	if err != nil {
		s.logger.Error("failed to get policies", slog.Any("error", err))
		return nil, fmt.Errorf("failed to get policies: %w", err)
	}
	s.logger.Debug("retrieved policies", slog.Int("count", len(policies)))
	return policies, nil
}

// AssignRole assigns a role to a user.
func (s *Service) AssignRole(ctx context.Context, userID uuid.UUID, role string) error {
	added, err := s.enforcer.AddRoleForUser(userID.String(), role)
	if err != nil {
		s.logger.Error("failed to assign role",
			slog.String("user_id", userID.String()),
			slog.String("role", role),
			slog.Any("error", err),
		)
		return fmt.Errorf("failed to assign role: %w", err)
	}

	if !added {
		s.logger.Warn("role already assigned",
			slog.String("user_id", userID.String()),
			slog.String("role", role),
		)
	}

	s.logger.Info("role assigned",
		slog.String("user_id", userID.String()),
		slog.String("role", role),
	)

	// Log role assignment
	_ = s.activityLogger.LogAction(ctx, activity.LogActionRequest{
		ResourceType: "role",
		Action:       activity.ActionAdminRoleAssign,
		ResourceID:   userID,
		Metadata: map[string]any{
			"role":    role,
			"user_id": userID.String(),
		},
	})

	return nil
}

// RemoveRole removes a role from a user.
func (s *Service) RemoveRole(ctx context.Context, userID uuid.UUID, role string) error {
	removed, err := s.enforcer.DeleteRoleForUser(userID.String(), role)
	if err != nil {
		s.logger.Error("failed to remove role",
			slog.String("user_id", userID.String()),
			slog.String("role", role),
			slog.Any("error", err),
		)
		return fmt.Errorf("failed to remove role: %w", err)
	}

	if !removed {
		s.logger.Warn("role not found",
			slog.String("user_id", userID.String()),
			slog.String("role", role),
		)
		return fmt.Errorf("role not found")
	}

	s.logger.Info("role removed",
		slog.String("user_id", userID.String()),
		slog.String("role", role),
	)

	// Log role removal
	_ = s.activityLogger.LogAction(ctx, activity.LogActionRequest{
		ResourceType: "role",
		Action:       activity.ActionAdminRoleRevoke,
		ResourceID:   userID,
		Metadata: map[string]any{
			"role":    role,
			"user_id": userID.String(),
		},
	})

	return nil
}

// GetUserRoles returns all roles assigned to a user.
func (s *Service) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]string, error) {
	roles, err := s.enforcer.GetRolesForUser(userID.String())
	if err != nil {
		s.logger.Error("failed to get user roles",
			slog.String("user_id", userID.String()),
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	s.logger.Debug("retrieved user roles",
		slog.String("user_id", userID.String()),
		slog.Int("count", len(roles)),
	)

	return roles, nil
}

// GetUsersForRole returns all users that have a specific role.
func (s *Service) GetUsersForRole(ctx context.Context, role string) ([]uuid.UUID, error) {
	users, err := s.enforcer.GetUsersForRole(role)
	if err != nil {
		s.logger.Error("failed to get users for role",
			slog.String("role", role),
			slog.Any("error", err),
		)
		return nil, fmt.Errorf("failed to get users for role: %w", err)
	}

	userIDs := make([]uuid.UUID, 0, len(users))
	for _, userStr := range users {
		userID, err := uuid.Parse(userStr)
		if err != nil {
			s.logger.Warn("invalid user ID in role mapping",
				slog.String("user_str", userStr),
				slog.String("role", role),
				slog.Any("error", err),
			)
			continue
		}
		userIDs = append(userIDs, userID)
	}

	s.logger.Debug("retrieved users for role",
		slog.String("role", role),
		slog.Int("count", len(userIDs)),
	)

	return userIDs, nil
}

// HasRole checks if a user has a specific role.
func (s *Service) HasRole(ctx context.Context, userID uuid.UUID, role string) (bool, error) {
	hasRole, err := s.enforcer.HasRoleForUser(userID.String(), role)
	if err != nil {
		s.logger.Error("failed to check user role",
			slog.String("user_id", userID.String()),
			slog.String("role", role),
			slog.Any("error", err),
		)
		return false, fmt.Errorf("failed to check user role: %w", err)
	}

	return hasRole, nil
}

// LoadPolicy reloads the policy from the database.
func (s *Service) LoadPolicy(ctx context.Context) error {
	if err := s.enforcer.LoadPolicy(); err != nil {
		s.logger.Error("failed to load policy", slog.Any("error", err))
		return fmt.Errorf("failed to load policy: %w", err)
	}

	s.logger.Info("policy reloaded successfully")
	return nil
}

// SavePolicy saves the current policy to the database.
func (s *Service) SavePolicy(ctx context.Context) error {
	if err := s.enforcer.SavePolicy(); err != nil {
		s.logger.Error("failed to save policy", slog.Any("error", err))
		return fmt.Errorf("failed to save policy: %w", err)
	}

	s.logger.Info("policy saved successfully")
	return nil
}

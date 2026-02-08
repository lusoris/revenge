package api

import (
	"context"
	"errors"

	"log/slog"

	"github.com/lusoris/revenge/internal/api/ogen"
	"github.com/lusoris/revenge/internal/service/rbac"
)

// ListPolicies lists all authorization policies (admin only).
func (h *Handler) ListPolicies(ctx context.Context) (ogen.ListPoliciesRes, error) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.ListPoliciesUnauthorized{}, nil
	}

	// Check if user is admin
	isAdmin, err := h.rbacService.HasRole(ctx, userID, "admin")
	if err != nil {
		h.logger.Error("failed to check admin role",
			slog.String("user_id", userID.String()),
			slog.Any("error",err),
		)
		return &ogen.ListPoliciesForbidden{
			Code:    500,
			Message: "Failed to check permissions",
		}, nil
	}

	if !isAdmin {
		return &ogen.ListPoliciesForbidden{
			Code:    403,
			Message: "Admin access required",
		}, nil
	}

	policies, err := h.rbacService.GetPolicies(ctx)
	if err != nil {
		h.logger.Error("failed to get policies",
			slog.Any("error",err),
		)
		return &ogen.ListPoliciesForbidden{
			Code:    500,
			Message: "Failed to get policies",
		}, nil
	}

	policyList := make([]ogen.Policy, 0, len(policies))
	for _, p := range policies {
		if len(p) >= 3 {
			policyList = append(policyList, ogen.Policy{
				Subject: p[0],
				Object:  p[1],
				Action:  p[2],
			})
		}
	}

	return &ogen.PolicyListResponse{
		Policies: policyList,
		Total:    int64(len(policyList)),
	}, nil
}

// AddPolicy adds a new authorization policy (admin only).
func (h *Handler) AddPolicy(ctx context.Context, req *ogen.PolicyRequest) (ogen.AddPolicyRes, error) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.AddPolicyUnauthorized{}, nil
	}

	// Check if user is admin
	isAdmin, err := h.rbacService.HasRole(ctx, userID, "admin")
	if err != nil {
		h.logger.Error("failed to check admin role",
			slog.String("user_id", userID.String()),
			slog.Any("error",err),
		)
		return &ogen.AddPolicyForbidden{
			Code:    500,
			Message: "Failed to check permissions",
		}, nil
	}

	if !isAdmin {
		return &ogen.AddPolicyForbidden{
			Code:    403,
			Message: "Admin access required",
		}, nil
	}

	if err := h.rbacService.AddPolicy(ctx, req.Subject, req.Object, req.Action); err != nil {
		h.logger.Error("failed to add policy",
			slog.String("subject", req.Subject),
			slog.String("object", req.Object),
			slog.String("action", req.Action),
			slog.Any("error",err),
		)
		return &ogen.AddPolicyForbidden{
			Code:    500,
			Message: "Failed to add policy",
		}, nil
	}

	return &ogen.Policy{
		Subject: req.Subject,
		Object:  req.Object,
		Action:  req.Action,
	}, nil
}

// RemovePolicy removes an authorization policy (admin only).
func (h *Handler) RemovePolicy(ctx context.Context, req *ogen.PolicyRequest) (ogen.RemovePolicyRes, error) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.RemovePolicyUnauthorized{}, nil
	}

	// Check if user is admin
	isAdmin, err := h.rbacService.HasRole(ctx, userID, "admin")
	if err != nil {
		h.logger.Error("failed to check admin role",
			slog.String("user_id", userID.String()),
			slog.Any("error",err),
		)
		return &ogen.RemovePolicyForbidden{
			Code:    500,
			Message: "Failed to check permissions",
		}, nil
	}

	if !isAdmin {
		return &ogen.RemovePolicyForbidden{
			Code:    403,
			Message: "Admin access required",
		}, nil
	}

	if err := h.rbacService.RemovePolicy(ctx, req.Subject, req.Object, req.Action); err != nil {
		h.logger.Warn("failed to remove policy",
			slog.String("subject", req.Subject),
			slog.String("object", req.Object),
			slog.String("action", req.Action),
			slog.Any("error",err),
		)
		return &ogen.RemovePolicyNotFound{}, nil
	}

	return &ogen.RemovePolicyNoContent{}, nil
}

// GetUserRoles gets all roles assigned to a user.
func (h *Handler) GetUserRoles(ctx context.Context, params ogen.GetUserRolesParams) (ogen.GetUserRolesRes, error) {
	_, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.Error{
			Code:    401,
			Message: "Authentication required",
		}, nil
	}

	targetUserID := params.UserId

	roles, err := h.rbacService.GetUserRoles(ctx, targetUserID)
	if err != nil {
		h.logger.Error("failed to get user roles",
			slog.String("user_id", targetUserID.String()),
			slog.Any("error",err),
		)
		return &ogen.Error{
			Code:    500,
			Message: "Failed to get user roles",
		}, nil
	}

	return &ogen.RoleListResponse{
		Roles: roles,
		Total: int64(len(roles)),
	}, nil
}

// AssignRole assigns a role to a user (admin only).
func (h *Handler) AssignRole(ctx context.Context, req *ogen.AssignRoleRequest, params ogen.AssignRoleParams) (ogen.AssignRoleRes, error) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.AssignRoleUnauthorized{}, nil
	}

	// Check if user is admin
	isAdmin, err := h.rbacService.HasRole(ctx, userID, "admin")
	if err != nil {
		h.logger.Error("failed to check admin role",
			slog.String("user_id", userID.String()),
			slog.Any("error",err),
		)
		return &ogen.AssignRoleForbidden{
			Code:    500,
			Message: "Failed to check permissions",
		}, nil
	}

	if !isAdmin {
		return &ogen.AssignRoleForbidden{
			Code:    403,
			Message: "Admin access required",
		}, nil
	}

	targetUserID := params.UserId

	if err := h.rbacService.AssignRole(ctx, targetUserID, req.Role); err != nil {
		h.logger.Error("failed to assign role",
			slog.String("user_id", targetUserID.String()),
			slog.String("role", req.Role),
			slog.Any("error",err),
		)
		return &ogen.AssignRoleForbidden{
			Code:    500,
			Message: "Failed to assign role",
		}, nil
	}

	// Return updated role list for the target user
	roles, err := h.rbacService.GetUserRoles(ctx, targetUserID)
	if err != nil {
		h.logger.Error("failed to get user roles after assign",
			slog.String("user_id", targetUserID.String()),
			slog.Any("error", err),
		)
		return &ogen.RoleListResponse{
			Roles: []string{req.Role},
			Total: 1,
		}, nil
	}

	return &ogen.RoleListResponse{
		Roles: roles,
		Total: int64(len(roles)),
	}, nil
}

// RemoveRole removes a role from a user (admin only).
func (h *Handler) RemoveRole(ctx context.Context, params ogen.RemoveRoleParams) (ogen.RemoveRoleRes, error) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.RemoveRoleUnauthorized{}, nil
	}

	// Check if user is admin
	isAdmin, err := h.rbacService.HasRole(ctx, userID, "admin")
	if err != nil {
		h.logger.Error("failed to check admin role",
			slog.String("user_id", userID.String()),
			slog.Any("error",err),
		)
		return &ogen.RemoveRoleForbidden{
			Code:    500,
			Message: "Failed to check permissions",
		}, nil
	}

	if !isAdmin {
		return &ogen.RemoveRoleForbidden{
			Code:    403,
			Message: "Admin access required",
		}, nil
	}

	targetUserID := params.UserId

	if err := h.rbacService.RemoveRole(ctx, targetUserID, params.Role); err != nil {
		h.logger.Warn("failed to remove role",
			slog.String("user_id", targetUserID.String()),
			slog.String("role", params.Role),
			slog.Any("error",err),
		)
		return &ogen.RemoveRoleNotFound{}, nil
	}

	return &ogen.RemoveRoleNoContent{}, nil
}

// ListRoles lists all available roles with their permissions (admin only).
func (h *Handler) ListRoles(ctx context.Context) (ogen.ListRolesRes, error) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.ListRolesUnauthorized{}, nil
	}

	// Check if user is admin
	isAdmin, err := h.rbacService.HasRole(ctx, userID, "admin")
	if err != nil {
		h.logger.Error("failed to check admin role",
			slog.String("user_id", userID.String()),
			slog.Any("error",err),
		)
		return &ogen.ListRolesForbidden{
			Code:    500,
			Message: "Failed to check permissions",
		}, nil
	}

	if !isAdmin {
		return &ogen.ListRolesForbidden{
			Code:    403,
			Message: "Admin access required",
		}, nil
	}

	roles, err := h.rbacService.ListRoles(ctx)
	if err != nil {
		h.logger.Error("failed to list roles", slog.Any("error",err))
		return &ogen.ListRolesForbidden{
			Code:    500,
			Message: "Failed to list roles",
		}, nil
	}

	roleDetails := make([]ogen.RoleDetail, 0, len(roles))
	for _, r := range roles {
		permissions := make([]ogen.Permission, 0, len(r.Permissions))
		for _, p := range r.Permissions {
			permissions = append(permissions, ogen.Permission{
				Resource: p.Resource,
				Action:   p.Action,
			})
		}

		detail := ogen.RoleDetail{
			Name:        r.Name,
			Permissions: permissions,
			IsBuiltIn:   r.IsBuiltIn,
			UserCount:   r.UserCount,
		}
		if r.Description != "" {
			detail.Description = ogen.NewOptString(r.Description)
		}
		roleDetails = append(roleDetails, detail)
	}

	return &ogen.RolesResponse{
		Roles: roleDetails,
		Total: int64(len(roleDetails)),
	}, nil
}

// GetRole gets a specific role with its permissions (admin only).
func (h *Handler) GetRole(ctx context.Context, params ogen.GetRoleParams) (ogen.GetRoleRes, error) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.GetRoleUnauthorized{}, nil
	}

	// Check if user is admin
	isAdmin, err := h.rbacService.HasRole(ctx, userID, "admin")
	if err != nil {
		h.logger.Error("failed to check admin role",
			slog.String("user_id", userID.String()),
			slog.Any("error",err),
		)
		return &ogen.GetRoleForbidden{
			Code:    500,
			Message: "Failed to check permissions",
		}, nil
	}

	if !isAdmin {
		return &ogen.GetRoleForbidden{
			Code:    403,
			Message: "Admin access required",
		}, nil
	}

	role, err := h.rbacService.GetRole(ctx, params.RoleName)
	if err != nil {
		if errors.Is(err, rbac.ErrRoleNotFound) {
			return &ogen.GetRoleNotFound{}, nil
		}
		h.logger.Error("failed to get role",
			slog.String("role", params.RoleName),
			slog.Any("error",err),
		)
		return &ogen.GetRoleForbidden{
			Code:    500,
			Message: "Failed to get role",
		}, nil
	}

	permissions := make([]ogen.Permission, 0, len(role.Permissions))
	for _, p := range role.Permissions {
		permissions = append(permissions, ogen.Permission{
			Resource: p.Resource,
			Action:   p.Action,
		})
	}

	detail := ogen.RoleDetail{
		Name:        role.Name,
		Permissions: permissions,
		IsBuiltIn:   role.IsBuiltIn,
		UserCount:   role.UserCount,
	}
	if role.Description != "" {
		detail.Description = ogen.NewOptString(role.Description)
	}

	return &detail, nil
}

// CreateRole creates a new custom role (admin only).
func (h *Handler) CreateRole(ctx context.Context, req *ogen.CreateRoleRequest) (ogen.CreateRoleRes, error) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.CreateRoleUnauthorized{}, nil
	}

	// Check if user is admin
	isAdmin, err := h.rbacService.HasRole(ctx, userID, "admin")
	if err != nil {
		h.logger.Error("failed to check admin role",
			slog.String("user_id", userID.String()),
			slog.Any("error",err),
		)
		return &ogen.CreateRoleForbidden{
			Code:    500,
			Message: "Failed to check permissions",
		}, nil
	}

	if !isAdmin {
		return &ogen.CreateRoleForbidden{
			Code:    403,
			Message: "Admin access required",
		}, nil
	}

	// Convert permissions
	permissions := make([]rbac.Permission, 0, len(req.Permissions))
	for _, p := range req.Permissions {
		permissions = append(permissions, rbac.Permission{
			Resource: p.Resource,
			Action:   p.Action,
		})
	}

	description := ""
	if req.Description.IsSet() {
		description = req.Description.Value
	}

	role, err := h.rbacService.CreateRole(ctx, req.Name, description, permissions)
	if err != nil {
		if errors.Is(err, rbac.ErrRoleAlreadyExists) {
			return &ogen.CreateRoleConflict{
				Code:    409,
				Message: "Role already exists",
			}, nil
		}
		h.logger.Error("failed to create role",
			slog.String("role", req.Name),
			slog.Any("error",err),
		)
		return &ogen.CreateRoleBadRequest{
			Code:    400,
			Message: err.Error(),
		}, nil
	}

	respPermissions := make([]ogen.Permission, 0, len(role.Permissions))
	for _, p := range role.Permissions {
		respPermissions = append(respPermissions, ogen.Permission{
			Resource: p.Resource,
			Action:   p.Action,
		})
	}

	detail := ogen.RoleDetail{
		Name:        role.Name,
		Permissions: respPermissions,
		IsBuiltIn:   role.IsBuiltIn,
		UserCount:   role.UserCount,
	}
	if role.Description != "" {
		detail.Description = ogen.NewOptString(role.Description)
	}

	return &detail, nil
}

// DeleteRole deletes a custom role (admin only, cannot delete built-in roles).
func (h *Handler) DeleteRole(ctx context.Context, params ogen.DeleteRoleParams) (ogen.DeleteRoleRes, error) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.DeleteRoleUnauthorized{}, nil
	}

	// Check if user is admin
	isAdmin, err := h.rbacService.HasRole(ctx, userID, "admin")
	if err != nil {
		h.logger.Error("failed to check admin role",
			slog.String("user_id", userID.String()),
			slog.Any("error",err),
		)
		return &ogen.DeleteRoleForbidden{
			Code:    500,
			Message: "Failed to check permissions",
		}, nil
	}

	if !isAdmin {
		return &ogen.DeleteRoleForbidden{
			Code:    403,
			Message: "Admin access required",
		}, nil
	}

	err = h.rbacService.DeleteRole(ctx, params.RoleName)
	if err != nil {
		if errors.Is(err, rbac.ErrRoleNotFound) {
			return &ogen.DeleteRoleNotFound{}, nil
		}
		if errors.Is(err, rbac.ErrBuiltInRole) {
			return &ogen.DeleteRoleBadRequest{
				Code:    400,
				Message: "Cannot delete built-in role",
			}, nil
		}
		if errors.Is(err, rbac.ErrRoleInUse) {
			return &ogen.DeleteRoleBadRequest{
				Code:    400,
				Message: "Cannot delete role that is assigned to users",
			}, nil
		}
		h.logger.Error("failed to delete role",
			slog.String("role", params.RoleName),
			slog.Any("error",err),
		)
		return &ogen.DeleteRoleForbidden{
			Code:    500,
			Message: "Failed to delete role",
		}, nil
	}

	return &ogen.DeleteRoleNoContent{}, nil
}

// UpdateRolePermissions updates all permissions for a role (admin only).
func (h *Handler) UpdateRolePermissions(ctx context.Context, req *ogen.UpdatePermissionsRequest, params ogen.UpdateRolePermissionsParams) (ogen.UpdateRolePermissionsRes, error) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.UpdateRolePermissionsUnauthorized{}, nil
	}

	// Check if user is admin
	isAdmin, err := h.rbacService.HasRole(ctx, userID, "admin")
	if err != nil {
		h.logger.Error("failed to check admin role",
			slog.String("user_id", userID.String()),
			slog.Any("error",err),
		)
		return &ogen.UpdateRolePermissionsForbidden{
			Code:    500,
			Message: "Failed to check permissions",
		}, nil
	}

	if !isAdmin {
		return &ogen.UpdateRolePermissionsForbidden{
			Code:    403,
			Message: "Admin access required",
		}, nil
	}

	// Convert permissions
	permissions := make([]rbac.Permission, 0, len(req.Permissions))
	for _, p := range req.Permissions {
		permissions = append(permissions, rbac.Permission{
			Resource: p.Resource,
			Action:   p.Action,
		})
	}

	role, err := h.rbacService.UpdateRolePermissions(ctx, params.RoleName, permissions)
	if err != nil {
		if errors.Is(err, rbac.ErrRoleNotFound) {
			return &ogen.UpdateRolePermissionsNotFound{}, nil
		}
		h.logger.Error("failed to update role permissions",
			slog.String("role", params.RoleName),
			slog.Any("error",err),
		)
		return &ogen.UpdateRolePermissionsBadRequest{
			Code:    400,
			Message: err.Error(),
		}, nil
	}

	respPermissions := make([]ogen.Permission, 0, len(role.Permissions))
	for _, p := range role.Permissions {
		respPermissions = append(respPermissions, ogen.Permission{
			Resource: p.Resource,
			Action:   p.Action,
		})
	}

	detail := ogen.RoleDetail{
		Name:        role.Name,
		Permissions: respPermissions,
		IsBuiltIn:   role.IsBuiltIn,
		UserCount:   role.UserCount,
	}
	if role.Description != "" {
		detail.Description = ogen.NewOptString(role.Description)
	}

	return &detail, nil
}

// ListPermissions lists all available permission combinations (admin only).
func (h *Handler) ListPermissions(ctx context.Context) (ogen.ListPermissionsRes, error) {
	userID, ok := h.getUserID(ctx)
	if !ok {
		return &ogen.ListPermissionsUnauthorized{}, nil
	}

	// Check if user is admin
	isAdmin, err := h.rbacService.HasRole(ctx, userID, "admin")
	if err != nil {
		h.logger.Error("failed to check admin role",
			slog.String("user_id", userID.String()),
			slog.Any("error",err),
		)
		return &ogen.ListPermissionsForbidden{
			Code:    500,
			Message: "Failed to check permissions",
		}, nil
	}

	if !isAdmin {
		return &ogen.ListPermissionsForbidden{
			Code:    403,
			Message: "Admin access required",
		}, nil
	}

	perms := h.rbacService.ListPermissions(ctx)

	permissions := make([]ogen.Permission, 0, len(perms))
	for _, p := range perms {
		permissions = append(permissions, ogen.Permission{
			Resource: p.Resource,
			Action:   p.Action,
		})
	}

	return &ogen.PermissionsResponse{
		Permissions: permissions,
		Total:       int64(len(permissions)),
	}, nil
}

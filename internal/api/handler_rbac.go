package api

import (
	"context"

	"go.uber.org/zap"

	"github.com/lusoris/revenge/internal/api/ogen"
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
			zap.String("user_id", userID.String()),
			zap.Error(err),
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
			zap.Error(err),
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
			zap.String("user_id", userID.String()),
			zap.Error(err),
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
			zap.String("subject", req.Subject),
			zap.String("object", req.Object),
			zap.String("action", req.Action),
			zap.Error(err),
		)
		return &ogen.AddPolicyForbidden{
			Code:    500,
			Message: "Failed to add policy",
		}, nil
	}

	return &ogen.AddPolicyCreated{}, nil
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
			zap.String("user_id", userID.String()),
			zap.Error(err),
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
			zap.String("subject", req.Subject),
			zap.String("object", req.Object),
			zap.String("action", req.Action),
			zap.Error(err),
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
			zap.String("user_id", targetUserID.String()),
			zap.Error(err),
		)
		return &ogen.Error{
			Code:    500,
			Message: "Failed to get user roles",
		}, nil
	}

	return &ogen.RoleListResponse{
		Roles: roles,
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
			zap.String("user_id", userID.String()),
			zap.Error(err),
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
			zap.String("user_id", targetUserID.String()),
			zap.String("role", req.Role),
			zap.Error(err),
		)
		return &ogen.AssignRoleForbidden{
			Code:    500,
			Message: "Failed to assign role",
		}, nil
	}

	return &ogen.AssignRoleCreated{}, nil
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
			zap.String("user_id", userID.String()),
			zap.Error(err),
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
			zap.String("user_id", targetUserID.String()),
			zap.String("role", params.Role),
			zap.Error(err),
		)
		return &ogen.RemoveRoleNotFound{}, nil
	}

	return &ogen.RemoveRoleNoContent{}, nil
}

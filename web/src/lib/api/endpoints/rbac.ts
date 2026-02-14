import { get, post, put, del } from '../client';
import type {
	RolesResponse,
	RoleDetail,
	CreateRoleRequest,
	Permission,
	PolicyListResponse
} from '../types';

// ─── Roles ───────────────────────────────────────────────────────────────────

export async function listRoles(): Promise<RolesResponse> {
	return get<RolesResponse>('/v1/rbac/roles');
}

export async function getRole(name: string): Promise<RoleDetail> {
	return get<RoleDetail>(`/v1/rbac/roles/${name}`);
}

export async function createRole(data: CreateRoleRequest): Promise<RoleDetail> {
	return post<RoleDetail>('/v1/rbac/roles', data);
}

export async function updateRole(
	name: string,
	data: Partial<CreateRoleRequest>
): Promise<RoleDetail> {
	return put<RoleDetail>(`/v1/rbac/roles/${name}`, data);
}

export async function deleteRole(name: string): Promise<void> {
	return del(`/v1/rbac/roles/${name}`);
}

// ─── Role Permissions ────────────────────────────────────────────────────────

export async function getRolePermissions(name: string): Promise<Permission[]> {
	const res = await get<{ permissions: Permission[] }>(`/v1/rbac/roles/${name}/permissions`);
	return res.permissions ?? [];
}

export async function updateRolePermissions(
	name: string,
	permissions: Permission[]
): Promise<void> {
	return put(`/v1/rbac/roles/${name}/permissions`, { permissions });
}

// ─── User Roles ──────────────────────────────────────────────────────────────

export async function getUserRoles(userId: string): Promise<string[]> {
	const res = await get<{ roles: string[] }>(`/v1/rbac/users/${userId}/roles`);
	return res.roles ?? [];
}

export async function assignUserRole(userId: string, role: string): Promise<void> {
	return post(`/v1/rbac/users/${userId}/roles`, { role });
}

export async function removeUserRole(userId: string, role: string): Promise<void> {
	return del(`/v1/rbac/users/${userId}/roles/${role}`);
}

// ─── Policies ────────────────────────────────────────────────────────────────

export async function listPolicies(): Promise<PolicyListResponse> {
	return get<PolicyListResponse>('/v1/rbac/policies');
}

export async function addPolicy(
	subject: string,
	object: string,
	action: string
): Promise<void> {
	return post('/v1/rbac/policies', { subject, object, action });
}

export async function removePolicy(
	subject: string,
	object: string,
	action: string
): Promise<void> {
	return del('/v1/rbac/policies', {
		params: { subject, object, action }
	});
}

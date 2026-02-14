import { get, post, put, del, patch } from '../client';
import type {
	User,
	AdminUserListResponse,
	AdminOIDCProvider,
	AdminOIDCProviderListResponse,
	CreateOIDCProviderRequest,
	ActivityLogListResponse,
	ActivityStats,
	RadarrStatus,
	SonarrStatus,
	QualityProfile,
	RootFolder
} from '../types';

// ─── Admin Users ─────────────────────────────────────────────────────────────

export async function adminListUsers(params?: {
	page?: number;
	page_size?: number;
	search?: string;
}): Promise<AdminUserListResponse> {
	return get<AdminUserListResponse>('/v1/admin/users', {
		params: params as Record<string, string>
	});
}

export async function adminGetUser(id: string): Promise<User> {
	return get<User>(`/v1/admin/users/${id}`);
}

export async function adminUpdateUser(
	id: string,
	data: Partial<User>
): Promise<User> {
	return put<User>(`/v1/admin/users/${id}`, data);
}

export async function adminDeleteUser(id: string): Promise<void> {
	return del(`/v1/admin/users/${id}`);
}

// ─── Admin OIDC Providers ────────────────────────────────────────────────────

export async function adminListOIDCProviders(): Promise<AdminOIDCProviderListResponse> {
	return get<AdminOIDCProviderListResponse>('/v1/admin/oidc/providers');
}

export async function adminGetOIDCProvider(id: string): Promise<AdminOIDCProvider> {
	return get<AdminOIDCProvider>(`/v1/admin/oidc/providers/${id}`);
}

export async function adminCreateOIDCProvider(
	data: CreateOIDCProviderRequest
): Promise<AdminOIDCProvider> {
	return post<AdminOIDCProvider>('/v1/admin/oidc/providers', data);
}

export async function adminUpdateOIDCProvider(
	id: string,
	data: Partial<CreateOIDCProviderRequest>
): Promise<AdminOIDCProvider> {
	return put<AdminOIDCProvider>(`/v1/admin/oidc/providers/${id}`, data);
}

export async function adminDeleteOIDCProvider(id: string): Promise<void> {
	return del(`/v1/admin/oidc/providers/${id}`);
}

// ─── Admin Activity ──────────────────────────────────────────────────────────

export async function adminGetActivityLogs(params?: {
	page?: number;
	page_size?: number;
	action?: string;
	user_id?: string;
}): Promise<ActivityLogListResponse> {
	return get<ActivityLogListResponse>('/v1/admin/activity', {
		params: params as Record<string, string>
	});
}

export async function adminGetActivityStats(): Promise<ActivityStats> {
	return get<ActivityStats>('/v1/admin/activity/stats');
}

// ─── Admin Integrations — Radarr ─────────────────────────────────────────────

export async function radarrGetStatus(): Promise<RadarrStatus> {
	return get<RadarrStatus>('/v1/admin/integrations/radarr/status');
}

export async function radarrSync(): Promise<{ message: string; status: string }> {
	return post('/v1/admin/integrations/radarr/sync');
}

export async function radarrGetQualityProfiles(): Promise<QualityProfile[]> {
	const res = await get<{ profiles: QualityProfile[] }>(
		'/v1/admin/integrations/radarr/quality-profiles'
	);
	return res.profiles ?? [];
}

export async function radarrGetRootFolders(): Promise<RootFolder[]> {
	const res = await get<{ folders: RootFolder[] }>(
		'/v1/admin/integrations/radarr/root-folders'
	);
	return res.folders ?? [];
}

// ─── Admin Integrations — Sonarr ─────────────────────────────────────────────

export async function sonarrGetStatus(): Promise<SonarrStatus> {
	return get<SonarrStatus>('/v1/admin/integrations/sonarr/status');
}

export async function sonarrSync(): Promise<{ message: string; status: string }> {
	return post('/v1/admin/integrations/sonarr/sync');
}

export async function sonarrGetQualityProfiles(): Promise<QualityProfile[]> {
	const res = await get<{ profiles: QualityProfile[] }>(
		'/v1/admin/integrations/sonarr/quality-profiles'
	);
	return res.profiles ?? [];
}

export async function sonarrGetRootFolders(): Promise<RootFolder[]> {
	const res = await get<{ folders: RootFolder[] }>(
		'/v1/admin/integrations/sonarr/root-folders'
	);
	return res.folders ?? [];
}

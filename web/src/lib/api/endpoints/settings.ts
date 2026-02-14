import { get, put } from '../client';
import type { ServerSetting, UserSetting } from '../types';

// ─── Server Settings (admin) ─────────────────────────────────────────────────

export async function getServerSettings(): Promise<ServerSetting[]> {
	const res = await get<{ settings: ServerSetting[] }>('/v1/settings/server');
	return res.settings ?? [];
}

export async function getServerSetting(key: string): Promise<ServerSetting> {
	return get<ServerSetting>(`/v1/settings/server/${key}`);
}

export async function updateServerSetting(
	key: string,
	value: string | number | boolean | Record<string, unknown>
): Promise<ServerSetting> {
	return put<ServerSetting>(`/v1/settings/server/${key}`, { value });
}

// ─── User Settings ───────────────────────────────────────────────────────────

export async function getUserSettings(): Promise<UserSetting[]> {
	const res = await get<{ settings: UserSetting[] }>('/v1/settings/user');
	return res.settings ?? [];
}

export async function getUserSetting(key: string): Promise<UserSetting> {
	return get<UserSetting>(`/v1/settings/user/${key}`);
}

export async function updateUserSetting(
	key: string,
	value: string | number | boolean | Record<string, unknown>
): Promise<UserSetting> {
	return put<UserSetting>(`/v1/settings/user/${key}`, { value });
}

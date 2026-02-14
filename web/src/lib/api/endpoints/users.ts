import { get, put, post, del } from '../client';
import type {
	User,
	UserUpdate,
	UserPreferences,
	Avatar,
	ChangePasswordRequest,
	OIDCUserLink,
	OIDCUserLinkListResponse
} from '../types';

export async function getCurrentUser(): Promise<User> {
	return get<User>('/v1/users/me');
}

export async function updateCurrentUser(data: UserUpdate): Promise<User> {
	return put<User>('/v1/users/me', data);
}

export async function getPreferences(): Promise<UserPreferences> {
	return get<UserPreferences>('/v1/users/me/preferences');
}

export async function updatePreferences(data: Partial<UserPreferences>): Promise<UserPreferences> {
	return put<UserPreferences>('/v1/users/me/preferences', data);
}

export async function uploadAvatar(file: File): Promise<Avatar> {
	const formData = new FormData();
	formData.append('avatar', file);
	return post<Avatar>('/v1/users/me/avatar', formData);
}

export async function deleteAvatar(): Promise<void> {
	return del('/v1/users/me/avatar');
}

export async function changePassword(data: ChangePasswordRequest): Promise<void> {
	return post('/v1/auth/change-password', data);
}

export async function forgotPassword(email: string): Promise<void> {
	return post('/v1/auth/forgot-password', { email });
}

export async function resetPassword(token: string, newPassword: string): Promise<void> {
	return post('/v1/auth/reset-password', { token, new_password: newPassword });
}

export async function listOIDCLinks(): Promise<OIDCUserLink[]> {
	const res = await get<OIDCUserLinkListResponse>('/v1/users/me/oidc');
	return res.links ?? [];
}

export async function unlinkOIDC(linkId: string): Promise<void> {
	return del(`/v1/users/me/oidc/${linkId}`);
}

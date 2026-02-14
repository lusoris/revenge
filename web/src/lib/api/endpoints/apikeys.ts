import { get, post, del } from '../client';
import type {
	APIKeyInfo,
	APIKeyListResponse,
	CreateAPIKeyRequest,
	CreateAPIKeyResponse
} from '../types';

export async function listAPIKeys(): Promise<APIKeyListResponse> {
	return get<APIKeyListResponse>('/v1/apikeys');
}

export async function getAPIKey(id: string): Promise<APIKeyInfo> {
	return get<APIKeyInfo>(`/v1/apikeys/${id}`);
}

export async function createAPIKey(data: CreateAPIKeyRequest): Promise<CreateAPIKeyResponse> {
	return post<CreateAPIKeyResponse>('/v1/apikeys', data);
}

export async function revokeAPIKey(id: string): Promise<void> {
	return del(`/v1/apikeys/${id}`);
}

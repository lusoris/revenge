import { get, post, setTokens, clearTokens } from '../client';
import type {
	LoginRequest,
	LoginResponse,
	RegisterRequest,
	User,
	OIDCProvider,
	OIDCCallbackResponse
} from '../types';

export async function login(req: LoginRequest): Promise<LoginResponse> {
	const res = await post<LoginResponse>('/v1/auth/login', req, { noAuth: true });
	setTokens(res.access_token, res.refresh_token, res.expires_in);
	return res;
}

export async function register(req: RegisterRequest): Promise<User> {
	return post<User>('/v1/auth/register', req, { noAuth: true });
}

export async function logout(): Promise<void> {
	try {
		await post<void>('/v1/auth/logout');
	} finally {
		clearTokens();
	}
}

export async function getCurrentUser(): Promise<User> {
	return get<User>('/v1/users/me');
}

export async function listOIDCProviders(): Promise<OIDCProvider[]> {
	const res = await get<{ providers: OIDCProvider[] }>('/v1/oidc/providers', { noAuth: true });
	return res.providers ?? [];
}

export async function handleOIDCCallback(
	provider: string,
	code: string,
	state: string
): Promise<OIDCCallbackResponse> {
	const res = await get<OIDCCallbackResponse>(`/v1/oidc/callback/${provider}`, {
		noAuth: true,
		params: { code, state }
	});
	if (res.access_token && res.refresh_token) {
		setTokens(res.access_token, res.refresh_token, res.expires_in);
	}
	return res;
}

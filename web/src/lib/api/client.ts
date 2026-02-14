// Typed fetch wrapper with auth, token refresh, and error handling.

import type { APIError } from './types';

// ─── Token State (module-level singleton) ────────────────────────────────────

let accessToken: string | null = null;
let refreshToken: string | null = null;
let tokenExpiresAt: number = 0;
let refreshPromise: Promise<boolean> | null = null;

// Subscribers notified on auth state changes
type AuthChangeCallback = (authenticated: boolean) => void;
const authChangeCallbacks: Set<AuthChangeCallback> = new Set();

export function onAuthChange(cb: AuthChangeCallback): () => void {
	authChangeCallbacks.add(cb);
	return () => authChangeCallbacks.delete(cb);
}

function notifyAuthChange(authenticated: boolean) {
	for (const cb of authChangeCallbacks) cb(authenticated);
}

export function setTokens(access: string, refresh: string, expiresIn: number) {
	accessToken = access;
	refreshToken = refresh;
	// Expire 30s early to avoid edge-case 401s
	tokenExpiresAt = Date.now() + (expiresIn - 30) * 1000;
	persistTokens();
	notifyAuthChange(true);
}

export function clearTokens() {
	accessToken = null;
	refreshToken = null;
	tokenExpiresAt = 0;
	localStorage.removeItem('revenge_refresh_token');
	notifyAuthChange(false);
}

export function getAccessToken(): string | null {
	return accessToken;
}

export function isAuthenticated(): boolean {
	return accessToken !== null;
}

function persistTokens() {
	if (refreshToken) {
		localStorage.setItem('revenge_refresh_token', refreshToken);
	}
}

export function restoreTokens(): string | null {
	return localStorage.getItem('revenge_refresh_token');
}

// ─── API Base ────────────────────────────────────────────────────────────────

const API_BASE = '/api';

// ─── Error Class ─────────────────────────────────────────────────────────────

export class ApiError extends Error {
	status: number;
	code: string;
	requestId?: string;
	details?: Record<string, string>;

	constructor(status: number, body: APIError) {
		super(body.message || body.error || `HTTP ${status}`);
		this.name = 'ApiError';
		this.status = status;
		this.code = body.error;
		this.requestId = body.request_id;
		this.details = body.details;
	}
}

// ─── Token Refresh ───────────────────────────────────────────────────────────

async function doRefresh(): Promise<boolean> {
	const rt = refreshToken ?? restoreTokens();
	if (!rt) return false;

	try {
		const res = await fetch(`${API_BASE}/v1/auth/refresh`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ refresh_token: rt })
		});

		if (!res.ok) {
			clearTokens();
			return false;
		}

		const data = await res.json();
		// /auth/refresh returns only access_token + expires_in (no new refresh_token)
		// Keep existing refresh token if the response doesn't include a new one
		setTokens(data.access_token, data.refresh_token ?? rt, data.expires_in);
		return true;
	} catch {
		clearTokens();
		return false;
	}
}

async function ensureToken(): Promise<boolean> {
	if (accessToken && Date.now() < tokenExpiresAt) return true;

	// Deduplicate concurrent refresh calls
	if (!refreshPromise) {
		refreshPromise = doRefresh().finally(() => {
			refreshPromise = null;
		});
	}
	return refreshPromise;
}

// ─── Core Fetch ──────────────────────────────────────────────────────────────

interface RequestOptions {
	method?: string;
	body?: unknown;
	params?: Record<string, string | number | boolean | undefined>;
	headers?: Record<string, string>;
	/** Skip auth header (for public endpoints like login) */
	noAuth?: boolean;
	/** Raw response (don't parse JSON) */
	raw?: boolean;
}

async function request<T>(path: string, opts: RequestOptions = {}): Promise<T> {
	const { method = 'GET', body, params, headers: extraHeaders, noAuth, raw } = opts;

	// Build URL with query params
	let url = `${API_BASE}${path}`;
	if (params) {
		const sp = new URLSearchParams();
		for (const [k, v] of Object.entries(params)) {
			if (v !== undefined && v !== null) sp.set(k, String(v));
		}
		const qs = sp.toString();
		if (qs) url += `?${qs}`;
	}

	// Ensure valid token for authenticated requests
	if (!noAuth) await ensureToken();

	const headers: Record<string, string> = {
		Accept: 'application/json',
		...extraHeaders
	};

	if (accessToken && !noAuth) {
		headers['Authorization'] = `Bearer ${accessToken}`;
	}

	if (body !== undefined && !(body instanceof FormData)) {
		headers['Content-Type'] = 'application/json';
	}

	const res = await fetch(url, {
		method,
		headers,
		body: body === undefined ? undefined : body instanceof FormData ? body : JSON.stringify(body)
	});

	// Handle 204 No Content
	if (res.status === 204) return undefined as T;

	// Handle errors
	if (!res.ok) {
		// 401 and we haven't retried refresh yet
		if (res.status === 401 && !noAuth && refreshToken) {
			const refreshed = await doRefresh();
			if (refreshed) {
				// Retry the request once
				return request(path, { ...opts, noAuth: false });
			}
		}

		let errorBody: APIError;
		try {
			errorBody = await res.json();
		} catch {
			errorBody = {
				error: 'unknown',
				message: res.statusText,
				status_code: res.status
			};
		}
		throw new ApiError(res.status, errorBody);
	}

	if (raw) return res as unknown as T;
	return res.json();
}

// ─── Public HTTP Methods ─────────────────────────────────────────────────────

export function get<T>(path: string, opts?: Omit<RequestOptions, 'method'>) {
	return request<T>(path, { ...opts, method: 'GET' });
}

export function post<T>(
	path: string,
	body?: unknown,
	opts?: Omit<RequestOptions, 'method' | 'body'>
) {
	return request<T>(path, { ...opts, method: 'POST', body });
}

export function put<T>(
	path: string,
	body?: unknown,
	opts?: Omit<RequestOptions, 'method' | 'body'>
) {
	return request<T>(path, { ...opts, method: 'PUT', body });
}

export function patch<T>(
	path: string,
	body?: unknown,
	opts?: Omit<RequestOptions, 'method' | 'body'>
) {
	return request<T>(path, { ...opts, method: 'PATCH', body });
}

export function del<T = void>(path: string, opts?: Omit<RequestOptions, 'method'>) {
	return request<T>(path, { ...opts, method: 'DELETE' });
}

// ─── Image Helper ────────────────────────────────────────────────────────────

export function imageUrl(
	type: 'poster' | 'backdrop' | 'profile' | 'logo',
	size: string,
	path: string
): string {
	if (!path) return '';
	// Strip leading slash if present
	const cleanPath = path.startsWith('/') ? path.slice(1) : path;
	return `${API_BASE}/v1/images/${type}/${size}/${cleanPath}`;
}

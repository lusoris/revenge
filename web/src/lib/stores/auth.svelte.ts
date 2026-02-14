// Auth store — Svelte 5 runes-based reactive auth state.

import type { User } from '$api/types';
import {
	setTokens,
	clearTokens,
	isAuthenticated as checkAuth,
	restoreTokens,
	onAuthChange
} from '$api/client';
import * as authApi from '$api/endpoints/auth';

// ─── Reactive State ──────────────────────────────────────────────────────────

let user = $state<User | null>(null);
let loading = $state(true);
let error = $state<string | null>(null);

// ─── Init (called once from root layout) ─────────────────────────────────────

export async function initAuth(): Promise<void> {
	loading = true;
	error = null;

	// Try to restore from persisted refresh token
	const rt = restoreTokens();
	if (!rt) {
		loading = false;
		return;
	}

	try {
		// The client will auto-refresh via ensureToken() inside get()
		user = await authApi.getCurrentUser();
	} catch {
		clearTokens();
		user = null;
	} finally {
		loading = false;
	}
}

// Listen for token clears (e.g. expired refresh) to sync state
onAuthChange((authenticated) => {
	if (!authenticated) {
		user = null;
	}
});

// ─── Actions ─────────────────────────────────────────────────────────────────

export async function login(username: string, password: string, totpCode?: string): Promise<User> {
	error = null;
	loading = true;
	try {
		const res = await authApi.login({ username, password, totp_code: totpCode });
		user = res.user;
		return res.user;
	} catch (e: unknown) {
		const msg = e instanceof Error ? e.message : 'Login failed';
		error = msg;
		throw e;
	} finally {
		loading = false;
	}
}

export async function register(
	username: string,
	email: string,
	password: string
): Promise<User> {
	error = null;
	loading = true;
	try {
		const newUser = await authApi.register({ username, email, password });
		return newUser;
	} catch (e: unknown) {
		const msg = e instanceof Error ? e.message : 'Registration failed';
		error = msg;
		throw e;
	} finally {
		loading = false;
	}
}

export async function logout(): Promise<void> {
	try {
		await authApi.logout();
	} finally {
		user = null;
		error = null;
	}
}

// ─── Getters ─────────────────────────────────────────────────────────────────

export function getAuth() {
	return {
		get user() {
			return user;
		},
		get loading() {
			return loading;
		},
		get error() {
			return error;
		},
		get isAuthenticated() {
			return user !== null && checkAuth();
		},
		get isAdmin() {
			return user?.is_admin ?? false;
		}
	};
}

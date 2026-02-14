import { get, del } from '../client';
import type { SessionInfo, SessionListResponse } from '../types';

export async function listSessions(): Promise<SessionListResponse> {
	return get<SessionListResponse>('/v1/sessions');
}

export async function getCurrentSession(): Promise<SessionInfo> {
	return get<SessionInfo>('/v1/sessions/current');
}

export async function revokeSession(sessionId: string): Promise<void> {
	return del(`/v1/sessions/${sessionId}`);
}

export async function revokeAllOtherSessions(): Promise<void> {
	return del('/v1/sessions');
}

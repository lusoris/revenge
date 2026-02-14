import { post, get, del } from '../client';
import type { PlaybackSession, StartPlaybackRequest } from '../types';

export async function startPlayback(req: StartPlaybackRequest): Promise<PlaybackSession> {
	return post<PlaybackSession>('/v1/playback/sessions', req);
}

export async function getPlaybackSession(sessionId: string): Promise<PlaybackSession> {
	return get<PlaybackSession>(`/v1/playback/sessions/${sessionId}`);
}

export async function stopPlayback(sessionId: string): Promise<void> {
	return del(`/v1/playback/sessions/${sessionId}`);
}

export async function heartbeat(
	sessionId: string,
	positionSeconds?: number
): Promise<void> {
	return post(`/v1/playback/sessions/${sessionId}/heartbeat`, {
		position_seconds: positionSeconds
	});
}

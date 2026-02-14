import { get, post, put, del } from '../client';
import type {
	TVSeries,
	TVShowListResponse,
	TVSeason,
	TVEpisode,
	TVEpisodeFile,
	TVContinueWatchingItem,
	EpisodeWatchProgress,
	CreditListResponse
} from '../types';

export interface ListTVShowsParams {
	limit?: number;
	offset?: number;
	order_by?: string;
	library_id?: string;
}

export async function listTVShows(params?: ListTVShowsParams): Promise<TVShowListResponse> {
	return get<TVShowListResponse>('/v1/tvshows', { params: params as Record<string, string> });
}

export async function getTVShow(id: string): Promise<TVSeries> {
	return get<TVSeries>(`/v1/tvshows/${id}`);
}

export async function getTVShowSeasons(id: string): Promise<TVSeason[]> {
	return get<TVSeason[]>(`/v1/tvshows/${id}/seasons`);
}

export async function getTVShowEpisodes(
	id: string,
	params?: { season_number?: number }
): Promise<TVEpisode[]> {
	return get<TVEpisode[]>(`/v1/tvshows/${id}/episodes`, {
		params: params as Record<string, string>
	});
}

export async function getTVShowCast(id: string): Promise<CreditListResponse> {
	return get<CreditListResponse>(`/v1/tvshows/${id}/cast`);
}

export async function getRecentlyAddedTVShows(params?: {
	limit?: number;
	offset?: number;
}): Promise<TVShowListResponse> {
	return get<TVShowListResponse>('/v1/tvshows/recently-added', {
		params: params as Record<string, string>
	});
}

export async function getTVContinueWatching(): Promise<TVContinueWatchingItem[]> {
	return get<TVContinueWatchingItem[]>('/v1/tvshows/continue-watching');
}

export async function getSeason(seasonId: string): Promise<TVSeason> {
	return get<TVSeason>(`/v1/tvshows/seasons/${seasonId}`);
}

export async function getSeasonEpisodes(seasonId: string): Promise<TVEpisode[]> {
	return get<TVEpisode[]>(`/v1/tvshows/seasons/${seasonId}/episodes`);
}

export async function getEpisode(episodeId: string): Promise<TVEpisode> {
	return get<TVEpisode>(`/v1/tvshows/episodes/${episodeId}`);
}

export async function getEpisodeFiles(episodeId: string): Promise<TVEpisodeFile[]> {
	return get<TVEpisodeFile[]>(`/v1/tvshows/episodes/${episodeId}/files`);
}

export async function getEpisodeProgress(episodeId: string): Promise<EpisodeWatchProgress> {
	return get<EpisodeWatchProgress>(`/v1/tvshows/episodes/${episodeId}/progress`);
}

export async function updateEpisodeProgress(
	episodeId: string,
	positionSeconds: number,
	durationSeconds: number
): Promise<EpisodeWatchProgress> {
	return put<EpisodeWatchProgress>(`/v1/tvshows/episodes/${episodeId}/progress`, {
		position_seconds: positionSeconds,
		duration_seconds: durationSeconds
	});
}

export async function deleteEpisodeProgress(episodeId: string): Promise<void> {
	return del(`/v1/tvshows/episodes/${episodeId}/progress`);
}

export async function markEpisodeWatched(episodeId: string): Promise<void> {
	return post(`/v1/tvshows/episodes/${episodeId}/watched`);
}

export async function getNextEpisode(showId: string): Promise<TVEpisode> {
	return get<TVEpisode>(`/v1/tvshows/${showId}/next-episode`);
}

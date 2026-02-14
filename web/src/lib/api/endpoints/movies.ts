import { get, post, del } from '../client';
import type {
	Movie,
	MovieListResponse,
	MovieFile,
	ContinueWatchingItem,
	MovieWatched,
	CreditListResponse
} from '../types';

export interface ListMoviesParams {
	limit?: number;
	offset?: number;
	order_by?: string;
	library_id?: string;
}

export async function listMovies(params?: ListMoviesParams): Promise<MovieListResponse> {
	return get<MovieListResponse>('/v1/movies', { params: params as Record<string, string> });
}

export async function getMovie(id: string): Promise<Movie> {
	return get<Movie>(`/v1/movies/${id}`);
}

export async function getMovieFiles(id: string): Promise<MovieFile[]> {
	return get<MovieFile[]>(`/v1/movies/${id}/files`);
}

export async function getMovieCast(id: string): Promise<CreditListResponse> {
	return get<CreditListResponse>(`/v1/movies/${id}/cast`);
}

export async function getMovieCrew(id: string): Promise<CreditListResponse> {
	return get<CreditListResponse>(`/v1/movies/${id}/crew`);
}

export async function getRecentlyAdded(params?: {
	limit?: number;
	offset?: number;
}): Promise<MovieListResponse> {
	return get<MovieListResponse>('/v1/movies/recently-added', {
		params: params as Record<string, string>
	});
}

export async function getTopRated(params?: {
	limit?: number;
	offset?: number;
	min_votes?: number;
}): Promise<MovieListResponse> {
	return get<MovieListResponse>('/v1/movies/top-rated', {
		params: params as Record<string, string>
	});
}

export async function getContinueWatching(): Promise<ContinueWatchingItem[]> {
	return get<ContinueWatchingItem[]>('/v1/movies/continue-watching');
}

export async function getWatchProgress(movieId: string): Promise<MovieWatched> {
	return get<MovieWatched>(`/v1/movies/${movieId}/progress`);
}

export async function updateWatchProgress(
	movieId: string,
	positionSeconds: number,
	durationSeconds: number
): Promise<MovieWatched> {
	return post<MovieWatched>(`/v1/movies/${movieId}/progress`, {
		position_seconds: positionSeconds,
		duration_seconds: durationSeconds
	});
}

export async function deleteWatchProgress(movieId: string): Promise<void> {
	return del(`/v1/movies/${movieId}/progress`);
}

export async function markAsWatched(movieId: string): Promise<void> {
	return post(`/v1/movies/${movieId}/watched`);
}

export async function getSimilarMovies(
	id: string,
	params?: { limit?: number }
): Promise<{ movies: Movie[] }> {
	return get(`/v1/movies/${id}/similar`, { params: params as Record<string, string> });
}

import { get } from '../client';
import type {
	SearchResults,
	TVShowSearchResults,
	MultiSearchResults,
	AutocompleteResults
} from '../types';

export interface SearchParams {
	q: string;
	page?: number;
	per_page?: number;
	sort_by?: string;
	filter_by?: string;
}

export async function searchMovies(params: SearchParams): Promise<SearchResults> {
	return get<SearchResults>('/v1/search/movies', {
		params: params as Record<string, string>
	});
}

export async function searchTVShows(params: SearchParams): Promise<TVShowSearchResults> {
	return get<TVShowSearchResults>('/v1/search/tvshows', {
		params: params as Record<string, string>
	});
}

export async function searchMulti(params: {
	q: string;
	limit?: number;
}): Promise<MultiSearchResults> {
	return get<MultiSearchResults>('/v1/search/multi', {
		params: params as Record<string, string>
	});
}

export async function autocompleteMovies(params: {
	q: string;
	limit?: number;
}): Promise<AutocompleteResults> {
	return get<AutocompleteResults>('/v1/search/movies/autocomplete', {
		params: params as Record<string, string>
	});
}

export async function autocompleteTVShows(params: {
	q: string;
	limit?: number;
}): Promise<AutocompleteResults> {
	return get<AutocompleteResults>('/v1/search/tvshows/autocomplete', {
		params: params as Record<string, string>
	});
}

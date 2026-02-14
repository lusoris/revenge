import { get } from '../client';
import type { Library, Genre } from '../types';

export async function listLibraries(): Promise<Library[]> {
	const res = await get<{ libraries: Library[] }>('/v1/libraries');
	return res.libraries ?? [];
}

export async function getLibrary(id: string): Promise<Library> {
	return get<Library>(`/v1/libraries/${id}`);
}

export async function listGenres(): Promise<Genre[]> {
	return get<Genre[]>('/v1/genres');
}

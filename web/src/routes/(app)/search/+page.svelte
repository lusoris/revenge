<script lang="ts">
	import { page } from '$app/state';
	import { createQuery } from '@tanstack/svelte-query';
	import { derived, writable } from 'svelte/store';
	import * as searchApi from '$api/endpoints/search';
	import MediaCard from '$components/media/MediaCard.svelte';
	import MediaGrid from '$components/media/MediaGrid.svelte';

	let query = $state(page.url.searchParams.get('q') ?? '');
	let debounced = $state(page.url.searchParams.get('q') ?? '');
	let timer: ReturnType<typeof setTimeout>;

	const debouncedStore = writable(debounced);
	$effect(() => {
		debouncedStore.set(debounced);
	});

	function onInput(e: Event) {
		query = (e.target as HTMLInputElement).value;
		clearTimeout(timer);
		timer = setTimeout(() => {
			debounced = query;
		}, 300);
	}

	const resultsOptions = derived(debouncedStore, ($d) => ({
		queryKey: ['search', 'multi', $d] as const,
		queryFn: () => searchApi.searchMulti({ q: $d, limit: 40 }),
		enabled: $d.length >= 2
	}));

	const autocompleteOptions = derived(debouncedStore, ($d) => ({
		queryKey: ['search', 'autocomplete', $d] as const,
		queryFn: () => searchApi.autocompleteMovies({ q: $d, limit: 8 }),
		enabled: $d.length >= 2 && $d.length < 4
	}));

	const results = createQuery(resultsOptions);
	const autocomplete = createQuery(autocompleteOptions);
</script>

<svelte:head>
	<title>Search — Revenge</title>
</svelte:head>

<div>
	<div class="mb-6">
		<input
			type="search"
			value={query}
			oninput={onInput}
			placeholder="Search movies and TV shows…"
			autofocus
			class="w-full rounded-xl border border-neutral-800 bg-neutral-900 px-4 py-3 text-base text-white placeholder-neutral-600 outline-none focus:border-neutral-600 focus:ring-1 focus:ring-neutral-600"
		/>
	</div>

	{#if debounced.length < 2}
		<div class="py-16 text-center">
			<p class="text-lg text-neutral-500">Type to search your library</p>
		</div>
	{:else if $results.isPending}
		<div class="flex justify-center py-16">
			<div class="h-8 w-8 animate-spin rounded-full border-2 border-neutral-700 border-t-white"></div>
		</div>
	{:else if $results.data}
		{@const movieHits = $results.data.movies?.hits ?? []}
		{@const tvHits = $results.data.tvshows?.hits ?? []}
		{@const totalHits = ($results.data.movies?.total_hits ?? 0) + ($results.data.tvshows?.total_hits ?? 0)}

		{#if totalHits === 0}
			<div class="py-16 text-center">
				<p class="text-lg text-neutral-400">No results for "{debounced}"</p>
			</div>
		{:else}
			<p class="mb-4 text-sm text-neutral-500">{totalHits} result{totalHits !== 1 ? 's' : ''}</p>

			{#if movieHits.length}
				<section class="mb-8">
					<h2 class="mb-3 text-lg font-semibold text-white">Movies</h2>
					<MediaGrid>
						{#each movieHits as hit (hit.document.id)}
							<MediaCard item={hit.document} type="movie" />
						{/each}
					</MediaGrid>
				</section>
			{/if}

			{#if tvHits.length}
				<section class="mb-8">
					<h2 class="mb-3 text-lg font-semibold text-white">TV Shows</h2>
					<MediaGrid>
						{#each tvHits as hit (hit.document.id)}
							<MediaCard item={hit.document} type="tvshow" />
						{/each}
					</MediaGrid>
				</section>
			{/if}
		{/if}
	{/if}
</div>

<script lang="ts">
	import { page } from '$app/state';
	import { createQuery } from '@tanstack/svelte-query';
	import * as searchApi from '$api/endpoints/search';
	import MediaCard from '$components/media/MediaCard.svelte';
	import MediaGrid from '$components/media/MediaGrid.svelte';

	let query = $state(page.url.searchParams.get('q') ?? '');
	let debounced = $state(page.url.searchParams.get('q') ?? '');
	let timer: ReturnType<typeof setTimeout>;

	function onInput(e: Event) {
		query = (e.target as HTMLInputElement).value;
		clearTimeout(timer);
		timer = setTimeout(() => {
			debounced = query;
		}, 300);
	}

	const results = createQuery(() => ({
		queryKey: ['search', 'multi', debounced],
		queryFn: () => searchApi.searchMulti({ q: debounced, limit: 40 }),
		enabled: debounced.length >= 2
	}));

	const autocomplete = createQuery(() => ({
		queryKey: ['search', 'autocomplete', debounced],
		queryFn: () => searchApi.autocompleteMovies({ q: debounced, limit: 8 }),
		enabled: debounced.length >= 2 && debounced.length < 4
	}));
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
		{@const { movies, tvshows, total_hits } = $results.data}

		{#if total_hits === 0}
			<div class="py-16 text-center">
				<p class="text-lg text-neutral-400">No results for "{debounced}"</p>
			</div>
		{:else}
			<p class="mb-4 text-sm text-neutral-500">{total_hits} result{total_hits !== 1 ? 's' : ''}</p>

			{#if movies.length}
				<section class="mb-8">
					<h2 class="mb-3 text-lg font-semibold text-white">Movies</h2>
					<MediaGrid>
						{#each movies as movie (movie.id)}
							<MediaCard item={movie} type="movie" />
						{/each}
					</MediaGrid>
				</section>
			{/if}

			{#if tvshows.length}
				<section class="mb-8">
					<h2 class="mb-3 text-lg font-semibold text-white">TV Shows</h2>
					<MediaGrid>
						{#each tvshows as show (show.id)}
							<MediaCard item={show} type="tvshow" />
						{/each}
					</MediaGrid>
				</section>
			{/if}
		{/if}
	{/if}
</div>

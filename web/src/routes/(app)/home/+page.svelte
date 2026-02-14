<script lang="ts">
	import * as moviesApi from '$api/endpoints/movies';
	import * as tvshowsApi from '$api/endpoints/tvshows';
	import MediaCard from '$components/media/MediaCard.svelte';
	import MediaRow from '$components/media/MediaRow.svelte';
	import { getAuth } from '$lib/stores/auth.svelte';
	import { createQuery } from '@tanstack/svelte-query';

	const auth = getAuth();

	const continueWatchingMovies = createQuery(() => ({
		queryKey: ['movies', 'continue-watching'],
		queryFn: () => moviesApi.getContinueWatching()
	}));

	const continueWatchingTV = createQuery(() => ({
		queryKey: ['tvshows', 'continue-watching'],
		queryFn: () => tvshowsApi.getTVContinueWatching()
	}));

	const recentMovies = createQuery(() => ({
		queryKey: ['movies', 'recently-added'],
		queryFn: () => moviesApi.getRecentlyAdded({ limit: 20 })
	}));

	const recentTVShows = createQuery(() => ({
		queryKey: ['tvshows', 'recently-added'],
		queryFn: () => tvshowsApi.getRecentlyAddedTVShows({ limit: 20 })
	}));
</script>

<svelte:head>
	<title>Home — Revenge</title>
</svelte:head>

<div class="space-y-6">
	<h1 class="text-2xl font-bold text-white">
		Welcome back, {auth.user?.display_name ?? auth.user?.username ?? 'User'}
	</h1>

	<!-- Continue Watching Movies -->
	{#if continueWatchingMovies.data?.length}
		<section class="mb-8">
			<h2 class="mb-3 text-lg font-semibold text-white">Continue Watching</h2>
			<div class="flex gap-3 overflow-x-auto pb-2">
				{#each continueWatchingMovies.data as cw (cw.movie.id)}
					<div class="w-36 flex-shrink-0 sm:w-40">
						<MediaCard item={cw.movie} type="movie" progress={cw.progress.percentage} />
					</div>
				{/each}
			</div>
		</section>
	{/if}

	<!-- Continue Watching TV -->
	{#if continueWatchingTV.data?.length}
		<section class="mb-8">
			<h2 class="mb-3 text-lg font-semibold text-white">Continue Watching — TV Shows</h2>
			<div class="flex gap-3 overflow-x-auto pb-2">
				{#each continueWatchingTV.data as cw (cw.series.id)}
					<div class="w-36 flex-shrink-0 sm:w-40">
						<MediaCard item={cw.series} type="tvshow" progress={cw.progress.percentage} />
					</div>
				{/each}
			</div>
		</section>
	{/if}

	<!-- Recently Added Movies -->
	{#if recentMovies.data}
		<MediaRow
			title="Recently Added Movies"
			items={recentMovies.data.items}
			type="movie"
			href="/movies"
		/>
	{/if}

	<!-- Recently Added TV Shows -->
	{#if recentTVShows.data}
		<MediaRow
			title="Recently Added TV Shows"
			items={recentTVShows.data.items}
			type="tvshow"
			href="/tvshows"
		/>
	{/if}

	<!-- Loading states -->
	{#if recentMovies.isPending && recentTVShows.isPending}
		<div class="flex justify-center py-12">
			<div
				class="h-8 w-8 animate-spin rounded-full border-2 border-neutral-700 border-t-white"
			></div>
		</div>
	{/if}

	<!-- Empty state -->
	{#if !recentMovies.isPending && !recentTVShows.isPending && !recentMovies.data?.items?.length && !recentTVShows.data?.items?.length}
		<div class="py-16 text-center">
			<p class="text-lg text-neutral-400">Your library is empty.</p>
			<p class="mt-1 text-sm text-neutral-600">
				Add media to a library to get started.
			</p>
		</div>
	{/if}
</div>

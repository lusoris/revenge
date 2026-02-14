<script lang="ts">
	import { createInfiniteQuery } from '@tanstack/svelte-query';
	import MediaCard from '$components/media/MediaCard.svelte';
	import MediaGrid from '$components/media/MediaGrid.svelte';
	import * as tvshowsApi from '$api/endpoints/tvshows';

	const PAGE_SIZE = 30;

	let sortBy = $state('created_at');

	const query = createInfiniteQuery(() => ({
		queryKey: ['tvshows', 'list', sortBy],
		queryFn: ({ pageParam = 0 }) =>
			tvshowsApi.listTVShows({ limit: PAGE_SIZE, offset: pageParam, order_by: sortBy }),
		getNextPageParam: (lastPage, allPages) => {
			const fetched = allPages.reduce((n, p) => n + p.tvshows.length, 0);
			return fetched < lastPage.total ? fetched : undefined;
		},
		initialPageParam: 0
	}));

	function loadMore() {
		if ($query.hasNextPage && !$query.isFetchingNextPage) {
			$query.fetchNextPage();
		}
	}

	const allShows = $derived(($query.data?.pages ?? []).flatMap((p) => p.tvshows));
	const total = $derived($query.data?.pages?.[0]?.total ?? 0);
</script>

<svelte:head>
	<title>TV Shows — Revenge</title>
</svelte:head>

<div>
	<div class="mb-6 flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold text-white">TV Shows</h1>
			{#if total > 0}
				<p class="text-sm text-neutral-500">{total} shows</p>
			{/if}
		</div>

		<select
			bind:value={sortBy}
			class="rounded-lg border border-neutral-800 bg-neutral-900 px-3 py-1.5 text-sm text-neutral-300 outline-none"
		>
			<option value="created_at">Recently Added</option>
			<option value="name">Title</option>
			<option value="first_air_date">First Aired</option>
			<option value="vote_average">Rating</option>
		</select>
	</div>

	{#if $query.isPending}
		<div class="flex justify-center py-16">
			<div
				class="h-8 w-8 animate-spin rounded-full border-2 border-neutral-700 border-t-white"
			></div>
		</div>
	{:else if $query.isError}
		<div class="py-16 text-center text-red-400">
			Failed to load TV shows. {$query.error?.message}
		</div>
	{:else}
		<MediaGrid>
			{#each allShows as show (show.id)}
				<MediaCard item={show} type="tvshow" />
			{/each}
		</MediaGrid>

		{#if $query.hasNextPage}
			<div class="mt-8 flex justify-center">
				<button
					onclick={loadMore}
					disabled={$query.isFetchingNextPage}
					class="rounded-lg border border-neutral-800 bg-neutral-900 px-6 py-2 text-sm text-neutral-300 transition-colors hover:bg-neutral-800 disabled:opacity-50"
				>
					{$query.isFetchingNextPage ? 'Loading…' : 'Load more'}
				</button>
			</div>
		{/if}
	{/if}
</div>

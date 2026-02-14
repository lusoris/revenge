<script lang="ts">
	import { imageUrl } from '$api/client';
	import type { Movie, TVSeries } from '$api/types';

	interface Props {
		item: Movie | TVSeries;
		type: 'movie' | 'tvshow';
		progress?: number; // 0-100
	}

	let { item, type, progress }: Props = $props();

	const title = $derived('title' in item ? item.title : item.name);
	const year = $derived(
		'release_date' in item
			? item.release_date?.slice(0, 4)
			: item.first_air_date?.slice(0, 4)
	);
	const href = $derived(type === 'movie' ? `/movies/${item.id}` : `/tvshows/${item.id}`);
	const poster = $derived(item.poster_path ? imageUrl('poster', 'w342', item.poster_path) : '');
</script>

<a {href} class="group relative block overflow-hidden rounded-lg">
	<!-- Poster -->
	<div class="aspect-[2/3] w-full bg-neutral-900">
		{#if poster}
			<img
				src={poster}
				alt={title}
				loading="lazy"
				class="h-full w-full object-cover transition-transform duration-300 group-hover:scale-105"
			/>
		{:else}
			<div class="flex h-full items-center justify-center text-neutral-700">
				<span class="text-4xl">{type === 'movie' ? '🎬' : '📺'}</span>
			</div>
		{/if}
	</div>

	<!-- Hover overlay -->
	<div
		class="absolute inset-0 bg-gradient-to-t from-black/80 via-transparent to-transparent opacity-0 transition-opacity group-hover:opacity-100"
	></div>

	<!-- Rating badge -->
	{#if item.vote_average && item.vote_average > 0}
		<div
			class="absolute right-2 top-2 rounded-md bg-black/70 px-1.5 py-0.5 text-xs font-medium text-yellow-400"
		>
			★ {item.vote_average.toFixed(1)}
		</div>
	{/if}

	<!-- Progress bar -->
	{#if progress !== undefined && progress > 0}
		<div class="absolute bottom-0 left-0 right-0 h-1 bg-neutral-800">
			<div
				class="h-full bg-red-500 transition-all"
				style="width: {progress}%"
			></div>
		</div>
	{/if}

	<!-- Title area -->
	<div class="mt-2 px-0.5">
		<p class="truncate text-sm font-medium text-neutral-200 group-hover:text-white">
			{title}
		</p>
		{#if year}
			<p class="text-xs text-neutral-500">{year}</p>
		{/if}
	</div>
</a>

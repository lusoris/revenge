<script lang="ts">
	import type { Movie, TVSeries } from '$api/types';
	import MediaCard from './MediaCard.svelte';
	import type { Snippet } from 'svelte';

	interface Props {
		title: string;
		items: (Movie | TVSeries)[];
		type: 'movie' | 'tvshow';
		href?: string;
	}

	let { title, items, type, href }: Props = $props();
</script>

{#if items.length > 0}
	<section class="mb-8">
		<div class="mb-3 flex items-center justify-between">
			<h2 class="text-lg font-semibold text-white">{title}</h2>
			{#if href}
				<a href={href} class="text-sm text-neutral-400 transition-colors hover:text-white">
					See all →
				</a>
			{/if}
		</div>

		<div class="flex gap-3 overflow-x-auto pb-2 scrollbar-hide">
			{#each items as item (item.id)}
				<div class="w-36 flex-shrink-0 sm:w-40">
					<MediaCard {item} {type} />
				</div>
			{/each}
		</div>
	</section>
{/if}

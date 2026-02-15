<script lang="ts">
	import { imageUrl } from '$api/client';
	import type { Credit } from '$api/types';
	import * as Carousel from '$components/ui/carousel';
	import type { CarouselAPI } from '$components/ui/carousel/context.js';
	import UserIcon from '@lucide/svelte/icons/user';

	interface Props {
		items: Credit[];
		title?: string;
	}

	let { items, title = 'Cast' }: Props = $props();

	let api = $state<CarouselAPI>();
	let canPrev = $state(false);
	let canNext = $state(false);

	$effect(() => {
		if (api) {
			canPrev = api.canScrollPrev();
			canNext = api.canScrollNext();
			api.on('select', () => {
				canPrev = api!.canScrollPrev();
				canNext = api!.canScrollNext();
			});
		}
	});
</script>

<section class="mb-8">
	<div class="mb-3 flex items-center justify-between">
		<h2 class="text-lg font-semibold text-white">{title}</h2>
		{#if items.length > 6}
			<div class="flex gap-1">
				<button
					onclick={() => api?.scrollPrev()}
					disabled={!canPrev}
					class="flex size-7 items-center justify-center rounded-full border border-neutral-700 text-neutral-400 transition-colors hover:border-neutral-500 hover:text-white disabled:opacity-30 disabled:hover:border-neutral-700 disabled:hover:text-neutral-400"
					aria-label="Previous cast members"
				>
					<svg class="size-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="M15.75 19.5 8.25 12l7.5-7.5" /></svg>
				</button>
				<button
					onclick={() => api?.scrollNext()}
					disabled={!canNext}
					class="flex size-7 items-center justify-center rounded-full border border-neutral-700 text-neutral-400 transition-colors hover:border-neutral-500 hover:text-white disabled:opacity-30 disabled:hover:border-neutral-700 disabled:hover:text-neutral-400"
					aria-label="Next cast members"
				>
					<svg class="size-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5"><path stroke-linecap="round" stroke-linejoin="round" d="m8.25 4.5 7.5 7.5-7.5 7.5" /></svg>
				</button>
			</div>
		{/if}
	</div>

	<Carousel.Root
		setApi={(emblaApi) => (api = emblaApi)}
		opts={{ align: 'start', dragFree: true }}
		class="w-full"
	>
		<Carousel.Content class="-ms-2">
			{#each items as person (person.id)}
				<Carousel.Item class="basis-28 ps-2 sm:basis-32">
					<div class="text-center">
						<div class="mx-auto aspect-[2/3] w-full overflow-hidden rounded-lg bg-neutral-800">
							{#if person.profile_path}
								<img
									src={imageUrl('profile', 'w185', person.profile_path)}
									alt={person.name}
									class="h-full w-full object-cover object-top"
								/>
							{:else}
								<div class="flex h-full items-center justify-center text-neutral-600">
									<UserIcon class="size-8" />
								</div>
							{/if}
						</div>
						<p class="mt-1.5 truncate text-xs font-medium text-neutral-300">{person.name}</p>
						{#if person.character}
							<p class="truncate text-xs text-neutral-500">{person.character}</p>
						{/if}
					</div>
				</Carousel.Item>
			{/each}
		</Carousel.Content>
	</Carousel.Root>
</section>

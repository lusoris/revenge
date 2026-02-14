<script lang="ts">
	import { imageUrl } from '$api/client';
	import * as playbackApi from '$api/endpoints/playback';
	import * as tvshowsApi from '$api/endpoints/tvshows';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { createQuery } from '@tanstack/svelte-query';

	const showId = $derived(page.params.id);

	const show = createQuery(() => ({
		queryKey: ['tvshow', showId],
		queryFn: () => tvshowsApi.getTVShow(showId)
	}));

	const seasons = createQuery(() => ({
		queryKey: ['tvshow', showId, 'seasons'],
		queryFn: () => tvshowsApi.getTVShowSeasons(showId)
	}));

	const cast = createQuery(() => ({
		queryKey: ['tvshow', showId, 'cast'],
		queryFn: () => tvshowsApi.getTVShowCast(showId)
	}));

	const nextEp = createQuery(() => ({
		queryKey: ['tvshow', showId, 'next-episode'],
		queryFn: () => tvshowsApi.getNextEpisode(showId),
		retry: false
	}));

	let selectedSeasonId = $state<string | null>(null);

	const episodes = createQuery(() => ({
		queryKey: ['tvshow', 'season', selectedSeasonId, 'episodes'],
		queryFn: () => (selectedSeasonId ? tvshowsApi.getSeasonEpisodes(selectedSeasonId) : Promise.resolve([])),
		enabled: !!selectedSeasonId
	}));

	// Auto-select first season when loaded
	$effect(() => {
		if (seasons.data?.length && !selectedSeasonId) {
			selectedSeasonId = seasons.data[0].id;
		}
	});

	const backdrop = $derived(
		show.data?.backdrop_path ? imageUrl('backdrop', 'w1280', show.data.backdrop_path) : ''
	);
	const poster = $derived(
		show.data?.poster_path ? imageUrl('poster', 'w500', show.data.poster_path) : ''
	);

	async function playEpisode(episodeId: string) {
		const filesData = await tvshowsApi.getEpisodeFiles(episodeId);
		if (!filesData?.length) return;
		const session = await playbackApi.startPlayback({
			episode_file_id: filesData[0].id
		});
		goto(`/play/${session.id}`);
	}

	async function playNext() {
		if (nextEp.data) {
			await playEpisode(nextEp.data.id);
		}
	}
</script>

<svelte:head>
	<title>{show.data?.title ?? 'TV Show'} — Revenge</title>
</svelte:head>

{#if show.isPending}
	<div class="flex justify-center py-16">
		<div class="h-8 w-8 animate-spin rounded-full border-2 border-neutral-700 border-t-white"></div>
	</div>
{:else if show.isError}
	<div class="py-16 text-center text-red-400">Failed to load TV show.</div>
{:else if show.data}
	{@const s = show.data}

	<!-- Backdrop -->
	<div class="relative -mx-4 -mt-6 mb-6 sm:-mx-6 lg:-mx-8">
		{#if backdrop}
			<div class="aspect-[21/9] w-full overflow-hidden">
				<img src={backdrop} alt="" class="h-full w-full object-cover" />
				<div class="absolute inset-0 bg-gradient-to-t from-neutral-950 via-neutral-950/60 to-transparent"></div>
			</div>
		{:else}
			<div class="h-48 bg-neutral-900"></div>
		{/if}

		<div class="absolute bottom-0 left-0 right-0 px-4 pb-6 sm:px-6 lg:px-8">
			<div class="flex gap-6">
				{#if poster}
					<div class="hidden w-32 flex-shrink-0 sm:block lg:w-40">
						<img src={poster} alt={s.title} class="rounded-lg shadow-2xl" />
					</div>
				{/if}

				<div class="flex-1">
					<h1 class="text-3xl font-bold text-white lg:text-4xl">{s.title}</h1>

					<div class="mt-2 flex flex-wrap items-center gap-3 text-sm text-neutral-400">
						{#if s.first_air_date}
							<span>{s.first_air_date.slice(0, 4)}{s.last_air_date ? `–${s.status === 'Ended' ? s.last_air_date.slice(0, 4) : ''}` : ''}</span>
						{/if}
						{#if s.number_of_seasons}
							<span>{s.number_of_seasons} season{s.number_of_seasons > 1 ? 's' : ''}</span>
						{/if}
						{#if s.vote_average}
							<span class="text-yellow-400">★ {s.vote_average.toFixed(1)}</span>
						{/if}
						{#if s.genres?.length}
							<span>{s.genres.map((g) => g.name).join(', ')}</span>
						{/if}
					</div>

					<div class="mt-4 flex gap-3">
						{#if nextEp.data}
							<button
								onclick={playNext}
								class="flex items-center gap-2 rounded-lg bg-white px-5 py-2.5 text-sm font-semibold text-black transition-colors hover:bg-neutral-200"
							>
								▶ S{nextEp.data.season_number}E{nextEp.data.episode_number}
							</button>
						{/if}
					</div>
				</div>
			</div>
		</div>
	</div>

	<!-- Overview -->
	{#if s.overview}
		<section class="mb-8">
			<p class="max-w-3xl leading-relaxed text-neutral-300">{s.overview}</p>
		</section>
	{/if}

	<!-- Seasons tabs + Episodes -->
	{#if seasons.data?.length}
		<section class="mb-8">
			<div class="mb-4 flex gap-2 overflow-x-auto pb-1">
				{#each seasons.data as season (season.id)}
					<button
						onclick={() => (selectedSeasonId = season.id)}
						class="flex-shrink-0 rounded-lg px-4 py-1.5 text-sm transition-colors {selectedSeasonId === season.id
							? 'bg-white text-black font-medium'
							: 'bg-neutral-900 text-neutral-400 hover:bg-neutral-800 hover:text-white'}"
					>
						{season.name ?? `Season ${season.season_number}`}
					</button>
				{/each}
			</div>

			{#if episodes.isPending}
				<div class="flex justify-center py-8">
					<div class="h-6 w-6 animate-spin rounded-full border-2 border-neutral-700 border-t-white"></div>
				</div>
			{:else if episodes.data?.length}
				<div class="space-y-2">
					{#each episodes.data as ep (ep.id)}
						<button
							onclick={() => playEpisode(ep.id)}
							class="flex w-full items-start gap-4 rounded-lg border border-neutral-800 bg-neutral-900/50 p-3 text-left transition-colors hover:bg-neutral-800/70"
						>
							{#if ep.still_path}
								<img
									src={imageUrl('backdrop', 'w300', ep.still_path)}
									alt=""
									class="h-20 w-36 flex-shrink-0 rounded object-cover"
									loading="lazy"
								/>
							{:else}
								<div class="flex h-20 w-36 flex-shrink-0 items-center justify-center rounded bg-neutral-800 text-neutral-700">
									📺
								</div>
							{/if}

							<div class="min-w-0 flex-1">
								<p class="text-sm font-medium text-white">
									{ep.episode_number}. {ep.name ?? `Episode ${ep.episode_number}`}
								</p>
								{#if ep.runtime}
									<p class="text-xs text-neutral-500">{ep.runtime}m</p>
								{/if}
								{#if ep.overview}
									<p class="mt-1 line-clamp-2 text-xs leading-relaxed text-neutral-500">
										{ep.overview}
									</p>
								{/if}
							</div>
						</button>
					{/each}
				</div>
			{:else}
				<p class="py-4 text-sm text-neutral-500">No episodes found.</p>
			{/if}
		</section>
	{/if}

	<!-- Cast -->
	{#if cast.data?.items?.length}
		<section class="mb-8">
			<h2 class="mb-3 text-lg font-semibold text-white">Cast</h2>
			<div class="flex gap-3 overflow-x-auto pb-2">
				{#each cast.data.items.slice(0, 20) as person (person.id)}
					<div class="w-24 flex-shrink-0 text-center">
						{#if person.profile_path}
							<img
								src={imageUrl('profile', 'w185', person.profile_path)}
								alt={person.name}
								class="mx-auto h-24 w-24 rounded-full object-cover"
								loading="lazy"
							/>
						{:else}
							<div class="mx-auto flex h-24 w-24 items-center justify-center rounded-full bg-neutral-800 text-neutral-600">
								👤
							</div>
						{/if}
						<p class="mt-1 truncate text-xs font-medium text-neutral-300">{person.name}</p>
						{#if person.character}
							<p class="truncate text-xs text-neutral-600">{person.character}</p>
						{/if}
					</div>
				{/each}
			</div>
		</section>
	{/if}
{/if}

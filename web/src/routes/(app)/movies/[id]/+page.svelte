<script lang="ts">
	import { page } from '$app/state';
	import { createQuery } from '@tanstack/svelte-query';
	import { imageUrl } from '$api/client';
	import * as moviesApi from '$api/endpoints/movies';
	import * as playbackApi from '$api/endpoints/playback';
	import { goto } from '$app/navigation';

	const movieId = $derived(page.params.id);

	const movie = createQuery(() => ({
		queryKey: ['movie', movieId],
		queryFn: () => moviesApi.getMovie(movieId)
	}));

	const files = createQuery(() => ({
		queryKey: ['movie', movieId, 'files'],
		queryFn: () => moviesApi.getMovieFiles(movieId)
	}));

	const cast = createQuery(() => ({
		queryKey: ['movie', movieId, 'cast'],
		queryFn: () => moviesApi.getMovieCast(movieId)
	}));

	const progress = createQuery(() => ({
		queryKey: ['movie', movieId, 'progress'],
		queryFn: () => moviesApi.getWatchProgress(movieId),
		retry: false // 404 if not started
	}));

	async function play() {
		const fileList = $files.data;
		if (!fileList?.length) return;

		const session = await playbackApi.startPlayback({
			movie_file_id: fileList[0].id,
			start_position_seconds: $progress.data?.position_seconds
		});
		goto(`/play/${session.id}`);
	}

	const backdrop = $derived(
		$movie.data?.backdrop_path ? imageUrl('backdrop', 'w1280', $movie.data.backdrop_path) : ''
	);
	const poster = $derived(
		$movie.data?.poster_path ? imageUrl('poster', 'w500', $movie.data.poster_path) : ''
	);
</script>

<svelte:head>
	<title>{$movie.data?.title ?? 'Movie'} — Revenge</title>
</svelte:head>

{#if $movie.isPending}
	<div class="flex justify-center py-16">
		<div
			class="h-8 w-8 animate-spin rounded-full border-2 border-neutral-700 border-t-white"
		></div>
	</div>
{:else if $movie.isError}
	<div class="py-16 text-center text-red-400">Failed to load movie.</div>
{:else if $movie.data}
	{@const m = $movie.data}

	<!-- Backdrop hero -->
	<div class="relative -mx-4 -mt-6 mb-6 sm:-mx-6 lg:-mx-8">
		{#if backdrop}
			<div class="aspect-[21/9] w-full overflow-hidden">
				<img src={backdrop} alt="" class="h-full w-full object-cover" />
				<div
					class="absolute inset-0 bg-gradient-to-t from-neutral-950 via-neutral-950/60 to-transparent"
				></div>
			</div>
		{:else}
			<div class="h-48 bg-neutral-900"></div>
		{/if}

		<!-- Metadata overlay -->
		<div class="absolute bottom-0 left-0 right-0 px-4 pb-6 sm:px-6 lg:px-8">
			<div class="flex gap-6">
				<!-- Poster -->
				{#if poster}
					<div class="hidden w-32 flex-shrink-0 sm:block lg:w-40">
						<img src={poster} alt={m.title} class="rounded-lg shadow-2xl" />
					</div>
				{/if}

				<div class="flex-1">
					<h1 class="text-3xl font-bold text-white lg:text-4xl">{m.title}</h1>

					<div class="mt-2 flex flex-wrap items-center gap-3 text-sm text-neutral-400">
						{#if m.release_date}
							<span>{m.release_date.slice(0, 4)}</span>
						{/if}
						{#if m.runtime}
							<span>{Math.floor(m.runtime / 60)}h {m.runtime % 60}m</span>
						{/if}
						{#if m.vote_average}
							<span class="text-yellow-400">★ {m.vote_average.toFixed(1)}</span>
						{/if}
						{#if m.genres?.length}
							<span>{m.genres.map((g) => g.name).join(', ')}</span>
						{/if}
					</div>

					<!-- Play button -->
					<div class="mt-4 flex gap-3">
						<button
							onclick={play}
							disabled={!$files.data?.length}
							class="flex items-center gap-2 rounded-lg bg-white px-5 py-2.5 text-sm font-semibold text-black transition-colors hover:bg-neutral-200 disabled:opacity-40"
						>
							▶ {$progress.data?.position_seconds ? 'Resume' : 'Play'}
						</button>
					</div>
				</div>
			</div>
		</div>
	</div>

	<!-- Overview -->
	{#if m.overview}
		<section class="mb-8">
			<p class="max-w-3xl leading-relaxed text-neutral-300">{m.overview}</p>
		</section>
	{/if}

	<!-- Cast -->
	{#if $cast.data?.credits?.length}
		<section class="mb-8">
			<h2 class="mb-3 text-lg font-semibold text-white">Cast</h2>
			<div class="flex gap-3 overflow-x-auto pb-2">
				{#each $cast.data.credits.slice(0, 20) as person (person.id)}
					<div class="w-24 flex-shrink-0 text-center">
						{#if person.profile_path}
							<img
								src={imageUrl('profile', 'w185', person.profile_path)}
								alt={person.name}
								class="mx-auto h-24 w-24 rounded-full object-cover"
								loading="lazy"
							/>
						{:else}
							<div
								class="mx-auto flex h-24 w-24 items-center justify-center rounded-full bg-neutral-800 text-neutral-600"
							>
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

	<!-- Files -->
	{#if $files.data?.length}
		<section class="mb-8">
			<h2 class="mb-3 text-lg font-semibold text-white">Files</h2>
			<div class="space-y-2">
				{#each $files.data as file (file.id)}
					<div
						class="flex items-center justify-between rounded-lg border border-neutral-800 bg-neutral-900/50 px-4 py-3"
					>
						<div>
							<p class="text-sm text-neutral-300">{file.resolution ?? 'Unknown'}</p>
							<p class="text-xs text-neutral-600">
								{file.codec ?? ''} · {file.audio_codec ?? ''}{file.audio_channels
									? ` ${file.audio_channels}ch`
									: ''}
								{#if file.file_size}
									· {(file.file_size / (1024 * 1024 * 1024)).toFixed(1)} GB
								{/if}
							</p>
						</div>
					</div>
				{/each}
			</div>
		</section>
	{/if}
{/if}

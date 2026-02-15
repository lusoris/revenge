<script lang="ts">
	import { imageUrl } from '$api/client';
	import * as searchApi from '$api/endpoints/search';
	import { goto } from '$app/navigation';
	import * as Dialog from '$components/ui/dialog';
	import FilmIcon from '@lucide/svelte/icons/film';
	import SearchIcon from '@lucide/svelte/icons/search';
	import TvIcon from '@lucide/svelte/icons/tv';
	import { createQuery } from '@tanstack/svelte-query';

	let open = $state(false);
	let query = $state('');
	let debounced = $state('');
	let autocompleteTerm = $state('');
	let timer: ReturnType<typeof setTimeout>;
	let acTimer: ReturnType<typeof setTimeout>;
	let selectedIndex = $state(0);
	let inputRef = $state<HTMLInputElement | null>(null);

	function onInput(e: Event) {
		query = (e.target as HTMLInputElement).value;
		selectedIndex = 0;

		// Autocomplete fires faster (100ms)
		clearTimeout(acTimer);
		acTimer = setTimeout(() => {
			autocompleteTerm = query;
		}, 100);

		// Full search fires slower (300ms)
		clearTimeout(timer);
		timer = setTimeout(() => {
			debounced = query;
		}, 300);
	}

	// Keyboard shortcut: Ctrl+K / Cmd+K
	function handleGlobalKeydown(e: KeyboardEvent) {
		if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
			e.preventDefault();
			open = !open;
		}
	}

	// Autocomplete (fast, lightweight)
	const autocomplete = createQuery(() => ({
		queryKey: ['autocomplete', autocompleteTerm] as const,
		queryFn: async () => {
			const [movies, tvshows] = await Promise.all([
				searchApi.autocompleteMovies({ q: autocompleteTerm, limit: 4 }),
				searchApi.autocompleteTVShows({ q: autocompleteTerm, limit: 4 })
			]);
			// Deduplicate and merge suggestions
			const seen = new Set<string>();
			const all: string[] = [];
			for (const s of [...(movies.suggestions ?? []), ...(tvshows.suggestions ?? [])]) {
				const key = s.toLowerCase();
				if (!seen.has(key)) {
					seen.add(key);
					all.push(s);
				}
			}
			return all.slice(0, 5);
		},
		enabled: autocompleteTerm.length >= 1
	}));

	// Full search results (richer, slower)
	const results = createQuery(() => ({
		queryKey: ['command-search', debounced] as const,
		queryFn: () => searchApi.searchMulti({ q: debounced, limit: 10 }),
		enabled: debounced.length >= 2
	}));

	interface ResultItem {
		id: string;
		title: string;
		year?: string;
		poster?: string;
		type: 'movie' | 'tvshow';
		kind: 'suggestion' | 'result';
	}

	// Autocomplete suggestions (shown while full results load)
	const suggestions = $derived(autocomplete.data ?? []);

	const items = $derived.by(() => {
		if (!results.data) return [];
		const arr: ResultItem[] = [];
		for (const hit of results.data.movies?.hits ?? []) {
			const doc = hit.document;
			arr.push({
				id: doc.id,
				title: doc.title ?? 'Unknown',
				year: doc.release_date?.slice(0, 4),
				poster: doc.poster_path ? imageUrl('poster', 'w185', doc.poster_path) : undefined,
				type: 'movie',
				kind: 'result'
			});
		}
		for (const hit of results.data.tvshows?.hits ?? []) {
			const doc = hit.document;
			arr.push({
				id: doc.id,
				title: doc.title ?? doc.name ?? 'Unknown',
				year: doc.first_air_date?.slice(0, 4),
				poster: doc.poster_path ? imageUrl('poster', 'w185', doc.poster_path) : undefined,
				type: 'tvshow',
				kind: 'result'
			});
		}
		return arr;
	});

	// Combined navigable list: suggestions first (if no full results yet), then media items
	const allNavigable = $derived.by(() => {
		// If full results are ready, just show those
		if (items.length > 0) return items;
		// Otherwise show autocomplete suggestions as navigable items that trigger a search
		return suggestions.map(
			(s): ResultItem => ({
				id: s,
				title: s,
				type: 'movie',
				kind: 'suggestion'
			})
		);
	});

	function navigate(item: ResultItem) {
		if (item.kind === 'suggestion') {
			// Fill input with suggestion and trigger full search
			query = item.title;
			debounced = item.title;
			autocompleteTerm = item.title;
			selectedIndex = 0;
			inputRef?.focus();
			return;
		}
		const path = item.type === 'movie' ? `/movies/${item.id}` : `/tvshows/${item.id}`;
		open = false;
		query = '';
		debounced = '';
		autocompleteTerm = '';
		goto(path);
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'ArrowDown') {
			e.preventDefault();
			selectedIndex = Math.min(selectedIndex + 1, allNavigable.length - 1);
		} else if (e.key === 'ArrowUp') {
			e.preventDefault();
			selectedIndex = Math.max(selectedIndex - 1, 0);
		} else if (e.key === 'Enter' && allNavigable[selectedIndex]) {
			e.preventDefault();
			navigate(allNavigable[selectedIndex]);
		}
	}

	function handleOpenChange(next: boolean) {
		open = next;
		if (!next) {
			query = '';
			debounced = '';
			autocompleteTerm = '';
			selectedIndex = 0;
		}
	}

	// Focus input when dialog opens
	$effect(() => {
		if (open && inputRef) {
			requestAnimationFrame(() => inputRef?.focus());
		}
	});
</script>

<svelte:window onkeydown={handleGlobalKeydown} />

<!-- Trigger button for sidebar -->
<button
	onclick={() => (open = true)}
	class="flex w-full items-center gap-3 rounded-lg px-3 py-2 text-sm text-neutral-400 transition-colors hover:bg-neutral-900 hover:text-white"
>
	<SearchIcon class="size-4" />
	<span class="flex-1 text-left">Search</span>
	<kbd
		class="hidden rounded border border-neutral-700 bg-neutral-800 px-1.5 py-0.5 text-[10px] font-medium text-neutral-500 lg:inline-block"
		>⌘K</kbd
	>
</button>

<!-- Command palette dialog -->
<Dialog.Root bind:open onOpenChange={handleOpenChange}>
	<Dialog.Content
		showCloseButton={false}
		class="top-[20%] translate-y-0 gap-0 overflow-hidden rounded-xl border-neutral-700 bg-neutral-900 p-0 shadow-2xl sm:max-w-xl"
	>
		<Dialog.Title class="sr-only">Search your library</Dialog.Title>
		<Dialog.Description class="sr-only"
			>Type to search for movies and TV shows in your library</Dialog.Description
		>

		<!-- Search input -->
		<div class="flex items-center border-b border-neutral-800 px-4">
			<SearchIcon class="size-4 shrink-0 text-neutral-500" />
			<input
				bind:this={inputRef}
				type="text"
				value={query}
				oninput={onInput}
				onkeydown={handleKeydown}
				placeholder="Search movies and TV shows…"
				class="flex-1 bg-transparent px-3 py-3.5 text-sm text-white placeholder-neutral-500 outline-none"
			/>
			{#if query}
				<button
					onclick={() => {
						query = '';
						debounced = '';
						autocompleteTerm = '';
						selectedIndex = 0;
						inputRef?.focus();
					}}
					class="text-xs text-neutral-500 hover:text-neutral-300"
				>
					Clear
				</button>
			{/if}
		</div>

		<!-- Results -->
		<div class="max-h-80 overflow-y-auto">
			{#if query.length < 1}
				<div class="px-4 py-8 text-center text-sm text-neutral-500">
					Type to search your library
				</div>
			{:else if allNavigable.length === 0 && results.isPending}
				<div class="flex justify-center py-8">
					<div
						class="h-5 w-5 animate-spin rounded-full border-2 border-neutral-700 border-t-white"
					></div>
				</div>
			{:else if allNavigable.length === 0 && debounced.length >= 2}
				<div class="px-4 py-8 text-center text-sm text-neutral-500">
					No results for "{debounced}"
				</div>
			{:else if allNavigable.length > 0}
				<div class="py-1">
					{#each allNavigable as item, i (item.kind + item.id)}
						<button
							onclick={() => navigate(item)}
							onmouseenter={() => (selectedIndex = i)}
							class="flex w-full items-center gap-3 px-4 py-2.5 text-left transition-colors {i ===
							selectedIndex
								? 'bg-neutral-800 text-white'
								: 'text-neutral-300 hover:bg-neutral-800/50'}"
						>
							{#if item.kind === 'suggestion'}
								<!-- Autocomplete suggestion -->
								<SearchIcon class="size-4 shrink-0 text-neutral-500" />
								<span class="flex-1 truncate text-sm">{item.title}</span>
								{#if i === selectedIndex}
									<kbd class="text-[10px] text-neutral-500">↵</kbd>
								{/if}
							{:else}
								<!-- Full media result -->
								{#if item.poster}
									<img
										src={item.poster}
										alt=""
										class="h-10 w-7 shrink-0 rounded object-cover"
									/>
								{:else}
									<div
										class="flex h-10 w-7 shrink-0 items-center justify-center rounded bg-neutral-800"
									>
										{#if item.type === 'movie'}
											<FilmIcon class="size-3.5 text-neutral-600" />
										{:else}
											<TvIcon class="size-3.5 text-neutral-600" />
										{/if}
									</div>
								{/if}

								<div class="min-w-0 flex-1">
									<p class="truncate text-sm font-medium">{item.title}</p>
									<p class="text-xs text-neutral-500">
										{item.type === 'movie' ? 'Movie' : 'TV Show'}
										{#if item.year}
											&middot; {item.year}
										{/if}
									</p>
								</div>

								{#if i === selectedIndex}
									<kbd class="text-[10px] text-neutral-500">↵</kbd>
								{/if}
							{/if}
						</button>
					{/each}
				</div>
			{/if}
		</div>

		<!-- Footer -->
		<div
			class="flex items-center gap-4 border-t border-neutral-800 px-4 py-2 text-[10px] text-neutral-600"
		>
			<span><kbd class="font-mono">↑↓</kbd> navigate</span>
			<span><kbd class="font-mono">↵</kbd> open</span>
			<span><kbd class="font-mono">esc</kbd> close</span>
		</div>
	</Dialog.Content>
</Dialog.Root>

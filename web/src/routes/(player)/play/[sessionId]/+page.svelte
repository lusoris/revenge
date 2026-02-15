<script lang="ts">
	import * as playbackApi from '$api/endpoints/playback';
	import { page } from '$app/state';
	import VideoPlayer from '$components/player/VideoPlayer.svelte';
	import { createQuery } from '@tanstack/svelte-query';

	const sessionId = $derived(page.params.sessionId);

	const session = createQuery(() => ({
		queryKey: ['playback', sessionId],
		queryFn: () => playbackApi.getPlaybackSession(sessionId)
	}));

	async function handleBack() {
		try {
			await playbackApi.stopPlayback(sessionId);
		} catch {}
		history.back();
	}

	function handleHeartbeat(position: number) {
		playbackApi.heartbeat(sessionId, position).catch(() => {});
	}
</script>

<svelte:head>
	<title>Playing — Revenge</title>
</svelte:head>

{#if session.isPending}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-black">
		<div
			class="h-10 w-10 animate-spin rounded-full border-2 border-neutral-700 border-t-white"
		></div>
	</div>
{:else if session.isError}
	<div class="fixed inset-0 z-50 flex flex-col items-center justify-center gap-4 bg-black">
		<p class="text-red-400">Failed to load playback session.</p>
		<button
			onclick={() => history.back()}
			class="rounded-lg bg-neutral-800 px-4 py-2 text-sm text-white hover:bg-neutral-700"
		>
			Go back
		</button>
	</div>
{:else if session.data}
	<VideoPlayer session={session.data} onBack={handleBack} onHeartbeat={handleHeartbeat} />
{/if}

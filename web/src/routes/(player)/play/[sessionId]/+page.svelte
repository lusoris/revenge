<script lang="ts">
	import * as playbackApi from '$api/endpoints/playback';
	import { page } from '$app/state';
	import { createQuery } from '@tanstack/svelte-query';
	import { onDestroy } from 'svelte';

	const sessionId = $derived(page.params.sessionId);

	const session = createQuery(() => ({
		queryKey: ['playback', sessionId],
		queryFn: () => playbackApi.getPlaybackSession(sessionId)
	}));

	let videoEl: HTMLVideoElement | undefined = $state();
	let hls: any = null;
	let heartbeatInterval: ReturnType<typeof setInterval>;
	let controlsVisible = $state(true);
	let controlsTimer: ReturnType<typeof setTimeout>;

	async function initPlayer(masterUrl: string) {
		if (!videoEl) return;

		// Safari supports HLS natively
		if (videoEl.canPlayType('application/vnd.apple.mpegurl')) {
			videoEl.src = masterUrl;
			videoEl.play();
			return;
		}

		const { default: Hls } = await import('hls.js');
		if (!Hls.isSupported()) {
			console.error('HLS not supported in this browser');
			return;
		}

		hls = new Hls({
			startPosition: session.data?.start_position_seconds ?? 0,
			maxBufferLength: 30,
			maxMaxBufferLength: 60
		});
		hls.loadSource(masterUrl);
		hls.attachMedia(videoEl);
		hls.on(Hls.Events.MANIFEST_PARSED, () => {
			videoEl?.play();
		});
	}

	function startHeartbeat() {
		heartbeatInterval = setInterval(() => {
			if (videoEl && !videoEl.paused) {
				playbackApi.heartbeat(sessionId, videoEl.currentTime).catch(() => {});
			}
		}, 30_000);
	}

	function showControls() {
		controlsVisible = true;
		clearTimeout(controlsTimer);
		controlsTimer = setTimeout(() => {
			if (videoEl && !videoEl.paused) controlsVisible = false;
		}, 3000);
	}

	async function stopAndGoBack() {
		try {
			if (videoEl) {
				await playbackApi.heartbeat(sessionId, videoEl.currentTime);
			}
			await playbackApi.stopPlayback(sessionId);
		} catch {}
		history.back();
	}

	function handleKeydown(e: KeyboardEvent) {
		if (!videoEl) return;
		switch (e.key) {
			case ' ':
			case 'k':
				e.preventDefault();
				videoEl.paused ? videoEl.play() : videoEl.pause();
				break;
			case 'ArrowLeft':
				e.preventDefault();
				videoEl.currentTime = Math.max(0, videoEl.currentTime - 10);
				break;
			case 'ArrowRight':
				e.preventDefault();
				videoEl.currentTime += 10;
				break;
			case 'f':
				e.preventDefault();
				document.fullscreenElement ? document.exitFullscreen() : videoEl.requestFullscreen();
				break;
			case 'Escape':
				e.preventDefault();
				stopAndGoBack();
				break;
			case 'm':
				e.preventDefault();
				videoEl.muted = !videoEl.muted;
				break;
			case 'ArrowUp':
				e.preventDefault();
				videoEl.volume = Math.min(1, videoEl.volume + 0.1);
				break;
			case 'ArrowDown':
				e.preventDefault();
				videoEl.volume = Math.max(0, videoEl.volume - 0.1);
				break;
		}
		showControls();
	}

	$effect(() => {
		if (session.data?.master_playlist_url && videoEl) {
			initPlayer(session.data.master_playlist_url);
			startHeartbeat();
		}
	});

	onDestroy(() => {
		clearInterval(heartbeatInterval);
		clearTimeout(controlsTimer);
		if (hls) {
			hls.destroy();
			hls = null;
		}
	});
</script>

<svelte:head>
	<title>Playing — Revenge</title>
</svelte:head>

<svelte:window onkeydown={handleKeydown} />

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
	class="fixed inset-0 z-50 bg-black"
	onmousemove={showControls}
	onclick={showControls}
>
	{#if session.isPending}
		<div class="flex h-full items-center justify-center">
			<div class="h-10 w-10 animate-spin rounded-full border-2 border-neutral-700 border-t-white"></div>
		</div>
	{:else if session.isError}
		<div class="flex h-full flex-col items-center justify-center gap-4">
			<p class="text-red-400">Failed to load playback session.</p>
			<button
				onclick={() => history.back()}
				class="rounded-lg bg-neutral-800 px-4 py-2 text-sm text-white"
			>
				Go back
			</button>
		</div>
	{:else}
		<!-- svelte-ignore a11y_media_has_caption -->
		<video
			bind:this={videoEl}
			class="h-full w-full"
			controls
			playsinline
			autoplay
		></video>

		<!-- Overlay back button -->
		<div
			class="pointer-events-none absolute inset-x-0 top-0 bg-gradient-to-b from-black/70 to-transparent px-4 py-3 transition-opacity duration-300 {controlsVisible ? 'opacity-100' : 'opacity-0'}"
		>
			<div class="pointer-events-auto flex items-center gap-3">
				<button
					onclick={stopAndGoBack}
					class="rounded-lg bg-black/50 px-3 py-1.5 text-sm text-white transition-colors hover:bg-black/80"
				>
					← Back
				</button>

			</div>
		</div>
	{/if}
</div>

<script lang="ts">
	import type { AudioTrack, PlaybackProfile, PlaybackSession, SubtitleTrack } from '$api/types';
	import * as DropdownMenu from '$components/ui/dropdown-menu';
	import * as Tooltip from '$components/ui/tooltip';
	import ArrowLeftIcon from '@lucide/svelte/icons/arrow-left';
	import CheckIcon from '@lucide/svelte/icons/check';
	import GaugeIcon from '@lucide/svelte/icons/gauge';
	import LanguagesIcon from '@lucide/svelte/icons/languages';
	import LoaderIcon from '@lucide/svelte/icons/loader';
	import MaximizeIcon from '@lucide/svelte/icons/maximize';
	import MessageSquareTextIcon from '@lucide/svelte/icons/message-square-text';
	import MinimizeIcon from '@lucide/svelte/icons/minimize';
	import PauseIcon from '@lucide/svelte/icons/pause';
	import PlayIcon from '@lucide/svelte/icons/play';
	import Volume1Icon from '@lucide/svelte/icons/volume-1';
	import Volume2Icon from '@lucide/svelte/icons/volume-2';
	import VolumeOffIcon from '@lucide/svelte/icons/volume-off';
	import { onDestroy, onMount } from 'svelte';

	interface Props {
		session: PlaybackSession;
		onBack: () => void;
		onHeartbeat: (position: number) => void;
	}

	let { session, onBack, onHeartbeat }: Props = $props();

	// ─── State ────────────────────────────────────────────────────────────────
	let videoEl: HTMLVideoElement | undefined = $state();
	let containerEl: HTMLDivElement | undefined = $state();
	let hls: any = $state(null);
	let hlsInitialized = false;

	// Playback state
	let playing = $state(false);
	let currentTime = $state(0);
	let duration = $state(session.duration_seconds || 0);
	let buffered = $state(0);
	let volume = $state(1);
	let muted = $state(false);
	let isFullscreen = $state(false);
	let loading = $state(true);
	let error = $state('');

	// Controls visibility
	let controlsVisible = $state(true);
	let controlsTimer: ReturnType<typeof setTimeout>;
	let seeking = $state(false);

	// Track state
	let currentAudioTrack = $state(0);
	let currentSubtitleTrack = $state(-1); // -1 = off
	let currentQuality = $state(-1); // -1 = auto

	// Heartbeat
	let heartbeatInterval: ReturnType<typeof setInterval>;

	// Track the pending play() promise so we never interrupt it with pause()
	let playPromise: Promise<void> | undefined;

	// ─── Derived ──────────────────────────────────────────────────────────────

	const formattedTime = $derived(formatTime(currentTime));
	const formattedDuration = $derived(formatTime(duration));
	const progressPercent = $derived(duration > 0 ? (currentTime / duration) * 100 : 0);
	const bufferedPercent = $derived(duration > 0 ? (buffered / duration) * 100 : 0);

	const qualityLabel = $derived.by(() => {
		if (currentQuality === -1) return 'Auto';
		const level = hls?.levels?.[currentQuality];
		return level ? `${level.height}p` : 'Auto';
	});

	// ─── HLS Init ─────────────────────────────────────────────────────────────

	async function initPlayer() {
		if (!videoEl || hlsInitialized) return;
		hlsInitialized = true;
		loading = true;

		const url = session.master_playlist_url;
		console.log('[Player] Initializing HLS player with URL:', url);

		try {
			const { default: Hls } = await import('hls.js');

			// Prefer HLS.js over native HLS — it handles codec filtering, level
			// switching, and fMP4 segments reliably across all browsers.
			// Native HLS (Safari) is only used as a fallback when MSE isn't available.
			if (Hls.isSupported()) {
				console.log('[Player] Using HLS.js (MSE)');
				const hlsInstance = new Hls({
				startPosition: session.start_position ?? 0,
				maxBufferLength: 6,
				maxMaxBufferLength: 6,
				enableWorker: true,
				lowLatencyMode: false,
				backBufferLength: Infinity,
				debug: false,
				// Skip navigator.mediaCapabilities.decodingInfo() for HEVC.
				// Chrome on Linux may report HEVC as unsupported via decodingInfo()
				// even when hardware VA-API decoding works fine. The fallback
				// MediaSource.isTypeSupported() check is sufficient.
				useMediaCapabilities: false,
				// Prefer HDR content when multiple renditions are available.
				videoPreference: { preferHDR: true },
				// Force start at the highest quality level (original/4K).
				// HLS.js sorts levels by bandwidth ascending, so the highest = last index.
				// profile count from session.profiles gives us the exact number of levels.
				startLevel: session.profiles.length - 1,
				capLevelToPlayerSize: false,
				abrEwmaDefaultEstimate: 100_000_000, // 100 Mbps — local/fast network assumption
				abrEwmaDefaultEstimateMax: 100_000_000,
			});

			hlsInstance.loadSource(url);
			hlsInstance.attachMedia(videoEl);

			// Track codec-incompatible levels to avoid retrying them
			const brokenLevels = new Set<number>();
			let mediaRecoveryAttempted = false;

			hlsInstance.on(Hls.Events.MANIFEST_PARSED, (_: any, data: any) => {
				const levels = data.levels ?? [];
				console.log('[Player] Manifest parsed, levels:', levels.length);
				if (levels.length > 0) {
					const highest = levels[levels.length - 1];
					console.log('[Player] Highest level:', highest?.height + 'p', highest?.bitrate, 'bps');
				}
				addSubtitleTracks(videoEl!);
				if (videoEl) {
					playPromise = videoEl.play();
					playPromise
						.then(() => { playPromise = undefined; })
						.catch((e: any) => {
							playPromise = undefined;
							console.warn('[Player] Autoplay blocked:', e.message);
						});
				}
			});

			hlsInstance.on(Hls.Events.ERROR, (_: any, data: any) => {
				const detail = data.details ?? '';
				const isFatal = data.fatal;

				// Subtitle load failures (404s) are non-fatal — just silence them.
				if (detail === 'subtitleTrackLoadError' || detail === 'subtitleLoadError') {
					console.warn('[Player] Subtitle load error (ignored):', data.url);
					return;
				}

				console.error('[Player] HLS error:', data.type, detail, 'fatal:', isFatal);

				// Handle buffer/codec errors by falling back to a lower level.
				// This covers HEVC on browsers without MSE HEVC support.
				const isCodecError = detail === 'bufferAppendingError' ||
					detail === 'bufferAppendError' ||
					detail === 'bufferAddCodecError';

				if (isCodecError) {
					// Determine which level actually failed — use data.level if available,
					// fall back to loadLevel/currentLevel, but never blacklist -1 (auto).
					let broken = data.level ?? data.frag?.level ?? hlsInstance.loadLevel ?? hlsInstance.currentLevel;
					if (broken < 0) broken = hlsInstance.levels ? hlsInstance.levels.length - 1 : 0;
					brokenLevels.add(broken);
					console.warn(`[Player] Level ${broken} codec incompatible, blacklisting`);

					// Find next lower working level (search downward from the broken one)
					const levels = hlsInstance.levels ?? [];
					let fallback = -1;
					for (let i = broken - 1; i >= 0; i--) {
						if (!brokenLevels.has(i)) {
							fallback = i;
							break;
						}
					}

					if (fallback >= 0) {
						console.log(`[Player] Falling back to level ${fallback} (${levels[fallback]?.height}p)`);
						// Recover MSE state then switch level
						hlsInstance.recoverMediaError();
						hlsInstance.currentLevel = fallback;
						hlsInstance.loadLevel = fallback;
						// Lock ABR to only use working levels
						hlsInstance.autoLevelCapping = fallback;
						return;
					}
					// No working levels left
					error = 'No compatible video quality available for this browser.';
					hlsInstance.destroy();
					return;
				}

				if (!isFatal) return;

				switch (data.type) {
					case Hls.ErrorTypes.NETWORK_ERROR:
						console.error('[Player] Network error, recovering...');
						hlsInstance.startLoad();
						break;
					case Hls.ErrorTypes.MEDIA_ERROR:
						if (!mediaRecoveryAttempted) {
							mediaRecoveryAttempted = true;
							console.error('[Player] Media error, recovering...');
							hlsInstance.recoverMediaError();
						} else {
							// Second media error — swap audio codec approach
							console.error('[Player] Media error persists, swapping codec...');
							hlsInstance.swapAudioCodec();
							hlsInstance.recoverMediaError();
						}
						break;
					default:
						error = `Playback failed: ${detail}`;
						hlsInstance.destroy();
						break;
				}
			});

			hlsInstance.on(Hls.Events.LEVEL_SWITCHED, (_: any, data: any) => {
				console.log('[Player] Level switched to:', data.level);
			});

			hls = hlsInstance;
			} else if (videoEl.canPlayType('application/vnd.apple.mpegurl')) {
				// Fallback: native HLS (Safari without MSE, or older iOS)
				console.log('[Player] Using native HLS (Safari fallback)');
				videoEl.src = url;
				addSubtitleTracks(videoEl);
				playPromise = videoEl.play();
				playPromise
					.then(() => { playPromise = undefined; loading = false; })
					.catch((e: any) => {
						playPromise = undefined;
						console.warn('[Player] Native autoplay blocked:', e.message);
						loading = false;
					});
			} else {
				error = 'HLS playback is not supported in this browser.';
				loading = false;
			}
		} catch (e) {
			console.error('[Player] Init failed:', e);
			error = 'Failed to initialize video player.';
			loading = false;
		}
	}

	function addSubtitleTracks(video: HTMLVideoElement) {
		// Add WebVTT subtitle tracks from session data
		for (const sub of session.subtitle_tracks ?? []) {
			if (!sub.url) continue;
			const track = document.createElement('track');
			track.kind = sub.is_forced ? 'forced' : 'subtitles';
			track.label = sub.title ?? sub.language ?? `Track ${sub.index}`;
			track.srclang = sub.language ?? 'und';
			track.src = sub.url;
			if (sub.is_default) track.default = true;
			video.appendChild(track);
		}
	}

	// ─── Controls ─────────────────────────────────────────────────────────────

	function togglePlay() {
		if (!videoEl) return;
		if (videoEl.paused) {
			playPromise = videoEl.play();
			playPromise
				.then(() => { playPromise = undefined; })
				.catch((e: any) => {
					playPromise = undefined;
					console.warn('[Player] Play failed:', e.message);
				});
		} else {
			// Wait for any pending play() to settle before pausing
			if (playPromise) {
				playPromise.then(() => videoEl?.pause()).catch(() => videoEl?.pause());
				playPromise = undefined;
			} else {
				videoEl.pause();
			}
		}
	}

	function seek(e: MouseEvent | TouchEvent) {
		if (!videoEl || !duration) return;
		const bar = (e.currentTarget as HTMLElement);
		const rect = bar.getBoundingClientRect();
		const clientX = 'touches' in e ? e.touches[0].clientX : e.clientX;
		const pct = Math.max(0, Math.min(1, (clientX - rect.left) / rect.width));
		videoEl.currentTime = pct * duration;
	}

	function setVolume(e: MouseEvent) {
		const bar = (e.currentTarget as HTMLElement);
		const rect = bar.getBoundingClientRect();
		const pct = Math.max(0, Math.min(1, (e.clientX - rect.left) / rect.width));
		volume = pct;
		if (videoEl) {
			videoEl.volume = pct;
			videoEl.muted = pct === 0;
			muted = pct === 0;
		}
	}

	function toggleMute() {
		if (!videoEl) return;
		muted = !muted;
		videoEl.muted = muted;
	}

	function toggleFullscreen() {
		if (!containerEl) return;
		if (document.fullscreenElement) {
			document.exitFullscreen();
		} else {
			containerEl.requestFullscreen();
		}
	}

	function setAudioTrack(index: number) {
		currentAudioTrack = index;
		if (hls) {
			hls.audioTrack = index;
		}
	}

	function setSubtitleTrack(index: number) {
		currentSubtitleTrack = index;
		if (videoEl?.textTracks) {
			for (let i = 0; i < videoEl.textTracks.length; i++) {
				videoEl.textTracks[i].mode = i === index ? 'showing' : 'hidden';
			}
		}
		// Also try hls.js subtitle track switching
		if (hls) {
			hls.subtitleTrack = index;
		}
	}

	function setQuality(index: number) {
		currentQuality = index;
		if (hls) {
			hls.currentLevel = index; // -1 = auto
		}
	}

	function showControls() {
		controlsVisible = true;
		clearTimeout(controlsTimer);
		controlsTimer = setTimeout(() => {
			if (videoEl && !videoEl.paused && !seeking) controlsVisible = false;
		}, 3000);
	}

	// ─── Keyboard ─────────────────────────────────────────────────────────────

	function handleKeydown(e: KeyboardEvent) {
		if (!videoEl) return;
		// Don't capture when a dropdown menu is open (typing in it)
		if ((e.target as HTMLElement)?.closest('[data-radix-popper-content-wrapper]')) return;

		switch (e.key) {
			case ' ':
			case 'k':
				e.preventDefault();
				togglePlay();
				break;
			case 'ArrowLeft':
				e.preventDefault();
				videoEl.currentTime = Math.max(0, videoEl.currentTime - 10);
				break;
			case 'ArrowRight':
				e.preventDefault();
				videoEl.currentTime += 10;
				break;
			case 'j':
				e.preventDefault();
				videoEl.currentTime = Math.max(0, videoEl.currentTime - 30);
				break;
			case 'l':
				e.preventDefault();
				videoEl.currentTime += 30;
				break;
			case 'f':
				e.preventDefault();
				toggleFullscreen();
				break;
			case 'Escape':
				e.preventDefault();
				if (document.fullscreenElement) {
					document.exitFullscreen();
				} else {
					onBack();
				}
				break;
			case 'm':
				e.preventDefault();
				toggleMute();
				break;
			case 'ArrowUp':
				e.preventDefault();
				volume = Math.min(1, volume + 0.1);
				videoEl.volume = volume;
				break;
			case 'ArrowDown':
				e.preventDefault();
				volume = Math.max(0, volume - 0.1);
				videoEl.volume = volume;
				break;
			case '0':
			case '1':
			case '2':
			case '3':
			case '4':
			case '5':
			case '6':
			case '7':
			case '8':
			case '9':
				e.preventDefault();
				videoEl.currentTime = (parseInt(e.key) / 10) * duration;
				break;
		}
		showControls();
	}

	// ─── Video Event Handlers ─────────────────────────────────────────────────

	function onTimeUpdate() {
		if (videoEl && !seeking) {
			currentTime = videoEl.currentTime;
		}
	}

	function onDurationChange() {
		if (videoEl && videoEl.duration && isFinite(videoEl.duration)) {
			duration = videoEl.duration;
		}
	}

	function onProgress() {
		if (videoEl && videoEl.buffered.length > 0) {
			buffered = videoEl.buffered.end(videoEl.buffered.length - 1);
		}
	}

	function onPlay() {
		playing = true;
		loading = false;
	}

	function onPause() {
		playing = false;
		controlsVisible = true;
	}

	function onWaiting() {
		loading = true;
	}

	function onCanPlay() {
		loading = false;
	}

	function onFullscreenChange() {
		isFullscreen = !!document.fullscreenElement;
	}

	// ─── Utilities ────────────────────────────────────────────────────────────

	function formatTime(seconds: number): string {
		if (!isFinite(seconds) || seconds < 0) return '0:00';
		const h = Math.floor(seconds / 3600);
		const m = Math.floor((seconds % 3600) / 60);
		const s = Math.floor(seconds % 60);
		if (h > 0) return `${h}:${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`;
		return `${m}:${s.toString().padStart(2, '0')}`;
	}

	function audioTrackLabel(track: AudioTrack): string {
		const parts: string[] = [];
		if (track.title) parts.push(track.title);
		else if (track.language) parts.push(track.language.toUpperCase());
		if (track.channels) {
			if (track.channels === 2) parts.push('Stereo');
			else if (track.channels === 6) parts.push('5.1');
			else if (track.channels === 8) parts.push('7.1');
			else parts.push(`${track.channels}ch`);
		}
		if (track.codec) parts.push(track.codec.toUpperCase());
		return parts.join(' · ') || `Track ${track.index + 1}`;
	}

	function subtitleTrackLabel(track: SubtitleTrack): string {
		const parts: string[] = [];
		if (track.title) parts.push(track.title);
		else if (track.language) parts.push(track.language.toUpperCase());
		if (track.is_forced) parts.push('Forced');
		return parts.join(' · ') || `Track ${track.index + 1}`;
	}

	function profileLabel(profile: PlaybackProfile): string {
		if (profile.is_original) return `Original (${profile.height}p)`;
		return `${profile.height}p`;
	}

	// ─── Lifecycle ────────────────────────────────────────────────────────────

	onMount(() => {
		console.log('[Player] onMount, videoEl:', !!videoEl, 'url:', session.master_playlist_url);
		if (session.master_playlist_url && videoEl) {
			initPlayer();
		}

		// Heartbeat every 30 seconds
		heartbeatInterval = setInterval(() => {
			if (videoEl && !videoEl.paused) {
				onHeartbeat(videoEl.currentTime);
			}
		}, 30_000);

		document.addEventListener('fullscreenchange', onFullscreenChange);
	});

	onDestroy(() => {
		clearInterval(heartbeatInterval);
		clearTimeout(controlsTimer);
		document.removeEventListener('fullscreenchange', onFullscreenChange);
		if (hls) {
			hls.destroy();
			hls = null;
		}
	});
</script>

<svelte:window onkeydown={handleKeydown} />

<Tooltip.Provider>
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
	bind:this={containerEl}
	class="fixed inset-0 z-50 bg-black"
	onmousemove={showControls}
	onclick={(e) => {
		// Only toggle play on direct clicks (not on controls)
		if (e.target === videoEl || (e.target as HTMLElement)?.classList.contains('video-backdrop')) {
			togglePlay();
		}
		showControls();
	}}
	ondblclick={(e) => {
		if (e.target === videoEl || (e.target as HTMLElement)?.classList.contains('video-backdrop')) {
			toggleFullscreen();
		}
	}}
>
	{#if error}
		<div class="flex h-full flex-col items-center justify-center gap-4">
			<p class="text-red-400">{error}</p>
			<button
				onclick={onBack}
				class="rounded-lg bg-neutral-800 px-4 py-2 text-sm text-white hover:bg-neutral-700"
			>
				Go back
			</button>
		</div>
	{:else}
		<!-- svelte-ignore a11y_media_has_caption -->
		<video
			bind:this={videoEl}
			class="h-full w-full object-contain"
			playsinline
			ontimeupdate={onTimeUpdate}
			ondurationchange={onDurationChange}
			onprogress={onProgress}
			onplay={onPlay}
			onpause={onPause}
			onwaiting={onWaiting}
			oncanplay={onCanPlay}
		></video>

		<!-- Click target / backdrop -->
		<div class="video-backdrop absolute inset-0"></div>

		<!-- Loading spinner (center) -->
		{#if loading && !error}
			<div class="pointer-events-none absolute inset-0 flex items-center justify-center">
				<LoaderIcon class="size-12 animate-spin text-white/80" />
			</div>
		{/if}

		<!-- Top bar: Back button + title -->
		<div
			class="absolute inset-x-0 top-0 z-10 bg-gradient-to-b from-black/80 to-transparent px-4 py-3 transition-opacity duration-300 {controlsVisible ? 'opacity-100' : 'pointer-events-none opacity-0'}"
		>
			<div class="flex items-center gap-3">
				<button
					onclick={onBack}
					class="flex items-center gap-1.5 rounded-lg bg-black/40 px-3 py-1.5 text-sm text-white backdrop-blur-sm transition-colors hover:bg-black/60"
				>
					<ArrowLeftIcon class="size-4" />
					Back
				</button>
			</div>
		</div>

		<!-- Bottom controls -->
		<div
			class="absolute inset-x-0 bottom-0 z-10 bg-gradient-to-t from-black/90 via-black/50 to-transparent transition-opacity duration-300 {controlsVisible ? 'opacity-100' : 'pointer-events-none opacity-0'}"
		>
			<!-- Progress bar -->
			<div class="group px-4">
				<!-- svelte-ignore a11y_no_static_element_interactions -->
				<div
					class="relative h-1 w-full cursor-pointer rounded-full bg-white/20 transition-all group-hover:h-1.5"
					onclick={seek}
					onmousedown={(e) => {
						seeking = true;
						seek(e);
						const onMove = (ev: MouseEvent) => {
							const bar = e.currentTarget as HTMLElement;
							const rect = bar.getBoundingClientRect();
							const pct = Math.max(0, Math.min(1, (ev.clientX - rect.left) / rect.width));
							currentTime = pct * duration;
							if (videoEl) videoEl.currentTime = currentTime;
						};
						const onUp = () => {
							seeking = false;
							window.removeEventListener('mousemove', onMove);
							window.removeEventListener('mouseup', onUp);
						};
						window.addEventListener('mousemove', onMove);
						window.addEventListener('mouseup', onUp);
					}}
				>
					<!-- Buffered -->
					<div
						class="absolute inset-y-0 left-0 rounded-full bg-white/30"
						style="width: {bufferedPercent}%"
					></div>
					<!-- Progress -->
					<div
						class="absolute inset-y-0 left-0 rounded-full bg-red-500"
						style="width: {progressPercent}%"
					></div>
					<!-- Thumb -->
					<div
						class="absolute top-1/2 -translate-y-1/2 size-3 rounded-full bg-red-500 opacity-0 shadow transition-opacity group-hover:opacity-100"
						style="left: calc({progressPercent}% - 6px)"
					></div>
				</div>
			</div>

			<!-- Control buttons -->
			<div class="flex items-center gap-1 px-4 py-2">
				<!-- Play/Pause -->
				<button
					onclick={togglePlay}
					class="rounded-md p-2 text-white transition-colors hover:bg-white/10"
				>
					{#if playing}
						<PauseIcon class="size-5" />
					{:else}
						<PlayIcon class="size-5" />
					{/if}
				</button>

				<!-- Volume -->
				<div class="group/vol flex items-center gap-1">
					<button
						onclick={toggleMute}
						class="rounded-md p-2 text-white transition-colors hover:bg-white/10"
					>
						{#if muted || volume === 0}
							<VolumeOffIcon class="size-5" />
						{:else if volume < 0.5}
							<Volume1Icon class="size-5" />
						{:else}
							<Volume2Icon class="size-5" />
						{/if}
					</button>
					<!-- svelte-ignore a11y_no_static_element_interactions -->
					<div
						class="w-0 overflow-hidden transition-all duration-200 group-hover/vol:w-20"
					>
						<div
							class="relative h-1 w-20 cursor-pointer rounded-full bg-white/20"
							onclick={setVolume}
						>
							<div
								class="absolute inset-y-0 left-0 rounded-full bg-white"
								style="width: {muted ? 0 : volume * 100}%"
							></div>
						</div>
					</div>
				</div>

				<!-- Time -->
				<span class="ml-2 select-none text-xs tabular-nums text-white/80">
					{formattedTime} / {formattedDuration}
				</span>

				<!-- Spacer -->
				<div class="flex-1"></div>

				<!-- Audio tracks -->
				{#if (session.audio_tracks?.length ?? 0) > 1}
					<DropdownMenu.Root>
						<Tooltip.Root>
							<Tooltip.Trigger>
								<DropdownMenu.Trigger
									class="rounded-md p-2 text-white transition-colors hover:bg-white/10"
								>
									<LanguagesIcon class="size-5" />
								</DropdownMenu.Trigger>
							</Tooltip.Trigger>
							<Tooltip.Content>Audio</Tooltip.Content>
						</Tooltip.Root>
						<DropdownMenu.Content
							class="max-h-64 min-w-48 overflow-y-auto border-neutral-700 bg-neutral-900"
							align="end"
							side="top"
						>
							<DropdownMenu.Label class="text-neutral-400">Audio Track</DropdownMenu.Label>
							<DropdownMenu.Separator />
							{#each session.audio_tracks ?? [] as track, i (track.index)}
								<DropdownMenu.Item
									onclick={() => setAudioTrack(i)}
									class="flex items-center gap-2 {currentAudioTrack === i ? 'text-white' : 'text-neutral-400'}"
								>
									{#if currentAudioTrack === i}
										<CheckIcon class="size-4 text-red-500" />
									{:else}
										<div class="size-4"></div>
									{/if}
									{audioTrackLabel(track)}
								</DropdownMenu.Item>
							{/each}
						</DropdownMenu.Content>
					</DropdownMenu.Root>
				{/if}

				<!-- Subtitles -->
				{#if (session.subtitle_tracks?.length ?? 0) > 0}
					<DropdownMenu.Root>
						<Tooltip.Root>
							<Tooltip.Trigger>
								<DropdownMenu.Trigger
									class="rounded-md p-2 text-white transition-colors hover:bg-white/10 {currentSubtitleTrack >= 0 ? 'text-yellow-400' : ''}"
								>
									<MessageSquareTextIcon class="size-5" />
								</DropdownMenu.Trigger>
							</Tooltip.Trigger>
							<Tooltip.Content>Subtitles</Tooltip.Content>
						</Tooltip.Root>
						<DropdownMenu.Content
							class="max-h-64 min-w-48 overflow-y-auto border-neutral-700 bg-neutral-900"
							align="end"
							side="top"
						>
							<DropdownMenu.Label class="text-neutral-400">Subtitles</DropdownMenu.Label>
							<DropdownMenu.Separator />
							<DropdownMenu.Item
								onclick={() => setSubtitleTrack(-1)}
								class="flex items-center gap-2 {currentSubtitleTrack === -1 ? 'text-white' : 'text-neutral-400'}"
							>
								{#if currentSubtitleTrack === -1}
									<CheckIcon class="size-4 text-red-500" />
								{:else}
									<div class="size-4"></div>
								{/if}
								Off
							</DropdownMenu.Item>
							{#each session.subtitle_tracks ?? [] as track, i (track.index)}
								<DropdownMenu.Item
									onclick={() => setSubtitleTrack(i)}
									class="flex items-center gap-2 {currentSubtitleTrack === i ? 'text-white' : 'text-neutral-400'}"
								>
									{#if currentSubtitleTrack === i}
										<CheckIcon class="size-4 text-red-500" />
									{:else}
										<div class="size-4"></div>
									{/if}
									{subtitleTrackLabel(track)}
								</DropdownMenu.Item>
							{/each}
						</DropdownMenu.Content>
					</DropdownMenu.Root>
				{/if}

				<!-- Quality -->
				{#if (session.profiles?.length ?? 0) > 1}
					<DropdownMenu.Root>
						<Tooltip.Root>
							<Tooltip.Trigger>
								<DropdownMenu.Trigger
									class="flex items-center gap-1 rounded-md px-2 py-1.5 text-xs font-medium text-white transition-colors hover:bg-white/10"
								>
									<GaugeIcon class="size-4" />
									<span>{qualityLabel}</span>
								</DropdownMenu.Trigger>
							</Tooltip.Trigger>
							<Tooltip.Content>Quality</Tooltip.Content>
						</Tooltip.Root>
						<DropdownMenu.Content
							class="min-w-40 border-neutral-700 bg-neutral-900"
							align="end"
							side="top"
						>
							<DropdownMenu.Label class="text-neutral-400">Quality</DropdownMenu.Label>
							<DropdownMenu.Separator />
							<DropdownMenu.Item
								onclick={() => setQuality(-1)}
								class="flex items-center gap-2 {currentQuality === -1 ? 'text-white' : 'text-neutral-400'}"
							>
								{#if currentQuality === -1}
									<CheckIcon class="size-4 text-red-500" />
								{:else}
									<div class="size-4"></div>
								{/if}
								Auto
							</DropdownMenu.Item>
							{#each session.profiles ?? [] as profile, i (profile.name)}
								<DropdownMenu.Item
									onclick={() => setQuality(i)}
									class="flex items-center gap-2 {currentQuality === i ? 'text-white' : 'text-neutral-400'}"
								>
									{#if currentQuality === i}
										<CheckIcon class="size-4 text-red-500" />
									{:else}
										<div class="size-4"></div>
									{/if}
									{profileLabel(profile)}
									<span class="ml-auto text-[10px] text-neutral-500">
										{(profile.bitrate / 1000).toFixed(0)} Mbps
									</span>
								</DropdownMenu.Item>
							{/each}
						</DropdownMenu.Content>
					</DropdownMenu.Root>
				{/if}

				<!-- Fullscreen -->
				<Tooltip.Root>
					<Tooltip.Trigger>
						<button
							onclick={toggleFullscreen}
							class="rounded-md p-2 text-white transition-colors hover:bg-white/10"
						>
							{#if isFullscreen}
								<MinimizeIcon class="size-5" />
							{:else}
								<MaximizeIcon class="size-5" />
							{/if}
						</button>
					</Tooltip.Trigger>
					<Tooltip.Content>Fullscreen (F)</Tooltip.Content>
				</Tooltip.Root>
			</div>
		</div>
	{/if}
</div>
</Tooltip.Provider>

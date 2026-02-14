<script lang="ts">
	import {
		radarrGetQualityProfiles, radarrGetRootFolders,
		radarrGetStatus, radarrSync,
		sonarrGetQualityProfiles, sonarrGetRootFolders,
		sonarrGetStatus, sonarrSync
	} from '$api/endpoints/admin';
	import { Badge } from '$components/ui/badge';
	import { Button } from '$components/ui/button';
	import * as Card from '$components/ui/card';
	import { Separator } from '$components/ui/separator';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';

	const queryClient = useQueryClient();

	const radarrStatus = createQuery(() => ({
		queryKey: ['admin', 'radarr', 'status'],
		queryFn: () => radarrGetStatus(),
		retry: false
	}));

	const radarrProfiles = createQuery(() => ({
		queryKey: ['admin', 'radarr', 'profiles'],
		queryFn: () => radarrGetQualityProfiles(),
		retry: false
	}));

	const radarrFolders = createQuery(() => ({
		queryKey: ['admin', 'radarr', 'folders'],
		queryFn: () => radarrGetRootFolders(),
		retry: false
	}));

	const sonarrStatus = createQuery(() => ({
		queryKey: ['admin', 'sonarr', 'status'],
		queryFn: () => sonarrGetStatus(),
		retry: false
	}));

	const sonarrProfiles = createQuery(() => ({
		queryKey: ['admin', 'sonarr', 'profiles'],
		queryFn: () => sonarrGetQualityProfiles(),
		retry: false
	}));

	const sonarrFolders = createQuery(() => ({
		queryKey: ['admin', 'sonarr', 'folders'],
		queryFn: () => sonarrGetRootFolders(),
		retry: false
	}));

	let radarrSyncing = $state(false);
	let sonarrSyncing = $state(false);

	async function handleRadarrSync() {
		radarrSyncing = true;
		try {
			await radarrSync();
			setTimeout(() => {
				queryClient.invalidateQueries({ queryKey: ['admin', 'radarr'] });
				radarrSyncing = false;
			}, 2000);
		} catch {
			radarrSyncing = false;
		}
	}

	async function handleSonarrSync() {
		sonarrSyncing = true;
		try {
			await sonarrSync();
			setTimeout(() => {
				queryClient.invalidateQueries({ queryKey: ['admin', 'sonarr'] });
				sonarrSyncing = false;
			}, 2000);
		} catch {
			sonarrSyncing = false;
		}
	}

	function formatBytes(bytes?: number): string {
		if (!bytes) return '—';
		const gb = bytes / 1024 / 1024 / 1024;
		return `${gb.toFixed(1)} GB`;
	}
</script>

<div class="grid gap-6 lg:grid-cols-2">
	<!-- Radarr -->
	<Card.Root class="border-neutral-800 bg-neutral-900">
		<Card.Header class="flex flex-row items-center justify-between">
			<div>
				<Card.Title class="text-white">Radarr</Card.Title>
				<Card.Description>Movie library integration</Card.Description>
			</div>
			{#if radarrStatus.data?.connected}
				<Badge class="bg-green-900 text-green-200">Connected</Badge>
			{:else}
				<Badge variant="destructive">Disconnected</Badge>
			{/if}
		</Card.Header>
		<Card.Content class="space-y-4">
			{#if radarrStatus.data}
				{@const s = radarrStatus.data}
				<div class="grid grid-cols-2 gap-3 text-sm">
					<div>
						<p class="text-neutral-500">Version</p>
						<p class="text-white">{s.version ?? '—'}</p>
					</div>
					<div>
						<p class="text-neutral-500">Total Movies</p>
						<p class="text-white">{s.syncStatus.totalMovies}</p>
					</div>
					<div>
						<p class="text-neutral-500">Last Sync</p>
						<p class="text-white">
							{s.syncStatus.lastSync
								? new Date(s.syncStatus.lastSync).toLocaleString()
								: 'Never'}
						</p>
					</div>
					<div>
						<p class="text-neutral-500">Sync Status</p>
						<p class="text-white">
							{s.syncStatus.isRunning ? '🔄 Running' : '✅ Idle'}
						</p>
					</div>
				</div>

				{#if s.syncStatus.lastSyncError}
					<p class="text-xs text-red-400">{s.syncStatus.lastSyncError}</p>
				{/if}

				<Separator class="bg-neutral-800" />

				<!-- Quality Profiles -->
				{#if radarrProfiles.data && radarrProfiles.data.length > 0}
					<div>
						<p class="mb-2 text-xs font-medium uppercase text-neutral-500">Quality Profiles</p>
						<div class="flex flex-wrap gap-1">
							{#each radarrProfiles.data as p}
								<Badge variant="outline" class="text-neutral-300">{p.name}</Badge>
							{/each}
						</div>
					</div>
				{/if}

				<!-- Root Folders -->
				{#if radarrFolders.data && radarrFolders.data.length > 0}
					<div>
						<p class="mb-2 text-xs font-medium uppercase text-neutral-500">Root Folders</p>
						{#each radarrFolders.data as f}
							<div class="flex items-center justify-between text-sm">
								<span class="text-neutral-300">{f.path}</span>
								<span class="text-neutral-500">{formatBytes(f.freeSpace)} free</span>
							</div>
						{/each}
					</div>
				{/if}
			{:else if radarrStatus.isLoading}
				<p class="text-sm text-neutral-500">Loading…</p>
			{:else if radarrStatus.isError}
				<p class="text-sm text-red-400">Failed to connect to Radarr</p>
			{/if}
		</Card.Content>
		<Card.Footer>
			<Button size="sm" onclick={handleRadarrSync} disabled={radarrSyncing}>
				{radarrSyncing ? 'Syncing…' : 'Sync Now'}
			</Button>
		</Card.Footer>
	</Card.Root>

	<!-- Sonarr -->
	<Card.Root class="border-neutral-800 bg-neutral-900">
		<Card.Header class="flex flex-row items-center justify-between">
			<div>
				<Card.Title class="text-white">Sonarr</Card.Title>
				<Card.Description>TV series library integration</Card.Description>
			</div>
			{#if sonarrStatus.data?.connected}
				<Badge class="bg-green-900 text-green-200">Connected</Badge>
			{:else}
				<Badge variant="destructive">Disconnected</Badge>
			{/if}
		</Card.Header>
		<Card.Content class="space-y-4">
			{#if sonarrStatus.data}
				{@const s = sonarrStatus.data}
				<div class="grid grid-cols-2 gap-3 text-sm">
					<div>
						<p class="text-neutral-500">Version</p>
						<p class="text-white">{s.version ?? '—'}</p>
					</div>
					<div>
						<p class="text-neutral-500">Total Series</p>
						<p class="text-white">{s.syncStatus.totalSeries}</p>
					</div>
					<div>
						<p class="text-neutral-500">Last Sync</p>
						<p class="text-white">
							{s.syncStatus.lastSync
								? new Date(s.syncStatus.lastSync).toLocaleString()
								: 'Never'}
						</p>
					</div>
					<div>
						<p class="text-neutral-500">Sync Status</p>
						<p class="text-white">
							{s.syncStatus.isRunning ? '🔄 Running' : '✅ Idle'}
						</p>
					</div>
				</div>

				{#if s.syncStatus.lastSyncError}
					<p class="text-xs text-red-400">{s.syncStatus.lastSyncError}</p>
				{/if}

				<Separator class="bg-neutral-800" />

				{#if sonarrProfiles.data && sonarrProfiles.data.length > 0}
					<div>
						<p class="mb-2 text-xs font-medium uppercase text-neutral-500">Quality Profiles</p>
						<div class="flex flex-wrap gap-1">
							{#each sonarrProfiles.data as p}
								<Badge variant="outline" class="text-neutral-300">{p.name}</Badge>
							{/each}
						</div>
					</div>
				{/if}

				{#if sonarrFolders.data && sonarrFolders.data.length > 0}
					<div>
						<p class="mb-2 text-xs font-medium uppercase text-neutral-500">Root Folders</p>
						{#each sonarrFolders.data as f}
							<div class="flex items-center justify-between text-sm">
								<span class="text-neutral-300">{f.path}</span>
								<span class="text-neutral-500">{formatBytes(f.freeSpace)} free</span>
							</div>
						{/each}
					</div>
				{/if}
			{:else if sonarrStatus.isLoading}
				<p class="text-sm text-neutral-500">Loading…</p>
			{:else if sonarrStatus.isError}
				<p class="text-sm text-red-400">Failed to connect to Sonarr</p>
			{/if}
		</Card.Content>
		<Card.Footer>
			<Button size="sm" onclick={handleSonarrSync} disabled={sonarrSyncing}>
				{sonarrSyncing ? 'Syncing…' : 'Sync Now'}
			</Button>
		</Card.Footer>
	</Card.Root>
</div>

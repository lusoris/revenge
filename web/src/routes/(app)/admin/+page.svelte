<script lang="ts">
	import { adminGetActivityStats, adminListUsers, radarrGetStatus, sonarrGetStatus } from '$api/endpoints/admin';
	import { Badge } from '$components/ui/badge';
	import * as Card from '$components/ui/card';
	import { createQuery } from '@tanstack/svelte-query';

	const radarrQuery = createQuery(() => ({
		queryKey: ['admin', 'radarr', 'status'],
		queryFn: () => radarrGetStatus(),
		retry: false
	}));

	const sonarrQuery = createQuery(() => ({
		queryKey: ['admin', 'sonarr', 'status'],
		queryFn: () => sonarrGetStatus(),
		retry: false
	}));

	const usersQuery = createQuery(() => ({
		queryKey: ['admin', 'users', 'count'],
		queryFn: () => adminListUsers({ page_size: 1 })
	}));

	const activityQuery = createQuery(() => ({
		queryKey: ['admin', 'activity', 'stats'],
		queryFn: () => adminGetActivityStats(),
		retry: false
	}));
</script>

<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
	<!-- Users -->
	<Card.Root class="border-neutral-800 bg-neutral-900">
		<Card.Header class="pb-2">
			<Card.Description>Total Users</Card.Description>
		</Card.Header>
		<Card.Content>
			<p class="text-3xl font-bold text-white">
				{usersQuery.data?.total ?? '—'}
			</p>
		</Card.Content>
	</Card.Root>

	<!-- Radarr -->
	<Card.Root class="border-neutral-800 bg-neutral-900">
		<Card.Header class="pb-2">
			<div class="flex items-center justify-between">
				<Card.Description>Radarr</Card.Description>
				{#if radarrQuery.data?.connected}
					<Badge class="bg-green-900 text-green-200">Connected</Badge>
				{:else if radarrQuery.isError}
					<Badge variant="destructive">Error</Badge>
				{:else}
					<Badge variant="outline" class="text-neutral-400">Loading</Badge>
				{/if}
			</div>
		</Card.Header>
		<Card.Content>
			{#if radarrQuery.data}
				<p class="text-3xl font-bold text-white">
					{radarrQuery.data.syncStatus.totalMovies}
				</p>
				<p class="text-xs text-neutral-500">movies · v{radarrQuery.data.version}</p>
			{:else}
				<p class="text-sm text-neutral-500">—</p>
			{/if}
		</Card.Content>
	</Card.Root>

	<!-- Sonarr -->
	<Card.Root class="border-neutral-800 bg-neutral-900">
		<Card.Header class="pb-2">
			<div class="flex items-center justify-between">
				<Card.Description>Sonarr</Card.Description>
				{#if sonarrQuery.data?.connected}
					<Badge class="bg-green-900 text-green-200">Connected</Badge>
				{:else if sonarrQuery.isError}
					<Badge variant="destructive">Error</Badge>
				{:else}
					<Badge variant="outline" class="text-neutral-400">Loading</Badge>
				{/if}
			</div>
		</Card.Header>
		<Card.Content>
			{#if sonarrQuery.data}
				<p class="text-3xl font-bold text-white">
					{sonarrQuery.data.syncStatus.totalSeries}
				</p>
				<p class="text-xs text-neutral-500">series · v{sonarrQuery.data.version}</p>
			{:else}
				<p class="text-sm text-neutral-500">—</p>
			{/if}
		</Card.Content>
	</Card.Root>

	<!-- Activity -->
	<Card.Root class="border-neutral-800 bg-neutral-900">
		<Card.Header class="pb-2">
			<Card.Description>Activity</Card.Description>
		</Card.Header>
		<Card.Content>
			{#if activityQuery.data}
				<p class="text-3xl font-bold text-white">{activityQuery.data.total_count}</p>
				<p class="text-xs text-neutral-500">
					{activityQuery.data.success_count} success · {activityQuery.data.failed_count} failed
				</p>
			{:else}
				<p class="text-sm text-neutral-500">—</p>
			{/if}
		</Card.Content>
	</Card.Root>
</div>

<script lang="ts">
	import { createQuery } from '@tanstack/svelte-query';
	import { derived, writable } from 'svelte/store';
	import { adminGetActivityLogs, adminGetActivityStats } from '$api/endpoints/admin';
	import { Button } from '$components/ui/button';
	import { Input } from '$components/ui/input';
	import * as Card from '$components/ui/card';
	import { Badge } from '$components/ui/badge';

	let currentPage = $state(1);
	let actionFilter = $state('');
	const pageSize = 25;

	const pageStore = writable(1);
	const filterStore = writable('');
	$effect(() => { pageStore.set(currentPage); });
	$effect(() => { filterStore.set(actionFilter); });

	const logsQuery = createQuery(
		derived([pageStore, filterStore], ([$p, $f]) => ({
			queryKey: ['admin', 'activity', 'logs', $p, $f],
			queryFn: () =>
				adminGetActivityLogs({
					page: $p,
					page_size: pageSize,
					action: $f || undefined
				})
		}))
	);

	const statsQuery = createQuery(
		derived(writable(null), () => ({
			queryKey: ['admin', 'activity', 'stats'],
			queryFn: () => adminGetActivityStats(),
			retry: false
		}))
	);

	function actionColor(action: string): string {
		if (action.includes('delete') || action.includes('revoke')) return 'bg-red-900 text-red-200';
		if (action.includes('create') || action.includes('register')) return 'bg-green-900 text-green-200';
		if (action.includes('update') || action.includes('change')) return 'bg-blue-900 text-blue-200';
		if (action.includes('login')) return 'bg-purple-900 text-purple-200';
		return 'bg-neutral-800 text-neutral-300';
	}
</script>

<div class="space-y-4">
	<!-- Stats Row -->
	{#if $statsQuery.data}
		<div class="grid grid-cols-3 gap-4">
			<Card.Root class="border-neutral-800 bg-neutral-900">
				<Card.Content class="py-4 text-center">
					<p class="text-2xl font-bold text-white">{$statsQuery.data.total_count}</p>
					<p class="text-xs text-neutral-500">Total Events</p>
				</Card.Content>
			</Card.Root>
			<Card.Root class="border-neutral-800 bg-neutral-900">
				<Card.Content class="py-4 text-center">
					<p class="text-2xl font-bold text-green-400">{$statsQuery.data.success_count}</p>
					<p class="text-xs text-neutral-500">Successful</p>
				</Card.Content>
			</Card.Root>
			<Card.Root class="border-neutral-800 bg-neutral-900">
				<Card.Content class="py-4 text-center">
					<p class="text-2xl font-bold text-red-400">{$statsQuery.data.failed_count}</p>
					<p class="text-xs text-neutral-500">Failed</p>
				</Card.Content>
			</Card.Root>
		</div>
	{/if}

	<!-- Logs Table -->
	<Card.Root class="border-neutral-800 bg-neutral-900">
		<Card.Header class="flex flex-row items-center justify-between gap-4">
			<Card.Title class="text-white">Activity Log</Card.Title>
			<Input
				placeholder="Filter by action…"
				bind:value={actionFilter}
				class="max-w-xs bg-neutral-800"
			/>
		</Card.Header>
		<Card.Content>
			{#if $logsQuery.isLoading}
				<p class="text-sm text-neutral-500">Loading…</p>
			{:else if $logsQuery.data}
				<div class="overflow-x-auto">
					<table class="w-full text-sm">
						<thead>
							<tr class="border-b border-neutral-800 text-left text-neutral-400">
								<th class="pb-3 pr-4 font-medium">Action</th>
								<th class="pb-3 pr-4 font-medium">User</th>
								<th class="pb-3 pr-4 font-medium">Resource</th>
								<th class="pb-3 pr-4 font-medium">Status</th>
								<th class="pb-3 pr-4 font-medium">IP</th>
								<th class="pb-3 font-medium">Time</th>
							</tr>
						</thead>
						<tbody class="divide-y divide-neutral-800">
							{#each $logsQuery.data.entries as entry}
								<tr>
									<td class="py-2 pr-4">
										<Badge class={actionColor(entry.action)}>{entry.action}</Badge>
									</td>
									<td class="py-2 pr-4 text-neutral-300">{entry.username ?? entry.user_id ?? '—'}</td>
									<td class="py-2 pr-4">
										{#if entry.resource_type}
											<span class="text-neutral-400">
												{entry.resource_type}
												{#if entry.resource_id}
													<span class="text-neutral-600">({entry.resource_id.slice(0, 8)}…)</span>
												{/if}
											</span>
										{:else}
											<span class="text-neutral-600">—</span>
										{/if}
									</td>
									<td class="py-2 pr-4">
										{#if entry.success}
											<span class="text-green-400">✓</span>
										{:else}
											<span class="text-red-400" title={entry.error_message}>✗</span>
										{/if}
									</td>
									<td class="py-2 pr-4 text-neutral-500">{entry.ip_address ?? '—'}</td>
									<td class="py-2 text-neutral-500">
										{new Date(entry.created_at).toLocaleString()}
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>

				{#if $logsQuery.data.entries.length === 0}
					<p class="py-8 text-center text-sm text-neutral-500">No activity logs found</p>
				{/if}

				<!-- Pagination -->
				{#if ($logsQuery.data.total ?? 0) > pageSize}
					<div class="mt-4 flex items-center justify-between">
						<p class="text-sm text-neutral-500">
							Page {currentPage} of {Math.ceil(($logsQuery.data.total ?? 0) / pageSize)}
						</p>
						<div class="flex gap-2">
							<Button
								variant="outline"
								size="sm"
								disabled={currentPage <= 1}
								onclick={() => (currentPage = Math.max(1, currentPage - 1))}
							>
								Previous
							</Button>
							<Button
								variant="outline"
								size="sm"
								disabled={currentPage >=
									Math.ceil(($logsQuery.data.total ?? 0) / pageSize)}
								onclick={() => currentPage++}
							>
								Next
							</Button>
						</div>
					</div>
				{/if}
			{/if}
		</Card.Content>
	</Card.Root>
</div>

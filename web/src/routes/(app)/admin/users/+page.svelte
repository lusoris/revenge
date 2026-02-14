<script lang="ts">
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { derived, writable } from 'svelte/store';
	import { adminListUsers, adminDeleteUser } from '$api/endpoints/admin';
	import { assignUserRole, removeUserRole, getUserRoles } from '$api/endpoints/rbac';
	import { Button } from '$components/ui/button';
	import { Input } from '$components/ui/input';
	import * as Card from '$components/ui/card';
	import { Badge } from '$components/ui/badge';
	import type { User } from '$api/types';

	const queryClient = useQueryClient();

	let search = $state('');
	let currentPage = $state(1);
	const pageSize = 20;

	const searchStore = writable('');
	const pageStore = writable(1);
	$effect(() => { searchStore.set(search); });
	$effect(() => { pageStore.set(currentPage); });

	const usersQuery = createQuery(
		derived([searchStore, pageStore], ([$s, $p]) => ({
			queryKey: ['admin', 'users', $s, $p],
			queryFn: () =>
				adminListUsers({
					page: $p,
					page_size: pageSize,
					search: $s || undefined
				})
		}))
	);

	async function handleDelete(user: User) {
		if (!confirm(`Delete user ${user.username}? This cannot be undone.`)) return;
		try {
			await adminDeleteUser(user.id);
			queryClient.invalidateQueries({ queryKey: ['admin', 'users'] });
		} catch (err: any) {
			alert(err.message);
		}
	}

	async function toggleAdmin(user: User) {
		try {
			if (user.is_admin) {
				await removeUserRole(user.id, 'admin');
			} else {
				await assignUserRole(user.id, 'admin');
			}
			queryClient.invalidateQueries({ queryKey: ['admin', 'users'] });
		} catch (err: any) {
			alert(err.message);
		}
	}
</script>

<Card.Root class="border-neutral-800 bg-neutral-900">
	<Card.Header class="flex flex-row items-center justify-between gap-4">
		<div>
			<Card.Title class="text-white">Users</Card.Title>
			<Card.Description>
				{$usersQuery.data?.total ?? 0} total users
			</Card.Description>
		</div>
		<Input
			placeholder="Search users…"
			bind:value={search}
			class="max-w-xs bg-neutral-800"
		/>
	</Card.Header>
	<Card.Content>
		{#if $usersQuery.isLoading}
			<p class="text-sm text-neutral-500">Loading users…</p>
		{:else if $usersQuery.data}
			<div class="overflow-x-auto">
				<table class="w-full text-sm">
					<thead>
						<tr class="border-b border-neutral-800 text-left text-neutral-400">
							<th class="pb-3 pr-4 font-medium">User</th>
							<th class="pb-3 pr-4 font-medium">Email</th>
							<th class="pb-3 pr-4 font-medium">Status</th>
							<th class="pb-3 pr-4 font-medium">Role</th>
							<th class="pb-3 font-medium">Created</th>
							<th class="pb-3 font-medium"></th>
						</tr>
					</thead>
					<tbody class="divide-y divide-neutral-800">
						{#each $usersQuery.data.users as user}
							<tr>
								<td class="py-3 pr-4">
									<div class="flex items-center gap-3">
										<div
											class="flex h-8 w-8 items-center justify-center rounded-full bg-neutral-800 text-xs font-medium text-neutral-300"
										>
											{(user.display_name ?? user.username).charAt(0).toUpperCase()}
										</div>
										<div>
											<p class="font-medium text-white">{user.display_name ?? user.username}</p>
											<p class="text-xs text-neutral-500">@{user.username}</p>
										</div>
									</div>
								</td>
								<td class="py-3 pr-4 text-neutral-300">{user.email}</td>
								<td class="py-3 pr-4">
									{#if user.is_active}
										<Badge class="bg-green-900 text-green-200">Active</Badge>
									{:else}
										<Badge variant="destructive">Disabled</Badge>
									{/if}
								</td>
								<td class="py-3 pr-4">
									{#if user.is_admin}
										<Badge class="bg-purple-900 text-purple-200">Admin</Badge>
									{:else}
										<Badge variant="outline" class="text-neutral-400">User</Badge>
									{/if}
								</td>
								<td class="py-3 text-neutral-400">
									{new Date(user.created_at).toLocaleDateString()}
								</td>
								<td class="py-3 text-right">
									<div class="flex justify-end gap-2">
										<Button variant="outline" size="sm" onclick={() => toggleAdmin(user)}>
											{user.is_admin ? 'Remove Admin' : 'Make Admin'}
										</Button>
										<Button variant="destructive" size="sm" onclick={() => handleDelete(user)}>
											Delete
										</Button>
									</div>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>

			<!-- Pagination -->
			{#if ($usersQuery.data.total ?? 0) > pageSize}
				<div class="mt-4 flex items-center justify-between">
					<p class="text-sm text-neutral-500">
						Page {currentPage} of {Math.ceil(($usersQuery.data.total ?? 0) / pageSize)}
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
							disabled={currentPage >= Math.ceil(($usersQuery.data.total ?? 0) / pageSize)}
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

<script lang="ts">
	import { createRole, deleteRole, listRoles } from '$api/endpoints/rbac';
	import { Badge } from '$components/ui/badge';
	import { Button } from '$components/ui/button';
	import * as Card from '$components/ui/card';
	import * as Dialog from '$components/ui/dialog';
	import { Input } from '$components/ui/input';
	import { Label } from '$components/ui/label';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';

	const queryClient = useQueryClient();

	const rolesQuery = createQuery(() => ({
		queryKey: ['admin', 'roles'],
		queryFn: () => listRoles()
	}));

	let showCreate = $state(false);
	let roleName = $state('');
	let roleDescription = $state('');
	let error = $state('');

	async function handleCreate() {
		if (!roleName.trim()) {
			error = 'Name is required';
			return;
		}
		try {
			await createRole({ name: roleName, description: roleDescription || undefined });
			queryClient.invalidateQueries({ queryKey: ['admin', 'roles'] });
			roleName = '';
			roleDescription = '';
			showCreate = false;
			error = '';
		} catch (err: any) {
			error = err.message;
		}
	}

	async function handleDelete(name: string) {
		if (!confirm(`Delete role "${name}"?`)) return;
		try {
			await deleteRole(name);
			queryClient.invalidateQueries({ queryKey: ['admin', 'roles'] });
		} catch (err: any) {
			alert(err.message);
		}
	}
</script>

<Card.Root class="border-neutral-800 bg-neutral-900">
	<Card.Header class="flex flex-row items-center justify-between">
		<div>
			<Card.Title class="text-white">Roles & Permissions</Card.Title>
			<Card.Description>
				{rolesQuery.data?.total ?? 0} roles configured
			</Card.Description>
		</div>
		<Dialog.Root bind:open={showCreate}>
			<Dialog.Trigger>
				{#snippet children({ props })}
					<Button size="sm" {...props}>Create Role</Button>
				{/snippet}
			</Dialog.Trigger>
			<Dialog.Content class="border-neutral-800 bg-neutral-900">
				<Dialog.Header>
					<Dialog.Title>Create Role</Dialog.Title>
					<Dialog.Description>Add a new role to the system</Dialog.Description>
				</Dialog.Header>
				<div class="space-y-4 py-4">
					{#if error}
						<p class="text-sm text-red-400">{error}</p>
					{/if}
					<div class="space-y-2">
						<Label>Name</Label>
						<Input bind:value={roleName} placeholder="editor" class="bg-neutral-800" />
					</div>
					<div class="space-y-2">
						<Label>Description</Label>
						<Input
							bind:value={roleDescription}
							placeholder="Can edit content"
							class="bg-neutral-800"
						/>
					</div>
				</div>
				<Dialog.Footer>
					<Button onclick={handleCreate}>Create</Button>
				</Dialog.Footer>
			</Dialog.Content>
		</Dialog.Root>
	</Card.Header>
	<Card.Content>
		{#if rolesQuery.isLoading}
			<p class="text-sm text-neutral-500">Loading roles…</p>
		{:else if rolesQuery.data}
			<div class="divide-y divide-neutral-800">
				{#each rolesQuery.data.roles as role}
					<div class="flex items-center justify-between py-3">
						<div>
							<div class="flex items-center gap-2">
								<p class="text-sm font-medium text-white">{role.name}</p>
								{#if role.is_built_in}
									<Badge variant="outline" class="text-neutral-400">Built-in</Badge>
								{/if}
							</div>
							{#if role.description}
								<p class="text-xs text-neutral-500">{role.description}</p>
							{/if}
							<p class="mt-0.5 text-xs text-neutral-500">
								{role.permissions.length} permissions · {role.user_count} users
							</p>
						</div>
						{#if !role.is_built_in}
							<Button
								variant="destructive"
								size="sm"
								onclick={() => handleDelete(role.name)}
							>
								Delete
							</Button>
						{/if}
					</div>
				{/each}
			</div>
			{#if rolesQuery.data.roles.length === 0}
				<p class="py-8 text-center text-sm text-neutral-500">No roles configured</p>
			{/if}
		{/if}
	</Card.Content>
</Card.Root>

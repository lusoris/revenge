<script lang="ts">
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { derived, writable } from 'svelte/store';
	import { listAPIKeys, createAPIKey, revokeAPIKey } from '$api/endpoints/apikeys';
	import { Button } from '$components/ui/button';
	import { Input } from '$components/ui/input';
	import { Label } from '$components/ui/label';
	import * as Card from '$components/ui/card';
	import * as Dialog from '$components/ui/dialog';
	import { Badge } from '$components/ui/badge';
	import type { APIKeyScope, CreateAPIKeyResponse } from '$api/types';

	const queryClient = useQueryClient();

	const keysQuery = createQuery(
		derived(writable(null), () => ({
			queryKey: ['apikeys'],
			queryFn: () => listAPIKeys()
		}))
	);

	let showCreate = $state(false);
	let keyName = $state('');
	let keyDescription = $state('');
	let selectedScopes: APIKeyScope[] = $state(['read']);
	let newKey: CreateAPIKeyResponse | null = $state(null);
	let error = $state('');

	async function handleCreate() {
		if (!keyName.trim()) {
			error = 'Name is required';
			return;
		}
		try {
			newKey = await createAPIKey({
				name: keyName,
				description: keyDescription || undefined,
				scopes: selectedScopes
			});
			queryClient.invalidateQueries({ queryKey: ['apikeys'] });
			error = '';
		} catch (err: any) {
			error = err.message;
		}
	}

	async function handleRevoke(id: string) {
		if (!confirm('Revoke this API key? This cannot be undone.')) return;
		try {
			await revokeAPIKey(id);
			queryClient.invalidateQueries({ queryKey: ['apikeys'] });
		} catch (err: any) {
			alert(err.message);
		}
	}

	function toggleScope(scope: APIKeyScope) {
		if (selectedScopes.includes(scope)) {
			selectedScopes = selectedScopes.filter((s) => s !== scope);
		} else {
			selectedScopes = [...selectedScopes, scope];
		}
	}

	function resetCreate() {
		keyName = '';
		keyDescription = '';
		selectedScopes = ['read'];
		newKey = null;
		error = '';
		showCreate = false;
	}
</script>

<Card.Root class="border-neutral-800 bg-neutral-900">
	<Card.Header class="flex flex-row items-center justify-between">
		<div>
			<Card.Title class="text-white">API Keys</Card.Title>
			<Card.Description>Manage personal API keys for programmatic access</Card.Description>
		</div>
		<Dialog.Root bind:open={showCreate} onOpenChange={(open) => { if (!open) resetCreate(); }}>
			<Dialog.Trigger>
				{#snippet children({ props })}
					<Button size="sm" {...props}>Create Key</Button>
				{/snippet}
			</Dialog.Trigger>
			<Dialog.Content class="border-neutral-800 bg-neutral-900">
				<Dialog.Header>
					<Dialog.Title>
						{newKey ? 'API Key Created' : 'Create API Key'}
					</Dialog.Title>
					<Dialog.Description>
						{newKey
							? 'Copy this key now. It will not be shown again.'
							: 'Create a new API key for programmatic access.'}
					</Dialog.Description>
				</Dialog.Header>

				{#if newKey}
					<div class="space-y-4 py-4">
						<code
							class="block break-all rounded bg-neutral-800 px-4 py-3 text-sm text-green-300"
						>
							{newKey.api_key}
						</code>
						<Button
							variant="outline"
							size="sm"
							onclick={() => navigator.clipboard.writeText(newKey?.api_key ?? '')}
						>
							Copy to Clipboard
						</Button>
					</div>
				{:else}
					<div class="space-y-4 py-4">
						{#if error}
							<p class="text-sm text-red-400">{error}</p>
						{/if}
						<div class="space-y-2">
							<Label>Name</Label>
							<Input bind:value={keyName} placeholder="My API Key" class="bg-neutral-800" />
						</div>
						<div class="space-y-2">
							<Label>Description (optional)</Label>
							<Input
								bind:value={keyDescription}
								placeholder="What is this key for?"
								class="bg-neutral-800"
							/>
						</div>
						<div class="space-y-2">
							<Label>Scopes</Label>
							<div class="flex gap-2">
								{#each ['read', 'write', 'admin'] as scope}
									<button
										class="rounded-md border px-3 py-1 text-sm transition-colors {selectedScopes.includes(scope as APIKeyScope)
											? 'border-white bg-white text-black'
											: 'border-neutral-700 text-neutral-400 hover:border-neutral-500'}"
										onclick={() => toggleScope(scope as APIKeyScope)}
									>
										{scope}
									</button>
								{/each}
							</div>
						</div>
					</div>
				{/if}

				<Dialog.Footer>
					{#if newKey}
						<Button onclick={resetCreate}>Done</Button>
					{:else}
						<Button onclick={handleCreate}>Create</Button>
					{/if}
				</Dialog.Footer>
			</Dialog.Content>
		</Dialog.Root>
	</Card.Header>
	<Card.Content>
		{#if $keysQuery.isLoading}
			<p class="text-sm text-neutral-500">Loading…</p>
		{:else if $keysQuery.data}
			<div class="divide-y divide-neutral-800">
				{#each $keysQuery.data.keys as key}
					<div class="flex items-center justify-between py-3">
						<div class="min-w-0 flex-1">
							<div class="flex items-center gap-2">
								<p class="text-sm font-medium text-white">{key.name}</p>
								{#each key.scopes as scope}
									<Badge variant="outline" class="text-xs text-neutral-400">{scope}</Badge>
								{/each}
								{#if !key.is_active}
									<Badge variant="destructive" class="text-xs">Revoked</Badge>
								{/if}
							</div>
							<p class="text-xs text-neutral-500">
								{key.key_prefix}···
								· Created {new Date(key.created_at).toLocaleDateString()}
								{#if key.last_used_at}
									· Last used {new Date(key.last_used_at).toLocaleDateString()}
								{/if}
								{#if key.expires_at}
									· Expires {new Date(key.expires_at).toLocaleDateString()}
								{/if}
							</p>
						</div>
						{#if key.is_active}
							<Button
								variant="destructive"
								size="sm"
								onclick={() => handleRevoke(key.id)}
							>
								Revoke
							</Button>
						{/if}
					</div>
				{/each}
			</div>
			{#if $keysQuery.data.keys.length === 0}
				<p class="py-8 text-center text-sm text-neutral-500">No API keys yet</p>
			{/if}
		{/if}
	</Card.Content>
</Card.Root>

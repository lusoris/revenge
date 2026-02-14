<script lang="ts">
	import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
	import { derived, writable } from 'svelte/store';
	import { getServerSettings, updateServerSetting } from '$api/endpoints/settings';
	import { Button } from '$components/ui/button';
	import { Input } from '$components/ui/input';
	import { Label } from '$components/ui/label';
	import { Switch } from '$components/ui/switch';
	import * as Card from '$components/ui/card';
	import * as Alert from '$components/ui/alert';
	import type { ServerSetting } from '$api/types';

	const queryClient = useQueryClient();

	const settingsQuery = createQuery(
		derived(writable(null), () => ({
			queryKey: ['admin', 'settings'],
			queryFn: () => getServerSettings()
		}))
	);

	let editingKey = $state<string | null>(null);
	let editValue = $state('');
	let successMsg = $state('');
	let errorMsg = $state('');

	function startEdit(setting: ServerSetting) {
		editingKey = setting.key;
		editValue = String(setting.value);
	}

	async function save(setting: ServerSetting) {
		try {
			let parsed: string | number | boolean = editValue;
			if (setting.data_type === 'integer') parsed = parseInt(editValue);
			if (setting.data_type === 'float') parsed = parseFloat(editValue);
			if (setting.data_type === 'boolean') parsed = editValue === 'true';

			await updateServerSetting(setting.key, parsed);
			queryClient.invalidateQueries({ queryKey: ['admin', 'settings'] });
			editingKey = null;
			successMsg = `Updated "${setting.key}"`;
			errorMsg = '';
		} catch (err: any) {
			errorMsg = err.message;
		}
	}

	function groupByCategory(settings: ServerSetting[]): Record<string, ServerSetting[]> {
		const groups: Record<string, ServerSetting[]> = {};
		for (const s of settings) {
			const cat = s.category ?? 'General';
			if (!groups[cat]) groups[cat] = [];
			groups[cat].push(s);
		}
		return groups;
	}
</script>

<div class="space-y-6">
	{#if successMsg}
		<Alert.Root class="border-green-800 bg-green-950 text-green-200">
			<Alert.Title>{successMsg}</Alert.Title>
		</Alert.Root>
	{/if}
	{#if errorMsg}
		<Alert.Root class="border-red-800 bg-red-950 text-red-200">
			<Alert.Title>{errorMsg}</Alert.Title>
		</Alert.Root>
	{/if}

	{#if $settingsQuery.isLoading}
		<p class="text-sm text-neutral-500">Loading settings…</p>
	{:else if $settingsQuery.data}
		{@const groups = groupByCategory($settingsQuery.data)}
		{#each Object.entries(groups) as [category, settings]}
			<Card.Root class="border-neutral-800 bg-neutral-900">
				<Card.Header>
					<Card.Title class="text-white">{category}</Card.Title>
				</Card.Header>
				<Card.Content>
					<div class="divide-y divide-neutral-800">
						{#each settings as setting}
							<div class="flex items-center justify-between py-3">
								<div class="min-w-0 flex-1 pr-4">
									<p class="text-sm font-medium text-white">{setting.key}</p>
									{#if setting.description}
										<p class="text-xs text-neutral-500">{setting.description}</p>
									{/if}
								</div>
								<div class="flex items-center gap-2">
									{#if editingKey === setting.key}
										{#if setting.data_type === 'boolean'}
											<Switch
												checked={editValue === 'true'}
												onCheckedChange={(v) => (editValue = String(v))}
											/>
										{:else}
											<Input
												bind:value={editValue}
												class="w-48 bg-neutral-800"
												type={setting.data_type === 'integer' || setting.data_type === 'float'
													? 'number'
													: 'text'}
											/>
										{/if}
										<Button size="sm" onclick={() => save(setting)}>Save</Button>
										<Button
											variant="outline"
											size="sm"
											onclick={() => (editingKey = null)}
										>
											Cancel
										</Button>
									{:else}
										<span class="text-sm text-neutral-300">
											{#if setting.is_secret}
												••••••
											{:else}
												{String(setting.value)}
											{/if}
										</span>
										<Button variant="outline" size="sm" onclick={() => startEdit(setting)}>
											Edit
										</Button>
									{/if}
								</div>
							</div>
						{/each}
					</div>
				</Card.Content>
			</Card.Root>
		{/each}

		{#if $settingsQuery.data.length === 0}
			<Card.Root class="border-neutral-800 bg-neutral-900">
				<Card.Content>
					<p class="py-8 text-center text-sm text-neutral-500">No server settings configured</p>
				</Card.Content>
			</Card.Root>
		{/if}
	{/if}
</div>

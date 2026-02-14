<script lang="ts">
	import { getPreferences, updatePreferences } from '$api/endpoints/users';
	import type { UserPreferences } from '$api/types';
	import * as Alert from '$components/ui/alert';
	import { Button } from '$components/ui/button';
	import * as Card from '$components/ui/card';
	import { Label } from '$components/ui/label';
	import { Switch } from '$components/ui/switch';
	import { createMutation, createQuery, useQueryClient } from '@tanstack/svelte-query';

	const queryClient = useQueryClient();

	const prefsQuery = createQuery(() => ({
		queryKey: ['preferences'],
		queryFn: () => getPreferences()
	}));

	let theme = $state<'light' | 'dark' | 'system'>('dark');
	let autoPlay = $state(true);
	let showAdult = $state(false);
	let showSpoilers = $state(false);
	let profileVisibility = $state<'public' | 'friends' | 'private'>('private');
	let successMsg = $state('');
	let errorMsg = $state('');

	$effect(() => {
		if (prefsQuery.data) {
			theme = prefsQuery.data.theme ?? 'dark';
			autoPlay = prefsQuery.data.auto_play_videos ?? true;
			showAdult = prefsQuery.data.show_adult_content ?? false;
			showSpoilers = prefsQuery.data.show_spoilers ?? false;
			profileVisibility = prefsQuery.data.profile_visibility ?? 'private';
		}
	});

	const updateMutation = createMutation(() => ({
		mutationFn: (data: Partial<UserPreferences>) => updatePreferences(data),
		onSuccess: () => {
			successMsg = 'Preferences saved';
			errorMsg = '';
			queryClient.invalidateQueries({ queryKey: ['preferences'] });
		},
		onError: (err: Error) => {
			errorMsg = err.message;
			successMsg = '';
		}
	}));

	function save() {
		updateMutation.mutate({
			theme,
			auto_play_videos: autoPlay,
			show_adult_content: showAdult,
			show_spoilers: showSpoilers,
			profile_visibility: profileVisibility
		});
	}
</script>

<div class="grid gap-6 lg:grid-cols-2">
	<Card.Root class="border-neutral-800 bg-neutral-900">
		<Card.Header>
			<Card.Title class="text-white">Display</Card.Title>
			<Card.Description>Theme and playback preferences</Card.Description>
		</Card.Header>
		<Card.Content class="space-y-4">
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

			<div class="space-y-2">
				<Label>Theme</Label>
				<div class="flex gap-2">
					{#each ['dark', 'light', 'system'] as t}
						<button
							class="rounded-md border px-4 py-2 text-sm transition-colors {theme === t
								? 'border-white bg-white text-black'
								: 'border-neutral-700 text-neutral-400 hover:border-neutral-500'}"
							onclick={() => (theme = t as typeof theme)}
						>
							{t.charAt(0).toUpperCase() + t.slice(1)}
						</button>
					{/each}
				</div>
			</div>

			<div class="flex items-center justify-between">
				<Label for="autoPlay">Auto-play videos</Label>
				<Switch id="autoPlay" checked={autoPlay} onCheckedChange={(v) => (autoPlay = v)} />
			</div>

			<div class="flex items-center justify-between">
				<Label for="showAdult">Show adult content</Label>
				<Switch id="showAdult" checked={showAdult} onCheckedChange={(v) => (showAdult = v)} />
			</div>

			<div class="flex items-center justify-between">
				<Label for="showSpoilers">Show spoilers</Label>
				<Switch
					id="showSpoilers"
					checked={showSpoilers}
					onCheckedChange={(v) => (showSpoilers = v)}
				/>
			</div>
		</Card.Content>
		<Card.Footer>
		<Button onclick={save} disabled={updateMutation.isPending}>
			{updateMutation.isPending ? 'Saving…' : 'Save Preferences'}
			</Button>
		</Card.Footer>
	</Card.Root>

	<Card.Root class="border-neutral-800 bg-neutral-900">
		<Card.Header>
			<Card.Title class="text-white">Privacy</Card.Title>
			<Card.Description>Control who can see your profile</Card.Description>
		</Card.Header>
		<Card.Content class="space-y-4">
			<div class="space-y-2">
				<Label>Profile Visibility</Label>
				<div class="flex gap-2">
					{#each [{ v: 'public', l: 'Public' }, { v: 'friends', l: 'Friends' }, { v: 'private', l: 'Private' }] as opt}
						<button
							class="rounded-md border px-4 py-2 text-sm transition-colors {profileVisibility === opt.v
								? 'border-white bg-white text-black'
								: 'border-neutral-700 text-neutral-400 hover:border-neutral-500'}"
							onclick={() => (profileVisibility = opt.v as typeof profileVisibility)}
						>
							{opt.l}
						</button>
					{/each}
				</div>
			</div>
		</Card.Content>
	</Card.Root>
</div>

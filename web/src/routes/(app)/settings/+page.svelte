<script lang="ts">
	import { changePassword, getCurrentUser, updateCurrentUser, uploadAvatar } from '$api/endpoints/users';
	import * as Alert from '$components/ui/alert';
	import { Button } from '$components/ui/button';
	import * as Card from '$components/ui/card';
	import { Input } from '$components/ui/input';
	import { Label } from '$components/ui/label';
	import { getAuth } from '$lib/stores/auth.svelte';
	import { createMutation, createQuery, useQueryClient } from '@tanstack/svelte-query';

	const auth = getAuth();
	const queryClient = useQueryClient();

	const userQuery = createQuery(() => ({
		queryKey: ['user', 'me'],
		queryFn: () => getCurrentUser()
	}));

	let displayName = $state(auth.user?.display_name ?? '');
	let email = $state(auth.user?.email ?? '');
	let timezone = $state(auth.user?.timezone ?? '');
	let successMsg = $state('');
	let errorMsg = $state('');

	let oldPassword = $state('');
	let newPassword = $state('');
	let confirmPassword = $state('');
	let pwSuccess = $state('');
	let pwError = $state('');

	$effect(() => {
		if (userQuery.data) {
			displayName = userQuery.data.display_name ?? '';
			email = userQuery.data.email ?? '';
			timezone = userQuery.data.timezone ?? '';
		}
	});

	const updateMutation = createMutation(() => ({
		mutationFn: (data: { email?: string; display_name?: string; timezone?: string }) =>
			updateCurrentUser(data),
		onSuccess: () => {
			successMsg = 'Profile updated successfully';
			errorMsg = '';
			queryClient.invalidateQueries({ queryKey: ['user', 'me'] });
		},
		onError: (err: Error) => {
			errorMsg = err.message;
			successMsg = '';
		}
	}));

	const passwordMutation = createMutation(() => ({
		mutationFn: (data: { old_password: string; new_password: string }) =>
			changePassword(data),
		onSuccess: () => {
			pwSuccess = 'Password changed successfully';
			pwError = '';
			oldPassword = '';
			newPassword = '';
			confirmPassword = '';
		},
		onError: (err: Error) => {
			pwError = err.message;
			pwSuccess = '';
		}
	}));

	function saveProfile() {
		updateMutation.mutate({ display_name: displayName, email, timezone });
	}

	function savePassword() {
		if (newPassword !== confirmPassword) {
			pwError = 'Passwords do not match';
			return;
		}
		if (newPassword.length < 8) {
			pwError = 'Password must be at least 8 characters';
			return;
		}
		passwordMutation.mutate({ old_password: oldPassword, new_password: newPassword });
	}

	let avatarInput: HTMLInputElement;

	async function handleAvatarUpload(e: Event) {
		const input = e.target as HTMLInputElement;
		const file = input.files?.[0];
		if (!file) return;
		try {
			await uploadAvatar(file);
			queryClient.invalidateQueries({ queryKey: ['user', 'me'] });
			successMsg = 'Avatar updated';
		} catch (err: any) {
			errorMsg = err.message ?? 'Failed to upload avatar';
		}
	}
</script>

<div class="grid gap-6 lg:grid-cols-2">
	<!-- Profile Card -->
	<Card.Root class="border-neutral-800 bg-neutral-900">
		<Card.Header>
			<Card.Title class="text-white">Profile</Card.Title>
			<Card.Description>Update your display name, email and avatar</Card.Description>
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

			<!-- Avatar -->
			<div class="flex items-center gap-4">
				{#if auth.user?.avatar_url}
					<img
						src={auth.user.avatar_url}
						alt="Avatar"
						class="h-16 w-16 rounded-full object-cover"
					/>
				{:else}
					<div
						class="flex h-16 w-16 items-center justify-center rounded-full bg-neutral-800 text-xl font-bold text-neutral-300"
					>
						{(auth.user?.display_name ?? auth.user?.username ?? '?').charAt(0).toUpperCase()}
					</div>
				{/if}
				<div>
					<input
						bind:this={avatarInput}
						type="file"
						accept="image/*"
						class="hidden"
						onchange={handleAvatarUpload}
					/>
					<Button variant="outline" size="sm" onclick={() => avatarInput.click()}>
						Change Avatar
					</Button>
				</div>
			</div>

			<div class="space-y-2">
				<Label for="username">Username</Label>
				<Input id="username" value={auth.user?.username ?? ''} disabled class="bg-neutral-800" />
			</div>

			<div class="space-y-2">
				<Label for="displayName">Display Name</Label>
				<Input id="displayName" bind:value={displayName} class="bg-neutral-800" />
			</div>

			<div class="space-y-2">
				<Label for="email">Email</Label>
				<Input id="email" type="email" bind:value={email} class="bg-neutral-800" />
			</div>

			<div class="space-y-2">
				<Label for="timezone">Timezone</Label>
				<Input id="timezone" bind:value={timezone} placeholder="America/New_York" class="bg-neutral-800" />
			</div>
		</Card.Content>
		<Card.Footer>
		<Button onclick={saveProfile} disabled={updateMutation.isPending}>
			{updateMutation.isPending ? 'Saving…' : 'Save Changes'}
			</Button>
		</Card.Footer>
	</Card.Root>

	<!-- Password Card -->
	<Card.Root class="border-neutral-800 bg-neutral-900">
		<Card.Header>
			<Card.Title class="text-white">Change Password</Card.Title>
			<Card.Description>Update your password</Card.Description>
		</Card.Header>
		<Card.Content class="space-y-4">
			{#if pwSuccess}
				<Alert.Root class="border-green-800 bg-green-950 text-green-200">
					<Alert.Title>{pwSuccess}</Alert.Title>
				</Alert.Root>
			{/if}
			{#if pwError}
				<Alert.Root class="border-red-800 bg-red-950 text-red-200">
					<Alert.Title>{pwError}</Alert.Title>
				</Alert.Root>
			{/if}

			<div class="space-y-2">
				<Label for="oldPassword">Current Password</Label>
				<Input id="oldPassword" type="password" bind:value={oldPassword} class="bg-neutral-800" />
			</div>

			<div class="space-y-2">
				<Label for="newPassword">New Password</Label>
				<Input id="newPassword" type="password" bind:value={newPassword} class="bg-neutral-800" />
			</div>

			<div class="space-y-2">
				<Label for="confirmPassword">Confirm Password</Label>
				<Input
					id="confirmPassword"
					type="password"
					bind:value={confirmPassword}
					class="bg-neutral-800"
				/>
			</div>
		</Card.Content>
		<Card.Footer>
		<Button onclick={savePassword} disabled={passwordMutation.isPending}>
			{passwordMutation.isPending ? 'Changing…' : 'Change Password'}
			</Button>
		</Card.Footer>
	</Card.Root>
</div>

<script lang="ts">
	import { goto } from '$app/navigation';
	import { register } from '$lib/stores/auth.svelte';

	let username = $state('');
	let email = $state('');
	let password = $state('');
	let confirmPassword = $state('');
	let submitting = $state(false);
	let errorMsg = $state('');
	let success = $state(false);

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		errorMsg = '';

		if (password !== confirmPassword) {
			errorMsg = 'Passwords do not match.';
			return;
		}

		if (password.length < 8) {
			errorMsg = 'Password must be at least 8 characters.';
			return;
		}

		submitting = true;
		try {
			await register(username, email, password);
			success = true;
			// Redirect to login after short delay
			setTimeout(() => goto('/login'), 1500);
		} catch (err: unknown) {
			errorMsg = err instanceof Error ? err.message : 'Registration failed';
		} finally {
			submitting = false;
		}
	}
</script>

<form onsubmit={handleSubmit} class="space-y-4">
	{#if errorMsg}
		<div class="rounded-md bg-red-500/10 px-4 py-3 text-sm text-red-400">
			{errorMsg}
		</div>
	{/if}

	{#if success}
		<div class="rounded-md bg-green-500/10 px-4 py-3 text-sm text-green-400">
			Account created! Redirecting to login…
		</div>
	{:else}
		<div>
			<label for="username" class="mb-1 block text-sm font-medium text-neutral-300">
				Username
			</label>
			<input
				id="username"
				type="text"
				bind:value={username}
				required
				autocomplete="username"
				class="w-full rounded-lg border border-neutral-800 bg-neutral-900 px-3 py-2 text-sm text-white placeholder-neutral-600 outline-none focus:border-neutral-600 focus:ring-1 focus:ring-neutral-600"
				placeholder="johndoe"
			/>
		</div>

		<div>
			<label for="email" class="mb-1 block text-sm font-medium text-neutral-300">Email</label>
			<input
				id="email"
				type="email"
				bind:value={email}
				required
				autocomplete="email"
				class="w-full rounded-lg border border-neutral-800 bg-neutral-900 px-3 py-2 text-sm text-white placeholder-neutral-600 outline-none focus:border-neutral-600 focus:ring-1 focus:ring-neutral-600"
				placeholder="you@example.com"
			/>
		</div>

		<div>
			<label for="password" class="mb-1 block text-sm font-medium text-neutral-300">
				Password
			</label>
			<input
				id="password"
				type="password"
				bind:value={password}
				required
				autocomplete="new-password"
				minlength="8"
				class="w-full rounded-lg border border-neutral-800 bg-neutral-900 px-3 py-2 text-sm text-white placeholder-neutral-600 outline-none focus:border-neutral-600 focus:ring-1 focus:ring-neutral-600"
				placeholder="••••••••"
			/>
		</div>

		<div>
			<label for="confirm" class="mb-1 block text-sm font-medium text-neutral-300">
				Confirm Password
			</label>
			<input
				id="confirm"
				type="password"
				bind:value={confirmPassword}
				required
				autocomplete="new-password"
				class="w-full rounded-lg border border-neutral-800 bg-neutral-900 px-3 py-2 text-sm text-white placeholder-neutral-600 outline-none focus:border-neutral-600 focus:ring-1 focus:ring-neutral-600"
				placeholder="••••••••"
			/>
		</div>

		<button
			type="submit"
			disabled={submitting}
			class="w-full rounded-lg bg-white px-4 py-2 text-sm font-medium text-black transition-colors hover:bg-neutral-200 disabled:opacity-50"
		>
			{submitting ? 'Creating account…' : 'Create account'}
		</button>

		<p class="text-center text-sm text-neutral-500">
			Already have an account?
			<a href="/login" class="text-white underline underline-offset-4 hover:text-neutral-300">
				Sign in
			</a>
		</p>
	{/if}
</form>

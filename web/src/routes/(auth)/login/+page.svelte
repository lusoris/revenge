<script lang="ts">
	import { goto } from '$app/navigation';
	import { login, getAuth } from '$lib/stores/auth.svelte';

	const auth = getAuth();

	let email = $state('');
	let password = $state('');
	let totpCode = $state('');
	let showTotp = $state(false);
	let submitting = $state(false);
	let errorMsg = $state('');

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		submitting = true;
		errorMsg = '';

		try {
			await login(email, password, showTotp ? totpCode : undefined);
			goto('/');
		} catch (err: unknown) {
			const msg = err instanceof Error ? err.message : 'Login failed';

			// If MFA required, show TOTP field
			if (msg.toLowerCase().includes('mfa') || msg.toLowerCase().includes('totp')) {
				showTotp = true;
				errorMsg = 'Enter your two-factor authentication code.';
			} else {
				errorMsg = msg;
			}
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
		<label for="password" class="mb-1 block text-sm font-medium text-neutral-300">Password</label>
		<input
			id="password"
			type="password"
			bind:value={password}
			required
			autocomplete="current-password"
			class="w-full rounded-lg border border-neutral-800 bg-neutral-900 px-3 py-2 text-sm text-white placeholder-neutral-600 outline-none focus:border-neutral-600 focus:ring-1 focus:ring-neutral-600"
			placeholder="••••••••"
		/>
	</div>

	{#if showTotp}
		<div>
			<label for="totp" class="mb-1 block text-sm font-medium text-neutral-300">2FA Code</label>
			<input
				id="totp"
				type="text"
				bind:value={totpCode}
				required
				autocomplete="one-time-code"
				inputmode="numeric"
				maxlength="6"
				class="w-full rounded-lg border border-neutral-800 bg-neutral-900 px-3 py-2 text-center text-sm tracking-widest text-white placeholder-neutral-600 outline-none focus:border-neutral-600 focus:ring-1 focus:ring-neutral-600"
				placeholder="000000"
			/>
		</div>
	{/if}

	<button
		type="submit"
		disabled={submitting}
		class="w-full rounded-lg bg-white px-4 py-2 text-sm font-medium text-black transition-colors hover:bg-neutral-200 disabled:opacity-50"
	>
		{submitting ? 'Signing in…' : 'Sign in'}
	</button>

	<p class="text-center text-sm text-neutral-500">
		Don't have an account?
		<a href="/register" class="text-white underline underline-offset-4 hover:text-neutral-300">
			Register
		</a>
	</p>
</form>

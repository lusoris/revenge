<script lang="ts">
	import { goto } from '$app/navigation';
	import { getAuth } from '$lib/stores/auth.svelte';
	import AppShell from '$components/layout/AppShell.svelte';

	const auth = getAuth();

	let { children } = $props();

	// Redirect to login if not authenticated (and not still loading)
	$effect(() => {
		if (!auth.loading && !auth.isAuthenticated) {
			goto('/login');
		}
	});
</script>

{#if auth.loading}
	<div class="flex min-h-screen items-center justify-center bg-neutral-950">
		<div class="text-center">
			<div
				class="mx-auto h-8 w-8 animate-spin rounded-full border-2 border-neutral-700 border-t-white"
			></div>
			<p class="mt-3 text-sm text-neutral-500">Loading…</p>
		</div>
	</div>
{:else if auth.isAuthenticated}
	<AppShell>
		{@render children()}
	</AppShell>
{/if}

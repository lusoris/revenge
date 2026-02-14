<script lang="ts">
	import { goto } from '$app/navigation';
	import { getAuth } from '$lib/stores/auth.svelte';
	import type { Snippet } from 'svelte';

	const auth = getAuth();
	let { children }: { children: Snippet } = $props();

	$effect(() => {
		if (!auth.loading && !auth.isAuthenticated) {
			goto('/login');
		}
	});
</script>

{#if auth.loading}
	<div class="flex min-h-screen items-center justify-center bg-black">
		<div class="h-8 w-8 animate-spin rounded-full border-2 border-neutral-700 border-t-white"></div>
	</div>
{:else if auth.isAuthenticated}
	{@render children()}
{/if}

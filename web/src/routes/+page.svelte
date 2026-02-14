<script lang="ts">
	import { goto } from '$app/navigation';
	import { getAuth } from '$lib/stores/auth.svelte';
	import { onMount } from 'svelte';

	const auth = getAuth();

	onMount(() => {
		// Wait for auth init to complete, then redirect
		const check = () => {
			if (auth.loading) {
				setTimeout(check, 50);
				return;
			}
			goto(auth.isAuthenticated ? '/home' : '/login', { replaceState: true });
		};
		check();
	});
</script>

<div class="flex min-h-screen items-center justify-center bg-neutral-950">
	<div
		class="h-8 w-8 animate-spin rounded-full border-2 border-neutral-700 border-t-white"
	></div>
</div>

<script lang="ts">
	import '../app.css';
	import { QueryClient, QueryClientProvider } from '@tanstack/svelte-query';
	import { onMount } from 'svelte';
	import { initAuth } from '$lib/stores/auth.svelte';

	let { children } = $props();

	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				staleTime: 1000 * 60 * 2,
				retry: 1,
				refetchOnWindowFocus: false
			}
		}
	});

	onMount(() => {
		initAuth();
	});
</script>

<svelte:head>
	<title>Revenge</title>
</svelte:head>

<QueryClientProvider client={queryClient}>
	{@render children()}
</QueryClientProvider>

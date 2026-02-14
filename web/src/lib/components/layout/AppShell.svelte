<script lang="ts">
	import Sidebar from '$components/layout/Sidebar.svelte';

	let { children } = $props();
</script>

<div class="min-h-screen bg-neutral-950">
	<Sidebar />

	<!-- Main content area — offset for desktop sidebar, bottom padding for mobile nav -->
	<main class="pb-16 lg:pl-56 lg:pb-0">
		<div class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
			{@render children()}
		</div>
	</main>
</div>

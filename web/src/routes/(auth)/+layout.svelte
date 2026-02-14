<script lang="ts">
	let { children } = $props();
</script>

<div class="flex min-h-screen items-center justify-center bg-neutral-950 p-4">
	<div class="w-full max-w-sm">
		<div class="mb-8 text-center">
			<h1 class="text-3xl font-bold tracking-tight text-white">Revenge</h1>
			<p class="mt-1 text-sm text-neutral-500">Self-hosted media server</p>
		</div>
		{@render children()}
	</div>
</div>

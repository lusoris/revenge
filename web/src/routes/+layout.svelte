<script>
	import '../app.css';

	let { children } = $props();
</script>

<svelte:head>
	<title>Revenge</title>
</svelte:head>

{@render children()}

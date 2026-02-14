<script lang="ts">
	import { page } from '$app/state';

	let { children } = $props();

	const tabs = [
		{ href: '/settings', label: 'Profile' },
		{ href: '/settings/security', label: 'Security' },
		{ href: '/settings/sessions', label: 'Sessions' },
		{ href: '/settings/apikeys', label: 'API Keys' },
		{ href: '/settings/preferences', label: 'Preferences' }
	];

	function isActive(href: string): boolean {
		return page.url.pathname === href;
	}
</script>

<div class="space-y-6">
	<div>
		<h1 class="text-2xl font-bold text-white">Settings</h1>
		<p class="mt-1 text-sm text-neutral-400">Manage your account and preferences</p>
	</div>

	<nav class="flex gap-1 border-b border-neutral-800">
		{#each tabs as tab}
			<a
				href={tab.href}
				class="border-b-2 px-4 py-2 text-sm font-medium transition-colors {isActive(tab.href)
					? 'border-white text-white'
					: 'border-transparent text-neutral-400 hover:border-neutral-600 hover:text-neutral-200'}"
			>
				{tab.label}
			</a>
		{/each}
	</nav>

	<div>
		{@render children()}
	</div>
</div>

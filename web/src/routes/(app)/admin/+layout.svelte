<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { getAuth } from '$lib/stores/auth.svelte';

	const auth = getAuth();

	let { children } = $props();

	const tabs = [
		{ href: '/admin', label: 'Dashboard' },
		{ href: '/admin/users', label: 'Users' },
		{ href: '/admin/integrations', label: 'Integrations' },
		{ href: '/admin/roles', label: 'Roles' },
		{ href: '/admin/activity', label: 'Activity' },
		{ href: '/admin/settings', label: 'Settings' }
	];

	function isActive(href: string): boolean {
		if (href === '/admin') return page.url.pathname === '/admin';
		return page.url.pathname.startsWith(href);
	}

	// Redirect non-admins
	$effect(() => {
		if (!auth.loading && !auth.isAdmin) {
			goto('/');
		}
	});
</script>

{#if auth.isAdmin}
	<div class="space-y-6">
		<div>
			<h1 class="text-2xl font-bold text-white">Administration</h1>
			<p class="mt-1 text-sm text-neutral-400">Manage server, users and integrations</p>
		</div>

		<nav class="flex gap-1 overflow-x-auto border-b border-neutral-800">
			{#each tabs as tab}
				<a
					href={tab.href}
					class="whitespace-nowrap border-b-2 px-4 py-2 text-sm font-medium transition-colors {isActive(tab.href)
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
{/if}

<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import { getAuth, logout } from '$lib/stores/auth.svelte';
	import { imageUrl } from '$api/client';

	const auth = getAuth();

	interface NavItem {
		href: string;
		label: string;
		icon: string;
	}

	const nav: NavItem[] = [
		{ href: '/', label: 'Home', icon: '🏠' },
		{ href: '/movies', label: 'Movies', icon: '🎬' },
		{ href: '/tvshows', label: 'TV Shows', icon: '📺' },
		{ href: '/search', label: 'Search', icon: '🔍' }
	];

	const adminNav: NavItem[] = [{ href: '/admin', label: 'Admin', icon: '⚙️' }];

	function isActive(href: string): boolean {
		if (href === '/') return page.url.pathname === '/';
		return page.url.pathname.startsWith(href);
	}

	async function handleLogout() {
		await logout();
		goto('/login');
	}
</script>

<!-- Desktop sidebar -->
<aside
	class="fixed inset-y-0 left-0 z-40 hidden w-56 flex-col border-r border-neutral-800 bg-neutral-950 lg:flex"
>
	<!-- Logo -->
	<div class="flex h-14 items-center px-4">
		<a href="/" class="text-lg font-bold tracking-tight text-white">Revenge</a>
	</div>

	<!-- Navigation -->
	<nav class="flex-1 space-y-1 overflow-y-auto px-2 py-2">
		{#each nav as item}
			<a
				href={item.href}
				class="flex items-center gap-3 rounded-lg px-3 py-2 text-sm transition-colors {isActive(item.href)
					? 'bg-neutral-800 text-white'
					: 'text-neutral-400 hover:bg-neutral-900 hover:text-white'}"
			>
				<span class="text-base">{item.icon}</span>
				{item.label}
			</a>
		{/each}

		{#if auth.isAdmin}
			<div class="my-3 border-t border-neutral-800"></div>
			{#each adminNav as item}
				<a
					href={item.href}
					class="flex items-center gap-3 rounded-lg px-3 py-2 text-sm transition-colors {isActive(item.href)
						? 'bg-neutral-800 text-white'
						: 'text-neutral-400 hover:bg-neutral-900 hover:text-white'}"
				>
					<span class="text-base">{item.icon}</span>
					{item.label}
				</a>
			{/each}
		{/if}
	</nav>

	<!-- User section -->
	<div class="border-t border-neutral-800 p-3">
		{#if auth.user}
			<div class="flex items-center gap-3">
				{#if auth.user.avatar_url}
					<img
						src={auth.user.avatar_url}
						alt={auth.user.display_name ?? auth.user.username}
						class="h-8 w-8 rounded-full object-cover"
					/>
				{:else}
					<div
						class="flex h-8 w-8 items-center justify-center rounded-full bg-neutral-800 text-xs font-medium text-neutral-300"
					>
						{(auth.user.display_name ?? auth.user.username).charAt(0).toUpperCase()}
					</div>
				{/if}
				<div class="min-w-0 flex-1">
					<p class="truncate text-sm font-medium text-white">
						{auth.user.display_name ?? auth.user.username}
					</p>
					<button
						onclick={handleLogout}
						class="text-xs text-neutral-500 transition-colors hover:text-red-400"
					>
						Sign out
					</button>
				</div>
			</div>
		{/if}
	</div>
</aside>

<!-- Mobile bottom nav -->
<nav
	class="fixed inset-x-0 bottom-0 z-40 flex border-t border-neutral-800 bg-neutral-950/95 backdrop-blur-sm lg:hidden"
>
	{#each nav as item}
		<a
			href={item.href}
			class="flex flex-1 flex-col items-center gap-0.5 py-2 text-xs transition-colors {isActive(item.href)
				? 'text-white'
				: 'text-neutral-500'}"
		>
			<span class="text-lg">{item.icon}</span>
			{item.label}
		</a>
	{/each}
	<button
		onclick={handleLogout}
		class="flex flex-1 flex-col items-center gap-0.5 py-2 text-xs text-neutral-500 transition-colors hover:text-red-400"
	>
		<span class="text-lg">👋</span>
		Logout
	</button>
</nav>

<script lang="ts">
	import { listSessions, revokeAllOtherSessions, revokeSession } from '$api/endpoints/sessions';
	import { Badge } from '$components/ui/badge';
	import { Button } from '$components/ui/button';
	import * as Card from '$components/ui/card';
	import { createQuery, useQueryClient } from '@tanstack/svelte-query';

	const queryClient = useQueryClient();

	const sessionsQuery = createQuery(() => ({
		queryKey: ['sessions'],
		queryFn: () => listSessions()
	}));

	async function handleRevoke(id: string) {
		if (!confirm('Revoke this session?')) return;
		try {
			await revokeSession(id);
			queryClient.invalidateQueries({ queryKey: ['sessions'] });
		} catch (err: any) {
			alert(err.message);
		}
	}

	async function handleRevokeAll() {
		if (!confirm('Revoke all other sessions? You will remain logged in.')) return;
		try {
			await revokeAllOtherSessions();
			queryClient.invalidateQueries({ queryKey: ['sessions'] });
		} catch (err: any) {
			alert(err.message);
		}
	}

	function parseUA(ua?: string): string {
		if (!ua) return 'Unknown device';
		// Simple UA parsing
		if (ua.includes('Firefox')) return 'Firefox';
		if (ua.includes('Chrome')) return 'Chrome';
		if (ua.includes('Safari')) return 'Safari';
		if (ua.includes('Edge')) return 'Edge';
		return ua.substring(0, 50);
	}
</script>

<Card.Root class="border-neutral-800 bg-neutral-900">
	<Card.Header class="flex flex-row items-center justify-between">
		<div>
			<Card.Title class="text-white">Active Sessions</Card.Title>
			<Card.Description>Devices where you are currently logged in</Card.Description>
		</div>
		<Button variant="destructive" size="sm" onclick={handleRevokeAll}>
			Revoke All Others
		</Button>
	</Card.Header>
	<Card.Content>
		{#if sessionsQuery.isLoading}
			<p class="text-sm text-neutral-500">Loading sessions…</p>
		{:else if sessionsQuery.data}
			<div class="divide-y divide-neutral-800">
				{#each sessionsQuery.data.sessions as session}
					<div class="flex items-center justify-between py-3">
						<div class="min-w-0 flex-1">
							<div class="flex items-center gap-2">
								<p class="text-sm font-medium text-white">
									{session.device_name ?? parseUA(session.user_agent)}
								</p>
								{#if session.is_current}
									<Badge class="bg-green-900 text-green-200">Current</Badge>
								{/if}
							</div>
							<p class="text-xs text-neutral-500">
								{session.ip_address ?? 'Unknown IP'}
								· Created {new Date(session.created_at).toLocaleDateString()}
								{#if session.last_active_at}
									· Last active {new Date(session.last_active_at).toLocaleString()}
								{/if}
							</p>
						</div>
						{#if !session.is_current}
							<Button
								variant="destructive"
								size="sm"
								onclick={() => handleRevoke(session.id)}
							>
								Revoke
							</Button>
						{/if}
					</div>
				{/each}
			</div>
			{#if sessionsQuery.data.sessions.length === 0}
				<p class="py-8 text-center text-sm text-neutral-500">No active sessions found</p>
			{/if}
		{/if}
	</Card.Content>
</Card.Root>

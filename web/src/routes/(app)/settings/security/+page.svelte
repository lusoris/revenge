<script lang="ts">
	import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
	import { derived, writable } from 'svelte/store';
	import {
		getMFAStatus,
		enableMFA,
		disableMFA,
		setupTOTP,
		verifyTOTP,
		removeTOTP,
		generateBackupCodes,
		regenerateBackupCodes,
		listWebAuthnCredentials,
		deleteWebAuthnCredential
	} from '$api/endpoints/mfa';
	import { Button } from '$components/ui/button';
	import { Input } from '$components/ui/input';
	import { Label } from '$components/ui/label';
	import * as Card from '$components/ui/card';
	import * as Alert from '$components/ui/alert';
	import { Badge } from '$components/ui/badge';
	import type { MFAStatus, TOTPSetup, BackupCodesResponse, WebAuthnCredentialInfo } from '$api/types';

	const queryClient = useQueryClient();

	const statusQuery = createQuery(
		derived(writable(null), () => ({
			queryKey: ['mfa', 'status'],
			queryFn: () => getMFAStatus()
		}))
	);

	const credentialsQuery = createQuery(
		derived(writable(null), () => ({
			queryKey: ['mfa', 'webauthn'],
			queryFn: () => listWebAuthnCredentials()
		}))
	);

	// TOTP setup state
	let showTOTPSetup = $state(false);
	let totpSetupData: TOTPSetup | null = $state(null);
	let totpCode = $state('');
	let totpError = $state('');
	let totpSuccess = $state('');

	// Backup codes state
	let backupCodes: string[] | null = $state(null);
	let showBackupCodes = $state(false);

	async function handleSetupTOTP() {
		try {
			totpSetupData = await setupTOTP();
			showTOTPSetup = true;
			totpError = '';
		} catch (err: any) {
			totpError = err.message;
		}
	}

	async function handleVerifyTOTP() {
		try {
			const res = await verifyTOTP(totpCode);
			if (res.success) {
				totpSuccess = 'TOTP enabled successfully';
				showTOTPSetup = false;
				totpCode = '';
				queryClient.invalidateQueries({ queryKey: ['mfa'] });
			}
		} catch (err: any) {
			totpError = err.message;
		}
	}

	async function handleRemoveTOTP() {
		if (!confirm('Remove TOTP? You will need to set it up again.')) return;
		try {
			await removeTOTP();
			queryClient.invalidateQueries({ queryKey: ['mfa'] });
		} catch (err: any) {
			totpError = err.message;
		}
	}

	async function handleGenerateBackupCodes() {
		try {
			const res = await generateBackupCodes();
			backupCodes = res.codes;
			showBackupCodes = true;
		} catch (err: any) {
			totpError = err.message;
		}
	}

	async function handleDeleteCredential(id: string) {
		if (!confirm('Delete this security key?')) return;
		try {
			await deleteWebAuthnCredential(id);
			queryClient.invalidateQueries({ queryKey: ['mfa', 'webauthn'] });
		} catch (err: any) {
			totpError = err.message;
		}
	}
</script>

<div class="grid gap-6 lg:grid-cols-2">
	<!-- MFA Status -->
	<Card.Root class="border-neutral-800 bg-neutral-900">
		<Card.Header>
			<Card.Title class="text-white">Two-Factor Authentication</Card.Title>
			<Card.Description>Add an extra layer of security to your account</Card.Description>
		</Card.Header>
		<Card.Content class="space-y-4">
			{#if totpError}
				<Alert.Root class="border-red-800 bg-red-950 text-red-200">
					<Alert.Title>{totpError}</Alert.Title>
				</Alert.Root>
			{/if}
			{#if totpSuccess}
				<Alert.Root class="border-green-800 bg-green-950 text-green-200">
					<Alert.Title>{totpSuccess}</Alert.Title>
				</Alert.Root>
			{/if}

			{#if $statusQuery.data}
				{@const s = $statusQuery.data}
				<div class="space-y-3">
					<div class="flex items-center justify-between">
						<span class="text-sm text-neutral-300">TOTP (Authenticator App)</span>
						{#if s.has_totp}
							<Badge class="bg-green-900 text-green-200">Enabled</Badge>
						{:else}
							<Badge variant="outline" class="text-neutral-400">Disabled</Badge>
						{/if}
					</div>
					<div class="flex items-center justify-between">
						<span class="text-sm text-neutral-300">Security Keys</span>
						<Badge variant="outline" class="text-neutral-400">{s.webauthn_count}</Badge>
					</div>
					<div class="flex items-center justify-between">
						<span class="text-sm text-neutral-300">Backup Codes</span>
						<Badge variant="outline" class="text-neutral-400">{s.unused_backup_codes} remaining</Badge>
					</div>
				</div>
			{:else if $statusQuery.isLoading}
				<p class="text-sm text-neutral-500">Loading…</p>
			{/if}
		</Card.Content>
		<Card.Footer class="flex gap-2">
			{#if $statusQuery.data?.has_totp}
				<Button variant="destructive" size="sm" onclick={handleRemoveTOTP}>Remove TOTP</Button>
			{:else}
				<Button size="sm" onclick={handleSetupTOTP}>Setup TOTP</Button>
			{/if}
			<Button variant="outline" size="sm" onclick={handleGenerateBackupCodes}>
				Generate Backup Codes
			</Button>
		</Card.Footer>
	</Card.Root>

	<!-- TOTP Setup -->
	{#if showTOTPSetup && totpSetupData}
		<Card.Root class="border-neutral-800 bg-neutral-900">
			<Card.Header>
				<Card.Title class="text-white">Setup Authenticator</Card.Title>
				<Card.Description>Scan this QR code with your authenticator app</Card.Description>
			</Card.Header>
			<Card.Content class="space-y-4">
				<div class="flex justify-center rounded-lg bg-white p-4">
					<img src={totpSetupData.qr_code} alt="TOTP QR Code" class="h-48 w-48" />
				</div>
				<div class="space-y-1">
					<p class="text-xs text-neutral-500">Or enter this secret manually:</p>
					<code class="block break-all rounded bg-neutral-800 px-3 py-2 text-xs text-neutral-300">
						{totpSetupData.secret}
					</code>
				</div>
				<div class="space-y-2">
					<Label for="totpCode">Verification Code</Label>
					<Input
						id="totpCode"
						bind:value={totpCode}
						placeholder="000000"
						maxlength={6}
						class="bg-neutral-800"
					/>
				</div>
			</Card.Content>
			<Card.Footer class="flex gap-2">
				<Button onclick={handleVerifyTOTP} disabled={totpCode.length !== 6}>Verify</Button>
				<Button variant="outline" onclick={() => (showTOTPSetup = false)}>Cancel</Button>
			</Card.Footer>
		</Card.Root>
	{/if}

	<!-- Backup Codes -->
	{#if showBackupCodes && backupCodes}
		<Card.Root class="border-neutral-800 bg-neutral-900">
			<Card.Header>
				<Card.Title class="text-white">Backup Codes</Card.Title>
				<Card.Description>Save these somewhere safe. Each code can only be used once.</Card.Description>
			</Card.Header>
			<Card.Content>
				<div class="grid grid-cols-2 gap-2">
					{#each backupCodes as code}
						<code class="rounded bg-neutral-800 px-3 py-2 text-center text-sm text-neutral-200">
							{code}
						</code>
					{/each}
				</div>
			</Card.Content>
			<Card.Footer>
				<Button variant="outline" size="sm" onclick={() => (showBackupCodes = false)}>
					Done
				</Button>
			</Card.Footer>
		</Card.Root>
	{/if}

	<!-- WebAuthn Credentials -->
	{#if $credentialsQuery.data && $credentialsQuery.data.credentials.length > 0}
		<Card.Root class="border-neutral-800 bg-neutral-900 lg:col-span-2">
			<Card.Header>
				<Card.Title class="text-white">Security Keys</Card.Title>
			</Card.Header>
			<Card.Content>
				<div class="divide-y divide-neutral-800">
					{#each $credentialsQuery.data.credentials as cred}
						<div class="flex items-center justify-between py-3">
							<div>
								<p class="text-sm font-medium text-white">{cred.name}</p>
								<p class="text-xs text-neutral-500">
									Added {new Date(cred.created_at).toLocaleDateString()}
									{#if cred.last_used_at}
										· Last used {new Date(cred.last_used_at).toLocaleDateString()}
									{/if}
								</p>
							</div>
							<Button
								variant="destructive"
								size="sm"
								onclick={() => handleDeleteCredential(cred.id)}
							>
								Remove
							</Button>
						</div>
					{/each}
				</div>
			</Card.Content>
		</Card.Root>
	{/if}
</div>

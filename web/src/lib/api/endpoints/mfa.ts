import { get, post, del, patch } from '../client';
import type {
	MFAStatus,
	TOTPSetup,
	TOTPVerifyResponse,
	BackupCodesResponse,
	WebAuthnCredentialInfo,
	WebAuthnCredentialsList
} from '../types';

// ─── MFA Status ──────────────────────────────────────────────────────────────

export async function getMFAStatus(): Promise<MFAStatus> {
	return get<MFAStatus>('/v1/mfa/status');
}

export async function enableMFA(): Promise<void> {
	return post('/v1/mfa/enable');
}

export async function disableMFA(): Promise<void> {
	return post('/v1/mfa/disable');
}

// ─── TOTP ────────────────────────────────────────────────────────────────────

export async function setupTOTP(accountName?: string): Promise<TOTPSetup> {
	return post<TOTPSetup>('/v1/mfa/totp/setup', { account_name: accountName });
}

export async function verifyTOTP(code: string): Promise<TOTPVerifyResponse> {
	return post<TOTPVerifyResponse>('/v1/mfa/totp/verify', { code });
}

export async function removeTOTP(): Promise<void> {
	return del('/v1/mfa/totp');
}

// ─── Backup Codes ────────────────────────────────────────────────────────────

export async function generateBackupCodes(): Promise<BackupCodesResponse> {
	return post<BackupCodesResponse>('/v1/mfa/backup-codes/generate');
}

export async function regenerateBackupCodes(): Promise<BackupCodesResponse> {
	return post<BackupCodesResponse>('/v1/mfa/backup-codes/regenerate');
}

// ─── WebAuthn ────────────────────────────────────────────────────────────────

export async function listWebAuthnCredentials(): Promise<WebAuthnCredentialsList> {
	return get<WebAuthnCredentialsList>('/v1/mfa/webauthn/credentials');
}

export async function beginWebAuthnRegistration(
	name?: string
): Promise<{ options: Record<string, unknown> }> {
	return post('/v1/mfa/webauthn/register/begin', { credential_name: name });
}

export async function finishWebAuthnRegistration(
	credential: Record<string, unknown>,
	name?: string
): Promise<{ success: boolean; message: string }> {
	return post('/v1/mfa/webauthn/register/finish', { credential, credential_name: name });
}

export async function renameWebAuthnCredential(id: string, name: string): Promise<void> {
	return patch(`/v1/mfa/webauthn/credentials/${id}`, { name });
}

export async function deleteWebAuthnCredential(id: string): Promise<void> {
	return del(`/v1/mfa/webauthn/credentials/${id}`);
}

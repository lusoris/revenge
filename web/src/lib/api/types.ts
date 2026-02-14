// API types — derived from OpenAPI schemas
// Only includes types needed for the current implementation phase.

// ─── Auth ────────────────────────────────────────────────────────────────────

export interface LoginRequest {
	username: string;
	password: string;
	totp_code?: string;
}

export interface LoginResponse {
	user: User;
	access_token: string;
	refresh_token: string;
	expires_in: number;
}

export interface RegisterRequest {
	username: string;
	email: string;
	password: string;
}

export interface RefreshSessionRequest {
	refresh_token: string;
}

export interface RefreshSessionResponse {
	access_token: string;
	refresh_token: string;
	expires_in: number;
}

export interface LogoutRequest {
	refresh_token?: string;
}

// ─── User ────────────────────────────────────────────────────────────────────

export interface User {
	id: string;
	username: string;
	email: string;
	display_name?: string;
	avatar_url?: string;
	is_admin: boolean;
	is_active: boolean;
	email_verified: boolean;
	locale?: string;
	timezone?: string;
	created_at: string;
	updated_at: string;
	last_login_at?: string;
}

export interface UserUpdate {
	email?: string;
	display_name?: string;
	timezone?: string;
}

export interface UserPreferences {
	user_id?: string;
	email_notifications?: { enabled?: boolean; frequency?: 'instant' | 'daily' | 'weekly' };
	push_notifications?: { enabled?: boolean };
	digest_notifications?: { enabled?: boolean; frequency?: 'daily' | 'weekly' | 'monthly' };
	profile_visibility?: 'public' | 'friends' | 'private';
	show_email?: boolean;
	show_activity?: boolean;
	theme?: 'light' | 'dark' | 'system';
	display_language?: string;
	content_language?: string;
	metadata_language?: string;
	show_adult_content?: boolean;
	show_spoilers?: boolean;
	auto_play_videos?: boolean;
}

export interface Avatar {
	id: string;
	user_id: string;
	file_path: string;
	file_size_bytes?: number;
	mime_type?: string;
	width?: number;
	height?: number;
	is_animated?: boolean;
	version: number;
	is_current: boolean;
	uploaded_at: string;
}

// ─── Health ──────────────────────────────────────────────────────────────────

export interface HealthCheck {
	name: string;
	status: 'healthy' | 'unhealthy';
	message?: string;
	details?: Record<string, HealthCheck>;
}

// ─── Movies ──────────────────────────────────────────────────────────────────

export interface Movie {
	id: string;
	title: string;
	original_title?: string;
	overview?: string;
	tagline?: string;
	release_date?: string;
	runtime?: number;
	vote_average?: number;
	vote_count?: number;
	poster_path?: string;
	backdrop_path?: string;
	tmdb_id?: number;
	imdb_id?: string;
	status?: string;
	genres?: Genre[];
	library_id: string;
	created_at: string;
	updated_at: string;
}

export interface MovieListResponse {
	items: Movie[];
	total: number;
	page?: number;
	page_size?: number;
}

export interface MovieFile {
	id: string;
	movie_id: string;
	file_path: string;
	file_size: number;
	duration?: number;
	resolution?: string;
	codec?: string;
	bitrate?: number;
	audio_channels?: number;
	audio_codec?: string;
	subtitle_tracks?: SubtitleTrack[];
	audio_tracks?: AudioTrack[];
}

export interface ContinueWatchingItem {
	movie: Movie;
	progress: MovieWatched;
}

export interface MovieWatched {
	movie_id: string;
	user_id: string;
	position_seconds: number;
	duration_seconds: number;
	percentage: number;
	completed: boolean;
	last_watched_at: string;
}

// ─── TV Shows ────────────────────────────────────────────────────────────────

export interface TVSeries {
	id: string;
	name: string;
	original_name?: string;
	overview?: string;
	first_air_date?: string;
	last_air_date?: string;
	vote_average?: number;
	vote_count?: number;
	poster_path?: string;
	backdrop_path?: string;
	tmdb_id?: number;
	imdb_id?: string;
	status?: string;
	number_of_seasons?: number;
	number_of_episodes?: number;
	genres?: Genre[];
	networks?: TVNetwork[];
	library_id: string;
	created_at: string;
	updated_at: string;
}

export interface TVShowListResponse {
	items: TVSeries[];
	total: number;
	page?: number;
	page_size?: number;
}

export interface TVSeason {
	id: string;
	series_id: string;
	season_number: number;
	name?: string;
	overview?: string;
	poster_path?: string;
	air_date?: string;
	episode_count?: number;
}

export interface TVEpisode {
	id: string;
	series_id: string;
	season_id: string;
	season_number: number;
	episode_number: number;
	name?: string;
	overview?: string;
	air_date?: string;
	runtime?: number;
	still_path?: string;
	vote_average?: number;
}

export interface TVEpisodeFile {
	id: string;
	episode_id: string;
	file_path: string;
	file_size: number;
	duration?: number;
	resolution?: string;
	codec?: string;
}

export interface TVContinueWatchingItem {
	series: TVSeries;
	episode: TVEpisode;
	progress: EpisodeWatchProgress;
}

export interface EpisodeWatchProgress {
	episode_id: string;
	user_id: string;
	position_seconds: number;
	duration_seconds: number;
	percentage: number;
	completed: boolean;
	last_watched_at: string;
}

// ─── Common ──────────────────────────────────────────────────────────────────

export interface Genre {
	id: number;
	name: string;
	movie_count?: number;
	tvshow_count?: number;
}

export interface TVNetwork {
	id: number;
	name: string;
	logo_path?: string;
}

export interface SubtitleTrack {
	index: number;
	language: string;
	title?: string;
	codec?: string;
	forced?: boolean;
	default?: boolean;
}

export interface AudioTrack {
	index: number;
	language: string;
	title?: string;
	codec?: string;
	channels?: number;
	default?: boolean;
}

export interface Credit {
	id: string;
	name: string;
	character?: string;
	department?: string;
	job?: string;
	profile_path?: string;
	order?: number;
}

export interface CreditListResponse {
	credits: Credit[];
	total: number;
}

// ─── Playback ────────────────────────────────────────────────────────────────

export interface StartPlaybackRequest {
	movie_file_id?: string;
	episode_file_id?: string;
	audio_track_index?: number;
	subtitle_track_index?: number;
	start_position_seconds?: number;
}

export interface PlaybackSession {
	id: string;
	master_playlist_url: string;
	profiles: PlaybackProfile[];
	audio_tracks: AudioTrack[];
	subtitle_tracks: SubtitleTrack[];
	duration_seconds: number;
	start_position_seconds: number;
	created_at: string;
}

export interface PlaybackProfile {
	name: string;
	resolution: string;
	bitrate: number;
	codec: string;
}

// ─── Search ──────────────────────────────────────────────────────────────────

export interface SearchHit {
	document: Movie;
	score?: number;
	highlights?: Record<string, string[]>;
}

export interface TVShowSearchHit {
	document: TVSeries;
	score?: number;
	highlights?: Record<string, string[]>;
}

export interface SearchResults {
	hits: SearchHit[];
	total_hits: number;
	total_pages: number;
	current_page: number;
}

export interface TVShowSearchResults {
	hits: TVShowSearchHit[];
	total_hits: number;
	total_pages: number;
	current_page: number;
}

export interface MultiSearchResults {
	movies?: SearchResults;
	tvshows?: TVShowSearchResults;
}

export interface AutocompleteResults {
	suggestions: AutocompleteSuggestion[];
}

export interface AutocompleteSuggestion {
	id: string;
	title: string;
	type: 'movie' | 'tvshow';
	year?: string;
	poster_path?: string;
}

// ─── Libraries ───────────────────────────────────────────────────────────────

export interface Library {
	id: string;
	name: string;
	type: 'movie' | 'tvshow';
	paths: string[];
	movie_count?: number;
	tvshow_count?: number;
	created_at: string;
	updated_at: string;
}

// ─── Sessions ────────────────────────────────────────────────────────────────

export interface SessionInfo {
	id: string;
	ip_address?: string;
	user_agent?: string;
	device_name?: string;
	created_at: string;
	last_active_at?: string;
	last_activity_at?: string;
	expires_at?: string;
	is_current: boolean;
	is_active?: boolean;
}

export interface SessionListResponse {
	sessions: SessionInfo[];
	total: number;
}

// ─── OIDC ────────────────────────────────────────────────────────────────────

export interface OIDCProvider {
	name: string;
	display_name?: string;
	is_default?: boolean;
}

export interface OIDCCallbackResponse {
	access_token: string;
	token_type: string;
	expires_in: number;
	refresh_token?: string;
	user?: User;
}

export interface OIDCUserLink {
	id: string;
	provider_name: string;
	provider_display_name: string;
	email?: string;
	last_login_at?: string;
	linked_at: string;
}

export interface OIDCUserLinkListResponse {
	links: OIDCUserLink[];
	total: number;
}

export interface OIDCAuthURLResponse {
	auth_url: string;
}

// ─── API Error ───────────────────────────────────────────────────────────────

export interface APIError {
	error: string;
	message: string;
	status_code: number;
	request_id?: string;
	details?: Record<string, string>;
}

// ─── Auth (password/email) ───────────────────────────────────────────────────

export interface ChangePasswordRequest {
	old_password: string;
	new_password: string;
}

export interface ForgotPasswordRequest {
	email: string;
}

export interface ResetPasswordRequest {
	token: string;
	new_password: string;
}

// ─── API Keys ────────────────────────────────────────────────────────────────

export type APIKeyScope = 'read' | 'write' | 'admin';

export interface APIKeyInfo {
	id: string;
	user_id: string;
	name: string;
	description?: string;
	key_prefix: string;
	scopes: APIKeyScope[];
	is_active: boolean;
	expires_at?: string;
	last_used_at?: string;
	created_at: string;
	updated_at: string;
}

export interface APIKeyListResponse {
	keys: APIKeyInfo[];
	total: number;
}

export interface CreateAPIKeyRequest {
	name: string;
	description?: string;
	scopes: APIKeyScope[];
	expires_at?: string;
}

export interface CreateAPIKeyResponse {
	id: string;
	name: string;
	key_prefix: string;
	scopes: string[];
	created_at: string;
	api_key: string;
	message: string;
}

// ─── MFA ─────────────────────────────────────────────────────────────────────

export interface MFAStatus {
	user_id: string;
	has_totp: boolean;
	webauthn_count: number;
	unused_backup_codes: number;
	require_mfa: boolean;
}

export interface TOTPSetup {
	secret: string;
	qr_code: string;
	url: string;
}

export interface TOTPVerifyResponse {
	success: boolean;
	message: string;
}

export interface BackupCodesResponse {
	codes: string[];
	count: number;
}

export interface WebAuthnCredentialInfo {
	id: string;
	name: string;
	created_at: string;
	last_used_at?: string;
	backup_eligible: boolean;
	backup_state: boolean;
	clone_detected: boolean;
}

export interface WebAuthnCredentialsList {
	credentials: WebAuthnCredentialInfo[];
	count: number;
}

// ─── RBAC ────────────────────────────────────────────────────────────────────

export interface Permission {
	resource: string;
	action: string;
}

export interface RoleDetail {
	name: string;
	description?: string;
	permissions: Permission[];
	is_built_in: boolean;
	user_count: number;
}

export interface RolesResponse {
	roles: RoleDetail[];
	total: number;
}

export interface CreateRoleRequest {
	name: string;
	description?: string;
	permissions?: Permission[];
}

export interface PolicyListResponse {
	policies: { subject: string; object: string; action: string }[];
	total: number;
}

// ─── Settings ────────────────────────────────────────────────────────────────

export type SettingDataType = 'string' | 'integer' | 'boolean' | 'float' | 'json';

export interface ServerSetting {
	key: string;
	value: string | number | boolean | Record<string, unknown>;
	description?: string;
	category?: string;
	data_type: SettingDataType;
	is_secret?: boolean;
	is_public?: boolean;
}

export interface UserSetting {
	user_id: string;
	key: string;
	value: string | number | boolean | Record<string, unknown>;
	description?: string;
	category?: string;
	data_type: SettingDataType;
}

// ─── Admin Users ─────────────────────────────────────────────────────────────

export interface AdminUserListResponse {
	users: User[];
	total: number;
}

// ─── Admin OIDC ──────────────────────────────────────────────────────────────

export interface ClaimMappings {
	username?: string;
	email?: string;
	name?: string;
	picture?: string;
	roles?: string;
}

export interface AdminOIDCProvider {
	id: string;
	name: string;
	display_name: string;
	provider_type: 'generic' | 'authentik' | 'keycloak';
	issuer_url: string;
	client_id: string;
	scopes: string[];
	claim_mappings?: ClaimMappings;
	role_mappings?: Record<string, string>;
	auto_create_users: boolean;
	update_user_info: boolean;
	allow_linking: boolean;
	is_enabled: boolean;
	is_default: boolean;
	created_at: string;
	updated_at: string;
}

export interface AdminOIDCProviderListResponse {
	providers: AdminOIDCProvider[];
	total: number;
}

export interface CreateOIDCProviderRequest {
	name: string;
	display_name: string;
	provider_type?: 'generic' | 'authentik' | 'keycloak';
	issuer_url: string;
	client_id: string;
	client_secret: string;
	scopes?: string[];
	claim_mappings?: ClaimMappings;
	auto_create_users?: boolean;
	is_enabled?: boolean;
	is_default?: boolean;
}

// ─── Admin Activity ──────────────────────────────────────────────────────────

export interface ActivityLogEntry {
	id: string;
	user_id?: string;
	username?: string;
	action: string;
	resource_type?: string;
	resource_id?: string;
	changes?: Record<string, unknown>;
	metadata?: Record<string, unknown>;
	ip_address?: string;
	user_agent?: string;
	success: boolean;
	error_message?: string;
	created_at: string;
}

export interface ActivityLogListResponse {
	entries: ActivityLogEntry[];
	total: number;
	page?: number;
	page_size?: number;
}

export interface ActivityStats {
	total_count: number;
	success_count: number;
	failed_count: number;
	oldest_entry?: string;
	newest_entry?: string;
}

// ─── Admin Integrations ─────────────────────────────────────────────────────

export interface RadarrStatus {
	connected: boolean;
	version?: string;
	instanceName?: string;
	startTime?: string;
	syncStatus: {
		isRunning: boolean;
		lastSync?: string;
		lastSyncError?: string;
		moviesAdded: number;
		moviesUpdated: number;
		moviesRemoved: number;
		totalMovies: number;
	};
}

export interface SonarrStatus {
	connected: boolean;
	version?: string;
	instanceName?: string;
	startTime?: string;
	syncStatus: {
		isRunning: boolean;
		lastSync?: string;
		lastSyncError?: string;
		seriesAdded: number;
		seriesUpdated: number;
		seriesRemoved: number;
		totalSeries: number;
		totalEpisodes?: number;
	};
}

export interface QualityProfile {
	id: number;
	name: string;
	upgradeAllowed?: boolean;
}

export interface RootFolder {
	id: number;
	path: string;
	accessible: boolean;
	freeSpace?: number;
}

// ─── Collections ─────────────────────────────────────────────────────────────

export interface MovieCollection {
	id: string;
	tmdb_collection_id?: number;
	name: string;
	overview?: string;
	poster_path?: string;
	backdrop_path?: string;
	created_at: string;
	updated_at: string;
}

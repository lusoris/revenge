// API types — derived from OpenAPI schemas
// Only includes types needed for the current implementation phase.

// ─── Auth ────────────────────────────────────────────────────────────────────

export interface LoginRequest {
	email: string;
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
	created_at: string;
	updated_at: string;
	last_login_at?: string;
}

export interface UserPreferences {
	theme?: 'dark' | 'light' | 'system';
	language?: string;
	subtitle_language?: string;
	audio_language?: string;
	default_quality?: string;
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
	movies: Movie[];
	total: number;
	limit: number;
	offset: number;
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
	tvshows: TVSeries[];
	total: number;
	limit: number;
	offset: number;
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

export interface SearchResults {
	hits: Movie[];
	total_hits: number;
	total_pages: number;
	current_page: number;
}

export interface TVShowSearchResults {
	hits: TVSeries[];
	total_hits: number;
	total_pages: number;
	current_page: number;
}

export interface MultiSearchResults {
	movies: Movie[];
	tvshows: TVSeries[];
	total_hits: number;
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
	user_id: string;
	ip_address?: string;
	user_agent?: string;
	created_at: string;
	last_active_at: string;
	is_current: boolean;
}

export interface SessionListResponse {
	sessions: SessionInfo[];
	total: number;
}

// ─── OIDC ────────────────────────────────────────────────────────────────────

export interface OIDCProvider {
	name: string;
	display_name?: string;
	icon_url?: string;
}

export interface OIDCCallbackResponse {
	access_token: string;
	token_type: string;
	expires_in: number;
	refresh_token?: string;
	user?: User;
}

// ─── API Error ───────────────────────────────────────────────────────────────

export interface APIError {
	error: string;
	message: string;
	status_code: number;
	request_id?: string;
	details?: Record<string, string>;
}

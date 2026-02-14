// tests/load/all_endpoints.js - Systematic test of all API endpoints
// Usage: k6 run --env PROFILE=gentle tests/load/all_endpoints.js
import { check, group, sleep } from 'k6';
import { Counter, Rate, Trend } from 'k6/metrics';
import { PROFILES, THRESHOLDS } from './config.js';
import { authGet, ensureUserPool, extractItems, getToken, login, randomFrom, vuUser } from './helpers.js';

// Custom metrics
const endpointErrors = new Counter('endpoint_errors');
const endpointLatency = new Trend('endpoint_latency');
const successRate = new Rate('success_rate');

// Select profile from env or default to smoke
const profileName = __ENV.PROFILE || 'smoke';
const profile = PROFILES[profileName] || PROFILES.smoke;
const threshold = THRESHOLDS[profileName] || THRESHOLDS.gentle;

export const options = {
    stages: profile.stages,
    thresholds: {
        ...threshold,
        'endpoint_errors': ['count<100'],
        'success_rate': ['rate>0.95'],
    },
};

// Global state for discovered IDs
let movieIds = [];
let tvshowIds = [];
let seasonIds = [];
let episodeIds = [];
let libraryIds = [];
let sessionIds = [];

export function setup() {
    // Create/ensure multi-user pool
    const userPool = ensureUserPool();

    // Login as first user to discover content IDs — token is NOT shared with VUs
    const token = login(userPool[0]);
    if (!token) {
        throw new Error('Failed to login during setup');
    }

    // Discover some IDs for parameterized endpoints
    const movies = authGet('/movies?limit=10', token);
    if (movies.status === 200) {
        const items = extractItems(movies);
        movieIds = items.map(m => m.id).filter(Boolean);
    }

    const tvshows = authGet('/tvshows?limit=10', token);
    if (tvshows.status === 200) {
        const items = extractItems(tvshows);
        tvshowIds = items.map(t => t.id).filter(Boolean);

        // Get seasons from first show
        if (tvshowIds.length > 0) {
            const seasons = authGet(`/tvshows/${tvshowIds[0]}/seasons`, token);
            if (seasons.status === 200) {
                const sItems = extractItems(seasons);
                seasonIds = sItems.map(s => s.id).filter(Boolean);
            }
        }
    }

    const libraries = authGet('/libraries', token);
    if (libraries.status === 200) {
        const items = extractItems(libraries);
        libraryIds = items.map(l => l.id).filter(Boolean);
    }

    const sessions = authGet('/sessions', token);
    if (sessions.status === 200) {
        const items = extractItems(sessions);
        sessionIds = items.map(s => s.id).filter(Boolean);
    }

    console.log(`Setup complete: ${movieIds.length} movies, ${tvshowIds.length} tvshows, ${libraryIds.length} libraries`);

    return { movieIds, tvshowIds, seasonIds, episodeIds, libraryIds, sessionIds, userPool };
}

export default function(data) {
    // Each VU gets its own user from the pool + own token
    const user = vuUser(data.userPool);
    const token = getToken(user);
    if (!token) {
        console.error('VU login failed, skipping iteration');
        return;
    }
    const movieId = randomFrom(data.movieIds) || '00000000-0000-0000-0000-000000000001';
    const tvshowId = randomFrom(data.tvshowIds) || '00000000-0000-0000-0000-000000000001';
    const seasonId = randomFrom(data.seasonIds) || '00000000-0000-0000-0000-000000000001';
    const libraryId = randomFrom(data.libraryIds) || '00000000-0000-0000-0000-000000000001';

    // Randomly select an endpoint group
    const groups = [
        () => testMovieEndpoints(token, movieId),
        () => testTvshowEndpoints(token, tvshowId, seasonId),
        () => testSearchEndpoints(token),
        () => testUserEndpoints(token),
        () => testMiscEndpoints(token, libraryId),
    ];

    const selectedGroup = groups[Math.floor(Math.random() * groups.length)];
    selectedGroup();

    sleep(Math.random() * 0.5 + 0.1); // 0.1-0.6s between requests
}

function recordMetric(res, name) {
    endpointLatency.add(res.timings.duration, { endpoint: name });
    // 2xx, 3xx, 404 (not found), 403 (forbidden) are all acceptable API responses
    const acceptable = (res.status >= 200 && res.status < 400) ||
                       res.status === 404 || res.status === 403;
    if (acceptable) {
        successRate.add(1);
    } else {
        successRate.add(0);
        endpointErrors.add(1, { endpoint: name, status: res.status });
    }
}

function testMovieEndpoints(token, movieId) {
    group('Movie Endpoints', () => {
        const endpoints = [
            { name: 'list_movies', fn: () => authGet('/movies', token) },
            { name: 'movies_recently_added', fn: () => authGet('/movies/recently-added', token) },
            { name: 'movies_top_rated', fn: () => authGet('/movies/top-rated', token) },
            { name: 'movies_continue_watching', fn: () => authGet('/movies/continue-watching', token) },
            { name: 'movies_stats', fn: () => authGet('/movies/stats', token) },
            { name: 'get_movie', fn: () => authGet(`/movies/${movieId}`, token) },
            { name: 'movie_files', fn: () => authGet(`/movies/${movieId}/files`, token) },
            { name: 'movie_cast', fn: () => authGet(`/movies/${movieId}/cast`, token) },
            { name: 'movie_crew', fn: () => authGet(`/movies/${movieId}/crew`, token) },
            { name: 'movie_genres', fn: () => authGet(`/movies/${movieId}/genres`, token) },
            { name: 'movie_similar', fn: () => authGet(`/movies/${movieId}/similar`, token) },
            { name: 'movie_progress', fn: () => authGet(`/movies/${movieId}/progress`, token) },
        ];

        const endpoint = endpoints[Math.floor(Math.random() * endpoints.length)];
        const res = endpoint.fn();
        recordMetric(res, endpoint.name);
        check(res, { [`${endpoint.name} ok`]: (r) => r.status >= 200 && r.status < 500 });
    });
}

function testTvshowEndpoints(token, tvshowId, seasonId) {
    group('TV Show Endpoints', () => {
        const endpoints = [
            { name: 'list_tvshows', fn: () => authGet('/tvshows', token) },
            { name: 'tvshows_recently_added', fn: () => authGet('/tvshows/recently-added', token) },
            { name: 'tvshows_continue_watching', fn: () => authGet('/tvshows/continue-watching', token) },
            { name: 'tvshows_stats', fn: () => authGet('/tvshows/stats', token) },
            { name: 'tvshows_episodes_recent', fn: () => authGet('/tvshows/episodes/recent', token) },
            { name: 'tvshows_episodes_upcoming', fn: () => authGet('/tvshows/episodes/upcoming', token) },
            { name: 'get_tvshow', fn: () => authGet(`/tvshows/${tvshowId}`, token) },
            { name: 'tvshow_seasons', fn: () => authGet(`/tvshows/${tvshowId}/seasons`, token) },
            { name: 'tvshow_episodes', fn: () => authGet(`/tvshows/${tvshowId}/episodes`, token) },
            { name: 'tvshow_cast', fn: () => authGet(`/tvshows/${tvshowId}/cast`, token) },
            { name: 'tvshow_crew', fn: () => authGet(`/tvshows/${tvshowId}/crew`, token) },
            { name: 'tvshow_genres', fn: () => authGet(`/tvshows/${tvshowId}/genres`, token) },
            { name: 'tvshow_watch_stats', fn: () => authGet(`/tvshows/${tvshowId}/watch-stats`, token) },
            { name: 'tvshow_next_episode', fn: () => authGet(`/tvshows/${tvshowId}/next-episode`, token) },
            { name: 'get_season', fn: () => authGet(`/tvshows/seasons/${seasonId}`, token) },
            { name: 'season_episodes', fn: () => authGet(`/tvshows/seasons/${seasonId}/episodes`, token) },
        ];

        const endpoint = endpoints[Math.floor(Math.random() * endpoints.length)];
        const res = endpoint.fn();
        recordMetric(res, endpoint.name);
        check(res, { [`${endpoint.name} ok`]: (r) => r.status >= 200 && r.status < 500 });
    });
}

function testSearchEndpoints(token) {
    group('Search Endpoints', () => {
        const queries = ['action', 'comedy', 'drama', 'test', 'movie', 'series', 'star'];
        const query = queries[Math.floor(Math.random() * queries.length)];

        const endpoints = [
            { name: 'search_movies', fn: () => authGet(`/search/movies?q=${query}`, token) },
            { name: 'search_movies_autocomplete', fn: () => authGet(`/search/movies/autocomplete?q=${query}`, token) },
            { name: 'search_movies_facets', fn: () => authGet('/search/movies/facets', token) },
            { name: 'search_tvshows', fn: () => authGet(`/search/tvshows?q=${query}`, token) },
            { name: 'search_tvshows_autocomplete', fn: () => authGet(`/search/tvshows/autocomplete?q=${query}`, token) },
            { name: 'search_tvshows_facets', fn: () => authGet('/search/tvshows/facets', token) },
            { name: 'search_multi', fn: () => authGet(`/search/multi?q=${query}`, token) },
        ];

        const endpoint = endpoints[Math.floor(Math.random() * endpoints.length)];
        const res = endpoint.fn();
        recordMetric(res, endpoint.name);
        check(res, { [`${endpoint.name} ok`]: (r) => r.status >= 200 && r.status < 500 });
    });
}

function testUserEndpoints(token) {
    group('User Endpoints', () => {
        const endpoints = [
            { name: 'get_me', fn: () => authGet('/users/me', token) },
            { name: 'get_preferences', fn: () => authGet('/users/me/preferences', token) },
            { name: 'get_oidc', fn: () => authGet('/users/me/oidc', token) },
            { name: 'list_sessions', fn: () => authGet('/sessions', token) },
            { name: 'current_session', fn: () => authGet('/sessions/current', token) },
            { name: 'mfa_status', fn: () => authGet('/mfa/status', token) },
            { name: 'list_apikeys', fn: () => authGet('/apikeys', token) },
            { name: 'user_settings', fn: () => authGet('/settings/user', token) },
        ];

        const endpoint = endpoints[Math.floor(Math.random() * endpoints.length)];
        const res = endpoint.fn();
        recordMetric(res, endpoint.name);
        check(res, { [`${endpoint.name} ok`]: (r) => r.status >= 200 && r.status < 500 });
    });
}

function testMiscEndpoints(token, libraryId) {
    group('Misc Endpoints', () => {
        const endpoints = [
            { name: 'list_genres', fn: () => authGet('/genres', token) },
            { name: 'list_libraries', fn: () => authGet('/libraries', token) },
            { name: 'get_library', fn: () => authGet(`/libraries/${libraryId}`, token) },
            { name: 'oidc_providers', fn: () => authGet('/oidc/providers', token) },
            { name: 'rbac_roles', fn: () => authGet('/rbac/roles', token) },
            { name: 'rbac_permissions', fn: () => authGet('/rbac/permissions', token) },
            { name: 'metadata_providers', fn: () => authGet('/metadata/providers', token) },
            { name: 'server_settings', fn: () => authGet('/settings/server', token) },
            { name: 'playback_sessions', fn: () => authGet('/playback/sessions', token) },
        ];

        const endpoint = endpoints[Math.floor(Math.random() * endpoints.length)];
        const res = endpoint.fn();
        recordMetric(res, endpoint.name);
        check(res, { [`${endpoint.name} ok`]: (r) => r.status >= 200 && r.status < 500 });
    });
}

// tests/load/realistic_usage.js - Simulates real user behavior patterns
// Usage: k6 run --env PROFILE=gentle tests/load/realistic_usage.js
import { check, group, sleep } from 'k6';
import { Counter, Trend } from 'k6/metrics';
import { PROFILES } from './config.js';
import {
    authGet,
    ensureUserPool,
    extractItems,
    getToken,
    login,
    randomFrom, sleepWithJitter, vuUser, weightedRandom
} from './helpers.js';

// Metrics
const scenarioCounter = new Counter('scenarios_completed');
const scenarioLatency = new Trend('scenario_latency');
const browsingActions = new Counter('browsing_actions');
const searchActions = new Counter('search_actions');

// Profile selection
const profileName = __ENV.PROFILE || 'smoke';
const profile = PROFILES[profileName] || PROFILES.smoke;

export const options = {
    thresholds: {
        http_req_duration: ['p(95)<1000', 'p(99)<3000'],
        http_req_failed: ['rate<0.05'],
        'scenarios_completed': ['count>10'],
    },
    scenarios: {
        browser: {
            executor: 'ramping-vus',
            stages: profile.stages,
            exec: 'browserSession',
        },
    },
};

// User scenarios with real behavior patterns
const SCENARIOS = [
    { name: 'browse_movies', weight: 25, fn: browseMovies },
    { name: 'browse_tvshows', weight: 25, fn: browseTvShows },
    { name: 'search_content', weight: 15, fn: searchContent },
    { name: 'continue_watching', weight: 15, fn: continueWatching },
    { name: 'account_management', weight: 10, fn: accountManagement },
    { name: 'discover_new', weight: 10, fn: discoverNew },
];

export function setup() {
    // Create/ensure multi-user pool
    const userPool = ensureUserPool();

    // Login as first user to discover content IDs — token is NOT shared with VUs
    const token = login(userPool[0]);
    if (!token) {
        throw new Error('Setup login failed');
    }

    // Pre-fetch some content IDs
    let movieIds = [];
    let tvshowIds = [];
    let seasonIds = [];

    const movies = authGet('/movies?limit=50', token);
    if (movies.status === 200) {
        movieIds = extractItems(movies).map(m => m.id).filter(Boolean);
    }

    const tvshows = authGet('/tvshows?limit=50', token);
    if (tvshows.status === 200) {
        tvshowIds = extractItems(tvshows).map(t => t.id).filter(Boolean);
    }

    // Get some seasons
    if (tvshowIds.length > 0) {
        for (let i = 0; i < Math.min(3, tvshowIds.length); i++) {
            const seasons = authGet(`/tvshows/${tvshowIds[i]}/seasons`, token);
            if (seasons.status === 200) {
                extractItems(seasons).forEach(s => {
                    if (s.id) seasonIds.push(s.id);
                });
            }
        }
    }

    console.log(`Setup: ${movieIds.length} movies, ${tvshowIds.length} shows, ${seasonIds.length} seasons`);
    return { movieIds, tvshowIds, seasonIds, userPool };
}

export function browserSession(data) {
    // Each VU gets its own user from the pool + own token via per-VU login cache
    const user = vuUser(data.userPool);
    const token = getToken(user);
    if (!token) {
        console.error('VU login failed, skipping iteration');
        return;
    }
    const vuData = { ...data, token };

    const scenario = weightedRandom(SCENARIOS);
    const start = Date.now();

    try {
        scenario.fn(vuData);
        scenarioCounter.add(1, { scenario: scenario.name });
        scenarioLatency.add(Date.now() - start, { scenario: scenario.name });
    } catch (e) {
        console.error(`Scenario ${scenario.name} failed: ${e.message}`);
    }

    // Natural pause between user actions
    sleep(sleepWithJitter(1, 0.5));
}

// Scenario: User browses movie catalog
function browseMovies(data) {
    group('Browse Movies', () => {
        const token = data.token;

        // 1. Load movie list
        let res = authGet('/movies?limit=20', token);
        check(res, { 'movies list': (r) => r.status === 200 });
        browsingActions.add(1);
        sleep(sleepWithJitter(0.5));

        // 2. Maybe check recently added
        if (Math.random() > 0.5) {
            res = authGet('/movies/recently-added?limit=10', token);
            check(res, { 'recently added': (r) => r.status === 200 });
            browsingActions.add(1);
            sleep(sleepWithJitter(0.3));
        }

        // 3. Click on a movie to see details
        const movieId = randomFrom(data.movieIds);
        if (movieId) {
            res = authGet(`/movies/${movieId}`, token);
            check(res, { 'movie detail': (r) => r.status === 200 || r.status === 404 });
            browsingActions.add(1);
            sleep(sleepWithJitter(0.8));

            // 4. Load related data (user reads details)
            const detailCalls = [
                () => authGet(`/movies/${movieId}/cast`, token),
                () => authGet(`/movies/${movieId}/crew`, token),
                () => authGet(`/movies/${movieId}/genres`, token),
            ];

            // User might check 1-3 of these
            const numChecks = Math.floor(Math.random() * 3) + 1;
            for (let i = 0; i < numChecks; i++) {
                const call = detailCalls[Math.floor(Math.random() * detailCalls.length)];
                res = call();
                browsingActions.add(1);
                sleep(sleepWithJitter(0.2));
            }

            // 5. Maybe check similar movies
            if (Math.random() > 0.6) {
                res = authGet(`/movies/${movieId}/similar`, token);
                check(res, { 'similar movies': (r) => r.status === 200 || r.status === 404 });
                browsingActions.add(1);
            }
        }
    });
}

// Scenario: User browses TV shows
function browseTvShows(data) {
    group('Browse TV Shows', () => {
        const token = data.token;

        // 1. Load show list
        let res = authGet('/tvshows?limit=20', token);
        check(res, { 'tvshows list': (r) => r.status === 200 });
        browsingActions.add(1);
        sleep(sleepWithJitter(0.5));

        // 2. Check upcoming episodes
        if (Math.random() > 0.5) {
            res = authGet('/tvshows/episodes/upcoming', token);
            check(res, { 'upcoming episodes': (r) => r.status === 200 });
            browsingActions.add(1);
            sleep(sleepWithJitter(0.3));
        }

        // 3. Click on a show
        const tvshowId = randomFrom(data.tvshowIds);
        if (tvshowId) {
            res = authGet(`/tvshows/${tvshowId}`, token);
            check(res, { 'show detail': (r) => r.status === 200 || r.status === 404 });
            browsingActions.add(1);
            sleep(sleepWithJitter(0.6));

            // 4. Load seasons
            res = authGet(`/tvshows/${tvshowId}/seasons`, token);
            check(res, { 'seasons': (r) => r.status === 200 || r.status === 404 });
            browsingActions.add(1);

            const seasons = extractItems(res);
            sleep(sleepWithJitter(0.3));

            // 5. Click on a season to see episodes
            if (seasons.length > 0) {
                const season = randomFrom(seasons);
                if (season && season.id) {
                    res = authGet(`/tvshows/seasons/${season.id}/episodes`, token);
                    check(res, { 'season episodes': (r) => r.status === 200 || r.status === 404 });
                    browsingActions.add(1);
                    sleep(sleepWithJitter(0.4));
                }
            }

            // 6. Check watch progress
            if (Math.random() > 0.5) {
                res = authGet(`/tvshows/${tvshowId}/watch-stats`, token);
                browsingActions.add(1);
            }
        }
    });
}

// Scenario: User searches for content
function searchContent(data) {
    group('Search Content', () => {
        const token = data.token;
        const searchTerms = ['action', 'comedy', 'drama', 'sci-fi', 'thriller', 'horror', 'romance', 'star', 'the', 'night'];
        const term = searchTerms[Math.floor(Math.random() * searchTerms.length)];

        // 1. Initial search (autocomplete as user types)
        let res = authGet(`/search/movies/autocomplete?q=${term.substring(0, 3)}`, token);
        check(res, { 'autocomplete': (r) => r.status === 200 });
        searchActions.add(1);
        sleep(sleepWithJitter(0.2));

        // 2. Full search
        res = authGet(`/search/movies?q=${term}&limit=20`, token);
        check(res, { 'movie search': (r) => r.status === 200 });
        searchActions.add(1);
        sleep(sleepWithJitter(0.3));

        // 3. Maybe search TV shows too
        if (Math.random() > 0.5) {
            res = authGet(`/search/tvshows?q=${term}&limit=20`, token);
            check(res, { 'tvshow search': (r) => r.status === 200 });
            searchActions.add(1);
            sleep(sleepWithJitter(0.3));
        }

        // 4. Multi-search
        if (Math.random() > 0.6) {
            res = authGet(`/search/multi?q=${term}`, token);
            check(res, { 'multi search': (r) => r.status === 200 });
            searchActions.add(1);
        }

        // 5. Check facets for filtering
        if (Math.random() > 0.7) {
            res = authGet('/search/movies/facets', token);
            check(res, { 'facets': (r) => r.status === 200 });
            searchActions.add(1);
        }
    });
}

// Scenario: User checks continue watching
function continueWatching(data) {
    group('Continue Watching', () => {
        const token = data.token;

        // 1. Check movie progress
        let res = authGet('/movies/continue-watching', token);
        check(res, { 'movies continue': (r) => r.status === 200 });
        browsingActions.add(1);
        sleep(sleepWithJitter(0.3));

        // 2. Check TV show progress
        res = authGet('/tvshows/continue-watching', token);
        check(res, { 'tvshows continue': (r) => r.status === 200 });
        browsingActions.add(1);
        sleep(sleepWithJitter(0.3));

        // 3. Check watch history
        if (Math.random() > 0.5) {
            res = authGet('/movies/watch-history', token);
            check(res, { 'watch history': (r) => r.status === 200 });
            browsingActions.add(1);
        }

        // 4. Check stats
        if (Math.random() > 0.6) {
            res = authGet('/movies/stats', token);
            check(res, { 'movie stats': (r) => r.status === 200 });
            res = authGet('/tvshows/stats', token);
            check(res, { 'tvshow stats': (r) => r.status === 200 });
            browsingActions.add(2);
        }
    });
}

// Scenario: User manages account
function accountManagement(data) {
    group('Account Management', () => {
        const token = data.token;

        // 1. Check profile
        let res = authGet('/users/me', token);
        check(res, { 'get me': (r) => r.status === 200 });
        sleep(sleepWithJitter(0.3));

        // 2. Check sessions
        res = authGet('/sessions', token);
        check(res, { 'sessions': (r) => r.status === 200 });
        sleep(sleepWithJitter(0.2));

        // 3. Check current session
        res = authGet('/sessions/current', token);
        check(res, { 'current session': (r) => r.status === 200 });
        sleep(sleepWithJitter(0.2));

        // 4. Check MFA status
        res = authGet('/mfa/status', token);
        check(res, { 'mfa status': (r) => r.status === 200 });
        sleep(sleepWithJitter(0.2));

        // 5. Check API keys
        res = authGet('/apikeys', token);
        check(res, { 'apikeys': (r) => r.status === 200 });
        sleep(sleepWithJitter(0.2));

        // 6. Check preferences
        res = authGet('/users/me/preferences', token);
        check(res, { 'preferences': (r) => r.status === 200 });

        // 7. Check OIDC links
        if (Math.random() > 0.7) {
            res = authGet('/users/me/oidc', token);
            check(res, { 'oidc links': (r) => r.status === 200 });
        }
    });
}

// Scenario: User discovers new content
function discoverNew(data) {
    group('Discover New', () => {
        const token = data.token;

        // 1. Check recently added movies
        let res = authGet('/movies/recently-added?limit=15', token);
        check(res, { 'new movies': (r) => r.status === 200 });
        browsingActions.add(1);
        sleep(sleepWithJitter(0.4));

        // 2. Check top rated
        res = authGet('/movies/top-rated?limit=15', token);
        check(res, { 'top rated': (r) => r.status === 200 });
        browsingActions.add(1);
        sleep(sleepWithJitter(0.4));

        // 3. Check recently added TV
        res = authGet('/tvshows/recently-added?limit=15', token);
        check(res, { 'new tvshows': (r) => r.status === 200 });
        browsingActions.add(1);
        sleep(sleepWithJitter(0.4));

        // 4. Check recent episodes
        res = authGet('/tvshows/episodes/recent?limit=15', token);
        check(res, { 'recent episodes': (r) => r.status === 200 });
        browsingActions.add(1);
        sleep(sleepWithJitter(0.3));

        // 5. Browse genres
        res = authGet('/genres', token);
        check(res, { 'genres': (r) => r.status === 200 });
        browsingActions.add(1);

        // 6. Maybe check libraries
        if (Math.random() > 0.6) {
            res = authGet('/libraries', token);
            check(res, { 'libraries': (r) => r.status === 200 });
            browsingActions.add(1);
        }
    });
}

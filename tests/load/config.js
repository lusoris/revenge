// tests/load/config.js - Shared configuration for k6 load tests
export const BASE_URL = __ENV.BASE_URL || 'http://localhost:8096';
export const API_BASE = `${BASE_URL}/api/v1`;

// Default test credentials (single-user fallback)
export const TEST_USER = {
    username: __ENV.TEST_USER || 'dbg_ext1',
    password: __ENV.TEST_PASS || 'TestPass123!',
};

// Admin credentials for admin-only endpoints
export const ADMIN_USER = {
    username: __ENV.ADMIN_USER || 'admin',
    password: __ENV.ADMIN_PASS || 'Admin123!',
};

// Number of load test users to create/use (distributed across VUs)
export const USER_POOL_SIZE = parseInt(__ENV.USER_POOL_SIZE || '20');

// Load test user pool — generated in setup(), distributed across VUs
// Each user: { username: 'loadtest_01', password: 'TestPass123!' }
export function buildUserPool(size = USER_POOL_SIZE) {
    const pool = [];
    for (let i = 1; i <= size; i++) {
        const idx = String(i).padStart(2, '0');
        pool.push({
            username: `loadtest_${idx}`,
            email: `lt${idx}@loadtest.local`,
            password: 'TestPass123!',
        });
    }
    return pool;
}

// Thresholds for different load profiles
// Note: http_req_failed counts non-2xx as failures, but 404/403 are valid API responses
export const THRESHOLDS = {
    gentle: {
        http_req_duration: ['p(95)<500', 'p(99)<1000'],
        http_req_failed: ['rate<0.30'], // Allow 404/403 responses
    },
    spike: {
        http_req_duration: ['p(95)<2000', 'p(99)<5000'],
        http_req_failed: ['rate<0.40'],
    },
    soak: {
        http_req_duration: ['p(95)<500', 'p(99)<1500'],
        http_req_failed: ['rate<0.30'],
    },
    stress: {
        http_req_duration: ['p(95)<3000'],
        http_req_failed: ['rate<0.50'],
    },
};

// Load profiles
export const PROFILES = {
    // Gentle ramp: 50 → 500 VUs over 5min
    gentle: {
        stages: [
            { duration: '30s', target: 50 },
            { duration: '1m', target: 100 },
            { duration: '1m', target: 200 },
            { duration: '1m', target: 350 },
            { duration: '1m', target: 500 },
            { duration: '30s', target: 0 },
        ],
    },
    // Spike: 0 → 1000 instantly
    spike: {
        stages: [
            { duration: '10s', target: 100 },
            { duration: '1s', target: 1000 },
            { duration: '1m', target: 1000 },
            { duration: '30s', target: 100 },
            { duration: '30s', target: 0 },
        ],
    },
    // Soak: 100 VUs for 30min
    soak: {
        stages: [
            { duration: '1m', target: 100 },
            { duration: '28m', target: 100 },
            { duration: '1m', target: 0 },
        ],
    },
    // Stress: keep adding until failure
    stress: {
        stages: [
            { duration: '30s', target: 100 },
            { duration: '1m', target: 300 },
            { duration: '1m', target: 500 },
            { duration: '1m', target: 800 },
            { duration: '1m', target: 1200 },
            { duration: '1m', target: 1500 },
            { duration: '1m', target: 2000 },
            { duration: '30s', target: 0 },
        ],
    },
    // Quick smoke test for CI
    smoke: {
        stages: [
            { duration: '10s', target: 5 },
            { duration: '30s', target: 10 },
            { duration: '10s', target: 0 },
        ],
    },
};

// Weighted endpoint categories for realistic traffic distribution
export const ENDPOINT_WEIGHTS = {
    // Most frequent - 40% of traffic
    highFrequency: [
        { method: 'GET', path: '/movies', weight: 15 },
        { method: 'GET', path: '/tvshows', weight: 15 },
        { method: 'GET', path: '/search/movies', weight: 5 },
        { method: 'GET', path: '/search/tvshows', weight: 5 },
    ],
    // Medium frequency - 30% of traffic
    mediumFrequency: [
        { method: 'GET', path: '/movies/recently-added', weight: 5 },
        { method: 'GET', path: '/tvshows/recently-added', weight: 5 },
        { method: 'GET', path: '/movies/continue-watching', weight: 5 },
        { method: 'GET', path: '/tvshows/continue-watching', weight: 5 },
        { method: 'GET', path: '/users/me', weight: 5 },
        { method: 'GET', path: '/sessions', weight: 5 },
    ],
    // Low frequency - 20% of traffic
    lowFrequency: [
        { method: 'GET', path: '/settings/user', weight: 3 },
        { method: 'GET', path: '/genres', weight: 3 },
        { method: 'GET', path: '/libraries', weight: 4 },
        { method: 'GET', path: '/apikeys', weight: 3 },
        { method: 'GET', path: '/mfa/status', weight: 2 },
        { method: 'GET', path: '/rbac/roles', weight: 2 },
        { method: 'GET', path: '/oidc/providers', weight: 3 },
    ],
    // Auth operations - 10% of traffic
    auth: [
        { method: 'POST', path: '/auth/refresh', weight: 7 },
        { method: 'GET', path: '/sessions/current', weight: 3 },
    ],
};

// tests/load/auth_stress.js - Stress test authentication system
// Tests login, token refresh, session management under load
// Usage: k6 run --env PROFILE=gentle tests/load/auth_stress.js
import { check, group, sleep } from 'k6';
import http from 'k6/http';
import { Counter, Rate, Trend } from 'k6/metrics';
import { API_BASE, PROFILES, TEST_USER } from './config.js';
import { randomString, sleepWithJitter } from './helpers.js';

// Metrics
const loginAttempts = new Counter('auth_login_attempts');
const loginSuccess = new Rate('auth_login_success');
const refreshAttempts = new Counter('auth_refresh_attempts');
const refreshSuccess = new Rate('auth_refresh_success');
const tokenLatency = new Trend('auth_token_latency');

const profileName = __ENV.PROFILE || 'smoke';
const profile = PROFILES[profileName] || PROFILES.smoke;

export const options = {
    stages: profile.stages,
    thresholds: {
        http_req_duration: ['p(95)<1000', 'p(99)<2000'],
        'auth_login_success': ['rate>0.95'],
        'auth_refresh_success': ['rate>0.95'],
        'auth_token_latency': ['p(95)<500'],
    },
};

export function setup() {
    // Verify basic auth works
    const res = http.post(
        `${API_BASE}/auth/login`,
        JSON.stringify(TEST_USER),
        { headers: { 'Content-Type': 'application/json' } }
    );

    if (res.status !== 200) {
        throw new Error(`Setup auth failed: ${res.status} ${res.body}`);
    }

    const body = JSON.parse(res.body);
    console.log('Auth setup verified');
    return { baseRefreshToken: body.refresh_token };
}

export default function(data) {
    // Randomly choose auth operation
    const operations = [
        { weight: 40, fn: () => testLogin() },
        { weight: 30, fn: () => testRefresh(data) },
        { weight: 20, fn: () => testSessionOps() },
        { weight: 10, fn: () => testMFAStatus() },
    ];

    // Weighted selection
    const total = operations.reduce((sum, op) => sum + op.weight, 0);
    let random = Math.random() * total;
    for (const op of operations) {
        random -= op.weight;
        if (random <= 0) {
            op.fn();
            break;
        }
    }

    sleep(sleepWithJitter(0.2, 0.1));
}

function testLogin() {
    group('Login Flow', () => {
        const start = Date.now();

        // Successful login
        let res = http.post(
            `${API_BASE}/auth/login`,
            JSON.stringify(TEST_USER),
            { headers: { 'Content-Type': 'application/json' }, tags: { name: 'login' } }
        );

        loginAttempts.add(1);
        const success = res.status === 200;
        loginSuccess.add(success ? 1 : 0);
        tokenLatency.add(Date.now() - start, { operation: 'login' });

        check(res, {
            'login returns 200': (r) => r.status === 200,
            'has access_token': (r) => {
                try {
                    return JSON.parse(r.body).access_token !== undefined;
                } catch { return false; }
            },
            'has refresh_token': (r) => {
                try {
                    return JSON.parse(r.body).refresh_token !== undefined;
                } catch { return false; }
            },
        });

        sleep(sleepWithJitter(0.1));

        // Test invalid login (should be fast rejection)
        if (Math.random() > 0.7) {
            const badRes = http.post(
                `${API_BASE}/auth/login`,
                JSON.stringify({ username: 'nonexistent', password: 'wrong' }),
                { headers: { 'Content-Type': 'application/json' }, tags: { name: 'login_invalid' } }
            );

            check(badRes, {
                'invalid login returns 401': (r) => r.status === 401,
            });
        }
    });
}

function testRefresh(data) {
    group('Token Refresh', () => {
        // First login to get fresh tokens
        const loginRes = http.post(
            `${API_BASE}/auth/login`,
            JSON.stringify(TEST_USER),
            { headers: { 'Content-Type': 'application/json' } }
        );

        if (loginRes.status !== 200) {
            console.error('Could not get tokens for refresh test');
            return;
        }

        const tokens = JSON.parse(loginRes.body);
        sleep(sleepWithJitter(0.1));

        // Refresh the token
        const start = Date.now();
        const res = http.post(
            `${API_BASE}/auth/refresh`,
            JSON.stringify({ refresh_token: tokens.refresh_token }),
            { headers: { 'Content-Type': 'application/json' }, tags: { name: 'refresh' } }
        );

        refreshAttempts.add(1);
        const success = res.status === 200;
        refreshSuccess.add(success ? 1 : 0);
        tokenLatency.add(Date.now() - start, { operation: 'refresh' });

        check(res, {
            'refresh returns 200': (r) => r.status === 200,
            'new access_token': (r) => {
                try {
                    return JSON.parse(r.body).access_token !== undefined;
                } catch { return false; }
            },
        });

        // Test invalid refresh token
        if (Math.random() > 0.8) {
            const badRes = http.post(
                `${API_BASE}/auth/refresh`,
                JSON.stringify({ refresh_token: 'invalid-token-' + randomString(16) }),
                { headers: { 'Content-Type': 'application/json' }, tags: { name: 'refresh_invalid' } }
            );

            check(badRes, {
                'invalid refresh rejected': (r) => r.status === 401 || r.status === 400,
            });
        }
    });
}

function testSessionOps() {
    group('Session Operations', () => {
        // Login first
        const loginRes = http.post(
            `${API_BASE}/auth/login`,
            JSON.stringify(TEST_USER),
            { headers: { 'Content-Type': 'application/json' } }
        );

        if (loginRes.status !== 200) return;

        const tokens = JSON.parse(loginRes.body);
        const headers = {
            'Authorization': `Bearer ${tokens.access_token}`,
            'Content-Type': 'application/json',
        };

        sleep(sleepWithJitter(0.1));

        // List sessions
        let res = http.get(`${API_BASE}/sessions`, { headers, tags: { name: 'list_sessions' } });
        check(res, { 'list sessions': (r) => r.status === 200 });

        // Get current session
        res = http.get(`${API_BASE}/sessions/current`, { headers, tags: { name: 'current_session' } });
        check(res, { 'current session': (r) => r.status === 200 });

        sleep(sleepWithJitter(0.1));

        // Session refresh
        res = http.post(`${API_BASE}/sessions/refresh`, null, { headers, tags: { name: 'session_refresh' } });
        check(res, { 'session refresh': (r) => r.status === 200 });
    });
}

function testMFAStatus() {
    group('MFA Status Check', () => {
        // Login first
        const loginRes = http.post(
            `${API_BASE}/auth/login`,
            JSON.stringify(TEST_USER),
            { headers: { 'Content-Type': 'application/json' } }
        );

        if (loginRes.status !== 200) return;

        const tokens = JSON.parse(loginRes.body);
        const headers = {
            'Authorization': `Bearer ${tokens.access_token}`,
            'Content-Type': 'application/json',
        };

        // Check MFA status
        const res = http.get(`${API_BASE}/mfa/status`, { headers, tags: { name: 'mfa_status' } });
        check(res, {
            'mfa status returns': (r) => r.status === 200,
            'has enabled field': (r) => {
                try {
                    const body = JSON.parse(r.body);
                    return body.enabled !== undefined;
                } catch { return false; }
            },
        });
    });
}

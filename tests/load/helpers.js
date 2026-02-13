// tests/load/helpers.js - Shared helper functions for k6 load tests
import { check } from 'k6';
import exec from 'k6/execution';
import http from 'k6/http';
import { API_BASE, TEST_USER } from './config.js';

// Cache tokens to avoid excessive logins
let cachedTokens = {};

/**
 * Generate a unique synthetic IP per VU.
 * Maps VU ID (1-based) to a 10.x.y.z address so each VU
 * is treated as a distinct client by the rate limiter.
 * @returns {string} IP like "10.0.0.1"
 */
export function vuIP() {
    const id = exec.vu.idInTest || 1;
    const a = (id >> 16) & 0xff;
    const b = (id >> 8) & 0xff;
    const c = id & 0xff;
    return `10.${a}.${b}.${c}`;
}

/**
 * Return an X-Forwarded-For header for the current VU.
 * Merge into any headers object to simulate distinct client IPs.
 */
export function forwardedHeaders(extra = {}) {
    return { 'X-Forwarded-For': vuIP(), ...extra };
}

/**
 * Login and get access token
 * @param {Object} user - User credentials {username, password}
 * @returns {string} Access token
 */
export function login(user = TEST_USER) {
    const cacheKey = user.username;
    if (cachedTokens[cacheKey]) {
        return cachedTokens[cacheKey];
    }

    const res = http.post(
        `${API_BASE}/auth/login`,
        JSON.stringify({
            username: user.username,
            password: user.password,
        }),
        {
            headers: forwardedHeaders({ 'Content-Type': 'application/json' }),
            tags: { name: 'auth_login' },
        }
    );

    const success = check(res, {
        'login successful': (r) => r.status === 200,
        'has access_token': (r) => {
            try {
                const body = JSON.parse(r.body);
                return body.access_token !== undefined;
            } catch {
                return false;
            }
        },
    });

    if (!success) {
        console.error(`Login failed for ${user.username}: ${res.status} ${res.body}`);
        return null;
    }

    const body = JSON.parse(res.body);
    cachedTokens[cacheKey] = body.access_token;
    return body.access_token;
}

/**
 * Get authorization headers with token
 * @param {string} token - Access token
 * @returns {Object} Headers object
 */
export function authHeaders(token) {
    return {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
        'X-Forwarded-For': vuIP(),
    };
}

/**
 * Make authenticated GET request
 * @param {string} path - API path (without /api/v1 prefix)
 * @param {string} token - Access token
 * @param {Object} params - Additional params
 */
export function authGet(path, token, params = {}) {
    return http.get(`${API_BASE}${path}`, {
        headers: authHeaders(token),
        tags: { name: path, ...params.tags },
        ...params,
    });
}

/**
 * Make authenticated POST request
 * @param {string} path - API path (without /api/v1 prefix)
 * @param {Object} body - Request body
 * @param {string} token - Access token
 * @param {Object} params - Additional params
 */
export function authPost(path, body, token, params = {}) {
    return http.post(`${API_BASE}${path}`, JSON.stringify(body), {
        headers: authHeaders(token),
        tags: { name: path, ...params.tags },
        ...params,
    });
}

/**
 * Make authenticated PUT request
 */
export function authPut(path, body, token, params = {}) {
    return http.put(`${API_BASE}${path}`, JSON.stringify(body), {
        headers: authHeaders(token),
        tags: { name: path, ...params.tags },
        ...params,
    });
}

/**
 * Make authenticated DELETE request
 */
export function authDelete(path, token, params = {}) {
    return http.del(`${API_BASE}${path}`, null, {
        headers: authHeaders(token),
        tags: { name: path, ...params.tags },
        ...params,
    });
}

/**
 * Standard response checks
 * @param {Object} res - HTTP response
 * @param {number} expectedStatus - Expected status code (default 200)
 */
export function checkResponse(res, expectedStatus = 200) {
    const checks = {};
    checks[`status is ${expectedStatus}`] = (r) => r.status === expectedStatus;
    checks['response time < 500ms'] = (r) => r.timings.duration < 500;
    return check(res, checks);
}

/**
 * Check if response is success (2xx)
 */
export function checkSuccess(res) {
    return check(res, {
        'status is 2xx': (r) => r.status >= 200 && r.status < 300,
    });
}

/**
 * Check if response is JSON
 */
export function checkJSON(res) {
    return check(res, {
        'is JSON': (r) => {
            const ct = r.headers['Content-Type'];
            return ct && ct.includes('application/json');
        },
    });
}

/**
 * Extract items from list response
 */
export function extractItems(res) {
    try {
        const body = JSON.parse(res.body);
        return body.items || body.data || body.movies || body.tvshows || body.sessions || [];
    } catch {
        return [];
    }
}

/**
 * Get random element from array
 */
export function randomFrom(arr) {
    if (!arr || arr.length === 0) return null;
    return arr[Math.floor(Math.random() * arr.length)];
}

/**
 * Generate random string
 */
export function randomString(length = 8) {
    const chars = 'abcdefghijklmnopqrstuvwxyz0123456789';
    let result = '';
    for (let i = 0; i < length; i++) {
        result += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    return result;
}

/**
 * Sleep with random jitter
 * @param {number} base - Base sleep time in seconds
 * @param {number} jitter - Max jitter in seconds
 */
export function sleepWithJitter(base, jitter = 0.5) {
    const actual = base + (Math.random() * jitter * 2 - jitter);
    return Math.max(0.1, actual);
}

/**
 * Weighted random selection
 * @param {Array} items - Array of {item, weight} objects
 */
export function weightedRandom(items) {
    const totalWeight = items.reduce((sum, i) => sum + i.weight, 0);
    let random = Math.random() * totalWeight;

    for (const item of items) {
        random -= item.weight;
        if (random <= 0) {
            return item;
        }
    }
    return items[items.length - 1];
}

/**
 * Login and return full token response (access + refresh tokens)
 * @param {Object} user - User credentials {username, password}
 * @returns {Object} {access_token, refresh_token} or null
 */
export function loginFull(user = TEST_USER) {
    const res = http.post(
        `${API_BASE}/auth/login`,
        JSON.stringify({
            username: user.username,
            password: user.password,
        }),
        {
            headers: forwardedHeaders({ 'Content-Type': 'application/json' }),
            tags: { name: 'auth_login' },
        }
    );

    if (res.status !== 200) {
        console.error(`Login failed for ${user.username}: ${res.status}`);
        return null;
    }

    return JSON.parse(res.body);
}

/**
 * Make request authenticated with an API key
 * @param {string} method - HTTP method (GET, POST, PUT, DELETE)
 * @param {string} path - API path
 * @param {string} apiKey - Raw API key
 * @param {Object} body - Optional request body
 * @param {Object} params - Additional params
 */
export function apiKeyRequest(method, path, apiKey, body = null, params = {}) {
    const headers = {
        'X-API-Key': apiKey,
        'Content-Type': 'application/json',
        'X-Forwarded-For': vuIP(),
    };
    const url = `${API_BASE}${path}`;
    const opts = { headers, tags: { name: path, ...params.tags }, ...params };

    switch (method.toUpperCase()) {
        case 'GET':
            return http.get(url, opts);
        case 'POST':
            return http.post(url, body ? JSON.stringify(body) : null, opts);
        case 'PUT':
            return http.put(url, body ? JSON.stringify(body) : null, opts);
        case 'DELETE':
            return http.del(url, null, opts);
        default:
            return http.get(url, opts);
    }
}

/**
 * Generate a random integer in [min, max]
 */
export function randomInt(min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

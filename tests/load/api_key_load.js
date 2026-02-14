// tests/load/api_key_load.js - Stress test API key authentication and caching
// Tests API key CRUD, auth-via-API-key performance, and cache efficiency
// Usage: k6 run --env PROFILE=gentle tests/load/api_key_load.js
import { check, group, sleep } from 'k6';
import { Counter, Rate, Trend } from 'k6/metrics';
import { PROFILES } from './config.js';
import {
    apiKeyRequest,
    authDelete,
    authGet,
    authPost,
    ensureUserPool,
    getToken,
    login,
    randomFrom,
    randomInt,
    randomString,
    sleepWithJitter,
    vuUser,
    weightedRandom
} from './helpers.js';

// Custom metrics
const apiKeyCreated = new Counter('apikey_created');
const apiKeyRevoked = new Counter('apikey_revoked');
const apiKeyAuthAttempts = new Counter('apikey_auth_attempts');
const apiKeyAuthSuccess = new Rate('apikey_auth_success');
const apiKeyCacheHitLatency = new Trend('apikey_cache_hit_latency');
const apiKeyColdLatency = new Trend('apikey_cold_latency');
const apiKeyCreateLatency = new Trend('apikey_create_latency');

const profileName = __ENV.PROFILE || 'smoke';
const profile = PROFILES[profileName] || PROFILES.smoke;

export const options = {
    stages: profile.stages,
    thresholds: {
        http_req_duration: ['p(95)<800', 'p(99)<2000'],
        http_req_failed: ['rate<0.20'],
        'apikey_auth_success': ['rate>0.85'],
        'apikey_cache_hit_latency': ['p(95)<100'],   // Cache hits should be fast
        'apikey_create_latency': ['p(95)<500'],
    },
};

const SCENARIOS = [
    { name: 'cached_auth', weight: 40, fn: cachedAuthPerformance },
    { name: 'lifecycle', weight: 25, fn: apiKeyLifecycle },
    { name: 'concurrent_keys', weight: 15, fn: concurrentKeys },
    { name: 'invalid_keys', weight: 10, fn: invalidKeyRejection },
    { name: 'key_management', weight: 10, fn: keyManagement },
];

export function setup() {
    // Create/ensure multi-user pool
    const userPool = ensureUserPool();

    const token = login(userPool[0]);
    if (!token) {
        throw new Error('Setup login failed');
    }

    // Pre-create a pool of API keys for cache testing
    const apiKeys = [];
    const keyIds = [];

    for (let i = 0; i < 5; i++) {
        const res = authPost('/apikeys', {
            name: `k6-load-pool-${i}-${randomString(4)}`,
            scopes: ['read'],
        }, token);

        if (res.status === 201) {
            try {
                const body = JSON.parse(res.body);
                const rawKey = body.api_key;
                const keyId = body.id;
                if (rawKey) {
                    apiKeys.push(rawKey);
                    if (keyId) keyIds.push(keyId);
                }
            } catch { /* ignore */ }
        }
    }

    // Warm the cache by making a request with each key
    for (const key of apiKeys) {
        apiKeyRequest('GET', '/users/me', key);
    }

    // Discover some content IDs for API-key-authed reads
    let movieIds = [];
    const movies = authGet('/movies?limit=20', token);
    if (movies.status === 200) {
        try {
            const body = JSON.parse(movies.body);
            movieIds = (body.items || body.movies || []).map(m => m.id).filter(Boolean);
        } catch { /* ignore */ }
    }

    console.log(`Setup: ${apiKeys.length} pre-created API keys, ${movieIds.length} movies`);
    return { apiKeys, keyIds, movieIds, userPool };
}

export default function (data) {
    // Each VU gets its own user from the pool + own token via per-VU login cache
    const user = vuUser(data.userPool);
    const token = getToken(user);
    if (!token) {
        console.error('VU login failed, skipping iteration');
        return;
    }
    const vuData = { ...data, token };

    const scenario = weightedRandom(SCENARIOS);
    try {
        scenario.fn(vuData);
    } catch (e) {
        console.error(`Scenario ${scenario.name} failed: ${e.message}`);
    }
    sleep(sleepWithJitter(0.3, 0.2));
}

// Scenario 1: Measure cached API key auth performance
// Uses pre-created keys (already cached) for repeated requests
function cachedAuthPerformance(data) {
    group('Cached API Key Auth', () => {
        const key = randomFrom(data.apiKeys);
        if (!key) return;

        // Multiple requests with the same key — should hit cache
        const endpoints = ['/users/me', '/movies', '/tvshows', '/genres', '/libraries'];
        const numRequests = randomInt(3, 8);

        for (let i = 0; i < numRequests; i++) {
            const endpoint = randomFrom(endpoints);
            const start = Date.now();

            const res = apiKeyRequest('GET', endpoint, key);
            const latency = Date.now() - start;

            apiKeyCacheHitLatency.add(latency);
            apiKeyAuthAttempts.add(1);

            const ok = res.status >= 200 && res.status < 300;
            apiKeyAuthSuccess.add(ok ? 1 : 0);

            check(res, {
                'cached auth ok': (r) => r.status >= 200 && r.status < 300,
            });

            sleep(sleepWithJitter(0.05, 0.02));
        }
    });
}

// Scenario 2: Full API key lifecycle (create → use → revoke)
function apiKeyLifecycle(data) {
    group('API Key Lifecycle', () => {
        const token = data.token;

        // 1. Create a new API key
        const createStart = Date.now();
        const createRes = authPost('/apikeys', {
            name: `k6-lifecycle-${randomString(6)}`,
            scopes: ['read', 'write'],
        }, token);

        apiKeyCreateLatency.add(Date.now() - createStart);

        const created = check(createRes, {
            'key created': (r) => r.status === 201,
        });

        if (!created) return;

        apiKeyCreated.add(1);
        let body;
        try {
            body = JSON.parse(createRes.body);
        } catch { return; }

        const rawKey = body.api_key;
        const keyId = body.id;
        if (!rawKey) return;

        sleep(sleepWithJitter(0.1));

        // 2. First use (cold — not cached yet)
        const coldStart = Date.now();
        let res = apiKeyRequest('GET', '/users/me', rawKey);
        apiKeyColdLatency.add(Date.now() - coldStart);
        apiKeyAuthAttempts.add(1);
        apiKeyAuthSuccess.add(res.status >= 200 && res.status < 300 ? 1 : 0);

        check(res, {
            'cold auth works': (r) => r.status >= 200 && r.status < 300,
        });

        sleep(sleepWithJitter(0.1));

        // 3. Second use (warm — should be cached)
        const warmStart = Date.now();
        res = apiKeyRequest('GET', '/movies', rawKey);
        apiKeyCacheHitLatency.add(Date.now() - warmStart);
        apiKeyAuthAttempts.add(1);
        apiKeyAuthSuccess.add(res.status >= 200 && res.status < 300 ? 1 : 0);

        check(res, {
            'warm auth works': (r) => r.status >= 200 && r.status < 300,
        });

        sleep(sleepWithJitter(0.1));

        // 4. Revoke key
        if (keyId) {
            const revokeRes = authDelete(`/apikeys/${keyId}`, token);
            check(revokeRes, {
                'key revoked': (r) => r.status === 204,
            });
            apiKeyRevoked.add(1);

            sleep(sleepWithJitter(0.1));

            // 5. Verify revoked key no longer works
            res = apiKeyRequest('GET', '/users/me', rawKey);
            check(res, {
                'revoked key rejected': (r) => r.status === 401 || r.status === 403,
            });
        }
    });
}

// Scenario 3: Multiple concurrent API keys from same user
function concurrentKeys(data) {
    group('Concurrent Keys', () => {
        const token = data.token;
        const createdKeys = [];

        // Create multiple keys rapidly
        const numKeys = randomInt(2, 4);
        for (let i = 0; i < numKeys; i++) {
            const res = authPost('/apikeys', {
                name: `k6-concurrent-${i}-${randomString(4)}`,
                scopes: ['read'],
            }, token);

            if (res.status === 201) {
                try {
                    const body = JSON.parse(res.body);
                    const rawKey = body.api_key;
                    const keyId = body.id;
                    if (rawKey) createdKeys.push({ rawKey, keyId });
                } catch { /* ignore */ }
                apiKeyCreated.add(1);
            }
        }

        // Use all keys concurrently (simulate multiple client integrations)
        for (const { rawKey } of createdKeys) {
            const res = apiKeyRequest('GET', '/movies?limit=5', rawKey);
            apiKeyAuthAttempts.add(1);
            apiKeyAuthSuccess.add(res.status >= 200 && res.status < 300 ? 1 : 0);
        }

        sleep(sleepWithJitter(0.2));

        // Clean up
        for (const { keyId } of createdKeys) {
            if (keyId) {
                authDelete(`/apikeys/${keyId}`, token);
                apiKeyRevoked.add(1);
            }
        }
    });
}

// Scenario 4: Invalid API key rejection (security test)
function invalidKeyRejection(data) {
    group('Invalid Key Rejection', () => {
        const invalidKeys = [
            `rvg_${randomString(32)}`,             // random plausible format
            randomString(40),                       // random gibberish
            '',                                     // empty
            'a'.repeat(200),                        // oversized key
            `rvg_${randomString(8)}_modified`,      // tampered key
        ];

        for (const key of invalidKeys) {
            if (!key) continue; // skip empty for header issues

            const res = apiKeyRequest('GET', '/users/me', key);
            apiKeyAuthAttempts.add(1);

            check(res, {
                'invalid key rejected': (r) => r.status === 401 || r.status === 403,
            });

            // Invalid keys should be rejected fast (no DB lookup if format wrong)
            check(res, {
                'fast rejection': (r) => r.timings.duration < 200,
            });
        }
    });
}

// Scenario 5: List and inspect existing keys
function keyManagement(data) {
    group('Key Management', () => {
        const token = data.token;

        // List all API keys
        let res = authGet('/apikeys', token);
        check(res, {
            'list keys': (r) => r.status === 200,
        });

        if (res.status === 200) {
            try {
                const body = JSON.parse(res.body);
                const keys = body.keys || [];

                // Inspect a random key
                if (keys.length > 0) {
                    const key = randomFrom(keys);
                    const keyId = key.id;
                    if (keyId) {
                        res = authGet(`/apikeys/${keyId}`, token);
                        check(res, {
                            'get key detail': (r) => r.status === 200,
                        });
                    }
                }
            } catch { /* ignore */ }
        }

        sleep(sleepWithJitter(0.2));
    });
}

export function teardown(data) {
    // Clean up pre-created pool keys
    const token = login(data.userPool[0]);
    for (const keyId of data.keyIds) {
        authDelete(`/apikeys/${keyId}`, token);
    }
    console.log(`Teardown: cleaned up ${data.keyIds.length} pool keys`);
}

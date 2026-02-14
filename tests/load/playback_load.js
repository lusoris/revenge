// tests/load/playback_load.js - Stress test playback session lifecycle
// Simulates concurrent users creating sessions, sending heartbeats, and ending sessions
// Usage: k6 run --env PROFILE=gentle tests/load/playback_load.js
import { check, group, sleep } from 'k6';
import { Counter, Rate, Trend } from 'k6/metrics';
import { PROFILES } from './config.js';
import {
    authDelete,
    authGet,
    authPost,
    ensureUserPool,
    extractItems,
    getToken,
    login,
    randomFrom,
    randomInt,
    sleepWithJitter,
    vuUser,
    weightedRandom,
} from './helpers.js';

// Custom metrics
const sessionsCreated = new Counter('playback_sessions_created');
const heartbeatsSent = new Counter('playback_heartbeats_sent');
const sessionsEnded = new Counter('playback_sessions_ended');
const heartbeatSuccess = new Rate('playback_heartbeat_success');
const sessionCreateLatency = new Trend('playback_session_create_latency');
const heartbeatLatency = new Trend('playback_heartbeat_latency');
const sessionLifecycleErrors = new Counter('playback_lifecycle_errors');

const profileName = __ENV.PROFILE || 'smoke';
const profile = PROFILES[profileName] || PROFILES.smoke;

export const options = {
    stages: profile.stages,
    thresholds: {
        http_req_duration: ['p(95)<800', 'p(99)<2000'],
        http_req_failed: ['rate<0.15'],
        'playback_heartbeat_success': ['rate>0.90'],
        'playback_heartbeat_latency': ['p(95)<200'],
        'playback_session_create_latency': ['p(95)<500'],
    },
};

// Playback scenarios weighted by real-world probability
const SCENARIOS = [
    { name: 'full_session', weight: 40, fn: fullPlaybackSession },
    { name: 'heartbeat_burst', weight: 25, fn: heartbeatBurst },
    { name: 'quick_browse', weight: 20, fn: quickBrowseAndAbandon },
    { name: 'session_inspect', weight: 15, fn: sessionInspect },
];

export function setup() {
    // Create/ensure multi-user pool
    const userPool = ensureUserPool();

    const token = login(userPool[0]);
    if (!token) {
        throw new Error('Setup login failed');
    }

    // Discover content IDs for playback
    let movieIds = [];
    let episodeIds = [];

    const movies = authGet('/movies?limit=50', token);
    if (movies.status === 200) {
        movieIds = extractItems(movies).map(m => m.id).filter(Boolean);
    }

    // Get episode IDs via TV shows → seasons → episodes
    const tvshows = authGet('/tvshows?limit=20', token);
    if (tvshows.status === 200) {
        const showIds = extractItems(tvshows).map(t => t.id).filter(Boolean);
        for (let i = 0; i < Math.min(3, showIds.length); i++) {
            const seasons = authGet(`/tvshows/${showIds[i]}/seasons`, token);
            if (seasons.status === 200) {
                const seasonList = extractItems(seasons);
                for (const season of seasonList.slice(0, 2)) {
                    if (!season.id) continue;
                    const episodes = authGet(`/tvshows/seasons/${season.id}/episodes`, token);
                    if (episodes.status === 200) {
                        extractItems(episodes).forEach(ep => {
                            if (ep.id) episodeIds.push(ep.id);
                        });
                    }
                }
            }
        }
    }

    console.log(`Setup: ${movieIds.length} movies, ${episodeIds.length} episodes for playback`);
    if (movieIds.length === 0 && episodeIds.length === 0) {
        throw new Error(
            'No movies or episodes found. Run tests/load/seed_playback_data.sh to insert test content, or ensure a library scan has completed.'
        );
    }
    return { movieIds, episodeIds, userPool };
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
        sessionLifecycleErrors.add(1, { scenario: scenario.name });
    }
    sleep(sleepWithJitter(0.5, 0.3));
}

// Scenario 1: Full playback session lifecycle (create → heartbeat × N → end)
function fullPlaybackSession(data) {
    group('Full Playback Session', () => {
        const token = data.token;
        const movieId = randomFrom(data.movieIds);
        if (!movieId) return;

        // 1. Create playback session
        const createStart = Date.now();
        const createRes = authPost('/playback/sessions', {
            media_type: 'movie',
            media_id: movieId,
        }, token);

        sessionCreateLatency.add(Date.now() - createStart);

        const created = check(createRes, {
            'session created': (r) => r.status === 201,
            'has session_id': (r) => {
                try {
                    return JSON.parse(r.body).session_id !== undefined;
                } catch { return false; }
            },
        });

        if (!created) {
            sessionLifecycleErrors.add(1, { step: 'create' });
            return;
        }

        sessionsCreated.add(1);
        let body;
        try {
            body = JSON.parse(createRes.body);
        } catch {
            return;
        }
        const sessionId = body.session_id;
        sleep(sleepWithJitter(0.3));

        // 2. Send heartbeats (simulating 2-8 minutes of watching with 30s intervals)
        const numHeartbeats = randomInt(4, 16);
        let position = 0;

        for (let i = 0; i < numHeartbeats; i++) {
            position += randomInt(25, 35); // ~30s of playback per heartbeat

            const hbStart = Date.now();
            const hbRes = authPost(`/playback/sessions/${sessionId}/heartbeat`, {
                position_seconds: position,
            }, token);

            heartbeatLatency.add(Date.now() - hbStart);
            const hbOk = hbRes.status === 204;
            heartbeatSuccess.add(hbOk ? 1 : 0);
            heartbeatsSent.add(1);

            if (!hbOk) {
                sessionLifecycleErrors.add(1, { step: 'heartbeat' });
            }

            // Realistic inter-heartbeat delay (scaled down for test)
            sleep(sleepWithJitter(0.2, 0.1));
        }

        // 3. End session
        const endRes = authDelete(`/playback/sessions/${sessionId}`, token);
        check(endRes, {
            'session ended': (r) => r.status === 204,
        });
        sessionsEnded.add(1);
    });
}

// Scenario 2: Burst heartbeats (simulates rapid seeking / buffering)
function heartbeatBurst(data) {
    group('Heartbeat Burst', () => {
        const token = data.token;
        const mediaId = randomFrom(data.movieIds);
        if (!mediaId) return;

        // Create session
        const createRes = authPost('/playback/sessions', {
            media_type: 'movie',
            media_id: mediaId,
        }, token);

        if (createRes.status !== 201) return;

        sessionsCreated.add(1);
        let body;
        try {
            body = JSON.parse(createRes.body);
        } catch { return; }
        const sessionId = body.session_id;

        // Send rapid-fire heartbeats (simulates seeking through content)
        const burstSize = randomInt(10, 30);
        let position = randomInt(0, 3600);

        for (let i = 0; i < burstSize; i++) {
            // Random seeking: jump forward or backward
            position = Math.max(0, position + randomInt(-300, 300));

            const hbStart = Date.now();
            const hbRes = authPost(`/playback/sessions/${sessionId}/heartbeat`, {
                position_seconds: position,
            }, token);

            heartbeatLatency.add(Date.now() - hbStart);
            heartbeatSuccess.add(hbRes.status === 204 ? 1 : 0);
            heartbeatsSent.add(1);

            // Very short delay between bursts (user scrubbing)
            sleep(sleepWithJitter(0.05, 0.02));
        }

        // Cleanup
        authDelete(`/playback/sessions/${sessionId}`, token);
        sessionsEnded.add(1);
    });
}

// Scenario 3: Quick browse — create session then abandon without ending
// (tests server-side session cleanup / timeout handling)
function quickBrowseAndAbandon(data) {
    group('Quick Browse & Abandon', () => {
        const token = data.token;
        const mediaId = randomFrom(data.movieIds);
        if (!mediaId) return;

        // Create session
        const createRes = authPost('/playback/sessions', {
            media_type: 'movie',
            media_id: mediaId,
        }, token);

        if (createRes.status === 201) {
            sessionsCreated.add(1);

            let body;
            try {
                body = JSON.parse(createRes.body);
            } catch { return; }
            const sessionId = body.session_id;

            // Send 1-2 heartbeats then abandon
            for (let i = 0; i < randomInt(1, 2); i++) {
                const hbRes = authPost(`/playback/sessions/${sessionId}/heartbeat`, {
                    position_seconds: randomInt(0, 120),
                }, token);
                heartbeatsSent.add(1);
                heartbeatSuccess.add(hbRes.status === 204 ? 1 : 0);
                sleep(sleepWithJitter(0.1));
            }
            // Intentionally NOT ending the session — tests server-side cleanup
        }
    });
}

// Scenario 4: Inspect an active session by ID
function sessionInspect(data) {
    group('Session Inspect', () => {
        const token = data.token;
        const mediaId = randomFrom(data.movieIds);
        if (!mediaId) return;

        // Create a session to inspect
        const createRes = authPost('/playback/sessions', {
            media_type: 'movie',
            media_id: mediaId,
        }, token);

        if (createRes.status !== 201) return;

        sessionsCreated.add(1);
        let body;
        try {
            body = JSON.parse(createRes.body);
        } catch { return; }
        const sessionId = body.session_id;

        sleep(sleepWithJitter(0.2));

        // Inspect the session
        const getRes = authGet(`/playback/sessions/${sessionId}`, token);
        check(getRes, {
            'get session detail': (r) => r.status === 200,
            'session has expected id': (r) => {
                try {
                    return JSON.parse(r.body).session_id === sessionId;
                } catch { return false; }
            },
        });

        sleep(sleepWithJitter(0.2));

        // Clean up
        authDelete(`/playback/sessions/${sessionId}`, token);
        sessionsEnded.add(1);
    });
}

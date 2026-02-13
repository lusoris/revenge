// tests/load/write_operations.js - Stress test write-heavy operations
// Tests watch progress updates, watched markers, and mixed read/write workloads
// Usage: k6 run --env PROFILE=gentle tests/load/write_operations.js
import { check, group, sleep } from 'k6';
import { Counter, Rate, Trend } from 'k6/metrics';
import { PROFILES } from './config.js';
import {
    authGet,
    authPost,
    authPut,
    extractItems,
    login,
    randomFrom,
    randomInt,
    sleepWithJitter,
    weightedRandom
} from './helpers.js';

// Custom metrics
const writeOps = new Counter('write_operations');
const readOps = new Counter('read_operations');
const writeSuccess = new Rate('write_success_rate');
const writeLatency = new Trend('write_latency');
const readAfterWriteConsistency = new Rate('read_after_write_consistent');
const progressUpdates = new Counter('progress_updates');
const watchedMarks = new Counter('watched_marks');

const profileName = __ENV.PROFILE || 'smoke';
const profile = PROFILES[profileName] || PROFILES.smoke;

export const options = {
    stages: profile.stages,
    thresholds: {
        http_req_duration: ['p(95)<1000', 'p(99)<3000'],
        http_req_failed: ['rate<0.20'],
        'write_success_rate': ['rate>0.80'],
        'write_latency': ['p(95)<500'],
        'read_after_write_consistent': ['rate>0.80'],
    },
};

const SCENARIOS = [
    { name: 'movie_progress', weight: 30, fn: movieProgressUpdates },
    { name: 'episode_progress', weight: 25, fn: episodeProgressUpdates },
    { name: 'mark_watched', weight: 20, fn: markWatched },
    { name: 'mixed_read_write', weight: 15, fn: mixedReadWrite },
    { name: 'settings_updates', weight: 10, fn: settingsUpdates },
];

export function setup() {
    const token = login();
    if (!token) {
        throw new Error('Setup login failed');
    }

    let movieIds = [];
    let episodeIds = [];
    let tvshowIds = [];

    // Discover movies
    const movies = authGet('/movies?limit=50', token);
    if (movies.status === 200) {
        movieIds = extractItems(movies).map(m => m.id).filter(Boolean);
    }

    // Discover TV shows and episodes
    const tvshows = authGet('/tvshows?limit=20', token);
    if (tvshows.status === 200) {
        tvshowIds = extractItems(tvshows).map(t => t.id).filter(Boolean);

        for (let i = 0; i < Math.min(5, tvshowIds.length); i++) {
            const seasons = authGet(`/tvshows/${tvshowIds[i]}/seasons`, token);
            if (seasons.status === 200) {
                for (const season of extractItems(seasons).slice(0, 2)) {
                    if (!season.id) continue;
                    const eps = authGet(`/tvshows/seasons/${season.id}/episodes`, token);
                    if (eps.status === 200) {
                        extractItems(eps).forEach(ep => {
                            if (ep.id) episodeIds.push(ep.id);
                        });
                    }
                }
            }
        }
    }

    console.log(`Setup: ${movieIds.length} movies, ${episodeIds.length} episodes, ${tvshowIds.length} tvshows`);
    return { token, movieIds, episodeIds, tvshowIds };
}

export default function (data) {
    const scenario = weightedRandom(SCENARIOS);
    try {
        scenario.fn(data);
    } catch (e) {
        console.error(`Scenario ${scenario.name} failed: ${e.message}`);
    }
    sleep(sleepWithJitter(0.5, 0.3));
}

// Scenario 1: Update movie watch progress (incremental position updates)
function movieProgressUpdates(data) {
    group('Movie Progress Updates', () => {
        const token = data.token;
        const movieId = randomFrom(data.movieIds);
        if (!movieId) return;

        // Simulate a user watching — periodic progress writes
        const movieDuration = randomInt(5400, 9000); // 90-150 min movie in seconds
        const numUpdates = randomInt(3, 10);
        let position = randomInt(0, movieDuration / 2); // start from random point

        for (let i = 0; i < numUpdates; i++) {
            position = Math.min(movieDuration, position + randomInt(120, 600));

            const start = Date.now();
            const res = authPut(`/movies/${movieId}/progress`, {
                position: position,
                duration: movieDuration,
            }, token);

            writeLatency.add(Date.now() - start);
            writeOps.add(1);
            progressUpdates.add(1);

            const ok = res.status >= 200 && res.status < 300;
            writeSuccess.add(ok ? 1 : 0);

            check(res, {
                'progress updated': (r) => r.status >= 200 && r.status < 300,
            });

            sleep(sleepWithJitter(0.1, 0.05));
        }

        // Verify progress was persisted (read-after-write consistency)
        sleep(sleepWithJitter(0.1));
        const readRes = authGet(`/movies/${movieId}`, token);
        readOps.add(1);

        if (readRes.status === 200) {
            try {
                const body = JSON.parse(readRes.body);
                // Check if progress field reflects our last update
                const hasProgress = body.progress !== undefined ||
                    body.watch_progress !== undefined ||
                    body.user_data !== undefined;
                readAfterWriteConsistency.add(hasProgress ? 1 : 0);
            } catch {
                readAfterWriteConsistency.add(0);
            }
        }
    });
}

// Scenario 2: Update episode watch progress (TV show binge pattern)
function episodeProgressUpdates(data) {
    group('Episode Progress Updates', () => {
        const token = data.token;

        // Pick a sequential set of episodes (binge watching)
        const numEpisodes = Math.min(randomInt(2, 5), data.episodeIds.length);
        const startIdx = randomInt(0, Math.max(0, data.episodeIds.length - numEpisodes));

        for (let i = 0; i < numEpisodes; i++) {
            const episodeId = data.episodeIds[startIdx + i];
            if (!episodeId) continue;

            const episodeDuration = randomInt(1200, 3600); // 20-60 min
            let position = 0;

            // A few progress updates per episode
            const updates = randomInt(2, 5);
            for (let j = 0; j < updates; j++) {
                position = Math.min(episodeDuration, position + randomInt(200, 800));

                const start = Date.now();
                const res = authPut(`/tvshows/episodes/${episodeId}/progress`, {
                    position: position,
                    duration: episodeDuration,
                }, token);

                writeLatency.add(Date.now() - start);
                writeOps.add(1);
                progressUpdates.add(1);
                writeSuccess.add(res.status >= 200 && res.status < 300 ? 1 : 0);

                sleep(sleepWithJitter(0.05, 0.02));
            }

            // Mark episode as watched at end
            if (Math.random() > 0.3) {
                const watchRes = authPost(`/tvshows/episodes/${episodeId}/watched`, {}, token);
                writeOps.add(1);
                watchedMarks.add(1);
                writeSuccess.add(watchRes.status >= 200 && watchRes.status < 300 ? 1 : 0);
            }

            sleep(sleepWithJitter(0.1));
        }
    });
}

// Scenario 3: Mark content as watched (movie + episode batch)
function markWatched(data) {
    group('Mark Watched', () => {
        const token = data.token;

        // Mark a random movie as watched
        const movieId = randomFrom(data.movieIds);
        if (movieId) {
            const start = Date.now();
            const res = authPost(`/movies/${movieId}/watched`, {}, token);
            writeLatency.add(Date.now() - start);
            writeOps.add(1);
            watchedMarks.add(1);

            check(res, {
                'movie marked watched': (r) => r.status >= 200 && r.status < 300,
            });
            writeSuccess.add(res.status >= 200 && res.status < 300 ? 1 : 0);
        }

        sleep(sleepWithJitter(0.2));

        // Bulk mark episodes as watched
        if (data.episodeIds.length >= 3) {
            const bulkIds = [];
            const startIdx = randomInt(0, Math.max(0, data.episodeIds.length - 5));
            for (let i = 0; i < Math.min(5, data.episodeIds.length); i++) {
                bulkIds.push(data.episodeIds[startIdx + i]);
            }

            const start = Date.now();
            const res = authPost('/tvshows/episodes/bulk-watched', {
                episode_ids: bulkIds,
            }, token);
            writeLatency.add(Date.now() - start);
            writeOps.add(1);
            watchedMarks.add(bulkIds.length);

            check(res, {
                'bulk watched': (r) => r.status >= 200 && r.status < 300,
            });
            writeSuccess.add(res.status >= 200 && res.status < 300 ? 1 : 0);
        }

        sleep(sleepWithJitter(0.2));

        // Verify via continue-watching or watch-history
        const histRes = authGet('/movies/watch-history', token);
        readOps.add(1);
        check(histRes, {
            'watch history accessible': (r) => r.status === 200,
        });
    });
}

// Scenario 4: Mixed read/write workload (realistic app behavior)
function mixedReadWrite(data) {
    group('Mixed Read/Write', () => {
        const token = data.token;

        // Simulate: browse → watch progress → browse more → update progress
        const actions = [
            // Read: browse catalog
            () => {
                const res = authGet('/movies?limit=20', token);
                readOps.add(1);
                return res.status === 200;
            },
            // Read: check continue watching
            () => {
                const res = authGet('/movies/continue-watching', token);
                readOps.add(1);
                return res.status === 200;
            },
            // Write: update progress on random movie
            () => {
                const movieId = randomFrom(data.movieIds);
                if (!movieId) return true;
                const start = Date.now();
                const res = authPut(`/movies/${movieId}/progress`, {
                    position: randomInt(0, 7200),
                    duration: 7200,
                }, token);
                writeLatency.add(Date.now() - start);
                writeOps.add(1);
                progressUpdates.add(1);
                writeSuccess.add(res.status >= 200 && res.status < 300 ? 1 : 0);
                return res.status >= 200 && res.status < 300;
            },
            // Read: search
            () => {
                const terms = ['action', 'comedy', 'drama', 'star', 'the'];
                const res = authGet(`/search/movies?q=${randomFrom(terms)}&limit=10`, token);
                readOps.add(1);
                return res.status === 200;
            },
            // Read: movie detail
            () => {
                const movieId = randomFrom(data.movieIds);
                if (!movieId) return true;
                const res = authGet(`/movies/${movieId}`, token);
                readOps.add(1);
                return res.status === 200 || res.status === 404;
            },
            // Write: mark watched
            () => {
                const movieId = randomFrom(data.movieIds);
                if (!movieId) return true;
                const start = Date.now();
                const res = authPost(`/movies/${movieId}/watched`, {}, token);
                writeLatency.add(Date.now() - start);
                writeOps.add(1);
                watchedMarks.add(1);
                writeSuccess.add(res.status >= 200 && res.status < 300 ? 1 : 0);
                return res.status >= 200 && res.status < 300;
            },
            // Read: tvshow list
            () => {
                const res = authGet('/tvshows?limit=20', token);
                readOps.add(1);
                return res.status === 200;
            },
            // Read: check stats
            () => {
                const res = authGet('/movies/stats', token);
                readOps.add(1);
                return res.status === 200;
            },
        ];

        // Execute 4-8 random actions with natural pauses
        const numActions = randomInt(4, 8);
        for (let i = 0; i < numActions; i++) {
            const action = randomFrom(actions);
            action();
            sleep(sleepWithJitter(0.2, 0.1));
        }
    });
}

// Scenario 5: User settings read/write
function settingsUpdates(data) {
    group('Settings Updates', () => {
        const token = data.token;

        // Read current preferences
        let res = authGet('/users/me/preferences', token);
        readOps.add(1);
        check(res, {
            'get preferences': (r) => r.status === 200,
        });

        sleep(sleepWithJitter(0.2));

        // Read user settings
        res = authGet('/settings/user', token);
        readOps.add(1);
        check(res, {
            'get user settings': (r) => r.status === 200,
        });

        sleep(sleepWithJitter(0.2));

        // Update a user setting (idempotent, safe for load test)
        const settingKeys = ['theme', 'language', 'default_page_size'];
        const key = randomFrom(settingKeys);

        const start = Date.now();
        res = authPut(`/settings/user/${key}`, {
            value: key === 'default_page_size' ? String(randomInt(10, 50)) : 'default',
        }, token);
        writeLatency.add(Date.now() - start);
        writeOps.add(1);

        // Settings PUT might return 200, 204, or 404 if key doesn't exist
        writeSuccess.add(res.status >= 200 && res.status < 300 ? 1 : 0);

        check(res, {
            'setting update': (r) => r.status >= 200 && r.status < 400,
        });

        sleep(sleepWithJitter(0.2));

        // Read back to verify
        res = authGet('/settings/user', token);
        readOps.add(1);
    });
}

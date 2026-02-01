# Documentation Completion TODO

> Last updated: 2026-02-01
> Status: 113 PLACEHOLDERs remaining / 146 total files (22% complete)

## How to Use This File

1. Check off items as you complete them: `- [x]`
2. Run `pytest tests/automation/test_yaml_validation.py::TestPlaceholderContent -v` to verify
3. Update thresholds in `tests/constants.py` as progress is made
4. Regenerate docs: `python scripts/automation/batch_regenerate.py`

---

## Priority 1: Architecture Foundation (5 files)

Core architecture docs - needed for understanding the entire system.

- [ ] `architecture/01_ARCHITECTURE.yaml` - System overview
- [ ] `architecture/02_DESIGN_PRINCIPLES.yaml` - Design philosophy
- [ ] `architecture/03_METADATA_SYSTEM.yaml` - Metadata handling
- [ ] `architecture/04_PLAYER_ARCHITECTURE.yaml` - Player subsystem
- [ ] `architecture/05_PLUGIN_ARCHITECTURE_DECISION.yaml` - Extension design

## Priority 2: Core Services (15 files)

Backend services that power all features.

- [ ] `services/AUTH.yaml` - Authentication service
- [ ] `services/USER.yaml` - User management
- [ ] `services/SESSION.yaml` - Session handling
- [ ] `services/RBAC.yaml` - Role-based access control
- [ ] `services/GRANTS.yaml` - Permission grants
- [ ] `services/APIKEYS.yaml` - API key management
- [ ] `services/OIDC.yaml` - OpenID Connect
- [ ] `services/LIBRARY.yaml` - Library management
- [ ] `services/METADATA.yaml` - Metadata service
- [ ] `services/SEARCH.yaml` - Search functionality
- [ ] `services/ACTIVITY.yaml` - Activity tracking
- [ ] `services/ANALYTICS.yaml` - Analytics
- [ ] `services/NOTIFICATION.yaml` - Notifications
- [ ] `services/SETTINGS.yaml` - User settings
- [ ] `services/FINGERPRINT.yaml` - Media fingerprinting

## Priority 3: Technical & Research (4 files)

Technical deep-dives and research docs.

- [ ] `technical/TECH_STACK.yaml` - Technology choices
- [ ] `research/USER_PAIN_POINTS_RESEARCH.yaml` - User research
- [ ] `research/UX_UI_RESOURCES.yaml` - UX/UI references

Meta docs (may not need wiki content):
- [ ] `01_DESIGN_DOC_TEMPLATE.yaml`
- [ ] `02_QUESTIONS_TO_DISCUSS.yaml`
- [ ] `03_DESIGN_DOCS_STATUS.yaml`

## Priority 4: Video Features (8 files)

Core media playback features.

- [ ] `features/video/MOVIE_MODULE.yaml`
- [ ] `features/video/TVSHOW_MODULE.yaml`
- [ ] `features/playback/WATCH_NEXT_CONTINUE_WATCHING.yaml`
- [ ] `features/playback/SKIP_INTRO.yaml`
- [ ] `features/playback/TRICKPLAY.yaml`
- [ ] `features/playback/SYNCPLAY.yaml`
- [ ] `features/playback/MEDIA_ENHANCEMENTS.yaml`
- [ ] `features/playback/RELEASE_CALENDAR.yaml`

## Priority 5: Shared Features (14 files)

Cross-cutting features used by multiple modules.

- [ ] `features/shared/LIBRARY_TYPES.yaml`
- [ ] `features/shared/ACCESS_CONTROLS.yaml`
- [ ] `features/shared/RBAC_CASBIN.yaml`
- [ ] `features/shared/CONTENT_RATING.yaml`
- [ ] `features/shared/NSFW_TOGGLE.yaml`
- [ ] `features/shared/I18N.yaml`
- [ ] `features/shared/SCROBBLING.yaml`
- [ ] `features/shared/REQUEST_SYSTEM.yaml`
- [ ] `features/shared/ANALYTICS_SERVICE.yaml`
- [ ] `features/shared/USER_EXPERIENCE_FEATURES.yaml`
- [ ] `features/shared/CLIENT_SUPPORT.yaml`
- [ ] `features/shared/WIKI_SYSTEM.yaml`
- [ ] `features/shared/NEWS_SYSTEM.yaml`
- [ ] `features/shared/TICKETING_SYSTEM.yaml`
- [ ] `features/shared/VOICE_CONTROL.yaml`

## Priority 6: Adult Content (6 files)

QAR module documentation.

- [ ] `features/adult/ADULT_CONTENT_SYSTEM.yaml`
- [ ] `features/adult/ADULT_METADATA.yaml`
- [ ] `features/adult/DATA_RECONCILIATION.yaml`
- [ ] `features/adult/GALLERY_MODULE.yaml`
- [ ] `features/adult/WHISPARR_STASHDB_SCHEMA.yaml`
- [ ] `features/livetv/LIVE_TV_DVR.yaml`

## Priority 7: Infrastructure Integrations (4 files)

Core infrastructure components.

- [ ] `integrations/infrastructure/POSTGRESQL.yaml`
- [ ] `integrations/infrastructure/DRAGONFLY.yaml`
- [ ] `integrations/infrastructure/TYPESENSE.yaml`
- [ ] `integrations/infrastructure/RIVER.yaml`

## Priority 8: Servarr Stack (5 files)

Media management integrations.

- [ ] `integrations/servarr/RADARR.yaml`
- [ ] `integrations/servarr/SONARR.yaml`
- [ ] `integrations/servarr/LIDARR.yaml`
- [ ] `integrations/servarr/WHISPARR.yaml`
- [ ] `integrations/servarr/CHAPTARR.yaml`

## Priority 9: Metadata Providers - Video (4 files)

- [ ] `integrations/metadata/video/TMDB.yaml`
- [ ] `integrations/metadata/video/THETVDB.yaml`
- [ ] `integrations/metadata/video/OMDB.yaml`
- [ ] `integrations/metadata/video/THEPOSTERDB.yaml`

## Priority 10: Metadata Providers - Music (4 files)

- [ ] `integrations/metadata/music/MUSICBRAINZ.yaml`
- [ ] `integrations/metadata/music/LASTFM.yaml`
- [ ] `integrations/metadata/music/SPOTIFY.yaml`
- [ ] `integrations/metadata/music/DISCOGS.yaml`

## Priority 11: Metadata Providers - Books (4 files)

- [ ] `integrations/metadata/books/OPENLIBRARY.yaml`
- [ ] `integrations/metadata/books/GOODREADS.yaml`
- [ ] `integrations/metadata/books/AUDIBLE.yaml`
- [ ] `integrations/metadata/books/HARDCOVER.yaml`

## Priority 12: Metadata Providers - Comics (3 files)

- [ ] `integrations/metadata/comics/COMICVINE.yaml`
- [ ] `integrations/metadata/comics/MARVEL_API.yaml`
- [ ] `integrations/metadata/comics/GRAND_COMICS_DATABASE.yaml`

## Priority 13: Scrobbling Services (5 files)

- [ ] `integrations/scrobbling/TRAKT.yaml`
- [ ] `integrations/scrobbling/LASTFM_SCROBBLE.yaml`
- [ ] `integrations/scrobbling/LISTENBRAINZ.yaml`
- [ ] `integrations/scrobbling/LETTERBOXD.yaml`
- [ ] `integrations/scrobbling/SIMKL.yaml`

## Priority 14: Anime Services (3 files)

- [ ] `integrations/anime/ANILIST.yaml`
- [ ] `integrations/anime/MYANIMELIST.yaml`
- [ ] `integrations/anime/KITSU.yaml`

## Priority 15: Auth Providers (4 files)

- [ ] `integrations/auth/AUTHENTIK.yaml`
- [ ] `integrations/auth/KEYCLOAK.yaml`
- [ ] `integrations/auth/AUTHELIA.yaml`
- [ ] `integrations/auth/GENERIC_OIDC.yaml`

## Priority 16: Casting & Live TV (5 files)

- [ ] `integrations/casting/CHROMECAST.yaml`
- [ ] `integrations/casting/DLNA.yaml`
- [ ] `integrations/livetv/TVHEADEND.yaml`
- [ ] `integrations/livetv/NEXTPVR.yaml`
- [ ] `integrations/livetv/ERSATZTV.yaml`

## Priority 17: Wiki Sources (6 files)

- [ ] `integrations/wiki/WIKIPEDIA.yaml`
- [ ] `integrations/wiki/FANDOM.yaml`
- [ ] `integrations/wiki/TVTROPES.yaml`
- [ ] `integrations/wiki/adult/IAFD.yaml`
- [ ] `integrations/wiki/adult/BABEPEDIA.yaml`
- [ ] `integrations/wiki/adult/BOOBPEDIA.yaml`

## Priority 18: Adult Metadata Providers (10 files)

- [ ] `integrations/metadata/adult/STASHDB.yaml`
- [ ] `integrations/metadata/adult/THEPORNDB.yaml`
- [ ] `integrations/metadata/adult/STASH.yaml`
- [ ] `integrations/metadata/adult/FREEONES.yaml`
- [ ] `integrations/metadata/adult/THENUDE.yaml`
- [ ] `integrations/metadata/adult/PORNHUB.yaml`
- [ ] `integrations/metadata/adult/TWITTER_X.yaml`
- [ ] `integrations/metadata/adult/INSTAGRAM.yaml`
- [ ] `integrations/metadata/adult/ONLYFANS.yaml`
- [ ] `integrations/metadata/adult/WHISPARR_V3_ANALYSIS.yaml`

## Priority 19: Transcoding (1 file)

- [ ] `integrations/transcoding/BLACKBEARD.yaml`

---

## Progress Tracking

Update these counts after each session:

| Date | Completed | Remaining | % Done |
|------|-----------|-----------|--------|
| 2026-02-01 | 33 | 113 | 22% |

---

## Quick Commands

```bash
# Run PLACEHOLDER test
pytest tests/automation/test_yaml_validation.py::TestPlaceholderContent -v

# Regenerate all docs
python scripts/automation/batch_regenerate.py

# Count remaining PLACEHOLDERs
grep -r "PLACEHOLDER" data/ --include="*.yaml" | wc -l
```

# Jellyfin Go - Verbesserungen gegenüber Original-Jellyfin

> Dieses Dokument beschreibt geplante Verbesserungen, die über das Original-Jellyfin hinausgehen.
> Inspiriert vom Rating-System-Pattern: Normalisierung, Domain-Trennung, saubere Relationen.

## Übersicht

| Verbesserung | Priorität | Aufwand | Status |
|--------------|-----------|---------|--------|
| [Genre-Domain-Trennung](#1-genre-domain-trennung) | P0 | Medium | ✅ DONE |
| [RBAC Permission System](#2-rbac-permission-system) | P1 | High | 🔴 TODO |
| [User Groups & Family Sharing](#3-user-groups--family-sharing) | P2 | Medium | 🔴 TODO |
| [Music Module Enhancement](#4-music-module-enhancement) | P1 | High | 🔴 TODO |
| [Library Scan Configuration](#5-library-scan-configuration) | P2 | Medium | 🔴 TODO |
| [Metadata Provider System](#6-metadata-provider-system) | P2 | Medium | 🔴 TODO |

---

## 1. Genre-Domain-Trennung

### ✅ IMPLEMENTED (Commit 6befd37c5)

### Problem im Original-Jellyfin
- Alle Genres in einem globalen Pool gespeichert (`TEXT[]` auf Items)
- Film-Genres erscheinen im Musik-Player (und umgekehrt)
- Hardcoded Workarounds im Code (`ExcludeItemTypes` in `Genre.GetTaggedItems()`)
- Keine Genre-Hierarchie (z.B. "Rock" → "Alternative Rock" → "Grunge")

### Implementierte Lösung
Domain-scoped Genres mit proper Relationen statt `TEXT[]`.

### Abgeschlossene Tasks

#### Database Migration `000015_genres` ✅
- [x] `genre_domain` ENUM: `movie`, `tv`, `music`, `book`, `podcast`, `game`
- [x] `genres` Tabelle mit Hierarchie-Support und External IDs
- [x] `media_item_genres` Junction-Table mit Source/Confidence
- [x] ~80 Seed-Genres für Movie, TV, Music Domains
- [x] Down-Migration

#### Domain Layer ✅
- [x] `internal/domain/genre.go`:
  - `GenreDomain` Type mit Konstanten
  - `Genre` Entity mit Parent/Children
  - `GenreRepository` Interface
  - `GenreService` Interface

#### Repository Layer ✅
- [x] `internal/infra/database/queries/genres.sql` (18+ Queries)
- [x] `internal/infra/database/repository/genre_repository.go`
- [x] `sqlc generate`

#### Service Layer ✅
- [x] `internal/service/genre/service.go`
- [x] `internal/service/genre/module.go` (fx)

#### API Layer ✅
- [x] `internal/api/handlers/genre.go`:
  - `GET /Genres` (mit `domain` Query-Param, `hierarchy` Mode)
  - `GET /Genres/{id}`
  - `GET /Genres/Search` (mit `domain`, `query`)
  - `POST /Genres` (Admin)
  - `PUT /Genres/{id}` (Admin)
  - `DELETE /Genres/{id}` (Admin)
  - `GET /Items/{itemId}/Genres`
  - `POST /Items/{itemId}/Genres/{genreId}` (Admin)
  - `DELETE /Items/{itemId}/Genres/{genreId}` (Admin)

### Offene Tasks für spätere Phase

#### Migration bestehender Daten
- [ ] Script: `TEXT[]` genres → `genres` + `media_item_genres`
- [ ] Domain-Erkennung basierend auf `media_items.type`
- [ ] `genres` Spalte auf `media_items` entfernen (nach Migration)

#### Tests (TODO)
- [ ] Repository Tests
- [ ] Service Tests
- [ ] Handler Tests
- [ ] Integration Test: Genre-Trennung zwischen Domains

---

## 2. RBAC Permission System

### Problem im Original-Jellyfin
- 30+ einzelne Boolean-Felder in `UserPolicy`:
  ```csharp
  public bool EnableMediaPlayback { get; set; }
  public bool EnableAudioPlaybackTranscoding { get; set; }
  public bool EnableVideoPlaybackTranscoding { get; set; }
  public bool EnableContentDeletion { get; set; }
  // ... 26 weitere
  ```
- Keine Rollen, keine Gruppen, keine Vererbung
- Jeder User wird individuell konfiguriert
- Kombinatorische Explosion bei neuen Features

### Lösung
Role-Based Access Control mit Permission-Vererbung und User-Overrides.

### Tasks

#### Database Migration `000015_user_permissions`
- [ ] `permission_categories` Tabelle:
  - `id UUID PRIMARY KEY`
  - `code VARCHAR(50) UNIQUE` (z.B. `playback`, `library`, `admin`)
  - `name VARCHAR(100)`
  - `description TEXT`
  - `sort_order INT`
- [ ] `permissions` Tabelle:
  - `id UUID PRIMARY KEY`
  - `category_id UUID REFERENCES permission_categories`
  - `code VARCHAR(100) UNIQUE` (z.B. `playback.video.transcode`)
  - `name VARCHAR(200)`
  - `description TEXT`
  - `default_value BOOLEAN`
- [ ] `roles` Tabelle:
  - `id UUID PRIMARY KEY`
  - `code VARCHAR(50) UNIQUE` (z.B. `admin`, `user`, `child`, `guest`)
  - `name VARCHAR(100)`
  - `description TEXT`
  - `is_system BOOLEAN` (built-in Rollen nicht löschbar)
- [ ] `role_permissions` Junction:
  - `role_id UUID`
  - `permission_id UUID`
  - `granted BOOLEAN`
- [ ] `user_roles` Junction:
  - `user_id UUID`
  - `role_id UUID`
  - `granted_at TIMESTAMPTZ`
  - `granted_by UUID`
- [ ] `user_permission_overrides` Tabelle:
  - `user_id UUID`
  - `permission_id UUID`
  - `granted BOOLEAN` (explizites Grant/Deny überschreibt Rolle)
- [ ] Seed-Daten:
  - Default Rollen: `admin`, `user`, `child`, `guest`
  - Default Permissions pro Kategorie
- [ ] Down-Migration

#### Domain Layer
- [ ] `internal/domain/permission.go`:
  - `Permission` Entity
  - `PermissionCategory` Entity
  - `Role` Entity
  - `UserPermission` (computed, mit Source: `role:admin` oder `override`)
  - `PermissionRepository` Interface
  - `PermissionService` Interface

#### Repository Layer
- [ ] `internal/infra/database/queries/permissions.sql`:
  - `GetPermissionByCode`
  - `ListPermissions`
  - `ListPermissionsByCategory`
  - `GetRole`
  - `ListRoles`
  - `GetRolePermissions`
  - `CreateRole`
  - `UpdateRole`
  - `DeleteRole`
  - `AssignRoleToUser`
  - `RemoveRoleFromUser`
  - `GetUserRoles`
  - `SetUserPermissionOverride`
  - `RemoveUserPermissionOverride`
  - `GetEffectiveUserPermissions` (computed)
- [ ] `internal/infra/database/repository/permission.go`

#### Service Layer
- [ ] `internal/service/permission/service.go`:
  - `HasPermission(ctx, userID, permissionCode) bool`
  - `GetEffectivePermissions(ctx, userID) []UserPermission`
  - `AssignRole(ctx, userID, roleCode) error`
  - `SetOverride(ctx, userID, permissionCode, granted) error`
- [ ] `internal/service/permission/module.go` (fx)

#### Middleware Update
- [ ] `internal/api/middleware/auth.go` erweitern:
  - `RequirePermission(permissionCode string) func(http.Handler) http.Handler`
  - Cache für User-Permissions (kurze TTL)

#### API Layer
- [ ] `internal/api/handlers/permission.go`:
  - `GET /Roles`
  - `GET /Roles/{id}`
  - `POST /Roles` (Admin)
  - `PUT /Roles/{id}` (Admin)
  - `DELETE /Roles/{id}` (Admin, nur non-system)
  - `GET /Users/{userId}/Permissions`
  - `GET /Users/{userId}/Roles`
  - `POST /Users/{userId}/Roles/{roleCode}` (Admin)
  - `DELETE /Users/{userId}/Roles/{roleCode}` (Admin)
  - `PUT /Users/{userId}/Permissions/{code}` (Admin, Override)
  - `DELETE /Users/{userId}/Permissions/{code}` (Admin, Remove Override)

#### Migration bestehender User
- [ ] `is_admin=true` → Rolle `admin` zuweisen
- [ ] `is_admin=false` → Rolle `user` zuweisen
- [ ] `is_admin` Spalte deprecaten (aber behalten für Kompatibilität)

#### Tests
- [ ] Permission Resolution Tests (Role + Override)
- [ ] Middleware Tests
- [ ] Service Tests
- [ ] Integration Test: Role-Zuweisung, Permission-Check

---

## 3. User Groups & Family Sharing

### Problem im Original-Jellyfin
- Keine Möglichkeit, mehrere User zu gruppieren
- Jede Library-Berechtigung einzeln pro User
- Kein Family-Sharing-Konzept
- Kein gemeinsamer Watch-History-Pool

### Lösung
User Groups mit shared Library Access und optionalem Watch-History-Sharing.

### Tasks

#### Database Migration `000019_user_groups`
- [ ] `user_groups` Tabelle:
  - `id UUID PRIMARY KEY`
  - `name VARCHAR(100)`
  - `owner_id UUID REFERENCES users`
  - `share_watch_history BOOLEAN DEFAULT false`
  - `created_at TIMESTAMPTZ`
- [ ] `user_group_members` Junction:
  - `group_id UUID`
  - `user_id UUID`
  - `can_manage BOOLEAN` (darf Mitglieder hinzufügen/entfernen)
  - `joined_at TIMESTAMPTZ`
- [ ] `library_access` Tabelle:
  - `id UUID PRIMARY KEY`
  - `library_id UUID REFERENCES libraries`
  - `user_id UUID` (NULL wenn group_id gesetzt)
  - `group_id UUID` (NULL wenn user_id gesetzt)
  - `access_level VARCHAR(20)`: `none`, `read`, `write`, `manage`
  - `CHECK` constraint: genau einer von user_id/group_id gesetzt
- [ ] Index auf `library_access(library_id, user_id)` und `(library_id, group_id)`

#### Domain Layer
- [ ] `internal/domain/group.go`:
  - `UserGroup` Entity
  - `GroupMember` Entity
  - `LibraryAccess` Entity
  - `LibraryAccessLevel` Type
  - `UserGroupRepository` Interface
  - `UserGroupService` Interface

#### Repository Layer
- [ ] `internal/infra/database/queries/groups.sql`:
  - `CreateGroup`
  - `GetGroup`
  - `ListUserGroups` (wo User Mitglied ist)
  - `AddMember`
  - `RemoveMember`
  - `ListGroupMembers`
  - `SetLibraryAccess`
  - `GetLibraryAccess` (für User, computed über Gruppen)
  - `ListAccessibleLibraries` (für User)

#### Service Layer
- [ ] `internal/service/group/service.go`
- [ ] Library Service erweitern: Access-Check über Groups

#### API Layer
- [ ] `internal/api/handlers/group.go`:
  - `GET /Groups` (eigene Gruppen)
  - `GET /Groups/{id}`
  - `POST /Groups`
  - `PUT /Groups/{id}`
  - `DELETE /Groups/{id}` (nur Owner)
  - `GET /Groups/{id}/Members`
  - `POST /Groups/{id}/Members/{userId}`
  - `DELETE /Groups/{id}/Members/{userId}`
  - `GET /Groups/{id}/Libraries` (Access-Konfiguration)
  - `PUT /Groups/{id}/Libraries/{libraryId}` (Set Access)

#### Tests
- [ ] Group CRUD Tests
- [ ] Library Access Resolution Tests
- [ ] Integration Test: Family-Szenario

---

## 4. Music Module Enhancement

### Problem im Original-Jellyfin
- Nur Basic MusicBrainz-Lookup
- Keine echten Artist-Entities (nur embedded auf Items)
- Keine Artist-Relationships (Bands, Mitglieder, Kollaborationen)
- Keine Lyrics-Integration
- Kein Audio-Fingerprinting (AcoustID)
- Multi-Disc-Alben werden schlecht gehandhabt
- Keine Music-spezifischen Playlists mit Smart-Features
- Kein Scrobbling/Play-History für Musik

### Lösung
Vollständiges Music-Domain mit echten Entities und Relationships.

### Tasks

#### Database Migration `000017_music_module`

##### Artists
- [ ] `artists` Tabelle:
  - `id UUID PRIMARY KEY`
  - `name VARCHAR(500)`
  - `sort_name VARCHAR(500)`
  - `disambiguation VARCHAR(500)` ("The Beatles (UK band)")
  - `biography TEXT`
  - `formed_date DATE`
  - `disbanded_date DATE`
  - `origin_country VARCHAR(2)` (ISO)
  - `artist_type VARCHAR(50)` (person, group, orchestra, choir)
  - `musicbrainz_id UUID UNIQUE`
  - `spotify_id VARCHAR(50)`
  - `lastfm_url VARCHAR(500)`
  - `discogs_id INT`
- [ ] `artist_relationships` Tabelle:
  - `artist_id UUID`
  - `related_artist_id UUID`
  - `relationship_type VARCHAR(50)` (member_of, collaboration, alias, tribute)
  - `start_date DATE`
  - `end_date DATE`

##### Albums
- [ ] `albums` Tabelle:
  - `id UUID PRIMARY KEY`
  - `library_id UUID REFERENCES libraries`
  - `title VARCHAR(500)`
  - `sort_title VARCHAR(500)`
  - `release_date DATE`
  - `release_type VARCHAR(50)` (album, ep, single, compilation, live, remix)
  - `total_tracks INT`
  - `total_discs INT DEFAULT 1`
  - `duration_ms BIGINT`
  - `musicbrainz_release_id UUID UNIQUE`
  - `musicbrainz_release_group_id UUID`
  - `spotify_id VARCHAR(50)`
  - `discogs_id INT`
  - `barcode VARCHAR(50)`
  - `original_release_date DATE`
  - `media_format VARCHAR(50)` (cd, vinyl, digital, cassette)
- [ ] `album_artists` Junction:
  - `album_id UUID`
  - `artist_id UUID`
  - `artist_type VARCHAR(50)` (primary, featured, remixer, producer)
  - `sort_order INT`

##### Tracks
- [ ] `tracks` Tabelle:
  - `id UUID PRIMARY KEY`
  - `album_id UUID REFERENCES albums` (nullable für Loose-Tracks)
  - `library_id UUID REFERENCES libraries`
  - `title VARCHAR(500)`
  - `sort_title VARCHAR(500)`
  - `duration_ms INT`
  - `track_number INT`
  - `disc_number INT DEFAULT 1`
  - `path TEXT NOT NULL`
  - `file_size BIGINT`
  - `audio_codec VARCHAR(50)` (flac, mp3, aac, opus)
  - `bitrate INT` (kbps)
  - `sample_rate INT` (Hz)
  - `bit_depth INT` (16, 24, 32)
  - `channels INT` (1, 2, 6)
  - `musicbrainz_recording_id UUID`
  - `musicbrainz_track_id UUID`
  - `isrc VARCHAR(20)`
  - `loudness_db DECIMAL(5,2)` (ReplayGain)
  - `bpm DECIMAL(5,2)`
  - `audio_fingerprint TEXT` (AcoustID)
  - `has_lyrics BOOLEAN DEFAULT false`
  - `lyrics_synced BOOLEAN DEFAULT false`
- [ ] `track_artists` Junction:
  - `track_id UUID`
  - `artist_id UUID`
  - `artist_type VARCHAR(50)` (primary, featured, remixer)
  - `sort_order INT`

##### Lyrics
- [ ] `lyrics` Tabelle:
  - `id UUID PRIMARY KEY`
  - `track_id UUID REFERENCES tracks`
  - `language VARCHAR(5)` (en, de, jp)
  - `is_synced BOOLEAN`
  - `content TEXT` (Plain oder LRC)
  - `source VARCHAR(100)` (embedded, lrclib, genius, manual)
  - `UNIQUE(track_id, language)`

##### Playlists
- [ ] `music_playlists` Tabelle:
  - `id UUID PRIMARY KEY`
  - `user_id UUID REFERENCES users`
  - `name VARCHAR(255)`
  - `description TEXT`
  - `is_public BOOLEAN`
  - `is_smart BOOLEAN`
  - `smart_rules JSONB`
  - `duration_ms BIGINT`
  - `track_count INT`
- [ ] `music_playlist_tracks` Junction:
  - `playlist_id UUID`
  - `track_id UUID`
  - `position INT`
  - `added_at TIMESTAMPTZ`

##### Play History / Scrobbling
- [ ] `music_play_history` Tabelle:
  - `id UUID PRIMARY KEY`
  - `user_id UUID`
  - `track_id UUID`
  - `played_at TIMESTAMPTZ`
  - `duration_played_ms INT`
  - `completed BOOLEAN` (>50% gespielt)
- [ ] Indices für schnelle Queries

#### Domain Layer
- [ ] `internal/domain/music.go`:
  - `Artist`, `ArtistRelationship`
  - `Album`, `AlbumArtist`
  - `Track`, `TrackArtist`
  - `Lyrics`
  - `MusicPlaylist`, `PlaylistTrack`
  - `PlayHistory`
  - Repository Interfaces
  - Service Interface

#### Repository Layer
- [ ] `internal/infra/database/queries/music_artists.sql`
- [ ] `internal/infra/database/queries/music_albums.sql`
- [ ] `internal/infra/database/queries/music_tracks.sql`
- [ ] `internal/infra/database/queries/music_playlists.sql`
- [ ] `internal/infra/database/queries/music_history.sql`
- [ ] Repository Implementierungen

#### Service Layer
- [ ] `internal/service/music/artist_service.go`
- [ ] `internal/service/music/album_service.go`
- [ ] `internal/service/music/track_service.go`
- [ ] `internal/service/music/playlist_service.go`
- [ ] `internal/service/music/history_service.go` (Scrobbling)
- [ ] `internal/service/music/module.go`

#### Metadata Providers
- [ ] `internal/service/music/providers/musicbrainz.go`
- [ ] `internal/service/music/providers/lastfm.go`
- [ ] `internal/service/music/providers/acoustid.go` (Fingerprinting)
- [ ] `internal/service/music/providers/lrclib.go` (Lyrics)

#### API Layer
- [ ] Artists:
  - `GET /Artists`
  - `GET /Artists/{id}`
  - `GET /Artists/{id}/Albums`
  - `GET /Artists/{id}/Tracks`
  - `GET /Artists/{id}/Similar`
  - `GET /Artists/{id}/Related` (Bands, Mitglieder)
- [ ] Albums:
  - `GET /Albums`
  - `GET /Albums/Recent`
  - `GET /Albums/{id}`
  - `GET /Albums/{id}/Tracks`
- [ ] Tracks:
  - `GET /Tracks/{id}`
  - `GET /Tracks/{id}/Lyrics`
  - `GET /Tracks/{id}/Similar`
  - `POST /Tracks/Identify` (Fingerprint-Lookup)
- [ ] Playlists:
  - `GET /Music/Playlists`
  - `POST /Music/Playlists`
  - `GET /Music/Playlists/{id}`
  - `PUT /Music/Playlists/{id}`
  - `DELETE /Music/Playlists/{id}`
  - `POST /Music/Playlists/{id}/Tracks`
  - `DELETE /Music/Playlists/{id}/Tracks/{trackId}`
  - `PUT /Music/Playlists/{id}/Tracks/Reorder`
- [ ] History:
  - `POST /Music/Scrobble` (Record play)
  - `GET /Music/History` (eigene)
  - `GET /Music/History/Recent`
  - `GET /Music/Stats` (Top Artists/Albums/Tracks)
- [ ] Radio/Mix:
  - `GET /Music/InstantMix/{trackId}`
  - `GET /Music/ArtistRadio/{artistId}`

#### Tests
- [ ] Artist CRUD Tests
- [ ] Album/Track Relationship Tests
- [ ] Playlist Tests
- [ ] Scrobbling Tests
- [ ] Integration Test: Full Music Library Scan

---

## 5. Library Scan Configuration

### Problem im Original-Jellyfin
- Schwache Content-Type-Enforcement
- Keine Per-Library Scan-Konfiguration
- Globale Einstellungen statt Library-spezifisch
- Kein feingranulares Pattern-Matching

### Lösung
Library-spezifische Scan-Konfiguration mit Patterns und Enforcement.

### Tasks

#### Database Migration `000018_library_scan_config`
- [ ] `library_scan_config` Tabelle:
  - `library_id UUID PRIMARY KEY REFERENCES libraries`
  - `include_patterns TEXT[]` (globs)
  - `exclude_patterns TEXT[]`
  - `use_folder_names BOOLEAN`
  - `enable_nfo_import BOOLEAN`
  - `prefer_embedded_metadata BOOLEAN`
  - `extract_chapters BOOLEAN`
  - `extract_thumbnails BOOLEAN`
  - `generate_intros BOOLEAN`
  - `enable_realtime_monitoring BOOLEAN`
  - `scan_on_startup BOOLEAN`
  - `enforce_content_type BOOLEAN`
  - `allowed_extensions TEXT[]`

#### Domain Layer
- [ ] `internal/domain/library.go` erweitern:
  - `LibraryScanConfig` Struct
  - `GetScanConfig`/`UpdateScanConfig` zum Interface

#### Repository Layer
- [ ] Queries für Scan-Config
- [ ] Repository-Methoden

#### Service Layer
- [ ] Library Service erweitern
- [ ] Scanner muss Config respektieren (wenn Scanner implementiert)

#### API Layer
- [ ] `GET /Library/VirtualFolders/{id}/ScanConfig`
- [ ] `PUT /Library/VirtualFolders/{id}/ScanConfig`

---

## 6. Metadata Provider System

### Problem im Original-Jellyfin
- Globale Provider-Konfiguration
- Keine Per-Library Provider-Auswahl
- Provider-Reihenfolge nicht flexibel

### Lösung
Provider-System mit Library-spezifischer Konfiguration.

### Tasks

#### Database Migration `000020_metadata_providers`
- [ ] `metadata_providers` Tabelle:
  - `id UUID PRIMARY KEY`
  - `code VARCHAR(50) UNIQUE` (tmdb, tvdb, musicbrainz, etc.)
  - `name VARCHAR(100)`
  - `provider_type VARCHAR(50)` (movie, tv, music, multi)
  - `is_enabled BOOLEAN`
  - `config JSONB`
  - `priority INT`
- [ ] `library_metadata_providers` Junction:
  - `library_id UUID`
  - `provider_id UUID`
  - `priority INT` (lower = higher priority)
  - `enabled BOOLEAN`
  - `config_override JSONB`
- [ ] Seed-Daten für Standard-Provider

#### Domain Layer
- [ ] `internal/domain/provider.go`:
  - `MetadataProvider` Entity
  - `LibraryProviderConfig` Entity
  - Interfaces

#### Service Layer
- [ ] Provider-Resolution nach Priority
- [ ] Fallback-Chain

#### API Layer
- [ ] `GET /Metadata/Providers`
- [ ] `GET /Library/VirtualFolders/{id}/Providers`
- [ ] `PUT /Library/VirtualFolders/{id}/Providers`

---

## Implementierungsreihenfolge

### Phase 1: Foundation Fixes (P0)
1. **Genre-Domain-Trennung** - Behebt das UX-Problem sofort
2. **Genre-Migration** - Bestehende Daten überführen

### Phase 2: User Management (P1)
3. **RBAC Permissions** - Foundation für alles Weitere
4. **Permission Middleware** - API absichern

### Phase 3: Content Enhancement (P1)
5. **Music Module** - Major Feature Gap schließen
6. **Music Metadata Providers** - MusicBrainz, Last.fm, AcoustID

### Phase 4: Configuration & Polish (P2)
7. **User Groups** - Family Sharing
8. **Library Scan Config** - Feintuning
9. **Metadata Provider System** - Flexibilität

---

## Notizen

### API-Kompatibilität
- Neue Endpoints parallel zu existierenden
- Alte Endpoints (Jellyfin-kompatibel) behalten
- Neue Features über `/v2/` Prefix oder Feature-Detection

### Migration Strategy
- Alle Migrations in Reihenfolge ausführbar
- Down-Migrations für Rollback
- Daten-Migrations als separate Schritte

### Testing Requirements
- Unit Tests für Domain/Service Layer
- Integration Tests für Repository
- E2E Tests für kritische Flows
- Performance Tests für Genre/Permission Resolution

---

## Referenzen

- Original Jellyfin Code: `Jellyfin.*`, `MediaBrowser.*`, `Emby.*` im Repo
- Rating System als Pattern-Vorlage: `internal/domain/rating.go`, `internal/service/rating/`
- MusicBrainz API: https://musicbrainz.org/doc/MusicBrainz_API
- AcoustID: https://acoustid.org/webservice
- LrcLib: https://lrclib.net/docs

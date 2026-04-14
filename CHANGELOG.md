# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
Each module is versioned independently following [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- **Repository**: Adopted golangci-lint v2 with 30+ linters as the standards baseline (goenvoy 2.0). The pure-stdlib gate is now enforced by `depguard`. Per-module `replace` directives have been removed in favour of versioned `arr` requires (ADR-0009).
- **Tooling**: Bulk-applied `t.Parallel()` to every test function and subtest; wrapped third-party errors with module-prefixed `fmt.Errorf("%s: %s: %w", ...)`; canonicalised HTTP header literals.

### Breaking — module path migration to `/v2`
The 8 modules below ship a major bump and require a new import path (`…/v2`).
Sister `arr/*` modules now require `arr/v2 v2.0.0` (8 modules: sonarr, radarr, lidarr, readarr, prowlarr, bazarr, whisparr, seerr).

| Module | Old path | New path |
|---|---|---|
| `arr` | `github.com/golusoris/goenvoy/arr` | `…/arr/v2` |
| `arr/mylar` | `…/arr/mylar` | `…/arr/mylar/v2` |
| `arr/prowlarr` | `…/arr/prowlarr` | `…/arr/prowlarr/v2` |
| `mediaserver/emby` | `…/mediaserver/emby` | `…/mediaserver/emby/v2` |
| `mediaserver/jellyfin` | `…/mediaserver/jellyfin` | `…/mediaserver/jellyfin/v2` |
| `mediaserver/tdarr` | `…/mediaserver/tdarr` | `…/mediaserver/tdarr/v2` |
| `metadata/book/googlebooks` | `…/metadata/book/googlebooks` | `…/metadata/book/googlebooks/v2` |
| `metadata/music/audiodb` | `…/metadata/music/audiodb` | `…/metadata/music/audiodb/v2` |

API breaks bundled with the v2 cut:
- **arr** (`arr/v2`): `DevelopmentConfigResource.ApiKey` → `APIKey`. JSON wire tag (`apiKey`) unchanged.
- **arr/mylar** (`arr/mylar/v2`): All DTO fields renamed to Go-idiomatic acronyms — `Comic.Id` → `ID`, `Issue.Id`/`Issue.ComicId` → `ID`/`ComicID` (and similarly for `Upcoming`, `WantedIssue`, `HistoryEntry`, `SearchResult`, `StoryArc`, `ReadList`, `Provider`). JSON wire tags unchanged.
- **arr/prowlarr** (`arr/prowlarr/v2`): `DevelopmentConfigResource.LogSql` → `LogSQL`. JSON wire tag (`logSql`) unchanged.
- **mediaserver/emby** (`mediaserver/emby/v2`): `ServerId` → `ServerID` on `SystemInfo`, `SessionInfo`, `User`. JSON wire tags unchanged.
- **mediaserver/jellyfin** (`mediaserver/jellyfin/v2`): `ServerId` → `ServerID` on the same three types. JSON wire tags unchanged.
- **mediaserver/tdarr** (`mediaserver/tdarr/v2`): `Node.Id`/`Worker.Id`/`DBFile.Id` → `ID`; `DBFile.LibraryId`/`ScanFilesData.LibraryId` → `LibraryID`. Method parameter renames `libraryId`/`nodeId`/`workerId` → `libraryID`/`nodeID`/`workerID` on `GetResStats`, `GetDBStatuses`, `ScanFiles`, `CancelWorkerItem`, `KillWorker` (call-site compatible). JSON wire tags unchanged.
- **metadata/book/googlebooks** (`metadata/book/googlebooks/v2`): `Volume.Id` → `ID`. JSON wire tag (`id`) unchanged.
- **metadata/music/audiodb** (`metadata/music/audiodb/v2`): `IdAlbum`/`IdArtist`/`IdTrack` → `IDAlbum`/`IDArtist`/`IDTrack` across `Artist`, `Album`, `Track`, `MusicVideo`. JSON wire tags unchanged.

## [v1.3.0] - 2026-04-12

### Added
- **MobyGames** (`metadata/game/mobygames`): Game metadata client — genres, platforms, groups, search, game details, screenshots, covers, recent and random games.
- **SteamGridDB** (`metadata/game/steamgriddb`): Custom game artwork client (grids, heroes, logos, icons) with bearer token auth.
- **RetroAchievements** (`metadata/game/retroachievements`): Achievement and ROM hash data — game details, extended info, hashes, console IDs.
- **ScreenScraper** (`metadata/game/screenscraper`): Comprehensive game metadata and media — game info, search, systems, genres, user info, infrastructure info. Supports both dev credentials and optional per-user auth via `WithUser`.
- **Hasheous** (`metadata/game/hasheous`): Hash-based ROM identification — lookup by MD5, SHA1, SHA256, CRC, or multi-hash POST; platform listing.
- **LaunchBox** (`metadata/game/launchbox`): LaunchBox XML database client — downloads and parses `Metadata.zip` into an in-memory store with game search, alternate names, images, and platform lookup.

### Fixed
- **arr** (`arr`): `NewBaseClient` now uses a cloned `http.DefaultTransport` instead of accepting the nil-transport default. This prevents `httptest.Server.Close()` — which calls `CloseIdleConnections` on the shared default transport — from breaking in-flight requests in other parallel tests.

## [v1.2.0] - 2026-04-10

### Added
- **BaseClient** (`metadata`): Shared HTTP client (`BaseClient`) with `DoRaw`, `DoRawURL`, `DoJSON`, `Get`, auth injection, and functional options — eliminates ~2,480 lines of duplicated HTTP boilerplate across all 27 metadata providers.
- **Letterboxd** (`metadata/video/letterboxd`): Social film discovery client with OAuth2 Bearer auth — 65 methods covering films, film collections, contributors, lists, log entries, members, comments, stories, search, news, and auth helpers — with 60 tests.

### Changed
- All 27 metadata providers now embed `*metadata.BaseClient` instead of maintaining independent HTTP plumbing.
- **Restructured** `metadata/movie/` → `metadata/video/` — all 7 video provider modules (TMDb, TheTVDB, Fanart.tv, OMDb, TVmaze, Letterboxd, OpenSubtitles) moved to new import paths under `metadata/video/`.

## [v1.1.0] - 2026-04-09

### Added
- **Arr** (`arr/*`): `HeadPing`, `UploadBackup`, and `GetRaw` methods added to all 13 *arr packages.
- **Trakt** (`metadata/tracking/trakt`): 140+ new methods — comments, notes, calendars, sync, lists, social, scrobble, users, people, certifications, countries, genres, languages, networks — with 171 tests.
- **Google Books** (`metadata/book/googlebooks`): 18 new methods — bookshelves (list, get), volumes in shelf, annotations (list, insert, delete, update, summary), user library (add, remove, clear, mark reading), series (get, membership) — with 22 tests.
- **Discogs** (`metadata/music/discogs`): 50+ new methods — release ratings, user identity/profile/submissions/contributions, user collection (folders, items, add/remove/rate), wantlist (add/remove), user lists, marketplace (listings, orders, fee, stats, price suggestions, inventory export) — with 59 tests.
- **TMDb** (`metadata/video/tmdb`): 80+ new methods — movie/TV extras (credits, images, videos, reviews, similar, recommendations, keywords, providers), TV seasons/episodes, person details, search (multi, keyword, company, collection), collections, account (lists, favorites, watchlist, ratings), lists (CRUD), certifications, watch providers, companies, keywords, changes, reviews, networks — with 121 tests.
- **TVDB** (`metadata/video/tvdb`): 70+ new methods — artwork (statuses, types), awards (categories), characters, companies (types), content ratings, countries, entity types, episodes, genders, genres, inspiration types, lists (extended, translations), movies (filter, slug, statuses, extended), people (types, extended, translations), search, seasons (types, extended, translations), series (filter, slug, statuses, episodes by language, extended, translations), source types, updates, user (info, favorites) — with 71 tests.
- **Simkl** (`metadata/tracking/simkl`): 23 new methods — ratings (add/remove), scrobble (start/pause/stop/checkin), sync (history, ratings, add-to-list, remove, watched), users (stats, recently watched), movie genres, random search, best filters — with 56 total tests.

### Changed
- **TMDb** (`metadata/video/tmdb`): Renamed `TMDbAvatar` → `UserAvatar` to avoid type-name stuttering.

### Fixed
- Lint fixes across 6 packages: godot (section comment formatting), gofmt, revive (comment format, stuttering names), gocritic (parameter type combining), unparam (constant parameters).

## [v1.0.0] - 2026-03-31

### Changed
- **Whisparr** (`arr/whisparr`): Renamed `ErosClient` → `ClientV3`, `NewEros` → `NewV3`, `ErosHistoryRecord` → `HistoryRecordV3`, `ErosParseResult` → `ParseResultV3`. Removes internal codename from public API.
- **TMDb** (`metadata/video/tmdb`): `DiscoverMovies` and `DiscoverTV` now accept `url.Values` instead of raw `string` for the `extraParams` parameter.

### Fixed
- **arr** (`arr`): `BaseClient.Delete` now accepts a request body, fixing 9 bulk-delete and editor-delete operations across Lidarr, Radarr, Sonarr, Readarr, and Whisparr v3 that silently sent empty DELETE requests.
- **Seerr** (`arr/seerr`): URL-encode `filter` parameter to prevent query string injection.
- **Mylar** (`arr/mylar`): Use `url.Values` for API key and command parameters instead of raw string concatenation.
- **Jackett** (`arr/jackett`): URL-escape indexer ID in Torznab search path to prevent path traversal.
- **Steam** (`metadata/game/steam`): URL-encode API key in query parameters.
- **TMDb** (`metadata/video/tmdb`): URL-escape `mediaType` and `timeWindow` path segments in `GetTrending` to prevent path traversal.
- **Kavita** (`mediaserver/kavita`): Add `sync.RWMutex` to protect JWT token from data races under concurrent use.
- 15 modules now use a 30-second HTTP client timeout instead of `http.DefaultClient` with no timeout: Navidrome, Komga, Kavita, Steam, RAWG, IGDB, OpenSubtitles, Open Library, Google Books, Last.fm, Spotify, ListenBrainz, Discogs, TheAudioDB, Deezer.

## [v0.1.0] - 2026-03-31

### Added
- **Jackett** (`arr/jackett`): Torznab/Newznab proxy client — search all/specific indexers, TV/movie/music/book search, capabilities, indexer management, server config.
- **NZBHydra2** (`arr/nzbhydra`): Meta NZB indexer client — Newznab search, TV/movie/book search, capabilities, statistics, search/download history, indexer statuses.
- **Spotify** (`metadata/music/spotify`): Music metadata via OAuth2 Bearer — search, artists, albums, tracks, audio features, new releases, categories, recommendations.
- **Deezer** (`metadata/music/deezer`): Music metadata (no auth) — search tracks/albums/artists, artist top tracks/albums/related, album tracks, genres, charts.
- **ListenBrainz** (`metadata/music/listenbrainz`): Listening data and statistics — submit listens, user listens/history, top artists/releases/recordings, listening activity, similar users.
- **IGDB** (`metadata/game/igdb`): Game metadata via Twitch OAuth2 — search games/companies, game details, platforms, genres, covers, screenshots, popular games.
- **RAWG** (`metadata/game/rawg`): Video game database — search games, game details/screenshots/trailers/DLC/series, platforms, genres, publishers, developers, tags, stores.
- **Steam** (`metadata/game/steam`): Steam Store and Web API — app details, featured games, app list, current players, app news, global achievements.
- New `metadata/game` parent module with shared game metadata types.

### Added
- **Tautulli** (`mediaserver/tautulli`): Plex monitoring client — activity, history, libraries, users, notifications, geo lookup.
- **Autobrr** (`arr/autobrr`): Automation client — filters, indexers, IRC networks, feeds, releases, logs, config.
- **Audiobookshelf** (`mediaserver/audiobookshelf`): Audiobook/podcast server client — libraries, items, users, sessions, search, collections.
- **Komga** (`mediaserver/komga`): Comic/manga server client with Basic Auth — libraries, series, books, collections, read lists, users.
- **Navidrome** (`mediaserver/navidrome`): Music server client using Subsonic/OpenSubsonic API — artists, albums, songs, playlists, search, scrobbling, starring.
- **OpenSubtitles** (`metadata/video/opensubtitles`): Subtitle search and download — search, features, languages, formats, user info, popular, latest.
- **Last.fm** (`metadata/music/lastfm`): Music metadata — artist/album/track info, similar artists, charts, tags, search.
- **Discogs** (`metadata/music/discogs`): Music database — releases, artists, labels, master releases, search.
- **Kavita** (`mediaserver/kavita`): Manga/comic/ebook reader with JWT auth — libraries, series, volumes, chapters, collections, reading lists, search.
- **Tdarr** (`mediaserver/tdarr`): Media transcoding server — status, nodes, workers, file search, resolution stats, scan management.
- **Mylar3** (`arr/mylar`): Comic book automation — comics, issues, wanted, history, story arcs, reading lists, providers.
- **FlareSolverr** (`arr/flaresolverr`): Cloudflare bypass proxy — GET/POST requests, session management.
- **TheAudioDB** (`metadata/music/audiodb`): Music metadata — artist/album/track search and lookup, music videos, discography, charts, trending.
- **Open Library** (`metadata/book/openlibrary`): Book metadata — search, works, editions, authors, subjects, ISBN lookup.
- **Google Books** (`metadata/book/googlebooks`): Book search and volume details — search with filters, volume retrieval.
- OAuth2 support for Trakt (device code, auth code, refresh, revoke).
- OAuth2 support for Simkl (device PIN, auth code exchange).
- OAuth2/PKCE support for MAL (authorization URL, code exchange, refresh).
- OAuth2 password grant and refresh for Kitsu.
- Bearer token support for AniList, Kitsu, MAL, Trakt, Simkl.
- TVDB automatic token refresh on 401 Unauthorized.
- Jellyfin MediaBrowser Authorization header (replaces X-Emby-* headers).
- 27 new OAuth2 tests across 6 modules.
- Shared types tests for anime, downloadclient, mediaserver, metadata parent packages.
- `go.work` for local multi-module development.
- `CONTRIBUTING.md`, `SECURITY.md`, `CHANGELOG.md`.
- GitHub Actions CI (test + lint matrix across all modules).
- GitHub Actions release workflow for tagged modules.
- `Makefile` with test-all, lint-all, vet-all, tidy-all, build-all, fmt-all targets.
- `.golangci.yml` with comprehensive linter configuration.

### Fixed
- SABnzbd `PauseItem`/`ResumeItem`/`SetSpeedLimit` missing `name` parameter.
- SABnzbd `ServerStats` fields changed from `string` to `int64` (bytes).
- Seerr `Search` URL-encodes query parameter.
- Trakt `errorAs` uses `errors.As` instead of manual type assertion.
- TVDB `doGet` formatting and `errors.As` for 401 detection.
- All example_test.go files updated to match actual API signatures.
- Lint fixes: `errors.New` for static strings, US English spelling.
- README corrected: removed unlisted services, added MusicBrainz/StashBox/TPDB.

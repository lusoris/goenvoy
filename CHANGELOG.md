# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
Each module is versioned independently following [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Letterboxd** (`metadata/movie/letterboxd`): Social film discovery client with OAuth2 Bearer auth — 65 methods covering films, film collections, contributors, lists, log entries, members, comments, stories, search, news, and auth helpers — with 60 tests.

## [v1.1.0] - 2026-04-09

### Added
- **Arr** (`arr/*`): `HeadPing`, `UploadBackup`, and `GetRaw` methods added to all 13 *arr packages.
- **Trakt** (`metadata/tracking/trakt`): 140+ new methods — comments, notes, calendars, sync, lists, social, scrobble, users, people, certifications, countries, genres, languages, networks — with 171 tests.
- **Google Books** (`metadata/book/googlebooks`): 18 new methods — bookshelves (list, get), volumes in shelf, annotations (list, insert, delete, update, summary), user library (add, remove, clear, mark reading), series (get, membership) — with 22 tests.
- **Discogs** (`metadata/music/discogs`): 50+ new methods — release ratings, user identity/profile/submissions/contributions, user collection (folders, items, add/remove/rate), wantlist (add/remove), user lists, marketplace (listings, orders, fee, stats, price suggestions, inventory export) — with 59 tests.
- **TMDb** (`metadata/movie/tmdb`): 80+ new methods — movie/TV extras (credits, images, videos, reviews, similar, recommendations, keywords, providers), TV seasons/episodes, person details, search (multi, keyword, company, collection), collections, account (lists, favorites, watchlist, ratings), lists (CRUD), certifications, watch providers, companies, keywords, changes, reviews, networks — with 121 tests.
- **TVDB** (`metadata/movie/tvdb`): 70+ new methods — artwork (statuses, types), awards (categories), characters, companies (types), content ratings, countries, entity types, episodes, genders, genres, inspiration types, lists (extended, translations), movies (filter, slug, statuses, extended), people (types, extended, translations), search, seasons (types, extended, translations), series (filter, slug, statuses, episodes by language, extended, translations), source types, updates, user (info, favorites) — with 71 tests.
- **Simkl** (`metadata/tracking/simkl`): 23 new methods — ratings (add/remove), scrobble (start/pause/stop/checkin), sync (history, ratings, add-to-list, remove, watched), users (stats, recently watched), movie genres, random search, best filters — with 56 total tests.

### Changed
- **TMDb** (`metadata/movie/tmdb`): Renamed `TMDbAvatar` → `UserAvatar` to avoid type-name stuttering.

### Fixed
- Lint fixes across 6 packages: godot (section comment formatting), gofmt, revive (comment format, stuttering names), gocritic (parameter type combining), unparam (constant parameters).

## [v1.0.0] - 2026-03-31

### Changed
- **Whisparr** (`arr/whisparr`): Renamed `ErosClient` → `ClientV3`, `NewEros` → `NewV3`, `ErosHistoryRecord` → `HistoryRecordV3`, `ErosParseResult` → `ParseResultV3`. Removes internal codename from public API.
- **TMDb** (`metadata/movie/tmdb`): `DiscoverMovies` and `DiscoverTV` now accept `url.Values` instead of raw `string` for the `extraParams` parameter.

### Fixed
- **arr** (`arr`): `BaseClient.Delete` now accepts a request body, fixing 9 bulk-delete and editor-delete operations across Lidarr, Radarr, Sonarr, Readarr, and Whisparr v3 that silently sent empty DELETE requests.
- **Seerr** (`arr/seerr`): URL-encode `filter` parameter to prevent query string injection.
- **Mylar** (`arr/mylar`): Use `url.Values` for API key and command parameters instead of raw string concatenation.
- **Jackett** (`arr/jackett`): URL-escape indexer ID in Torznab search path to prevent path traversal.
- **Steam** (`metadata/game/steam`): URL-encode API key in query parameters.
- **TMDb** (`metadata/movie/tmdb`): URL-escape `mediaType` and `timeWindow` path segments in `GetTrending` to prevent path traversal.
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
- **OpenSubtitles** (`metadata/movie/opensubtitles`): Subtitle search and download — search, features, languages, formats, user info, popular, latest.
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

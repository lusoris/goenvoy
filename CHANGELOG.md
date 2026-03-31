# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
Each module is versioned independently following [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
- **Gotify** (`notification/gotify`): Push notification server — messages, applications, clients, users, health.
- **TheAudioDB** (`metadata/music/audiodb`): Music metadata — artist/album/track search and lookup, music videos, discography, charts, trending.
- **Open Library** (`metadata/book/openlibrary`): Book metadata — search, works, editions, authors, subjects, ISBN lookup.
- **Google Books** (`metadata/book/googlebooks`): Book search and volume details — search with filters, volume retrieval.
- New `notification/` category with shared notification types.
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

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).
Each module is versioned independently following [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
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

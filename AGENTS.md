# Agent guide — goenvoy

> Cross-tool context for Claude Code, Cursor, Aider, Codex, Continue, and other coding assistants.
> **Read this before suggesting changes.** Then read the per-category `AGENTS.md` for the area you're touching.

## What this repo is

`goenvoy` is a **multi-module monorepo** of pure-stdlib Go HTTP-API clients — 69 modules covering 63+ services across arr-stack, metadata, download clients, media servers, and anime. Each service is its own Go module at `github.com/golusoris/goenvoy/<category>/<service>`, independently versioned via `<category>/<service>/vX.Y.Z` tags.

## Hard rules

1. **Never add a non-stdlib import.** Pure stdlib is load-bearing (ADR-0001). `net/http`, `encoding/json`, `encoding/xml`, `crypto/*`, `context`, `net/url`, `net/http/httptest` only. `depguard` in `.golangci.yml` enforces. If you *think* you need a dependency, write an ADR first.
2. **Every module exports `New(baseURL, apiKey string, opts ...Option) (*Client, error)`.** Exceptions go in the module's `AGENTS.md` (OAuth2 device-code, public-key-only APIs).
3. **Every I/O method takes `context.Context` as the first argument.** No exceptions.
4. **Every module defines an `APIError` struct implementing `error`.** Status code, message, response body.
5. **Every HTTP call closes the response body.** `bodyclose` enforces.
6. **All errors wrap** — `fmt.Errorf("<module>: <op>: %w", err)`. `wrapcheck` enforces at boundaries.
7. **Never log secrets.** API keys / bearer tokens / refresh tokens must not appear in error messages, wrapped errors, or test fixtures.
8. **Every merged commit: 0 lint · 0 gosec · 0 govulncheck · race-green** — **per module**. `//nolint` requires a same-line justification comment.
9. **Breaking a public API requires a `Migration:` footer** in the commit body with before/after Go snippets. CI `apidiff` runs against the last tag.

See [.workingdir/PRINCIPLES.md](.workingdir/PRINCIPLES.md) (the adapted §2 contract) for the full coding / security / testing rules.

## Repository layout

```
goenvoy/
├── arr/              # shared *arr types + base client
│   ├── sonarr/       # Sonarr v3 client
│   ├── radarr/       # Radarr v3 client
│   ├── lidarr/       # Lidarr client
│   ├── readarr/      # Readarr client
│   ├── whisparr/     # Whisparr v2 + v3 (Eros) client
│   ├── prowlarr/     # Prowlarr client
│   ├── bazarr/       # Bazarr client
│   ├── seerr/        # Seerr client
│   ├── autobrr/      # Autobrr client
│   ├── mylar/        # Mylar3 client
│   ├── flaresolverr/ # FlareSolverr client
│   ├── jackett/      # Jackett client
│   └── nzbhydra/     # NZBHydra2 client
├── metadata/         # shared metadata types (Rating, Image, Person, ...)
│   ├── video/        # TMDb, TheTVDB, Fanart.tv, OMDb, TVmaze, Letterboxd, OpenSubtitles
│   ├── anime/        # AniList, Kitsu, AniDB, MAL
│   ├── music/        # MusicBrainz, Last.fm, Discogs, TheAudioDB, Spotify, Deezer, ListenBrainz
│   ├── tracking/     # Trakt, Simkl
│   ├── adult/        # StashBox, TPDB
│   ├── book/         # Google Books, Open Library
│   └── game/         # IGDB, RAWG, Steam, MobyGames, SteamGridDB, RetroAchievements, ScreenScraper, Hasheous, LaunchBox
├── downloadclient/   # shared download types + Downloader interface
│   ├── qbit/         # qBittorrent WebUI client
│   ├── transmission/ # Transmission RPC client
│   ├── deluge/       # Deluge JSON-RPC client
│   ├── rtorrent/     # rTorrent XMLRPC client
│   ├── sabnzbd/      # SABnzbd client
│   └── nzbget/       # NZBGet JSON-RPC client
├── mediaserver/      # shared media server types
│   ├── plex/         # Plex Media Server client
│   ├── jellyfin/     # Jellyfin client
│   ├── emby/         # Emby client
│   ├── tautulli/     # Tautulli client
│   ├── audiobookshelf/ # Audiobookshelf client
│   ├── komga/        # Komga client
│   ├── navidrome/    # Navidrome client
│   ├── kavita/       # Kavita client
│   ├── stash/        # StashApp GraphQL client
│   └── tdarr/        # Tdarr client
└── anime/            # shared anime types
    └── shoko/        # Shoko Server client
```

## Common tasks

| Task | Claude Code skill / command |
|---|---|
| Add a new service client | `/add-service-client <category> <service>` |
| Add a method to an existing client | `/add-service-method <module> <MethodName>` |
| Bump a module's version | `/bump-module <module> <level>` |
| Release a module | `/release-module <module> <version>` |
| Audit + refresh pinned API docs | `/audit-service-docs` |

## CI gates (per module, in matrix)

Every PR runs, per module:

- `golangci-lint run` — 30+ linters, config at `.golangci.yml`.
- `gosec ./...` — SARIF uploaded to code-scanning.
- `govulncheck ./...`.
- `go test -race -count=1 -coverprofile=...` — coverage ≥ 70%.
- `apidiff` vs the previous `<module>/vX.Y.Z` tag — breaking changes require a `BREAKING CHANGE:` / `Migration:` footer.
- `go vet` + `go build`.

And repo-wide:

- PR title matches Conventional Commits 1.0.
- CodeQL (Go, `security-extended,security-and-quality`).
- OSSF Scorecard.

## Pinned upstream-API docs

`docs/upstream/<service>.md` records the canonical API-doc URL, the version this module targets, and the last-verified date. Check these before suggesting a new method — public docs may be ahead or behind.

## When in doubt

Read [.workingdir/PRINCIPLES.md](.workingdir/PRINCIPLES.md) for the full contract, then the per-category `AGENTS.md` for the area you're touching.

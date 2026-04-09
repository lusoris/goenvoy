# goenvoy

[![CI](https://github.com/lusoris/goenvoy/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/lusoris/goenvoy/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white)](https://go.dev)

A collection of Go client libraries for media automation, metadata, and download service APIs — **63 modules** covering 57+ services.

Designed as a **multi-module monorepo** — import only what you need, from a single library to an entire category.

## Features

- **Pure stdlib** — zero external dependencies, just `net/http` and `encoding/json`
- **Context-aware** — every method takes `context.Context` for cancellation and timeouts
- **Functional options** — `WithHTTPClient`, `WithTimeout`, `WithHeader` and service-specific options
- **Type-safe errors** — each module defines an `APIError` with status code, message, and body
- **OAuth2 support** — device code, authorization code, PKCE, and refresh flows where applicable
- **Thoroughly tested** — httptest-based tests, race detector enabled, no live API calls

## Categories

| Category | Module | Services |
|---|---|---|
| **Arr Stack** | `github.com/lusoris/goenvoy/arr` | Sonarr, Radarr, Lidarr, Readarr, Whisparr (v2+v3), Prowlarr, Bazarr, Seerr, Autobrr, Mylar3, FlareSolverr, Jackett, NZBHydra2 |
| **Metadata** | `github.com/lusoris/goenvoy/metadata` | TMDb, TheTVDB, Fanart.tv, OMDb, TVmaze, Letterboxd, AniList, Kitsu, AniDB, MAL, Trakt, Simkl, MusicBrainz, StashBox, TPDB, OpenSubtitles, Last.fm, Discogs, TheAudioDB, Open Library, Google Books, Spotify, Deezer, ListenBrainz, IGDB, RAWG, Steam |
| **Download Clients** | `github.com/lusoris/goenvoy/downloadclient` | qBittorrent, Transmission, Deluge, rTorrent, SABnzbd, NZBGet |
| **Media Servers** | `github.com/lusoris/goenvoy/mediaserver` | Plex, Jellyfin, Emby, Tautulli, Audiobookshelf, Komga, Navidrome, Kavita, Stash, Tdarr |
| **Anime** | `github.com/lusoris/goenvoy/anime` | Shoko Server |

## Install

Import a specific service library:

```go
go get github.com/lusoris/goenvoy/arr/sonarr
```

Or import shared category types:

```go
go get github.com/lusoris/goenvoy/arr
```

## Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/lusoris/goenvoy/arr/sonarr"
)

func main() {
    client, err := sonarr.New("http://localhost:8989", "your-api-key")
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Get system status
    status, err := client.GetSystemStatus(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s version %s\n", status.AppName, status.Version)

    // Get all series
    series, err := client.GetAllSeries(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Total series: %d\n", len(series))
}
```

Each service module follows the same pattern: `New(baseURL, apiKey) → typed methods with context`.

## Structure

Each category has a base module with shared types, plus sub-modules for individual services:

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
│   ├── movie/        # TMDb, TheTVDB, Fanart.tv, OMDb, TVmaze, Letterboxd, OpenSubtitles
│   ├── anime/        # AniList, Kitsu, AniDB, MAL
│   ├── music/        # MusicBrainz, Last.fm, Discogs, TheAudioDB, Spotify, Deezer, ListenBrainz
│   ├── tracking/     # Trakt, Simkl
│   ├── adult/        # StashBox, TPDB
│   ├── book/         # Google Books, Open Library
│   └── game/         # IGDB, RAWG, Steam
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
├── anime/            # shared anime types
│   └── shoko/        # Shoko Server client
```

## Development

Requires Go 1.26+. See [CONTRIBUTING.md](CONTRIBUTING.md) for full details.

```bash
# Set up workspace (local dev, links all 63 modules)
go work init && find . -name 'go.mod' -not -path './.workingdir/*' -exec dirname {} \; | xargs go work use

# Run all tests
make test-all

# Lint all modules
make lint-all

# Tidy all modules
make tidy-all

# Format all modules
make fmt-all
```

## License

[MIT](LICENSE)
